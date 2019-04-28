package http

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go-common/app/interface/main/app-wall/dao/padding"
	"go-common/app/interface/main/app-wall/model"
	"go-common/library/ecode"
	log "go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

// wallShikeClick
func wallShikeClick(c *bm.Context) {
	params := c.Request.Form
	// ip uint32, cid, mac, idfa, cb
	appid := params.Get("appid")
	mac := params.Get("mac")
	idfa := params.Get("idfa")
	cb := params.Get("cb")
	cid := model.ChannelShike
	if appid != model.GdtIOSAppID {
		log.Error("gdt click wrong appid(%s)", appid)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if cid == "" || idfa == "" || cb == "" {
		log.Error("arg cid(%s) mac(%s) idfa(%s) callback(%s) is empty", cid, idfa, cb)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	ip := model.InetAtoN(metadata.String(c, metadata.RemoteIP))
	if err := offerSvc.Click(c, ip, cid, mac, idfa, cb, time.Now()); err != nil {
		log.Error("offerSvc.Click error(%v)", err)
		c.JSON(nil, ecode.NotModified)
		return
	}
	res := map[string]interface{}{}
	res["code"] = ecode.OK
	returnDataJSON(c, res, nil)
}

// wallGdtClick
func wallGdtClick(c *bm.Context) {
	params := c.Request.Form
	appID := params.Get("appid")
	muid := params.Get("muid")
	var now time.Time
	ts, _ := strconv.ParseInt(params.Get("click_time"), 10, 64)
	if ts != 0 {
		now = time.Unix(ts, 0)
	} else {
		now = time.Now()
	}
	clickid := params.Get("click_id")
	adid := params.Get("advertiser_id")
	appType := params.Get("app_type")
	if id, ok := model.AppIDGdt[appType]; !ok || appID != id {
		log.Error("wallGdtClick app_type(%s) and appid(%s) is illegal", appType, appID)
		res := map[string]interface{}{"ret": ecode.RequestErr}
		returnDataJSON(c, res, nil)
		return
	}
	if _, ok := model.ChannelGdt[adid]; !ok {
		log.Error("wallGdtClick advertiser_id(%s) is illegal", adid)
		res := map[string]interface{}{"ret": ecode.RequestErr}
		returnDataJSON(c, res, nil)
		return
	}
	ip := model.InetAtoN(metadata.String(c, metadata.RemoteIP))
	if appType == model.TypeIOS {
		if err := offerSvc.Click(c, ip, adid, "", muid, clickid, now); err != nil {
			log.Error("wallGdtClick %+v", err)
			res := map[string]interface{}{"ret": ecode.ServerErr}
			returnDataJSON(c, res, nil)
			return
		}
	} else if appType == model.TypeAndriod {
		if err := offerSvc.ANClick(c, adid, muid, "", "", clickid, ip, now); err != nil {
			log.Error("wallGdtClick %+v", err)
			res := map[string]interface{}{"ret": ecode.ServerErr}
			returnDataJSON(c, res, nil)
			return
		}
	}
	res := map[string]interface{}{"ret": ecode.OK}
	returnDataJSON(c, res, nil)
}

// wallTestActive
func wallTestActive(c *bm.Context) {
	params := c.Request.Form
	idfa := params.Get("idfa")
	if idfa == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	ip := model.InetAtoN(metadata.String(c, metadata.RemoteIP))
	if err := offerSvc.Active(c, ip, "0", "", "", idfa, "test by haogw", time.Now()); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	res := map[string]interface{}{
		"code": ecode.OK,
	}
	returnDataJSON(c, res, nil)
}

// wallActive
func wallActive(c *bm.Context) {
	var (
		_aesKey = []byte("iHiRzNVIy8mvc4HnKPKeiJEP90zycllH")
		_aesIv  = []byte("w4gQHf5M7RQdBq2U")
	)
	body := c.Request.Form.Get("body")
	if body == "" {
		log.Error("body is empty")
		c.JSON(nil, ecode.RequestErr)
		return
	}
	log.Info("body before base64 (%v)", body)
	paramBody := c.Request.Form
	appkey := paramBody.Get("appkey")
	appver := paramBody.Get("appver")
	build := paramBody.Get("build")
	device := paramBody.Get("device")
	filtered := paramBody.Get("filtered")
	mobiApp := paramBody.Get("mobi_app")
	sign := paramBody.Get("sign")
	log.Info("appkey: (%v)\n appver: (%v)\n build: (%v)\n device: (%v)\n filtered: (%v)\n mobi_app: (%v)\n sign: (%v)", appkey, appver, build, device, filtered, mobiApp, sign)
	bs, err := base64.StdEncoding.DecodeString(body)
	if err != nil {
		log.Error("base64.StdEncoding.DecodeString(%s) error(%v)", body, err)
		c.JSON(nil, ecode.SignCheckErr)
		return
	}
	bs, err = offerSvc.CBCDecrypt(bs, _aesKey, _aesIv, padding.PKCS5)
	if err != nil {
		log.Error("aes.CBCDecrypt(%s) error(%v)", base64.StdEncoding.EncodeToString(bs), err)
		c.JSON(nil, ecode.SignCheckErr)
		return
	}
	params, err := url.ParseQuery(string(bs))
	if err != nil {
		log.Error("url.ParseQuery(%s) error(%v)", bs, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	rmac := params.Get("rmac")
	mac := params.Get("mac")
	idfa := params.Get("idfa")
	name := params.Get("name")
	log.Info("rmac: (%v)\n mac: (%v)\n idfa: (%v)\n name: (%v)\n", rmac, mac, idfa, name)
	if idfa == "" || name == "" || strings.Contains(idfa, "-") || strings.Contains(rmac, ":") || strings.Contains(mac, ":") {
		log.Error("mac(%s) rmac(%s) idfa(%s) format error, or idfa(%s) name(%s) is empty", mac, rmac, idfa, idfa, name)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	ip := model.InetAtoN(metadata.String(c, metadata.RemoteIP))
	var (
		mid int64
	)
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	log.Info("mid: (%d)", mid)
	err = offerSvc.Active(c, ip, strconv.FormatInt(mid, 10), rmac, mac, idfa, name, time.Now())
	if err != nil {
		log.Error("offerSvc.Active(%d, %d, %s, %s, %s, %s) error(%v)", ip, mid, rmac, mac, idfa, name, err)
		c.JSON(nil, ecode.NotModified)
		return
	}
	log.Info("wall offer active idfa(%s)", idfa)
	res := map[string]interface{}{
		"code": ecode.OK,
	}
	returnDataJSON(c, res, nil)
}

func wallExist(c *bm.Context) {
	params := c.Request.Form
	// ip uint32, cid, mac, idfa, cb
	appID := params.Get("appid")
	idfa := params.Get("idfa")
	if (appID != model.GdtIOSAppID && appID != model.GdtAndroidAppID) || idfa == "" {
		log.Error("arg cid(%s) idfa(%s) is empty", appID, idfa)
		c.JSON(nil, ecode.RequestErr)
		return
	}

	exist, _ := offerSvc.Exists(c, strings.Replace(idfa, "-", "", -1))
	const (
		ret = `{"%s":"%d"}`
	)
	e := 1
	if !exist {
		e = 0
	}
	bs := []byte(fmt.Sprintf(ret, idfa, e))
	if _, err := c.Writer.Write(bs); err != nil {
		log.Error("w.Write error(%v)", err)
	}
	res := map[string]interface{}{
		"code": ecode.OK,
	}
	returnDataJSON(c, res, nil)
}

func wallDotinappClick(c *bm.Context) {
	params := c.Request.Form
	appid := params.Get("appid")
	// pid := params.Get("pid")
	// subid := params.Get("subid")
	// ua := params.Get("ua")
	idfa := params.Get("idfa")
	clickid := params.Get("clickid")
	pid := model.ChannelDontin
	if appid != model.GdtIOSAppID {
		log.Error("dotinapp click wrong appid(%s)", appid)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if pid == "" || idfa == "" || clickid == "" {
		log.Error("pid(%s) idfa(%s) clickid(%s) is empty", pid, idfa, clickid)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	ip := model.InetAtoN(metadata.String(c, metadata.RemoteIP))
	if err := offerSvc.Click(c, ip, pid, "", idfa, clickid, time.Now()); err != nil {
		log.Error("offerSvc.Click error(%v)", err)
		c.JSON(nil, ecode.NotModified)
		return
	}
	res := map[string]interface{}{
		"code": ecode.OK,
	}
	returnDataJSON(c, res, nil)
}

func wallToutiaoClick(c *bm.Context) {
	const (
		OSAndroid = "0"
		// OSIOS     = "1"
	)
	params := c.Request.Form
	os := params.Get("os")
	if os != OSAndroid {
		log.Error("wallToutiaoClick os(%s) is illegal", os)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	imei := params.Get("imei")
	androidid := params.Get("androidid")
	if imei == "" && androidid == "" {
		log.Error("wallToutiaoClick imei(%s) and androidid(%s) is illegal", imei, androidid)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	cb := params.Get("callback_url")
	if cb == "" {
		log.Error("wallToutiaoClick callback_url(%s) is illegal", cb)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	mac := params.Get("mac")
	ip := model.InetAtoN(params.Get("ip"))
	var now time.Time
	ts, _ := strconv.ParseInt(params.Get("timestamp"), 10, 64)
	if ts != 0 {
		now = time.Unix(0, ts*1e6)
	} else {
		now = time.Now()
	}
	if os == OSAndroid {
		if err := offerSvc.ANClick(c, model.ChannelToutiao, imei, androidid, mac, cb, ip, now); err != nil {
			log.Error("%+v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}
}

func wallActive2(c *bm.Context) {
	params := c.Request.Form
	os := params.Get("os")
	if _, ok := model.AppIDGdt[os]; !ok {
		log.Error("wallActive2 os(%s) is illegal", os)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if os != model.TypeAndriod {
		log.Error("wallActive2 os(%s) is illegal", os)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	imei := params.Get("imei")
	androidid := params.Get("androidid")
	if imei == "" && androidid == "" {
		log.Error("wallActive2 imei(%s) and androidid(%s) is illegal", imei, androidid)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	mac := params.Get("mac")
	err := offerSvc.ANActive(c, imei, androidid, mac, time.Now())
	if err != nil {
		log.Error("wallActive2 %+v", err)
	}
	c.JSON(nil, err)
}
