package swagger

import (
	"github.com/supergiant/control/pkg/user"
)

// authParams contains signin parameters.
// swagger:parameters auth
type authParams struct {
	// in:body
	Req user.AuthRequest
}

// userParams contains user details.
// swagger:parameters registerRoot createUser
type userParams struct {
	// in:body
	Req user.User
}

// authResponse contains an auth header.
// swagger:response authResponse
type authResponse struct {
	// Authorization header
	Authorization string `json:"Authorization"`
	// AccessControlExposeHeaders header
	AccessControlExposeHeaders string `json:"Access-Control-Expose-Headers"`
}

// coldstartResponse tells if there is any user exists.
// swagger:response coldstartResponse
type coldstartResponse struct {
	// in:body
	ColdstartResponse user.ColdstartResponse
}
