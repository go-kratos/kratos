package agent

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"go-common/app/service/main/bns/agent/backend"
	"go-common/app/service/main/bns/conf"
	"go-common/library/conf/env"
	"go-common/library/log"
	"go-common/library/stat/prom"
)

// Agent easyns agent
type Agent struct {
	// plugable easyns server backend
	backend backend.Backend

	// agent id, now it is os hostname
	agentID string

	// agent cfg
	cfg *conf.Config

	// agent local cache, distributed by region, zone, env
	//caches cache.LocalCaches

	// httpAddrs are the addresses per protocol the HTTP server binds to
	httpAddrs []conf.ProtoAddr

	// httpServers provides the HTTP API on various endpoints
	httpServers []*HTTPServer

	// dnsAddr is the address the DNS server binds to
	dnsAddrs []conf.ProtoAddr

	// dnsServer provides the DNS API
	dnsServers []*DNSServer

	// wgServers is the wait group for all HTTP servers
	wgServers sync.WaitGroup
}

// New agent
func New(c *conf.Config) (*Agent, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("get hostname as agent id failed: %s", err)
	}

	httpAddrs, err := c.HTTPAddrs()
	if err != nil {
		return nil, fmt.Errorf("invalid HTTP bind address: %s", err)
	}

	dnsAddrs, err := c.DNSAddrs()
	if err != nil {
		return nil, fmt.Errorf("invalid DNS bind address: %s", err)
	}
	//	var caches cache.LocalCaches
	a := &Agent{
		cfg:       c,
		agentID:   hostname,
		httpAddrs: httpAddrs,
		dnsAddrs:  dnsAddrs,
	}

	return a, nil
}

// Start agent
func (a *Agent) Start() error {
	// start DNS servers
	if err := a.listenAndServeDNS(); err != nil {
		return err
	}

	// listen HTTP
	httpln, err := a.listenHTTP(a.httpAddrs)
	if err != nil {
		return err
	}

	// initial backend
	a.backend, err = backend.New(a.cfg.Backend.Backend, a.cfg.Backend.Config)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// make sure backend is available
	if err = a.backend.Ping(ctx); err != nil {
		return err
	}

	// start serving http servers
	for _, l := range httpln {
		srv := NewHTTPServer(l.Addr().String(), a)
		if err := a.serveHTTP(l, srv); err != nil {
			return err
		}
		a.httpServers = append(a.httpServers, srv)
	}

	return nil
}

// Query name
func (a *Agent) Query(name string) ([]*backend.Instance, error) {
	target, sel, err := backend.ParseName(name, a.DefaultSel())
	if err != nil {
		log.Error("dns: parse name failed! name: %s, err: %s", name, err)
		return nil, err
	}
	// TODO
	inss, err := a.backend.Query(context.Background(), target, sel, backend.Metadata{})
	if err != nil {
		prom.BusinessErrCount.Incr("bns:query")
	}
	return inss, err
}

func (a *Agent) listenAndServeDNS() error {
	notif := make(chan conf.ProtoAddr, len(a.dnsAddrs))
	for _, p := range a.dnsAddrs {
		p := p // capture loop var

		// create server
		svr, err := NewDNSServer(a, a.cfg.DNS)
		if err != nil {
			return err
		}
		a.dnsServers = append(a.dnsServers, svr)

		// start server
		a.wgServers.Add(1)
		go func() {
			defer a.wgServers.Done()

			err := svr.ListenAndServe(p.Net, p.Addr, func() { notif <- p })
			if err != nil && !strings.Contains(err.Error(), "accept") {
				log.Error("agent: Error starting DNS server %s (%s): %v", p.Addr, p.Net, err)
			}
		}()
	}

	// wait for servers to be up
	timeout := time.After(time.Second)
	for range a.dnsAddrs {
		select {
		case p := <-notif:
			log.Info("agent: Started DNS server %s (%s)", p.Addr, p.Net)
			continue
		case <-timeout:
			return fmt.Errorf("agent: timeout starting DNS servers")
		}
	}
	return nil
}

func (a *Agent) listenHTTP(addrs []conf.ProtoAddr) ([]net.Listener, error) {
	var ln []net.Listener
	for _, p := range addrs {
		var l net.Listener
		var err error

		switch {
		case p.Net == "unix":
			l, err = a.listenSocket(p.Addr)
		case p.Net == "tcp" && p.Proto == "http":
			l, err = net.Listen("tcp", p.Addr)
		default:
			return nil, fmt.Errorf("%s:%s listener not supported", p.Net, p.Proto)
		}

		if err != nil {
			for _, l := range ln {
				l.Close()
			}
			return nil, err
		}

		if tcpl, ok := l.(*net.TCPListener); ok {
			l = &tcpKeepAliveListener{tcpl}
		}

		ln = append(ln, l)
	}
	return ln, nil
}

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by NewHttpServer so dead TCP connections
// eventually go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(30 * time.Second)
	return tc, nil
}

func (a *Agent) listenSocket(path string) (net.Listener, error) {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		log.Warn("agent: Replacing socket %q\n", path)
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("error removing socket file: %s", err)
	}
	l, err := net.Listen("unix", path)
	if err != nil {
		return nil, err
	}

	// TODO: set file permissions.
	return l, nil
}

func (a *Agent) serveHTTP(l net.Listener, srv *HTTPServer) error {
	srv.proto = "http"
	notif := make(chan string)
	a.wgServers.Add(1)
	go func() {
		defer a.wgServers.Done()
		notif <- srv.Addr
		err := srv.Serve(l)
		if err != nil && err != http.ErrServerClosed {
			log.Error("agent: Error starting http service: %s\n", err.Error())
		}
	}()

	select {
	case addr := <-notif:
		log.Info("agent: Started HTTP Server on %s\n", addr)
		return nil
	case <-time.After(time.Second):
		return fmt.Errorf("agent: timeout starting HTTP servers")
	}
}

// ShutDown agent
func (a *Agent) ShutDown(ctx context.Context) error {
	errmsg := make([]string, 0)
	for _, hsrv := range a.httpServers {
		if err := hsrv.Shutdown(ctx); err != nil {
			errmsg = append(errmsg, err.Error())
		}
	}
	for _, dsrv := range a.dnsServers {
		if err := dsrv.Shutdown(); err != nil {
			errmsg = append(errmsg, err.Error())
		}
	}
	err := a.backend.Close(ctx)
	if err != nil {
		errmsg = append(errmsg, err.Error())
	}
	if len(errmsg) > 0 {
		return fmt.Errorf("%s", strings.Join(errmsg, "\n"))
	}
	return nil
}

// DefaultSel default selector from config
func (a *Agent) DefaultSel() backend.Selector {
	return backend.Selector{
		Env:    env.DeployEnv,
		Region: env.Region,
		Zone:   env.Zone,
	}
}
