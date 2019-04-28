package view

import (
	"context"
	"fmt"
	"time"

	"go-common/app/interface/main/app-intl/model"
	"go-common/app/interface/main/app-intl/model/bangumi"
	"go-common/app/interface/main/app-intl/model/manager"
	"go-common/app/interface/main/app-intl/model/view"
	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	location "go-common/app/service/main/location/model"
	thumbup "go-common/app/service/main/thumbup/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
	"go-common/library/text/translate/chinese"
)

var (
	_rate = map[int]int64{15: 464, 16: 464, 32: 1028, 48: 1328, 64: 2192, 74: 3192, 80: 3192, 112: 6192, 116: 6192, 66: 1820}
)

const (
	_dmformat       = "http://comment.bilibili.com/%d.xml"
	_qn480          = 32
	_qnAndroidBuild = 5325000
	_qnIosBuild     = 8170
)

// initReqUser init Req User
func (s *Service) initReqUser(c context.Context, v *view.View, mid int64, plat int8, build int) {
	// owner ext
	g, ctx := errgroup.WithContext(c)
	if v.Author.Mid > 0 {
		g.Go(func() (err error) {
			v.OwnerExt.OfficialVerify.Type = -1
			card, err := s.accDao.Card3(ctx, v.Author.Mid)
			if err != nil {
				log.Error("%+v", err)
				err = nil
				return
			}
			if card != nil {
				otp := -1
				odesc := ""
				if card.Official.Role != 0 {
					if card.Official.Role <= 2 {
						otp = 0
					} else {
						otp = 1
					}
					odesc = card.Official.Title
				}
				v.OwnerExt.OfficialVerify.Type = otp
				v.OwnerExt.OfficialVerify.Desc = odesc
				v.OwnerExt.Vip.Type = int(card.Vip.Type)
				v.OwnerExt.Vip.VipStatus = int(card.Vip.Status)
				v.OwnerExt.Vip.DueDate = card.Vip.DueDate
				v.Author.Name = card.Name
				v.Author.Face = card.Face
			}
			return
		})
		g.Go(func() (err error) {
			stat, err := s.relDao.Stat(c, v.Author.Mid)
			if err != nil {
				log.Error("%+v", err)
				err = nil
				return
			}
			if stat != nil {
				v.OwnerExt.Fans = int(stat.Follower)
			}
			return
		})
		g.Go(func() error {
			if ass, err := s.assDao.Assist(ctx, v.Author.Mid); err != nil {
				log.Error("%+v", err)
			} else {
				v.OwnerExt.Assists = ass
			}
			return nil
		})
	}
	// req user
	v.ReqUser = &view.ReqUser{Favorite: 0, Attention: -999, Like: 0, Dislike: 0}
	// check req user
	if mid > 0 {
		g.Go(func() error {
			if is, _ := s.favDao.IsFav(ctx, mid, v.Aid); is {
				v.ReqUser.Favorite = 1
			}
			return nil
		})
		g.Go(func() error {
			res, err := s.thumbupDao.HasLike(ctx, mid, _businessLike, []int64{v.Aid})
			if err != nil {
				log.Error("%+v", err)
				return nil
			}
			if typ, ok := res[v.Aid]; ok {
				if typ == thumbup.StateLike {
					v.ReqUser.Like = 1
				} else if typ == thumbup.StateDislike {
					v.ReqUser.Dislike = 1
				}
			}
			return nil
		})
		g.Go(func() (err error) {
			res, err := s.coinDao.ArchiveUserCoins(ctx, v.Aid, mid, _avTypeAv)
			if err != nil {
				log.Error("%+v", err)
				err = nil
			}
			if res != nil && res.Multiply > 0 {
				v.ReqUser.Coin = 1
			}
			return
		})
		if v.Author.Mid > 0 {
			g.Go(func() error {
				fl, err := s.accDao.Following3(ctx, mid, v.Author.Mid)
				if err != nil {
					log.Error("%+v", err)
					return nil
				}
				if fl {
					v.ReqUser.Attention = 1
				}
				return nil
			})
		}
	}
	if err := g.Wait(); err != nil {
		log.Error("%+v", err)
	}
}

// initRelateCMTag is.
func (s *Service) initRelateCMTag(c context.Context, v *view.View, plat int8, build int, mid int64, buvid, mobiApp, device, network, adExtra, from string, now time.Time, isTW bool) {
	var (
		rls []*view.Relate
		mr  *manager.Relate
		err error
	)
	tids := s.initTag(c, v, mid)
	g, ctx := errgroup.WithContext(c)
	g.Go(func() (err error) {
		if mid > 0 || buvid != "" {
			if rls, v.UserFeature, v.ReturnCode, err = s.newRcmdRelate(ctx, plat, v.Aid, mid, buvid, mobiApp, from, build); err != nil {
				log.Error("s.newRcmdRelate(%d) error(%+v)", v.Aid, err)
			}
		}
		if len(rls) == 0 {
			rls, err = s.dealRcmdRelate(ctx, plat, v.Aid)
			log.Warn("s.dealRcmdRelate aid(%d) mid(%d) buvid(%s)", v.Aid, mid, buvid)
			return
		}
		v.IsRec = 1
		log.Warn("s.newRcmdRelate returncode(%s) aid(%d) mid(%d) buvid(%s)", v.ReturnCode, v.Aid, mid, buvid)
		return
	})
	g.Go(func() (err error) {
		mr = s.relateCache(ctx, plat, build, now, v.Aid, tids, v.TypeID)
		return
	})
	if err = g.Wait(); err != nil {
		log.Error("%+v", err)
		return
	}
	if len(rls) == 0 {
		s.prom.Incr("zero_relates")
		return
	}
	var (
		r  *view.Relate
		rm map[int]*view.Relate
	)
	aidm := map[int64]struct{}{}
	if r, err = s.dealManagerRelate(c, plat, mr); err != nil {
		log.Error("%+v", err)
	}
	if r != nil {
		// 相关推荐的第一位是运营插入位
		rm = map[int]*view.Relate{0: r}
		aidm[r.Aid] = struct{}{}
	}
	// 权重：详情页稿件>相关推荐运营位稿件>AI相关推荐稿件
	for _, rl := range rls {
		// AI相关推荐稿件不能和详情页稿件重复
		if rl.Aid == v.Aid {
			continue
		}
		if len(rm) != 0 {
			for {
				// 相关推荐运营位稿件不能和详情页稿件重复
				if r, ok := rm[len(v.Relates)]; ok && r.Aid != v.Aid {
					v.Relates = append(v.Relates, r)
					continue
				}
				// AI相关推荐稿件不能和相关推荐运营位稿件重复
				if _, ok := aidm[rl.Aid]; !ok {
					v.Relates = append(v.Relates, rl)
				}
				break
			}
		} else {
			v.Relates = append(v.Relates, rl)
		}
	}
	if isTW {
		for _, rl := range v.Relates {
			rl.Title = chinese.Convert(c, rl.Title)
		}
	}
}

// initMovie is.
func (s *Service) initMovie(c context.Context, v *view.View, mid int64, build int, mobiApp, device string, nMovie bool) (err error) {
	s.pHit.Incr("is_movie")
	var m *bangumi.Movie
	if m, err = s.banDao.Movie(c, v.Aid, mid, build, mobiApp, device); err != nil || m == nil {
		log.Error("%+v", err)
		err = ecode.NothingFound
		s.pMiss.Incr("err_is_PGC")
		return
	}
	if v.Rights.HD5 == 1 && m.PayUser.Status == 0 && !s.checkVIP(c, mid) {
		v.Rights.HD5 = 0
	}
	if len(m.List) == 0 {
		err = ecode.NothingFound
		return
	}
	vps := make([]*view.Page, 0, len(m.List))
	for _, l := range m.List {
		vp := &view.Page{
			Page3: &archive.Page3{Cid: l.Cid, Page: int32(l.Page), From: l.Type, Part: l.Part, Vid: l.Vid},
		}
		vps = append(vps, vp)
	}
	m.List = nil
	// view
	v.Pages = vps
	v.Rights.Download = int32(m.AllowDownload)
	m.AllowDownload = 0
	v.Rights.Bp = 0
	if nMovie {
		v.Movie = m
		v.Desc = m.Season.Evaluate
	}
	return
}

// initPGC is.
func (s *Service) initPGC(c context.Context, v *view.View, mid int64, build int, mobiApp, device string) (err error) {
	s.pHit.Incr("is_PGC")
	var season *bangumi.Season
	if season, err = s.banDao.PGC(c, v.Aid, mid, build, mobiApp, device); err != nil {
		log.Error("%+v", err)
		err = ecode.NothingFound
		s.pMiss.Incr("err_is_PGC")
		return
	}
	if season != nil {
		if season.Player != nil {
			if len(v.Pages) != 0 {
				if season.Player.Cid != 0 {
					v.Pages[0].Cid = season.Player.Cid
				}
				if season.Player.From != "" {
					v.Pages[0].From = season.Player.From
				}
				if season.Player.Vid != "" {
					v.Pages[0].Vid = season.Player.Vid
				}
			}
			season.Player = nil
		}
		if season.AllowDownload == "1" {
			v.Rights.Download = 1
		} else {
			v.Rights.Download = 0
		}
		if season.SeasonID != "" {
			season.AllowDownload = ""
			v.Season = season
		}
	}
	if v.Rights.HD5 == 1 && !s.checkVIP(c, mid) {
		v.Rights.HD5 = 0
	}
	v.Rights.Bp = 0
	return
}

// initPages is.
func (s *Service) initPages(c context.Context, vs *view.ViewStatic, ap []*archive.Page3) {
	pages := make([]*view.Page, 0, len(ap))
	for _, v := range ap {
		page := &view.Page{}
		metas := make([]*view.Meta, 0, 4)
		for q, r := range _rate {
			meta := &view.Meta{
				Quality: q,
				Size:    int64(float64(r*v.Duration) * 1.1 / 8.0),
			}
			metas = append(metas, meta)
		}
		if vs.AttrVal(archive.AttrBitIsBangumi) == archive.AttrYes {
			v.From = "bangumi"
		}
		page.Page3 = v
		page.Metas = metas
		page.DMLink = fmt.Sprintf(_dmformat, v.Cid)
		pages = append(pages, page)
	}
	vs.Pages = pages
}

// initDownload is.
func (s *Service) initDownload(c context.Context, v *view.View, mid int64, cdnIP string) (err error) {
	var download int64
	if v.AttrVal(archive.AttrBitLimitArea) == archive.AttrYes {
		if download, err = s.ipLimit(c, mid, v.Aid, cdnIP); err != nil {
			return
		}
	} else {
		download = location.AllowDown
	}
	if download == location.ForbiddenDown {
		v.Rights.Download = int32(download)
		return
	}
	for _, p := range v.Pages {
		if p.From == "qq" {
			download = location.ForbiddenDown
			break
		}
	}
	v.Rights.Download = int32(download)
	return
}

// initContributions is.
func (s *Service) initContributions(c context.Context, v *view.View, isTW bool) {
	const _count = 5
	if v.ReqUser != nil && v.ReqUser.Attention == 1 {
		return
	}
	var hasAudio bool
	for _, page := range v.Pages {
		if page.Audio != nil {
			hasAudio = true
			break
		}
	}
	if hasAudio || v.OwnerExt.Archives < 20 {
		return
	}
	aids, err := s.arcDao.ViewContributeCache(c, v.Author.Mid)
	if err != nil {
		log.Error("%+v", err)
		return
	}
	if len(aids) < _count+1 {
		return
	}
	as, err := s.arcDao.Archives(c, aids)
	if err != nil {
		log.Error("%+v", err)
		return
	}
	if len(as) == 0 {
		return
	}
	ctbt := make([]*view.Contribution, 0, len(as))
	for _, aid := range aids {
		if a, ok := as[aid]; ok && a.IsNormal() {
			if a.Aid != v.Aid {
				if isTW {
					a.Title = chinese.Convert(c, a.Title)
				}
				vc := &view.Contribution{Aid: a.Aid, Title: a.Title, Pic: a.Pic, Author: a.Author, Stat: a.Stat, CTime: a.PubDate}
				ctbt = append(ctbt, vc)
			}
		}
	}
	if len(ctbt) > _count {
		ctbt = ctbt[:_count]
	}
	if len(ctbt) == _count {
		v.Contributions = ctbt
	}
}

// initAudios is.
func (s *Service) initAudios(c context.Context, v *view.View) {
	pLen := len(v.Pages)
	if pLen == 0 || pLen > 100 {
		return
	}
	if pLen > 50 {
		pLen = 50
	}
	cids := make([]int64, 0, len(v.Pages[:pLen]))
	for _, p := range v.Pages[:pLen] {
		cids = append(cids, p.Cid)
	}
	vam, err := s.audioDao.AudioByCids(c, cids)
	if err != nil {
		log.Error("%+v", err)
		return
	}
	if len(vam) != 0 {
		for _, p := range v.Pages[:pLen] {
			if va, ok := vam[p.Cid]; ok {
				p.Audio = va
			}
		}
		if len(v.Pages) == 1 {
			if va, ok := vam[v.Pages[0].Cid]; ok {
				v.Audio = va
			}
		}
	}
}

// initTag is.
func (s *Service) initTag(c context.Context, v *view.View, mid int64) (tids []int64) {
	tags, err := s.tagDao.ArcTags(c, v.Aid, mid)
	if err != nil {
		log.Error("%+v", err)
		return
	}
	tids = make([]int64, 0, len(tags))
	for _, tag := range tags {
		tids = append(tids, tag.ID)
	}
	v.Tag = tags
	return
}

// initDM is.
func (s *Service) initDM(c context.Context, v *view.View) {
	const (
		_dmTypeAv    = 1
		_dmPlatMobie = 1
	)
	pLen := len(v.Pages)
	if pLen == 0 || pLen > 100 {
		return
	}
	if pLen > 50 {
		pLen = 50
	}
	cids := make([]int64, 0, len(v.Pages[:pLen]))
	for _, p := range v.Pages[:pLen] {
		cids = append(cids, p.Cid)
	}
	res, err := s.dmDao.SubjectInfos(c, _dmTypeAv, _dmPlatMobie, cids...)
	if err != nil {
		log.Error("%+v", err)
		return
	}
	if len(res) == 0 {
		return
	}
	for _, p := range v.Pages[:pLen] {
		if r, ok := res[p.Cid]; ok {
			p.DM = r
		}
	}
}

// dealRcmdRelate is.
func (s *Service) dealRcmdRelate(c context.Context, plat int8, aid int64) (rls []*view.Relate, err error) {
	if rls, err = s.arcDao.RelatesCache(c, aid); err != nil {
		return
	}
	if len(rls) != 0 {
		return
	}
	s.prom.Incr("need_relates")
	var aids []int64
	if aids, err = s.arcDao.RelateAids(c, aid); err != nil {
		return
	}
	if len(aids) == 0 {
		return
	}
	var as map[int64]*api.Arc
	if as, err = s.arcDao.Archives(c, aids); err != nil {
		return
	}
	for _, aid := range aids {
		if a, ok := as[aid]; ok {
			if s.overseaCheck(archive.BuildArchive3(a), plat) || !a.IsNormal() {
				continue
			}
			r := &view.Relate{}
			r.FromAv(a, "", "", nil)
			rls = append(rls, r)
		}
	}
	if len(rls) != 0 {
		// 如果是繁体区会修改relate的title，addcache是异步操作，需要深度拷贝，避免并发读写的panci
		rels := make([]*view.Relate, 0, len(rls))
		for _, rl := range rls {
			r := &view.Relate{}
			*r = *rl
			rels = append(rels, r)
		}
		s.arcDao.AddRelatesCache(aid, rels)
	}
	return
}

// dealManagerRelate is.
func (s *Service) dealManagerRelate(c context.Context, plat int8, mr *manager.Relate) (r *view.Relate, err error) {
	if mr == nil || mr.Param < 1 {
		return
	}
	switch mr.Goto {
	case model.GotoAv:
		var a *api.Arc
		if a, err = s.arcDao.Archive(c, mr.Param); err != nil {
			return
		}
		if a != nil {
			r = &view.Relate{}
			r.FromOperate(mr, a, model.FromOperation)
		}
	}
	return
}

// newRcmdRelate is.
func (s *Service) newRcmdRelate(c context.Context, plat int8, aid, mid int64, buvid, mobiApp, from string, build int) (rls []*view.Relate, userFeature, returnCode string, err error) {
	recData, userFeature, returnCode, err := s.arcDao.NewRelateAids(c, aid, mid, build, buvid, from, plat)
	if err != nil || len(recData) == 0 {
		return
	}
	var (
		aids []int64
		arcm map[int64]*api.Arc
	)
	for _, rec := range recData {
		switch rec.Goto {
		case model.GotoAv:
			aids = append(aids, rec.Oid)
		}
	}
	if len(aids) == 0 {
		return
	}
	if arcm, err = s.arcDao.Archives(c, aids); err != nil {
		log.Error("%+v", err)
		return
	}
	if len(arcm) == 0 {
		return
	}
	for _, rec := range recData {
		switch rec.Goto {
		case model.GotoAv:
			arc, ok := arcm[rec.Oid]
			if !ok || s.overseaCheck(archive.BuildArchive3(arc), plat) || !arc.IsNormal() {
				continue
			}
			r := &view.Relate{AvFeature: rec.AvFeature, Source: rec.Source}
			r.FromAv(arc, "", rec.TrackID, nil)
			rls = append(rls, r)
		}
	}
	return
}
