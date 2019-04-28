package search

import (
	"context"
	"go-common/library/sync/errgroup"
	"go-common/library/xstr"
	"net/http"
	"net/url"
	"strconv"

	mdlSearch "go-common/app/interface/main/tv/model/search"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

// UserSearch search user .
func (d *Dao) UserSearch(ctx context.Context, arg *mdlSearch.UserSearch) (user []*mdlSearch.User, err error) {
	var (
		params = url.Values{}
		req    *http.Request
		ip     = metadata.String(ctx, metadata.RemoteIP)
	)
	params.Set("platform", arg.Platform)
	params.Set("mobi_app", arg.MobiAPP)
	params.Set("build", arg.Build)
	params.Set("keyword", arg.Keyword)
	params.Set("page", strconv.Itoa(arg.Page))
	params.Set("pagesize", strconv.Itoa(arg.Pagesize))
	params.Set("userid", strconv.FormatInt(arg.MID, 10))
	params.Set("order", arg.Order)
	params.Set("main_ver", "v3") // 支持赛事搜索
	params.Set("search_type", "bili_user")
	params.Set("user_type", strconv.Itoa(arg.UserType)) // 用户类型
	params.Set("highlight", strconv.Itoa(arg.Highlight))
	params.Set("order_sort", strconv.Itoa(arg.OrderSort))
	params.Set("from_source", arg.FromSource)
	params.Set("bili_user_vl", strconv.Itoa(d.cfgWild.BiliUserVl))
	// new request
	if req, err = d.client.NewRequest("GET", d.userSearch, ip, params); err != nil {
		log.Error("[wild.UserSearch] d.client.NewRequest url(%s) error(%v)", d.userSearch, err)
		return
	}
	req.Header.Set("Buvid", arg.Buvid)
	var res struct {
		Code  int               `json:"code"`
		SeID  string            `json:"seid"`
		Pages int               `json:"numPages"`
		List  []*mdlSearch.User `json:"result"`
	}
	if err = d.client.Do(ctx, req, &res); err != nil {
		log.Error("[wild.UserSearch] d.client.Do url(%s) error(%v)", d.userSearch, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		err = ecode.Int(res.Code)
		log.Error("[wild.UserSearch] url(%s) error(%v)", d.userSearch, err)
		return
	}
	user = res.List
	return
}

// SearchAllWild wild search all .
func (d *Dao) SearchAllWild(ctx context.Context, arg *mdlSearch.UserSearch) (user *mdlSearch.Search, err error) {
	var (
		req *http.Request
		ip  = metadata.String(ctx, metadata.RemoteIP)
	)
	params := url.Values{}
	user = &mdlSearch.Search{}
	params.Set("build", arg.Build)
	params.Set("keyword", arg.Keyword)
	params.Set("main_ver", "v3")
	params.Set("mobi_app", arg.MobiAPP)
	params.Set("device", arg.Device)
	params.Set("userid", strconv.FormatInt(arg.MID, 10))
	params.Set("tids", strconv.Itoa(arg.RID))
	params.Set("highlight", strconv.Itoa(arg.Highlight))
	params.Set("page", strconv.Itoa(arg.Page))
	params.Set("pagesize", strconv.Itoa(arg.Pagesize))
	params.Set("bili_user_num", strconv.Itoa(d.cfgWild.BiliUserNum))
	params.Set("bili_user_vl", strconv.Itoa(d.cfgWild.BiliUserVl))
	params.Set("user_num", strconv.Itoa(d.cfgWild.UserNum))
	params.Set("user_video_limit", strconv.Itoa(d.cfgWild.UserVideoLimit))
	params.Set("platform", arg.Platform)
	// params.Set("duration", strconv.Itoa(arg.Duration)) // 视频时长筛选，默认是0
	params.Set("order", arg.Order)
	params.Set("search_type", "all")
	params.Set("from_source", arg.FromSource)
	params.Set("media_bangumi_num", strconv.Itoa(arg.SeasonNum))
	params.Set("movie_num", strconv.Itoa(arg.MovieNum))
	params.Set("is_new_pgc", "1") // 新番剧
	params.Set("media_ft_num", strconv.Itoa(arg.MovieNum))
	// params.Set("flow_need", "1")      // 混排
	// params.Set("query_rec_need", "1") // 搜索结果推荐
	// new request
	if req, err = d.client.NewRequest("GET", d.userSearch, ip, params); err != nil {
		log.Error("d.client.NewRequest URI(%s) error(%v)", d.userSearch, err)
		return
	}
	req.Header.Set("Buvid", arg.Buvid)
	if err = d.client.Do(ctx, req, user); err != nil {
		log.Error("[wild.SearchAllWild] d.client.Do() url(%s) error(%v)", d.userSearch, err)
	}
	return
}

// card bangumi card .
func (d *Dao) cardInfo(c context.Context, mid int64, sids []int64) (s map[string]*mdlSearch.Card, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("season_ids", xstr.JoinInts(sids))
	var res struct {
		Code   int                        `json:"code"`
		Result map[string]*mdlSearch.Card `json:"result"`
	}
	if err = d.client.Get(c, d.card, ip, params, &res); err != nil {
		log.Error("d.client.NewRequest url(%s) error(%v)", d.resultURL, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		err = ecode.Int(res.Code)
		log.Error("[wild] cardInfo error(%v)", err)
		return
	}
	s = res.Result
	return
}

// PgcSearch .
func (d *Dao) PgcSearch(c context.Context, arg *mdlSearch.UserSearch) (st *mdlSearch.TypeSearch, err error) {
	var (
		req            *http.Request
		ip             = metadata.String(c, metadata.RemoteIP)
		seasonIDs      []int64
		bangumis       map[string]*mdlSearch.Card
		items1, items2 []*mdlSearch.Item
	)
	params := url.Values{}
	params.Set("build", arg.Build)
	params.Set("keyword", arg.Keyword)
	params.Set("main_ver", "v3")
	params.Set("mobi_app", arg.MobiAPP)
	params.Set("device", arg.Device)
	params.Set("userid", strconv.FormatInt(arg.MID, 10))
	params.Set("highlight", strconv.Itoa(arg.Highlight))
	params.Set("page", strconv.Itoa(arg.Page))
	params.Set("pagesize", strconv.Itoa(arg.Pagesize))
	params.Set("platform", arg.Platform)
	params.Set("order", arg.Order)
	params.Set("search_type", "all")
	params.Set("from_source", arg.FromSource)
	params.Set("media_bangumi_num", strconv.Itoa(arg.SeasonNum))
	params.Set("media_ft_num", strconv.Itoa(arg.MovieNum))
	params.Set("is_new_pgc", "1")
	if req, err = d.client.NewRequest("GET", d.userSearch, ip, params); err != nil {
		log.Error("d.client.NewRequest url(%s) error(%v)", d.userSearch, err)
		return
	}
	req.Header.Set("Buvid", arg.Buvid)
	res := &mdlSearch.Search{}
	if err = d.client.Do(c, req, res); err != nil {
		log.Error("[wild.PgcSearch] d.client.Do url(%s) error(%v)", d.userSearch, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		err = ecode.Int(res.Code)
		log.Error("[wild.PgcSearch] code(%d) error(%v)", res.Code, err)
		return
	}
	for _, v := range res.Result.MediaBangumi {
		seasonIDs = append(seasonIDs, v.SeasonID)
	}
	for _, v := range res.Result.MediaFt {
		seasonIDs = append(seasonIDs, v.SeasonID)
	}
	if len(seasonIDs) > 0 {
		if bangumis, err = d.cardInfo(c, arg.MID, seasonIDs); err != nil {
			log.Error("[wild.PgcSearch] MovieByType2 %+v", err)
			return
		}
	}
	if len(bangumis) > 0 {
		group := new(errgroup.Group)
		group.Go(func() error {
			items1 = make([]*mdlSearch.Item, 0, len(res.Result.MediaBangumi))
			for _, v := range res.Result.MediaBangumi {
				si := &mdlSearch.Item{}
				si.FromMedia(v, "", mdlSearch.GotoMovie, bangumis)
				items1 = append(items1, si)
			}
			return nil
		})
		group.Go(func() error {
			items2 = make([]*mdlSearch.Item, 0, len(res.Result.MediaFt))
			for _, v := range res.Result.MediaFt {
				si := &mdlSearch.Item{}
				si.FromMedia(v, "", mdlSearch.GotoMovie, bangumis)
				items2 = append(items2, si)
			}
			return nil
		})
		if err = group.Wait(); err != nil {
			log.Error("[wild.PgcSearch] group.Wait() is error(%v)", err)
			return
		}
	}
	items1 = append(items1, items2...)
	st = &mdlSearch.TypeSearch{TrackID: res.Trackid, Pages: res.Page, Total: res.Total, Items: items1}
	return
}
