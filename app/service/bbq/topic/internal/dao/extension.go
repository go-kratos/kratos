package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"go-common/app/service/bbq/topic/api"
	"go-common/library/cache/redis"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_selectExtension = "select svid, content from extension where svid in (%s)"
	_insertExtension = "insert ignore into extension (`svid`,`type`,`content`) values (?,?,?)"
)
const (
	_videoExtensionKey = "ext:%d"
)

// RawVideoExtension 从mysql获取extension
func (d *Dao) RawVideoExtension(ctx context.Context, svids []int64) (res map[int64]*api.VideoExtension, err error) {
	res = make(map[int64]*api.VideoExtension)
	if len(svids) == 0 {
		return
	}

	querySQL := fmt.Sprintf(_selectExtension, xstr.JoinInts(svids))
	rows, err := d.db.Query(ctx, querySQL)
	if err != nil {
		log.Errorw(ctx, "log", "get extension error", "err", err, "sql", querySQL)
		return
	}
	defer rows.Close()

	var svid int64
	var content string
	for rows.Next() {
		if err = rows.Scan(&svid, &content); err != nil {
			log.Errorw(ctx, "log", "get extension from mysql fail", "sql", querySQL)
			return
		}
		// 由于数据库中的数据和缓存中还不太一样，因此这里需要对db读取的数据进行额外处理
		var extension api.Extension
		json.Unmarshal([]byte(content), &extension.TitleExtra)
		// TODO：check
		log.V(10).Infow(ctx, "log", "unmarshal content", "result", extension)
		data, _ := json.Marshal(&extension)
		res[svid] = &api.VideoExtension{Svid: svid, Extension: string(data)}
	}
	log.V(1).Infow(ctx, "log", "get extension", "req", svids, "rsp_size", len(res))
	return
}

// CacheVideoExtension 从缓存获取extension
func (d *Dao) CacheVideoExtension(ctx context.Context, svids []int64) (res map[int64]*api.VideoExtension, err error) {
	res = make(map[int64]*api.VideoExtension)

	conn := d.redis.Get(ctx)
	defer conn.Close()
	for _, svid := range svids {
		conn.Send("GET", fmt.Sprintf(_videoExtensionKey, svid))
	}
	conn.Flush()
	var data string
	for _, svid := range svids {
		if data, err = redis.String(conn.Receive()); err != nil {
			if err == redis.ErrNil {
				err = nil
			} else {
				log.Errorv(ctx, log.KV("event", "redis_get"), log.KV("svid", svid))
			}
			continue
		}
		extension := new(api.VideoExtension)
		extension.Svid = svid
		extension.Extension = data
		res[extension.Svid] = extension
	}
	log.Infov(ctx, log.KV("event", "redis_get"), log.KV("row_num", len(res)))
	return
}

// AddCacheVideoExtension 添加extension缓存
func (d *Dao) AddCacheVideoExtension(ctx context.Context, extensions map[int64]*api.VideoExtension) (err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	for svid, value := range extensions {
		conn.Send("SET", fmt.Sprintf(_videoExtensionKey, svid), value.Extension, "EX", d.topicExpire)
	}
	conn.Flush()
	for i := 0; i < len(extensions); i++ {
		conn.Receive()
	}
	log.Infov(ctx, log.KV("event", "redis_set"), log.KV("row_num", len(extensions)))
	return
}

// DelCacheVideoExtension 删除extension缓存
func (d *Dao) DelCacheVideoExtension(ctx context.Context, svid int64) {
	var key = fmt.Sprintf(_videoExtensionKey, svid)
	conn := d.redis.Get(ctx)
	defer conn.Close()
	conn.Do("DEL", key)
}

// InsertExtension 插入extension到db
func (d *Dao) InsertExtension(ctx context.Context, svid int64, extensionType int64, extension *api.Extension) (rowsAffected int64, err error) {
	data, _ := json.Marshal(extension.TitleExtra)
	res, err := d.db.Exec(ctx, _insertExtension, svid, extensionType, string(data))
	if err != nil {
		log.Errorw(ctx, "log", "insert extension db fail", "svid", svid, "extension_type", extensionType, "extension", extensionType)
		return
	}
	rowsAffected, tmpErr := res.RowsAffected()
	if tmpErr != nil {
		log.Warnw(ctx, "log", "get rows affected fail")
	}
	d.DelCacheVideoExtension(ctx, svid)
	return
}
