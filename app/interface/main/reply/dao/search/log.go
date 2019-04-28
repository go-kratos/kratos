package search

import (
	"context"
	"net/url"
	"strconv"
	"time"

	"go-common/app/interface/main/reply/model/adminlog"
	"go-common/library/log"
	"go-common/library/xstr"
)

type searchAdminLog struct {
	ReplyID  int64  `json:"rpid"`     // 评论id
	AdminID  int64  `json:"adminid"`  // 操作人
	State    int32  `json:"state"`    // 操作人身份
	ReplyMid int64  `json:"replymid"` // 评论人
	CTime    string `json:"ctime"`    // 删除时间
}

// LogPaginate paginating the admin logs for size of 'pageSize', and returning the number of reporting, the number of admin logs delete by administrator
func (dao *Dao) LogPaginate(c context.Context, oid int64, tp int, states []int64, curPage, pageSize int, startTime string, now time.Time) (logs []*adminlog.AdminLog, replyCount, reportCount, pageCount, total int64, err error) {
	params := url.Values{}
	params.Set("appid", "replylog")
	params.Set("oid", strconv.FormatInt(oid, 10))
	params.Set("type", strconv.Itoa(tp))
	params.Set("delstats", xstr.JoinInts(states))
	params.Set("start_time", startTime)
	params.Set("pagesize", strconv.Itoa(pageSize))
	params.Set("page", strconv.Itoa(curPage))
	var res struct {
		Code        int               `json:"code"`
		Logs        []*searchAdminLog `json:"result"`
		ReplyCount  int64             `json:"adminDeletedNum"`
		ReportCount int64             `json:"reportNum"`
		Page        int64             `json:"page"`
		PageSize    int64             `json:"pagesize"`
		PageCount   int64             `json:"pagecount"`
		Total       int64             `json:"total"`
	}
	if err = dao.httpCli.Get(c, dao.logURL, "", params, &res); err != nil {
		log.Error("adminlog url(%v),err (%v)", dao.logURL+"?"+params.Encode(), err)
		return
	}
	if res.Logs == nil {
		logs = make([]*adminlog.AdminLog, 0)
	}
	for _, log := range res.Logs {
		var (
			tmp = &adminlog.AdminLog{}
		)
		tmp.ReplyID = log.ReplyID
		tmp.State = log.State
		tmp.AdminID = log.AdminID
		tmp.CTime = log.CTime
		tmp.ReplyMid = log.ReplyMid
		logs = append(logs, tmp)
	}
	replyCount = res.ReplyCount
	reportCount = res.ReportCount
	pageCount = res.PageCount
	total = res.Total
	return
}
