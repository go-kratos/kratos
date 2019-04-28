package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/interface/openplatform/article/conf"
	"go-common/app/interface/openplatform/article/http"
	artrpc "go-common/app/interface/openplatform/article/rpc/server"
	"go-common/app/interface/openplatform/article/service"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	"go-common/library/net/trace"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		log.Error("conf.Init() error(%+v)", err)
		panic(err)
	}
	log.Init(conf.Conf.Log)
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	defer log.Close()
	log.Info("article start")
	// server init
	ecode.Init(conf.Conf.Ecode)
	svr := service.New(conf.Conf)
	http.Init(conf.Conf, svr)
	rpcSvr := artrpc.New(svr)
	// signal handler
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("article get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			log.Info("article exit")
			rpcSvr.Close()
			svr.Close()
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
