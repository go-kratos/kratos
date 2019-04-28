package http

import (
	"encoding/json"
	"go-common/app/service/video/stream-mng/service"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"io/ioutil"
)

// streamValidate 流鉴权接口
func streamValidate(c *bm.Context) {
	var vp service.ValidateParams
	switch c.Request.Header.Get("Content-Type") {
	case "application/x-www-form-urlencoded":
		if len(c.Request.PostForm) == 0 {
			c.Set("output_data", "stream_valid_err = empty post body")
			c.JSONMap(map[string]interface{}{"message": "empty post body"}, ecode.RequestErr)
			c.Abort()
			return
		}
		vp.Key = c.Request.PostFormValue("key")
		vp.StreamName = c.Request.PostFormValue("stream_name")
		vp.Src = c.Request.PostFormValue("src")
		vp.Type = json.Number(c.Request.PostFormValue("type"))
	default:
		defer c.Request.Body.Close()

		b, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.Set("output_data", err.Error())
			c.JSONMap(map[string]interface{}{"message": err.Error()}, ecode.RequestErr)
			c.Abort()
			return
		}

		if len(b) == 0 {
			c.Set("output_data", "stream_valid_err = empty params")
			c.JSONMap(map[string]interface{}{"message": "empty params"}, ecode.RequestErr)
			c.Abort()
			return
		}

		err = json.Unmarshal(b, &vp)

		if err != nil {
			c.Set("output_data", err.Error())
			c.JSONMap(map[string]interface{}{"message": err.Error()}, ecode.RequestErr)
			c.Abort()
			return
		}
	}

	c.Set("input_params", vp)

	permission, err := srv.CheckStreamValidate(c, &vp, false)
	if err != nil {
		c.Set("output_data", err.Error())
		if err.Error() == "room is closed" {
			c.JSONMap(map[string]interface{}{"message": err.Error()}, ecode.LimitExceed)
		} else {
			c.JSONMap(map[string]interface{}{"message": err.Error()}, ecode.RequestErr)
		}
		c.Abort()
		return
	}
	c.Set("output_data", permission)
	c.JSONMap(map[string]interface{}{"data": map[string]int{"permission": permission}}, nil)
}
