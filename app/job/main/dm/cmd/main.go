package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"go-common/app/job/main/dm/conf"
	"go-common/app/job/main/dm/http"
	"go-common/app/job/main/dm/service"
	"go-common/library/log"
	"go-common/library/net/trace"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		log.Error("conf.Init() error(%v)", err)
		panic(err)
	}
	// init log
	log.Init(conf.Conf.Xlog)
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	defer log.Close()
	// service init
	svc := service.New(conf.Conf)
	http.Init(conf.Conf, svc)
	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("dm-job get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			log.Info("dm-job exit")
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
