package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-common/app/service/video/stream-mng/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"io/ioutil"
)

func openNotify(c *bm.Context) {
	p, err := parseNotifyBody(c)
	if err != nil {
		// log.Errorv(c, log.KV("log", fmt.Sprintf("open_notify_err = %v", err)))
		c.JSONMap(map[string]interface{}{"msg": err.Error()}, ecode.RequestErr)
		c.Set("output_data", err.Error())
		return
	}

	err = srv.StreamingNotify(c, p, true)
	if err != nil {
		// log.Errorv(c, log.KV("log", fmt.Sprintf("open_notify_err = %v", err)))
		c.JSONMap(map[string]interface{}{"msg": err.Error()}, ecode.RequestErr)
		c.Set("output_data", err.Error())
		// c.Abort()
		return
	}

	c.Set("output_data", "success")
	c.JSONMap(map[string]interface{}{"msg": "success"}, ecode.OK)
}

func closeNotify(c *bm.Context) {
	p, err := parseNotifyBody(c)
	if err != nil {
		c.JSONMap(map[string]interface{}{"msg": err.Error()}, ecode.RequestErr)
		c.Set("output_data", err.Error())
		return
	}

	err = srv.StreamingNotify(c, p, false)
	if err != nil {
		c.JSONMap(map[string]interface{}{"msg": err.Error()}, ecode.RequestErr)
		c.Set("output_data", err.Error())
		return
	}

	c.Set("output_data", "success")
	c.JSONMap(map[string]interface{}{"msg": "success"}, ecode.OK)
}

func parseNotifyBody(c *bm.Context) (*model.StreamingNotifyParam, error) {
	// log.Info("%v %v", c.Request.Header.Get("Content-Type"), c.Request.Form)
	switch c.Request.Header.Get("Content-Type") {
	case "application/x-www-form-urlencoded":
		// log.Info("%+v", c.Request.PostForm)
		if len(c.Request.PostForm) == 0 {
			return nil, errors.New("empty post body")
		}
		p := &model.StreamingNotifyParam{}
		p.Key = c.Request.PostFormValue("key")
		p.Sign = c.Request.PostFormValue("sign")
		p.SRC = c.Request.PostFormValue("src")
		p.StreamName = c.Request.PostFormValue("stream_name")
		p.SUID = c.Request.PostFormValue("suid")
		p.TS = json.Number(c.Request.PostFormValue("ts"))
		p.Type = json.Number(c.Request.PostFormValue("type"))
		// log.Info("%+v", p)
		// log.Infov(c, log.KV("log", fmt.Sprintf("notify_input = %+v", p)))
		c.Set("input_params", *p)
		return p, nil
	default:
		defer c.Request.Body.Close()
		b, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			return nil, err
		}

		if len(b) == 0 {
			return nil, errors.New("empty body")
		}

		var snp model.StreamingNotifyParam
		err = json.Unmarshal(b, &snp)
		if err != nil {
			log.Errorv(c, log.KV("log", fmt.Sprintf("notify_parse_body_error = %v", err)))
			return &snp, errors.New("invalid json body")
		}
		// log.Infov(c, log.KV("log", fmt.Sprintf("notify_input = %+v", snp)))
		c.Set("input_params", snp)
		return &snp, nil
	}
}
