package http

import (
	"encoding/json"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"go-common/app/service/openplatform/abtest/model"
	"go-common/app/service/openplatform/abtest/model/validator"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func versionID(c *bm.Context) {
	params := new(validator.VerionParams)
	if err := c.Bind(params); err != nil {
		return
	}
	c.JSON(abSvr.VersionID(c, params.Group))
}

func version(c *bm.Context) {
	var (
		key      = c.Request.Form.Get("key")
		verStr   = c.Request.Form.Get("version")
		ver      = &model.Version{}
		groupStr = c.Request.Form.Get("group")
		appKey   = c.Request.Form.Get("appkey")
		group    int
		err      error
	)
	if verStr != "" {
		if err = json.Unmarshal([]byte(verStr), ver); err != nil {
			log.Warn("[http.stra|version] json.Unmarshal(%s) err: %v", verStr, err)
		}
	}
	if groupStr != "" {
		if group, err = strconv.Atoi(groupStr); err != nil {
			log.Warn("[http.stra|version] strconv.Atoi(group:%s) err: %v", groupStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	c.JSON(abSvr.Version(c, group, key, ver, appKey))
}

func listAb(c *bm.Context) {
	var (
		params = new(validator.ListParams)
		data   []*model.AB
		total  int
		err    error
	)
	if err = c.Bind(params); err != nil {
		return
	}
	if err = checkStatus(params.Mstatus); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if data, total, err = abSvr.ListAb(c, params.Pn, params.Ps, params.Mstatus, params.Group); err != nil {
		log.Error("[http.stra|version] abSvr.ListAb(%d,%d) err: %v", params.Pn, params.Ps, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]interface{}{
		"result": data,
		"total":  total,
	}, nil)
}

func addAb(c *bm.Context) {
	var (
		params = new(validator.AddAbParams)
		ab     = &model.AB{}
		err    error
	)
	if err = c.Bind(params); err != nil {
		return
	}
	// params
	if err = json.Unmarshal([]byte(params.Data), ab); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if ab.Group == 0 {
		ab.Group = params.Group
	}
	if ab.Group == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if params.Group != 0 && ab.Group != params.Group {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if !ab.Stra.Check() {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	rand.Seed(time.Now().Unix())
	c.JSON(abSvr.AddAb(c, ab))
}

func updateAb(c *bm.Context) {
	var (
		err   error
		ab    *model.AB
		srcAb model.AB
	)
	params := new(validator.UpdateAbParams)
	if err = c.Bind(params); err != nil {
		return
	}
	if ab, err = abSvr.Ab(c, params.ID, params.Group); err != nil {
		c.JSON(nil, err)
		return
	}
	if params.Group != 0 && ab.Group != params.Group {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	srcAb = *ab
	if err = json.Unmarshal([]byte(params.Data), ab); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if !ab.Stra.Check() {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	ab.Version = srcAb.Version
	ab.Status = srcAb.Status
	ab.Group = srcAb.Group
	c.JSON(abSvr.UpdateAb(c, params.ID, ab))
}

func updateStatus(c *bm.Context) {
	var (
		ab  *model.AB
		err error
	)
	params := new(validator.UpdateStatusAbParams)
	if err = c.Bind(params); err != nil {
		return
	}
	if params.Status > 3 || params.Status < 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if ab, err = abSvr.Ab(c, params.ID, params.Group); err != nil {
		c.JSON(nil, err)
		return
	}
	if params.Group != 0 && ab.Group != params.Group {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(abSvr.UpdateStatus(c, params.ID, params.Status, params.Modifier, ab.Group))
}

func deleteAb(c *bm.Context) {
	var (
		ab  *model.AB
		err error
	)
	params := new(validator.DelAbParams)
	if err = c.Bind(params); err != nil {
		return
	}
	if ab, err = abSvr.Ab(c, params.ID, params.Group); err != nil {
		c.JSON(nil, err)
		return
	}
	if params.Group != 0 && ab.Group != params.Group {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(abSvr.DeleteAb(c, params.ID))
}

func checkStatus(s string) (err error) {
	strs := strings.Split(s, ",")
	for _, a := range strs {
		if _, err = strconv.Atoi(a); err != nil {
			return
		}
	}
	return
}

func total(c *bm.Context) {
	c.JSON(abSvr.Total(c))
}
