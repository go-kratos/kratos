package apollo

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/philchia/agollo"

	"github.com/bilibili/kratos/pkg/conf/paladin"
)

var (
	_            paladin.Client = &apollo{}
	defaultValue                = ""
)

type apolloWatcher struct {
	keys []string // in apollo, they're called namespaces
	C    chan paladin.Event
}

func newApolloWatcher(keys []string) *apolloWatcher {
	return &apolloWatcher{keys: keys, C: make(chan paladin.Event, 5)}
}

func (aw *apolloWatcher) HasKey(key string) bool {
	if len(aw.keys) == 0 {
		return true
	}
	for _, k := range aw.keys {
		if k == key {
			return true
		}
	}
	return false
}

func (aw *apolloWatcher) Handle(event paladin.Event) {
	select {
	case aw.C <- event:
	default:
		log.Printf("paladin: event channel full discard ns %s update event", event.Key)
	}
}

// apollo is apollo config client.
type apollo struct {
	client   *agollo.Client
	values   *paladin.Map
	wmu      sync.RWMutex
	watchers map[*apolloWatcher]struct{}
}

// Config is apollo config client config.
type Config struct {
	AppID      string   `json:"app_id"`
	Cluster    string   `json:"cluster"`
	CacheDir   string   `json:"cache_dir"`
	MetaAddr   string   `json:"meta_addr"`
	Namespaces []string `json:"namespaces"`
}

type apolloDriver struct{}

var (
	confAppID, confCluster, confCacheDir, confMetaAddr, confNamespaces string
)

func init() {
	addApolloFlags()
	paladin.Register(PaladinDriverApollo, &apolloDriver{})
}

func addApolloFlags() {
	flag.StringVar(&confAppID, "apollo.appid", "", "apollo app id")
	flag.StringVar(&confCluster, "apollo.cluster", "", "apollo cluster")
	flag.StringVar(&confCacheDir, "apollo.cachedir", "/tmp", "apollo cache dir")
	flag.StringVar(&confMetaAddr, "apollo.metaaddr", "", "apollo meta server addr, e.g. localhost:8080")
	flag.StringVar(&confNamespaces, "apollo.namespaces", "", "subscribed apollo namespaces, comma separated, e.g. app.yml,mysql.yml")
}

func buildConfigForApollo() (c *Config, err error) {
	if appidFromEnv := os.Getenv("APOLLO_APP_ID"); appidFromEnv != "" {
		confAppID = appidFromEnv
	}
	if confAppID == "" {
		err = errors.New("invalid apollo appid, pass it via APOLLO_APP_ID=xxx with env or --apollo.appid=xxx with flag")
		return
	}
	if clusterFromEnv := os.Getenv("APOLLO_CLUSTER"); clusterFromEnv != "" {
		confCluster = clusterFromEnv
	}
	if confCluster == "" {
		err = errors.New("invalid apollo cluster, pass it via APOLLO_CLUSTER=xxx with env or --apollo.cluster=xxx with flag")
		return
	}
	if cacheDirFromEnv := os.Getenv("APOLLO_CACHE_DIR"); cacheDirFromEnv != "" {
		confCacheDir = cacheDirFromEnv
	}
	if metaAddrFromEnv := os.Getenv("APOLLO_META_ADDR"); metaAddrFromEnv != "" {
		confMetaAddr = metaAddrFromEnv
	}
	if confMetaAddr == "" {
		err = errors.New("invalid apollo meta addr, pass it via APOLLO_META_ADDR=xxx with env or --apollo.metaaddr=xxx with flag")
		return
	}
	if namespacesFromEnv := os.Getenv("APOLLO_NAMESPACES"); namespacesFromEnv != "" {
		confNamespaces = namespacesFromEnv
	}
	namespaceNames := strings.Split(confNamespaces, ",")
	if len(namespaceNames) == 0 {
		err = errors.New("invalid apollo namespaces, pass it via APOLLO_NAMESPACES=xxx with env or --apollo.namespaces=xxx with flag")
		return
	}
	c = &Config{
		AppID:      confAppID,
		Cluster:    confCluster,
		CacheDir:   confCacheDir,
		MetaAddr:   confMetaAddr,
		Namespaces: namespaceNames,
	}
	return
}

// New new an apollo config client.
// it watches apollo namespaces changes and updates local cache.
// BTW, in our context, namespaces in apollo means keys in paladin.
func (ad *apolloDriver) New() (paladin.Client, error) {
	c, err := buildConfigForApollo()
	if err != nil {
		return nil, err
	}
	return ad.new(c)
}

func (ad *apolloDriver) new(conf *Config) (paladin.Client, error) {
	if conf == nil {
		err := errors.New("invalid apollo conf")
		return nil, err
	}
	client := agollo.NewClient(&agollo.Conf{
		AppID:          conf.AppID,
		Cluster:        conf.Cluster,
		NameSpaceNames: conf.Namespaces, // these namespaces will be subscribed at init
		CacheDir:       conf.CacheDir,
		IP:             conf.MetaAddr,
	})
	err := client.Start()
	if err != nil {
		return nil, err
	}
	a := &apollo{
		client:   client,
		values:   new(paladin.Map),
		watchers: make(map[*apolloWatcher]struct{}),
	}
	raws, err := a.loadValues(conf.Namespaces)
	if err != nil {
		return nil, err
	}
	a.values.Store(raws)
	// watch namespaces by default.
	a.WatchEvent(context.TODO(), conf.Namespaces...)
	go a.watchproc(conf.Namespaces)
	return a, nil
}

// loadValues load values from apollo namespaces to values
func (a *apollo) loadValues(keys []string) (values map[string]*paladin.Value, err error) {
	values = make(map[string]*paladin.Value, len(keys))
	for _, k := range keys {
		if values[k], err = a.loadValue(k); err != nil {
			return
		}
	}
	return
}

// loadValue load value from apollo namespace content to value
func (a *apollo) loadValue(key string) (*paladin.Value, error) {
	content := a.client.GetNameSpaceContent(key, defaultValue)
	return paladin.NewValue(content, content), nil
}

// reloadValue reload value by key and send event
func (a *apollo) reloadValue(key string) (err error) {
	// NOTE: in some case immediately read content from client after receive event
	// will get old content due to cache, sleep 100ms make sure get correct content.
	time.Sleep(100 * time.Millisecond)
	var (
		value    *paladin.Value
		rawValue string
	)
	value, err = a.loadValue(key)
	if err != nil {
		return
	}
	rawValue, err = value.Raw()
	if err != nil {
		return
	}
	raws := a.values.Load()
	raws[key] = value
	a.values.Store(raws)
	a.wmu.RLock()
	n := 0
	for w := range a.watchers {
		if w.HasKey(key) {
			n++
			// FIXME(Colstuwjx): check change event and send detail type like EventAdd\Update\Delete.
			w.Handle(paladin.Event{Event: paladin.EventUpdate, Key: key, Value: rawValue})
		}
	}
	a.wmu.RUnlock()
	log.Printf("paladin: reload config: %s events: %d\n", key, n)
	return
}

// apollo config daemon to watch remote apollo notifications
func (a *apollo) watchproc(keys []string) {
	events := a.client.WatchUpdate()
	for {
		select {
		case event := <-events:
			if err := a.reloadValue(event.Namespace); err != nil {
				log.Printf("paladin: load key: %s error: %s, skipped", event.Namespace, err)
			}
		}
	}
}

// Get return value by key.
func (a *apollo) Get(key string) *paladin.Value {
	return a.values.Get(key)
}

// GetAll return value map.
func (a *apollo) GetAll() *paladin.Map {
	return a.values
}

// WatchEvent watch with the specified keys.
func (a *apollo) WatchEvent(ctx context.Context, keys ...string) <-chan paladin.Event {
	aw := newApolloWatcher(keys)
	err := a.client.SubscribeToNamespaces(keys...)
	if err != nil {
		log.Printf("subscribe namespaces %v failed, %v", keys, err)
		return aw.C
	}
	a.wmu.Lock()
	a.watchers[aw] = struct{}{}
	a.wmu.Unlock()
	return aw.C
}

// Close close watcher.
func (a *apollo) Close() (err error) {
	if err = a.client.Stop(); err != nil {
		return
	}
	a.wmu.RLock()
	for w := range a.watchers {
		close(w.C)
	}
	a.wmu.RUnlock()
	return
}
