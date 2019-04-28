package dao

import (
	"context"
	"fmt"
	"strconv"

	"go-common/app/admin/main/reply/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_prefixIdx        = "i_"
	_prefixNewRootIdx = "ri_"
	_prefixAuditIdx   = "ai_%d_%d"
	// 针对大忽悠时间被删除评论的人的MID
	_prefixAdminDelMid = "mid_%d"

	// f_{折叠类型，根评论还是评论区}_{评论区ID或者根评论ID}
	_foldedReplyFmt = "f_%s_%d"

	// dialog
	_prefixDialogIdx = "d_%d"

	_maxCount = 20000
)

func keyFolderIdx(kind string, ID int64) string {
	return fmt.Sprintf(_foldedReplyFmt, kind, ID)
}

func keyDelMid(mid int64) string {
	return fmt.Sprintf(_prefixAdminDelMid, mid)
}

// KeyMainIdx ...
func keyMainIdx(oid int64, tp, sort int32) string {
	if oid > _oidOverflow {
		return fmt.Sprintf("%s_%d_%d_%d", _prefixIdx, oid, tp, sort)
	}
	return _prefixIdx + strconv.FormatInt((oid<<16)|(int64(tp)<<8)|int64(sort), 10)
}

// keyDialogIdx ...
func keyDialogIdx(dialog int64) string {
	return fmt.Sprintf(_prefixDialogIdx, dialog)
}

// KeyRootIdx ...
func keyRootIdx(root int64) string {
	return _prefixNewRootIdx + strconv.FormatInt(root, 10)
}

func keyIdx(oid int64, tp, sort int32) string {
	if oid > _oidOverflow {
		return fmt.Sprintf("%s_%d_%d_%d", _prefixIdx, oid, tp, sort)
	}
	return _prefixIdx + strconv.FormatInt((oid<<16)|(int64(tp)<<8)|int64(sort), 10)
}

func keyNewRootIdx(rpID int64) string {
	return _prefixNewRootIdx + strconv.FormatInt(rpID, 10)
}

func keyAuditIdx(oid int64, tp int32) string {
	return fmt.Sprintf(_prefixAuditIdx, oid, tp)
}

func (d *Dao) ExsistsDelMid(c context.Context, mid int64) (ok bool, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := keyDelMid(mid)
	if ok, err = redis.Bool(conn.Do("EXISTS", key)); err != nil {
		log.Error("conn.Do(EXISTS, %s) error(%v)", key, err)
	}
	return
}

func (d *Dao) SetDelMid(c context.Context, mid int64) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	// 15天内
	if err = conn.Send("SETEX", keyDelMid(mid), 86400*15, 1); err != nil {
		log.Error("redis SETEX(%s) error(%v)", keyDelMid(mid), err)
	}
	return
}

// ExpireIndex set expire time for index.
func (d *Dao) ExpireIndex(c context.Context, oid int64, typ, sort int32) (ok bool, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	if ok, err = redis.Bool(conn.Do("EXPIRE", keyIdx(oid, typ, sort), d.redisExpire)); err != nil {
		log.Error("conn.Do(EXPIRE) error(%v)", err)
	}
	return
}

// ExpireNewChildIndex set expire time for root's index.
func (d *Dao) ExpireNewChildIndex(c context.Context, root int64) (ok bool, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	if ok, err = redis.Bool(conn.Do("EXPIRE", keyNewRootIdx(root), d.redisExpire)); err != nil {
		log.Error("conn.Do(EXPIRE) error(%v)", err)
	}
	return
}

// TopChildReply ...
func (d *Dao) TopChildReply(c context.Context, root, child int64) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	if _, err = conn.Do("ZADD", keyNewRootIdx(root), 0, child); err != nil {
		return
	}
	return
}

// CountReplies get count of reply.
func (d *Dao) CountReplies(c context.Context, oid int64, tp, sort int32) (count int, err error) {
	key := keyIdx(oid, tp, sort)
	conn := d.redis.Get(c)
	defer conn.Close()
	if count, err = redis.Int(conn.Do("ZCARD", key)); err != nil {
		log.Error("CountReplies error(%v)", err)
	}
	return
}

// MinScore get the lowest score from sorted set
func (d *Dao) MinScore(c context.Context, oid int64, tp int32, sort int32) (score int32, err error) {
	key := keyIdx(oid, tp, sort)
	conn := d.redis.Get(c)
	defer conn.Close()
	values, err := redis.Values(conn.Do("ZRANGE", key, 0, 0, "WITHSCORES"))
	if err != nil {
		log.Error("conn.Do(ZREVRANGE, %s) error(%v)", key, err)
		return
	}
	if len(values) != 2 {
		err = fmt.Errorf("redis zrange items(%v) length not 2", values)
		return
	}
	var id int64
	redis.Scan(values, &id, &score)
	return
}

// AddFloorIndex add index by floor.
func (d *Dao) AddFloorIndex(c context.Context, rp *model.Reply) (err error) {
	min, err := d.MinScore(c, rp.Oid, rp.Type, model.SortByFloor)
	if err != nil {
		log.Error("s.dao.Redis.AddFloorIndex failed , oid(%d) type(%d) err(%v)", rp.Oid, rp.Type, err)
	} else if rp.Floor <= min {
		return
	}

	key := keyIdx(rp.Oid, rp.Type, model.SortByFloor)
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("ZADD", key, rp.Floor, rp.ID); err != nil {
		log.Error("conn.Send error(%v)", err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.redisExpire); err != nil {
		log.Error("conn.Send error(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	if _, err = conn.Receive(); err != nil {
		log.Error("conn.Receive error(%v)", err)
		return
	}
	return
}

// AddCountIndex add index by count.
func (d *Dao) AddCountIndex(c context.Context, rp *model.Reply) (err error) {
	if rp.IsTop() {
		return
	}
	var count int
	if count, err = d.CountReplies(c, rp.Oid, rp.Type, model.SortByCount); err != nil {
		return
	} else if count >= _maxCount {
		var min int32
		if min, err = d.MinScore(c, rp.Oid, rp.Type, model.SortByCount); err != nil {
			return
		}
		if rp.RCount <= min {
			return
		}
	}

	key := keyIdx(rp.Oid, rp.Type, model.SortByCount)
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("ZADD", key, int64(rp.RCount)<<32|(int64(rp.Floor)&0xFFFFFFFF), rp.ID); err != nil {
		log.Error("conn.Send error(%v)", err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.redisExpire); err != nil {
		log.Error("conn.Send error(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	if _, err = conn.Receive(); err != nil {
		log.Error("conn.Receive error(%v)", err)
		return
	}
	return
}

// AddLikeIndex add index by like.
func (d *Dao) AddLikeIndex(c context.Context, rp *model.Reply, rpt *model.Report) (err error) {
	if rp.IsTop() {
		return
	}
	var rptCn int32
	if rpt != nil {
		rptCn = rpt.Count
	}
	score := int64((float32(rp.Like+model.WeightLike) / float32(rp.Hate+model.WeightHate+rptCn)) * 100)
	score = score<<32 | (int64(rp.RCount) & 0xFFFFFFFF)
	var count int
	if count, err = d.CountReplies(c, rp.Oid, rp.Type, model.SortByLike); err != nil {
		return
	} else if count >= _maxCount {
		var min int32
		if min, err = d.MinScore(c, rp.Oid, rp.Type, model.SortByLike); err != nil {
			return
		}
		if score <= int64(min) {
			return
		}
	}

	key := keyIdx(rp.Oid, rp.Type, model.SortByLike)
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("ZADD", key, score, rp.ID); err != nil {
		log.Error("conn.Send error(%v)", err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.redisExpire); err != nil {
		log.Error("conn.Send error(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	if _, err = conn.Receive(); err != nil {
		log.Error("conn.Receive() error(%v)", err)
		return
	}
	return
}

// DelIndexBySort delete index by sort.
func (d *Dao) DelIndexBySort(c context.Context, rp *model.Reply, sort int32) (err error) {
	key := keyIdx(rp.Oid, rp.Type, sort)
	conn := d.redis.Get(c)
	defer conn.Close()
	if _, err = conn.Do("ZREM", key, rp.ID); err != nil {
		log.Error("conn.Do(ZREM) error(%v)", err)
	}
	return
}

// DelReplyIndex delete reply index.
func (d *Dao) DelReplyIndex(c context.Context, rp *model.Reply) (err error) {
	var (
		key string
		n   int
	)
	conn := d.redis.Get(c)
	defer conn.Close()
	if rp.Root == 0 {
		key = keyIdx(rp.Oid, rp.Type, model.SortByFloor)
		err = conn.Send("ZREM", key, rp.ID)
		key = keyIdx(rp.Oid, rp.Type, model.SortByCount)
		err = conn.Send("ZREM", key, rp.ID)
		key = keyIdx(rp.Oid, rp.Type, model.SortByLike)
		err = conn.Send("ZREM", key, rp.ID)
		n += 3
	} else {
		if rp.Dialog != 0 {
			key = keyDialogIdx(rp.Dialog)
			err = conn.Send("ZREM", key, rp.ID)
			n++
		}
		key = keyNewRootIdx(rp.Root)
		err = conn.Send("ZREM", key, rp.ID)
		n++
	}
	if err != nil {
		log.Error("conn.Send error(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < n; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive error(%v)", err)
			return
		}
	}
	return
}

// AddNewChildIndex add root reply index by floor.
func (d *Dao) AddNewChildIndex(c context.Context, rp *model.Reply) (err error) {
	key := keyNewRootIdx(rp.Root)
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("ZADD", key, rp.Floor, rp.ID); err != nil {
		log.Error("conn.Send error(%v)", err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.redisExpire); err != nil {
		log.Error("conn.Send error(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	if _, err = conn.Receive(); err != nil {
		log.Error("conn.Receive() error(%v)", err)
		return
	}
	return
}

// DelAuditIndex delete audit reply cache.
func (d *Dao) DelAuditIndex(c context.Context, rp *model.Reply) (err error) {
	key := keyAuditIdx(rp.Oid, rp.Type)
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("ZREM", key, rp.ID); err != nil {
		log.Error("conn.Send(ZREM %s) error(%v)", key, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	if _, err = conn.Receive(); err != nil {
		log.Error("conn.Receive() error(%v)", err)
		return
	}
	return
}

// RemReplyFromRedis ...
func (d *Dao) RemReplyFromRedis(c context.Context, keyMapping map[string][]int64) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	count := 0
	for key, rpIDs := range keyMapping {
		if len(rpIDs) == 0 {
			continue
		}
		args := make([]interface{}, 0, 1+len(rpIDs))
		args = append(args, key)
		for _, rpID := range rpIDs {
			args = append(args, rpID)
		}
		if err = conn.Send("ZREM", args...); err != nil {
			log.Error("conn.Send(ZREM %s) error(%v)", key, err)
			return
		}
		count++
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < count; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// ExpireFolder ...
func (d *Dao) ExpireFolder(c context.Context, kind string, ID int64) (ok bool, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	if ok, err = redis.Bool(conn.Do("EXPIRE", keyFolderIdx(kind, ID), d.redisExpire)); err != nil {
		log.Error("conn.Do(EXPIRE) error(%v)", err)
	}
	return
}

// AddFolder ...
func (d *Dao) AddFolder(c context.Context, keyMapping map[string][]*model.Reply) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	count := 0
	for key, rps := range keyMapping {
		var args []interface{}
		args = append(args, key)
		for _, rp := range rps {
			args = append(args, rp.Floor)
			args = append(args, rp.ID)
		}
		if err = conn.Send("ZADD", args...); err != nil {
			log.Error("conn.Send(ZADD %s) error(%v)", key, err)
			return
		}
		count++
		if err = conn.Send("EXPIRE", key, d.redisExpire); err != nil {
			return
		}
		count++
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < count; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// RemFolder ...
func (d *Dao) RemFolder(c context.Context, kind string, ID, rpID int64) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := keyFolderIdx(kind, ID)
	if _, err = conn.Do("ZREM", key, rpID); err != nil {
		log.Error("conn.Do(ZREM) error(%v)", err)
	}
	return
}
