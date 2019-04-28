package v1

import (
	v1pb "go-common/app/interface/live/web-room/api/http/v1"
	"go-common/app/interface/live/web-room/conf"
	xCaptcha "go-common/app/service/live/xcaptcha/api/grpc/v1"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"

	"context"
)

// CaptchaService struct
type CaptchaService struct {
	conf     *conf.Config
	xCaptcha *xCaptcha.Client
	// optionally add other properties here, such as dao
	// dao *dao.Dao
}

//NewCaptchaService init
func NewCaptchaService(c *conf.Config) (s *CaptchaService) {
	s = &CaptchaService{
		conf: c,
	}
	client, err := xCaptcha.NewClient(c.XCaptcha)
	if err != nil {
		log.Error("[web-room][captcha][new XCaptcha error] init XCaptcha error, %+v", err)
	}
	s.xCaptcha = client
	return s
}

// captcha 相关服务

// Create implementation
// 创建验证码
func (s *CaptchaService) Create(ctx context.Context, req *v1pb.CreateCaptchaReq) (resp *v1pb.CreateCaptchaResp, err error) {
	resp = &v1pb.CreateCaptchaResp{}
	uid, _ := metadata.Value(ctx, "mid").(int64)
	if uid <= 0 {
		err = ecode.Error(ecode.NoLogin, "未登录")
		return
	}
	XCaptchaReq := &xCaptcha.XCreateCaptchaReq{
		Type:       req.GetType(),
		ClientType: req.GetClientType(),
		Height:     req.GetHeight(),
		Width:      req.GetWidth(),
		Uid:        uid,
		ClientIp:   metadata.String(ctx, metadata.RemoteIP),
	}
	if s.xCaptcha == nil {
		err = ecode.Error(ecode.ServerErr, "创建验证码失败")
		return
	}
	XCaptchaResp, err := s.xCaptcha.Create(ctx, XCaptchaReq)
	if err != nil || XCaptchaResp == nil {
		return
	}

	resp = &v1pb.CreateCaptchaResp{
		Type:    XCaptchaResp.Type,
		Geetest: &v1pb.GeeTest{},
		Image:   &v1pb.Image{},
	}
	resp.Geetest.Challenge = XCaptchaResp.Geetest.Challenge
	resp.Geetest.Gt = XCaptchaResp.Geetest.Gt
	resp.Image.Token = XCaptchaResp.Image.Token
	resp.Image.Content = XCaptchaResp.Image.Content
	resp.Image.Tips = XCaptchaResp.Image.Tips
	return
}

// Verify implementation
// 校验接口 `midware:"auth" method:"POST"`
func (s *CaptchaService) Verify(ctx context.Context, req *v1pb.VerifyReq) (resp *v1pb.VerifyResp, err error) {
	resp = &v1pb.VerifyResp{}
	uid, _ := metadata.Value(ctx, "mid").(int64)
	if uid <= 0 {
		err = ecode.Error(ecode.NoLogin, "未登录")
		return
	}
	if req.GetAnti() == "" {
		err = ecode.Error(ecode.ParamInvalid, "参数错误")
		return
	}
	XCaptchaReq := &xCaptcha.CheckReq{
		Anti:     req.GetAnti(),
		Uid:      uid,
		ClientIp: metadata.String(ctx, metadata.RemoteIP),
	}
	if s.xCaptcha == nil {
		err = ecode.Error(ecode.ServerErr, "创建验证码失败")
		return
	}
	XCaptchaResp, err := s.xCaptcha.Check(ctx, XCaptchaReq)
	if err != nil || XCaptchaResp == nil {
		return
	}
	resp.Type = XCaptchaResp.Type
	resp.Token = XCaptchaResp.Token
	return
}
