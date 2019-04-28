package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/service/main/card/conf"
	grpc "go-common/app/service/main/card/server/grpc"
	"go-common/app/service/main/card/server/http"
	"go-common/app/service/main/card/service"
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
	log.Info("card-service start")
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	svc := service.New(conf.Conf)
	http.Init(conf.Conf, svc)
	ws := grpc.New(conf.Conf.WardenServer, svc)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			svc.Close()
			ws.Shutdown(context.Background())
			log.Info("card-service exit")
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
