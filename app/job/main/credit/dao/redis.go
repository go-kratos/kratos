package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"go-common/app/job/main/credit/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_voteOpIdx    = "vo_%d_%d"
	_caseOpIdx    = "caseop_"
	_blockIdx     = "bl_%d_%d"
	_grantCaseKey = "gr_ca_li_v2"
)

func voteIndexKey(cid int64, otype int8) string {
	return fmt.Sprintf(_voteOpIdx, otype, cid)
}

func caseIndexKey(cid int64) string {
	return _caseOpIdx + strconv.FormatInt(cid, 10)
}

func blockIndexKey(otype, btype int64) string {
	return fmt.Sprintf(_blockIdx, otype, btype)
}

// DelCaseIdx DEL case opinion idx.
func (d *Dao) DelCaseIdx(c context.Context, cid int64) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	if _, err = conn.Do("DEL", caseIndexKey(cid)); err != nil {
		log.Error("del case idx err(%v)", err)
		return
	}
	return
}

// DelBlockedInfoIdx ZREM block info idx.
func (d *Dao) DelBlockedInfoIdx(c context.Context, bl *model.BlockedInfo) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	conn.Send("ZREM", blockIndexKey(bl.OriginType, bl.BlockedType), bl.ID)
	conn.Send("ZREM", blockIndexKey(0, -1), bl.ID)
	conn.Send("ZREM", blockIndexKey(0, bl.BlockedType), bl.ID)
	conn.Send("ZREM", blockIndexKey(bl.OriginType, -1), bl.ID)
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush err(%v)", err)
		return
	}
	for i := 0; i < 4; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
		}
	}
	return
}

// AddBlockInfoIdx ZADD block info idx.
func (d *Dao) AddBlockInfoIdx(c context.Context, bl *model.BlockedInfo) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	var mtime time.Time
	if mtime, err = time.ParseInLocation("2006-01-02 15:04:05", bl.MTime, time.Local); err != nil {
		log.Error("time.ParseInLocation err(%v)", err)
		return
	}
	conn.Send("ZADD", blockIndexKey(bl.OriginType, bl.BlockedType), mtime.Unix(), bl.ID)
	conn.Send("ZADD", blockIndexKey(0, -1), mtime.Unix(), bl.ID)
	conn.Send("ZADD", blockIndexKey(0, bl.BlockedType), mtime.Unix(), bl.ID)
	conn.Send("ZADD", blockIndexKey(bl.OriginType, -1), mtime.Unix(), bl.ID)
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush err(%v)", err)
		return
	}
	for i := 0; i < 4; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
		}
	}
	return
}

// DelVoteIdx DEL case opinion idx.
func (d *Dao) DelVoteIdx(c context.Context, cid int64) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("DEL", voteIndexKey(cid, 1)); err != nil {
		log.Error("del case idx err(%v)", err)
		return
	}
	if err = conn.Send("DEL", voteIndexKey(cid, 2)); err != nil {
		log.Error("del case idx err(%v)", err)
		return
	}
	conn.Flush()
	for i := 0; i < 2; i++ {
		conn.Receive()
	}
	return
}

// SetGrantCase set grant case ids.
func (d *Dao) SetGrantCase(c context.Context, mcases map[int64]*model.SimCase) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	args := redis.Args{}.Add(_grantCaseKey)
	for cid, mcase := range mcases {
		var bs []byte
		bs, err = json.Marshal(mcase)
		if err != nil {
			log.Error("json.Marshal(%+v) error(%v)", mcase, err)
			err = nil
			continue
		}
		args = args.Add(cid).Add(string(bs))
	}
	if _, err = conn.Do("HMSET", args...); err != nil {
		log.Error("conn.Send(HMSET,%v) error(%v)", args, err)
	}
	return
}

// DelGrantCase del grant case id.
func (d *Dao) DelGrantCase(c context.Context, cids []int64) (err error) {
	var args = []interface{}{_grantCaseKey}
	conn := d.redis.Get(c)
	defer conn.Close()
	for _, cid := range cids {
		args = append(args, cid)
	}
	if _, err = conn.Do("HDEL", args...); err != nil {
		log.Error("conn.Send(HDEL,%s) err(%v)", _grantCaseKey, err)
	}
	return
}

// TotalGrantCase get length of grant case ids.
func (d *Dao) TotalGrantCase(c context.Context) (count int, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	if count, err = redis.Int(conn.Do("HLEN", _grantCaseKey)); err != nil {
		if err != redis.ErrNil {
			log.Error("conn.Do(HLEN, %s) error(%v)", _grantCaseKey, err)
			return
		}
		err = nil
	}
	return
}

// GrantCases get granting case cids.
func (d *Dao) GrantCases(c context.Context) (cids []int64, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	var ms map[string]string
	if ms, err = redis.StringMap(conn.Do("HGETALL", _grantCaseKey)); err != nil {
		if err == redis.ErrNil {
			err = nil
		}
		return
	}
	for m, s := range ms {
		if s == "" {
			continue
		}
		cid, err := strconv.ParseInt(m, 10, 64)
		if err != nil {
			log.Error("strconv.ParseInt(%s) error(%v)", m, err)
			err = nil
			continue
		}
		cids = append(cids, cid)
	}
	return
}
