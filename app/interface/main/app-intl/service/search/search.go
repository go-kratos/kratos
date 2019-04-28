package search

import (
	"context"
	"time"

	tag "go-common/app/interface/main/app-interface/model/tag"
	"go-common/app/interface/main/app-intl/conf"
	accdao "go-common/app/interface/main/app-intl/dao/account"
	arcdao "go-common/app/interface/main/app-intl/dao/archive"
	artdao "go-common/app/interface/main/app-intl/dao/article"
	bgmdao "go-common/app/interface/main/app-intl/dao/bangumi"
	resdao "go-common/app/interface/main/app-intl/dao/resource"
	srchdao "go-common/app/interface/main/app-intl/dao/search"
	tagdao "go-common/app/interface/main/app-intl/dao/tag"
	"go-common/app/interface/main/app-intl/model"
	"go-common/app/interface/main/app-intl/model/bangumi"
	"go-common/app/interface/main/app-intl/model/search"
	article "go-common/app/interface/openplatform/article/model"
	account "go-common/app/service/main/account/model"
	"go-common/app/service/main/archive/api"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

var (
	_emptyResult = &search.Result{
		NavInfo: []*search.NavInfo{},
		Page:    0,
	}
)

// Service is search service
type Service struct {
	c       *conf.Config
	srchDao *srchdao.Dao
	accDao  *accdao.Dao
	arcDao  *arcdao.Dao
	artDao  *artdao.Dao
	// artDao     *artdao.Dao
	resDao *resdao.Dao
	tagDao *tagdao.Dao
	bgmDao *bgmdao.Dao
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
		c:       c,
		srchDao: srchdao.New(c),
		accDao:  accdao.New(c),
		arcDao:  arcdao.New(c),
		artDao:  artdao.New(c),
		// artDao:             artdao.New(c),
		resDao:             resdao.New(c),
		tagDao:             tagdao.New(c),
		bgmDao:             bgmdao.New(c),
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
func (s *Service) Search(c context.Context, mid, zoneid int64, mobiApp, device, platform, buvid, keyword, duration, order, filtered, lang, fromSource, recommend string, plat int8, rid, highlight, build, pn, ps int, now time.Time) (res *search.Result, err error) {
	var (
		aids      []int64
		am        map[int64]*api.Arc
		owners    []int64
		follows   map[int64]bool
		seasonIDs []int64
		bangumis  map[string]*bangumi.Card
	)
	var (
		seasonNum int
		movieNum  int
	)
	seasonNum = s.seasonNum
	movieNum = s.movieNum
	all, code, err := s.srchDao.Search(c, mid, zoneid, mobiApp, device, platform, buvid, keyword, duration, order, filtered, fromSource, recommend, plat, seasonNum, movieNum, s.upUserNum, s.uvLimit, s.userNum, s.userVideoLimit, s.biliUserNum, s.biliUserVideoLimit, rid, highlight, build, pn, ps, now)
	if err != nil {
		log.Error("%+v", err)
		return
	}
	if code == model.ForbidCode || code == model.NoResultCode {
		res = _emptyResult
		return
	}
	res = &search.Result{}
	res.Trackid = all.Trackid
	res.Page = all.Page
	res.Array = all.FlowPlaceholder
	res.Attribute = all.Attribute
	res.NavInfo = s.convertNav(all, plat, build, lang)
	if len(all.FlowResult) != 0 {
		var item []*search.Item
		for _, v := range all.FlowResult {
			switch v.Type {
			case search.TypeUser, search.TypeBiliUser:
				owners = append(owners, v.User.Mid)
				for _, vr := range v.User.Res {
					aids = append(aids, vr.Aid)
				}
			case search.TypeVideo:
				aids = append(aids, v.Video.ID)
			case search.TypeMediaBangumi, search.TypeMediaFt:
				seasonIDs = append(seasonIDs, v.Media.SeasonID)
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
		if len(aids) != 0 {
			g.Go(func() (err error) {
				if am, err = s.arcDao.Archives(ctx, aids); err != nil {
					log.Error("%+v", err)
					err = nil
				}
				return
			})
		}
		if len(seasonIDs) != 0 {
			g.Go(func() (err error) {
				if bangumis, err = s.bgmDao.Card(ctx, mid, seasonIDs); err != nil {
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
		if all.SuggestKeyword != "" && pn == 1 {
			i := &search.Item{Title: all.SuggestKeyword, Goto: model.GotoSuggestKeyWord}
			item = append(item, i)
		}
		for _, v := range all.FlowResult {
			i := &search.Item{TrackID: v.TrackID, LinkType: v.LinkType, Position: v.Position}
			switch v.Type {
			case search.TypeVideo:
				i.FromVideo(v.Video, am[v.Video.ID])
			case search.TypeMediaBangumi:
				i.FromMedia(v.Media, "", model.GotoBangumi, bangumis)
			case search.TypeMediaFt:
				i.FromMedia(v.Media, "", model.GotoMovie, bangumis)
			case search.TypeSpecial:
				i.FromOperate(v.Operate, model.GotoSpecial)
			case search.TypeBanner:
				i.FromOperate(v.Operate, model.GotoBanner)
			case search.TypeUser:
				if follows[v.User.Mid] {
					i.Attentions = 1
				}
				i.FromUser(v.User, am)
			case search.TypeBiliUser:
				if follows[v.User.Mid] {
					i.Attentions = 1
				}
				i.FromUpUser(v.User, am)
			case search.TypeSpecialS:
				i.FromOperate(v.Operate, model.GotoSpecialS)
			case search.TypeQuery:
				i.Title = v.TypeName
				i.FromQuery(v.Query)
			case search.TypeConverge:
				var (
					avids, artids []int64
					avm           map[int64]*api.Arc
					artm          map[int64]*article.Meta
				)
				for _, c := range v.Operate.ContentList {
					switch c.Type {
					case 0:
						avids = append(avids, c.ID)
					case 2:
						artids = append(artids, c.ID)
					}
				}
				g, ctx := errgroup.WithContext(c)
				if len(aids) != 0 {
					g.Go(func() (err error) {
						if avm, err = s.arcDao.Archives(ctx, avids); err != nil {
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
				i.FromConverge(v.Operate, avm, artm)
			case search.TypeTwitter:
				i.FromTwitter(v.Twitter)
			}
			if i.Goto != "" {
				item = append(item, i)
			}
		}
		res.Item = item
		if all.EggInfo != nil {
			res.EasterEgg = &search.EasterEgg{ID: all.EggInfo.Source, ShowCount: all.EggInfo.ShowCount}
		}
		return
	}
	// archive
	for _, v := range all.Result.Video {
		aids = append(aids, v.ID)
	}
	if duration == "0" && order == "totalrank" && rid == 0 {
		for _, v := range all.Result.Movie {
			if v.Type == "movie" {
				aids = append(aids, v.Aid)
			}
		}
	}
	if pn == 1 {
		for _, v := range all.Result.User {
			for _, vr := range v.Res {
				aids = append(aids, vr.Aid)
			}
		}
		for _, v := range all.Result.BiliUser {
			for _, vr := range v.Res {
				aids = append(aids, vr.Aid)
			}
			owners = append(owners, v.Mid)
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
	if len(aids) != 0 {
		g.Go(func() (err error) {
			if am, err = s.arcDao.Archives(ctx, aids); err != nil {
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
	return
}

// SearchByType is tag bangumi movie upuser video search
func (s *Service) SearchByType(c context.Context, mid, zoneid int64, mobiApp, device, platform, buvid, sType, keyword, filtered, order string, plat int8, build, highlight, categoryID, userType, orderSort, pn, ps int, now time.Time) (res *search.TypeSearch, err error) {
	switch sType {
	case "upper":
		if res, err = s.upper(c, mid, keyword, mobiApp, device, platform, buvid, filtered, order, s.biliUserVideoLimit, highlight, build, userType, orderSort, pn, ps, now); err != nil {
			return
		}
	case "article":
		if res, err = s.article(c, mid, zoneid, highlight, keyword, mobiApp, device, platform, buvid, filtered, order, sType, plat, categoryID, build, pn, ps, now); err != nil {
			return
		}
	case "season2":
		if res, err = s.srchDao.Season2(c, mid, keyword, mobiApp, device, platform, buvid, highlight, build, pn, ps); err != nil {
			return
		}
	case "movie2":
		if res, err = s.srchDao.MovieByType2(c, mid, keyword, mobiApp, device, platform, buvid, highlight, build, pn, ps); err != nil {
			return
		}
	case "tag":
		if res, err = s.channel(c, mid, keyword, mobiApp, platform, buvid, device, order, sType, build, pn, ps, highlight); err != nil {
			return
		}
	}
	if res == nil {
		res = &search.TypeSearch{Items: []*search.Item{}}
	}
	return
}

// Suggest3 for search suggest
func (s *Service) Suggest3(c context.Context, mid int64, platform, buvid, keyword string, build, highlight int, mobiApp string, now time.Time) (res *search.SuggestionResult3) {
	var (
		suggest *search.Suggest3
		err     error
		aids    []int64
		am      map[int64]*api.Arc
	)
	res = &search.SuggestionResult3{}
	if suggest, err = s.srchDao.Suggest3(c, mid, platform, buvid, keyword, build, highlight, mobiApp, now); err != nil {
		log.Error("%+v", err)
		return
	}
	for _, v := range suggest.Result {
		if v.TermType == search.SuggestionJump {
			if v.SubType == search.SuggestionAV {
				aids = append(aids, v.Ref)
			}
		}
	}
	g, ctx := errgroup.WithContext(c)
	if len(aids) != 0 {
		g.Go(func() (err error) {
			if am, err = s.arcDao.Archives(ctx, aids); err != nil {
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
		si.FromSuggest3(v, am)
		res.List = append(res.List, si)
	}
	res.TrackID = suggest.TrackID
	return
}

// convertNav deal with old search pageinfo to new.
func (s *Service) convertNav(all *search.Search, plat int8, build int, lang string) (nis []*search.NavInfo) {
	const (
		_showHide = 0
	)
	var (
		season  = "番剧"
		upper   = "用户"
		movie   = "影视"
		article = "专栏"
	)
	if lang == model.Hant {
		season = "番劇"
		upper = "UP主"
		movie = "影視"
		article = "專欄"
	}
	nis = make([]*search.NavInfo, 0, 4)
	// season
	// media season
	if all.PageInfo.MediaBangumi != nil {
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
	// upper
	if all.PageInfo.BiliUser != nil {
		var nav = &search.NavInfo{
			Name:  upper,
			Total: all.PageInfo.BiliUser.NumResults,
			Pages: all.PageInfo.BiliUser.Pages,
			Type:  2,
		}
		nis = append(nis, nav)
	}
	// media movie
	if all.PageInfo.MediaFt != nil {
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
	if all.PageInfo.Article != nil {
		var nav = &search.NavInfo{
			Name:  article,
			Total: all.PageInfo.Article.NumResults,
			Pages: all.PageInfo.Article.Pages,
			Type:  6,
		}
		nis = append(nis, nav)
	}
	return
}

// upper search for upper
func (s *Service) upper(c context.Context, mid int64, keyword, mobiApp, device, platform, buvid, filtered, order string, biliUserVL, highlight, build, userType, orderSort, pn, ps int, now time.Time) (res *search.TypeSearch, err error) {
	var (
		owners  []int64
		follows map[int64]bool
	)
	if res, err = s.srchDao.Upper(c, mid, keyword, mobiApp, device, platform, buvid, filtered, order, biliUserVL, highlight, build, userType, orderSort, pn, ps, now); err != nil {
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
		g, ctx := errgroup.WithContext(c)
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
	if res != nil && len(res.Items) != 0 {
		return
	}
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
	return
}

// channel search for channel
func (s *Service) channel(c context.Context, mid int64, keyword, mobiApp, platform, buvid, device, order, sType string, build, pn, ps, highlight int) (res *search.TypeSearch, err error) {
	var (
		g          *errgroup.Group
		ctx        context.Context
		tags       []int64
		tagMyInfos []*tag.Tag
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
