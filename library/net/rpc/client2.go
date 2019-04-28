package rpc

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"sync/atomic"
	"time"

	"go-common/library/conf/env"
	"go-common/library/log"
	"go-common/library/naming"
	"go-common/library/naming/discovery"
	"go-common/library/net/netutil/breaker"
	xtime "go-common/library/time"
)

const (
	scheme          = "gorpc"
	_policySharding = "sharding"
)

// ClientConfig rpc client config.
type ClientConfig struct {
	Policy  string
	Zone    string
	Cluster string
	Color   string
	Timeout xtime.Duration
	Breaker *breaker.Config
}

// Client2 support for load balancing and service discovery.
type Client2 struct {
	c        *ClientConfig
	appID    string
	dis      naming.Resolver
	balancer atomic.Value
}

// NewDiscoveryCli new discovery client.
func NewDiscoveryCli(appID string, cf *ClientConfig) (c *Client2) {
	if cf == nil {
		cf = &ClientConfig{Timeout: xtime.Duration(300 * time.Millisecond)}
	} else if cf.Timeout <= 0 {
		cf.Timeout = xtime.Duration(300 * time.Millisecond)
	}
	c = &Client2{
		c:     cf,
		appID: appID,
		dis:   discovery.Build(appID),
	}
	var pools = make(map[string]*Client)
	fmt.Printf("开始创建：%s 的gorpc client，等待从discovery拉取节点：%s\n", c.appID, time.Now().Format("2006-01-02 15:04:05"))
	event := c.dis.Watch()
	select {
	case _, ok := <-event:
		if ok {
			c.disc(pools)
			fmt.Printf("结束创建：%s 的gorpc client，从discovery拉取节点和创建成功：%s\n", c.appID, time.Now().Format("2006-01-02 15:04:05"))
		} else {
			panic("刚启动就从discovery拉到了关闭的event")
		}
	case <-time.After(10 * time.Second):
		fmt.Printf("失败创建：%s 的gorpc client，竟然从discovery拉取节点超时了：%s\n", c.appID, time.Now().Format("2006-01-02 15:04:05"))
		if env.DeployEnv == env.DeployEnvProd {
			panic("刚启动就从discovery拉节点超时，请检查配置或联系Discovery维护者")
		}
	}
	go c.discproc(event, pools)
	return
}

// Boardcast boardcast all rpc client.
func (c *Client2) Boardcast(ctx context.Context, serviceMethod string, args interface{}, reply interface{}) (err error) {
	var (
		ok bool
		b  balancer
	)
	if b, ok = c.balancer.Load().(balancer); ok {
		if err = b.Boardcast(ctx, serviceMethod, args, reply); err != ErrNoClient {
			return
		}
	}
	return nil
}

// Call invokes the named function, waits for it to complete, and returns its error status.
// this include rpc.Client.Call method, and takes a timeout.
func (c *Client2) Call(ctx context.Context, serviceMethod string, args interface{}, reply interface{}) (err error) {
	var (
		ok bool
		b  balancer
	)
	if b, ok = c.balancer.Load().(balancer); ok {
		if err = b.Call(ctx, serviceMethod, args, reply); err != ErrNoClient {
			return
		}
	}
	stats.Incr(serviceMethod, "no_rpc_client")
	return ErrNoClient
}

func (c *Client2) removeAndClose(pools, dcs map[string]*Client) {
	if len(dcs) == 0 {
		return
	}
	// after rpc timeout(double duration), close no used cliens
	if c.c != nil {
		to := c.c.Timeout
		time.Sleep(2 * time.Duration(to))
	}
	for key, cli := range dcs {
		delete(pools, key)
		cli.Close()
	}
}

func (c *Client2) discproc(event <-chan struct{}, pools map[string]*Client) {
	for {
		if _, ok := <-event; ok {
			c.disc(pools)
			continue
		}
		return
	}
}

func (c *Client2) disc(pools map[string]*Client) (err error) {
	var (
		weights   int64
		key       string
		i, j, idx int
		nodes     map[string]struct{}
		dcs       map[string]*Client
		blc       balancer
		cli       *Client
		cs, wcs   []*Client
		svr       *naming.Instance
	)
	insMap, ok := c.dis.Fetch(context.Background())
	if !ok {
		log.Error("discovery fetch instance fail(%s)", c.appID)
		return
	}
	zone := env.Zone
	if c.c.Zone != "" {
		zone = c.c.Zone
	}
	tinstance, ok := insMap[zone]
	if !ok {
		for _, value := range insMap {
			tinstance = value
			break
		}
	}
	instance := make([]*naming.Instance, 0, len(tinstance))
	for _, svr := range tinstance {
		nsvr := new(naming.Instance)
		*nsvr = *svr
		cluster := svr.Metadata[naming.MetaCluster]
		if c.c.Cluster != "" && c.c.Cluster != cluster {
			continue
		}
		instance = append(instance, nsvr)
	}
	log.Info("discovery get  %d instances ", len(instance))
	if len(instance) > 0 {
		nodes = make(map[string]struct{}, len(instance))
		cs = make([]*Client, 0, len(instance))
		dcs = make(map[string]*Client, len(pools))
		svrWeights := make([]int, 0, len(instance))
		weights = 0
		for _, svr = range instance {
			weight, err := strconv.ParseInt(svr.Metadata["weight"], 10, 64)
			if err != nil {
				weight = 10
			}
			key = svr.Hostname
			nodes[key] = struct{}{}

			var addr string
			if cli, ok = pools[key]; !ok {
				for _, saddr := range svr.Addrs {
					u, err := url.Parse(saddr)
					if err == nil && u.Scheme == scheme {
						addr = u.Host
					}
				}
				if addr == "" {
					log.Warn("net/rpc: invalid rpc address(%s,%s,%v) found!", svr.AppID, svr.Hostname, svr.Addrs)
					continue
				}
				cli = Dial(addr, c.c.Timeout, c.c.Breaker)
				pools[key] = cli
			}
			svrWeights = append(svrWeights, int(weight))
			weights += weight // calc all weight
			log.Info("new cli %+v instance info %+v", addr, svr)
			cs = append(cs, cli)
		}
		// delete old nodes
		for key, cli = range pools {
			if _, ok = nodes[key]; !ok {
				log.Info("syncproc will delete node: %s", key)
				dcs[key] = cli
			}
		}
		// new client slice by weights
		wcs = make([]*Client, 0, weights)
		for i, j = 0, 0; i < int(weights); j++ { // j++ means next svr
			idx = j % len(cs)
			if svrWeights[idx] > 0 {
				i++ // i++ means all weights must fill wrrClis
				svrWeights[idx]--
				wcs = append(wcs, cs[idx])
			}
		}
		switch c.c.Policy {
		case _policySharding:
			blc = &sharding{
				pool:   wcs,
				weight: int64(weights),
				server: int64(len(instance)),
			}
			log.Info("discovery syncproc sharding weights:%d size:%d raw:%d", weights, weights, len(instance))
		default:
			blc = &wrr{
				pool:   wcs,
				weight: int64(weights),
				server: int64(len(instance)),
			}
			log.Info("discovery %s syncproc wrr weights:%d size:%d raw:%d", c.appID, weights, weights, len(instance))
		}
		c.balancer.Store(blc)
		c.removeAndClose(pools, dcs)
	}
	return
}
