package dao

import (
	"context"
	"go-common/library/log"
	"strconv"
)

// PubMessage pub to databus struct
type PubMessage struct {
	RoomId  int64  `json:"room_id"`
	Uid     int64  `json:"uid"`
	Ip      string `json:"ip"`
	Action  int    `json:"action"`   // 1create 2verify
	ReqType int64  `json:"req_type"` // request captcha type 0image 1geetest
	ResType int64  `json:"res_type"` // response captcha type 0image 1geetest
	ResCode int64  `json:"res_code"` // 0success 1 failed
}

// Pub databus publish
func (d *Dao) Pub(ctx context.Context, message PubMessage) (err error) {
	if err = d.captchaAnti.Send(ctx, strconv.FormatInt(message.Uid, 10), message); err != nil {
		log.Error("[XCaptcha][DataBus] call for publish error, err:%v, msg:%v", err, message)
	}
	return
}
