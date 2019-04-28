package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"go-common/app/service/main/msm/conf"
	"go-common/app/service/main/msm/http"
	"go-common/app/service/main/msm/service"
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
	log.Info("msm-service start")
	svr := service.New(conf.Conf)
	http.Init(conf.Conf, svr)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("msm-service get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			svr.Close()
			log.Info("msm-service exit")
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
