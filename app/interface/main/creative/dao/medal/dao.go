package medal

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/model/medal"
	"go-common/library/ecode"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
)

const (
	_checkStatus = "/fans_medal/v1/medal/if_can_open"
	_openMedal   = "/fans_medal/v1/medal/open"
	_getMedal    = "/fans_medal/v1/medal/get"
	_recentFans  = "/fans_medal/v1/fans_medal/get_recent_fans_list"
	_checkMedal  = "/fans_medal/v1/medal/check_open"
	_fansRank    = "/fans_medal/v1/medal/top_medal_fans_list"
	_rename      = "/fans_medal/v1/medal/rename"
)

// Dao  define
type Dao struct {
	c *conf.Config
	// http
	client *httpx.Client
	// uri
	checkStatusURI string
	openMedalURI   string
	getMedalURI    string
	recentFansURI  string
	checkMedalURI  string
	fansRankURI    string
	renameURI      string
}

// New init dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:              c,
		client:         httpx.NewClient(c.HTTPClient.Slow),
		checkStatusURI: c.Host.Live + _checkStatus,
		openMedalURI:   c.Host.Live + _openMedal,
		getMedalURI:    c.Host.Live + _getMedal,
		recentFansURI:  c.Host.Live + _recentFans,
		checkMedalURI:  c.Host.Live + _checkMedal,
		fansRankURI:    c.Host.Live + _fansRank,
		renameURI:      c.Host.Live + _rename,
	}
	return
}

// Status check if up can open medal
func (d *Dao) Status(c context.Context, mid int64) (st *medal.Status, err error) {
	var (
		req    *http.Request
		params = url.Values{}
	)
	params.Set("uid", strconv.FormatInt(mid, 10))
	params.Set("source", "2")
	if req, err = http.NewRequest("GET", d.checkStatusURI+"?"+params.Encode(), nil); err != nil {
		log.Error("Status url(%s) error(%v)", d.checkStatusURI+"?"+params.Encode(), err)
		err = ecode.CreativeFansMedalErr
		return
	}
	req.Header.Set("X-BiliLive-UID", strconv.FormatInt(mid, 10))
	var res struct {
		Code int           `json:"code"`
		Data *medal.Status `json:"data"`
	}
	if err = d.client.Do(c, req, &res); err != nil {
		log.Error("d.Status url(%s) res(%+v) err(%v)", d.checkStatusURI+"?"+params.Encode(), res, err)
		err = ecode.CreativeFansMedalErr
		return
	}
	if res.Code != 0 {
		log.Error("Status url(%s) res(%v)", d.checkStatusURI+"?"+params.Encode(), res)
		err = ecode.Int(res.Code)
		return
	}
	st = res.Data
	return
}

// OpenMedal open medal for up
func (d *Dao) OpenMedal(c context.Context, mid int64, name string) (err error) {
	var (
		req    *http.Request
		params = url.Values{}
	)
	params.Set("uid", strconv.FormatInt(mid, 10))
	params.Set("medal_name", name)
	params.Set("source", "2")
	if req, err = http.NewRequest("POST", d.openMedalURI, strings.NewReader(params.Encode())); err != nil {
		log.Error("OpenMedal url(%s) error(%v)", d.openMedalURI+"?"+params.Encode(), err)
		err = ecode.CreativeFansMedalErr
		return
	}
	req.Header.Set("X-BiliLive-UID", strconv.FormatInt(mid, 10))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	var res struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	if err = d.client.Do(c, req, &res); err != nil {
		log.Error("d.OpenMedal url(%s) res(%+v) err(%v)", d.openMedalURI+"?"+params.Encode(), res, err)
		err = ecode.CreativeFansMedalErr
		return
	}
	if res.Code != 0 {
		log.Error("OpenMedal url(%s) res(%v)", d.openMedalURI+"?"+params.Encode(), res)
		err = ecode.Int(res.Code)
		return
	}
	return
}

// CheckMedal check medal name for up
func (d *Dao) CheckMedal(c context.Context, mid int64, name string) (valid int, err error) {
	var (
		req    *http.Request
		params = url.Values{}
	)
	params.Set("uid", strconv.FormatInt(mid, 10))
	params.Set("medal_name", name)
	params.Set("source", "2")
	if req, err = http.NewRequest("POST", d.checkMedalURI, strings.NewReader(params.Encode())); err != nil {
		log.Error("checkMedalURI url(%s) error(%v)", d.checkMedalURI+"?"+params.Encode(), err)
		err = ecode.CreativeFansMedalErr
		return
	}
	req.Header.Set("X-BiliLive-UID", strconv.FormatInt(mid, 10))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	var res struct {
		Code int `json:"code"`
		Data struct {
			Enable bool `json:"enable"`
		} `json:"data"`
	}
	if err = d.client.Do(c, req, &res); err != nil {
		log.Error("d.checkMedalURI url(%s) res(%+v) err(%v)", d.checkMedalURI+"?"+params.Encode(), res, err)
		err = ecode.CreativeFansMedalErr
		return
	}
	if res.Code != 0 {
		log.Error("checkMedalURI url(%s) res(%v)", d.checkMedalURI+"?"+params.Encode(), res)
		err = ecode.Int(res.Code)
		return
	}
	if !res.Data.Enable {
		valid = 1
	}
	return
}

//Medal get medal
func (d *Dao) Medal(c context.Context, mid int64) (mdl *medal.Medal, err error) {
	var (
		req    *http.Request
		params = url.Values{}
	)
	params.Set("uid", strconv.FormatInt(mid, 10))
	params.Set("source", "2")
	if req, err = http.NewRequest("GET", d.getMedalURI+"?"+params.Encode(), nil); err != nil {
		log.Error("Medal url(%s) error(%v)", d.getMedalURI+"?"+params.Encode(), err)
		err = ecode.CreativeFansMedalErr
		return
	}
	req.Header.Set("X-BiliLive-UID", strconv.FormatInt(mid, 10))
	var res struct {
		Code int          `json:"code"`
		Data *medal.Medal `json:"data"`
	}
	if err = d.client.Do(c, req, &res); err != nil {
		log.Error("Medal url(%s) response(%+v) error(%v)", d.getMedalURI+"?"+params.Encode(), res, err)
		err = ecode.CreativeFansMedalErr
		return
	}
	if res.Code != 0 {
		log.Error("Medal url(%s) res(%v)", d.getMedalURI+"?"+params.Encode(), res)
		err = ecode.Int(res.Code)
		return
	}
	mdl = res.Data
	return
}

//RecentFans get a list of recent fans
func (d *Dao) RecentFans(c context.Context, mid int64) (fans []*medal.RecentFans, err error) {
	var (
		req    *http.Request
		params = url.Values{}
	)
	params.Set("uid", strconv.FormatInt(mid, 10))
	params.Set("source", "2")
	params.Set("size", "10")
	if req, err = http.NewRequest("GET", d.recentFansURI+"?"+params.Encode(), nil); err != nil {
		log.Error("RecentFans url(%s) error(%v)", d.recentFansURI+"?"+params.Encode(), err)
		err = ecode.CreativeFansMedalErr
		return
	}
	req.Header.Set("X-BiliLive-UID", strconv.FormatInt(mid, 10))
	var res struct {
		Code int                 `json:"code"`
		Data []*medal.RecentFans `json:"data"`
	}
	if err = d.client.Do(c, req, &res); err != nil {
		log.Error("RecentFans url(%s) response(%+v) error(%v)", d.recentFansURI+"?"+params.Encode(), res, err)
		err = ecode.CreativeFansMedalErr
		return
	}
	if res.Code != 0 {
		log.Error("RecentFans url(%s) res(%v)", d.recentFansURI+"?"+params.Encode(), res)
		err = ecode.Int(res.Code)
		return
	}
	fans = res.Data
	return
}

//Rank get fans rank list
func (d *Dao) Rank(c context.Context, mid int64) (rank []*medal.FansRank, err error) {
	var (
		req    *http.Request
		params = url.Values{}
	)
	params.Set("uid", strconv.FormatInt(mid, 10))
	if req, err = http.NewRequest("GET", d.fansRankURI+"?"+params.Encode(), nil); err != nil {
		log.Error("Rank url(%s) error(%v)", d.fansRankURI+"?"+params.Encode(), err)
		err = ecode.CreativeFansMedalErr
		return
	}
	req.Header.Set("X-BiliLive-UID", strconv.FormatInt(mid, 10))
	var res struct {
		Code int               `json:"code"`
		Data []*medal.FansRank `json:"data"`
	}
	if err = d.client.Do(c, req, &res); err != nil {
		log.Error("Rank url(%s) response(%+v) error(%v)", d.fansRankURI+"?"+params.Encode(), res, err)
		err = ecode.CreativeFansMedalErr
		return
	}
	if res.Code != 0 {
		log.Error("Rank url(%s) res(%v)", d.fansRankURI+"?"+params.Encode(), res)
		err = ecode.Int(res.Code)
		return
	}
	rank = res.Data
	return
}

//Rename rename medal name
func (d *Dao) Rename(c context.Context, mid int64, name, ak, ck string) (err error) {
	var (
		req    *http.Request
		params = url.Values{}
	)
	params.Set("uid", strconv.FormatInt(mid, 10))
	params.Set("medal_name", name)
	params.Set("source", "2")
	if ak != "" {
		params.Set("access_key", ak)
		params.Set("platform", "ios")
	} else {
		params.Set("platform", "pc")
	}
	if req, err = http.NewRequest("POST", d.renameURI, strings.NewReader(params.Encode())); err != nil {
		log.Error("Rename url(%s) error(%v)", d.renameURI+"?"+params.Encode(), err)
		err = ecode.CreativeFansMedalErr
		return
	}
	req.Header.Set("X-BiliLive-UID", strconv.FormatInt(mid, 10))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", ck)
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Do(c, req, &res); err != nil {
		log.Error("Rename url(%s) response(%+v) error(%v)", d.renameURI+"?"+params.Encode(), res, err)
		err = ecode.CreativeFansMedalErr
		return
	}
	if res.Code != 0 {
		log.Error("Rename url(%s) res(%v)", d.renameURI+"?"+params.Encode(), res)
		err = ecode.Int(res.Code)
		return
	}
	return
}
