package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"go-common/app/interface/main/app-resource/conf"
	"go-common/app/interface/main/app-resource/grpc"
	"go-common/app/interface/main/app-resource/http"
	absvr "go-common/app/interface/main/app-resource/service/abtest"
	auditsvr "go-common/app/interface/main/app-resource/service/audit"
	broadcastsvr "go-common/app/interface/main/app-resource/service/broadcast"
	domainsvr "go-common/app/interface/main/app-resource/service/domain"
	guidesvc "go-common/app/interface/main/app-resource/service/guide"
	modulesvr "go-common/app/interface/main/app-resource/service/module"
	"go-common/app/interface/main/app-resource/service/notice"
	"go-common/app/interface/main/app-resource/service/param"
	pingsvr "go-common/app/interface/main/app-resource/service/ping"
	pluginsvr "go-common/app/interface/main/app-resource/service/plugin"
	showsvr "go-common/app/interface/main/app-resource/service/show"
	sidesvr "go-common/app/interface/main/app-resource/service/sidebar"
	"go-common/app/interface/main/app-resource/service/splash"
	staticsvr "go-common/app/interface/main/app-resource/service/static"
	"go-common/app/interface/main/app-resource/service/version"
	whitesvr "go-common/app/interface/main/app-resource/service/white"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/trace"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		log.Error("conf.Init() error(%v)", err)
		panic(err)
	}
	// init log
	log.Init(conf.Conf.Log)
	defer log.Close()
	log.Info("app-resource start")
	// init trace
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	// ecode init
	ecode.Init(conf.Conf.Ecode)
	// service init
	svr := initService(conf.Conf)
	http.Init(conf.Conf, svr)
	grpcSvr, err := grpc.New(nil, svr)
	if err != nil {
		panic(err)
	}
	// init pprof conf.Conf.Perf
	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("app-resource get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			grpcSvr.Shutdown(context.TODO())
			log.Info("app-resource exit")
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}

// initService init services.
func initService(c *conf.Config) (svr *http.Server) {
	svr = &http.Server{
		AuthSvc: auth.New(nil),
		// init self service,
		PgSvr:        pluginsvr.New(c),
		PingSvr:      pingsvr.New(c),
		SideSvr:      sidesvr.New(c),
		VerSvc:       version.New(c),
		ParamSvc:     param.New(c),
		NtcSvc:       notice.New(c),
		SplashSvc:    splash.New(c),
		AuditSvc:     auditsvr.New(c),
		AbSvc:        absvr.New(c),
		ModuleSvc:    modulesvr.New(c),
		GuideSvc:     guidesvc.New(c),
		StaticSvc:    staticsvr.New(c),
		DomainSvc:    domainsvr.New(c),
		BroadcastSvc: broadcastsvr.New(c),
		WhiteSvc:     whitesvr.New(c),
		ShowSvc:      showsvr.New(c),
	}
	return
}
