package service

import (
	"context"
	"time"

	"go-common/app/interface/main/web/model"
	arcmdl "go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	coinmdl "go-common/app/service/main/coin/api"
	thumbup "go-common/app/service/main/thumbup/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

var _emptyCoinArcList = make([]*model.CoinArc, 0)

// Coins get archive User added coins.
func (s *Service) Coins(c context.Context, mid, aid int64) (res *model.ArchiveUserCoins, err error) {
	var rs *coinmdl.ItemUserCoinsReply
	if rs, err = s.coinClient.ItemUserCoins(c, &coinmdl.ItemUserCoinsReq{Mid: mid, Aid: aid, Business: model.CoinArcBusiness}); err != nil {
		log.Error("s.coinClient.ItemUserCoins(%d,%d) error(%v)", mid, aid, err)
		return
	}
	res = new(model.ArchiveUserCoins)
	if rs != nil {
		res.Multiply = rs.Number
	}
	return
}

// AddCoin add coin to archive.
func (s *Service) AddCoin(c context.Context, aid, mid, upID, multiply, avtype int64, business, ck, ua, refer string, now time.Time, selectLike int) (like bool, err error) {
	var (
		pubTime int64
		typeID  int32
		maxCoin int64 = 2
		ip            = metadata.String(c, metadata.RemoteIP)
	)
	switch avtype {
	case model.CoinAddArcType:
		var a *arcmdl.ArcReply
		if a, err = s.arcClient.Arc(c, &arcmdl.ArcRequest{Aid: aid}); err != nil {
			log.Error("s.arcClient.Arc(%v) error(%v)", aid, err)
			return
		}
		if !a.Arc.IsNormal() {
			err = ecode.ArchiveNotExist
			return
		}
		if a.Arc.Copyright == int32(archive.CopyrightCopy) {
			maxCoin = 1
		}
		upID = a.Arc.Author.Mid
		typeID = a.Arc.TypeID
		pubTime = int64(a.Arc.PubDate)
	case model.CoinAddArtType:
		maxCoin = 1
	}
	arg := &coinmdl.AddCoinReq{
		IP:       ip,
		Mid:      mid,
		Upmid:    upID,
		MaxCoin:  maxCoin,
		Aid:      aid,
		Business: business,
		Number:   multiply,
		Typeid:   typeID,
		PubTime:  pubTime,
	}
	if _, err = s.coinClient.AddCoin(c, arg); err == nil && avtype == model.CoinAddArcType && selectLike == 1 {
		if err = s.thumbup.Like(c, &thumbup.ArgLike{Mid: mid, UpMid: upID, Business: _businessLike, MessageID: aid, Type: thumbup.TypeLike, RealIP: ip}); err != nil {
			log.Error("AddCoin s.thumbup.Like  mid(%d) upID(%d) aid(%d) error(%+v)", mid, upID, aid, err)
			err = nil
		} else {
			like = true
		}
	}
	return
}

// CoinExp get coin exp today
func (s *Service) CoinExp(c context.Context, mid int64) (exp int64, err error) {
	var todayExp *coinmdl.TodayExpReply
	if todayExp, err = s.coinClient.TodayExp(c, &coinmdl.TodayExpReq{Mid: mid}); err != nil {
		log.Error("CoinExp s.coinClient.TodayExp mid(%d) error(%v)", mid, err)
		err = nil
		return
	}
	exp = todayExp.Exp
	return
}

// CoinList get coin list.
func (s *Service) CoinList(c context.Context, mid int64, pn, ps int) (list []*model.CoinArc, count int, err error) {
	var (
		coinReply *coinmdl.ListReply
		aids      []int64
		arcsReply *arcmdl.ArcsReply
	)
	if coinReply, err = s.coinClient.List(c, &coinmdl.ListReq{Mid: mid, Business: model.CoinArcBusiness, Ts: time.Now().Unix()}); err != nil {
		log.Error("CoinList s.coinClient.List(%d) error(%v)", mid, err)
		err = nil
		list = _emptyCoinArcList
		return
	}
	existAids := make(map[int64]int64, len(coinReply.List))
	afVideos := make(map[int64]*coinmdl.ModelList, len(coinReply.List))
	for _, v := range coinReply.List {
		if _, ok := existAids[v.Aid]; ok {
			afVideos[v.Aid].Number += v.Number
			continue
		}
		afVideos[v.Aid] = v
		aids = append(aids, v.Aid)
		existAids[v.Aid] = v.Aid
	}
	count = len(aids)
	start := (pn - 1) * ps
	end := pn * ps
	switch {
	case start > count:
		aids = aids[:0]
	case end >= count:
		aids = aids[start:]
	default:
		aids = aids[start:end]
	}
	if len(aids) == 0 {
		list = _emptyCoinArcList
		return
	}
	if arcsReply, err = s.arcClient.Arcs(c, &arcmdl.ArcsRequest{Aids: aids}); err != nil {
		log.Error("CoinList s.arcClient.Arcs(%v) error(%v)", aids, err)
		err = nil
		list = _emptyCoinArcList
		return
	}
	for _, aid := range aids {
		if arc, ok := arcsReply.Arcs[aid]; ok && arc.IsNormal() {
			if arc.Access >= 10000 {
				arc.Stat.View = 0
			}
			if item, ok := afVideos[aid]; ok {
				list = append(list, &model.CoinArc{Arc: arc, Coins: item.Number, Time: item.Ts})
			}
		}
	}
	if len(list) == 0 {
		list = _emptyCoinArcList
	}
	return
}
