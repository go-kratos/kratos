package model

import (
	"net/http"

	"go-common/library/ecode"
	"go-common/library/log"
)

//ClusterItems cluster items
type ClusterItems struct {
	Items map[string]interface{} `json:"items"`
}

//RmTokenPost get token
type RmTokenPost = struct {
	APIToken   string `json:"api_token"`
	PlatformID string `json:"platform_id"`
}

//ClusterResponse json body
type ClusterResponse struct {
	ExcludeDataResponse
	Count int                   `json:"count"`
	Data  *ClusterResponseItems `json:"data"`
}

//ClusterResponseItems json body
type ClusterResponseItems struct {
	Items []*ClusterResponseItemsSon `json:"items"`
}

//ClusterResponseItemsSon json body
type ClusterResponseItemsSon struct {
	Configuration      ClusterConfiguration      `json:"configuration"`
	Resources          ClusterPoolResourcesUsage `json:"resources"`
	PoolResourcesUsage PoolResourcesUsageSon     `json:"pool_resources_usage"`
}

//PoolResourcesUsageSon json body
type PoolResourcesUsageSon struct {
	Ep     ClusterPoolResourcesUsage `json:"ep"`
	Melloi ClusterPoolResourcesUsage `json:"melloi"`
	Public ClusterPoolResourcesUsage `json:"public"`
}

//ClusterConf visit PaaS
type ClusterConf struct {
	TestHost        string
	UatHost         string
	QueryJobCPUHost string
}

// ClusterRmTokenResponse query token response.
type ClusterRmTokenResponse struct {
	ExcludeDataResponse
	Data AuthData `json:"data"`
}

// ExcludeDataResponse no data field response.
type ExcludeDataResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// CheckStatus check status.
func (e *ExcludeDataResponse) CheckStatus() (err error) {
	if e.Status >= http.StatusMultipleChoices || e.Status < http.StatusOK {
		log.Error("The status(%d) of paas may represent a request error(%s)", e.Status, e.Message)
		err = ecode.MelloiPaasRequestErr
	}
	return
}

//AuthData AuthData stuct
type AuthData struct {
	Token      string `json:"token"`
	PlatformID string `json:"platform_id"`
	UserName   string `json:"user_name"`
	Secret     string `json:"secret"`
	Expired    int    `json:"expired"`
	Admin      bool   `json:"admin"`
}

//ClusterPoolResourcesUsage ClusterPoolResourcesUsage stuct
type ClusterPoolResourcesUsage struct {
	CPURequest        int `json:"cpu_request"`
	CPULimit          int `json:"cpu_limit"`
	CPUCapacity       int `json:"cpu_capacity"`
	CPUAllocatable    int `json:"cpu_allocatable"`
	MemoryRequest     int `json:"memory_request"`
	MemoryLimit       int `json:"memory_limit"`
	MemoryCapacity    int `json:"memory_capacity"`
	MemoryAllocatable int `json:"memory_allocatable"`
	PodRunning        int `json:"pod_running"`
	PodTotal          int `json:"pod_total"`
	PodCapacity       int `json:"pod_capacity"`
	PodAllocatable    int `json:"pod_allocatable"`
	NodesNum          int `json:"nodes_num"`
}

//ClusterConfiguration ClusterConfiguration stuct
type ClusterConfiguration struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Desc            string `json:"desc"`
	Region          string `json:"region"`
	Zone            string `json:"zone"`
	Idc             string `json:"idc"`
	Envs            string `json:"envs"`
	CloudProvider   string `json:"cloud_provider"`
	Orchestrator    string `json:"orchestrator"`
	APITarget       string `json:"api_target"`
	Register        string `json:"register"`
	DefaultRegistry string `json:"default_registry"`
	Ctime           string `json:"ctime"`
	Mtime           string `json:"mtime"`
}
