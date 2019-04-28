package service

import (
	"context"

	"go-common/app/interface/main/web/model"
	arcmdl "go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	coinmdl "go-common/app/service/main/coin/api"
	favmdl "go-common/app/service/main/favorite/model"
	thumbup "go-common/app/service/main/thumbup/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"
)

const (
	_businessLike = "archive"
)

// Like like archive.
func (s *Service) Like(c context.Context, aid, mid int64, like int8) (upperID int64, err error) {
	var (
		arcReply *arcmdl.ArcReply
		ip       = metadata.String(c, metadata.RemoteIP)
	)
	if arcReply, err = s.arcClient.Arc(c, &arcmdl.ArcRequest{Aid: aid}); err != nil {
		log.Error("s.arcClient.Arc(%d) error(%v)", aid, err)
		return
	}
	arc := arcReply.Arc
	if !arc.IsNormal() {
		err = ecode.ArchiveNotExist
		return
	}
	upperID = arc.Author.Mid
	err = s.thumbup.Like(c, &thumbup.ArgLike{Mid: mid, UpMid: upperID, Business: _businessLike, MessageID: aid, Type: like, RealIP: ip})
	return
}

// LikeTriple like & coin & fav
func (s *Service) LikeTriple(c context.Context, aid, mid int64) (res *model.TripleRes, err error) {
	var (
		arcReply *arcmdl.ArcReply
		ip       = metadata.String(c, metadata.RemoteIP)
	)
	res = new(model.TripleRes)
	maxCoin := int64(1)
	multiply := int64(1)
	if arcReply, err = s.arcClient.Arc(c, &arcmdl.ArcRequest{Aid: aid}); err != nil {
		log.Error("s.arcClient.Arc(%d) error(%v)", aid, err)
		return
	}
	a := arcReply.Arc
	if !a.IsNormal() {
		err = ecode.ArchiveNotExist
		return
	}
	if a.Copyright == int32(archive.CopyrightOriginal) {
		maxCoin = 2
		multiply = 2
	}
	res.UpID = a.Author.Mid
	eg := errgroup.Group{}
	eg.Go(func() (err error) {
		if multiply == 2 {
			if userCoins, e := s.coinClient.UserCoins(c, &coinmdl.UserCoinsReq{Mid: mid}); e != nil {
				log.Error("s.coinClient.UserCoins error(%v)", e)
			} else if userCoins != nil {
				if userCoins.Count < 1 {
					return
				}
				if userCoins.Count < 2 {
					multiply = 1
				}
			}
		}
		cArg := &coinmdl.AddCoinReq{
			IP:       ip,
			Mid:      mid,
			Upmid:    a.Author.Mid,
			MaxCoin:  maxCoin,
			Aid:      aid,
			Business: model.CoinArcBusiness,
			Number:   multiply,
			Typeid:   a.TypeID,
			PubTime:  int64(a.PubDate),
		}
		if _, err = s.coinClient.AddCoin(c, cArg); err != nil {
			log.Error("s.coinClient.AddCoin error(%v)", err)
			err = nil
			if arcUserCoins, e := s.coinClient.ItemUserCoins(c, &coinmdl.ItemUserCoinsReq{Mid: mid, Aid: aid, Business: model.CoinArcBusiness}); e != nil {
				log.Error("s.coinClient.ItemUserCoins error(%v)", e)
			} else {
				if arcUserCoins != nil && arcUserCoins.Number > 0 {
					res.Coin = true
				}
			}
		} else {
			res.Multiply = multiply
			res.Anticheat = true
			res.Coin = true
		}
		return
	})
	eg.Go(func() (err error) {
		var isFav bool
		if isFav, err = s.fav.IsFav(context.Background(), &favmdl.ArgIsFav{Type: favmdl.TypeVideo, Mid: mid, Oid: aid, RealIP: ip}); err != nil {
			log.Error("s.fav.IsFav error(%v)", err)
			err = nil
		} else if isFav {
			res.Fav = true
			return
		}
		fArg := &favmdl.ArgAdd{Type: favmdl.TypeVideo, Mid: mid, Oid: aid, Fid: 0, RealIP: ip}
		if err = s.fav.Add(c, fArg); err != nil {
			if ecode.FavVideoExist.Equal(err) {
				res.Fav = true
				return
			}
			log.Error("s.fav.Add error(%v)", err)
			err = nil
		} else {
			res.Fav = true
			res.Anticheat = true
		}
		return
	})
	eg.Go(func() (err error) {
		if err = s.thumbup.Like(c, &thumbup.ArgLike{Mid: mid, UpMid: res.UpID, Business: _businessLike, MessageID: aid, Type: thumbup.TypeLike, RealIP: ip}); err != nil {
			if ecode.ThumbupDupLikeErr.Equal(err) {
				res.Like = true
				return
			}
			log.Error("s.thumbup.Like error(%v)", err)
			err = nil
		} else {
			res.Like = true
			res.Anticheat = true
		}
		return
	})
	eg.Wait()
	return
}

// HasLike get if has like.
func (s *Service) HasLike(c context.Context, aid, mid int64) (like int8, err error) {
	var (
		data map[int64]int8
		ip   = metadata.String(c, metadata.RemoteIP)
	)
	if data, err = s.thumbup.HasLike(c, &thumbup.ArgHasLike{Business: _businessLike, MessageIDs: []int64{aid}, Mid: mid, RealIP: ip}); err != nil {
		log.Error("s.thumbup.HasLike aid(%d) mid(%d) error(%v)", aid, mid, err)
		return
	}
	if v, ok := data[aid]; ok {
		like = v
	}
	return
}
