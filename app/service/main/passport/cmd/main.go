package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/service/main/passport/conf"
	"go-common/app/service/main/passport/http"
	rpc "go-common/app/service/main/passport/rpc/server"
	"go-common/app/service/main/passport/service"
	"go-common/library/log"
	xrpc "go-common/library/net/rpc"
	"go-common/library/net/trace"
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

	var rpcSvr *xrpc.Server
	if conf.Conf.Switch.RPC {
		rpcSvr = rpc.New(conf.Conf, svr)
	}

	http.Init(conf.Conf, svr)
	// signal handler
	log.Info("passport-service start")
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("passport-service get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Info("passport-service exit")
			if conf.Conf.Switch.RPC {
				rpcSvr.Close()
				time.Sleep(time.Second * 2)
			}
			svr.Close()
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
