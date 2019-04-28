package main

import (
	"flag"
	"os"

	"go-common/app/admin/main/laser/conf"
	"go-common/app/admin/main/laser/http"
	"go-common/library/log"
	"go-common/library/os/signal"
	"go-common/library/syscall"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Log)
	defer log.Close()
	http.Init(conf.Conf)
	log.Info("laser-admin start")
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("laser-admin get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			log.Info("laser-admin exit")
			return
		case syscall.SIGHUP:
			// TODO reload
		default:
			return
		}
	}

}
