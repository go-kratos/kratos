package main

import (
	"flag"
	"go-common/app/service/live/resource/sdk"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/service/live/gift/conf"
	"go-common/app/service/live/gift/server/grpc"
	"go-common/app/service/live/gift/server/http"
	"go-common/app/service/live/gift/service"
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
	log.Info("gift-service start")
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	ecode.Init(conf.Conf.Ecode)
	svc := service.New(conf.Conf)
	http.Init(conf.Conf, svc)
	grpc.Init(conf.Conf)
	titansSdk.Init(conf.Conf.Titan)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			svc.Close()
			log.Info("gift-service exit")
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
