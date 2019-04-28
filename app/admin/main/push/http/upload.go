package http

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/main/push/conf"
	"go-common/app/admin/main/push/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func upload(c *bm.Context) {
	var (
		err error
		req = c.Request
	)
	req.ParseMultipartForm(1024 * 1024 * 1024) // 1G
	fileName := req.FormValue("filename")
	if fileName == "" {
		log.Error("filename is empty")
		c.JSON(nil, ecode.RequestErr)
		return
	}
	file, header, err := req.FormFile("file")
	if err != nil {
		log.Error("req.FormFile() error(%v)", err)
		c.JSON(nil, err)
		return
	}
	defer file.Close()
	bs, err := ioutil.ReadAll(file)
	if err != nil {
		log.Error("ioutil.ReadAll() error(%v)", err)
		c.JSON(nil, err)
		return
	}
	typ, _ := strconv.Atoi(req.FormValue("type"))
	if typ == model.UploadTypeMid {
		if err = pushSrv.CheckUploadMid(c, bs); err != nil {
			c.JSON(nil, err)
			return
		}
	} else if typ == model.UploadTypeToken {
		if err = pushSrv.CheckUploadToken(c, bs); err != nil {
			c.JSON(nil, err)
			return
		}
	} else {
		log.Error("type(%d) invalid", typ)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	dir := fmt.Sprintf("%s/%s", strings.TrimSuffix(conf.Conf.Cfg.MountDir, "/"), time.Now().Format("20060102"))
	path := fmt.Sprintf("%s/%x", dir, md5.Sum([]byte(fileName)))
	if err = pushSrv.Upload(c, dir, path, bs); err != nil {
		log.Error("upload file file(%s) error(%v)", path, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(struct {
		Name string `json:"name"`
		Path string `json:"path"`
	}{
		Name: header.Filename,
		Path: path,
	}, nil)
}

func upimg(ctx *bm.Context) {
	f, h, err := ctx.Request.FormFile("file")
	if err != nil {
		log.Error("upimg error(%v)", err)
		ctx.JSON(nil, err)
		return
	}
	defer f.Close()
	url, err := pushSrv.Upimg(ctx, f, h)
	if err != nil {
		ctx.JSON(nil, err)
		return
	}
	ctx.JSON(map[string]string{"url": url}, nil)
}
