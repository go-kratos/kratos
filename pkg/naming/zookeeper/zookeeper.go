package zookeeper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"path"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bilibili/kratos/pkg/log"
	"github.com/bilibili/kratos/pkg/naming"
	xtime "github.com/bilibili/kratos/pkg/time"
	"github.com/go-zookeeper/zk"
)

// Config is zookeeper config.
type Config struct {
	Root      string         `json:"root"`
	Endpoints []string       `json:"endpoints"`
	Timeout   xtime.Duration `json:"timeout"`
}

var (
	_once    sync.Once
	_builder naming.Builder

	// ErrDuplication is a register duplication err
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

type appInfo struct {
	resolver map[*Resolve]struct{}
	ins      atomic.Value
	zkb      *Zookeeper
	once     sync.Once
}

// Resolve zookeeper resolver.
type Resolve struct {
	id    string
	event chan struct{}
	zkb   *Zookeeper
}

// Zookeeper is a zookeeper client Builder.
// path: /{root}/{appid}/{ip} -> json(instance)
type Zookeeper struct {
	c          *Config
	cli        *zk.Conn
	connEvent  <-chan zk.Event
	ctx        context.Context
	cancelFunc context.CancelFunc

	mutex    sync.RWMutex
	apps     map[string]*appInfo
	registry map[string]struct{}
}

// New is new a zookeeper builder.
func New(c *Config) (zkb *Zookeeper, err error) {
	if c.Timeout == 0 {
		c.Timeout = xtime.Duration(time.Second)
	}
	if len(c.Endpoints) == 0 {
		errInfo := fmt.Sprintf("zookeeper New failed, endpoints is null")
		log.Error(errInfo)
		return nil, errors.New(errInfo)
	}

	zkConn, connEvent, err := zk.Connect(c.Endpoints, time.Duration(c.Timeout))
	if err != nil {
		log.Error(fmt.Sprintf("zk Connect err:(%v)", err))
		return
	}
	log.Info(fmt.Sprintf("zk Connect ok!"))

	ctx, cancel := context.WithCancel(context.Background())
	zkb = &Zookeeper{
		c:          c,
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
func (z *Zookeeper) Build(appid string, options ...naming.BuildOpt) naming.Resolver {
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
	})
	return r
}

// Scheme return zookeeper's scheme.
func (z *Zookeeper) Scheme() string {
	return "zookeeper"
}

// Register is register instance.
func (z *Zookeeper) Register(ctx context.Context, ins *naming.Instance) (cancelFunc context.CancelFunc, err error) {
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
				log.Info("watch zkClient state, connEvent:(%+v)", connEvent)
				if connEvent.State == zk.StateHasSession {
					if err = z.register(ctx, ins); err != nil {
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

func (z *Zookeeper) createPath(paths string) error {
	var (
		lastPath = "/"
		seps     = strings.Split(paths, "/")
	)
	for _, part := range seps {
		if part == "" {
			continue
		}
		lastPath = path.Join(lastPath, part)
		ok, _, err := z.cli.Exists(lastPath)
		if err != nil {
			return err
		}
		if ok {
			continue
		}
		ret, err := z.cli.Create(lastPath, nil, 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			log.Warn(fmt.Sprintf("createPath, fail to Create node:(%s). error:(%v)", paths, err))
		} else {
			log.Info(fmt.Sprintf("createPath, succeed to Create node:(%s). retStr:(%s)", paths, ret))
		}
	}
	return nil
}

func (z *Zookeeper) registerPeerServer(nodePath string, ins *naming.Instance) (err error) {
	var (
		str string
	)
	val, err := json.Marshal(ins)
	if err != nil {
		return
	}
	log.Info(fmt.Sprintf("registerPeerServer, ins after json.Marshal:(%v)", string(val)))
	ok, _, err := z.cli.Exists(nodePath)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	str, err = z.cli.Create(nodePath, val, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err != nil {
		log.Warn(fmt.Sprintf("registerPeerServer, fail to Create node:%s. error:(%v)", nodePath, err))
	} else {
		log.Info(fmt.Sprintf("registerPeerServer, succeed to Create node:%s. retStr:(%s)", nodePath, str))
	}
	return
}

// register is register instance to zookeeper.
func (z *Zookeeper) register(ctx context.Context, ins *naming.Instance) (err error) {
	log.Info("zookeeper register enter, instance Addrs:(%v)", ins.Addrs)

	prefix := z.keyPrefix(ins.AppID)
	if err = z.createPath(prefix); err != nil {
		log.Warn(fmt.Sprintf("register, fail to createPath node error:(%v)", err))
	}
	for _, addr := range ins.Addrs {
		u, err := url.Parse(addr)
		if err != nil {
			continue
		}
		// grpc://127.0.0.1:8000 to 127.0.0.1
		nodePath := prefix + "/" + strings.SplitN(u.Host, ":", 2)[0]
		if err = z.registerPeerServer(nodePath, ins); err != nil {
			log.Warn(fmt.Sprintf("registerServer, fail to RegisterPeerServer node:%s error:(%v)", addr, err))
		} else {
			log.Info(fmt.Sprintf("registerServer, succeed to RegistServer node."))
		}
	}
	return nil
}

func (z *Zookeeper) unregister(ins *naming.Instance) (err error) {
	log.Info("zookeeper unregister enter, instance Addrs:(%v)", ins.Addrs)
	prefix := z.keyPrefix(ins.AppID)
	for _, addr := range ins.Addrs {
		u, err := url.Parse(addr)
		if err != nil {
			continue
		}
		// grpc://127.0.0.1:8000 to 127.0.0.1
		nodePath := prefix + "/" + strings.SplitN(u.Host, ":", 2)[0]
		exists, _, err := z.cli.Exists(nodePath)
		if err != nil {
			log.Error("zk.Conn.Exists node:(%v), error:(%v)", nodePath, err)
			continue
		}
		if exists {
			_, s, err := z.cli.Get(nodePath)
			if err != nil {
				log.Error("zk.Conn.Get node:(%s), error:(%v)", nodePath, err)
				continue
			}
			if err = z.cli.Delete(nodePath, s.Version); err != nil {
				log.Error("zk.Conn.Delete node:(%s), error:(%v)", nodePath, err)
				continue
			}
		}

		log.Info(fmt.Sprintf("unregister, client.Delete:(%v), appid:(%v), hostname:(%v) success", nodePath, ins.AppID, ins.Hostname))
	}
	return
}

func (z *Zookeeper) keyPrefix(appID string) string {
	return path.Join(z.c.Root, appID)
}

// Close stop all running process including zk fetch and register.
func (z *Zookeeper) Close() error {
	z.cancelFunc()
	return nil
}

func (a *appInfo) watch(appID string) {
	_ = a.fetchstore(appID)
	go func() {
		prefix := a.zkb.keyPrefix(appID)
		for {
			log.Info(fmt.Sprintf("zk ChildrenW enter, prefix:(%v)", prefix))
			snapshot, _, event, err := a.zkb.cli.ChildrenW(prefix)
			if err != nil {
				log.Error("zk ChildrenW fail to watch:%s error:(%v)", prefix, err)
				time.Sleep(time.Second)
				_ = a.fetchstore(appID)
				continue
			}
			log.Info(fmt.Sprintf("zk ChildrenW ok, prefix:%s snapshot:(%v)", prefix, snapshot))
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
	prefix := a.zkb.keyPrefix(appID)
	childs, _, err := a.zkb.cli.Children(prefix)
	if err != nil {
		log.Error(fmt.Sprintf("fetchstore, fail to get Children of node:(%v), error:(%v)", prefix, err))
		return
	}
	log.Info(fmt.Sprintf("fetchstore, ok to get Children of node:(%v), childs:(%v)", prefix, childs))
	ins := &naming.InstancesInfo{
		Instances: make(map[string][]*naming.Instance, 0),
	}
	for _, child := range childs {
		nodePath := prefix + "/" + child
		resp, _, err := a.zkb.cli.Get(nodePath)
		if err != nil {
			log.Error("zookeeper: fetch client.Get(%s) error:(%v)", nodePath, err)
			return err
		}
		in := new(naming.Instance)
		if err = json.Unmarshal(resp, in); err != nil {
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
