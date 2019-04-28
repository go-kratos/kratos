package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/job/main/passport/conf"
	"go-common/app/job/main/passport/http"
	"go-common/app/job/main/passport/service"
	"go-common/library/log"
	"go-common/library/net/trace"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Xlog)
	defer log.Close()
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	// service init
	srv := service.New(conf.Conf)
	http.Init(conf.Conf, srv)
	// signal handler
	log.Info("passport-job start")
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("passport-job get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			srv.Close()
			time.Sleep(time.Second * 2)
			log.Info("passport-job exit")
			return
		case syscall.SIGHUP:
			// TODO reload
		default:
			return
		}
	}
}
