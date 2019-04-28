package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"go-common/app/infra/config/conf"
	"go-common/app/infra/config/http"
	"go-common/app/infra/config/rpc/server"
	"go-common/app/infra/config/service/v1"
	"go-common/app/infra/config/service/v2"
	"go-common/library/conf/env"
	"go-common/library/log"
	"go-common/library/naming"
	"go-common/library/naming/discovery"
	xip "go-common/library/net/ip"
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
	// service init
	svr2 := v2.New(conf.Conf)
	svr := v1.New(conf.Conf)
	rpcSvr := rpc.New(conf.Conf, svr, svr2)
	http.Init(conf.Conf, svr, svr2, rpcSvr)
	// start discovery register
	var (
		err    error
		cancel context.CancelFunc
	)
	if env.IP == "" {
		ip := xip.InternalIP()
		hn, _ := os.Hostname()
		dis := discovery.New(nil)
		ins := &naming.Instance{
			Zone:     env.Zone,
			Env:      env.DeployEnv,
			AppID:    "config.service",
			Hostname: hn,
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
		log.Info("config-service get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			if cancel != nil {
				cancel()
			}
			rpcSvr.Close()
			svr.Close()
			log.Info("config-service exit")
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
