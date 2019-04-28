package http

import (
	"io"
	"strconv"
	"strings"

	"go-common/app/infra/config/conf"
	"go-common/app/infra/config/service/v1"
	"go-common/app/infra/config/service/v2"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/antispam"
	v "go-common/library/net/http/blademaster/middleware/verify"

	"github.com/dgryski/go-farm"
)

var (
	cnf      *conf.Config
	verify   *v.Verify
	confSvc  *v1.Service
	confSvc2 *v2.Service
	anti     *antispam.Antispam
)

// Init init.
func Init(c *conf.Config, s *v1.Service, s2 *v2.Service, rpcCloser io.Closer) {
	initService(c)
	verify = v.New(c.Verify)
	cnf = c
	confSvc = s
	confSvc2 = s2
	engine := bm.DefaultServer(c.BM)
	innerRouter(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start() error(%v)", err)
		panic(err)
	}
}

// innerRouter init inner router.
func innerRouter(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)
	b := e.Group("/", verify.Verify)
	noAuth := e.Group("/")
	{
		v1 := b.Group("v1/config/")
		{
			v1.GET("host/infos", hosts)
			v1.POST("host/clear", clearhost)
			v1.POST("push", push)
		}
		{
			noAuth.GET("v1/config/versions", versions)
			noAuth.GET("v1/config/builds", builds)
			noAuth.GET("v1/config/check", check)
			noAuth.GET("v1/config/get", config)
			noAuth.GET("v1/config/get2", configN)
			noAuth.GET("v1/config/file.so", file)
			noAuth.GET("v1/config/version/ing", versionIng)
			noAuth.POST("v1/config/config/add", addConfigs)
			noAuth.POST("v1/config/config/copy", copyConfigs)
			noAuth.POST("v1/config/config/update", updateConfigs)

			noAuth.GET("config/v2/versions", versions2)
			noAuth.GET("config/v2/builds", builds2)
			noAuth.GET("config/v2/check", check2)
			noAuth.GET("config/v2/get", setMid, anti.ServeHTTP, config2)
			noAuth.GET("config/v2/file.so", file2)
			noAuth.GET("config/v2/latest", latest)
		}
		v2 := b.Group("config/v2/")
		{
			v2.POST("host/clear", clearhost2)
		}
	}
}
func setMid(c *bm.Context) {
	var (
		token   string
		service string
		query   = c.Request.URL.Query()
		hash    uint64
	)
	service = query.Get("service")
	if service == "" {
		token = query.Get("token")
		if token == "" {
			c.JSON(nil, ecode.RequestErr)
			c.Abort()
			return
		}
		hash = farm.Hash64([]byte(token))
	} else {
		arrs := strings.Split(service, "_")
		if len(arrs) != 3 {
			c.JSON(nil, ecode.RequestErr)
			c.Abort()
			return
		}
		_, err := strconv.ParseInt(arrs[0], 10, 64)
		if err != nil {
			c.JSON(nil, ecode.RequestErr)
			c.Abort()
			return
		}
		hash = farm.Hash64([]byte(service))
	}
	c.Set("mid", int64(hash))
}

func initService(c *conf.Config) {
	anti = antispam.New(c.Antispam)
}
