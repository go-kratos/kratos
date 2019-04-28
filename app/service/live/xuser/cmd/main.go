package main

import (
	"context"
	"flag"
	"fmt"
	"go-common/app/service/live/xuser/server/grpc"
	"go-common/library/net/trace"
	"go-common/library/queue/databus/report"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/service/live/xuser/conf"
	"go-common/app/service/live/xuser/server/http"
	"go-common/app/service/live/xuser/service"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Log)
	defer log.Close()
	log.Info("xuser-service start")
	if conf.TraceInit {
		trace.Init(conf.Conf.Tracer)
		defer trace.Close()
	}
	report.InitUser(conf.Conf.Report)
	ecode.Init(conf.Conf.Ecode)
	svc := service.New(conf.Conf)
	http.Init(conf.Conf, svc)

	// start grpc server
	svr, err := grpc.New(svc)
	if err != nil {
		panic(fmt.Sprintf("start xuser grpc server fail! %s", err))
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			svc.Close()
			if svr != nil {
				svr.Shutdown(context.Background())
			}
			log.Info("xuser-service exit")
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
