package dao

import (
	"context"
	"net"
	"net/http"
	"time"

	"fmt"
	"go-common/app/service/bbq/search/conf"
	"go-common/library/cache/redis"
	xsql "go-common/library/database/sql"
	"gopkg.in/olivere/elastic.v5"
)

// Dao dao
type Dao struct {
	c          *conf.Config
	redis      *redis.Pool
	db         *xsql.DB
	esPool     map[string]*elastic.Client
	httpClient []*http.Client
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:          c,
		redis:      redis.NewPool(c.Redis),
		db:         xsql.NewMySQL(c.MySQL),
		esPool:     newEsPool(c.Es),
		httpClient: dao.createHTTPClient(),
	}
	dao.createESIndex(_bbqEsName, _videoIndex, _videoMapping)
	return
}

// newEsPool new es cluster action
func newEsPool(esInfo map[string]*conf.Es) (esCluster map[string]*elastic.Client) {
	esCluster = make(map[string]*elastic.Client)
	for esName, e := range esInfo {
		client, err := elastic.NewClient(elastic.SetURL(e.Addr...), elastic.SetSniff(false))
		if err != nil {
			panic(fmt.Sprintf("es:集群连接失败, cluster: %s, %v", esName, err))
		}
		esCluster[esName] = client
	}
	return
}

// createESIndex 初始化es索引
func (d *Dao) createESIndex(esName, index, mapping string) {
	exist, err := d.esPool[esName].IndexExists(index).Do(context.Background())
	if err != nil {
		panic(fmt.Sprintf("check if index exists, name(%s) error (%v)", esName, err))
	}
	if exist {
		return
	}
	if _, err = d.esPool[esName].CreateIndex(index).Body(mapping).Do(context.Background()); err != nil {
		panic(fmt.Sprintf("create index, name(%s) error (%v)", esName, err))
	}
}

// createHTTPClient .
func (d *Dao) createHTTPClient() []*http.Client {
	clients := make([]*http.Client, 0)
	for i := 0; i < 10; i++ {
		client := &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 100,
				IdleConnTimeout:     30 * time.Second,
			},
			Timeout: 2 * time.Second,
		}
		clients = append(clients, client)
	}
	return clients
}

// Close close the resource.
func (d *Dao) Close() {
	d.redis.Close()
	d.db.Close()
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) error {
	// TODO: if you need use mc,redis, please add
	return d.db.Ping(c)
}
