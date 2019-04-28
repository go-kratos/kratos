package model

// DPTask data platform task
type DPTask struct {
	Task
	DPParams
}

// DPParams data platform params
type DPParams struct {
	Age             int              `form:"age" json:"age"`
	Sex             int              `form:"sex" json:"sex"`
	IsUp            int              `form:"is_up" json:"is_up"`
	IsFormalMember  int              `form:"is_formal_member" json:"is_formal_member"`
	UserActiveDay   int              `form:"user_active_day" json:"user_active_day"`
	UserNewDay      int              `form:"user_new_day" json:"user_new_day"`
	UserSilentDay   int              `form:"user_silent_day" json:"user_silent_day"`
	Area            []int            `form:"area,split" json:"-"`
	AreaStr         string           `json:"area"`
	Level           []int            `form:"level,split" json:"-"`
	LevelStr        string           `json:"level"`
	Platforms       []int            `form:"platforms,split" json:"-"`
	PlatformStr     string           `json:"platforms"`
	Like            []int            `form:"like,split" json:"-"`
	LikeStr         string           `json:"like"`
	Channel         []string         `form:"channel,split" json:"-"`
	ChannelStr      string           `json:"channel"`
	VipExpireStr    string           `form:"vip_expire" json:"-"`
	VipExpires      []*VipExpire     `json:"vip_expire"`
	AttentionStr    string           `form:"self_attention" json:"-"`
	Attentions      []*SelfAttention `json:"self_attention"`
	AttentionsType  int              `form:"self_attention_type" json:"self_attention_type"`
	ActivePeriodStr string           `form:"active" json:"-"`
	ActivePeriods   []*ActivePeriod  `json:"active"`
	ActivePeriod    int
}

// SelfAttention 自选关注
type SelfAttention struct {
	Type    int    `json:"type"`
	Include string `json:"include"`
	Exclude string `json:"exclude"`
}

// VipExpire 大会员过期时间
type VipExpire struct {
	Begin string `json:"begin"`
	End   string `json:"end"`
}

// ActivePeriod 活跃时间段
type ActivePeriod struct {
	Period     int    `json:"period"`
	PushTime   string `json:"push_time"`
	ExpireTime string `json:"expire_time"`
}

// DPCondition data platform condition
type DPCondition struct {
	ID        int64
	Task      int64
	Job       string
	Type      int
	Condition string
	SQL       string
	Status    int
	StatusURL string
	File      string
}

// TableName .
func (c *DPCondition) TableName() string {
	return "push_dataplatform_conditions"
}
