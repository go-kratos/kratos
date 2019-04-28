package elec

import (
	"context"
	"net/url"
	"strconv"

	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_showURL     = "/internal/member/show"
	_arcOpenURI  = "/internal/archice/partin"
	_arcCloseURI = "/internal/archice/exit"
)

// ArcShow return archive elec show, contains rank.
func (d *Dao) ArcShow(c context.Context, mid, aid int64, ip string) (show bool, err error) {
	params := url.Values{}
	params.Set("upmid", strconv.FormatInt(mid, 10))
	params.Set("aid", strconv.FormatInt(aid, 10))
	params.Set("nolist", "0")
	var res struct {
		Code int `json:"code"`
		Data struct {
			Show bool `json:"show"`
		} `json:"data"`
	}
	if err = d.client.Get(c, d.showURI, ip, params, &res); err != nil {
		log.Error("elec url(%s) error(%v)", d.showURI+"?"+params.Encode(), err)
		err = ecode.CreativeElecErr
		return
	}
	log.Info("ArcShow d.showURI url(%s)|res(%+v)", d.showURI+"?"+params.Encode(), res)
	if res.Code != 0 {
		log.Error("d.client.Get(%s) code(%d) error(%v)", d.showURI+"?"+params.Encode(), res.Code, err)
		err = ecode.Int(res.Code)
		return
	}
	show = res.Data.Show
	return
}

// ArcUpdate arc open or close elec.
func (d *Dao) ArcUpdate(c context.Context, mid, aid int64, openElec int8, ip string) (err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("aid", strconv.FormatInt(aid, 10))
	var url string
	if openElec == 1 {
		url = d.arcOpenURL
	} else if openElec == 0 {
		url = d.arcCloseURL
	}
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Post(c, url, ip, params, &res); err != nil {
		log.Error("d.client.Do uri(%s) aid(%d) mid(%d) orderID(%d) code(%d) error(%v)", url+"?"+params.Encode(), mid, aid, openElec, res.Code, err)
		err = ecode.CreativeElecErr
		return
	}
	log.Info("dealElec ArcUpdate url(%s)", url+"?"+params.Encode())
	if res.Code != 0 {
		log.Error("arc elec update state  url(%s) res(%v); mid(%d), aid(%d), ip(%s), code(%d), error(%v)", url, res, mid, aid, ip, res.Code, err)
		err = ecode.Int(res.Code)
		return
	}
	return
}
