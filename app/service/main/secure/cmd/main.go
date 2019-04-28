package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/service/main/secure/conf"
	"go-common/app/service/main/secure/http"
	rpc "go-common/app/service/main/secure/rpc/server"
	"go-common/app/service/main/secure/service"
	"go-common/library/log"
	xrpc "go-common/library/net/rpc"
	"go-common/library/net/trace"
)

var (
	srv    *service.Service
	rpcSvr *xrpc.Server
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		log.Error("conf.Init() error(%v)", err)
		panic(err)
	}
	log.Init(conf.Conf.Log)
	defer log.Close()
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	log.Info("credit-timer start")
	srv = service.New(conf.Conf)
	rpcSvr = rpc.New(conf.Conf, srv)
	http.Init(srv)
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
			rpcSvr.Close()
			time.Sleep(time.Second)
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
