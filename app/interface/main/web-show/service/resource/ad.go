package resource

import (
	"context"
	"math/rand"

	"go-common/app/interface/main/web-show/dao/ad"
	resmdl "go-common/app/interface/main/web-show/model/resource"
	account "go-common/app/service/main/account/model"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

var (
	_emptyVideoAds = []*resmdl.VideoAD{}
)

// VideoAd get videoad by aid
func (s *Service) VideoAd(c context.Context, arg *resmdl.ArgAid) (res []*resmdl.VideoAD) {
	arg.IP = metadata.String(c, metadata.RemoteIP)
	if arg.Mid > 0 {
		// ignore error
		var (
			resPro *account.Card
			err    error
		)
		if resPro, err = s.user(c, arg.Mid, arg.IP); err == nil {
			if s.normalVip(c, resPro) {
				return
			}
		}
		// NOTE cache?
		if isBp := s.bangumiDao.IsBp(c, arg.Mid, arg.Aid, arg.IP); isBp {
			log.Info("mid(%d) aid(%d) is bp", arg.Mid, arg.IP)
			res = _emptyVideoAds
			return
		}
	}
	if res = s.videoAdByAid(arg.Aid); len(res) == 0 {
		res = _emptyVideoAds
	}
	return
}

func (s *Service) user(c context.Context, mid int64, ip string) (resPro *account.Card, err error) {
	arg := &account.ArgMid{
		Mid: mid,
	}
	resPro, err = s.accRPC.Card3(c, arg)
	if err != nil {
		ad.PromError("accRPC.Info2", "s.accRPC.Info2() err(%v)", err)
		log.Error("s.accRPC.Info2() err(%v)", err)
	}
	return
}

// checkVip check normal vip
func (s *Service) normalVip(c context.Context, pro *account.Card) bool {
	if pro.Vip.Type != 0 && pro.Vip.Status == 1 {
		return true
	}
	return false
}

func (s *Service) videoAdByAid(aid int64) (res []*resmdl.VideoAD) {
	ss := s.videoCache[aid]
	l := len(ss)
	if l == 0 {
		return
	}
	// NOTE this means StrategyOnly
	if l == 1 {
		res = ss[0]
		return
	}
	// NOTE this means StrategyShare
	res = ss[rand.Intn(l)]
	return
}

// loadVideoAd load videoad to cache
func (s *Service) loadVideoAd() (err error) {
	ads, err := s.resdao.VideoAds(context.Background())
	if err != nil {
		log.Error("s.resdao.VideoAds error(%v)", err)
		return
	}
	tmp := make(map[int64][][]*resmdl.VideoAD)
	for aid, vads := range ads {
		if len(vads) < 1 {
			continue
		}
		if vads[0].Strategy == resmdl.StrategyOnly || vads[0].Strategy == resmdl.StrategyRank {
			tmp[aid] = append(tmp[aid], vads)
		} else if vads[0].Strategy == resmdl.StrategyShare {
			for _, vad := range vads {
				tmp[aid] = append(tmp[aid], []*resmdl.VideoAD{vad})
			}
		}
	}
	s.videoCache = tmp
	return
}
