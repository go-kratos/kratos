package http

import (
	"go-common/app/service/main/passport-sns/api"
	"go-common/app/service/main/passport-sns/conf"
	"go-common/app/service/main/passport-sns/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	svr *service.Service
	vfy *verify.Verify
)

// Init init config
func Init(c *conf.Config) {
	svr = service.New(c)
	vfy = verify.New(c.Verify)
	e := bm.DefaultServer(c.BM)
	e.Inject("/x/internal/passport-sns/", vfy.Verify)
	// 生成工具还不支持，需要后续优化
	e.Ping(func(c *bm.Context) {})
	api.RegisterPassportSNSBMServer(e, svr)
	if err := e.Start(); err != nil {
		log.Error("bm.DefaultServer error(%v)", err)
	}
}
