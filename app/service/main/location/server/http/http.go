package http

import (
	"io"

	"go-common/app/service/main/location/conf"
	"go-common/app/service/main/location/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	svr    *service.Service
	vfySvc *verify.Verify
)

// Init init http.
func Init(c *conf.Config, s *service.Service, rpcCloser io.Closer) {
	svr = s
	vfySvc = verify.New(c.Verify)
	// init external router
	engineInner := bm.DefaultServer(c.BM.Inner)
	innerRouter(engineInner)
	// init internal server
	if err := engineInner.Start(); err != nil {
		log.Error("engineInner.Start() error(%v) | config(%v)", err, c)
		panic(err)
	}
}

// innerRouter init inner router.
func innerRouter(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)
	loc := e.Group("/x/location")
	// old api will delete soon
	loc.GET("/check", vfySvc.Verify, tmpInfo)
	loc.GET("/zone", vfySvc.Verify, tmpInfos)
	// ip info
	loc.GET("/info", vfySvc.Verify, info)
	loc.GET("/infos", vfySvc.Verify, infos)
	loc.GET("/info/complete", vfySvc.Verify, infoComplete)
	loc.GET("/infos/complete", vfySvc.Verify, infosComplete)
	loc.GET("/anonym", vfySvc.Verify, anonym)
	// area limit api
	zl := loc.Group("/zlimit")
	zl.GET("/pgc/in", vfySvc.Verify, pgcZone)
	zl.GET("/archive", vfySvc.Verify, auth)
	zl.GET("/archive2", vfySvc.Verify, archive2)
	zl.GET("/group", vfySvc.Verify, authGID)
	zl.GET("/groups", vfySvc.Verify, authGIDs)

	mng := loc.Group("/manager")
	mng.GET("/flushCache", vfySvc.Verify, flushCache)
}
