package zookeeper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bilibili/kratos/pkg/log"
	"github.com/bilibili/kratos/pkg/naming"
	"github.com/go-zookeeper/zk"
)

type Config struct {
	// Endpoints is a list of URLs.
	Endpoints []string `json:"endpoints"`
}

var (
	_once    sync.Once
	_builder naming.Builder

	//ErrDuplication is a register duplication err
	ErrDuplication = errors.New("zookeeper: instance duplicate registration")
)

// Builder return default zookeeper resolver builder.
func Builder(c *Config) naming.Builder {
	_once.Do(func() {
		_builder, _ = New(c)
	})
	return _builder
}

// Build register resolver into default zookeeper.
func Build(c *Config, id string) naming.Resolver {
	return Builder(c).Build(id)
}

// ZookeeperBuilder is a zookeeper client Builder
type ZookeeperBuilder struct {
	cli        *zk.Conn
	connEvent  <-chan zk.Event
	ctx        context.Context
	cancelFunc context.CancelFunc

	mutex    sync.RWMutex
	apps     map[string]*appInfo
	registry map[string]struct{}
}

type appInfo struct {
	resolver map[*Resolve]struct{}
	ins      atomic.Value
	zkb      *ZookeeperBuilder
	once     sync.Once
}

// Resolve zookeeper resolver.
type Resolve struct {
	id    string
	event chan struct{}
	zkb   *ZookeeperBuilder
}

// New is new a zookeeper builder
func New(c *Config) (zkb *ZookeeperBuilder, err error) {
	//example: endpointSli = []string{"192.168.1.78:2181", "192.168.1.79:2181", "192.168.1.80:2181"}
	if len(c.Endpoints) == 0 {
		errInfo := fmt.Sprintf("zookeeper New failed, endpoints is null")
		log.Error(errInfo)
		return nil, errors.New(errInfo)
	}

	zkConn, connEvent, err := zk.Connect(c.Endpoints, 5*time.Second)
	if err != nil {
		log.Error(fmt.Sprintf("zk Connect err:(%v)", err))
		return
	} else {
		log.Info(fmt.Sprintf("zk Connect ok!"))
	}

	ctx, cancel := context.WithCancel(context.Background())
	zkb = &ZookeeperBuilder{
		cli:        zkConn,
		connEvent:  connEvent,
		ctx:        ctx,
		cancelFunc: cancel,
		apps:       map[string]*appInfo{},
		registry:   map[string]struct{}{},
	}
	return
}

// Build zookeeper resovler builder.
func (z *ZookeeperBuilder) Build(appid string) naming.Resolver {
	r := &Resolve{
		id:    appid,
		zkb:   z,
		event: make(chan struct{}, 1),
	}
	z.mutex.Lock()
	app, ok := z.apps[appid]
	if !ok {
		app = &appInfo{
			resolver: make(map[*Resolve]struct{}),
			zkb:      z,
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
		log.Info("zookeeper: AddWatch(%s) already watch(%v)", appid, ok)
	})
	return r
}

// Scheme return zookeeper's scheme
func (z *ZookeeperBuilder) Scheme() string {
	return "zookeeper"
}

// Register is register instance
func (z *ZookeeperBuilder) Register(ctx context.Context, ins *naming.Instance) (cancelFunc context.CancelFunc, err error) {
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
		for {
			select {
			case connEvent := <-z.connEvent:
				log.Warn("watch zkClient state, connEvent:(%v)", connEvent)
				if connEvent.State == zk.StateHasSession {
					log.Warn("watch zkClient state, state is StateHasSession...")
					err = z.register(ctx, ins)
					if err != nil {
						log.Warn(fmt.Sprintf("watch zkClient state, fail to register node error:(%v)", err))
						continue
					}
				}
			case <-ctx.Done():
				ch <- struct{}{}
				return
			}
		}
	}()
	return
}

func (z *ZookeeperBuilder) registerPerServer(name string) (err error) {
	var (
		str string
	)

	str, err = z.cli.Create(name, nil, 0, zk.WorldACL(zk.PermAll))
	if err != nil {
		log.Warn(fmt.Sprintf("registerPerServer, fail to Create node:(%s). err:(%v)", name, err))
	} else {
		log.Info(fmt.Sprintf("registerPerServer, succeed to Create node:(%s). retStr:(%s)", name, str))
	}

	return
}

func (z *ZookeeperBuilder) registerEphServer(name, host string, ins *naming.Instance) (err error) {
	var (
		str string
	)

	val, _ := json.Marshal(ins)
	log.Info(fmt.Sprintf("registerEphServer, ins after json.Marshal:(%v)", string(val)))

	str, err = z.cli.Create(name+host, val, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err != nil {
		log.Warn(fmt.Sprintf("registerEphServer, fail to Create node:%s. err:(%v)", name+host, err))
	} else {
		log.Info(fmt.Sprintf("registerEphServer, succeed to Create node:%s. retStr:(%s)", name+host, str))
	}

	return
}

// register 注册zookeeper节点
func (z *ZookeeperBuilder) register(ctx context.Context, ins *naming.Instance) (err error) {
	log.Info("zookeeper register enter, instance Addrs:(%v)", ins.Addrs)
	prefix := z.keyPrefix(ins)

	err = z.registerPerServer(prefix)
	if err != nil {
		log.Warn(fmt.Sprintf("register, fail to registerPerServer node error:(%v)", err))
	}

	for _, val := range ins.Addrs {
		err = z.registerEphServer(prefix, "/"+val, ins)
		if err != nil {
			log.Warn(fmt.Sprintf("registerServer, fail to RegisterEphServer node error:(%v)", err))
		} else {
			log.Info(fmt.Sprintf("registerServer, succeed to RegistServer node."))
		}
	}

	return nil
}

// unregister 删除zookeeper中节点信息
func (z *ZookeeperBuilder) unregister(ins *naming.Instance) (err error) {
	log.Info("zookeeper unregister enter, instance Addrs:(%v)", ins.Addrs)
	prefix := z.keyPrefix(ins)

	for _, val := range ins.Addrs {
		strNode := prefix + "/" + val
		exists, _, err := z.cli.Exists(strNode)
		if err != nil {
			log.Error("zk.Conn.Exists node:(%v), error:(%s)", strNode, err.Error())
			return err
		}
		if exists {
			_, s, err := z.cli.Get(strNode)
			if err != nil {
				log.Error("zk.Conn.Get node:(%s), error:(%s)", strNode, err.Error())
				return err
			}
			return z.cli.Delete(strNode, s.Version)
		}

		log.Info(fmt.Sprintf("unregister, client.Delete:(%v), appid:(%v), hostname:(%v) success", strNode, ins.AppID, ins.Hostname))
	}

	return
}

func (z *ZookeeperBuilder) keyPrefix(ins *naming.Instance) string {
	return fmt.Sprintf("/%s", ins.AppID)
}

// Close stop all running process including zk fetch and register
func (z *ZookeeperBuilder) Close() error {
	z.cancelFunc()
	return nil
}

func (a *appInfo) watch(appID string) {
	_ = a.fetchstore(appID)
	prefix := fmt.Sprintf("/%s", appID)

	go func() {
		for {
			log.Info(fmt.Sprintf("zk ChildrenW enter, prefix:(%v)", prefix))
			snapshot, _, event, err := a.zkb.cli.ChildrenW(prefix)
			if err != nil {
				continue
			}

			log.Info(fmt.Sprintf("zk ChildrenW ok, snapshot:(%v)", snapshot))
			for ev := range event {
				log.Info(fmt.Sprintf("zk ChildrenW ok, prefix:(%v), event Path:(%v), Type:(%v)", prefix, ev.Path, ev.Type))
				if ev.Type == zk.EventNodeChildrenChanged {
					_ = a.fetchstore(appID)
				}
			}
		}
	}()
}

func (a *appInfo) fetchstore(appID string) (err error) {
	prefix := fmt.Sprintf("/%s", appID)
	strNode := ""
	childs, _, err := a.zkb.cli.Children(prefix)
	if err != nil {
		log.Error(fmt.Sprintf("fetchstore, fail to get Children of node:(%v), err:(%v)", prefix, err))
	} else {
		log.Info(fmt.Sprintf("fetchstore, ok to get Children of node:(%v), childs:(%v)", prefix, childs))
	}

	ins := &naming.InstancesInfo{
		Instances: make(map[string][]*naming.Instance, 0),
	}

	//for range childs
	for _, child := range childs {
		strNode = prefix + "/" + child
		resp, _, err := a.zkb.cli.Get(strNode)
		if err != nil {
			log.Error("zookeeper: fetch client.Get(%s) error(%v)", strNode, err)
			return err
		}

		in := new(naming.Instance)
		err = json.Unmarshal(resp, in)
		if err != nil {
			return err
		}
		ins.Instances[in.Zone] = append(ins.Instances[in.Zone], in)

	}
	a.store(ins)

	return nil
}

func (a *appInfo) store(ins *naming.InstancesInfo) {

	a.ins.Store(ins)
	a.zkb.mutex.RLock()
	for rs := range a.resolver {
		select {
		case rs.event <- struct{}{}:
		default:
		}
	}
	a.zkb.mutex.RUnlock()
}

// Watch watch instance.
func (r *Resolve) Watch() <-chan struct{} {
	return r.event
}

// Fetch fetch resolver instance.
func (r *Resolve) Fetch(ctx context.Context) (ins *naming.InstancesInfo, ok bool) {
	r.zkb.mutex.RLock()
	app, ok := r.zkb.apps[r.id]
	r.zkb.mutex.RUnlock()
	if ok {
		ins, ok = app.ins.Load().(*naming.InstancesInfo)
		return
	}
	return
}

// Close close resolver.
func (r *Resolve) Close() error {
	r.zkb.mutex.Lock()
	if app, ok := r.zkb.apps[r.id]; ok && len(app.resolver) != 0 {
		delete(app.resolver, r)
	}
	r.zkb.mutex.Unlock()
	return nil
}
