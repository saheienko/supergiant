package swagger

// taskIDParam is used to identify a workflow task.
// swagger:parameters getTask restartTask streamLogs
type taskIDParam struct {
	// in:path
	// required: true
	TaskID string `json:"taskID"`
}

// taskResponse contains representations of a workflow task.
// swagger:response taskResponse
type taskResponse struct {
	// in:body
	Task string
}
