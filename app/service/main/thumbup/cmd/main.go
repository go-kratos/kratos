package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/service/main/thumbup/conf"
	rpc "go-common/app/service/main/thumbup/server/gorpc"
	grpc "go-common/app/service/main/thumbup/server/grpc"
	"go-common/app/service/main/thumbup/server/http"
	"go-common/app/service/main/thumbup/service"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	"go-common/library/net/trace"
	"go-common/library/queue/databus/report"
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
	log.Info("thumbup start")
	ecode.Init(conf.Conf.Ecode)
	report.InitManager(nil)
	// server init
	svr := service.New(conf.Conf)
	http.Init(conf.Conf, svr)
	rpcSvr := rpc.New(conf.Conf, svr)
	grpcSvr := grpc.New(conf.Conf.GRPC, svr)
	// signal handler
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("thumbup get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Info("thumbup exit")
			rpcSvr.Close()
			grpcSvr.Shutdown(context.TODO())
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
