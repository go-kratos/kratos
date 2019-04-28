package http

import (
	"io/ioutil"
	"net/http"
	"strconv"

	usrmdl "go-common/app/service/main/usersuit/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"

	"github.com/pkg/errors"
)

func updateFace(c *bm.Context) {
	var (
		//ip      = c.RemoteIP()
		mid, ok = c.Get("mid")
	)
	defer c.Request.Form.Del("face") // 防止日志不出现
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	c.Request.ParseMultipartForm(32 << 20)
	face, err := func() ([]byte, error) {
		fs := c.Request.FormValue("face")
		if fs != "" {
			log.Info("Succeeded to parse face file from form value: mid: %d, length: %d", mid, len(fs))
			return []byte(fs), nil
		}
		log.Warn("Failed to parse face file from form value, fallback to form file: mid: %d", mid)
		f, _, err := c.Request.FormFile("face")
		if err != nil {
			return nil, errors.Wrapf(err, "parse face form file: mid: %d", mid)
		}
		defer f.Close()
		data, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, errors.Wrapf(err, "read face form file: mid: %d", mid)
		}
		if len(data) <= 0 {
			return nil, errors.Wrapf(err, "form file data: mid: %d, length: %d", mid, len(data))
		}
		log.Info("Succeeded to parse file from form file: mid: %d, length: %d", mid, len(data))
		return data, nil
	}()
	if err != nil {
		log.Error("Failed to parse face file: mid: %d: %+v", mid, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	log.Info("Succeeded to parse face data: mid: %d, face-length: %d", mid, len(face))
	if len(face) > 2*1024*1024 {
		c.JSON(nil, ecode.UpdateFaceSize)
		return
	}
	ftype := http.DetectContentType(face)
	if ftype != "image/jpeg" && ftype != "image/png" && ftype != "image/jp2" {
		c.JSON(nil, ecode.UpdateFaceFormat)
		return
	}
	c.JSON(memberSvc.UpdateFace(c, mid.(int64), face, ftype))
}

func updateSex(c *bm.Context) {
	var (
		err error
		sex int64
		//ip      = c.RemoteIP()
		params  = c.Request.Form
		mid, ok = c.Get("mid")
	)
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	if sex, err = strconv.ParseInt(params.Get("sex"), 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, memberSvc.UpdateSex(c, mid.(int64), sex))
}

func updateSign(c *bm.Context) {
	var (
		//ip      = c.RemoteIP()
		params  = c.Request.Form
		mid, ok = c.Get("mid")
	)
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	sign := params.Get("user_sign")
	c.JSON(nil, memberSvc.UpdateSign(c, mid.(int64), sign))
}

func updateBirthday(c *bm.Context) {
	var (
		birthday string
		//ip       = c.RemoteIP()
		params  = c.Request.Form
		mid, ok = c.Get("mid")
	)
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	if birthday = params.Get("birthday"); len(birthday) == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, memberSvc.UpdateBirthday(c, mid.(int64), birthday))
}

func updateUname(c *bm.Context) {
	var (
		uname string
		//ip      = c.RemoteIP()
		params  = c.Request.Form
		mid, ok = c.Get("mid")
	)
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	if uname = params.Get("uname"); len(uname) == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, memberSvc.UpdateName(c, mid.(int64), uname, params.Get("appkey")))
}

func nickFree(c *bm.Context) {
	var (
		mid, ok = c.Get("mid")
		//ip      = c.RemoteIP()
	)
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	c.JSON(memberSvc.NickFree(c, mid.(int64)))
}

func pendantEquip(c *bm.Context) {
	var (
		mid, ok = c.Get("mid")
		params  = c.Request.Form
		err     error
		pid     int64
		status  int64
		source  int64
	)
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	if pid, err = strconv.ParseInt(params.Get("pid"), 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if status, err = strconv.ParseInt(params.Get("status"), 10, 64); err != nil || (status != usrmdl.PendantPutOn && status != usrmdl.PendantPickOff) {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// source 挂件来源 可选， 默认0， 0未知 1 背包挂件 2大会员挂件
	source = usrmdl.ParseSource(params.Get("source"))

	c.JSON(nil, usSvc.Equip(c, mid.(int64), pid, int8(status), source))
}
