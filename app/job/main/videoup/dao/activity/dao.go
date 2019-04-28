package activity

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/job/main/videoup/conf"
	"go-common/app/job/main/videoup/model/archive"
	"go-common/library/log"
	xhttp "go-common/library/net/http/blademaster"
)

// Dao is message dao.
type Dao struct {
	c      *conf.Config
	client *xhttp.Client
	AddUri string
	UpUri  string
}

// New new a activity dao.
func New(c *conf.Config) (d *Dao) {
	// http://act.bilibili.com/api/likes/video/add
	d = &Dao{
		c:      c,
		client: xhttp.NewClient(c.HTTPClient),
		AddUri: c.Host.Act + "/api/likes/video/add/",
		UpUri:  c.Host.Act + "/api/likes/upbyaid/",
	}
	return
}

// AddVideo add video to activity.
func (d *Dao) AddVideo(c context.Context, a *archive.Archive, missionID int64) (err error) {
	params := url.Values{}
	params.Set("aid", strconv.FormatInt(a.Aid, 10))
	params.Set("mid", strconv.FormatInt(a.Mid, 10))
	params.Set("message", a.Title)
	params.Set("image", a.Cover)
	params.Set("type", strconv.FormatInt(int64(a.TypeID), 10))
	var res struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	if err = d.client.Post(c, d.AddUri+strconv.FormatInt(missionID, 10), "", params, &res); err != nil {
		log.Error("d.client.Post error(%v)", err)
		return
	}
	if res.Code != 0 {
		log.Error("url(%s) res code(%d) or res.result(%v)", d.AddUri+strconv.FormatInt(missionID, 10)+"?"+params.Encode(), res.Code, res.Msg)
	}
	return
}

// UpVideo update video to activity.
func (d *Dao) UpVideo(c context.Context, a *archive.Archive, missionID int64) (err error) {
	params := url.Values{}
	params.Set("aid", strconv.FormatInt(a.Aid, 10))
	params.Set("mission_id", strconv.FormatInt(missionID, 10))
	params.Set("state", "-1")
	var res struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	if err = d.client.Post(c, d.UpUri+strconv.FormatInt(missionID, 10), "", params, &res); err != nil {
		log.Error("d.client.Post error(%v)", err)
		return
	}
	if res.Code != 0 {
		log.Error("url(%s) res code(%d) or res.result(%v)", d.UpUri+strconv.FormatInt(missionID, 10)+"?"+params.Encode(), res.Code, res.Msg)
	}
	return
}
