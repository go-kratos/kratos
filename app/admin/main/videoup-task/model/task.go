package model

import (
	"time"

	"go-common/library/log"
	"go-common/library/xstr"
)

var (
	// TookTypeMinute video task took time in 1 minute
	TookTypeMinute = int8(1)
	// TookTypeHalfHour video task took time in 10 minutes
	TookTypeHalfHour = int8(2)
	// TaskStateUnclaimed video task belongs to nobody
	TaskStateUnclaimed = int8(0)
	// TaskStateUntreated video task not submit
	TaskStateUntreated = int8(1)
	// TaskStateCompleted video task completed
	TaskStateCompleted = int8(2)
	// TaskStateDelayed video task delayed
	TaskStateDelayed = int8(3)
	// TaskStateClosed video task closed
	TaskStateClosed = int8(4)

	// PoolForFirst 一审任务池
	PoolForFirst = int8(0)
	// PoolForSecond 二审任务池
	PoolForSecond = int8(1)

	// TypeRealTime 实时任务
	TypeRealTime = int8(0)
	// TypeDispatched 已分发任务
	TypeDispatched = int8(1)
	// TypeFinished 结束任务
	TypeFinished = int8(2)
	// TypeDelay 延时任务
	TypeDelay = int8(3)
	// TypeReview 复审任务
	TypeReview = int8(5)
	// TypeUpDelete 已删除任务
	TypeUpDelete = int8(6)
	// TypeSpecialWait 特殊停滞任务
	TypeSpecialWait = int8(7)

	// SubjectForNormal 普通任务
	SubjectForNormal = int8(0) //normal task subject
	// SubjectForTask 指派任务
	SubjectForTask = int8(1) //specified task subject

	_taskdispatchstate = map[int8]struct{}{
		TypeRealTime:    {},
		TypeDispatched:  {},
		TypeFinished:    {},
		TypeDelay:       {},
		TypeReview:      {},
		TypeUpDelete:    {},
		TypeSpecialWait: {},
	}
	// WLVConf 默认值
	WLVConf = &WeightVC{
		MaxWeight:  int64(200000), //最大权重值
		SubRelease: int64(18),     //指派再释放的任务
		//特殊任务参数
		Slv1: int64(8),  // 普通用户>=1W粉
		Slv2: int64(10), // 普通用户>=10W粉
		Slv3: int64(12), // 优质用户<1W粉
		Slv4: int64(15), // 优质用户>=1W粉
		Slv5: int64(18), // 优质用户>=10W粉
		Slv6: int64(6),  // 高危用户>=10W粉
		Slv7: int64(0),  // 其他高危
		//普通任务参数
		Nlv1: int64(3),  // 等待时长 9-15
		Nlv2: int64(6),  // 等待时长 15-27
		Nlv3: int64(9),  // 等待时长 27-45
		Nlv4: int64(12), // 等待时长 >45
		Nlv5: int64(0),  // 等待时长 <=9
		//定时任务参数
		Tlv1: int64(3),  // 距离发布2h-4h
		Tlv2: int64(9),  // 距离发布1-2h
		Tlv3: int64(21), // 距离发布 <1h
		Tlv4: int64(0),  // 距离发布 > 4h
	}
)

const (
	// ActionHandsUP 0签入
	ActionHandsUP = int8(0)
	// ActionHandsOFF 1签出
	ActionHandsOFF = int8(1)
	// ActionSubmit 2提交
	ActionSubmit = int8(2)
	// ActionDelay  3延迟
	ActionDelay = int8(3)
	// ActionClose  4关闭
	ActionClose = int8(4)
	//ActionOldSubmit 5旧一审提交
	ActionOldSubmit = int8(5)
	//ActionTaskDelete 10任务删除
	ActionTaskDelete = int8(10)
	//ActionDispatch 分配
	ActionDispatch = int8(6)
	//ActionRelease 释放（拒审）
	ActionRelease = int8(7)

	// WConfMid 按照mid配置权重
	WConfMid = int8(0)
	// WConfTaskID 按照taskid配置权重
	WConfTaskID = int8(1)
	// WConfType 按照分区配置权重
	WConfType = int8(2)
	// WConfUpFrom 按照投稿来源配置权重
	WConfUpFrom = int8(3)
	// WConfRelease 指派任务释放
	WConfRelease = int8(4)
	// TimeFormatSec 时间格式化
	TimeFormatSec = "2006-01-02 15:04:05"

	// TaskLeader 组长
	TaskLeader int8 = 1
	// TaskMember 组员
	TaskMember int8 = 2
)

// MemberStat .
type MemberStat struct {
	UID           int64 `json:"uid"`
	DispatchCount int64 `json:"dispatch"`
	ReleaseCount  int64 `json:"release"`
	SubmitCount   int64 `json:"submit"`       // 总提交数
	OSubmitCount  int64 `json:"oldsubmit"`    // action=5
	NSubmitCount  int64 `json:"newsubmit"`    // action=2
	BelongCount   int64 `json:"belong"`       // 总推送数 分配到的-释法的
	PassCount     int64 `json:"pass"`         // 提交并且通过的数量
	NormalCount   int64 `json:"normalCount"`  // 普通任务数量
	SubjectCount  int64 `json:"specialCount"` // 指派任务数量

	InTime       string  `json:"inTime"`       // 最近开始时间
	QuitTime     string  `json:"quitTime"`     // 最近退出时间
	CompleteRate string  `json:"completeRate"` // 处理率
	PassRate     string  `json:"passRate"`     // 通过率
	SumDu        int64   `json:"-"`
	SumDuration  string  `json:"videoDuration"` // 累积处理视频总时长
	AvgUt        float64 `json:"-"`
	AvgUtime     string  `json:"avgTime"` // 单任务平均耗时
}

// CfWeightDesc 权重配置文字描述
func CfWeightDesc(radio int8) (desc string) {
	switch radio {
	case WConfMid:
		desc = "mid配置"
	case WConfTaskID:
		desc = "taskid配置"
	case WConfType:
		desc = "分区配置"
	case WConfUpFrom:
		desc = "投稿来源"
	case WConfRelease:
		desc = "指派释放"
	default:
		desc = "其他配置"
	}
	return
}

// IsDispatch 判断任务状态
func IsDispatch(st int8) bool {
	if _, ok := _taskdispatchstate[st]; ok {
		return true
	}
	return false
}

// ParseWeightConf 解析权重配置
func ParseWeightConf(twc *WeightConf, uid int64, uname string) (mcases map[int64]*WCItem, IsTaskID bool, err error) {
	var (
		ids []int64
	)
	mcases = make(map[int64]*WCItem)
	if ids, err = xstr.SplitInts(twc.Ids); err != nil {
		log.Error("ParseWeightConfig Config(%v) parse error(%v) Idlist(%s)", twc, err)
		return nil, false, err
	}

	for _, id := range ids {
		wci := &WCItem{
			CID:    id,
			Radio:  twc.Radio,
			Rule:   twc.Rule,
			Weight: twc.Weight,
			Uname:  uname,
			Desc:   twc.Desc,
			Bt:     twc.Bt,
			Et:     twc.Et,
			Mtime:  NewFormatTime(time.Now()),
		}
		if twc.Radio == WConfTaskID {
			IsTaskID = true
		}

		mcases[id] = wci
	}
	return
}

// WeightVC Weight Value Config 权重分值配置
type WeightVC struct {
	MaxWeight  int64 `json:"maxweight" form:"maxweight" default:"20000"`
	SubRelease int64 `json:"subrelease" form:"subrelease" default:"18"`
	Slv1       int64 `json:"slv1" form:"slv1" default:"8"`
	Slv2       int64 `json:"slv2" form:"slv2" default:"10"`
	Slv3       int64 `json:"slv3" form:"slv3" default:"12"`
	Slv4       int64 `json:"slv4" form:"slv4" default:"15"`
	Slv5       int64 `json:"slv5" form:"slv5" default:"18"`
	Slv6       int64 `json:"slv6" form:"slv6" default:"6"`
	Slv7       int64 `json:"slv7" form:"slv7" default:"0"`
	Nlv1       int64 `json:"nlv1" form:"nlv1" default:"3"`
	Nlv2       int64 `json:"nlv2" form:"nlv2" default:"6"`
	Nlv3       int64 `json:"nlv3" form:"nlv3" default:"9"`
	Nlv4       int64 `json:"nlv4" form:"nlv4" default:"12"`
	Nlv5       int64 `json:"nlv5" form:"nlv5" default:"0"`
	Tlv1       int64 `json:"tlv1" form:"tlv1" default:"3"`
	Tlv2       int64 `json:"tlv2" form:"tlv2" default:"9"`
	Tlv3       int64 `json:"tlv3" form:"tlv3" default:"21"`
	Tlv4       int64 `json:"tlv4" form:"tlv4" default:"0"`
}

// Task 审核任务
type Task struct {
	ID      int64      `json:"id"`
	Pool    int8       `json:"pool"`
	Subject int8       `json:"subject"`
	AdminID int64      `json:"adminid"`
	Aid     int64      `json:"aid"`
	Cid     int64      `json:"cid"`
	UID     int64      `json:"uid"`
	State   int8       `json:"state"`
	UTime   int64      `json:"utime"`
	CTime   FormatTime `json:"ctime"`
	MTime   FormatTime `json:"mtime"`
	DTime   FormatTime `json:"dtime"`
	GTime   FormatTime `json:"gtime"`
	PTime   FormatTime `json:"ptime"`
	Weight  int64      `json:"weight"`
	Mid     int64      `json:"mid"`
}

// TaskWeightLog 权重变更日志
type TaskWeightLog struct {
	TaskID    int64      `json:"taskid"`
	Mid       int64      `json:"mid"`
	Weight    int64      `json:"weight"`
	CWeight   int64      `json:"cweight"`
	NWeight   int64      `json:"nweight"`
	SWeight   int64      `json:"sweight"`
	TWeight   int64      `json:"tweight"`
	Uptime    FormatTime `json:"uptime"`
	Creator   string     `json:"creator"`   //创作者
	UpSpecial []int8     `json:"upspecial"` //标记是否优质，劣质用户
	Fans      int64      `json:"fans"`      //粉丝数
	Wait      float64    `json:"wait"`      //等待时长
	Ptime     string     `json:"ptime,omitempty"`
	CfItems   []*WCItem  `json:"cfitems,omitempty"`
	Desc      string     `json:"desc,omitempty"` // 配置描述
}

// TaskPriority 任务参数
type TaskPriority struct {
	TaskID  int64      `json:"taskid"`
	Weight  int64      `json:"weight"` //权重总值
	State   int8       `json:"state"`  //任务状态
	Mid     int64      `json:"mid"`
	Special int8       `json:"special"` //特殊任务
	Ctime   FormatTime `json:"ctime"`   //任务生成时间
	Ptime   FormatTime `json:"ptime"`   //定时发布时间
	CfItems []*WCItem  `json:"cfitems,omitempty"`

	Fans      int64  `json:"fans"`     //粉丝数
	AccFailed bool   `json:"accfaild"` //账号查询是否失败
	UpGroups  []int8 `json:"ugs"`      //分组
	UpFrom    int8   `json:"upfrom"`   //来源
	TypeID    int16  `json:"typeid"`   //分区
}

// WeightConf 任务权重配置
type WeightConf struct {
	Radio  int8       `form:"radio"`                       // 0,mid，1，taskid，2，分区, 3, 投稿来源
	Ids    string     `form:"ids"  validate:"required"`    // id列表，逗号分隔
	Rule   int8       `form:"rule"`                        // 0,动态权重，1，静态权重
	Weight int64      `form:"weight"  validate:"required"` // 配置的权重
	Desc   string     `form:"desc"  validate:"required"`   // 描述信息
	Bt     FormatTime `form:"bt"`                          //配置生效开始时间
	Et     FormatTime `form:"et"`                          //配置生效结束时间
}

// WCItem task weight config item
type WCItem struct {
	Radio    int8       `json:"radio"`
	ID       int64      `json:"id,omitempty"`
	CID      int64      `json:"cid"` // config id 四种配置通用
	UID      int64      `json:"uid,omitempty"`
	Uname    string     `json:"user,omitempty"`
	TypeName string     `json:"typename,omitempty"`
	UpFrom   string     `json:"upfrom,omitempty"`
	Rule     int8       `json:"rule"`
	State    int8       `json:"state"`
	Weight   int64      `json:"weight,omitempty"`
	Mtime    FormatTime `json:"mtime,omitempty"`
	Desc     string     `json:"desc,omitempty"`
	FileName string     `json:"filename,omitempty"`
	Title    string     `json:"title,omitempty"`
	Vid      int64      `json:"vid,omitempty"`
	Creator  string     `json:"creator,omitempty"`
	Fans     int64      `json:"fans,omitempty"`
	Bt       FormatTime `json:"bt,omitempty"`
	Et       FormatTime `json:"et,omitempty"`
}

// Confs 权重配置筛选参数
type Confs struct {
	Radio    int8       `form:"radio" default:"1"`
	Cid      int64      `form:"cid" default:"-1"`
	Operator string     `form:"operator"`
	Bt       FormatTime `form:"bt"`
	Et       FormatTime `form:"et"`
	Rule     int8       `form:"rule" default:"-1"`
	State    int        `form:"state"`
	Pn       int        `form:"page" default:"1"`
	Ps       int        `form:"ps" default:"20"`
}

// TaskTook 一审耗时
type TaskTook struct {
	ID     int64     `json:"id"`
	M90    int       `json:"m90"`
	M80    int       `json:"m80"`
	M60    int       `json:"m60"`
	M50    int       `json:"m50"`
	TypeID int8      `json:"type"`
	Ctime  time.Time `json:"ctime"`
	Mtime  time.Time `json:"-"`
}

// AuthRole 一审任务角色
type AuthRole struct {
	ID       int64     `json:"id"`
	UID      int64     `json:"uid"`
	Role     int8      `json:"role"`
	UserName string    `json:"username"`
	NickName string    `json:"nickname"`
	Ctime    time.Time `json:"ctime"`
	Mtime    time.Time `json:"mtime"`
}

// Consumers 组员信息
type Consumers struct {
	ID       int64      `json:"id"`
	UID      int64      `json:"uid"`
	UserName string     `json:"username"`
	State    int8       `json:"state"`
	Ctime    FormatTime `json:"ctime"`
	Mtime    FormatTime `json:"mtime"`
	LastOut  string     `json:"lastout,omitempty"`
}

// ConsumerLog 组员日志
type ConsumerLog struct {
	UID    int64  `json:"uid"`
	Uname  string `json:"uname"`
	Action int8   `json:"action"`
	Ctime  string `json:"ctime"`
	Desc   string `json:"desc"`
}

// InQuit 组员日志
type InQuit struct {
	Date    string `json:"date"`
	UID     int64  `json:"uid"`
	Uname   string `json:"uname"`
	InTime  string `json:"inTime"`
	OutTime string `json:"quitTime"`
}

// TaskForLog 释放任务
type TaskForLog struct {
	ID      int64
	Cid     int64
	Subject int8
	Mtime   time.Time
}
