package feed

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	cdm "go-common/app/interface/main/app-card/model"
	"go-common/app/interface/main/app-card/model/card/ai"
	"go-common/app/interface/main/app-card/model/card/audio"
	"go-common/app/interface/main/app-card/model/card/bangumi"
	"go-common/app/interface/main/app-card/model/card/banner"
	"go-common/app/interface/main/app-card/model/card/cm"
	"go-common/app/interface/main/app-card/model/card/live"
	"go-common/app/interface/main/app-card/model/card/operate"
	"go-common/app/interface/main/app-card/model/card/rank"
	"go-common/app/interface/main/app-card/model/card/show"
	"go-common/app/interface/main/app-feed/model"
	"go-common/app/interface/main/app-feed/model/feed"
	bustag "go-common/app/interface/main/tag/model"
	article "go-common/app/interface/openplatform/article/model"
	account "go-common/app/service/main/account/model"
	"go-common/app/service/main/archive/model/archive"
	locmdl "go-common/app/service/main/location/model"
	relation "go-common/app/service/main/relation/model"
	episodegrpc "go-common/app/service/openplatform/pgc-season/api/grpc/episode/v1"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"
)

const (
	_cardAdAv    = 1
	_cardAdWeb   = 2
	_cardAdWebS  = 3
	_cardAdLarge = 7
	_feedgroups  = "tianma2.0_autoplay_card"
)

// Index is
func (s *Service) Index(c context.Context, mid int64, plat int8, build int, buvid, network, mobiApp, device, platform, openEvent string, loginEvent int, idx int64, pull bool, now time.Time, bannerHash, adExtra string, qn int, interest string, style, flush, fnver, fnval, autoplayCard int) (is []*feed.Item, userFeature json.RawMessage, isRcmd, newUser bool, code, clean int, autoPlayInfoc string, err error) {
	var (
		ris       []*ai.Item
		adm       map[int]*cm.AdInfo
		adAidm    map[int64]struct{}
		hasBanner bool
		bs        []*banner.Banner
		version   string
		blackAidm map[int64]struct{}
		adInfom   map[int]*cm.AdInfo
		follow    *operate.Follow
		autoPlay  int
		ip        = metadata.String(c, metadata.RemoteIP)
		info      *locmdl.Info
	)
	//abtest================
	// if mid > 0 && mid%20 == 19 {
	// 	clean = 1
	// } else {
	// 	clean = 0
	// }
	clean = 0
	if ab, ok := s.abtestCache[_feedgroups]; ok {
		if ab.AbTestIn(buvid + _feedgroups) {
			switch autoplayCard {
			case 0, 1, 2, 3:
				autoPlay = 1
			default:
				autoPlay = 2
			}
		} else {
			autoPlay = 2
		}
	} else {
		switch autoplayCard {
		case 1, 3:
			autoPlay = 1
		default:
			autoPlay = 2
		}
	}
	autoPlayInfoc = fmt.Sprintf("%d|%d", autoPlay, autoplayCard)
	if info, err = s.loc.Info(c, ip); err != nil {
		log.Warn("s.loc.Info(%v) error(%v)", ip, err)
		err = nil
	}
	//abtest================
	group := s.group(mid, buvid)
	g, ctx := errgroup.WithContext(c)
	g.Go(func() error {
		ris, userFeature, isRcmd, newUser, code = s.indexRcmd(ctx, plat, build, buvid, mid, group, loginEvent, 0, info, interest, network, style, -1, flush, autoPlayInfoc, now)
		return nil
	})
	// 暂停实验
	// if !((group == 18 || group == 19) && style == 3) {
	g.Go(func() (err error) {
		if adm, adAidm, err = s.indexAd(ctx, plat, build, buvid, mid, network, mobiApp, device, openEvent, info, now, adExtra, style); err != nil {
			log.Error("%+v", err)
			err = nil
		}
		return
	})
	// }
	g.Go(func() (err error) {
		if hasBanner, bs, version, err = s.indexBanner(ctx, plat, build, buvid, mid, loginEvent, bannerHash, network, mobiApp, device, "", adExtra); err != nil {
			log.Error("%+v", err)
			err = nil
		}
		return
	})
	g.Go(func() (err error) {
		if blackAidm, err = s.BlackList(ctx, mid); err != nil {
			log.Error("%+v", err)
			err = nil
		}
		return
	})
	g.Go(func() (err error) {
		if follow, err = s.SearchFollow(ctx, platform, mobiApp, device, buvid, build, mid); err != nil {
			log.Error("%+v", err)
			err = nil
		}
		return
	})
	if err = g.Wait(); err != nil {
		return
	}
	ris, adInfom = s.mergeItem(c, mid, ris, adm, adAidm, hasBanner, blackAidm, plat, follow)
	is, isRcmd, err = s.dealItem(c, mid, plat, build, buvid, ris, bs, version, isRcmd, network, mobiApp, device, openEvent, idx, pull, qn, now, adExtra, adInfom, fnver, fnval, autoPlay, follow)
	return
}

// Dislike is.
func (s *Service) Dislike(c context.Context, mid, id int64, buvid, gt string, reasonID, cmreasonID, feedbackID, upperID, rid, tagID int64, adcb string, now time.Time) (err error) {
	if gt == model.GotoAv {
		s.blk.AddBlacklist(mid, id)
	}
	return s.rcmd.PubDislike(c, buvid, gt, id, mid, reasonID, cmreasonID, feedbackID, upperID, rid, tagID, adcb, now)
}

// DislikeCancel is.
func (s *Service) DislikeCancel(c context.Context, mid, id int64, buvid, gt string, reasonID, cmreasonID, feedbackID, upperID, rid, tagID int64, adcb string, now time.Time) (err error) {
	if gt == model.GotoAv {
		s.blk.DelBlacklist(mid, id)
	}
	return s.rcmd.PubDislikeCancel(c, buvid, gt, id, mid, reasonID, cmreasonID, feedbackID, upperID, rid, tagID, adcb, now)
}

func (s *Service) indexRcmd(c context.Context, plat int8, build int, buvid string, mid int64, group int, loginEvent, parentMode int, zone *locmdl.Info, interest, network string, style int, column cdm.ColumnStatus, flush int, autoPlay string, now time.Time) (is []*ai.Item, userFeature json.RawMessage, isRcmd, newUser bool, code int) {
	count := s.indexCount(plat)
	if buvid != "" || mid != 0 {
		var (
			err    error
			zoneID int64
		)
		if zone != nil {
			zoneID = zone.ZoneID
		}
		if is, userFeature, code, newUser, err = s.rcmd.Recommend(c, plat, buvid, mid, build, loginEvent, parentMode, 0, zoneID, group, interest, network, style, column, flush, autoPlay, now); err != nil {
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

func (s *Service) indexAd(c context.Context, plat int8, build int, buvid string, mid int64, network, mobiApp, device, openEvent string, zone *locmdl.Info, now time.Time, adExtra string, style int) (adm map[int]*cm.AdInfo, adAidm map[int64]struct{}, err error) {
	var advert *cm.Ad
	resource := s.adResource(plat, build)
	if resource == 0 {
		return
	}
	//  兼容老的style逻辑，3为新单列，上报给商业产品的参数定义为：1 单列 2双列
	if style == 3 {
		style = 1
	}
	var country, province, city string
	if zone != nil {
		country = zone.Country
		province = zone.Province
		city = zone.City
	}
	if advert, err = s.ad.Ad(c, mid, build, buvid, []int64{resource}, country, province, city, network, mobiApp, device, openEvent, adExtra, style, now); err != nil {
		return
	}
	if advert == nil || len(advert.AdsInfo) == 0 {
		return
	}
	if adsInfo, ok := advert.AdsInfo[resource]; ok {
		adm = make(map[int]*cm.AdInfo, len(adsInfo))
		adAidm = make(map[int64]struct{}, len(adsInfo))
		for source, info := range adsInfo {
			if info == nil {
				continue
			}
			var adInfo *cm.AdInfo
			if info.AdInfo != nil {
				adInfo = info.AdInfo
				adInfo.RequestID = advert.RequestID
				adInfo.Resource = resource
				adInfo.Source = source
				adInfo.IsAd = info.IsAd
				adInfo.IsAdLoc = true
				adInfo.CmMark = info.CmMark
				adInfo.Index = info.Index
				adInfo.CardIndex = info.CardIndex
				adInfo.ClientIP = advert.ClientIP
				if adInfo.CreativeID != 0 && adInfo.CardType == _cardAdAv {
					adAidm[adInfo.CreativeContent.VideoID] = struct{}{}
				}
			} else {
				adInfo = &cm.AdInfo{RequestID: advert.RequestID, Resource: resource, Source: source, IsAdLoc: true, IsAd: info.IsAd, CmMark: info.CmMark, Index: info.Index, CardIndex: info.CardIndex, ClientIP: advert.ClientIP}
			}
			adm[adInfo.CardIndex-1] = adInfo
		}
	}
	return
}

func (s *Service) indexBanner(c context.Context, plat int8, build int, buvid string, mid int64, loginEvent int, hash, network, mobiApp, device, openEvent, adExtra string) (has bool, bs []*banner.Banner, version string, err error) {
	const (
		_androidBanBannerHash = 515009
		_iphoneBanBannerHash  = 6120
		_ipadBanBannerHash    = 6160
	)
	if (plat == model.PlatAndroid && build > _androidBanBannerHash) || (plat == model.PlatIPhone && build > _iphoneBanBannerHash) || (plat == model.PlatIPad && build > _ipadBanBannerHash) || loginEvent != 0 {
		if bs, version, err = s.banners(c, plat, build, mid, buvid, network, mobiApp, device, openEvent, adExtra, ""); err != nil {
			return
		} else if loginEvent != 0 {
			has = true
		} else if version != "" {
			has = hash != version
		}
	}
	return
}

func (s *Service) mergeItem(c context.Context, mid int64, rs []*ai.Item, adm map[int]*cm.AdInfo, adAidm map[int64]struct{}, hasBanner bool, blackAids map[int64]struct{}, plat int8, follow *operate.Follow) (is []*ai.Item, adInfom map[int]*cm.AdInfo) {
	if len(rs) == 0 {
		return
	}
	const (
		cardIndex     = 7
		cardIndexIPad = 17
		cardOffset    = 2
	)
	if hasBanner {
		rs = append([]*ai.Item{&ai.Item{Goto: model.GotoBanner}}, rs...)
		for index, ad := range adm {
			if ((model.IsIPad(plat) && index <= cardIndexIPad) || index <= cardIndex) && (ad.CardType == _cardAdWeb || ad.CardType == _cardAdLarge) {
				ad.CardIndex = ad.CardIndex + cardOffset
			}
		}
	}
	if follow != nil {
		followPos := s.c.Feed.Index.FollowPosition
		if followPos-1 >= 0 && followPos-1 <= len(rs) {
			rs = append(rs[:followPos-1], append([]*ai.Item{&ai.Item{ID: follow.ID, Goto: model.GotoSearchSubscribe}}, rs[followPos-1:]...)...)
		}
	}
	is = make([]*ai.Item, 0, len(rs)+len(adm))
	adInfom = make(map[int]*cm.AdInfo, len(adm))
	var existsBanner, existsAdWeb bool
	for _, r := range rs {
		for {
			if ad, ok := adm[len(is)]; ok {
				if ad.CreativeID != 0 {
					var item *ai.Item
					if ad.CardType == _cardAdAv {
						item = &ai.Item{ID: ad.CreativeContent.VideoID, Goto: model.GotoAdAv, Ad: ad}
					} else if ad.CardType == _cardAdWeb {
						item = &ai.Item{Goto: model.GotoAdWeb, Ad: ad}
						existsAdWeb = true
					} else if ad.CardType == _cardAdWebS {
						item = &ai.Item{Goto: model.GotoAdWebS, Ad: ad}
					} else if ad.CardType == _cardAdLarge {
						item = &ai.Item{Goto: model.GotoAdLarge, Ad: ad}
					} else {
						b, _ := json.Marshal(ad)
						log.Error("ad---%s", b)
						break
					}
					is = append(is, item)
					continue
				} else {
					adInfom[len(is)] = ad
				}
			}
			break
		}
		if r.Goto == model.GotoAv {
			if _, ok := blackAids[r.ID]; ok {
				continue
			} else if _, ok := s.blackCache[r.ID]; ok {
				continue
			}
			if _, ok := adAidm[r.ID]; ok {
				continue
			}
		} else if r.Goto == model.GotoBanner {
			if existsBanner {
				continue
			} else {
				existsBanner = true
			}
		} else if r.Goto == model.GotoRank && existsAdWeb {
			continue
		} else if r.Goto == model.GotoLogin && mid != 0 {
			continue
		}
		is = append(is, r)
	}
	return
}

func (s *Service) dealItem(c context.Context, mid int64, plat int8, build int, buvid string, rs []*ai.Item, bs []*banner.Banner, version string, isRcmd bool, network, mobiApp, device, openEvent string, idx int64, pull bool, qn int, now time.Time, adExtra string, adInfom map[int]*cm.AdInfo, fnver, fnval, autoPlay int, follow *operate.Follow) (is []*feed.Item, isAI bool, err error) {
	if len(rs) == 0 {
		is = _emptyItem
		return
	}
	var (
		aids, tids, roomIDs, sids, metaIDs, shopIDs, audioIDs []int64
		upIDs, avUpIDs, rmUpIDs, mtUpIDs                      []int64
		seasonIDs                                             []int32
		ranks                                                 []*rank.Rank
		am                                                    map[int64]*archive.ArchiveWithPlayer
		tagm                                                  map[int64]*bustag.Tag
		follows                                               map[int64]bool
		rm                                                    map[int64]*live.Room
		sm                                                    map[int64]*bangumi.Season
		hasBangumiRcmd                                        bool
		update                                                *bangumi.Update
		atm                                                   map[int64]*article.Meta
		scm                                                   map[int64]*show.Shopping
		aum                                                   map[int64]*audio.Audio
		hasBanner                                             bool
		card                                                  map[int64]*account.Card
		upStatm                                               map[int64]*relation.Stat
		arcOK                                                 bool
		seasonCards                                           map[int32]*episodegrpc.EpisodeCardsProto
	)
	isAI = isRcmd
	convergem := map[int64]*operate.Converge{}
	downloadm := map[int64]*operate.Download{}
	liveUpm := map[int64][]*live.Card{}
	followm := map[int64]*operate.Follow{}
	for _, r := range rs {
		switch r.Goto {
		case model.GotoAv, model.GotoAdAv, model.GotoPlayer, model.GotoUpRcmdAv:
			if r.ID != 0 {
				aids = append(aids, r.ID)
			}
			if r.Tid != 0 {
				tids = append(tids, r.Tid)
			}
		case model.GotoLive, model.GotoPlayerLive:
			if r.ID != 0 {
				roomIDs = append(roomIDs, r.ID)
			}
		case model.GotoBangumi:
			if r.ID != 0 {
				sids = append(sids, r.ID)
			}
		case model.GotoPGC:
			if r.ID != 0 {
				seasonIDs = append(seasonIDs, int32(r.ID))
			}
		case model.GotoRank:
			card, aid := s.RankCard(model.PlatIPhone)
			ranks = card
			aids = append(aids, aid...)
		case model.GotoBangumiRcmd:
			hasBangumiRcmd = true
		case model.GotoBanner:
			hasBanner = true
		case model.GotoConverge:
			if card, ok := s.convergeCache[r.ID]; ok {
				for _, item := range card.Items {
					switch item.Goto {
					case model.GotoAv:
						if item.Pid != 0 {
							aids = append(aids, item.Pid)
						}
					case model.GotoLive:
						if item.Pid != 0 {
							roomIDs = append(roomIDs, item.Pid)
						}
					case model.GotoArticle:
						if item.Pid != 0 {
							metaIDs = append(metaIDs, item.Pid)
						}
					}
				}
				convergem[r.ID] = card
			}
		case model.GotoGameDownloadS:
			if card, ok := s.downloadCache[r.ID]; ok {
				downloadm[r.ID] = card
			}
		case model.GotoArticleS:
			if r.ID != 0 {
				metaIDs = append(metaIDs, r.ID)
			}
		case model.GotoShoppingS:
			if r.ID != 0 {
				shopIDs = append(shopIDs, r.ID)
			}
		case model.GotoAudio:
			if r.ID != 0 {
				audioIDs = append(audioIDs, r.ID)
			}
		case model.GotoLiveUpRcmd:
			if r.ID != 0 {
				if cs, ok := s.liveCardCache[r.ID]; ok {
					for _, c := range cs {
						upIDs = append(upIDs, c.UID)
					}
				}
			}
		case model.GotoSubscribe:
			if r.ID != 0 {
				if card, ok := s.followCache[r.ID]; ok {
					for _, item := range card.Items {
						switch item.Goto {
						case cdm.GotoMid:
							if item.Pid != 0 {
								upIDs = append(upIDs, item.Pid)
							}
						case cdm.GotoTag:
							if item.Pid != 0 {
								tids = append(tids, item.Pid)
							}
						}
					}
					followm[r.ID] = card
				}
			}
		case model.GotoChannelRcmd:
			if r.ID != 0 {
				if card, ok := s.followCache[r.ID]; ok {
					if card.Pid != 0 {
						aids = append(aids, card.Pid)
					}
					if card.Tid != 0 {
						tids = append(tids, card.Tid)
					}
					followm[r.ID] = card
				}
			}
		case model.GotoSearchSubscribe:
			if follow != nil {
				followm[follow.ID] = follow
				for _, item := range follow.Items {
					upIDs = append(upIDs, item.Pid)
				}
			}
		}
	}
	g, ctx := errgroup.WithContext(c)
	if len(aids) != 0 {
		g.Go(func() (err error) {
			if am, err = s.ArchivesWithPlayer(ctx, aids, qn, mobiApp, fnver, fnval, 0, build); err != nil {
				return
			}
			arcOK = true
			for _, a := range am {
				avUpIDs = append(avUpIDs, a.Author.Mid)
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
	if len(roomIDs) != 0 {
		g.Go(func() (err error) {
			if rm, err = s.lv.AppMRoom(ctx, roomIDs); err != nil {
				log.Error("%+v", err)
				err = nil
			}
			for _, r := range rm {
				rmUpIDs = append(rmUpIDs, r.UID)
			}
			return
		})
	}
	if len(sids) != 0 {
		g.Go(func() (err error) {
			if sm, err = s.bgm.Seasons(ctx, sids, now); err != nil {
				log.Error("%+v", err)
				err = nil
			}
			return
		})
	}
	if len(seasonIDs) != 0 {
		g.Go(func() (err error) {
			if seasonCards, err = s.bgm.CardsInfoReply(ctx, seasonIDs); err != nil {
				log.Error("%+v", err)
				err = nil
			}
			return
		})
	}
	// TODO DEL
	// if hasBangumiRcmd && mid != 0 {
	if hasBangumiRcmd {
		g.Go(func() (err error) {
			if update, err = s.bgm.Updates(ctx, mid, now); err != nil {
				log.Error("%+v", err)
				err = nil
			}
			return
		})
	}
	if hasBanner && version == "" {
		g.Go(func() (err error) {
			if bs, version, err = s.banners(ctx, plat, build, mid, buvid, network, mobiApp, device, openEvent, adExtra, ""); err != nil {
				log.Error("%+v", err)
				err = nil
			}
			return
		})
	}
	if len(metaIDs) != 0 {
		g.Go(func() (err error) {
			if atm, err = s.art.Articles(ctx, metaIDs); err != nil {
				log.Error("%+v", err)
				err = nil
			}
			for _, at := range atm {
				if at.Author != nil {
					mtUpIDs = append(mtUpIDs, at.Author.Mid)
				}
			}
			return
		})
	}
	if len(shopIDs) != 0 {
		g.Go(func() (err error) {
			if scm, err = s.show.Card(ctx, shopIDs); err != nil {
				log.Error("%+v", err)
				err = nil
			}
			return
		})
	}
	if len(audioIDs) != 0 {
		g.Go(func() (err error) {
			if aum, err = s.audio.Audios(ctx, audioIDs); err != nil {
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
				if card, err = s.acc.Cards3(ctx, upIDs); err != nil {
					log.Error("%+v", err)
					err = nil
				}
				return
			})
			g.Go(func() (err error) {
				if upStatm, err = s.rel.Stats(ctx, upIDs); err != nil {
					log.Error("%+v", err)
					err = nil
				}
				return
			})
			if mid != 0 {
				g.Go(func() error {
					follows = s.acc.Relations3(ctx, upIDs, mid)
					return nil
				})
			}
		}
		g.Wait()
	}
	isAI = isAI && arcOK
	// init feed items
	is = make([]*feed.Item, 0, len(rs))
	var (
		smallCardCnt  int
		middleCardCnt int
	)
	ip := metadata.String(c, metadata.RemoteIP)
	adm := map[int]*feed.Item{}
	isIpad := plat == model.PlatIPad
	for _, r := range rs {
		il := len(is)
		i := &feed.Item{AI: r}
		i.FromRcmd(r)
		switch r.Goto {
		case model.GotoAv, model.GotoUpRcmdAv:
			a, ok := am[r.ID]
			if !ok && !arcOK {
				a = &archive.ArchiveWithPlayer{Archive3: r.Archive}
			}
			isOsea := model.IsOverseas(plat)
			if a != nil && a.Archive3 != nil && a.IsNormal() && (!isOsea || (isOsea && a.AttrVal(archive.AttrBitOverseaLock) == 0)) {
				i.FromPlayerAv(a)
				if arcOK {
					if info, ok := tagm[r.Tid]; ok {
						i.Tag = &feed.Tag{TagID: info.ID, TagName: info.Name, IsAtten: info.IsAtten, Count: &feed.TagCount{Atten: info.Count.Atten}}
					}
				} else if r.Tag != nil {
					i.Tag = &feed.Tag{TagID: r.Tag.ID, TagName: r.Tag.Name}
				}
				i.FromDislikeReason(plat, build)
				i.FromRcmdReason(r.RcmdReason)
				if follows[i.Mid] {
					i.IsAtten = 1
				}
				if card, ok := card[i.Mid]; ok {
					if card.Official.Role != 0 {
						i.Official = &feed.OfficialInfo{Role: card.Official.Role, Title: card.Official.Title, Desc: card.Official.Desc}
					}
				}
				// for GotoUpRcmdAv
				i.Goto = r.Goto
				if i.Goto == model.GotoUpRcmdAv {
					// TODO 等待开启
					// percent := i.Like / (i.Like + i.Dislike) * 100
					// if percent != 0 {
					// 	i.Desc = strconv.Itoa(percent) + "%的人推荐"
					// }
					i.Desc = ""
				}
				is = append(is, i)
				smallCardCnt++
			}
		case model.GotoLive:
			if r, ok := rm[r.ID]; ok {
				i.FromLive(r)
				if card, ok := card[i.Mid]; ok {
					if card.Official.Role != 0 {
						i.Official = &feed.OfficialInfo{Role: card.Official.Role, Title: card.Official.Title, Desc: card.Official.Desc}
					}
				}
				if stat, ok := upStatm[i.Mid]; ok {
					i.Fans = stat.Follower
				}
				if follows[i.Mid] {
					i.IsAtten = 1
				}
				if i.Goto != "" {
					is = append(is, i)
					smallCardCnt++
				}
			}
		case model.GotoBangumi:
			if s, ok := sm[r.ID]; ok {
				i.FromSeason(s)
				is = append(is, i)
				smallCardCnt++
			}
		case model.GotoPGC:
			if s, ok := seasonCards[int32(r.ID)]; ok {
				i.FromPGCSeason(s)
				is = append(is, i)
				smallCardCnt++
			}
		case model.GotoLogin:
			i.FromLogin()
			is = append(is, i)
			smallCardCnt++
		case model.GotoAdAv:
			if r.Ad != nil {
				if a, ok := am[r.ID]; ok && model.AdAvIsNormal(a) {
					i.FromAdAv(r.Ad, a)
					if follows[i.Mid] {
						i.IsAtten = 1
					}
					if card, ok := card[i.Mid]; ok {
						if card.Official.Role != 0 {
							i.Official = &feed.OfficialInfo{Role: card.Official.Role, Title: card.Official.Title, Desc: card.Official.Desc}
						}
					}
					i.ClientIP = ip
					adm[i.CardIndex-1] = i
				}
			}
		case model.GotoAdWebS:
			if r.Ad != nil {
				i.FromAdWebS(r.Ad)
				i.ClientIP = ip
				adm[i.CardIndex-1] = i
			}
		case model.GotoAdWeb:
			if r.Ad != nil {
				i.FromAdWeb(r.Ad)
				i.ClientIP = ip
				adm[i.CardIndex-1] = i
			}
		case model.GotoAdLarge:
			if r.Ad != nil {
				i.FromAdLarge(r.Ad)
				i.ClientIP = ip
				adm[i.CardIndex-1] = i
			}
		case model.GotoSpecial:
			if sc, ok := s.specialCache[r.ID]; ok {
				i.FromSpecial(sc.ID, sc.Title, sc.Cover, sc.Desc, sc.ReValue, sc.ReType, sc.Badge, sc.Size)
			}
			if i.Goto != "" {
				if !isIpad {
					if smallCardCnt%2 != 0 {
						is = swapTwoItem(is, i)
					} else {
						is = append(is, i)
					}
				} else {
					if (smallCardCnt+middleCardCnt*2)%2 != 0 {
						is = swapTwoItem(is, i)
					} else {
						is = append(is, i)
					}
					middleCardCnt++
				}
			}
		case model.GotoSpecialS:
			if sc, ok := s.specialCache[r.ID]; ok {
				i.FromSpecialS(sc.ID, sc.Title, sc.Cover, sc.SingleCover, sc.Desc, sc.ReValue, sc.ReType, sc.Badge)
			}
			if i.Goto != "" {
				if !isIpad {
					is = append(is, i)
					smallCardCnt++
				}
			}
		case model.GotoRank:
			i.FromRank(ranks, am)
			if i.Goto != "" {
				if !isIpad {
					if smallCardCnt%2 != 0 {
						is = swapTwoItem(is, i)
					} else {
						is = append(is, i)
					}
				} else {
					if (smallCardCnt+middleCardCnt*2)%2 != 0 {
						is = swapTwoItem(is, i)
					} else {
						is = append(is, i)
					}
					middleCardCnt++
				}
			}
		case model.GotoBangumiRcmd:
			if mid != 0 && update != nil && update.Updates != 0 {
				i.FromBangumiRcmd(update)
				if !isIpad {
					if smallCardCnt%2 != 0 {
						is = swapTwoItem(is, i)
					} else {
						is = append(is, i)
					}
				} else {
					is = append(is, i)
					smallCardCnt++
				}
			}
		case model.GotoBanner:
			if len(bs) != 0 {
				i.FromBanner(bs, version)
				if !isIpad {
					if smallCardCnt%2 != 0 {
						is = swapTwoItem(is, i)
					} else {
						is = append(is, i)
					}
				} else {
					switch (smallCardCnt + middleCardCnt*2) % 4 {
					case 0:
						is = append(is, i)
					case 1:
						is = swapTwoItem(is, i)
					case 2:
						switch is[len(is)-1].Goto {
						case model.GotoRank, model.GotoAdWeb, model.GotoAdLarge:
							is = swapTwoItem(is, i)
						default:
							is = swapThreeItem(is, i)
						}
					case 3:
						is = swapThreeItem(is, i)
					}
				}
			}
		case model.GotoConverge:
			if cc, ok := convergem[r.ID]; ok {
				i.FromConverge(cc, am, rm, atm)
				if i.Goto != "" {
					if !isIpad {
						if smallCardCnt%2 != 0 {
							is = swapTwoItem(is, i)
						} else {
							is = append(is, i)
						}
					}
				}
			}
		case model.GotoGameDownloadS:
			if gd, ok := downloadm[r.ID]; ok {
				i.FromGameDownloadS(gd, plat, build)
				if i.Goto != "" {
					if !isIpad {
						is = append(is, i)
						smallCardCnt++
					}
				}
			}
		case model.GotoArticleS:
			if m, ok := atm[r.ID]; ok {
				i.FromArticleS(m)
				if card, ok := card[i.Mid]; ok {
					if card.Official.Role != 0 {
						i.Official = &feed.OfficialInfo{Role: card.Official.Role, Title: card.Official.Title, Desc: card.Official.Desc}
					}
				}
				if i.Goto != "" {
					if !isIpad {
						is = append(is, i)
						smallCardCnt++
					}
				}
			}
		case model.GotoShoppingS:
			if c, ok := scm[r.ID]; ok {
				i.FromShoppingS(c)
				if i.Goto != "" {
					if !isIpad {
						is = append(is, i)
						smallCardCnt++
					}
				}
			}
		case model.GotoAudio:
			if au, ok := aum[r.ID]; ok {
				i.FromAudio(au)
				is = append(is, i)
				smallCardCnt++
			}
		case model.GotoPlayer:
			if a, ok := am[r.ID]; ok {
				i.FromPlayer(a)
				if i.Goto != "" {
					if info, ok := tagm[r.Tid]; ok {
						i.Tag = &feed.Tag{TagID: info.ID, TagName: info.Name, IsAtten: info.IsAtten, Count: &feed.TagCount{Atten: info.Count.Atten}}
					}
					if follows[i.Mid] {
						i.IsAtten = 1
					}
					if card, ok := card[i.Mid]; ok {
						if card.Official.Role != 0 {
							i.Official = &feed.OfficialInfo{Role: card.Official.Role, Title: card.Official.Title, Desc: card.Official.Desc}
						}
					}
					i.FromDislikeReason(plat, build)
					if !isIpad {
						if smallCardCnt%2 != 0 {
							is = swapTwoItem(is, i)
						} else {
							is = append(is, i)
						}
					}
				}
			}
		case model.GotoPlayerLive:
			if r, ok := rm[r.ID]; ok {
				i.FromPlayerLive(r)
				if i.Goto != "" {
					if follows[i.Mid] {
						i.IsAtten = 1
					}
					if card, ok := card[i.Mid]; ok {
						if card.Official.Role != 0 {
							i.Official = &feed.OfficialInfo{Role: card.Official.Role, Title: card.Official.Title, Desc: card.Official.Desc}
						}
					}
					if stat, ok := upStatm[i.Mid]; ok {
						i.Fans = stat.Follower
					}
					if !isIpad {
						if smallCardCnt%2 != 0 {
							is = swapTwoItem(is, i)
						} else {
							is = append(is, i)
						}
					}
				}
			}
		case model.GotoSubscribe, model.GotoSearchSubscribe:
			if c, ok := followm[r.ID]; ok {
				if !isIpad {
					i.FromSubscribe(c, card, follows, upStatm, tagm)
					if i.Goto != "" {
						if smallCardCnt%2 != 0 {
							is = swapTwoItem(is, i)
						} else {
							is = append(is, i)
						}
					}
				}
			}
		case model.GotoChannelRcmd:
			if c, ok := followm[r.ID]; ok {
				if !isIpad {
					i.FromChannelRcmd(c, am, tagm)
					if i.Goto != "" {
						if !isIpad {
							is = append(is, i)
							smallCardCnt++
						}
					}
				}
			}
		case model.GotoLiveUpRcmd:
			if c, ok := liveUpm[r.ID]; ok {
				if !isIpad {
					i.FromLiveUpRcmd(r.ID, c, card)
					if i.Goto != "" {
						if smallCardCnt%2 != 0 {
							is = swapTwoItem(is, i)
						} else {
							is = append(is, i)
						}
					}
				}
			}
		default:
			log.Warn("unexpected goto(%s) %+v", r.Goto, r)
			continue
		}
		if ad, ok := adm[il]; ok {
			switch ad.Goto {
			case model.GotoAdAv, model.GotoAdWebS:
				is = append(is, ad)
				smallCardCnt++
			case model.GotoAdWeb, model.GotoAdLarge:
				if !isIpad {
					if smallCardCnt%2 != 0 {
						is = swapTwoItem(is, ad)
					} else {
						is = append(is, ad)
					}
				} else {
					if (smallCardCnt+middleCardCnt*2)%2 != 0 {
						is = swapTwoItem(is, ad)
					} else {
						is = append(is, ad)
					}
					middleCardCnt++
				}
			}
		}
	}
	if !isIpad {
		is = is[:len(is)-smallCardCnt%2]
	} else {
		switch (smallCardCnt + middleCardCnt*2) % 4 {
		case 1:
			is = is[:len(is)-1]
		case 2:
			if isMiddleCard(is[len(is)-1].Goto) {
				is = is[:len(is)-1]
			} else {
				is = is[:len(is)-2]
			}
		case 3:
			if isMiddleCard(is[len(is)-1].Goto) {
				is = is[:len(is)-2]
			} else {
				is = is[:len(is)-3]
			}
		}
	}
	rl := len(is)
	if rl == 0 {
		is = _emptyItem
		return
	}
	if idx == 0 {
		idx = now.Unix()
	}
	for i, r := range is {
		if pull {
			r.Idx = idx + int64(rl-i)
		} else {
			r.Idx = idx - int64(i+1)
		}
		if ad, ok := adInfom[i]; ok {
			r.SrcID = ad.Source
			r.RequestID = ad.RequestID
			r.IsAdLoc = ad.IsAdLoc
			r.IsAd = ad.IsAd
			r.CmMark = ad.CmMark
			r.AdIndex = ad.Index
			r.ClientIP = ip
			r.CardIndex = i + 1
		} else if r.IsAd {
			r.CardIndex = i + 1
		}
		if i == 0 {
			r.AutoplayCard = autoPlay
		}
	}
	return
}

func (s *Service) adResource(plat int8, build int) (resource int64) {
	const (
		_androidBanAd = 500001
	)
	if plat == model.PlatIPhone || plat == model.PlatIPhoneB || (plat == model.PlatAndroid && build >= _androidBanAd) || plat == model.PlatIPad {
		resource = s.cmResourceMap[plat]
	}
	return
}

func swapTwoItem(rs []*feed.Item, i *feed.Item) (is []*feed.Item) {
	rs[len(rs)-1].Idx, i.Idx = i.Idx, rs[len(rs)-1].Idx
	is = append(rs, rs[len(rs)-1])
	is[len(is)-2] = i
	return
}

func swapThreeItem(rs []*feed.Item, i *feed.Item) (is []*feed.Item) {
	rs[len(rs)-1].Idx, i.Idx = i.Idx, rs[len(rs)-1].Idx
	rs[len(rs)-2].Idx, rs[len(is)-1].Idx = rs[len(rs)-1].Idx, rs[len(rs)-2].Idx
	is = append(rs, rs[len(rs)-1])
	is[len(is)-2] = i
	is[len(is)-3], is[len(is)-2] = is[len(is)-2], is[len(is)-3]
	return
}

func isMiddleCard(gt string) bool {
	return gt == model.GotoRank || gt == model.GotoAdWeb || gt == model.GotoPlayer ||
		gt == model.GotoPlayerLive || gt == model.GotoConverge || gt == model.GotoSpecial || gt == model.GotoAdLarge || gt == model.GotoLiveUpRcmd
}

func (s *Service) indexCount(plat int8) (count int) {
	if plat == model.PlatIPad {
		count = s.c.Feed.Index.IPadCount
	} else {
		count = s.c.Feed.Index.Count
	}
	return
}
