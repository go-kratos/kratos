package feed

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	cdm "go-common/app/interface/main/app-card/model"
	"go-common/app/interface/main/app-card/model/bplus"
	"go-common/app/interface/main/app-card/model/card"
	"go-common/app/interface/main/app-card/model/card/ai"
	"go-common/app/interface/main/app-card/model/card/audio"
	"go-common/app/interface/main/app-card/model/card/bangumi"
	"go-common/app/interface/main/app-card/model/card/banner"
	"go-common/app/interface/main/app-card/model/card/cm"
	"go-common/app/interface/main/app-card/model/card/live"
	"go-common/app/interface/main/app-card/model/card/operate"
	"go-common/app/interface/main/app-card/model/card/show"
	"go-common/app/interface/main/app-feed/model"
	"go-common/app/interface/main/app-feed/model/feed"
	tag "go-common/app/interface/main/tag/model"
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

const _qn480 = 32

var (
	_cardAdAvm = map[int]struct{}{
		1: struct{}{},
	}
	_cardAdWebm = map[int]struct{}{
		2:  struct{}{},
		7:  struct{}{},
		20: struct{}{},
	}
	_cardAdWebSm = map[int]struct{}{
		3:  struct{}{},
		26: struct{}{},
	}
	_followMode = &feed.FollowMode{
		Title: "当前为首页推荐 - 关注模式（内测版）",
		Option: []*feed.Option{
			{Title: "通用模式", Desc: "开启后，推荐你可能感兴趣的内容", Value: 0},
			{Title: "关注模式（内测版）", Desc: "开启后，仅显示关注UP主更新的视频", Value: 1},
		},
		ToastMessage: "关注UP主的内容已经看完啦，请稍后再试",
	}
)

func (s *Service) Index2(c context.Context, buvid string, mid int64, plat int8, param *feed.IndexParam, style int, now time.Time) (is []card.Handler, config *feed.Config, infoc *feed.Infoc, err error) {
	var (
		rs        []*ai.Item
		adm       map[int]*cm.AdInfo
		adAidm    map[int64]struct{}
		banners   []*banner.Banner
		version   string
		blackAidm map[int64]struct{}
		adInfom   map[int]*cm.AdInfo
		follow    *operate.Card
		info      *locmdl.Info
	)
	ip := metadata.String(c, metadata.RemoteIP)
	config = s.indexConfig(c, plat, buvid, mid, param)
	if config.FollowMode == nil {
		param.RecsysMode = 0
	}
	noCache := param.RecsysMode == 1
	followMode := config.FollowMode != nil
	infoc = &feed.Infoc{}
	infoc.AutoPlayInfoc = fmt.Sprintf("%d|%d", config.AutoplayCard, param.AutoPlayCard)
	if info, err = s.loc.Info(c, ip); err != nil {
		log.Warn("s.loc.Info(%v) error(%v)", ip, err)
		err = nil
	}
	group := s.group(mid, buvid)
	if !s.c.Feed.Index.Abnormal || followMode {
		g, ctx := errgroup.WithContext(c)
		g.Go(func() error {
			rs, infoc.UserFeature, infoc.IsRcmd, infoc.NewUser, infoc.Code = s.indexRcmd2(ctx, plat, buvid, mid, param, group, info, style, infoc.AutoPlayInfoc, noCache, now)
			return nil
		})
		g.Go(func() (err error) {
			if banners, version, err = s.indexBanner2(ctx, plat, buvid, mid, param); err != nil {
				log.Error("%+v", err)
				err = nil
			}
			return
		})
		if param.RecsysMode == 0 {
			g.Go(func() (err error) {
				if adm, adAidm, err = s.indexAd2(ctx, plat, buvid, mid, param, info, style, now); err != nil {
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
				if follow, err = s.SearchFollow2(ctx, param.Platform, param.MobiApp, param.Device, buvid, param.Build, mid); err != nil {
					log.Error("%+v", err)
					err = nil
				}
				return
			})
		}
		if err = g.Wait(); err != nil {
			log.Error("%+v", err)
			return
		}
		if param.RecsysMode == 1 {
			var tmp []*ai.Item
			for _, r := range rs {
				if r.Goto == model.GotoBanner {
					continue
				}
				tmp = append(tmp, r)
			}
			if len(tmp) == 0 {
				is = []card.Handler{}
				return
			}
		}
		rs, adInfom = s.mergeItem2(c, plat, mid, rs, adm, adAidm, banners, version, blackAidm, follow, followMode)
	} else {
		count := s.indexCount(plat)
		rs = s.recommendCache(count)
		log.Warn("feed index show disaster recovery data len(%d)", len(is))
	}
	if config.AutoplayCard == 1 && cdm.Columnm[param.Column] == cdm.ColumnSvrSingle {
		param.Qn = _qn480
	}
	is, infoc.IsRcmd = s.dealItem2(c, mid, buvid, plat, rs, param, infoc.IsRcmd, noCache, followMode, follow, now)
	s.dealAdLoc(is, param, adInfom, now)
	return
}

func (s *Service) indexConfig(c context.Context, plat int8, buvid string, mid int64, param *feed.IndexParam) (config *feed.Config) {
	config = &feed.Config{}
	config.Column = cdm.Columnm[param.Column]
	// if mid > 0 && mid%20 == 19 {
	// 	config.FeedCleanAbtest = 1
	// } else {
	// 	config.FeedCleanAbtest = 0
	// }
	config.FeedCleanAbtest = 0
	if !model.IsIPad(plat) {
		if ab, ok := s.abtestCache[_feedgroups]; ok {
			if ab.AbTestIn(buvid + _feedgroups) {
				switch param.AutoPlayCard {
				case 0, 1, 2, 3:
					config.AutoplayCard = 1
				default:
					config.AutoplayCard = 2
				}
			} else {
				config.AutoplayCard = 2
			}
		} else {
			switch param.AutoPlayCard {
			case 1, 3:
				config.AutoplayCard = 1
			default:
				config.AutoplayCard = 2
			}
		}
	} else {
		// ipad 不允许自动播放
		config.AutoplayCard = 2
	}
	if mid < 1 {
		return
	}
	if _, ok := s.autoplayMidsCache[mid]; ok && param.AutoPlayCard != 4 {
		config.AutoplayCard = 1
	}
	if _, ok := s.followModeList[mid]; ok {
		tmpConfig := &feed.FollowMode{}
		if s.c.Feed.Index.FollowMode == nil {
			*tmpConfig = *_followMode
		} else {
			*tmpConfig = *s.c.Feed.Index.FollowMode
		}
		if param.RecsysMode != 1 {
			tmpConfig.ToastMessage = ""
		}
		config.FollowMode = tmpConfig
	}
	return
}

func (s *Service) indexRcmd2(c context.Context, plat int8, buvid string, mid int64, param *feed.IndexParam, group int, zone *locmdl.Info, style int, autoPlay string, noCache bool, now time.Time) (is []*ai.Item, userFeature json.RawMessage, isRcmd, newUser bool, code int) {
	count := s.indexCount(plat)
	if buvid != "" || mid > 0 {
		var (
			err    error
			zoneID int64
		)
		if zone != nil {
			zoneID = zone.ZoneID
		}
		if is, userFeature, code, newUser, err = s.rcmd.Recommend(c, plat, buvid, mid, param.Build, param.LoginEvent, param.ParentMode, param.RecsysMode, zoneID, group, param.Interest, param.Network, style, param.Column, param.Flush, autoPlay, now); err != nil {
			log.Error("%+v", err)
		}
		if noCache {
			isRcmd = true
			return
		}
		if len(is) != 0 {
			isRcmd = true
		}
		var fromCache bool
		if len(is) == 0 && mid > 0 && !ecode.ServiceUnavailable.Equal(err) {
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

func (s *Service) indexAd2(c context.Context, plat int8, buvid string, mid int64, param *feed.IndexParam, zone *locmdl.Info, style int, now time.Time) (adm map[int]*cm.AdInfo, adAidm map[int64]struct{}, err error) {
	var advert *cm.Ad
	resource := s.adResource(plat, param.Build)
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
	if advert, err = s.ad.Ad(c, mid, param.Build, buvid, []int64{resource}, country, province, city, param.Network, param.MobiApp, param.Device, param.OpenEvent, param.AdExtra, style, now); err != nil {
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

func (s *Service) indexBanner2(c context.Context, plat int8, buvid string, mid int64, param *feed.IndexParam) (banners []*banner.Banner, version string, err error) {
	hash := param.BannerHash
	if param.LoginEvent != 0 {
		hash = ""
	}
	banners, version, err = s.banners(c, plat, param.Build, mid, buvid, param.Network, param.MobiApp, param.Device, param.OpenEvent, param.AdExtra, hash)
	return
}

func (s *Service) mergeItem2(c context.Context, plat int8, mid int64, rs []*ai.Item, adm map[int]*cm.AdInfo, adAidm map[int64]struct{}, banners []*banner.Banner, version string, blackAids map[int64]struct{}, follow *operate.Card, followMode bool) (is []*ai.Item, adInfom map[int]*cm.AdInfo) {
	if len(rs) == 0 {
		return
	}
	const (
		cardIndex     = 7
		cardIndexIPad = 17
		cardOffset    = 2
	)
	if len(banners) != 0 {
		rs = append([]*ai.Item{&ai.Item{Goto: model.GotoBanner, Banners: banners, Version: version}}, rs...)
		for index, ad := range adm {
			if _, ok := _cardAdWebm[ad.CardType]; ok && ((model.IsIPad(plat) && index <= cardIndexIPad) || index <= cardIndex) {
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
	var existsAdWeb bool
	for _, r := range rs {
		for {
			if ad, ok := adm[len(is)]; ok {
				if ad.CreativeID != 0 {
					var item *ai.Item
					if _, ok := _cardAdAvm[ad.CardType]; ok {
						item = &ai.Item{ID: ad.CreativeContent.VideoID, Goto: model.GotoAdAv, Ad: ad}
					} else if _, ok := _cardAdWebm[ad.CardType]; ok {
						item = &ai.Item{Goto: model.GotoAdWeb, Ad: ad}
						existsAdWeb = true
					} else if _, ok := _cardAdWebSm[ad.CardType]; ok {
						item = &ai.Item{Goto: model.GotoAdWebS, Ad: ad}
					} else {
						b, _ := json.Marshal(ad)
						log.Error("ad---%s", b)
						break
					}
					is = append(is, item)
					continue
				} else {
					adInfom[ad.CardIndex-1] = ad
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
		} else if r.Goto == model.GotoBanner && len(is) != 0 {
			// banner 必须在第一位
			continue
		} else if r.Goto == model.GotoRank && existsAdWeb {
			continue
		} else if r.Goto == model.GotoLogin && mid > 0 {
			continue
		} else if r.Goto == model.GotoFollowMode && !followMode {
			continue
		}
		is = append(is, r)
	}
	return
}

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
			h.Get().AdInfo.CardIndex = i + 1
		}
	}
}

func (s *Service) dealItem2(c context.Context, mid int64, buvid string, plat int8, rs []*ai.Item, param *feed.IndexParam, isRcmd, noCache, followMode bool, follow *operate.Card, now time.Time) (is []card.Handler, isAI bool) {
	if len(rs) == 0 {
		is = []card.Handler{}
		return
	}
	var (
		aids, tids, roomIDs, sids, metaIDs, shopIDs, audioIDs, picIDs []int64
		seasonIDs                                                     []int32
		upIDs, avUpIDs, rmUpIDs, mtUpIDs                              []int64
		am                                                            map[int64]*archive.ArchiveWithPlayer
		tagm                                                          map[int64]*tag.Tag
		rm                                                            map[int64]*live.Room
		sm                                                            map[int64]*bangumi.Season
		hasUpdate, getBanner                                          bool
		update                                                        *bangumi.Update
		metam                                                         map[int64]*article.Meta
		shopm                                                         map[int64]*show.Shopping
		audiom                                                        map[int64]*audio.Audio
		cardm                                                         map[int64]*account.Card
		statm                                                         map[int64]*relation.Stat
		moe                                                           *bangumi.Moe
		isAtten                                                       map[int64]int8
		arcOK                                                         bool
		rank                                                          *operate.Card
		seasonm                                                       map[int32]*episodegrpc.EpisodeCardsProto
		banners                                                       []*banner.Banner
		version                                                       string
		picm                                                          map[int64]*bplus.Picture
	)
	convergem := map[int64]*operate.Card{}
	followm := map[int64]*operate.Card{}
	downloadm := map[int64]*operate.Card{}
	specialm := map[int64]*operate.Card{}
	liveUpm := map[int64][]*live.Card{}
	isAI = isRcmd
	for _, r := range rs {
		if r == nil {
			continue
		}
		switch r.Goto {
		case model.GotoBanner:
			if len(r.Banners) != 0 {
				banners = r.Banners
				version = r.Version
			} else {
				getBanner = true
			}
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
			os, aid := s.RankCard(plat)
			rank = &operate.Card{}
			rank.FromRank(os)
			aids = append(aids, aid...)
		case model.GotoBangumiRcmd:
			hasUpdate = true
		case model.GotoConverge:
			cardm, aid, roomID, metaID := s.convergeCard(c, 3, r.ID)
			for id, card := range cardm {
				convergem[id] = card
			}
			aids = append(aids, aid...)
			roomIDs = append(roomIDs, roomID...)
			metaIDs = append(metaIDs, metaID...)
		case model.GotoGameDownloadS:
			cardm := s.downloadCard(c, r.ID)
			for id, card := range cardm {
				downloadm[id] = card
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
			cardm, upID := s.liveUpRcmdCard(c, r.ID)
			for id, card := range cardm {
				liveUpm[id] = card
			}
			upIDs = append(upIDs, upID...)
		case model.GotoSubscribe:
			cardm, upID, tid := s.subscribeCard(c, r.ID)
			for id, card := range cardm {
				followm[id] = card
			}
			upIDs = append(upIDs, upID...)
			tids = append(tids, tid...)
		case model.GotoSearchSubscribe:
			if follow != nil {
				followm[follow.ID] = follow
				for _, item := range follow.Items {
					upIDs = append(upIDs, item.ID)
				}
			}
		case model.GotoChannelRcmd:
			cardm, aid, tid := s.channelRcmdCard(c, r.ID)
			for id, card := range cardm {
				followm[id] = card
			}
			aids = append(aids, aid...)
			tids = append(tids, tid...)
		case model.GotoSpecial, model.GotoSpecialS:
			cardm := s.specialCard(c, r.ID)
			for id, card := range cardm {
				specialm[id] = card
			}
		case model.GotoPicture:
			if r.ID != 0 {
				picIDs = append(picIDs, r.ID)
			}
			if r.RcmdReason != nil && r.RcmdReason.Style == 4 {
				upIDs = append(upIDs, r.RcmdReason.FollowedMid)
			}
		}
	}
	g, ctx := errgroup.WithContext(c)
	if getBanner {
		g.Go(func() (err error) {
			if banners, version, err = s.banners(ctx, plat, param.Build, mid, buvid, param.Network, param.MobiApp, param.Device, param.OpenEvent, param.AdExtra, ""); err != nil {
				log.Error("%+v", err)
				err = nil
			}
			return
		})
	}
	if len(aids) != 0 {
		g.Go(func() (err error) {
			if am, err = s.ArchivesWithPlayer(ctx, aids, param.Qn, param.MobiApp, param.Fnver, param.Fnval, param.ForceHost, param.Build); err != nil {
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
			if seasonm, err = s.bgm.CardsInfoReply(ctx, seasonIDs); err != nil {
				log.Error("%+v", err)
				err = nil
			}
			return
		})
	}
	if hasUpdate && mid > 0 {
		g.Go(func() (err error) {
			/*
				{
				    "code": 0,
				    "message": "success",
				    "result": {
				        "title": "小埋。。。",
				        "square_cover": "http://i0.hdslb.com/bfs/bangumi/dd2281c9f1c44e07c835e488ce1e1bae36f533e3.jpg",
				        "updates": 67
				    }
				}
			*/
			if update, err = s.bgm.Updates(ctx, mid, now); err != nil {
				log.Error("%+v", err)
				err = nil
			}
			return
		})
	}
	if len(metaIDs) != 0 {
		g.Go(func() (err error) {
			if metam, err = s.art.Articles(ctx, metaIDs); err != nil {
				log.Error("%+v", err)
				err = nil
			}
			for _, meta := range metam {
				if meta.Author != nil {
					mtUpIDs = append(mtUpIDs, meta.Author.Mid)
				}
			}
			return
		})
	}
	if len(shopIDs) != 0 {
		g.Go(func() (err error) {
			if shopm, err = s.show.Card(ctx, shopIDs); err != nil {
				log.Error("%+v", err)
				err = nil
			}
			return
		})
	}
	if len(audioIDs) != 0 {
		g.Go(func() (err error) {
			if audiom, err = s.audio.Audios(ctx, audioIDs); err != nil {
				log.Error("%+v", err)
				err = nil
			}
			return
		})
	}
	if len(picIDs) != 0 {
		g.Go(func() (err error) {
			if picm, err = s.bplus.DynamicDetail(ctx, picIDs...); err != nil {
				log.Error("%+v", err)
				err = nil
			}
			return
		})
	}
	// 萌战下线
	// if mid > 0 {
	// 	g.Go(func() (err error) {
	// 		if moe, err = s.bgm.FollowPull(ctx, mid, mobiApp, device, now); err != nil {
	// 			log.Error("%+v", err)
	// 			err = nil
	// 		}
	// 		return
	// 	})
	// }
	if err := g.Wait(); err != nil {
		log.Error("%+v", err)
		if noCache {
			is = []card.Handler{}
			return
		}
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
			if mid > 0 && param.RecsysMode == 0 {
				g.Go(func() error {
					isAtten = s.acc.IsAttention(ctx, upIDs, mid)
					return nil
				})
			}
		}
		g.Wait()
	}
	isAI = isAI && arcOK
	if moe != nil {
		moePos := s.c.Feed.Index.MoePosition
		if moePos-1 >= 0 && moePos-1 <= len(rs) {
			rs = append(rs[:moePos-1], append([]*ai.Item{&ai.Item{ID: moe.ID, Goto: model.GotoMoe}}, rs[moePos-1:]...)...)
		}
	}
	var cardTotal int
	is = make([]card.Handler, 0, len(rs))
	insert := map[int]card.Handler{}
	for _, r := range rs {
		if r == nil {
			continue
		}
		var (
			main     interface{}
			cardType cdm.CardType
		)
		op := &operate.Card{}
		op.From(cdm.CardGt(r.Goto), r.ID, r.Tid, plat, param.Build)
		// 卡片展示点赞数实验
		// if mid%20 == 11 && ((plat == model.PlatIPhone && param.Build >= 8290) || (plat == model.PlatAndroid && param.Build >= 5360000)) {
		// 	op.FromSwitch(cdm.SwitchFeedIndexLike)
		// }
		// 变化卡片类型
		switch r.Goto {
		case model.GotoSpecialS, model.GotoGameDownloadS, model.GotoShoppingS:
			if r.Style == 2 {
				cardType = cdm.LargeCoverV1
			}
		case model.GotoPicture:
			if p, ok := picm[r.ID]; ok {
				switch cdm.Columnm[param.Column] {
				case cdm.ColumnSvrSingle:
					if len(p.Imgs) < 3 {
						cardType = cdm.OnePicV1
					} else {
						cardType = cdm.ThreePicV1
					}
				case cdm.ColumnSvrDouble:
					if len(p.Imgs) < 3 {
						// 版本过滤5.37为新卡片
						if (plat == model.PlatIPhone && param.Build > 8300) || (plat == model.PlatAndroid && param.Build > 5365000) {
							cardType = cdm.OnePicV2
						} else {
							cardType = cdm.SmallCoverV2
						}
					} else {
						cardType = cdm.ThreePicV2
					}
				default:
					continue
				}
			} else {
				continue
			}
		case model.GotoInterest:
			switch cdm.Columnm[param.Column] {
			case cdm.ColumnSvrSingle:
				cardType = cdm.OptionsV1
			case cdm.ColumnSvrDouble:
				cardType = cdm.OptionsV2
			default:
				continue
			}
		case model.GotoFollowMode:
			cardType = cdm.Select
		default:
		}
		h := card.Handle(plat, cdm.CardGt(r.Goto), cardType, param.Column, r, tagm, isAtten, statm, cardm)
		if h == nil {
			continue
		}
		switch r.Goto {
		case model.GotoAv, model.GotoUpRcmdAv, model.GotoPlayer:
			if !arcOK {
				if r.Archive != nil {
					am = map[int64]*archive.ArchiveWithPlayer{r.Archive.Aid: &archive.ArchiveWithPlayer{Archive3: r.Archive}}
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
			if plat == model.PlatIPhone && param.Build > 8290 || plat == model.PlatAndroid && param.Build > 5365000 {
				op.Switch = cdm.SwitchCooperationShow
			} else {
				op.Switch = cdm.SwitchCooperationHide
			}
		case model.GotoLive, model.GotoPlayerLive:
			main = rm
		case model.GotoBangumi:
			main = sm
		case model.GotoPGC:
			main = seasonm
		case model.GotoLogin:
			op.FromLogin(r.ID)
		case model.GotoSpecial, model.GotoSpecialS:
			op = specialm[r.ID]
		case model.GotoRank:
			main = map[cdm.Gt]interface{}{cdm.GotoAv: am}
			op = rank
		case model.GotoBangumiRcmd:
			main = update
		case model.GotoBanner:
			op.FromBanner(banners, version)
		case model.GotoConverge:
			main = map[cdm.Gt]interface{}{cdm.GotoAv: am, cdm.GotoLive: rm, cdm.GotoArticle: metam}
			op = convergem[r.ID]
		case model.GotoGameDownloadS:
			op = downloadm[r.ID]
		case model.GotoArticleS:
			main = metam
		case model.GotoShoppingS:
			main = shopm
		case model.GotoAudio:
			main = audiom
		case model.GotoChannelRcmd:
			main = am
			op = followm[r.ID]
		case model.GotoSubscribe, model.GotoSearchSubscribe:
			op = followm[r.ID]
		case model.GotoLiveUpRcmd:
			main = liveUpm
		case model.GotoMoe:
			main = moe
		case model.GotoPicture:
			main = picm
		case model.GotoAdAv:
			main = am
			op.FromAdAv(r.Ad)
		case model.GotoAdWebS, model.GotoAdWeb:
			main = r.Ad
		case model.GotoInterest:
			main = s.c.Feed.Index.Interest
		case model.GotoFollowMode:
			var (
				title  string
				desc   string
				button []string
			)
			if s.c.Feed.Index.FollowMode != nil && s.c.Feed.Index.FollowMode.Card != nil {
				title = s.c.Feed.Index.FollowMode.Card.Title
				desc = s.c.Feed.Index.FollowMode.Card.Desc
				button = s.c.Feed.Index.FollowMode.Card.Button
			}
			op.FromFollowMode(title, desc, button)
		default:
			log.Warn("unexpected goto(%s) %+v", r.Goto, r)
			continue
		}
		if op != nil {
			op.Plat = plat
			op.Build = param.Build
		}
		h.From(main, op)
		// 卡片不正常要continue
		if !h.Get().Right {
			continue
		}
		switch r.Goto {
		case model.GotoAdAv, model.GotoAdWebS, model.GotoAdWeb:
			// 判断结果列表长度，如果列表的末尾不是广告位，则放到插入队列里
			if len(is) != r.Ad.CardIndex-1 {
				insert[r.Ad.CardIndex-1] = h
				// 插入队列后一定要continue，否则就直接加到队列末尾了
				continue
			}
		}
		is, cardTotal = s.appendItem(plat, is, h, param.Column, cardTotal)
		// 从插入队列里获取广告
		if h, ok := insert[len(is)]; ok {
			is, cardTotal = s.appendItem(plat, is, h, param.Column, cardTotal)
		}
	}
	// 双列末尾卡片去空窗
	if !model.IsIPad(plat) {
		if cdm.Columnm[param.Column] == cdm.ColumnSvrDouble {
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

func (s *Service) appendItem(plat int8, rs []card.Handler, h card.Handler, column cdm.ColumnStatus, cardTotal int) (is []card.Handler, total int) {
	h.Get().ThreePointFrom()
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

func (s *Service) Converge(c context.Context, mid int64, plat int8, param *feed.ConvergeParam, now time.Time) (is []card.Handler, converge *operate.Card, err error) {
	cardm, _, _, _ := s.convergeCard(c, 0, param.ID)
	converge, ok := cardm[param.ID]
	if !ok {
		is = []card.Handler{}
		return
	}
	rs := make([]*ai.Item, 0, len(converge.Items))
	for _, item := range converge.Items {
		rs = append(rs, &ai.Item{ID: item.ID, Goto: string(item.CardGoto)})
	}
	indexParam := &feed.IndexParam{
		MobiApp:   param.MobiApp,
		Device:    param.Device,
		Build:     param.Build,
		Qn:        param.Qn,
		Fnver:     param.Fnver,
		Fnval:     param.Fnval,
		ForceHost: param.ForceHost,
	}
	is, _ = s.dealItem2(c, mid, "", plat, rs, indexParam, false, true, false, nil, now)
	for _, item := range is {
		// 运营tab页没有不感兴趣
		item.Get().ThreePointWatchLater()
	}
	return
}
