package dao

import (
	"context"
	"fmt"

	"go-common/app/interface/main/dm2/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_fmtSubtitle          = "s_subtitle_%d_%d"
	_fmtVideoSubtitle     = "s_video_%d_%d"
	_fmtSubtitleDraft     = "s_draft_%v_%v_%v_%v"
	_fmtSubtitleSubject   = "s_subtitle_allow_%d"
	_fmtSubtitleReportTag = "s_subtitle_report_%d_%d" // s_subtitle_report_bid_rid

)

func (d *Dao) subtitleKey(oid int64, subtitleID int64) string {
	return fmt.Sprintf(_fmtSubtitle, oid, subtitleID)
}

func (d *Dao) subtitleVideoKey(oid int64, tp int32) string {
	return fmt.Sprintf(_fmtVideoSubtitle, oid, tp)
}

func (d *Dao) subtitleDraftKey(oid int64, tp int32, mid int64, lan uint8) string {
	return fmt.Sprintf(_fmtSubtitleDraft, oid, tp, mid, lan)
}

func (d *Dao) subtitleSubjectKey(aid int64) string {
	return fmt.Sprintf(_fmtSubtitleSubject, aid)
}

func (d *Dao) subtitleReportTagKey(bid, rid int64) string {
	return fmt.Sprintf(_fmtSubtitleReportTag, bid, rid)
}

// SetVideoSubtitleCache .
func (d *Dao) SetVideoSubtitleCache(c context.Context, oid int64, tp int32, res *model.VideoSubtitleCache) (err error) {
	var (
		item *memcache.Item
		conn = d.subtitleMc.Get(c)
		key  = d.subtitleVideoKey(oid, tp)
	)
	defer conn.Close()
	item = &memcache.Item{
		Key:        key,
		Object:     res,
		Flags:      memcache.FlagJSON | memcache.FlagGzip,
		Expiration: d.subtitleMcExpire,
	}
	if err = conn.Set(item); err != nil {
		log.Error("conn.Set(%v) error(%v)", item, err)
	}
	return
}

// VideoSubtitleCache .
func (d *Dao) VideoSubtitleCache(c context.Context, oid int64, tp int32) (res *model.VideoSubtitleCache, err error) {
	var (
		item *memcache.Item
		conn = d.subtitleMc.Get(c)
		key  = d.subtitleVideoKey(oid, tp)
	)
	defer conn.Close()
	if item, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			res = nil
			return
		}
		log.Error("memcache.Get(%s) error(%v)", key, err)
		return
	}
	if err = conn.Scan(item, &res); err != nil {
		log.Error("mc.Scan() error(%v)", err)
	}
	return
}

// DelVideoSubtitleCache .
func (d *Dao) DelVideoSubtitleCache(c context.Context, oid int64, tp int32) (err error) {
	var (
		key  = d.subtitleVideoKey(oid, tp)
		conn = d.subtitleMc.Get(c)
	)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("memcache.Delete(%s) error(%v)", key, err)
		}
	}
	return
}

// SubtitleDraftCache .
func (d *Dao) SubtitleDraftCache(c context.Context, oid int64, tp int32, mid int64, lan uint8) (subtitle *model.Subtitle, err error) {
	var (
		item *memcache.Item
		conn = d.subtitleMc.Get(c)
		key  = d.subtitleDraftKey(oid, tp, mid, lan)
	)
	defer conn.Close()
	if item, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("memcache.Get(%s) error(%v)", key, err)
		}
		return
	}
	if err = conn.Scan(item, &subtitle); err != nil {
		log.Error("mc.Scan() error(%v)", err)
	}
	return
}

// SetSubtitleDraftCache .
func (d *Dao) SetSubtitleDraftCache(c context.Context, subtitle *model.Subtitle) (err error) {
	var (
		item *memcache.Item
		conn = d.subtitleMc.Get(c)
		key  = d.subtitleDraftKey(subtitle.Oid, subtitle.Type, subtitle.Mid, subtitle.Lan)
	)
	defer conn.Close()
	item = &memcache.Item{
		Key:        key,
		Object:     subtitle,
		Flags:      memcache.FlagJSON | memcache.FlagGzip,
		Expiration: d.subtitleMcExpire,
	}
	if err = conn.Set(item); err != nil {
		log.Error("conn.Set(%v) error(%v)", item, err)
	}
	return
}

// DelSubtitleDraftCache .
func (d *Dao) DelSubtitleDraftCache(c context.Context, oid int64, tp int32, mid int64, lan uint8) (err error) {
	var (
		key  = d.subtitleDraftKey(oid, tp, mid, lan)
		conn = d.subtitleMc.Get(c)
	)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("memcache.Delete(%s) error(%v)", key, err)
		}
	}
	return
}

// SubtitlesCache .
func (d *Dao) SubtitlesCache(c context.Context, oid int64, subtitleIds []int64) (res map[int64]*model.Subtitle, missed []int64, err error) {
	var (
		conn          = d.subtitleMc.Get(c)
		keys          []string
		subtitleIDMap = make(map[string]int64)
	)
	res = make(map[int64]*model.Subtitle)
	defer conn.Close()
	for _, subtitleID := range subtitleIds {
		k := d.subtitleKey(oid, subtitleID)
		if _, ok := subtitleIDMap[k]; !ok {
			keys = append(keys, k)
			subtitleIDMap[k] = subtitleID
		}
	}
	rs, err := conn.GetMulti(keys)
	if err != nil {
		log.Error("conn.GetMulti(%v) error(%v)", keys, err)
		return
	}
	for k, r := range rs {
		st := &model.Subtitle{}
		if err = conn.Scan(r, st); err != nil {
			log.Error("conn.Scan(%s) error(%v)", r.Value, err)
			err = nil
			continue
		}
		res[subtitleIDMap[k]] = st
		// delete hit key
		delete(subtitleIDMap, k)
	}
	// missed key
	missed = make([]int64, 0, len(subtitleIDMap))
	for _, subtitleID := range subtitleIDMap {
		missed = append(missed, subtitleID)
	}
	return
}

// SubtitleCache .
func (d *Dao) SubtitleCache(c context.Context, oid int64, subtitleID int64) (subtitle *model.Subtitle, err error) {
	var (
		item *memcache.Item
		conn = d.subtitleMc.Get(c)
		key  = d.subtitleKey(oid, subtitleID)
	)
	defer conn.Close()
	if item, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("memcache.Get(%s) error(%v)", key, err)
		}
		return
	}
	if err = conn.Scan(item, &subtitle); err != nil {
		log.Error("mc.Scan() error(%v)", err)
	}
	return
}

// SetSubtitleCache .
func (d *Dao) SetSubtitleCache(c context.Context, subtitle *model.Subtitle) (err error) {
	var (
		item *memcache.Item
		conn = d.subtitleMc.Get(c)
		key  = d.subtitleKey(subtitle.Oid, subtitle.ID)
	)
	defer conn.Close()
	item = &memcache.Item{
		Key:        key,
		Object:     subtitle,
		Flags:      memcache.FlagJSON | memcache.FlagGzip,
		Expiration: d.subtitleMcExpire,
	}
	if err = conn.Set(item); err != nil {
		log.Error("conn.Set(%v) error(%v)", item, err)
	}
	return
}

// DelSubtitleCache .
func (d *Dao) DelSubtitleCache(c context.Context, oid int64, subtitleID int64) (err error) {
	var (
		key  = d.subtitleKey(oid, subtitleID)
		conn = d.subtitleMc.Get(c)
	)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("memcache.Delete(%s) error(%v)", key, err)
		}
	}
	return
}

// SetSubtitleSubjectCache .
func (d *Dao) SetSubtitleSubjectCache(c context.Context, subtitleSubject *model.SubtitleSubject) (err error) {
	var (
		key  = d.subtitleSubjectKey(subtitleSubject.Aid)
		conn = d.subtitleMc.Get(c)
		item *memcache.Item
	)
	defer conn.Close()
	item = &memcache.Item{
		Key:        key,
		Object:     subtitleSubject,
		Flags:      memcache.FlagJSON | memcache.FlagGzip,
		Expiration: d.subtitleMcExpire,
	}
	if err = conn.Set(item); err != nil {
		log.Error("conn.Set(%v) error(%v)", item, err)
	}
	return
}

// SubtitleSubjectCache .
func (d *Dao) SubtitleSubjectCache(c context.Context, aid int64) (subtitleSubject *model.SubtitleSubject, err error) {
	var (
		item *memcache.Item
		conn = d.subtitleMc.Get(c)
		key  = d.subtitleSubjectKey(aid)
	)
	defer conn.Close()
	if item, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("memcache.Get(%s) error(%v)", key, err)
		}
		return
	}
	if err = conn.Scan(item, &subtitleSubject); err != nil {
		log.Error("mc.Scan() error(%v)", err)
	}
	return
}

// DelSubtitleSubjectCache .
func (d *Dao) DelSubtitleSubjectCache(c context.Context, aid int64) (err error) {
	var (
		key  = d.subtitleSubjectKey(aid)
		conn = d.subtitleMc.Get(c)
	)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("memcache.Delete(%s) error(%v)", key, err)
		}
	}
	return
}

// SubtitleWorlFlowTagCache .
func (d *Dao) SubtitleWorlFlowTagCache(c context.Context, bid, rid int64) (data []*model.WorkFlowTag, err error) {
	var (
		item *memcache.Item
		key  = d.subtitleReportTagKey(bid, rid)
		conn = d.subtitleMc.Get(c)
	)
	defer conn.Close()
	if item, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("memcache.Get(%s) error(%v)", key, err)
		}
		return
	}
	if err = conn.Scan(item, &data); err != nil {
		log.Error("mc.Scan() error(%v)", err)
	}
	return
}

// SetSubtitleWorlFlowTagCache .
func (d *Dao) SetSubtitleWorlFlowTagCache(c context.Context, bid, rid int64, data []*model.WorkFlowTag) (err error) {
	var (
		key  = d.subtitleReportTagKey(bid, rid)
		conn = d.subtitleMc.Get(c)
		item *memcache.Item
	)
	defer conn.Close()
	if len(data) == 0 {
		return
	}
	item = &memcache.Item{
		Key:        key,
		Object:     data,
		Flags:      memcache.FlagJSON | memcache.FlagGzip,
		Expiration: d.subtitleMcExpire,
	}
	if err = conn.Set(item); err != nil {
		log.Error("conn.Set(%v) error(%v)", item, err)
	}
	return
}
