package main

import (
	"context"
	"flag"
	"fmt"
	"go-common/app/service/live/xuserex/server/http"
	"go-common/app/service/live/xuserex/service"
	"os"
	"os/signal"
	"syscall"

	"go-common/app/service/live/resource/sdk"
	"go-common/app/service/live/xuserex/conf"
	"go-common/app/service/live/xuserex/server/grpc"
	ecode "go-common/library/ecode/tip"
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
	log.Info("[xuserex] start")
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	ecode.Init(conf.Conf.Ecode)
	http.Init(conf.Conf)
	svc := service.New(conf.Conf)
	svr, err := grpc.New(svc)
	if err != nil {
		panic(fmt.Sprintf("start xuser grpc server fail! %s", err))
	}

	titansSdk.Init(conf.Conf.Titan)
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
			log.Info("[xuserex] exit")
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
