package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/service/main/tv/internal/conf"
	"go-common/app/service/main/tv/internal/server/grpc"
	"go-common/app/service/main/tv/internal/server/http"
	"go-common/app/service/main/tv/internal/service"
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
	log.Info("tv-service start")
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
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
			log.Info("tv-service exit")
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
