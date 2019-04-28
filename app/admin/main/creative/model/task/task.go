package task

const (
	//StateDel for normal state
	StateDel = int8(-1)
	//StateNormal for normal state
	StateNormal = int8(0)
	//StateHide for hide state
	StateHide = int8(1)

	//TaskManagement 任务管理
	TaskManagement = uint8(1)
	//AchievementManagement 成就管理
	AchievementManagement = uint8(2)

	//LogClientTask 日志服务类型
	LogClientTask = 301
)

const (
	_ uint8 = iota
	//TaskTypeNewcomer 新手任务
	TaskTypeNewcomer
	//TaskTypeAdvanced 进阶任务
	TaskTypeAdvanced
	//TaskTypeMonthly 月常任务
	TaskTypeMonthly
)

var (
	//TaskRootNameMap 管理分类, 1-任务管理、2-成就管理
	TaskRootNameMap = map[uint8]string{
		TaskManagement:        "任务管理",
		AchievementManagement: "成就管理",
	}

	//TaskGroupNameMap 任务分类, 1-新手任务、2-进阶任务、3-月常任务; 4-互动成就、5-投稿成就、6-行为成就、7-高级成就
	TaskGroupNameMap = map[uint8]string{
		TaskTypeNewcomer: "新手任务",
		TaskTypeAdvanced: "进阶任务",
		TaskTypeMonthly:  "月常任务",
		4:                "互动成就",
		5:                "投稿成就",
		6:                "行为成就",
		7:                "高级成就",
	}
)

//CheckRootType check task root type.
func CheckRootType(ty uint8) bool {
	if ty == TaskTypeNewcomer || ty == TaskTypeAdvanced || ty == TaskTypeMonthly {
		return true
	}
	return false
}

var (
	//TargetMap for target show
	TargetMap = map[int8]string{
		1:  "开放浏览的稿件",
		2:  "分享自己视频的次数",
		3:  "创作学院的观看记录",
		4:  "所有avid的获得评论数",
		5:  "所有avid获得分享数",
		6:  "所有avid的获得收藏数",
		7:  "所有avid的获得硬币数",
		8:  "所有avid获得点赞数",
		9:  "所有avid的获得弹幕数",
		10: "粉丝数",
		11: "水印开关为打开状态",
		12: "关注列表含有“哔哩哔哩创作中心”",
		13: "用手机投稿上传视频",
		14: "开放浏览的稿件",
		15: "任意avid的获得点击量",
		16: "任意avid的评论",
		17: "任意avid的获得分享数",
		18: "任意avid的获得收藏数",
		19: "任意avid的获得硬币数",
		20: "任意avid的获得点赞数",
		21: "任意avid的获得弹幕数",
		22: "激励计划状态为已开通",
		23: "粉丝勋章为开启状态",
	}
)

//TableName get table name
func (tg *TaskGroup) TableName() string {
	return "newcomers_task_group"
}

//TableName get table name
func (tgr *TaskGroupReward) TableName() string {
	return "newcomers_grouptask_reward"
}

//TaskGroup for task group.
type TaskGroup struct {
	ID        int64           `gorm:"column:id" form:"id" json:"id"`
	Rank      int64           `gorm:"column:rank" form:"rank" json:"rank"`
	State     int8            `gorm:"column:state" form:"state" json:"state"` //-1-删除, 0-正常, 1-隐藏
	RootType  uint8           `gorm:"column:root_type"  form:"root_type" json:"root_type"`
	Type      int8            `gorm:"column:type"  form:"type" json:"type"`
	CTime     string          `gorm:"column:ctime" form:"ctime" json:"-"`
	MTime     string          `gorm:"column:mtime" form:"mtime" json:"-"`
	RewardIDs string          `gorm:"-" form:"reward_ids" json:"-"`
	Comment   string          `gorm:"-" form:"comment" json:"comment"`
	Tasks     []*Task         `json:"tasks"`
	Reward    []*RewardResult `json:"reward"`
}

//TaskGroupReward for task group relation reward.
type TaskGroupReward struct {
	ID          int64  `gorm:"column:id" form:"id" json:"id"`
	TaskGroupID int64  `gorm:"column:task_group_id"  form:"task_group_id" json:"task_group_id"`
	RewardID    int64  `gorm:"column:reward_id" form:"reward_id" json:"reward_id"`
	State       int8   `gorm:"column:state"  form:"state" json:"state"`
	CTime       string `gorm:"column:ctime" form:"ctime" json:"-"`
	MTime       string `gorm:"column:mtime" form:"mtime" json:"-"`
}

//OrderTask for task or task group order.
type OrderTask struct {
	ID         int64 `form:"id"  validate:"required"`
	Rank       int64 `form:"rank" validate:"required"`
	SwitchID   int64 `form:"switch_id"  validate:"required"`
	SwitchRank int64 `form:"switch_rank" validate:"required"`
}

//RewardResult for task group relation reward result.
type RewardResult struct {
	RewardID   int64  `json:"reward_id"`
	RewardName string `json:"reward_name"`
}

//TableName get table name
func (t *Task) TableName() string {
	return "newcomers_task"
}

//TableName get table name
func (tr *TaskReward) TableName() string {
	return "newcomers_task_reward"
}

// Task for def task struct.
type Task struct {
	ID          int64           `gorm:"column:id" form:"id" json:"id"`
	GroupID     int64           `gorm:"column:group_id" form:"group_id" json:"group_id"`
	Title       string          `gorm:"column:title" form:"title" json:"title"`
	Desc        string          `gorm:"column:desc" form:"desc" json:"desc"`
	Comment     string          `gorm:"column:comment" form:"comment" json:"comment"`
	Type        int8            `gorm:"column:type" form:"type" json:"type"`
	State       int8            `gorm:"column:state" form:"state" json:"state"`
	TargetType  int8            `gorm:"column:target_type" form:"target_type" json:"target_type"`
	TargetValue int32           `gorm:"column:target_value" form:"target_value" json:"target_value"`
	Rank        int64           `gorm:"column:rank" form:"rank" json:"rank"`
	Extra       string          `gorm:"column:extra" form:"extra" json:"extra"`             //跳转链接等附加信息,json格式
	FanRange    string          `gorm:"column:fan_range" form:"fan_range" json:"fan_range"` //粉丝范围, json格式
	UpTime      string          `gorm:"column:up_time" form:"up_time" json:"up_time"`       //月常活动任务-上线时间
	DownTime    string          `gorm:"column:down_time" form:"down_time" json:"down_time"` //月常活动任务-下线时间
	CTime       string          `gorm:"column:ctime" form:"ctime" json:"-"`
	MTime       string          `gorm:"column:mtime" form:"mtime" json:"-"`
	RewardIDs   string          `gorm:"-" form:"reward_ids" json:"-"`
	Reward      []*RewardResult `json:"reward"`
}

//TaskReward for task relation reward.
type TaskReward struct {
	ID       int64  `gorm:"column:id"    form:"id" json:"id"`
	TaskID   int64  `gorm:"column:task_id"  form:"task_id" json:"task_id"`
	RewardID int64  `gorm:"column:reward_id" form:"reward_id" json:"reward_id"`
	State    int8   `gorm:"column:state"  form:"state" json:"state"`
	Comment  string `gorm:"column:comment" form:"comment" json:"comment"`
	CTime    string `gorm:"column:ctime" form:"ctime" json:"-"`
	MTime    string `gorm:"column:mtime" form:"mtime" json:"-"`
}

//TableName get table name
func (r *Reward) TableName() string {
	return "newcomers_reward"
}

//Reward for task reward.
type Reward struct {
	ID         int64     `gorm:"column:id"    form:"id" json:"id"`
	ParentID   int64     `gorm:"column:parent_id"  form:"parent_id" json:"parent_id"`
	Type       int8      `gorm:"column:type"  form:"type" json:"type"`
	State      int8      `gorm:"column:state" form:"state" json:"state"`
	IsActive   int8      `gorm:"column:is_active" form:"is_active" json:"is_active"`
	Name       string    `gorm:"column:name" form:"name" json:"name"`
	Logo       string    `gorm:"column:logo" form:"logo" json:"logo"`
	Comment    string    `gorm:"column:comment" form:"comment" json:"comment"`
	UnlockLogo string    `gorm:"column:unlock_logo" form:"unlock_logo" json:"unlock_logo"` //奖励未解锁, logo url
	NameExtra  string    `gorm:"column:name_extra" form:"name_extra" json:"name_extra"`    //支持奖励名称展示,json格式
	PrizeID    string    `gorm:"column:prize_id" form:"prize_id" json:"prize_id"`          //业务方奖品id
	PrizeUnit  int8      `gorm:"column:prize_unit" form:"prize_unit" json:"prize_unit"`    //奖品单位
	Expire     int16     `gorm:"column:expire"  form:"expire" json:"expire"`               //有效期 单位天
	CTime      string    `gorm:"column:ctime" form:"ctime" json:"-"`
	MTime      string    `gorm:"column:mtime" form:"mtime" json:"-"`
	Children   []*Reward `json:"children,omitempty"`
}

//TableName get table name
func (gf *GiftReward) TableName() string {
	return "newcomers_gift_reward"
}

//GiftReward for task gift reward.
type GiftReward struct {
	ID        int64           `gorm:"column:id"    form:"id" json:"id"`
	RootType  uint8           `gorm:"column:root_type"  form:"root_type" json:"root_type"`
	TaskType  int64           `gorm:"column:task_type"  form:"task_type" json:"task_type"`
	RewardID  int64           `gorm:"column:reward_id"  form:"reward_id" json:"-"`
	State     int8            `gorm:"column:state"  form:"state" json:"state"`
	Comment   string          `gorm:"column:comment" form:"comment" json:"comment"`
	CTime     string          `gorm:"column:ctime" form:"ctime" json:"-"`
	MTime     string          `gorm:"column:mtime" form:"mtime" json:"-"`
	RewardIDs string          `gorm:"-" form:"reward_ids" json:"-"`
	Reward    []*RewardResult `json:"reward"`
}

// LogParam for manager.
type LogParam struct {
	UID     int64       `json:"uid"`
	UName   string      `json:"uname"`
	Action  string      `json:"action"`
	OID     int64       `json:"oid"`
	OIDs    string      `json:"oids"`
	OName   string      `json:"oname"`
	OState  int8        `json:"ostate"`
	Content interface{} `json:"content"`
}
