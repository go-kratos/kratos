package dao

import (
	"context"
	"fmt"
	"hash/crc32"
	"strconv"

	"go-common/app/interface/main/dm2/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

var (
	_lockvalue = []byte("1")
)

const (
	_prefixPubLock     = "dm_pub_"
	_prefixCharPubLock = "dm_pub_char_"
	_prefixXML         = "dm_xml_"
	_prefixSub         = "s_"
	_prefixAjax        = "dm_ajax_"
	_prefixDMJudge     = "dm_judge_%d_%d_%d"
	_prefixDMLimit     = "dm_limit_"
	_prefixAdvanceCmt  = "cache:AdvanceComment:%d@%d.%s"
	_prefixAdvLock     = "adv_lock_"
	_prefixAdvance     = "cache:AdvanceComment:"
	_prefixHis         = "dm_his_%d_%d_%d"
	_prefixHisIdx      = "dm_his_idx_%d_%d_%s"
	_prefixDMMask      = "dm_mask_%d_%d_%d"
)

func keyMsgPubLock(mid, color, rnd int64, mode, fontsize int32, ip, msg string) string {
	crcStr := fmt.Sprintf("%d_%s_%s_%d_%d_%d_%d", mid, msg, ip, fontsize, color, mode, rnd)
	return _prefixPubLock + fmt.Sprint(crc32.ChecksumIEEE([]byte(crcStr)))
}

func keyOidPubLock(mid, oid int64, ip string) string {
	crcStr := fmt.Sprintf("%d_%s_%d", mid, ip, oid)
	return _prefixPubLock + fmt.Sprint(crc32.ChecksumIEEE([]byte(crcStr)))
}

func keyPubCntLock(mid, color int64, mode, fontsize int32, ip, msg string) string {
	crcStr := fmt.Sprintf("%d_%s_%s_%d_%d_%d", mid, msg, ip, fontsize, color, mode)
	return _prefixPubLock + fmt.Sprint(crc32.ChecksumIEEE([]byte(crcStr)))
}

func keyCharPubLock(mid, oid int64) string {
	return _prefixCharPubLock + fmt.Sprintf("%d_%d", mid, oid)
}

func keyXML(oid int64) string {
	return _prefixXML + strconv.FormatInt(oid, 10)
}

func keySubject(tp int32, oid int64) string {
	return _prefixSub + fmt.Sprintf("%d_%d", tp, oid)
}

func keyAjax(oid int64) string {
	return _prefixAjax + strconv.FormatInt(oid, 10)
}

func keyJudge(tp int8, oid, dmid int64) string {
	return fmt.Sprintf(_prefixDMJudge, tp, oid, dmid)
}

func keyDMLimitMid(mid int64) string {
	return _prefixDMLimit + strconv.FormatInt(mid, 10)
}

func keyAdvanceCmt(mid, oid int64, mode string) string {
	return fmt.Sprintf(_prefixAdvanceCmt, mid, oid, mode)
}

func keyAdvLock(mid, cid int64) string {
	return _prefixAdvLock + strconv.FormatInt(mid, 10) + "_" + strconv.FormatInt(cid, 10)
}

func keyHistory(tp int32, oid, timestamp int64) string {
	return fmt.Sprintf(_prefixHis, tp, oid, timestamp)
}

func keyHistoryIdx(tp int32, oid int64, month string) string {
	return fmt.Sprintf(_prefixHisIdx, tp, oid, month)
}

func keyDMMask(tp int32, oid int64, plat int8) string {
	return fmt.Sprintf(_prefixDMMask, tp, oid, plat)
}

// SubjectCache get subject from memcache.
func (d *Dao) SubjectCache(c context.Context, tp int32, oid int64) (sub *model.Subject, err error) {
	var (
		conn = d.dmMC.Get(c)
		key  = keySubject(tp, oid)
		rp   *memcache.Item
	)
	defer conn.Close()
	if rp, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			PromCacheMiss("dm_subject", 1)
			sub = nil
			err = nil
		} else {
			log.Error("mc.Get(%s) error(%v)", key, err)
		}
		return
	}
	sub = &model.Subject{}
	PromCacheHit("dm_subject", 1)
	if err = conn.Scan(rp, &sub); err != nil {
		log.Error("mc.Scan() error(%v)", err)
	}
	return
}

// SubjectsCache multi get subject from memcache.
func (d *Dao) SubjectsCache(c context.Context, tp int32, oids []int64) (cached map[int64]*model.Subject, missed []int64, err error) {
	var (
		conn   = d.dmMC.Get(c)
		keys   []string
		oidMap = make(map[string]int64)
	)
	cached = make(map[int64]*model.Subject)
	defer conn.Close()
	for _, oid := range oids {
		k := keySubject(tp, oid)
		if _, ok := oidMap[k]; !ok {
			keys = append(keys, k)
			oidMap[k] = oid
		}
	}
	rs, err := conn.GetMulti(keys)
	if err != nil {
		log.Error("conn.GetMulti(%v) error(%v)", keys, err)
		return
	}
	for k, r := range rs {
		sub := &model.Subject{}
		if err = conn.Scan(r, sub); err != nil {
			log.Error("conn.Scan(%s) error(%v)", r.Value, err)
			err = nil
			continue
		}
		cached[oidMap[k]] = sub
		// delete hit key
		delete(oidMap, k)
	}
	// missed key
	missed = make([]int64, 0, len(oidMap))
	for _, oid := range oidMap {
		missed = append(missed, oid)
	}
	PromCacheHit("dm_subjects", int64(len(cached)))
	PromCacheMiss("dm_subjects", int64(len(missed)))
	return
}

// AddSubjectCache add subject cache.
func (d *Dao) AddSubjectCache(c context.Context, sub *model.Subject) (err error) {
	var (
		conn = d.dmMC.Get(c)
		key  = keySubject(sub.Type, sub.Oid)
	)
	defer conn.Close()
	item := &memcache.Item{
		Key:        key,
		Object:     sub,
		Flags:      memcache.FlagJSON,
		Expiration: d.subjectExpire,
	}
	if err = conn.Set(item); err != nil {
		log.Error("conn.Set(%v) error(%v)", item, err)
	}
	return
}

// DelXMLCache delete xml content.
func (d *Dao) DelXMLCache(c context.Context, oid int64) (err error) {
	conn := d.dmMC.Get(c)
	if err = conn.Delete(keyXML(oid)); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Delete(%s) error(%v)", keyXML(oid), err)
		}
	}
	conn.Close()
	return
}

// AddXMLCache add xml content to memcache.
func (d *Dao) AddXMLCache(c context.Context, oid int64, value []byte) (err error) {
	conn := d.dmMC.Get(c)
	defer conn.Close()
	item := &memcache.Item{
		Key:        keyXML(oid),
		Value:      value,
		Expiration: d.dmExpire,
	}
	if err = conn.Set(item); err != nil {
		log.Error("conn.Set(%s) error(%v)", keyXML(oid), err)
	}
	return
}

// XMLCache get xml content.
func (d *Dao) XMLCache(c context.Context, oid int64) (data []byte, err error) {
	key := keyXML(oid)
	conn := d.dmMC.Get(c)
	defer conn.Close()
	item, err := conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			PromCacheMiss("dm_xml", 1)
		} else {
			log.Error("mc.Get(%s) error(%v)", key, err)
		}
		return
	}
	PromCacheHit("dm_xml", 1)
	data = item.Value
	return
}

// AjaxDMCache get ajax dm from memcache.
func (d *Dao) AjaxDMCache(c context.Context, oid int64) (msgs []string, err error) {
	conn := d.dmMC.Get(c)
	defer conn.Close()
	key := keyAjax(oid)
	msgs = make([]string, 0)
	item, err := conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			PromCacheMiss("dm_ajax", 1)
		} else {
			log.Error("conn.Get(%s) error(%v)", key, err)
		}
		return
	}
	PromCacheHit("dm_ajax", 1)
	if err = conn.Scan(item, &msgs); err != nil {
		log.Error("conn.Scan(%v) error(%v)", item, err)
	}
	return
}

// AddAjaxDMCache set ajax dm to memcache.
func (d *Dao) AddAjaxDMCache(c context.Context, oid int64, msgs []string) (err error) {
	conn := d.dmMC.Get(c)
	defer conn.Close()
	key := keyAjax(oid)
	item := &memcache.Item{Key: key, Object: msgs, Flags: memcache.FlagJSON, Expiration: d.ajaxExpire}
	if err = conn.Set(item); err != nil {
		log.Error("conn.Set(%v) error(%v)", item, err)
	}
	return
}

// SetDMJudgeCache dm fengjiwei list
func (d *Dao) SetDMJudgeCache(c context.Context, tp int8, oid, dmid int64, l *model.JudgeDMList) (err error) {
	key := keyJudge(tp, oid, dmid)
	conn := d.dmMC.Get(c)
	defer conn.Close()
	item := memcache.Item{
		Key:        key,
		Object:     l,
		Expiration: 86400 * 30,
		Flags:      memcache.FlagJSON,
	}
	if err = conn.Set(&item); err != nil {
		log.Error("DMJudge:mc.Set(%v) error(%v)", item, err)
	}
	return
}

// DMJudgeCache memcache cache of dm judge list.
func (d *Dao) DMJudgeCache(c context.Context, tp int8, oid, dmid int64) (l *model.JudgeDMList, err error) {
	key := keyJudge(tp, oid, dmid)
	conn := d.dmMC.Get(c)
	defer conn.Close()
	item, err := conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			l = nil
			PromCacheMiss("dm_judge", 1)
		} else {
			log.Error("mc.Get(key:%s) error(%v)", key, err)
		}
		return
	}
	PromCacheHit("dm_judge", 1)
	if err = conn.Scan(item, &l); err != nil {
		log.Error("conn.Scan(%v) error(%v)", item, err)
	}
	return
}

// AddMsgPubLock set publock into memcache
func (d *Dao) AddMsgPubLock(c context.Context, mid, color, rnd int64, mode, fontsize int32, ip, msg string) (err error) {
	conn := d.dmMC.Get(c)
	item := memcache.Item{
		Key:        keyMsgPubLock(mid, color, rnd, mode, fontsize, ip, msg),
		Value:      _lockvalue,
		Expiration: 300,
		Flags:      memcache.FlagRAW,
	}
	if err = conn.Set(&item); err != nil {
		log.Error("mc.Set(%v) error(%v)", item, err)
	}
	conn.Close()
	return
}

// MsgPublock get publock
func (d *Dao) MsgPublock(c context.Context, mid, color, rnd int64, mode, fontsize int32, ip, msg string) (cached bool, err error) {
	conn := d.dmMC.Get(c)
	defer conn.Close()
	key := keyMsgPubLock(mid, color, rnd, mode, fontsize, ip, msg)
	if _, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			cached = false
			err = nil
		} else {
			log.Error("mc.Get(%s) error(%v)", key, err)
		}
		return
	}
	cached = true
	return
}

// AddOidPubLock set publock into memcache
func (d *Dao) AddOidPubLock(c context.Context, mid, oid int64, ip string) (err error) {
	conn := d.dmMC.Get(c)
	item := memcache.Item{
		Key:        keyOidPubLock(mid, oid, ip),
		Value:      _lockvalue,
		Expiration: 4,
		Flags:      memcache.FlagRAW,
	}
	if err = conn.Set(&item); err != nil {
		log.Error("mc.Set(%v) error(%v)", item, err)
	}
	conn.Close()
	return
}

// OidPubLock get publock
func (d *Dao) OidPubLock(c context.Context, mid, oid int64, ip string) (cached bool, err error) {
	conn := d.dmMC.Get(c)
	defer conn.Close()
	key := keyOidPubLock(mid, oid, ip)
	if _, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			cached = false
			err = nil
		} else {
			log.Error("mc.Get(%s) error(%v)", key, err)
		}
		return
	}
	cached = true
	return
}

// AddDMLimitCache add Dmlimit in cache
func (d *Dao) AddDMLimitCache(c context.Context, mid int64, limiter *model.Limiter) (err error) {
	var (
		conn = d.dmMC.Get(c)
		key  = keyDMLimitMid(mid)
	)
	defer conn.Close()
	item := &memcache.Item{
		Key:        key,
		Object:     limiter,
		Flags:      memcache.FlagJSON,
		Expiration: 600,
	}
	if err = conn.Set(item); err != nil {
		log.Error("conn.Set(%v) error(%v)", item, err)
	}
	return
}

// DMLimitCache get dm limit from memcache.
func (d *Dao) DMLimitCache(c context.Context, mid int64) (limiter *model.Limiter, err error) {
	var (
		conn = d.dmMC.Get(c)
		key  = keyDMLimitMid(mid)
		rp   *memcache.Item
	)
	limiter = &model.Limiter{}
	defer conn.Close()
	if rp, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			limiter = nil
			err = nil
		} else {
			log.Error("mc.Get(%s) error(%v)", key, err)
		}
		return
	}
	if err = conn.Scan(rp, &limiter); err != nil {
		log.Error("mc.Scan() error(%v)", err)
	}
	return
}

// AddAdvanceCmtCache return
func (d *Dao) AddAdvanceCmtCache(c context.Context, oid, mid int64, mode string, adv *model.AdvanceCmt) (err error) {
	var (
		conn = d.filterMC.Get(c)
		key  = keyAdvanceCmt(mid, oid, mode)
	)
	defer conn.Close()
	item := &memcache.Item{
		Key:        key,
		Object:     adv,
		Flags:      memcache.FlagJSON,
		Expiration: d.filterMCExpire,
	}
	if err = conn.Set(item); err != nil {
		log.Error("conn.Set(%v) error(%v)", item, err)
	}
	return
}

// AdvanceCmtCache return advance comment from memcache.
func (d *Dao) AdvanceCmtCache(c context.Context, oid, mid int64, mode string) (adv *model.AdvanceCmt, err error) {
	var (
		conn = d.filterMC.Get(c)
		key  = keyAdvanceCmt(mid, oid, mode)
		rp   *memcache.Item
	)
	defer conn.Close()
	adv = &model.AdvanceCmt{}
	if rp, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			adv = nil
			err = nil
			PromCacheMiss("dm_advance", 1)
		} else {
			log.Error("mc.Get(%s) error(%v)", key, err)
		}
		return
	}
	PromCacheHit("dm_advance", 1)
	if err = conn.Scan(rp, &adv); err != nil {
		log.Error("mc.Scan() error(%v)", err)
	}
	return
}

// AddAdvanceLock 购买高级弹幕锁
func (d *Dao) AddAdvanceLock(c context.Context, mid, cid int64) (succeed bool) {
	var (
		key  = keyAdvLock(mid, cid)
		conn = d.filterMC.Get(c)
	)
	defer conn.Close()
	item := &memcache.Item{
		Key:        key,
		Value:      []byte("3"),
		Expiration: 60,
	}
	if err := conn.Add(item); err != nil {
		succeed = false
		log.Error("conn.Add(%s) error(%v)", key, err)
	} else {
		succeed = true
	}
	return
}

// DelAdvanceLock 删除购买高级弹幕锁
func (d *Dao) DelAdvanceLock(c context.Context, mid, cid int64) (err error) {
	var (
		key  = keyAdvLock(mid, cid)
		conn = d.filterMC.Get(c)
	)
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Delete(%s) error(%v)", key, err)
		}
	}
	conn.Close()
	return
}

// DelAdvCache delete from cache.
func (d *Dao) DelAdvCache(c context.Context, mid, cid int64, mode string) (err error) {
	var (
		key = _prefixAdvance + strconv.FormatInt(mid, 10) + "@" + strconv.FormatInt(cid, 10) + "." + mode
	)
	conn := d.filterMC.Get(c)
	err = conn.Delete(key)
	conn.Close()
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("memcache.Delete(%s) error(%v)", key, err)
		}
	}
	return
}

// AddHistoryCache add dm history to memcache.
func (d *Dao) AddHistoryCache(c context.Context, tp int32, oid, timestamp int64, value []byte) (err error) {
	var (
		conn = d.dmMC.Get(c)
		key  = keyHistory(tp, oid, timestamp)
	)
	defer conn.Close()
	item := &memcache.Item{
		Key:        key,
		Value:      value,
		Expiration: d.historyExpire,
	}
	if err = conn.Set(item); err != nil {
		log.Error("conn.Set(%s) error(%v)", key, err)
	}
	return
}

// HistoryCache history cache.
func (d *Dao) HistoryCache(c context.Context, tp int32, oid, timestamp int64) (data []byte, err error) {
	var (
		conn = d.dmMC.Get(c)
		key  = keyHistory(tp, oid, timestamp)
	)
	defer conn.Close()
	item, err := conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			PromCacheMiss("dm_history", 1)
		} else {
			log.Error("mc.Get(%s) error(%v)", key, err)
		}
		return
	}
	PromCacheHit("dm_history", 1)
	data = item.Value
	return
}

// AddHisIdxCache add dm history date index to memcache.
func (d *Dao) AddHisIdxCache(c context.Context, tp int32, oid int64, month string, dates []string) (err error) {
	var (
		conn = d.dmMC.Get(c)
		key  = keyHistoryIdx(tp, oid, month)
	)
	defer conn.Close()
	item := &memcache.Item{
		Key:        key,
		Object:     dates,
		Flags:      memcache.FlagJSON,
		Expiration: d.historyExpire,
	}
	if err = conn.Set(item); err != nil {
		log.Error("conn.Set(%s) error(%v)", key, err)
	}
	return
}

// HistoryIdxCache get history date index.
func (d *Dao) HistoryIdxCache(c context.Context, tp int32, oid int64, month string) (dates []string, err error) {
	var (
		conn = d.dmMC.Get(c)
		key  = keyHistoryIdx(tp, oid, month)
	)
	defer conn.Close()
	item, err := conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			PromCacheMiss("dm_history_index", 1)
		} else {
			log.Error("mc.Get(%s) error(%v)", key, err)
		}
		return
	}
	PromCacheHit("dm_history_index", 1)
	if err = conn.Scan(item, &dates); err != nil {
		log.Error("conn.Scan(%+v) error(%v)", item, err)
	}
	return
}

// DMMaskCache get dm mask cache
func (d *Dao) DMMaskCache(c context.Context, tp int32, oid int64, plat int8) (mask *model.Mask, err error) {
	var (
		conn = d.dmMC.Get(c)
		key  = keyDMMask(tp, oid, plat)
		item *memcache.Item
	)
	defer conn.Close()
	if item, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			mask = nil
			PromCacheMiss("dm_mask", 1)
		} else {
			log.Error("mc.Get(%s) error(%v)", key, err)
		}
		return
	}
	mask = &model.Mask{}
	PromCacheHit("dm_mask", 1)
	if err = conn.Scan(item, &mask); err != nil {
		log.Error("conn.Scan(%+v) error(%v)", item, err)
	}
	return
}

// AddMaskCache add dm mask cache
func (d *Dao) AddMaskCache(c context.Context, tp int32, mask *model.Mask) (err error) {
	var (
		conn = d.dmMC.Get(c)
		key  = keyDMMask(tp, mask.Cid, mask.Plat)
	)
	defer conn.Close()
	item := &memcache.Item{
		Key:        key,
		Object:     mask,
		Flags:      memcache.FlagJSON,
		Expiration: d.dmMaskExpire,
	}
	if err = conn.Set(item); err != nil {
		log.Error("conn.Set(%+v) error(%v)", item, err)
	}
	return
}
