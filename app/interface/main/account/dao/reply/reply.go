package reply

import (
	"context"
	xhttp "net/http"
	"net/url"
	"strconv"

	"go-common/app/interface/main/account/conf"
	"go-common/app/interface/main/account/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

var (
	_replyHistoryURI  = "/x/internal/v2/reply/record"
	_activityPagesURI = "/activity/pages"
)

// Dao dao
type Dao struct {
	c             *conf.Config
	client        *bm.Client
	replyHistory  string
	activityPages string
}

// New Dao
func New(c *conf.Config) (d *Dao) {
	return &Dao{
		c:             c,
		client:        bm.NewClient(c.HTTPClient.Normal),
		replyHistory:  c.Host.API + _replyHistoryURI,
		activityPages: c.Host.WWW + _activityPagesURI,
	}
}

// ReplyHistoryList reply history list
func (d *Dao) ReplyHistoryList(c context.Context, mid int64, stime, etime, order, sort string, pn, ps int64, accessKey, cookie, ip string) (rhl *model.ReplyHistory, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("stime", stime)
	params.Set("etime", etime)
	params.Set("order", order)
	params.Set("sort", sort)
	params.Set("pn", strconv.FormatInt(pn, 10))
	params.Set("ps", strconv.FormatInt(ps, 10))
	params.Set("access_key", accessKey)
	req, err := d.client.NewRequest(xhttp.MethodGet, d.replyHistory, ip, params)
	if err != nil {
		return
	}
	var res struct {
		Code int `json:"code"`
		Data struct {
			Page struct {
				Num   int `json:"num"`
				Size  int `json:"size"`
				Total int `json:"total"`
			} `json:"page"`
			Records []struct {
				ID      int           `json:"id"`
				Oid     int64         `json:"oid"`
				Type    int64         `json:"type"`
				Floor   int           `json:"floor"`
				Like    int           `json:"like"`
				Rcount  int           `json:"rcount"`
				Mid     int64         `json:"mid"`
				State   int           `json:"state"`
				Message string        `json:"message"`
				Ctime   string        `json:"ctime"`
				Members []*model.Info `json:"members"`
			} `json:"records"`
		} `json:"data"`
		Message string `json:"message"`
		TTL     int    `json:"ttl"`
	}
	if err = d.client.Do(c, req, &res); err != nil {
		log.Error("member interface reply request reply history list failed, err(%v)", err)
		return
	}
	if res.Code != 0 {
		log.Error("member interface reply request reply history list code(%d), err(%v)", res.Code, err)
		err = ecode.Int(res.Code)
		return
	}
	rhl = &model.ReplyHistory{
		Page:    res.Data.Page,
		Records: make([]*model.Record, 0),
	}
	for _, v := range res.Data.Records {
		tme := make([]*model.Member, 0)
		for _, vt := range v.Members {
			m, _ := strconv.ParseInt(vt.Mid, 10, 64)
			tmp := &model.Member{
				Mid:   m,
				Uname: vt.Name,
			}
			tme = append(tme, tmp)
		}
		rhlt := &model.Record{
			ID:      v.ID,
			Oid:     v.Oid,
			OidStr:  strconv.FormatInt(v.Oid, 10),
			Type:    v.Type,
			Floor:   v.Floor,
			Like:    v.Like,
			Rcount:  v.Rcount,
			Mid:     v.Mid,
			State:   v.State,
			Message: v.Message,
			Ctime:   v.Ctime,
			Members: tme,
		}
		rhl.Records = append(rhl.Records, rhlt)
	}
	return
}

// ActivityPages activity pages api
func (d *Dao) ActivityPages(c context.Context, mid int64, aids []int64, accessKey, cookie, ip string) (at map[int64]*model.RecordAppend, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("pids", xstr.JoinInts(aids))
	params.Set("all", "isOne")
	params.Set("access_key", accessKey)
	req, err := d.client.NewRequest(xhttp.MethodGet, d.activityPages, ip, params)
	if err != nil {
		return
	}
	var res struct {
		Code int `json:"code"`
		Data struct {
			List []struct {
				ID    int64  `json:"id"`
				Name  string `json:"name"`
				PcURL string `json:"pc_url"`
			} `json:"list"`
		} `json:"data"`
	}
	at = make(map[int64]*model.RecordAppend)
	if err = d.client.Do(c, req, &res); err != nil {
		log.Error("member interface reply request activity failed, err(%v)", err)
		return
	}
	if res.Code != 0 {
		log.Error("member interface reply request activity code != 0, err(%v)", err)
		err = ecode.Int(res.Code)
		return
	}
	for _, v := range res.Data.List {
		at[v.ID] = &model.RecordAppend{
			Title: v.Name,
			URL:   v.PcURL,
		}
	}
	return
}
