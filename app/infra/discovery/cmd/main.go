package main

import (
	"flag"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/infra/discovery/conf"
	"go-common/app/infra/discovery/http"
	"go-common/app/infra/discovery/service"
	"go-common/library/log"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		log.Error("conf.Init() error(%v)", err)
		panic(err)
	}
	// init log
	log.Init(conf.Conf.Log)
	defer log.Close()
	log.Info("discovery start")
	// service init
	rand.Seed(time.Now().UnixNano())
	svc, cancel := service.New(conf.Conf)
	http.Init(conf.Conf, svc)
	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("discovery get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			cancel()
			time.Sleep(time.Second)
			log.Info("discovery exit")
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
