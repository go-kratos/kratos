package newcomer

import (
	"go-common/library/time"
)

const (
	//RewardCanActivate reward receive state 0-可激活 >1-已激活不可点击>2-已过期不可点击
	RewardCanActivate int8 = iota
	//RewardActivatedNotClick reward activated state 1-已激活不可点击
	RewardActivatedNotClick
	//RewardExpireNotClick reward activated state 2-已过期不可点击
	RewardExpireNotClick
)

const (
	//DefualtTaskType 0-默认任务
	DefualtTaskType int8 = iota
	// NewcomerTaskType 1-新手任务
	NewcomerTaskType
	// AdvancedTaskType 2-进阶任务
	AdvancedTaskType
	// MonthTaskType 3-月常任务
	MonthTaskType
)

const (
	_ int8 = iota
	// Bcoin 1-B币券
	Bcoin
	// BigMember      2-大会员服务
	BigMember
	// MemberBuy      3-会员购
	MemberBuy
	// IncentivePlan 4-激励计划
	IncentivePlan
	// PersonalCenter 5-个人中心
	PersonalCenter
)

const (
	//UserTaskLevel0   未解锁任务
	UserTaskLevel0 int8 = iota
	//UserTaskLevel01  只解锁新手任务
	UserTaskLevel01
	//UserTaskLevel02 解锁新手与进阶任务
	UserTaskLevel02
)

const (
	//FreezeState 任务或奖励被冻结状态
	FreezeState = -1
	//NormalState 任务或奖励正常状态
	NormalState = 0
	//HiddenState 任务或奖励隐藏状态
	HiddenState = 1

	//RewardBaseType 基础奖励
	RewardBaseType = 0
	//RewardGiftType 礼包奖励
	RewardGiftType = 1

	//NoBindTask 用户未绑定任务
	NoBindTask = -1
	//BindTask 用户已绑定任务
	BindTask = 0

	//TaskIncomplete 任务未完成
	TaskIncomplete = -1
	//TaskCompleted 任务完成
	TaskCompleted = 0

	//RewardNotAvailable 奖励不可领取
	RewardNotAvailable = -1
	//RewardAvailable 奖励可领取
	RewardAvailable = 0
	//RewardReceived 奖励已领取
	RewardReceived = 1
	//RewardUnlock 奖励未解锁
	RewardUnlock = 2

	//RewardNeedActivate 奖励可激活
	RewardNeedActivate = 1
	//RewardNoneedActivate 奖励不可激活
	RewardNoneedActivate = 0

	//FromWeb web端
	FromWeb = 1
	//FromH5 h5端
	FromH5 = 2
)

const (
	_ int8 = iota
	//TargetType001 该UID下开放浏览的稿件≥1
	TargetType001
	//TargetType002 该UID分享自己视频的次数≥1
	TargetType002
	//TargetType003 该UID在创作学院的观看记录≥1
	TargetType003
	//TargetType004 该UID下所有avid的获得评论数≥3
	TargetType004
	//TargetType005 该UID下所有avid获得分享数≥3
	TargetType005
	//TargetType006 该UID的所有avid的获得收藏数≥5
	TargetType006
	//TargetType007 该UID下所有avid的获得硬币数≥5
	TargetType007
	//TargetType008 该UID下所有avid获得点赞数≥5
	TargetType008
	//TargetType009 该UID下所有avid的获得弹幕数≥5
	TargetType009
	//TargetType010 该UID的粉丝数≥10
	TargetType010
	//TargetType011 任务完成期间该UID的水印开关为打开状态
	TargetType011
	//TargetType012 该UID的关注列表含有“哔哩哔哩创作中心”
	TargetType012
	//TargetType013 用手机投稿上传视频
	TargetType013
	//TargetType014 该UID下开放浏览的稿件≥5
	TargetType014
	//TargetType015 该UID下任意avid的获得点击量≥1000
	TargetType015
	//TargetType016 该UID下任意avid的评论≥30
	TargetType016
	//TargetType017 该UID下任意avid的获得分享数≥10
	TargetType017
	//TargetType018 该UID下任意avid的获得收藏数≥30
	TargetType018
	//TargetType019 该UID下任意avid的获得硬币数≥50
	TargetType019
	//TargetType020 该UID下任意avid的获得点赞数≥50
	TargetType020
	//TargetType021 该UID下任意avid的获得弹幕数≥50
	TargetType021
	//TargetType022 该UID的粉丝数≥1000
	TargetType022
	//TargetType023 该UID的激励计划状态为已开通
	TargetType023
	//TargetType024 该UID粉丝勋章为开启状态
	TargetType024
)

const (
	_ int8 = iota
	//ArcUpCount UpCount get archives count
	ArcUpCount
	//AcaPlayCount get all play achive count.
	AcaPlayCount
	//DataUpStat get up stat from hbase
	DataUpStat
	//AccProfileWithStat get account
	AccProfileWithStat
	//WmWaterMark get watermark.
	WmWaterMark
	//AccRelation get all relation state.
	AccRelation
	//DataUpArchiveStat 获取最高播放/评论/弹幕/...数
	DataUpArchiveStat
	//OrderGrowAccountState 获取up主状态 type 类型 0 视频 2 专栏 3 素材.
	OrderGrowAccountState
	//MedalCheckMedal get medal
	MedalCheckMedal
)

const (
	//MsgFinishedCount 发送未完成任务状态
	MsgFinishedCount = 1
	//MsgForWaterMark 发送用户设置水印消息
	MsgForWaterMark = 1
	//MsgForAcademyFavVideo 发送用户已在创作学院观看过自己喜欢的视频的消息
	MsgForAcademyFavVideo = 2
	//MsgForGrowAccount 发送用户已在参加激励计划的消息
	MsgForGrowAccount = 3
	//MsgForOpenFansMedal 成功开通粉丝勋章
	MsgForOpenFansMedal = 4
)

var (
	// TaskRedirectMap task map for app
	TaskRedirectMap = map[string]map[int8][]string{
		"android": {
			TargetType001: []string{"去投稿", "bilibili://uper/user_center/add_archive/"},
			TargetType002: []string{"去分享", "activity://uper/manuscript-list/"},
			TargetType003: []string{"前往", "https://member.bilibili.com/college?from=task"},
			TargetType004: []string{"前往", "activity://uper/manuscript-list/"},
			TargetType005: []string{"前往", "activity://uper/manuscript-list/"},
			TargetType006: []string{"前往", "activity://uper/manuscript-list/"},
			TargetType007: []string{"前往", "activity://uper/manuscript-list/"},
			TargetType008: []string{"前往", "activity://uper/manuscript-list/"},
			TargetType009: []string{"前往", "activity://uper/manuscript-list/"},
			TargetType010: []string{"前往", "activity://uper/manuscript-list/"},
			TargetType011: []string{"去设置", "https://member.bilibili.com/studio/gabriel/watermark"},
			TargetType012: []string{"去关注", ""},
			TargetType013: []string{"去投稿", "bilibili://uper/user_center/add_archive/"},
			TargetType014: []string{"去投稿", "bilibili://uper/user_center/add_archive/"},
			TargetType015: []string{"前往", "activity://uper/manuscript-list/"},
			TargetType016: []string{"前往", "activity://uper/manuscript-list/"},
			TargetType017: []string{"前往", "activity://uper/manuscript-list/"},
			TargetType018: []string{"前往", "activity://uper/manuscript-list/"},
			TargetType019: []string{"前往", "activity://uper/manuscript-list/"},
			TargetType020: []string{"前往", "activity://uper/manuscript-list/"},
			TargetType021: []string{"前往", "activity://uper/manuscript-list/"},
			TargetType022: []string{"前往", "https://member.bilibili.com/studio/gabriel/fans-manage/overview"},
			TargetType023: []string{"去加入", "https://member.bilibili.com/studio/up-allowance-h5#/"},
			TargetType024: []string{"去开通", "https://member.bilibili.com/studio/gabriel/fans-manage/medal"},
		},
		"ios": {TargetType001: []string{"去投稿", "/uper/user_center/add_archive/"},
			TargetType002: []string{"去分享", "/uper/user_center/archive_list"},
			TargetType003: []string{"前往", "https://member.bilibili.com/college?from=task"},
			TargetType004: []string{"前往", "/uper/user_center/archive_list"},
			TargetType005: []string{"前往", "/uper/user_center/archive_list"},
			TargetType006: []string{"前往", "/uper/user_center/archive_list"},
			TargetType007: []string{"前往", "/uper/user_center/archive_list"},
			TargetType008: []string{"前往", "/uper/user_center/archive_list"},
			TargetType009: []string{"前往", "/uper/user_center/archive_list"},
			TargetType010: []string{"前往", "/uper/user_center/archive_list"},
			TargetType011: []string{"去设置", "https://member.bilibili.com/studio/gabriel/watermark"},
			TargetType012: []string{"去关注", ""},
			TargetType013: []string{"去投稿", "/uper/user_center/add_archive/"},
			TargetType014: []string{"去投稿", "/uper/user_center/add_archive/"},
			TargetType015: []string{"前往", "/uper/user_center/archive_list"},
			TargetType016: []string{"前往", "/uper/user_center/archive_list"},
			TargetType017: []string{"前往", "/uper/user_center/archive_list"},
			TargetType018: []string{"前往", "/uper/user_center/archive_list"},
			TargetType019: []string{"前往", "/uper/user_center/archive_list"},
			TargetType020: []string{"前往", "/uper/user_center/archive_list"},
			TargetType021: []string{"前往", "/uper/user_center/archive_list"},
			TargetType022: []string{"前往", "https://member.bilibili.com/studio/gabriel/fans-manage/overview"},
			TargetType023: []string{"去加入", "https://member.bilibili.com/studio/up-allowance-h5#/"},
			TargetType024: []string{"去开通", "https://member.bilibili.com/studio/gabriel/fans-manage/medal"},
		},
	}

	// H5RedirectMap task map for app
	H5RedirectMap = map[string]map[int8][]string{
		"android": {
			TargetType001: []string{"去投稿", "bilibili://uper/user_center/add_archive/"},
			TargetType002: []string{"去分享", "bilibili://uper/user_center/manuscript-list/"},
			TargetType003: []string{"前往", "https://member.bilibili.com/college?from=task"},
			TargetType004: []string{"前往", "bilibili://uper/user_center/manuscript-list/"},
			TargetType005: []string{"前往", "bilibili://uper/user_center/manuscript-list/"},
			TargetType006: []string{"前往", "bilibili://uper/user_center/manuscript-list/"},
			TargetType007: []string{"前往", "bilibili://uper/user_center/manuscript-list/"},
			TargetType008: []string{"前往", "bilibili://uper/user_center/manuscript-list/"},
			TargetType009: []string{"前往", "bilibili://uper/user_center/manuscript-list/"},
			TargetType010: []string{"前往", "bilibili://uper/user_center/manuscript-list/"},
			TargetType011: []string{"去设置", "https://member.bilibili.com/studio/gabriel/watermark"},
			TargetType012: []string{"去关注", "去关注"},
			TargetType013: []string{"去投稿", "bilibili://uper/user_center/add_archive/"},
			TargetType014: []string{"去投稿", "bilibili://uper/user_center/add_archive/"},
			TargetType015: []string{"前往", "bilibili://uper/user_center/manuscript-list/"},
			TargetType016: []string{"前往", "bilibili://uper/user_center/manuscript-list/"},
			TargetType017: []string{"前往", "bilibili://uper/user_center/manuscript-list/"},
			TargetType018: []string{"前往", "bilibili://uper/user_center/manuscript-list/"},
			TargetType019: []string{"前往", "bilibili://uper/user_center/manuscript-list/"},
			TargetType020: []string{"前往", "bilibili://uper/user_center/manuscript-list/"},
			TargetType021: []string{"前往", "bilibili://uper/user_center/manuscript-list/"},
			TargetType022: []string{"前往", "https://member.bilibili.com/studio/gabriel/fans-manage/overview"},
			TargetType023: []string{"去加入", "https://member.bilibili.com/studio/up-allowance-h5#/"},
			TargetType024: []string{"去开通", "https://member.bilibili.com/studio/gabriel/fans-manage/medal"},
		},
		"ios": {TargetType001: []string{"去投稿", "/uper/user_center/add_archive/"},
			TargetType002: []string{"去分享", "/uper/user_center/archive_list"},
			TargetType003: []string{"前往", "https://member.bilibili.com/college?from=task"},
			TargetType004: []string{"前往", "/uper/user_center/archive_list"},
			TargetType005: []string{"前往", "/uper/user_center/archive_list"},
			TargetType006: []string{"前往", "/uper/user_center/archive_list"},
			TargetType007: []string{"前往", "/uper/user_center/archive_list"},
			TargetType008: []string{"前往", "/uper/user_center/archive_list"},
			TargetType009: []string{"前往", "/uper/user_center/archive_list"},
			TargetType010: []string{"前往", "/uper/user_center/archive_list"},
			TargetType011: []string{"去设置", "https://member.bilibili.com/studio/gabriel/watermark"},
			TargetType012: []string{"去关注", "去关注"},
			TargetType013: []string{"去投稿", "/uper/user_center/add_archive/"},
			TargetType014: []string{"去投稿", "/uper/user_center/add_archive/"},
			TargetType015: []string{"前往", "/uper/user_center/archive_list"},
			TargetType016: []string{"前往", "/uper/user_center/archive_list"},
			TargetType017: []string{"前往", "/uper/user_center/archive_list"},
			TargetType018: []string{"前往", "/uper/user_center/archive_list"},
			TargetType019: []string{"前往", "/uper/user_center/archive_list"},
			TargetType020: []string{"前往", "/uper/user_center/archive_list"},
			TargetType021: []string{"前往", "/uper/user_center/archive_list"},
			TargetType022: []string{"前往", "https://member.bilibili.com/studio/gabriel/fans-manage/overview"},
			TargetType023: []string{"去加入", "https://member.bilibili.com/studio/up-allowance-h5#/"},
			TargetType024: []string{"去开通", "https://member.bilibili.com/studio/gabriel/fans-manage/medal"},
		},
	}

	// TaskGroupTipMap taskGroup tips for h5
	TaskGroupTipMap = map[int8]map[int64]string{
		RewardNotAvailable: {
			1: "快迈出你的第一步吧~~",
			2: "数据会在完成任务的第二天上午12:00进行核实哦。",
			3: "数据会在完成任务的第二天上午12:00进行核实哦。",
			4: "完成全部新手任务就可以解锁大礼包哦～",
			5: "数据会在完成任务的第二天上午12:00进行核实哦。",
			6: "数据会在完成任务的第二天上午12:00进行核实哦。",
			7: "数据会在完成任务的第二天上午12:00进行核实哦。",
			8: "完成全部任务就可以解锁大礼包哦～",
		},
		RewardAvailable: {
			1: "会员购优惠券领取后就即时生效了哦～",
			2: "B币券领取后就即时生效了哦～",
			3: "大会员代金券领取后就即时生效了哦～",
			4: "会员购优惠券领取后就即时生效了哦～",
			5: "会员购优惠券领取后就即时生效了哦～",
			6: "大会员代金券领取后就即时生效了哦～",
			7: "B币券领取后就即时生效了哦～",
			8: "双倍激励卡领取后需激活才可使用哦～",
		},
		RewardReceived: {
			1: "可以在我的奖品查看领奖记录哦～",
			2: "可以在我的奖品查看领奖记录哦～",
			3: "可以在我的奖品查看领奖记录哦～",
			4: "可以在我的奖品查看领奖记录哦～",
			5: "可以在我的奖品查看领奖记录哦～",
			6: "可以在我的奖品查看领奖记录哦～",
			7: "可以在我的奖品查看领奖记录哦～",
			8: "可以在我的奖品查看领奖记录哦～",
		},
		RewardUnlock: {
			1: "完成全部新手任务就可以解锁大礼包哦～",
			2: "完成全部新手任务就可以解锁大礼包哦～",
			3: "完成全部新手任务就可以解锁大礼包哦～",
			4: "完成全部新手任务就可以解锁大礼包哦～",
			5: "完成全部新手任务就可以解锁大礼包哦～",
			6: "完成全部新手任务就可以解锁大礼包哦～",
			7: "完成全部新手任务就可以解锁大礼包哦～",
			8: "完成全部新手任务就可以解锁大礼包哦～",
		},
	}

	// GiftTipMap gift tips for h5
	GiftTipMap = map[int8]map[int8]string{
		RewardNotAvailable: {
			1: "完成全部新手任务马上就能领头像挂件了呢～",
			2: "完成全部进阶任务马上就能领头像挂件了呢～",
		},
		RewardAvailable: {
			1: "头像挂件领取后即时生效哦～",
			2: "头像挂件领取后即时生效哦～",
		},
		RewardReceived: {
			1: "可以去我的奖品查看领奖记录哦～",
			2: "可以去我的奖品查看领奖记录哦～",
		},
		//RewardUnlock:{
		//	1:"",
		//	2:"再完成n个任务就能领取了呢",
		//},
	}
)

// Task for def task struct.
type Task struct {
	ID           int64     `json:"id"`
	GroupID      int64     `json:"-"`
	Type         int8      `json:"type"`
	State        int8      `json:"-"`
	Title        string    `json:"title"`
	Desc         string    `json:"desc"`
	Comment      string    `json:"-"`
	TargetType   int8      `json:"-"`
	TargetValue  int       `json:"-"`
	CompleteSate int8      `json:"complete_state"`
	Label        string    `json:"label,omitempty"`
	Redirect     string    `json:"redirect,omitempty"`
	Rank         int64     `json:"-"`
	Extra        string    `json:"extra"`
	FanRange     string    `json:"-"`
	UpTime       time.Time `json:"-"`
	DownTime     time.Time `json:"-"`
	Online       int8      `json:"-"`
	CTime        time.Time `json:"-"`
	MTime        time.Time `json:"-"`
}

// AppTasks for def task struct.
type AppTasks struct {
	ID       int64  `json:"id"`
	Type     int8   `json:"type"`
	Title    string `json:"title"`
	Label    string `json:"label"`
	Redirect string `json:"redirect"`
}

//TaskGroup for  newcomer & advanced tasks
type TaskGroup struct {
	Tasks      []*Task `json:"tasks"`
	GroupID    int64   `json:"group_id"`
	RewardID   []int64 `json:"reward_id"`
	Completed  int64   `json:"completed"`
	Incomplete int64   `json:"incomplete"`
}

// TaskList for def task list.
type TaskList struct {
	TaskGroups      []*TaskGroup `json:"task_groups"`
	TotalCompleted  int64        `json:"total_completed"`
	TotalIncomplete int64        `json:"total_incomplete"`
}

// Reward for def reward struct
type Reward struct {
	ID         int64     `json:"id"`
	ParentID   int64     `json:"parent_id"`
	Type       int8      `json:"type"`
	State      int8      `json:"state"`
	IsActive   int8      `json:"is_active"`
	PriceID    string    `json:"price_id"`
	PrizeUnit  int       `json:"prize_unit"`
	Expire     int       `json:"expire"`
	Name       string    `json:"name"`
	Logo       string    `json:"logo"`
	Comment    string    `json:"comment"`
	UnlockLogo string    `json:"unlock_logo"`
	NameExtra  string    `json:"name_extra"`
	CTime      time.Time `json:"-"`
	MTime      time.Time `json:"-"`
}

// TaskReward def to combine task and reward data structures
type TaskReward struct {
	Mid int64

	//task data
	TaskID           int64
	TaskGroupID      int64
	TaskTitle        string
	TaskDesc         string
	TaskType         int8
	TaskState        int8
	TaskCompleteSate int8
	Label            string
	Redirect         string

	//reward data
	RewardID       int64
	RewardParentID int64
	RewardName     string
	RewardLogo     string
	RewardType     int8
	RewardState    int8
	RewardPriceID  string
}

// TaskKind for newcomer & advanced & monthly task classification
type TaskKind struct {
	Type      int8  `json:"type"`
	State     int8  `json:"state"`
	Completed int64 `json:"completed"`
	Total     int64 `json:"total"`
}

//TaskRewardGroup for  newcomer & advanced tasks
type TaskRewardGroup struct {
	GroupID     int64     `json:"group_id"`
	Tasks       []*Task   `json:"tasks"`
	Rewards     []*Reward `json:"rewards"`
	RewardState int8      `json:"reward_state"` //  -1-不可领取 , 0-可领取 , 1-已领取
	Completed   int64     `json:"completed"`
	Total       int64     `json:"total"`
	TaskType    int8      `json:"task_type,omitempty"`
	Tip         string    `json:"tip,omitempty"`
}

// TaskGift for def struct
type TaskGift struct {
	State   int8      `json:"state"` //  -1-不可领取 ，0-可领取 , 1-已领取
	Type    int8      `json:"type,omitempty"`
	Rewards []*Reward `json:"rewards"`
	Tip     string    `json:"tip,omitempty"`
}

// TaskRewardList for def task list.
type TaskRewardList struct {
	TaskReceived int8               `json:"task_received"` // -1-未领取任务，0-已领取任务
	TaskType     int8               `json:"task_type"`
	TaskKinds    []*TaskKind        `json:"task_kinds"`
	TaskGroups   []*TaskRewardGroup `json:"task_groups"`
	TaskGift     []*TaskGift        `json:"task_gift"`
}

// RewardReceive for def reward receive records.
type RewardReceive struct {
	ID          int64     `json:"id"`
	MID         int64     `json:"mid"`
	TaskGiftID  int64     `json:"task_gift_id"`
	TaskGroupID int64     `json:"task_group_id"`
	RewardID    int64     `json:"reward_id"`
	RewardType  int8      `json:"reward_type"`
	State       int8      `json:"state"`
	ReceiveTime time.Time `json:"receive_time"`
	CTime       time.Time `json:"ctime"`
	MTime       time.Time `json:"mtime"`
	ExpireTime  time.Time `json:"expire_time"`
	RewardName  string    `json:"reward_name"`
}

// RewardReceiveGroup for reward receive group
type RewardReceiveGroup struct {
	Count          int              `json:"count"`
	RewardType     int8             `json:"reward_type"`
	RewardTypeName string           `json:"reward_type_name"`
	RewardTypeLogo string           `json:"reward_type_logo"`
	Comment        string           `json:"comment"`
	Items          []*RewardReceive `json:"items"`
}

// UserTask for def user task struct.
type UserTask struct {
	ID           int64     `json:"id"`
	MID          int64     `json:"mid"`
	TaskID       int64     `json:"task_id"`
	TaskGroupID  int64     `json:"task_group_id"`
	TaskType     int8      `json:"task_type"`
	State        int8      `json:"state"`
	TaskBindTime time.Time `json:"task_bind_time"`
	CTime        time.Time `json:"ctime"`
	MTime        time.Time `json:"mtime"`
}

// IndexNewcomer for index show
type IndexNewcomer struct {
	TaskReceived int8    `json:"task_received"`
	SubZero      bool    `json:"sub_zero"`
	NoReceive    int     `json:"no_receive"`
	Tasks        []*Task `json:"tasks"`
}

// AppIndexNewcomer for index show
type AppIndexNewcomer struct {
	TaskReceived int8        `json:"task_received"`
	H5URL        string      `json:"h5_url"`
	AppTasks     []*AppTasks `json:"tasks"`
}

// CheckTaskStateReq check task state req by creative-job grpc client.
type CheckTaskStateReq struct {
	MID    int64
	TaskID int64
}

// TaskGroupReward for def task-group-reward
type TaskGroupReward struct {
	ID          int64     `json:"id"`
	TaskGroupID int64     `json:"task_group_id"`
	RewardID    int64     `json:"reward_id"`
	State       int8      `json:"state"`
	Comment     string    `json:"comment"`
	CTime       time.Time `json:"ctime"`
	MTime       time.Time `json:"mtime"`
}

// GiftReward for gift reward
type GiftReward struct {
	ID       int64     `json:"id"`
	RootType int8      `json:"root_type"`
	TaskType int8      `json:"task_type"`
	RewardID int64     `json:"reward_id"`
	State    int8      `json:"state"`
	Comment  string    `json:"comment"`
	CTime    time.Time `json:"ctime"`
	MTime    time.Time `json:"mtime"`
}

//TaskMsg for newcomer task finish notify.
type TaskMsg struct {
	MID       int64 `json:"mid"`
	Count     int64 `json:"count"`
	From      int   `json:"from"`
	TimeStamp int64 `json:"timestamp"`
}

// H5TaskRewardList for def task list.
type H5TaskRewardList struct {
	TaskReceived int8               `json:"task_received"` // -1-未领取任务，0-已领取任务
	TaskGroups   []*TaskRewardGroup `json:"task_groups"`
	TaskGift     []*TaskGift        `json:"task_gifts"`
}

//PubTask for def struct
type PubTask struct {
	ID    int64  `json:"id"`
	Type  int8   `json:"type"`
	Title string `json:"title"`
	Desc  string `json:"desc"`
	State int8   `json:"state"`
}

//PubTaskList for def struct
type PubTaskList struct {
	TaskReceived int8       `json:"task_received"`
	Tasks        []*PubTask `json:"tasks"`
}

// TaskGroupEntity for def struct
type TaskGroupEntity struct {
	ID       int64     `json:"id"`
	Rank     int64     `json:"rank"`
	State    int8      `json:"state"`
	RootType int8      `json:"root_type"`
	Type     int8      `json:"type"`
	Online   int8      `json:"online"`
	CTime    time.Time `json:"ctime"`
	MTime    time.Time `json:"mtime"`
}

// TaskRewardEntity for def struct
type TaskRewardEntity struct {
	ID       int64     `json:"id"`
	TaskID   int64     `json:"task_id"`
	RewardID int64     `json:"reward_id"`
	State    int8      `json:"state"`
	Type     int8      `json:"type"`
	Comment  string    `json:"comment"`
	CTime    time.Time `json:"ctime"`
	MTime    time.Time `json:"mtime"`
}

// RewardReceive2 for def reward receive records.
type RewardReceive2 struct {
	ID          int64     `json:"id"`
	MID         int64     `json:"mid"`
	OID         int64     `json:"oid"`
	Type        int8      `json:"type"`
	RewardID    int64     `json:"reward_id"`
	RewardType  int8      `json:"reward_type"`
	State       int8      `json:"state"`
	ReceiveTime time.Time `json:"receive_time"`
	CTime       time.Time `json:"ctime"`
	MTime       time.Time `json:"mtime"`
	RewardName  string    `json:"reward_name"`
}
