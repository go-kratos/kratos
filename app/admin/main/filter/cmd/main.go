package main

import (
	"flag"
	"os"
	"runtime"
	"time"

	"go-common/app/admin/main/filter/conf"
	"go-common/app/admin/main/filter/http"
	"go-common/app/admin/main/filter/service"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	"go-common/library/net/trace"
	"go-common/library/os/signal"
	"go-common/library/syscall"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	if err := conf.Init(); err != nil {
		log.Error("conf.Init() error(%v)", err)
		panic(err)
	}
	// init log
	log.Init(conf.Conf.Log)
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	defer log.Close()
	log.Info("filter-admin start")
	// service init
	svc := service.New(conf.Conf)
	http.Init(svc)
	ecode.Init(conf.Conf.Ecode)
	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("filter-admin get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			svc.Close()
			time.Sleep(2 * time.Second)
			log.Info("filter-admin exit")
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
