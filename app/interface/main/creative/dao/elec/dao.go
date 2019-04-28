package elec

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/model/elec"
	"go-common/library/ecode"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

const (
	_userStatURL     = "/internal/member/info"
	_arcStatURL      = "/internal/member/show"
	_elecRelationURI = "/api/query.elec.relation.do"
)

// Dao  define
type Dao struct {
	c *conf.Config
	// http
	client *httpx.Client
	// uri
	userStatURI     string
	arcStatURI      string
	elecRelationURI string
	// user
	userInfoURL string
	userJoinURL string
	userExitURL string
	// arc
	arcOpenURL  string
	arcCloseURL string
	notifyURL   string
	// status
	getStatusURL string
	setStatusURL string
	// rank
	recentRankURL  string
	currentRankURL string
	totalRankURL   string
	// money
	dailyBillURL string
	balanceURL   string
	// recent elec for app
	recentElecURL string
	// elec remark
	remarkListURL   string
	remarkDetailURL string
	remarkURL       string
}

// New init dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:               c,
		client:          httpx.NewClient(c.HTTPClient.Slow),
		userStatURI:     c.Host.Elec + _userStatURL,
		arcStatURI:      c.Host.Elec + _arcStatURL,
		elecRelationURI: c.Host.Elec + _elecRelationURI,
		// api elec.
		userInfoURL:     c.Host.Elec + _userInfoURI,
		userJoinURL:     c.Host.Elec + _userJoinURI,
		userExitURL:     c.Host.Elec + _userExitURI,
		arcOpenURL:      c.Host.Elec + _arcOpenURI,
		arcCloseURL:     c.Host.Elec + _arcCloseURI,
		notifyURL:       c.Host.Elec + _notifyURI,
		getStatusURL:    c.Host.Elec + _getStatusURI,
		setStatusURL:    c.Host.Elec + _setStatusURI,
		recentRankURL:   c.Host.Elec + _recentRankURI,
		currentRankURL:  c.Host.Elec + _currentRankURI,
		totalRankURL:    c.Host.Elec + _totalRankURI,
		dailyBillURL:    c.Host.Elec + _dailyBillURI,
		balanceURL:      c.Host.Elec + _balanceURI,
		recentElecURL:   c.Host.Elec + _recentElecURI,
		remarkListURL:   c.Host.Elec + _remarkListURI,
		remarkDetailURL: c.Host.Elec + _remarkDetailURI,
		remarkURL:       c.Host.Elec + _remarkURI,
	}
	return
}

// UserState get user elec state.
func (d *Dao) UserState(c context.Context, mid int64, ip string) (data *elec.UserState, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int `json:"code"`
		Data struct {
			Mid    int    `json:"mid"`
			State  int    `json:"state"`
			Reason string `json:"reason"`
			Count  int    `json:"count"`
			CTime  string `json:"ctime"`
		} `json:"data"`
	}
	if err = d.client.Get(c, d.userInfoURL, ip, params, &res); err != nil {
		log.Error("elec url(%s) response(%v) error(%v)", d.userInfoURL+"?"+params.Encode(), res, err)
		err = ecode.CreativeElecErr
		return
	}
	log.Info("UserState d.userInfoURL url(%s), code(%d)", d.userInfoURL+"?"+params.Encode(), res.Code)
	if res.Code != 0 {
		log.Error("elec url(%s) res(%v)", d.userInfoURL, res)
		err = ecode.CreativeElecErr
		return
	}
	data = &elec.UserState{
		Mid:    strconv.Itoa(res.Data.Mid),
		State:  strconv.Itoa(res.Data.State),
		Count:  strconv.Itoa(res.Data.Count),
		CTime:  res.Data.CTime,
		Reason: res.Data.Reason,
	}
	return
}

// ArchiveState get arc elec state.
func (d *Dao) ArchiveState(c context.Context, aid, mid int64, ip string) (data *elec.ArcState, err error) {
	params := url.Values{}
	params.Set("upmid", strconv.FormatInt(mid, 10))
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("aid", strconv.FormatInt(aid, 10))
	params.Set("nolist", "1")
	var res struct {
		Code int            `json:"code"`
		Data *elec.ArcState `json:"data"`
	}
	if err = d.client.Get(c, d.arcStatURI, ip, params, &res); err != nil {
		log.Error("elec url(%s) response(%v) error(%v)", d.arcStatURI+"?"+params.Encode(), res, err)
		err = ecode.CreativeElecErr
		return
	}
	log.Info("ArchiveState d.arcStatURI url(%s), code(%d)", d.arcStatURI+"?"+params.Encode(), res.Code)
	if res.Code != 0 {
		log.Error("elec url(%s) res(%v)", d.arcStatURI, res)
		err = ecode.CreativeElecErr
		return
	}
	data = res.Data
	return
}

// ElecRelation check if multi user charged.
func (d *Dao) ElecRelation(c context.Context, mid int64, mids []int64, ip string) (chargeMap map[int64]int, err error) {
	params := url.Values{}
	params.Set("act", "appkey")
	params.Set("type", "json")
	params.Set("up_mid", strconv.Itoa(int(mid)))
	params.Set("mids", xstr.JoinInts(mids))
	var res struct {
		Code int               `json:"code"`
		Data *elec.EleRelation `json:"data"`
	}
	if err = d.client.Get(c, d.elecRelationURI, ip, params, &res); err != nil {
		log.Error("elec url(%s) response(%v) error(%v)", d.elecRelationURI+"?"+params.Encode(), res, err)
		err = ecode.CreativeElecErr
		return
	}
	if res.Code != 0 {
		log.Error("elec url(%s) res(%v)", d.elecRelationURI, res)
		err = ecode.CreativeElecErr
		return
	}
	chargeMap = map[int64]int{}
	for _, v := range res.Data.RetList {
		if !v.IsElec {
			chargeMap[v.Mid] = 0
		} else {
			chargeMap[v.Mid] = 1
		}
	}
	return
}
