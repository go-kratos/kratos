package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/service/main/vip/conf"
	"go-common/app/service/main/vip/http"
	rpc "go-common/app/service/main/vip/rpc/server"
	grpc "go-common/app/service/main/vip/server/grpc"
	"go-common/app/service/main/vip/service"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
)

func main() {

	flag.Parse()
	// init conf,log,trace,stat,perf.
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Log)
	defer log.Close()
	// service init
	svr := service.New(conf.Conf)
	//ecode init FIXME
	ecode.Init(conf.Conf.Ecode)
	// rpc server init
	rpcSvr := rpc.New(conf.Conf, svr)
	http.Init(conf.Conf, svr)
	ws := grpc.New(conf.Conf.WardenServer, svr)
	// signal handler
	log.Info("vip-service start")
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("vip-service get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			rpcSvr.Close()
			ws.Shutdown(context.Background())
			time.Sleep(time.Second * 2)
			log.Info("vip-service exit")
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
