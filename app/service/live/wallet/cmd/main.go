package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"go-common/app/service/live/wallet/conf"
	"go-common/app/service/live/wallet/http"
	"go-common/app/service/live/wallet/service"
	"go-common/library/log"
	"go-common/library/net/trace"
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
	log.Info("/live-wallet start")
	// init trace
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	// service init
	svr := service.New(conf.Conf)
	http.Init(conf.Conf, svr)
	// init pprof conf.Conf.Perf
	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("/live-wallet get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Info("/live-wallet exit")
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
