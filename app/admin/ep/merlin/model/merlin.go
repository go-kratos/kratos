package model

import (
	"context"
	"strconv"
	"time"

	"go-common/library/ecode"
)

// Env Env.
type Env struct {
	ClusterID int64 `json:"cluster_id"`
	NetworkID int64 `json:"network_id"`
}

// Node Node.
type Node struct {
	BusinessUnit string `json:"business_unit"`
	Project      string `json:"project"`
	App          string `json:"app"`
	TreeID       int64  `json:"tree_id"`
}

// GenMachinesRequest Gen Machines Request.
type GenMachinesRequest struct {
	Env
	PaasMachine
	PaasMachineSystem
	Nodes   []*Node `json:"nodes"`
	Comment string  `json:"comment"`
	Amount  int     `json:"amount"`
}

// Verify verify GenMachinesRequest.
func (g *GenMachinesRequest) Verify() error {
	if g.Amount < 1 {
		return ecode.MerlinInvalidMachineAmountErr
	}
	l := len(g.Nodes)
	if l < 1 || l > 10 {
		return ecode.MerlinInvalidNodeAmountErr
	}
	return nil
}

// Mutator  Mutator.
func (g *GenMachinesRequest) Mutator(m *Machine) {
	m.BusinessUnit = g.Nodes[0].BusinessUnit
	m.Project = g.Nodes[0].Project
	m.App = g.Nodes[0].App
	m.ClusterID = g.ClusterID
	m.NetworkID = g.NetworkID
	m.Comment = g.Comment
}

// ToMachineNode to machine node.
func (g *GenMachinesRequest) ToMachineNode(mID int64) (treeNodes []*MachineNode) {
	for _, node := range g.Nodes {
		treeNodes = append(treeNodes, &MachineNode{
			MachineID:    mID,
			BusinessUnit: node.BusinessUnit,
			Project:      node.Project,
			App:          node.App,
			TreeID:       node.TreeID,
		})
	}
	return
}

// PaasMachineSystem Paas Machine System.
type PaasMachineSystem struct {
	Command   string         `json:"command"`
	Envs      []*EnvVariable `json:"envs"`
	HostAlias []*Host        `json:"host_alias"`
}

// Host Host.
type Host struct {
	IP        string   `json:"ip"`
	Hostnames []string `json:"hostnames"`
}

// EnvVariable EnvVariable.
type EnvVariable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// GenMachine GenMachine.
type GenMachine struct {
	Machine
	IP          string `json:"ip"`
	ClusterName string `json:"cluster_name"`
	NetworkName string `json:"network_name"`
}

// MachineDetail Machine Detail.
type MachineDetail struct {
	Machine
	PaasMachineDetail
	Nodes       []*MachineNode `json:"nodes"`
	NetworkName string         `json:"network_name"`
	Name        string         `json:"name"`
	IsSnapShot  bool           `json:"is_support_snapshot"`
}

// PaginateMachine Paginate Machine.
type PaginateMachine struct {
	Total    int64        `json:"total"`
	PageNum  int          `json:"page_num"`
	PageSize int          `json:"page_size"`
	Machines []GenMachine `json:"machines"`
}

// PaginateMachineLog Paginate Machine Log.
type PaginateMachineLog struct {
	Total       int64               `json:"total"`
	PageNum     int                 `json:"page_num"`
	PageSize    int                 `json:"page_size"`
	MachineLogs []*AboundMachineLog `json:"machine_logs"`
}

// PaginateMobileMachineLog Paginate Mobile Machine Log.
type PaginateMobileMachineLog struct {
	Total       int64                     `json:"total"`
	PageNum     int                       `json:"page_num"`
	PageSize    int                       `json:"page_size"`
	MachineLogs []*AboundMobileMachineLog `json:"machine_logs"`
}

// PaginateMobileMachineLendOutLog Paginate Mobile Machine Lend Out Log.
type PaginateMobileMachineLendOutLog struct {
	Total                 int64                   `json:"total"`
	PageNum               int                     `json:"page_num"`
	PageSize              int                     `json:"page_size"`
	MachineLendOutRecords []*MachineLendOutRecord `json:"machine_lendout_records"`
}

// MachineLendOutRecord Machine Lend Out Record.
type MachineLendOutRecord struct {
	MachineID    int64     `json:"machine_id"`
	Lender       string    `json:"lender"`
	Status       int       `json:"status"`
	LendTime     time.Time `json:"lend_time"`
	EnableReturn bool      `json:"enable_return"`
}

// PaginateMobileMachineErrorLog Paginate Mobile Machine Error Log.
type PaginateMobileMachineErrorLog struct {
	Total       int64                    `json:"total"`
	PageNum     int                      `json:"page_num"`
	PageSize    int                      `json:"page_size"`
	MachineLogs []*MobileMachineErrorLog `json:"machine_logs"`
}

// PaginateApplicationRecord Paginate Application Record.
type PaginateApplicationRecord struct {
	Total              int64                `json:"total"`
	PageNum            int                  `json:"page_num"`
	PageSize           int                  `json:"page_size"`
	ApplicationRecords []*ApplicationRecord `json:"application_records"`
}

// ToPaasGenMachineRequest to Paas GenMachine Request.
func (g GenMachinesRequest) ToPaasGenMachineRequest(machineLimitRatio float32) *PaasGenMachineRequest {
	var (
		pms     = make([]PaasMachine, g.Amount)
		treeIDs []int64
		fn      = g.Nodes[0]
	)
	g.CPULimit = g.CPURequest * CPURatio
	g.MemoryLimit = g.MemoryRequest * MemoryRatio

	g.CPURequest = g.CPURequest * CPURatio / machineLimitRatio
	g.MemoryRequest = g.MemoryRequest * MemoryRatio / machineLimitRatio

	for i := 0; i < g.Amount; i++ {
		pms[i] = g.PaasMachine
		pms[i].Name = g.Name + "-" + strconv.Itoa(i+1)
		pms[i].PaasMachineSystem = g.PaasMachineSystem
	}
	for _, node := range g.Nodes {
		treeIDs = append(treeIDs, node.TreeID)
	}
	return &PaasGenMachineRequest{
		BusinessUnit: fn.BusinessUnit,
		Project:      fn.Project,
		App:          fn.App,
		TreeIDs:      treeIDs,
		Env:          g.Env,
		Machines:     pms,
	}
}

// Pagination Pagination.
type Pagination struct {
	PageSize int `form:"page_size" json:"page_size"`
	PageNum  int `form:"page_num" json:"page_num"`
}

// Verify verify the value of pageNum and pageSize.
func (p *Pagination) Verify() error {
	if p.PageNum < 0 {
		return ecode.MerlinIllegalPageNumErr
	} else if p.PageNum == 0 {
		p.PageNum = DefaultPageNum
	}
	if p.PageSize < 0 {
		return ecode.MerlinIllegalPageSizeErr
	} else if p.PageSize == 0 {
		p.PageSize = DefaultPageSize
	}
	return nil
}

// QueryMachineRequest Query Machine Request.
type QueryMachineRequest struct {
	Pagination
	TreeNode
	MachineName string `form:"machine_name"`
	Username    string `form:"username"`
	Requester   string
}

// QueryMachineLogRequest Query Machine Log Request.
type QueryMachineLogRequest struct {
	Pagination
	MachineID   int64  `form:"machine_id"`
	MachineName string `form:"machine_name"`
	OperateUser string `form:"operate_user"`
	OperateType string `form:"operate_type"`
}

// QueryMobileMachineLogRequest Query Mobile Machine Log Request.
type QueryMobileMachineLogRequest struct {
	Pagination
	MachineID   int64  `form:"machine_id"`
	Serial      string `form:"serial"`
	OperateUser string `form:"operate_user"`
	OperateType string `form:"operate_type"`
}

// QueryMobileMachineErrorLogRequest Query Mobile Machine Error Log Request.
type QueryMobileMachineErrorLogRequest struct {
	Pagination
	MachineID int64 `form:"machine_id"`
}

// AboundMachineLog Abound Machine Log.
type AboundMachineLog struct {
	MachineLog
	Name string `json:"machine_name"`
}

// AboundMobileMachineLog Abound mobile Machine Log.
type AboundMobileMachineLog struct {
	MobileMachineLog
	Serial string `json:"serial"`
}

// ApplyEndTimeRequest Apply End Time Request.
type ApplyEndTimeRequest struct {
	MachineID    int64  `json:"machine_id"`
	ApplyEndTime string `json:"apply_end_time"`
	Auditor      string `json:"auditor"`
}

// AuditEndTimeRequest Audit End Time Request.
type AuditEndTimeRequest struct {
	AuditID     int64  `json:"audit_id"`
	AuditResult bool   `json:"audit_result"`
	Comment     string `json:"comment"`
}

// BeforeDelMachineFunc Before DelMachine Func.
type BeforeDelMachineFunc func(c context.Context, id int64, username string) error

// MachineStatusResponse Machine Status Response.
type MachineStatusResponse struct {
	Initialized  string   `json:"initialized"`
	PodScheduled string   `json:"pod_scheduled"`
	Ready        string   `json:"ready"`
	SynTree      string   `json:"syn_tree"`
	RetryCount   int      `json:"retry_count"`
	Log          string   `json:"log"`
	MachineEvent []string `json:"events"`
}

// CreatingMachineStatus Creating Machine Status.
func (msr *MachineStatusResponse) CreatingMachineStatus() int {
	if msr.Initialized == False {
		return CreatingMachineInMerlin
	}
	if msr.PodScheduled == False {
		return InitializeMachineInMerlin
	}
	if msr.Ready == False {
		return ReadyMachineInMerlin
	}
	if msr.SynTree == False {
		return SynTreeMachineInMerlin
	}
	return BootMachineInMerlin
}

// FailedMachineStatus Failed Machine Status.
func (msr *MachineStatusResponse) FailedMachineStatus() int {
	if msr.Initialized == False {
		return InitializedFailedMachineInMerlin
	}
	if msr.PodScheduled == False {
		return PodScheduledFailedMachineInMerlin
	}
	if msr.Ready == False {
		return ReadyFailedMachineInMerlin
	}
	if msr.SynTree == False {
		return SynTreeFailedMachineInMerlin
	}
	return 0
}

// InstanceMachineStatusResponse Instance Machine Status Response.
func InstanceMachineStatusResponse(machineStatus int) *MachineStatusResponse {
	switch machineStatus {
	case ImmediatelyFailedMachineInMerlin:
		return &MachineStatusResponse{}
	case InitializedFailedMachineInMerlin:
		return &MachineStatusResponse{Initialized: False, PodScheduled: False, Ready: False, SynTree: False, RetryCount: 0}
	case PodScheduledFailedMachineInMerlin:
		return &MachineStatusResponse{Initialized: True, PodScheduled: False, Ready: False, SynTree: False, RetryCount: 0}
	case ReadyFailedMachineInMerlin:
		return &MachineStatusResponse{Initialized: True, PodScheduled: True, Ready: False, SynTree: False, RetryCount: 0}
	case SynTreeFailedMachineInMerlin:
		return &MachineStatusResponse{Initialized: True, PodScheduled: True, Ready: True, SynTree: False, RetryCount: 0}
	case BootMachineInMerlin:
		return &MachineStatusResponse{Initialized: True, PodScheduled: True, Ready: True, SynTree: True, RetryCount: 0}
	case ShutdownMachineInMerlin:
		return &MachineStatusResponse{Initialized: True, PodScheduled: True, Ready: True, SynTree: True, RetryCount: 0}
	}
	return nil
}

// TreeNode tree node struct for merlin.
type TreeNode struct {
	BusinessUnit string `form:"business_unit"`
	Project      string `form:"project"`
	App          string `form:"app"`
}

// VerifyFieldValue verify that all fields is not empty.
func (t TreeNode) VerifyFieldValue() (err error) {
	if t.BusinessUnit == "" || t.Project == "" || t.App == "" {
		err = ecode.MerlinShouldTreeFullPath
	}
	return
}

// TreePath join all fields.
func (t *TreeNode) TreePath() string {
	return t.BusinessUnit + "." + t.Project + "." + t.App
}

// TreePathWithoutEmptyField join all fields except empty.
func (t *TreeNode) TreePathWithoutEmptyField() (treePath string) {
	if t.BusinessUnit == "" {
		return
	}
	treePath = "bilibili." + t.BusinessUnit
	if t.Project == "" {
		return
	}
	treePath = treePath + "." + t.Project
	if t.App == "" {
		return
	}
	treePath = treePath + "." + t.App
	return
}

// UpdateMachineNodeRequest request struct for updating machine node.
type UpdateMachineNodeRequest struct {
	MachineID int64          `json:"machine_id"`
	Nodes     []*MachineNode `json:"nodes"`
}

// VerifyNodes verify the Nodes field of UpdateMachineNodeRequest
func (u *UpdateMachineNodeRequest) VerifyNodes() error {
	l := len(u.Nodes)
	if l < 1 || l > 10 {
		return ecode.MerlinInvalidNodeAmountErr
	}
	return nil
}

// ToMachineNodes convert to machine nodes with injecting machine id.
func (u UpdateMachineNodeRequest) ToMachineNodes() []*MachineNode {
	for _, n := range u.Nodes {
		n.MachineID = u.MachineID
	}
	return u.Nodes
}

// QueryMobileDeviceRequest Query Device Farm Request.
type QueryMobileDeviceRequest struct {
	Pagination
	MobileID  int64  `json:"mobile_id"`
	Serial    string `json:"serial"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	OwnerName string `json:"owner_name"`
	CPU       string `json:"cpu"`
	Version   string `json:"version"`
	Mode      string `json:"mode"`
	Type      int    `json:"type"`
	State     string `json:"state"`
	Usage     int    `json:"usage"`
	Online    bool   `json:"online"`
}

// PaginateMobileMachines Paginate Device.
type PaginateMobileMachines struct {
	Total          int64                    `json:"total"`
	PageNum        int                      `json:"page_num"`
	PageSize       int                      `json:"page_size"`
	MobileMachines []*MobileMachineResponse `json:"mobile_devices"`
}

// QueryMachine2ImageLogRequest Query Machine to  Image Log Request.
type QueryMachine2ImageLogRequest struct {
	Pagination
	MachineID int64 `form:"machine_id"`
}

// PaginateHubImageLog Paginate Hub Image Log.
type PaginateHubImageLog struct {
	Total        int64          `json:"total"`
	PageNum      int            `json:"page_num"`
	PageSize     int            `json:"page_size"`
	HubImageLogs []*HubImageLog `json:"hub_image_logs"`
}

// ImageConfiguration Image Configuration.
type ImageConfiguration struct {
	ImageFullName string `json:"image_full_name"`
	PaasMachineSystem
}
