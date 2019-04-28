package assist

import (
	"context"
	"go-common/app/interface/main/creative/model/assist"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/xstr"
	"net/url"
	"strconv"
	"strings"
)

const (
	// api
	_addAssistURI        = "/x/internal/assist/add"
	_delAssistURI        = "/x/internal/assist/del"
	_getAssistInfoURI    = "/x/internal/assist/info"
	_getAssistLogsURI    = "/x/internal/assist/logs"
	_addAssistLogURI     = "/x/internal/assist/log/add"
	_getAssistURI        = "/x/internal/assist/assists"
	_getAssistLogInfoURI = "/x/internal/assist/log/info"
	_revocAssistLogURI   = "/x/internal/assist/log/cancel"
	_getAssistLogObjURI  = "/x/internal/assist/log/obj"
	_getAssistStatURI    = "/x/internal/assist/stat"
)

// Assists get all Assists from assist service.
func (d *Dao) Assists(c context.Context, mid int64, ip string) (assists []*assist.Assist, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int              `json:"code"`
		Data []*assist.Assist `json:"data"`
	}
	if err = d.client.Get(c, d.assistListURL, ip, params, &res); err != nil {
		log.Error("d.client.Get(%s,%s,%s) err(%v)", d.assistListURL, ip, params.Encode(), err)
		err = ecode.CreativeAssistErr
		return
	}
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("d.client.Get(%s,%s,%s) err(%v)", d.assistListURL, ip, params.Encode(), err)
		return
	}
	assists = res.Data
	return
}

// AssistLog get assist log info from assist service
func (d *Dao) AssistLog(c context.Context, mid, assistMid, logID int64, ip string) (assistLog *assist.AssistLog, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("assist_mid", strconv.FormatInt(assistMid, 10))
	params.Set("log_id", strconv.FormatInt(logID, 10))
	var res struct {
		Code int               `json:"code"`
		Data *assist.AssistLog `json:"data"`
	}
	if err = d.client.Get(c, d.assistLogInfoURL, ip, params, &res); err != nil {
		log.Error("d.client.Get(%s,%s,%s) err(%v)", d.assistLogInfoURL, ip, params.Encode(), err)
		err = ecode.CreativeAssistErr
		return
	}
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("d.client.Get(%s,%s,%s) err(%v)", d.assistLogInfoURL, ip, params.Encode(), err)
		return
	}
	assistLog = res.Data
	return
}

// AssistLogs get logs from assist service.
func (d *Dao) AssistLogs(c context.Context, mid, assistMid, pn, ps, stime, etime int64, ip string) (logs []*assist.AssistLog, pager map[string]int64, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("assist_mid", strconv.FormatInt(assistMid, 10))
	params.Set("pn", strconv.FormatInt(pn, 10))
	params.Set("ps", strconv.FormatInt(ps, 10))
	params.Set("stime", strconv.FormatInt(stime, 10))
	params.Set("etime", strconv.FormatInt(etime, 10))
	var res struct {
		Code  int                 `json:"code"`
		Data  []*assist.AssistLog `json:"data"`
		Pager map[string]int64    `json:"pager"`
	}
	if err = d.client.Get(c, d.assistLogsURL, ip, params, &res); err != nil {
		log.Error("d.client.Get(%s,%s,%s) err(%v)", d.assistLogsURL, ip, params.Encode(), err)
		err = ecode.CreativeAssistErr
		return
	}
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("d.client.Get(%s,%s,%s) err(%v)", d.assistLogsURL, ip, params.Encode(), err)
		return
	}
	for _, v := range res.Data {
		if v.Type == 3 && v.Action == 8 {
			detailT3A8 := strings.Split(v.Detail, ":")
			if len(detailT3A8) > 1 {
				v.Detail = " 用户UID:" + detailT3A8[len(detailT3A8)-1]
			}
		}
		if v.Type == 3 && v.Action == 9 {
			detailT3A9 := strings.Split(v.Detail, ":")
			if len(detailT3A9) > 1 {
				v.Detail = " 用户UID:" + detailT3A9[len(detailT3A9)-1]
			}
		}
	}
	logs = res.Data
	pager = res.Pager
	return
}

// AddAssist add assist
func (d *Dao) AddAssist(c context.Context, mid, assistMid int64, ip, upUname string) (err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("up_uname", upUname)
	params.Set("assist_mid", strconv.FormatInt(assistMid, 10))
	var res struct {
		Code int `json:"code"`
	}
	log.Info("AddOrDelAssist d.client.Post(%s,%s,%s) err(%v)", d.assistAddURL, ip, params.Encode(), err)
	if err = d.client.Post(c, d.assistAddURL, ip, params, &res); err != nil {
		log.Error("d.client.Post(%s,%s,%s) err(%v)", d.assistAddURL, ip, params.Encode(), err)
		err = ecode.CreativeAssistErr
		return
	}
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("d.client.Post(%s,%s,%s) err(%v)", d.assistAddURL, ip, params.Encode(), err)
		return
	}
	return
}

// DelAssist cancel assist
func (d *Dao) DelAssist(c context.Context, mid, assistMid int64, ip, upUname string) (err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("assist_mid", strconv.FormatInt(assistMid, 10))
	params.Set("up_uname", upUname)
	var res struct {
		Code int `json:"code"`
	}
	log.Info("AddOrDelAssist d.client.Post(%s,%s,%s) err(%v)", d.assistDelURL, ip, params.Encode(), err)
	if err = d.client.Post(c, d.assistDelURL, ip, params, &res); err != nil {
		log.Error("d.client.Post(%s,%s,%s) err(%v)", d.assistDelURL, ip, params.Encode(), err)
		err = ecode.CreativeAssistErr
		return
	}
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("d.client.Post(%s,%s,%s) err(%v)", d.assistDelURL, ip, params.Encode(), err)
		return
	}
	return
}

// RevocAssistLog calcel assistlog action
func (d *Dao) RevocAssistLog(c context.Context, mid, assistMid, logID int64, ip string) (err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("assist_mid", strconv.FormatInt(assistMid, 10))
	params.Set("log_id", strconv.FormatInt(logID, 10))
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Post(c, d.assistLogRevocURL, ip, params, &res); err != nil {
		log.Error("d.client.Post(%s,%s,%s) err(%v)", d.assistLogRevocURL, ip, params.Encode(), err)
		err = ecode.CreativeAssistErr
		return
	}
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("d.client.Post(%s,%s,%s) err(%v)", d.assistLogRevocURL, ip, params.Encode(), err)
		return
	}
	return
}

// Stat get assists stat
func (d *Dao) Stat(c context.Context, mid int64, assistMids []int64, ip string) (stat map[int64]map[int8]map[int8]int, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("assmids", xstr.JoinInts(assistMids))
	var res struct {
		Code int                             `json:"code"`
		Data map[int64]map[int8]map[int8]int `json:"data"`
	}
	if err = d.client.Get(c, d.assistStatURL, ip, params, &res); err != nil {
		log.Error("d.client.Get(%s,%s,%s) err(%v)", d.assistStatURL, ip, params.Encode(), err)
		err = ecode.CreativeAssistErr
		return
	}
	if res.Code != 0 {
		log.Error("d.client.Get(%s,%s,%s) err(%v)", d.assistStatURL, ip, params.Encode(), err)
		err = ecode.Int(res.Code)
		return
	}
	stat = res.Data
	return
}

// Info  check if is assist
func (d *Dao) Info(c context.Context, mid, assistMid int64, ip string) (assist int8, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("assist_mid", strconv.FormatInt(assistMid, 10))
	params.Set("type", "1")
	var res struct {
		Code int `json:"code"`
		Data struct {
			Assist int8 `json:"assist"`
		} `json:"data"`
	}
	if err = d.client.Get(c, d.assistInfoURL, ip, params, &res); err != nil {
		log.Error("d.client.Post(%s,%s,%s) err(%v)", d.assistInfoURL, ip, params.Encode(), err)
		err = ecode.CreativeAssistErr
		return
	}
	if res.Code != 0 {
		log.Error("d.client.Post(%s,%s,%s) err(%v)", d.assistInfoURL, ip, params.Encode(), err)
		err = ecode.Int(res.Code)
		return
	}
	assist = res.Data.Assist
	return
}

// AssistLogObj get assist log info from assist service
func (d *Dao) AssistLogObj(c context.Context, tp, act int8, mid, objID int64) (assLog *assist.AssistLog, err error) {
	params := url.Values{}
	params.Set("type", strconv.FormatInt(int64(tp), 10))
	params.Set("action", strconv.FormatInt(int64(act), 10))
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("object_id", strconv.FormatInt(objID, 10))
	var res struct {
		Code int               `json:"code"`
		Data *assist.AssistLog `json:"data"`
	}
	if err = d.client.Get(c, d.assistLogObjURL, "", params, &res); err != nil {
		log.Error("d.client.Get(%s,%s) err(%v)", d.assistLogObjURL, params.Encode(), err)
		err = ecode.CreativeAssistErr
		return
	}
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("d.client.Get(%s,%s) err(%v)", d.assistLogObjURL, params.Encode(), err)
		return
	}
	assLog = res.Data
	return
}
