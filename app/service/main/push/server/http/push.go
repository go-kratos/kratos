package http

import (
	"context"
	"strconv"
	"time"

	"go-common/app/service/main/push/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	xtime "go-common/library/time"
	"go-common/library/xstr"
)

func push(c *bm.Context) {
	params := c.Request.Form
	appID, _ := strconv.ParseInt(params.Get("app_id"), 10, 64)
	if appID < 1 {
		c.JSON(nil, ecode.RequestErr)
		log.Error("app_id is wrong: %s", params.Get("app_id"))
		return
	}
	platform := params.Get("platform")
	alertTitle := params.Get("alert_title")
	if alertTitle != "" {
		res, err := pushSrv.Filter(c, alertTitle)
		if err == nil && res != alertTitle {
			log.Error("alertTitle(%s) contains invalid content", alertTitle)
			c.JSON(nil, ecode.PushSensitiveWordsErr)
			return
		}
	}
	alertBody := params.Get("alert_body")
	if alertBody == "" {
		c.JSON(nil, ecode.RequestErr)
		log.Error("alert_body is empty")
		return
	}
	res, err := pushSrv.Filter(c, alertBody)
	if err == nil && res != alertBody {
		log.Error("alertBody(%s) contains invalid content", alertBody)
		c.JSON(nil, ecode.PushSensitiveWordsErr)
		return
	}
	linkType, _ := strconv.Atoi(params.Get("link_type"))
	if linkType < 1 {
		c.JSON(nil, ecode.RequestErr)
		log.Error("link_type is wrong: %s", params.Get("link_type"))
		return
	}
	linkValue := params.Get("link_value")
	expireTime, _ := strconv.ParseInt(params.Get("expire_time"), 10, 64)
	if expireTime == 0 {
		expireTime = time.Now().Add(7 * 24 * time.Hour).Unix()
	}
	builds := params.Get("builds")
	sound, vibration := model.SwitchOn, model.SwitchOn
	if params.Get("sound") != "" {
		if sd, _ := strconv.Atoi(params.Get("sound")); sd == model.SwitchOff {
			sound = model.SwitchOff
		}
	}
	if params.Get("vibration") != "" {
		if vr, _ := strconv.Atoi(params.Get("vibration")); vr == model.SwitchOff {
			vibration = model.SwitchOff
		}
	}
	passThrough, _ := strconv.Atoi(params.Get("pass_through"))
	if passThrough != model.SwitchOn {
		passThrough = model.SwitchOff
	}
	mid := params.Get("mid")
	if mid == "" {
		log.Error("mid is empty", mid)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	mids, err := xstr.SplitInts(mid)
	if err != nil {
		log.Error("parse mid(%s) error(%v)", mid, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	task := &model.Task{
		ID:          model.TempTaskID(),
		Job:         model.JobName(time.Now().UnixNano(), alertBody, linkValue, ""),
		APPID:       appID,
		Platform:    model.SplitInts(platform),
		Title:       alertTitle,
		Summary:     alertBody,
		LinkType:    int8(linkType),
		LinkValue:   linkValue,
		PushTime:    xtime.Time(time.Now().Unix()),
		ExpireTime:  xtime.Time(expireTime),
		Build:       model.ParseBuild(builds),
		Sound:       sound,
		Vibration:   vibration,
		PassThrough: passThrough,
		Group:       params.Get("group"),
		ImageURL:    params.Get("image_url"),
	}
	go pushSrv.Pushs(context.Background(), task, mids)
	c.JSON(nil, nil)
}

func singlePush(c *bm.Context) {
	params := c.Request.Form
	appID, _ := strconv.ParseInt(params.Get("app_id"), 10, 64)
	if appID < 1 {
		c.JSON(nil, ecode.RequestErr)
		log.Error("app_id is wrong: %s", params.Get("app_id"))
		return
	}
	businessID, _ := strconv.ParseInt(params.Get("business_id"), 10, 64)
	if businessID < 1 {
		c.JSON(nil, ecode.RequestErr)
		log.Error("business_id is wrong: %s", params.Get("business_id"))
		return
	}
	token := params.Get("token")
	if token == "" {
		c.JSON(nil, ecode.RequestErr)
		log.Error("token is empty")
		return
	}
	platform := params.Get("platform")
	alertTitle := params.Get("alert_title")
	if alertTitle != "" {
		res, err := pushSrv.Filter(c, alertTitle)
		if err == nil && res != alertTitle {
			log.Error("alertTitle(%s) contains invalid content", alertTitle)
			c.JSON(nil, ecode.PushSensitiveWordsErr)
			return
		}
	}
	alertBody := params.Get("alert_body")
	if alertBody == "" {
		c.JSON(nil, ecode.RequestErr)
		log.Error("alert_body is empty")
		return
	}
	res, err := pushSrv.Filter(c, alertBody)
	if err == nil && res != alertBody {
		log.Error("alertBody(%s) contains invalid content", alertBody)
		c.JSON(nil, ecode.PushSensitiveWordsErr)
		return
	}
	linkType, _ := strconv.Atoi(params.Get("link_type"))
	if linkType < 1 {
		c.JSON(nil, ecode.RequestErr)
		log.Error("link_type is wrong: %s", params.Get("link_type"))
		return
	}
	linkValue := params.Get("link_value")
	expireTime, _ := strconv.ParseInt(params.Get("expire_time"), 10, 64)
	if expireTime == 0 {
		expireTime = time.Now().Add(7 * 24 * time.Hour).Unix()
	}
	builds := params.Get("builds")
	sound, vibration := model.SwitchOn, model.SwitchOn
	if params.Get("sound") != "" {
		if sd, _ := strconv.Atoi(params.Get("sound")); sd == model.SwitchOff {
			sound = model.SwitchOff
		}
	}
	if params.Get("vibration") != "" {
		if vr, _ := strconv.Atoi(params.Get("vibration")); vr == model.SwitchOff {
			vibration = model.SwitchOff
		}
	}
	passThrough, _ := strconv.Atoi(params.Get("pass_through"))
	if passThrough != model.SwitchOn {
		passThrough = model.SwitchOff
	}
	mid, _ := strconv.ParseInt(params.Get("mid"), 10, 64)
	if mid < 1 {
		c.JSON(nil, ecode.RequestErr)
		log.Error("mid is wrong (%d)", mid)
		return
	}
	task := &model.Task{
		ID:          model.TempTaskID(),
		Job:         model.JobName(time.Now().UnixNano(), alertBody, linkValue, ""),
		BusinessID:  businessID,
		APPID:       appID,
		Platform:    model.SplitInts(platform),
		Title:       alertTitle,
		Summary:     alertBody,
		LinkType:    int8(linkType),
		LinkValue:   linkValue,
		PushTime:    xtime.Time(time.Now().Unix()),
		ExpireTime:  xtime.Time(expireTime),
		Build:       model.ParseBuild(builds),
		Sound:       sound,
		Vibration:   vibration,
		PassThrough: passThrough,
		Group:       params.Get("group"),
		ImageURL:    params.Get("image_url"),
	}
	go pushSrv.SinglePush(context.Background(), token, task, mid)
	c.JSON(nil, nil)
}

func testToken(c *bm.Context) {
	params := c.Request.Form
	appID, _ := strconv.ParseInt(params.Get("app_id"), 10, 64)
	if appID < 1 {
		c.JSON(nil, ecode.RequestErr)
		log.Error("app_id is wrong: %s", params.Get("app_id"))
		return
	}
	alertTitle := params.Get("alert_title")
	if alertTitle == "" {
		alertTitle = model.DefaultMessageTitle
	}
	res, err := pushSrv.Filter(c, alertTitle)
	if err == nil && res != alertTitle {
		log.Error("alertTitle(%s) contains invalid content", alertTitle)
		c.JSON(nil, ecode.PushSensitiveWordsErr)
		return
	}
	alertBody := params.Get("alert_body")
	if alertBody == "" {
		c.JSON(nil, ecode.RequestErr)
		log.Error("alert_body is empty")
		return
	}
	res, err = pushSrv.Filter(c, alertBody)
	if err == nil && res != alertBody {
		log.Error("alertBody(%s) contains invalid content", alertBody)
		c.JSON(nil, ecode.PushSensitiveWordsErr)
		return
	}
	token := params.Get("token")
	if token == "" {
		c.JSON(nil, ecode.RequestErr)
		log.Error("token is empty")
		return
	}
	linkType, _ := strconv.Atoi(params.Get("link_type"))
	if linkType < 1 {
		c.JSON(nil, ecode.RequestErr)
		log.Error("link_type is wrong: %s", params.Get("link_type"))
		return
	}
	linkValue := params.Get("link_value")
	expireTime, _ := strconv.ParseInt(params.Get("expire_time"), 10, 64)
	if expireTime == 0 {
		expireTime = time.Now().Add(7 * 24 * time.Hour).Unix()
	}
	sound, vibration := model.SwitchOn, model.SwitchOn
	if params.Get("sound") != "" {
		if sd, _ := strconv.Atoi(params.Get("sound")); sd == model.SwitchOff {
			sound = model.SwitchOff
		}
	}
	if params.Get("vibration") != "" {
		if vr, _ := strconv.Atoi(params.Get("vibration")); vr == model.SwitchOff {
			vibration = model.SwitchOff
		}
	}
	passThrough, _ := strconv.Atoi(params.Get("pass_through"))
	if passThrough != model.SwitchOn {
		passThrough = model.SwitchOff
	}
	info := &model.PushInfo{
		TaskID:      model.TempTaskID(),
		Job:         model.JobName(time.Now().UnixNano(), alertBody, linkValue, ""),
		APPID:       appID,
		Title:       alertTitle,
		Summary:     alertBody,
		LinkType:    int8(linkType),
		LinkValue:   linkValue,
		ExpireTime:  xtime.Time(expireTime),
		PassThrough: passThrough,
		Sound:       sound,
		Vibration:   vibration,
		ImageURL:    params.Get("image_url"),
	}
	go pushSrv.TestToken(context.Background(), info, token)
	c.JSON(nil, nil)
}
