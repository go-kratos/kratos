package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"go-common/app/job/main/workflow/conf"
	"go-common/app/job/main/workflow/http"
	"go-common/app/job/main/workflow/service"
	"go-common/library/log"
	"go-common/library/queue/databus/report"
)

var (
	srv *service.Service
)

func main() {
	flag.Parse()
	report.InitManager(nil)
	if err := conf.Init(); err != nil {
		log.Error("conf.Init() error(%v)", err)
		panic(err)
	}
	log.Init(conf.Conf.Log)
	defer log.Close()
	log.Info("workflow-job start")
	srv = service.New(conf.Conf)
	http.Init(conf.Conf, srv)
	signalHandler()
}

func signalHandler() {
	var (
		err error
		ch  = make(chan os.Signal, 1)
	)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		si := <-ch
		switch si {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Info("get a signal %s, stop the consume process", si.String())
			if err = srv.Close(); err != nil {
				log.Error("srv close consumer error(%v)", err)
			}
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
