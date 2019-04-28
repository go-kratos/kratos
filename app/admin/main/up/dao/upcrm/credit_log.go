package upcrm

import (
	"context"
	"fmt"
	"net/url"

	"go-common/app/admin/main/up/model/upcrmmodel"
	"go-common/library/log"
)

const (
	creditLogURL = "http://api.bilibili.co/x/internal/upcredit/log/get"
)

//
//type ArgGetLogHistory struct {
//	Mid      int64        `params:"mid;Required"`
//	FromDate systime.Time `params:"from_date"`
//	ToDate   systime.Time `params:"to_date"`
//	Limit    int          `params:"limit" default:"20"`
//}

//LogList it's log list
type LogList struct {
	LogList []upcrmmodel.SimpleCreditLogWithContent `json:"log_list"`
}

//CreditLogHTTPResult it's log result from http server
type CreditLogHTTPResult struct {
	Code int     `json:"code"`
	Data LogList `json:"data"`
	Msg  string  `json:"message"`
}

//GetCreditLog get credit log from upcredit server
func (d *Dao) GetCreditLog(mid int64, limit int) (result []upcrmmodel.SimpleCreditLogWithContent, err error) {
	params := url.Values{}
	params.Set("mid", fmt.Sprintf("%d", mid))
	//params.Set("from_date", fromDate.Time().Format("2006-01-02"))
	//params.Set("to_date", toDate.Time().Format("2006-01-02"))
	params.Set("limit", fmt.Sprintf("%d", limit))

	var c = context.Background()
	var httpResult = CreditLogHTTPResult{}
	err = d.httpClient.Get(c, creditLogURL, "", params, &httpResult)
	if err != nil {
		log.Error("get credit log http fail, request=%s?%s, err=%+v", creditLogURL, params.Encode(), err)
		return
	}

	log.Info("get credit log http ok, result=%+v", httpResult)
	result = httpResult.Data.LogList
	return
}
