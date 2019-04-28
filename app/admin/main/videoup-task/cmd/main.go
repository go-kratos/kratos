package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"go-common/app/admin/main/videoup-task/conf"
	"go-common/app/admin/main/videoup-task/http"
	"go-common/app/admin/main/videoup-task/service"
	"go-common/library/log"
	"go-common/library/net/trace"
	"go-common/library/os/signal"
	"go-common/library/queue/databus/report"
	"go-common/library/syscall"
)

func main() {
	var err error
	//conf init
	flag.Parse()
	if err = conf.Init(); err != nil {
		panic(err)
	}
	//log init
	log.Init(conf.Conf.Xlog)
	defer log.Close()
	fmt.Printf("conf(%+v)", conf.Conf)
	report.InitManager(conf.Conf.ManagerReport)

	//trace init
	trace.Init(nil)
	defer trace.Close()

	//http init
	srv := service.New(conf.Conf)
	http.Init(conf.Conf, srv)

	//signal notify to change service behavior
	sch := make(chan os.Signal, 1)
	signal.Notify(sch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-sch
		log.Info("videoup-task-admin got a signal %s", s.String())

		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Info("videoup-task-admin is closed")
			srv.Close()
			time.Sleep(time.Second * 1)
			return
		case syscall.SIGHUP:
			//reload
		default:
			return
		}
	}
}
