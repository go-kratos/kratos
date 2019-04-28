package http

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"

	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func upload(c *bm.Context) {
	imageFile, _, err := c.Request.FormFile("file")
	if err != nil {
		log.Error("upload err(%v)", err)
		httpCode(c, err)
		return
	}
	defer imageFile.Close()
	bs, err := ioutil.ReadAll(imageFile)
	if err != nil {
		log.Error("ioutil.ReadAll err(%v)", err)
		httpCode(c, err)
		return
	}
	filetype := http.DetectContentType(bs)
	// var extension string
	switch filetype {
	case "image/jpeg", "image/jpg", "image/gif", "image/png", "application/pdf":
	default:
		log.Warn("unknown filetype(%s) ", filetype)
		return
	}
	//重新格式化文件名
	local, err := svc.Upload(c, "", filetype, time.Now().Unix(), bytes.NewReader(bs))
	if err != nil {
		log.Error("svc.Upload error(%v)", err)
		httpCode(c, err)
		return
	}
	httpData(c, local, nil)
}
