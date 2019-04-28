package http

import (
	"go-common/app/interface/bbq/app-bbq/api/http/v1"
	bm "go-common/library/net/http/blademaster"
	"strings"

	"github.com/json-iterator/go"
)

func appSetting(c *bm.Context) {
	args := &v1.AppSettingRequest{}
	if err := c.Bind(args); err != nil {
		return
	}
	data, err := srv.AppSetting(c, args)
	if err != nil {
		c.JSON(nil, err)
		return
	}

	if strings.ToLower(args.Base.Client) == "ios" {
		c.JSON(data, err)
		return
	}

	resp := struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		TTL     int         `json:"ttl"`
		Data    interface{} `json:"data,omitempty"`
	}{
		Code:    0,
		Message: "0",
		TTL:     1,
		Data:    data,
	}

	str, _ := jsoniter.MarshalToString(resp)
	c.String(0, strings.Replace(str, "\\\\", "\\", -1))
}

func appPackage(c *bm.Context) {
	args := struct {
		Lastest int `json:"lastest" form:"lastest"`
	}{}
	c.Bind(&args)
	c.JSON(srv.AppPackage(c, args.Lastest))
}
