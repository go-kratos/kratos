package paladin

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"sync"
	"time"

	"go-common/library/conf/env"
	"go-common/library/ecode"
	"go-common/library/log"
	xip "go-common/library/net/ip"
	"go-common/library/net/netutil"

	"github.com/pkg/errors"
)

const (
	_apiGet   = "http://%s/config/v2/get?%s"
	_apiCheck = "http://%s/config/v2/check?%s"

	_maxLoadRetries = 3
)

var (
	_ Client = &sven{}

	svenHost    string
	svenVersion string
	svenPath    string
	svenToken   string
	svenAppoint string
	svenTreeid  string

	_debug bool
)

func init() {
	flag.StringVar(&svenHost, "conf_host", os.Getenv("CONF_HOST"), `config api host.`)
	flag.StringVar(&svenVersion, "conf_version", os.Getenv("CONF_VERSION"), `app version.`)
	flag.StringVar(&svenPath, "conf_path", os.Getenv("CONF_PATH"), `config file path.`)
	flag.StringVar(&svenToken, "conf_token", os.Getenv("CONF_TOKEN"), `config token.`)
	flag.StringVar(&svenAppoint, "conf_appoint", os.Getenv("CONF_APPOINT"), `config appoint.`)
	flag.StringVar(&svenTreeid, "tree_id", os.Getenv("TREE_ID"), `tree id.`)

	if env.DeployEnv == env.DeployEnvDev {
		_debug = true
	}
}

type watcher struct {
	keys []string
	ch   chan Event
}

func newWatcher(keys []string) *watcher {
	return &watcher{keys: keys, ch: make(chan Event, 5)}
}

func (w *watcher) HasKey(key string) bool {
	if len(w.keys) == 0 {
		return true
	}
	for _, k := range w.keys {
		if k == key {
			return true
		}
	}
	return false
}

func (w *watcher) Handle(event Event) {
	select {
	case w.ch <- event:
	default:
		log.Error("paladin: discard event:%+v", event)
	}
}

func (w *watcher) Chan() <-chan Event {
	return w.ch
}

func (w *watcher) Close() {
	close(w.ch)
}

// sven is sven config client.
type sven struct {
	values   *Map
	wmu      sync.RWMutex
	watchers map[*watcher]struct{}

	httpCli *http.Client
	backoff *netutil.BackoffConfig
}

// NewSven new a config client.
func NewSven() (Client, error) {
	s := &sven{
		values:   new(Map),
		watchers: make(map[*watcher]struct{}),
		httpCli:  &http.Client{Timeout: 60 * time.Second},
		backoff: &netutil.BackoffConfig{
			MaxDelay:  5 * time.Second,
			BaseDelay: 1.0 * time.Second,
			Factor:    1.6,
			Jitter:    0.2,
		},
	}
	if err := s.checkEnv(); err != nil {
		return nil, err
	}
	ver, err := s.load()
	if err != nil {
		return nil, err
	}
	go s.watchproc(ver)
	return s, nil
}

func (s *sven) checkEnv() error {
	if svenHost == "" || svenVersion == "" || svenPath == "" || svenToken == "" || svenTreeid == "" {
		return fmt.Errorf("config env invalid. conf_host(%s) conf_version(%s) conf_path(%s) conf_token(%s) conf_appoint(%s) tree_id(%s)", svenHost, svenVersion, svenPath, svenToken, svenAppoint, svenTreeid)
	}
	return nil
}

// Get return value by key.
func (s *sven) Get(key string) *Value {
	return s.values.Get(key)
}

// GetAll return value map.
func (s *sven) GetAll() *Map {
	return s.values
}

// WatchEvent watch with the specified keys.
func (s *sven) WatchEvent(ctx context.Context, keys ...string) <-chan Event {
	w := newWatcher(keys)
	s.wmu.Lock()
	s.watchers[w] = struct{}{}
	s.wmu.Unlock()
	return w.Chan()
}

// Close close watcher.
func (s *sven) Close() (err error) {
	s.wmu.RLock()
	for w := range s.watchers {
		w.Close()
	}
	s.wmu.RUnlock()
	return
}

func (s *sven) fireEvent(event Event) {
	s.wmu.RLock()
	for w := range s.watchers {
		if w.HasKey(event.Key) {
			w.Handle(event)
		}
	}
	s.wmu.RUnlock()
}

func (s *sven) load() (ver int64, err error) {
	var (
		v  *version
		cs []*content
	)
	if v, err = s.check(-1); err != nil {
		log.Error("paladin: s.check(-1) error(%v)", err)
		return
	}
	for i := 0; i < _maxLoadRetries; i++ {
		if cs, err = s.config(v); err == nil {
			all := make(map[string]*Value, len(cs))
			for _, v := range cs {
				all[v.Name] = &Value{val: v.Config, raw: v.Config}
			}
			s.values.Store(all)
			return v.Version, nil
		}
		log.Error("paladin: s.config(%v) error(%v)", ver, err)
		time.Sleep(s.backoff.Backoff(i))
	}
	return 0, err
}

func (s *sven) watchproc(ver int64) {
	var retry int
	for {
		v, err := s.check(ver)
		if err != nil {
			if ecode.NotModified.Equal(err) {
				time.Sleep(time.Second)
				continue
			}
			log.Error("paladin: s.check(%d) error(%v)", ver, err)
			retry++
			time.Sleep(s.backoff.Backoff(retry))
			continue
		}
		cs, err := s.config(v)
		if err != nil {
			log.Error("paladin: s.config(%v) error(%v)", ver, err)
			retry++
			time.Sleep(s.backoff.Backoff(retry))
			continue
		}
		all := s.values.Load()
		news := make(map[string]*Value, len(cs))
		for _, v := range cs {
			if _, ok := all[v.Name]; !ok {
				go s.fireEvent(Event{Event: EventAdd, Key: v.Name, Value: v.Config})
			} else if v.Config != "" {
				go s.fireEvent(Event{Event: EventUpdate, Key: v.Name, Value: v.Config})
			} else {
				go s.fireEvent(Event{Event: EventRemove, Key: v.Name, Value: v.Config})
			}
			news[v.Name] = &Value{val: v.Config, raw: v.Config}
		}
		for k, v := range all {
			if _, ok := news[k]; !ok {
				news[k] = v
			}
		}
		s.values.Store(news)
		ver = v.Version
		retry = 0
	}
}

type version struct {
	Version int64   `json:"version"`
	Diffs   []int64 `json:"diffs"`
}

type config struct {
	Version int64  `json:"version"`
	Content string `json:"content"`
	Md5     string `json:"md5"`
}

type content struct {
	Cid    int64  `json:"cid"`
	Name   string `json:"name"`
	Config string `json:"config"`
}

func (s *sven) check(ver int64) (v *version, err error) {
	params := newParams()
	params.Set("version", strconv.FormatInt(ver, 10))
	params.Set("appoint", svenAppoint)
	var res struct {
		Code int      `json:"code"`
		Data *version `json:"data"`
	}
	uri := fmt.Sprintf(_apiCheck, svenHost, params.Encode())
	if _debug {
		fmt.Printf("paladin: check(%d) uri(%s)\n", ver, uri)
	}
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return
	}
	resp, err := s.httpCli.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = errors.Errorf("paladin: httpCli.GET(%s) error(%d)", params.Encode(), resp.StatusCode)
		return
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if err = json.Unmarshal(b, &res); err != nil {
		return
	}
	if ec := ecode.Int(res.Code); !ec.Equal(ecode.OK) {
		err = ec
		return
	}
	if res.Data == nil {
		err = errors.Errorf("paladin: http version is nil. params(%s)", params.Encode())
		return
	}
	v = res.Data
	return
}

func (s *sven) config(ver *version) (cts []*content, err error) {
	ids, _ := json.Marshal(ver.Diffs)
	params := newParams()
	params.Set("version", strconv.FormatInt(ver.Version, 10))
	params.Set("ids", string(ids))
	var res struct {
		Code int     `json:"code"`
		Data *config `json:"data"`
	}
	uri := fmt.Sprintf(_apiGet, svenHost, params.Encode())
	if _debug {
		fmt.Printf("paladin: config(%+v) uri(%s)\n", ver, uri)
	}
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return
	}
	resp, err := s.httpCli.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = errors.Errorf("paladin: httpCli.GET(%s) error(%d)", params.Encode(), resp.StatusCode)
		return
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if err = json.Unmarshal(b, &res); err != nil {
		return
	}
	if !ecode.Int(res.Code).Equal(ecode.OK) || res.Data == nil {
		err = errors.Errorf("paladin: http config is nil. params(%s) ecode(%d)", params.Encode(), res.Code)
		return
	}
	if err = json.Unmarshal([]byte(res.Data.Content), &cts); err != nil {
		return
	}
	for _, c := range cts {
		if err = ioutil.WriteFile(path.Join(svenPath, c.Name), []byte(c.Config), 0644); err != nil {
			return
		}
	}
	return
}

func newParams() url.Values {
	params := url.Values{}
	params.Set("service", serviceName())
	params.Set("build", svenVersion)
	params.Set("token", svenToken)
	params.Set("hostname", env.Hostname)
	params.Set("ip", ipAddr())
	return params
}

func ipAddr() string {
	if env.IP != "" {
		return env.IP
	}
	return xip.InternalIP()
}

func serviceName() string {
	return fmt.Sprintf("%s_%s_%s", svenTreeid, env.DeployEnv, env.Zone)
}
