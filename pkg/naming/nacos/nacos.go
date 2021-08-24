package nacos

import (
	"context"
	"errors"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-kratos/kratos/pkg/log"
	"github.com/go-kratos/kratos/pkg/naming"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

const (
	_statusUP = "1"
	//Time units is second
	registerTTL        = 90
	defaultDialTimeout = 30

	NacosMetaRegion  = "region"
	NacosMetaZone    = "zone"
	NacosMetaEnv     = "env"
	NacosMetaColor   = "color"
	NacosMetaScheme  = "scheme"
	NacosMetaVersion = "version"
)

// Config is nacos config.
type Config struct {
	ServerConfigs []constant.ServerConfig
	ClientConfig  constant.ClientConfig
}

var (
	_once    sync.Once
	_builder *Nacos

	// ErrDuplication is a register duplication err
	ErrDuplication = errors.New("nacos: instance duplicate registration")
)

// Builder return default nacos resolver builder.
func Builder(c *Config) *Nacos {
	_once.Do(func() {
		_builder, _ = newNacos(c)
	})
	return _builder
}

// Build register resolver into default nacos.
func Build(c *Config, id string, options ...naming.BuildOpt) naming.Resolver {
	return Builder(c).Build(id, options...)
}

func Register(c *Config, ctx context.Context, ins *naming.Instance) (cancelFunc context.CancelFunc, err error) {
	return Builder(c).Register(ctx, ins)
}

type appInfo struct {
	resolver map[*Resolve]struct{}
	ins      atomic.Value
	nab      *Nacos
	once     sync.Once
}

// Resolve nacos resolver.
type Resolve struct {
	id    string
	event chan struct{}
	nab   *Nacos
	opt   *naming.BuildOptions
}

// Nacos is a nacos client Builder.
// path: /{root}/{appid}/{ip} -> json(instance)
type Nacos struct {
	c          *Config
	cli        naming_client.INamingClient
	ctx        context.Context
	cancelFunc context.CancelFunc

	mutex    sync.RWMutex
	apps     map[string]*appInfo
	registry map[string]struct{}
}

// New is new a nacos builder.
func newNacos(c *Config) (nab *Nacos, err error) {
	naConn, err := clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": c.ServerConfigs,
		"clientConfig":  c.ClientConfig,
	})
	if err != nil {
		errInfo := fmt.Sprintf("nacos client create failed")
		log.Error(errInfo)
		return nil, errors.New(errInfo)
	}

	log.Info(fmt.Sprintf("nacos create ok!"))

	ctx, cancel := context.WithCancel(context.Background())
	nab = &Nacos{
		c:          c,
		cli:        naConn,
		ctx:        ctx,
		cancelFunc: cancel,
		apps:       map[string]*appInfo{},
		registry:   map[string]struct{}{},
	}
	return
}

// Build nacos resovler builder.
func (z *Nacos) Build(appid string, options ...naming.BuildOpt) naming.Resolver {
	r := &Resolve{
		id:    appid,
		nab:   z,
		event: make(chan struct{}, 1),
		opt:   new(naming.BuildOptions),
	}
	for _, opt := range options {
		opt.Apply(r.opt)
	}
	z.mutex.Lock()
	app, ok := z.apps[appid]
	if !ok {
		app = &appInfo{
			resolver: make(map[*Resolve]struct{}),
			nab:      z,
		}
		z.apps[appid] = app
	}
	app.resolver[r] = struct{}{}
	z.mutex.Unlock()
	if ok {
		select {
		case r.event <- struct{}{}:
		default:
		}
	}
	app.once.Do(func() {
		go app.watch(appid)
	})
	return r
}

// Scheme return nacos's scheme.
func (z *Nacos) Scheme() string {
	return "nacos"
}

// Register is register instance.
func (z *Nacos) Register(ctx context.Context, ins *naming.Instance) (cancelFunc context.CancelFunc, err error) {
	z.mutex.Lock()
	if _, ok := z.registry[ins.AppID]; ok {
		err = ErrDuplication
	} else {
		z.registry[ins.AppID] = struct{}{}
	}
	z.mutex.Unlock()
	if err != nil {
		return
	}
	ctx, cancel := context.WithCancel(z.ctx)
	if err = z.register(ctx, ins); err != nil {
		z.mutex.Lock()
		delete(z.registry, ins.AppID)
		z.mutex.Unlock()
		cancel()
		return
	}

	ch := make(chan struct{}, 1)
	cancelFunc = context.CancelFunc(func() {
		cancel()
		<-ch
	})
	go func() {
		ticker := time.NewTicker(time.Duration(registerTTL/3) * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err = z.register(ctx, ins); err != nil {
					log.Warn(fmt.Sprintf("watch zkClient state, fail to register node error:(%v)", err))
					continue
				}
				return
			case <-ctx.Done():
				_ = z.unregister(ins)
				ch <- struct{}{}
				return
			}
		}
	}()
	return
}

// register is register instance to nacos.
func (z *Nacos) register(ctx context.Context, ins *naming.Instance) (err error) {
	log.Info("nacos register enter, instance Addrs:(%v)", ins.Addrs)
	for _, addr := range ins.Addrs {
		u, err := url.Parse(addr)
		if err != nil {
			continue
		}
		// grpc://127.0.0.1:8000 to 127.0.0.1
		ip := strings.SplitN(u.Host, ":", 2)[0]
		port, _ := strconv.Atoi(u.Port())

		var weight int64
		if weight, _ = strconv.ParseInt(ins.Metadata[naming.MetaWeight], 10, 64); weight <= 0 {
			weight = 10
		}
		metadata := map[string]string{
			NacosMetaRegion:  ins.Region,
			NacosMetaZone:    ins.Zone,
			NacosMetaEnv:     ins.Env,
			NacosMetaVersion: ins.Version,
			NacosMetaScheme:  u.Scheme,
			NacosMetaColor:   ins.Metadata[naming.MetaColor],
		}
		_, err = z.cli.RegisterInstance(vo.RegisterInstanceParam{
			Ip:          ip,
			Port:        uint64(port),
			Tenant:      "",
			Weight:      float64(weight),
			Enable:      true,
			Healthy:     true,
			Metadata:    metadata,
			ServiceName: ins.AppID,
			Ephemeral:   true,
		})

		if err != nil {
			log.Warn(fmt.Sprintf("registerServer, fail to RegisterPeerServer node:%s error:(%v)", addr, err))
		} else {
			log.Info(fmt.Sprintf("registerServer, succeed to RegistServer node."))
		}
	}
	return nil
}

func (z *Nacos) unregister(ins *naming.Instance) (err error) {
	log.Info("nacos unregister enter, instance Addrs:(%v)", ins.Addrs)
	for _, addr := range ins.Addrs {
		u, err := url.Parse(addr)
		if err != nil {
			continue
		}
		ip := strings.SplitN(u.Host, ":", 2)[0]
		port, _ := strconv.Atoi(u.Port())

		_, err = z.cli.RegisterInstance(vo.RegisterInstanceParam{
			Ip:          ip,
			Port:        uint64(port),
			ServiceName: ins.AppID,
			Ephemeral:   true,
		})
		if err != nil {
			log.Error("nacos.Conn.Delete node:(%s), error:(%v)", addr, err)
			continue
		}
		log.Info(fmt.Sprintf("unregister, client.Delete:(%v), appid:(%v), hostname:(%v) success", addr, ins.AppID, ins.Hostname))
	}
	return
}

// Close stop all running process including nacos fetch and register.
func (z *Nacos) Close() error {
	z.cancelFunc()
	//
	z.mutex.RLock()
	for appID, info := range z.apps {
		param := &vo.SubscribeParam{
			ServiceName:       appID,
			SubscribeCallback: info.subscribeCallback,
		}
		err := z.cli.Unsubscribe(param)
		if err != nil {
			log.Error("Unsubscribe error(%v)", err)
			return err
		}
	}
	z.mutex.RUnlock()
	return nil
}

func (a *appInfo) watch(appID string) {
	_ = a.fetchstore(appID)
	go func() {
		for {
			log.Info(fmt.Sprintf("nacos watch enter, serverName:(%v)", appID))
			param := &vo.SubscribeParam{
				ServiceName:       appID,
				SubscribeCallback: a.subscribeCallback,
			}
			err := a.nab.cli.Subscribe(param)
			if err != nil {
				log.Error("nacos Subscribe fail to watch:%s error:(%v)", appID, err)
				time.Sleep(time.Second)
				_ = a.fetchstore(appID)
				continue
			}
			log.Info(fmt.Sprintf("nacos Subscribe ok, serverName:%s", appID))
			return
		}
	}()
}

func (a *appInfo) subscribeCallback(services []model.SubscribeService, err error) {
	for rs := range a.resolver {
		log.Info(fmt.Sprintf("nacos subcallback,appid:(%s),services:(%v)", rs.id, services))
		_ = a.fetchstore(rs.id)
	}
}

func (a *appInfo) fetchstore(appID string) (err error) {
	param := vo.SelectInstancesParam{
		Clusters:    nil,
		ServiceName: appID,
		GroupName:   "",
		HealthyOnly: true,
	}
	instances, err := a.nab.cli.SelectInstances(param)
	if err != nil {
		log.Error(fmt.Sprintf("fetchstore, fail to get Children of node:(%v), error:(%v)", appID, err))
		return
	}
	log.Info(fmt.Sprintf("fetchstore, ok to get Children of node:(%v), childs:(%v)", appID, instances))
	ins := &naming.InstancesInfo{
		Instances: make(map[string][]*naming.Instance, 0),
	}
	for _, instance := range instances {
		addr := fmt.Sprintf("%s://%s:%d", instance.Metadata[NacosMetaScheme], instance.Ip, instance.Port)
		insMetadata := map[string]string{
			naming.MetaWeight:  strconv.Itoa(int(instance.Weight)),
			naming.MetaColor:   instance.Metadata[NacosMetaColor],
			naming.MetaCluster: instance.ClusterName,
			naming.MetaZone:    instance.Metadata[NacosMetaZone],
		}
		in := &naming.Instance{
			Region:   instance.Metadata["region"],
			Zone:     instance.Metadata["zone"],
			Env:      instance.Metadata["env"],
			AppID:    instance.ServiceName,
			Hostname: "",
			Addrs:    []string{addr},
			Version:  instance.Metadata["version"],
			LastTs:   0,
			Metadata: insMetadata,
			Status: func() int64 {
				if instance.Enable {
					if instance.Healthy {
						return 1
					} else {
						return 2
					}
				} else {
					return 3
				}
				return 4
			}(),
		}
		ins.Instances[in.Zone] = append(ins.Instances[in.Zone], in)

	}
	a.store(ins)
	return nil
}

func (a *appInfo) store(ins *naming.InstancesInfo) {
	a.ins.Store(ins)
	a.nab.mutex.RLock()
	for rs := range a.resolver {
		select {
		case rs.event <- struct{}{}:
		default:
		}
	}
	a.nab.mutex.RUnlock()
}

// Watch watch instance.
func (r *Resolve) Watch() <-chan struct{} {
	return r.event
}

// Fetch fetch resolver instance.
func (r *Resolve) Fetch(ctx context.Context) (ins *naming.InstancesInfo, ok bool) {
	r.nab.mutex.RLock()
	app, ok := r.nab.apps[r.id]
	r.nab.mutex.RUnlock()
	if ok {
		var appIns *naming.InstancesInfo
		appIns, ok = app.ins.Load().(*naming.InstancesInfo)
		if !ok {
			return
		}
		ins = new(naming.InstancesInfo)
		ins.LastTs = appIns.LastTs
		ins.Scheduler = appIns.Scheduler
		if r.opt.Filter != nil {
			ins.Instances = r.opt.Filter(appIns.Instances)
		} else {
			ins.Instances = make(map[string][]*naming.Instance)
			for zone, in := range appIns.Instances {
				ins.Instances[zone] = in
			}
		}
		if r.opt.Scheduler != nil {
			ins.Instances[r.opt.ClientZone] = r.opt.Scheduler(ins)
		}
		if r.opt.Subset != nil && r.opt.SubsetSize != 0 {
			for zone, inss := range ins.Instances {
				ins.Instances[zone] = r.opt.Subset(inss, r.opt.SubsetSize)
			}
		}
	}
	return
}

// Close close resolver.
func (r *Resolve) Close() error {
	r.nab.mutex.Lock()
	if app, ok := r.nab.apps[r.id]; ok && len(app.resolver) != 0 {
		delete(app.resolver, r)
	}
	r.nab.mutex.Unlock()
	return nil
}
