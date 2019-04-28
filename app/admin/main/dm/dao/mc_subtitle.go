package dao

import (
	"context"
	"fmt"

	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_fmtSubtitle        = "s_subtitle_%d_%d"
	_fmtVideoSubtitle   = "s_video_%d_%d"
	_fmtSubtitleDraft   = "s_draft_%v_%v_%v_%v"
	_fmtSubtitleSubject = "s_subtitle_allow_%d"
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

// DelVideoSubtitleCache .
func (d *Dao) DelVideoSubtitleCache(c context.Context, oid int64, tp int32) (err error) {
	var (
		key  = d.subtitleVideoKey(oid, tp)
		conn = d.subtitleMC.Get(c)
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

// DelSubtitleDraftCache .
func (d *Dao) DelSubtitleDraftCache(c context.Context, oid int64, tp int32, mid int64, lan uint8) (err error) {
	var (
		key  = d.subtitleDraftKey(oid, tp, mid, lan)
		conn = d.subtitleMC.Get(c)
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

// DelSubtitleCache .
func (d *Dao) DelSubtitleCache(c context.Context, oid int64, subtitleID int64) (err error) {
	var (
		key  = d.subtitleKey(oid, subtitleID)
		conn = d.subtitleMC.Get(c)
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

// DelSubtitleSubjectCache .
func (d *Dao) DelSubtitleSubjectCache(c context.Context, aid int64) (err error) {
	var (
		key  = d.subtitleSubjectKey(aid)
		conn = d.subtitleMC.Get(c)
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
