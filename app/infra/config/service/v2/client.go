package v2

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go-common/app/infra/config/model"
	"go-common/library/ecode"
	"go-common/library/log"
	xtime "go-common/library/time"
)

const (
	_buildVerKey  = "%d_%s"
	_pushKey      = "%s_%s"
	_cacheKey     = "%s_%d"
	_pushForceKey = "%s_%s_%s"
	_tagForceKey  = "%d_%s_%d_tagforce"
	_tagIDkey     = "%d_%s_tagID"
	_lastForceKey = "%d_%s_lastForce"
)

func lastForceKey(svr int64, bver string) string {
	return fmt.Sprintf(_tagIDkey, svr, bver)
}

func tagIDkey(svr int64, bver string) string {
	return fmt.Sprintf(_tagIDkey, svr, bver)
}

// tagForceKey...
func tagForceKey(svr int64, bver string, tagID int64) string {
	return fmt.Sprintf(_tagForceKey, svr, bver, tagID)
}

// pushForceKey push
func pushForceKey(svr, host, ip string) string {
	return fmt.Sprintf(_pushForceKey, svr, host, ip)
}

// PushKey push sub id
func pushKey(svr, host string) string {
	return fmt.Sprintf(_pushKey, svr, host)
}

// buildVerKey version mapping key
func buildVerKey(svr int64, bver string) string {
	return fmt.Sprintf(_buildVerKey, svr, bver)
}

// cacheKey config cache key
func cacheKey(svr string, ver int64) string {
	return fmt.Sprintf(_cacheKey, svr, ver)
}

// genConfig generate config
func genConfig(ver int64, values []*model.Value) (conf *model.Content, err error) {
	var b []byte
	if b, err = json.Marshal(values); err != nil {
		return
	}
	mb := md5.Sum(b)
	conf = &model.Content{
		Version: ver,
		Md5:     hex.EncodeToString(mb[:]),
		Content: string(b),
	}
	return
}

// Push version to clients & generate config caches
func (s *Service) Push(c context.Context, svr *model.Service) (err error) {
	var (
		hosts  []*model.Host
		values []*model.Value
		conf   *model.Content
		app    *model.App
		ids    []int64
		force  int8
	)
	if ids, err = s.dao.ConfIDs(svr.Version); err != nil {
		return
	}
	if values, err = s.dao.ConfigsByIDs(ids); err != nil {
		return
	}

	if len(values) == 0 {
		err = fmt.Errorf("config values is empty. svr:%s, host:%s, buildVer:%s, ver:%d", svr.Name, svr.Host, svr.BuildVersion, svr.Version)
		log.Error("%v", err)
		return
	}
	// compatible old version sdk
	if conf, err = genConfig(svr.Version, values); err != nil {
		log.Error("get config value:%s error(%v) ", values, err)
		return
	}
	cacheName := cacheKey(svr.Name, svr.Version)
	if err = s.dao.SetFile(cacheName, conf); err != nil {
		return
	}

	if app, err = s.app(svr.Name); err != nil {
		return
	}
	tag := s.setTag(app.ID, svr.BuildVersion, &cacheTag{Tag: svr.Version, ConfIDs: ids})
	s.setTagID(app.ID, tag.C.Tag, svr.BuildVersion)
	if force, err = s.dao.TagForce(svr.Version); err != nil {
		return
	}
	if force == 1 {
		s.setLFVForces(svr.Version, app.ID, svr.BuildVersion)
	}
	// push hosts
	if hosts, err = s.dao.Hosts(c, svr.Name); err != nil {
		log.Error("get hosts error. svr:%s, buildVer:%s, ver:%d", svr.Name, svr.BuildVersion, svr.Version)
		err = nil
		return
	}
	for _, h := range hosts {
		if h.State == model.HostOnline {
			pushKey := pushKey(h.Service, h.Name)
			if ok := s.pubEvent(pushKey, &model.Diff{Version: svr.Version, Diffs: tag.diff()}); ok {
				log.Info("s.events.Pub(%s, %d) ok: %t", pushKey, conf.Version, ok)
			}
		}
	}
	return
}

// Config return config content.
func (s *Service) Config(c context.Context, svrName, token string, ver int64, ids []int64) (conf *model.Content, err error) {
	var (
		values, all []*model.Value
	)
	if _, err = s.appAuth(c, svrName, token); err != nil {
		return
	}
	cacheName := cacheKey(svrName, ver)
	if conf, err = s.dao.File(cacheName); err == nil {
		if len(ids) == 0 {
			return
		}
		if err = json.Unmarshal([]byte(conf.Content), &all); err != nil {
			return
		}
	} else {
		if all, err = s.values(ver); err != nil {
			return
		}
		if conf, err = genConfig(ver, all); err != nil {
			return
		}
		cacheName := cacheKey(svrName, ver)
		if err = s.dao.SetFile(cacheName, conf); err != nil {
			return
		}
		if len(ids) == 0 {
			return
		}
	}
	if len(all) == 0 {
		err = ecode.NothingFound
		return
	}
	for _, v := range all {
		for _, id := range ids {
			if v.ConfigID == id {
				values = append(values, v)
			}
		}
	}
	if len(values) == 0 {
		log.Error("Config(%s,%v) error", svrName, ids)
		err = ecode.NothingFound
		return
	}
	if conf, err = genConfig(ver, values); err != nil {
		log.Error("get config value:%s error(%v) ", values, err)
		return
	}
	return
}

// File get one file content.
func (s *Service) File(c context.Context, svr *model.Service) (val string, err error) {
	var (
		curVer, appID int64
		all           []*model.Value
		conf          *model.Content
		tag           *curTag
	)
	if appID, err = s.appAuth(c, svr.Name, svr.Token); err != nil {
		return
	}
	if svr.Version != model.UnknownVersion {
		curVer = svr.Version
	} else {
		// get current version, return if has new config version
		if tag, err = s.curTag(appID, svr.BuildVersion); err != nil {
			return
		}
		curVer = tag.cur()
	}
	cacheName := cacheKey(svr.Name, svr.Version)
	if conf, err = s.dao.File(cacheName); err == nil {
		if err = json.Unmarshal([]byte(conf.Content), &all); err != nil {
			return
		}
	} else {
		if all, err = s.values(curVer); err != nil {
			return
		}
		if conf, err = genConfig(curVer, all); err != nil {
			return
		}
		cacheName := cacheKey(svr.Name, curVer)
		if err = s.dao.SetFile(cacheName, conf); err != nil {
			return
		}
	}
	for _, v := range all {
		if v.Name == svr.File {
			val = v.Config
			break
		}
	}
	return
}

// CheckVersion check client version.
func (s *Service) CheckVersion(c context.Context, rhost *model.Host, token string) (evt chan *model.Diff, err error) {
	var (
		appID, tagID int64
		tag          *curTag
		ForceVersion int64
		tagForce     int8
		lastForce    int64
	)
	if appID, err = s.appAuth(c, rhost.Service, token); err != nil {
		return
	}
	// set heartbeat
	rhost.HeartbeatTime = xtime.Time(time.Now().Unix())
	s.dao.SetHost(c, rhost, rhost.Service)
	evt = make(chan *model.Diff, 1)

	//请开始你的表演
	tagForce, err = s.tagForce(appID, rhost.BuildVersion)
	if err != nil {
		return
	}
	rhost.Force = tagForce
	s.dao.SetHost(c, rhost, rhost.Service)
	if rhost.ConfigVersion > 0 {
		// force has a higher priority than appoint
		if ForceVersion, err = s.curForce(rhost, appID); err != nil {
			return
		}
		if ForceVersion > 0 {
			rhost.ForceVersion = ForceVersion
			s.dao.SetHost(c, rhost, rhost.Service)
			if ForceVersion != rhost.ConfigVersion {
				evt <- &model.Diff{Version: ForceVersion}
			}
			return
		}
		if lastForce, err = s.lastForce(appID, rhost.ConfigVersion, rhost.BuildVersion); err != nil {
			return
		}
		if rhost.ConfigVersion <= lastForce {
			if rhost.ConfigVersion != lastForce {
				evt <- &model.Diff{Version: lastForce}
			}
			return
		}
	}
	if rhost.Appoint > 0 {
		if rhost.Appoint != rhost.ConfigVersion {
			evt <- &model.Diff{Version: rhost.Appoint}
		}
		return
	}
	//结束表演

	// get current version, return if has new config version
	if tag, err = s.curTag(appID, rhost.BuildVersion); err != nil {
		return
	}
	if tagID = tag.cur(); tagID == 0 {
		err = ecode.NothingFound
		return
	}
	if tagID != rhost.ConfigVersion {
		ver := &model.Diff{Version: tagID}
		if rhost.ConfigVersion == tag.old() {
			ver.Diffs = tag.diff()
		}
		evt <- ver
		return
	}
	pushKey := pushKey(rhost.Service, rhost.Name)
	s.eLock.Lock()
	s.events[pushKey] = evt
	s.eLock.Unlock()
	return
}

// values get configs by tag id.
func (s *Service) values(tagID int64) (values []*model.Value, err error) {
	var (
		ids []int64
	)
	if ids, err = s.dao.ConfIDs(tagID); err != nil {
		return
	}
	values, err = s.dao.ConfigsByIDs(ids)
	return
}

// appAuth check app is auth
func (s *Service) appAuth(c context.Context, svr, token string) (appID int64, err error) {
	app, err := s.app(svr)
	if err != nil {
		return
	}
	if app.Token != token {
		err = ecode.AccessDenied
		return
	}
	appID = app.ID
	return
}

// Hosts return client hosts.
func (s *Service) Hosts(c context.Context, svr string) (hosts []*model.Host, err error) {
	return s.dao.Hosts(c, svr)
}

// Builds all builds
func (s *Service) Builds(c context.Context, svr string) (builds []string, err error) {
	var (
		app *model.App
	)
	if app, err = s.app(svr); err != nil {
		return
	}
	return s.dao.BuildsByAppID(app.ID)
}

// VersionSuccess return client versions which configuration is complete
func (s *Service) VersionSuccess(c context.Context, svr, bver string) (versions *model.Versions, err error) {

	var (
		vers []*model.ReVer
		ver  int64
		app  *model.App
	)
	if app, err = s.app(svr); err != nil {
		return
	}
	if ver, err = s.dao.TagID(app.ID, bver); err != nil {
		log.Error("BuildVersion(%v) error(%v)", svr, err)
		return
	}
	if vers, err = s.dao.Tags(app.ID); err != nil {
		return
	}
	versions = &model.Versions{
		Version: vers,
		DefVer:  ver,
	}
	return
}

// ClearHost clear service hosts.
func (s *Service) ClearHost(c context.Context, svr string) (err error) {
	return s.dao.ClearHost(c, svr)
}

// pubEvent publish a event to chan.
func (s *Service) pubEvent(key string, evt *model.Diff) (ok bool) {
	s.eLock.RLock()
	c, ok := s.events[key]
	s.eLock.RUnlock()
	if ok {
		c <- evt
	}
	return
}

// Unsub unsub a event.
func (s *Service) Unsub(svr, host string) {
	key := pushKey(svr, host)
	s.eLock.Lock()
	delete(s.events, key)
	s.eLock.Unlock()
}

func (s *Service) curTag(appID int64, build string) (tag *curTag, err error) {
	var (
		ok       bool
		tagID    int64
		ids      []int64
		tagForce int8
	)
	// get current version, return if has new config version
	verKey := buildVerKey(appID, build)
	s.tLock.RLock()
	tag, ok = s.tags[verKey]
	s.tLock.RUnlock()
	if !ok {
		if tagID, err = s.dao.TagID(appID, build); err != nil {
			return
		}
		if ids, err = s.dao.ConfIDs(tagID); err != nil {
			return
		}
		if tagForce, err = s.dao.TagForce(tagID); err != nil {
			return
		}
		tag = s.setTag(appID, build, &cacheTag{Tag: tagID, ConfIDs: ids, Force: tagForce})
	}
	return
}

func (s *Service) setTag(appID int64, bver string, cTag *cacheTag) (nTag *curTag) {
	var (
		oTag *curTag
		ok   bool
	)
	verKey := buildVerKey(appID, bver)
	nTag = &curTag{C: cTag}
	s.tLock.Lock()
	oTag, ok = s.tags[verKey]
	if ok && oTag.C != nil {
		nTag.O = oTag.C
	}
	s.tags[verKey] = nTag
	s.tLock.Unlock()
	return
}

func (s *Service) app(svr string) (app *model.App, err error) {
	var (
		ok     bool
		treeID int64
	)
	s.aLock.RLock()
	app, ok = s.apps[svr]
	s.aLock.RUnlock()
	if !ok {
		arrs := strings.Split(svr, "_")
		if len(arrs) != 3 {
			err = ecode.RequestErr
			return
		}
		if treeID, err = strconv.ParseInt(arrs[0], 10, 64); err != nil {
			return
		}
		if app, err = s.dao.AppByTree(arrs[2], arrs[1], treeID); err != nil {
			log.Error("Token(%v) error(%v)", svr, err)
			return
		}
		s.aLock.Lock()
		s.apps[svr] = app
		s.aLock.Unlock()
	}
	return app, nil
}

//SetToken set token.
func (s *Service) SetToken(svr, token string) (err error) {
	var (
		ok        bool
		treeID    int64
		app, nApp *model.App
	)
	s.aLock.RLock()
	app, ok = s.apps[svr]
	s.aLock.RUnlock()
	if ok {
		nApp = &model.App{ID: app.ID, Name: app.Name, Token: token}
	} else {
		arrs := strings.Split(svr, "_")
		if len(arrs) != 3 {
			err = ecode.RequestErr
			return
		}
		if treeID, err = strconv.ParseInt(arrs[0], 10, 64); err != nil {
			return
		}
		if nApp, err = s.dao.AppByTree(arrs[2], arrs[1], treeID); err != nil {
			log.Error("Token(%v) error(%v)", svr, err)
			return
		}
	}
	s.aLock.Lock()
	s.apps[svr] = nApp
	s.aLock.Unlock()
	return
}

//TmpBuilds get builds.
func (s *Service) TmpBuilds(svr, env string) (builds []string, err error) {
	var (
		apps   []*model.DBApp
		appIDs []int64
	)
	switch env {
	case "10":
		env = "dev"
	case "11":
		env = "fat1"
	case "13":
		env = "uat"
	case "14":
		env = "pre"
	case "3":
		env = "prod"
	default:
	}
	if apps, err = s.dao.AppsByNameEnv(svr, env); err != nil {
		return
	}
	if len(apps) == 0 {
		err = ecode.NothingFound
		return
	}
	for _, app := range apps {
		appIDs = append(appIDs, app.ID)
	}
	return s.dao.BuildsByAppIDs(appIDs)
}

// AppService ...
func (s *Service) AppService(zone, env, token string) (service string, err error) {
	var (
		ok  bool
		key string
		app *model.App
	)
	key = fmt.Sprintf("%s_%s_%s", zone, env, token)
	s.bLock.RLock()
	app, ok = s.services[key]
	s.bLock.RUnlock()
	if !ok {
		app, err = s.dao.AppGet(zone, env, token)
		if err != nil {
			log.Error("AppService error(%v)", err)
			return
		}
		s.bLock.Lock()
		s.services[key] = app
		s.bLock.Unlock()
	}
	service = fmt.Sprintf("%d_%s_%s", app.TreeID, app.Env, app.Zone)
	return
}

// Force version to clients & generate config caches
func (s *Service) Force(c context.Context, svr *model.Service, hostnames map[string]string, sType int8) (err error) {
	var (
		values []*model.Value
		ids    []int64
	)
	if sType == 1 { // push 1 clear other
		if ids, err = s.dao.ConfIDs(svr.Version); err != nil {
			return
		}
		if values, err = s.dao.ConfigsByIDs(ids); err != nil {
			return
		}
		if len(values) == 0 {
			err = fmt.Errorf("config values is empty. svr:%s, host:%s, buildVer:%s, ver:%d", svr.Name, svr.Host, svr.BuildVersion, svr.Version)
			log.Error("%v", err)
			return
		}
	}
	for key, val := range hostnames {
		pushForceKey := pushForceKey(svr.Name, key, val)
		s.pubForce(pushForceKey, svr.Version)
		log.Info("s.events.Pub(%s, %d)", pushForceKey, svr.Version)
	}
	return
}

func (s *Service) curForce(rhost *model.Host, appID int64) (version int64, err error) {
	var (
		ok bool
	)
	// get force version
	pushForceKey := pushForceKey(rhost.Service, rhost.Name, rhost.IP)
	s.fLock.RLock()
	version, ok = s.forces[pushForceKey]
	s.fLock.RUnlock()
	if !ok {
		if version, err = s.dao.Force(appID, rhost.Name); err != nil {
			return
		}
		s.fLock.Lock()
		s.forces[pushForceKey] = version
		s.fLock.Unlock()
	}
	return
}

// pubEvent publish a forces.
func (s *Service) pubForce(key string, version int64) {
	s.fLock.Lock()
	s.forces[key] = version
	s.fLock.Unlock()
}

func (s *Service) tagForce(appID int64, build string) (force int8, err error) {
	var (
		ok    bool
		tagID int64
	)
	key := tagIDkey(appID, build)
	s.tagIDLock.RLock()
	tagID, ok = s.tagID[key]
	s.tagIDLock.RUnlock()
	if !ok {
		if tagID, err = s.dao.TagID(appID, build); err != nil {
			return
		}
		s.setTagID(appID, tagID, build)
	}

	verKey := tagForceKey(appID, build, tagID)
	s.tfLock.RLock()
	force, ok = s.forceType[verKey]
	s.tfLock.RUnlock()
	if !ok {
		if force, err = s.dao.TagForce(tagID); err != nil {
			return
		}
		s.tfLock.Lock()
		s.forceType[verKey] = force
		s.tfLock.Unlock()
	}
	return
}

func (s *Service) setTagID(appID, tagID int64, build string) {
	key := tagIDkey(appID, build)
	s.tagIDLock.Lock()
	s.tagID[key] = tagID
	s.tagIDLock.Unlock()
}

func (s *Service) lastForce(appID, version int64, build string) (lastForce int64, err error) {
	var ok bool
	var buildID int64
	key := lastForceKey(appID, build)
	s.lfvLock.RLock()
	lastForce, ok = s.lfvforces[key]
	s.lfvLock.RUnlock()
	if !ok {
		if buildID, err = s.dao.BuildID(appID, build); err != nil {
			return
		}
		if lastForce, err = s.dao.LastForce(appID, buildID); err != nil {
			return
		}
		s.lfvLock.Lock()
		s.lfvforces[key] = lastForce
		s.lfvLock.Unlock()
	}
	return
}

func (s *Service) setLFVForces(lastForce, appID int64, build string) {
	key := lastForceKey(appID, build)
	s.lfvLock.Lock()
	s.lfvforces[key] = lastForce
	s.lfvLock.Unlock()
}

//CheckLatest ...
func (s *Service) CheckLatest(c context.Context, rhost *model.Host, token string) (ver int64, err error) {
	var (
		appID, tagID int64
		tag          *curTag
	)
	if appID, err = s.appAuth(c, rhost.Service, token); err != nil {
		return
	}
	// get current version, return if has new config version
	if tag, err = s.curTag(appID, rhost.BuildVersion); err != nil {
		return
	}
	if tagID = tag.cur(); tagID == 0 {
		err = ecode.NothingFound
		return
	}
	ver = tagID
	return
}

// ConfigCheck return config content.
func (s *Service) ConfigCheck(c context.Context, svrName, token string, ver int64, ids []int64) (conf *model.Content, err error) {
	var (
		values, all []*model.Value
		appID       int64
		tagAll      *model.DBTag
	)
	if appID, err = s.appAuth(c, svrName, token); err != nil {
		return
	}
	if tagAll, err = s.dao.TagAll(ver); err != nil {
		return
	}
	if appID != tagAll.AppID {
		err = ecode.RequestErr
		return
	}
	cacheName := cacheKey(svrName, ver)
	if conf, err = s.dao.File(cacheName); err == nil {
		if len(ids) == 0 {
			return
		}
		if err = json.Unmarshal([]byte(conf.Content), &all); err != nil {
			return
		}
	} else {
		if all, err = s.values(ver); err != nil {
			return
		}
		if conf, err = genConfig(ver, all); err != nil {
			return
		}
		cacheName := cacheKey(svrName, ver)
		if err = s.dao.SetFile(cacheName, conf); err != nil {
			return
		}
		if len(ids) == 0 {
			return
		}
	}
	if len(all) == 0 {
		err = ecode.NothingFound
		return
	}
	for _, v := range all {
		for _, id := range ids {
			if v.ConfigID == id {
				values = append(values, v)
			}
		}
	}
	if len(values) == 0 {
		log.Error("Config(%s,%v) error", svrName, ids)
		err = ecode.NothingFound
		return
	}
	if conf, err = genConfig(ver, values); err != nil {
		log.Error("get config value:%s error(%v) ", values, err)
		return
	}
	return
}
