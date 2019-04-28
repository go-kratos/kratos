package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/service/main/rank/conf"
	rpc "go-common/app/service/main/rank/server/gorpc"
	"go-common/app/service/main/rank/server/http"
	"go-common/app/service/main/rank/service"
	"go-common/library/conf/env"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	"go-common/library/naming"
	"go-common/library/naming/discovery"
	xip "go-common/library/net/ip"
	"go-common/library/net/trace"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Log)
	defer log.Close()
	log.Info("start")
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	ecode.Init(conf.Conf.Ecode)
	// int service
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
			AppID: env.AppID,
			Addrs: []string{
				"http://" + ip + ":" + env.HTTPPort,
				"gorpc://" + ip + ":" + env.GORPCPort,
			},
		}
		cancel, err = dis.Register(context.Background(), ins)
		if err != nil {
			panic(err)
		}
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			if cancel != nil {
				cancel()
			}
			rpcSvr.Close()
			time.Sleep(time.Second * 2)
			svr.Close()
			log.Info("rank-service exit")
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
