package v1

import (
	"encoding/json"

	v1pb "go-common/app/interface/live/app-room/api/http/v1"
	risk "go-common/app/service/live/live_riskcontrol/api/grpc/v1"
	xcaptcha "go-common/app/service/live/xcaptcha/api/grpc/v1"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

type rbody struct {
	Roomid   int64  `json:"roomid"`
	Msg      string `json:"msg" `
	Rnd      string `json:"rnd" `
	Fontsize int64  `json:"fontsize"`
	Mode     int64  `json:"mode" `
	Color    int64  `json:"color"`
	Bubble   int64  `json:"bubble"`
}

func isriskcontrol(ctx *bm.Context, uid int64, r *v1pb.SendDMReq) (forbid bool, err *ecode.Status) {
	req := &risk.GetForbiddenReq{
		Uid:    uid,
		Uri:    "/xlive/app-room/v1/dM/sendmsg",
		Ip:     metadata.String(ctx, metadata.RemoteIP),
		Method: "POST",
		Header: make(map[string]string),
	}
	for k := range ctx.Request.Header {
		req.Header[k] = ctx.Request.Header.Get(k)
	}
	rb := &rbody{
		Roomid:   r.Roomid,
		Msg:      r.Msg,
		Rnd:      r.Rnd,
		Fontsize: r.Fontsize,
		Mode:     r.Mode,
		Color:    r.Color,
		Bubble:   r.Bubble,
	}
	jb, _ := json.Marshal(rb)
	req.Body = string(jb)

	resp, rerr := riskClient.GetForbidden(ctx, req)
	if rerr != nil {
		log.Error("DM: riskcontrol err:%+v", rerr)
		return false, nil
	}

	switch resp.IsForbidden {
	case 0:
		return false, nil
	case 1:
		return true, ecode.Error(400, "访问被拒绝")
	case 2:
		return true, ecode.Error(1990000, "need a second time verify")
	}
	return false, nil
}

func checkVerify(ctx *bm.Context, anti string, uid int64, roomid int64) bool {
	req := &xcaptcha.XVerifyReq{
		Uid:      uid,
		ClientIp: metadata.String(ctx, metadata.RemoteIP),
		XAnti:    anti,
		RoomId:   roomid,
	}

	if _, err := xcaptchaClient.Verify(ctx, req); err != nil {
		return false
	}

	return true
}
