package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"go-common/app/job/main/passport-sns/conf"
	"go-common/app/job/main/passport-sns/http"
	"go-common/library/log"
	"go-common/library/net/trace"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		log.Error("conf.Init() error(%v)", err)
		panic(err)
	}
	log.Init(conf.Conf.Xlog)
	defer log.Close()
	log.Info("passport-sns-job start")
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	http.Init(conf.Conf)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("passport-sns-job get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Info("passport-sns-job exit")
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
