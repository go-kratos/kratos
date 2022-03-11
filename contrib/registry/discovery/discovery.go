package discovery

import (
	"context"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"

	"github.com/SeeMusic/kratos/v2/log"
)

type Discovery struct {
	config     *Config
	once       sync.Once
	ctx        context.Context
	cancelFunc context.CancelFunc
	httpClient *resty.Client

	node    atomic.Value
	nodeIdx uint64

	mutex       sync.RWMutex
	apps        map[string]*appInfo
	registry    map[string]struct{}
	lastHost    string
	cancelPolls context.CancelFunc

	logger log.Logger
}

type appInfo struct {
	resolver map[*Resolve]struct{}
	zoneIns  atomic.Value
	lastTs   int64 // latest timestamp
}

// New construct a Discovery instance which implements registry.Registrar,
// registry.Discovery and registry.Watcher.
func New(c *Config, logger log.Logger) *Discovery {
	if logger == nil {
		logger = log.NewStdLogger(os.Stdout)
		logger = log.With(logger,
			"registry.pluginName", "Discovery",
			"ts", log.DefaultTimestamp,
			"caller", log.DefaultCaller,
		)
	}
	if c == nil {
		c = new(Config)
	}
	if err := fixConfig(c); err != nil {
		panic(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	d := &Discovery{
		config:     c,
		ctx:        ctx,
		cancelFunc: cancel,
		apps:       map[string]*appInfo{},
		registry:   map[string]struct{}{},
		logger:     logger,
	}

	d.httpClient = resty.New().
		SetTimeout(40 * time.Second)

	// Discovery self found and watch
	r := d.resolveBuild(_discoveryAppID)
	event := r.Watch()
	_, ok := <-event
	if !ok {
		panic("Discovery watch self failed")
	}
	discoveryIns, ok := r.Fetch(context.Background())
	if ok {
		d.newSelf(discoveryIns.Instances)
	}
	go d.selfProc(r, event)

	return d
}

// Close stop all running process including Discovery and register
func (d *Discovery) Close() error {
	d.cancelFunc()
	return nil
}

func (d *Discovery) Logger() *log.Helper {
	return log.NewHelper(d.logger)
}

// selfProc start a goroutine to refresh Discovery self registration information.
func (d *Discovery) selfProc(resolver *Resolve, event <-chan struct{}) {
	for {
		_, ok := <-event
		if !ok {
			return
		}
		zones, ok := resolver.Fetch(context.Background())
		if ok {
			d.newSelf(zones.Instances)
		}
	}
}

// newSelf
func (d *Discovery) newSelf(zones map[string][]*discoveryInstance) {
	ins, ok := zones[d.config.Zone]
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
	var olds int
	for _, n := range nodes {
		if node, ok := d.node.Load().([]string); ok {
			for _, o := range node {
				if o == n {
					olds++
					break
				}
			}
		}
	}
	if len(nodes) == olds {
		return
	}

	rand.Shuffle(len(nodes), func(i, j int) {
		nodes[i], nodes[j] = nodes[j], nodes[i]
	})
	d.node.Store(nodes)
}

// resolveBuild Discovery resolver builder.
func (d *Discovery) resolveBuild(appID string) *Resolve {
	r := &Resolve{
		id:    appID,
		d:     d,
		event: make(chan struct{}, 1),
	}

	d.mutex.Lock()
	app, ok := d.apps[appID]
	if !ok {
		app = &appInfo{
			resolver: make(map[*Resolve]struct{}),
		}
		d.apps[appID] = app
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

	d.Logger().Debugf("Discovery: AddWatch(%s) already watch(%v)", appID, ok)
	d.once.Do(func() {
		go d.serverProc()
	})
	return r
}

func (d *Discovery) serverProc() {
	defer d.Logger().Debug("Discovery serverProc quit")

	var (
		retry  int
		ctx    context.Context
		cancel context.CancelFunc
	)

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
			d.switchNode()
		default:
		}

		apps, err := d.polls(ctx)
		if err != nil {
			d.switchNode()
			if ctx.Err() == context.Canceled {
				ctx = nil
				continue
			}
			time.Sleep(time.Second)
			retry++
			continue
		}
		retry = 0
		d.broadcast(apps)
	}
}

func (d *Discovery) pickNode() string {
	nodes, ok := d.node.Load().([]string)
	if !ok || len(nodes) == 0 {
		return d.config.Nodes[rand.Intn(len(d.config.Nodes))]
	}
	return nodes[atomic.LoadUint64(&d.nodeIdx)%uint64(len(nodes))]
}

func (d *Discovery) switchNode() {
	atomic.AddUint64(&d.nodeIdx, 1)
}

// renew an instance with Discovery
func (d *Discovery) renew(ctx context.Context, ins *discoveryInstance) (err error) {
	// d.Logger().Debug("Discovery:renew renew calling")

	d.mutex.RLock()
	c := d.config
	d.mutex.RUnlock()

	res := new(discoveryCommonResp)
	uri := fmt.Sprintf(_renewURL, d.pickNode())

	// construct parameters to renew
	p := newParams(d.config)
	p.Set(_paramKeyAppID, ins.AppID)

	// send request to Discovery server.
	if _, err = d.httpClient.R().
		SetContext(ctx).
		SetQueryParamsFromValues(p).
		SetResult(&res).
		Post(uri); err != nil {
		d.switchNode()
		d.Logger().Errorf("Discovery: renew client.Get(%v) env(%s) appid(%s) hostname(%s) error(%v)",
			uri, c.Env, ins.AppID, c.Host, err)
		return
	}

	if res.Code != _codeOK {
		err = fmt.Errorf("discovery.renew failed ErrorCode: %d", res.Code)
		if res.Code == _codeNotFound {
			if err = d.register(ctx, ins); err != nil {
				err = errors.Wrap(err, "Discovery.renew instance, and failed to register ins")
			}
			return
		}

		d.Logger().Errorf(
			"Discovery: renew client.Get(%v) env(%s) appid(%s) hostname(%s) code(%v)",
			uri, c.Env, ins.AppID, c.Host, res.Code,
		)
	}

	return
}

// cancel Remove the registered instance from Discovery
func (d *Discovery) cancel(ins *discoveryInstance) (err error) {
	d.mutex.RLock()
	config := d.config
	d.mutex.RUnlock()

	res := new(discoveryCommonResp)
	uri := fmt.Sprintf(_cancelURL, d.pickNode())
	p := newParams(d.config)
	p.Set(_paramKeyAppID, ins.AppID)

	// request
	// send request to Discovery server.
	if _, err = d.httpClient.R().
		SetContext(context.TODO()).
		SetQueryParamsFromValues(p).
		SetResult(&res).
		Post(uri); err != nil {
		d.switchNode()
		d.Logger().Errorf("Discovery cancel client.Get(%v) env(%s) appid(%s) hostname(%s) error(%v)",
			uri, config.Env, ins.AppID, config.Host, err)
		return
	}

	// handle response error
	if res.Code != _codeOK {
		if res.Code == _codeNotFound {
			return nil
		}

		d.Logger().Warnf("Discovery cancel client.Get(%v)  env(%s) appid(%s) hostname(%s) code(%v)",
			uri, config.Env, ins.AppID, config.Host, res.Code)
		err = fmt.Errorf("ErrorCode: %d", res.Code)
		return
	}

	return
}

func (d *Discovery) broadcast(apps map[string]*disInstancesInfo) {
	for appID, v := range apps {
		var count int
		// v maybe nil in old version(less than v1.1) Discovery,check incase of panic
		if v == nil {
			continue
		}
		for zone, ins := range v.Instances {
			if len(ins) == 0 {
				delete(v.Instances, zone)
			}
			count += len(ins)
		}
		if count == 0 {
			continue
		}
		d.mutex.RLock()
		app, ok := d.apps[appID]
		d.mutex.RUnlock()
		if ok {
			app.lastTs = v.LastTs
			app.zoneIns.Store(v)
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

func (d *Discovery) polls(ctx context.Context) (apps map[string]*disInstancesInfo, err error) {
	var (
		lastTss = make([]int64, 0, 4)
		appIDs  = make([]string, 0, 16)
		host    = d.pickNode()
		changed bool
	)
	if host != d.lastHost {
		d.lastHost = host
		changed = true
	}

	d.mutex.RLock()
	config := d.config
	for k, v := range d.apps {
		if changed {
			v.lastTs = 0
		}
		appIDs = append(appIDs, k)
		lastTss = append(lastTss, v.lastTs)
	}
	d.mutex.RUnlock()

	// if there is no app, polls just return.
	if len(appIDs) == 0 {
		return
	}

	uri := fmt.Sprintf(_pollURL, host)
	res := new(discoveryPollsResp)

	// params
	p := newParams(nil)
	p.Set(_paramKeyEnv, config.Env)
	p.Set(_paramKeyHostname, config.Host)
	for _, appID := range appIDs {
		p.Add(_paramKeyAppID, appID)
	}
	for _, ts := range lastTss {
		p.Add("latest_timestamp", strconv.FormatInt(ts, 10))
	}

	// request
	reqURI := uri + "?" + p.Encode()
	if _, err = d.httpClient.R().
		SetContext(ctx).
		SetQueryParamsFromValues(p).
		SetResult(res).Get(uri); err != nil {
		d.switchNode()
		d.Logger().Errorf("Discovery: client.Get(%s) error(%+v)", reqURI, err)
		return nil, err
	}

	if res.Code != _codeOK {
		if res.Code != _codeNotModified {
			d.Logger().Errorf("Discovery: client.Get(%s) get error code(%d)", reqURI, res.Code)
		}
		err = fmt.Errorf("discovery.polls failed ErrCode: %d", res.Code)
		return
	}

	for _, app := range res.Data {
		if app.LastTs == 0 {
			err = ErrServerError
			d.Logger().Errorf("Discovery: client.Get(%s) latest_timestamp is 0, instances:(%+v)", reqURI, res.Data)
			return
		}
	}

	d.Logger().Debugf("Discovery: successfully polls(%s) instances (%+v)", reqURI, res.Data)
	apps = res.Data
	return
}

// Resolve Discovery resolver.
type Resolve struct {
	id    string
	event chan struct{}
	d     *Discovery
}

// Watch instance.
func (r *Resolve) Watch() <-chan struct{} {
	return r.event
}

// Fetch resolver instance.
func (r *Resolve) Fetch(ctx context.Context) (ins *disInstancesInfo, ok bool) {
	r.d.mutex.RLock()
	app, ok := r.d.apps[r.id]
	r.d.mutex.RUnlock()
	if ok {
		var appIns *disInstancesInfo
		appIns, ok = app.zoneIns.Load().(*disInstancesInfo)
		if !ok {
			return
		}
		ins = new(disInstancesInfo)
		ins.LastTs = appIns.LastTs
		ins.Scheduler = appIns.Scheduler
		ins.Instances = make(map[string][]*discoveryInstance)
		for zone, in := range appIns.Instances {
			ins.Instances[zone] = in
		}
		//if r.opt.Filter != nil {
		//	ins.Instances = r.opt.Filter(appIns.Instances)
		//} else {
		//	ins.Instances = make(map[string][]*discoveryInstance)
		//	for zone, in := range appIns.Instances {
		//		ins.Instances[zone] = in
		//	}
		//}
		//if r.opt.scheduler != nil {
		//	ins.Instances[r.opt.ClientZone] = r.opt.scheduler(ins)
		//}
		//if r.opt.Subset != nil && r.opt.SubsetSize != 0 {
		//	for zone, inss := range ins.Instances {
		//		ins.Instances[zone] = r.opt.Subset(inss, r.opt.SubsetSize)
		//	}
		//}
	}
	return
}

// Close resolver
func (r *Resolve) Close() error {
	r.d.mutex.Lock()
	if app, ok := r.d.apps[r.id]; ok && len(app.resolver) != 0 {
		delete(app.resolver, r)
		// TODO: delete app from builder
	}
	r.d.mutex.Unlock()
	return nil
}
