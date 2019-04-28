package view

import (
	"context"
	"fmt"

	"go-common/app/interface/main/app-view/model/elec"
	"go-common/app/interface/main/app-view/model/view"
	"go-common/app/service/main/archive/model/archive"
	thumbup "go-common/app/service/main/thumbup/model"
	"go-common/library/conf/env"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

// CheckAccess check
func (s *Service) CheckAccess(mid int64) bool {
	if !s.BnjIsGrey {
		return true
	}
	_, ok := s.BnjWhiteMid[mid]
	if !ok && env.DeployEnv == env.DeployEnvProd {
		log.Error("mid(%d) env(%s) not allow", mid, env.DeployEnv)
		return false
	}
	return true
}

func (s *Service) initBnjPages(c context.Context, ps []*archive.Page3) (pages []*view.Page) {
	for _, v := range ps {
		page := &view.Page{}
		metas := make([]*view.Meta, 0, 4)
		for q, r := range _rate {
			meta := &view.Meta{
				Quality: q,
				Size:    int64(float64(r*v.Duration) * 1.1 / 8.0),
			}
			metas = append(metas, meta)
		}
		page.Page3 = v
		page.Metas = metas
		page.DMLink = fmt.Sprintf(_dmformat, v.Cid)
		pages = append(pages, page)
	}
	return
}

// initBnjReqUser is
func (s *Service) initBnjReqUser(c context.Context, authorMid, aid, mid int64) (reqUser *view.ReqUser, err error) {
	reqUser = &view.ReqUser{Favorite: 0, Attention: -999, Like: 0, Dislike: 0}
	if mid == 0 {
		return
	}
	g := errgroup.Group{}
	g.Go(func() error {
		var is bool
		if is, _ = s.favDao.IsFav(context.Background(), mid, aid); is {
			reqUser.Favorite = 1
		}
		return nil
	})
	g.Go(func() error {
		res, err := s.thumbupDao.HasLike(context.Background(), mid, _businessLike, []int64{aid})
		if err != nil {
			log.Error("s.thumbupDao.HasLike err(%+v)", err)
			return nil
		}
		if res.States == nil {
			return nil
		}
		if typ, ok := res.States[aid]; ok {
			if typ.State == thumbup.StateLike {
				reqUser.Like = 1
			} else if typ.State == thumbup.StateDislike {
				reqUser.Dislike = 1
			}
		}
		return nil
	})
	g.Go(func() (err error) {
		res, err := s.coinDao.ArchiveUserCoins(context.Background(), aid, mid, _avTypeAv)
		if err != nil {
			log.Error("%+v", err)
			err = nil
		}
		if res != nil && res.Multiply > 0 {
			reqUser.Coin = 1
		}
		return
	})
	if authorMid > 0 {
		g.Go(func() error {
			fl, err := s.accDao.Following3(context.Background(), mid, authorMid)
			if err != nil {
				log.Error("%+v", err)
				return nil
			}
			if fl {
				reqUser.Attention = 1
			}
			return nil
		})
	}
	g.Wait()
	return
}

// Bnj2019 is
func (s *Service) Bnj2019(c context.Context, mid int64, relateID int64) (bnj *view.BnjMain, err error) {
	if s.BnjMainView == nil || !s.BnjMainView.IsNormal() {
		err = ecode.NothingFound
		return
	}
	bnj = new(view.BnjMain)
	bnj.ElecSmallText = s.c.Bnj2019.ElecSmallText
	bnj.ElecBigText = s.c.Bnj2019.ElecBigText
	bnj.Archive3 = s.BnjMainView.Archive3
	bnj.ReqUser = &view.ReqUser{}
	bnj.Elec = s.BnjElecInfo
	bnj.Pages = s.initBnjPages(c, s.BnjMainView.Pages)
	bnj.ReqUser, _ = s.initBnjReqUser(c, bnj.Author.Mid, bnj.Aid, mid)
	bnj.PlayerIcon = s.playerIcon
	bnj.Elec = s.BnjElecInfo
	for _, a := range s.BnjLists {
		relate := &view.BnjItem{
			Aid:       a.Aid,
			Cid:       a.FirstCid,
			Tid:       a.TypeID,
			Pic:       a.Pic,
			Copyright: a.Copyright,
			PubDate:   a.PubDate,
			Title:     a.Title,
			Desc:      a.Desc,
			Stat:      a.Stat,
			Duration:  a.Duration,
			Author:    a.Author,
			Dimension: a.Dimension,
			Rights:    a.Rights,
		}
		if relate.Aid == s.c.Bnj2019.AdAv {
			relate.IsAd = 1
		}
		if relate.Aid == relateID {
			relate.Pages = s.initBnjPages(c, a.Pages)
			relate.ReqUser, _ = s.initBnjReqUser(c, a.Author.Mid, a.Aid, mid)
		}
		bnj.Relates = append(bnj.Relates, relate)
	}
	return
}

// BnjList is
func (s *Service) BnjList(c context.Context, mid int64) (list *view.BnjList, err error) {
	list = new(view.BnjList)
	for _, item := range s.BnjLists {
		list.Item = append(list.Item, &view.BnjItem{
			Aid:      item.Aid,
			Cid:      item.FirstCid,
			Pic:      item.Pic,
			Duration: item.Duration,
			IsAd:     0,
			Author:   item.Author,
		})
	}
	return
}

// BnjItem is
func (s *Service) BnjItem(c context.Context, aid, mid int64) (item *view.BnjItem, err error) {
	var v *archive.View3
	if aid == s.BnjMainView.Aid {
		v = s.BnjMainView
	} else {
		for _, l := range s.BnjLists {
			if aid == l.Aid {
				v = l
				break
			}
		}
	}
	if v == nil || !v.IsNormal() {
		err = ecode.NothingFound
		return
	}
	item = &view.BnjItem{
		Aid:       v.Aid,
		Cid:       v.FirstCid,
		Tid:       v.TypeID,
		Pic:       v.Pic,
		Copyright: v.Copyright,
		PubDate:   v.PubDate,
		Title:     v.Title,
		Desc:      v.Desc,
		Stat:      v.Stat,
		Duration:  v.Duration,
		Author:    v.Author,
		Dimension: v.Dimension,
		Rights:    v.Rights,
	}
	if item.Aid == s.c.Bnj2019.AdAv {
		item.IsAd = 1
	}
	item.ReqUser, _ = s.initBnjReqUser(c, v.Author.Mid, v.Aid, mid)
	item.Pages = s.initBnjPages(c, v.Pages)
	return
}

func (s *Service) loadBnj2019Infos() (err error) {
	var (
		aids      []int64
		avm       map[int64]*archive.View3
		list      []*archive.View3
		mainView  *archive.View3
		elec      *elec.Info
		whiteMid  = make(map[int64]struct{})
		liveMids  []int64
		bnjStatus int
	)
	if bnjStatus, liveMids, err = s.liveDao.Bnj2019Conf(context.Background()); err != nil {
		log.Error("%+v", err)
	} else {
		log.Info("got live bnj2019 mids(%v)", liveMids)
		for _, mid := range liveMids {
			whiteMid[mid] = struct{}{}
		}
	}
	if bnjStatus == 1 {
		s.BnjIsGrey = true
	} else {
		s.BnjIsGrey = false
	}
	// TODO live mids
	for _, mid := range s.c.Bnj2019.WhiteMids {
		whiteMid[mid] = struct{}{}
	}
	s.BnjWhiteMid = whiteMid
	aids = append(aids, s.c.Bnj2019.MainAid)
	aids = append(aids, s.c.Bnj2019.AidList...)
	if avm, err = s.arcDao.ViewsRPC(context.Background(), aids); err != nil {
		log.Error("bnj s.arcDao.Archives(%v) error(%v)", aids, err)
		return
	}
	mainView, ok := avm[s.c.Bnj2019.MainAid]
	if !ok {
		log.Error("bnj main archive(%d) not exist", s.c.Bnj2019.MainAid)
		return
	}
	mainView.Rights.Elec = 1
	mainView.Rights.Download = 1
	s.BnjMainView = mainView
	for _, aid := range s.c.Bnj2019.AidList {
		a, ok := avm[aid]
		if !ok {
			log.Error("bnj list has no aid(%d)", aid)
			continue
		}
		if !a.IsNormal() {
			log.Error("bnj list aid(%d) not open(%d)", aid, a.State)
			continue
		}
		a.Rights.Elec = 1
		a.Rights.Download = 1
		list = append(list, a)
	}
	if len(list) == 0 {
		log.Error("list is zero")
		return
	}
	s.BnjLists = list
	if elec, err = s.elcDao.TotalInfo(context.Background(), mainView.Author.Mid, mainView.Aid); err == nil {
		s.BnjElecInfo = elec
		s.BnjElecInfo.Total += s.c.Bnj2019.FakeElec
	} else {
		log.Error("s.elecDao.TotalInfo(%d,%d) error(%v)", mainView.Author.Mid, mainView.Aid, err)
	}
	return
}
