package http

import (
	"net/http"

	"go-common/app/service/main/account-recovery/conf"
	"go-common/app/service/main/account-recovery/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	srv     *service.Service
	vfy     *verify.Verify
	authSvr *permit.Permit
)

// Init init
func Init(c *conf.Config) {
	srv = service.New(c)
	vfy = verify.New(c.Verify)
	authSvr = permit.New(c.Auth)
	engine := bm.DefaultServer(c.BM)
	router(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start() error(%v)", err)
		panic(err)
	}
}

func router(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)
	g := e.Group("/x/account-recovery")
	{
		g.POST("/query", queryAccount)
		g.POST("/commit", commitInfo)
		g.POST("/getCaptchaMail", getCaptchaMail)
		g.GET("/captcha/token", webToken)
		g.POST("/captcha/verify", verifyCode)
		g.POST("/file/upload", bm.CORS(), fileUpload)
	}
	adminGroup := e.Group("/x/admin/account-recovery", authSvr.Permit("ACCOUNT_RECOVERY_NORMAL"))
	{
		adminGroup.GET("/queryCon", queryConWithBusiness(""))
		adminGroup.GET("/queryCon/account", authSvr.Permit("ACCOUNT_RECOVERY_ACCOUNT"), queryConWithBusiness("account"))
		adminGroup.GET("/queryCon/game", authSvr.Permit("ACCOUNT_RECOVERY_GAME"), queryConWithBusiness("game"))
		// todo judge,batchJudge也需要分配权限
		adminGroup.POST("/judge", judge)
		adminGroup.POST("/batchJudge", batchJudge)
		adminGroup.POST("/gameList", gameList)
	}

	internalGroup := e.Group("/x/internal/account-recovery", vfy.Verify)
	{
		internalGroup.POST("/compare", compareInfo)
		internalGroup.POST("/sendMail", sendMail)
	}
}

func ping(c *bm.Context) {
	if err := srv.Ping(c); err != nil {
		log.Error("ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func register(c *bm.Context) {
	c.JSON(map[string]interface{}{}, nil)
}
