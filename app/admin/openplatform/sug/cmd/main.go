package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"go-common/app/admin/openplatform/sug/conf"
	"go-common/app/admin/openplatform/sug/http"
	"go-common/app/admin/openplatform/sug/service"
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

	log.Info("sug-season start")
	ecode.Init(conf.Conf.Ecode)
	s := service.New(conf.Conf)
	http.Init(conf.Conf, s)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("open-sug get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Info("open-sug exit")
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
