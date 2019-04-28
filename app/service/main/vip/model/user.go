package model

// VipInfoBoResp vipinfo bo resp.
type VipInfoBoResp struct {
	Mid                  int64 `json:"mid"`
	VipType              int32 `json:"vip_type"`
	PayType              int32 `json:"pay_type"`
	PayChannelID         int32 `json:"pay_channel_id"`
	VipStatus            int32 `json:"vip_status"`
	VipStartTime         int64 `json:"vip_start_time"`
	VipOverdueTime       int64 `json:"vip_overdue_time"`
	AnnualVipOverdueTime int64 `json:"annual_vip_overdue_time"`
	VipRecentTime        int64 `json:"vip_recent_time"`
	AutoRenewed          int32 `json:"auto_renewed"`
	IosOverdueTime       int64 `json:"ios_overdue_time"`
}
