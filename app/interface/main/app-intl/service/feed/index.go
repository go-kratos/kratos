package feed

import (
	"context"
	"encoding/json"
	"time"

	cdm "go-common/app/interface/main/app-card/model"
	"go-common/app/interface/main/app-card/model/card"
	"go-common/app/interface/main/app-card/model/card/ai"
	"go-common/app/interface/main/app-card/model/card/banner"
	"go-common/app/interface/main/app-card/model/card/cm"
	"go-common/app/interface/main/app-card/model/card/operate"
	"go-common/app/interface/main/app-intl/model"
	"go-common/app/interface/main/app-intl/model/feed"
	tag "go-common/app/interface/main/tag/model"
	account "go-common/app/service/main/account/model"
	"go-common/app/service/main/archive/model/archive"
	locmdl "go-common/app/service/main/location/model"
	relation "go-common/app/service/main/relation/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"
	"go-common/library/text/translate/chinese"
)

// Index is.
func (s *Service) Index(c context.Context, buvid string, mid int64, plat int8, param *feed.IndexParam, now time.Time, style int) (is []card.Handler, userFeature json.RawMessage, isRcmd, newUser bool, code int, autoPlay, clean int8, autoPlayInfoc string, err error) {
	var (
		rs        []*ai.Item
		adm       map[int]*cm.AdInfo
		adAidm    map[int64]struct{}
		banners   []*banner.Banner
		version   string
		blackAidm map[int64]struct{}
		adInfom   map[int]*cm.AdInfo
		follow    *operate.Card
		ip        = metadata.String(c, metadata.RemoteIP)
		info      *locmdl.Info
		isTW      = model.TWLocale(param.Locale)
	)
	// 国际版不做abtest
	clean = 0
	autoPlay = 2
	group := s.group(mid, buvid)
	if info, err = s.loc.Info(c, ip); err != nil {
		log.Warn("s.loc.Info(%v) error(%v)", ip, err)
		err = nil
	}
	if !s.c.Feed.Index.Abnormal {
		g, ctx := errgroup.WithContext(c)
		g.Go(func() error {
			rs, userFeature, isRcmd, newUser, code = s.indexRcmd(ctx, plat, param.Build, buvid, mid, group, param.LoginEvent, info, param.Interest, param.Network, style, param.Column, param.Flush, autoPlayInfoc, now)
			if isTW {
				for _, r := range rs {
					if r.RcmdReason != nil {
						r.RcmdReason.Content = chinese.Convert(ctx, r.RcmdReason.Content)
					}
				}
			}
			return nil
		})
		g.Go(func() (err error) {
			if blackAidm, err = s.BlackList(ctx, mid); err != nil {
				log.Error("%+v", err)
				err = nil
			}
			return
		})
		if err = g.Wait(); err != nil {
			return
		}
		rs, adInfom = s.mergeItem(c, mid, rs, adm, adAidm, banners, version, blackAidm, plat, follow)
	} else {
		count := s.indexCount(plat)
		rs = s.recommendCache(count)
		log.Warn("feed index show disaster recovery data len(%d)", len(is))
	}
	is, isRcmd, err = s.dealItem(c, param.Column, mid, buvid, plat, param.Build, rs, isRcmd, param.MobiApp, param.Device, param.Network, param.OpenEvent, param.AdExtra, param.Qn, param.Fnver, param.Fnval, follow, isTW, now)
	s.dealAdLoc(is, param, adInfom, now)
	return
}

// indexRcmd is.
func (s *Service) indexRcmd(c context.Context, plat int8, build int, buvid string, mid int64, group int, loginEvent int, info *locmdl.Info, interest, network string, style int, column cdm.ColumnStatus, flush int, autoPlay string, now time.Time) (is []*ai.Item, userFeature json.RawMessage, isRcmd, newUser bool, code int) {
	count := s.indexCount(plat)
	if buvid != "" || mid != 0 {
		var (
			err    error
			zoneID int64
		)
		if info != nil {
			zoneID = info.ZoneID
		}
		if is, userFeature, code, newUser, err = s.rcmd.Recommend(c, plat, buvid, mid, build, loginEvent, zoneID, group, interest, network, style, column, flush, autoPlay, now); err != nil {
			log.Error("%+v", err)
		} else if len(is) != 0 {
			isRcmd = true
		}
		var fromCache bool
		if len(is) == 0 && mid != 0 && !ecode.ServiceUnavailable.Equal(err) {
			if is, err = s.indexCache(c, mid, count); err != nil {
				log.Error("%+v", err)
			}
			if len(is) != 0 {
				s.pHit.Incr("index_cache")
			} else {
				s.pMiss.Incr("index_cache")
			}
			fromCache = true
		}
		if len(is) == 0 || (fromCache && len(is) < count) {
			is = s.recommendCache(count)
		}
	} else {
		is = s.recommendCache(count)
	}
	return
}

// mergeItem is.
func (s *Service) mergeItem(c context.Context, mid int64, rs []*ai.Item, adm map[int]*cm.AdInfo, adAidm map[int64]struct{}, banners []*banner.Banner, version string, blackAids map[int64]struct{}, plat int8, follow *operate.Card) (is []*ai.Item, adInfom map[int]*cm.AdInfo) {
	if len(rs) == 0 {
		return
	}
	if len(banners) != 0 {
		rs = append([]*ai.Item{{Goto: model.GotoBanner, Banners: banners, Version: version}}, rs...)
	}
	is = make([]*ai.Item, 0, len(rs)+len(adm))
	adInfom = make(map[int]*cm.AdInfo, len(adm))
	for _, r := range rs {
		if r.Goto == model.GotoAv {
			if _, ok := blackAids[r.ID]; ok {
				continue
			} else if _, ok := s.blackCache[r.ID]; ok {
				continue
			}
			if _, ok := adAidm[r.ID]; ok {
				continue
			}
		} else if r.Goto == model.GotoBanner && len(is) != 0 {
			// banner 必须在第一位
			continue
		} else if r.Goto == model.GotoLogin && mid != 0 {
			continue
		}
		is = append(is, r)
	}
	return
}

// dealAdLoc is.
func (*Service) dealAdLoc(is []card.Handler, param *feed.IndexParam, adInfom map[int]*cm.AdInfo, now time.Time) {
	il := len(is)
	if il == 0 {
		return
	}
	if param.Idx < 1 {
		param.Idx = now.Unix()
	}
	for i, h := range is {
		if param.Pull {
			h.Get().Idx = param.Idx + int64(il-i)
		} else {
			h.Get().Idx = param.Idx - int64(i+1)
		}
		if ad, ok := adInfom[i]; ok {
			h.Get().AdInfo = ad
		} else if h.Get().AdInfo != nil {
			h.Get().AdInfo.CardIndex = i
		}
	}
}

// dealItem is.
func (s *Service) dealItem(c context.Context, column cdm.ColumnStatus, mid int64, buvid string, plat int8, build int, rs []*ai.Item, isRcmd bool, mobiApp, device, network, openEvent, adExtra string, qn, fnver, fnval int, follow *operate.Card, isTW bool, now time.Time) (is []card.Handler, isAI bool, err error) {
	if len(rs) == 0 {
		is = []card.Handler{}
		return
	}
	var (
		aids, tids                       []int64
		upIDs, avUpIDs, rmUpIDs, mtUpIDs []int64
		am                               map[int64]*archive.ArchiveWithPlayer
		tagm                             map[int64]*tag.Tag
		rank                             *operate.Card
		cardm                            map[int64]*account.Card
		statm                            map[int64]*relation.Stat
		isAtten                          map[int64]int8
		arcOK                            bool
	)
	followm := map[int64]*operate.Card{}
	isAI = isRcmd
	for _, r := range rs {
		if r == nil {
			continue
		}
		if isTW && r.RcmdReason != nil {
			r.RcmdReason.Content = chinese.Convert(c, r.RcmdReason.Content)
		}
		switch r.Goto {
		case model.GotoAv, model.GotoPlayer, model.GotoUpRcmdAv:
			if r.ID != 0 {
				aids = append(aids, r.ID)
			}
			if r.Tid != 0 {
				tids = append(tids, r.Tid)
			}
		case model.GotoRank:
			os, aid := s.RankCard(plat)
			rank = &operate.Card{}
			rank.FromRank(os)
			aids = append(aids, aid...)
		case model.GotoChannelRcmd:
			cardm, aid, tid := s.channelRcmdCard(c, r.ID)
			for id, card := range cardm {
				followm[id] = card
			}
			aids = append(aids, aid...)
			tids = append(tids, tid...)
		}
	}
	g, ctx := errgroup.WithContext(c)
	if len(aids) != 0 {
		g.Go(func() (err error) {
			if am, err = s.ArchivesWithPlayer(ctx, aids, qn, mobiApp, fnver, fnval); err != nil {
				return
			}
			arcOK = true
			for _, a := range am {
				avUpIDs = append(avUpIDs, a.Author.Mid)
				if isTW {
					out := chinese.Converts(ctx, a.Title, a.Desc, a.TypeName)
					a.Title = out[a.Title]
					a.Desc = out[a.Desc]
					a.TypeName = out[a.TypeName]
				}
			}
			return
		})
	}
	if len(tids) != 0 {
		g.Go(func() (err error) {
			if tagm, err = s.tg.InfoByIDs(ctx, mid, tids); err != nil {
				log.Error("%+v", err)
				err = nil
			}
			return
		})
	}
	if err = g.Wait(); err != nil {
		log.Error("%+v", err)
		if isRcmd {
			count := s.indexCount(plat)
			rs = s.recommendCache(count)
		}
	} else {
		upIDs = append(upIDs, avUpIDs...)
		upIDs = append(upIDs, rmUpIDs...)
		upIDs = append(upIDs, mtUpIDs...)
		g, ctx = errgroup.WithContext(c)
		if len(upIDs) != 0 {
			g.Go(func() (err error) {
				if cardm, err = s.acc.Cards3(ctx, upIDs); err != nil {
					log.Error("%+v", err)
					err = nil
				}
				return
			})
			g.Go(func() (err error) {
				if statm, err = s.rel.Stats(ctx, upIDs); err != nil {
					log.Error("%+v", err)
					err = nil
				}
				return
			})
			if mid != 0 {
				g.Go(func() error {
					isAtten = s.acc.IsAttention(ctx, upIDs, mid)
					return nil
				})
			}
		}
		g.Wait()
	}
	isAI = isAI && arcOK
	var cardTotal int
	is = make([]card.Handler, 0, len(rs))
	for _, r := range rs {
		if r == nil {
			continue
		}
		var (
			main     interface{}
			cardType cdm.CardType
		)
		op := &operate.Card{}
		op.From(cdm.CardGt(r.Goto), r.ID, r.Tid, plat, build)
		h := card.Handle(plat, cdm.CardGt(r.Goto), cardType, column, r, tagm, isAtten, statm, cardm)
		if h == nil {
			continue
		}
		switch r.Goto {
		case model.GotoAv, model.GotoPlayer, model.GotoUpRcmdAv:
			if !arcOK {
				if r.Archive != nil {
					if isTW {
						out := chinese.Converts(c, r.Archive.Title, r.Archive.Desc, r.Archive.TypeName, r.Archive.Author.Name)
						r.Archive.Title = out[r.Archive.Title]
						r.Archive.Desc = out[r.Archive.Desc]
						r.Archive.TypeName = out[r.Archive.TypeName]
					}
					am = map[int64]*archive.ArchiveWithPlayer{r.Archive.Aid: {Archive3: r.Archive}}
				}
				if r.Tag != nil {
					tagm = map[int64]*tag.Tag{r.Tag.ID: r.Tag}
					op.Tid = r.Tag.ID
				}
			}
			if a, ok := am[r.ID]; ok && (a.AttrVal(archive.AttrBitOverseaLock) == 0 || !model.IsOverseas(plat)) {
				main = am
				op.TrackID = r.TrackID
			}
		case model.GotoRank:
			main = map[cdm.Gt]interface{}{cdm.GotoAv: am}
			op = rank
		case model.GotoChannelRcmd:
			main = am
		case model.GotoLogin:
			op.FromLogin(r.ID)
		default:
			log.Warn("unexpected goto(%s) %+v", r.Goto, r)
			continue
		}
		h.From(main, op)
		// 卡片不正常要continue
		if !h.Get().Right {
			continue
		}
		is, cardTotal = s.appendItem(plat, is, h, column, cardTotal)
	}
	// 双列末尾卡片去空窗
	if !model.IsIPad(plat) {
		if cdm.Columnm[column] == cdm.ColumnSvrDouble {
			is = is[:len(is)-cardTotal%2]
		}
	} else {
		// 复杂的ipad去空窗逻辑
		if cardTotal%4 == 3 {
			if is[len(is)-2].Get().CardLen == 2 {
				is = is[:len(is)-2]
			} else {
				is = is[:len(is)-3]
			}
		} else if cardTotal%4 == 2 {
			if is[len(is)-1].Get().CardLen == 2 {
				is = is[:len(is)-1]
			} else {
				is = is[:len(is)-2]
			}
		} else if cardTotal%4 == 1 {
			is = is[:len(is)-1]
		}
	}
	if len(is) == 0 {
		is = []card.Handler{}
		return
	}
	return
}

// appendItem is.
func (s *Service) appendItem(plat int8, rs []card.Handler, h card.Handler, column cdm.ColumnStatus, cardTotal int) (is []card.Handler, total int) {
	h.Get().ThreePointFrom()
	// 国际版暂不支持稿件反馈
	if h.Get().ThreePoint != nil {
		h.Get().ThreePoint.Feedbacks = nil
	}
	if !model.IsIPad(plat) {
		// 双列大小卡换位去空窗
		if cdm.Columnm[column] == cdm.ColumnSvrDouble {
			// 通栏卡
			if h.Get().CardLen == 0 {
				if cardTotal%2 == 1 {
					is = card.SwapTwoItem(rs, h)
				} else {
					is = append(rs, h)
				}
			} else {
				is = append(rs, h)
			}
		} else {
			is = append(rs, h)
		}
	} else {
		// ipad卡片不展示标签
		h.Get().DescButton = nil
		// ipad大小卡换位去空窗
		if h.Get().CardLen == 0 {
			// 通栏卡
			if cardTotal%4 == 3 {
				is = card.SwapFourItem(rs, h)
			} else if cardTotal%4 == 2 {
				is = card.SwapThreeItem(rs, h)
			} else if cardTotal%4 == 1 {
				is = card.SwapTwoItem(rs, h)
			} else {
				is = append(rs, h)
			}
		} else if h.Get().CardLen == 2 {
			// 半栏卡
			if cardTotal%4 == 3 {
				is = card.SwapTwoItem(rs, h)
			} else if cardTotal%4 == 2 {
				is = append(rs, h)
			} else if cardTotal%4 == 1 {
				is = card.SwapTwoItem(rs, h)
			} else {
				is = append(rs, h)
			}
		} else {
			is = append(rs, h)
		}
	}
	total = cardTotal + h.Get().CardLen
	return
}

// indexCount is.
func (s *Service) indexCount(plat int8) (count int) {
	if plat == model.PlatIPad {
		count = s.c.Feed.Index.IPadCount
	} else {
		count = s.c.Feed.Index.Count
	}
	return
}
