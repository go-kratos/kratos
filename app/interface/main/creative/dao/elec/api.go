package elec

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/dao/tool"
	"go-common/app/interface/main/creative/model/elec"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_userInfoURI     = "/internal/member/info"
	_userJoinURI     = "/internal/member/elec/partin"
	_userExitURI     = "/internal/member/elec/exit"
	_arcOpenURI      = "/internal/archice/partin"
	_arcCloseURI     = "/internal/archice/exit"
	_notifyURI       = "/internal/notify/info"
	_getStatusURI    = "/api/user/queryset/v2"
	_setStatusURI    = "/api/user/modifyset/v2"
	_recentRankURI   = "/api/query.recent.do"
	_currentRankURI  = "/api/query.rank.do"
	_totalRankURI    = "/api/query.total.rank.do"
	_dailyBillURI    = "/api/query.daily.bill.do"
	_balanceURI      = "/api/query.wallet.balance.do"
	_recentElecURI   = "/api/recent/elec"
	_remarkListURI   = "/api/elec/remark/list"
	_remarkDetailURI = "/api/elec/remake/detail"
	_remarkURI       = "/api/remake/reply"
)

// UserInfo get user elec info.
func (d *Dao) UserInfo(c context.Context, mid int64, ip string) (st *elec.UserInfo, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int `json:"code"`
		Data *elec.UserInfo
	}
	if err = d.client.Get(c, d.userInfoURL, ip, params, &res); err != nil {
		log.Error("elec url(%s) response(%v) error(%v)", d.userInfoURL+"?"+params.Encode(), res, err)
		err = ecode.CreativeElecErr
		return
	}
	if res.Code != 0 {
		log.Error("elec url(%s) res(%v)", d.userInfoURL, res)
		err = ecode.CreativeElecErr
		return
	}
	st = res.Data
	return
}

// UserUpdate join or exit elec.
func (d *Dao) UserUpdate(c context.Context, mid int64, st int8, ip string) (u *elec.UserInfo, err error) {
	params := url.Values{}
	params.Set("appkey", conf.Conf.App.Key)
	params.Set("appsecret", conf.Conf.App.Secret)
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	// url
	var (
		query, _ = tool.Sign(params)
		url      string
	)
	if st == 1 {
		url = d.userJoinURL + "?" + query
	} else if st == 2 {
		url = d.userExitURL + "?" + query
	}
	// new requests
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Error("http.NewRequest(%s) error(%v); mid(%d), ip(%s)", url, err, mid, ip)
		err = ecode.CreativeElecErr
		return
	}
	req.Header.Set("X-BACKEND-BILI-REAL-IP", ip)
	var res struct {
		Code int `json:"code"`
		Data *elec.UserInfo
	}
	if err = d.client.Do(c, req, &res); err != nil {
		log.Error("d.client.Do(%s) error(%v); mid(%d), ip(%s)", url, err, mid, ip)
		err = ecode.CreativeElecErr
		return
	}
	if res.Code != 0 {
		log.Error("user elec update state url(%s) res(%v); mid(%d), ip(%s), code(%d)", url, res, mid, ip, res.Code)
		err = ecode.CreativeElecErr
	}
	u = res.Data
	return
}

// ArcUpdate arc open or close elec.
func (d *Dao) ArcUpdate(c context.Context, mid, aid int64, st int8, ip string) (err error) {
	params := url.Values{}
	params.Set("appkey", conf.Conf.App.Key)
	params.Set("appsecret", conf.Conf.App.Secret)
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("aid", strconv.FormatInt(aid, 10))
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	// url
	var (
		query, _ = tool.Sign(params)
		url      string
	)
	if st == 1 {
		url = d.arcOpenURL + "?" + query
	} else if st == 2 {
		url = d.arcCloseURL + "?" + query
	}
	// new request
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Error("http.NewRequest(%s) error(%v); mid(%d), aid(%d), ip(%s)", url, err, mid, aid, ip)
		err = ecode.CreativeElecErr
		return
	}
	req.Header.Set("X-BACKEND-BILI-REAL-IP", ip)
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Do(c, req, &res); err != nil {
		log.Error("d.client.Do(%s) error(%v); mid(%d), aid(%d), ip(%s)", url, err, mid, aid, ip)
		err = ecode.CreativeElecErr
		return
	}
	if res.Code != 0 {
		log.Error("arc elec update state  url(%s) res(%v); mid(%d), aid(%d), ip(%s)", url, res, mid, aid, ip)
		err = ecode.CreativeElecErr
	}
	return
}

// Notify get up-to-date notice.
func (d *Dao) Notify(c context.Context, ip string) (nt *elec.Notify, err error) {
	params := url.Values{}
	var res struct {
		Code int `json:"code"`
		Data *elec.Notify
	}
	if err = d.client.Get(c, d.notifyURL, ip, params, &res); err != nil {
		log.Error("elec url(%s) response(%v) error(%v)", d.notifyURL+"?"+params.Encode(), res, err)
		err = ecode.CreativeElecErr
		return
	}
	if res.Code != 0 {
		log.Error("elec notify url(%s) res(%v)", d.notifyURL, res)
		err = ecode.CreativeElecErr
		return
	}
	nt = res.Data
	return
}

// Status get elec setting status.
func (d *Dao) Status(c context.Context, mid int64, ip string) (st *elec.Status, err error) {
	params := url.Values{}
	params.Set("act", "appkey")
	params.Set("mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int `json:"code"`
		Data struct {
			Info *elec.Status `json:"info"`
		} `json:"data"`
	}
	if err = d.client.Get(c, d.getStatusURL, ip, params, &res); err != nil {
		log.Error("elec url(%s) response(%v) error(%v)", d.getStatusURL+"?"+params.Encode(), res, err)
		err = ecode.CreativeElecErr
		return
	}
	if res.Code != 0 {
		log.Error("elec url(%s) res(%v)", d.getStatusURL, res)
		err = ecode.CreativeElecErr
		return
	}
	st = res.Data.Info
	return
}

// UpStatus update elec setting status.
func (d *Dao) UpStatus(c context.Context, mid int64, spday int, ip string) (err error) {
	params := url.Values{}
	params.Set("type", "json")
	params.Set("act", "appkey")
	params.Set("appkey", conf.Conf.App.Key)
	params.Set("appsecret", conf.Conf.App.Secret)
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("display_specialday", strconv.Itoa(spday))
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	var (
		query, _ = tool.Sign(params)
		url      string
	)
	url = d.setStatusURL + "?" + query
	// new request
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Error("http.NewRequest(%s) error(%v); mid(%d), ip(%s)", url, err, mid, ip)
		err = ecode.CreativeElecErr
		return
	}
	req.Header.Set("X-BACKEND-BILI-REAL-IP", ip)
	var res struct {
		Code int `json:"code"`
		Data struct {
			Ret int `json:"ret"`
		} `json:"data"`
	}
	if err = d.client.Do(c, req, &res); err != nil {
		log.Error("d.client.Do(%s) error(%v); mid(%d), ip(%s)", url, err, mid, ip)
		err = ecode.CreativeElecErr
		return
	}
	if res.Code != 0 {
		log.Error("user elec update setting url(%s) res(%v); mid(%d), ip(%s)", url, res, mid, ip)
		err = ecode.CreativeElecErr
	}
	return
}

// RecentRank get recent rank.
func (d *Dao) RecentRank(c context.Context, mid, size int64, ip string) (rec []*elec.Rank, err error) {
	params := url.Values{}
	params.Set("type", "json")
	params.Set("act", "appkey")
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("size", strconv.FormatInt(size, 10))
	var res struct {
		Code int `json:"code"`
		Data struct {
			List []*elec.Rank `json:"list"`
		} `json:"data"`
	}
	if err = d.client.Get(c, d.recentRankURL, ip, params, &res); err != nil {
		log.Error("elec url(%s) response(%v) error(%v)", d.recentRankURL+"?"+params.Encode(), res, err)
		err = ecode.CreativeElecErr
		return
	}
	if res.Code != 0 {
		log.Error("elec url(%s) res(%v)", d.recentRankURL, res)
		err = ecode.CreativeElecErr
		return
	}
	rec = res.Data.List
	return
}

// CurrentRank get current rank.
func (d *Dao) CurrentRank(c context.Context, mid int64, ip string) (cur []*elec.Rank, err error) {
	params := url.Values{}
	params.Set("type", "json")
	params.Set("act", "appkey")
	params.Set("mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int `json:"code"`
		Data struct {
			List []*elec.Rank `json:"list"`
		} `json:"data"`
	}
	if err = d.client.Get(c, d.currentRankURL, ip, params, &res); err != nil {
		log.Error("elec url(%s) response(%v) error(%v)", d.currentRankURL+"?"+params.Encode(), res, err)
		err = ecode.CreativeElecErr
		return
	}
	if res.Code != 0 {
		log.Error("elec url(%s) res(%v)", d.currentRankURL, res)
		err = ecode.CreativeElecErr
		return
	}
	cur = res.Data.List
	return
}

// TotalRank get total rank.
func (d *Dao) TotalRank(c context.Context, mid int64, ip string) (tol []*elec.Rank, err error) {
	params := url.Values{}
	params.Set("type", "json")
	params.Set("act", "appkey")
	params.Set("mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int `json:"code"`
		Data struct {
			List []*elec.Rank `json:"list"`
		} `json:"data"`
	}
	if err = d.client.Get(c, d.totalRankURL, ip, params, &res); err != nil {
		log.Error("elec url(%s) response(%v) error(%v)", d.totalRankURL+"?"+params.Encode(), res, err)
		err = ecode.CreativeElecErr
		return
	}
	if res.Code != 0 {
		log.Error("elec url(%s) res(%v)", d.totalRankURL, res)
		err = ecode.CreativeElecErr
		return
	}
	tol = res.Data.List
	return
}

// DailyBill daily settlement.
func (d *Dao) DailyBill(c context.Context, mid int64, pn, ps int, begin, end, ip string) (bl *elec.BillList, err error) {
	params := url.Values{}
	params.Set("type", "json")
	params.Set("act", "appkey")
	params.Set("appkey", conf.Conf.App.Key)
	params.Set("appsecret", conf.Conf.App.Secret)
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("page_no", strconv.Itoa(pn))
	params.Set("page_size", strconv.Itoa(ps))
	params.Set("begin_time", begin)
	params.Set("end_time", end)
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	// url
	var (
		query, _ = tool.Sign(params)
		url      string
	)
	url = d.dailyBillURL + "?" + query
	// new request
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Error("http.NewRequest(%s) error(%v); mid(%d), ip(%s)", url, err, mid, ip)
		err = ecode.CreativeElecErr
		return
	}
	req.Header.Set("X-BACKEND-BILI-REAL-IP", ip)
	var res struct {
		Code int            `json:"code"`
		Data *elec.BillList `json:"data"`
	}
	if err = d.client.Do(c, req, &res); err != nil {
		log.Error("d.client.Do(%s) error(%v); mid(%d), ip(%s)", url, err, mid, ip)
		err = ecode.CreativeElecErr
		return
	}
	if res.Code != 0 {
		log.Error("user elec daily bill url(%s) res(%v); mid(%d), ip(%s)", url, res, mid, ip)
		err = ecode.CreativeElecErr
	}
	bl = res.Data
	return
}

// Balance get battery balance.
func (d *Dao) Balance(c context.Context, mid int64, ip string) (bal *elec.Balance, err error) {
	params := url.Values{}
	params.Set("type", "json")
	params.Set("act", "appkey")
	params.Set("appkey", conf.Conf.App.Key)
	params.Set("appsecret", conf.Conf.App.Secret)
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	var (
		query, _ = tool.Sign(params)
		url      string
	)
	url = d.balanceURL + "?" + query
	// new request
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Error("http.NewRequest(%s) error(%v); mid(%d), ip(%s)", url, err, mid, ip)
		err = ecode.CreativeElecErr
		return
	}
	req.Header.Set("X-BACKEND-BILI-REAL-IP", ip)
	var res struct {
		Code int           `json:"code"`
		Data *elec.Balance `json:"data"`
	}
	if err = d.client.Do(c, req, &res); err != nil {
		log.Error("d.client.Do(%s) error(%v); mid(%d), ip(%s)", url, err, mid, ip)
		err = ecode.CreativeElecErr
		return
	}
	if res.Code != 0 {
		log.Error("user elec balance url(%s) res(%v); mid(%d), ip(%s)", url, res, mid, ip)
		err = ecode.CreativeElecErr
	}
	bal = res.Data
	return
}

// RecentElec get aid  & elec_num.
func (d *Dao) RecentElec(c context.Context, mid int64, pn, ps int, ip string) (rec *elec.RecentElecList, err error) {
	params := url.Values{}
	params.Set("type", "json")
	params.Set("act", "appkey")
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("pn", strconv.Itoa(pn))
	params.Set("ps", strconv.Itoa(ps))
	var res struct {
		Code int                  `json:"code"`
		Data *elec.RecentElecList `json:"data"`
	}
	if err = d.client.Get(c, d.recentElecURL, ip, params, &res); err != nil {
		log.Error("elec url(%s) response(%v) error(%v)", d.recentElecURL+"?"+params.Encode(), res, err)
		err = ecode.CreativeElecErr
		return
	}
	if res.Code != 0 {
		log.Error("elec url(%s) res(%v)", d.recentElecURL, res)
		err = ecode.CreativeElecErr
		return
	}
	rec = res.Data
	return
}

// RemarkList get remark list.
func (d *Dao) RemarkList(c context.Context, mid int64, pn, ps int, begin, end, ip string) (rec *elec.RemarkList, err error) {
	params := url.Values{}
	params.Set("type", "json")
	params.Set("act", "appkey")
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("pn", strconv.Itoa(pn))
	params.Set("ps", strconv.Itoa(ps))
	params.Set("start_time", begin)
	params.Set("end_time", end)
	var res struct {
		Code int              `json:"code"`
		Data *elec.RemarkList `json:"data"`
	}
	if err = d.client.Get(c, d.remarkListURL, ip, params, &res); err != nil {
		log.Error("elec url(%s) response(%v) error(%v)", d.remarkListURL+"?"+params.Encode(), res, err)
		err = ecode.CreativeElecErr
		return
	}
	if res.Code != 0 {
		log.Error("elec url(%s) res(%v)", d.remarkListURL, res)
		err = ecode.CreativeElecErr
		return
	}
	rec = res.Data
	return
}

// RemarkDetail get remark detail.
func (d *Dao) RemarkDetail(c context.Context, mid, id int64, ip string) (re *elec.Remark, err error) {
	params := url.Values{}
	params.Set("type", "json")
	params.Set("act", "appkey")
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("id", strconv.FormatInt(id, 10))
	var res struct {
		Code int `json:"code"`
		Data struct {
			Info *elec.Remark `json:"info"`
		} `json:"data"`
	}
	if err = d.client.Get(c, d.remarkDetailURL, ip, params, &res); err != nil {
		log.Error("elec url(%s) response(%v) error(%v)", d.remarkDetailURL+"?"+params.Encode(), res, err)
		err = ecode.CreativeElecErr
		return
	}
	if res.Code != 0 {
		log.Error("elec url(%s) res(%v)", d.remarkDetailURL, res)
		err = ecode.CreativeElecErr
		return
	}
	re = res.Data.Info
	return
}

// Remark  reply a msg.
func (d *Dao) Remark(c context.Context, mid, id int64, msg, ak, ck, ip string) (status int, err error) {
	params := url.Values{}
	params.Set("type", "json")
	params.Set("act", "appkey")
	params.Set("access_key", ak)
	params.Set("appkey", conf.Conf.App.Key)
	params.Set("appsecret", conf.Conf.App.Secret)
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("id", strconv.FormatInt(id, 10))
	params.Set("message", msg)
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	var (
		query, _ = tool.Sign(params)
		url      string
	)
	url = d.remarkURL + "?" + query
	// new request
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Error("http.NewRequest(%s) error(%v); mid(%d), ip(%s)", url, err, mid, ip)
		err = ecode.CreativeElecErr
		return
	}
	req.Header.Set("X-BACKEND-BILI-REAL-IP", ip)
	req.Header.Set("Cookie", ck)
	var res struct {
		Code int `json:"code"`
		Data struct {
			Status int `json:"status"`
		} `json:"data"`
	}
	if err = d.client.Do(c, req, &res); err != nil {
		log.Error("d.client.Do(%s) error(%v); mid(%d), ip(%s)", url, err, mid, ip)
		err = ecode.CreativeElecErr
		return
	}
	if res.Code != 0 {
		log.Error("user elec daily bill url(%s) res(%v); mid(%d), ip(%s)", url, res, mid, ip)
		if res.Code == 61001 || res.Code == 61002 {
			err = ecode.Int(res.Code)
		} else {
			err = ecode.CreativeElecErr
		}
	}
	status = res.Data.Status
	return
}
