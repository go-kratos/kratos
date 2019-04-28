package http

import (
	"context"

	"go-common/app/interface/main/mcn/conf"
	"go-common/app/interface/main/mcn/model/mcnmodel"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/http/blademaster"
)

// func bmHTTPErrorWithMsg(c *blademaster.Context, err error, msg string) {
// 	if c.IsAborted() {
// 		return
// 	}
// 	c.Error = err
// 	bcode := ecode.Cause(err)
// 	if msg == "" {
// 		msg = err.Error()
// 	}
// 	c.Render(http.StatusOK, render.JSON{
// 		Code:    bcode.Code(),
// 		Message: msg,
// 		Data:    nil,
// 	})
// }

// service的函数原型
type serviceFunc func(context context.Context, arg interface{}) (res interface{}, err error)

// response writer
type responseWriter func(c *blademaster.Context, arg interface{}, res interface{}, err error)

type argParser func(c *blademaster.Context, arg interface{}) (err error)

func argGetParser(c *blademaster.Context, arg interface{}) (err error) {
	err = c.Bind(arg)
	return
}

// argPostJSONParser .
// func argPostJSONParser(c *blademaster.Context, arg interface{}) (err error) {
// 	respBody, _ := ioutil.ReadAll(c.Request.Body)
// 	if err = json.Unmarshal(respBody, arg); err != nil {
// 		log.Error("json unmarshal fail, err=%v", err)
// 		return
// 	}
// 	if err = binding.Validator.ValidateStruct(arg); err != nil {
// 		log.Error("binding.Validator.ValidateStruct(%+v) error(%+v)", arg, err)
// 	}
// 	return
// }

func jsonWriter(c *blademaster.Context, arg interface{}, res interface{}, err error) {
	c.JSON(res, err)
}

// func csvWriter(c *blademaster.Context, arg interface{}, res interface{}, err error) {
// 	formater, ok := res.(model.CsvFormatter)
// 	if !ok {
// 		log.Error("res cannot convert CsvFommater, res type=%s", reflect.TypeOf(res).Name())
// 		c.String(ecode.ServerErr.Code(), "res cannot convert CsvFommater")
// 		return
// 	}

// 	c.Writer.Header().Set("Content-Type", "application/csv")
// 	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=\"%s\"", formater.GetFileName()))

// 	var buf = &bytes.Buffer{}
// 	var csvWriter = csv.NewWriter(buf)
// 	formater.ToCsv(csvWriter)
// 	csvWriter.Flush()
// 	c.Writer.Write(buf.Bytes())
// }

// func decideWriter(c *blademaster.Context, arg interface{}, res interface{}, err error) {
// 	var writer responseWriter
// 	var decider, ok = arg.(model.ExportArgInterface)
// 	if ok {
// 		switch decider.ExportFormat() {
// 		case "json":
// 			writer = jsonWriter
// 		case "csv":
// 			writer = csvWriter
// 		}
// 	}

// 	if writer != nil {
// 		writer(c, arg, res, err)
// 	} else {
// 		jsonWriter(c, arg, res, err)
// 	}
// }

type preBindFuncType func(*blademaster.Context) error
type preHandleFuncType func(*blademaster.Context, interface{}) error

func httpGenerateFunc(arg interface{}, sfunc serviceFunc, description string, parser argParser, writer responseWriter, preBindFuncs []preBindFuncType, preHandleFuncs []preHandleFuncType) func(*blademaster.Context) {
	return func(c *blademaster.Context) {
		var res interface{}
		var err error
		//var errMsg string
	exitswitch:
		switch {
		default:
			for _, f := range preBindFuncs {
				err = f(c)
				if err != nil {
					break exitswitch
				}
			}
			if err = parser(c, arg); err != nil {
				log.Error("%s, request argument bind fail, err=%v", description, err)
				//errMsg = fmt.Sprintf("wrong argument, %s", err.Error())
				err = ecode.RequestErr
				break
			}
			for _, f := range preHandleFuncs {
				err = f(c, arg)
				if err != nil {
					break exitswitch
				}
			}
			var scoreRes, e = sfunc(c, arg)
			err = e
			if e != nil {
				//errMsg = err.Error()
				log.Error("%s query fail, req=%+v, err=%+v", description, arg, err)
				break
			}
			log.Info("%s query ok, req=%+v, result=%+v", description, arg, scoreRes)
			res = scoreRes
		}
		if err != nil {
			//bmHTTPErrorWithMsg(c, err, errMsg)
			if c.IsAborted() {
				return
			}
			c.JSON(nil, err)
		} else {
			writer(c, arg, res, err)
		}
	}
}

func httpGetFunc(arg interface{}, sfunc serviceFunc, description string, preBindFuncs []preBindFuncType, preHandleFunc []preHandleFuncType) func(*blademaster.Context) {
	return httpGetFuncWithWriter(arg, sfunc, description, jsonWriter, preBindFuncs, preHandleFunc)
}

func httpGetFuncWithWriter(arg interface{}, sfunc serviceFunc, description string, writer responseWriter, preBindFuncs []preBindFuncType, preHandleFunc []preHandleFuncType) func(*blademaster.Context) {
	return httpGenerateFunc(arg, sfunc, description, argGetParser, writer, preBindFuncs, preHandleFunc)
}

func getMid(bc *blademaster.Context) int64 {
	mid, exists := bc.Get("mid")
	if !exists {
		return 0
	}
	realmid, ok := mid.(int64)
	if !ok {
		return 0
	}
	return realmid
}

// 根据cookie中的mid来设置请求中的mid，请求的arg必需是mcnmodel.CookieMidInterface
func getCookieMid(c *blademaster.Context, arg interface{}) preBindFuncType {
	return func(c *blademaster.Context) (err error) {
		var cookieInterface, ok = arg.(mcnmodel.CookieMidInterface)
		if !ok || cookieInterface == nil {
			return nil
		}
		var mid = getMid(c)
		if mid == 0 {
			err = ecode.AccessDenied
			return
		}
		cookieInterface.SetMid(mid)

		return
	}
}

func cheatReq(c *blademaster.Context, arg interface{}) (err error) {
	var mid = getMid(c)
	if mid == 0 {
		return
	}
	if conf.Conf.Other.IsWhiteList(mid) {
		var cheatInterface, ok = arg.(mcnmodel.CheatInterface)
		if ok && cheatInterface != nil {
			if cheatInterface.Cheat() {
				log.Warn("cheat happend, uri=%s, req=%+v", c.Request.RequestURI, arg)
			}
		}
	}
	return
}
