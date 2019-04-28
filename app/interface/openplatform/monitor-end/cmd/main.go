package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/interface/openplatform/monitor-end/conf"
	"go-common/app/interface/openplatform/monitor-end/http"
	"go-common/app/interface/openplatform/monitor-end/service"
	"go-common/library/log"
	"go-common/library/net/trace"
)

func main() {
	flag.Parse()
	conf.Init()
	log.Init(conf.Conf.Log)
	defer log.Close()
	trace.Init(nil)
	defer trace.Close()
	svr := service.New(conf.Conf)
	defer svr.Close()
	http.Init(conf.Conf, svr)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			time.Sleep(time.Second) // 休眠1s
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
