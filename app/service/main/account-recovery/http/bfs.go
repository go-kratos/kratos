package http

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"go-common/app/service/main/account-recovery/model"
	"go-common/library/database/bfs"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var _allowedExt = map[string]struct{}{
	".zip": {},
	".rar": {},
	".7z":  {},
}

func fileUpload(c *bm.Context) {
	defer c.Request.Form.Del("file") // 防止日志不出现
	c.Request.ParseMultipartForm(32 << 20)
	recoveryFile, fileName, err := func() ([]byte, string, error) {
		f, fh, err := c.Request.FormFile("file")
		if err != nil {
			return nil, "", errors.Wrapf(err, "parse file form file: ")
		}
		defer f.Close()
		data, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, "", errors.Wrapf(err, "read file form file:")
		}
		if len(data) <= 0 {
			return nil, "", errors.Wrapf(err, "form file data: length: %d", len(data))
		}
		log.Info("Succeeded to parse file from form file: length: %d", len(data))
		return data, fh.Filename, nil
	}()
	if err != nil {
		log.Error("Failed to parse file file: %+v", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	log.Info("Succeeded to parse recoveryFile data: recoveryFile-length: %d", len(recoveryFile))

	//限制文件大小
	if len(recoveryFile) > 10*1024*1024 {
		log.Error("account-recovery: file is to large(%v)", len(recoveryFile))
		c.JSON(nil, ecode.FileTooLarge)
		return
	}

	//限制文件类型 *.zip, *.rar, *.7z
	ext := filepath.Ext(fileName)
	_, allowed := _allowedExt[ext]
	if !allowed {
		c.JSON(nil, ecode.BfsUploadFileContentTypeIllegal)
		return
	}

	request := &bfs.Request{
		Bucket:   "account",
		Dir:      "recovery",
		File:     recoveryFile,
		Filename: fmt.Sprintf("%s%s", uuid4(), ext),
	}
	bfsClient := bfs.New(nil)
	location, err := bfsClient.Upload(c, request)
	if err != nil {
		log.Error("err(%+v)", err)
		c.JSON(nil, err)
		return
	}

	fileURL := model.BuildFileURL(location)
	data := map[string]interface{}{
		"url":     location,
		"fileURL": fileURL,
	}
	c.JSON(data, nil)
}

func uuid4() string {
	return uuid.New().String()
}
