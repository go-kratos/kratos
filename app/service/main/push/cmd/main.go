package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/service/main/push/conf"
	pushrpc "go-common/app/service/main/push/server/gorpc"
	"go-common/app/service/main/push/server/grpc"
	"go-common/app/service/main/push/server/http"
	"go-common/app/service/main/push/service"
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
	log.Init(conf.Conf.Log)
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	defer log.Close()
	log.Info("push-service start")
	ecode.Init(conf.Conf.Ecode)
	svr := service.New(conf.Conf)
	rpcSrv := pushrpc.New(conf.Conf, svr)
	grpcSvr := grpc.New(conf.Conf.GRPC, svr)
	http.Init(conf.Conf, svr)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("push-service get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			rpcSrv.Close()
			grpcSvr.Shutdown(context.Background())
			svr.Close()
			log.Info("push-service exit")
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
