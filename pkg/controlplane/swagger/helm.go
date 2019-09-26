package swagger

import (
	"github.com/supergiant/control/pkg/model"
	"k8s.io/helm/pkg/repo"
)

// repoNameParam is used to identify a helm repository.
// swagger:parameters getRepo updateRepo deleteRepo listCharts getChart
type repoNameParam struct {
	// in:path
	// required: true
	RepoName string `json:"repoName"`
}

// chartNameParam is used to identify a helm chart.
// swagger:parameters getChart
type chartNameParam struct {
	// in:path
	// required: true
	chartName string `json:"chartName"`
}

// repoParam contains a helm repository parameters.
// swagger:parameters createRepo updateRepo
type repoParam struct {
	// in:body
	Body repo.Entry
}

// repositoryResponse contains representations of a helm repository.
// swagger:response repositoryResponse
type repositoryResponse struct {
	// in:body
	Task model.ReleaseInfo
}

// listReposResponse contains a list of helm repos.
// swagger:response listReposResponse
type listReposResponse struct {
	// in:body
	Task []model.ReleaseInfo
}

// listChartsResponse contains a list of helm charts.
// swagger:response listChartsResponse
type listChartsResponse struct {
	// in:body
	Task model.ChartInfo
}

// chartResponse contains representations of a helm chart.
// swagger:response chartResponse
type chartResponse struct {
	// in:body
	Chart model.ChartData
}
