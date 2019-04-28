package main

import (
	"flag"
	"go-common/library/queue/databus/report"
	"os"
	"os/signal"
	"time"

	"go-common/app/service/live/userexp/conf"
	"go-common/app/service/live/userexp/http"
	"go-common/app/service/live/userexp/service"
	"go-common/library/log"
	"go-common/library/net/trace"
	"go-common/library/syscall"
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
	// report
	report.InitUser(conf.Conf.Report)
	log.Info("live-userexp-service start")
	// service init
	svr := service.New(conf.Conf)
	http.Init(conf.Conf, svr)
	// signal handler
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSTOP)
	for {
		s := <-c
		log.Info("live-userexp-service get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			time.Sleep(time.Second * 2)
			log.Info("live-userexp-service exit")
			return
		case syscall.SIGHUP:
			log.Info("TODO: reload for syscall.SIGHUP")
			return
		default:
			return
		}
	}
}
