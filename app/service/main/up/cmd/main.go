package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/service/main/up/conf"
	"go-common/app/service/main/up/server/gorpc"
	"go-common/app/service/main/up/server/grpc"
	"go-common/app/service/main/up/server/http"
	"go-common/app/service/main/up/service"
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
	// init log
	log.Init(conf.Conf.Xlog)
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	defer log.Close()
	log.SetFormat("[%D %T] [%L] [%S] %M")
	log.Info("up-servicestart")
	svr := service.New(conf.Conf)
	ecode.Init(conf.Conf.Ecode)
	// service init
	http.Init(conf.Conf, svr)
	rpcSvr := gorpc.New(conf.Conf, svr)
	grpcSvr := grpc.New(conf.Conf.GRPCServer, svr)
	report.InitManager(conf.Conf.ManagerReport)
	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("up-service get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			rpcSvr.Close()
			grpcSvr.Shutdown(context.TODO())
			time.Sleep(time.Second * 2)
			svr.Close()
			log.Info("up-service exit")
			return
		case syscall.SIGHUP:
			// TODO reload
		default:
			return
		}
	}
}
