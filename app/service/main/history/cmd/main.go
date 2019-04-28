package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/service/main/history/conf"
	"go-common/app/service/main/history/server/grpc"
	"go-common/app/service/main/history/server/http"
	"go-common/app/service/main/history/service"
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
	log.Info("start")
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	svr := service.New(conf.Conf)
	http.Init(conf.Conf, svr)
	grpcSvr := grpc.New(conf.Conf.GRPC, svr)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			grpcSvr.Shutdown(context.TODO())
			time.Sleep(time.Second * 1)
			log.Info("exit")
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
