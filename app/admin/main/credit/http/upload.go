package http

import (
	"context"
	"crypto/md5"
	"encoding/csv"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func upload(c *bm.Context) {
	imageFile, imageHeader, err := c.Request.FormFile("file")
	if err != nil {
		log.Error("upload err(%v)", err)
		httpCode(c, err)
		return
	}
	defer imageFile.Close()
	//读取前512个字节用于判断文件类型
	firstImageBytes := make([]byte, 512)
	_, err = imageFile.Read(firstImageBytes)
	if err != nil {
		log.Error("imageFile.Read error(%v)", err)
		httpCode(c, err)
		return
	}
	md5Checksum := md5.Sum(firstImageBytes)
	extensionMatcher := regexp.MustCompile(`\\.\\w+$`)
	imageName := extensionMatcher.ReplaceAllString(imageHeader.Filename, "")
	filetype := http.DetectContentType(firstImageBytes)
	var extension string
	switch filetype {
	case "image/jpeg", "image/jpg":
		extension = "jpg"
	case "image/gif":
		extension = "gif"
	case "image/png":
		extension = "png"
	case "application/pdf":
		extension = "pdf"
	default:
		log.Warn("unknown filetype(%s) ", filetype)
		return
	}
	imageName = url.PathEscape(imageName)
	//重新格式化文件名
	uploadFilePath := fmt.Sprintf("%x-%v.%v", md5Checksum, imageName, extension)
	local, err := creSvc.Upload(c, uploadFilePath, extension, time.Now().Unix(), imageFile)
	if err != nil {
		log.Error("creSvc.Upload error(%v)", err)
		httpCode(c, err)
		return
	}
	httpData(c, local, nil)
}

func annualCoins(c *bm.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		log.Error("upload err(%v)", err)
		httpCode(c, err)
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
	fmids := creSvc.AnnualCoins(context.Background(), reader)
	httpData(c, fmids, nil)
}
