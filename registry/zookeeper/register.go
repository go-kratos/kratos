package registry

import (
	"context"
	"encoding/json"
	"path"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-zookeeper/zk"
)

var (
	_ registry.Registrar = &Registry{}
	_ registry.Discovery = &Registry{}
)

// Option is etcd registry option.
type Option func(o *options)

type options struct {
	ctx      context.Context
	rootPath string
	timeout  time.Duration
}

// WithContext with registry context.
func WithContext(ctx context.Context) Option {
	return func(o *options) { o.ctx = ctx }
}

// WithRootPath with registry root path.
func WithRootPath(path string) Option {
	return func(o *options) { o.rootPath = path }
}

// WithTimeout with registry timeout.
func WithTimeout(timeout time.Duration) Option {
	return func(o *options) { o.timeout = timeout }
}

// Registry is consul registry
type Registry struct {
	opts     *options
	conn     *zk.Conn
	lock     sync.Mutex
	registry map[string]*serviceSet
}

func New(zkServers []string, opts ...Option) (*Registry, error) {
	options := &options{
		ctx:      context.Background(),
		rootPath: "/microservices",
		timeout:  time.Second * 5,
	}
	for _, o := range opts {
		o(options)
	}
	conn, _, err := zk.Connect(zkServers, options.timeout)
	if err != nil {
		return nil, err
	}
	return &Registry{
		opts:     options,
		conn:     conn,
		registry: make(map[string]*serviceSet),
	}, err
}

func (r *Registry) Register(ctx context.Context, service *registry.ServiceInstance) error {
	var data []byte
	var err error
	if err := r.ensureName(r.opts.rootPath, []byte("")); err != nil {
		return err
	}
	serviceNamePath := path.Join(r.opts.rootPath, service.Name)
	if err = r.ensureName(serviceNamePath, []byte("")); err != nil {
		return err
	}
	if data, err = json.Marshal(service); err != nil {
		return err
	}
	servicePath := path.Join(serviceNamePath, service.ID)
	if err = r.ensureName(servicePath, data); err != nil {
		return err
	}
	return nil
}

// Deregister registry service to zookeeper.
func (r *Registry) Deregister(ctx context.Context, service *registry.ServiceInstance) error {
	ch := make(chan error, 1)
	servicePath := path.Join(r.opts.rootPath, service.Name, service.ID)
	go func() {
		err := r.conn.Delete(servicePath, -1)
		ch <- err
	}()
	var err error
	select {
	case <-ctx.Done():
		err = ctx.Err()
	case err = <-ch:
	}
	return err
}

// GetService get services from zookeeper
func (r *Registry) GetService(ctx context.Context, serviceName string) ([]*registry.ServiceInstance, error) {
	serviceNamePath := path.Join(r.opts.rootPath, serviceName)
	servicesID, _, err := r.conn.Children(serviceNamePath)
	if err != nil {
		return nil, err
	}
	var items []*registry.ServiceInstance
	for _, service := range servicesID {
		var item = &registry.ServiceInstance{}
		servicePath := path.Join(serviceNamePath, service)
		serviceInstanceByte, _, err := r.conn.Get(servicePath)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(serviceInstanceByte, item); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *Registry) Watch(ctx context.Context, serviceName string) (registry.Watcher, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	set, ok := r.registry[serviceName]
	if !ok {
		set = &serviceSet{
			watcher:     make(map[*watcher]struct{}, 0),
			services:    &atomic.Value{},
			serviceName: serviceName,
		}
		r.registry[serviceName] = set
	}
	// 初始化watcher
	w := &watcher{
		event: make(chan struct{}, 1),
	}
	w.ctx, w.cancel = context.WithCancel(context.Background())
	w.set = set
	set.lock.Lock()
	set.watcher[w] = struct{}{}
	set.lock.Unlock()
	ss, _ := set.services.Load().([]*registry.ServiceInstance)
	if len(ss) > 0 {
		// 如果services有值需要推送给watcher，否则watch的时候可能会永远阻塞拿不到初始的数据
		w.event <- struct{}{}
	}

	// 放在最后是为了防止漏推送
	if !ok {
		go r.resolve(set)
	}
	return w, nil
}

func (r *Registry) resolve(ss *serviceSet) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	services, err := r.GetService(ctx, ss.serviceName)
	cancel()
	if err == nil && len(services) > 0 {
		ss.broadcast(services)
	}
}

// ensureName ensure node exists, if not exist, create and set data
func (r *Registry) ensureName(path string, data []byte) error {
	exists, _, err := r.conn.Exists(path)
	if err != nil {
		return err
	}
	if !exists {
		_, err := r.conn.Create(path, data, 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			return err
		}
	}
	return nil
}
