package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

var (
	infile  string
	appname string
	host    string
)
var usage = func() {
	fmt.Fprintf(os.Stderr, "Usage of cache admin import conf tool:\n")
	flag.PrintDefaults()
}

func main() {
	flag.StringVar(&appname, "app", "", "server app name")
	flag.StringVar(&host, "h", "http://sven.bilibili.co", "sven host")
	flag.StringVar(&infile, "in", "twemproxy.yml", "server yml file")
	flag.Usage = usage
	flag.Parse()
	if appname == "" {
		panic("app name cannot be nil")
	}
	b, err := ioutil.ReadFile(infile)
	if err != nil {
		panic(err)
	}
	confs := make(map[string]server)
	err = yaml.Unmarshal(b, &confs)
	if err != nil {
		panic(err)
	}

	// one service can't use 5168 memcaches at the same time.
	redisPort := 26379
	mcPort := 21211

	for name, conf := range confs {
		var (
			addr      string
			cacheType string
		)

		if conf.Redis {
			addr = fmt.Sprintf("0.0.0.0:%d", redisPort)
			cacheType = "redis"
			redisPort++
		} else {
			addr = fmt.Sprintf("0.0.0.0:%d", mcPort)
			cacheType = "memcache"
			mcPort++
		}

		addCluster(appname, name, addr, cacheType, conf.AutoEjectHosts, conf.ServerFailureLimit)
		var nodes []*node
		for _, n := range conf.Servers {
			ss := strings.Split(n, " ")
			idx := strings.LastIndex(ss[0], ":")
			weight, _ := strconv.ParseInt(ss[0][idx+1:], 10, 64)
			nodes = append(nodes, &node{Addr: ss[0][:idx], Weigth: weight, Alias: ss[1]})
		}
		addNode(name, nodes)
	}
}

type server struct {
	AutoEjectHosts     bool     `yaml:"auto_eject_hosts"`
	Backlog            int      `yaml:"backlog"`
	Distribution       string   `yaml:"distribution"`
	Hash               string   `yaml:"hash"`
	Listen             string   `yaml:"listen"`
	Preconnect         bool     `yaml:"preconnect"`
	Timeout            int      `yaml:"timeout"`
	Redis              bool     `yaml:"redis"`
	ServerConnections  int      `yaml:"server_connections"`
	ServerFailureLimit int      `yaml:"server_failure_limit"`
	ServerRetryTimeout int      `yaml:"server_retry_timeout"`
	Servers            []string `yaml:"servers"`
}

const (
	readTimeout  = "1000"
	writeTimeout = "1000"
	conn         = "20"
	dialTimeout  = "100"
	proto        = "tcp"
)

var (
	addCluURL  = "%s/x/admin/cache/cluster/add"
	addNodeURL = "%s/x/admin/cache/cluster/node/modify"
)

func addCluster(app, cluster, addr, tp string, reject bool, fail int) {
	params := url.Values{}
	params.Add("dail_timeout", dialTimeout)
	params.Add("ping_fail_limit", strconv.FormatInt(int64(fail), 10))
	params.Add("ping_auto_reject", strconv.FormatBool(reject))
	params.Add("hash_distribution", "ketama")
	params.Add("read_timeout", readTimeout)
	params.Add("write_timeout", writeTimeout)
	params.Add("node_conn", conn)
	params.Add("type", tp)
	params.Add("app_id", app)
	params.Add("name", cluster)
	params.Add("hash_method", "fnv1a_64")
	params.Add("listen_proto", proto)
	params.Add("listen_addr", addr)
	url := fmt.Sprintf(addCluURL, host)
	req, err := http.NewRequest("POST", url, strings.NewReader(params.Encode()))
	if err != nil {
		fmt.Printf("add cluster req err %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("add cluster err", err)
		return
	}
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("add cluster req %v resp %v err %v\n", params.Encode(), string(bs), err)
}

type node struct {
	Addr   string `json:"addr"`
	Weigth int64  `json:"weight"`
	Alias  string `json:"alias"`
}

func addNode(name string, nodes []*node) {
	params := url.Values{}
	bs, err := json.Marshal(nodes)
	if err != nil {
		fmt.Println("node marshal err", err)
		return
	}
	params.Add("action", "1")
	params.Add("name", name)
	params.Add("nodes", string(bs))
	url := fmt.Sprintf(addNodeURL, host)
	req, err := http.NewRequest("POST", url, strings.NewReader(params.Encode()))
	if err != nil {
		fmt.Printf("add node req err %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("add node err", err)
		return
	}
	defer resp.Body.Close()
	bs, err = ioutil.ReadAll(resp.Body)
	fmt.Printf("add node req %v resp %v err %v\n", params.Encode(), string(bs), err)
}
