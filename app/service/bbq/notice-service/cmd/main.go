package main

import (
	"context"
	"flag"
	"go-common/app/service/bbq/notice-service/internal/server/grpc"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/service/bbq/notice-service/internal/conf"
	"go-common/app/service/bbq/notice-service/internal/server/http"
	"go-common/app/service/bbq/notice-service/internal/service"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	"go-common/library/net/trace"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Log)
	defer log.Close()
	log.Info("notice-service-service start")
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	ecode.Init(conf.Conf.Ecode)
	svc := service.New(conf.Conf)
	grpcServ := grpc.New(conf.Conf.GRPC, svc)
	http.Init(conf.Conf, svc)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			svc.Close()
			grpcServ.Shutdown(context.Background())
			log.Info("notice-service-service exit")
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
