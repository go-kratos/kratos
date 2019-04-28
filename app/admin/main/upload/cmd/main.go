package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"go-common/app/admin/main/upload/conf"
	"go-common/app/admin/main/upload/http"
	"go-common/app/admin/main/upload/service"
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
	log.Info("upload start, listening: %s \n", conf.Conf.BM.Addr)
	// service init
	svc := service.New(conf.Conf)
	http.Init(conf.Conf, svc)
	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("upload get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Info("upload exit")
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
