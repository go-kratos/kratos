package archive

import (
	"context"
	"go-common/app/interface/main/videoup/conf"
	upapi "go-common/app/service/main/up/api/v1"
	"go-common/library/cache/redis"
	bm "go-common/library/net/http/blademaster"
	"time"
)

// Dao is archive dao.
type Dao struct {
	c *conf.Config
	// http
	httpR    *bm.Client
	httpW    *bm.Client
	UpClient upapi.UpClient
	// redis
	redis       *redis.Pool
	redisExpire int32
	// uri
	viewURI        string
	addURI         string
	editURI        string
	typesURI       string
	descFormatURI  string
	tagUpURI       string
	staffConfigURI string
	applyStaffs    string

	// ad check
	porderConfigURL string
	gameListURL     string
}

const (
	_descFormatURL = "/videoup/desc/format"
	_porderConfig  = "/videoup/porder/config/list"
	_gameList      = "/game/list"
	_staffConfURI  = "/x/internal/creative/staff/config"
)

// New new a dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c: c,
		//filename redis
		redis:       redis.NewPool(c.Redis.Videoup.Config),
		redisExpire: int32(time.Duration(c.Redis.Videoup.Expire) / time.Second),
		// http client
		httpR: bm.NewClient(c.HTTPClient.Read),
		httpW: bm.NewClient(c.HTTPClient.Write),
		// uri
		viewURI:        c.Host.Archive + _viewURL,
		addURI:         c.Host.Archive + _addURL,
		editURI:        c.Host.Archive + _editURL,
		typesURI:       c.Host.Archive + _typesURL,
		descFormatURI:  c.Host.Archive + _descFormatURL,
		tagUpURI:       c.Host.Archive + _tagUpURL,
		staffConfigURI: c.Host.APICo + _staffConfURI,
		applyStaffs:    c.Host.Archive + _applyStaffs,
		// ad
		porderConfigURL: c.Host.Archive + _porderConfig,
		gameListURL:     c.Game.OpenHost + _gameList,
	}
	var err error
	if d.UpClient, err = upapi.NewClient(c.UpClient); err != nil {
		panic(err)
	}
	return d
}

// Ping ping success.
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.pingRedis(c); err != nil {
		return
	}
	return
}

// Close close resource.
func (d *Dao) Close() {
	if d.redis != nil {
		d.redis.Close()
	}
}
