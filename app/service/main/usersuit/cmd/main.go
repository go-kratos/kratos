package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/service/main/usersuit/conf"
	"go-common/app/service/main/usersuit/http"
	rpc "go-common/app/service/main/usersuit/rpc/server"
	"go-common/app/service/main/usersuit/service"
	"go-common/library/log"
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
	rpcSvr := rpc.New(conf.Conf, svr)
	http.Init(conf.Conf, svr)
	// signal handler
	log.Info("usersuit-service start")
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("usersuit-service get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			rpcSvr.Close()
			time.Sleep(time.Second * 2)
			svr.Close()
			log.Info("usersuit-service exit")
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
