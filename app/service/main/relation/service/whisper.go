package service

import (
	"context"
	"sort"
	"time"

	"go-common/app/service/main/relation/model"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
)

// Whispers get user's whisper list.
func (s *Service) Whispers(c context.Context, mid int64) (fs []*model.Following, err error) {
	if fs, err = s.followings(c, mid); err != nil {
		return
	}
	fs = model.Filter(fs, model.AttrWhisper)
	sort.Sort(model.SortFollowings(fs))
	return
}

// AddWhisper add whisper.
func (s *Service) AddWhisper(c context.Context, mid, fid int64, src uint8, ric map[string]string) (err error) {
	var (
		a, ra   uint32
		na, nra uint32
		at      uint32
		st      *model.Stat
		tx      *sql.Tx
		n       = new(model.Stat)
		rn      = new(model.Stat)
		now     = time.Now()
		audit   *model.Audit
		monitor bool
	)
	if mid <= 0 || fid <= 0 {
		return
	}
	if mid == fid {
		err = ecode.RelFollowSelfBanned
		return
	}
	if monitor, err = s.Monitor(c, fid); err != nil {
		return
	} else if monitor {
		return
	}
	realIP := ric[RelInfocIP]
	if audit, err = s.Audit(c, mid, realIP); err != nil {
		log.Error("s.Audit.mid(%d) error(%v) return(%v)", mid, err, audit)
		return
	}
	if audit.Blocked {
		err = ecode.UserDisabled
		return
	}
	if audit.Rank == UserRank && !audit.BindTel {
		err = ecode.RelFollowReachTelLimit
		return
	}

	if err = s.initStat(c, mid, fid); err != nil {
		return
	}
	if tx, err = s.dao.BeginTran(c); err != nil {
		return
	}
	defer func() {
		if err != nil {
			if err1 := tx.Rollback(); err1 != nil {
				log.Error("tx.Rollback() error(%v)", err)
			}
			return
		}
		if err = tx.Commit(); err != nil {
			log.Error("tx.Commit() error(%v)", err)
		}
	}()
	if _, _, err = s.txStat(c, tx, mid, fid); err != nil {
		return
	}
	if a, err = s.dao.Relation(c, mid, fid); err != nil {
		return
	}
	if ra, err = s.dao.Relation(c, fid, mid); err != nil {
		return
	}
	at = model.Attr(a)
	switch at {
	case model.AttrBlack:
		err = ecode.RelFollowAlreadyBlack
		return
	case model.AttrWhisper:
		return
	}
	na = model.AttrWhisper
	switch at {
	case model.AttrFriend:
		nra = model.UnsetAttr(ra, model.AttrFriend)
		if _, err = s.dao.TxSetFollowing(c, tx, mid, fid, na, src, model.StatusOK, now); err != nil {
			return
		}
		if _, err = s.dao.TxSetFollowing(c, tx, fid, mid, nra, src, model.StatusOK, now); err != nil {
			return
		}
		if _, err = s.dao.TxSetFollower(c, tx, mid, fid, na, src, model.StatusOK, now); err != nil {
			return
		}
		if _, err = s.dao.TxSetFollower(c, tx, fid, mid, nra, src, model.StatusOK, now); err != nil {
			return
		}
		if _, err = s.dao.TxDelTagUser(c, tx, mid, fid); err != nil {
			return
		}
		n.Following = -1
		n.Whisper = 1
	case model.AttrFollowing:
		if _, err = s.dao.TxSetFollowing(c, tx, mid, fid, na, src, model.StatusOK, now); err != nil {
			return
		}
		if _, err = s.dao.TxSetFollower(c, tx, mid, fid, na, src, model.StatusOK, now); err != nil {
			return
		}
		if _, err = s.dao.TxDelTagUser(c, tx, mid, fid); err != nil {
			return
		}
		n.Following = -1
		n.Whisper = 1
	case model.AttrNoRelation:
		if st, err = s.dao.TxStat(c, tx, mid); err != nil {
			return
		} else if st != nil && st.Count() >= s.c.Relation.MaxFollowingLimit {
			err = ecode.RelFollowReachMaxLimit
			return
		}
		if _, err = s.dao.TxAddFollowing(c, tx, mid, fid, na, src, now); err != nil {
			return
		}
		if _, err = s.dao.TxAddFollower(c, tx, mid, fid, na, src, now); err != nil {
			return
		}
		n.Whisper = 1
		rn.Follower = 1
	}
	if !n.Empty() {
		if _, err = s.dao.TxAddStat(c, tx, mid, n, now); err != nil {
			return
		}
	}
	if !rn.Empty() {
		_, err = s.dao.TxAddStat(c, tx, fid, rn, now)
	}
	s.RelationInfoc(mid, fid, now.Unix(), realIP, ric[RelInfocSid], ric[RelInfocBuvid], addWhisperURL, ric[RelInfocReferer], ric[RelInfocUA], src)

	// log to report
	l := &model.RelationLog{
		Mid:         mid,
		Fid:         fid,
		Ts:          now.Unix(),
		Ip:          realIP,
		Source:      uint32(src),
		FromAttr:    a,
		ToAttr:      na,
		FromRevAttr: ra,
		ToRevAttr:   nra,
		Content: map[string]string{
			"sid":        ric[RelInfocSid],
			"buvid":      ric[RelInfocBuvid],
			"url":        addWhisperURL,
			"referer":    ric[RelInfocReferer],
			"user-agent": ric[RelInfocUA],
		},
	}
	s.dao.AddWhisperLog(c, l)
	return
}

// DelWhisper del whisper.
func (s *Service) DelWhisper(c context.Context, mid, fid int64, src uint8, ric map[string]string) (err error) {
	var (
		a      uint32
		tx     *sql.Tx
		n      = new(model.Stat)
		rn     = new(model.Stat)
		now    = time.Now()
		realIP = ric[RelInfocIP]
	)
	if mid <= 0 || fid <= 0 {
		return
	}
	if mid == fid {
		err = ecode.RelFollowSelfBanned
		return
	}
	if err = s.initStat(c, mid, fid); err != nil {
		return
	}
	if tx, err = s.dao.BeginTran(c); err != nil {
		return
	}
	defer func() {
		if err != nil {
			if err1 := tx.Rollback(); err1 != nil {
				log.Error("tx.Rollback() error(%v)", err1)
			}
			return
		}
		if err = tx.Commit(); err != nil {
			log.Error("tx.Commit() error(%v)", err)
		}
	}()
	if _, _, err = s.txStat(c, tx, mid, fid); err != nil {
		return
	}
	if a, err = s.dao.Relation(c, mid, fid); err != nil {
		return
	}
	if model.AttrWhisper != model.Attr(a) {
		s.CompareAndDelCache(c, mid, fid, a)
		// err = ecode.RelFollowAttrNotSet
		log.Warn("Invalid state between %d and %d with attribute: %d", mid, fid, a)
		return
	}
	if _, err = s.dao.TxSetFollowing(c, tx, mid, fid, model.AttrNoRelation, src, model.StatusDel, now); err != nil {
		return
	}
	if _, err = s.dao.TxSetFollower(c, tx, mid, fid, model.AttrNoRelation, src, model.StatusDel, now); err != nil {
		return
	}
	n.Whisper = -1
	rn.Follower = -1
	if !n.Empty() {
		if _, err = s.dao.TxAddStat(c, tx, mid, n, now); err != nil {
			return
		}
	}
	if !rn.Empty() {
		_, err = s.dao.TxAddStat(c, tx, fid, rn, now)
	}
	s.RelationInfoc(mid, fid, now.Unix(), ric[RelInfocIP], ric[RelInfocSid], ric[RelInfocBuvid], delWhisperURL, ric[RelInfocReferer], ric[RelInfocUA], src)

	// log to report
	l := &model.RelationLog{
		Mid:         mid,
		Fid:         fid,
		Ts:          now.Unix(),
		Ip:          realIP,
		Source:      uint32(src),
		FromAttr:    a,
		ToAttr:      model.AttrNoRelation,
		FromRevAttr: 0, // no means on whisper
		ToRevAttr:   0, // no means on whisper
		Content: map[string]string{
			"sid":        ric[RelInfocSid],
			"buvid":      ric[RelInfocBuvid],
			"url":        delWhisperURL,
			"referer":    ric[RelInfocReferer],
			"user-agent": ric[RelInfocUA],
		},
	}
	s.dao.DelWhisperLog(c, l)
	return
}
