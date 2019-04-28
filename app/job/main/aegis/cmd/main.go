package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/job/main/aegis/conf"
	"go-common/app/job/main/aegis/server/http"
	"go-common/app/job/main/aegis/service"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	"go-common/library/net/trace"
	"go-common/library/queue/databus/report"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Log)
	defer log.Close()
	log.Info("aegis-job start")
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	if !conf.Conf.Debug {
		ecode.Init(conf.Conf.Ecode)
	}

	// report init
	report.InitManager(nil)
	svc := service.New(conf.Conf)
	http.Init(svc, conf.Conf)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			svc.Close()
			log.Info("aegis-job exit")
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
