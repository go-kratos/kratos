package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/service/main/coupon/conf"
	"go-common/app/service/main/coupon/http"
	rpc "go-common/app/service/main/coupon/rpc/server"
	grpc "go-common/app/service/main/coupon/server/grpc"
	"go-common/app/service/main/coupon/service"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	xrpc "go-common/library/net/rpc"
)

var (
	svc    *service.Service
	rpcSvr *xrpc.Server
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
	log.Info("coupon start")
	// ecode init
	ecode.Init(conf.Conf.Ecode)
	// service init
	svc = service.New(conf.Conf)
	// rpc init
	rpcSvr = rpc.New(conf.Conf, svc)
	// http init
	http.Init(conf.Conf, svc)
	// grpc
	ws := grpc.New(conf.Conf.WardenServer, svc)
	// init pprof conf.Conf.Perf
	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("coupon get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			rpcSvr.Close()
			ws.Shutdown(context.Background())
			time.Sleep(time.Second * 1)
			log.Info("coupon exit")
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
