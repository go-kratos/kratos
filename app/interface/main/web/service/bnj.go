package service

import (
	"context"
	"time"

	"go-common/app/interface/main/web/model"
	accmdl "go-common/app/service/main/account/api"
	arcmdl "go-common/app/service/main/archive/api"
	coinmdl "go-common/app/service/main/coin/api"
	favmdl "go-common/app/service/main/favorite/model"
	thumbup "go-common/app/service/main/thumbup/model"
	"go-common/library/conf/env"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"
)

func (s *Service) checkBnjAccess(mid int64) bool {
	if s.c.Bnj2019.Open {
		return true
	}
	if env.DeployEnv == env.DeployEnvProd && len(s.bnjGrayUids) > 0 {
		if mid == 0 {
			return false
		}
		if _, ok := s.bnjGrayUids[mid]; !ok {
			return false
		}
	}
	return true
}

// Bnj2019Aids get bnj aids.3
func (s Service) Bnj2019Aids(c context.Context) []int64 {
	aids := s.c.Bnj2019.BnjListAids
	aids = append(aids, s.c.Bnj2019.BnjMainAid)
	return aids
}

// Timeline get timeline.
func (s *Service) Timeline(c context.Context, mid int64) (data []*model.Timeline, err error) {
	if !s.checkBnjAccess(mid) {
		err = ecode.AccessDenied
		return
	}
	for _, v := range s.c.Bnj2019.Timeline {
		data = append(data, &model.Timeline{
			Name:    v.Name,
			Start:   v.Start.Unix(),
			End:     v.End.Unix(),
			Cover:   v.Cover,
			H5Cover: v.H5Cover,
		})
	}
	return
}

// Bnj2019 get bnj2019 data.
func (s *Service) Bnj2019(c context.Context, mid int64) (data *model.Bnj2019, err error) {
	if !s.checkBnjAccess(mid) {
		err = ecode.AccessDenied
		return
	}
	if s.bnj2019View == nil || !s.bnj2019View.Arc.IsNormal() {
		err = ecode.NothingFound
		return
	}
	data = &model.Bnj2019{
		Bnj2019View: s.bnj2019View,
		Elec:        s.BnjElecInfo,
		Related:     s.bnj2019List,
		ReqUser:     &model.ReqUser{},
	}
	if len(data.Related) == 0 {
		data.Related = make([]*model.Bnj2019Related, 0)
	}
	if mid > 0 {
		authorMid := s.bnj2019View.Author.Mid
		aid := s.bnj2019View.Aid
		ip := metadata.String(c, metadata.RemoteIP)
		group, errCtx := errgroup.WithContext(c)
		// attention
		group.Go(func() error {
			if resp, e := s.accClient.Relation3(errCtx, &accmdl.RelationReq{Mid: mid, Owner: authorMid, RealIp: ip}); e != nil {
				log.Error("Bnj2019 s.accClient.Relation3(%d,%d,%s) error(%v)", mid, authorMid, ip, e)
			} else if resp != nil {
				data.ReqUser.Attention = resp.Following
			}
			return nil
		})
		// favorite
		group.Go(func() error {
			if resp, e := s.fav.IsFav(errCtx, &favmdl.ArgIsFav{Type: favmdl.TypeVideo, Mid: mid, Oid: aid, RealIP: ip}); e != nil {
				log.Error("Bnj2019 s.fav.IsFav(%d,%d,%s) error(%v)", mid, aid, ip, e)
			} else {
				data.ReqUser.Favorite = resp
			}
			return nil
		})
		// like
		group.Go(func() error {
			if resp, e := s.thumbup.HasLike(errCtx, &thumbup.ArgHasLike{Business: _businessLike, MessageIDs: []int64{aid}, Mid: mid, RealIP: ip}); e != nil {
				log.Error("Bnj2019 s.thumbup.HasLike(%d,%d,%s) error %v", aid, mid, ip, e)
			} else if resp != nil {
				if v, ok := resp[aid]; ok {
					switch v {
					case thumbup.StateLike:
						data.ReqUser.Like = true
					case thumbup.StateDislike:
						data.ReqUser.Dislike = true
					}
				}
			}
			return nil
		})
		// coin
		group.Go(func() error {
			if resp, e := s.coinClient.ItemUserCoins(errCtx, &coinmdl.ItemUserCoinsReq{Mid: mid, Aid: aid, Business: model.CoinArcBusiness}); e != nil {
				log.Error("Bnj2019 s.coinClient.ItemUserCoins(%d,%d,%s) error %v", mid, aid, ip, e)
			} else if resp != nil {
				data.ReqUser.Coin = resp.Number
			}
			return nil
		})
		group.Wait()
	}
	return
}

func (s *Service) bnj2019proc() {
	// bnj gray uid
	go func() {
		for {
			time.Sleep(time.Duration(s.c.Bnj2019.BnjTick))
			if mids, err := s.dao.Bnj2019Conf(context.Background()); err != nil {
				log.Error("bnj2019proc s.dao.Bnj2019Conf error(%v)", err)
				continue
			} else {
				tmp := make(map[int64]struct{}, len(mids))
				if len(mids) > 0 {
					for _, mid := range mids {
						tmp[mid] = struct{}{}
					}
				}
				s.bnjGrayUids = tmp
			}
		}
	}()
	// main arc
	go func() {
		for {
			time.Sleep(time.Duration(s.c.Bnj2019.BnjTick))
			if s.c.Bnj2019.BnjMainAid == 0 {
				continue
			}
			if viewReply, err := s.arcClient.View(context.Background(), &arcmdl.ViewRequest{Aid: s.c.Bnj2019.BnjMainAid}); err != nil {
				log.Error("bnj2019proc main s.arcClient.View(%d) error(%v)", s.c.Bnj2019.BnjMainAid, err)
				continue
			} else if viewReply != nil {
				s.bnj2019View = &model.Bnj2019View{Arc: viewReply.Arc, Pages: viewReply.Pages}
				// elec
				if elec, err := s.dao.ElecShow(context.Background(), viewReply.Arc.Author.Mid, viewReply.Arc.Aid, 0); err == nil {
					s.BnjElecInfo = elec
					s.BnjElecInfo.TotalCount += s.c.Bnj2019.FakeElec
				} else {
					log.Error("bnj2019proc s.dao.ElecShow(%d,%d) error(%v)", viewReply.Arc.Author.Mid, viewReply.Arc.Aid, err)
				}
			}
		}
	}()
	// live arc
	go func() {
		for {
			time.Sleep(time.Duration(s.c.Bnj2019.BnjTick))
			if s.c.Bnj2019.LiveAid == 0 {
				continue
			}
			if arcReply, err := s.arcClient.Arc(context.Background(), &arcmdl.ArcRequest{Aid: s.c.Bnj2019.LiveAid}); err != nil {
				log.Error("bnj2019proc live arc s.arcClient.Arc(%d) error(%v)", s.c.Bnj2019.LiveAid, err)
				continue
			} else if arcReply != nil {
				s.bnj2019LiveArc = arcReply
			}
		}
	}()
	// list arc
	go func() {
		for {
			time.Sleep(time.Duration(s.c.Bnj2019.BnjTick))
			if len(s.c.Bnj2019.BnjListAids) == 0 {
				continue
			}
			if viewsReply, err := s.arcClient.Views(context.Background(), &arcmdl.ViewsRequest{Aids: s.c.Bnj2019.BnjListAids}); err != nil {
				log.Error("bnj2019proc list s.arcClient.Views(%v) error(%v)", s.c.Bnj2019.BnjListAids, err)
				continue
			} else {
				var tmpList []*model.Bnj2019Related
				for _, aid := range s.c.Bnj2019.BnjListAids {
					if view, ok := viewsReply.Views[aid]; ok && view.Arc.IsNormal() {
						item := &model.Bnj2019Related{Arc: view.Arc, Pages: view.Pages}
						tmpList = append(tmpList, item)
					}
				}
				s.bnj2019List = tmpList
			}
		}
	}()
}
