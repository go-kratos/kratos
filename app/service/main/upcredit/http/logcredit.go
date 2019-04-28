package http

import (
	"encoding/json"
	"fmt"
	"go-common/app/service/main/upcredit/model/calculator"
	"go-common/app/service/main/upcredit/model/upcrmmodel"
	"go-common/app/service/main/upcredit/service"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/http/blademaster"
	"io/ioutil"
	"time"
)

func logCredit(c *blademaster.Context) {
	var data interface{}
	var err error
	var errMsg string
	var r = new(upcrmmodel.ArgCreditLogAdd)
	switch {
	default:
		var body, _ = ioutil.ReadAll(c.Request.Body)
		if err = json.Unmarshal(body, r); err != nil {
			log.Error("request argument json decode fail, err=%v", err)
			errMsg = fmt.Sprintf("wrong argument, %s, body={%s}", err.Error(), string(body))
			err = ecode.RequestErr
			break
		}

		err = Svc.LogCredit(c, r)
	}
	if err != nil {
		service.BmHTTPErrorWithMsg(c, err, errMsg)
	} else {
		c.JSON(data, err)
	}
}

func recalc(c *blademaster.Context) {
	var data interface{}
	var err error
	var errMsg string
	var r = new(struct {
		TableNum int    `form:"tablenum"`
		CalcDate string `form:"date"`
		AllTable bool   `form:"all_table" default:"false"`
	})
	switch {
	default:
		if err = c.Bind(r); err != nil {
			log.Error("param error")
			err = ecode.RequestErr
			break
		}

		var calc = calculator.New(Svc.CreditScoreInputChan)
		var date = time.Time{}
		if r.CalcDate != "" {
			date, _ = time.Parse("2006-01-02", r.CalcDate)
		} else {
			date = time.Now()
		}
		if r.AllTable {
			Svc.CalcSvc.AddCalcJob(date)
		} else {
			go calc.CalcLogTable(r.TableNum, date, nil)
		}
		log.Info("start calculate process, req=%+v", r)
	}
	if err != nil {
		service.BmHTTPErrorWithMsg(c, err, errMsg)
	} else {
		c.JSON(data, err)
	}
}

func scoreGet(c *blademaster.Context) {
	var data interface{}
	var err error
	var errMsg string
	var r = new(upcrmmodel.ArgMidDate)
	switch {
	default:
		if err = c.Bind(r); err != nil {
			log.Error("param error")
			err = ecode.RequestErr
			break
		}
		var arg = upcrmmodel.GetScoreParam{
			Mid:       r.Mid,
			ScoreType: r.ScoreType,
		}
		var now = time.Now()
		if r.Days > 60 {
			r.Days = 60
		}

		if r.FromDate != "" {
			arg.FromDate, err = time.Parse(upcrmmodel.DateStr, r.FromDate)
			if err != nil {
				errMsg = err.Error()
				err = ecode.RequestErr
				break
			}
		}
		if r.ToDate != "" {
			arg.ToDate, err = time.Parse(upcrmmodel.DateStr, r.FromDate)
			if err != nil {
				errMsg = err.Error()
				err = ecode.RequestErr
				break
			}
		}

		if r.Days != 0 {
			var y, m, d = now.Date()
			arg.ToDate = time.Date(y, m, d, 0, 0, 0, 0, now.Location())
			arg.FromDate = arg.ToDate.AddDate(0, 0, -r.Days)
		}
		var result, e = Svc.GetCreditScore(c, &arg)
		err = e
		if err != nil {
			log.Error("fail to get credit score, req=%+v, err=%+v", arg, err)
			break
		}
		data = map[string]interface{}{
			"score_list": result,
		}
		log.Info("get credit score, req=%+v, datalen=%d", r, len(result))
	}
	if err != nil {
		service.BmHTTPErrorWithMsg(c, err, errMsg)
	} else {
		c.JSON(data, err)
	}
}

func logGet(c *blademaster.Context) {
	var data interface{}
	var err error
	var errMsg string
	var r = new(upcrmmodel.ArgGetLogHistory)
	switch {
	default:
		if err = c.Bind(r); err != nil {
			log.Error("param error")
			errMsg = "param error"
			err = ecode.RequestErr
			break
		}

		var result, e = Svc.GetCreditLog(c, r)
		err = e
		if err != nil {
			errMsg = err.Error()
			log.Error("fail to get credit log, req=%+v, err=%+v", r, err)
			break
		}
		data = map[string]interface{}{
			"log_list": result,
		}
		log.Info("get credit log ok, req=%+v, datalen=%d", r, len(result))
	}
	if err != nil {
		service.BmHTTPErrorWithMsg(c, err, errMsg)
	} else {
		c.JSON(data, err)
	}
}

func calcSection(c *blademaster.Context) {
	var data interface{}
	var err error
	var errMsg string
	var r = new(struct {
		CalcDate string `form:"date"`
	})
	switch {
	default:
		if err = c.Bind(r); err != nil {
			log.Error("param error")
			err = ecode.RequestErr
			break
		}

		var date = time.Time{}
		if r.CalcDate != "" {
			date, _ = time.Parse("2006-01-02", r.CalcDate)
		} else {
			date = time.Now().AddDate(0, 0, -1)
		}
		var job = &service.CalcStatisticJob{
			ID:   1,
			Date: date,
			Svc:  Svc.CalcSvc,
		}
		Svc.CalcSvc.JobChannel <- job
		log.Info("start calculate process, req=%+v, job=%+v", r, job)
	}
	if err != nil {
		service.BmHTTPErrorWithMsg(c, err, errMsg)
	} else {
		c.JSON(data, err)
	}
}

func test(c *blademaster.Context) {
	Svc.Test()
}
