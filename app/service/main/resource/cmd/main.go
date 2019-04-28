package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/service/main/resource/conf"
	"go-common/app/service/main/resource/http"
	rpc "go-common/app/service/main/resource/rpc/server"
	grpc "go-common/app/service/main/resource/server/grpc"
	"go-common/app/service/main/resource/service"
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
	log.Init(conf.Conf.XLog)
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	defer log.Close()
	log.Info("resource-service start")
	// service init
	svr := service.New(conf.Conf)
	rpcSvr := rpc.New(conf.Conf, svr)
	grpcSvr := grpc.New(nil, svr)
	http.Init(conf.Conf, svr)
	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("resource-service get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			rpcSvr.Close()
			grpcSvr.Shutdown(context.TODO())
			time.Sleep(time.Second * 2)
			log.Info("resource-service exit")
			return
		case syscall.SIGHUP:
			// TODO reload
		default:
			return
		}
	}
}
