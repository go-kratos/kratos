package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/interface/main/history/conf"
	"go-common/app/interface/main/history/http"
	rpc "go-common/app/interface/main/history/server/gorpc"
	"go-common/app/interface/main/history/server/grpc"
	"go-common/app/interface/main/history/service"
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
	log.Init(conf.Conf.Xlog)
	defer log.Close()
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	log.Info("history start")
	ecode.Init(conf.Conf.Ecode)
	report.InitUser(conf.Conf.Report)
	svr := service.New(conf.Conf)
	rpcSvr := rpc.New(conf.Conf, svr)
	grpcSvr := grpc.New(conf.Conf.GRPC, svr)
	http.Init(conf.Conf, svr)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("history get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Info("history exit")
			grpcSvr.Shutdown(context.TODO())
			time.Sleep(2 * time.Second)
			rpcSvr.Close()
			svr.Close()
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
