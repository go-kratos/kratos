package main

import (
	"context"
	"flag"
	"go-common/app/service/live/broadcast-proxy/conf"
	"go-common/app/service/live/broadcast-proxy/server"
	"go-common/app/service/live/broadcast-proxy/service"
	"go-common/library/log"
	"net/http"
	"net/http/pprof"
	"runtime"
	"time"
)

func RunPprofServer(addr string) {
	pprofServer := http.NewServeMux()
	pprofServer.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	pprofServer.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	pprofServer.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	pprofServer.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	pprofServer.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	go func() {
		if e := http.ListenAndServe(addr, pprofServer); e != nil {
			log.Error("pprof server error ListenAndServe addr:%s,error:%+v", addr, e)
		}
		defer func() {
			if e := recover(); e != nil {
				log.Error("expected panic from pprof server,error:%+v", e)
			}
		}()
	}()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var confTomlFile string
	flag.StringVar(&confTomlFile, "conf", "", "config file for broadcast proxy")
	flag.Parse()

	conf, err := conf.NewBroadcastProxyConfig(confTomlFile)
	if err != nil {
		panic(err)
	}

	log.Init(conf.Log)
	log.Info("Broadcast Proxy Service:%s", time.Now().String())
	log.Info("Broadcast Proxy Config:%+v", conf)

	proxy, err := server.NewBroadcastProxy(conf.Backend.BackendServer, conf.Backend.ProbePath,
		conf.Backend.MaxIdleConnsPerHost, conf.Backend.ProbeSample)
	if err != nil {
		panic(err)
	}
	defer proxy.Close()

	dispatcher, err := server.NewCometDispatcher(conf.Ipip, conf.Dispatch, conf.Sven)
	if err != nil {
		panic(err)
	}
	defer dispatcher.Close()

	httpService, err := server.NewBroadcastService(conf.Http.Address, proxy, dispatcher)
	if err != nil {
		panic(err)
	}
	defer httpService.Close()

	grpcService, err := service.NewGrpcService(proxy, dispatcher)
	if err != nil {
		panic(err)
	}
	defer grpcService.Shutdown(context.Background())

	RunPprofServer(conf.Perf)

	quit := make(chan struct{})
	<-quit
}
