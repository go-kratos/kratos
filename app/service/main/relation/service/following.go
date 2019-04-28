package service

import (
	"context"
	"sort"
	"time"

	"go-common/app/service/main/relation/model"
	"go-common/app/service/main/relation/model/sets"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/stat/prom"
)

// Followings get user's following list.
func (s *Service) Followings(c context.Context, mid int64) (fs []*model.Following, err error) {
	if fs, err = s.followings(c, mid); err != nil {
		return
	}
	fs = model.Filter(fs, model.AttrFollowing)
	sort.Sort(model.SortFollowings(fs))
	return
}

// Attentions return user's following and whisper list.
func (s *Service) Attentions(c context.Context, mid int64) (fs []*model.Following, err error) {
	if fs, err = s.followings(c, mid); err != nil {
		return
	}
	fs = model.Filter(fs, model.AttrFollowing|model.AttrWhisper)
	sort.Sort(model.SortFollowings(fs))
	return
}

func (s *Service) followings(c context.Context, mid int64) (fs []*model.Following, err error) {
	if mid <= 0 {
		return
	}
	var mc, redis = true, true

	prom.CacheHit.Incr("followings_mc")
	if fs, err = s.dao.FollowingCache(c, mid); err != nil {
		mc = false
		err = nil
	} else if fs != nil {
		return
	}
	prom.CacheMiss.Incr("followings_mc")

	prom.CacheHit.Incr("followings_redis")
	if fs, err = s.dao.FollowingsCache(c, mid); err != nil {
		redis = false
		err = nil
	} else if len(fs) > 0 {
		if mc {
			s.addCache(func() {
				s.dao.SetFollowingCache(context.TODO(), mid, fs)
			})
		}
		return
	}
	prom.CacheMiss.Incr("followings_redis")

	if fs, err = s.dao.Followings(c, mid); err != nil {
		return
	} else if len(fs) == 0 {
		fs = _emptyFollowings
	} else {
		var (
			ts map[int64][]int64
		)
		if ts, err = s.dao.UserTag(c, mid); err != nil {
			return
		}
		for _, f := range fs {
			if tags, ok := ts[f.Mid]; ok {
				f.Tag = tags
				for _, id := range f.Tag {
					if id == -10 {
						f.Special = 1
						break
					}
				}
			}
		}
	}
	if redis || mc {
		s.addCache(func() {
			if redis {
				s.dao.SetFollowingsCache(context.TODO(), mid, fs)
			}
			if mc {
				s.dao.SetFollowingCache(context.TODO(), mid, fs)
			}
		})
	}
	return
}

// AddFollowing add following.
func (s *Service) AddFollowing(c context.Context, mid, fid int64, src uint8, ric map[string]string) (err error) {
	var (
		a, ra   uint32
		na, nra uint32
		at      uint32
		st      *model.Stat
		tx      *sql.Tx
		friend  = false
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
				log.Error("tx.Rollback() error(%v)", err1)
			}
			return
		}
		if err = tx.Commit(); err != nil {
			log.Error("tx.Commit() error(%v)", err)
		}
	}()
	if st, _, err = s.txStat(c, tx, mid, fid); err != nil {
		return
	} else if st == nil {
		log.Error("s.txStat(%d %d) error(%v)", mid, fid, err)
		err = ecode.ServerErr
		return
	}
	if a, err = s.dao.Relation(c, mid, fid); err != nil {
		return
	}
	if ra, err = s.dao.Relation(c, fid, mid); err != nil {
		return
	}
	na = model.AttrFollowing
	nra = ra
	switch model.Attr(ra) {
	case model.AttrFollowing:
		friend = true
		na = model.SetAttr(na, model.AttrFriend)
		nra = model.SetAttr(ra, model.AttrFriend)
	}
	at = model.Attr(a)
	switch at {
	case model.AttrBlack:
		err = ecode.RelFollowAlreadyBlack
		return
	case model.AttrFriend, model.AttrFollowing:
		return
	case model.AttrWhisper:
		if friend {
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
		} else {
			if _, err = s.dao.TxSetFollowing(c, tx, mid, fid, na, src, model.StatusOK, now); err != nil {
				return
			}
			if _, err = s.dao.TxSetFollower(c, tx, mid, fid, na, src, model.StatusOK, now); err != nil {
				return
			}
		}
		n.Whisper = -1
		n.Following = 1
	case model.AttrNoRelation:
		if st.Count() >= s.c.Relation.MaxFollowingLimit {
			err = ecode.RelFollowReachMaxLimit
			return
		}
		if friend {
			if _, err = s.dao.TxAddFollowing(c, tx, mid, fid, na, src, now); err != nil {
				return
			}
			if _, err = s.dao.TxSetFollowing(c, tx, fid, mid, nra, src, model.StatusOK, now); err != nil {
				return
			}
			if _, err = s.dao.TxAddFollower(c, tx, mid, fid, na, src, now); err != nil {
				return
			}
			if _, err = s.dao.TxSetFollower(c, tx, fid, mid, nra, src, model.StatusOK, now); err != nil {
				return
			}
		} else {
			if _, err = s.dao.TxAddFollowing(c, tx, mid, fid, na, src, now); err != nil {
				return
			}
			if _, err = s.dao.TxAddFollower(c, tx, mid, fid, na, src, now); err != nil {
				return
			}
		}
		n.Following = 1
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
	s.RelationInfoc(mid, fid, now.Unix(), realIP, ric[RelInfocSid], ric[RelInfocBuvid], addFollowingURL, ric[RelInfocReferer], ric[RelInfocUA], src)

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
			"url":        addFollowingURL,
			"referer":    ric[RelInfocReferer],
			"user-agent": ric[RelInfocUA],
		},
	}
	s.dao.AddFollowingLog(c, l)
	// 后续逻辑
	s.addCache(func() {
		s.onAddFollowing(context.Background(), mid, fid)
	})
	return
}

// DelFollowing del following.
func (s *Service) DelFollowing(c context.Context, mid, fid int64, src uint8, ric map[string]string) (err error) {
	var (
		a           uint32
		ra, na, nra uint32
		tx          *sql.Tx
		friend      = false
		n           = new(model.Stat)
		rn          = new(model.Stat)
		now         = time.Now()
		realIP      = ric[RelInfocIP]
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
	switch model.Attr(a) {
	case model.AttrFriend:
		friend = true
	case model.AttrFollowing:
	default:
		s.CompareAndDelCache(c, mid, fid, a)
		// err = ecode.RelFollowAttrNotSet
		log.Warn("Invalid state between %d and %d with attribute: %d", mid, fid, a)
		return
	}
	na = model.AttrNoRelation
	if friend {
		nra = model.AttrFollowing
		if _, err = s.dao.TxSetFollowing(c, tx, mid, fid, na, src, model.StatusDel, now); err != nil {
			return
		}
		if _, err = s.dao.TxSetFollowing(c, tx, fid, mid, nra, src, model.StatusOK, now); err != nil {
			return
		}
		if _, err = s.dao.TxSetFollower(c, tx, mid, fid, na, src, model.StatusDel, now); err != nil {
			return
		}
		if _, err = s.dao.TxSetFollower(c, tx, fid, mid, nra, src, model.StatusOK, now); err != nil {
			return
		}
	} else {
		if _, err = s.dao.TxSetFollowing(c, tx, mid, fid, na, src, model.StatusDel, now); err != nil {
			return
		}
		if _, err = s.dao.TxSetFollower(c, tx, mid, fid, na, src, model.StatusDel, now); err != nil {
			return
		}
	}
	n.Following = -1
	rn.Follower = -1
	if !n.Empty() {
		if _, err = s.dao.TxAddStat(c, tx, mid, n, now); err != nil {
			return
		}
	}
	if !rn.Empty() {
		if _, err = s.dao.TxAddStat(c, tx, fid, rn, now); err != nil {
			return
		}
	}
	_, err = s.dao.TxDelTagUser(c, tx, mid, fid)
	s.RelationInfoc(mid, fid, now.Unix(), ric[RelInfocIP], ric[RelInfocSid], ric[RelInfocBuvid], delFollowingURL, ric[RelInfocReferer], ric[RelInfocUA], src)

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
			"url":        addFollowingURL,
			"referer":    ric[RelInfocReferer],
			"user-agent": ric[RelInfocUA],
		},
	}
	s.dao.DelFollowingLog(c, l)
	// 后续逻辑
	s.addCache(func() {
		s.onDelFollowing(context.Background(), mid, fid)
	})
	return
}

// DelFollowingCache delete following cache.
func (s *Service) DelFollowingCache(c context.Context, mid int64) (err error) {
	if err = s.dao.DelFollowingsCache(c, mid); err != nil {
		return
	}
	err = s.dao.DelFollowingCache(c, mid)
	return
}

// UpdateFollowingCache update following cache.
func (s *Service) UpdateFollowingCache(c context.Context, mid int64, following *model.Following) (err error) {
	if following.Attribute == model.AttrNoRelation {
		err = s.dao.DelFollowing(c, mid, following)
	} else {
		err = s.dao.AddFollowingCache(c, mid, following)
	}
	if err != nil {
		return
	}
	err = s.dao.DelFollowingCache(c, mid)
	return
}

func (s *Service) onAddFollowing(ctx context.Context, mid, fid int64) {
	// 最近有新粉丝通知逻辑
	if err := s.dao.AddRctFollower(ctx, mid, fid); err != nil {
		log.Error("Failed to add recent follower: mid: %d fid: %d: %+v", mid, fid, err)
		return
	}
	if err := s.dao.SetRctFollowerNotify(ctx, fid, true); err != nil {
		log.Error("Failed to set recent follower notify: fid: %d flag: true: %+v", mid, err)
		return
	}
}

func (s *Service) onDelFollowing(ctx context.Context, mid, fid int64) {
	// 最近有新粉丝通知逻辑
	if err := s.dao.DelRctFollower(ctx, mid, fid); err != nil {
		log.Error("Failed to del recent follower: mid: %d fid: %d: %+v", mid, fid, err)
		return
	}
	count, err := s.dao.RctFollowerCount(ctx, fid)
	if err != nil {
		log.Error("Failed to get recent follower count: fid: %d: %+v", fid, err)
		return
	}
	if count > 0 {
		return
	}
	if err := s.dao.SetRctFollowerNotify(ctx, fid, false); err != nil {
		log.Error("Failed to set recent follower notify: fid: %d flag: false: %+v", fid, err)
		return
	}
}

// CompareAndDelCache is
func (s *Service) CompareAndDelCache(ctx context.Context, mid, fid int64, inAttr uint32) {
	fr, err := s.Relation(ctx, mid, fid)
	if err != nil {
		return
	}

	cacheAttr := model.AttrNoRelation
	if fr != nil {
		cacheAttr = model.Attr(fr.Attribute)
	}
	refAttr := model.Attr(inAttr)
	if cacheAttr == refAttr {
		return
	}
	log.Warn("The relation attribute is inconsistent: mid: %d, fid: %d, reference: %d, cache: %d", mid, fid, refAttr, cacheAttr)
	s.DelFollowingCache(ctx, mid)
	s.DelFollowerCache(ctx, mid)
}

// SameFollowings get users' same following list.
func (s *Service) SameFollowings(c context.Context, arg *model.ArgSameFollowing) ([]*model.Following, error) {
	flw1, err := s.Followings(c, arg.Mid1)
	if err != nil {
		return nil, err
	}
	flw2, err := s.Followings(c, arg.Mid2)
	if err != nil {
		return nil, err
	}

	smids1 := sets.NewInt64(asMids(flw1)...)
	smids2 := sets.NewInt64(asMids(flw2)...)
	sameMids := smids1.Intersection(smids2)

	// zhuangsusu: 以第一个用户的关注时间的倒序来展示给第二个用户看
	flw1Map := asFollowingMap(flw1)
	flw2Map := asFollowingMap(flw2)
	sorted := make([]*model.Following, 0, sameMids.Len())
	for mid := range sameMids {
		sorted = append(sorted, flw1Map[mid])
	}
	sort.Sort(model.SortFollowings(sorted))
	result := make([]*model.Following, 0, sameMids.Len())
	for _, f := range sorted {
		result = append(result, flw2Map[f.Mid])
	}
	return result, nil
}

func asMids(fs []*model.Following) []int64 {
	mids := make([]int64, 0, len(fs))
	for _, f := range fs {
		mids = append(mids, f.Mid)
	}
	return mids
}

func asFollowingMap(fs []*model.Following) map[int64]*model.Following {
	flwMap := make(map[int64]*model.Following, len(fs))
	for _, f := range fs {
		flwMap[f.Mid] = f
	}
	return flwMap
}
