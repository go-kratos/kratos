package dao

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"go-common/app/admin/main/block/conf"
	"go-common/app/admin/main/block/model"
	"go-common/library/ecode"

	"github.com/pkg/errors"
)

// BlackhouseBlock .
func (d *Dao) BlackhouseBlock(c context.Context, p *model.ParamBatchBlock) (err error) {
	midStrs := make([]string, len(p.MIDs))
	for i, mid := range p.MIDs {
		midStrs[i] = fmt.Sprintf("%d", mid)
	}
	params := url.Values{}
	params.Set("mids", strings.Join(midStrs, ","))
	params.Set("oper_id", fmt.Sprintf("%d", p.AdminID))
	params.Set("operator_name", p.AdminName)
	params.Set("blocked_days", fmt.Sprintf("%d", p.Duration))
	switch p.Action {
	case model.BlockActionForever:
		params.Set("blocked_forever", "1")
		params.Set("punish_type", "3")
	default:
		params.Set("blocked_forever", "0")
		params.Set("punish_type", "2")
	}
	params.Set("blocked_remark", p.Comment)
	params.Set("origin_type", fmt.Sprintf("%d", p.Area))
	params.Set("punish_time", fmt.Sprintf("%d", time.Now().Unix()))
	params.Set("reason_type", fmt.Sprintf("%d", parseReasonType(p.Reason)))

	var res struct {
		Code int `json:"code"`
	}
	if err = d.httpClient.Post(c, conf.Conf.Property.BlackHouseURL, "", params, &res); err != nil {
		err = errors.WithStack(err)
		return
	}
	if res.Code != 0 {
		err = errors.WithStack(ecode.Int(res.Code))
		return
	}
	return
}

func parseReasonType(msg string) (t int) {
	switch msg {
	case "刷屏":
		t = 1
	case "抢楼":
		t = 2
	case "发布色情低俗信息":
		t = 3
	case "发布赌博诈骗信息":
		t = 4
	case "发布违禁相关信息", "发布违禁信息":
		t = 5
	case "发布垃圾广告信息":
		t = 6
	case "发布人身攻击言论":
		t = 7
	case "发布侵犯他人隐私信息":
		t = 8
	case "发布引战言论":
		t = 9
	case "发布剧透信息":
		t = 10
	case "恶意添加无关标签":
		t = 11
	case "恶意删除他人标签":
		t = 12
	case "发布色情信息":
		t = 13
	case "发布低俗信息":
		t = 14
	case "发布暴力血腥信息":
		t = 15
	case "涉及恶意投稿行为":
		t = 16
	case "发布非法网站信息":
		t = 17
	case "发布传播不实信息":
		t = 18
	case "发布怂恿教唆信息":
		t = 19
	case "恶意刷屏":
		t = 20
	case "账号违规":
		t = 21
	case "恶意抄袭":
		t = 22
	case "冒充自制原创":
		t = 23
	case "发布青少年不良内容":
		t = 24
	case "破坏网络安全":
		t = 25
	case "发布虚假误导信息":
		t = 26
	case "仿冒官方认证账号":
		t = 27
	case "发布不适宜内容":
		t = 28
	case "违反运营规则":
		t = 29
	case "恶意创建话题":
		t = 30
	case "发布违规抽奖":
		t = 31
	default:
		t = 0
	}
	return
}

// SendSysMsg send sys msg.
func (d *Dao) SendSysMsg(c context.Context, code string, mids []int64, title string, content string, remoteIP string) (err error) {
	params := url.Values{}
	params.Set("mc", code)
	params.Set("title", title)
	params.Set("data_type", "4")
	params.Set("context", content)
	params.Set("mid_list", midsToParam(mids))
	var res struct {
		Code int `json:"code"`
		Data *struct {
			Status int8   `json:"status"`
			Remark string `json:"remark"`
		} `json:"data"`
	}
	if err = d.httpClient.Post(c, conf.Conf.Property.MSGURL, remoteIP, params, &res); err != nil {
		err = errors.WithStack(err)
		return
	}
	if res.Code != 0 {
		err = errors.WithStack(ecode.Int(res.Code))
		return
	}
	return
}

func midsToParam(mids []int64) (str string) {
	strs := make([]string, 0, len(mids))
	for _, mid := range mids {
		strs = append(strs, fmt.Sprintf("%d", mid))
	}
	return strings.Join(strs, ",")
}

// TelInfo tel info.
func (d *Dao) TelInfo(c context.Context, mid int64) (tel string, err error) {
	params := url.Values{}
	params.Set("mid", fmt.Sprintf("%d", mid))
	var resp struct {
		Code int `json:"code"`
		Data struct {
			Mid      int64  `json:"mid"`
			Tel      string `json:"tel"`
			JoinIP   string `json:"join_ip"`
			JoinTime int64  `json:"join_time"`
		} `json:"data"`
	}
	if err = d.httpClient.Get(c, conf.Conf.Property.TelURL, "", params, &resp); err != nil {
		err = errors.Wrapf(err, "telinfo : %d", mid)
		return
	}
	if resp.Code != 0 {
		err = errors.Errorf("telinfo url(%s) res(%+v) err(%+v)", conf.Conf.Property.TelURL+"?"+params.Encode(), resp, ecode.Int(resp.Code))
		return
	}
	tel = resp.Data.Tel
	return
}

// MailInfo .
func (d *Dao) MailInfo(c context.Context, mid int64) (mail string, err error) {
	params := url.Values{}
	params.Set("mid", fmt.Sprintf("%d", mid))
	var resp struct {
		Code int    `json:"code"`
		Data string `json:"data"`
	}
	if err = d.httpClient.Get(c, conf.Conf.Property.MailURL, "", params, &resp); err != nil {
		err = errors.Wrapf(err, "mailinfo : %d", mid)
		return
	}
	if resp.Code != 0 {
		err = errors.Errorf("mailinfo url(%s) res(%+v) err(%+v)", conf.Conf.Property.MailURL+"?"+params.Encode(), resp, ecode.Int(resp.Code))
		return
	}
	mail = resp.Data
	return
}
