package etcd

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bilibili/kratos/pkg/log"
	"github.com/bilibili/kratos/pkg/naming"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"google.golang.org/grpc"
)

var (
	//etcdPrefix is a etcd globe key prefix
	endpoints  string
	etcdPrefix string

	//Time units is second
	registerTTL        = 90
	defaultDialTimeout = 30
)

var (
	_once    sync.Once
	_builder naming.Builder
	//ErrDuplication is a register duplication err
	ErrDuplication = errors.New("etcd: instance duplicate registration")
)

func init() {
	addFlag(flag.CommandLine)
}

func addFlag(fs *flag.FlagSet) {
	// env
	fs.StringVar(&endpoints, "etcd.endpoints", os.Getenv("ETCD_ENDPOINTS"), "etcd.endpoints is etcd endpoints. value: 127.0.0.1:2379,127.0.0.2:2379 etc.")
	fs.StringVar(&etcdPrefix, "etcd.prefix", defaultString("ETCD_PREFIX", "kratos_etcd"), "etcd globe key prefix or use ETCD_PREFIX env variable. value etcd_prefix etc.")
}

func defaultString(env, value string) string {
	v := os.Getenv(env)
	if v == "" {
		return value
	}
	return v
}

// Builder return default etcd resolver builder.
func Builder(c *clientv3.Config) naming.Builder {
	_once.Do(func() {
		_builder, _ = New(c)
	})
	return _builder
}

// Build register resolver into default etcd.
func Build(c *clientv3.Config, id string) naming.Resolver {
	return Builder(c).Build(id)
}

// EtcdBuilder is a etcd clientv3 EtcdBuilder
type EtcdBuilder struct {
	cli        *clientv3.Client
	ctx        context.Context
	cancelFunc context.CancelFunc

	mutex    sync.RWMutex
	apps     map[string]*appInfo
	registry map[string]struct{}
}
type appInfo struct {
	resolver map[*Resolve]struct{}
	ins      atomic.Value
	e        *EtcdBuilder
	once     sync.Once
}

// Resolve etch resolver.
type Resolve struct {
	id    string
	event chan struct{}
	e     *EtcdBuilder
	opt   *naming.BuildOptions
}

// New is new a etcdbuilder
func New(c *clientv3.Config) (e *EtcdBuilder, err error) {
	if c == nil {
		if endpoints == "" {
			panic(fmt.Errorf("invalid etcd config endpoints:%+v", endpoints))
		}
		c = &clientv3.Config{
			Endpoints:   strings.Split(endpoints, ","),
			DialTimeout: time.Second * time.Duration(defaultDialTimeout),
			DialOptions: []grpc.DialOption{grpc.WithBlock()},
		}
	}
	cli, err := clientv3.New(*c)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	e = &EtcdBuilder{
		cli:        cli,
		ctx:        ctx,
		cancelFunc: cancel,
		apps:       map[string]*appInfo{},
		registry:   map[string]struct{}{},
	}
	return
}

// Build disovery resovler builder.
func (e *EtcdBuilder) Build(appid string, opts ...naming.BuildOpt) naming.Resolver {
	r := &Resolve{
		id:    appid,
		e:     e,
		event: make(chan struct{}, 1),
		opt:   new(naming.BuildOptions),
	}
	e.mutex.Lock()
	app, ok := e.apps[appid]
	if !ok {
		app = &appInfo{
			resolver: make(map[*Resolve]struct{}),
			e:        e,
		}
		e.apps[appid] = app
	}
	app.resolver[r] = struct{}{}
	e.mutex.Unlock()
	if ok {
		select {
		case r.event <- struct{}{}:
		default:
		}
	}

	app.once.Do(func() {
		go app.watch(appid)
		log.Info("etcd: AddWatch(%s) already watch(%v)", appid, ok)
	})
	return r
}

// Scheme return etcd's scheme
func (e *EtcdBuilder) Scheme() string {
	return "etcd"

}

// Register is register instance
func (e *EtcdBuilder) Register(ctx context.Context, ins *naming.Instance) (cancelFunc context.CancelFunc, err error) {
	e.mutex.Lock()
	if _, ok := e.registry[ins.AppID]; ok {
		err = ErrDuplication
	} else {
		e.registry[ins.AppID] = struct{}{}
	}
	e.mutex.Unlock()
	if err != nil {
		return
	}
	ctx, cancel := context.WithCancel(e.ctx)
	if err = e.register(ctx, ins); err != nil {
		e.mutex.Lock()
		delete(e.registry, ins.AppID)
		e.mutex.Unlock()
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
				_ = e.register(ctx, ins)
			case <-ctx.Done():
				_ = e.unregister(ins)
				ch <- struct{}{}
				return
			}
		}
	}()
	return
}

//注册和续约公用一个操作
func (e *EtcdBuilder) register(ctx context.Context, ins *naming.Instance) (err error) {
	prefix := e.keyPrefix(ins)
	val, _ := json.Marshal(ins)

	ttlResp, err := e.cli.Grant(context.TODO(), int64(registerTTL))
	if err != nil {
		log.Error("etcd: register client.Grant(%v) error(%v)", registerTTL, err)
		return err
	}
	_, err = e.cli.Put(ctx, prefix, string(val), clientv3.WithLease(ttlResp.ID))
	if err != nil {
		log.Error("etcd: register client.Put(%v) appid(%s) hostname(%s) error(%v)",
			prefix, ins.AppID, ins.Hostname, err)
		return err
	}
	return nil
}
func (e *EtcdBuilder) unregister(ins *naming.Instance) (err error) {
	prefix := e.keyPrefix(ins)

	if _, err = e.cli.Delete(context.TODO(), prefix); err != nil {
		log.Error("etcd: unregister client.Delete(%v) appid(%s) hostname(%s) error(%v)",
			prefix, ins.AppID, ins.Hostname, err)
	}
	log.Info("etcd: unregister client.Delete(%v)  appid(%s) hostname(%s) success",
		prefix, ins.AppID, ins.Hostname)
	return
}

func (e *EtcdBuilder) keyPrefix(ins *naming.Instance) string {
	return fmt.Sprintf("/%s/%s/%s", etcdPrefix, ins.AppID, ins.Hostname)
}

// Close stop all running process including etcdfetch and register
func (e *EtcdBuilder) Close() error {
	e.cancelFunc()
	return nil
}
func (a *appInfo) watch(appID string) {
	_ = a.fetchstore(appID)
	prefix := fmt.Sprintf("/%s/%s/", etcdPrefix, appID)
	rch := a.e.cli.Watch(a.e.ctx, prefix, clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			if ev.Type == mvccpb.PUT || ev.Type == mvccpb.DELETE {
				_ = a.fetchstore(appID)
			}
		}
	}
}

func (a *appInfo) fetchstore(appID string) (err error) {
	prefix := fmt.Sprintf("/%s/%s/", etcdPrefix, appID)
	resp, err := a.e.cli.Get(a.e.ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		log.Error("etcd: fetch client.Get(%s) error(%+v)", prefix, err)
		return err
	}

	ins, err := a.paserIns(resp)
	if err != nil {
		return err
	}
	a.store(ins)
	return nil
}
func (a *appInfo) store(ins *naming.InstancesInfo) {

	a.ins.Store(ins)
	a.e.mutex.RLock()
	for rs := range a.resolver {
		select {
		case rs.event <- struct{}{}:
		default:
		}
	}
	a.e.mutex.RUnlock()
}

func (a *appInfo) paserIns(resp *clientv3.GetResponse) (ins *naming.InstancesInfo, err error) {
	ins = &naming.InstancesInfo{
		Instances: make(map[string][]*naming.Instance, 0),
	}
	for _, ev := range resp.Kvs {
		in := new(naming.Instance)

		err := json.Unmarshal(ev.Value, in)
		if err != nil {
			return nil, err
		}
		ins.Instances[in.Zone] = append(ins.Instances[in.Zone], in)
	}
	return ins, nil
}

// Watch watch instance.
func (r *Resolve) Watch() <-chan struct{} {
	return r.event
}

// Fetch fetch resolver instance.
func (r *Resolve) Fetch(ctx context.Context) (ins *naming.InstancesInfo, ok bool) {
	r.e.mutex.RLock()
	app, ok := r.e.apps[r.id]
	r.e.mutex.RUnlock()
	if ok {
		ins, ok = app.ins.Load().(*naming.InstancesInfo)
		return
	}
	return
}

// Close close resolver.
func (r *Resolve) Close() error {
	r.e.mutex.Lock()
	if app, ok := r.e.apps[r.id]; ok && len(app.resolver) != 0 {
		delete(app.resolver, r)
	}
	r.e.mutex.Unlock()
	return nil
}
