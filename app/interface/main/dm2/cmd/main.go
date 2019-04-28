package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/interface/main/dm2/conf"
	"go-common/app/interface/main/dm2/http"
	rpc "go-common/app/interface/main/dm2/rpc/server"
	"go-common/app/interface/main/dm2/service"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	"go-common/library/net/trace"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		log.Error("conf.Init() error(%v)", err)
		panic(err)
	}
	// init log
	log.Init(conf.Conf.Xlog)
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	defer log.Close()
	// service init
	svc := service.New(conf.Conf)
	http.Init(conf.Conf, svc)
	ecode.Init(conf.Conf.Ecode)
	rpcSvc := rpc.New(conf.Conf, svc)
	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("dm2 get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			rpcSvc.Close()
			time.Sleep(1 * time.Second)
			log.Info("dm2 interface exit")
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
