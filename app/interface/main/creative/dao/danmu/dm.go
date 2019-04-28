package danmu

import (
	"context"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/dao/tool"
	"go-common/app/interface/main/creative/model/danmu"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_dmSearchURI   = "/x/internal/v2/dm/search"
	_dmEditURI     = "/x/internal/v2/dm/edit/state"
	_dmTransferURI = "/x/internal/dm/up/transfer"
	_dmPoolURI     = "/x/internal/v2/dm/edit/pool"
	_dmDistriURI   = "/x/internal/v2/dm/distribution"
	_dmRecentURI   = "/x/internal/v2/dm/recent"
)

// List fn
func (d *Dao) List(c context.Context, cid, mid int64, page, size int, order, pool, midStr, ip string) (dmList *danmu.DmList, err error) {
	params := url.Values{}
	params.Set("type", "1")
	params.Set("oid", strconv.FormatInt(cid, 10))
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("pn", strconv.FormatInt(int64(page), 10))
	params.Set("ps", strconv.FormatInt(int64(size), 10))
	params.Set("order", order)
	params.Set("pool", pool)
	if len(midStr) > 0 {
		midss, _ := xstr.SplitInts(midStr)
		midsa := make([]int64, 0)
		for _, v := range midss {
			if v > 0 {
				midsa = append(midsa, v)
			}
		}
		if len(midsa) > 0 {
			params.Set("mids", xstr.JoinInts(midsa))
		}
	}

	var res struct {
		Code           int                   `json:"code"`
		SearchDMResult *danmu.SearchDMResult `json:"data"`
	}
	dmList = &danmu.DmList{
		List: make([]*danmu.MemberDM, 0),
	}
	if err = d.client.Get(c, d.dmSearchURL, ip, params, &res); err != nil {
		log.Error("d.DmSearch Get(%s,%s,%s) err(%v)", d.dmSearchURL, ip, params.Encode(), err)
		err = ecode.CreativeDanmuErr
		return
	}
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("d.DmSearch Get(%s,%s,%s) err(%v)|code(%d)", d.dmSearchURL, ip, params.Encode(), err, res.Code)
		if err == ecode.NothingFound {
			log.Error("d.DmSearch Get NothingFound (%s,%s,%s) err(%v)|code(%d)", d.dmSearchURL, ip, params.Encode(), err, res.Code)
			err = nil
		}
		return
	}
	if res.SearchDMResult != nil {
		dmList.Page = res.SearchDMResult.Page.Num
		dmList.Size = res.SearchDMResult.Page.Size
		dmList.TotalItems = res.SearchDMResult.Page.Total
		dmList.TotalPages = int(math.Ceil(float64(res.SearchDMResult.Page.Total) / float64(res.SearchDMResult.Page.Size)))
		for _, v := range res.SearchDMResult.Result {
			list := &danmu.MemberDM{
				ID:       v.ID,
				FontSize: v.FontSize,
				Color:    v.Color,
				Mode:     v.Mode,
				Msg:      v.Msg,
				Oid:      v.Oid,
				Mid:      v.Mid,
				Playtime: float64(v.Progress) / 1000,
				Pool:     v.Pool,
				Ctime:    v.Ctime,
				Attrs:    v.Attrs,
			}
			if d.isProtect(v.Attrs, 1) {
				list.State = 2
			}
			dmList.List = append(dmList.List, list)
		}
	}
	return
}

// Edit fn
func (d *Dao) Edit(c context.Context, mid, cid int64, state int8, dmids []int64, ip string) (err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("oid", strconv.FormatInt(cid, 10))
	params.Set("type", "1")
	params.Set("state", strconv.FormatInt(int64(state), 10))
	params.Set("dmids", xstr.JoinInts(dmids))
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Post(c, d.dmEditURL, ip, params, &res); err != nil {
		log.Error("d.DmEdit.Post(%s,%s,%s) err(%v)", d.dmEditURL, ip, params.Encode(), err)
		err = ecode.CreativeDanmuErr
		return
	}
	log.Info("d.DmEdit.Post res.Code (%s,%s,%s) err(%v)|code(%d)", d.dmEditURL, ip, params.Encode(), err, res.Code)
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("d.DmEdit.Post res.Code (%s,%s,%s) err(%v)|code(%d)", d.dmEditURL, ip, params.Encode(), err, res.Code)
		return
	}
	return
}

// Transfer fn
func (d *Dao) Transfer(c context.Context, mid, fromCID, toCID int64, offset float64, ak, ck, ip string) (err error) {
	params := url.Values{}
	params.Set("from", strconv.FormatInt(fromCID, 10))
	params.Set("to", strconv.FormatInt(toCID, 10))
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("offset", strconv.FormatFloat(offset, 'f', 3, 64))
	params.Set("appkey", conf.Conf.App.Key)
	params.Set("appsecret", conf.Conf.App.Secret)
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	var (
		query, _ = tool.Sign(params)
		url      string
	)
	url = d.dmTransferURL + "?" + query
	var res struct {
		Code int `json:"code"`
	}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Error("http.NewRequest(%s) error(%v)|mid(%d),fromCID(%d),toCID(%d),ak(%s),ck(%s),ip(%s)", url, err, mid, fromCID, toCID, ak, ck, ip)
		err = ecode.CreativeDanmuErr
		return
	}
	req.Header.Set("X-BACKEND-BILI-REAL-IP", ip)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", ck)
	if err = d.client.Do(c, req, &res); err != nil {
		log.Error("d.client.Do(%s) error(%v)|mid(%d),fromCID(%d),toCID(%d),ak(%s),ck(%s),ip(%s)", url, err, mid, fromCID, toCID, ak, ck, ip)
		err = ecode.CreativeDanmuErr
		return
	}
	log.Info("dao Transfer(%s)mid(%d),fromCID(%d),toCID(%d),ak(%s),ck(%s),ip(%s)", url, mid, fromCID, toCID, ak, ck, ip)
	if res.Code != 0 {
		log.Error("dm Transfer(%s) res(%v); mid(%d), ip(%s), code(%d)", url, res, mid, ip, res.Code)
		err = ecode.Int(res.Code)
		return
	}
	return
}

// UpPool fn
func (d *Dao) UpPool(c context.Context, mid, cid int64, dmids []int64, pool int8) (err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("oid", strconv.FormatInt(cid, 10))
	params.Set("type", "1")
	params.Set("pool", strconv.FormatInt(int64(pool), 10))
	params.Set("dmids", xstr.JoinInts(dmids))
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Post(c, d.dmPoolURL, "", params, &res); err != nil {
		log.Error("d.Pool.Post(%s,%s,%s) err(%v)", d.dmPoolURL, "", params.Encode(), err)
		err = ecode.CreativeDanmuErr
		return
	}
	log.Info("d.Pool.Post res.Code (%s,%s,%s) err(%v)|code(%d)", d.dmPoolURL, "", params.Encode(), err, res.Code)
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("d.DmEdit.Post res.Code (%s,%s,%s) err(%v)|code(%d)", d.dmPoolURL, "", params.Encode(), err, res.Code)
		return
	}
	return
}

// Distri fn
func (d *Dao) Distri(c context.Context, mid, cid int64, ip string) (distri map[int64]int64, err error) {
	params := url.Values{}
	params.Set("oid", strconv.FormatInt(cid, 10))
	params.Set("type", "1")
	params.Set("interval", "1")
	var res struct {
		Code   int             `json:"code"`
		Distri map[int64]int64 `json:"data"`
	}
	res.Distri = make(map[int64]int64)
	if err = d.client.Get(c, d.dmDistriURL, ip, params, &res); err != nil {
		log.Error("d.DmDistri.Get(%s,%s,%s) err(%v)", d.dmDistriURL, ip, params.Encode(), err)
		err = ecode.CreativeDanmuErr
		return
	}
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("d.DmDistri.Get(%s,%s,%s) err(%v)|code(%d)", d.dmDistriURL, ip, params.Encode(), err, res.Code)
		return
	}
	distri = res.Distri
	return
}

// Recent fn
func (d *Dao) Recent(c context.Context, mid, pn, ps int64, ip string) (dmRecent *danmu.DmRecent, aids []int64, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("pn", strconv.FormatInt(pn, 10))
	params.Set("ps", strconv.FormatInt(ps, 10))
	dmRecent = &danmu.DmRecent{
		List: make([]*danmu.Recent, 0),
	}
	aids = make([]int64, 0)
	var res struct {
		Code         int                 `json:"code"`
		ResNewRecent *danmu.ResNewRecent `json:"data"`
	}
	if err = d.client.Get(c, d.dmRecentURL, ip, params, &res); err != nil {
		log.Error("d.DmRecent Get(%s,%s,%s) err(%v)", d.dmRecentURL, ip, params.Encode(), err)
		err = ecode.CreativeDanmuErr
		return
	}
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("d.DmRecent Get(%s,%s,%s) err(%v)|code(%d)", d.dmRecentURL, ip, params.Encode(), err, res.Code)
		return
	}
	if res.ResNewRecent != nil {
		dmRecent.Page = res.ResNewRecent.Page.Pn
		dmRecent.Size = res.ResNewRecent.Page.Ps
		dmRecent.TotalItems = res.ResNewRecent.Page.Total
		dmRecent.TotalPages = int(math.Ceil(float64(res.ResNewRecent.Page.Total) / float64(res.ResNewRecent.Page.Ps)))
		for _, v := range res.ResNewRecent.Result {
			list := &danmu.Recent{
				ID:       v.ID,
				Aid:      v.Aid,
				Type:     v.Type,
				Oid:      v.Oid,
				Mid:      v.Mid,
				Msg:      v.Msg,
				Playtime: float64(v.Progress) / 1000,
				FontSize: v.FontSize,
				Color:    v.Color,
				Mode:     v.Mode,
				Pool:     v.Pool,
				Title:    v.Title,
				Ctime:    v.Ctime,
				Mtime:    v.Ctime,
				Attrs:    v.Attrs,
			}
			if d.isProtect(v.Attrs, 1) {
				list.State = 2
			}
			dmRecent.List = append(dmRecent.List, list)
			aids = append(aids, v.Aid)
		}
	}
	return
}

func (d *Dao) isProtect(attrs string, num int64) bool {
	if len(attrs) == 0 {
		return false
	}
	attrInts, _ := xstr.SplitInts(attrs)
	for _, v := range attrInts {
		if v == num {
			return true
		}
	}
	return false
}
