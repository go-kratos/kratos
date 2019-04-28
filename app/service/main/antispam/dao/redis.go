package dao

import (
	"context"
	"fmt"

	"go-common/app/service/main/antispam/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_regexpsKey           = "regexps"
	_localCountsKey       = "resource_id:%d:keyword_id:%d:local_limit_counts"
	_totalCountsKey       = "keyword_id:%d:total_counts"
	_globalCountsKey      = "keyword_id:%d:global_limit_counts"
	_keywordsSenderIDsKey = "keyword_id:%d:sender_ids"
	_rulesKey             = "rule:area:%s:limit_type:%s"
	_areaSendersKey       = "area:%s:sender_id:%d"
)

func sendersKey(keywordID int64) string {
	return fmt.Sprintf(_keywordsSenderIDsKey, keywordID)
}

func areaSendersKey(area string, senderID int64) string {
	return fmt.Sprintf(_areaSendersKey, area, senderID)
}

func totalCountsKey(keywordID int64) string {
	return fmt.Sprintf(_totalCountsKey, keywordID)
}

func localCountsKey(keywordID, oid int64) string {
	return fmt.Sprintf(_localCountsKey, oid, keywordID)
}

func globalCountsKey(keywordID int64) string {
	return fmt.Sprintf(_globalCountsKey, keywordID)
}

func rulesKey(area, limitType string) string {
	return fmt.Sprintf(_rulesKey, area, limitType)
}

// pingRedis check redis connection
func (d *Dao) pingRedis(c context.Context) (err error) {
	conn := d.redis.Get(c)
	_, err = conn.Do("SET", "PING", "PONG")
	conn.Close()
	return
}

// CntSendersCache .
func (d *Dao) CntSendersCache(c context.Context, keywordID int64) (cnt int64, err error) {
	var (
		key  = sendersKey(keywordID)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if cnt, err = redis.Int64(conn.Do("ZCARD", key)); err != nil {
		log.Error("redis.Int64(conn.Do(ZCARD, %s)) error(%v)", key, err)
	}
	return
}

// GlobalLocalLimitCache .
func (d *Dao) GlobalLocalLimitCache(c context.Context, keywordID, oid int64) ([]int64, error) {
	var (
		globalKey = globalCountsKey(keywordID)
		localKey  = localCountsKey(keywordID, oid)
		conn      = d.redis.Get(c)
	)
	defer conn.Close()
	if err := conn.Send("GET", globalKey); err != nil {
		log.Error("%v", err)
		return nil, err
	}
	if err := conn.Send("GET", localKey); err != nil {
		log.Error("%v", err)
		return nil, err
	}
	if err := conn.Flush(); err != nil {
		log.Error("%v", err)
		return nil, err
	}
	counts := make([]int64, 0)
	for i := 0; i < 2; i++ {
		count, err := redis.Int64(conn.Receive())
		if err == nil || err == redis.ErrNil {
			counts = append(counts, count)
			continue
		}
		log.Error("%v", err)
		return nil, err
	}
	return counts, nil
}

// IncrGlobalLimitCache .
func (d *Dao) IncrGlobalLimitCache(c context.Context, keywordID int64) (int64, error) {
	var (
		key  = globalCountsKey(keywordID)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	count, err := redis.Int64(conn.Do("INCR", key))
	if err != nil {
		log.Error("%v", err)
		return 0, err
	}
	return count, nil
}

// IncrLocalLimitCache .
func (d *Dao) IncrLocalLimitCache(c context.Context, keywordID, oid int64) (int64, error) {
	var (
		key  = localCountsKey(keywordID, oid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	count, err := redis.Int64(conn.Do("INCR", key))
	if err != nil {
		log.Error("%v", err)
		return 0, err
	}
	return count, nil
}

// LocalLimitExpire .
func (d *Dao) LocalLimitExpire(c context.Context, keywordID, oid, dur int64) error {
	var (
		key  = localCountsKey(keywordID, oid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if _, err := conn.Do("EXPIRE", key, dur); err != nil {
		log.Error("%v", err)
		return err
	}
	return nil
}

// GlobalLimitExpire .
func (d *Dao) GlobalLimitExpire(c context.Context, keywordID, dur int64) error {
	var (
		key  = globalCountsKey(keywordID)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if _, err := conn.Do("EXPIRE", key, dur); err != nil {
		log.Error("%v", err)
		return err
	}
	return nil
}

// DelRegexpCache .
func (d *Dao) DelRegexpCache(c context.Context) error {
	conn := d.redis.Get(c)
	defer conn.Close()
	if _, err := conn.Do("DEL", _regexpsKey); err != nil {
		log.Error("%v", err)
		return err
	}
	return nil
}

// DelRulesCache .
func (d *Dao) DelRulesCache(c context.Context, area, limitType string) error {
	var (
		key  = rulesKey(area, limitType)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if _, err := conn.Do("DEL", key); err != nil {
		log.Error("%v", err)
		return err
	}
	return nil
}

// AreaSendersExpire .
func (d *Dao) AreaSendersExpire(c context.Context, area string, senderID, dur int64) error {
	var (
		key  = areaSendersKey(area, senderID)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if _, err := conn.Do("EXPIRE", key, dur); err != nil {
		log.Error("%v", err)
		return err
	}
	return nil
}

// IncrAreaSendersCache .
func (d *Dao) IncrAreaSendersCache(c context.Context, area string, senderID int64) (int64, error) {
	var (
		key  = areaSendersKey(area, senderID)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	count, err := redis.Int64(conn.Do("INCR", key))
	if err != nil {
		log.Error("%v", err)
		return 0, err
	}
	return count, nil
}

// AllSendersCache .
func (d *Dao) AllSendersCache(c context.Context, keywordID int64) ([]string, error) {
	var (
		key  = sendersKey(keywordID)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	r, err := redis.Strings(conn.Do("ZRANGEBYSCORE", key, "-inf", "+inf"))
	if err != nil {
		log.Error("%v", err)
		return nil, err
	}
	return r, nil
}

// SendersCache .
func (d *Dao) SendersCache(c context.Context, keywordID, limit, offset int64) ([]string, error) {
	var (
		key  = sendersKey(keywordID)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	r, err := redis.Strings(conn.Do("ZRANGEBYSCORE", key, "-inf", "+inf", "LIMIT", limit, offset))
	if err != nil {
		log.Error("%v", err)
		return nil, err
	}
	return r, nil
}

// TotalLimitExpire .
func (d *Dao) TotalLimitExpire(c context.Context, keywordID, dur int64) error {
	var (
		key  = totalCountsKey(keywordID)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if _, err := conn.Do("EXPIRE", key, dur); err != nil {
		log.Error("%v", err)
		return err
	}
	return nil
}

// IncrTotalLimitCache .
func (d *Dao) IncrTotalLimitCache(c context.Context, keywordID int64) (int64, error) {
	var (
		key  = totalCountsKey(keywordID)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	count, err := redis.Int64(conn.Do("INCR", key))
	if err != nil {
		log.Error("%v", err)
		return 0, err
	}
	return count, nil
}

// ZaddSendersCache insert into sortedset and return total counts of sorted set
func (d *Dao) ZaddSendersCache(c context.Context, keywordID, score, senderID int64) (int64, error) {
	var (
		key  = sendersKey(keywordID)
		val  = fmt.Sprintf("%d", senderID)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	_, err := redis.Int64(conn.Do("ZADD", key, score, val))
	if err != nil {
		log.Error("%v", err)
		return 0, err
	}
	r, err := redis.Int64(conn.Do("ZCARD", key))
	if err != nil {
		log.Error("%v", err)
		return 0, err
	}
	return r, nil
}

// ZremSendersCache return the number of memebers removed from the sorted set
func (d *Dao) ZremSendersCache(c context.Context, keywordID int64, senderIDStr string) (int64, error) {
	var (
		key  = sendersKey(keywordID)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	r, err := redis.Int64(conn.Do("ZREM", key, senderIDStr))
	if err != nil {
		log.Error("%v", err)
		return 0, err
	}
	return r, nil
}

// DelKeywordRelatedCache .
func (d *Dao) DelKeywordRelatedCache(c context.Context, ks []*model.Keyword) error {
	var conn = d.redis.Get(c)
	defer conn.Close()
	for _, v := range ks {
		if err := conn.Send("DEL", totalCountsKey(v.ID)); err != nil {
			log.Error("%v", err)
			return err
		}
		if err := conn.Send("DEL", sendersKey(v.ID)); err != nil {
			log.Error("%v", err)
			return err
		}
	}
	if err := conn.Flush(); err != nil {
		log.Error("%v", err)
		return err
	}
	for i := 0; i < len(ks)*2; i++ {
		if _, err := conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return err
		}
	}
	return nil
}

// DelCountRelatedCache .
func (d *Dao) DelCountRelatedCache(c context.Context, k *model.Keyword) error {
	var conn = d.redis.Get(c)
	defer conn.Close()
	if err := conn.Send("DEL", globalCountsKey(k.ID)); err != nil {
		log.Error("%v", err)
		return err
	}
	if err := conn.Send("DEL", localCountsKey(k.ID, k.SenderID)); err != nil {
		log.Error("%v", err)
		return err
	}
	if err := conn.Send("DEL", sendersKey(k.ID)); err != nil {
		log.Error("%v", err)
		return err
	}
	if err := conn.Flush(); err != nil {
		log.Error("%v", err)
		return err
	}
	for i := 0; i < 3; i++ {
		if _, err := conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return err
		}
	}
	return nil
}
