package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/interface/live/web-ucenter/conf"
	"go-common/app/interface/live/web-ucenter/server/http"
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
	log.Info("web-ucenter-interface start")
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	ecode.Init(conf.Conf.Ecode)
	http.Init(conf.Conf)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Info("web-ucenter-interface exit")
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
