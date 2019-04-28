package service

import (
	"context"
	"time"

	"go-common/app/service/main/thumbup/dao"
	"go-common/app/service/main/thumbup/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// Like like
func (s *Service) Like(c context.Context, business string, mid, originID, messageID int64, likeType int8, upMid int64) (stat model.Stats, err error) {
	var businessID int64
	if businessID, err = s.CheckBusinessOrigin(business, originID); err != nil {
		return
	}
	if (likeType < 0) || (likeType > 4) {
		err = ecode.RequestErr
		return
	}
	oldState, err := s.dao.LikeState(c, mid, businessID, originID, messageID)
	if err == nil {
		if oldState == model.StateBlank {
			switch likeType {
			case model.TypeCancelLike:
				err = ecode.ThumbupCancelLikeErr
			case model.TypeCancelDislike:
				err = ecode.ThumbupCancelDislikeErr
			}
		} else if oldState == model.StateLike {
			switch likeType {
			case model.TypeLike:
				err = ecode.ThumbupDupLikeErr
			case model.TypeCancelDislike:
				err = ecode.ThumbupCancelDislikeErr
			}
		} else if oldState == model.StateDislike {
			switch likeType {
			case model.TypeCancelLike:
				err = ecode.ThumbupCancelLikeErr
			case model.TypeDislike:
				err = ecode.ThumbupDupDislikeErr
			}
		} else {
			dao.PromError("点赞表状态异常")
			log.Error("service.Like(mid:%v bid:%v oid: %v, messageID: %v, state: %v)", mid, businessID, originID, messageID, oldState)
			return
		}
		if err != nil {
			return
		}
	}
	err = s.dbus.Do(c, func(c context.Context) {
		msg := &model.LikeMsg{UpMid: upMid, Mid: mid, Type: likeType, LikeTime: time.Now(), Business: business, OriginID: originID, MessageID: messageID}
		s.dao.PubLikeDatabus(c, msg)
	})
	if err != nil {
		return
	}
	// todo fix
	if business == "reply" {
		return
	}
	stats, err := s.Stats(c, business, originID, []int64{messageID})
	if err == nil && stats != nil && stats[messageID] != nil {
		stat = *stats[messageID]
	}
	stat = calculateCount(stat, likeType)
	return
}

func calculateCount(stat model.Stats, typ int8) model.Stats {
	var likesCount, dislikesCount int64
	switch typ {
	case model.TypeLike:
		likesCount = 1
	case model.TypeCancelLike:
		likesCount = -1
	case model.TypeDislike:
		dislikesCount = 1
	case model.TypeCancelDislike:
		dislikesCount = -1
	case model.TypeLikeReverse:
		likesCount = -1
		dislikesCount = 1
	case model.TypeDislikeReverse:
		likesCount = 1
		dislikesCount = -1
	}
	stat.Likes += likesCount
	stat.Dislikes += dislikesCount
	return stat
}

func (s *Service) updateStatCache(c context.Context, businessID, originID int64, stat *model.Stats) (err error) {
	if stat == nil {
		return
	}
	var ok bool
	if originID == 0 {
		err = s.dao.AddStatsCache(c, businessID, stat)
	} else {
		if ok, err = s.dao.ExpireHashStatsCache(c, businessID, originID); ok {
			err = s.dao.AddHashStatsCache(c, businessID, originID, stat)
		}
	}
	return
}
