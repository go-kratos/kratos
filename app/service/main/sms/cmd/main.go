package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/service/main/sms/conf"
	gorpc "go-common/app/service/main/sms/server/gorpc"
	grpc "go-common/app/service/main/sms/server/grpc"
	"go-common/app/service/main/sms/server/http"
	"go-common/app/service/main/sms/service"
	"go-common/library/log"
	"go-common/library/net/trace"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Xlog)
	defer log.Close()
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	svr := service.New(conf.Conf)
	gorpcSvr := gorpc.New(conf.Conf, svr)
	grpcSvr := grpc.New(conf.Conf.GRPC, svr)
	http.Init(conf.Conf, svr)
	log.Info("sms-service start")
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("sms-service get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			gorpcSvr.Close()
			grpcSvr.Shutdown(context.Background())
			svr.Close()
			time.Sleep(time.Second * 2)
			log.Info("sms-service exit")
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
