package dao

import (
	"context"
	"net/url"
	"strconv"
	"time"

	"go-common/app/admin/main/videoup-task/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_searchURL = "/x/admin/search/log"
)

// OutTime 退出时间,es的group by查询,最大1000条
func (d *Dao) OutTime(c context.Context, ids []int64) (mcases map[int64][]interface{}, err error) {
	mcases = make(map[int64][]interface{})
	params := url.Values{}
	params.Set("appid", "log_audit_group")
	params.Set("group", "uid")
	params.Set("uid", xstr.JoinInts(ids))
	params.Set("business", strconv.Itoa(model.LogClientConsumer))
	params.Set("action", strconv.Itoa(int(model.ActionHandsOFF)))
	params.Set("ps", strconv.Itoa(len(ids)))
	res := &model.SearchLogResult{}
	if err = d.hclient.Get(c, d.c.Host.Search+_searchURL, "", params, &res); err != nil {
		log.Error("log_audit_group d.hclient.Get error(%v)", err)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("log_audit_group ecode:%v", res.Code)
		return
	}
	for _, item := range res.Data.Result {
		mcases[item.UID] = []interface{}{item.Ctime}
	}
	log.Info("log_audit_group get: %s params:%s ret:%v", _searchURL, params.Encode(), res)
	return
}

// InQuitList 登入登出日志
func (d *Dao) InQuitList(c context.Context, uids []int64, bt, et string) (l []*model.InQuit, err error) {
	params := url.Values{}
	params.Set("appid", "log_audit")
	params.Set("business", strconv.Itoa(model.LogClientConsumer))
	if len(uids) > 0 {
		params.Set("uid", xstr.JoinInts(uids))
	}
	if len(bt) > 0 && len(et) > 0 {
		params.Set("ctime_from", bt)
		params.Set("ctime_to", et)
	}
	params.Set("order", "ctime")
	params.Set("sort", "desc")
	params.Set("ps", "10000")

	res := &model.SearchLogResult{}
	if err = d.hclient.Get(c, d.c.Host.Search+_searchURL, "", params, res); err != nil {
		log.Error("InQuitList d.hclient.Get error(%v)", err)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("InQuitList ecode:%v", res.Code)
		return
	}

	mapHelp := make(map[int64]*model.InQuit)
	for i := len(res.Data.Result) - 1; i >= 0; i-- {
		item := res.Data.Result[i]
		if item.Action == "0" {
			ctime, _ := time.Parse(model.TimeFormatSec, item.Ctime)
			iqlog := &model.InQuit{
				Date:   ctime.Format("2006-01-02"),
				UID:    item.UID,
				Uname:  item.Uname,
				InTime: ctime.Format("15:04:05"),
			}
			mapHelp[item.UID] = iqlog
			l = append([]*model.InQuit{iqlog}, l[:]...)
		}
		if item.Action == "1" {
			if iqlog, ok := mapHelp[item.UID]; ok {
				ctime, _ := time.Parse(model.TimeFormatSec, item.Ctime)
				if date := ctime.Format("2006-01-02"); date == iqlog.Date {
					iqlog.OutTime = ctime.Format("15:04:05")
				} else {
					iqlog.OutTime = ctime.Format(model.TimeFormatSec)
				}
			}
		}
	}

	return
}
