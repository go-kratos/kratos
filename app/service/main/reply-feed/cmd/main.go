package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"go-common/app/service/main/reply-feed/conf"
	"go-common/app/service/main/reply-feed/server/grpc"
	"go-common/app/service/main/reply-feed/server/http"
	"go-common/app/service/main/reply-feed/service"
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
	log.Info("reply-feed service start")
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	ecode.Init(conf.Conf.Ecode)
	svc := service.New(conf.Conf)
	http.Init(conf.Conf, svc)
	grpc.New(conf.Conf.GRPC, svc)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			svc.Close()
			log.Info("reply-feed service exit")
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
