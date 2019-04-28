package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	server "go-common/app/service/main/broadcast/server/grpc"
	"go-common/app/service/main/broadcast/server/http"
	"go-common/app/service/main/broadcast/service"
	"go-common/library/conf/env"
	"go-common/library/conf/paladin"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	"go-common/library/naming"
	"go-common/library/naming/discovery"
	"go-common/library/net/ip"
)

const (
	ver = "v1.4.4"
)

func main() {
	flag.Parse()
	if err := paladin.Init(); err != nil {
		panic(err)
	}
	var (
		ac struct {
			Discovery *discovery.Config
		}
	)
	if err := paladin.Get("application.toml").UnmarshalTOML(&ac); err != nil {
		if err != paladin.ErrNotExist {
			panic(err)
		}
	}
	log.Init(nil)
	defer log.Close()
	log.Info("broadcast-service %s start", ver)
	// use internal discovery
	dis := discovery.New(ac.Discovery)
	// new a service
	srv := service.New(dis)
	ecode.Init(nil)
	http.Init(srv)
	// grpc server
	rpcSrv, rpcPort := server.New(srv)
	rpcSrv.Start()
	// register discovery
	var (
		err    error
		cancel context.CancelFunc
	)
	if env.IP == "" {
		ipAddr := ip.InternalIP()
		// broadcast discovery
		ins := &naming.Instance{
			Zone:     env.Zone,
			Env:      env.DeployEnv,
			Hostname: env.Hostname,
			AppID:    "push.service.broadcast",
			Addrs: []string{
				"grpc://" + ipAddr + ":" + rpcPort,
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
		log.Info("broadcast-service get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Info("broadcast-service %s exit", ver)
			if cancel != nil {
				cancel()
			}
			rpcSrv.Shutdown(context.Background())
			time.Sleep(time.Second * 2)
			srv.Close()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
