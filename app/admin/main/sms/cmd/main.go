package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/admin/main/sms/conf"
	"go-common/app/admin/main/sms/server/http"
	"go-common/app/admin/main/sms/service"
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
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	log.Info("sms-admin start")
	srv := service.New(conf.Conf)
	http.Init(conf.Conf, srv)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("sms-admin get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			srv.Close()
			log.Info("sms-admin exit")
			time.Sleep(1 * time.Second)
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
