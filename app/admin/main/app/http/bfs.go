package http

import (
	"io/ioutil"
	"net/http"

	"go-common/app/admin/main/app/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/render"
)

func clientUpCover(c *bm.Context) {
	req := c.Request
	file, _, err := req.FormFile("cover")
	if err != nil {
		log.Error("c.Request().FormFile(\"file\") error(%v) | ", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	bs, err := ioutil.ReadAll(file)
	file.Close()
	if err != nil {
		log.Error("ioutil.ReadAll(c.Request().Body) error(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	ftype := http.DetectContentType(bs)
	if model.IsCoverType(ftype) {
		log.Error("filetype not allow file type(%s)", ftype)
		renderErrMsg(c, ecode.RequestErr.Code(), "文件上传错误：图片类型错误")
		return
	}
	url, err := bfsSvc.ClientUpCover(c, ftype, bs)
	if err != nil {
		code := ecode.RequestErr
		renderErrMsg(c, code.Code(), "文件上传错误："+code.Message())
		return
	}
	data := map[string]interface{}{
		"url": url,
	}
	c.Render(http.StatusOK, render.MapJSON(data))
}

func renderErrMsg(c *bm.Context, code int, msg string) {
	data := map[string]interface{}{
		"code":    code,
		"message": msg,
	}
	c.Render(http.StatusOK, render.MapJSON(data))
}
