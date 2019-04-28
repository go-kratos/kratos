package model

// PaasJobResponse create job response in paas.
type PaasJobResponse struct {
	ExcludeDataResponse
	Data interface{} `json:"data"`
}

//Job job json
type Job struct {
	Name        string `json:"name"`
	CPU         int    `json:"cpu"`
	Memory      int    `json:"memory"`
	Parallelism int    `json:"parallelism"`
	FileName    string `json:"file_name"`
	ResJtl      string `json:"res_jtl"`
	ResLog      string `json:"res_log"`
	JmeterLog   string `json:"jmeter_log"`
	EnvInfo     string `json:"env_info"`
	JarPath     string `json:"jar_path"`
	Command     string `json:"command"`
}

// CleanableDocker docker clearable container list
type CleanableDocker struct {
	Name string `json:"name"`
}

//PaasQueryAndDelJob query and del machines request.
type PaasQueryAndDelJob struct {
	BusinessUnit string `json:"business_unit"`
	Project      string `json:"project"`
	App          string `json:"app"`
	Env          string `json:"env"`
	Name         string `json:"name"`
	ClusterID    int    `json:"cluster_id"`
	TreeID       int    `json:"tree_id"`
}

// PaasJobDetail machine detail.
type PaasJobDetail struct {
	BusinessUnit   string `json:"business_unit"`
	Project        string `json:"project"`
	App            string `json:"app"`
	Env            string `json:"env"`
	Name           string `json:"name"`
	Image          string `json:"image"`
	ImageVersion   string `json:"image_version"`
	Volumes        string `json:"volumes"`
	CPURequest     int    `json:"cpu_request"`
	CPULimit       int    `json:"cpu_limit"`
	MemoryRequest  int    `json:"memory_request"`
	Command        string `json:"command"`
	ResourcePoolID string `json:"resource_pool_id"`
	Parallelism    int    `json:"parallelism"`
	Completions    int    `json:"completions"`
	RetriesLimit   int    `json:"retries_limit"`
	NetworkID      int    `json:"network_id"`
	ClusterID      int    `json:"cluster_id"`
	TreeID         int    `json:"tree_id"`
	HostInfo       string `json:"host_info"`
	EnvInfo        string `json:"env_info"`
}

// PaasJobQueryStatus machine detail.
type PaasJobQueryStatus struct {
	ExcludeDataResponse
	Data PaasJobQueryData `json:"data"`
}

// PaasJobQueryData machine detail.
type PaasJobQueryData struct {
	StartTime      string                 `json:"start_time"`
	CompletionTime string                 `json:"completion_time"`
	ActiveNum      int                    `json:"active_num"`
	SucceededNum   int                    `json:"succeeded_num"`
	FailedNum      int                    `json:"failed_num"`
	Conditions     PaasJobQueryConditions `json:"conditions"`
	Pods           []PodInfo              `json:"pods"`
}

//PaasQueryJobCPUPostDetail query job cpu detail
type PaasQueryJobCPUPostDetail struct {
	Action     string `json:"Action"`
	PublicKey  string `json:"PublicKey"`
	Signature  int    `json:"Signature"`
	DataSource string `json:"DataSource"`
	Query      string `json:"Query"`
}

//PaasQueryJobCPUResult paas query cpu result
type PaasQueryJobCPUResult struct {
	ReqID   string      `json:"ReqId"`
	Action  string      `json:"Action"`
	RetCode int         `json:"RetCode"`
	Data    []CPUResult `json:"Data"`
}

//CPUResult cpu result
type CPUResult struct {
	JobMetric JobMetric     `json:"metric"`
	Value     []interface{} `json:"value"`
}

//JobMetric job metric
type JobMetric struct {
	ContainerEnvAppID     string `json:"container_env_app_id"`
	ContainerEnvDeployEnv string `json:"container_env_deploy_env"`
	ContainerEnvPodCon    string `json:"container_env_pod_container"`
	ContainerEnvPodName   string `json:"container_env_pod_name"`
	Job                   string `json:"job"`
	Pro                   string `json:"pro"`
}

// PodInfo pod info
type PodInfo struct {
	AppID             string      `json:"app_id"`
	AppType           string      `json:"app_type"`
	ContainerID       string      `json:"container_id"`
	ContainerStatuses interface{} `json:"container_statuses"`
	CreateTime        string      `json:"create_time"`
	DeployEnv         string      `json:"deploy_env"`
	DiscoveryStatus   interface{} `json:"discovery_status"`
	Health            string      `json:"health"`
	HostIP            string      `json:"host_ip"`
	Image             string      `json:"image"`
	IP                string      `json:"ip"`
	Lables            interface{} `json:"lables"`
	Name              string      `json:"name"`
	Namespace         string      `json:"namespace"`
	Port              interface{} `json:"port"`
	StartTime         string      `json:"start_time"`
	Status            string      `json:"status"`
}

// PaasJobQueryConditions machine detail.
type PaasJobQueryConditions struct {
	Complete string `json:"complete"`
}
