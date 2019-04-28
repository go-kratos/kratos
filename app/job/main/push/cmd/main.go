package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/job/main/push/conf"
	"go-common/app/job/main/push/http"
	"go-common/app/job/main/push/service"
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
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	log.Info("push-job start")
	svc := service.New(conf.Conf)
	http.Init(conf.Conf, svc)
	listenSignals(svc)
}

func listenSignals(svc *service.Service) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("push-job get a signal: %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			svc.Close()
			log.Info("push-job stop")
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
			// TODO: reload
		default:
			return
		}
	}
}
