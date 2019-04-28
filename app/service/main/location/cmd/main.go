package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/service/main/location/conf"
	rpc "go-common/app/service/main/location/rpc/server"
	"go-common/app/service/main/location/server/grpc"
	"go-common/app/service/main/location/server/http"
	"go-common/app/service/main/location/service"
	"go-common/library/log"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/rpc/warden/resolver/livezk"
	"go-common/library/net/trace"
)

const (
	discoveryID = "location.service"
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
	// service init
	svr := service.New(conf.Conf)
	rpcSvr := rpc.New(conf.Conf, svr)
	// warden init
	var err error
	var cancelzk context.CancelFunc = func() {}
	var grpcSvr *warden.Server
	if conf.Conf.WardenServer != nil {
		grpcSvr = grpc.New(conf.Conf.WardenServer, svr)
		if conf.Conf.LiveZK != nil {
			if cancelzk, err = livezk.Register(conf.Conf.LiveZK, conf.Conf.WardenServer.Addr, discoveryID); err != nil {
				panic(err)
			}
		}
	}
	// http init
	http.Init(conf.Conf, svr, rpcSvr)
	log.Info("location-service start")
	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for s := range c {
		log.Info("location-service get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			cancelzk()
			if grpcSvr != nil {
				grpcSvr.Shutdown(context.Background())
			}
			rpcSvr.Close()
			time.Sleep(2 * time.Second)
			log.Info("location-service exit")
			return
		case syscall.SIGHUP:
			// TODO reload
		default:
			return
		}
	}
}
