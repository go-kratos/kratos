package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/admin/main/workflow/http"
	"go-common/app/admin/main/workflow/service"
	"go-common/library/conf/paladin"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	"go-common/library/queue/databus/report"
)

func main() {
	flag.Parse()
	if err := paladin.Init(); err != nil {
		panic(err)
	}
	// init log
	log.Init(nil)
	defer log.Close()
	log.Info("workflow-admin start")
	svc := service.New()
	http.Init(svc)
	// report
	report.InitManager(nil)
	// ecode init
	ecode.Init(nil)
	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("workflow-admin get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Info("workflow-admin exit")
			time.Sleep(1 * time.Second)
			svc.Close()
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
