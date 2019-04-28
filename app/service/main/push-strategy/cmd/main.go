package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/service/main/push-strategy/conf"
	"go-common/app/service/main/push-strategy/http"
	"go-common/app/service/main/push-strategy/service"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	"go-common/library/net/trace"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		log.Error("conf.Init() error(%v)", err)
		panic(err)
	}
	log.Init(conf.Conf.Log)
	defer log.Close()
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	log.Info("push-strategystart")
	ecode.Init(conf.Conf.Ecode)
	srv := service.New(conf.Conf)
	http.Init(conf.Conf, srv)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("push-strategy get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			srv.Close()
			time.Sleep(1 * time.Second)
			log.Info("push-strategyexit")
			return
		case syscall.SIGHUP:
			// TODO reload
		default:
			return
		}
	}
}
