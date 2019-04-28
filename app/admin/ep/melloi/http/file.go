package http

import (
	"bytes"
	"io"
	"mime/multipart"
	"path"
	"strconv"
	"time"

	"go-common/app/admin/ep/melloi/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)

const (
	_upLoadSuccess = "上传成功"
	_upLoadFail    = "上传失败"
)

func upload(c *bm.Context) {
	uploadParam := model.UploadParam{}
	if err := c.BindWith(&uploadParam, binding.Form); err != nil {
		c.JSON(nil, err)
		return
	}

	formFile, header, err := c.Request.FormFile("file")
	if err != nil {
		log.Error("Get form file failed,error(%v)", err)
		c.JSON(nil, err)
		return
	}
	var upLoadResultMap = make(map[string]string)
	upLoadResultMap["result"] = _upLoadSuccess
	scriptID, scriptPath, err := srv.Upload(c, &uploadParam, formFile, header)
	if err != nil {
		log.Error("Write file failed, error(%v)", err)
		upLoadResultMap["result"] = _upLoadFail
		c.JSON(upLoadResultMap, err)
		return
	}
	upLoadResultMap["script_id"] = strconv.Itoa(scriptID)
	upLoadResultMap["script_path"] = scriptPath
	upLoadResultMap["full_path"] = path.Join(scriptPath, header.Filename)
	c.JSON(upLoadResultMap, nil)
}

func uploadDependProto(c *bm.Context) {
	var (
		upParam  = &model.UploadParam{}
		header   *multipart.FileHeader
		formFile multipart.File
		err      error
	)

	if err = c.BindWith(upParam, binding.Form); err != nil {
		return
	}
	if formFile, header, err = c.Request.FormFile("file"); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(srv.UploadAndGetProtoInfo(c, upParam, formFile, header))
}

func compileProtoFile(c *bm.Context) {
	uploadParam := model.UploadParam{}
	if err := c.BindWith(&uploadParam, binding.JSON); err != nil {
		c.JSON(nil, err)
		return
	}
	fileName := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	c.JSON(srv.CompileProtoFile(uploadParam.ScriptPath, fileName))
}

func uploadImg(c *bm.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		log.Error("get IMG file error (%v)", err)
		c.JSON(nil, err)
		return
	}
	defer file.Close()
	buff := new(bytes.Buffer)
	if _, err = io.Copy(buff, file); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(srv.UploadImg(c, buff.Bytes(), header.Filename))
}

func readFile(c *bm.Context) {
	file := c.Request.Form.Get("file")
	limit := c.Request.Form.Get("limit")
	var upLoadResultMap = make(map[string]string)
	data, err := srv.ReadFile(c, file, limit)
	if err != nil {
		return
	}
	upLoadResultMap["data"] = data
	c.JSON(upLoadResultMap, nil)
}

func downloadFile(c *bm.Context) {
	filePath := c.Request.Form.Get("file_path")
	c.JSON(nil, srv.DownloadFile(c, filePath, c.Writer))
}

func isFileExists(c *bm.Context) {
	fileName := c.Request.Form.Get("file_name")
	exists, err := srv.IsFileExists(c, fileName)
	var ResultMap = make(map[string]string)
	ResultMap["fileExists"] = strconv.FormatBool(exists)
	c.JSON(ResultMap, err)
}
