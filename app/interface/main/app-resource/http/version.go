package http

import (
	"strconv"

	"go-common/app/interface/main/app-resource/model"
	"go-common/app/interface/main/app-resource/model/version"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"

	"github.com/golang/protobuf/proto"
)

// getVersion get version
func getVersion(c *bm.Context) {
	params := c.Request.Form
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	plat := model.Plat(mobiApp, device)
	data, err := verSvc.Version(plat)
	c.JSON(data, err)
}

// versionUpdate get versionUpdate
func versionUpdate(c *bm.Context) {
	var (
		params = c.Request.Form
		header = c.Request.Header
	)
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	plat := model.Plat(mobiApp, device)
	buildStr := params.Get("build")
	channel := params.Get("channel")
	sdkStr := params.Get("sdkint")
	platModel := params.Get("model")
	oldID := params.Get("old_id")
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	buvid := header.Get(_headerBuvid)
	// mobiApp not equal to android or mobiApp is null,default android
	if !model.IsAndroid(plat) {
		plat = model.PlatAndroid
	}
	if (plat == model.PlatAndroid) && (build >= 591000 && build <= 599000) {
		plat = model.PlatAndroidB
	}
	data, err := verSvc.VersionUpdate(build, plat, buvid, sdkStr, channel, platModel, oldID)
	if err != nil {
		c.JSON(nil, ecode.NotModified)
		return
	}
	c.JSON(data, nil)
}

// versionUpdate get versionUpdate
func versionUpdatePb(c *bm.Context) {
	var (
		params = c.Request.Form
		header = c.Request.Header
	)
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	plat := model.Plat(mobiApp, device)
	buildStr := params.Get("build")
	channel := params.Get("channel")
	sdkStr := params.Get("sdkint")
	platModel := params.Get("model")
	oldID := params.Get("old_id")
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	buvid := header.Get(_headerBuvid)
	// mobiApp not equal to android or mobiApp is null,default android
	if !model.IsAndroid(plat) {
		plat = model.PlatAndroid
	}
	data, err := verSvc.VersionUpdate(build, plat, buvid, sdkStr, channel, platModel, oldID)
	if err != nil {
		c.JSON(nil, ecode.NotModified)
		return
	}
	size, _ := strconv.Atoi(data.Size)
	resPb := &version.VerUpdate{
		Ver:     *proto.String(data.Version),
		Build:   *proto.Int(data.Build),
		Info:    *proto.String(data.Desc),
		Size:    *proto.Int(size),
		Url:     *proto.String(data.Url),
		Hash:    *proto.String(data.MD5),
		Policy:  *proto.Int(data.Policy),
		IsForce: *proto.Int(data.IsForce),
		Mtime:   *proto.Int64(data.Mtime.Time().Unix()),
	}
	c.JSON(resPb, nil)
}

func versionSo(c *bm.Context) {
	params := c.Request.Form
	name := params.Get("name")
	seedStr := params.Get("seed")
	buildStr := params.Get("build")
	sdkStr := params.Get("sdkint")
	model := params.Get("model")
	seed, err := strconv.Atoi(seedStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	sdkint, _ := strconv.Atoi(sdkStr)
	data, err := verSvc.VersionSo(build, seed, sdkint, name, model)
	if err != nil {
		c.JSON(nil, ecode.NotModified)
		return
	}
	c.JSON(data, nil)
}

// versionRn get versionUpdate
func versionRn(c *bm.Context) {
	params := c.Request.Form
	deploymentKey := params.Get("deployment_key")
	bundleID := params.Get("bundle_id")
	version := params.Get("base_version")
	data, err := verSvc.VersionRn(version, deploymentKey, bundleID)
	c.JSON(data, err)
}
