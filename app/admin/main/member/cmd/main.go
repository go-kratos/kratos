package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"go-common/app/admin/main/member/conf"
	"go-common/app/admin/main/member/http"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	"go-common/library/net/trace"
	"go-common/library/queue/databus/report"
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
	// report log
	report.InitManager(conf.Conf.ManagerReport)
	log.Info("member start")
	// init trace
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	// ecode init
	ecode.Init(conf.Conf.Ecode)
	// service init
	http.Init(conf.Conf)
	// init pprof conf.Conf.Perf
	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("member get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Info("member exit")
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
