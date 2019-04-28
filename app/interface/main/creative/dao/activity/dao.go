package activity

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
	"strings"

	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/model/activity"
	"go-common/library/ecode"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
)

const (
	_activityAllList    = "/activity/list/videoall"
	_activityUpdate     = "/api/likes/upbyaid/%d"
	_activitySubject    = "/activity/subject/%d"
	_activityProtocol   = "/activity/protocol/%d"
	_actOnlineByTypeURI = "/activity/online/by/type"
	_actURLType16       = "https://www.bilibili.com/blackboard/x/activity-tougao-h5/detail?from=snap&id="
	_actLike            = "/activity/likes/list/%d"
)

// Dao  define
type Dao struct {
	c *conf.Config
	// http
	client *httpx.Client
	// uri
	ActAllListURI      string
	ActUpdateURI       string
	ActSubjectURI      string
	ActProtocolURI     string
	ActOnlineByTypeURL string
	ActLikeURI         string
}

// New init dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:                  c,
		client:             httpx.NewClient(c.HTTPClient.Normal),
		ActAllListURI:      c.Host.Activity + _activityAllList,
		ActUpdateURI:       c.Host.Act + _activityUpdate,
		ActSubjectURI:      c.Host.Activity + _activitySubject,
		ActProtocolURI:     c.Host.Act + _activityProtocol,
		ActOnlineByTypeURL: c.Host.Activity + _actOnlineByTypeURI,
		ActLikeURI:         c.Host.Activity + _actLike,
	}
	return
}

// MissionOnlineByTid fn
func (d *Dao) MissionOnlineByTid(c context.Context, tid, plat int16) (mm []*activity.ActWithTP, err error) {
	var res struct {
		Code int                   `json:"code"`
		Data []*activity.ActWithTP `json:"data"`
	}
	mm = make([]*activity.ActWithTP, 0)
	params := url.Values{}
	params.Set("type", strconv.Itoa(int(tid)))
	params.Set("plat", strconv.Itoa(int(plat)))
	if err = d.client.Get(c, d.ActOnlineByTypeURL, "", params, &res); err != nil {
		log.Error("videoup ActOnlineByTypeURL error(%v) | ActOnlineByTypeURL(%s)", err, d.ActOnlineByTypeURL+"?"+params.Encode())
		err = ecode.CreativeActivityErr
		return
	}
	if res.Code != 0 {
		log.Error("videoup ActOnlineByTypeURL res.Code nq zero error(%v) | ActOnlineByTypeURL(%s) res(%v)", res.Code, d.ActOnlineByTypeURL+"?"+params.Encode(), res)
		err = ecode.CreativeActivityErr
		return
	}
	for _, m := range res.Data {
		miss := &activity.ActWithTP{
			ID:       m.ID,
			Name:     m.Name,
			Hot:      m.Hot,
			Protocol: m.Protocol,
			Types:    m.Types,
		}
		if len(m.Tags) > 0 {
			miss.Tags = strings.Split(m.Tags, ",")[0]
		} else {
			miss.Tags = m.Name
		}
		if m.Type == 16 && len(m.ActURL) == 0 {
			miss.ActURL = _actURLType16 + strconv.FormatInt(miss.ID, 10)
		} else {
			miss.ActURL = m.ActURL
		}
		mm = append(mm, miss)
	}
	return
}

// Activities get activity list.
func (d *Dao) Activities(c context.Context) (act []*activity.Activity, err error) {
	var (
		url string
		res struct {
			Code int                  `json:"code"`
			Data []*activity.Activity `json:"data"`
		}
	)
	url = d.ActAllListURI
	if err = d.client.Get(c, url, "", nil, &res); err != nil {
		log.Error("ActivityList url(%s) response(%v) error(%v)", url, res, err)
		err = ecode.CreativeActivityErr
		return
	}
	if res.Code != 0 {
		log.Error("ActivityList url(%s) res(%v)", url, res)
		err = ecode.CreativeActivityErr
		return
	}
	for _, v := range res.Data {
		if v.Type == 16 && len(v.ActURL) == 0 {
			v.ActURL = _actURLType16 + strconv.FormatInt(v.ID, 10)
		}
	}
	act = res.Data
	return
}

// Unbind update the aid a status of this activity
func (d *Dao) Unbind(c context.Context, aid, missionID int64, ip string) (err error) {
	params := url.Values{}
	params.Set("mission_id", strconv.FormatInt(missionID, 10))
	params.Set("aid", strconv.FormatInt(aid, 10))
	params.Set("state", "-1")
	ms := json.RawMessage{}
	if err = d.client.RESTfulPost(c, d.ActUpdateURI, ip, params, &ms, missionID); err != nil {
		log.Error("ActUpdateURI url(%s) missionID(%d) error(%v)", d.ActUpdateURI, missionID, err)
		err = ecode.CreativeActivityErr
		return
	}
	log.Info("d.UpdateByAid url(%s) params(%s) res(%s)", d.ActUpdateURI, params.Encode(), string(ms))
	return
}

//Subject get any exist activity by missionID
func (d *Dao) Subject(c context.Context, missionID int64) (sub *activity.Activity, err error) {
	var res struct {
		Code int               `json:"code"`
		Data *activity.Subject `json:"data"`
	}
	if err = d.client.RESTfulGet(c, d.ActSubjectURI, "", url.Values{}, &res, missionID); err != nil {
		log.Error("ActSubjectURI url(%s) missionID(%d) error(%v)", d.ActSubjectURI, missionID, err)
		err = ecode.CreativeActivityErr
		return
	}
	if res.Code != 0 {
		log.Error("ActSubjectURI url(%s) res(%v)", d.ActSubjectURI, res)
		err = ecode.CreativeActivityErr
		return
	}
	var ID int
	if ID, err = strconv.Atoi(res.Data.ID); err != nil {
		log.Error("ActSubjectURI url(%s) res(%v)", d.ActSubjectURI, res)
		err = ecode.CreativeActivityErr
		return
	}
	sub = &activity.Activity{
		Name: res.Data.Name,
		ID:   int64(ID),
	}
	return
}

// Protocol fn
func (d *Dao) Protocol(c context.Context, missionID int64) (p *activity.Protocol, err error) {
	var res struct {
		Code int                `json:"code"`
		Data *activity.Protocol `json:"data"`
	}
	if err = d.client.RESTfulGet(c, d.ActProtocolURI, "", url.Values{}, &res, missionID); err != nil {
		log.Error("ActProtocolURI url(%s) missionID(%d) error(%v)", d.ActProtocolURI, missionID, err)
		err = ecode.CreativeActivityErr
		return
	}
	if res.Code != 0 {
		log.Error("ActProtocolURI url(%s) res(%v)", d.ActProtocolURI, res)
		err = ecode.CreativeActivityErr
		return
	}
	p = res.Data
	return
}

//Likes fn
func (d *Dao) Likes(c context.Context, missionID int64) (likeCnt int, err error) {
	var res struct {
		Code int            `json:"code"`
		Data *activity.Like `json:"data"`
	}
	if err = d.client.RESTfulGet(c, d.ActLikeURI, "", url.Values{}, &res, missionID); err != nil {
		log.Error("ActLikeURI url(%s) missionID(%d) error(%v)", d.ActLikeURI, missionID, err)
		err = ecode.CreativeActivityErr
		return
	}
	if res.Code != 0 {
		log.Error("ActLikeURI url(%s) res(%v)", d.ActLikeURI, res)
		err = ecode.CreativeActivityErr
		return
	}
	if res.Data != nil {
		likeCnt = res.Data.Count
	}
	return
}
