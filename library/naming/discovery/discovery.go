package discovery

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"go-common/library/conf/env"
	"go-common/library/ecode"
	"go-common/library/exp/feature"
	"go-common/library/log"
	"go-common/library/naming"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/netutil"
	"go-common/library/net/netutil/breaker"
	xtime "go-common/library/time"
	"go-common/library/xstr"
)

const (
	_registerURL = "http://%s/discovery/register"
	_setURL      = "http://%s/discovery/set"
	_cancelURL   = "http://%s/discovery/cancel"
	_renewURL    = "http://%s/discovery/renew"

	_pollURL  = "http://%s/discovery/polls"
	_nodesURL = "http://%s/discovery/nodes"

	_registerGap = 30 * time.Second

	_statusUP = "1"
)

const (
	_appid = "infra.discovery"
)

var (
	_ naming.Builder  = &Discovery{}
	_ naming.Registry = &Discovery{}

	_selfDiscoveryFeatrue feature.Feature = "discovery.self"
	_discoveryFeatures                    = map[feature.Feature]feature.Spec{
		_selfDiscoveryFeatrue: {Default: false},
	}

	// ErrDuplication duplication treeid.
	ErrDuplication = errors.New("discovery: instance duplicate registration")
)

func init() {
	feature.DefaultGate.Add(_discoveryFeatures)
}

// Config discovery configures.
type Config struct {
	Nodes  []string
	Key    string
	Secret string
	Region string
	Zone   string
	Env    string
	Host   string
}

type appData struct {
	ZoneInstances map[string][]*naming.Instance `json:"zone_instances"`
	LastTs        int64                         `json:"latest_timestamp"`
}

// Discovery is discovery client.
type Discovery struct {
	once       sync.Once
	conf       *Config
	ctx        context.Context
	cancelFunc context.CancelFunc
	httpClient *bm.Client

	mutex       sync.RWMutex
	apps        map[string]*appInfo
	registry    map[string]struct{}
	lastHost    string
	cancelPolls context.CancelFunc
	idx         uint64
	node        atomic.Value
	delete      chan *appInfo
}

type appInfo struct {
	zoneIns  atomic.Value
	resolver map[*Resolver]struct{}
	lastTs   int64 // latest timestamp
}

func fixConfig(c *Config) {
	if len(c.Nodes) == 0 {
		c.Nodes = []string{"api.bilibili.co"}
	}
	if env.Region != "" {
		c.Region = env.Region
	}
	if env.Zone != "" {
		c.Zone = env.Zone
	}
	if env.DeployEnv != "" {
		c.Env = env.DeployEnv
	}
	if env.Hostname != "" {
		c.Host = env.Hostname
	} else {
		c.Host, _ = os.Hostname()
	}
}

var (
	once              sync.Once
	_defaultDiscovery *Discovery
)

func initDefault() {
	once.Do(func() {
		_defaultDiscovery = New(nil)
	})
}

// Builder return default discvoery resolver builder.
func Builder() naming.Builder {
	if _defaultDiscovery == nil {
		initDefault()
	}
	return _defaultDiscovery
}

// Build register resolver into default discovery.
func Build(id string) naming.Resolver {
	if _defaultDiscovery == nil {
		initDefault()
	}
	return _defaultDiscovery.Build(id)
}

// New new a discovery client.
func New(c *Config) (d *Discovery) {
	if c == nil {
		c = &Config{
			Nodes:  []string{"discovery.bilibili.co", "api.bilibili.co"},
			Key:    "discovery",
			Secret: "discovery",
		}
	}
	fixConfig(c)
	ctx, cancel := context.WithCancel(context.Background())
	d = &Discovery{
		ctx:        ctx,
		cancelFunc: cancel,
		conf:       c,
		apps:       map[string]*appInfo{},
		registry:   map[string]struct{}{},
		delete:     make(chan *appInfo, 10),
	}
	// httpClient
	cfg := &bm.ClientConfig{
		App: &bm.App{
			Key:    c.Key,
			Secret: c.Secret,
		},
		Dial:    xtime.Duration(3 * time.Second),
		Timeout: xtime.Duration(40 * time.Second),
		Breaker: &breaker.Config{
			Window:  100,
			Sleep:   3,
			Bucket:  10,
			Ratio:   0.5,
			Request: 100,
		},
	}
	d.httpClient = bm.NewClient(cfg)
	if feature.DefaultGate.Enabled(_selfDiscoveryFeatrue) {
		resolver := d.Build(_appid)
		event := resolver.Watch()
		_, ok := <-event
		if !ok {
			panic("discovery watch failed")
		}
		ins, ok := resolver.Fetch(context.Background())
		if ok {
			d.newSelf(ins)
		}
		go d.selfproc(resolver, event)
	}
	return
}

func (d *Discovery) selfproc(resolver naming.Resolver, event <-chan struct{}) {
	for {
		_, ok := <-event
		if !ok {
			return
		}
		zones, ok := resolver.Fetch(context.Background())
		if ok {
			d.newSelf(zones)
		}
	}
}

func (d *Discovery) newSelf(zones map[string][]*naming.Instance) {
	ins, ok := zones[d.conf.Zone]
	if !ok {
		return
	}
	var nodes []string
	for _, in := range ins {
		for _, addr := range in.Addrs {
			u, err := url.Parse(addr)
			if err == nil && u.Scheme == "http" {
				nodes = append(nodes, u.Host)
			}
		}
	}
	// diff old nodes
	olds, ok := d.node.Load().([]string)
	if ok {
		var diff int
		for _, n := range nodes {
			for _, o := range olds {
				if o == n {
					diff++
					break
				}
			}
		}
		if len(nodes) == diff {
			return
		}
	}
	// FIXME: we should use rand.Shuffle() in golang 1.10
	Shuffle(len(nodes), func(i, j int) {
		nodes[i], nodes[j] = nodes[j], nodes[i]
	})
	d.node.Store(nodes)
}

// Build disovery resovler builder.
func (d *Discovery) Build(appid string) naming.Resolver {
	r := &Resolver{
		id:    appid,
		d:     d,
		event: make(chan struct{}, 1),
	}
	d.mutex.Lock()
	app, ok := d.apps[appid]
	if !ok {
		app = &appInfo{
			resolver: make(map[*Resolver]struct{}),
		}
		d.apps[appid] = app
		cancel := d.cancelPolls
		if cancel != nil {
			cancel()
		}
	}
	app.resolver[r] = struct{}{}
	d.mutex.Unlock()
	if ok {
		select {
		case r.event <- struct{}{}:
		default:
		}
	}
	log.Info("disocvery: AddWatch(%s) already watch(%v)", appid, ok)
	d.once.Do(func() {
		go d.serverproc()
	})
	return r
}

// Scheme return discovery's scheme
func (d *Discovery) Scheme() string {
	return "discovery"
}

// Resolver discveory resolver.
type Resolver struct {
	id    string
	event chan struct{}
	d     *Discovery
}

// Watch watch instance.
func (r *Resolver) Watch() <-chan struct{} {
	return r.event
}

// Fetch fetch resolver instance.
func (r *Resolver) Fetch(c context.Context) (ins map[string][]*naming.Instance, ok bool) {
	r.d.mutex.RLock()
	app, ok := r.d.apps[r.id]
	r.d.mutex.RUnlock()
	if ok {
		ins, ok = app.zoneIns.Load().(map[string][]*naming.Instance)
		return
	}
	return
}

// Close close resolver.
func (r *Resolver) Close() error {
	r.d.mutex.Lock()
	if app, ok := r.d.apps[r.id]; ok && len(app.resolver) != 0 {
		delete(app.resolver, r)
		// TODO: delete app from builder
	}
	r.d.mutex.Unlock()
	return nil
}

func (d *Discovery) pickNode() string {
	nodes, ok := d.node.Load().([]string)
	if !ok || len(nodes) == 0 {
		return d.conf.Nodes[d.idx%uint64(len(d.conf.Nodes))]
	}
	return nodes[d.idx%uint64(len(nodes))]
}

func (d *Discovery) switchNode() {
	atomic.AddUint64(&d.idx, 1)
}

// Reload reload the config
func (d *Discovery) Reload(c *Config) {
	fixConfig(c)
	d.mutex.Lock()
	d.conf = c
	d.mutex.Unlock()
}

// Close stop all running process including discovery and register
func (d *Discovery) Close() error {
	d.cancelFunc()
	return nil
}

// Register Register an instance with discovery and renew automatically
func (d *Discovery) Register(c context.Context, ins *naming.Instance) (cancelFunc context.CancelFunc, err error) {
	d.mutex.Lock()
	if _, ok := d.registry[ins.AppID]; ok {
		err = ErrDuplication
	} else {
		d.registry[ins.AppID] = struct{}{}
	}
	d.mutex.Unlock()
	if err != nil {
		return
	}
	if err = d.register(c, ins); err != nil {
		d.mutex.Lock()
		delete(d.registry, ins.AppID)
		d.mutex.Unlock()
		return
	}
	ctx, cancel := context.WithCancel(d.ctx)
	ch := make(chan struct{}, 1)
	cancelFunc = context.CancelFunc(func() {
		cancel()
		<-ch
	})
	go func() {
		ticker := time.NewTicker(_registerGap)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := d.renew(ctx, ins); err != nil && ecode.NothingFound.Equal(err) {
					d.register(ctx, ins)
				}
			case <-ctx.Done():
				d.cancel(ins)
				ch <- struct{}{}
				return
			}
		}
	}()
	return
}

// Set set ins status and metadata.
func (d *Discovery) Set(ins *naming.Instance) error {
	return d.set(context.Background(), ins)
}

// cancel Remove the registered instance from discovery
func (d *Discovery) cancel(ins *naming.Instance) (err error) {
	d.mutex.RLock()
	conf := d.conf
	d.mutex.RUnlock()

	res := new(struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	})
	uri := fmt.Sprintf(_cancelURL, d.pickNode())
	params := d.newParams(conf)
	params.Set("appid", ins.AppID)
	// request
	if err = d.httpClient.Post(context.Background(), uri, "", params, &res); err != nil {
		d.switchNode()
		log.Error("discovery cancel client.Get(%v) env(%s) appid(%s) hostname(%s) error(%v)",
			uri, conf.Env, ins.AppID, conf.Host, err)
		return
	}
	if ec := ecode.Int(res.Code); !ec.Equal(ecode.OK) {
		log.Warn("discovery cancel client.Get(%v)  env(%s) appid(%s) hostname(%s) code(%v)",
			uri, conf.Env, ins.AppID, conf.Host, res.Code)
		err = ec
		return
	}
	log.Info("discovery cancel client.Get(%v)  env(%s) appid(%s) hostname(%s) success",
		uri, conf.Env, ins.AppID, conf.Host)
	return
}

// register Register an instance with discovery
func (d *Discovery) register(ctx context.Context, ins *naming.Instance) (err error) {
	d.mutex.RLock()
	conf := d.conf
	d.mutex.RUnlock()

	var metadata []byte
	if ins.Metadata != nil {
		if metadata, err = json.Marshal(ins.Metadata); err != nil {
			log.Error("discovery:register instance Marshal metadata(%v) failed!error(%v)", ins.Metadata, err)
		}
	}
	res := new(struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	})
	uri := fmt.Sprintf(_registerURL, d.pickNode())
	params := d.newParams(conf)
	params.Set("appid", ins.AppID)
	params.Set("addrs", strings.Join(ins.Addrs, ","))
	params.Set("version", ins.Version)
	params.Set("status", _statusUP)
	params.Set("metadata", string(metadata))
	if err = d.httpClient.Post(ctx, uri, "", params, &res); err != nil {
		d.switchNode()
		log.Error("discovery: register client.Get(%v)  zone(%s) env(%s) appid(%s) addrs(%v) error(%v)",
			uri, conf.Zone, conf.Env, ins.AppID, ins.Addrs, err)
		return
	}
	if ec := ecode.Int(res.Code); !ec.Equal(ecode.OK) {
		log.Warn("discovery: register client.Get(%v)  env(%s) appid(%s) addrs(%v)  code(%v)",
			uri, conf.Env, ins.AppID, ins.Addrs, res.Code)
		err = ec
		return
	}
	log.Info("discovery: register client.Get(%v) env(%s) appid(%s) addrs(%s) success",
		uri, conf.Env, ins.AppID, ins.Addrs)
	return
}

// rset set  instance info with discovery
func (d *Discovery) set(ctx context.Context, ins *naming.Instance) (err error) {
	d.mutex.RLock()
	conf := d.conf
	d.mutex.RUnlock()
	res := new(struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	})
	uri := fmt.Sprintf(_setURL, d.pickNode())
	params := d.newParams(conf)
	params.Set("appid", ins.AppID)
	params.Set("version", ins.Version)
	params.Set("status", strconv.FormatInt(ins.Status, 10))
	if ins.Metadata != nil {
		var metadata []byte
		if metadata, err = json.Marshal(ins.Metadata); err != nil {
			log.Error("discovery:set instance Marshal metadata(%v) failed!error(%v)", ins.Metadata, err)
		}
		params.Set("metadata", string(metadata))
	}
	if err = d.httpClient.Post(ctx, uri, "", params, &res); err != nil {
		d.switchNode()
		log.Error("discovery: set client.Get(%v)  zone(%s) env(%s) appid(%s) addrs(%v) error(%v)",
			uri, conf.Zone, conf.Env, ins.AppID, ins.Addrs, err)
		return
	}
	if ec := ecode.Int(res.Code); !ec.Equal(ecode.OK) {
		log.Warn("discovery: set client.Get(%v)  env(%s) appid(%s) addrs(%v)  code(%v)",
			uri, conf.Env, ins.AppID, ins.Addrs, res.Code)
		err = ec
		return
	}
	log.Info("discovery: set client.Get(%v) env(%s) appid(%s) addrs(%s) success",
		uri+"?"+params.Encode(), conf.Env, ins.AppID, ins.Addrs)
	return
}

// renew Renew an instance with discovery
func (d *Discovery) renew(ctx context.Context, ins *naming.Instance) (err error) {
	d.mutex.RLock()
	conf := d.conf
	d.mutex.RUnlock()

	res := new(struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	})
	uri := fmt.Sprintf(_renewURL, d.pickNode())
	params := d.newParams(conf)
	params.Set("appid", ins.AppID)
	if err = d.httpClient.Post(ctx, uri, "", params, &res); err != nil {
		d.switchNode()
		log.Error("discovery: renew client.Get(%v)  env(%s) appid(%s) hostname(%s) error(%v)",
			uri, conf.Env, ins.AppID, conf.Host, err)
		return
	}
	if ec := ecode.Int(res.Code); !ec.Equal(ecode.OK) {
		err = ec
		if ec.Equal(ecode.NothingFound) {
			return
		}
		log.Error("discovery: renew client.Get(%v) env(%s) appid(%s) hostname(%s) code(%v)",
			uri, conf.Env, ins.AppID, conf.Host, res.Code)
		return
	}
	return
}

func (d *Discovery) serverproc() {
	var (
		retry  int
		update bool
		ctx    context.Context
		cancel context.CancelFunc
	)
	bc := netutil.DefaultBackoffConfig
	ticker := time.NewTicker(time.Minute * 30)
	defer ticker.Stop()
	for {
		if ctx == nil {
			ctx, cancel = context.WithCancel(d.ctx)
			d.mutex.Lock()
			d.cancelPolls = cancel
			d.mutex.Unlock()
		}
		select {
		case <-d.ctx.Done():
			return
		case <-ticker.C:
			update = true
		default:
		}
		if !feature.DefaultGate.Enabled(_selfDiscoveryFeatrue) {
			nodes, ok := d.node.Load().([]string)
			if !ok || len(nodes) == 0 || update {
				update = false
				tnodes := d.nodes()
				if len(tnodes) == 0 {
					time.Sleep(bc.Backoff(retry))
					retry++
					continue
				}
				retry = 0
				// FIXME: we should use rand.Shuffle() in golang 1.10
				Shuffle(len(tnodes), func(i, j int) {
					tnodes[i], tnodes[j] = tnodes[j], tnodes[i]
				})
				d.node.Store(tnodes)
			}
		}
		apps, err := d.polls(ctx, d.pickNode())
		if err != nil {
			d.switchNode()
			if ctx.Err() == context.Canceled {
				ctx = nil
				continue
			}
			time.Sleep(bc.Backoff(retry))
			retry++
			continue
		}
		retry = 0
		d.broadcast(apps)
	}
}

func (d *Discovery) nodes() (nodes []string) {
	res := new(struct {
		Code int `json:"code"`
		Data []struct {
			Addr string `json:"addr"`
		} `json:"data"`
	})
	uri := fmt.Sprintf(_nodesURL, d.pickNode())
	if err := d.httpClient.Get(d.ctx, uri, "", nil, res); err != nil {
		d.switchNode()
		log.Error("discovery: consumer client.Get(%v)error(%+v)", uri, err)
		return
	}
	if ec := ecode.Int(res.Code); !ec.Equal(ecode.OK) {
		log.Error("discovery: consumer client.Get(%v) error(%v)", uri, res.Code)
		return
	}
	if len(res.Data) == 0 {
		log.Warn("discovery: get nodes(%s) failed,no nodes found!", uri)
		return
	}
	nodes = make([]string, 0, len(res.Data))
	for i := range res.Data {
		nodes = append(nodes, res.Data[i].Addr)
	}
	return
}

func (d *Discovery) polls(ctx context.Context, host string) (apps map[string]appData, err error) {
	var (
		lastTs  []int64
		appid   []string
		changed bool
	)
	if host != d.lastHost {
		d.lastHost = host
		changed = true
	}
	d.mutex.RLock()
	conf := d.conf
	for k, v := range d.apps {
		if changed {
			v.lastTs = 0
		}
		appid = append(appid, k)
		lastTs = append(lastTs, v.lastTs)
	}
	d.mutex.RUnlock()
	if len(appid) == 0 {
		return
	}
	uri := fmt.Sprintf(_pollURL, host)
	res := new(struct {
		Code int                `json:"code"`
		Data map[string]appData `json:"data"`
	})
	params := url.Values{}
	params.Set("env", conf.Env)
	params.Set("hostname", conf.Host)
	params.Set("appid", strings.Join(appid, ","))
	params.Set("latest_timestamp", xstr.JoinInts(lastTs))
	if err = d.httpClient.Get(ctx, uri, "", params, res); err != nil {
		log.Error("discovery: client.Get(%s) error(%+v)", uri+"?"+params.Encode(), err)
		return
	}
	if ec := ecode.Int(res.Code); !ec.Equal(ecode.OK) {
		if !ec.Equal(ecode.NotModified) {
			log.Error("discovery: client.Get(%s) get error code(%d)", uri+"?"+params.Encode(), res.Code)
			err = ec
		}
		return
	}
	info, _ := json.Marshal(res.Data)
	for _, app := range res.Data {
		if app.LastTs == 0 {
			err = ecode.ServerErr
			log.Error("discovery: client.Get(%s) latest_timestamp is 0,instances:(%s)", uri+"?"+params.Encode(), info)
			return
		}
	}
	log.Info("discovery: polls uri(%s)", uri+"?"+params.Encode())
	log.Info("discovery: successfully polls(%s) instances (%s)", uri+"?"+params.Encode(), info)
	apps = res.Data
	return
}

func (d *Discovery) broadcast(apps map[string]appData) {
	for id, v := range apps {
		var count int
		for zone, ins := range v.ZoneInstances {
			if len(ins) == 0 {
				delete(v.ZoneInstances, zone)
			}
			count += len(ins)
		}
		if count == 0 {
			continue
		}
		d.mutex.RLock()
		app, ok := d.apps[id]
		d.mutex.RUnlock()
		if ok {
			app.lastTs = v.LastTs
			app.zoneIns.Store(v.ZoneInstances)
			d.mutex.RLock()
			for rs := range app.resolver {
				select {
				case rs.event <- struct{}{}:
				default:
				}
			}
			d.mutex.RUnlock()
		}
	}
}

func (d *Discovery) newParams(conf *Config) url.Values {
	params := url.Values{}
	params.Set("region", conf.Region)
	params.Set("zone", conf.Zone)
	params.Set("env", conf.Env)
	params.Set("hostname", conf.Host)
	return params
}
