package http

import (
	grpc "go-common/app/service/bbq/push/api/grpc/v1"
	"go-common/app/service/bbq/push/api/http/v1"
	"go-common/app/service/bbq/push/conf"
	"go-common/app/service/bbq/push/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"

	"github.com/json-iterator/go"

	"github.com/pkg/errors"
)

var (
	srv *service.Service
	vfy *verify.Verify
)

// Init init
func Init(c *conf.Config, svc *service.Service) {
	srv = svc
	vfy = verify.New(c.Verify)
	engine := bm.DefaultServer(c.BM)
	route(engine)
	if err := engine.Start(); err != nil {
		log.Error("bm Start error(%v)", err)
		panic(err)
	}
}

func route(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)
	g := e.Group("/bbq/internal/push")
	{
		g.POST("/notice", notice)
		g.POST("/message", message)
	}
}

func ping(c *bm.Context) {
	c.String(0, "pong")
}

func register(c *bm.Context) {
	c.JSON(map[string]interface{}{}, nil)
}

func notice(c *bm.Context) {
	args := &v1.NoticeRequest{}
	if err := c.Bind(args); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	devs := []*grpc.Device{
		{
			RegisterID: args.RegID,
			Platform:   args.Platform,
			SDK:        args.SDK,
			SendNo:     0,
		},
	}

	ext := make(map[string]string)
	ext["schema"] = args.Schema
	ext["callback"] = args.Callback
	extBytes, _ := jsoniter.Marshal(ext)
	req := &grpc.NotificationRequest{
		Dev: devs,
		Body: &grpc.NotificationBody{
			Title:   args.Title,
			Content: args.Content,
			Extra:   string(extBytes),
		},
	}
	resp, err := srv.Notification(c, req)
	if err != nil {
		log.Errorv(c, log.KV("error", err))
	}
	c.String(0, resp.String())
}

func message(c *bm.Context) {
	args := &v1.MessageRequest{}
	if err := c.Bind(args); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	devs := []*grpc.Device{
		{
			RegisterID: args.RegID,
			Platform:   args.Platform,
			SDK:        args.SDK,
			SendNo:     0,
		},
	}

	ext := make(map[string]string)
	ext["schema"] = args.Schema
	ext["callback"] = args.Callback
	extBytes, _ := jsoniter.Marshal(ext)
	req := &grpc.MessageRequest{
		Dev: devs,
		Body: &grpc.MessageBody{
			Title:       args.Title,
			Content:     args.Content,
			ContentType: args.ContentType,
			Extra:       string(extBytes),
		},
	}
	resp, err := srv.Message(c, req)
	if err != nil {
		log.Errorv(c, log.KV("error", err))
	}
	c.String(0, resp.String())
}
