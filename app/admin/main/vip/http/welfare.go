package http

import (
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"go-common/app/admin/main/vip/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

const (
	_maxFileSize = 1048576
)

func welfareTypeSave(c *bm.Context) {
	arg := new(struct {
		ID   int    `form:"id"`
		Name string `form:"name" validate:"required"`
	})
	if err := c.Bind(arg); err != nil {
		return
	}
	username, _ := c.Get("username")
	c.JSON(nil, vipSvc.WelfareTypeSave(arg.ID, arg.Name, username.(string)))
}

func welfareTypeState(c *bm.Context) {
	arg := new(struct {
		ID int `form:"id" validate:"required"`
	})
	if err := c.Bind(arg); err != nil {
		return
	}
	username, _ := c.Get("username")
	c.JSON(nil, vipSvc.WelfareTypeState(c, arg.ID, username.(string)))
}

func welfareTypeList(c *bm.Context) {
	c.JSON(vipSvc.WelfareTypeList())
}

func welfareSave(c *bm.Context) {
	const _redirectThridPage = 3

	wfReq := new(model.WelfareReq)
	if err := c.Bind(wfReq); err != nil {
		return
	}
	if wfReq.UsageForm == _redirectThridPage && wfReq.ReceiveUri == "" {
		c.JSON(nil, ecode.VipWelfareRequestErr)
		return
	}
	username, _ := c.Get("username")
	c.JSON(nil, vipSvc.WelfareSave(username.(string), wfReq))
}

func welfareState(c *bm.Context) {
	arg := new(struct {
		ID int `form:"id" validate:"required"`
	})
	if err := c.Bind(arg); err != nil {
		return
	}
	username, _ := c.Get("username")
	c.JSON(nil, vipSvc.WelfareState(arg.ID, username.(string)))
}

func welfareList(c *bm.Context) {
	arg := new(struct {
		TID int `form:"tid"`
	})
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(vipSvc.WelfareList(arg.TID))
}

func welfareBatchUpload(c *bm.Context) {
	arg := new(struct {
		WID      int    `form:"wid"`
		Filename string `form:"filename"`
		Vtime    int    `form:"vtime"`
	})
	if err := c.Bind(arg); err != nil {
		return
	}
	fileBody, _, err := getFileBody(c, "file")
	if err != nil {
		c.JSON(nil, err)
		return
	}
	username, _ := c.Get("username")
	c.JSON(nil, vipSvc.WelfareBatchUpload(fileBody, arg.Filename, username.(string), arg.WID, arg.Vtime))
}

func welfareBatchList(c *bm.Context) {
	arg := new(struct {
		WID int `form:"wid" validate:"required"`
	})
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(vipSvc.WelfareBatchList(arg.WID))
}

func welfareBatchState(c *bm.Context) {
	arg := new(struct {
		ID int `form:"id" validate:"required"`
	})
	if err := c.Bind(arg); err != nil {
		return
	}
	username, _ := c.Get("username")
	c.JSON(nil, vipSvc.WelfareBatchState(c, arg.ID, username.(string)))
}

func getFileBody(c *bm.Context, name string) (body []byte, filetype string, err error) {
	var file multipart.File
	if file, _, err = c.Request.FormFile(name); err != nil {
		if err == http.ErrMissingFile {
			err = nil
			return
		}
		err = ecode.RequestErr
		return
	}
	if file == nil {
		return
	}
	defer file.Close()
	if body, err = ioutil.ReadAll(file); err != nil {
		err = ecode.RequestErr
		return
	}
	filetype = http.DetectContentType(body)
	if filetype != "text/plain; charset=utf-8" {
		err = ecode.VipFileTypeErr
		return
	}
	if len(body) > _maxFileSize {
		err = ecode.FileTooLarge
		return
	}
	return
}
