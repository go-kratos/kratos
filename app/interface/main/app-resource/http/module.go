package http

import (
	"encoding/json"
	"strconv"
	"time"

	mdl "go-common/app/interface/main/app-resource/model/module"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func list(c *bm.Context) {
	var (
		params                            = c.Request.Form
		build, sysver, level, scale, arch int
		err                               error
		env                               string
	)
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	platform := params.Get("platform")
	env = params.Get("env")
	if env != mdl.EnvRelease && env != mdl.EnvTest && env != mdl.EnvDefault {
		env = mdl.EnvRelease
	}
	build, _ = strconv.Atoi(params.Get("build"))
	sysver, _ = strconv.Atoi(params.Get("sysver"))
	level, _ = strconv.Atoi(params.Get("level"))
	scale, _ = strconv.Atoi(params.Get("scale"))
	arch, _ = strconv.Atoi(params.Get("arch"))
	rPoolName := params.Get("resource_pool_name")
	// NOTE: don't ask way rPoolName coded "pink", idiot demand! fuck!!!
	// rPoolName := "pink"
	// params
	verlist := params.Get("verlist")
	var versions []*mdl.Versions
	if verlist != "" {
		if err = json.Unmarshal([]byte(verlist), &versions); err != nil {
			log.Error("http list() json.Unmarshal(%s) mobile_app(%s) device(%s) build(%d) error(%v)", verlist, mobiApp, device, build, err)
		}
	}
	data := moduleSvc.List(c, mobiApp, device, platform, rPoolName, env, build, sysver, level, scale, arch, versions, time.Now())
	c.JSON(data, nil)
}

func module(c *bm.Context) {
	var (
		params                            = c.Request.Form
		build, sysver, level, scale, arch int
		err                               error
		env                               string
	)
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	platform := params.Get("platform")
	env = params.Get("env")
	if env != mdl.EnvRelease && env != mdl.EnvTest && env != mdl.EnvDefault {
		env = mdl.EnvRelease
	}
	build, _ = strconv.Atoi(params.Get("build"))
	sysver, _ = strconv.Atoi(params.Get("sysver"))
	level, _ = strconv.Atoi(params.Get("level"))
	scale, _ = strconv.Atoi(params.Get("scale"))
	arch, _ = strconv.Atoi(params.Get("arch"))
	rPoolName := params.Get("resource_pool_name")
	// NOTE: don't ask way rPoolName coded "pink", idiot demand! fuck!!!
	// rPoolName := "pink"
	rName := params.Get("resource_name")
	verStr := params.Get("ver")
	ver, _ := strconv.Atoi(verStr)
	if rPoolName == "" || rName == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	data, err := moduleSvc.Resource(c, mobiApp, device, platform, rPoolName, rName, env, ver, build, sysver, level, scale, arch, time.Now())
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}
