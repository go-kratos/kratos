package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/service/main/favorite/conf"
	"go-common/app/service/main/favorite/server/gorpc"
	gserver "go-common/app/service/main/favorite/server/grpc"
	"go-common/app/service/main/favorite/server/http"
	"go-common/app/service/main/favorite/service"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	"go-common/library/net/trace"
)

func main() {
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
	log.Info("favorite start")
	ecode.Init(conf.Conf.Ecode)
	// service init
	svc := service.New(conf.Conf)
	rpcSvr := rpc.New(conf.Conf, svc)
	grpcSvr := gserver.New(conf.Conf.WardenServer, svc)
	if _, err := grpcSvr.Start(); err != nil {
		panic(err)
	}
	http.Init(conf.Conf, svc)

	// init signal
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		si := <-ch
		switch si {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Info("get a signal %s, stop the consume process", si.String())
			rpcSvr.Close()
			grpcSvr.Shutdown(context.Background())
			time.Sleep(time.Second * 2)
			svc.Close()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
