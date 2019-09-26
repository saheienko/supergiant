package swagger

import (
	"github.com/supergiant/control/pkg/kube"
	"github.com/supergiant/control/pkg/model"
	"github.com/supergiant/control/pkg/profile"
	"github.com/supergiant/control/pkg/provisioner"
	"github.com/supergiant/control/pkg/workflows/steps"
	v1 "k8s.io/api/core/v1"
	"k8s.io/helm/pkg/proto/hapi/release"
)

type TaskMap map[string]string

// kubeIDParam identifies a kube model.
// swagger:parameters getKube deleteKube getKubeconfig installRelease listReleases getRelease deleteRelease listTasks addNodes listNodes deleteNode getNodesMetrics getClusterMetrics listServices restartProvisioning applyResources
type kubeIDParam struct {
	// in:path
	KubeID string `json:"kubeID"`
}

// unameParam identifies a kubernetes user.
// swagger:parameters getKubeconfig
type unameParam struct {
	// in:path
	Uname string `json:"uname"`
}

// releaseNameParam identifies a helm release.
// swagger:parameters getRelease deleteRelease
type releaseNameParam struct {
	// in:path
	ReleaseName string `json:"releaseName"`
}

// nodeNameParam identifies a kubernetes node.
// swagger:parameters deleteNode
type nodeNameParam struct {
	// in:path
	NodeName string `json:"nodeName"`
}

// listNodesQueryParam identifies a node role.
// swagger:parameters listNodes
type listNodesQueryParam struct {
	// in:query
	Role string `json:"role"`
}

// kubeBodyParam contains a kubernetes cluster parameters.
// swagger:parameters createKube
type kubeBodyParam struct {
	// in:body
	Body model.Kube
}

// provisionBodyParam contains a cluster provisioning parameters.
// swagger:parameters provision
type provisionBodyParam struct {
	// in:body
	Body provisioner.ProvisionRequest
}

// importBodyParam contains a cluster import parameters.
// swagger:parameters importKube
type importBodyParam struct {
	// in:body
	Body kube.ImportRequest
}

// installReleaseBodyParam contains a helm release parameters.
// swagger:parameters installRelease
type installReleaseBodyParam struct {
	// in:body
	Body steps.InstallAppConfig
}

// addNodesBodyParam contains a list of node profiles.
// swagger:parameters addNodes
type addNodesBodyParam struct {
	// in:body
	Body []profile.NodeProfile
}

// provisionResponse contains a cluster provisioning metadata.
// swagger:response provisionResponse
type provisionResponse struct {
	// in:body
	Provision provisioner.ProvisionResponse
}

// listKubesResponse contains a list of kube models.
// swagger:response listKubesResponse
type listKubesResponse struct {
	// in:body
	KubeList []model.Kube
}

// importKubeResponse contains an id of a cluster.
// swagger:response importKubeResponse
type importKubeResponse struct {
	// in:body
	Body kube.ImportResponse
}

// kubeResponse contains an id of a cluster.
// swagger:response kubeResponse
type kubeResponse struct {
	// in:body
	Kube model.Kube
}

// kubeconfigResponse contains an kubeconfig for a cluster.
// swagger:response kubeconfigResponse
type kubeconfigResponse struct {
	// in:body
	Kube string
}

// installReleaseResponse contains a task id.
// swagger:response installReleaseResponse
type installReleaseResponse struct {
	// in:body
	Kube kube.InstallReleaseResp
}

// listReleasesResponse contains a list of helm releases.
// swagger:response listReleasesResponse
type listReleasesResponse struct {
	// in:body
	Kube []model.ReleaseInfo
}

// releaseResponse contains a list of helm releases.
// swagger:response releaseResponse
type releaseResponse struct {
	// in:body
	Kube release.Release
}

// deleteReleaseResponse contains a helm release details.
// swagger:response deleteReleaseResponse
type deleteReleaseResponse struct {
	// in:body
	Kube model.ReleaseInfo
}

// listTasksResponse contains a list of provisioning tasks.
// swagger:response listTasksResponse
type listTasksResponse struct {
	// in:body
	Kube []kube.TaskDTO
}

// addNodesResponse contains a list of provisioning tasks.
// swagger:response addNodesResponse
type addNodesResponse struct {
	// in:body
	Tasks []string
}

// listNodesResponse contains a list of provisioning tasks.
// swagger:response listNodesResponse
type listNodesResponse struct {
	// in:body
	Tasks []v1.Node
}

// nodeMetricsResponse contains kubernetes metrics for nodes.
// swagger:response nodeMetricsResponse
type nodeMetricsResponse struct {
	// in:body
	Metrics map[string]map[string]interface{}
}

// clusterMetricsResponse contains kubernetes cluster metrics.
// swagger:response clusterMetricsResponse
type clusterMetricsResponse struct {
	// in:body
	Metrics map[string]interface{}
}

// listServicesResponse contains a list of service for proxying.
// swagger:response listServicesResponse
type listServicesResponse struct {
	// in:body
	Services []kube.ServiceInfo
}
