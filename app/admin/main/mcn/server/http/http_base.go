package http

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"go-common/app/admin/main/mcn/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
	"go-common/library/net/http/blademaster/render"

	"github.com/pkg/errors"
)

func bmHTTPErrorWithMsg(c *blademaster.Context, err error) {
	if c.IsAborted() {
		return
	}
	c.Error = err
	bcode, ok := errors.Cause(err).(ecode.Codes)
	var msg string
	if !ok {
		msg = err.Error()
		bcode = ecode.String(err.Error())
	} else {
		msg = bcode.Message()
	}

	c.Render(http.StatusOK, render.JSON{
		Code:    bcode.Code(),
		Message: msg,
		Data:    nil,
	})
}

func checkCookieFun(c *blademaster.Context) (err error) {
	_, _, err = checkCookie(c)
	return
}

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
func argPostJSONParser(c *blademaster.Context, arg interface{}) (err error) {
	respBody, _ := ioutil.ReadAll(c.Request.Body)
	if err = json.Unmarshal(respBody, arg); err != nil {
		log.Error("json unmarshal fail, err=%v", err)
		return
	}
	if err = binding.Validator.ValidateStruct(arg); err != nil {
		log.Error("binding.Validator.ValidateStruct(%+v) error(%+v)", arg, err)
	}
	return
}

func jsonWriter(c *blademaster.Context, arg interface{}, res interface{}, err error) {
	c.JSON(res, err)
}

func csvWriter(c *blademaster.Context, arg interface{}, res interface{}, err error) {
	formater, ok := res.(model.CsvFormatter)
	if !ok {
		log.Error("res cannot convert CsvFommater, res type=%s", reflect.TypeOf(res).Name())
		c.String(ecode.ServerErr.Code(), "res cannot convert CsvFommater")
		return
	}

	c.Writer.Header().Set("Content-Type", "application/csv")
	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=\"%s\"", formater.GetFileName()))

	var buf = &bytes.Buffer{}
	var csvWriter = csv.NewWriter(buf)
	formater.ToCsv(csvWriter)
	csvWriter.Flush()
	c.Writer.Write(buf.Bytes())
}

func decideWriter(c *blademaster.Context, arg interface{}, res interface{}, err error) {
	var writer responseWriter
	var decider, ok = arg.(model.ExportArgInterface)
	if ok {
		switch decider.ExportFormat() {
		case model.ResponeModelJSON:
			writer = jsonWriter
		case model.ResponeModelCSV:
			writer = csvWriter
		}
	}

	if writer != nil {
		writer(c, arg, res, err)
	} else {
		jsonWriter(c, arg, res, err)
	}
}

func httpGenerateFunc(arg interface{}, sfunc serviceFunc, description string, parser argParser, writer responseWriter, preFuncs ...func(*blademaster.Context) error) func(*blademaster.Context) {
	return func(c *blademaster.Context) {
		var res interface{}
		var err error
	exitswitch:
		switch {
		default:
			for _, f := range preFuncs {
				err = f(c)
				if err != nil {
					break exitswitch
				}
			}
			if err = parser(c, arg); err != nil {
				log.Error("%s, request argument bind fail, err=%v", description, err)
				err = ecode.RequestErr
				break
			}

			var scoreRes, e = sfunc(c, arg)
			err = e
			if e != nil {
				log.Error("%s query fail, req=%+v, err=%+v", description, arg, err)
				break
			}
			log.Info("%s query ok, req=%+v, result=%+v", description, arg, scoreRes)
			res = scoreRes
		}
		if err != nil {
			bmHTTPErrorWithMsg(c, err)
		} else {
			writer(c, arg, res, err)
		}
	}
}

func httpPostJSONCookie(arg interface{}, sfunc serviceFunc, description string, preFuncs ...func(*blademaster.Context) error) func(*blademaster.Context) {
	preFuncs = append(preFuncs, checkCookieFun)
	return httpGenerateFunc(arg, sfunc, description, argPostJSONParser, jsonWriter, preFuncs...)
}

// func httpPostFormCookie(arg interface{}, sfunc serviceFunc, description string, preFuncs ...func(*blademaster.Context) error) func(*blademaster.Context) {
// 	preFuncs = append(preFuncs, checkCookieFun)
// 	return httpGenerateFunc(arg, sfunc, description, argGetParser, jsonWriter, preFuncs...)
// }

func httpGetFunCheckCookie(arg interface{}, sfunc serviceFunc, description string, preFuncs ...func(*blademaster.Context) error) func(*blademaster.Context) {
	preFuncs = append(preFuncs, checkCookieFun)
	return httpGenerateFunc(arg, sfunc, description, argGetParser, jsonWriter, preFuncs...)
}

// func httpGetFunc(arg interface{}, sfunc serviceFunc, description string) func(*blademaster.Context) {
// 	return httpGetFuncWithWriter(arg, sfunc, description, jsonWriter)
// }

func httpGetFuncWithWriter(arg interface{}, sfunc serviceFunc, description string, writer responseWriter) func(*blademaster.Context) {
	return httpGenerateFunc(arg, sfunc, description, argGetParser, writer)
}

func httpGetWriterByExport(arg interface{}, sfunc serviceFunc, description string) func(*blademaster.Context) {
	return httpGetFuncWithWriter(arg, sfunc, description, decideWriter)
}

//func offlineactivityAdd(c *blademaster.Context) {
//	httpPostFunCheckCookie(
//		new(offlineactivity.AddActivityArg),
//		func(context context.Context, arg interface{}) (res interface{}, err error) {

//			return svr.AddOfflineActivity(context, arg.(*offlineactivity.AddActivityArg))
//		},
//		"offlineactivityAdd")(c)
//}
