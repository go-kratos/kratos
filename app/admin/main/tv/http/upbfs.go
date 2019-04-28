package http

import (
	"io/ioutil"
	"net/http"
	"time"

	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

const maxSize = 1024 * 1024 * 20

func upbfs(c *bm.Context) {
	var req = c.Request
	// read the file
	req.ParseMultipartForm(maxSize)
	log.Info("Request Info: %v, %v, %v", req.PostForm, req.Form, req.MultipartForm)
	file, _, err := req.FormFile("file")
	if err != nil {
		renderErrMsg(c, ecode.RequestErr.Code(), "文件为空")
		return
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Error("resource uploadFile.ReadAll error(%v)", err)
		return
	}
	// parse file, get type, size
	ftype := http.DetectContentType(content)
	if ftype != "image/jpeg" && ftype != "image/png" && ftype != "image/webp" && ftype != "image/gif" {
		log.Error("filetype not allow file type(%s)", ftype)
		renderErrMsg(c, ecode.RequestErr.Code(), "检查文件类型，需为图片")
		return
	}
	fsize := len(content)
	if fsize > maxSize {
		renderErrMsg(c, ecode.RequestErr.Code(), "文件过大，不支持超过20M的文件")
		return
	}
	// upload file to BFS
	c.JSON(tvSrv.Upload(c, "", ftype, time.Now().Unix(), content))
}
