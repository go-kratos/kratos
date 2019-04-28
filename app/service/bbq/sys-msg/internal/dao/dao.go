package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"go-common/library/cache"
	"strconv"

	"go-common/library/cache/redis"
	xsql "go-common/library/database/sql"
	"go-common/library/log"

	"go-common/app/service/bbq/sys-msg/api/v1"
	"go-common/app/service/bbq/sys-msg/internal/conf"
)

const (
	_selectSQL    = "select id, type, sender, receiver, jump_url, text, ctime, state from sys_msg where id in (%s)"
	_insertSQL    = "insert into sys_msg (`type`,`sender`,`receiver`,`jump_url`,`text`) values (?,?,?,?)"
	_redisKey     = "sys:msg:%d"
	_redisExpireS = 600
)

//go:generate $GOPATH/src/go-common/app/tool/cache/gen
type _cache interface {
	// cache: -batch=50 -max_group=10 -batch_err=break -nullcache=&v1.SysMsg{Id:0} -check_null_code=$==nil||$.Id==0
	SysMsg(c context.Context, ids []int64) (map[int64]*v1.SysMsg, error)
}

// Dao dao
type Dao struct {
	c     *conf.Config
	cache *cache.Cache
	redis *redis.Pool
	db    *xsql.DB
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:     c,
		cache: cache.New(1, 1024),
		redis: redis.NewPool(c.Redis),
		db:    xsql.NewMySQL(c.MySQL),
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
	d.redis.Close()
	d.db.Close()
}

// Ping dao ping
func (d *Dao) Ping(ctx context.Context) error {
	// TODO: add mc,redis... if you use
	return d.db.Ping(ctx)
}

// RawSysMsg 获取系统消息
func (d *Dao) RawSysMsg(ctx context.Context, ids []int64) (res map[int64]*v1.SysMsg, err error) {
	if len(ids) == 0 {
		return
	}
	res = make(map[int64]*v1.SysMsg)

	querySQL := fmt.Sprintf(_selectSQL, intJoin(ids, ","))
	log.V(1).Infov(ctx, log.KV("sql", querySQL))
	rows, err := d.db.Query(ctx, querySQL)
	if err != nil {
		log.Errorv(ctx, log.KV("log", "query mysql sys msg fail"), log.KV("sql", querySQL))
		return
	}
	defer rows.Close()

	//"select id, type, sender, receiver, text, ctime from sys_msg where state = 0 and id in (%s)"
	for rows.Next() {
		var msg v1.SysMsg
		if err = rows.Scan(&msg.Id, &msg.Type, &msg.Sender, &msg.Receiver, &msg.JumpUrl, &msg.Text, &msg.Ctime, &msg.State); err != nil {
			log.Errorv(ctx, log.KV("log", "scan mysql sys msg fail"), log.KV("sql", querySQL))
			return
		}
		res[msg.Id] = &msg
	}

	log.V(1).Infov(ctx, log.KV("log", "get sys msg from mysql"), log.KV("req_size", len(ids)), log.KV("rsp_size", len(res)))
	return
}

// CreateSysMsg 创建系统消息
func (d *Dao) CreateSysMsg(ctx context.Context, msg *v1.SysMsg) (err error) {
	result, err := d.db.Exec(ctx, _insertSQL, msg.Type, msg.Sender, msg.Receiver, msg.JumpUrl, msg.Text)
	if err != nil {
		log.Errorv(ctx, log.KV("log", "exec mysql fail: create sys msg"), log.KV("sql", _insertSQL), log.KV("msg", msg.String()))
		return
	}
	msgID, _ := result.LastInsertId()
	d.DelCacheSysMsg(ctx, msgID)
	return
}

func intJoin(raw []int64, split string) (res string) {
	for i, v := range raw {
		if i != 0 {
			res += split
		}
		res += strconv.FormatInt(v, 10)
	}
	return
}

// CacheSysMsg .
func (d *Dao) CacheSysMsg(ctx context.Context, ids []int64) (res map[int64]*v1.SysMsg, err error) {
	res = make(map[int64]*v1.SysMsg)
	conn := d.redis.Get(ctx)
	defer conn.Close()

	for _, id := range ids {
		conn.Send("GET", fmt.Sprintf(_redisKey, id))
	}
	conn.Flush()
	for _, id := range ids {
		var by []byte
		by, err = redis.Bytes(conn.Receive())
		if err == redis.ErrNil {
			err = nil
			log.V(1).Infov(ctx, log.KV("log", "get sys msg nil from redis"), log.KV("id", id))
			continue
		}
		var msg v1.SysMsg
		if err = json.Unmarshal(by, &msg); err != nil {
			log.Errorv(ctx, log.KV("log", "unmarshal sys msg fail: str="+string(by)))
			return
		}
		res[id] = &msg
	}
	return
}

// DelCacheSysMsg 删除sys_msg缓存
func (d *Dao) DelCacheSysMsg(ctx context.Context, msgID int64) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	redisKey := fmt.Sprintf(_redisKey, msgID)
	conn.Do("DEL", redisKey)
	log.V(1).Infov(ctx, log.KV("log", "del redis_key: "+redisKey))
}

// AddCacheSysMsg 添加sys_msg缓存
func (d *Dao) AddCacheSysMsg(ctx context.Context, msg map[int64]*v1.SysMsg) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	for id, val := range msg {
		b, _ := json.Marshal(val)
		conn.Send("SETEX", fmt.Sprintf(_redisKey, id), _redisExpireS, b)
	}
	conn.Flush()
	for range msg {
		conn.Receive()
	}
}
