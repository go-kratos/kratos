package jpush

// Option .
type Option struct {
	SendNo             int   `json:"sendno,omitempty"`
	TimeLive           int   `json:"time_to_live,omitempty"`
	ApnsProduction     bool  `json:"apns_production"`
	OverrideMsgID      int64 `json:"override_msg_id,omitempty"`
	BigPushDuration    int   `json:"big_push_duration,omitempty"`
	ReturnInvalidToken bool  `json:"return_invalid_rid,omitempty"` // 是否同步返回无效的token
}

// SetSendno .
func (o *Option) SetSendno(no int) {
	o.SendNo = no
}

// SetTimelive .
func (o *Option) SetTimelive(timelive int) {
	o.TimeLive = timelive
}

// SetOverrideMsgID .
func (o *Option) SetOverrideMsgID(id int64) {
	o.OverrideMsgID = id
}

// SetApns .
func (o *Option) SetApns(apns bool) {
	o.ApnsProduction = apns
}

// SetBigPushDuration .
func (o *Option) SetBigPushDuration(dur int) {
	o.BigPushDuration = dur
}

// SetReturnInvalidToken .
func (o *Option) SetReturnInvalidToken(onoff bool) {
	o.ReturnInvalidToken = onoff
}
