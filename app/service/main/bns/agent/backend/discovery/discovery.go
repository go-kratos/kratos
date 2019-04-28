package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"go-common/library/log"
	"go-common/library/naming"
	"go-common/library/stat/prom"

	"go-common/app/service/main/bns/agent/backend"
)

func init() {
	backend.Registry("discovery", New)
}

const (
	// NodeStatusUP Ready to receive register
	NodeStatusUP NodeStatus = iota
	// NodeStatusLost lost with each other
	NodeStatusLost
	defaultCacheExpireIn int64 = 10
)

// NodeStatus Status of instance
type NodeStatus int

// ServerNode backend servier node status
type ServerNode struct {
	Addr   string     `json:"addr"`
	Zone   string     `json:"zone"`
	Status NodeStatus `json:"status"`
}

// InstanceList discovery instance list
type InstanceList struct {
	Instances       []naming.Instance `json:"instances"`
	LatestTimestamp int64             `json:"latest_timestamp"`
}

// InstanceMetadata discovery instance metadata
type InstanceMetadata struct {
	Provider interface{} `json:"provider"`
}

type discoveryResp struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

type discoveryFetchResp struct {
	discoveryResp `json:",inline"`
	Data          *InstanceList `json:"data"`
}

type discoveryNodeResp struct {
	discoveryResp `json:",inline"`
	Data          []ServerNode `json:"data"`
}

var urlParse = url.Parse

// New discovery backend
func New(config map[string]interface{}) (backend.Backend, error) {
	if config == nil {
		return nil, fmt.Errorf("discovery require url, secret, appkey")
	}
	var url, secret, appKey string
	var ok bool

	discoveryURL := os.Getenv("DISCOVERY_URL")
	if url, ok = config["url"].(string); !ok && discoveryURL == "" {
		return nil, fmt.Errorf("discovery require `url`")
	}
	// use env DISCOVERY_URL overwrite config
	if discoveryURL != "" {
		url = discoveryURL
	}

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}

	refreshInterval := time.Minute
	if second, ok := config["refresh_interval"].(int); ok {
		refreshInterval = time.Duration(second) * time.Second
	}
	cacheExpireIn := defaultCacheExpireIn
	if second, ok := config["cacheExpireIn"].(int); ok {
		cacheExpireIn = int64(second)
	}
	u, err := urlParse(url)
	if err != nil {
		return nil, err
	}
	dis := &discovery{
		client:         http.DefaultClient,
		scheme:         u.Scheme,
		host:           u.Host,
		discoveryHosts: []string{u.Host},
		secret:         secret,
		appKey:         appKey,
		refreshTick:    time.NewTicker(refreshInterval),
		cacheMap:       make(map[string]*cacheNode),
		cacheExpireIn:  cacheExpireIn,
	}
	go dis.daemon()
	return dis, nil
}

type cacheNode struct {
	expired int64
	data    []*backend.Instance
}

func (c *cacheNode) IsExpired() bool {
	return time.Now().Unix() > c.expired
}

func newCacheNode(data []*backend.Instance, expireIn int64) *cacheNode {
	return &cacheNode{expired: time.Now().Unix() + expireIn, data: data}
}

type discovery struct {
	client *http.Client

	secret string
	appKey string

	host   string
	scheme string

	discoveryHosts []string
	rmx            sync.RWMutex

	cachermx      sync.RWMutex
	cacheMap      map[string]*cacheNode
	cacheExpireIn int64

	refreshTick *time.Ticker
}

var _ backend.Backend = &discovery{}

func (d *discovery) Ping(ctx context.Context) error {
	_, err := d.Nodes(ctx)
	return err
}

func (d *discovery) Close(ctx context.Context) error {
	d.refreshTick.Stop()
	return nil
}

func (d *discovery) daemon() {
	for range d.refreshTick.C {
		log.V(10).Info("refresh discovery nodes ...")
		nodes, err := d.Nodes(context.Background())
		if err != nil {
			log.Error("refresh discovery nodes error %s", err)
			continue
		}
		hosts := make([]string, 0, len(nodes))
		for i := range nodes {
			if nodes[i].Status == NodeStatusUP {
				hosts = append(hosts, nodes[i].Addr)
			}
		}
		d.rmx.Lock()
		d.discoveryHosts = hosts
		d.rmx.Unlock()
		log.V(10).Info("new discovery nodes list %v", hosts)
	}
}

func (d *discovery) Nodes(ctx context.Context) ([]ServerNode, error) {
	req, err := http.NewRequest(http.MethodGet, "/discovery/nodes", nil)
	if err != nil {
		return nil, err
	}

	resp, err := d.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	nodesResp := &discoveryNodeResp{}
	if err := json.NewDecoder(resp.Body).Decode(nodesResp); err != nil {
		return nil, err
	}

	if nodesResp.Code < 0 || nodesResp.Data == nil || len(nodesResp.Data) == 0 {
		return nil, fmt.Errorf("no data found, err: %s", nodesResp.Message)
	}

	log.V(10).Info("responsed data: %v", nodesResp.Data)
	return nodesResp.Data, nil
}

func (d *discovery) do(req *http.Request) (resp *http.Response, err error) {
	req.URL.Scheme = d.scheme
	req.Host = d.host
	req.URL.Scheme = d.scheme
	req.Host = d.host
	d.rmx.RLock()
	hosts := d.discoveryHosts
	d.rmx.RUnlock()
	for _, host := range shuffle(hosts) {
		req.URL.Host = host
		resp, err = d.client.Do(req)
		log.V(5).Info("request discovery.. request: %s", req.URL)
		if err == nil {
			return
		}
		log.Error("request discovery %s err %s, try next", host, err)
	}
	return
}

// Query appid
func (d *discovery) Query(ctx context.Context, target backend.Target, sel backend.Selector, md backend.Metadata) ([]*backend.Instance, error) {
	// TODO: parallel query
	key := target.Name + sel.String()
	if data, ok := d.fromCache(key); ok {
		prom.CacheHit.Incr("bns:discovery_mem_cache_hit")
		return data, nil
	}
	prom.CacheMiss.Incr("bns:discovery_mem_cache_miss")
	instanceList, err := d.fetch(ctx, target.Name, sel, md)
	if err != nil {
		return nil, err
	}
	data := copyInstance(instanceList)
	if len(data) != 0 {
		d.setCache(key, data)
	}
	return data, nil
}

func (d *discovery) fromCache(key string) ([]*backend.Instance, bool) {
	d.cachermx.RLock()
	defer d.cachermx.RUnlock()
	node, ok := d.cacheMap[key]
	if !ok || node.IsExpired() {
		return nil, false
	}
	return node.data, true
}

func (d *discovery) setCache(key string, data []*backend.Instance) {
	d.cachermx.Lock()
	defer d.cachermx.Unlock()
	d.cacheMap[key] = newCacheNode(data, d.cacheExpireIn)
}

func (d *discovery) fetch(ctx context.Context, discoveryID string, sel backend.Selector, md backend.Metadata) (*InstanceList, error) {
	params := url.Values{}
	params.Add("appid", discoveryID)
	params.Add("env", sel.Env)
	params.Add("region", sel.Region)
	params.Add("hostname", md.ClientHost)
	params.Add("zone", sel.Zone)
	params.Add("status", "1")

	if md.LatestTimestamps != "" {
		params.Add("latest_timestamp", md.LatestTimestamps)
	} else {
		params.Add("latest_timestamp", "0")
	}

	payload := params.Encode()
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/discovery/fetch?%s", payload), nil)
	if err != nil {
		return nil, err
	}

	resp, err := d.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	fetchResp := &discoveryFetchResp{}
	if err := json.NewDecoder(resp.Body).Decode(fetchResp); err != nil {
		return nil, err
	}

	if fetchResp.Code < 0 || fetchResp.Data == nil || fetchResp.Data.Instances == nil || len(fetchResp.Data.Instances) == 0 {
		return nil, fmt.Errorf("no data found, err: %s", fetchResp.Message)
	}

	log.V(10).Info("fetchResponsed data: %v", fetchResp.Data)

	return fetchResp.Data, nil
}

func shuffle(hosts []string) []string {
	for i := 0; i < len(hosts); i++ {
		j := rand.Intn(len(hosts) - i)
		hosts[i], hosts[i+j] = hosts[i+j], hosts[i]
	}
	return hosts
}

func copyInstance(il *InstanceList) (inss []*backend.Instance) {
copyloop:
	for _, in := range il.Instances {
		out := &backend.Instance{
			DiscoveryID: in.AppID,
			Env:         in.Env,
			Hostname:    in.Hostname,
			Zone:        in.Zone,
		}
		for _, addr := range in.Addrs {
			ip, err := ipFromURI(addr)
			if err == nil {
				out.IPAddr = ip
				inss = append(inss, out)
				continue copyloop
			}
			log.Error("extract ip from addr %s error: %s", addr, err)
		}
		log.Error("can't found any ip for discoveryID: %s", in.AppID)
	}
	return
}

// extract ip from uri
func ipFromURI(uri string) (net.IP, error) {
	var hostport string
	if u, err := url.Parse(uri); err != nil {
		hostport = uri
	} else {
		hostport = u.Host
	}
	if strings.ContainsRune(hostport, ':') {
		host, _, err := net.SplitHostPort(hostport)
		if err != nil {
			return net.IPv4zero, err
		}
		return net.ParseIP(host), nil
	}
	return net.ParseIP(hostport), nil
}
