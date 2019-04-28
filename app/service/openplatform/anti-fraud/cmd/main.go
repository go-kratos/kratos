package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"go-common/app/service/openplatform/anti-fraud/conf"
	"go-common/app/service/openplatform/anti-fraud/server/http"
	"go-common/library/log"
	"go-common/library/net/trace"
)

func main() {
	flag.Parse()

	if err := conf.Init(); err != nil {
		panic(err)
	}

	log.Init(conf.Conf.Log)
	defer log.Close()

	trace.Init(nil)
	defer trace.Close()
	log.Info("anti_fraud start")
	// http init
	http.Init(conf.Conf)

	// signal handler
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("anti-fraud get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			log.Info("anti-fraud exit")
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
