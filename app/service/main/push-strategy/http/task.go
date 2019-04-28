package http

import (
	"net/url"
	"strconv"
	"strings"
	"time"

	pushmdl "go-common/app/service/main/push/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	xtime "go-common/library/time"
)

func addTask(c *bm.Context) {
	var (
		id   int64
		err  error
		req  = c.Request
		auth = req.Header.Get("Authorization")
	)
	req.ParseMultipartForm(500 * 1024 * 1024) // 500M
	appID, _ := strconv.ParseInt(req.FormValue("app_id"), 10, 64)
	if appID < 1 {
		log.Error("app_id is wrong: %s", req.FormValue("app_id"))
		c.JSON(nil, ecode.RequestErr)
		return
	}
	businessID, _ := strconv.ParseInt(req.FormValue("business_id"), 10, 64)
	if businessID < 1 {
		log.Error("business_id is wrong: %s", req.FormValue("business_id"))
		c.JSON(nil, ecode.RequestErr)
		return
	}
	platform := req.FormValue("platform")
	am, err := url.ParseQuery(auth)
	if err != nil {
		log.Error("parse Authorization(%s) error(%v)", auth, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	token := am.Get("token")
	if token == "" {
		log.Error("token is empty")
		c.JSON(nil, ecode.RequestErr)
		return
	}
	alertTitle := req.FormValue("alert_title")
	alertBody := req.FormValue("alert_body")
	if alertBody == "" {
		log.Error("alert_body is empty")
		c.JSON(nil, ecode.RequestErr)
		return
	}
	mids := req.FormValue("mids")
	if mids == "" {
		log.Error("mids is empty")
		c.JSON(nil, ecode.RequestErr)
		return
	}
	linkType, _ := strconv.Atoi(req.FormValue("link_type"))
	if linkType < 1 {
		log.Error("link_type is wrong: %s", req.FormValue("link_type"))
		c.JSON(nil, ecode.RequestErr)
		return
	}
	linkValue := req.FormValue("link_value")
	expireTime, _ := strconv.ParseInt(req.FormValue("expire_time"), 10, 64)
	if expireTime == 0 {
		expireTime = time.Now().Add(3 * 24 * time.Hour).Unix()
	}
	pushTime, _ := strconv.ParseInt(req.FormValue("push_time"), 10, 64)
	if pushTime == 0 {
		pushTime = time.Now().Unix()
	}
	passThrough, _ := strconv.Atoi(req.FormValue("pass_throught"))
	builds := req.FormValue("builds")
	group := req.FormValue("group")
	uuid := req.FormValue("uuid")
	imageURL := req.FormValue("image_url")
	task := &pushmdl.Task{
		Job:         pushmdl.JobName(time.Now().UnixNano(), alertBody, linkValue, group),
		Type:        pushmdl.TaskTypeStrategyMid,
		APPID:       appID,
		BusinessID:  businessID,
		Platform:    pushmdl.SplitInts(platform),
		Title:       alertTitle,
		Summary:     alertBody,
		LinkType:    int8(linkType),
		LinkValue:   linkValue,
		Build:       pushmdl.ParseBuild(builds),
		PassThrough: passThrough,
		PushTime:    xtime.Time(pushTime),
		ExpireTime:  xtime.Time(expireTime),
		Status:      pushmdl.TaskStatusPretreatmentPrepared,
		Group:       group,
		ImageURL:    imageURL,
		Extra:       &pushmdl.TaskExtra{Group: group},
	}
	log.Info("http add task(%d) uuid(%s) business(%d) link_value(%s) mids(%d)", task.Job, uuid, task.BusinessID, task.LinkValue, len(strings.Split(mids, ",")))
	if id, err = srv.AddTask(c, uuid, token, task, mids); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(id, nil)
}
