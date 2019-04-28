package http

import (
	"context"
	"strconv"
	"time"

	pushmdl "go-common/app/service/main/push/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	xtime "go-common/library/time"
)

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
		alertTitle = "哔哩哔哩消息"
	}
	alertBody := params.Get("alert_body")
	if alertBody == "" {
		c.JSON(nil, ecode.RequestErr)
		log.Error("alert_body is empty")
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
	sound, vibration := pushmdl.SwitchOn, pushmdl.SwitchOn
	if params.Get("sound") != "" {
		if sd, _ := strconv.Atoi(params.Get("sound")); sd == pushmdl.SwitchOff {
			sound = pushmdl.SwitchOff
		}
	}
	if params.Get("vibration") != "" {
		if vr, _ := strconv.Atoi(params.Get("vibration")); vr == pushmdl.SwitchOff {
			vibration = pushmdl.SwitchOff
		}
	}
	passThrough, _ := strconv.Atoi(params.Get("pass_through"))
	if passThrough != pushmdl.SwitchOn {
		passThrough = pushmdl.SwitchOff
	}
	img := params.Get("image_url")
	info := &pushmdl.PushInfo{
		TaskID:      pushmdl.TempTaskID(),
		APPID:       appID,
		Title:       alertTitle,
		Summary:     alertBody,
		LinkType:    int8(linkType),
		LinkValue:   linkValue,
		ExpireTime:  xtime.Time(expireTime),
		PassThrough: passThrough,
		Sound:       sound,
		Vibration:   vibration,
		ImageURL:    img,
	}
	c.JSON(nil, pushSrv.TestToken(context.Background(), info, token))
}
