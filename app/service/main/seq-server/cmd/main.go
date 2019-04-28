package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/service/main/seq-server/conf"
	"go-common/app/service/main/seq-server/http"
	rpc "go-common/app/service/main/seq-server/rpc/server"
	"go-common/app/service/main/seq-server/service"
	"go-common/library/conf/env"
	"go-common/library/log"
	"go-common/library/naming"
	"go-common/library/naming/discovery"
	xip "go-common/library/net/ip"
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
	defer log.Close()
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	log.Info("seq-server start")
	// service init
	svr := service.New(conf.Conf)
	rpcSvr := rpc.New(conf.Conf, svr)
	http.Init(conf.Conf, svr)
	// start discovery register
	var (
		err    error
		cancel context.CancelFunc
	)
	if env.IP == "" {
		ip := xip.InternalIP()
		dis := discovery.New(nil)
		ins := &naming.Instance{
			Zone:  env.Zone,
			Env:   env.DeployEnv,
			AppID: "seq.server",
			Addrs: []string{
				"http://" + ip + ":" + env.HTTPPort,
				"gorpc://" + ip + ":" + env.GORPCPort,
			},
		}
		if cancel, err = dis.Register(context.Background(), ins); err != nil {
			panic(err)
		}
	}
	// end discovery register
	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("seq-server get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			if cancel != nil {
				cancel()
			}
			rpcSvr.Close()
			time.Sleep(time.Second * 2)
			log.Info("seq-server exit")
			return
		case syscall.SIGHUP:
			// TODO reload
		default:
			return
		}
	}
}
