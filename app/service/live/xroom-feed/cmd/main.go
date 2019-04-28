package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/service/live/xroom-feed/internal/server/grpc"
	"go-common/app/service/live/xroom-feed/internal/server/http"
	"go-common/app/service/live/xroom-feed/internal/service"
	"go-common/library/conf/paladin"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
)

func main() {
	flag.Parse()
	if err := paladin.Init(); err != nil {
		panic(err)
	}
	log.Init(nil) // debug flag: log.dir={path}
	defer log.Close()
	log.Info("xroom-feed-service start")
	ecode.Init(nil)
	svc := service.New()
	grpcSrv := grpc.New(svc)
	httpSrv := http.New(svc)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			ctx, _ := context.WithTimeout(context.Background(), 35*time.Second)
			grpcSrv.Shutdown(ctx)
			httpSrv.Shutdown(ctx)
			log.Info("xroom-feed-service exit")
			svc.Close()
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
