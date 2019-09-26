package swagger

import "github.com/supergiant/control/pkg/message"

// errorResponse represents an error response.
// swagger:response errorResponse
type errorResponse struct {
	// in:body
	Body message.Message
}

// emptyResponse represents an empty http response.
// swagger:response emptyResponse
type emptyResponse struct {
	// in:body
	Body string
}
