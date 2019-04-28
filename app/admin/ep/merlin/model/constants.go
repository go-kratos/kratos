package model

//merlin machine status.
const (
	// ImmediatelyFailedMachineInMerlin Paas immediately returns failed, when you create some new machines.
	ImmediatelyFailedMachineInMerlin = -200
	// InitializedFailedMachineInMerlin Scheduled detected that paas failed to execute Initialized task.
	InitializedFailedMachineInMerlin = -201
	// PodScheduledFailedMachineInMerlin Scheduled detected that paas failed to execute PodScheduled task.
	PodScheduledFailedMachineInMerlin = -202
	// ReadyFailedMachineInMerlin Scheduled detected that paas failed to execute Ready task.
	ReadyFailedMachineInMerlin = -203
	// ReadyFailedMachineInMerlin Scheduled detected that paas failed to sync the node of tree service.
	SynTreeFailedMachineInMerlin = -204
	// CreateTagFailedMachineInMerlin detected that merlin failed to create tag.
	CreateTagFailedMachineInMerlin = -205

	// RemovedMachineInMerlin the user removed the machine.
	RemovedMachineInMerlin = -100

	// CreatingMachineInMerlin Paas is creating the machine now.
	CreatingMachineInMerlin = 0
	// InitializeMachineInMerlin Paas is executing Initialize task.
	InitializeMachineInMerlin = 1
	// InitializeMachineInMerlin Paas is executing PodScheduled task.
	PodScheduledMachineInMerlin = 2
	// InitializeMachineInMerlin Paas is executing Ready task.
	ReadyMachineInMerlin = 3
	// InitializeMachineInMerlin Paas is syncing the node of tree service.
	SynTreeMachineInMerlin = 4

	// BootMachineInMerlin The machine is turned on.
	BootMachineInMerlin = 100

	// ShutdownMachineInMerlin The machine is off state.
	ShutdownMachineInMerlin = 200
)

// paas return response status.
const (
	// CreateFailedMachineInPaas Paas created the machine failed
	CreateFailedMachineInPaas = 0
	// CreatingMachineInPass Paas is creating the machine now
	CreatingMachineInPass = 1

	// SuccessDeletePaasMachines success deleted paas machine
	SuccessDeletePaasMachines = 1
)

// pagination.
const (
	DefaultPageSize = 5
	DefaultPageNum  = 1
)

// snapshot status
const (
	SnapshotInit    = "快照初始化"
	SnapshotDoing   = "快照进行中"
	SnapshotSuccess = "快照已完成"
	SnapShotFailed  = "快照失败"
)

// machine log.
const (
	GenForMachineLog      = "创建"
	DeleteForMachineLog   = "删除"
	TransferForMachineLog = "转移"

	OperationSuccessForMachineLog = "成功"
	OperationFailedForMachineLog  = "失败"

	MBStartLog    = "移动设备启动"
	MBShutDownLog = "移动设备关闭"
	MBBindLog     = "移动设备绑定"
	MBReleaseLog  = "移动设备释放"
	MBLendOutLog  = "移动设备借出"
	MBReturnLog   = "移动设备归还"
)

// mobile machine usage.
const (
	MBInUse     = 1
	MBFree      = 2
	MBNoConnect = 3
)

// mobile machine usage.
const (
	MBOnline  = 1  //在线
	MBOffline = -1 //离线
	MBHostDel = -2 //删除
)

// is Simulator or RealMachine.
const (
	MBSimulator = 1 //虚拟机
	MBReal      = 0 //真机
)

// is real machine on site or not.
const (
	MBOnSite  = 0 //归还
	MBLendOut = 1 //借出
)

// machine suffix.
const (
	MachinePodNameSuffix = "-0"
)

// delay log.
const (
	DelayMachineEndTime       = "手动延期"
	CancelDelayMachineEndTime = "取消延期"
	AuditDelayMachineEndTime  = "审批延期"
)

// delay status.
const (
	DelayStatusInit    = 0 //延期状态初始化
	DelayStatusAuto    = 1 //可自动延期
	DelayStatusApply   = 2 //可申请延期
	DelayStatusDisable = 3 //不可申请延期
)

// apply status.
const (
	ApplyDelayInit    = "申请延期中"
	ApplyDelayCancel  = "申请延期取消"
	ApplyDelayApprove = "申请延期通过"
	ApplyDelayDecline = "申请延期驳回"
)

// task.
const (
	DeleteMachine = "DeleteMachine" // 删除机器
)

// task status.
const (
	TaskInit    = 0  //未开始
	TaskDone    = 1  //已执行成功
	TaskFailed  = 2  //已执行失败
	TaskDeleted = -1 //任务删除
)

// machine expired status.
const (
	MailTypeMachineWillExpired      = 1  //机器将要过期
	MailTypeMachineDeleted          = 2  //机器已被删除
	MailTypeMachineTransfer         = 5  //机器转移
	MailTypeTaskDeleteMachineFailed = 11 //删除机器任务失败

	MailTypeApplyDelayMachineEndTime = 3 //申请延长机器过期时间
	MailTypeAuditDelayMachineEndTime = 4 //审核延长机器过期时间

)

// image operate type.
const (
	ImageNoSnapshot    = 0
	ImagePullAndPush   = 1
	ImagePull          = 2
	ImageTag           = 3
	ImagePush          = 4
	ImageMachine2Image = 5
)

// image operate err
const (
	ImageSuccess  = 0
	ImageInit     = -1
	ImagePullErr  = -2
	ImageReTagErr = -3
	ImagePushErr  = -4
)

// time format.
const (
	TimeFormat = "2006-01-02 15:04:05"
)

// bool str.
const (
	False = "False"
	True  = "True"
)

// tree role admin.
const (
	TreeRoleAdmin = 1
)

// image status.
const (
	AliveImageStatus   = 1
	DeletedImageStatus = 2
)

// other
const (
	Success = "success"
)

// machine ratio
const (
	CPURatio    = 1000
	MemoryRatio = 1024
)
