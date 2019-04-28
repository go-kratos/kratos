package dao

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"go-common/app/infra/discovery/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_evictThreshold = int64(90 * time.Second)
	_evictCeiling   = int64(3600 * time.Second)
)

// Registry handles replication of all operations to peer Discovery nodes to keep them all in sync.
type Registry struct {
	appm  map[string]*model.Apps // appid-env -> apps
	aLock sync.RWMutex

	conns map[string]map[string]*conn // zone.env.appid-> host
	cLock sync.RWMutex

	gd *Guard
}

// conn the poll chan contains consumer.
type conn struct {
	ch         chan map[string]*model.InstanceInfo // TODO(felix): increase
	arg        *model.ArgPolls
	latestTime int64
	count      int
}

// newConn new consumer chan.
func newConn(ch chan map[string]*model.InstanceInfo, latestTime int64, arg *model.ArgPolls) *conn {
	return &conn{ch: ch, latestTime: latestTime, arg: arg, count: 1}
}

// NewRegistry new register.
func NewRegistry() (r *Registry) {
	r = &Registry{
		appm:  make(map[string]*model.Apps),
		conns: make(map[string]map[string]*conn),
		gd:    new(Guard),
	}
	go r.proc()
	return
}

func (r *Registry) newapps(appid, env string) (a *model.Apps, ok bool) {
	key := appsKey(appid, env)
	r.aLock.Lock()
	if a, ok = r.appm[key]; !ok {
		a = model.NewApps()
		r.appm[key] = a
	}
	r.aLock.Unlock()
	return
}

func (r *Registry) apps(appid, env, zone string) (as []*model.App, a *model.Apps, ok bool) {
	key := appsKey(appid, env)
	r.aLock.RLock()
	a, ok = r.appm[key]
	r.aLock.RUnlock()
	if ok {
		as = a.App(zone)
	}
	return
}

func appsKey(appid, env string) string {
	// NOTE disocvery 不区分具体环境
	if appid == model.AppID {
		return appid
	}
	return fmt.Sprintf("%s-%s", appid, env)
}

func (r *Registry) newApp(ins *model.Instance) (a *model.App) {
	as, _ := r.newapps(ins.Appid, ins.Env)
	a, _ = as.NewApp(ins.Zone, ins.Appid, ins.Treeid, ins.LatestTimestamp)
	return
}

// Register a new instance.
func (r *Registry) Register(ins *model.Instance, latestTime int64) (err error) {
	a := r.newApp(ins)
	i, ok := a.NewInstance(ins, latestTime)
	if ok {
		r.gd.incrExp()
	}
	// NOTE: make sure free poll before update appid latest timestamp.
	r.broadcast(i.Env, i.Appid, a)
	return
}

// Renew marks the given instance of the given app name as renewed, and also marks whether it originated from replication.
func (r *Registry) Renew(arg *model.ArgRenew) (i *model.Instance, ok bool) {
	a, _, _ := r.apps(arg.Appid, arg.Env, arg.Zone)
	if len(a) == 0 {
		return
	}
	if i, ok = a[0].Renew(arg.Hostname); !ok {
		return
	}
	r.gd.incrFac()
	return
}

// Cancel cancels the registration of an instance.
func (r *Registry) Cancel(arg *model.ArgCancel) (i *model.Instance, ok bool) {
	if i, ok = r.cancel(arg.Zone, arg.Env, arg.Appid, arg.Hostname, arg.LatestTimestamp); !ok {
		return
	}
	r.gd.decrExp()
	return
}

func (r *Registry) cancel(zone, env, appid, hostname string, latestTime int64) (i *model.Instance, ok bool) {
	var l int
	a, as, _ := r.apps(appid, env, zone)
	if len(a) == 0 {
		return
	}
	if i, l, ok = a[0].Cancel(hostname, latestTime); !ok {
		return
	}
	as.UpdateLatest(latestTime)
	if l == 0 {
		if a[0].Len() == 0 {
			as.Del(zone)
		}
	}
	if len(as.App("")) == 0 {
		r.aLock.Lock()
		delete(r.appm, appsKey(appid, env))
		r.aLock.Unlock()
	}
	r.broadcast(env, appid, a[0]) // NOTE: make sure free poll before update appid latest timestamp.
	return
}

// FetchAll fetch all instances of all the families.
func (r *Registry) FetchAll() (im map[string][]*model.Instance) {
	ass := r.allapp()
	im = make(map[string][]*model.Instance)
	for _, as := range ass {
		for _, a := range as.App("") {
			im[a.AppID] = append(im[a.AppID], a.Instances()...)
		}
	}
	return
}

// Fetch fetch all instances by appid.
func (r *Registry) Fetch(zone, env, appid string, latestTime int64, status uint32) (info *model.InstanceInfo, err error) {
	key := appsKey(appid, env)
	r.aLock.RLock()
	a, ok := r.appm[key]
	r.aLock.RUnlock()
	if !ok {
		err = ecode.NothingFound
		return
	}
	info, err = a.InstanceInfo(zone, latestTime, status)
	return
}

// Polls hangs request and then write instances when that has changes, or return NotModified.
func (r *Registry) Polls(arg *model.ArgPolls) (ch chan map[string]*model.InstanceInfo, new bool, err error) {
	var (
		ins = make(map[string]*model.InstanceInfo, len(arg.Treeid))
		in  *model.InstanceInfo
	)
	if len(arg.Appid) != len(arg.LatestTimestamp) {
		arg.LatestTimestamp = make([]int64, len(arg.Appid))
	}
	for i := range arg.Appid {
		in, err = r.Fetch(arg.Zone, arg.Env, arg.Appid[i], arg.LatestTimestamp[i], model.InstanceStatusUP)
		if err == ecode.NothingFound {
			log.Error("Polls region(%s) zone(%s) env(%s) appid(%s) error(%v)", arg.Region, arg.Zone, arg.Env, arg.Appid[i], err)
			return
		}
		if err == nil {
			if len(arg.Treeid) != 0 {
				ins[strconv.FormatInt(arg.Treeid[i], 10)] = in
			} else {
				ins[arg.Appid[i]] = in
			}
			new = true
		}
	}
	if new {
		ch = make(chan map[string]*model.InstanceInfo, 1)
		ch <- ins
		return
	}
	r.cLock.Lock()
	for i := range arg.Appid {
		k := appsKey(arg.Appid[i], arg.Env)
		if _, ok := r.conns[k]; !ok {
			r.conns[k] = make(map[string]*conn, 1)
		}
		connection, ok := r.conns[k][arg.Hostname]
		if !ok {
			if ch == nil {
				ch = make(chan map[string]*model.InstanceInfo, 5) // NOTE: there maybe have more than one connection on the same hostname!!!
			}
			connection = newConn(ch, arg.LatestTimestamp[i], arg)
			log.Info("Polls from(%s) new connection(%d)", arg.Hostname, connection.count)
		} else {
			connection.count++ // NOTE: there maybe have more than one connection on the same hostname!!!
			if ch == nil {
				ch = connection.ch
			}
			log.Info("Polls from(%s) reuse connection(%d)", arg.Hostname, connection.count)
		}
		r.conns[k][arg.Hostname] = connection
	}
	r.cLock.Unlock()
	return
}

// Polling get polling clients.
func (r *Registry) Polling(arg *model.ArgPolling) (resp []string, err error) {
	r.cLock.RLock()
	conns, ok := r.conns[appsKey(arg.Appid, arg.Env)]
	if !ok {
		r.cLock.RUnlock()
		return
	}
	resp = make([]string, 0, len(conns))
	for host := range conns {
		resp = append(resp, host)
	}
	r.cLock.RUnlock()
	return
}

// broadcast on poll by chan.
// NOTE: make sure free poll before update appid latest timestamp.
func (r *Registry) broadcast(env, appid string, a *model.App) {
	key := appsKey(appid, env)
	r.cLock.Lock()
	defer r.cLock.Unlock()
	conns, ok := r.conns[key]
	if !ok {
		return
	}
	delete(r.conns, key)
	for _, conn := range conns {
		ii, _ := r.Fetch(conn.arg.Zone, env, appid, 0, model.InstanceStatusUP) // TODO(felix): latesttime!=0 increase
		var key string
		if len(conn.arg.Treeid) != 0 {
			key = strconv.FormatInt(a.Treeid, 10)
		} else {
			key = a.AppID
		}
		for i := 0; i < conn.count; i++ { // NOTE: there maybe have more than one connection on the same hostname!!!
			select {
			case conn.ch <- map[string]*model.InstanceInfo{key: ii}: // NOTE: if chan is full, means no poller.
				log.Info("broadcast to(%s) success(%d)", conn.arg.Hostname, i+1)
			case <-time.After(time.Millisecond * 500):
				log.Info("broadcast to(%s) failed(%d) maybe chan full", conn.arg.Hostname, i+1)
			}
		}
	}
}

// Set Set the status of instance by hostnames.
func (r *Registry) Set(c context.Context, arg *model.ArgSet) (ok bool) {
	a, _, _ := r.apps(arg.Appid, arg.Env, arg.Zone)
	if len(a) == 0 {
		return
	}
	if ok = a[0].Set(arg); !ok {
		return
	}
	r.broadcast(arg.Env, arg.Appid, a[0])
	return
}

func (r *Registry) allapp() (ass []*model.Apps) {
	r.aLock.RLock()
	ass = make([]*model.Apps, 0, len(r.appm))
	for _, as := range r.appm {
		ass = append(ass, as)
	}
	r.aLock.RUnlock()
	return
}

// reset expect renews, count the renew of all app, one app has two expect remews in minute.
func (r *Registry) resetExp() {
	cnt := int64(0)
	for _, p := range r.allapp() {
		for _, a := range p.App("") {
			cnt += int64(a.Len())
		}
	}
	r.gd.setExp(cnt)
}

func (r *Registry) proc() {
	tk := time.Tick(1 * time.Minute)
	tk2 := time.Tick(15 * time.Minute)
	for {
		select {
		case <-tk:
			r.gd.updateFac()
			r.evict()
		case <-tk2:
			r.resetExp()
		}
	}
}

func (r *Registry) evict() {
	protect := r.gd.ok()
	// We collect first all expired items, to evict them in random order. For large eviction sets,
	// if we do not that, we might wipe out whole apps before self preservation kicks in. By randomizing it,
	// the impact should be evenly distributed across all applications.
	var eis []*model.Instance
	var registrySize int
	// all projects
	ass := r.allapp()
	for _, as := range ass {
		for _, a := range as.App("") {
			registrySize += a.Len()
			is := a.Instances()
			for _, i := range is {
				delta := time.Now().UnixNano() - i.RenewTimestamp
				if (!protect && delta > _evictThreshold) || delta > _evictCeiling {
					eis = append(eis, i)
				}
			}
		}
	}
	// To compensate for GC pauses or drifting local time, we need to use current registry size as a base for
	// triggering self-preservation. Without that we would wipe out full registry.
	eCnt := len(eis)
	registrySizeThreshold := int(float64(registrySize) * _percentThreshold)
	evictionLimit := registrySize - registrySizeThreshold
	if eCnt > evictionLimit {
		eCnt = evictionLimit
	}
	if eCnt == 0 {
		return
	}
	for i := 0; i < eCnt; i++ {
		// Pick a random item (Knuth shuffle algorithm)
		next := i + rand.Intn(len(eis)-i)
		eis[i], eis[next] = eis[next], eis[i]
		ei := eis[i]
		r.cancel(ei.Zone, ei.Env, ei.Appid, ei.Hostname, time.Now().UnixNano())
	}
}

// DelConns delete conn of host in appid
func (r *Registry) DelConns(arg *model.ArgPolls) {
	r.cLock.Lock()
	for i := range arg.Appid {
		k := appsKey(arg.Appid[i], arg.Env)
		conns, ok := r.conns[k]
		if !ok {
			log.Warn("DelConn key(%s) not found", k)
			continue
		}
		if connection, ok := conns[arg.Hostname]; ok {
			if connection.count > 1 {
				log.Info("DelConns from(%s) count decr(%d)", arg.Hostname, connection.count)
				connection.count--
			} else {
				log.Info("DelConns from(%s) delete(%d)", arg.Hostname, connection.count)
				delete(conns, arg.Hostname)
			}
		}
	}
	r.cLock.Unlock()
}
