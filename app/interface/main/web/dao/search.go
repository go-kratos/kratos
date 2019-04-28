package dao

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"go-common/app/interface/main/web/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/xstr"
)

const (
	_searchVer       = "v3"
	_searchPlatform  = "web"
	_searchUpRecType = "up_rec"
)

// SearchAll search all data.
func (d *Dao) SearchAll(c context.Context, mid int64, arg *model.SearchAllArg, buvid, ua, typ string) (res *model.Search, err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params = setSearchParam(params, mid, model.SearchTypeAll, arg.Keyword, _searchPlatform, arg.FromSource, buvid, ip)
	params.Set("duration", strconv.Itoa(arg.Duration))
	params.Set("page", strconv.Itoa(arg.Pn))
	params.Set("tids", strconv.Itoa(arg.Rid))
	if typ == model.WxSearchType {
		params.Set("highlight", strconv.Itoa(arg.Highlight))
		for k, v := range model.SearchDefaultArg[model.WxSearchTypeAll] {
			params.Set(k, strconv.Itoa(v))
		}
	} else {
		for k, v := range model.SearchDefaultArg[model.SearchTypeAll] {
			params.Set(k, strconv.Itoa(v))
		}
		params.Set("single_column", strconv.Itoa(arg.SingleColumn))
	}
	var req *http.Request
	if req, err = d.httpSearch.NewRequest(http.MethodGet, d.searchURL, ip, params); err != nil {
		return
	}
	req.Header.Set("browser-info", ua)
	res = new(model.Search)
	if err = d.httpSearch.Do(c, req, &res); err != nil {
		log.Error("Search d.httpSearch.Get(%s) error(%v)", d.searchURL, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("Search d.httpSearch.Get(%s) code(%d) error", d.searchURL, res.Code)
		err = ecode.Int(res.Code)
	}
	return
}

// SearchVideo search season data.
func (d *Dao) SearchVideo(c context.Context, mid int64, arg *model.SearchTypeArg, buvid, ua string) (res *model.SearchTypeRes, err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params = setSearchTypeParam(params, mid, model.SearchTypeVideo, buvid, ip, arg)
	params.Set("duration", strconv.Itoa(arg.Duration))
	params.Set("order", arg.Order)
	params.Set("from_source", arg.FromSource)
	params.Set("tids", strconv.FormatInt(arg.Rid, 10))
	params.Set("page", strconv.Itoa(arg.Pn))
	for k, v := range model.SearchDefaultArg[model.SearchTypeVideo] {
		params.Set(k, strconv.Itoa(v))
	}
	var req *http.Request
	if req, err = d.httpSearch.NewRequest(http.MethodGet, d.searchURL, ip, params); err != nil {
		return
	}
	req.Header.Set("browser-info", ua)
	res = new(model.SearchTypeRes)
	if err = d.httpSearch.Do(c, req, &res); err != nil {
		log.Error("SearchVideo d.httpSearch.Get(%s) error(%v)", d.searchURL, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("SearchVideo d.httpSearch.Get(%s) code(%d) error", d.searchURL, res.Code)
		err = ecode.Int(res.Code)
	}
	return
}

// SearchBangumi search bangumi data.
func (d *Dao) SearchBangumi(c context.Context, mid int64, arg *model.SearchTypeArg, buvid, ua string) (res *model.SearchTypeRes, err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params = setSearchTypeParam(params, mid, model.SearchTypeBangumi, buvid, ip, arg)
	params.Set("duration", strconv.Itoa(arg.Duration))
	params.Set("order", arg.Order)
	params.Set("page", strconv.Itoa(arg.Pn))
	for k, v := range model.SearchDefaultArg[model.SearchTypeBangumi] {
		params.Set(k, strconv.Itoa(v))
	}
	var req *http.Request
	if req, err = d.httpSearch.NewRequest(http.MethodGet, d.searchURL, ip, params); err != nil {
		return
	}
	req.Header.Set("browser-info", ua)
	res = new(model.SearchTypeRes)
	if err = d.httpSearch.Do(c, req, &res); err != nil {
		log.Error("SearchBangumi d.httpSearch.Get(%s) error(%v)", d.searchURL, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("SearchBangumi d.httpSearch.Get(%s) code(%d) error", d.searchURL, res.Code)
		err = ecode.Int(res.Code)
	}
	return
}

// SearchPGC search pgc(movie) data.
func (d *Dao) SearchPGC(c context.Context, mid int64, arg *model.SearchTypeArg, buvid, ua string) (res *model.SearchTypeRes, err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params = setSearchTypeParam(params, mid, model.SearchTypePGC, buvid, ip, arg)
	params.Set("page", strconv.Itoa(arg.Pn))
	for k, v := range model.SearchDefaultArg[model.SearchTypePGC] {
		params.Set(k, strconv.Itoa(v))
	}
	var req *http.Request
	if req, err = d.httpSearch.NewRequest(http.MethodGet, d.searchURL, ip, params); err != nil {
		return
	}
	req.Header.Set("browser-info", ua)
	res = new(model.SearchTypeRes)
	if err = d.httpSearch.Do(c, req, &res); err != nil {
		log.Error("SearchPGC d.httpSearch.Get(%s) error(%v)", d.searchURL, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("SearchPGC d.httpSearch.Get(%s) code(%d) error", d.searchURL, res.Code)
		err = ecode.Int(res.Code)
	}
	return
}

// SearchLive search live data.
func (d *Dao) SearchLive(c context.Context, mid int64, arg *model.SearchTypeArg, buvid, ua string) (res *model.SearchTypeRes, err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params = setSearchTypeParam(params, mid, model.SearchTypeLive, buvid, ip, arg)
	params.Set("page", strconv.Itoa(arg.Pn))
	params.Set("order", arg.Order)
	for k, v := range model.SearchDefaultArg[model.SearchTypeLive] {
		params.Set(k, strconv.Itoa(v))
	}
	var req *http.Request
	if req, err = d.httpSearch.NewRequest(http.MethodGet, d.searchURL, ip, params); err != nil {
		return
	}
	req.Header.Set("browser-info", ua)
	res = new(model.SearchTypeRes)
	if err = d.httpSearch.Do(c, req, &res); err != nil {
		log.Error("SearchVideo d.httpSearch.Get(%s) error(%v)", d.searchURL, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("SearchVideo d.httpSearch.Get(%s) code(%d) error", d.searchURL, res.Code)
		err = ecode.Int(res.Code)
	}
	return
}

// SearchLiveRoom search live data.
func (d *Dao) SearchLiveRoom(c context.Context, mid int64, arg *model.SearchTypeArg, buvid, ua string) (res *model.SearchTypeRes, err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params = setSearchTypeParam(params, mid, model.SearchTypeLiveRoom, buvid, ip, arg)
	params.Set("page", strconv.Itoa(arg.Pn))
	params.Set("order", arg.Order)
	for k, v := range model.SearchDefaultArg[model.SearchTypeLiveRoom] {
		params.Set(k, strconv.Itoa(v))
	}
	var req *http.Request
	if req, err = d.httpSearch.NewRequest(http.MethodGet, d.searchURL, ip, params); err != nil {
		return
	}
	req.Header.Set("browser-info", ua)
	res = new(model.SearchTypeRes)
	if err = d.httpSearch.Do(c, req, &res); err != nil {
		log.Error("SearchLiveRoom d.httpSearch.Get(%s) error(%v)", d.searchURL, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("SearchLiveRoom d.httpSearch.Get(%s) code(%d) error", d.searchURL, res.Code)
		err = ecode.Int(res.Code)
	}
	return
}

// SearchLiveUser search live user data.
func (d *Dao) SearchLiveUser(c context.Context, mid int64, arg *model.SearchTypeArg, buvid, ua string) (res *model.SearchTypeRes, err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params = setSearchTypeParam(params, mid, model.SearchTypeLiveUser, buvid, ip, arg)
	params.Set("page", strconv.Itoa(arg.Pn))
	params.Set("order", arg.Order)
	for k, v := range model.SearchDefaultArg[model.SearchTypeLiveUser] {
		params.Set(k, strconv.Itoa(v))
	}
	var req *http.Request
	if req, err = d.httpSearch.NewRequest(http.MethodGet, d.searchURL, ip, params); err != nil {
		return
	}
	req.Header.Set("browser-info", ua)
	res = new(model.SearchTypeRes)
	if err = d.httpSearch.Do(c, req, &res); err != nil {
		log.Error("SearchVideo d.httpSearch.Get(%s) error(%v)", d.searchURL, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("SearchVideo d.httpSearch.Get(%s) code(%d) error", d.searchURL, res.Code)
		err = ecode.Int(res.Code)
	}
	return
}

// SearchArticle search article.
func (d *Dao) SearchArticle(c context.Context, mid int64, arg *model.SearchTypeArg, buvid, ua string) (res *model.SearchTypeRes, err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params = setSearchTypeParam(params, mid, model.SearchTypeArticle, buvid, ip, arg)
	params.Set("category_id", strconv.FormatInt(arg.CategoryID, 10))
	params.Set("page", strconv.Itoa(arg.Pn))
	params.Set("order", arg.Order)
	for k, v := range model.SearchDefaultArg[model.SearchTypeArticle] {
		params.Set(k, strconv.Itoa(v))
	}
	var req *http.Request
	if req, err = d.httpSearch.NewRequest(http.MethodGet, d.searchURL, ip, params); err != nil {
		return
	}
	req.Header.Set("browser-info", ua)
	res = new(model.SearchTypeRes)
	if err = d.httpSearch.Do(c, req, &res); err != nil {
		log.Error("SearchArticle d.httpSearch.Get(%s) error(%v)", d.searchURL, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("SearchArticle d.httpSearch.Get(%s) code(%d) error", d.searchURL, res.Code)
		err = ecode.Int(res.Code)
	}
	return
}

// SearchSpecial search special data.
func (d *Dao) SearchSpecial(c context.Context, mid int64, arg *model.SearchTypeArg, buvid, ua string) (res *model.SearchTypeRes, err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params = setSearchTypeParam(params, mid, model.SearchTypeSpecial, buvid, ip, arg)
	params.Set("page", strconv.Itoa(arg.Pn))
	params.Set("vp_num", strconv.Itoa(arg.VpNum))
	for k, v := range model.SearchDefaultArg[model.SearchTypeSpecial] {
		params.Set(k, strconv.Itoa(v))
	}
	var req *http.Request
	if req, err = d.httpSearch.NewRequest(http.MethodGet, d.searchURL, ip, params); err != nil {
		return
	}
	req.Header.Set("browser-info", ua)
	res = new(model.SearchTypeRes)
	if err = d.httpSearch.Do(c, req, &res); err != nil {
		log.Error("SearchSpecial d.httpSearch.Get(%s) error(%v)", d.searchURL, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("SearchSpecial d.httpSearch.Get(%s) code(%d) error", d.searchURL, res.Code)
		err = ecode.Int(res.Code)
	}
	return
}

// SearchTopic search topic data.
func (d *Dao) SearchTopic(c context.Context, mid int64, arg *model.SearchTypeArg, buvid, ua string) (res *model.SearchTypeRes, err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params = setSearchTypeParam(params, mid, model.SearchTypeTopic, buvid, ip, arg)
	params.Set("page", strconv.Itoa(arg.Pn))
	for k, v := range model.SearchDefaultArg[model.SearchTypeTopic] {
		params.Set(k, strconv.Itoa(v))
	}
	if arg.Highlight > 0 {
		params.Set("highlight", strconv.Itoa(arg.Highlight))
	}
	var req *http.Request
	if req, err = d.httpSearch.NewRequest(http.MethodGet, d.searchURL, ip, params); err != nil {
		return
	}
	req.Header.Set("browser-info", ua)
	res = new(model.SearchTypeRes)
	if err = d.httpSearch.Do(c, req, &res); err != nil {
		log.Error("SearchVideo d.httpSearch.Get(%s) error(%v)", d.searchURL, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("SearchVideo d.httpSearch.Get(%s) code(%d) error", d.searchURL, res.Code)
		err = ecode.Int(res.Code)
	}
	return
}

// SearchUser search user data.
func (d *Dao) SearchUser(c context.Context, mid int64, arg *model.SearchTypeArg, buvid, ua string) (res *model.SearchTypeRes, err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params = setSearchTypeParam(params, mid, model.SearchTypeUser, buvid, ip, arg)
	params.Set("page", strconv.Itoa(arg.Pn))
	params.Set("user_type", strconv.Itoa(arg.UserType))
	params.Set("bili_user_vl", strconv.Itoa(arg.BiliUserVl))
	params.Set("order_sort", strconv.Itoa(arg.OrderSort))
	params.Set("order", arg.Order)
	for k, v := range model.SearchDefaultArg[model.SearchTypeUser] {
		params.Set(k, strconv.Itoa(v))
	}
	var req *http.Request
	if req, err = d.httpSearch.NewRequest(http.MethodGet, d.searchURL, ip, params); err != nil {
		return
	}
	req.Header.Set("browser-info", ua)
	res = new(model.SearchTypeRes)
	if err = d.httpSearch.Do(c, req, &res); err != nil {
		log.Error("SearchUser d.httpSearch.Get(%s) error(%v)", d.searchURL, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("SearchUser d.httpSearch.Get(%s) code(%d) error", d.searchURL, res.Code)
		err = ecode.Int(res.Code)
	}
	return
}

// SearchPhoto search photo data.
func (d *Dao) SearchPhoto(c context.Context, mid int64, arg *model.SearchTypeArg, buvid, ua string) (res *model.SearchTypeRes, err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params = setSearchTypeParam(params, mid, model.SearchTypePhoto, buvid, ip, arg)
	params.Set("category_id", strconv.FormatInt(arg.CategoryID, 10))
	params.Set("page", strconv.Itoa(arg.Pn))
	params.Set("order", arg.Order)
	for k, v := range model.SearchDefaultArg[model.SearchTypePhoto] {
		params.Set(k, strconv.Itoa(v))
	}
	if arg.Highlight > 0 {
		params.Set("highlight", strconv.Itoa(arg.Highlight))
	}
	var req *http.Request
	if req, err = d.httpSearch.NewRequest(http.MethodGet, d.searchURL, ip, params); err != nil {
		return
	}
	req.Header.Set("browser-info", ua)
	res = new(model.SearchTypeRes)
	if err = d.httpSearch.Do(c, req, &res); err != nil {
		log.Error("SearchPhoto d.httpSearch.Get(%s) error(%v)", d.searchURL, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("SearchPhoto d.httpSearch.Get(%s) code(%d) error", d.searchURL, res.Code)
		err = ecode.Int(res.Code)
	}
	return
}

// SearchRec search recommend data.
func (d *Dao) SearchRec(c context.Context, mid int64, pn, ps int, keyword, fromSource, buvid, ua string) (res *model.SearchRec, err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params = setSearchParam(params, mid, "", keyword, _searchPlatform, fromSource, buvid, ip)
	params.Set("page", strconv.Itoa(pn))
	params.Set("pagesize", strconv.Itoa(ps))
	var req *http.Request
	if req, err = d.httpSearch.NewRequest(http.MethodGet, d.searchRecURL, ip, params); err != nil {
		return
	}
	req.Header.Set("browser-info", ua)
	res = new(model.SearchRec)
	if err = d.httpSearch.Do(c, req, &res); err != nil {
		log.Error("Search d.httpSearch.Get(%s) error(%v)", d.searchRecURL, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("Search d.httpSearch.Do(%s) code error(%d)", d.searchRecURL, res.Code)
		err = ecode.Int(res.Code)
	}
	return
}

// SearchDefault get search default word.
func (d *Dao) SearchDefault(c context.Context, mid int64, fromSource, buvid, ua string) (data *model.SearchDefault, err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params.Set("main_ver", _searchVer)
	params.Set("platform", _searchPlatform)
	params.Set("clientip", ip)
	params.Set("userid", strconv.FormatInt(mid, 10))
	params.Set("search_type", "default")
	params.Set("from_source", fromSource)
	params.Set("buvid", buvid)
	var req *http.Request
	if req, err = d.httpSearch.NewRequest(http.MethodGet, d.searchDefaultURL, ip, params); err != nil {
		return
	}
	req.Header.Set("browser-info", ua)
	var res struct {
		Code   int    `json:"code"`
		SeID   string `json:"seid"`
		Tips   string `json:"recommend_tips"`
		Result []struct {
			ID       int64  `json:"id"`
			Name     string `json:"name"`
			ShowName string `json:"show_name"`
			Type     string `json:"type"`
		} `json:"result"`
	}
	if err = d.httpSearch.Do(c, req, &res); err != nil {
		log.Error("Search d.httpSearch.Get(%s) error(%v)", d.searchDefaultURL, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("Search d.httpSearch.Do(%s) code error(%d)", d.searchDefaultURL, res.Code)
		err = ecode.Int(res.Code)
	}
	if len(res.Result) == 0 {
		err = ecode.NothingFound
		return
	}
	data = &model.SearchDefault{}
	for _, v := range res.Result {
		data.Trackid = res.SeID
		data.ID = v.ID
		data.ShowName = v.ShowName
		data.Name = v.Name
	}
	return
}

// UpRecommend .
func (d *Dao) UpRecommend(c context.Context, mid int64, arg *model.SearchUpRecArg, buvid string) (rs []*model.SearchUpRecRes, trackID string, err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params.Set("userid", strconv.FormatInt(mid, 10))
	params.Set("service_area", arg.ServiceArea)
	params.Set("rec_type", _searchUpRecType)
	params.Set("platform", arg.Platform)
	params.Set("clientip", ip)
	params.Set("pagesize", strconv.Itoa(arg.Ps))
	params.Set("buvid", buvid)
	if arg.MobiApp != "" {
		params.Set("mobi_app", arg.MobiApp)
	}
	if arg.Device != "" {
		params.Set("device", arg.Device)
	}
	if arg.Build != 0 {
		params.Set("build", strconv.FormatInt(arg.Build, 10))
	}
	if arg.ContextID != 0 {
		params.Set("context_id", strconv.FormatInt(arg.ContextID, 10))
	}
	if len(arg.MainTids) > 0 {
		params.Set("main_tids", xstr.JoinInts(arg.MainTids))
	}
	if len(arg.SubTids) > 0 {
		params.Set("sub_tids", xstr.JoinInts(arg.SubTids))
	}
	var res struct {
		Code    int                     `json:"code"`
		Trackid string                  `json:"trackid"`
		Data    []*model.SearchUpRecRes `json:"data"`
	}
	if err = d.httpSearch.Get(c, d.searchUpRecURL, ip, params, &res); err != nil {
		log.Error("UpRecommend d.httpSearch.Get(%s) error(%v)", d.searchUpRecURL, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("UpRecommend d.httpSearch.Do(%s) code error(%d)", d.searchUpRecURL, res.Code)
		err = ecode.Int(res.Code)
		return
	}
	rs = res.Data
	trackID = res.Trackid
	return
}

// SearchEgg search egg.
func (d *Dao) SearchEgg(c context.Context) (data []*model.SearchEgg, err error) {
	var (
		ip  = metadata.String(c, metadata.RemoteIP)
		res struct {
			Code int                `json:"code"`
			Data []*model.SearchEgg `json:"data"`
		}
	)
	if err = d.httpSearch.Get(c, d.searchEggURL, ip, url.Values{}, &res); err != nil {
		log.Error("SearchEgg d.httpSearch.Get(%s) error(%v)", d.searchEggURL, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("SearchEgg d.httpSearch.Do(%s) code error(%d)", d.searchEggURL, res.Code)
		err = ecode.Int(res.Code)
		return
	}
	data = res.Data
	return
}

func setSearchParam(param url.Values, mid int64, searchType, keyword, platform, fromSource, buvid, ip string) url.Values {
	param.Set("main_ver", _searchVer)
	if searchType != "" {
		param.Set("search_type", searchType)
	}
	param.Set("platform", platform)
	param.Set("keyword", keyword)
	param.Set("from_source", fromSource)
	param.Set("userid", strconv.FormatInt(mid, 10))
	param.Set("buvid", buvid)
	param.Set("clientip", ip)
	return param
}

func setSearchTypeParam(param url.Values, mid int64, searchType, buvid, ip string, arg *model.SearchTypeArg) url.Values {
	param.Set("main_ver", _searchVer)
	if searchType != "" {
		param.Set("search_type", searchType)
	}
	param.Set("platform", arg.Platform)
	param.Set("keyword", arg.Keyword)
	param.Set("from_source", arg.FromSource)
	param.Set("userid", strconv.FormatInt(mid, 10))
	param.Set("single_column", strconv.Itoa(arg.SingleColumn))
	param.Set("buvid", buvid)
	param.Set("clientip", ip)
	return param
}
