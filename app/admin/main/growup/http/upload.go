package http

import (
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"path"
	"strings"
	"time"

	"bytes"
	"crypto/md5"
	"fmt"
	"go-common/app/admin/main/growup/conf"
	"go-common/app/admin/main/growup/dao"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/render"
	"io"
	"os"
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
	if file, header, err = c.Request.FormFile("file"); err != nil {
		c.JSON(nil, ecode.RequestErr)
		log.Error("c.Request().FormFile(\"file\") error(%v)", err)
		return
	}
	defer file.Close()
	fileName = header.Filename
	fileTpye = strings.TrimPrefix(path.Ext(fileName), ".")
	if body, err = ioutil.ReadAll(file); err != nil {
		c.JSON(nil, ecode.RequestErr)
		log.Error("ioutil.ReadAll(c.Request().Body) error(%v)", err)
		return
	}
	if location, err = svr.Upload(c, "", fileTpye, time.Now(), body); err != nil {
		c.JSON(nil, err)
		return
	}

	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code":    0,
		"message": "0",
		"data": struct {
			URL string `json:"url"`
		}{
			location,
		},
	}))
}

const (
	uploadFilePath = "/data/uploadfiles/"
)

type localUploadResult struct {
	FileName string `json:"filename"`
}
type uploadResult struct {
	localUploadResult
	MemberCount int64 `json:"member_count"`
}

// upload
func uploadLocal(c *bm.Context) {
	var (
		file   multipart.File
		header *multipart.FileHeader
		result uploadResult
		err    error
	)
	switch {
	default:
		if file, header, err = c.Request.FormFile("file"); err != nil {
			log.Error("c.Request().FormFile(\"file\") error(%v)", err)
			break
		}
		if header.Size >= int64(conf.Conf.Bfs.MaxFileSize) {
			log.Error("file is too big, filesize=%d, expected<=%d", conf.Conf.Bfs.MaxFileSize, header.Size)
			err = fmt.Errorf("文件过大，需要小于%dkb", conf.Conf.Bfs.MaxFileSize/1024)
			break
		}
		var ext = path.Ext(header.Filename)
		if ext != ".csv" && ext != ".txt" {
			log.Error("only csv or txt supported")
			err = fmt.Errorf("只支持csv或txt格式文件，以逗号或换行分隔")
			break
		}

		defer file.Close()
		// 创建保存文件
		if _, e := os.Stat(uploadFilePath); os.IsNotExist(e) {
			os.MkdirAll(uploadFilePath, 0777)
		}
		var h = md5.New()
		h.Write([]byte(header.Filename))
		var filenameMd5 = fmt.Sprintf("%x", h.Sum(nil))
		var desfilename = filenameMd5 + ext
		var destFile *os.File
		destFile, err = os.Create(uploadFilePath + desfilename)
		if err != nil {
			log.Error("Create failed: %s\n", err)
			break
		}

		defer destFile.Close()
		var membuf = bytes.NewBuffer(nil)
		_, err = io.Copy(membuf, file)
		if err != nil {
			log.Error("err copy file, err=%s", err)
			break
		}
		var midList = dao.ParseMidsFromString(membuf.String())
		result.MemberCount = int64(len(midList))

		_, err = destFile.Write(membuf.Bytes())
		if err != nil {
			log.Error("write file error, err=%s", err)
			break
		}
		result.FileName = desfilename
	}

	if err != nil {
		bmHTTPErrorWithMsg(c, err, err.Error())
	} else {
		c.JSON(&result, err)
	}
}
