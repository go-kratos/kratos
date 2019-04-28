package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/service/main/feed/conf"
	"go-common/app/service/main/feed/http"
	rpc "go-common/app/service/main/feed/rpc/server"
	"go-common/app/service/main/feed/service"
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
	http.Init(conf.Conf, svr)
	rpcSvr := rpc.New(conf.Conf, svr)
	// signal handler
	log.Info("feed-service start")
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("feed-service get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Info("feed-service exit")
			rpcSvr.Close()
			time.Sleep(time.Second * 2)
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
