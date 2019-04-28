package http

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"github.com/siddontang/go-mysql/mysql"
	"go-common/app/admin/main/up/dao/upcrm"
	"go-common/app/admin/main/up/model/upcrmmodel"
	"go-common/app/admin/main/up/service"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/http/blademaster"
	"time"
)

func getTitleFields(compareType int) []string {
	var compareTitle = ""
	switch compareType {
	case upcrmmodel.CompareType7day:
		compareTitle = "7日前"
	case upcrmmodel.CompareType30day:
		compareTitle = "30日前"
	case upcrmmodel.CompareTypeMonthFirstDay:
		compareTitle = "本月1号"
	}
	return []string{
		"分数段",
		"昨日",
		"占比",
		compareTitle,
		"占比",
	}
}

func getScoreName(scoreType int) string {
	var name = "分数"
	switch scoreType {
	case upcrm.ScoreTypePr:
		name = "影响力分"
	case upcrm.ScoreTypeQuality:
		name = "质量分"
	case upcrm.ScoreTypeCredit:
		name = "信用分"
	}
	return name
}

func getScoreQueryContentField(result *upcrmmodel.ScoreQueryResult, index int) []string {
	if result == nil {
		return nil
	}
	var fields []string
	if index >= len(result.XAxis) {
		return nil
	}
	fields = append(fields, result.XAxis[index])
	if index < len(result.YAxis) {
		fields = append(fields, fmt.Sprintf("%d", result.YAxis[index].Value))
		fields = append(fields, fmt.Sprintf("%0.2f%%", float32(result.YAxis[index].Percent)/100.0))
	} else {
		fields = append(fields, "-", "-")
	}
	if index < len(result.CompareAxis) {
		fields = append(fields, fmt.Sprintf("%d", result.CompareAxis[index].Value))
		fields = append(fields, fmt.Sprintf("%0.2f%%", float32(result.CompareAxis[index].Percent)/100.0))
	} else {
		fields = append(fields, "-", "-")
	}
	return fields
}

func crmScoreQuery(c *blademaster.Context) {
	var arg = new(upcrmmodel.ScoreQueryArgs)
	var res interface{}
	var err error
	var errMsg string
	var scoreRes upcrmmodel.ScoreQueryResult
	switch {
	default:
		if err = c.Bind(arg); err != nil {
			log.Error("request argument bind fail, err=%v", err)
			errMsg = fmt.Sprintf("wrong argument, %s", err.Error())
			err = ecode.RequestErr
			break
		}

		scoreRes, err = Svc.Crmservice.ScoreQuery(c, arg)
		if err != nil {
			errMsg = err.Error()
			log.Error("score query fail, req=%+v, err=%+v", arg, err)
			break
		}
		log.Info("score query ok, req=%+v, result=%+v", arg, scoreRes)
		res = scoreRes
	}

	if err != nil {
		service.BmHTTPErrorWithMsg(c, err, errMsg)
	} else {
		if arg.Export == "csv" {
			c.Writer.Header().Set("Content-Type", "application/csv")
			c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=\"%s_%s.csv\"", getScoreName(arg.ScoreType), time.Now().Format(mysql.TimeFormat)))

			var buf = &bytes.Buffer{}
			var csvWriter = csv.NewWriter(buf)
			csvWriter.Write(getTitleFields(arg.CompareType))
			for i := 0; i < len(scoreRes.XAxis); i++ {
				csvWriter.Write(getScoreQueryContentField(&scoreRes, i))
			}
			csvWriter.Flush()
			c.Writer.Write(buf.Bytes())
		} else {
			c.JSON(res, err)
		}
	}
}

func crmScoreQueryUp(c *blademaster.Context) {
	httpQueryFunc(new(upcrmmodel.ScoreQueryUpArgs),
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return Svc.Crmservice.ScoreQueryUp(context, arg.(*upcrmmodel.ScoreQueryUpArgs))
		},
		"ScoreQueryUp")(c)
}

func crmScoreQueryUpHistory(c *blademaster.Context) {
	httpQueryFunc(new(upcrmmodel.ScoreQueryUpHistoryArgs),
		// 由于不支持泛型，所以这里只能再包一层
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return Svc.Crmservice.ScoreQueryUpHistory(context, arg.(*upcrmmodel.ScoreQueryUpHistoryArgs))
		},
		"ScoreQueryUpHistory")(c)
}

func crmPlayQueryInfo(c *blademaster.Context) {
	httpQueryFunc(new(upcrmmodel.PlayQueryArgs),
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return Svc.Crmservice.PlayQueryInfo(context, arg.(*upcrmmodel.PlayQueryArgs))
		},
		"PlayQueryInfo")(c)
}

func crmInfoQueryUp(c *blademaster.Context) {
	httpQueryFunc(new(upcrmmodel.InfoQueryArgs),
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return Svc.Crmservice.UpBaseInfoQuery(context, arg.(*upcrmmodel.InfoQueryArgs))
		},
		"QueryBaseUpInfo")(c)
}

func crmCreditLogQueryUp(c *blademaster.Context) {
	httpQueryFunc(new(upcrmmodel.CreditLogQueryArgs),
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return Svc.Crmservice.CreditLogQueryUp(context, arg.(*upcrmmodel.CreditLogQueryArgs))
		},
		"CreditLogQueryUp")(c)
}

func crmRankQueryList(c *blademaster.Context) {
	httpQueryFunc(new(upcrmmodel.UpRankQueryArgs),
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return Svc.Crmservice.UpRankQueryList(context, arg.(*upcrmmodel.UpRankQueryArgs))
		},
		"UpRankQueryList")(c)
}

func crmInfoAccountInfo(c *blademaster.Context) {
	httpQueryFunc(new(upcrmmodel.InfoAccountInfoArgs),
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return Svc.Crmservice.UpAccountInfo(context, arg.(*upcrmmodel.InfoAccountInfoArgs))
		},
		"InfoAccountInfo")(c)
}

func crmInfoSearch(c *blademaster.Context) {
	httpPostFunc(new(upcrmmodel.InfoSearchArgs),
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return Svc.Crmservice.UpInfoSearch(context, arg.(*upcrmmodel.InfoSearchArgs))
		},
		"UpInfoSearch")(c)
}

func testGetViewBase(c *blademaster.Context) {
	httpQueryFunc(new(upcrmmodel.TestGetViewBaseArgs),
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return Svc.Crmservice.TestGetViewBase(context, arg.(*upcrmmodel.TestGetViewBaseArgs))
		},
		"TestGetViewBase")(c)
}

func crmQueryUpInfoWithViewerData(c *blademaster.Context) {
	httpQueryFunc(new(upcrmmodel.UpInfoWithViewerArg),
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return Svc.Crmservice.QueryUpInfoWithViewerData(context, arg.(*upcrmmodel.UpInfoWithViewerArg))
		},
		"QueryUpInfoWithViewerData")(c)
}
