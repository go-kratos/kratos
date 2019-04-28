package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/interface/main/tag/conf"
	"go-common/app/interface/main/tag/http"
	rpc "go-common/app/interface/main/tag/rpc/server"
	"go-common/app/interface/main/tag/server/grpc"
	"go-common/app/interface/main/tag/service"
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
	defer log.Close()
	log.Info("tag-interface start")
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	svr := service.New(conf.Conf)
	ecode.Init(conf.Conf.Ecode)
	http.Init(conf.Conf, svr)
	rpcSvr := rpc.Init(conf.Conf, svr)
	grpcSvr := grpc.New(conf.Conf.GRPCServer, svr)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("tag-interface get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			rpcSvr.Close()
			grpcSvr.Shutdown(context.TODO())
			time.Sleep(time.Second * 2)
			svr.Close()
			log.Info("tag-interface exit")
			return
		case syscall.SIGHUP:
			// TODO reload
		default:
			return
		}
	}
}
