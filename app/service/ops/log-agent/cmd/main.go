package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"
	"context"
	"flag"
	"fmt"
	"net/http"

	"go-common/app/service/ops/log-agent/conf"
	"go-common/library/log"
	"go-common/app/service/ops/log-agent/pipeline"
	"go-common/app/service/ops/log-agent/pipeline/hostlogcollector"
	"go-common/app/service/ops/log-agent/pipeline/dockerlogcollector"
	"go-common/app/service/ops/log-agent/pkg/limit"
	"go-common/app/service/ops/log-agent/pkg/flowmonitor"
	"go-common/app/service/ops/log-agent/pkg/httpstream"
	"go-common/app/service/ops/log-agent/pkg/lancermonitor"
	"go-common/app/service/ops/log-agent/pkg/lancerroute"
	"go-common/library/conf/env"
	xip "go-common/library/net/ip"
	"go-common/library/naming/discovery"
	"go-common/library/naming"
)

const AppVersion = "2.1.0"

type Agent struct {
	limit      *limit.Limit
	httpstream *httpstream.HttpStream
}

func main() {
	var (
		err         error
		ctx, cancel = context.WithCancel(context.Background())
		//cancel context.CancelFunc
	)
	version := flag.Bool("v", false, "show version and exit")
	flag.Parse()
	if *version {
		fmt.Println(AppVersion)
		os.Exit(0)
	}

	if err = conf.Init(); err != nil {
		panic(err)
	}

	agent := new(Agent)

	// set context
	ctx = context.WithValue(ctx, "GlobalConfig", conf.Conf)
	ctx = context.WithValue(ctx, "MetaPath", conf.Conf.HostLogCollector.MetaPath)

	// init xlog
	conf.Conf.Log.Stdout = true // ooooo just for test
	log.Init(conf.Conf.Log)
	defer log.Close()
	log.Info("log agent [version: %s] start", AppVersion)

	// resource limit by cgroup
	if conf.Conf.Limit != nil {
		conf.Conf.Limit.AppName = "log-agent"
		if agent.limit, err = limit.LimitRes(conf.Conf.Limit); err != nil {
			log.Warn("resource limit disabled: %s", err)
		}
	}

	// /debug/vars
	if conf.Conf.DebugAddr != "" {
		go http.ListenAndServe(conf.Conf.DebugAddr, nil)
	}

	// init lancer route
	err = lancerroute.InitLancerRoute()
	if err != nil {
		panic(err)
	}

	// httpstream
	if agent.httpstream, err = httpstream.NewHttpStream(conf.Conf.HttpStream); err != nil {
		log.Warn("httpstream disabled: %s", err)
	}

	// start pipeline management
	err = pipeline.InitPipelineMng(ctx)
	if err != nil {
		panic(err)
	}

	// start host log collector
	err = hostlogcollector.InitHostLogCollector(ctx, conf.Conf.HostLogCollector)
	if err != nil {
		panic(err)
	}

	// start docker log collector
	err = dockerlogcollector.InitDockerLogCollector(ctx, conf.Conf.DockerLogCollector)
	if err != nil {
		log.Error("failed to start docker log collector: %s", err)
	}

	// flow monitor
	if conf.Conf.Flowmonitor != nil {
		if err = flowmonitor.InitFlowMonitor(conf.Conf.Flowmonitor); err != nil {
			panic(err)
		}
	}

	// lancer monitor
	if conf.Conf.LancerMonitor != nil {
		if _, err = lancermonitor.InitLancerMonitor(conf.Conf.LancerMonitor); err != nil {
			panic(err)
		}
	}

	// start discovery register
	if env.IP == "" {
		ip := xip.InternalIP()
		dis := discovery.New(conf.Conf.Discovery)
		ins := &naming.Instance{
			Zone:     env.Zone,
			Env:      env.DeployEnv,
			AppID:    env.AppID,
			Hostname: env.Hostname,
			Version:  AppVersion,
			Addrs: []string{
				ip,
			},
		}
		_, err = dis.Register(ctx, ins)
		if err != nil {
			panic(err)
		}
	}

	// signal
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSTOP)
	for {
		s := <-ch
		log.Info("agent get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT, syscall.SIGHUP:
			if cancel != nil {
				cancel()
			}
			time.Sleep(time.Second)
			return
		default:
			return
		}
	}
}
