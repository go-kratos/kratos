package model

// Pool .
type Pool struct {
	Id          int64  `json:"id"`
	CoinId      int64  `json:"coin_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	StartTime   int64  `json:"start_time"`
	EndTime     int64  `json:"end_time"`
	Status      int64  `json:"status"`
	IsBottom    int64  `json:"is_bottom"`
}

// Coin .
type Coin struct {
	Id        int64  `json:"id"`
	Title     string `json:"title"`
	GiftType  int64  `json:"gift_type"`
	ChangeNum int64  `json:"change_num"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
	Status    int64  `json:"status"`
}

// UserInfo .
type UserInfo struct {
	Uid           int64 `json:"uid"`
	NormalScore   int64 `json:"normal_score"`
	ColorfulScore int64 `json:"colorful_score"`
}

// GiftMsg .
type GiftMsg struct {
	// by notify
	MsgContent string `form:"msg_content"`
}

// CoinConfig .
type CoinConfig struct {
	CoinId         int64 `json:"coin_id"`
	Type           int64 `json:"type"`
	AreaV2ParentId int64 `json:"area_v2_parent_id"`
	AreaV2Id       int64 `json:"area_v2_id"`
	GiftId         int64 `json:"gift_id"`
	IsAll          int64 `json:"is_all"`
}

// PoolPrize .
type PoolPrize struct {
	Id          int64  `json:"id"`
	PoolId      int64  `json:"pool_id"`
	Type        int64  `json:"type"`
	Num         int64  `json:"num"`
	ObjectId    int64  `json:"object_id"`
	Expire      int64  `json:"expire"`
	WebUrl      string `json:"web_url"`
	MobileUrl   string `json:"mobile_url"`
	Description string `json:"description"`
	JumpUrl     string `json:"jump_url"`
	ProType     int64  `json:"pro_type"`
	Chance      int64  `json:"chance"`
	LoopNum     int64  `json:"loop_num"`
	LimitNum    int64  `json:"limit_num"`
	Weight      int64  `json:"weight"`
}

// AddCapsule AddCapsule
type AddCapsule struct {
	Uid    int64  `json:"uid"`
	Type   string `json:"type"`
	CoinId int64  `json:"coin_id"`
	Num    int64  `json:"num"`
	Source string `json:"source"`
	MsgId  string `json:"msg_id"`
}

// ExtraData .
type ExtraData struct {
	Id        int64  `json:"id"`
	Uid       int64  `json:"uid"`
	Type      string `json:"type"`
	ItemValue int64  `json:"item_value"`
	ItemExtra string `json:"item_extra"`
}
