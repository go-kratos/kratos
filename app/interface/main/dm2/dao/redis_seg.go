package dao

import (
	"context"
	"encoding/xml"
	"fmt"

	"go-common/app/interface/main/dm2/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_keyIdx        = "i_%d_%d_%d_%d" // normal dm segment sortedset(s_type_oid_cnt_n, ctime, dmid)
	_keyIdxSub     = "s_%d_%d"       // subtitle dm sortedset(s_type_oid, progress, dmid)
	_keyIdxSpe     = "spe_%d_%d"     // special dm sortedset(spe_type_oid, progress,dmid)
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

func keyIdxSpecial(tp int32, oid int64) string {
	return fmt.Sprintf(_keyIdxSpe, tp, oid)
}

// DMIDSubtitlesCache return subtitle dm ids.
func (d *Dao) DMIDSubtitlesCache(c context.Context, tp int32, oid int64, ps, pe, limit int64) (dmids []int64, err error) {
	var (
		conn   = d.dmSegRds.Get(c)
		keySub = keyIdxSub(tp, oid)
	)
	defer conn.Close()
	if dmids, err = redis.Int64s(conn.Do("ZRANGEBYSCORE", keySub, ps, pe, "LIMIT", 0, limit)); err != nil {
		log.Error("conn.DO(ZRANGEBYSCORE %s) error(%v)", keySub, err)
	}
	if len(dmids) > 0 {
		PromCacheHit("dm_seg_dmid", 1)
	} else {
		PromCacheMiss("dm_seg_dmid", 1)
	}
	return
}

// DMIDCache return dm index id.
func (d *Dao) DMIDCache(c context.Context, tp int32, oid int64, cnt, n, limit int64) (dmids []int64, err error) {
	var (
		conn   = d.dmSegRds.Get(c)
		keyIdx = keyIdx(tp, oid, cnt, n)
	)
	defer conn.Close()
	if dmids, err = redis.Int64s(conn.Do("ZRANGE", keyIdx, 0, -1)); err != nil {
		log.Error("DMIDSPCache.conn.DO(ZRANGEBYSCORE %s) error(%v)", keyIdx, err)
	}
	if len(dmids) > 0 {
		PromCacheHit("dm_seg_dmid", 1)
	} else {
		PromCacheMiss("dm_seg_dmid", 1)
	}
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
	PromCacheMiss("dmid_content", int64(len(missed)))
	PromCacheHit("dmid_content", int64(len(dmids)-len(missed)))
	return
}

// IdxContentCacheV2 get elems info by dmid.
func (d *Dao) IdxContentCacheV2(c context.Context, tp int32, oid int64, dmids []int64) (elems []*model.Elem, missed []int64, err error) {
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
		if err == redis.ErrNil {
			err = nil
			missed = dmids
		} else {
			log.Error("conn.Do(HMGET %v) error(%v)", args, err)
		}
		return
	}
	for k, dmid = range dmids {
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
	PromCacheMiss("dmid_elem", int64(len(missed)))
	PromCacheHit("dmid_elem", int64(len(dmids)-len(missed)))
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

// AddIdxContentCache add index content cache to redis.
func (d *Dao) AddIdxContentCache(c context.Context, tp int32, oid int64, dms []*model.DM, realname bool) (err error) {
	var (
		key  string
		conn = d.dmSegRds.Get(c)
	)
	defer conn.Close()
	for _, dm := range dms {
		key = keyIdxContent(tp, oid)
		if err = conn.Send("HSET", key, dm.ID, dm.ToXMLSeg(realname)); err != nil {
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

// DMIDSpecialsCache return special dmids
func (d *Dao) DMIDSpecialsCache(c context.Context, tp int32, oid int64) (dmids []int64, err error) {
	var (
		conn   = d.dmSegRds.Get(c)
		keySpe = keyIdxSpecial(tp, oid)
	)
	defer conn.Close()
	if dmids, err = redis.Int64s(conn.Do("ZRANGE", keySpe, 0, -1)); err != nil {
		log.Error("conn.DO(ZRANGE %s) error(%v)", keySpe, err)
	}
	if len(dmids) > 0 {
		PromCacheHit("dm_spe_dmid", 1)
	} else {
		PromCacheMiss("dm_spe_dmid", 1)
	}
	return
}
