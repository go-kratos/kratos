package dao

import (
	"context"
	"net/url"
	"strconv"
	"time"

	"go-common/app/job/openplatform/article/model"
	"go-common/library/cache/redis"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_gameSyncURL = "http://line3-h5-mobile-api.biligame.com/h5/internal/article/sync"
	_typeAdd     = "1"
	_typeUpdate  = "2"
	_typeDel     = "3"
	_gameKey     = "artjob_game_mids"
)

// GameSync game sync
func (d *Dao) GameSync(c context.Context, action string, cvid int64) (err error) {
	params := url.Values{}
	params.Set("timestamp", strconv.FormatInt(time.Now().Unix()*1000, 10))
	var tp string
	if action == model.ActInsert {
		tp = _typeAdd
	} else if action == model.ActUpdate {
		tp = _typeUpdate
	} else {
		tp = _typeDel
	}
	params.Set("type", tp)
	params.Set("article_id", strconv.FormatInt(cvid, 10))
	resp := struct {
		Code int
	}{}
	if err = d.gameHTTPClient.Post(c, _gameSyncURL, "", params, &resp); err != nil {
		log.Error("game: d.gameHTTPClient.Post(%s) error(%+v)", _gameSyncURL+params.Encode(), err)
		PromError("game:同步数据")
		return
	}
	if resp.Code != 0 {
		err = ecode.Int(resp.Code)
		log.Error("game: d.gameHTTPClient.Get(%s) code: %v error(%+v)", _gameSyncURL, resp.Code, err)
		PromError("game:同步数据")
		return
	}
	log.Info("game: dao.GameSync success action: %v cvid: %v", action, cvid)
	return
}

// CacheGameList .
func (d *Dao) CacheGameList(c context.Context) (mids []int64, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	if mids, err = redis.Int64s(conn.Do("ZRANGE", _gameKey, 0, -1)); err != nil {
		PromError("redis:游戏列表缓存")
		log.Error("conn.Zrange(%s) error(%+v)", _gameKey, err)
	}
	return
}

// AddCacheGameList .
func (d *Dao) AddCacheGameList(c context.Context, mids []int64) (err error) {
	var (
		key   = _gameKey
		conn  = d.redis.Get(c)
		count int
	)
	defer conn.Close()
	if len(mids) == 0 {
		return
	}
	if err = conn.Send("DEL", key); err != nil {
		PromError("redis:删除游戏列表缓存")
		log.Error("conn.Send(DEL, %s) error(%+v)", key, err)
		return
	}
	count++
	for i, mid := range mids {
		score := i
		if err = conn.Send("ZADD", key, "CH", score, mid); err != nil {
			PromError("redis:增加游戏列表缓存")
			log.Error("conn.Send(ZADD, %s, %d, %v) error(%+v)", key, score, mid, err)
			return
		}
		count++
	}
	if err = conn.Send("EXPIRE", key, d.gameCacheExpire); err != nil {
		PromError("redis:游戏列表缓存设定过期")
		log.Error("conn.Send(EXPIRE, %s, %d) error(%+v)", key, d.gameCacheExpire, err)
		return
	}
	count++
	if err = conn.Flush(); err != nil {
		PromError("redis:增加游戏列表缓存flush")
		log.Error("conn.Flush error(%+v)", err)
		return
	}
	for i := 0; i < count; i++ {
		if _, err = conn.Receive(); err != nil {
			PromError("redis:增加游戏列表缓存receive")
			log.Error("conn.Receive error(%+v)", err)
			return
		}
	}
	return
}
