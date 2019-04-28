package http

import (
	"go-common/app/interface/main/tv/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var (
	getForm = map[string]bool{
		"mobi_app": true,
		"build":    true,
		"channel":  true,
		"seed":     true,
		"sdkint":   false,
		"model":    false,
		"old_id":   false,
	}
)

func validatePostData(c *bm.Context, params map[string]bool) bool {
	for k, v := range params {
		if v {
			if c.Request.Form.Get(k) == "" {
				log.Error("The field is required(%s)", k)
				return false
			}
		}
	}
	return true
}

func upgrade(c *bm.Context) {
	var (
		req = c.Request.Form
		ver = &model.VerUpdate{
			MobiApp: req.Get("mobi_app"),
			Build:   int(parseInt(req.Get("build"))),
			Channel: req.Get("channel"),
			Seed:    int(parseInt(req.Get("seed"))),
			Sdkint:  int(parseInt(req.Get("sdkint"))),
			Model:   req.Get("model"),
			OldID:   req.Get("old_id"),
		}
	)
	if !validatePostData(c, getForm) {
		c.JSONMap(map[string]interface{}{
			"code":    ecode.RequestErr,
			"message": "lack of fields",
		}, nil)
		return
	}
	result, errCode, err := gobSvc.VerUpdate(c, ver)
	if err != nil {
		log.Error("[VerUpdate] Load App Upgrade Data, Err: %v", err)
		c.JSONMap(map[string]interface{}{
			"code":    errCode,
			"message": "load data fail",
		}, nil)
		return
	}
	c.JSON(result, nil)
}
