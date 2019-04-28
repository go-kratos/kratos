package dao

import (
	"context"
	"go-common/app/service/live/captcha/api/liverpc/v0"
	"go-common/app/service/live/captcha/api/liverpc/v1"
	"go-common/library/log"
)

// LiveCreate call liveRpc for create captcha
func (d *Dao) LiveCreate(ctx context.Context, width int64, height int64) (resp *v1.CaptchaCreateResp_Data, err error) {
	resp = &v1.CaptchaCreateResp_Data{}
	rpcReq := &v1.CaptchaCreateReq{
		Width:  width,
		Height: height,
	}
	rpcResp, err := d.liveCaptcha.V1Captcha.Create(ctx, rpcReq)
	if err != nil {
		log.Error("[XCaptcha][Create][call liveCaptcha] call error, error:%s", err.Error())
		return
	}
	if rpcResp.Data == nil {
		log.Error("[XCaptcha][create][call liveCaptcha] return error, response.data is nil")
		return
	}
	code := rpcResp.Code
	msg := rpcResp.Msg
	if code != 0 {
		log.Error("[XCaptcha][Create][call liveCaptcha] create error, code:%s, msg:%s, data:%s", code, msg, resp)
		return
	}
	resp = rpcResp.Data
	return
}

// LiveCheck liveCaptcha check
func (d *Dao) LiveCheck(ctx context.Context, token string, phrase string) (resp int64, err error) {
	resp = -400
	rpcReq := &v0.CaptchaCheckReq{
		Token:  token,
		Phrase: phrase,
	}
	rpcResp, err := d.liveCaptcha.V0Captcha.Check(ctx, rpcReq)
	if err != nil {
		log.Error("[XCaptcha][verify][call liveCaptcha] call error, error:%s", err.Error())
		return
	}
	if rpcResp == nil {
		log.Error("[XCaptcha][verify][call liveCaptcha] return error, response is nil")
		return
	}
	resp = rpcResp.Code
	return
}
