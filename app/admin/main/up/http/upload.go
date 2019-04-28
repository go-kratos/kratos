package http

import (
	"io/ioutil"
	"mime/multipart"
	"path"
	"strings"

	"go-common/app/admin/main/up/conf"
	"go-common/app/admin/main/up/service"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/http/blademaster"
	"regexp"
	"time"
)

//由数字、26个英文字母或者下划线组成的字符串
var reg = regexp.MustCompile(`^\w+$`)

// upload
func upload(c *blademaster.Context) {
	var (
		fileTpye string
		file     multipart.File
		header   *multipart.FileHeader
		fileName string
		body     []byte
		location string
		err      error
		res      interface{}
		errMsg   string
	)

	switch {
	default:
		if file, header, err = c.Request.FormFile("file"); err != nil {
			errMsg = err.Error()
			err = ecode.RequestErr
			log.Error("c.Request().FormFile(\"file\") error(%v)", err)
			break
		}
		defer file.Close()
		fileName = header.Filename
		fileTpye = strings.TrimPrefix(path.Ext(fileName), ".")
		if body, err = ioutil.ReadAll(file); err != nil {
			errMsg = err.Error()
			err = ecode.RequestErr
			log.Error("ioutil.ReadAll(c.Request().Body) error(%v)", err)
			break
		}
		// 如果不符合规则，就不用文件名
		if !reg.MatchString(fileName) {
			fileName = ""
		}
		if location, err = Svc.Upload(c, fileName, fileTpye, time.Now(), body, conf.Conf.BfsConf); err != nil {
			errMsg = err.Error()
			break
		}

		res = struct {
			URL string `json:"url"`
		}{
			location,
		}
	}

	if err != nil {
		service.BmHTTPErrorWithMsg(c, err, errMsg)
	} else {
		c.JSON(res, err)
	}
}
