package search

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go-common/app/interface/main/app-interface/conf"
	accdao "go-common/app/interface/main/app-interface/dao/account"
	arcdao "go-common/app/interface/main/app-interface/dao/archive"
	artdao "go-common/app/interface/main/app-interface/dao/article"
	bangumidao "go-common/app/interface/main/app-interface/dao/bangumi"
	bplusdao "go-common/app/interface/main/app-interface/dao/bplus"
	livedao "go-common/app/interface/main/app-interface/dao/live"
	resdao "go-common/app/interface/main/app-interface/dao/resource"
	srchdao "go-common/app/interface/main/app-interface/dao/search"
	tagdao "go-common/app/interface/main/app-interface/dao/tag"
	"go-common/app/interface/main/app-interface/model"
	"go-common/app/interface/main/app-interface/model/bangumi"
	"go-common/app/interface/main/app-interface/model/banner"
	"go-common/app/interface/main/app-interface/model/bplus"
	"go-common/app/interface/main/app-interface/model/live"
	"go-common/app/interface/main/app-interface/model/search"
	tagmdl "go-common/app/interface/main/app-interface/model/tag"
	article "go-common/app/interface/openplatform/article/model"
	account "go-common/app/service/main/account/model"
	"go-common/app/service/main/archive/api"
	resmdl "go-common/app/service/main/resource/model"
	seasongrpc "go-common/app/service/openplatform/pgc-season/api/grpc/season/v1"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

// const for search
const (
	_oldAndroid = 514000
	_oldIOS     = 6090

	IPhoneSearchResourceID  = 2447
	AndroidSearchResourceID = 2450
	IPadSearchResourceID    = 2811
)

var (
	_emptyItem   = []*search.Item{}
	_emptyResult = &search.Result{
		NavInfo: []*search.NavInfo{},
		Page:    0,
		Items: search.ResultItems{
			Season:   _emptyItem,
			Upper:    _emptyItem,
			Movie:    _emptyItem,
			Archive:  _emptyItem,
			LiveRoom: _emptyItem,
			LiveUser: _emptyItem,
		},
	}
)

// Service is search service
type Service struct {
	c          *conf.Config
	srchDao    *srchdao.Dao
	accDao     *accdao.Dao
	arcDao     *arcdao.Dao
	liveDao    *livedao.Dao
	artDao     *artdao.Dao
	resDao     *resdao.Dao
	tagDao     *tagdao.Dao
	bangumiDao *bangumidao.Dao
	bplusDao   *bplusdao.Dao
	// config
	seasonNum          int
	movieNum           int
	seasonShowMore     int
	movieShowMore      int
	upUserNum          int
	uvLimit            int
	userNum            int
	userVideoLimit     int
	biliUserNum        int
	biliUserVideoLimit int
	iPadSearchBangumi  int
	iPadSearchFt       int
}

// New is search service initial func
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:                  c,
		srchDao:            srchdao.New(c),
		accDao:             accdao.New(c),
		liveDao:            livedao.New(c),
		arcDao:             arcdao.New(c),
		artDao:             artdao.New(c),
		resDao:             resdao.New(c),
		tagDao:             tagdao.New(c),
		bangumiDao:         bangumidao.New(c),
		bplusDao:           bplusdao.New(c),
		seasonNum:          c.Search.SeasonNum,
		movieNum:           c.Search.MovieNum,
		seasonShowMore:     c.Search.SeasonMore,
		movieShowMore:      c.Search.MovieMore,
		upUserNum:          c.Search.UpUserNum,
		uvLimit:            c.Search.UVLimit,
		userNum:            c.Search.UpUserNum,
		userVideoLimit:     c.Search.UVLimit,
		biliUserNum:        c.Search.BiliUserNum,
		biliUserVideoLimit: c.Search.BiliUserVideoLimit,
		iPadSearchBangumi:  c.Search.IPadSearchBangumi,
		iPadSearchFt:       c.Search.IPadSearchFt,
	}
	return
}

// Search get all type search data.
func (s *Service) Search(c context.Context, mid, zoneid int64, mobiApp, device, platform, buvid, keyword, duration, order, filtered, lang, fromSource, recommend, parent string, plat int8, rid, highlight, build, pn, ps, isQuery int, old bool, now time.Time) (res *search.Result, err error) {
	const (
		_newIPhonePGC      = 6500
		_newAndroidPGC     = 519010
		_newIPhoneSearch   = 6500
		_newAndroidSearch  = 5215000
		_newAndroidBSearch = 591200
	)
	var (
		newPGC, flow, isNewTwitter bool
		avids                      []int64
		avm                        map[int64]*api.Arc
		owners                     []int64
		follows                    map[int64]bool
		roomIDs                    []int64
		lm                         map[int64]*live.RoomInfo
		seasonIDs                  []int64
		bangumis                   map[string]*bangumi.Card
		//tagSeasonIDs []int32
		tagBangumis    map[int32]*seasongrpc.CardInfoProto
		tags           []int64
		tagMyInfos     []*tagmdl.Tag
		dynamicIDs     []int64
		dynamicDetails map[int64]*bplus.Detail
		accInfos       map[int64]*account.Info
		cooperation    bool
	)
	// android 概念版 591205
	if (plat == model.PlatAndroid && build >= _newAndroidPGC && build != 591205) || (plat == model.PlatIPhone && build >= _newIPhonePGC && build != 7140) || (plat == model.PlatAndroidB && build >= _newAndroidBSearch) || (plat == model.PlatIPad && build >= search.SearchNewIPad) || (plat == model.PlatIpadHD && build >= search.SearchNewIPadHD) || model.IsIPhoneB(plat) {
		newPGC = true
	}
	// 处理一个ios概念版是 7140，是否需要过滤
	if (plat == model.PlatAndroid && build >= _newAndroidSearch) || (plat == model.PlatIPhone && build >= _newIPhoneSearch && build != 7140) || (plat == model.PlatAndroidB && build >= _newAndroidBSearch) || model.IsIPhoneB(plat) {
		flow = true
	}
	var (
		seasonNum int
		movieNum  int
	)
	if (plat == model.PlatIPad && build >= search.SearchNewIPad) || (plat == model.PlatIpadHD && build >= search.SearchNewIPadHD) {
		seasonNum = s.iPadSearchBangumi
		movieNum = s.iPadSearchFt
	} else {
		seasonNum = s.seasonNum
		movieNum = s.movieNum
	}
	all, code, err := s.srchDao.Search(c, mid, zoneid, mobiApp, device, platform, buvid, keyword, duration, order, filtered, fromSource, recommend, parent, plat, seasonNum, movieNum, s.upUserNum, s.uvLimit, s.userNum, s.userVideoLimit, s.biliUserNum, s.biliUserVideoLimit, rid, highlight, build, pn, ps, isQuery, old, now, newPGC, flow)
	if err != nil {
		log.Error("%+v", err)
		return
	}
	if (model.IsAndroid(plat) && build > s.c.SearchBuildLimit.NewTwitterAndroid) || (model.IsIPhone(plat) && build > s.c.SearchBuildLimit.NewTwitterIOS) {
		isNewTwitter = true
	}
	if code == model.ForbidCode || code == model.NoResultCode {
		res = _emptyResult
		err = nil
		return
	}
	res = &search.Result{}
	res.Trackid = all.Trackid
	res.Page = all.Page
	res.Array = all.FlowPlaceholder
	res.Attribute = all.Attribute
	res.NavInfo = s.convertNav(all, plat, build, lang, old, newPGC)
	if len(all.FlowResult) != 0 {
		var item []*search.Item
		for _, v := range all.FlowResult {
			switch v.Type {
			case search.TypeUser, search.TypeBiliUser:
				owners = append(owners, v.User.Mid)
				for _, vr := range v.User.Res {
					avids = append(avids, vr.Aid)
				}
				roomIDs = append(roomIDs, v.User.RoomID)
			case search.TypeVideo:
				avids = append(avids, v.Video.ID)
			case search.TypeLive:
				roomIDs = append(roomIDs, v.Live.RoomID)
			case search.TypeMediaBangumi, search.TypeMediaFt:
				seasonIDs = append(seasonIDs, v.Media.SeasonID)
			case search.TypeStar:
				if v.Star.MID != 0 {
					owners = append(owners, v.Star.MID)
				}
				if v.Star.TagID != 0 {
					tags = append(tags, v.Star.TagID)
				}
			case search.TypeArticle:
				owners = append(owners, v.Article.Mid)
			case search.TypeChannel:
				tags = append(tags, v.Channel.TagID)
				if len(v.Channel.Values) > 0 {
					for _, vc := range v.Channel.Values {
						switch vc.Type {
						case search.TypeVideo:
							if vc.Video != nil {
								avids = append(avids, vc.Video.ID)
							}
							//case search.TypeLive:
							//	if vc.Live != nil {
							//		roomIDs = append(roomIDs, vc.Live.RoomID)
							//	}
							//case search.TypeMediaBangumi, search.TypeMediaFt:
							//	if vc.Media != nil {
							//		tagSeasonIDs = append(tagSeasonIDs, int32(vc.Media.SeasonID))
							//	}
						}
					}
				}
			case search.TypeTwitter:
				dynamicIDs = append(dynamicIDs, v.Twitter.ID)
			}
		}
		g, ctx := errgroup.WithContext(c)
		if len(owners) != 0 {
			if mid > 0 {
				g.Go(func() error {
					follows = s.accDao.Relations3(ctx, owners, mid)
					return nil
				})
			}
			g.Go(func() (err error) {
				if accInfos, err = s.accDao.Infos3(ctx, owners); err != nil {
					log.Error("%v", err)
					err = nil
				}
				return
			})
		}
		if len(avids) != 0 {
			g.Go(func() (err error) {
				if avm, err = s.arcDao.Archives2(ctx, avids); err != nil {
					log.Error("%+v", err)
					err = nil
				}
				return
			})
		}
		if len(roomIDs) != 0 {
			g.Go(func() (err error) {
				if lm, err = s.liveDao.LiveByRIDs(ctx, roomIDs); err != nil {
					log.Error("%+v", err)
					err = nil
				}
				return
			})
		}
		if len(seasonIDs) != 0 {
			g.Go(func() (err error) {
				if bangumis, err = s.bangumiDao.Card(ctx, mid, seasonIDs); err != nil {
					log.Error("%+v", err)
					err = nil
				}
				return
			})
		}
		//if len(tagSeasonIDs) != 0 {
		//	g.Go(func() (err error) {
		//		if tagBangumis, err = s.bangumiDao.Cards(ctx, tagSeasonIDs); err != nil {
		//			log.Error("%+v", err)
		//			err = nil
		//		}
		//		return
		//	})
		//}
		if len(tags) != 0 {
			g.Go(func() (err error) {
				if tagMyInfos, err = s.tagDao.TagInfos(ctx, tags, mid); err != nil {
					log.Error("%v \n", err)
					err = nil
				}
				return
			})
		}
		if len(dynamicIDs) != 0 {
			g.Go(func() (err error) {
				if dynamicDetails, err = s.bplusDao.DynamicDetails(ctx, dynamicIDs, "search"); err != nil {
					log.Error("%v \n", err)
					err = nil
				}
				return
			})
		}
		if err = g.Wait(); err != nil {
			log.Error("%+v", err)
			return
		}
		if all.SuggestKeyword != "" && pn == 1 {
			i := &search.Item{Title: all.SuggestKeyword, Goto: model.GotoSuggestKeyWord, SugKeyWordType: 1}
			item = append(item, i)
		} else if all.CrrQuery != "" && pn == 1 {
			if (model.IsAndroid(plat) && build > s.c.SearchBuildLimit.QueryCorAndroid) || (model.IsIPhone(plat) && build > s.c.SearchBuildLimit.QueryCorIOS) {
				i := &search.Item{Title: fmt.Sprintf("已匹配%q的搜索结果", all.CrrQuery), Goto: model.GotoSuggestKeyWord, SugKeyWordType: 2}
				item = append(item, i)
			}
		}
		for _, v := range all.FlowResult {
			i := &search.Item{TrackID: v.TrackID, LinkType: v.LinkType, Position: v.Position}
			switch v.Type {
			case search.TypeVideo:
				if (model.IsAndroid(plat) && build > s.c.SearchBuildLimit.CooperationAndroid) || (model.IsIPhone(plat) && build > s.c.SearchBuildLimit.CooperationIOS) {
					cooperation = true
				}
				i.FromVideo(v.Video, avm[v.Video.ID], cooperation)
			case search.TypeLive:
				i.FromLive(v.Live, lm[v.Live.RoomID])
			case search.TypeMediaBangumi:
				i.FromMedia(v.Media, "", model.GotoBangumi, bangumis)
			case search.TypeMediaFt:
				i.FromMedia(v.Media, "", model.GotoMovie, bangumis)
			case search.TypeArticle:
				i.FromArticle(v.Article, accInfos[v.Article.Mid])
			case search.TypeSpecial:
				i.FromOperate(v.Operate, model.GotoSpecial)
			case search.TypeBanner:
				i.FromOperate(v.Operate, model.GotoBanner)
			case search.TypeUser:
				if follows[v.User.Mid] {
					i.Attentions = 1
				}
				i.FromUser(v.User, avm, lm[v.User.RoomID])
			case search.TypeBiliUser:
				if follows[v.User.Mid] {
					i.Attentions = 1
				}
				i.FromUpUser(v.User, avm, lm[v.User.RoomID])
			case search.TypeSpecialS:
				i.FromOperate(v.Operate, model.GotoSpecialS)
			case search.TypeGame:
				i.FromGame(v.Game)
			case search.TypeQuery:
				i.Title = v.TypeName
				i.FromQuery(v.Query)
			case search.TypeComic:
				i.FromComic(v.Comic)
			case search.TypeConverge:
				var (
					aids, rids, artids []int64
					am                 map[int64]*api.Arc
					rm                 map[int64]*live.Room
					artm               map[int64]*article.Meta
				)
				for _, c := range v.Operate.ContentList {
					switch c.Type {
					case 0:
						aids = append(aids, c.ID)
					case 1:
						rids = append(rids, c.ID)
					case 2:
						artids = append(artids, c.ID)
					}
				}
				g, ctx := errgroup.WithContext(c)
				if len(aids) != 0 {
					g.Go(func() (err error) {
						if am, err = s.arcDao.Archives2(ctx, aids); err != nil {
							log.Error("%+v", err)
							err = nil
						}
						return
					})
				}
				if len(rids) != 0 {
					g.Go(func() (err error) {
						if rm, err = s.liveDao.AppMRoom(ctx, rids); err != nil {
							log.Error("%+v", err)
							err = nil
						}
						return
					})
				}
				if len(artids) != 0 {
					g.Go(func() (err error) {
						if artm, err = s.artDao.Articles(ctx, artids); err != nil {
							log.Error("%+v", err)
							err = nil
						}
						return
					})
				}
				if err = g.Wait(); err != nil {
					log.Error("%+v", err)
					continue
				}
				i.FromConverge(v.Operate, am, rm, artm)
			case search.TypeTwitter:
				i.FromTwitter(v.Twitter, dynamicDetails, s.c.SearchDynamicSwitch.IsUP, s.c.SearchDynamicSwitch.IsCount, isNewTwitter)
			case search.TypeStar:
				if v.Star.TagID != 0 {
					i.URIType = search.StarChannel
					for _, myInfo := range tagMyInfos {
						if myInfo != nil && myInfo.TagID == v.Star.TagID {
							i.IsAttention = myInfo.IsAtten
							break
						}
					}
				} else if v.Star.MID != 0 {
					i.URIType = search.StarSpace
					if follows[v.Star.MID] {
						i.IsAttention = 1
					}
				}
				i.FromStar(v.Star)
			case search.TypeTicket:
				i.FromTicket(v.Ticket)
			case search.TypeProduct:
				i.FromProduct(v.Product)
			case search.TypeSpecialerGuide:
				i.FromSpecialerGuide(v.SpecialerGuide)
			case search.TypeChannel:
				i.FromChannel(v.Channel, avm, tagBangumis, lm, tagMyInfos)
			}
			if i.Goto != "" {
				item = append(item, i)
			}
		}
		res.Item = item
		if plat == model.PlatAndroid && build < search.SearchEggInfoAndroid {
			return
		}
		if all.EggInfo != nil {
			res.EasterEgg = &search.EasterEgg{ID: all.EggInfo.Source, ShowCount: all.EggInfo.ShowCount}
		}
		return
	}
	var items []*search.Item
	if all.SuggestKeyword != "" && pn == 1 {
		res.Items.SuggestKeyWord = &search.Item{Title: all.SuggestKeyword, Goto: model.GotoSuggestKeyWord}
	}
	// archive
	for _, v := range all.Result.Video {
		avids = append(avids, v.ID)
	}
	if duration == "0" && order == "totalrank" && rid == 0 {
		for _, v := range all.Result.Movie {
			if v.Type == "movie" {
				avids = append(avids, v.Aid)
			}
		}
	}
	if pn == 1 {
		for _, v := range all.Result.User {
			for _, vr := range v.Res {
				avids = append(avids, vr.Aid)
			}
		}
		if old {
			for _, v := range all.Result.UpUser {
				for _, vr := range v.Res {
					avids = append(avids, vr.Aid)
				}
				owners = append(owners, v.Mid)
				roomIDs = append(roomIDs, v.RoomID)
			}
		} else {
			for _, v := range all.Result.BiliUser {
				for _, vr := range v.Res {
					avids = append(avids, vr.Aid)
				}
				owners = append(owners, v.Mid)
				roomIDs = append(roomIDs, v.RoomID)
			}
		}
	}
	if model.IsOverseas(plat) {
		for _, v := range all.Result.LiveRoom {
			roomIDs = append(roomIDs, v.RoomID)
		}
		for _, v := range all.Result.LiveUser {
			roomIDs = append(roomIDs, v.RoomID)
		}
	}
	g, ctx := errgroup.WithContext(c)
	if len(owners) != 0 {
		if mid > 0 {
			g.Go(func() error {
				follows = s.accDao.Relations3(ctx, owners, mid)
				return nil
			})
		}
	}
	if len(avids) != 0 {
		g.Go(func() (err error) {
			if avm, err = s.arcDao.Archives2(ctx, avids); err != nil {
				log.Error("%+v", err)
				err = nil
			}
			return
		})
	}
	if len(roomIDs) != 0 {
		g.Go(func() (err error) {
			if lm, err = s.liveDao.LiveByRIDs(ctx, roomIDs); err != nil {
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
	if duration == "0" && order == "totalrank" && rid == 0 {
		var promptBangumi, promptFt string
		// season
		bangumi := all.Result.Bangumi
		items = make([]*search.Item, 0, len(bangumi))
		for _, v := range bangumi {
			si := &search.Item{}
			if (model.IsAndroid(plat) && build <= _oldAndroid) || (model.IsIPhone(plat) && build <= _oldIOS) {
				si.FromSeason(v, model.GotoBangumi)
			} else {
				si.FromSeason(v, model.GotoBangumiWeb)
			}
			items = append(items, si)
		}
		res.Items.Season = items
		// movie
		movie := all.Result.Movie
		items = make([]*search.Item, 0, len(movie))
		for _, v := range movie {
			si := &search.Item{}
			si.FromMovie(v, avm)
			items = append(items, si)
		}
		res.Items.Movie = items
		// season2
		mb := all.Result.MediaBangumi
		items = make([]*search.Item, 0, len(mb))
		for k, v := range mb {
			si := &search.Item{}
			if ((plat == model.PlatIPad && build >= search.SearchNewIPad) || (plat == model.PlatIpadHD && build >= search.SearchNewIPadHD)) && (k == len(mb)-1) && all.PageInfo.MediaBangumi.NumResults > s.iPadSearchBangumi {
				promptBangumi = fmt.Sprintf("查看全部番剧 ( %d ) >", all.PageInfo.MediaBangumi.NumResults)
			}
			si.FromMedia(v, promptBangumi, model.GotoBangumi, nil)
			items = append(items, si)
		}
		res.Items.Season2 = items
		// movie2
		mf := all.Result.MediaFt
		items = make([]*search.Item, 0, len(mf))
		for k, v := range mf {
			si := &search.Item{}
			if ((plat == model.PlatIPad && build >= search.SearchNewIPad) || (plat == model.PlatIpadHD && build >= search.SearchNewIPadHD)) && (k == len(mf)-1) && all.PageInfo.MediaFt.NumResults > s.iPadSearchFt {
				promptFt = fmt.Sprintf("查看全部影视 ( %d ) >", all.PageInfo.MediaFt.NumResults)
			}
			si.FromMedia(v, promptFt, model.GotoMovie, nil)
			si.Goto = model.GotoAv
			items = append(items, si)
		}
		res.Items.Movie2 = items
	}
	if pn == 1 {
		// upper + user
		var tmp []*search.User
		if old {
			tmp = all.Result.UpUser
		} else {
			tmp = all.Result.BiliUser
		}
		items = make([]*search.Item, 0, len(tmp)+len(all.Result.User))
		for _, v := range all.Result.User {
			si := &search.Item{}
			si.FromUser(v, avm, lm[v.RoomID])
			if follows[v.Mid] {
				si.Attentions = 1
			}
			items = append(items, si)
		}
		if len(items) == 0 {
			for _, v := range tmp {
				si := &search.Item{}
				si.FromUpUser(v, avm, lm[v.RoomID])
				if follows[v.Mid] {
					si.Attentions = 1
				}
				if old {
					si.IsUp = true
				}
				items = append(items, si)
			}
		}
		res.Items.Upper = items
	}
	items = make([]*search.Item, 0, len(all.Result.Video))
	for _, v := range all.Result.Video {
		si := &search.Item{}
		si.FromVideo(v, avm[v.ID], cooperation)
		items = append(items, si)
	}
	res.Items.Archive = items
	// live room
	if model.IsOverseas(plat) {
		items = make([]*search.Item, 0, len(all.Result.LiveRoom))
		for _, v := range all.Result.LiveRoom {
			si := &search.Item{}
			si.FromLive(v, lm[v.RoomID])
			items = append(items, si)
		}
		res.Items.LiveRoom = items
		// live user
		items = make([]*search.Item, 0, len(all.Result.LiveUser))
		for _, v := range all.Result.LiveUser {
			si := &search.Item{}
			si.FromLive(v, lm[v.RoomID])
			items = append(items, si)
		}
		res.Items.LiveUser = items
	}
	return
}

// SearchByType is tag bangumi movie upuser video search
func (s *Service) SearchByType(c context.Context, mid, zoneid int64, mobiApp, device, platform, buvid, sType, keyword, filtered, order string, plat int8, build, highlight, categoryID, userType, orderSort, pn, ps int, old bool, now time.Time) (res *search.TypeSearch, err error) {
	switch sType {
	case "season":
		if res, err = s.srchDao.Season(c, mid, zoneid, keyword, mobiApp, device, platform, buvid, filtered, plat, build, pn, ps, now); err != nil {
			return
		}
	case "upper":
		if res, err = s.upper(c, mid, keyword, mobiApp, device, platform, buvid, filtered, order, s.biliUserVideoLimit, highlight, build, userType, orderSort, pn, ps, old, now); err != nil {
			return
		}
	case "movie":
		if !model.IsOverseas(plat) {
			if res, err = s.srchDao.MovieByType(c, mid, zoneid, keyword, mobiApp, device, platform, buvid, filtered, plat, build, pn, ps, now); err != nil {
				return
			}
		}
	case "live_room", "live_user":
		if res, err = s.srchDao.LiveByType(c, mid, zoneid, keyword, mobiApp, device, platform, buvid, filtered, order, sType, plat, build, pn, ps, now); err != nil {
			return
		}
	case "article":
		if res, err = s.article(c, mid, zoneid, highlight, keyword, mobiApp, device, platform, buvid, filtered, order, sType, plat, categoryID, build, pn, ps, now); err != nil {
			return
		}
	case "season2":
		if (model.IsAndroid(plat) && build <= s.c.SearchBuildLimit.PGCHighLightAndroid) || (model.IsIOS(plat) && build <= s.c.SearchBuildLimit.PGCHighLightIOS) {
			highlight = 0
		}
		if res, err = s.srchDao.Season2(c, mid, keyword, mobiApp, device, platform, buvid, highlight, build, pn, ps); err != nil {
			return
		}
	case "movie2":
		if !model.IsOverseas(plat) {
			if (model.IsAndroid(plat) && build <= s.c.SearchBuildLimit.PGCHighLightAndroid) || (model.IsIOS(plat) && build <= s.c.SearchBuildLimit.PGCHighLightIOS) {
				highlight = 0
			}
			if res, err = s.srchDao.MovieByType2(c, mid, keyword, mobiApp, device, platform, buvid, highlight, build, pn, ps); err != nil {
				return
			}
		}
	case "tag":
		if res, err = s.channel(c, mid, keyword, mobiApp, platform, buvid, device, order, sType, build, pn, ps, highlight); err != nil {
			return
		}
	case "video":
		if res, err = s.srchDao.Video(c, mid, keyword, mobiApp, device, platform, buvid, highlight, build, pn, ps); err != nil {
			return
		}
	}
	if res == nil {
		res = &search.TypeSearch{Items: []*search.Item{}}
	}
	return
}

// SearchLive is search live
func (s *Service) SearchLive(c context.Context, mid int64, mobiApp, platform, buvid, device, sType, keyword, order string, build, pn, ps int) (res *search.TypeSearch, err error) {
	if res, err = s.srchDao.Live(c, mid, keyword, mobiApp, platform, buvid, device, order, sType, build, pn, ps); err != nil {
		return
	}
	if res == nil {
		res = &search.TypeSearch{Items: []*search.Item{}}
	}
	return
}

// SearchLiveAll is search live
func (s *Service) SearchLiveAll(c context.Context, mid int64, mobiApp, platform, buvid, device, sType, keyword, order string, build, pn, ps int) (res *search.TypeSearchLiveAll, err error) {
	var (
		g         *errgroup.Group
		ctx       context.Context
		uid       int64
		owners    []int64
		glorys    []*live.Glory
		follows   map[int64]bool
		userInfos map[int64]map[string]*live.Exp
	)
	if res, err = s.srchDao.LiveAll(c, mid, keyword, mobiApp, platform, buvid, device, order, sType, build, pn, ps); err != nil {
		return
	}
	if res.Master != nil {
		for _, item := range res.Master.Items {
			uid = item.Mid
			owners = append(owners, uid)
			break
		}
	}
	if len(owners) != 0 {
		g, ctx = errgroup.WithContext(c)
		if uid > 0 {
			g.Go(func() error {
				follows = s.accDao.Relations3(ctx, owners, mid)
				return nil
			})
			g.Go(func() error {
				glorys, _ = s.liveDao.Glory(ctx, uid)
				return nil
			})
			g.Go(func() error {
				userInfos, _ = s.liveDao.UserInfo(ctx, owners)
				return nil
			})
		}
		if err = g.Wait(); err != nil {
			log.Error("%+v", err)
			return
		}
		for _, m := range res.Master.Items {
			if follows[uid] {
				m.IsAttention = 1
			}
			m.Glory = &search.Glory{
				Title: "主播荣誉",
				Total: len(glorys),
				Items: make([]*search.Item, 0, len(glorys)),
			}
			if userInfo, ok := userInfos[m.Mid]; ok {
				if u, ok := userInfo["exp"]; ok {
					if u != nil || u.Master != nil {
						m.Level = u.Master.Level
						m.LevelColor = u.Master.Color
					}
				}
			}
			for _, glory := range glorys {
				if glory.GloryInfo != nil {
					item := &search.Item{
						Title: glory.GloryInfo.Name,
						Cover: glory.GloryInfo.Cover,
					}
					m.Glory.Items = append(m.Glory.Items, item)
				}
			}
		}
	}
	if res == nil {
		res = &search.TypeSearchLiveAll{Master: &search.TypeSearch{Items: []*search.Item{}}, Room: &search.TypeSearch{Items: []*search.Item{}}}
	}
	return
}

// channel search for channel
func (s *Service) channel(c context.Context, mid int64, keyword, mobiApp, platform, buvid, device, order, sType string, build, pn, ps, highlight int) (res *search.TypeSearch, err error) {
	var (
		g          *errgroup.Group
		ctx        context.Context
		tags       []int64
		tagMyInfos []*tagmdl.Tag
	)
	if res, err = s.srchDao.Channel(c, mid, keyword, mobiApp, platform, buvid, device, order, sType, build, pn, ps, highlight); err != nil {
		return
	}
	if res == nil || len(res.Items) == 0 {
		return
	}
	tags = make([]int64, 0, len(res.Items))
	for _, item := range res.Items {
		tags = append(tags, item.ID)
	}
	if len(tags) != 0 {
		g, ctx = errgroup.WithContext(c)
		if mid > 0 {
			g.Go(func() error {
				tagMyInfos, _ = s.tagDao.TagInfos(ctx, tags, mid)
				return nil
			})
		}
		if err = g.Wait(); err != nil {
			log.Error("%+v", err)
			return
		}
		for _, item := range res.Items {
			for _, myInfo := range tagMyInfos {
				if myInfo != nil && myInfo.TagID == item.ID {
					item.IsAttention = myInfo.IsAtten
					break
				}
			}
		}
	}
	return
}

// upper search for upper
func (s *Service) upper(c context.Context, mid int64, keyword, mobiApp, device, platform, buvid, filtered, order string, biliUserVL, highlight, build, userType, orderSort, pn, ps int, old bool, now time.Time) (res *search.TypeSearch, err error) {
	var (
		g       *errgroup.Group
		ctx     context.Context
		owners  []int64
		follows map[int64]bool
	)
	if res, err = s.srchDao.Upper(c, mid, keyword, mobiApp, device, platform, buvid, filtered, order, biliUserVL, highlight, build, userType, orderSort, pn, ps, old, now); err != nil {
		return
	}
	if res == nil || len(res.Items) == 0 {
		return
	}
	owners = make([]int64, 0, len(res.Items))
	for _, item := range res.Items {
		owners = append(owners, item.Mid)
	}
	if len(owners) != 0 {
		g, ctx = errgroup.WithContext(c)
		if mid > 0 {
			g.Go(func() error {
				follows = s.accDao.Relations3(ctx, owners, mid)
				return nil
			})
		}
		if err = g.Wait(); err != nil {
			log.Error("%+v", err)
			return
		}
		for _, item := range res.Items {
			if follows[item.Mid] {
				item.Attentions = 1
			}
		}
	}
	return
}

// article search for article
func (s *Service) article(c context.Context, mid, zoneid int64, highlight int, keyword, mobiApp, device, platform, buvid, filtered, order, sType string, plat int8, categoryID, build, pn, ps int, now time.Time) (res *search.TypeSearch, err error) {
	if res, err = s.srchDao.ArticleByType(c, mid, zoneid, keyword, mobiApp, device, platform, buvid, filtered, order, sType, plat, categoryID, build, highlight, pn, ps, now); err != nil {
		log.Error("%+v", err)
		return
	}
	if res != nil && len(res.Items) > 0 {
		var mids []int64
		for _, v := range res.Items {
			mids = append(mids, v.Mid)
		}
		var infom map[int64]*account.Info
		if infom, err = s.accDao.Infos3(c, mids); err != nil {
			log.Error("%+v", err)
			err = nil
			return
		}
		for _, item := range res.Items {
			if info, ok := infom[item.Mid]; ok {
				item.Name = info.Name
			}
		}
	}
	return
}

// HotSearch is hot word search
func (s *Service) HotSearch(c context.Context, buvid string, mid int64, build, limit int, mobiApp, device, platform string, now time.Time) (res *search.Hot) {
	var err error
	if res, err = s.srchDao.HotSearch(c, buvid, mid, build, limit, mobiApp, device, platform, now); err != nil {
		log.Error("%+v", err)
	}
	if res != nil {
		res.TrackID = res.SeID
		res.SeID = ""
		res.Code = 0
	} else {
		res = &search.Hot{}
	}
	return
}

// Suggest for search suggest
func (s *Service) Suggest(c context.Context, mid int64, buvid, keyword string, build int, mobiApp, device string, now time.Time) (res *search.Suggestion) {
	var (
		suggest *search.Suggest
		err     error
	)
	res = &search.Suggestion{}
	if suggest, err = s.srchDao.Suggest(c, mid, buvid, keyword, build, mobiApp, device, now); err != nil {
		log.Error("%+v", err)
		return
	}
	if suggest != nil {
		res.UpUser = suggest.Result.Accurate.UpUser
		res.Bangumi = suggest.Result.Accurate.Bangumi
		for _, v := range suggest.Result.Tag {
			res.Suggest = append(res.Suggest, v.Value)
		}
		res.TrackID = suggest.Stoken
	}
	return
}

// Suggest2 for search suggest
func (s *Service) Suggest2(c context.Context, mid int64, platform, buvid, keyword string, build int, mobiApp string, now time.Time) (res *search.Suggestion2) {
	var (
		suggest *search.Suggest2
		err     error
		avids   []int64
		avm     map[int64]*api.Arc
		roomIDs []int64
		lm      map[int64]*live.RoomInfo
	)
	res = &search.Suggestion2{}
	if suggest, err = s.srchDao.Suggest2(c, mid, platform, buvid, keyword, build, mobiApp, now); err != nil {
		log.Error("%+v", err)
		return
	}
	if suggest.Result != nil {
		for _, v := range suggest.Result.Tag {
			if v.SpID == search.SuggestionJump {
				if v.Type == search.SuggestionAV {
					avids = append(avids, v.Ref)
				}
				if v.Type == search.SuggestionLive {
					roomIDs = append(roomIDs, v.Ref)
				}
			}
		}
		g, ctx := errgroup.WithContext(c)
		if len(avids) != 0 {
			g.Go(func() (err error) {
				if avm, err = s.arcDao.Archives2(ctx, avids); err != nil {
					log.Error("%+v", err)
					err = nil
				}
				return
			})
		}
		if len(roomIDs) != 0 {
			g.Go(func() (err error) {
				if lm, err = s.liveDao.LiveByRIDs(ctx, roomIDs); err != nil {
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
		for _, v := range suggest.Result.Tag {
			si := &search.Item{}
			si.FromSuggest2(v, avm, lm)
			res.List = append(res.List, si)
		}
		res.TrackID = suggest.Stoken
	}
	return
}

// Suggest3 for search suggest
func (s *Service) Suggest3(c context.Context, mid int64, platform, buvid, keyword, device string, build, highlight int, mobiApp string, now time.Time) (res *search.SuggestionResult3) {
	var (
		suggest *search.Suggest3
		err     error
		avids   []int64
		avm     map[int64]*api.Arc
		roomIDs []int64
		lm      map[int64]*live.RoomInfo
	)
	res = &search.SuggestionResult3{}
	if suggest, err = s.srchDao.Suggest3(c, mid, platform, buvid, keyword, device, build, highlight, mobiApp, now); err != nil {
		log.Error("%+v", err)
		return
	}
	for _, v := range suggest.Result {
		if v.TermType == search.SuggestionJump {
			if v.SubType == search.SuggestionAV {
				avids = append(avids, v.Ref)
			}
			if v.SubType == search.SuggestionLive {
				roomIDs = append(roomIDs, v.Ref)
			}
		}
	}
	g, ctx := errgroup.WithContext(c)
	if len(avids) != 0 {
		g.Go(func() (err error) {
			if avm, err = s.arcDao.Archives2(ctx, avids); err != nil {
				log.Error("%+v", err)
				err = nil
			}
			return
		})
	}
	if len(roomIDs) != 0 {
		g.Go(func() (err error) {
			if lm, err = s.liveDao.LiveByRIDs(ctx, roomIDs); err != nil {
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
	for _, v := range suggest.Result {
		si := &search.Item{}
		si.FromSuggest3(v, avm, lm)
		res.List = append(res.List, si)
	}
	res.TrackID = suggest.TrackID
	return
}

// User for search uer
func (s *Service) User(c context.Context, mid int64, buvid, mobiApp, device, platform, keyword, filtered, order, fromSource string, highlight, build, userType, orderSort, pn, ps int, now time.Time) (res *search.UserResult) {
	res = &search.UserResult{}
	user, err := s.srchDao.User(c, mid, keyword, mobiApp, device, platform, buvid, filtered, order, fromSource, highlight, build, userType, orderSort, pn, ps, now)
	if err != nil {
		log.Error("%+v", err)
		return
	}
	if len(user) == 0 {
		return
	}
	res.Items = make([]*search.Item, 0, len(user))
	for _, u := range user {
		res.Items = append(res.Items, &search.Item{Mid: u.Mid, Name: u.Name, Face: u.Pic})
	}
	return
}

// convertNav deal with old search pageinfo to new.
func (s *Service) convertNav(all *search.Search, plat int8, build int, lang string, old, newPGC bool) (nis []*search.NavInfo) {
	const (
		_showHide          = 0
		_oldAndroidArticle = 515009
	)
	var (
		season   = "番剧"
		live     = "直播"
		upper    = "用户"
		movie    = "影视"
		liveroom = "直播间"
		liveuser = "主播"
		article  = "专栏"
	)
	if old {
		upper = "UP主"
	}
	if lang == model.Hant {
		season = "番劇"
		live = "直播"
		upper = "UP主"
		movie = "影視"
		liveroom = "直播间"
		liveuser = "主播"
		article = "專欄"
	}
	nis = make([]*search.NavInfo, 0, 4)
	// season
	if !newPGC && all.PageInfo.Bangumi != nil {
		var nav = &search.NavInfo{
			Name:  season,
			Total: all.PageInfo.Bangumi.NumResults,
			Pages: all.PageInfo.Bangumi.Pages,
			Type:  1,
		}
		if all.PageInfo.Bangumi.NumResults > s.seasonNum {
			nav.Show = s.seasonShowMore
		} else {
			nav.Show = _showHide
		}
		nis = append(nis, nav)
	}
	// media season
	if newPGC && all.PageInfo.MediaBangumi != nil {
		var nav = &search.NavInfo{
			Name:  season,
			Total: all.PageInfo.MediaBangumi.NumResults,
			Pages: all.PageInfo.MediaBangumi.Pages,
			Type:  7,
		}
		if all.PageInfo.MediaBangumi.NumResults > s.seasonNum {
			nav.Show = s.seasonShowMore
		} else {
			nav.Show = _showHide
		}
		nis = append(nis, nav)
	}
	// live
	if (model.IsAndroid(plat) && build > search.SearchLiveAllAndroid) || (model.IsIPhone(plat) && build > search.SearchLiveAllIOS) || ((plat == model.PlatIPad && build >= search.SearchNewIPad) || (plat == model.PlatIpadHD && build >= search.SearchNewIPadHD)) || model.IsIPhoneB(plat) {
		if all.PageInfo.LiveAll != nil {
			var nav = &search.NavInfo{
				Name:  live,
				Total: all.PageInfo.LiveAll.NumResults,
				Pages: all.PageInfo.LiveAll.Pages,
				Type:  4,
			}
			nis = append(nis, nav)
		}
	} else {
		if all.PageInfo.LiveRoom != nil {
			var nav = &search.NavInfo{
				Name:  live,
				Total: all.PageInfo.LiveRoom.NumResults,
				Pages: all.PageInfo.LiveRoom.Pages,
				Type:  4,
			}
			nis = append(nis, nav)
		}
	}
	// upper
	if old {
		if all.PageInfo.UpUser != nil {
			var nav = &search.NavInfo{
				Name:  upper,
				Total: all.PageInfo.UpUser.NumResults,
				Pages: all.PageInfo.UpUser.Pages,
				Type:  2,
			}
			nis = append(nis, nav)
		}
	} else {
		if all.PageInfo.BiliUser != nil {
			var nav = &search.NavInfo{
				Name:  upper,
				Total: all.PageInfo.BiliUser.NumResults,
				Pages: all.PageInfo.BiliUser.Pages,
				Type:  2,
			}
			nis = append(nis, nav)
		}
	}
	// movie
	if !newPGC && all.PageInfo.Film != nil {
		var nav = &search.NavInfo{
			Name:  movie,
			Total: all.PageInfo.Film.NumResults,
			Pages: all.PageInfo.Film.Pages,
			Type:  3,
		}
		if all.PageInfo.Movie != nil && all.PageInfo.Movie.NumResults > s.movieNum {
			nav.Show = s.movieShowMore
		} else {
			nav.Show = _showHide
		}
		nis = append(nis, nav)
	}
	// media movie
	if newPGC && all.PageInfo.MediaFt != nil {
		var nav = &search.NavInfo{
			Name:  movie,
			Total: all.PageInfo.MediaFt.NumResults,
			Pages: all.PageInfo.MediaFt.Pages,
			Type:  8,
		}
		if all.PageInfo.MediaFt.NumResults > s.movieNum {
			nav.Show = s.movieShowMore
		} else {
			nav.Show = _showHide
		}
		nis = append(nis, nav)
	}
	if model.IsOverseas(plat) {
		// live room
		if all.PageInfo.LiveRoom != nil {
			var nav = &search.NavInfo{
				Name:  liveroom,
				Total: all.PageInfo.LiveRoom.NumResults,
				Pages: all.PageInfo.LiveRoom.Pages,
				Type:  4,
			}
			nis = append(nis, nav)
		}
		if all.PageInfo.LiveUser != nil {
			var nav = &search.NavInfo{
				Name:  liveuser,
				Total: all.PageInfo.LiveUser.NumResults,
				Pages: all.PageInfo.LiveUser.Pages,
				Type:  5,
			}
			nis = append(nis, nav)
		}
	} else {
		if all.PageInfo.Article != nil {
			if (model.IsIPhone(plat) && build > _oldIOS) || (model.IsAndroid(plat) && build > _oldAndroidArticle) || model.IsIPhoneB(plat) {
				var nav = &search.NavInfo{
					Name:  article,
					Total: all.PageInfo.Article.NumResults,
					Pages: all.PageInfo.Article.Pages,
					Type:  6,
				}
				nis = append(nis, nav)
			}
		}
	}
	return
}

// RecommendNoResult search when no result
func (s *Service) RecommendNoResult(c context.Context, platform, mobiApp, device, buvid, keyword string, build, pn, ps int, mid int64) (res *search.NoResultRcndResult, err error) {
	if res, err = s.srchDao.RecommendNoResult(c, platform, mobiApp, device, buvid, keyword, build, pn, ps, mid); err != nil {
		log.Error("%+v", err)
	}
	return
}

// Recommend search recommend
func (s *Service) Recommend(c context.Context, mid int64, build, from, show int, buvid, platform, mobiApp, device string) (res *search.RecommendResult, err error) {
	if res, err = s.srchDao.Recommend(c, mid, build, from, show, buvid, platform, mobiApp, device); err != nil {
		log.Error("%+v", err)
	}
	return
}

// DefaultWords search for default words
func (s *Service) DefaultWords(c context.Context, mid int64, build, from int, buvid, platform, mobiApp, device string) (res *search.DefaultWords, err error) {
	if res, err = s.srchDao.DefaultWords(c, mid, build, from, buvid, platform, mobiApp, device); err != nil {
		log.Error("%+v", err)
	}
	return
}

// Resource for rsource
func (s *Service) Resource(c context.Context, mobiApp, device, network, buvid, adExtra string, build int, plat int8, mid int64) (res []*banner.Banner, err error) {
	var (
		bnsm  map[int][]*resmdl.Banner
		resID int
	)
	if model.IsAndroid(plat) {
		resID = AndroidSearchResourceID
	} else if model.IsIPhone(plat) {
		resID = IPhoneSearchResourceID
	} else if model.IsIPad(plat) {
		resID = IPadSearchResourceID
	}
	if bnsm, err = s.resDao.Banner(c, mobiApp, device, network, "", buvid, adExtra, strconv.Itoa(resID), build, plat, mid); err != nil {
		return
	}
	// only one position
	for _, rb := range bnsm[resID] {
		b := &banner.Banner{}
		b.ChangeBanner(rb)
		res = append(res, b)
		break
	}
	return
}

// RecommendPre search at pre-page.
func (s *Service) RecommendPre(c context.Context, platform, mobiApp, device, buvid string, build, ps int, mid int64) (res *search.RecommendPreResult, err error) {
	if res, err = s.srchDao.RecommendPre(c, platform, mobiApp, device, buvid, build, ps, mid); err != nil {
		log.Error("%+v", err)
	}
	return
}

// SearchEpisodes search PGC episodes
func (s *Service) SearchEpisodes(c context.Context, mid, ssID int64) (res []*search.Item, err error) {
	var (
		seasonIDs []int64
		bangumis  map[string]*bangumi.Card
	)
	seasonIDs = []int64{ssID}
	if bangumis, err = s.bangumiDao.Card(c, mid, seasonIDs); err != nil {
		log.Error("%+v", err)
		return
	}
	if bangumi, ok := bangumis[strconv.FormatInt(ssID, 10)]; ok {
		for _, v := range bangumi.Episodes {
			tmp := &search.Item{
				Param:  strconv.Itoa(int(v.ID)),
				Index:  v.Index,
				Badges: v.Badges,
			}
			tmp.URI = model.FillURI(model.GotoEP, tmp.Param, nil)
			res = append(res, tmp)
		}
	}
	return
}
