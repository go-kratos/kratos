package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"go-common/app/job/main/appstatic/conf"
	"go-common/app/job/main/appstatic/http"
	"go-common/app/job/main/appstatic/service"
	"go-common/library/log"
	"go-common/library/net/trace"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Log)
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	defer log.Close()
	log.Info("appstatic-job start")
	srv := service.New(conf.Conf)
	http.Init(conf.Conf, srv)
	initSignal(srv)
}

func initSignal(srv *service.Service) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("appstatic-job get a signal: %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			if err := srv.Close(); err != nil {
				log.Error("srv close consumer error(%v)", err)
			}
			return
		case syscall.SIGHUP:
			// TODO: reload
		default:
			return
		}
	}
}
