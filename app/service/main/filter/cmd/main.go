package main

import (
	"context"
	"flag"
	"os"
	"runtime"
	"time"

	"go-common/app/service/main/filter/conf"
	"go-common/app/service/main/filter/http"
	rpc "go-common/app/service/main/filter/rpc/server"
	"go-common/app/service/main/filter/server/grpc"
	"go-common/app/service/main/filter/service"
	"go-common/library/conf/env"
	"go-common/library/log"
	"go-common/library/naming"
	"go-common/library/naming/discovery"
	xip "go-common/library/net/ip"
	"go-common/library/net/trace"
	"go-common/library/os/signal"
	"go-common/library/syscall"
)

func main() {
	flag.Parse()
	err := conf.Init()
	if err != nil {
		panic(err)
	}
	runtime.GOMAXPROCS(runtime.NumCPU() * conf.Conf.Property.GoMaxProce)
	// init log
	log.Init(conf.Conf.Log)
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	defer log.Close()
	log.Info("filter-service start")
	svc := service.New()
	rpcSvr := rpc.New(conf.Conf, svc)
	ws := grpc.New(conf.Conf.WardenServer, svc)
	http.Init(svc)
	// start discovery register
	var cancel context.CancelFunc
	if env.IP == "" {
		ip := xip.InternalIP()
		hn, _ := os.Hostname()
		dis := discovery.New(nil)
		ins := &naming.Instance{
			Zone:     env.Zone,
			Env:      env.DeployEnv,
			AppID:    "filter.service",
			Hostname: hn,
			Addrs: []string{
				"http://" + ip + ":" + env.HTTPPort,
				"gorpc://" + ip + ":" + env.GORPCPort,
				"grpc://" + ip + ":" + env.GRPCPort,
			},
		}
		if cancel, err = dis.Register(context.Background(), ins); err != nil {
			panic(err)
		}
	}
	// end discovery register

	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSTOP)
	for {
		s := <-c
		log.Info("filter-service get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			if cancel != nil {
				cancel()
			}
			rpcSvr.Close()
			ws.Shutdown(context.Background())
			time.Sleep(2 * time.Second)
			log.Info("filter-service exit")
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
