package http

import (
	"context"
	"encoding/json"
	"fmt"
	"go-common/app/admin/main/up/service"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/http/blademaster"
	"io/ioutil"
)

// service的函数原型
type serviceFunc func(context context.Context, arg interface{}) (res interface{}, err error)

// 由于不支持泛型，写得比较难看
// 很多重复的代码用下面来代替
func httpQueryFunc(arg interface{}, sfunc serviceFunc, description string) (httpFunc func(c *blademaster.Context)) {
	httpFunc = func(c *blademaster.Context) {
		//var arg = new(upcrmmodel.ScoreQueryUpHistoryArgs)
		var res interface{}
		var err error
		var errMsg string
		switch {
		default:
			if err = c.Bind(arg); err != nil {
				log.Error("%s, request argument bind fail, err=%v", description, err)
				errMsg = fmt.Sprintf("wrong argument, %s", err.Error())
				err = ecode.RequestErr
				break
			}

			var scoreRes, e = sfunc(c, arg)
			err = e
			if e != nil {
				errMsg = err.Error()
				log.Error("%s query fail, req=%+v, err=%+v", description, arg, err)
				break
			}
			log.Info("%s query ok, req=%+v, result=%+v", description, arg, scoreRes)
			res = scoreRes
		}

		if err != nil {
			service.BmHTTPErrorWithMsg(c, err, errMsg)
		} else {
			c.JSON(res, err)
		}
	}
	return
}

func httpPostFunc(arg interface{}, sfunc serviceFunc, description string) (httpFunc func(c *blademaster.Context)) {
	httpFunc = func(c *blademaster.Context) {
		//var arg = new(upcrmmodel.ScoreQueryUpHistoryArgs)
		var res interface{}
		var err error
		var errMsg string
		switch {
		default:
			respBody, _ := ioutil.ReadAll(c.Request.Body)
			if err = json.Unmarshal(respBody, arg); err != nil {
				log.Error("%s, json unmarshal fail, err=%v", description, err)
				errMsg = fmt.Sprintf("wrong argument, %s", err.Error())
				err = ecode.RequestErr
				break
			}

			scoreRes, e := sfunc(c, arg)
			err = e
			if e != nil {
				errMsg = err.Error()
				log.Error("%s query fail, req=%+v, err=%+v", description, arg, err)
				break
			}
			log.Info("%s query ok, req=%+v, result=%+v", description, arg, scoreRes)
			res = scoreRes
		}

		if err != nil {
			service.BmHTTPErrorWithMsg(c, err, errMsg)
		} else {
			c.JSON(res, err)
		}
	}
	return
}
