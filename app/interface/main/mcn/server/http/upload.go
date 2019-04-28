package http

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"path"
	"time"

	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

const (
	maxFileSize = 20 << 20
)

func upload(c *bm.Context) {
	c.Request.ParseMultipartForm(maxFileSize)
	imageFile, header, err := c.Request.FormFile("file")
	if err != nil {
		log.Error("upload fail, no field file")
		c.JSON(nil, err)
		return
	}
	if header.Size > maxFileSize {
		c.JSON(nil, ecode.MCNContractFileSize)
		return
	}
	fileExt := path.Ext(header.Filename)
	if fileExt == "" {
		c.JSON(nil, ecode.MCNUnknownFileExt)
		return
	}
	defer imageFile.Close()
	bs, err := ioutil.ReadAll(imageFile)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	filetype := http.DetectContentType(bs)
	switch filetype {
	case
		"image/jpeg",
		"image/jpg",
		"image/png",
		"application/pdf":
	case "application/octet-stream", "application/zip":
		switch fileExt[1:] {
		case "doc":
			filetype = "application/doc"
		case "docx":
			filetype = "application/docx"
		case "docm":
			filetype = "application/docm"
		}
	case "application/msword":
		filetype = "application/doc"
	case "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
		filetype = "application/docx"
	case "application/vnd.ms-word.document.macroEnabled.12":
		filetype = "application/docm"
	default:
		c.JSON(nil, ecode.MCNUnknownFileTypeErr)
		return
	}
	c.JSON(srv.Upload(c, "", filetype, time.Now().Unix(), bytes.NewReader(bs)))
}
