package service

import (
	"context"
	"time"

	"go-common/app/interface/main/growup/model"

	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	xtime "go-common/library/time"
)

// JoinColumn join up to creative column
func (s *Service) JoinColumn(c context.Context, mid int64, accountType, signType int) (err error) {
	id, err := s.dao.Blocked(c, mid)
	if err != nil {
		log.Error("s.dao.GetBlocked mid(%d) error(%v)", mid, err)
		return
	}
	if id != 0 {
		log.Info("mid(%d) is blocked", mid)
		return ecode.GrowupDisabled
	}

	// get up view
	ip := metadata.String(c, metadata.RemoteIP)
	stat, err := s.dao.ArticleStat(c, mid, ip)
	if err != nil {
		log.Error("s.dao.ArticleStat mid(%d) error(%v)", mid, err)
		return
	}
	if stat.View < s.conf.Threshold.LimitArticleView {
		log.Info("mid(%d) view(%s) not reach standard", mid, stat.View)
		return ecode.GrowupDisabled
	}

	// get up nickname
	card, err := s.dao.Card(c, mid)
	if err != nil {
		log.Error("s.dao.Card(%d) error(%v)", mid, err)
		return
	}

	fans, err := s.dao.Fans(c, mid)
	if err != nil {
		return
	}
	state, err := s.dao.GetAccountState(c, "up_info_column", mid)
	if err != nil {
		return
	}

	// if account state is 2 3 4 5 6 7 return
	if state >= 2 && state < 8 {
		return
	}

	now := xtime.Time(time.Now().Unix())
	// sign_type: 1.basic; 2.first publish; 0:default.
	v := &model.UpInfo{
		MID:            mid,
		Nickname:       card.Name,
		AccountType:    accountType,
		MainCategory:   0,
		Fans:           fans,
		AccountState:   2,
		SignType:       signType,
		ApplyAt:        now,
		TotalPlayCount: stat.View,
	}

	_, err = s.dao.InsertUpInfo(c, "up_info_column", "total_view_count", v)
	return
}
