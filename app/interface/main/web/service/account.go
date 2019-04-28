package service

import (
	"context"

	"go-common/app/interface/main/web/model"
	artmdl "go-common/app/interface/openplatform/article/model"
	accmdl "go-common/app/service/main/account/api"
	"go-common/app/service/main/archive/model/archive"
	relmdl "go-common/app/service/main/relation/model"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"
)

const _cardBakCacheRand = 10

// Attentions get attention list.
func (s *Service) Attentions(c context.Context, mid int64) (rs []int64, err error) {
	var (
		attentions []*relmdl.Following
		remoteIP   = metadata.String(c, metadata.RemoteIP)
	)
	if attentions, err = s.relation.Followings(c, &relmdl.ArgMid{Mid: mid, RealIP: remoteIP}); err != nil {
		log.Error("s.relation.Followings(%d,%s) error %v", mid, remoteIP, err)
	} else {
		rs = make([]int64, 0)
		for _, v := range attentions {
			rs = append(rs, v.Mid)
		}
	}
	return
}

// Card get card relation archive count data.
func (s *Service) Card(c context.Context, mid, loginID int64, topPhoto, article bool) (rs *model.Card, err error) {
	var (
		cardReply                                          *accmdl.CardReply
		relResp                                            *accmdl.RelationReply
		card                                               *model.AccountCard
		space                                              *model.Space
		arcCount                                           int
		group                                              *errgroup.Group
		infoErr, statErr, spaceErr, relErr, upcErr, artErr error
		cacheRs                                            *model.Card
		remoteIP                                           = metadata.String(c, metadata.RemoteIP)
	)
	relation := &accmdl.RelationReply{}
	stat := &relmdl.Stat{}
	upArts := &artmdl.UpArtMetas{}
	group, errCtx := errgroup.WithContext(c)
	card = new(model.AccountCard)
	card.Attentions = make([]int64, 0)
	group.Go(func() error {
		if cardReply, infoErr = s.accClient.Card3(errCtx, &accmdl.MidReq{Mid: mid}); infoErr != nil {
			log.Error("s.accClient.Card3(%d,%s) error %v", mid, remoteIP, infoErr)
		} else {
			card.FromCard(cardReply.Card)
		}
		return nil
	})
	group.Go(func() error {
		if stat, statErr = s.relation.Stat(errCtx, &relmdl.ArgMid{Mid: mid, RealIP: remoteIP}); statErr != nil {
			log.Error("s.relation.Stat(%d) error(%v)", mid, statErr)
		} else {
			card.Fans = int(stat.Follower)
			card.Attention = int(stat.Following)
			card.Friend = int(stat.Following)
		}
		return nil
	})
	if topPhoto {
		group.Go(func() error {
			space, spaceErr = s.dao.TopPhoto(errCtx, mid)
			return nil
		})
	}
	if loginID > 0 {
		group.Go(func() error {
			if relResp, relErr = s.accClient.Relation3(errCtx, &accmdl.RelationReq{Mid: loginID, Owner: mid, RealIp: remoteIP}); relErr != nil {
				log.Error("s.accClient.Relation3(%d,%d,%s) error %v", loginID, mid, remoteIP, relErr)
			} else if relResp != nil {
				relation = relResp
			}
			return nil
		})
	}
	group.Go(func() error {
		if arcCount, upcErr = s.arc.UpCount2(errCtx, &archive.ArgUpCount2{Mid: mid}); upcErr != nil {
			log.Error("s.arc.UpCount2(%d) error %v", mid, upcErr)
		}
		return nil
	})
	if article {
		group.Go(func() error {
			if upArts, artErr = s.art.UpArtMetas(errCtx, &artmdl.ArgUpArts{Mid: mid, Pn: _samplePn, Ps: _samplePs, RealIP: remoteIP}); artErr != nil {
				log.Error("s.art.UpArtMetas(%d) error(%v)", mid, artErr)
			}
			if upArts == nil {
				upArts = &artmdl.UpArtMetas{Count: 0}
			}
			return nil
		})
	}
	group.Wait()
	addCache := true
	if infoErr != nil || (topPhoto && spaceErr != nil) || (loginID > 0 && relErr != nil) || upcErr != nil {
		if cacheRs, err = s.dao.CardBakCache(c, mid); err != nil {
			addCache = false
			log.Error("s.dao.CardBakCache(%d) error (%v)", mid, err)
			err = nil
		} else if cacheRs != nil {
			if infoErr != nil {
				card = cacheRs.Card
			}
			if statErr != nil {
				stat = &relmdl.Stat{Follower: cacheRs.Follower}
			}
			if topPhoto && spaceErr != nil {
				space = cacheRs.Space
			}
			if loginID > 0 && relErr != nil {
				relation = &accmdl.RelationReply{Following: cacheRs.Following}
			}
			if upcErr != nil {
				arcCount = cacheRs.ArchiveCount
			}
			if artErr != nil {
				upArts = &artmdl.UpArtMetas{Count: cacheRs.ArticleCount}
			}
		}
		if topPhoto && space == nil {
			space = &model.Space{SImg: s.c.DefaultTop.SImg, LImg: s.c.DefaultTop.LImg}
		}
	}
	rs = &model.Card{
		Card:         card,
		Space:        space,
		Following:    relation.Following,
		ArchiveCount: arcCount,
		ArticleCount: upArts.Count,
		Follower:     stat.Follower,
	}
	if addCache {
		s.cache.Do(c, func(c context.Context) {
			if s.r.Intn(_cardBakCacheRand) == 1 {
				s.dao.SetCardBakCache(c, mid, rs)
			}
		})
	}
	return
}
