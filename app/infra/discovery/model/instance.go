package model

import (
	"encoding/json"
	"strconv"
	"sync"
	"time"

	"go-common/library/ecode"
	"go-common/library/log"
)

// InstanceStatus Status of instance
// type InstanceStatus uint32

const (
	// InstanceStatusUP Ready to receive traffic
	InstanceStatusUP = uint32(1)
	// InstancestatusWating Intentionally shutdown for traffic
	InstancestatusWating = uint32(1) << 1
)

func (i *Instance) filter(status uint32) bool {
	return status&i.Status > 0
}

// Action Replicate type of node
type Action int

const (
	// Register Replicate the add action to all nodes
	Register Action = iota
	// Renew Replicate the heartbeat action to all nodes
	Renew
	// Cancel Replicate the cancel action to all nodes
	Cancel
	// Weight Replicate the Weight action to all nodes
	Weight
	// Delete Replicate the Delete action to all nodes
	Delete
	// Status Replicate the Status action to all nodes
	Status
)

// Instance holds information required for registration with
// <Discovery Server> and to be discovered by other components.
type Instance struct {
	Region   string            `json:"region"`
	Zone     string            `json:"zone"`
	Env      string            `json:"env"`
	Appid    string            `json:"appid"`
	Treeid   int64             `json:"treeid"`
	Hostname string            `json:"hostname"`
	HTTP     string            `json:"http"`
	RPC      string            `json:"rpc"`
	Version  string            `json:"version"`
	Metadata map[string]string `json:"metadata"`
	Addrs    []string          `json:"addrs"`
	// Status enum instance status
	Status uint32 `json:"status"`

	// timestamp
	RegTimestamp   int64 `json:"reg_timestamp"`
	UpTimestamp    int64 `json:"up_timestamp"` // NOTE: It is latest timestamp that status becomes UP.
	RenewTimestamp int64 `json:"renew_timestamp"`
	DirtyTimestamp int64 `json:"dirty_timestamp"`

	LatestTimestamp int64 `json:"latest_timestamp"`
}

// NewInstance new a instance.
func NewInstance(arg *ArgRegister) (i *Instance) {
	now := time.Now().UnixNano()
	i = &Instance{
		Region:          arg.Region,
		Zone:            arg.Zone,
		Env:             arg.Env,
		Appid:           arg.Appid,
		Treeid:          arg.Treeid,
		Hostname:        arg.Hostname,
		HTTP:            arg.HTTP,
		RPC:             arg.RPC,
		Version:         arg.Version,
		Status:          arg.Status,
		Addrs:           arg.Addrs,
		RegTimestamp:    now,
		UpTimestamp:     now,
		LatestTimestamp: now,
		RenewTimestamp:  now,
		DirtyTimestamp:  now,
	}
	i.Metadata = make(map[string]string)
	if arg.Metadata != "" {
		if err := json.Unmarshal([]byte(arg.Metadata), &i.Metadata); err != nil {
			log.Error("json unmarshal metadata err %v", err)
		}
	}
	return
}

// InstanceInfo the info get by consumer.
type InstanceInfo struct {
	Instances          []*Instance            `json:"instances"`
	ZoneInstances      map[string][]*Instance `json:"zone_instances"`
	LatestTimestamp    int64                  `json:"latest_timestamp"`
	LatestTimestampStr string                 `json:"latest_timestamp_str"`
}

// Apps app distinguished by zone
type Apps struct {
	apps            map[string]*App
	lock            sync.RWMutex
	latestTimestamp int64
}

// NewApps return new Apps.
func NewApps() *Apps {
	return &Apps{
		apps: make(map[string]*App),
	}
}

// NewApp news a app by appid. If ok=false, returns the app of already exist.
func (p *Apps) NewApp(zone, appid string, treeid, lts int64) (a *App, new bool) {
	p.lock.Lock()
	a, ok := p.apps[zone]
	if !ok {
		a = NewApp(zone, appid, treeid)
		p.apps[zone] = a
	}
	if lts <= p.latestTimestamp {
		// insure increase
		lts = p.latestTimestamp + 1
	}
	p.latestTimestamp = lts
	p.lock.Unlock()
	new = !ok
	return
}

// App get app by zone.
func (p *Apps) App(zone string) (as []*App) {
	p.lock.RLock()
	if zone != "" {
		a, ok := p.apps[zone]
		if !ok {
			p.lock.RUnlock()
			return
		}
		as = []*App{a}
	} else {
		for _, a := range p.apps {
			as = append(as, a)
		}
	}
	p.lock.RUnlock()
	return
}

// Del del app by zone.
func (p *Apps) Del(zone string) {
	p.lock.Lock()
	delete(p.apps, zone)
	p.lock.Unlock()
}

// InstanceInfo return slice of instances.if up is true,return all status instance else return up status instance
func (p *Apps) InstanceInfo(zone string, latestTime int64, status uint32) (ci *InstanceInfo, err error) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	if latestTime >= p.latestTimestamp {
		err = ecode.NotModified
		return
	}
	ci = &InstanceInfo{
		LatestTimestamp:    p.latestTimestamp,
		LatestTimestampStr: strconv.FormatInt(p.latestTimestamp/int64(time.Second), 10),
		ZoneInstances:      make(map[string][]*Instance),
	}
	var ok bool
	for z, app := range p.apps {
		if zone == "" || z == zone {
			ok = true
			as := app.Instances()
			if len(as) == 0 {
				continue
			}
			instance := make([]*Instance, 0, len(as))
			for _, i := range as {
				// if up is false return all status instance
				if i.filter(status) {
					// if i.Status == InstanceStatusUP && i.LatestTimestamp > latestTime { // TODO(felix): increase
					ni := new(Instance)
					*ni = *i
					instance = append(instance, ni)
				}
			}
			ci.Instances = append(ci.Instances, instance...)
			ci.ZoneInstances[z] = instance
		}
	}
	if !ok {
		err = ecode.NothingFound
	} else if len(ci.Instances) == 0 {
		err = ecode.NotModified
	}
	return
}

// UpdateLatest update LatestTimestamp.
func (p *Apps) UpdateLatest(latestTime int64) {
	if latestTime <= p.latestTimestamp {
		// insure increase
		latestTime = p.latestTimestamp + 1
	}
	p.latestTimestamp = latestTime
}

// App Instances distinguished by hostname
type App struct {
	AppID           string
	Treeid          int64
	Zone            string
	instances       map[string]*Instance
	latestTimestamp int64
	lock            sync.RWMutex
}

// NewApp new App.
func NewApp(zone, appid string, treeid int64) (a *App) {
	a = &App{
		Treeid:    treeid,
		AppID:     appid,
		Zone:      zone,
		instances: make(map[string]*Instance),
	}
	return
}

// Instances return slice of instances.
func (a *App) Instances() (is []*Instance) {
	a.lock.RLock()
	is = make([]*Instance, 0, len(a.instances))
	for _, i := range a.instances {
		ni := new(Instance)
		*ni = *i
		is = append(is, ni)
	}
	a.lock.RUnlock()
	return
}

// NewInstance new a instance.
func (a *App) NewInstance(ni *Instance, latestTime int64) (i *Instance, ok bool) {
	i = new(Instance)
	a.lock.Lock()
	oi, ok := a.instances[ni.Hostname]
	if ok {
		ni.UpTimestamp = oi.UpTimestamp
		if ni.DirtyTimestamp < oi.DirtyTimestamp {
			log.Warn("register exist(%v) dirty timestamp over than caller(%v)", oi, ni)
			ni = oi
		}
	}
	a.instances[ni.Hostname] = ni
	a.updateLatest(latestTime)
	*i = *ni
	a.lock.Unlock()
	ok = !ok
	return
}

// Renew new a instance.
func (a *App) Renew(hostname string) (i *Instance, ok bool) {
	i = new(Instance)
	a.lock.Lock()
	defer a.lock.Unlock()
	oi, ok := a.instances[hostname]
	if !ok {
		return
	}
	oi.RenewTimestamp = time.Now().UnixNano()
	*i = *oi
	return
}

func (a *App) updateLatest(latestTime int64) {
	if latestTime <= a.latestTimestamp {
		// insure increase
		latestTime = a.latestTimestamp + 1
	}
	a.latestTimestamp = latestTime
}

// Cancel cancel a instance.
func (a *App) Cancel(hostname string, latestTime int64) (i *Instance, l int, ok bool) {
	i = new(Instance)
	a.lock.Lock()
	defer a.lock.Unlock()
	oi, ok := a.instances[hostname]
	if !ok {
		return
	}
	delete(a.instances, hostname)
	l = len(a.instances)
	oi.LatestTimestamp = latestTime
	a.updateLatest(latestTime)
	*i = *oi
	return
}

// Len returns the length of instances.
func (a *App) Len() (l int) {
	a.lock.RLock()
	l = len(a.instances)
	a.lock.RUnlock()
	return
}

// Set set new status,metadata of instance .
func (a *App) Set(changes *ArgSet) (ok bool) {
	a.lock.Lock()
	defer a.lock.Unlock()
	var (
		dst     *Instance
		setTime int64
	)
	if changes.SetTimestamp == 0 {
		setTime = time.Now().UnixNano()
	}
	for i, hostname := range changes.Hostname {
		if dst, ok = a.instances[hostname]; !ok {
			log.Error("Set hostname(%s) not found", hostname)
			return
		}
		if len(changes.Status) != 0 {
			if uint32(changes.Status[i]) != InstanceStatusUP && uint32(changes.Status[i]) != InstancestatusWating {
				log.Error("SetStatus change status(%d) is error", changes.Status[i])
				ok = false
				return
			}
			dst.Status = uint32(changes.Status[i])
			if dst.Status == InstanceStatusUP {
				dst.UpTimestamp = setTime
			}
		}
		if len(changes.Metadata) != 0 {
			metadata := make(map[string]string)
			if err := json.Unmarshal([]byte(changes.Metadata[i]), &metadata); err != nil {
				log.Error("set change metadata err %s", changes.Metadata[i])
				ok = false
				return
			}
			dst.Metadata = metadata
		}
		dst.LatestTimestamp = setTime
		dst.DirtyTimestamp = setTime
	}
	a.updateLatest(setTime)
	return
}
