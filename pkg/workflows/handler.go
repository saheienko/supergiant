package workflows

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/control/pkg/model"
	"github.com/supergiant/control/pkg/runner"
	"github.com/supergiant/control/pkg/runner/ssh"
	"github.com/supergiant/control/pkg/storage"
	"github.com/supergiant/control/pkg/util"
	"github.com/supergiant/control/pkg/workflows/steps"
)

type cloudAccountGetter interface {
	Get(context.Context, string) (*model.CloudAccount, error)
}

type TaskHandler struct {
	runnerFactory func(config ssh.Config) (runner.Runner, error)
	getTail       func(string) (*tail.Tail, error)

	cloudAccGetter cloudAccountGetter
	repository     storage.Interface
	getWriter      func(string) (io.WriteCloser, error)
}

type RunTaskRequest struct {
	WorkflowName string       `json:"workflowName"`
	Cfg          steps.Config `json:"config"`
}

type BuildTaskRequest struct {
	StepNames []string     `json:"stepNames"`
	Cfg       steps.Config `json:"config"`
	SshConfig ssh.Config   `json:"sshConfig"`
}

type TaskResponse struct {
	ID string `json:"id"`
}

func NewTaskHandler(repository storage.Interface, runnerFactory func(config ssh.Config) (runner.Runner, error), getter cloudAccountGetter, logDir string) *TaskHandler {
	return &TaskHandler{
		runnerFactory:  runnerFactory,
		repository:     repository,
		cloudAccGetter: getter,
		getWriter:      util.GetWriterFunc(logDir),
		getTail: func(id string) (*tail.Tail, error) {
			t, err := tail.TailFile(path.Join(logDir, util.MakeFileName(id)),
				tail.Config{
					Follow:    true,
					MustExist: true,
					Location: &tail.SeekInfo{
						Offset: 0,
						Whence: io.SeekStart,
					},
					Logger:      tail.DiscardingLogger,
					MaxLineSize: 160,
				})

			if err != nil {
				return nil, err
			}

			return t, nil
		},
	}
}

func (h *TaskHandler) Register(m *mux.Router) {
	// swagger:route GET /v1/api/tasks/{taskID} tasks getTask
	//
	// Get a task.
	//
	// Responses:
	// default: errorResponse
	// 200: taskResponse
	//
	m.HandleFunc("/tasks/{id}", h.GetTask).Methods(http.MethodGet)

	// swagger:route POST /v1/api/tasks/{taskID}/restart tasks restartTask
	//
	// Restart a task.
	//
	// Responses:
	// default: errorResponse
	// 202: emptyResponse
	//
	m.HandleFunc("/tasks/{id}/restart", h.RestartTask).Methods(http.MethodPost)

	// swagger:route POST /v1/api/tasks/{taskID}/logs tasks streamLogs
	//
	// Set up a websocket connection for logs streaming.
	//
	// Responses:
	// default: errorResponse
	// 200: emptyResponse
	//
	m.HandleFunc("/tasks/{id}/logs", h.StreamLogs).Methods(http.MethodGet)
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]

	if !ok {
		http.Error(w, "need id of task", http.StatusBadRequest)
		return
	}

	data, err := h.repository.Get(r.Context(), Prefix, id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

func (h *TaskHandler) RestartTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]

	if !ok {
		http.Error(w, "need id of task", http.StatusBadRequest)
		return
	}

	logrus.Debugf("get task %s", id)
	data, err := h.repository.Get(r.Context(), Prefix, id)

	if err != nil {
		logrus.Debugf("task %s not found", id)
		http.NotFound(w, r)
		return
	}

	task, err := DeserializeTask(data, h.repository)

	if err != nil {
		logrus.Debugf("error deserializing task %s %v", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fileName := util.MakeFileName(id)
	writer, err := h.getWriter(fileName)

	if err != nil {
		http.Error(w, fmt.Sprintf("get writer %v", err), http.StatusInternalServerError)
		logrus.Errorf("Get writer %v", err)
		return
	}

	task.Run(context.Background(), *task.Config, writer)
	w.WriteHeader(http.StatusAccepted)
}

func (h *TaskHandler) StreamLogs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]

	if !ok {
		http.Error(w, "need id of task", http.StatusBadRequest)
		return
	}

	var upgrader = websocket.Upgrader{
		HandshakeTimeout: time.Second * 10,
		WriteBufferSize:  1024,
		ReadBufferSize:   0,
		// TODO(stgleb): Do something more safe in future
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	t, err := h.getTail(id)

	if os.IsNotExist(err) {
		http.NotFound(w, r)
		logrus.Errorf("Not found %s", util.MakeFileName(id))
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logrus.Errorf("Open file %s for tail %v", util.MakeFileName(id), err)
		return
	}

	c, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logrus.Errorf("Upgrade connection %v", err)
		return
	}

	go func() {
		pingTicker := time.NewTicker(time.Second * 60)

		for {
			select {
			case line := <-t.Lines:
				c.SetWriteDeadline(time.Now().Add(time.Second * 10))
				err = c.WriteMessage(websocket.TextMessage, []byte(line.Text))

				// Do not log this error, since client can simply disconnect
				if err != nil {
					return
				}
			case <-pingTicker.C:
				c.SetWriteDeadline(time.Now().Add(time.Second * 10))
				if err := c.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
					return
				}
			}
		}
	}()
}
