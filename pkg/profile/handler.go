package profile

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pborman/uuid"
	"github.com/sirupsen/logrus"
	"gopkg.in/asaskevich/govalidator.v8"

	"github.com/supergiant/control/pkg/sgerrors"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Register(r *mux.Router) {
	// swagger:route GET /v1/api/kubeprofiles/{kubeprofileID} profiles getProfile
	//
	// Get a profile.
	//
	// Responses:
	// default: errorResponse
	// 200: profileResponse
	//
	r.HandleFunc("/kubeprofiles/{id}", h.GetProfile).Methods(http.MethodGet)

	// swagger:route POST /v1/api/kubeprofiles profiles createProfile
	//
	// Create a profile.
	//
	// Responses:
	// default: errorResponse
	// 201: emptyResponse
	//
	r.HandleFunc("/kubeprofiles", h.CreateProfile).Methods(http.MethodPost)

	// swagger:route GET /v1/api/kubeprofiles profiles listProfiles
	//
	// List profiles.
	//
	// Responses:
	// default: errorResponse
	// 200: listProfilesResponse
	//
	r.HandleFunc("/kubeprofiles", h.GetProfiles).Methods(http.MethodGet)
}

func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	profileId := vars["id"]

	kubeProfile, err := h.service.Get(r.Context(), profileId)
	if err != nil {
		if sgerrors.IsNotFound(err) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(kubeProfile); err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) CreateProfile(w http.ResponseWriter, r *http.Request) {
	profile := &Profile{}

	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	profile.ID = uuid.NewUUID().String()

	ok, err := govalidator.ValidateStruct(profile)
	if !ok {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.Create(r.Context(), profile); err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetProfiles(w http.ResponseWriter, r *http.Request) {
	profiles, err := h.service.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(profiles); err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
