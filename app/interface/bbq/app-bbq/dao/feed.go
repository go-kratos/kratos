package dao

import (
	"context"
	"fmt"
	"go-common/app/interface/bbq/app-bbq/model"
	"go-common/library/cache/redis"
	"go-common/library/time"
	"go-common/library/xstr"
	"strconv"

	"go-common/library/log"
)

// FetchAvailableOutboxList 根据提供的up主mids去获取比lastSvID还小的svid
// @param CursorID 			函数并不判断SvID正确性，由调用方保证
// @param cursorNext	表明fetch的方向，同时会影响排序顺序
// 							true：	则有如下条件sql"... and id < {{svid}} order by id desc..."
//							false：	则有如下条件sql"... and id > {{svid}} order by id asc..."
// @return svIDs		注意：svid的返回有顺序
func (d *Dao) FetchAvailableOutboxList(c context.Context, fetchStates string, mids []int64, cursorNext bool, cursorSvID int64, cursorPubtime time.Time, size int) (svIDs []int64, err error) {
	if len(mids) == 0 {
		return
	}
	compareSymbol := string(">=")
	orderDirection := "asc"
	if cursorNext {
		compareSymbol = "<="
		orderDirection = "desc"
	}
	midStr := xstr.JoinInts(mids)
	// 多个字段order by，每个字段都要指定desc，否则不带的字段为asc
	querySQL := fmt.Sprintf("select svid, pubtime from video where mid in (%s) and state in (%s) and "+
		"pubtime %s ? order by pubtime %s, svid %s limit %d",
		midStr, fetchStates, compareSymbol, orderDirection, orderDirection, size)
	log.V(1).Infov(c, log.KV("event", "mysql_select"), log.KV("table", "video"),
		log.KV("sql", querySQL))
	rows, err := d.db.Query(c, querySQL, cursorPubtime)
	if err != nil {
		log.Errorv(c, log.KV("event", "mysql_select"), log.KV("table", "video"),
			log.KV("sql", querySQL))
		return
	}
	defer rows.Close()
	var svID int64
	var pubTIme time.Time
	conflict := bool(true)
	for rows.Next() {
		if err = rows.Scan(&svID, &pubTIme); err != nil {
			log.Errorv(c, log.KV("event", "mysql_scan"), log.KV("table", "video"),
				log.KV("sql", querySQL))
			return
		}
		// 为了解决同一个pubtime的冲突问题
		if pubTIme == cursorPubtime && conflict {
			if svID == cursorSvID {
				conflict = false
			}
			continue
		}
		svIDs = append(svIDs, svID)
	}
	log.Infov(c, log.KV("event", "mysql_select"), log.KV("table", "video"),
		log.KV("desc", "fetch outbox"), log.KV("size", len(svIDs)))
	return
}

// SetMIDLastPubtime redis设置
func (d *Dao) SetMIDLastPubtime(c context.Context, mid int64, pubtime int64) (err error) {
	key := fmt.Sprintf(model.CacheKeyLastPubtime, mid)
	conn := d.redis.Get(c)
	defer conn.Close()
	_, err = conn.Do("set", key, pubtime, "ex", model.CacheExpireLastPubtime)
	if err != nil {
		log.Errorv(c, log.KV("event", "redis_set"), log.KV("key", key), log.KV("value", pubtime))
	}
	return
}

// GetMIDLastPubtime 获取该用户feed中的最新浏览svid
func (d *Dao) GetMIDLastPubtime(c context.Context, mid int64) (pubtime int64, err error) {
	pubtime = 0
	key := fmt.Sprintf(model.CacheKeyLastPubtime, mid)
	conn := d.redis.Get(c)
	defer conn.Close()
	var data []byte
	if data, err = redis.Bytes(conn.Do("get", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
			log.V(1).Infov(c, log.KV("event", "redis_get"), log.KV("key", key), log.KV("result", "not_found"))
		} else {
			log.Errorv(c, log.KV("event", "redis_get"), log.KV("key", key))
		}
		return
	}

	pubtime, err = strconv.ParseInt(string(data), 10, 0)
	if err != nil {
		log.Errorv(c, log.KV("event", "redis_get"), log.KV("key", key), log.KV("value", data))
	}
	return
}
