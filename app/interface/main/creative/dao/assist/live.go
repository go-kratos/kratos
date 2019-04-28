package assist

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"go-common/app/interface/main/creative/model/assist"
	"go-common/library/ecode"
	"go-common/library/log"
	gtime "go-common/library/time"
)

const (
	// api
	_liveStatus      = "/user/v1/UserRoom/isAnchor"
	_liveAddAssist   = "/live_user/v1/RoomAdmin/add"
	_liveDelAssist   = "/live_user/v1/RoomAdmin/del"
	_liveRevocBanned = "/live_user/v1/RoomSilent/del"
	_liveAssists     = "/live_user/v1/RoomAdmin/get_by_anchor"
	_liveCheckAssist = "/live_user/v1/RoomAdmin/is_admin"
)

// LiveStatus check if user opened live room. 0:yes 500401: no
func (d *Dao) LiveStatus(c context.Context, mid int64, ip string) (ok int8, err error) {
	var req *http.Request
	params := url.Values{}
	params.Set("uid", strconv.FormatInt(mid, 10))
	if req, err = http.NewRequest("GET", d.liveStatusURL+"?"+params.Encode(), nil); err != nil {
		log.Error("LiveStatus url(%s) error(%v)", d.liveStatusURL+"?"+params.Encode(), err)
		err = ecode.CreativeLiveErr
		return
	}
	req.Header.Set("X-BiliLive-UID", strconv.FormatInt(mid, 10))
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Do(c, req, &res); err != nil {
		log.Error("LiveStatus url(%s) response(%+v) error(%v)", d.liveStatusURL+"?"+params.Encode(), res, err)
		err = ecode.CreativeLiveErr
		return
	}
	if res.Code == 500401 {
		return
	}
	if res.Code != 0 {
		log.Error("LiveStatus url(%s) res(%v)", d.liveStatusURL+"?"+params.Encode(), res)
		err = ecode.Int(res.Code)
		return
	}
	ok = 1
	return
}

// LiveAddAssist add assist in live room   120013: 直播间不存在
func (d *Dao) LiveAddAssist(c context.Context, mid, assistMid int64, cookie, ip string) (err error) {
	var req *http.Request
	params := url.Values{}
	params.Set("admin", strconv.FormatInt(assistMid, 10))
	params.Set("anchor_id", strconv.FormatInt(mid, 10))
	if req, err = http.NewRequest("POST", d.liveAddAssistURL+"?"+params.Encode(), nil); err != nil {
		log.Error("LiveAddAssist url(%s) error(%v)", d.liveAddAssistURL+"?"+params.Encode(), err)
		err = ecode.CreativeLiveErr
		return
	}
	req.Header.Set("X-BiliLive-UID", strconv.FormatInt(mid, 10))
	req.Header.Set("Cookie", cookie)
	var res struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err = d.client.Do(c, req, &res); err != nil {
		log.Error("LiveAddAssist url(%s) response(%+v) error(%v)", d.liveAddAssistURL+"?"+params.Encode(), res, err)
		err = ecode.CreativeLiveErr
		return
	}
	if res.Code == 120013 {
		log.Error("LiveAddAssist url(%s) mid(%d) not opened(%v)", d.liveAddAssistURL+"?"+params.Encode(), mid, 120013)
		err = ecode.CreativeLiveNotOpenErr
		return
	}
	if res.Code != 0 {
		log.Error("LiveAddAssist url(%s) res(%v)", d.liveAddAssistURL+"?"+params.Encode(), res)
		err = ecode.Int(res.Code)
		return
	}
	return
}

// LiveDelAssist del assist in live room
func (d *Dao) LiveDelAssist(c context.Context, mid, assistMid int64, cookie, ip string) (err error) {
	var req *http.Request
	params := url.Values{}
	params.Set("admin", strconv.FormatInt(assistMid, 10))
	params.Set("anchor_id", strconv.FormatInt(mid, 10))
	if req, err = http.NewRequest("POST", d.liveDelAssistURL+"?"+params.Encode(), nil); err != nil {
		log.Error("LiveDelAssist url(%s) error(%v)", d.liveDelAssistURL+"?"+params.Encode(), err)
		err = ecode.CreativeLiveErr
		return
	}
	req.Header.Set("X-BiliLive-UID", strconv.FormatInt(mid, 10))
	req.Header.Set("Cookie", cookie)
	var res struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err = d.client.Do(c, req, &res); err != nil {
		log.Error("LiveDelAssist url(%s) response(%+v) error(%v)", d.liveDelAssistURL+"?"+params.Encode(), res, err)
		err = ecode.CreativeLiveErr
		return
	}
	if res.Code == 120013 {
		log.Error("LiveDelAssist url(%s) mid(%d) not opened(%v)", d.liveDelAssistURL+"?"+params.Encode(), mid, 120013)
		err = ecode.CreativeLiveNotOpenErr
		return
	}
	if res.Code != 0 {
		log.Error("LiveDelAssist url(%s) res(%v)", d.liveDelAssistURL+"?"+params.Encode(), res)
		err = ecode.Int(res.Code)
		return
	}
	return
}

// LiveBannedRevoc recover assist banned state
func (d *Dao) LiveBannedRevoc(c context.Context, mid int64, banID, cookie, ip string) (err error) {
	var req *http.Request
	params := url.Values{}
	params.Set("uid", strconv.FormatInt(mid, 10))
	params.Set("id", banID)
	if req, err = http.NewRequest("POST", d.liveRevocBannedURL+"?"+params.Encode(), nil); err != nil {
		log.Error("liveRevocBanned url(%s) error(%v)", d.liveRevocBannedURL+"?"+params.Encode(), err)
		err = ecode.CreativeLiveErr
		return
	}
	req.Header.Set("X-BiliLive-UID", strconv.FormatInt(mid, 10))
	req.Header.Set("Cookie", cookie)
	var res struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err = d.client.Do(c, req, &res); err != nil {
		log.Error("liveRevocBanned url(%s) response(%+v) error(%v)", d.liveRevocBannedURL+"?"+params.Encode(), res, err)
		err = ecode.CreativeLiveErr
		return
	}
	if res.Code != 0 {
		log.Error("liveRevocBanned url(%s) res(%v)", d.liveRevocBannedURL+"?"+params.Encode(), res)
		err = ecode.Int(res.Code)
		return
	}
	return
}

// LiveAssists get assists from live
func (d *Dao) LiveAssists(c context.Context, mid int64, ip string) (assists []*assist.LiveAssist, err error) {
	var req *http.Request
	params := url.Values{}
	params.Set("anchor_id", strconv.FormatInt(mid, 10))
	if req, err = http.NewRequest("GET", d.liveAssistsURL+"?"+params.Encode(), nil); err != nil {
		log.Error("LiveAssists url(%s) error(%v)", d.liveAssistsURL+"?"+params.Encode(), err)
		err = ecode.CreativeLiveErr
		return
	}
	req.Header.Set("X-BiliLive-UID", strconv.FormatInt(mid, 10))
	var res struct {
		Code int                  `json:"code"`
		Data []*assist.LiveAssist `json:"data"`
	}
	if err = d.client.Do(c, req, &res); err != nil {
		log.Error("LiveAssists url(%s) response(%+v) error(%v)", d.liveAssistsURL+"?"+params.Encode(), res, err)
		err = ecode.CreativeLiveErr
		return
	}
	if res.Code == 120013 {
		err = nil
		return
	}
	if res.Code != 0 {
		log.Error("LiveAssists url(%s) res(%v)", d.liveAssistsURL+"?"+params.Encode(), res)
		err = ecode.Int(res.Code)
		return
	}
	assists = res.Data
	// set ctime
	for _, ass := range assists {
		intTime, err := time.Parse("2006-01-02 15:04:05", ass.Datetime)
		if err == nil {
			var ctime gtime.Time
			ctime.Scan(intTime)
			ass.CTime = ctime - 8*3600 // adjust timezone for fe
		}
	}
	return
}

// LiveCheckAssist check assist in live room
func (d *Dao) LiveCheckAssist(c context.Context, mid, assistMid int64, ip string) (isAss int8, err error) {
	var req *http.Request
	params := url.Values{}
	params.Set("uid", strconv.FormatInt(assistMid, 10))
	params.Set("anchor_id", strconv.FormatInt(mid, 10))
	if req, err = http.NewRequest("POST", d.liveCheckAssURL+"?"+params.Encode(), nil); err != nil {
		log.Error("LiveCheckAssist url(%s) error(%v)", d.liveCheckAssURL+"?"+params.Encode(), err)
		err = ecode.CreativeLiveErr
		return
	}
	req.Header.Set("X-BiliLive-UID", strconv.FormatInt(mid, 10))
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Do(c, req, &res); err != nil {
		log.Error("LiveCheckAssist url(%s) response(%+v) error(%v)", d.liveCheckAssURL+"?"+params.Encode(), res, err)
		err = ecode.CreativeLiveErr
		return
	}
	if res.Code == 120013 {
		log.Error("LiveCheckAssist url(%s) mid(%d) not opened(%v)", d.liveCheckAssURL+"?"+params.Encode(), mid, 120013)
		err = ecode.CreativeLiveNotOpenErr
		return
	}
	if res.Code != 0 {
		log.Error("LiveCheckAssist url(%s) res(%v)", d.liveCheckAssURL+"?"+params.Encode(), res)
		err = ecode.Int(res.Code)
		return
	}
	isAss = 1
	return
}
