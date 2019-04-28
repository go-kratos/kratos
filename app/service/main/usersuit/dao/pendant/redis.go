package pendant

import (
	"context"
	"strconv"

	"encoding/json"

	"go-common/app/service/main/usersuit/model"
	"go-common/library/cache/redis"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	_pendantPKG   = "pkg_" // key of
	_pendantEquip = "pe_"
)

func keyEquip(mid int64) string {
	return _pendantEquip + strconv.FormatInt(mid, 10)
}

// encode
func (d *Dao) encode(mid, pid, expires, tp int64, status, isVIP int32, pendant *model.Pendant) (res []byte, err error) {
	ft := &model.PendantPackage{Mid: mid, Pid: pid, Expires: expires, Type: tp, Status: status, IsVIP: isVIP, Pendant: pendant}
	return json.Marshal(ft)
}

// decode
func (d *Dao) decode(src []byte, v *model.PendantPackage) (err error) {
	return json.Unmarshal(src, v)
}

// AddPKGCache set package cache.
func (d *Dao) AddPKGCache(c context.Context, mid int64, info []*model.PendantPackage) (err error) {
	var (
		key  = _pendantPKG + strconv.FormatInt(mid, 10)
		args = redis.Args{}.Add(key)
	)
	for i := 0; i < len(info); i++ {
		var ef []byte
		if ef, err = d.encode(info[i].Mid, info[i].Pid, info[i].Expires, info[i].Type, info[i].Status, info[i].IsVIP, info[i].Pendant); err != nil {
			return
		}
		args = args.Add(i, ef)
	}
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("DEL", key); err != nil {
		log.Error("conn.Send(DEL, %s) error(%v)", key, err)
		return
	}
	if err = conn.Send("HMSET", args...); err != nil {
		log.Error("conn.Send(HMSET, %s) error(%v)", key, err)
		return
	}

	if err = conn.Send("EXPIRE", key, d.pendantExpire); err != nil {
		log.Error("conn.Send(EXPIRE, %s) error(%v)", key, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush() error(%v)", err)
		return
	}
	for i := 0; i < 3; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() %d error(%v)", i+1, err)
			break
		}
	}
	return
}

// PKGCache get package cache.
func (d *Dao) PKGCache(c context.Context, mid int64) (info []*model.PendantPackage, err error) {
	var (
		key = _pendantPKG + strconv.FormatInt(mid, 10)
		tmp = make(map[string]string, len(info))
	)
	conn := d.redis.Get(c)
	defer conn.Close()
	if tmp, err = redis.StringMap(conn.Do("HGETALL", key)); err != nil {
		return
	}
	if err == nil && len(tmp) > 0 {
		for i := 0; i < len(tmp); i++ {
			s := strconv.FormatInt(int64(i), 10)
			vf := &model.PendantPackage{}
			vf.Pendant = &model.Pendant{}
			if err = d.decode([]byte(tmp[s]), vf); err != nil {
				return
			}
			info = append(info, &model.PendantPackage{
				Mid:     vf.Mid,
				Pid:     vf.Pid,
				Expires: vf.Expires,
				Type:    vf.Type,
				Status:  vf.Status,
				IsVIP:   vf.IsVIP,
				Pendant: vf.Pendant,
			})
		}
	}
	return
}

// DelPKGCache del package cache
func (d *Dao) DelPKGCache(c context.Context, mid int64) (err error) {
	key := _pendantPKG + strconv.FormatInt(mid, 10)
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("DEL", key); err != nil {
		log.Error("conn.Send(DEL, %s) error(%v)", key, err)
		return
	}
	return
}

// equipCache return pendant info cache
func (d *Dao) equipCache(c context.Context, mid int64) (info *model.PendantEquip, err error) {
	var (
		item []byte
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if item, err = redis.Bytes(conn.Do("GET", keyEquip(mid))); err != nil {
		if err == redis.ErrNil {
			err = nil
		}
		return
	}
	if err = json.Unmarshal(item, &info); err != nil {
		log.Error("json.Unmarshal(%v) err(%v)", item, err)
	}
	return
}

// equipsCache obtain equips from redis .
func (d *Dao) equipsCache(c context.Context, mids []int64) (map[int64]*model.PendantEquip, []int64, error) {
	var (
		err  error
		bss  [][]byte
		key  string
		args = redis.Args{}
		conn = d.redis.Get(c)
	)

	for _, v := range mids {
		key = keyEquip(v)
		args = args.Add(key)
	}
	defer conn.Close()
	if bss, err = redis.ByteSlices(conn.Do("MGET", args...)); err != nil {
		if err == redis.ErrNil {
			return nil, nil, nil
		}
		log.Error("Failed mget equip: keys: %+v: %+v", args, err)
		return nil, nil, err
	}
	info := make(map[int64]*model.PendantEquip, len(mids))
	for _, bs := range bss {
		if bs == nil {
			continue
		}
		pe := &model.PendantEquip{}
		if err = json.Unmarshal(bs, pe); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", string(bs), err)
			err = nil
			continue
		}
		info[pe.Mid] = pe
	}
	missed := make([]int64, 0, len(mids))
	for _, mid := range mids {
		if _, ok := info[mid]; !ok {
			missed = append(missed, mid)
		}
	}
	return info, missed, nil
}

// AddEquipCache set pendant info cache
func (d *Dao) AddEquipCache(c context.Context, mid int64, info *model.PendantEquip) (err error) {
	var (
		key    = keyEquip(mid)
		values []byte
		conn   = d.redis.Get(c)
	)
	defer conn.Close()
	if values, err = json.Marshal(info); err != nil {
		return
	}
	if err = conn.Send("SET", keyEquip(mid), values); err != nil {
		log.Error("conn.Send(SET, %s, %d) error(%v)", key, values, err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.pendantExpire); err != nil {
		log.Error("conn.Send(Expire, %s, %d) error(%v)", key, d.pendantExpire, err)
		return
	}
	if err = conn.Flush(); err != nil {
		err = errors.Wrap(err, "conn.Send Flush")
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			err = errors.Wrap(err, "conn.Send conn.Receive()")
			return
		}
	}
	return
}

// AddEquipsCache mset equips info to caache .
func (d *Dao) AddEquipsCache(c context.Context, equips map[int64]*model.PendantEquip) (err error) {
	var (
		bs      []byte
		key     string
		keys    []string
		argsMid = redis.Args{}
		conn    = d.redis.Get(c)
	)
	defer conn.Close()
	for _, v := range equips {
		if bs, err = json.Marshal(v); err != nil {
			log.Error("json.Marshal err(%v)", err)
			continue
		}
		key = keyEquip(v.Mid)
		keys = append(keys, key)
		argsMid = argsMid.Add(key).Add(string(bs))
	}
	if err = conn.Send("MSET", argsMid...); err != nil {
		err = errors.Wrap(err, "conn.Send(MSET) error")
		return
	}
	count := 1
	for _, v := range keys {
		count++
		if err = conn.Send("EXPIRE", v, d.pendantExpire); err != nil {
			err = errors.Wrap(err, "conn.Send error")
			return
		}
	}
	if err = conn.Flush(); err != nil {
		err = errors.Wrap(err, "conn.Send Flush")
		return
	}
	for i := 0; i < count; i++ {
		if _, err = conn.Receive(); err != nil {
			err = errors.Wrap(err, "conn.Send conn.Receive()")
			return
		}
	}
	return
}

// DelEquipCache set pendant info cache
func (d *Dao) DelEquipCache(c context.Context, mid int64) (err error) {
	key := keyEquip(mid)
	conn := d.redis.Get(c)
	defer conn.Close()
	if _, err = conn.Do("DEL", key); err != nil {
		log.Error("conn.Do(DEL, %s) error(%v)", key, err)
	}
	return
}

// DelEquipsCache del batch equip cache .
func (d *Dao) DelEquipsCache(c context.Context, mids []int64) (err error) {
	var (
		args = redis.Args{}
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	for _, v := range mids {
		args = args.Add(keyEquip(v))
	}
	if _, err = conn.Do("DEL", args...); err != nil {
		log.Error("conn.Do(DEL, %s) error(%v)", args, err)
	}
	return
}
