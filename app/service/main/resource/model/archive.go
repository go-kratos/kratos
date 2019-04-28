package model

// resource archive const
const (
	// StateOpen 开放浏览
	StateOpen = int8(0)
	// StateOrange 橙色通过
	StateOrange = int8(1)
	// StateForbidWait 待审
	StateForbidWait = int8(-1)
	// StateForbidRecycle 被打回
	StateForbidRecycle = int8(-2)
	// StateForbidPolice 网警锁定
	StateForbidPolice = int8(-3)
	// StateForbidLock 被锁定
	StateForbidLock = int8(-4)
	// StateForbidFackLock 管理员锁定（可浏览）
	StateForbidFackLock = int8(-5)
	// StateForbidFixed 修复待审
	StateForbidFixed = int8(-6)
	// StateForbidLater 暂缓审核
	StateForbidLater = int8(-7)
	// StateForbidPatched 补档待审
	StateForbidPatched = int8(-8)
	// StateForbidWaitXcode 等待转码
	StateForbidWaitXcode = int8(-9)
	// StateForbidAdminDelay 延迟审核
	StateForbidAdminDelay = int8(-10)
	// StateForbidFixing 视频源待修
	StateForbidFixing = int8(-11)
	// StateForbidStorageFail 转储失败
	StateForbidStorageFail = int8(-12)
	// StateForbidOnlyComment 允许评论待审
	StateForbidOnlyComment = int8(-13)
	// StateForbidTmpRecicle 临时回收站
	StateForbidTmpRecicle = int8(-14)
	// StateForbidDispatch 分发中
	StateForbidDispatch = int8(-15)
	// StateForbidXcodeFail 转码失败
	StateForbidXcodeFail = int8(-16)
	// StateForbitUpLoad 创建未提交
	StateForbitUpLoad = int8(-20) // NOTE:spell body can judge to change state
	// StateForbidSubmit 创建已提交
	StateForbidSubmit = int8(-30)
	// StateForbidUserDelay 定时发布
	StateForbidUserDelay = int8(-40)
	// StateForbidUpDelete 用户删除
	StateForbidUpDelete = int8(-100)

	// resource apply
	ApplyFirstAudit   = 0  // 待一审
	ApplySecondAudit  = 1  // 待二审
	ApplyNoAssignment = 2  // 未投放
	ApplyAssignment   = 3  // 已投放
	ApplyReject       = -1 // 已驳回
	ApplyRecall       = -2 // 已撤回
)

// Archive archive struct
type Archive struct {
	ID        int64  `json:"id"`
	Mid       int64  `json:"mid"`
	TypeID    int16  `json:"typeid"`
	HumanRank int    `json:"humanrank"`
	Duration  int    `json:"duration"`
	Title     string `json:"title"`
	Cover     string `json:"cover"`
	Content   string `json:"content"`
	Tag       string `json:"tag"`
	Attribute int32  `json:"attribute"`
	Copyright int8   `json:"copyright"`
	AreaLimit int8   `json:"arealimit"`
	State     int8   `json:"state"`
	Author    string `json:"author"`
	Access    int    `json:"access"`
	Forward   int    `json:"forward"`
	PubTime   string `json:"pubtime"`
	Reason    string `json:"reject_reason"`
	Round     int8   `json:"round"`
	CTime     string `json:"ctime"`
	MTime     string `json:"mtime"`
}
