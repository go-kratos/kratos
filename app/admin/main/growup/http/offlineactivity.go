package http

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"go-common/app/admin/main/growup/dao/shell"
	"go-common/app/admin/main/growup/model"
	"go-common/app/admin/main/growup/model/offlineactivity"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/render"
	"io/ioutil"
	"net/http"
	"reflect"
)

func bmHTTPErrorWithMsg(c *blademaster.Context, err error, msg string) {
	if c.IsAborted() {
		return
	}
	c.Error = err
	bcode := ecode.Cause(err)
	if msg == "" {
		msg = err.Error()
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
func argPosJSONParser(c *blademaster.Context, arg interface{}) (err error) {
	respBody, _ := ioutil.ReadAll(c.Request.Body)
	if err = json.Unmarshal(respBody, arg); err != nil {
		log.Error("json unmarshal fail, err=%v", err)
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
		case "json":
			writer = jsonWriter
		case "csv":
			writer = csvWriter
		}
	}

	if writer != nil {
		writer(c, arg, res, err)
	} else {
		jsonWriter(c, arg, res, err)
	}
}

func httpGetFuncWithWriter(arg interface{}, sfunc serviceFunc, description string, writer responseWriter) func(*blademaster.Context) {
	return httpGenerateFunc(arg, sfunc, description, argGetParser, writer)
}

func httpGenerateFunc(arg interface{}, sfunc serviceFunc, description string, parser argParser, writer responseWriter, preFuncs ...func(*blademaster.Context) error) func(*blademaster.Context) {
	return func(c *blademaster.Context) {
		var res interface{}
		var err error
		var errMsg string
	exitswitch:
		switch {
		default:
			for _, f := range preFuncs {
				err = f(c)
				if err != nil {
					log.Error("request err=%s, arg=%+v", err, arg)
					break exitswitch
				}
			}

			if err = parser(c, arg); err != nil {
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
			bmHTTPErrorWithMsg(c, err, errMsg)
		} else {
			writer(c, arg, res, err)
		}
	}
}

func httpPostFunCheckCookie(arg interface{}, sfunc serviceFunc, description string, preFuncs ...func(*blademaster.Context) error) func(*blademaster.Context) {
	preFuncs = append(preFuncs, checkCookieFun)
	return httpGenerateFunc(arg, sfunc, description, argPosJSONParser, jsonWriter, preFuncs...)
}

// 贝壳回调
func offlineactivityShellCallback(c *blademaster.Context) {
	var err error
	var v = new(shell.OrderCallbackParam)
	if e := c.Bind(v); e != nil {
		err = e
		log.Error("parse arg error, err=%s", err)
		return
	}
	var result = "SUCCESS"
	if err = svr.ShellCallback(c, v); err != nil {
		log.Error("shell call back err, err=%s, arg=%v", err, v)
		result = "FAIL"
	}
	c.String(0, result)
}

func offlineactivityAdd(c *blademaster.Context) {
	httpPostFunCheckCookie(
		new(offlineactivity.AddActivityArg),
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return svr.AddOfflineActivity(context, arg.(*offlineactivity.AddActivityArg))
		},
		"offlineactivityAdd")(c)
}

func offlineactivityPreAdd(c *blademaster.Context) {
	httpPostFunCheckCookie(
		new(offlineactivity.AddActivityArg),
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return svr.PreAddOfflineActivity(context, arg.(*offlineactivity.AddActivityArg))
		},
		"offlienactivityPreAdd")(c)
}

func offlineactivityQueryActivity(c *blademaster.Context) {
	httpGetFuncWithWriter(
		new(offlineactivity.QueryActivityByIDArg),
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return svr.OfflineActivityQueryActivity(context, arg.(*offlineactivity.QueryActivityByIDArg))
		},
		"offlineactivityQueryActivity",
		decideWriter)(c)
}
func offlineactivityQueryUpBonusSummary(c *blademaster.Context) {
	httpGetFuncWithWriter(
		new(offlineactivity.QueryUpBonusByMidArg),
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return svr.OfflineActivityQueryUpBonusSummary(context, arg.(*offlineactivity.QueryUpBonusByMidArg))
		},

		"offlineactivityQueryUpBonus",
		decideWriter)(c)
}

func offlineactivityQueryUpBonusActivity(c *blademaster.Context) {
	httpGetFuncWithWriter(
		new(offlineactivity.QueryUpBonusByMidArg),
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return svr.OfflineActivityQueryUpBonusByActivity(context, arg.(*offlineactivity.QueryUpBonusByMidArg))
		},
		"offlineactivityQueryUpBonusActivity",
		decideWriter)(c)
}

func offlineactivityQueryMonth(c *blademaster.Context) {
	httpGetFuncWithWriter(
		new(offlineactivity.QueryActvityMonthArg),
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return svr.OfflineActivityQueryActivityByMonth(context, arg.(*offlineactivity.QueryActvityMonthArg))
		},
		"offlineactivityQueryMonth",
		decideWriter)(c)
}
