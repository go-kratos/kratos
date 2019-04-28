package dao

import (
	"context"
	"encoding/json"
	"go-common/app/interface/bbq/wechat/internal/conf"
	"go-common/app/interface/bbq/wechat/internal/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
	"io/ioutil"
	"net/http"
)

// Dao dao
type Dao struct {
	c *conf.Config
	// cache      *fanout.Fanout
	redis *redis.Pool
}

// type _cache interface {
// 	cache: -nullcache=&model.InviteCode{DeviceID:""} -check_null_code=$==nil||$.DeviceID==""
// 	InviteCode(c context.Context, deviceID string) (*model.InviteCode, error)
// }

// var 常量
var (
	api    string
	jsapi  string
	appid  string
	secret string
)

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c: c,
		// cache:      fanout.New("cache", fanout.Worker(1), fanout.Buffer(1024)),
		redis: redis.NewPool(c.Redis),
	}
	api = c.URLs.Weixin
	jsapi = c.URLs.Jsapi
	appid = c.Weixin.AppID
	secret = c.Weixin.Secret
	return
}

// Close close the resource.
func (d *Dao) Close() {
	d.redis.Close()
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) (err error) {
	conn := d.redis.Get(c)
	_, err = conn.Do("SET", "PING", "PONG")
	conn.Close()
	log.V(1).Infov(c, log.KV("event", "redis_ping"))
	return
}

// TokenGet Get Token
func (d *Dao) TokenGet(c context.Context) (token string, err error) {
	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		log.Error("weixin request error (%v)", err)
	}

	q := req.URL.Query()
	q.Add("grant_type", "client_credential")
	q.Add("appid", appid)
	q.Add("secret", secret)

	req.URL.RawQuery = q.Encode()

	var resp *http.Response
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Error("weixin response error (%v)", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	res := new(model.WXToken)
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Error("weixin token parse error (%v)", err)
	}
	token, err = d.TicketGet(c, res.Token, res.Expires)
	return
}

// TicketGet Get Ticket
func (d *Dao) TicketGet(c context.Context, access string, expires int) (token string, err error) {
	req, err := http.NewRequest("GET", jsapi, nil)
	if err != nil {
		log.Error("ticket request error (%v)", err)
	}

	q := req.URL.Query()
	q.Add("access_token", access)
	q.Add("type", "jsapi")

	req.URL.RawQuery = q.Encode()

	var resp *http.Response
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Error("ticket response error (%v)", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	res := new(model.WXTicket)
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Error("ticket parse error (%v)", err)
	}
	token = res.Ticket
	err = d.TokenUpdate(c, token, expires)
	return
}

// TokenUpdate redis设置
func (d *Dao) TokenUpdate(c context.Context, token string, expires int) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	_, err = conn.Do("setex", appid, expires, token)
	if err != nil {
		log.Errorv(c, log.KV("event", "redis_set"), log.KV("key", appid), log.KV("value", token))
	}
	return
}

// TokenGetLast 获取最新token
func (d *Dao) TokenGetLast(c context.Context) (token string, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	var data []byte
	if data, err = redis.Bytes(conn.Do("get", appid)); err != nil {
		if err == redis.ErrNil {
			err = nil
			log.V(1).Infov(c, log.KV("event", "redis_get"), log.KV("key", appid), log.KV("result", "not_found"))
		} else {
			log.Errorv(c, log.KV("event", "redis_get"), log.KV("key", appid))
		}
		return
	}

	token = string(data)
	return
}
