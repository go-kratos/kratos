package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/service/main/coin/conf"
	"go-common/app/service/main/coin/server/gorpc"
	grpc "go-common/app/service/main/coin/server/grpc"
	"go-common/app/service/main/coin/server/http"
	"go-common/app/service/main/coin/service"
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
	defer log.Close()
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	ecode.Init(nil)
	log.Info("coin-service start")
	report.InitUser(conf.Conf.UserReport)
	// service init
	svr := service.New(conf.Conf)
	rpcSvr := rpc.New(conf.Conf, svr)
	grpcSvr := grpc.New(conf.Conf.GRPC, svr)
	http.Init(conf.Conf, svr)
	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("coin-service get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			grpcSvr.Shutdown(context.TODO())
			rpcSvr.Close()
			svr.Close()
			time.Sleep(time.Second * 1)
			log.Info("coin-service exit")
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
