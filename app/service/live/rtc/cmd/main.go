package main

import (
	"flag"
	"go-common/app/service/live/rtc/internal/server/grpc"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/service/live/rtc/internal/conf"
	"go-common/app/service/live/rtc/internal/server/http"
	"go-common/app/service/live/rtc/internal/service"
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
	log.Info("rtc-service start")
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	ecode.Init(conf.Conf.Ecode)
	grpc.New(conf.Conf)
	svc := service.New(conf.Conf)
	http.Init(conf.Conf, svc)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			svc.Close()
			log.Info("rtc-service exit")
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
