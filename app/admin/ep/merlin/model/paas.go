package model

import (
	"encoding/json"
	"net/http"

	"go-common/library/ecode"
	"go-common/library/log"
)

// PaasConf conf of paas.
type PaasConf struct {
	Host              string
	Token             string
	MachineTimeout    string
	MachineLimitRatio float32
}

// PaasMachine machine in paas.
type PaasMachine struct {
	PaasMachineSystem
	Name          string  `json:"name"`
	Image         string  `json:"image"`
	CPURequest    float32 `json:"cpu_request"`
	MemoryRequest float32 `json:"memory_request"`

	CPULimit    float32 `json:"cpu_limit"`
	MemoryLimit float32 `json:"memory_limit"`

	DiskRequest    int    `json:"disk_request"`
	VolumnMount    string `json:"volumn_mount"`
	Snapshot       bool   `json:"snapshot"`
	ForcePullImage bool   `json:"force_pull_image"`
}

// ToMachine convert PaasMachine to Machine.
func (pm *PaasMachine) ToMachine(u string, gmr *GenMachinesRequest) (m *Machine) {
	m = &Machine{
		Username: u,
		Name:     pm.Name,
		PodName:  pm.Name + MachinePodNameSuffix,
		Status:   CreatingMachineInMerlin,
	}
	gmr.Mutator(m)
	return
}

// PaasGenMachineRequest create machine request in paas.
type PaasGenMachineRequest struct {
	Env
	BusinessUnit string        `json:"business_unit"`
	Project      string        `json:"project"`
	App          string        `json:"app"`
	TreeIDs      []int64       `json:"tree_id"`
	Machines     []PaasMachine `json:"machines"`
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
		err = ecode.MerlinPaasRequestErr
	}
	return
}

// Instance instance.
type Instance struct {
	InstanceName string `json:"instance_name"`
}

// CreateInstance create instance.
type CreateInstance struct {
	Instance
	InstanceCreateStatus int `json:"instance_create_status"`
}

// ReleaseInstance release instance.
type ReleaseInstance struct {
	Instance
	InstanceReleaseStatus int `json:"instance_release_status"`
}

// PaasGenMachineResponse create machine response in paas.
type PaasGenMachineResponse struct {
	ExcludeDataResponse
	Data []*CreateInstance `json:"data"`
}

// PaasDelMachineResponse delete machine response in paas.
type PaasDelMachineResponse struct {
	ExcludeDataResponse
	Data ReleaseInstance `json:"data"`
}

// PaasSnapshotMachineResponse Paas Snapshot Machine Response.
type PaasSnapshotMachineResponse struct {
	ExcludeDataResponse
	Data string `json:"data"`
}

// DetailCondition detail condition.
type DetailCondition struct {
	Initialized  string `json:"Initialized"`
	PodScheduled string `json:"PodScheduled"`
	Ready        string `json:"Ready"`
}

// MachineStatus machine status.
type MachineStatus struct {
	Condition       string          `json:"condition"`
	Message         string          `json:"message"`
	DetailCondition DetailCondition `json:"detail_conditions"`
	InstanceIP      string          `json:"instance_ip"`
	RestartCount    int             `json:"restart_count"`
	Log             string          `json:"log"`
	MachineEvent    MachineEvent    `json:"events"`
}

// MachineEvent MachineEvent.
type MachineEvent struct {
	MachinesItems []*MachinesItem `json:"Items"`
}

// MachinesItem MachinesItem.
type MachinesItem struct {
	Kind      string `json:"Kind"`
	Name      string `json:"Name"`
	Namespace string `json:"Namespace"`
	Reason    string `json:"Reason"`
	Message   string `json:"Message"`
	Count     int    `json:"Count"`
	LastTime  string `json:"LastTime"`
}

// ToMachineStatusResponse convert  MachineStatus to MachineStatusResponse.
func (ms *MachineStatus) ToMachineStatusResponse() (msr *MachineStatusResponse) {
	var events []string
	for _, item := range ms.MachineEvent.MachinesItems {
		mJSON, _ := json.Marshal(item)
		events = append(events, string(mJSON))
	}

	return &MachineStatusResponse{
		Initialized:  ms.DetailCondition.Initialized,
		PodScheduled: ms.DetailCondition.PodScheduled,
		Ready:        ms.DetailCondition.Ready,
		Log:          ms.Log,
		MachineEvent: events,
	}
}

// PaasQueryMachineStatusResponse query machine status response.
type PaasQueryMachineStatusResponse struct {
	ExcludeDataResponse
	Data MachineStatus `json:"data"`
}

// PaasQueryClustersResponse query cluster response.
type PaasQueryClustersResponse struct {
	ExcludeDataResponse
	Data Clusters `json:"data"`
}

// PaasQueryClusterResponse query cluster response.
type PaasQueryClusterResponse struct {
	ExcludeDataResponse
	Data *Cluster `json:"data"`
}

// PaasMachineDetail machine detail.
type PaasMachineDetail struct {
	Condition string `json:"condition"`
	Name      string `json:"name"`
	Image     string `json:"image"`

	CPURequest    float32 `json:"cpu_request"`
	MemoryRequest float32 `json:"memory_request"`

	CPULimit    float32 `json:"cpu_limit"`
	MemoryLimit float32 `json:"memory_limit"`

	DiskRequest int    `json:"disk_request"`
	VolumnMount string `json:"volumn_mount"`
	ClusterName string `json:"cluster_name"`
	Env         string `json:"env"`
	IP          string `json:"ip"`
	PaasMachineSystem
}

// ConvertUnits convert units of cpu and memory.
func (p PaasMachineDetail) ConvertUnits() PaasMachineDetail {
	p.CPURequest = p.CPURequest / CPURatio
	p.MemoryRequest = p.MemoryRequest / MemoryRatio
	p.CPULimit = p.CPULimit / CPURatio
	p.MemoryLimit = p.MemoryLimit / MemoryRatio
	return p
}

// PaasQueryMachineResponse query machine response.
type PaasQueryMachineResponse struct {
	ExcludeDataResponse
	Data PaasMachineDetail `json:"data"`
}

// Clusters clusters.
type Clusters struct {
	Items []*Cluster `json:"items"`
}

// Network network.
type Network struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	Subnet   string  `json:"subnet"`
	Capacity float64 `json:"capacity"`
}

// Resource resource.
type Resource struct {
	CPUUsage    float64 `json:"cpu_usage"`    //集群总体CPU使用率
	MemUsage    float64 `json:"mem_usage"`    //集群总体内存使用率
	PodTotal    int     `json:"pod_total"`    //集群实例总数
	PodCapacity int     `json:"pod_capacity"` //集群实例容量
	NodesNum    int     `json:"nodes_num"`    //集群节点数量
}

// Cluster cluster.
type Cluster struct {
	ID                int64     `json:"id"`
	Name              string    `json:"name"`
	IsSupportSnapShot bool      `json:"is_support_snapshot"`
	Desc              string    `json:"desc"`
	IDc               string    `json:"idc"`
	Networks          []Network `json:"networks"`
	Resources         Resource  `json:"resources"`
}

// Verify verify cluster
func (c *Cluster) Verify() error {
	if c.Name == "" || len(c.Networks) < 1 || c.Networks[0].Name == "" {
		return ecode.MerlinInvalidClusterErr
	}
	return nil
}

// ToMachine convert CreateInstance to Machine.
func (i *CreateInstance) ToMachine(u string, gmr *GenMachinesRequest) (m *Machine) {
	var status int
	switch i.InstanceCreateStatus {
	case CreateFailedMachineInPaas:
		status = ImmediatelyFailedMachineInMerlin
	case CreatingMachineInPass:
		status = CreatingMachineInMerlin
	}
	m = &Machine{
		Username: u,
		Name:     i.InstanceName,
		PodName:  i.InstanceName + MachinePodNameSuffix,
		Status:   status,
	}
	gmr.Mutator(m)
	return
}

// PaasQueryAndDelMachineRequest query and del machines request.
type PaasQueryAndDelMachineRequest struct {
	BusinessUnit string `json:"business_unit"`
	Project      string `json:"project"`
	App          string `json:"app"`
	ClusterID    int64  `json:"cluster_id"`
	Name         string `json:"name"`
}

// PaasAuthRequest auth request.
type PaasAuthRequest struct {
	APIToken   string `json:"api_token"`
	PlatformID string `json:"platform_id"`
}

// PaasAuthResponse auth response.
type PaasAuthResponse struct {
	ExcludeDataResponse
	Data PaasAuthInfo `json:"data"`
}

// PaasAuthInfo auth information.
type PaasAuthInfo struct {
	Token      string `json:"token"`
	PlatformID string `json:"platform_id"`
	Username   string `json:"user_name"`
	Secret     string `json:"secret"`
	Expired    int64  `json:"expired"`
	Admin      bool   `json:"admin"`
}

// PaasUpdateMachineNodeRequest update machine request.
type PaasUpdateMachineNodeRequest struct {
	PaasQueryAndDelMachineRequest
	TreeID []int64 `json:"tree_id"`
}

// NewPaasUpdateMachineNodeRequest new PaasUpdateMachineNodeRequest
func NewPaasUpdateMachineNodeRequest(p *PaasQueryAndDelMachineRequest, ns []*MachineNode) *PaasUpdateMachineNodeRequest {
	var treeIDs []int64
	for _, n := range ns {
		treeIDs = append(treeIDs, n.TreeID)
	}
	return &PaasUpdateMachineNodeRequest{
		TreeID: treeIDs,
		PaasQueryAndDelMachineRequest: *p,
	}
}

// PaasUpdateMachineNodeResponse update machine response.
type PaasUpdateMachineNodeResponse struct {
	ExcludeDataResponse
	Data string `json:"data"`
}
