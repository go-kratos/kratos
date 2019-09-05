package livezk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"go-common/library/log"
	"go-common/library/naming"
	xtime "go-common/library/time"

	"github.com/pkg/errors"
	"github.com/samuel/go-zookeeper/zk"
)

const (
	basePath = "/live/service"
	scheme   = "grpc"
)

// Zookeeper Server&Client settings.
type Zookeeper struct {
	Root         string
	Addrs        []string
	Timeout      xtime.Duration
	PullInterval xtime.Duration
}

// New new live zookeeper registry
func New(config *Zookeeper) (naming.Registry, error) {
	lz := &livezk{
		zkConfig: config,
	}
	var err error
	lz.zkConn, lz.zkEvent, err = zk.Connect(config.Addrs, time.Duration(config.Timeout))
	go lz.eventproc()
	return lz, err
}

// NewResolveBuilder is
func NewResolveBuilder(config *Zookeeper) (naming.Builder, error) {
	if config.PullInterval < 10 {
		config.PullInterval = xtime.Duration(10 * time.Second)
	}
	zr := &resolvBuilder{
		zkConfig: config,
		apps:     map[string]*appInfo{},
	}
	conn, event, err := zk.Connect(config.Addrs, time.Duration(config.Timeout))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	zr.zkConn = conn
	zr.zkEvent = event
	go zr.eventproc()
	return zr, nil
}

type zkIns struct {
	Group       string `json:"group"`
	LibVersion  string `json:"lib_version"`
	StartupTime string `json:"startup_time"`
}

func newZkInsData(ins *naming.Instance) ([]byte, error) {
	zi := &zkIns{
		// TODO group support
		Group:       "default",
		LibVersion:  ins.Version,
		StartupTime: time.Now().Format("2006-01-02 15:04:05"),
	}
	return json.Marshal(zi)
}

var _ naming.Registry = &livezk{}
var _ naming.Builder = &resolvBuilder{}
var _ naming.Resolver = &resolver{}

// livezk live service zookeeper registry
type livezk struct {
	zkConfig *Zookeeper
	zkConn   *zk.Conn
	zkEvent  <-chan zk.Event
}

type resolvBuilder struct {
	zkConfig *Zookeeper
	zkConn   *zk.Conn
	zkEvent  <-chan zk.Event

	mutex sync.RWMutex
	apps  map[string]*appInfo
}

type resolver struct {
	id    string
	event chan struct{}

	pullInterval time.Duration
	zkEvent      <-chan zk.Event
	rb           *resolvBuilder
}

type appInfo struct {
	zoneIns  atomic.Value
	resolver map[*resolver]struct{}
}

func (l *livezk) Register(ctx context.Context, ins *naming.Instance) (cancel context.CancelFunc, err error) {
	nodePath := path.Join(l.zkConfig.Root, basePath, ins.AppID)
	if err = l.createAll(nodePath); err != nil {
		return
	}
	var rpc string
	for _, addr := range ins.Addrs {
		u, ue := url.Parse(addr)
		if ue == nil && u.Scheme == scheme {
			rpc = u.Host
			break
		}
	}
	if rpc == "" {
		err = errors.New("no GRPC addr")
		return
	}

	dataPath := path.Join(nodePath, rpc)
	data, err := newZkInsData(ins)
	if err != nil {
		return nil, err
	}
	_, err = l.zkConn.Create(dataPath, data, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err != nil {
		return nil, err
	}
	return func() {
		l.unregister(dataPath)
	}, nil
}

func (l *livezk) Close() error {
	l.zkConn.Close()
	return nil
}

func (l *livezk) createAll(nodePath string) (err error) {
	seps := strings.Split(nodePath, "/")
	lastPath := "/"
	ok := false
	for _, part := range seps {
		if part == "" {
			continue
		}
		lastPath = path.Join(lastPath, part)
		if ok, _, err = l.zkConn.Exists(lastPath); err != nil {
			return err
		} else if ok {
			continue
		}
		if _, err = l.zkConn.Create(lastPath, nil, 0, zk.WorldACL(zk.PermAll)); err != nil {
			return
		}
	}
	return
}

func (l *livezk) eventproc() {
	for event := range l.zkEvent {
		// TODO handle zookeeper event
		log.Info("zk event: err: %+v, path: %s, server: %s, state: %s, type: %s",
			event.Err, event.Path, event.Server, event.State, event.Type)
	}
}

func (l *livezk) unregister(dataPath string) error {
	return l.zkConn.Delete(dataPath, -1)
}

func (zr *resolver) Close() error {
	zr.rb.mutex.Lock()
	if app, ok := zr.rb.apps[zr.id]; ok && len(app.resolver) != 0 {
		delete(app.resolver, zr)
		// TODO: delete app from builder
	}
	zr.rb.mutex.Unlock()
	return nil
}

func (zr *resolver) Fetch(ctx context.Context) (map[string][]*naming.Instance, bool) {
	zr.rb.mutex.RLock()
	app, ok := zr.rb.apps[zr.id]
	zr.rb.mutex.RUnlock()
	if ok {
		ins, ok := app.zoneIns.Load().(map[string][]*naming.Instance)
		return ins, ok
	}
	return nil, false
}

func (zr *resolver) Watch() <-chan struct{} {
	return zr.event
}

func (rb *resolvBuilder) Build(appid string) naming.Resolver {
	inss, zkEvent, err := rb.resolvOne(appid)
	if err != nil {
		panic(err)
	}
	r := &resolver{
		id:           appid,
		rb:           rb,
		zkEvent:      zkEvent,
		pullInterval: time.Duration(rb.zkConfig.PullInterval),
		event:        make(chan struct{}, 1),
	}
	_, ok := rb.apps[appid]
	if !ok {
		rb.apps[appid] = &appInfo{
			resolver: map[*resolver]struct{}{},
		}
	}
	ai := rb.apps[appid]
	ai.zoneIns.Store(inss)
	ai.resolver[r] = struct{}{}
	r.event <- struct{}{}
	go r.eventproc()
	return r
}

func (rb *resolvBuilder) Scheme() string {
	return "livezk"
}

func (rb *resolvBuilder) resolvOne(appid string) (map[string][]*naming.Instance, <-chan zk.Event, error) {
	appPath := path.Join(basePath, appid)
	addrs, _, event, err := rb.zkConn.ChildrenW(appPath)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	inss := make(map[string][]*naming.Instance, len(addrs))
	for _, addr := range addrs {
		addrPath := path.Join(appPath, addr)
		bs, _, err := rb.zkConn.Get(addrPath)
		if err != nil {
			log.Error("Failed to get zk value: %s: %+v", addrPath, errors.WithStack(err))
			continue
		}
		zi := new(zkIns)
		if err := json.Unmarshal(bs, zi); err != nil {
			log.Error("Failed to unmarshal zk instance: %s: %+v", string(bs), errors.WithStack(err))
			continue
		}
		ins := &naming.Instance{
			Zone:     zi.Group,
			Addrs:    []string{fmt.Sprintf("bplusrpc://%s", addr)},
			Version:  zi.LibVersion,
			Hostname: addr, // using addr to hostname
			Metadata: map[string]string{
				"group":        zi.Group,
				"startup_time": zi.StartupTime,
			},
		}
		_, ok := inss[ins.Zone]
		if !ok {
			inss[ins.Zone] = make([]*naming.Instance, 0)
		}
		inss[ins.Zone] = append(inss[ins.Zone], ins)
	}
	return inss, event, nil
}

func (rb *resolvBuilder) eventproc() {
	for event := range rb.zkEvent {
		log.Info("resolv builder zk event: err: %+v, path: %s, server: %s, state: %s, type: %s",
			event.Err, event.Path, event.Server, event.State, event.Type)
	}
}

func (zr *resolver) eventproc() {
	var event zk.Event
	for {
		select {
		case event = <-zr.zkEvent:
			log.Info("resolver zk event: %v", event)
		case <-time.After(time.Duration(zr.pullInterval)):
		}

		zr.rb.mutex.RLock()
		app, ok := zr.rb.apps[zr.id]
		zr.rb.mutex.RUnlock()
		if !ok {
			log.Warn("Ignore this event due to app: %s is not watched: event: %+v", zr.id, event)
			continue
		}

		instances, zkEvent, err := zr.rb.resolvOne(zr.id)
		if err != nil {
			log.Error("Error while attempting to read [%s] service instances after connection disruption: %v", zr.id, err)
			continue
		}
		app.zoneIns.Store(instances)
		for resolver := range app.resolver {
			select {
			case resolver.event <- struct{}{}:
			default:
			}
		}
		zr.zkEvent = zkEvent
	}
}
