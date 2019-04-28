package v1

import (
	v1pb "go-common/app/interface/live/web-room/api/http/v1"
	"go-common/app/interface/live/web-room/conf"
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
	//dmClient 弹幕clietn
	dmClient dmrpc.DMClient
	//riskClient 风控Client
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

	if riskClient, err = risk.NewClient(c.DM); err != nil {
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
		Platform: "web",
		Msgtype:  0,
		Bubble:   req.Bubble,
		Lancer: &dmrpc.Lancer{
			Build:     0,
			Buvid:     "",
			Refer:     ctx.Request.Header.Get("Referer"),
			UserAgent: ctx.Request.Header.Get("User-Agent"),
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

// SendMsg implementation
// `method:"POST"`
func (s *DMService) SendMsg(ctx *bm.Context, req *v1pb.SendDMReq) (resp *v1pb.SendMsgResp, err error) {
	uid, ok := ctx.Get("mid")
	if !ok {
		err = ecode.Error(1003218, "未登录")
		return nil, err
	}

	//验证码
	if req.Anti != "" {
		result := checkVerify(ctx, req.Anti, uid.(int64), req.Roomid)
		if !result {
			return nil, ecode.Error(1990001, "验证码验证失败")
		}
		return sendMsg(ctx, req, uid.(int64))
	}

	//风控校验
	ifb, ferr := isriskcontrol(ctx, uid.(int64), req)
	if ifb {
		return nil, ferr
	}
	//发送弹幕
	return sendMsg(ctx, req, uid.(int64))
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
