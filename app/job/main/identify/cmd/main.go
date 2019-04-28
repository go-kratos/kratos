package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/job/main/identify/conf"
	"go-common/app/job/main/identify/http"
	"go-common/app/job/main/identify/service"
	"go-common/library/log"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		log.Error("conf.Init() error(%v)", err)
		panic(err)
	}
	log.Init(conf.Conf.Log)
	defer log.Close()
	log.Info("identify-job start")
	// init
	svc := service.New(conf.Conf)
	http.Init(conf.Conf, svc)
	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("identify job get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			svc.Close()
			time.Sleep(time.Second)
			log.Info("identify job exit")
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
