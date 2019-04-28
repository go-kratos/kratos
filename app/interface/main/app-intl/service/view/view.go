package view

import (
	"context"
	"time"

	"go-common/app/interface/main/app-intl/model"
	"go-common/app/interface/main/app-intl/model/view"
	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"
	"go-common/library/text/translate/chinese"
)

const (
	// _descLen is.
	_descLen = 250
	// _avTypeAv is.
	_avTypeAv = 1
	// _businessLike is.
	_businessLike = "archive"
)

// View all view data.
func (s *Service) View(c context.Context, mid, aid, movieID int64, plat int8, build int, ak, mobiApp, device, buvid, cdnIP, network, adExtra, from string, now time.Time, locale string) (v *view.View, err error) {
	if v, err = s.ViewPage(c, mid, aid, movieID, plat, build, ak, mobiApp, device, cdnIP, true, now, locale); err != nil {
		ip := metadata.String(c, metadata.RemoteIP)
		if err == ecode.AccessDenied || err == ecode.NothingFound {
			log.Warn("s.ViewPage() mid(%d) aid(%d) movieID(%d) plat(%d) ak(%s) ip(%s) cdn_ip(%s) error(%v)", mid, aid, movieID, plat, ak, ip, cdnIP, err)
		} else {
			log.Error("s.ViewPage() mid(%d) aid(%d) movieID(%d) plat(%d) ak(%s) ip(%s) cdn_ip(%s) error(%v)", mid, aid, movieID, plat, ak, ip, cdnIP, err)
		}
		return
	}
	isTW := model.TWLocale(locale)
	g, ctx := errgroup.WithContext(c)
	g.Go(func() (err error) {
		v.VIPActive = s.vipActiveCache[view.VIPActiveView]
		return
	})
	g.Go(func() (err error) {
		s.initReqUser(ctx, v, mid, plat, build)
		if v.AttrVal(archive.AttrBitIsPGC) != archive.AttrYes {
			s.initContributions(ctx, v, isTW)
		}
		return
	})
	g.Go(func() (err error) {
		s.initRelateCMTag(ctx, v, plat, build, mid, buvid, mobiApp, device, network, adExtra, from, now, isTW)
		return
	})
	if v.AttrVal(archive.AttrBitIsPGC) != archive.AttrYes {
		g.Go(func() (err error) {
			s.initDM(ctx, v)
			return
		})
		g.Go(func() (err error) {
			s.initAudios(ctx, v)
			return
		})
		if len([]rune(v.Desc)) > _descLen {
			g.Go(func() (err error) {
				if desc, _ := s.arcDao.Description(ctx, v.Aid); desc != "" {
					v.Desc = desc
				}
				return
			})
		}
	}
	if err = g.Wait(); err != nil {
		log.Error("%+v", err)
	}
	return
}

// ViewPage view page data.
func (s *Service) ViewPage(c context.Context, mid, aid, movieID int64, plat int8, build int, ak, mobiApp, device, cdnIP string, nMovie bool, now time.Time, locale string) (v *view.View, err error) {
	if aid == 0 && movieID == 0 {
		err = ecode.NothingFound
		return
	}
	var (
		vs         *view.ViewStatic
		vp         *archive.View3
		seasoninfo map[int64]int64
		ok         bool
	)
	if movieID != 0 {
		if seasoninfo, err = s.banDao.SeasonidAid(c, movieID, now); err != nil {
			log.Error("%+v", err)
			err = ecode.NothingFound
			return
		}
		if aid, ok = seasoninfo[movieID]; !ok || aid == 0 {
			err = ecode.NothingFound
			return
		}
		var a *api.Arc
		if a, err = s.arcDao.Archive3(c, aid); err != nil {
			log.Error("%+v", err)
			err = ecode.NothingFound
			return
		}
		if a == nil {
			err = ecode.NothingFound
			return
		}
		vs = &view.ViewStatic{Archive3: archive.BuildArchive3(a)}
		s.prom.Incr("from_movieID")
	} else {
		if vp, err = s.arcDao.ViewCache(c, aid); err != nil {
			log.Error("%+v", err)
		}
		if vp == nil || vp.Archive3 == nil || len(vp.Pages) == 0 || vp.AttrVal(archive.AttrBitIsMovie) == archive.AttrYes {
			if vp, err = s.arcDao.View3(c, aid); err != nil {
				log.Error("%+v", err)
				err = ecode.NothingFound
				return
			}
		}
		if vp == nil || vp.Archive3 == nil || len(vp.Pages) == 0 {
			err = ecode.NothingFound
			return
		}
		vs = &view.ViewStatic{Archive3: vp.Archive3}
		s.initPages(c, vs, vp.Pages)
		s.prom.Incr("from_aid")
	}
	isTW := model.TWLocale(locale)
	if isTW {
		out := chinese.Converts(c, vs.Title, vs.Desc)
		vs.Title = out[vs.Title]
		vs.Desc = out[vs.Desc]
	}
	// TODO fuck chanpin
	vs.Stat.DisLike = 0
	if s.overseaCheck(vs.Archive3, plat) {
		err = ecode.AreaLimit
		return
	}
	// check region area limit
	if err = s.areaLimit(c, plat, int16(vs.TypeID)); err != nil {
		return
	}
	v = &view.View{ViewStatic: vs, DMSeg: 1, PlayerIcon: s.playerIcon}
	if v.AttrVal(archive.AttrBitIsPGC) != archive.AttrYes {
		// check access
		if err = s.checkAceess(c, mid, v.Aid, int(v.State), int(v.Access), ak); err != nil {
			// archive is ForbitFixed and Transcoding and StateForbitDistributing need analysis history body .
			if v.State != archive.StateForbidFixed {
				return
			}
			err = nil
		}
		if v.Access > 0 {
			v.Stat.View = 0
		}
	}
	g, ctx := errgroup.WithContext(c)
	if mid != 0 {
		g.Go(func() (err error) {
			v.History, _ = s.arcDao.Progress(ctx, v.Aid, mid)
			return
		})
	}
	if v.AttrVal(archive.AttrBitIsPGC) == archive.AttrYes {
		if v.AttrVal(archive.AttrBitIsMovie) != archive.AttrYes {
			g.Go(func() error {
				return s.initPGC(ctx, v, mid, build, mobiApp, device)
			})
		} else {
			g.Go(func() error {
				return s.initMovie(ctx, v, mid, build, mobiApp, device, nMovie)
			})
		}
	} else {
		g.Go(func() (err error) {
			if err = s.initDownload(ctx, v, mid, cdnIP); err != nil {
				ip := metadata.String(ctx, metadata.RemoteIP)
				log.Error("aid(%d) mid(%d) ip(%s) cdn_ip(%s) error(%+v)", v.Aid, mid, ip, cdnIP, err)
			}
			return
		})
	}
	if err = g.Wait(); err != nil {
		log.Error("%+v", err)
	}
	return
}
