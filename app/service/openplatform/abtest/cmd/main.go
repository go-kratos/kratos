package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/service/openplatform/abtest/conf"
	"go-common/app/service/openplatform/abtest/http"
	"go-common/app/service/openplatform/abtest/service"
	"go-common/library/log"
	"go-common/library/net/trace"
)

func main() {
	flag.Parse()
	// init conf,log,trace,stat,perf.
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Log)
	defer log.Close()
	trace.Init(nil)
	defer trace.Close()
	// service init
	svr := service.New(conf.Conf)

	http.Init(conf.Conf, svr)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("abtest-service get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			time.Sleep(time.Second * 2)
			log.Info("abtest-service exit")
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
