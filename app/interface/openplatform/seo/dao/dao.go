package dao

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"go-common/app/interface/openplatform/seo/conf"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_pro  = "pro"
	_item = "item"
)

// Dao dao
type Dao struct {
	c  *conf.Config
	mc *memcache.Pool
}

// New init memcache
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:  c,
		mc: memcache.NewPool(c.Memcache),
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
	d.mc.Close()
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.pingMC(c); err != nil {
		log.Error("pingMC error(%+v)", err)
		return
	}
	return
}

// pingMc ping memcache
func (d *Dao) pingMC(c context.Context) (err error) {
	con := d.mc.Get(c)
	defer con.Close()
	item := memcache.Item{Key: "ping", Value: []byte{1}, Expiration: 1}
	err = con.Set(&item)
	return
}

// getUrl get page url
// page: pro, item
func getUrl(id int, name string, bot bool) string {
	p := conf.GetPage(name)
	if p == nil {
		return ""
	}
	url := p.Bfs
	if !bot {
		url = p.Url
	}
	return fmt.Sprintf(url, id)
}

// getKey get page cache key
// name: pro, item
// return key: pro:bot:1001, pro:app:1001
// return key: item:bot:1001, item:app:1001
func getKey(id int, name string, bot bool) string {
	ua := "bot"
	if !bot {
		ua = "app"
	}
	return fmt.Sprintf("%s:%s:%d", name, ua, id)
}

// GetFile get page from file
func (d *Dao) GetFile(c context.Context, path string) (res []byte, err error) {
	log.Info(path)
	res, err = ioutil.ReadFile(path)
	return
}

// GetUrl get page from url
func (d *Dao) GetUrl(c context.Context, url string) (res []byte, err error) {
	log.Info(url)
	r, err := http.Get(url)
	if err != nil {
		return
	}
	defer r.Body.Close()
	if err == nil {
		res, err = ioutil.ReadAll(r.Body)
	}
	return
}

// GetCache get page from cache
func (d *Dao) GetCache(c context.Context, key string) (res []byte, err error) {
	log.Info(key)
	con := d.mc.Get(c)
	defer con.Close()

	item, err := con.Get(key)
	if err != nil {
		return
	}
	err = con.Scan(item, &res)
	return
}

// AddCache add page to cache
func (d *Dao) AddCache(c context.Context, key string, val []byte) (err error) {
	log.Info(key)
	item := &memcache.Item{
		Key:        key,
		Value:      val,
		Flags:      memcache.FlagRAW,
		Expiration: conf.Conf.Seo.Expire,
	}
	con := d.mc.Get(c)
	defer con.Close()
	if err = con.Set(item); err != nil {
		log.Error("key(%s) error(%v)", key, err)
	}
	return
}

// DelCache delete page cache
func (d *Dao) DelCache(c context.Context, key string) (err error) {
	con := d.mc.Get(c)
	defer con.Close()
	return con.Delete(key)
}

// ClearCache clear all page cache
func (d *Dao) ClearCache(c context.Context) (err error) {
	return
}
