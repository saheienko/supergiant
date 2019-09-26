package swagger

import (
	"github.com/supergiant/control/pkg/account"
	"github.com/supergiant/control/pkg/model"
)

// accountNameParam is used to identify a cloud account.
// swagger:parameters getAccount updateAccount deleteAccount getAccountRegions getAccountZones getAccountTypes
type accountNameParam struct {
	// in:path
	// required: true
	AccountName string `json:"accountName"`
}

// accountParam contains a cloud account parameters.
// swagger:parameters createAccount updateAccount
type accountParam struct {
	// in:body
	Body model.CloudAccount
}

// regionParam is used to identify a cloud provider region.
// swagger:parameters getAccountZones getAccountTypes
type regionParam struct {
	// in:path
	// required: true
	Region string `json:"region"`
}

// azParam is used to identify a cloud provider availability zone.
// swagger:parameters getAccountTypes
type azParam struct {
	// in:path
	// required: true
	AZ string `json:"az"`
}

// listAccountsResponse contains a list of cloud accounts.
// swagger:response listAccountsResponse
type listAccountsResponse struct {
	// in:body
	AccountList []model.CloudAccount
}

// accountResponse contains representations of a cloud accounts.
// swagger:response accountResponse
type accountResponse struct {
	// in:body
	Account model.CloudAccount
}

// accountRegionsResponse returns a list of supported regions for a given account.
// swagger:response accountRegionsResponse
type accountRegionsResponse struct {
	// in:body
	Regions account.RegionSizes
}

// accountZonesResponse contains a list of supported availability zones for a given account.
// swagger:response accountZonesResponse
type accountZonesResponse struct {
	// in:body
	Zones []string
}

// accountTypesResponse contains a list of supported machine types for a given account.
// swagger:response accountTypesResponse
type accountTypesResponse struct {
	// in:body
	Types []string
}
