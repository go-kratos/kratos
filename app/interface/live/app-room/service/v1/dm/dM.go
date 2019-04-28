package v1

import (
	v1pb "go-common/app/interface/live/app-room/api/http/v1"
	"go-common/app/interface/live/app-room/conf"
	dmrpc "go-common/app/service/live/live-dm/api/grpc/v1"
	risk "go-common/app/service/live/live_riskcontrol/api/grpc/v1"
	xcaptcha "go-common/app/service/live/xcaptcha/api/grpc/v1"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"strings"
)

// DMService struct
type DMService struct {
	conf *conf.Config
}

var (
	//弹幕client
	dmClient dmrpc.DMClient
	//风控client
	riskClient risk.IsForbiddenClient
	//验证码client
	xcaptchaClient xcaptcha.XCaptchaClient
)

//NewDMService init
func NewDMService(c *conf.Config) (s *DMService) {
	s = &DMService{
		conf: c,
	}

	var err error
	if dmClient, err = dmrpc.NewClient(c.DM); err != nil {
		panic(err)
	}
	if riskClient, err = risk.NewClient(c.Risk); err != nil {
		panic(err)
	}
	if xcaptchaClient, err = xcaptcha.NewClient(c.VerifyConf); err != nil {
		panic(err)
	}
	return s
}

func sendMsg(ctx *bm.Context, req *v1pb.SendDMReq, uid int64) (resp *v1pb.SendMsgResp, err error) {
	resp = &v1pb.SendMsgResp{}
	ip := metadata.String(ctx, metadata.RemoteIP)

	var ck = make([]string, 0, 10)
	for _, v := range ctx.Request.Cookies() {
		if v.Name != "SESSDATA" {
			ck = append(ck, v.Name+"="+v.Value)
		}
	}

	var dmReq = &dmrpc.SendMsgReq{
		Uid:      uid,
		Roomid:   req.Roomid,
		Msg:      req.Msg,
		Rnd:      req.Rnd,
		Ip:       ip,
		Fontsize: req.Fontsize,
		Mode:     req.Mode,
		Platform: req.Platform,
		Msgtype:  0,
		Bubble:   req.Bubble,
		Lancer: &dmrpc.Lancer{
			Build:     req.Build,
			Buvid:     ctx.Request.Header.Get("Buvid"),
			UserAgent: ctx.Request.Header.Get("User-Agent"),
			Refer:     ctx.Request.Header.Get("Referer"),
			Cookie:    strings.Join(ck, ";"),
		},
	}
	gresp, gerr := dmClient.SendMsg(ctx, dmReq)
	if gerr != nil {
		log.Error("DM GRPC ERR: %v", gerr)
		err = ecode.Error(1003218, "系统正在维护中,请稍后尝试")
		return nil, err
	}
	if gresp.IsLimit {
		err = ecode.Error(ecode.Code(gresp.Code), gresp.LimitMsg)
		return nil, err
	}

	return resp, nil
}

//验证码风控
func verifyRisk(ctx *bm.Context, uid int64, req *v1pb.SendDMReq) (resp *v1pb.SendMsgResp, err error) {
	//验证码
	if req.Anti != "" {
		result := checkVerify(ctx, req.Anti, uid, req.Roomid)
		if !result {
			return nil, ecode.Error(1990001, "验证码验证失败")
		}
		return sendMsg(ctx, req, uid)
	}

	//风控校验
	ifb, ferr := isriskcontrol(ctx, uid, req)
	if ifb {
		return nil, ferr
	}
	return sendMsg(ctx, req, uid)
}

// SendMsg implementation
// `method:"POST"`
func (s *DMService) SendMsg(ctx *bm.Context, req *v1pb.SendDMReq) (resp *v1pb.SendMsgResp, err error) {
	//获取UID
	uid, ok := ctx.Get("mid")
	if !ok {
		err = ecode.Error(1003218, "未登录")
		return nil, err
	}
	uid64, ok := uid.(int64)
	if !ok {
		log.Error("DM: mid error")
		err = ecode.Error(1003218, "未登录")
		return nil, err
	}
	device, ok := ctx.Get("device")
	if !ok {
		log.Error("DM: Get device error")
		return sendMsg(ctx, req, uid64)
	}
	devices, ok := device.(*bm.Device)
	if !ok {
		log.Error("DM: device error")
		return sendMsg(ctx, req, uid64)
	}

	//验证码版本控制
	if (devices.RawMobiApp == "android" && devices.Build >= 5360000) ||
		(devices.RawMobiApp == "iphone" && devices.Build >= 8290) {
		return verifyRisk(ctx, uid64, req)
	}
	//发送弹幕
	return sendMsg(ctx, req, uid64)
}

// GetHistory implementation
// `method:"POST"`
func (s *DMService) GetHistory(ctx *bm.Context, req *v1pb.HistoryReq) (resp *v1pb.HistoryResp, err error) {
	resp = &v1pb.HistoryResp{}
	var hreq = &dmrpc.HistoryReq{
		Roomid: req.Roomid,
	}
	gresp, err := dmClient.GetHistory(ctx, hreq)
	if err != nil {
		log.Error("DM GRPC ERR: %v", err)
		return
	}
	resp.Admin = gresp.Admin
	resp.Room = gresp.Room
	return
}
