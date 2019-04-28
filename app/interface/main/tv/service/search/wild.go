package search

import (
	"context"

	mdlSearch "go-common/app/interface/main/tv/model/search"
	v1 "go-common/app/service/main/archive/api"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

const (
	_showHide         = 0
	season            = "番剧"
	upper             = "用户"
	movie             = "影视"
	_searchType       = "all"
	_mobiAPP          = "app"
	_bangumiType      = 1
	_biliUserType     = 2
	_filmType         = 3
	_mediaBangumiType = 7
	_mediaFtType      = 8
)

// UserSearch search user .
func (s *Service) UserSearch(ctx context.Context, arg *mdlSearch.UserSearch) (res []*mdlSearch.User, err error) {
	if res, err = s.dao.UserSearch(ctx, arg); err != nil {
		log.Error("s.dao.UserSearch error(%v)", err)
	}
	if len(res) == 0 {
		res = make([]*mdlSearch.User, 0)
	}
	return
}

// SearchAll search all .
func (s *Service) SearchAll(ctx context.Context, arg *mdlSearch.UserSearch) (res *mdlSearch.ResultAll, err error) {
	var (
		user    = &mdlSearch.Search{}
		avm     map[int64]*v1.Arc
		avids   []int64
		items   []*mdlSearch.Item
		wildCfg = s.conf.Wild.WildSearch
	)
	arg.SeasonNum = wildCfg.SeasonNum
	arg.MovieNum = wildCfg.MovieNum
	arg.SearchType = _searchType
	arg.MobiAPP = _mobiAPP
	if user, err = s.dao.SearchAllWild(ctx, arg); err != nil {
		log.Error(" s.dao.SearchAllWild error(%v)", err)
	}
	res = &mdlSearch.ResultAll{}
	if user == nil {
		return
	}
	res.Trackid = user.Trackid
	res.Page = user.Page
	res.Attribute = user.Attribute
	nis := make([]*mdlSearch.NavInfo, 0, 4)
	// season
	if user.PageInfo.Bangumi != nil {
		var nav = &mdlSearch.NavInfo{
			Name:  season,
			Total: user.PageInfo.Bangumi.NumResult,
			Pages: user.PageInfo.Bangumi.Pages,
			Type:  _bangumiType,
		}
		if user.PageInfo.Bangumi.NumResult > wildCfg.SeasonNum {
			nav.Show = wildCfg.SeasonMore
		} else {
			nav.Show = _showHide
		}
		nis = append(nis, nav)
	}
	// media season
	if user.PageInfo.MediaBangumi != nil {
		var nav = &mdlSearch.NavInfo{
			Name:  season,
			Total: user.PageInfo.MediaBangumi.NumResult,
			Pages: user.PageInfo.MediaBangumi.Pages,
			Type:  _mediaBangumiType,
		}
		if user.PageInfo.MediaBangumi.NumResult > wildCfg.SeasonNum {
			nav.Show = wildCfg.SeasonMore
		} else {
			nav.Show = _showHide
		}
		nis = append(nis, nav)
	}
	// upper
	if user.PageInfo.BiliUser != nil {
		var nav = &mdlSearch.NavInfo{
			Name:  upper,
			Total: user.PageInfo.BiliUser.NumResult,
			Pages: user.PageInfo.BiliUser.Pages,
			Type:  _biliUserType,
		}
		nis = append(nis, nav)
	}
	// movie
	if user.PageInfo.Film != nil {
		var nav = &mdlSearch.NavInfo{
			Name:  movie,
			Total: user.PageInfo.Film.NumResult,
			Pages: user.PageInfo.Film.Pages,
			Type:  _filmType,
		}
		if user.PageInfo.Movie != nil && user.PageInfo.Movie.NumResult > wildCfg.MovieNum {
			nav.Show = wildCfg.MovieMore
		} else {
			nav.Show = _showHide
		}
		nis = append(nis, nav)
	}
	// media movie
	if user.PageInfo.MediaFt != nil {
		var nav = &mdlSearch.NavInfo{
			Name:  movie,
			Total: user.PageInfo.MediaFt.NumResult,
			Pages: user.PageInfo.MediaFt.Pages,
			Type:  _mediaFtType,
		}
		if user.PageInfo.MediaFt.NumResult > wildCfg.MovieNum {
			nav.Show = wildCfg.MovieMore
		} else {
			nav.Show = _showHide
		}
		nis = append(nis, nav)
	}
	res.NavInfo = nis
	// archive
	for _, v := range user.Result.Video {
		avids = append(avids, v.ID)
	}
	for _, v := range user.Result.Movie {
		if v.Type == "movie" {
			avids = append(avids, v.Aid)
		}
	}
	if arg.Page == 1 {
		for _, v := range user.Result.User {
			for _, vr := range v.Res {
				avids = append(avids, vr.Aid)
			}
		}
		for _, v := range user.Result.BiliUser {
			for _, vr := range v.Res {
				avids = append(avids, vr.Aid)
			}
		}
	}
	group := new(errgroup.Group)
	if len(avids) != 0 {
		group.Go(func() (err error) {
			if avm, err = s.arcDao.Archives(ctx, avids); err != nil {
				log.Error("%+v", err)
				err = nil
			}
			return
		})
	}
	if err = group.Wait(); err != nil {
		return
	}

	// item add data .
	var promptBangumi, promptFt string
	// season
	bangumi := user.Result.Bangumi
	items = make([]*mdlSearch.Item, 0, len(bangumi))
	for _, v := range bangumi {
		si := &mdlSearch.Item{}
		si.FromSeason(v, mdlSearch.GotoBangumiWeb)
		items = append(items, si)
	}
	if len(res.Items.Season) > 5 && arg.RID == 0 {
		res.Items.Season = res.Items.Season[:5]
		res.Items.Season = items
	} else {
		res.Items.Season = items
	}
	// movie
	movie := user.Result.Movie
	items = make([]*mdlSearch.Item, 0, len(movie))
	for _, v := range movie {
		si := &mdlSearch.Item{}
		si.FromMovie(v, avm)
		items = append(items, si)
	}
	res.Items.Movie = items
	// season2
	mb := user.Result.MediaBangumi
	// movie2
	mf := user.Result.MediaFt
	items = make([]*mdlSearch.Item, 0, len(mb)+len(mf))
	for _, v := range mb {
		si := &mdlSearch.Item{}
		si.FromMedia(v, promptBangumi, mdlSearch.GotoBangumi, nil)
		items = append(items, si)
	}
	for _, v := range mf {
		si := &mdlSearch.Item{}
		si.FromMedia(v, promptFt, mdlSearch.GotoMovie, nil)
		si.Goto = mdlSearch.GotoAv
		items = append(items, si)
	}
	if len(res.Items.Season2) > 5 && arg.RID == 0 {
		res.Items.Season2 = res.Items.Season2[:5]
		res.Items.Season2 = items
	} else {
		res.Items.Season2 = items
	}

	items = make([]*mdlSearch.Item, 0, len(user.Result.Video))
	for _, v := range user.Result.Video {
		si := &mdlSearch.Item{}
		si.FromVideo(v, avm[v.ID])
		items = append(items, si)
	}
	res.Items.Archive = items
	return
}

// PgcSearch search .
func (s *Service) PgcSearch(ctx context.Context, arg *mdlSearch.UserSearch) (res *mdlSearch.TypeSearch, err error) {
	var (
		wildCfg = s.conf.Wild.WildSearch
	)
	arg.SeasonNum = wildCfg.SeasonNum
	arg.MovieNum = wildCfg.MovieNum
	arg.SearchType = _searchType
	arg.MobiAPP = _mobiAPP
	if res, err = s.dao.PgcSearch(ctx, arg); err != nil {
		log.Error("[wild.PgcSearch] s.dao.PgcSearch error(%v)", err)
		return
	}
	if len(res.Items) <= 0 {
		res.Items = make([]*mdlSearch.Item, 0)
	}
	return
}
