package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/job/bbq/recall/internal/conf"
	"go-common/app/job/bbq/recall/internal/service"
	"go-common/library/log"
	"go-common/library/net/trace"
)

var (
	_serviceName string
)

func init() {
	flag.StringVar(&_serviceName, "service", "", "run service name")
}

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Log)
	defer log.Close()
	log.Info("recall-job start")
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	svc := service.New(conf.Conf)
	defer svc.Close()

	if _serviceName != "" {
		svc.RunSrv(_serviceName)
	} else {
		svc.InitCron()
		deamon()
	}
}

func deamon() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Info("recall-job exit")
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
