package jpush

// Option .
type Option struct {
	SendNo          int   `json:"sendno,omitempty"`
	TimeLive        int   `json:"time_to_live,omitempty"`
	ApnsProduction  bool  `json:"apns_production"`
	OverrideMsgID   int64 `json:"override_msg_id,omitempty"`
	BigPushDuration int   `json:"big_push_duration,omitempty"`
}
