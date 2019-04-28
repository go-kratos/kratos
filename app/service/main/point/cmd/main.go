package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/service/main/point/conf"
	rpc "go-common/app/service/main/point/rpc/server"
	grpc "go-common/app/service/main/point/server/grpc"
	"go-common/app/service/main/point/server/http"
	"go-common/app/service/main/point/service"
	ecode "go-common/library/ecode/tip"

	"go-common/library/log"
)

var svc *service.Service

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		log.Error("conf.Init() error(%v)", err)
		panic(err)
	}
	// init log
	log.Init(conf.Conf.Log)
	defer log.Close()
	log.Info("point start")
	// ecode init
	ecode.Init(conf.Conf.Ecode)
	// service init
	svc = service.New(conf.Conf)
	// rpc init
	rpcSvr := rpc.New(conf.Conf, svc)
	ws := grpc.New(conf.Conf.WardenServer, svc)
	// service init
	http.Init(svc)
	// init pprof conf.Conf.Perf
	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("point get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			rpcSvr.Close()
			ws.Shutdown(context.Background())
			time.Sleep(time.Second * 2)
			log.Info("point exit")
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
