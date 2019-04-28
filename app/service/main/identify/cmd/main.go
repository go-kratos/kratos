package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/service/main/identify/conf"
	"go-common/app/service/main/identify/server/grpc"
	"go-common/app/service/main/identify/server/http"
	"go-common/app/service/main/identify/service"
	"go-common/library/log"
	// "go-common/library/net/rpc/warden/resolver/livezk"
	"go-common/library/net/trace"
)

const (
// discoveryID = "passport.service.identify"
)

func main() {
	flag.Parse()
	// init conf,log,trace,stat,perf.
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Xlog)
	defer log.Close()
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	// service init
	svr := service.New(conf.Conf)
	http.Init(conf.Conf, svr)

	// init warden server
	ws := grpc.New(conf.Conf.WardenServer, svr)
	// 先主站内部和chenzhihui测试可用，再对外提供
	// cancel, err := livezk.Register(conf.Conf.LiveZK, conf.Conf.WardenServer.Addr, discoveryID)
	// if err != nil {
	// 	panic(err)
	// }

	// signal handler
	log.Info("identify-service start")
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("identify-service get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			// cancel()
			ws.Shutdown(context.Background())
			time.Sleep(time.Second * 2)
			svr.Close()
			log.Info("identify-service exit")
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
