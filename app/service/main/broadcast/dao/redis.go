package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"go-common/app/service/main/broadcast/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_prefixMidServer    = "mid_%d" // mid -> key:server
	_prefixKeyServer    = "key_%s" // key -> server
	_prefixServerOnline = "ol_%s"  // server -> online
	_keyServers         = "servers"
	_keyMigrateRooms    = "migrate_rooms"
	_keyMigrateServers  = "migrate_servers"
)

var (
	_redisExpire = int32(time.Minute * 30 / time.Second)
)

func keyMidServer(mid int64) string {
	return fmt.Sprintf(_prefixMidServer, mid)
}

func keyKeyServer(key string) string {
	return fmt.Sprintf(_prefixKeyServer, key)
}

func keyServerOnline(key string) string {
	return fmt.Sprintf(_prefixServerOnline, key)
}

// pingRedis check redis connection.
func (d *Dao) pingRedis(c context.Context) (err error) {
	conn := d.redis.Get(c)
	_, err = conn.Do("SET", "PING", "PONG")
	conn.Close()
	return
}

// AddMapping add a mapping.
// Mapping:
//	mid -> key_server
//	key -> server
func (d *Dao) AddMapping(c context.Context, mid int64, key, server string) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	var n = 2
	if mid > 0 {
		if err = conn.Send("HSET", keyMidServer(mid), key, server); err != nil {
			log.Error("conn.Send(HSET %d,%s,%s) error(%v)", mid, server, key, err)
			return
		}
		if err = conn.Send("EXPIRE", keyMidServer(mid), _redisExpire); err != nil {
			log.Error("conn.Send(EXPIRE %d,%s,%s) error(%v)", mid, key, server, err)
			return
		}
		n += 2
	}
	if err = conn.Send("SET", keyKeyServer(key), server); err != nil {
		log.Error("conn.Send(HSET %d,%s,%s) error(%v)", mid, server, key, err)
		return
	}
	if err = conn.Send("EXPIRE", keyKeyServer(key), _redisExpire); err != nil {
		log.Error("conn.Send(EXPIRE %d,%s,%s) error(%v)", mid, key, server, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush() error(%v)", err)
		return
	}
	for i := 0; i < n; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// ExpireMapping expire a mapping.
func (d *Dao) ExpireMapping(c context.Context, mid int64, key string) (has bool, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	var n = 1
	if mid > 0 {
		if err = conn.Send("EXPIRE", keyMidServer(mid), _redisExpire); err != nil {
			log.Error("conn.Send(EXPIRE %d,%s) error(%v)", mid, key, err)
			return
		}
		n++
	}
	if err = conn.Send("EXPIRE", keyKeyServer(key), _redisExpire); err != nil {
		log.Error("conn.Send(EXPIRE %d,%s) error(%v)", mid, key, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush() error(%v)", err)
		return
	}
	for i := 0; i < n; i++ {
		if has, err = redis.Bool(conn.Receive()); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// DelMapping del a mapping.
func (d *Dao) DelMapping(c context.Context, mid int64, key, server string) (has bool, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	n := 1
	if mid > 0 {
		if err = conn.Send("HDEL", keyMidServer(mid), key); err != nil {
			log.Error("conn.Send(HDEL %d,%s,%s) error(%v)", mid, key, server, err)
			return
		}
		n++
	}
	if err = conn.Send("DEL", keyKeyServer(key)); err != nil {
		log.Error("conn.Send(HDEL %d,%s,%s) error(%v)", mid, key, server, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush() error(%v)", err)
		return
	}
	for i := 0; i < n; i++ {
		if has, err = redis.Bool(conn.Receive()); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// ServersByKeys get a server by key.
func (d *Dao) ServersByKeys(c context.Context, keys []string) (res []string, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	var args []interface{}
	for _, key := range keys {
		args = append(args, keyKeyServer(key))
	}
	if res, err = redis.Strings(conn.Do("MGET", args...)); err != nil {
		log.Error("conn.Do(MGET %v) error(%v)", args, err)
	}
	return
}

// KeysByMids get a key server by mid.
func (d *Dao) KeysByMids(c context.Context, mids []int64) (ress map[string]string, olMids []int64, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	ress = make(map[string]string)
	for _, mid := range mids {
		if err = conn.Send("HGETALL", keyMidServer(mid)); err != nil {
			log.Error("conn.Do(HGETALL %d) error(%v)", mid, err)
			return
		}
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush() error(%v)", err)
		return
	}
	for idx := 0; idx < len(mids); idx++ {
		var (
			res map[string]string
		)
		if res, err = redis.StringMap(conn.Receive()); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
		if len(res) > 0 {
			olMids = append(olMids, mids[idx])
		}
		for k, v := range res {
			ress[k] = v
		}
	}
	return
}

// AddServerOnline add server online.
func (d *Dao) AddServerOnline(c context.Context, server string, sharding int32, online *model.Online) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	b, _ := json.Marshal(online)
	key := keyServerOnline(server)
	if err = conn.Send("HSET", key, strconv.FormatInt(int64(sharding), 10), b); err != nil {
		log.Error("conn.Send(SET %s,%d) error(%v)", key, sharding, err)
		return
	}
	if err = conn.Send("EXPIRE", key, _redisExpire); err != nil {
		log.Error("conn.Send(EXPIRE %s) error(%v)", key, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush() error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// ServerOnline get a server online.
func (d *Dao) ServerOnline(c context.Context, server string, shard int) (online *model.Online, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := keyServerOnline(server)
	hashKey := fmt.Sprint(shard)
	b, err := redis.Bytes(conn.Do("HGET", key, hashKey))
	if err != nil {
		if err != redis.ErrNil {
			log.Error("conn.Do(HGET %s %s) error(%v)", key, hashKey, err)
		} else {
			err = nil
		}
		return
	}
	online = new(model.Online)
	if err = json.Unmarshal(b, online); err != nil {
		log.Error("serverOnline json.Unmarshal(%s) error(%v)", b, err)
	}
	return
}

// DelServerOnline del a server online.
func (d *Dao) DelServerOnline(c context.Context, server string, shard int) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := keyServerOnline(server)
	hashKey := fmt.Sprint(shard)
	if _, err = conn.Do("HDEL", key, hashKey); err != nil {
		log.Error("conn.Do(DEL %s) error(%v)", key, err)
	}
	return
}

// SetServers set servers info.
func (d *Dao) SetServers(c context.Context, srvs []*model.ServerInfo) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	b, _ := json.Marshal(srvs)
	if _, err = conn.Do("SET", _keyServers, b); err != nil {
		log.Error("conn.Do(SET %s,%s) error(%v)", _keyServers, b, err)
	}
	return
}

// Servers return servers.
func (d *Dao) Servers(c context.Context) (srvs []*model.ServerInfo, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	b, err := redis.Bytes(conn.Do("GET", _keyServers))
	if err != nil {
		if err != redis.ErrNil {
			log.Error("conn.Do(GET %s) error(%v)", _keyServers, err)
		} else {
			err = nil
		}
		return
	}
	if err = json.Unmarshal(b, &srvs); err != nil {
		log.Error("MigrateServers json.Unmarshal(%s) error(%v)", b, err)
	}
	return
}

// MigrateServers migrate servers.
func (d *Dao) MigrateServers(c context.Context) (conns, ips int64, err error) {
	var servers struct {
		Conns int64 `json:"conn_count"`
		IPs   int64 `json:"ip_count"`
	}
	conn := d.redis.Get(c)
	defer conn.Close()
	b, err := redis.Bytes(conn.Do("GET", _keyMigrateServers))
	if err != nil {
		if err != redis.ErrNil {
			log.Error("conn.Do(GET %s) error(%v)", _keyMigrateServers, err)
		} else {
			err = nil
		}
		return
	}
	if err = json.Unmarshal(b, &servers); err != nil {
		log.Error("MigrateServers json.Unmarshal(%s) error(%v)", b, err)
		return
	}
	conns = servers.Conns
	ips = servers.IPs
	return
}

// MigrateRooms migrate rooms.
func (d *Dao) MigrateRooms(c context.Context, shard int) (rooms map[string]int32, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	b, err := redis.Bytes(conn.Do("HGET", _keyMigrateRooms, fmt.Sprint(shard)))
	if err != nil {
		if err != redis.ErrNil {
			log.Error("conn.Do(HGET %s,%d) error(%v)", _keyMigrateRooms, shard, err)
		} else {
			err = nil
		}
		return
	}
	rooms = make(map[string]int32)
	if err = json.Unmarshal(b, &rooms); err != nil {
		log.Error("migrateRooms json.Unmarshal() error(%v)", err)
	}
	return
}
