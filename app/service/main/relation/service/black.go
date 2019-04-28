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

// Blacks get black list.
func (s *Service) Blacks(c context.Context, mid int64) (fs []*model.Following, err error) {
	if mid <= 0 {
		return
	}
	stat, err := s.Stat(c, mid)
	if err != nil || stat.Black == 0 {
		return
	}
	if fs, err = s.followings(c, mid); err != nil {
		return
	}
	fs = model.Filter(fs, model.AttrBlack)
	sort.Sort(model.SortFollowings(fs))
	return
}

// AddBlack add black.
func (s *Service) AddBlack(c context.Context, mid int64, fid int64, src uint8, ric map[string]string) (err error) {
	var (
		a, ra   uint32
		na, nra uint32
		at      uint32
		st      *model.Stat
		tx      *sql.Tx
		n       = new(model.Stat)
		rn      = new(model.Stat)
		now     = time.Now()
		realIP  = ric[RelInfocIP]
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
	if ra, err = s.dao.Relation(c, fid, mid); err != nil {
		return
	}
	at = model.Attr(a)
	if st, err = s.dao.TxStat(c, tx, mid); err != nil {
		return
	} else if st != nil && st.BlackCount() >= s.c.Relation.MaxBlackLimit {
		err = ecode.RelBlackReachMaxLimit
		return
	}
	switch at {
	case model.AttrBlack:
		return
	case model.AttrNoRelation:
		if _, err = s.dao.TxAddFollowing(c, tx, mid, fid, model.AttrBlack, src, now); err != nil {
			return
		}
		if _, err = s.dao.TxAddFollower(c, tx, mid, fid, model.AttrBlack, src, now); err != nil {
			return
		}
	case model.AttrFriend:
		nra = model.UnsetAttr(ra, model.AttrFriend)
		na = model.AttrBlack
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
	case model.AttrFollowing:
		if _, err = s.dao.TxAddFollowing(c, tx, mid, fid, model.AttrBlack, src, now); err != nil {
			return
		}
		if _, err = s.dao.TxAddFollower(c, tx, mid, fid, model.AttrBlack, src, now); err != nil {
			return
		}
		if _, err = s.dao.TxDelTagUser(c, tx, mid, fid); err != nil {
			return
		}
	case model.AttrWhisper:
		if _, err = s.dao.TxAddFollowing(c, tx, mid, fid, model.AttrBlack, src, now); err != nil {
			return
		}
		if _, err = s.dao.TxAddFollower(c, tx, mid, fid, model.AttrBlack, src, now); err != nil {
			return
		}
	}
	n.Black = 1
	switch at {
	case model.AttrFriend, model.AttrFollowing:
		n.Following = -1
		rn.Follower = -1
	case model.AttrWhisper:
		n.Whisper = -1
		rn.Follower = -1
	}
	if !n.Empty() {
		if _, err = s.dao.TxAddStat(c, tx, mid, n, now); err != nil {
			return
		}
	}
	if !rn.Empty() {
		_, err = s.dao.TxAddStat(c, tx, fid, rn, now)
	}
	s.RelationInfoc(mid, fid, now.Unix(), ric[RelInfocIP], ric[RelInfocSid], ric[RelInfocBuvid], addBlackURL, ric[RelInfocReferer], ric[RelInfocUA], src)

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
			"url":        addBlackURL,
			"referer":    ric[RelInfocReferer],
			"user-agent": ric[RelInfocUA],
		},
	}
	s.dao.AddBlackLog(c, l)
	return
}

// DelBlack del black.
func (s *Service) DelBlack(c context.Context, mid, fid int64, src uint8, ric map[string]string) (err error) {
	var (
		a, ra    uint32
		na, nra  uint32
		nat, rat uint32
		tx       *sql.Tx
		friend   = false
		n        = new(model.Stat)
		rn       = new(model.Stat)
		status   = model.StatusOK
		now      = time.Now()
		realIP   = ric[RelInfocIP]
	)
	// if mid == fid {
	// 	err = ecode.RelFollowSelfBanned
	// 	return
	// }
	if mid <= 0 || fid <= 0 {
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
	if ra, err = s.dao.Relation(c, fid, mid); err != nil {
		return
	}

	if model.AttrBlack != model.Attr(a) {
		s.CompareAndDelCache(c, mid, fid, a)
		// err = ecode.RelFollowAttrNotSet
		log.Warn("Invalid state between %d and %d with attribute: %d", mid, fid, a)
		return
	}
	n.Black = -1
	na = model.AttrNoRelation
	nra = ra
	// no recover relation
	//na = model.UnsetAttr(a, model.AttrBlack)
	//nat = model.Attr(na)
	// switch nat {
	// case model.AttrFollowing:
	// 	n.Following = 1
	// 	rn.Follower = 1
	// case model.AttrWhisper:
	// 	n.Whisper = 1
	// 	rn.Follower = 1
	// }
	rat = model.Attr(ra)
	switch rat {
	case model.AttrBlack:
		na = model.AttrNoRelation
		status = model.StatusDel
		n.Following = 0
		n.Whisper = 0
		rn.Follower = 0
	case model.AttrFollowing:
		if nat == model.AttrFollowing {
			na = model.SetAttr(na, model.AttrFriend)
			nra = model.SetAttr(nra, model.AttrFriend)
			friend = true
		}
	}
	if friend {
		if _, err = s.dao.TxSetFollowing(c, tx, mid, fid, na, src, status, now); err != nil {
			return
		}
		if _, err = s.dao.TxSetFollowing(c, tx, fid, mid, nra, src, status, now); err != nil {
			return
		}
		if _, err = s.dao.TxSetFollower(c, tx, mid, fid, na, src, status, now); err != nil {
			return
		}
		if _, err = s.dao.TxSetFollower(c, tx, fid, mid, nra, src, status, now); err != nil {
			return
		}
	} else {
		if _, err = s.dao.TxSetFollowing(c, tx, mid, fid, na, src, status, now); err != nil {
			return
		}
		if _, err = s.dao.TxSetFollower(c, tx, mid, fid, na, src, status, now); err != nil {
			return
		}
	}
	if !n.Empty() {
		if _, err = s.dao.TxAddStat(c, tx, mid, n, now); err != nil {
			return
		}
	}
	if !rn.Empty() {
		_, err = s.dao.TxAddStat(c, tx, fid, rn, now)
	}
	s.RelationInfoc(mid, fid, now.Unix(), ric[RelInfocIP], ric[RelInfocSid], ric[RelInfocBuvid], delBlackURL, ric[RelInfocReferer], ric[RelInfocUA], src)

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
			"url":        delBlackURL,
			"referer":    ric[RelInfocReferer],
			"user-agent": ric[RelInfocUA],
		},
	}
	s.dao.DelBlackLog(c, l)
	return
}
