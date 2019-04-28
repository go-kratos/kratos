package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Machine Machine.
type Machine struct {
	ID           int64     `json:"id" gorm:"column:id"`
	Name         string    `json:"name" gorm:"column:name"`
	PodName      string    `json:"pod_name,omitempty" gorm:"column:pod_name"`
	Status       int       `json:"status" gorm:"column:status"`
	Username     string    `json:"username" gorm:"column:username"`
	BusinessUnit string    `json:"business_unit,omitempty" gorm:"column:business_unit"`
	Project      string    `json:"project,omitempty" gorm:"column:project"`
	App          string    `json:"app,omitempty" gorm:"column:app"`
	ClusterID    int64     `json:"cluster_id,omitempty" gorm:"column:cluster_id"`
	NetworkID    int64     `json:"network_id,omitempty" gorm:"column:network_id"`
	Ctime        time.Time `json:"ctime" gorm:"column:ctime;default:current_timestamp"`
	Utime        time.Time `json:"utime" gorm:"column:mtime;default:current_timestamp on update current_timestamp"`
	UpdateBy     string    `json:"update_by" gorm:"column:update_by"`
	EndTime      time.Time `json:"end_time" gorm:"column:end_time"`
	Comment      string    `json:"comment" gorm:"column:comment"`
	DelayStatus  int       `json:"delay_status" gorm:"column:delay_status"`
}

// AfterCreate After Create.
func (m *Machine) AfterCreate(db *gorm.DB) (err error) {
	if err = db.Model(m).Where("name = ?", m.Name).Find(&m).Error; err != nil {
		return
	}
	m.EndTime = m.Ctime.AddDate(0, 3, 0)
	if err = db.Model(&Machine{}).Where("id = ?", m.ID).Update("end_time", m.EndTime).Error; err != nil {
		return
	}
	return
}

// AfterCreate After Create.
func (h *HubImageLog) AfterCreate(db *gorm.DB) (err error) {
	err = db.Model(h).Where("imagetag = ?", h.ImageTag).Find(&h).Error
	return
}

// IsFailed Is Failed.
func (m *Machine) IsFailed() bool {
	return m.Status >= ImmediatelyFailedMachineInMerlin && m.Status < RemovedMachineInMerlin
}

// IsDeleted Is Deleted.
func (m *Machine) IsDeleted() bool {
	return m.Status >= RemovedMachineInMerlin && m.Status < CreatingMachineInMerlin
}

// IsCreating Is Creating.
func (m *Machine) IsCreating() bool {
	return m.Status >= CreatingMachineInMerlin && m.Status < BootMachineInMerlin
}

// IsBoot Is Boot.
func (m *Machine) IsBoot() bool {
	return m.Status >= BootMachineInMerlin && m.Status < ShutdownMachineInMerlin
}

// IsShutdown Is Shutdown.
func (m *Machine) IsShutdown() bool {
	return m.Status >= ShutdownMachineInMerlin && m.Status < 300
}

// ToTreeNode return Tree node.
func (m *Machine) ToTreeNode() *TreeNode {
	return &TreeNode{
		BusinessUnit: m.BusinessUnit,
		Project:      m.Project,
		App:          m.App,
	}
}

// ToMachineLog generate a machine log struct.
func (m *Machine) ToMachineLog() (ml *MachineLog) {
	ml = &MachineLog{
		OperateType: GenForMachineLog,
		Username:    m.Username,
		MachineID:   m.ID,
	}
	if m.Status == CreatingMachineInMerlin {
		ml.OperateResult = OperationSuccessForMachineLog
	} else if m.Status == ImmediatelyFailedMachineInMerlin {
		ml.OperateResult = OperationFailedForMachineLog
	}
	return
}

// MachineLog Machine Log.
type MachineLog struct {
	ID            int64     `json:"-" gorm:"column:id"`
	Username      string    `json:"username" gorm:"column:username"`
	MachineID     int64     `json:"machine_id" gorm:"column:machine_id"`
	OperateType   string    `json:"operate_type" gorm:"column:operation_type"`
	OperateResult string    `json:"operate_result" gorm:"column:operation_result"`
	OperateTime   time.Time `json:"operate_time" gorm:"column:ctime;default:current_timestamp"`
	UTime         time.Time `json:"-" gorm:"column:mtime;default:current_timestamp on update current_timestamp"`
}

// MobileMachineLog Mobile Machine Log.
type MobileMachineLog struct {
	ID            int64     `json:"-" gorm:"column:id"`
	Username      string    `json:"username" gorm:"column:username"`
	MachineID     int64     `json:"machine_id" gorm:"column:machine_id"`
	OperateType   string    `json:"operate_type" gorm:"column:operation_type"`
	OperateResult string    `json:"operate_result" gorm:"column:operation_result"`
	OperateTime   time.Time `json:"operate_time" gorm:"column:ctime;default:current_timestamp"`
	UTime         time.Time `json:"-" gorm:"column:mtime;default:current_timestamp on update current_timestamp"`
}

// MobileMachineErrorLog Mobile Machine Error Log.
type MobileMachineErrorLog struct {
	ID           int64     `json:"id" gorm:"column:id"`
	MachineID    int64     `json:"machine_id" gorm:"column:machine_id"`
	SerialName   string    `json:"serial" gorm:"column:serial"`
	ErrorMessage string    `json:"error_message" gorm:"column:error_message"`
	ErrorCode    int       `json:"error_code" gorm:"column:error_code"`
	CTime        time.Time `json:"create_time" gorm:"column:ctime;default:current_timestamp"`
	UTime        time.Time `json:"-" gorm:"column:mtime;default:current_timestamp on update current_timestamp"`
}

// Snapshot  Snapshot.
type Snapshot struct {
	ID        int64     `gorm:"column:id"`
	Name      string    `gorm:"column:name"`
	MachineID int64     `gorm:"column:machine_id"`
	UserID    int64     `gorm:"column:user_id"`
	Ctime     time.Time `json:"ctime" gorm:"column:ctime;default:current_timestamp"`
	Utime     time.Time `json:"utime" gorm:"column:mtime;default:current_timestamp on update current_timestamp"`
	UpdateBy  int       `gorm:"column:update_by"`
	Comment   string    `gorm:"column:comment"`
}

// SnapshotLog Snapshot Log.
type SnapshotLog struct {
	ID            int64     `gorm:"column:id"`
	UserID        int64     `gorm:"column:user_id"`
	SnapshotID    int64     `gorm:"column:snapshot_id"`
	OperateType   string    `gorm:"column:operation_type"`
	OperateResult int       `gorm:"column:operation_result"`
	OperateTime   time.Time `gorm:"column:operation_time;default:current_timestamp"`
}

// Task Task.
type Task struct {
	ID          int64     `gorm:"column:id"`
	TYPE        string    `gorm:"column:type"`
	ExecuteTime time.Time `gorm:"column:execute_time"`
	MachineID   int64     `gorm:"column:machine_id"`
	Status      int       `gorm:"column:status"`
	Ctime       time.Time `gorm:"column:ctime;default:current_timestamp"`
	UTime       time.Time `gorm:"column:mtime;default:current_timestamp on update current_timestamp"`
}

// User User.
type User struct {
	ID    int64     `json:"id" gorm:"auto_increment;primary_key;column:id"`
	Name  string    `json:"username" gorm:"column:name"`
	EMail string    `json:"email" gorm:"column:email"`
	CTime time.Time `gorm:"column:ctime;default:current_timestamp"`
	UTime time.Time `gorm:"column:mtime;default:current_timestamp on update current_timestamp"`
}

// Image Image.
type Image struct {
	ID          int64     `json:"id" gorm:"auto_increment;primary_key;column:id"`
	Name        string    `json:"name" gorm:"varchar(100);column:name"`
	Status      int       `json:"status" gorm:"not null;column:status"`
	OS          string    `json:"os" gorm:"not null;column:os"`
	Version     string    `json:"version" gorm:"not null;column:version"`
	Description string    `json:"description" gorm:"column:description"`
	CreatedBy   string    `json:"created_by" gorm:"column:created_by"`
	UpdatedBy   string    `json:"updated_by" gorm:"column:updated_by"`
	Ctime       time.Time `json:"ctime" gorm:"column:ctime;default:current_timestamp"`
	Utime       time.Time `json:"utime" gorm:"column:mtime;default:current_timestamp on update current_timestamp"`
}

// MachinePackage MachinePackage.
type MachinePackage struct {
	ID              int64     `json:"id" gorm:"column:id"`
	Name            string    `json:"name" gorm:"column:name"`
	CPUCore         int       `json:"cpu_request" gorm:"column:cpu_core"`
	Memory          int       `json:"memory_request" gorm:"column:memory"`
	StorageCapacity int       `json:"disk_request" gorm:"column:storage_capacity"`
	CTime           time.Time `json:"ctime" gorm:"column:ctime;default:current_timestamp"`
	UTime           time.Time `json:"utime" gorm:"column:mtime;default:current_timestamp on update current_timestamp"`
}

// MailLog MailLog.
type MailLog struct {
	ID           int64     `gorm:"column:id"`
	ReceiverName string    `gorm:"column:receiver_name"`
	MailType     int       `gorm:"column:mail_type"`
	SendHead     string    `gorm:"column:send_head"`
	SendContext  string    `gorm:"column:send_context"`
	SendTime     time.Time `gorm:"column:ctime;default:current_timestamp"`
	UTime        time.Time `gorm:"column:mtime;default:current_timestamp on update current_timestamp"`
}

// ToPaasQueryAndDelMachineRequest To Paas Query And Del Machine Request.
func (m *Machine) ToPaasQueryAndDelMachineRequest() (pqadmr *PaasQueryAndDelMachineRequest) {
	pqadmr = &PaasQueryAndDelMachineRequest{}
	pqadmr.Name = m.Name
	pqadmr.BusinessUnit = m.BusinessUnit
	pqadmr.Project = m.Project
	pqadmr.App = m.App
	pqadmr.ClusterID = m.ClusterID
	return
}

// ApplicationRecord ApplicationRecord.
type ApplicationRecord struct {
	ID           int64     `json:"id" gorm:"column:id"`
	Applicant    string    `json:"applicant" gorm:"column:applicant"`
	MachineID    int64     `json:"machine_id" gorm:"column:machine_id"`
	ApplyEndTime time.Time `json:"apply_end_time" gorm:"column:apply_end_time"`
	Status       string    `json:"status" gorm:"column:status"`
	Auditor      string    `json:"auditor" gorm:"column:auditor"`
	CTime        time.Time `json:"ctime" gorm:"column:ctime;default:current_timestamp"`
	UTime        time.Time `json:"utime" gorm:"column:mtime;default:current_timestamp on update current_timestamp"`
}

// MachineNode the node is associated with machine.
type MachineNode struct {
	ID           int64     `json:"id" gorm:"column:id"`
	MachineID    int64     `json:"machine_id" gorm:"column:machine_id"`
	BusinessUnit string    `json:"business_unit" gorm:"column:business_unit"`
	Project      string    `json:"project" gorm:"column:project"`
	App          string    `json:"app" gorm:"column:app"`
	TreeID       int64     `json:"tree_id,omitempty" gorm:"column:tree_id"`
	CTime        time.Time `json:"create_time" gorm:"column:ctime;default:current_timestamp"`
}

// HubImageLog Hub Image Log
type HubImageLog struct {
	ID          int64     `json:"id" gorm:"column:id"`
	MachineID   int64     `json:"machine_id" gorm:"column:machine_id"`
	UserName    string    `json:"username" gorm:"column:username"`
	ImageSrc    string    `json:"image_src" gorm:"column:imagesrc"`
	ImageTag    string    `json:"image_tag" gorm:"column:imagetag"`
	Status      int       `json:"status" gorm:"column:status"`
	OperateType int       `json:"operate_type" gorm:"column:operate_type"`
	CTime       time.Time `json:"create_time" gorm:"column:ctime;default:current_timestamp"`
	UTime       time.Time `json:"update_time" gorm:"column:mtime;default:current_timestamp on update current_timestamp"`
}

// MobileMachine the node is associated with MobileMachine.
type MobileMachine struct {
	ID           int64     `json:"id" gorm:"column:id"`
	Serial       string    `json:"serial" gorm:"column:serial"`
	Name         string    `json:"name" gorm:"column:name"`
	CPU          string    `json:"cpu" gorm:"column:cpu"`
	Version      string    `json:"version" gorm:"column:version"`
	Mode         string    `json:"mode" gorm:"column:mode"`
	State        string    `json:"state" gorm:"column:state"`
	Host         string    `json:"host" gorm:"column:host"`
	CTime        time.Time `json:"create_time" gorm:"column:ctime;default:current_timestamp"`
	MTime        time.Time `json:"update_time" gorm:"column:mtime;default:current_timestamp"`
	LastBindTime time.Time `json:"last_bind_time" gorm:"column:last_bind_time;default:current_timestamp"`
	OwnerName    string    `json:"owner_name" gorm:"column:owner_name"`
	Username     string    `json:"username" gorm:"column:username"`
	Type         int       `json:"type" gorm:"column:type"`
	EndTime      time.Time `json:"end_time" gorm:"column:end_time"`
	Alias        string    `json:"alias" gorm:"column:alias"`
	Comment      string    `json:"comment" gorm:"column:comment"`
	WsURL        string    `json:"wsurl" gorm:"column:wsurl"`
	UploadURL    string    `json:"upload_url" gorm:"column:upload_url"`
	Action       int       `json:"action" gorm:"column:action"`
	IsLendOut    int       `json:"is_lendout" gorm:"column:is_lendout"`
	UUID         string    `json:"uuid" gorm:"column:uuid"`
}

// MobileImage Mobile Image.
type MobileImage struct {
	ID       int64     `json:"id" gorm:"column:id"`
	Mode     string    `json:"mode" gorm:"column:mode"`
	CTime    time.Time `json:"ctime" gorm:"column:ctime;default:current_timestamp"`
	MTime    time.Time `json:"mtime" gorm:"column:mtime;default:current_timestamp"`
	ImageSrc string    `json:"image_src" gorm:"column:image_src"`
}

// MobileSyncLog MobileSyncLog.
type MobileSyncLog struct {
	ID        int64     `json:"id" gorm:"column:id"`
	UUID      string    `json:"uuid" gorm:"column:uuid"`
	AddCnt    int       `json:"add_count" gorm:"column:add_count"`
	UpdateCnt int       `json:"update_count" gorm:"column:update_count"`
	DeleteCnt int       `json:"delete_count" gorm:"column:delete_count"`
	TotalCnt  int       `json:"total_count" gorm:"column:total_count"`
	Status    int       `json:"status" gorm:"column:status"`
	CTime     time.Time `json:"ctime" gorm:"column:ctime;default:current_timestamp"`
	MTime     time.Time `json:"mtime" gorm:"column:mtime;default:current_timestamp"`
}

// MobileCategory MobileCategory.
type MobileCategory struct {
	CPUs     []string `json:"cpus"`
	Versions []string `json:"versions"`
	Modes    []string `json:"modes"`
	States   []string `json:"states"`
	Types    []int    `json:"types"`
	Usages   []int    `json:"usages"`
}

// SnapshotRecord Snapshot Record
type SnapshotRecord struct {
	ID        int64     `json:"id" gorm:"column:id"`
	MachineID int64     `json:"machine_id" gorm:"column:machine_id"`
	Username  string    `json:"username" gorm:"column:username"`
	ImageName string    `json:"image_name" gorm:"column:image_name"`
	Status    string    `json:"status" gorm:"column:status"`
	Ctime     time.Time `json:"ctime" gorm:"column:ctime;default:current_timestamp"`
	MTime     time.Time `json:"mtime" gorm:"column:mtime;default:current_timestamp on update current_timestamp"`
}

// HubImageConf Hub Image Conf
type HubImageConf struct {
	ID        int64     `json:"id" gorm:"column:id"`
	ImageName string    `json:"image_name" gorm:"column:image_name"`
	UpdateBy  string    `json:"update_by" gorm:"column:update_by"`
	Command   string    `json:"command" gorm:"column:command"`
	Envs      string    `json:"environments" gorm:"column:environments"`
	Hosts     string    `json:"hosts" gorm:"column:hosts"`
	CTime     time.Time `json:"create_time" gorm:"column:ctime;default:current_timestamp"`
	UTime     time.Time `json:"update_time" gorm:"column:mtime;default:current_timestamp on update current_timestamp"`
}
