package dao

import (
	"context"
	"fmt"
	"net/url"

	"go-common/app/admin/main/dm/model"
	"go-common/library/log"
)

const (
	_sendNotify       = "/api/notify/send.user.notify.do"
	_msgReporterTitle = "举报处理结果通知"
	_msgPosterTitle   = "弹幕违规处理通知"
	_msgReporterKey   = "1_6_4"
	_msgPosterKey     = "1_6_5"

	_msgSubtitleTitle   = "字幕状态变更"
	_msgSubtitleUpKey   = "1_6_8"
	_msgSubtitleUserKey = "1_6_7"
)

type msgReturn struct {
	Code int64       `json:"code"`
	Ts   interface{} `json:"ts"`
	Data *struct {
		Mc         string  `json:"mc"`
		DataType   int8    `json:"data_type"`
		TotalCount int64   `json:"total_count"`
		ErrorCount int64   `json:"error_count"`
		ErrorMids  []int64 `json:"error_mid_list"`
	} `json:"data"`
}

// SendMsgToReporter send message
func (d *Dao) SendMsgToReporter(c context.Context, rptMsg *model.ReportMsg) (err error) {
	var (
		res = &msgReturn{}
	)
	params := url.Values{}
	params.Set("mc", _msgReporterKey)
	params.Set("title", _msgReporterTitle)
	params.Set("data_type", "4")
	params.Set("context", d.createReportContent(rptMsg))
	params.Set("mid_list", rptMsg.Uids)
	err = d.httpCli.Post(c, d.sendNotifyURI, "", params, res)
	if err != nil {
		log.Error("d.httpCli.Post(%s) error(%v)", d.sendNotifyURI+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		err = fmt.Errorf("return code:%d", res.Code)
		log.Error("d.httpCli.Post(%s) error(%v)", d.sendNotifyURI+"?"+params.Encode(), err)
	}
	return
}

func (d *Dao) createReportContent(rptMsg *model.ReportMsg) (content string) {
	if rptMsg.State == model.StatSecondAutoDelete ||
		rptMsg.State == model.StatFirstDelete ||
		rptMsg.State == model.StatSecondDelete {
		if rptMsg.Block != 0 { // 如果未封禁
			content = fmt.Sprintf(model.RptTemplate["del"], rptMsg.Title, rptMsg.Aid, rptMsg.Msg, "，该用户已被封禁", model.BlockReason[rptMsg.BlockReason])
		} else {
			content = fmt.Sprintf(model.RptTemplate["del"], rptMsg.Title, rptMsg.Aid, rptMsg.Msg, "", model.AdminRptReason[rptMsg.RptReason])
		}
	} else {
		content = fmt.Sprintf(model.RptTemplate["ignore"], rptMsg.Title, rptMsg.Aid, rptMsg.Msg)
	}
	return
}

// SendMsgToPoster send message
func (d *Dao) SendMsgToPoster(c context.Context, rptMsg *model.ReportMsg) (err error) {
	var (
		res = &msgReturn{}
	)
	params := url.Values{}
	params.Set("mc", _msgPosterKey)
	params.Set("title", _msgPosterTitle)
	params.Set("data_type", "4")
	msgContent, err := d.createPosterContent(rptMsg)
	if err != nil {
		log.Error("d.SendMsgToPoster error(%v)", err)
		return
	}
	params.Set("context", msgContent)
	params.Set("mid_list", rptMsg.Uids)
	err = d.httpCli.Post(c, d.sendNotifyURI, "", params, res)
	if err != nil {
		log.Error("d.httpCli.Post(%s) error(%v)", d.sendNotifyURI+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		err = fmt.Errorf("return code:%d", res.Code)
		log.Error("d.httpCli.Post(%s) error(%v)", d.sendNotifyURI+"?"+params.Encode(), err)
	}
	return
}

func (d *Dao) createPosterContent(rptMsg *model.ReportMsg) (content string, err error) {
	var (
		part1, part2, tmpl string
	)
	if rptMsg.Block > 0 {
		part1 = fmt.Sprintf("，并被封禁%d天", rptMsg.Block)
		part2 = model.BlockReason[rptMsg.BlockReason]
		tmpl, err = model.PosterBlockMsg(rptMsg.BlockReason)
		if err != nil {
			log.Error("reportModel.PosterBlockMsg error(%v)", err)
			return
		}
		content = fmt.Sprintf(tmpl, rptMsg.Title, rptMsg.Aid, rptMsg.Msg, part1, part2)
	} else if rptMsg.Block == -1 {
		part1 = "，并被永久封禁"
		part2 = model.BlockReason[rptMsg.BlockReason]
		tmpl, err = model.PosterBlockMsg(rptMsg.BlockReason)
		if err != nil {
			log.Error("reportModel.PosterBlockMsg error(%v)", err)
			return
		}
		content = fmt.Sprintf(tmpl, rptMsg.Title, rptMsg.Aid, rptMsg.Msg, part1, part2)
	} else {
		part1 = ""
		part2 = model.AdminRptReason[rptMsg.RptReason]
		tmpl, err = model.PosterAdminRptMsg(rptMsg.RptReason)
		if err != nil {
			log.Error("report.PosterAdminRptMsg error(%v)", err)
			return
		}
		content = fmt.Sprintf(tmpl, rptMsg.Title, rptMsg.Aid, rptMsg.Msg, part1, part2)
	}
	return
}

// SendMsgToSubtitleUp .
func (d *Dao) SendMsgToSubtitleUp(c context.Context, arg *model.NotifySubtitleUp) (err error) {
	var (
		res = &msgReturn{}
	)
	params := url.Values{}
	params.Set("mc", _msgSubtitleUpKey)
	params.Set("title", _msgSubtitleTitle)
	params.Set("data_type", "4")
	params.Set("context", arg.Msg())
	params.Set("mid_list", fmt.Sprint(arg.Mid))
	err = d.httpCli.Post(c, d.sendNotifyURI, "", params, res)
	if err != nil {
		log.Error("d.httpCli.Post(%s) error(%v)", d.sendNotifyURI+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		err = fmt.Errorf("return code:%d", res.Code)
		log.Error("d.httpCli.Post(%s) error(%v)", d.sendNotifyURI+"?"+params.Encode(), err)
	}
	return
}

// SendMsgToSubtitleUser .
func (d *Dao) SendMsgToSubtitleUser(c context.Context, arg *model.NotifySubtitleUser) (err error) {
	var (
		res = &msgReturn{}
	)
	params := url.Values{}
	params.Set("mc", _msgSubtitleUserKey)
	params.Set("title", _msgSubtitleTitle)
	params.Set("data_type", "4")
	params.Set("context", arg.Msg())
	params.Set("mid_list", fmt.Sprint(arg.Mid))
	err = d.httpCli.Post(c, d.sendNotifyURI, "", params, res)
	if err != nil {
		log.Error("d.httpCli.Post(%s) error(%v)", d.sendNotifyURI+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		err = fmt.Errorf("return code:%d", res.Code)
		log.Error("d.httpCli.Post(%s) error(%v)", d.sendNotifyURI+"?"+params.Encode(), err)
	}
	return
}
