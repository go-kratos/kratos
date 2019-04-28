package dao

import (
	"context"
	"fmt"
	"strconv"

	model "go-common/app/interface/main/credit/model"
	"go-common/library/cache/redis"
)

const (
	_voteOpIdx = "vo_%d_%d"
	_caseOpIdx = "caseop_"
)

func voteIndexKey(cid int64, otype int8) string {
	return fmt.Sprintf(_voteOpIdx, otype, cid)
}

func caseIndexKey(cid int64) string {
	return _caseOpIdx + strconv.FormatInt(cid, 10)
}

// VoteOpIdxCache get vote opinion index from cache.
func (d *Dao) VoteOpIdxCache(c context.Context, cid, start, end int64, otype int8) (ids []int64, err error) {
	var (
		key  = voteIndexKey(cid, otype)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	ids, err = redis.Int64s(conn.Do("LRANGE", key, start, end))
	return
}

// ExpireVoteIdx expire vote idx.
func (d *Dao) ExpireVoteIdx(c context.Context, cid int64, otype int8) (ok bool, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	ok, err = redis.Bool(conn.Do("EXPIRE", voteIndexKey(cid, otype), d.redisExpire))
	return
}

// LenVoteIdx get lenth of vote index.
func (d *Dao) LenVoteIdx(c context.Context, cid int64, otype int8) (count int, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	count, err = redis.Int(conn.Do("LLEN", voteIndexKey(cid, otype)))
	return
}

// CaseOpIdxCache get case opinion index from cache.
func (d *Dao) CaseOpIdxCache(c context.Context, cid, start, end int64) (ids []int64, err error) {
	var (
		key  = caseIndexKey(cid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	ids, err = redis.Int64s(conn.Do("ZREVRANGE", key, start, end))
	return
}

// LenCaseIdx get lenth of vote index.
func (d *Dao) LenCaseIdx(c context.Context, cid int64) (count int, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	count, err = redis.Int(conn.Do("ZCARD", caseIndexKey(cid)))
	return
}

// ExpireCaseIdx expire case index cache.
func (d *Dao) ExpireCaseIdx(c context.Context, cid int64) (ok bool, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	ok, err = redis.Bool(conn.Do("EXPIRE", caseIndexKey(cid), d.redisExpire))
	return
}

// LoadVoteOpIdxs load vote opinion index into cache.
func (d *Dao) LoadVoteOpIdxs(c context.Context, cid int64, otype int8, idx []int64) (err error) {
	var (
		ok   bool
		key  = voteIndexKey(cid, otype)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if ok, err = redis.Bool(conn.Do("EXPIRE", key, d.redisExpire)); ok {
		return
	}
	for _, id := range idx {
		if err = conn.Send("LPUSH", key, id); err != nil {
			return
		}
	}
	if err = conn.Send("EXPIRE", key, d.redisExpire); err != nil {
		return
	}
	if err = conn.Flush(); err != nil {
		return
	}
	for i := 0; i < len(idx)+1; i++ {
		_, err = conn.Receive()
	}
	return
}

// LoadCaseIdxs load case opinion index into redis.
func (d *Dao) LoadCaseIdxs(c context.Context, cid int64, ops []*model.Opinion) (err error) {
	key := caseIndexKey(cid)
	conn := d.redis.Get(c)
	defer conn.Close()
	for _, op := range ops {
		if err = conn.Send("ZADD", key, op.Like-op.Hate, op.OpID); err != nil {
			return
		}
	}
	if err = conn.Send("EXPIRE", key, d.redisExpire); err != nil {
		return
	}
	if err = conn.Flush(); err != nil {
		return
	}
	for i := 0; i < len(ops)+1; i++ {
		_, err = conn.Receive()
	}
	return
}

// DelCaseIdx DEL case opinion idx.
func (d *Dao) DelCaseIdx(c context.Context, cid int64) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	_, err = conn.Do("DEL", caseIndexKey(cid))
	return
}

// DelVoteIdx DEL case opinion idx.
func (d *Dao) DelVoteIdx(c context.Context, cid int64) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("DEL", voteIndexKey(cid, 1)); err != nil {
		return
	}
	if err = conn.Send("DEL", voteIndexKey(cid, 2)); err != nil {
		return
	}
	conn.Flush()
	for i := 0; i < 2; i++ {
		conn.Receive()
	}
	return
}
