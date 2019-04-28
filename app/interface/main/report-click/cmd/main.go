package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/interface/main/report-click/conf"
	"go-common/app/interface/main/report-click/http"
	"go-common/library/log"
	"go-common/library/net/trace"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		log.Error("conf.Init() error(%v)", err)
		panic(err)
	}
	log.Init(conf.Conf.Xlog)
	defer log.Close()
	log.Info("report-click start")
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	engine := http.New(conf.Conf)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("report-click get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
			defer cancel()
			engine.Shutdown(ctx)
			log.Info("report-click exit")
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
