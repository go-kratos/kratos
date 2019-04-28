package service

import (
	"context"
	"time"

	"go-common/app/service/main/relation/model"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/stat/prom"
)

const (
	_dailyFollowerNotifyLimit = 2
)

// Followers get follower list.
func (s *Service) Followers(c context.Context, mid int64) (fs []*model.Following, err error) {
	if mid <= 0 {
		return
	}
	var mc = true

	prom.CacheHit.Incr("followers")
	if fs, err = s.dao.FollowerCache(c, mid); err != nil {
		err = nil
		mc = false
	} else if fs != nil {
		return
	}
	prom.CacheMiss.Incr("followers")

	if fs, err = s.dao.Followers(c, mid); err != nil {
		return
	} else if len(fs) == 0 {
		fs = _emptyFollowings
	}
	if mc {
		s.addCache(func() {
			s.dao.SetFollowerCache(context.TODO(), mid, fs)
		})
	}
	return
}

// DelFollower del follower.
func (s *Service) DelFollower(c context.Context, mid, fid int64, src uint8, ric map[string]string) (err error) {
	if mid <= 0 || fid <= 0 {
		return
	}
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
		// s.CompareAndDelCache(c, mid, fid, a)
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
		_, err = s.dao.TxAddStat(c, tx, fid, rn, now)
	}
	_, err = s.dao.TxDelTagUser(c, tx, mid, fid)
	s.RelationInfoc(mid, fid, now.Unix(), ric[RelInfocIP], ric[RelInfocSid], ric[RelInfocBuvid], delFollowerURL, ric[RelInfocReferer], ric[RelInfocUA], src)

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
			"url":        delFollowerURL,
			"referer":    ric[RelInfocReferer],
			"user-agent": ric[RelInfocUA],
		},
	}
	s.dao.DelFollowerLog(c, l)
	return
}

// DelFollowerCache delete follower cache.
func (s *Service) DelFollowerCache(c context.Context, mid int64) (err error) {
	err = s.dao.DelFollowerCache(c, mid)
	return
}

// Unread is
// 展示最近有未知晓的新粉丝
func (s *Service) Unread(c context.Context, fid int64) (bool, error) {
	if shouldNotified := s.shouldNotified(c, fid); !shouldNotified {
		return false, nil
	}

	flag, err := s.dao.RctFollowerNotify(c, fid)
	if err != nil {
		return false, err
	}
	count, err := s.dao.RctFollowerCount(c, fid)
	if err != nil {
		return false, err
	}
	notify := false
	if flag && count > 0 {
		notify = true
	}
	s.addCache(func() {
		s.dao.SetRctFollowerNotify(context.Background(), fid, false)
		// 真的看见红点了再加一
		if notify {
			s.dao.IncrTodayNotifyCount(context.Background(), fid)
		}
	})
	return notify, nil
}

// UnreadCount unread count.
func (s *Service) UnreadCount(c context.Context, fid int64) (int64, error) {
	count, err := s.dao.RctFollowerCount(c, fid)
	if err != nil {
		return 0, err
	}
	if count < 0 {
		count = 0
	}
	s.addCache(func() {
		s.ResetUnreadCount(context.Background(), fid)
	})
	return count, nil
}

// ResetUnread is
// 重置未知晓的新粉丝
func (s *Service) ResetUnread(ctx context.Context, fid int64) error {
	if err := s.dao.SetRctFollowerNotify(ctx, fid, false); err != nil {
		log.Error("Failed to reset recent follower notify with fid: %d: %+v", fid, err)
	}
	if err := s.dao.IncrTodayNotifyCount(ctx, fid); err != nil {
		log.Error("Failed to incr today notify count with fid: %d: %+v", fid, err)
	}
	return nil
}

// ResetUnreadCount is
func (s *Service) ResetUnreadCount(ctx context.Context, fid int64) error {
	if err := s.dao.EmptyRctFollower(ctx, fid); err != nil {
		log.Error("Failed to empty recent follower with fid: %d: %+v", fid, err)
	}
	if err := s.dao.SetRctFollowerNotify(ctx, fid, false); err != nil {
		log.Error("Failed to reset recent follower with fid: %d: %+v", fid, err)
	}
	return nil
}

// FollowerNotifySetting get follower-notify setting
func (s *Service) FollowerNotifySetting(c context.Context, arg *model.ArgMid) (followerNotify *model.FollowerNotifySetting, err error) {
	var (
		enabled  bool
		addCache = true
	)
	followerNotify, err = s.dao.GetFollowerNotifyCache(c, arg.Mid)
	if err != nil {
		addCache = false
		err = nil
	}
	if followerNotify != nil {
		prom.CacheHit.Incr("FollowerNotify")
		return
	}
	prom.CacheMiss.Incr("FollowerNotify")

	if enabled, err = s.dao.FollowerNotifySetting(c, arg.Mid); err != nil {
		return
	}
	followerNotify = &model.FollowerNotifySetting{
		Mid:     arg.Mid,
		Enabled: enabled,
	}

	if !addCache {
		return
	}
	// 异步的更新缓存
	s.addCache(func() {
		s.dao.SetFollowerNotifyCache(context.TODO(), arg.Mid, followerNotify)
	})
	return
}

// DisableFollowerNotify disable follower-notify setting
func (s *Service) DisableFollowerNotify(c context.Context, arg *model.ArgMid) (err error) {
	if _, err = s.dao.DisableFollowerNotify(c, arg.Mid); err != nil {
		return
	}
	s.dao.DelFollowerNotifyCache(context.TODO(), arg.Mid)
	return
}

// EnableFollowerNotify enable follower-notify setting
func (s *Service) EnableFollowerNotify(c context.Context, arg *model.ArgMid) (err error) {
	if _, err = s.dao.EnableFollowerNotify(c, arg.Mid); err != nil {
		return
	}
	s.dao.DelFollowerNotifyCache(context.TODO(), arg.Mid)
	return
}

func (s *Service) shouldNotified(c context.Context, mid int64) (shouldNotified bool) {
	var (
		notifyCount    int64
		err            error
		followerNotify *model.FollowerNotifySetting
	)
	// 得到用户新粉丝消息提醒设置，如果被禁用，则不再提醒
	if followerNotify, err = s.FollowerNotifySetting(c, &model.ArgMid{Mid: mid}); err != nil {
		log.Error("Failed to get follower notify setting: fid: %d: %+v", mid, err)
		return true
	}
	if !followerNotify.Enabled {
		return false
	}
	// 得到当日新粉丝提醒数量, 如果大于等于限制值，则不再提醒
	if notifyCount, err = s.dao.TodayNotifyCountCache(c, mid); err != nil {
		log.Error("Failed to get notify count: fid: %d: %+v", mid, err)
		return true
	}
	if notifyCount >= _dailyFollowerNotifyLimit {
		return false
	}
	return true
}
