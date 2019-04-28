package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"go-common/app/interface/main/broadcast/conf"
	"go-common/app/interface/main/broadcast/server"
	wardenServer "go-common/app/interface/main/broadcast/server/grpc"
	"go-common/app/interface/main/broadcast/server/http"
	"go-common/library/conf/env"
	"go-common/library/conf/paladin"
	"go-common/library/log"
	"go-common/library/naming"
	"go-common/library/naming/discovery"
	"go-common/library/net/ip"
	"go-common/library/net/rpc/warden/resolver"
)

const (
	ver = "v1.2.0"
)

func main() {
	flag.Parse()
	if err := paladin.Init(); err != nil {
		panic(err)
	}
	if err := paladin.Watch("broadcast.toml", conf.Conf); err != nil {
		panic(err)
	}
	runtime.GOMAXPROCS(conf.Conf.Broadcast.MaxProc)
	log.Init(conf.Conf.Log)
	defer log.Close()
	log.Info("broadcast %s start", ver)
	rand.Seed(time.Now().UTC().UnixNano())
	dis := discovery.New(conf.Conf.Discovery)
	resolver.Register(dis)
	// server
	srv := server.NewServer(conf.Conf)
	if err := server.InitWhitelist(conf.Conf.Whitelist); err != nil {
		panic(err)
	}
	if err := server.InitTCP(srv, conf.Conf.TCP.Bind, conf.Conf.Broadcast.MaxProc); err != nil {
		panic(err)
	}
	if err := server.InitWebsocket(srv, conf.Conf.WebSocket.Bind, conf.Conf.Broadcast.MaxProc); err != nil {
		panic(err)
	}
	if conf.Conf.WebSocket.TLSOpen {
		if err := server.InitWebsocketWithTLS(srv, conf.Conf.WebSocket.TLSBind, conf.Conf.WebSocket.CertFile, conf.Conf.WebSocket.PrivateFile, runtime.NumCPU()); err != nil {
			panic(err)
		}
	}
	// v1
	if conf.Conf.Broadcast.OpenPortV1 {
		if err := server.InitTCPV1(srv, conf.Conf.TCP.BindV1, conf.Conf.Broadcast.MaxProc); err != nil {
			panic(err)
		}
		if err := server.InitWebsocketV1(srv, conf.Conf.WebSocket.BindV1, conf.Conf.Broadcast.MaxProc); err != nil {
			panic(err)
		}
		if conf.Conf.WebSocket.TLSOpen {
			if err := server.InitWebsocketWithTLSV1(srv, conf.Conf.WebSocket.TLSBindV1, conf.Conf.WebSocket.CertFile, conf.Conf.WebSocket.PrivateFile, runtime.NumCPU()); err != nil {
				panic(err)
			}
		}
	}
	// http & grpc
	http.Init(conf.Conf)
	rpcSrv := wardenServer.New(conf.Conf, srv)
	var (
		err       error
		disCancel context.CancelFunc
	)
	if env.IP == "" {
		addr := ip.InternalIP()
		_, port, _ := net.SplitHostPort(conf.Conf.WardenServer.Addr)
		ins := &naming.Instance{
			Region:   env.Region,
			Zone:     env.Zone,
			Env:      env.DeployEnv,
			Hostname: env.Hostname,
			Status:   1,
			AppID:    "push.interface.broadcast",
			Addrs: []string{
				"grpc://" + addr + ":" + port,
			},
			Metadata: map[string]string{
				"weight":      os.Getenv("WEIGHT"),
				"ip_addrs":    os.Getenv("IP_ADDRS"),
				"ip_addrs_v6": os.Getenv("IP_ADDRS_V6"),
				"offline":     os.Getenv("OFFLINE"),
				"overseas":    os.Getenv("OVERSEAS"),
			},
		}
		disCancel, err = dis.Register(context.Background(), ins)
		if err != nil {
			panic(err)
		}
		go func() {
			for {
				var (
					err     error
					conns   int
					ips     = make(map[string]struct{})
					roomIPs = make(map[string]struct{})
				)
				for _, bucket := range srv.Buckets() {
					for ip := range bucket.IPCount() {
						ips[ip] = struct{}{}
					}
					for ip := range bucket.RoomIPCount() {
						roomIPs[ip] = struct{}{}
					}
					conns += bucket.ChannelCount()
				}
				ins.Metadata["conns"] = fmt.Sprint(conns)
				ins.Metadata["ips"] = fmt.Sprint(len(ips))
				ins.Metadata["room_ips"] = fmt.Sprint(len(roomIPs))
				if err = dis.Set(ins); err != nil {
					log.Error("dis.Set(%+v) error(%v)", ins, err)
					time.Sleep(time.Second)
					continue
				}
				time.Sleep(time.Duration(conf.Conf.Broadcast.ServerTick))
			}
		}()
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("broadcast get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Info("broadcast %s exit", ver)
			if disCancel != nil {
				disCancel()
			}
			srv.Close()
			rpcSrv.Shutdown(context.Background())
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
