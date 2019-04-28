package http

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"

	"go-common/app/service/main/push/dao/huawei"
	"go-common/app/service/main/push/dao/jpush"
	"go-common/app/service/main/push/dao/mi"
	"go-common/app/service/main/push/dao/oppo"
	pushmdl "go-common/app/service/main/push/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var mobis = map[string]int{
	"android":       pushmdl.MobiAndroid,
	"android_i":     pushmdl.MobiAndroid,
	"iphone":        pushmdl.MobiIPhone,
	"ipad":          pushmdl.MobiIPad,
	"android_comic": pushmdl.MobiAndroidComic,
}

func huaweiCallback(c *bm.Context) {
	var (
		err error
		res = make(map[string]interface{})
	)
	defer func() {
		c.JSONMap(res, err)
	}()
	bs, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		res["errno"] = ecode.RequestErr
		res["errmsg"] = "param error"
		log.Error("huawei callback param(%s) error", bs)
		return
	}
	cb := new(huawei.Callback)
	if err = json.Unmarshal(bs, &cb); err != nil {
		res["errno"] = ecode.RequestErr
		res["errmsg"] = "param error"
		log.Error("huawei callback param(%s) error(%v)", bs, err)
		return
	}
	if err = pushSrv.CallbackHuawei(c, cb); err != nil {
		res["errno"] = err
		res["errmsg"] = err.Error()
		log.Error("huawei callback param(%s) error(%v)", bs, err)
		return
	}
	res["errno"] = ecode.OK
	res["errmsg"] = "success"
}

func miCallback(c *bm.Context) {
	param := c.Request.Form
	d := param.Get("data")
	if d == "" {
		c.JSON(nil, ecode.RequestErr)
		log.Error("xiaomi callback param empty")
		return
	}
	m := map[string]*mi.Callback{}
	if err := json.Unmarshal([]byte(d), &m); err != nil {
		c.JSON(nil, ecode.RequestErr)
		log.Error("xiaomi callback param(%s) error(%v)", d, err)
		return
	}
	c.JSON(nil, pushSrv.CallbackXiaomi(c, m))
}

func miRegidCallback(c *bm.Context) {
	params := c.Request.Form
	appID := params.Get("app_id")
	if appID == "" {
		c.JSON(nil, ecode.RequestErr)
		log.Error("app_id is empty")
		return
	}
	appVer := params.Get("app_version")
	if appVer == "" {
		c.JSON(nil, ecode.RequestErr)
		log.Error("app_version is empty")
		return
	}
	appPkg := params.Get("app_pkg")
	if appPkg == "" {
		c.JSON(nil, ecode.RequestErr)
		log.Error("app_pkg is empty")
		return
	}
	regid := params.Get("regid")
	if regid == "" {
		c.JSON(nil, ecode.RequestErr)
		log.Error("regid is empty")
		return
	}
	auth := strings.TrimPrefix(c.Request.Header.Get("Authorization"), "key=")
	if auth == "" {
		c.JSON(nil, ecode.RequestErr)
		log.Error("auth is empty")
		return
	}
	cb := &mi.RegidCallback{
		AppID:     appID,
		AppVer:    appVer,
		AppPkg:    appPkg,
		AppSecret: auth,
		Regid:     regid,
	}
	c.JSON(nil, pushSrv.CallbackXiaomiRegid(c, cb))
}

func oppoCallback(c *bm.Context) {
	bs, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		log.Error("huawei callback param(%s) error", bs)
		return
	}
	var cb []*oppo.Callback
	if err = json.Unmarshal(bs, &cb); err != nil {
		c.JSON(nil, ecode.RequestErr)
		log.Error("oppo callback param(%s) error(%v)", bs, err)
		return
	}
	param := c.Request.Form
	task := param.Get("task")
	if task == "" {
		c.JSON(nil, ecode.RequestErr)
		log.Error("oppo callback task empty")
		return
	}
	c.JSON(nil, pushSrv.CallbackOppo(c, task, cb))
}

func jpushCallback(ctx *bm.Context) {
	bs, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Error("jpush batch callback param(%s) error(%v)", bs, err)
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	var cbs []*jpush.CallbackReply
	if err = json.Unmarshal(bs, &cbs); err != nil {
		log.Error("jpush batch callback param(%s) error(%v)", bs, err)
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	ctx.JSON(nil, pushSrv.CallbackJpush(ctx, cbs))
}

func iOSCallback(c *bm.Context) {
	params := c.Request.Form
	task := params.Get("task")
	if task == "" {
		c.JSON(nil, ecode.RequestErr)
		log.Error("task is empty")
		return
	}
	p := params.Get("mobi_app")
	if p == "" {
		c.JSON(nil, ecode.RequestErr)
		log.Error("mobi_app is empty")
		return
	}
	token := params.Get("tid")
	if token == "" {
		c.JSON(nil, ecode.RequestErr)
		log.Error("tid is empty")
		return
	}
	c.JSON(nil, pushSrv.CallbackIOS(c, task, token, mobis[p]))
}

func androidCallback(c *bm.Context) {
	// Warn: 刚开始接极光时，极光没有回执功能，只能客户端来做送达回执，数据不准确
	// 如今极光提供了送达回执的功能，数据不能重复记录，所以这里直接return
	// 目前这个方法只有极光送达回执在调用(调用方是客户端)，直接停用
	c.JSON(nil, nil)
}

func clickCallback(c *bm.Context) {
	var (
		err      error
		pid      int
		platform int
		params   = c.Request.Form
	)
	task := params.Get("task")
	if task == "" {
		c.JSON(nil, ecode.RequestErr)
		log.Error("task is empty")
		return
	}
	app, _ := strconv.ParseInt(params.Get("app"), 10, 64)
	if app <= 0 {
		c.JSON(nil, ecode.RequestErr)
		log.Error("app error(%s)", params.Get("app"))
		return
	}
	mid, _ := strconv.ParseInt(params.Get("mid"), 10, 64)
	p := params.Get("mobi_app")
	if p == "" {
		c.JSON(nil, ecode.RequestErr)
		log.Error("mobi_app is empty")
		return
	}
	pid, ok := mobis[p]
	if !ok {
		c.JSON(nil, ecode.RequestErr)
		log.Error("mobi_app error(%s)", params.Get("mobi_app"))
		return
	}
	sdk, _ := strconv.Atoi(params.Get("push_sdk"))
	if sdk == 0 {
		c.JSON(nil, ecode.RequestErr)
		log.Error("push_sdk error(%s)", params.Get("push_sdk"))
		return
	}
	switch sdk {
	case pushmdl.PushSDKApns:
		if pid == pushmdl.MobiIPhone {
			platform = pushmdl.PlatformIPhone
		} else if pid == pushmdl.MobiIPad {
			platform = pushmdl.PlatformIPad
		}
	case pushmdl.PushSDKXiaomi:
		platform = pushmdl.PlatformXiaomi
	case pushmdl.PushSDKHuawei:
		platform = pushmdl.PlatformHuawei
	case pushmdl.PushSDKOppo:
		platform = pushmdl.PlatformOppo
	case pushmdl.PushSDKJpush:
		platform = pushmdl.PlatformJpush
	default:
		c.JSON(nil, ecode.RequestErr)
		log.Error("invalid push_sdk value(%d)", sdk)
		return
	}
	token, err := url.QueryUnescape(params.Get("token"))
	if err != nil || token == "" {
		c.JSON(nil, ecode.RequestErr)
		log.Error("token(%s) error(%v)", params.Get("token"), err)
		return
	}
	buvid := c.Request.Header.Get("Buvid")
	if buvid == "" {
		if buvid, err = url.QueryUnescape(params.Get("buvid")); buvid == "" {
			c.JSON(nil, ecode.RequestErr)
			log.Error("buvid(%s) error(%v)", params.Get("buvid"), err)
			return
		}
	}
	cb := &pushmdl.Callback{
		Task:     task,
		APP:      app,
		Platform: platform,
		Mid:      mid,
		Pid:      pid,
		Token:    token,
		Buvid:    buvid,
		Click:    1,
	}
	c.JSON(nil, pushSrv.CallbackClick(c, cb))
}
