package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func upload(c *bm.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		log.Error("FormFile err(%v)", err)
		httpCode(c, fmt.Sprintf("File Upload FormFile Error:(%v)", err), ecode.RequestErr)
		return
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Error(" uploadFile.ReadAll error(%v)", err)
		httpCode(c, fmt.Sprintf("File ioutil.ReadAll Error:(%v)", err), ecode.RequestErr)
		return
	}
	filetype := http.DetectContentType(content)
	switch filetype {
	case "image/jpeg", "image/jpg":
	case "image/gif":
	case "image/png":
	case "application/zip":
	default:
		httpCode(c, fmt.Sprintf("not allow filetype(%s)", filetype), ecode.RequestErr)
		log.Warn("not allow filetype(%s) ", filetype)
		return
	}
	local, err := svc.Upload(c, "", filetype, time.Now().Unix(), content)
	if err != nil {
		log.Error("svc.Upload error(%v)", err)
		httpCode(c, fmt.Sprintf("svc.Upload error:(%v)", err), ecode.RequestErr)
		return
	}
	c.JSON(local, nil)
}
