package http

import (
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"time"

	"go-common/app/admin/main/apm/model/ut"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func upload(c *bm.Context) {
	var (
		file      multipart.File
		body      []byte
		htmlURL   string
		reportURL string
		dataURL   string
		err       error
		pkg       = new(ut.PkgAnls)
		files     []*ut.File
		res       = &ut.UploadRes{}
		header    *multipart.FileHeader
	)
	c.Request.ParseMultipartForm(32 << 20)
	if err = c.Bind(res); err != nil {
		return
	}
	log.Info("ut.upload(%d) start! current_time(%d)", res.MergeID, time.Now().Unix())
	defer log.Info("ut.upload(%d) finished. current_time(%d)", res.MergeID, time.Now().Unix())
	if file, _, err = c.Request.FormFile("report_file"); err != nil {
		c.JSON(nil, ecode.RequestErr)
		log.Error("report request upload err (%v)", err)
		return
	}
	defer file.Close()
	if body, err = ioutil.ReadAll(file); err != nil {
		c.JSON(nil, ecode.RequestErr)
		log.Error("ioutil.ReadAll(c.Request().Body) error(%v)", err)
		return
	}
	if body, err = apmSvc.ParseContent(c, body); err != nil {
		c.JSON(err.Error(), err)
		return
	}
	if pkg, err = apmSvc.CalcCount(c, body); err != nil {
		c.JSON(nil, err)
		return
	}
	if pkg.Assertions == 0 {
		c.JSON("no result", nil)
		return
	}
	if reportURL, err = apmSvc.Upload(c, "json", time.Now().Unix(), body); err != nil {
		c.JSON(nil, err)
		return
	}
	if file, header, err = c.Request.FormFile("data_file"); err == nil && header.Size > 0 {
		defer file.Close()
		if body, err = ioutil.ReadAll(file); err != nil {
			c.JSON(nil, ecode.RequestErr)
			log.Error("Upload data request error(%v)", err)
			return
		}
		if files, err = apmSvc.CalcCountFiles(c, res, body); err != nil {
			c.JSON(nil, err)
			log.Error("Upload data calcCount error(%v)", err)
			return
		}
		if dataURL, err = apmSvc.Upload(c, "text/plain", time.Now().Unix(), body); err != nil {
			c.JSON(nil, err)
			return
		}
	}
	if file, _, err = c.Request.FormFile("html_file"); err != nil {
		c.JSON(nil, ecode.RequestErr)
		log.Error("html request upload err (%v)", err)
		return
	}
	defer file.Close()
	if body, err = ioutil.ReadAll(file); err != nil {
		c.JSON(nil, ecode.RequestErr)
		log.Error("html read error(%v)", err)
		return
	}
	if htmlURL, err = apmSvc.Upload(c, "html", time.Now().Unix(), body); err != nil {
		c.JSON(nil, err)
		return
	}
	if err = apmSvc.AddUT(c, pkg, files, res, dataURL, reportURL, htmlURL); err != nil {
		c.JSON(nil, err)
		return
	}
	// update ut_app has_ut = 1 && converage
	if err = apmSvc.UpdateUTApp(c, pkg); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

// uploadApp upload path to ut_app
func uploadApp(c *bm.Context) {
	var (
		file multipart.File
		body []byte
		apps []*ut.App
		err  error
	)
	c.Request.ParseMultipartForm(32 << 20)
	if file, _, err = c.Request.FormFile("path_file"); err != nil {
		c.JSON(nil, ecode.RequestErr)
		log.Error("path request upload error(%v)", err)
		return
	}
	defer file.Close()
	if body, err = ioutil.ReadAll(file); err != nil {
		c.JSON(nil, ecode.RequestErr)
		log.Error("path_file read error(%v)", err)
		return
	}
	if err = json.Unmarshal(body, &apps); err != nil {
		log.Error("json.Unmarshal error(%v)", err)
		return
	}
	if err = apmSvc.AddUTApp(c, apps); err != nil {
		c.JSON(nil, ecode.RequestErr)
		log.Error("apmSvc.AddUtApp error(%v)", err)
		return
	}
	c.JSON(nil, err)
}
