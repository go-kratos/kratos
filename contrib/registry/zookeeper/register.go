package zookeeper

import (
	"context"
	"path"
	"time"

	"github.com/go-zookeeper/zk"
	"golang.org/x/sync/singleflight"

	"github.com/go-kratos/kratos/v2/registry"
)

var (
	_ registry.Registrar = &Registry{}
	_ registry.Discovery = &Registry{}
)

// Option is etcd registry option.
type Option func(o *options)

type options struct {
	namespace string
	user      string
	password  string
}

// WithRootPath with registry root path.
func WithRootPath(path string) Option {
	return func(o *options) { o.namespace = path }
}

// WithDigestACL with registry password.
func WithDigestACL(user string, password string) Option {
	return func(o *options) {
		o.user = user
		o.password = password
	}
}

// Registry is consul registry
type Registry struct {
	opts *options
	conn *zk.Conn

	group singleflight.Group
}

func New(conn *zk.Conn, opts ...Option) *Registry {
	options := &options{
		namespace: "/microservices",
	}
	for _, o := range opts {
		o(options)
	}
	return &Registry{
		opts: options,
		conn: conn,
	}
}

func (r *Registry) Register(ctx context.Context, service *registry.ServiceInstance) error {
	var (
		data []byte
		err  error
	)
	if err = r.ensureName(r.opts.namespace, []byte(""), 0); err != nil {
		return err
	}
	serviceNamePath := path.Join(r.opts.namespace, service.Name)
	if err = r.ensureName(serviceNamePath, []byte(""), 0); err != nil {
		return err
	}
	if data, err = marshal(service); err != nil {
		return err
	}
	servicePath := path.Join(serviceNamePath, service.ID)
	if err = r.ensureName(servicePath, data, zk.FlagEphemeral); err != nil {
		return err
	}
	go r.reRegister(servicePath, data)
	return nil
}

// Deregister registry service to zookeeper.
func (r *Registry) Deregister(ctx context.Context, service *registry.ServiceInstance) error {
	ch := make(chan error, 1)
	servicePath := path.Join(r.opts.namespace, service.Name, service.ID)
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
	instances, err, _ := r.group.Do(serviceName, func() (interface{}, error) {
		serviceNamePath := path.Join(r.opts.namespace, serviceName)
		servicesID, _, err := r.conn.Children(serviceNamePath)
		if err != nil {
			return nil, err
		}
		items := make([]*registry.ServiceInstance, 0, len(servicesID))
		for _, service := range servicesID {
			servicePath := path.Join(serviceNamePath, service)
			serviceInstanceByte, _, err := r.conn.Get(servicePath)
			if err != nil {
				return nil, err
			}
			item, err := unmarshal(serviceInstanceByte)
			if err != nil {
				return nil, err
			}
			items = append(items, item)
		}
		return items, nil
	})
	if err != nil {
		return nil, err
	}
	return instances.([]*registry.ServiceInstance), nil
}

func (r *Registry) Watch(ctx context.Context, serviceName string) (registry.Watcher, error) {
	prefix := path.Join(r.opts.namespace, serviceName)
	return newWatcher(ctx, prefix, serviceName, r.conn)
}

// ensureName ensure node exists, if not exist, create and set data
func (r *Registry) ensureName(path string, data []byte, flags int32) error {
	exists, stat, err := r.conn.Exists(path)
	if err != nil {
		return err
	}
	// ephemeral nodes handling after restart
	// fixes a race condition if the server crashes without using CreateProtectedEphemeralSequential()
	if flags&zk.FlagEphemeral == zk.FlagEphemeral {
		err = r.conn.Delete(path, stat.Version)
		if err != nil && err != zk.ErrNoNode {
			return err
		}
		exists = false
	}
	if !exists {
		if len(r.opts.user) > 0 && len(r.opts.password) > 0 {
			_, err = r.conn.Create(path, data, flags, zk.DigestACL(zk.PermAll, r.opts.user, r.opts.password))
		} else {
			_, err = r.conn.Create(path, data, flags, zk.WorldACL(zk.PermAll))
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// reRegister re-register data node info when bad connection recovered
func (r *Registry) reRegister(path string, data []byte) {
	sessionID := r.conn.SessionID()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for range ticker.C {
		cur := r.conn.SessionID()
		// sessionID changed
		if cur > 0 && sessionID != cur {
			// re-ensureName
			if err := r.ensureName(path, data, zk.FlagEphemeral); err != nil {
				return
			}
			sessionID = cur
		}
	}
}
