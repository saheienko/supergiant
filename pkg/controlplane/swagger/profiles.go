package swagger

import (
	"github.com/supergiant/control/pkg/profile"
)

// kubeprofileIDParam is used to identify a kube profile model.
// swagger:parameters getProfile
type kubeprofileIDParam struct {
	// in:path
	// required: true
	KubeprofileID string `json:"kubeprofileID"`
}

// profileParam contains a kube profile parameters.
// swagger:parameters createProfile
type profileParam struct {
	// in:body
	Body profile.Profile
}

// profileResponse contains representations of a kube profile.
// swagger:response profileResponse
type profileResponse struct {
	// in:body
	AccountList profile.Profile
}

// listProfilesResponse returns all kube profiles.
// swagger:response listProfilesResponse
type listProfilesResponse struct {
	// in:body
	AccountList []profile.Profile
}
