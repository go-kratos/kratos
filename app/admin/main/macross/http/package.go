package http

import (
	"fmt"
	"go-common/app/admin/main/macross/conf"
	"go-common/app/admin/main/macross/model/package"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"path/filepath"
	"strings"
)

func packageUpload(c *bm.Context) {
	var err = c.Request.ParseMultipartForm(1 << 30)
	res := map[string]interface{}{}
	res["message"] = "success"

	if err != nil {
		log.Error("c.Request.ParseMultipartForm() error(%v)", err)
		res["message"] = fmt.Sprintf("%v", err)
		c.JSONMap(res, err)
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		log.Error("c.Request.FormFile() error(%v)", err)
		res["message"] = fmt.Sprintf("%v", err)
		c.JSONMap(res, err)
		return
	}
	defer file.Close()

	var (
		clientType = strings.ToLower(c.Request.FormValue("client_type"))
		appName    = c.Request.FormValue("app_name")
		pipelineID = c.Request.FormValue("pipeline_id")
		apkName    = c.Request.FormValue("apk_name")
		channel    = c.Request.FormValue("channel")
		saveDir    = filepath.Join(clientType, appName, pipelineID)
		pkgInfo    upload.PkgInfo
	)
	if clientType == "" || appName == "" || pipelineID == "" {
		errMsg := "client_type, app_name, pipeline_id can not be null"
		log.Error(errMsg)
		res["message"] = errMsg
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	if clientType != "ios" && clientType != "android" {
		errMsg := "client_type must be 'ios' or 'android'"
		log.Error(errMsg)
		res["message"] = errMsg
		c.JSONMap(res, ecode.RequestErr)
		return
	}

	pkgInfo.FileName = header.Filename
	pkgInfo.SaveDir = filepath.Join(conf.Conf.Property.Package.SavePath, saveDir)
	pkgInfo.ClientType = clientType
	pkgInfo.Channel = channel
	pkgInfo.ApkName = apkName

	err = svr.PackageUpload(file, pkgInfo)
	if err != nil {
		res["message"] = fmt.Sprintf("%v", err)
		c.JSONMap(res, err)
		return
	}
	c.JSONMap(res, nil)
}

func packageList(c *bm.Context) {
	var (
		clientType = strings.ToLower(c.Request.FormValue("client_type"))
		appName    = c.Request.Form.Get("app_name")
		pipelineID = c.Request.Form.Get("pipeline_id")
		saveDir    = filepath.Join(clientType, appName, pipelineID)
	)
	res := map[string]interface{}{}
	res["message"] = "success"

	saveDir = filepath.Join(conf.Conf.Property.Package.SavePath, saveDir)
	if clientType == "" || appName == "" || pipelineID == "" {
		errMsg := "client_type, app_name, pipeline_id can not be null"
		log.Error(errMsg)
		res["message"] = errMsg
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	if clientType != "ios" && clientType != "android" {
		errMsg := "client_type must be 'ios' or 'android'"
		log.Error(errMsg)
		res["message"] = errMsg
		c.JSONMap(res, ecode.RequestErr)
		return
	}

	fileList, err := svr.PackageList(saveDir)
	if err != nil {
		res["message"] = fmt.Sprintf("%v", err)
		c.JSONMap(res, err)
		return
	}

	c.JSON(fileList, nil)
}
