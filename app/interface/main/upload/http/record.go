package http

import (
	"bytes"
	"io"
	"net/http"
	"strconv"

	"go-common/app/interface/main/upload/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
	"go-common/library/net/http/blademaster/render"
)

const (
	_defaultDistance = 1
)

// genImageUpload .
func genImageUpload(c *bm.Context) {
	params := c.Request.Form
	uploadKey := params.Get("upload_key")
	wmKey := params.Get("wm_key")
	wmText := params.Get("wm_text")
	if len(wmText) > 20 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	distance, err := strconv.Atoi(params.Get("distance"))
	if err != nil {
		distance = _defaultDistance
	}
	vertical, err := strconv.ParseBool(params.Get("wm_vertical"))
	if err != nil {
		vertical = true
	}
	c.JSON(uploadSvr.GenImageUpload(c, uploadKey, wmKey, wmText, distance, vertical))
}

// uploadImagePublic .
func uploadImagePublic(c *bm.Context) {
	if err := c.Request.ParseMultipartForm(model.MaxUploadSize); err != nil {
		c.JSON(nil, ecode.BfsUploadFileTooLarge)
		return
	}
	params := c.Request.Form
	uploadKey := params.Get("upload_key")
	uploadToken := params.Get("upload_token")
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		log.Error("upload.UploadImage.file.illegal,err:(%v)", err.Error())
		c.JSON(nil, ecode.RequestErr)
		return
	}
	defer file.Close()
	buf := new(bytes.Buffer)
	if _, err = io.Copy(buf, file); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	result, err := uploadSvr.Upload(c, uploadKey, uploadToken, header.Header.Get("Content-Type"), buf.Bytes())
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(result, nil)
}

// internalUpload upload by key and sign.
func internalUpload(c *bm.Context) {
	var (
		err error
		mid int64
	)
	up := new(model.UploadParam)
	if err = c.BindWith(up, binding.FormMultipart); err != nil {
		return
	}
	up.WMInit()
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		log.Error("upload.UploadImage.file.illegal,err::%v", err.Error())
		c.JSON(nil, ecode.FileNotExists)
		return
	}
	defer file.Close()
	buf := new(bytes.Buffer)
	if _, err = io.Copy(buf, file); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(uploadSvr.UploadRecord(c, model.UploadInternal, mid, up, buf.Bytes()))
}

// internalUploadImage .
func internalUploadImage(c *bm.Context) {
	var (
		err error
		mid int64
	)
	up := new(model.UploadParam)
	if err = c.BindWith(up, binding.FormMultipart); err != nil {
		return
	}
	up.WMInit()
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		log.Error("upload.UploadImage.file.illegal,err::%v", err.Error())
		c.JSON(nil, ecode.FileNotExists)
		return
	}
	defer file.Close()
	buf := new(bytes.Buffer)
	if _, err = io.Copy(buf, file); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	result, err := uploadSvr.UploadRecord(c, model.UploadInternal, mid, up, buf.Bytes())
	if err != nil {
		log.Error("uploadSvr.UploadRecord(%d,%+v) error(%v)", mid, up, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(result, nil)
}

// internalUploadAdminImage .
func internalUploadAdminImage(c *bm.Context) {
	var err error
	up := new(model.UploadParam)
	if err = c.BindWith(up, binding.FormMultipart); err != nil {
		return
	}
	up.WMInit()
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		log.Error("upload.UploadImage.file.illegal,err::%v", err.Error())
		c.JSON(nil, ecode.RequestErr)
		return
	}
	defer file.Close()
	buf := new(bytes.Buffer)
	if _, err = io.Copy(buf, file); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	result, err := uploadSvr.UploadAdminRecord(c, model.UploadInternalAdmin, up, buf.Bytes())
	if err != nil {
		log.Error("uploadSrv.Upload(%v,%v,%v) error(%v)", up.Bucket, up.FileName, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(result, nil)
}

// uploadMobileImage .
func uploadMobileImage(c *bm.Context) {
	var (
		err error
		mid int64
	)
	up := new(model.UploadParam)
	if err = c.BindWith(up, binding.FormMultipart); err != nil {
		return
	}
	up.WMInit()
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if mid <= 0 {
		c.Render(http.StatusOK, render.JSON{
			Code:    ecode.UserNotExist.Code(),
			Message: "invalid or not exist mid",
			Data:    nil,
		})
		return
	}
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		log.Error("upload.UploadMobileImage.file.illegal,err::%v", err.Error())
		c.JSON(nil, ecode.FileNotExists)
		return
	}
	defer file.Close()
	buf := new(bytes.Buffer)
	if _, err = io.Copy(buf, file); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	result, err := uploadSvr.UploadRecord(c, model.UploadApp, mid, up, buf.Bytes())
	if err != nil {
		log.Error("uploadSrv.UploadRecord(%v,%v,%v) error(%v)", mid, up.Bucket, up.FileName, err)
		c.JSON(nil, err)
		return
	}
	log.Info("app/upload param (%+v) result (%+v)", up, result)
	c.JSON(result, nil)
}

// uploadWebImage .
func uploadWebImage(c *bm.Context) {
	var (
		err error
		mid int64
	)
	up := new(model.UploadParam)
	if err = c.BindWith(up, binding.FormMultipart); err != nil {
		return
	}
	up.WMInit()
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if mid <= 0 {
		c.Render(http.StatusOK, render.JSON{
			Code:    ecode.UserNotExist.Code(),
			Message: "invalid or not exist mid",
			Data:    nil,
		})
		return
	}
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		log.Error("upload.UploadWebImage.file.illegal,err::%v", err.Error())
		c.JSON(nil, ecode.FileNotExists)
		return
	}
	defer file.Close()
	buf := new(bytes.Buffer)
	if _, err = io.Copy(buf, file); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	result, err := uploadSvr.UploadRecord(c, model.UploadWeb, mid, up, buf.Bytes())
	if err != nil {
		log.Error("uploadSrv.UploadRecord(%v,%v,%v) error(%v)", mid, up.Bucket, up.FileName, err)
		c.JSON(nil, err)
		return
	}
	log.Info("web/upload param (%+v) result (%+v)", up, result)
	c.JSON(result, nil)
}
