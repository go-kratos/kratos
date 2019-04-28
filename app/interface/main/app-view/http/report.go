package http

import (
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"

	"go-common/app/interface/main/app-view/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func copyWriter(c *bm.Context) {
	params := c.Request.Form
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	lang := params.Get("lang")
	aid, err := strconv.ParseInt(params.Get("aid"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	plat := model.Plat(mobiApp, device)
	c.JSON(reportSvr.CopyWriter(c, aid, plat, lang))
}

func addReport(c *bm.Context) {
	var mid int64
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	params := c.Request.Form
	ak := params.Get("access_key")
	reason := params.Get("reason")
	pics := params.Get("pics")
	aid, err := strconv.ParseInt(params.Get("aid"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	tp, err := strconv.Atoi(params.Get("type"))
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if tp < 0 || tp > 10 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, reportSvr.AddReport(c, mid, aid, tp, ak, reason, pics))
}

func upload(c *bm.Context) {
	var (
		fileType string
		body     []byte
		file     multipart.File
		err      error
	)
	if file, _, err = c.Request.FormFile("file"); err != nil {
		c.JSON(nil, ecode.RequestErr)
		log.Error("c.Request.FormFile(\"file\") error(%v)", err)
		return
	}
	defer file.Close()
	if body, err = ioutil.ReadAll(file); err != nil {
		c.JSON(nil, ecode.RequestErr)
		log.Error("ioutil.ReadAll(c.Request.Body) error(%v)", err)
		return
	}
	fileType = http.DetectContentType(body)
	url, err := reportSvr.Upload(c, fileType, body)
	c.JSON(struct {
		URL string `json:"url"`
	}{url}, err)
}
