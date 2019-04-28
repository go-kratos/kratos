package http

import (
	"io/ioutil"
	"net/http"

	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/render"
)

func clientUpload(c *bm.Context) {
	var (
		req = c.Request
		md5 string
		url string
	)
	req.ParseMultipartForm(int64(bfsSvc.BfsMaxSize))
	file, _, err := req.FormFile("file")
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
	if md5, err = bfsSvc.FileMd5(bs); err != nil {
		log.Error("bfsSvc.FileMd5 error(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	ftype := http.DetectContentType(bs)
	//if model.IsCoverType(ftype) {
	//	log.Error("filetype not allow file type(%s)", ftype)
	//	renderErrMsg(c, ecode.RequestErr.Code(), "文件上传错误：图片类型错误")
	//	return
	//}
	if url, err = bfsSvc.ClientUpCover(c, ftype, bs); err != nil {
		log.Error("bfsSvc.ClientUpCover error(%v)", err)
		c.JSON("文件上传错误："+err.Error(), ecode.RequestErr)
		return
	}
	data := map[string]interface{}{
		"url":  url,
		"md5":  md5,
		"size": len(bs),
	}
	c.Render(http.StatusOK, render.MapJSON(data))
}
