package main

import (
	"flag"
	"os"
	"time"

	"go-common/app/service/main/figure/conf"
	"go-common/app/service/main/figure/http"
	rpc "go-common/app/service/main/figure/rpc/server"
	"go-common/app/service/main/figure/service"
	"go-common/library/log"
	"go-common/library/net/trace"
	"go-common/library/os/signal"
	"go-common/library/syscall"
)

var svr *service.Service

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		log.Error("conf.Init() error(%v)", err)
		panic(err)
	}
	log.Init(conf.Conf.Log)
	defer log.Close()
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	// service init
	svr = service.New(conf.Conf)
	http.Init(svr)
	rpcSvr := rpc.New(conf.Conf, svr)
	log.Info("figure-service start")
	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSTOP)
	for {
		s := <-c
		log.Info("figure-service get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			rpcSvr.Close()
			time.Sleep(time.Second * 2)
			log.Info("figure-service exit")
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
