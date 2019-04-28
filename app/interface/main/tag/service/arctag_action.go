package service

import (
	"context"
	"time"

	"go-common/app/interface/main/tag/model"
	account "go-common/app/service/main/account/model"
	rpcModel "go-common/app/service/main/tag/model"
	"go-common/library/ecode"
)

// Like user action like tag .
func (s *Service) Like(c context.Context, mid, aid, tid int64, now time.Time) (err error) {
	if _, ok := s.channelMap[tid]; ok {
		return ecode.ChannelNoLike
	}
	var card *account.Card
	if card, err = s.dao.UserCard(c, mid); err != nil {
		return
	}
	if card.Level < int32(s.c.Tag.ArcTagLikeLevel) {
		return ecode.TagArcLikeLevelLower
	}
	if card.Silence != model.UserBannedNone {
		return ecode.TagArcAccountBlocked
	}
	count, err := s.dao.SpamCache(c, mid, model.SpamLike)
	if err != nil {
		return
	}
	if count >= s.c.Tag.ArcTagLikeMaxNum {
		return ecode.TagArcTagLikeMaxFre
	}
	if _, err = s.normalArchive(c, aid); err != nil {
		return
	}
	s.dao.IncrSpamCache(c, mid, model.SpamLike)
	return s.likeService(c, mid, tid, aid, rpcModel.ResTypeArchive)
}

// Hate user action Hate2 tag .
func (s *Service) Hate(c context.Context, mid, aid, tid int64, now time.Time) (err error) {
	if _, ok := s.channelMap[tid]; ok {
		return ecode.ChannelNoHate
	}
	var card *account.Card
	if card, err = s.dao.UserCard(c, mid); err != nil {
		return
	}
	if card.Level < int32(s.c.Tag.ArcTagHateLevel) {
		return ecode.TagArcHateLevelLower
	}
	if card.Silence != model.UserBannedNone {
		return ecode.TagArcAccountBlocked
	}
	count, err := s.dao.SpamCache(c, mid, model.SpamHate)
	if err != nil {
		return
	}
	if count >= s.c.Tag.ArcTagHateMaxNum {
		return ecode.TagArcTagLikeMaxFre
	}
	if _, err = s.normalArchive(c, aid); err != nil {
		return
	}
	s.dao.IncrSpamCache(c, mid, model.SpamHate)
	return s.hateService(c, mid, tid, aid, rpcModel.ResTypeArchive)
}
