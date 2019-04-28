package v1

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"go-common/app/infra/config/model"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	xtime "go-common/library/time"
)

const (
	_buildVerKey = "%s_%s_%s"
	_pushKey     = "%s_%s_%s"
	_cacheKey    = "%s_%s_%d_%s"
	_cacheKey2   = "%s_%s_%d_%s_2"
	_fileKey     = "%s_%d"
)

var (
	addBusiness    = "前端接口api添加配置版本"
	addInfo        = "添加版本:%d"
	copyBusiness   = "前端接口api拷贝配置版本"
	copyInfo       = "拷贝版本:%d，新的版本:%d"
	updateBusiness = "前端接口api更新配置"
	updateInfo     = "更新版本:%d"
)

// PushKey push sub id
func pushKey(svr, host, env string) string {
	return fmt.Sprintf(_pushKey, svr, host, env)
}

// buildVerKey version mapping key
func buildVerKey(svr, bver, env string) string {
	return fmt.Sprintf(_buildVerKey, svr, bver, env)
}

// cacheKey config cache key
func cacheKey(svr, bver, env string, ver int64) string {
	return fmt.Sprintf(_cacheKey, svr, bver, ver, env)
}

// cacheKey config cache key
func cacheKey2(svr, bver, env string, ver int64) string {
	return fmt.Sprintf(_cacheKey2, svr, bver, ver, env)
}

// fileKey
func fileKey(filename string, ver int64) string {
	return fmt.Sprintf(_fileKey, filename, ver)
}

// tokenKey
func tokenKey(svr, env string) string {
	return fmt.Sprintf("%s_%s", svr, env)
}

// genConfig generate config
func genConfig(ver int64, cs []*model.NSValue) (conf *model.Content, err error) {
	var b []byte
	data := make(map[string]string)
	for _, c := range cs {
		data[c.Name] = c.Config
	}
	if b, err = json.Marshal(data); err != nil {
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

// genConfig2 generate config
func genConfig2(ver int64, cs []*model.NSValue, ns map[int64]string) (conf *model.Content, err error) {
	var (
		b  []byte
		v  string
		ok bool
		s  *model.Namespace
	)
	nsc := make(map[string]*model.Namespace)
	for _, c := range cs {
		if v, ok = ns[c.NamespaceID]; !ok && c.NamespaceID != 0 {
			continue
		}
		if s, ok = nsc[v]; !ok {
			s = &model.Namespace{Name: v, Data: map[string]string{}}
			nsc[v] = s
		}
		s.Data[c.Name] = c.Config
	}
	if b, err = json.Marshal(nsc); err != nil {
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
		hosts      []*model.Host
		values     []*model.NSValue
		conf       *model.Content
		conf2      *model.Content
		namespaces map[int64]string
	)
	if values, err = s.dao.Values(c, svr.Version); err != nil {
		return
	}
	if namespaces, err = s.dao.Namespaces(c, svr.Version); err != nil {
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
	cacheKey := cacheKey(svr.Name, svr.BuildVersion, svr.Env, svr.Version)
	if err = s.dao.SetFile(cacheKey, conf); err != nil {
		log.Error("set confCashe error. svr:%s, buildVer:%s, ver:%d", svr.Name, svr.BuildVersion, svr.Env)
		err = nil
	}
	if conf2, err = genConfig2(svr.Version, values, namespaces); err != nil {
		log.Error("get config2 value:%s error(%v) ", values, err)
		return
	}
	cacheKey2 := cacheKey2(svr.Name, svr.BuildVersion, svr.Env, svr.Version)
	if err = s.dao.SetFile(cacheKey2, conf2); err != nil {
		log.Error("set confCashe2 error. svr:%s, buildVer:%s, ver:%d", svr.Name, svr.BuildVersion, svr.Env)
		err = nil
	}
	s.setVersion(svr.Name, svr.BuildVersion, svr.Env, svr.Version)
	// push hosts
	if hosts, err = s.dao.Hosts(c, svr.Name, svr.Env); err != nil {
		log.Error("get hosts error. svr:%s, buildVer:%s, ver:%d", svr.Name, svr.BuildVersion, svr.Version)
		err = nil
		return
	}
	for _, h := range hosts {
		if h.State == model.HostOnline {
			pushKey := pushKey(h.Service, h.Name, svr.Env)
			if ok := s.pubEvent(pushKey, &model.Version{Version: conf.Version}); ok {
				log.Info("s.events.Pub(%s, %d) ok: %t", pushKey, conf.Version, ok)
			}
		}
	}
	return
}

// Config return config content.
func (s *Service) Config(c context.Context, svr *model.Service) (conf *model.Content, err error) {
	var values []*model.NSValue
	if err = s.appAuth(c, svr.Name, svr.Env, svr.Token); err != nil {
		return
	}
	cacheName := cacheKey(svr.Name, svr.BuildVersion, svr.Env, svr.Version)
	if conf, err = s.dao.File(cacheName); err == nil {
		return
	}
	if values, err = s.dao.Values(c, svr.Version); err != nil {
		return
	}
	if len(values) == 0 {
		err = fmt.Errorf("config values is empty. svr:%s, host:%s, buildVer:%s, ver:%d", svr.Name, svr.Host, svr.BuildVersion, svr.Version)
		log.Error("%v", err)
		return
	}
	if conf, err = genConfig(svr.Version, values); err != nil {
		log.Error("get config value:%s error(%v) ", values, err)
		return
	}
	if err = s.dao.SetFile(cacheName, conf); err != nil {
		err = nil
	}
	return
}

// Config2 return config content.
func (s *Service) Config2(c context.Context, svr *model.Service) (conf *model.Content, err error) {
	var (
		values     []*model.NSValue
		namespaces map[int64]string
	)
	if err = s.appAuth(c, svr.Name, svr.Env, svr.Token); err != nil {
		return
	}
	cacheName := cacheKey2(svr.Name, svr.BuildVersion, svr.Env, svr.Version)
	if conf, err = s.dao.File(cacheName); err == nil {
		return
	}
	if namespaces, err = s.dao.Namespaces(c, svr.Version); err != nil {
		return
	}
	if values, err = s.dao.Values(c, svr.Version); err != nil {
		return
	}
	if len(values) == 0 {
		err = fmt.Errorf("config values is empty. svr:%s, host:%s, buildVer:%s, ver:%d", svr.Name, svr.Host, svr.BuildVersion, svr.Version)
		log.Error("%v", err)
		return
	}
	if conf, err = genConfig2(svr.Version, values, namespaces); err != nil {
		log.Error("get config value:(%s) error(%v) ", values, err)
		return
	}
	if err = s.dao.SetFile(cacheName, conf); err != nil {
		err = nil
	}
	return
}

// File get one file content.
func (s *Service) File(c context.Context, svr *model.Service) (val string, err error) {
	var (
		curVer int64
		ok     bool
	)
	if err = s.appAuth(c, svr.Name, svr.Env, svr.Token); err != nil {
		return
	}
	if svr.Version != model.UnknownVersion {
		curVer = svr.Version
	} else {
		curVer, ok = s.version(svr.Name, svr.BuildVersion, svr.Env)
		if !ok {
			if curVer, err = s.dao.BuildVersion(c, svr.Name, svr.BuildVersion, svr.Env); err != nil {
				log.Error("BuildVersion(%v) error(%v)", svr, err)
				return
			}
			s.setVersion(svr.Name, svr.BuildVersion, svr.Env, curVer)
		}
	}
	fKey := fileKey(svr.File, curVer)
	if val, err = s.dao.FileStr(fKey); err == nil {
		return
	}
	if val, err = s.dao.Value(c, svr.File, curVer); err != nil {
		log.Error("Value(%v) error(%v)", svr.File, err)
		return
	}
	s.dao.SetFileStr(fKey, val)
	return
}

// CheckVersion check client version.
func (s *Service) CheckVersion(c context.Context, rhost *model.Host, env, token string) (evt chan *model.Version, err error) {
	var (
		curVer int64
	)
	if err = s.appAuth(c, rhost.Service, env, token); err != nil {
		return
	}
	// set heartbeat
	rhost.HeartbeatTime = xtime.Time(time.Now().Unix())
	if err = s.dao.SetHost(c, rhost, rhost.Service, env); err != nil {
		err = nil
	}
	evt = make(chan *model.Version, 1)
	if rhost.Appoint > 0 {
		if rhost.Appoint != rhost.ConfigVersion {
			evt <- &model.Version{Version: rhost.Appoint}
		}
		return
	}
	// get current version, return if has new config version
	if curVer, err = s.curVer(c, rhost.Service, rhost.BuildVersion, env); err != nil {
		return
	}

	if curVer == model.UnknownVersion {
		err = ecode.NothingFound
		return
	}
	if curVer != rhost.ConfigVersion {
		evt <- &model.Version{Version: curVer}
		return
	}
	pushKey := pushKey(rhost.Service, rhost.Name, env)
	s.eLock.Lock()
	s.events[pushKey] = evt
	s.eLock.Unlock()
	return
}

// AppAuth check app is auth
func (s *Service) appAuth(c context.Context, svr, env, token string) (err error) {
	var (
		dbToken  string
		ok       bool
		tokenKey = tokenKey(svr, env)
	)
	s.eLock.RLock()
	dbToken, ok = s.token[tokenKey]
	s.eLock.RUnlock()
	if !ok {
		if dbToken, err = s.dao.Token(c, svr, env); err != nil {
			log.Error("Token(%v,%v) error(%v)", svr, env, err)
			return
		}
		s.SetToken(c, svr, env, dbToken)
	}

	if dbToken != token {
		err = ecode.AccessDenied
	}
	return
}

// SetToken update Token
func (s *Service) SetToken(c context.Context, svr, env, token string) {
	tokenKey := tokenKey(svr, env)
	s.eLock.Lock()
	s.token[tokenKey] = token
	s.eLock.Unlock()
}

// Hosts return client hosts.
func (s *Service) Hosts(c context.Context, svr, env string) (hosts []*model.Host, err error) {
	return s.dao.Hosts(c, svr, env)
}

// VersionSuccess return client versions which configuration is complete
func (s *Service) VersionSuccess(c context.Context, svr, env, bver string) (versions *model.Versions, err error) {
	var (
		vers []*model.ReVer
		ver  int64
	)
	if vers, err = s.dao.Versions(c, svr, env, model.ConfigEnd); err != nil {
		log.Error("Versions(%v,%v,%v) error(%v)", svr, env, bver, err)
		return
	}
	if ver, err = s.dao.BuildVersion(c, svr, bver, env); err != nil {
		log.Error("BuildVersion(%v) error(%v)", svr, err)
		return
	}
	versions = &model.Versions{
		Version: vers,
		DefVer:  ver,
	}
	return
}

// VersionIng return client versions which configuration is creating
func (s *Service) VersionIng(c context.Context, svr, env string) (vers []int64, err error) {
	var (
		res []*model.ReVer
	)
	if res, err = s.dao.Versions(c, svr, env, model.ConfigIng); err != nil {
		log.Error("Versions(%v,%v) error(%v)", svr, env, err)
		return
	}
	vers = make([]int64, 0)
	for _, reVer := range res {
		vers = append(vers, reVer.Version)
	}
	return
}

// Builds all builds
func (s *Service) Builds(c context.Context, svr, env string) (builds []string, err error) {
	return s.dao.Builds(c, svr, env)
}

// AddConfigs insert config into db.
func (s *Service) AddConfigs(c context.Context, svr, env, token, user string, data map[string]string) (err error) {
	var (
		svrID int64
		ver   int64
	)
	if err = s.appAuth(c, svr, env, token); err != nil {
		return
	}
	if svrID, err = s.dao.ServiceID(c, svr, env); err != nil {
		return
	}
	var tx *sql.Tx
	if tx, err = s.dao.BeginTran(c); err != nil {
		log.Error("begin tran error(%v)", err)
		return
	}
	if ver, err = s.dao.TxInsertVer(tx, svrID, user); err != nil {
		tx.Rollback()
		return
	}
	if len(data) != 0 {
		if err = s.dao.TxInsertValues(c, tx, ver, user, data); err != nil {
			tx.Rollback()
			return
		}
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit error(%v)", err)
		return
	}
	s.dao.InsertLog(c, user, addBusiness, fmt.Sprintf(addInfo, ver))
	return
}

// CopyConfigs copy config in newVer.
func (s *Service) CopyConfigs(c context.Context, svr, env, token, user string, build string) (ver int64, err error) {
	var (
		svrID  int64
		curVer int64
		values []*model.NSValue
	)
	if err = s.appAuth(c, svr, env, token); err != nil {
		return
	}
	if curVer, err = s.curVer(c, svr, build, env); err != nil {
		return
	}
	if values, err = s.dao.Values(c, curVer); err != nil {
		return
	}
	data := make(map[string]string)
	for _, c := range values {
		data[c.Name] = c.Config
	}
	if svrID, err = s.dao.ServiceID(c, svr, env); err != nil {
		return
	}
	var tx *sql.Tx
	if tx, err = s.dao.BeginTran(c); err != nil {
		log.Error("begin tran error(%v)", err)
		return
	}
	if ver, err = s.dao.TxInsertVer(tx, svrID, user); err != nil {
		tx.Rollback()
		return
	}
	if err = s.dao.TxInsertValues(c, tx, ver, user, data); err != nil {
		tx.Rollback()
		return
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit error(%v)", err)
		return
	}
	s.dao.InsertLog(c, user, copyBusiness, fmt.Sprintf(copyInfo, curVer, ver))
	return
}

// UpdateConfigs update config.
func (s *Service) UpdateConfigs(c context.Context, svr, env, token, user string, ver int64, data map[string]string) (err error) {
	var (
		values  []*model.NSValue
		addData = make(map[string]string)
		udata   = make(map[string]string)
	)
	if err = s.appAuth(c, svr, env, token); err != nil {
		return
	}
	if len(data) == 0 {
		return
	}
	if values, err = s.dao.Values(c, ver); err != nil {
		return
	}
	if len(values) == 0 {
		return ecode.NothingFound
	}
	oldData := make(map[string]string)
	for _, c := range values {
		oldData[c.Name] = c.Config
	}
	for k, v := range data {
		if _, ok := oldData[k]; ok {
			udata[k] = v
		} else {
			addData[k] = v
		}
	}
	var tx *sql.Tx
	if tx, err = s.dao.BeginTran(c); err != nil {
		log.Error("begin tran error(%v)", err)
		return
	}
	if len(addData) != 0 {
		if err = s.dao.TxInsertValues(c, tx, ver, user, addData); err != nil {
			tx.Rollback()
			return
		}
	}
	if len(udata) != 0 {
		if err = s.dao.TxUpdateValues(tx, ver, user, udata); err != nil {
			tx.Rollback()
			return
		}
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit error(%v)", err)
		return
	}
	s.dao.InsertLog(c, user, updateBusiness, fmt.Sprintf(updateInfo, ver))
	return
}

func (s *Service) version(svr, bver, env string) (ver int64, ok bool) {
	verKey := buildVerKey(svr, bver, env)
	s.vLock.RLock()
	ver, ok = s.versions[verKey]
	s.vLock.RUnlock()
	return
}

func (s *Service) setVersion(svr, bver, env string, ver int64) {
	verKey := buildVerKey(svr, bver, env)
	s.vLock.Lock()
	s.versions[verKey] = ver
	s.vLock.Unlock()
}

// ClearHost clear service hosts.
func (s *Service) ClearHost(c context.Context, svr, env string) (err error) {
	return s.dao.ClearHost(c, svr, env)
}

// pubEvent publish a event to chan.
func (s *Service) pubEvent(key string, evt *model.Version) (ok bool) {
	s.eLock.RLock()
	c, ok := s.events[key]
	s.eLock.RUnlock()
	if ok {
		c <- evt
	}
	return
}

// Unsub unsub a event.
func (s *Service) Unsub(svr, host, env string) {
	key := pushKey(svr, host, env)
	s.eLock.Lock()
	delete(s.events, key)
	s.eLock.Unlock()
}

func (s *Service) curVer(c context.Context, svr, build, env string) (ver int64, err error) {
	var ok bool
	// get current version, return if has new config version
	ver, ok = s.version(svr, build, env)
	if !ok {
		if ver, err = s.dao.BuildVersion(c, svr, build, env); err != nil {
			return
		}
		s.setVersion(svr, build, env, ver)
	}
	return
}
