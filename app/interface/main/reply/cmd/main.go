package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"go-common/app/interface/main/reply/conf"
	"go-common/app/interface/main/reply/http"
	ecode "go-common/library/ecode/tip"
	"go-common/library/exp/feature"
	"go-common/library/log"
	"go-common/library/net/trace"
	"go-common/library/queue/databus/report"
)

func main() {
	feature.DefaultGate.AddFlag(flag.CommandLine)
	flag.Parse()
	if err := conf.Init(); err != nil {
		log.Error("conf.Init() error(%v)", err)
		panic(err)
	}
	log.Init(conf.Conf.XLog)
	ecode.Init(conf.Conf.Ecode)
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	defer log.Close()
	log.Info("reply start")
	// init report agent
	report.InitUser(conf.Conf.UserReport)
	// init http
	http.Init(conf.Conf)
	// init pprof
	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("reply get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Info("reply exit")
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
