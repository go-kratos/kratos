package model

// VipUserInfoMsg vip binlog msg.
type VipUserInfoMsg struct {
	ID                   int64  `json:"id"`
	Mid                  int64  `json:"mid"`
	Ver                  int64  `json:"ver"`
	VipType              int8   `json:"vip_type"`
	VipPayType           int8   `json:"vip_pay_type"`
	PayChannelID         int64  `json:"pay_channel_id"`
	VipStatus            int32  `json:"vip_status"`
	VipStartTime         string `json:"vip_start_time"`
	VipRecentTime        string `json:"vip_recent_time"`
	VipOverdueTime       string `json:"vip_overdue_time"`
	AnnualVipOverdueTime string `json:"annual_vip_overdue_time"`
	IosOverdueTime       string `json:"ios_overdue_time"`
}

// VipReq vip params.
type VipReq struct {
	Mid            int64
	VipType        int8
	VipStatus      int32
	VipOverdueTime int64
}
