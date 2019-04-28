package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/infra/databus/conf"
	"go-common/app/infra/databus/http"
	"go-common/app/infra/databus/service"
	"go-common/app/infra/databus/tcp"
	"go-common/library/log"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	// init log
	log.Init(conf.Conf.Log)
	defer log.Close()
	log.Info("databus start")
	// service init
	svc := service.New(conf.Conf)
	http.Init(conf.Conf, svc)
	tcp.Init(conf.Conf, svc)
	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("databus get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Info("databus exit")
			tcp.Close()
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
