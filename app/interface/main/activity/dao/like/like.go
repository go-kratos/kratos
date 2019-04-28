package like

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"go-common/app/interface/main/activity/model/like"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/database/elastic"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	xtime "go-common/library/time"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_selLikeSQL        = "SELECT id,wid FROM likes where state = 1 AND sid = ? ORDER BY type"
	_likeSQL           = "SELECT id,sid,type,mid,wid,state,stick_top,ctime,mtime FROM likes WHERE id = ? and state = 1"
	_likeMoreLidSQL    = "SELECT id,sid,type,mid,wid,state,stick_top,ctime,mtime FROM likes WHERE id > ?  order by id asc limit 1000"
	_likesBySidSQL     = "SELECT id,sid,type,mid,wid,state,stick_top,ctime,mtime FROM likes WHERE id > ? and sid = ? and state = 1 order by id asc limit 1000"
	_likesSQL          = "SELECT id,sid,type,mid,wid,state,stick_top,ctime,mtime FROM likes WHERE  id IN (%s) and state = 1"
	_likeListSQL       = "SELECT id,wid,ctime FROM likes WHERE state = 1 AND sid = ? ORDER BY id DESC"
	_likeMaxIDSQL      = "SELECT id FROM likes ORDER BY id DESC limit 1"
	_keyLikeTagFmt     = "l_t_%d_%d"
	_keyLikeTagCntsFmt = "l_t_cs_%d"
	_keyLikeRegionFmt  = "l_r_%d_%d"
	// likeAPI ip frequence key the old is ddos:like:ip:%s
	_keyIPRequestFmt = "go:ddos:l:ip:%s"
	// the cache set of like order by ctime the old is bilibili-activity:ctime:%d
	_keyLikeListCtimeFmt  = "go:bl-a:ctime:%d"
	_keyLikeListRandomFmt = "go:bl-a:random:%d"
	// the cache set of like type order by ctime
	_keyLikeListTypeCtimeFmt = "go:b:a:t:%d:%d"
	// storyKing LikeAct cache
	_keyStoryDilyLikeFmt = "go:s:d:m:%s:%d:%d"
	// storyKing each likeAct cahce
	_keyStoryEachLikeFmt = "go:s:ea:m:%s:%d:%d:%d"
	// es index
	_activity = "activity"
	// EsOrderLikes archive center likes.
	EsOrderLikes = "likes"
	// EsOrderCoin archive center coin .
	EsOrderCoin = "coin"
	// EsOrderReply archive center reply.
	EsOrderReply = "reply"
	// EsOrderShare  archive center share.
	EsOrderShare = "share"
	// EsOrderClick archive center click
	EsOrderClick = "click"
	// EsOrderDm archive center  dm
	EsOrderDm = "dm"
	// EsOrderFav archive center fav
	EsOrderFav = "fav"
	// ActOrderLike activity list like order.
	ActOrderLike = "like"
	// ActOrderCtime activity list ctime order.
	ActOrderCtime = "ctime"
	// ActOrderRandom order random .
	ActOrderRandom = "random"
)

// ipRequestKey .
func ipRequestKey(ip string) string {
	return fmt.Sprintf(_keyIPRequestFmt, ip)
}

func likeListCtimeKey(sid int64) string {
	return fmt.Sprintf(_keyLikeListCtimeFmt, sid)
}

func likeListRandomKey(sid int64) string {
	return fmt.Sprintf(_keyLikeListRandomFmt, sid)
}

func likeListTypeCtimeKey(types int, sid int64) string {
	return fmt.Sprintf(_keyLikeListTypeCtimeFmt, types, sid)
}

func keyLikeTag(sid, tagID int64) string {
	return fmt.Sprintf(_keyLikeTagFmt, sid, tagID)
}

func keyLikeTagCounts(sid int64) string {
	return fmt.Sprintf(_keyLikeTagCntsFmt, sid)
}

func keyLikeRegion(sid int64, regionID int16) string {
	return fmt.Sprintf(_keyLikeRegionFmt, sid, regionID)
}

func keyStoryLikeKey(sid, mid int64, daily string) string {
	return fmt.Sprintf(_keyStoryDilyLikeFmt, daily, sid, mid)
}

func keyStoryEachLike(sid, mid, lid int64, daily string) string {
	return fmt.Sprintf(_keyStoryEachLikeFmt, daily, sid, mid, lid)
}

// LikeTypeList dao sql.
func (dao *Dao) LikeTypeList(c context.Context, sid int64) (ns []*like.Like, err error) {
	rows, err := dao.db.Query(c, _selLikeSQL, sid)
	if err != nil {
		log.Error("LikeTypeList dao.db.Query error(%v)", err)
		return
	}
	ns = make([]*like.Like, 0)
	defer rows.Close()
	for rows.Next() {
		n := &like.Like{}
		if err = rows.Scan(&n.ID, &n.Wid); err != nil {
			log.Error("row.Scan error(%v)", err)
			return
		}
		ns = append(ns, n)
	}
	if err = rows.Err(); err != nil {
		log.Error("row.Scan row error(%v)", err)
	}
	return
}

// LikeList dao sql
func (dao *Dao) LikeList(c context.Context, sid int64) (ns []*like.Item, err error) {
	rows, err := dao.db.Query(c, _likeListSQL, sid)
	if err != nil {
		log.Error("LikeList dao.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		n := new(like.Item)
		if err = rows.Scan(&n.ID, &n.Wid, &n.Ctime); err != nil {
			log.Error("row.Scan error(%v)", err)
			return
		}
		ns = append(ns, n)
	}
	if err = rows.Err(); err != nil {
		log.Error("row.Scan row error(%v)", err)
	}
	return
}

// RawLikes get likes by wid.
func (dao *Dao) RawLikes(c context.Context, ids []int64) (data map[int64]*like.Item, err error) {
	rows, err := dao.db.Query(c, fmt.Sprintf(_likesSQL, xstr.JoinInts(ids)))
	if err != nil {
		log.Error("Likes dao.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	data = make(map[int64]*like.Item)
	for rows.Next() {
		res := &like.Item{}
		if err = rows.Scan(&res.ID, &res.Sid, &res.Type, &res.Mid, &res.Wid, &res.State, &res.StickTop, &res.Ctime, &res.Mtime); err != nil {
			log.Error("Likes row.Scan error(%v)", err)
			return
		}
		data[res.ID] = res
	}
	if err = rows.Err(); err != nil {
		log.Error("Likes row.Scan row error(%v)", err)
	}

	return
}

// LikeTagCache get like tag cache.
func (dao *Dao) LikeTagCache(c context.Context, sid, tagID int64, start, end int) (likes []*like.Item, err error) {
	var values []interface{}
	key := keyLikeTag(sid, tagID)
	conn := dao.redis.Get(c)
	defer conn.Close()
	if values, err = redis.Values(conn.Do("ZREVRANGE", key, start, end)); err != nil {
		log.Error("LikeTagCache conn.Do(ZREVRANGE, %s) error(%v)", key, err)
		return
	} else if len(values) == 0 {
		return
	}
	for len(values) > 0 {
		var bs []byte
		if values, err = redis.Scan(values, &bs); err != nil {
			log.Error("LikeRegionCache redis.Scan(%v) error(%v)", values, err)
			return
		}
		like := new(like.Item)
		if err = json.Unmarshal(bs, &like); err != nil {
			log.Error("LikeRegionCache conn.Do(ZRANGE, %s) error(%v)", key, err)
			continue
		}
		if like.ID > 0 {
			likes = append(likes, like)
		}
	}
	return
}

// LikeTagCnt get like tag cnt.
func (dao *Dao) LikeTagCnt(c context.Context, sid, tagID int64) (count int, err error) {
	key := keyLikeTag(sid, tagID)
	conn := dao.redis.Get(c)
	defer conn.Close()
	if count, err = redis.Int(conn.Do("ZCARD", key)); err != nil {
		log.Error("LikeRegionCnt conn.Do(ZCARD, %s) error(%v)", key, err)
	}
	return
}

// SetLikeTagCache set like tag cache no expire.
func (dao *Dao) SetLikeTagCache(c context.Context, sid, tagID int64, likes []*like.Item) (err error) {
	var bs []byte
	key := keyLikeTag(sid, tagID)
	conn := dao.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("DEL", key); err != nil {
		log.Error("SetLikeTagCache conn.Send(DEL, %s) error(%v)", key, err)
		return
	}
	args := redis.Args{}.Add(key)
	for _, v := range likes {
		if bs, err = json.Marshal(v); err != nil {
			log.Error("SetLikeTagCache json.Marshal() error(%v)", err)
			return
		}
		args = args.Add(v.Ctime).Add(bs)
	}
	if err = conn.Send("ZADD", args...); err != nil {
		log.Error("conn.Send(ZADD, %s) error(%v)", key, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
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

// LikeRegionCache get like region cache.
func (dao *Dao) LikeRegionCache(c context.Context, sid int64, regionID int16, start, end int) (likes []*like.Item, err error) {
	var values []interface{}
	key := keyLikeRegion(sid, regionID)
	conn := dao.redis.Get(c)
	defer conn.Close()
	if values, err = redis.Values(conn.Do("ZREVRANGE", key, start, end)); err != nil {
		log.Error("LikeRegionCache conn.Do(ZREVRANGE, %s) error(%v)", key, err)
		return
	} else if len(values) == 0 {
		return
	}
	for len(values) > 0 {
		var bs []byte
		if values, err = redis.Scan(values, &bs); err != nil {
			log.Error("LikeRegionCache redis.Scan(%v) error(%v)", values, err)
			return
		}
		like := new(like.Item)
		if err = json.Unmarshal(bs, &like); err != nil {
			log.Error("LikeRegionCache conn.Do(ZREVRANGE, %s) error(%v)", key, err)
			continue
		}
		if like.ID > 0 {
			likes = append(likes, like)
		}
	}
	return
}

// LikeRegionCnt get like region cnt.
func (dao *Dao) LikeRegionCnt(c context.Context, sid int64, regionID int16) (count int, err error) {
	key := keyLikeRegion(sid, regionID)
	conn := dao.redis.Get(c)
	defer conn.Close()
	if count, err = redis.Int(conn.Do("ZCARD", key)); err != nil {
		log.Error("LikeRegionCnt conn.Do(ZCARD, %s) error(%v)", key, err)
	}
	return
}

// SetLikeRegionCache set like region cache.
func (dao *Dao) SetLikeRegionCache(c context.Context, sid int64, regionID int16, likes []*like.Item) (err error) {
	var bs []byte
	key := keyLikeRegion(sid, regionID)
	conn := dao.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("DEL", key); err != nil {
		log.Error("SetLikeTagCache conn.Send(DEL, %s) error(%v)", key, err)
		return
	}
	args := redis.Args{}.Add(key)
	for _, v := range likes {
		if bs, err = json.Marshal(v); err != nil {
			log.Error("SetLikeRegionCache json.Marshal() error(%v)", err)
			return
		}
		args = args.Add(v.Ctime).Add(bs)
	}
	if err = conn.Send("ZADD", args...); err != nil {
		log.Error("conn.Send(ZADD, %s) error(%v)", key, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
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

// SetTagLikeCountsCache .
func (dao *Dao) SetTagLikeCountsCache(c context.Context, sid int64, counts map[int64]int32) (err error) {
	key := keyLikeTagCounts(sid)
	conn := dao.redis.Get(c)
	defer conn.Close()
	args := redis.Args{}.Add(key)
	for tagID, count := range counts {
		args = args.Add(tagID).Add(count)
	}
	if _, err = conn.Do("HMSET", args...); err != nil {
		log.Error("SetLikeCountsCache conn.Do(HMSET) key(%s) error(%v)", key, err)
	}
	return
}

// TagLikeCountsCache get tag like counts cache.
func (dao *Dao) TagLikeCountsCache(c context.Context, sid int64, tagIDs []int64) (counts map[int64]int32, err error) {
	if len(tagIDs) == 0 {
		return
	}
	key := keyLikeTagCounts(sid)
	conn := dao.redis.Get(c)
	defer conn.Close()
	args := redis.Args{}.Add(key).AddFlat(tagIDs)
	var tmpCounts []int
	if tmpCounts, err = redis.Ints(conn.Do("HMGET", args...)); err != nil {
		log.Error("redis.Ints(HMGET) key(%s) args(%v) error(%v)", key, args, err)
		return
	}
	if len(tmpCounts) != len(tagIDs) {
		return
	}
	counts = make(map[int64]int32, len(tagIDs))
	for i, tagID := range tagIDs {
		counts[tagID] = int32(tmpCounts[i])
	}
	return
}

// RawLike get like by id .
func (dao *Dao) RawLike(c context.Context, id int64) (res *like.Item, err error) {
	res = new(like.Item)
	row := dao.db.QueryRow(c, _likeSQL, id)
	if err = row.Scan(&res.ID, &res.Sid, &res.Type, &res.Mid, &res.Wid, &res.State, &res.StickTop, &res.Ctime, &res.Mtime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			err = errors.Wrap(err, "LikeByID:QueryRow")
		}
	}
	return
}

// LikeListMoreLid get likes data with like.id greater than lid
func (dao *Dao) LikeListMoreLid(c context.Context, lid int64) (res []*like.Item, err error) {
	var rows *sql.Rows
	if rows, err = dao.db.Query(c, _likeMoreLidSQL, lid); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			err = errors.Wrap(err, "LikeListMoreLid:dao.db.Query()")
		}
		return
	}
	defer rows.Close()
	res = make([]*like.Item, 0, 1000)
	for rows.Next() {
		a := &like.Item{}
		if err = rows.Scan(&a.ID, &a.Sid, &a.Type, &a.Mid, &a.Wid, &a.State, &a.StickTop, &a.Ctime, &a.Mtime); err != nil {
			err = errors.Wrap(err, "LikeListMoreLid:rows.Scan()")
			return
		}
		res = append(res, a)
	}
	if err = rows.Err(); err != nil {
		err = errors.Wrap(err, "LikeListMoreLid: rows.Err()")
	}
	return
}

// LikesBySid get sid all likes .
func (dao *Dao) LikesBySid(c context.Context, lid, sid int64) (res []*like.Item, err error) {
	var rows *sql.Rows
	if rows, err = dao.db.Query(c, _likesBySidSQL, lid, sid); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			err = errors.Wrap(err, "LikesBySid:dao.db.Query()")
		}
		return
	}
	defer rows.Close()
	res = make([]*like.Item, 0, 1000)
	for rows.Next() {
		a := &like.Item{}
		if err = rows.Scan(&a.ID, &a.Sid, &a.Type, &a.Mid, &a.Wid, &a.State, &a.StickTop, &a.Ctime, &a.Mtime); err != nil {
			err = errors.Wrap(err, "LikesBySid:rows.Scan()")
			return
		}
		res = append(res, a)
	}
	if err = rows.Err(); err != nil {
		err = errors.Wrap(err, "LikesBySid:rows.Err()")
	}
	return
}

// IPReqquestCheck check ip has ben used or not .
func (dao *Dao) IPReqquestCheck(c context.Context, ip string) (val int, err error) {
	var (
		mcKey = ipRequestKey(ip)
		conn  = dao.mc.Get(c)
		item  *memcache.Item
	)
	defer conn.Close()
	if item, err = conn.Get(mcKey); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			val = 0
		} else {
			err = errors.Wrap(err, "IPReqquestCheck:conn.Get() error")
		}
		return
	}
	if err = conn.Scan(item, &val); err != nil {
		err = errors.Wrap(err, "IPReqquestCheck:conn.Scan() ")
	}
	return
}

// SetIPRequest set ip has been used
func (dao *Dao) SetIPRequest(c context.Context, ip string) (err error) {
	var (
		conn = dao.mc.Get(c)
		item = &memcache.Item{
			Key:        ipRequestKey(ip),
			Expiration: dao.mcLikeIPExpire,
			Flags:      memcache.FlagRAW,
			Value:      []byte("1"),
		}
	)
	defer conn.Close()
	if err = conn.Set(item); err != nil {
		err = errors.Wrap(err, "SetIPRequest:conn.Set()")
	}
	return
}

// LikeCtime .
func (dao *Dao) LikeCtime(c context.Context, sid int64, start, end int) (res []int64, err error) {
	var (
		conn = dao.redis.Get(c)
		key  = likeListCtimeKey(sid)
	)
	defer conn.Close()
	if res, err = redis.Int64s(conn.Do("ZREVRANGE", key, start, end)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			err = errors.Wrap(err, "conn.Do(ZREVRANGE)")
		}
	}
	return
}

// LikeRandom .
func (dao *Dao) LikeRandom(c context.Context, sid int64, start, end int) (res []int64, err error) {
	var (
		conn = dao.redis.Get(c)
		key  = likeListRandomKey(sid)
	)
	defer conn.Close()
	if res, err = redis.Int64s(conn.Do("ZREVRANGE", key, start, end)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			err = errors.Wrap(err, "conn.Do(ZREVRANGE)")
		}
	}
	return
}

// LikeRandomCount .
func (dao *Dao) LikeRandomCount(c context.Context, sid int64) (res int64, err error) {
	var (
		conn = dao.redis.Get(c)
		key  = likeListRandomKey(sid)
	)
	defer conn.Close()
	if res, err = redis.Int64(conn.Do("ZCARD", key)); err != nil {
		log.Error("LikeRandomCount conn.Do(ZCARD, %s) error(%v)", key, err)
	}
	return
}

// SetLikeRandom .
func (dao *Dao) SetLikeRandom(c context.Context, sid int64, ids []int64) (err error) {
	var (
		conn = dao.redis.Get(c)
		key  = likeListRandomKey(sid)
	)
	defer conn.Close()
	if len(ids) == 0 {
		return
	}
	args := redis.Args{}.Add(key)
	for k, v := range ids {
		args = args.Add(k + 1).Add(v)
	}
	if err = conn.Send("ZADD", args...); err != nil {
		err = errors.Wrap(err, "conn.Send(ZADD)")
		return
	}
	if err = conn.Send("EXPIRE", key, dao.randomExpire); err != nil {
		err = errors.Wrap(err, "conn.Send(EXPIRE)")
		return
	}
	if err = conn.Flush(); err != nil {
		err = errors.Wrap(err, "conn.Flush()")
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			err = errors.Wrapf(err, "conn.Receive(%d)", i)
			return
		}
	}
	return
}

// LikeCount .
func (dao *Dao) LikeCount(c context.Context, sid int64) (res int64, err error) {
	var (
		conn = dao.redis.Get(c)
		key  = likeListCtimeKey(sid)
	)
	defer conn.Close()
	if res, err = redis.Int64(conn.Do("ZCARD", key)); err != nil {
		log.Error("LikeCount conn.Do(ZCARD, %s) error(%v)", key, err)
	}
	return
}

// LikeListCtime set like list by ctime.
func (dao *Dao) LikeListCtime(c context.Context, sid int64, items []*like.Item) (err error) {
	var (
		conn = dao.redis.Get(c)
		key  = likeListCtimeKey(sid)
		max  = 0
	)
	defer conn.Close()
	args := redis.Args{}.Add(key)
	for _, v := range items {
		args = args.Add(v.Ctime).Add(v.ID)
		if v.Type != 0 {
			typeKey := likeListTypeCtimeKey(v.Type, sid)
			typeArgs := redis.Args{}.Add(typeKey).Add(v.Ctime).Add(v.ID)
			if err = conn.Send("ZADD", typeArgs...); err != nil {
				log.Error("LikeListCtime:conn.Send(%v) error(%v)", v, err)
				return
			}
			max++
		}
	}
	if err = conn.Send("ZADD", args...); err != nil {
		log.Error("LikeListCtime:conn.Send(%v) error(%v)", items, err)
		return
	}
	max++
	if err = conn.Flush(); err != nil {
		log.Error("LikeListCtime:conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < max; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("LikeListCtime:conn.Receive(%d) error(%v)", i, err)
			return
		}
	}
	return
}

//DelLikeListCtime delete likeList Ctime cache .
func (dao *Dao) DelLikeListCtime(c context.Context, sid int64, items []*like.Item) (err error) {
	var (
		conn = dao.redis.Get(c)
		key  = likeListCtimeKey(sid)
		max  = 0
	)
	defer conn.Close()
	args := redis.Args{}.Add(key)
	for _, v := range items {
		args = args.Add(v.ID)
		if v.Type != 0 {
			typeKey := likeListTypeCtimeKey(v.Type, sid)
			if err = conn.Send("ZREM", typeKey, v.ID); err != nil {
				log.Error("DelLikeListCtime:conn.Send(%v) error(%v)", v, err)
				return
			}
			max++
		}
	}
	if err = conn.Send("ZREM", args...); err != nil {
		log.Error("DelLikeListCtime:conn.Send(%v) error(%v)", args, err)
		return
	}
	max++
	if err = conn.Flush(); err != nil {
		log.Error("DelLikeListCtime:conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < max; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("DelLikeListCtime:conn.Receive(%d) error(%v)", i, err)
			return
		}
	}
	return
}

// LikeMaxID get likes last id .
func (dao *Dao) LikeMaxID(c context.Context) (res *like.Item, err error) {
	res = new(like.Item)
	rows := dao.db.QueryRow(c, _likeMaxIDSQL)
	if err = rows.Scan(&res.ID); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			err = errors.Wrap(err, "LikeMaxID:QueryRow")
		}
	}
	return
}

// GroupItemData like data.
func (dao *Dao) GroupItemData(c context.Context, sid int64, ck string) (data []*like.GroupItem, err error) {
	var req *http.Request
	if req, err = dao.client.NewRequest(http.MethodGet, fmt.Sprintf(dao.likeItemURL, sid), metadata.String(c, metadata.RemoteIP), url.Values{}); err != nil {
		return
	}
	req.Header.Set("Cookie", ck)
	var res struct {
		Code int `json:"code"`
		Data struct {
			List []*like.GroupItem `json:"list"`
		} `json:"data"`
	}
	if err = dao.client.Do(c, req, &res, dao.likeItemURL); err != nil {
		err = errors.Wrapf(err, "LikeData dao.client.Do sid(%d)", sid)
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrapf(ecode.Int(res.Code), "LikeData sid(%d)", sid)
		return
	}
	data = res.Data.List
	return
}

// RawSourceItemData get source data.
func (dao *Dao) RawSourceItemData(c context.Context, sid int64) (sids []int64, err error) {
	var res struct {
		Code int `json:"code"`
		Data struct {
			List []*struct {
				Data struct {
					Sid string `json:"sid"`
				} `json:"data"`
			} `json:"list"`
		} `json:"data"`
	}
	if err = dao.client.RESTfulGet(c, dao.sourceItemURL, metadata.String(c, metadata.RemoteIP), url.Values{}, &res, sid); err != nil {
		err = errors.Wrapf(err, "LikeData dao.client.RESTfulGet sid(%d)", sid)
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrapf(ecode.Int(res.Code), "LikeData sid(%d)", sid)
		return
	}
	for _, v := range res.Data.List {
		if sid, e := strconv.ParseInt(v.Data.Sid, 10, 64); e != nil {
			continue
		} else {
			sids = append(sids, sid)
		}
	}
	return
}

// SourceItem get source data json raw message.
func (dao *Dao) SourceItem(c context.Context, sid int64) (source json.RawMessage, err error) {
	var res struct {
		Code int             `json:"code"`
		Data json.RawMessage `json:"data"`
	}
	if err = dao.client.RESTfulGet(c, dao.sourceItemURL, metadata.String(c, metadata.RemoteIP), url.Values{}, &res, sid); err != nil {
		err = errors.Wrapf(err, "LikeData dao.client.RESTfulGet sid(%d)", sid)
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrapf(ecode.Int(res.Code), "LikeData sid(%d)", sid)
		return
	}
	source = res.Data
	return
}

// StoryLikeSum .
func (dao *Dao) StoryLikeSum(c context.Context, sid, mid int64) (res int64, err error) {
	var (
		conn = dao.redis.Get(c)
		now  = time.Now().Format("2006-01-02")
		key  = keyStoryLikeKey(sid, mid, now)
	)
	defer conn.Close()
	if res, err = redis.Int64(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
			res = -1
		} else {
			err = errors.Wrap(err, "redis.Do(get)")
		}
	}
	return
}

// IncrStoryLikeSum .
func (dao *Dao) IncrStoryLikeSum(c context.Context, sid, mid int64, score int64) (res int64, err error) {
	var (
		conn = dao.redis.Get(c)
		now  = time.Now().Format("2006-01-02")
		key  = keyStoryLikeKey(sid, mid, now)
	)
	defer conn.Close()
	if res, err = redis.Int64(conn.Do("INCRBY", key, score)); err != nil {
		err = errors.Wrap(err, "redis.Do(get)")
	}
	return
}

// SetLikeSum .
func (dao *Dao) SetLikeSum(c context.Context, sid, mid int64, sum int64) (err error) {
	var (
		conn = dao.redis.Get(c)
		now  = time.Now().Format("2006-01-02")
		key  = keyStoryLikeKey(sid, mid, now)
		res  bool
	)
	defer conn.Close()
	if res, err = redis.Bool(conn.Do("SETNX", key, sum)); err != nil {
		err = errors.Wrap(err, "redis.Bool(SETNX)")
		return
	}
	if res {
		conn.Do("EXPIRE", key, 86400)
	} else {
		err = errors.New("redis.Bool(SETNX) res false")
	}
	return
}

// StoryEachLikeSum  .
func (dao *Dao) StoryEachLikeSum(c context.Context, sid, mid, lid int64) (res int64, err error) {
	var (
		conn = dao.redis.Get(c)
		now  = time.Now().Format("2006-01-02")
		key  = keyStoryEachLike(sid, mid, lid, now)
	)
	defer conn.Close()
	if res, err = redis.Int64(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
			res = -1
		} else {
			err = errors.Wrap(err, "redis.Do(get)")
		}
	}
	return
}

// IncrStoryEachLikeAct .
func (dao *Dao) IncrStoryEachLikeAct(c context.Context, sid, mid, lid int64, score int64) (res int64, err error) {
	var (
		conn = dao.redis.Get(c)
		now  = time.Now().Format("2006-01-02")
		key  = keyStoryEachLike(sid, mid, lid, now)
	)
	defer conn.Close()
	if res, err = redis.Int64(conn.Do("INCRBY", key, score)); err != nil {
		err = errors.Wrap(err, "redis.Do(get)")
	}
	return
}

// SetEachLikeSum .
func (dao *Dao) SetEachLikeSum(c context.Context, sid, mid, lid int64, sum int64) (err error) {
	var (
		conn = dao.redis.Get(c)
		now  = time.Now().Format("2006-01-02")
		key  = keyStoryEachLike(sid, mid, lid, now)
		res  bool
	)
	defer conn.Close()
	if res, err = redis.Bool(conn.Do("SETNX", key, sum)); err != nil {
		err = errors.Wrap(err, "redis.Bool(SETNX)")
		return
	}
	if res {
		conn.Do("EXPIRE", key, 86400)
	} else {
		err = errors.New("redis.Bool(SETNX) res false")
	}
	return
}

// ListFromES .
func (dao *Dao) ListFromES(c context.Context, sid int64, order string, ps, pn int, seed int64) (res *like.ListInfo, err error) {
	actResult := new(struct {
		Result []struct {
			ID    int64      `json:"id"`
			Wid   int64      `json:"wid"`
			Ctime xtime.Time `json:"ctime"`
			Sid   int64      `json:"sid"`
			Type  int        `json:"type"`
			Mid   int64      `json:"mid"`
			State int        `json:"state"`
			Mtime xtime.Time `json:"mtime"`
			Likes int64      `json:"likes"`
			Click int64      `json:"click"`
			Coin  int64      `json:"coin"`
			Share int64      `json:"share"`
			Reply int64      `json:"reply"`
			Dm    int64      `json:"dm"`
			Fav   int64      `json:"fav"`
		} `json:"result"`
		Page *like.Page `json:"page"`
	})
	req := dao.es.NewRequest(_activity).Index(_activity).WhereEq("sid", sid).WhereEq("state", 1).Ps(ps).Pn(pn)
	if order != "" {
		req.Order(order, elastic.OrderDesc)
	}
	if seed > 0 {
		req.OrderRandomSeed(time.Unix(seed, 0).Format("2006-01-02 15:04:05"))
	}
	req.Fields("id", "sid", "wid", "mid", "type", "ctime", "mtime", "state", "click", "likes", "coin", "share", "reply", "dm", "fav")
	if err = req.Scan(c, &actResult); err != nil {
		err = errors.Wrap(err, "req.Scan")
		return
	}
	if len(actResult.Result) == 0 {
		return
	}
	res = &like.ListInfo{Page: actResult.Page, List: make([]*like.List, 0, len(actResult.Result))}
	for _, v := range actResult.Result {
		a := &like.List{
			Likes: v.Likes,
			Click: v.Click,
			Coin:  v.Coin,
			Share: v.Share,
			Reply: v.Reply,
			Dm:    v.Dm,
			Fav:   v.Fav,
			Item: &like.Item{
				ID:    v.ID,
				Wid:   v.Wid,
				Ctime: v.Ctime,
				Sid:   v.Sid,
				Type:  v.Type,
				Mid:   v.Mid,
				State: v.State,
				Mtime: v.Mtime,
			},
		}
		res.List = append(res.List, a)
	}
	return
}

// MultiTags .
func (dao *Dao) MultiTags(c context.Context, wids []int64) (tagList map[int64][]string, err error) {
	if len(wids) == 0 {
		return
	}
	var res struct {
		Code int                   `json:"code"`
		Data map[int64][]*like.Tag `json:"data"`
	}
	params := url.Values{}
	params.Set("aids", xstr.JoinInts(wids))
	if err = dao.client.Get(c, dao.tagURL, "", params, &res); err != nil {
		log.Error("MultiTags:dao.client.Get(%s) error(%+v)", dao.tagURL, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		err = ecode.Int(res.Code)
		return
	}
	tagList = make(map[int64][]string, len(res.Data))
	for k, v := range res.Data {
		if len(v) == 0 {
			continue
		}
		tagList[k] = make([]string, 0, len(v))
		for _, val := range v {
			tagList[k] = append(tagList[k], val.Name)
		}
	}
	return
}

// OidInfoFromES .
func (dao *Dao) OidInfoFromES(c context.Context, oids []int64, sType int) (res map[int64]*like.Item, err error) {
	actResult := new(struct {
		Result []struct {
			ID    int64      `json:"id"`
			Wid   int64      `json:"wid"`
			Ctime xtime.Time `json:"ctime"`
			Sid   int64      `json:"sid"`
			Type  int        `json:"type"`
			Mid   int64      `json:"mid"`
			State int        `json:"state"`
			Mtime xtime.Time `json:"mtime"`
			Likes int64      `json:"likes"`
			Click int64      `json:"click"`
			Coin  int64      `json:"coin"`
			Share int64      `json:"share"`
			Reply int64      `json:"reply"`
			Dm    int64      `json:"dm"`
			Fav   int64      `json:"fav"`
		} `json:"result"`
		Page *like.Page `json:"page"`
	})
	req := dao.es.NewRequest(_activity).Index(_activity).WhereIn("wid", oids).WhereEq("type", sType)
	req.Fields("id", "sid", "wid", "mid", "type", "ctime", "mtime", "state")
	if err = req.Scan(c, &actResult); err != nil {
		err = errors.Wrap(err, "req.Scan")
		return
	}
	if len(actResult.Result) == 0 {
		return
	}
	res = make(map[int64]*like.Item, len(actResult.Result))
	for _, v := range actResult.Result {
		res[v.Wid] = &like.Item{
			ID:    v.ID,
			Wid:   v.Wid,
			Ctime: v.Ctime,
			Sid:   v.Sid,
			Type:  v.Type,
			Mid:   v.Mid,
			State: v.State,
			Mtime: v.Mtime,
		}
	}
	return
}
