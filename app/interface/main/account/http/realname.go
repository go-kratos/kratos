package http

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"go-common/app/interface/main/account/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

// modify modify user relation.
func realnameStatus(c *bm.Context) {
	var (
		err    error
		mid, _ = c.Get("mid")
	)
	var resData struct {
		Status int8 `json:"status"`
	}
	if resData.Status, err = realnameSvc.Status(c, mid.(int64)); err != nil {
		log.Error("%+v", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(resData, nil)
}

func realnameApplyStatus(c *bm.Context) {
	var (
		err    error
		mid, _ = c.Get("mid")
	)
	var resData struct {
		Status   int8   `json:"status"`
		Remark   string `json:"remark"`
		Realname string `json:"realname"`
		Card     string `json:"card"`
	}
	if resData.Status, resData.Remark, resData.Realname, resData.Card, err = realnameSvc.ApplyStatus(c, mid.(int64)); err != nil {
		log.Error("%+v", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(resData, nil)
}

func realnameCardTypes(c *bm.Context) {
	var (
		err    error
		params = c.Request.Form

		platform = params.Get("platform")
		buildStr = params.Get("build")
		mobiapp  = params.Get("mobi_app")
		device   = params.Get("device")
		build    int
	)
	if build, err = strconv.Atoi(buildStr); err != nil {
		log.Error("%+v", errors.WithStack(err))
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var resData []*model.RealnameCardType
	if resData, err = realnameSvc.CardTypes(c, platform, mobiapp, device, build); err != nil {
		log.Error("%+v", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(resData, nil)
}

func realnameCardTypesV2(c *bm.Context) {
	var (
		err error
	)
	var resData []*model.RealnameCardType
	if resData, err = realnameSvc.CardTypesV2(c); err != nil {
		log.Error("%+v", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(resData, nil)
}

func realnameCountryList(c *bm.Context) {
	var (
		err error
	)
	var resData []*model.RealnameCountry
	if resData, err = realnameSvc.CountryList(c); err != nil {
		log.Error("%+v", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(resData, nil)
}

func realnameTelCapture(c *bm.Context) {
	var (
		err    error
		mid, _ = c.Get("mid")
	)
	if err = realnameSvc.TelCapture(c, mid.(int64)); err != nil {
		log.Error("%+v", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func realnameTelInfo(c *bm.Context) {
	var (
		err    error
		mid, _ = c.Get("mid")
	)
	var resData struct {
		Tel string `json:"tel"`
	}
	if resData.Tel, err = realnameSvc.TelInfo(c, mid.(int64)); err != nil {
		log.Error("%+v", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(resData, nil)
}

func realnameApply(c *bm.Context) {
	var (
		err    error
		mid, _ = c.Get("mid")
		params = c.Request.Form

		realname      = params.Get("real_name")
		cardTypeStr   = params.Get("card_type")
		cardType      int
		cardNum       = params.Get("card_num")
		countryIDStr  = params.Get("country")
		countryID     int
		captureStr    = params.Get("capture")
		capture       int
		handIMGToken  = params.Get("img1_token")
		frontIMGToken = params.Get("img2_token")
		backIMGToken  = params.Get("img3_token")
	)
	if cardType, err = strconv.Atoi(cardTypeStr); err != nil {
		log.Error("%+v", errors.WithStack(err))
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if countryIDStr == "" {
		countryID = 0 // 默认0：中国
	} else {
		if countryID, err = strconv.Atoi(countryIDStr); err != nil {
			log.Error("%+v", errors.WithStack(err))
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if countryID < 0 {
		countryID = 0
	}

	if capture, err = strconv.Atoi(captureStr); err != nil {
		log.Error("%+v", errors.WithStack(err))
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = realnameSvc.Apply(c, mid.(int64), realname, cardType, cardNum, countryID, capture, handIMGToken, frontIMGToken, backIMGToken); err != nil {
		log.Error("%+v")
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func realnameUpload(c *bm.Context) {
	var (
		mid, _ = c.Get("mid")
	)
	defer c.Request.Form.Del("img") // 防止日志不出现
	c.Request.ParseMultipartForm(32 << 20)
	imgBytes, err := func() ([]byte, error) {
		img := c.Request.FormValue("img")
		if img != "" {
			log.Info("Succeeded to parse img file from form value: mid: %d, length: %d", mid, len(img))
			return []byte(img), nil
		}
		log.Warn("Failed to parse img file from form value, fallback to form file: mid: %d", mid)
		f, _, err := c.Request.FormFile("img")
		if err != nil {
			return nil, errors.Wrapf(err, "parse img form file: mid: %d", mid)
		}
		defer f.Close()
		data, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, errors.Wrapf(err, "read img form file: mid: %d", mid)
		}
		if len(data) <= 0 {
			return nil, errors.Wrapf(err, "form file data: mid: %d, length: %d", mid, len(data))
		}
		log.Info("Succeeded to parse file from form file: mid: %d, length: %d", mid, len(data))
		return data, nil
	}()
	if err != nil {
		log.Error("Failed to parse realname upload file: mid: %d: %+v", mid, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var resData struct {
		SRC string `json:"token"`
	}
	if resData.SRC, err = realnameSvc.Upload(c, mid.(int64), imgBytes); err != nil {
		log.Error("%+v", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(resData, nil)
}

func realnamePreview(c *bm.Context) {
	var (
		img    []byte
		err    error
		mid, _ = c.Get("mid")

		params = c.Request.Form
		src    = params.Get("src")
	)
	if img, err = realnameSvc.Preview(c, mid.(int64), src); err != nil {
		log.Error("%+v", err)
		c.JSON(nil, err)
		return
	}
	c.Writer.Header().Set("Content-Type", http.DetectContentType(img))
	c.JSON(img, err)
}

// alipay api

func realnameChannel(c *bm.Context) {
	c.JSON(realnameSvc.Channel(c))
}

func realnameCaptcha(c *bm.Context) {
	var (
		mid, _    = c.Get("mid")
		userAgent = c.Request.UserAgent()
		ip        = metadata.String(c, metadata.RemoteIP)
		err       error
	)
	var resp struct {
		URL    string `json:"url"`
		Remote int    `json:"remote"`
	}
	resp.URL, resp.Remote, err = realnameSvc.CaptchaGTRegister(c, mid.(int64), ip, geetestClientType(userAgent))
	c.JSON(resp, err)
}

func realnameCaptchaRefresh(c *bm.Context) {
	var (
		err       error
		mid, _    = c.Get("mid")
		userAgent = c.Request.UserAgent()
		ip        = metadata.String(c, metadata.RemoteIP)
		v         = &model.ParamRealnameCaptchaGTRefresh{}
	)
	if err = c.Bind(v); err != nil {
		return
	}
	var resp struct {
		CaptureType int `json:"captcha_type"`
		CaptureInfo struct {
			Success   int    `json:"success"`
			GT        string `json:"gt"`
			Challenge string `json:"challenge"`
		} `json:"captcha_info"`
	}
	resp.CaptureType = 1
	resp.CaptureInfo.Challenge, resp.CaptureInfo.GT, resp.CaptureInfo.Success, err = realnameSvc.CaptchaGTRefresh(c, mid.(int64), ip, geetestClientType(userAgent), v.Hash)
	c.JSON(resp, err)
}

func realnameCaptchaConfirm(c *bm.Context) {
	var (
		err       error
		mid, _    = c.Get("mid")
		userAgent = c.Request.UserAgent()
		ip        = metadata.String(c, metadata.RemoteIP)
		v         = &model.ParamRealnameCaptchaGTCheck{}
	)
	if err = c.Bind(v); err != nil {
		return
	}
	if v.Remote == 1 && len(v.Challenge) != 34 {
		err = ecode.RequestErr
		return
	}
	c.JSON(realnameSvc.CaptchaGTValidate(c, mid.(int64), ip, geetestClientType(userAgent), v))
}

func realnameAlipayApply(c *bm.Context) {
	var (
		err    error
		mid, _ = c.Get("mid")
		v      = &model.ParamRealnameAlipayApply{}
	)
	if err = c.Bind(v); err != nil {
		return
	}
	c.JSON(realnameSvc.AlipayApply(c, mid.(int64), v))
}

func realnameAlipayConfirm(c *bm.Context) {
	var (
		mid, _ = c.Get("mid")
	)
	c.JSON(realnameSvc.AlipayConfirm(c, mid.(int64)))
}

func geetestClientType(userAgent string) string {
	return "h5"
}
