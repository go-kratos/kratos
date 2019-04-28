package http

import (
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"path"
	"strings"
	"time"

	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

// upload
func upload(c *bm.Context) {
	var (
		fileTpye string
		file     multipart.File
		header   *multipart.FileHeader
		fileName string
		body     []byte
		location string
		err      error
	)
	// res := c.Result()
	if file, header, err = c.Request.FormFile("file"); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	defer file.Close()
	fileName = header.Filename
	fileTpye = strings.TrimPrefix(path.Ext(fileName), ".")
	if body, err = ioutil.ReadAll(file); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	location, err = feedbackSvr.Upload(c, "", fileTpye, time.Now(), body)
	c.JSON(struct {
		URL string `json:"url"`
	}{location}, err)
}

// uploadFile
func uploadFile(c *bm.Context) {
	var (
		req                = c.Request
		fileTpye, location string
		body               []byte
		err                error
	)
	if body, err = ioutil.ReadAll(req.Body); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	req.Body.Close()
	fileTpye = http.DetectContentType(body)
	location, err = feedbackSvr.Upload(c, "", fileTpye, time.Now(), body)
	c.JSON(struct {
		URL string `json:"url"`
	}{location}, err)
}
