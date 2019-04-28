package dao

import (
	"context"
	"encoding/xml"
	"fmt"

	"go-common/app/job/main/dm2/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_keyIdx        = "i_%d_%d_%d_%d" // normal dm segment sortedset(s_type_oid_cnt_n, ctime, dmid)
	_keyIdxSub     = "s_%d_%d"       // subtitle dm sortedset(s_type_oid, progress, dmid)
	_keyIdxContent = "c_%d_%d"       // dm content hash(d_type_oid, dmid, xml)
)

func keyIdx(tp int32, oid, cnt, n int64) string {
	return fmt.Sprintf(_keyIdx, tp, oid, cnt, n)
}

// keyIdxSub return dm idx key.
func keyIdxSub(tp int32, oid int64) string {
	return fmt.Sprintf(_keyIdxSub, tp, oid)
}

// keyIdxContent return key of different dm.
func keyIdxContent(tp int32, oid int64) string {
	return fmt.Sprintf(_keyIdxContent, tp, oid)
}

// ExpireDMID set expire time of index.
func (d *Dao) ExpireDMID(c context.Context, tp int32, oid, cnt, n int64) (ok bool, err error) {
	key := keyIdx(tp, oid, cnt, n)
	conn := d.dmSegRds.Get(c)
	if ok, err = redis.Bool(conn.Do("EXPIRE", key, d.dmSegExpire)); err != nil {
		log.Error("conn.Do(EXPIRE %s) error(%v)", key, err)
	}
	conn.Close()
	return
}

// DMIDCache return dm ids.
func (d *Dao) DMIDCache(c context.Context, tp int32, oid int64, cnt, n, limit int64) (dmids []int64, err error) {
	var (
		conn = d.dmSegRds.Get(c)
		key  = keyIdx(tp, oid, cnt, n)
	)
	defer conn.Close()
	if dmids, err = redis.Int64s(conn.Do("ZRANGE", key, 0, -1)); err != nil {
		log.Error("DMIDSPCache.conn.DO(ZRANGEBYSCORE %s) error(%v)", key, err)
	}
	return
}

// AddDMIDCache add dmid(normal and special) to segment redis.
func (d *Dao) AddDMIDCache(c context.Context, tp int32, oid, cnt, n int64, dmids ...int64) (err error) {
	key := keyIdx(tp, oid, cnt, n)
	conn := d.dmSegRds.Get(c)
	defer conn.Close()
	for _, dmid := range dmids {
		if err = conn.Send("ZADD", key, dmid, dmid); err != nil {
			log.Error("conn.Send(ZADD %s) error(%v)", key, err)
			return
		}
	}
	if err = conn.Send("EXPIRE", key, d.dmSegExpire); err != nil {
		log.Error("conn.Send(EXPIRE %s) error(%v)", key, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush() error(%v)", err)
		return
	}
	for i := 0; i < len(dmids)+1; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// DelDMIDCache delete dm segment cache.
func (d *Dao) DelDMIDCache(c context.Context, tp int32, oid, cnt, n int64) (err error) {
	key := keyIdx(tp, oid, cnt, n)
	conn := d.dmSegRds.Get(c)
	if _, err = conn.Do("DEL", key); err != nil {
		log.Error("conn.Do(DEL %s) err(%v)", key, err)
	}
	conn.Close()
	return
}

// ExpireDMIDSubtitle set expire time of subtitle dmid.
func (d *Dao) ExpireDMIDSubtitle(c context.Context, tp int32, oid int64) (ok bool, err error) {
	key := keyIdxSub(tp, oid)
	conn := d.dmSegRds.Get(c)
	if ok, err = redis.Bool(conn.Do("EXPIRE", key, d.dmSegExpire)); err != nil {
		log.Error("conn.Do(EXPIRE %s) error(%v)", key, err)
	}
	conn.Close()
	return
}

// DMIDSubtitleCache get subtitle dmid.
func (d *Dao) DMIDSubtitleCache(c context.Context, tp int32, oid int64, ps, pe, limit int64) (dmids []int64, err error) {
	var (
		conn = d.dmSegRds.Get(c)
		key  = keyIdxSub(tp, oid)
	)
	defer conn.Close()
	if dmids, err = redis.Int64s(conn.Do("ZRANGEBYSCORE", key, ps, pe, "LIMIT", 0, limit)); err != nil {
		log.Error("conn.DO(ZRANGEBYSCORE %s) error(%v)", key, err)
	}
	return
}

// AddDMIDSubtitleCache add subtitle dmid to redis.
func (d *Dao) AddDMIDSubtitleCache(c context.Context, tp int32, oid int64, dms ...*model.DM) (err error) {
	key := keyIdxSub(tp, oid)
	conn := d.dmSegRds.Get(c)
	defer conn.Close()
	for _, dm := range dms {
		if err = conn.Send("ZADD", key, dm.Progress, dm.ID); err != nil {
			log.Error("conn.Send(ZADD %s) error(%v)", key, err)
			return
		}
	}
	if err = conn.Send("EXPIRE", key, d.dmSegExpire); err != nil {
		log.Error("conn.Send(EXPIRE %s) error(%v)", key, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush() error(%v)", err)
		return
	}
	for i := 0; i < len(dms)+1; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// DelDMIDSubtitleCache delete subtitle dmid cache.
func (d *Dao) DelDMIDSubtitleCache(c context.Context, tp int32, oid int64) (err error) {
	key := keyIdxSub(tp, oid)
	conn := d.dmSegRds.Get(c)
	if _, err = conn.Do("DEL", key); err != nil {
		log.Error("conn.Do(DEL %s) error(%v)", key, err)
	}
	conn.Close()
	return
}

// AddIdxContentCaches add index content cache to redis.
func (d *Dao) AddIdxContentCaches(c context.Context, tp int32, oid int64, dms ...*model.DM) (err error) {
	var (
		conn = d.dmSegRds.Get(c)
		key  = keyIdxContent(tp, oid)
	)
	defer conn.Close()
	for _, dm := range dms {
		if err = conn.Send("HSET", key, dm.ID, dm.ToXMLSeg()); err != nil {
			log.Error("conn.Send(HSET %s,%v) error(%v)", key, dm, err)
			return
		}
	}
	if err = conn.Send("EXPIRE", key, d.dmSegExpire); err != nil {
		log.Error("conn.Send(EXPIRE %s) error(%v)", key, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush() error(%v)", err)
		return
	}
	for i := 0; i <= len(dms); i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// DelIdxContentCaches del index content cache.
func (d *Dao) DelIdxContentCaches(c context.Context, tp int32, oid int64, dmids ...int64) (err error) {
	key := keyIdxContent(tp, oid)
	conn := d.dmSegRds.Get(c)
	args := []interface{}{key}
	for _, dmid := range dmids {
		args = append(args, dmid)
	}
	if _, err = conn.Do("HDEL", args...); err != nil {
		log.Error("conn.Do(HDEL %s) error(%v)", key, err)
	}
	conn.Close()
	return
}

// IdxContentCache get xml info by dmid.
func (d *Dao) IdxContentCache(c context.Context, tp int32, oid int64, dmids []int64) (res []byte, missed []int64, err error) {
	var (
		k      int
		dmid   int64
		values [][]byte
		key    = keyIdxContent(tp, oid)
		args   = []interface{}{key}
	)
	for _, dmid = range dmids {
		args = append(args, dmid)
	}
	conn := d.dmSegRds.Get(c)
	defer conn.Close()
	if values, err = redis.ByteSlices(conn.Do("HMGET", args...)); err != nil {
		log.Error("conn.Do(HMGET %v) error(%v)", args, err)
		if err == redis.ErrNil {
			return nil, nil, nil
		}
		return
	}
	for k, dmid = range dmids {
		if len(values[k]) == 0 {
			missed = append(missed, dmid)
			continue
		}
		res = append(res, values[k]...)
	}
	return
}

// IdxContentCacheV2 get elems info by dmid.
func (d *Dao) IdxContentCacheV2(c context.Context, tp int32, oid int64, dmids []int64) (elems []*model.Elem, missed []int64, err error) {
	var (
		values [][]byte
		key    = keyIdxContent(tp, oid)
		args   = []interface{}{key}
	)
	for _, dmid := range dmids {
		args = append(args, dmid)
	}
	conn := d.dmSegRds.Get(c)
	defer conn.Close()
	if values, err = redis.ByteSlices(conn.Do("HMGET", args...)); err != nil {
		log.Error("conn.Do(HMGET %v) error(%v)", args, err)
		if err == redis.ErrNil {
			return nil, nil, nil
		}
		return
	}
	for k, dmid := range dmids {
		if len(values[k]) == 0 {
			missed = append(missed, dmid)
			continue
		}
		elem, err := d.xmlToElem(values[k])
		if err != nil {
			missed = append(missed, dmid)
			continue
		}
		elems = append(elems, elem)
	}
	return
}

// 在缓存过渡期将<d p="弹幕ID,弹幕属性,播放时间,弹幕模式,字体大小,颜色,发送时间,弹幕池,用户hash id">弹幕内容</d>
// 装换为 model.Elem结构
func (d *Dao) xmlToElem(data []byte) (e *model.Elem, err error) {
	var v struct {
		XMLName   xml.Name `xml:"d"`
		Attribute string   `xml:"p,attr"`
		Content   string   `xml:",chardata"`
	}
	if err = xml.Unmarshal(data, &v); err != nil {
		return
	}
	e = &model.Elem{Content: v.Content, Attribute: v.Attribute}
	return
}
