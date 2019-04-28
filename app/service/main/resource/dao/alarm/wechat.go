package alarm

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go-common/app/service/main/resource/model"
	"go-common/library/log"
)

var (
	_states = map[int8]string{
		int8(0):  "开放浏览",
		int8(-1): "待审",
		int8(-2): "打回稿件回收站",
		int8(-3): "网警锁定删除",
		int8(-4): "锁定稿件",
		// -5:   "锁定稿件开放浏览",
		int8(-6): "修复待审",
		int8(-7): "暂缓审核",
		// -8:   "补档待审",
		int8(-9):  "等待转码",
		int8(-10): "延迟发布",
		int8(-11): "视频源待修",
		// -12:  "上传失败",
		int8(-13): "允许评论待审",
		// -14:  "临时回收站",
		int8(-15):  "分发中",
		int8(-16):  "转码失败",
		int8(-30):  "创建已提交",
		int8(-40):  "UP主定时发布",
		int8(-100): "UP主删除",
	}
	_codes = map[int]string{
		404: "页面未找到",
		502: "服务端异常",
		504: "SLB超时",
	}
)

func stateDescribe(state int8) string {
	des, ok := _states[state]
	if ok {
		return des
	}
	return strconv.Itoa(int(state))
}

const (
	_warnTitle = `
	【线上投放告警】您有一个投放发生了异常，请尽快确认是否影响投放。

	%v
以上内容确认后，如有异常请联系相关人员手动处理。`

	_offLineTitle = `
	【线上投放下线告警】您有一个投放内容变不可见状态，已自动下线。请尽快补充投放。

	%v
以上投放所涉及的排期申请已经置为 <申请未投放>。`

	_archiveContent = `
	投放ID：%d
	稿件ID：%d
	投放标题：%s
	投放位置：%s (位置ID：%d)
	投放时间段：%s - %s 
	当前稿件状态：%s

	`

	_urlContent = `
	投放ID：%d
	URL：%v
	投放标题：%s
	投放位置：%s (位置ID：%d)
	投放时间段：%s - %s 
	错误信息：%v(错误码: %v)

	`
)

type wxParams struct {
	Username  string `json:"username"`
	Content   string `json:"content"`
	Token     string `json:"token"`
	Timestamp int64  `json:"timestamp"`
	Sign      string `json:"signature"`
}
type resp struct {
	Status int64  `json:"status"`
	Msg    string `json:"msg"`
}

// SendWeChart send message to QYWX
func (d *Dao) SendWeChart(c context.Context, ns int8, userName string, res []*model.ResWarnInfo, titleType string) (err error) {
	var (
		users       = append(d.c.WeChantUsers, userName)
		newStateStr = stateDescribe(ns)
		contents    []string
	)
	for _, re := range res {
		stime := re.STime.Time().Format("2006-01-02 15:04:05")
		etime := re.ETime.Time().Format("2006-01-02 15:04:05")
		contents = append(contents, fmt.Sprintf(_archiveContent, re.AssignmentID, re.AID, re.AssignmentName, re.ResourceName, re.ResourceID, stime, etime, newStateStr))
	}
	params := url.Values{}
	params.Set("username", strings.Join(users, ","))
	if titleType == "warn" {
		params.Set("content", fmt.Sprintf(_warnTitle, strings.Join(contents, "")))
	} else {
		params.Set("content", fmt.Sprintf(_offLineTitle, strings.Join(contents, "")))
	}
	params.Set("token", d.c.WeChatToken)
	params.Set("timestamp", strconv.FormatInt(time.Now().Unix(), 10))
	mh := md5.Sum([]byte(params.Encode() + d.c.WeChatSecret))
	params.Set("signature", hex.EncodeToString(mh[:]))
	p := &wxParams{
		Username: params.Get("username"),
		Content:  params.Get("content"),
		Token:    params.Get("token"),
		Sign:     params.Get("signature"),
	}
	p.Timestamp, _ = strconv.ParseInt(params.Get("timestamp"), 10, 64)
	bs, _ := json.Marshal(p)
	payload := strings.NewReader(string(bs))
	req, _ := http.NewRequest("POST", "http://bap.bilibili.co/api/v1/message/add", payload)
	req.Header.Add("content-type", "application/json; charset=utf-8")
	v := &resp{}
	if err = d.httpClient.Do(context.TODO(), req, v); err != nil {
		log.Error("s.httpClient.Do error(%v)", err)
	}
	return
}

// sendWeChartURL send message to QYWX
func (d *Dao) sendWeChartURL(c context.Context, code int, userName string, res []*model.ResWarnInfo) (err error) {
	var (
		users    = append(d.c.WeChantUsers, userName)
		codeInfo string
		contents []string
	)
	if codeInfo = _codes[code]; codeInfo == "" {
		codeInfo = "未知错误，请手动确认"
	}
	for _, re := range res {
		stime := re.STime.Time().Format("2006-01-02 15:04:05")
		etime := re.ETime.Time().Format("2006-01-02 15:04:05")
		contents = append(contents, fmt.Sprintf(_urlContent, re.AssignmentID, re.URL, re.AssignmentName, re.ResourceName, re.ResourceID, stime, etime, codeInfo, code))
	}
	params := url.Values{}
	params.Set("username", strings.Join(users, ","))
	params.Set("content", fmt.Sprintf(_warnTitle, strings.Join(contents, "")))
	params.Set("token", d.c.WeChatToken)
	params.Set("timestamp", strconv.FormatInt(time.Now().Unix(), 10))
	mh := md5.Sum([]byte(params.Encode() + d.c.WeChatSecret))
	params.Set("signature", hex.EncodeToString(mh[:]))
	p := &wxParams{
		Username: params.Get("username"),
		Content:  params.Get("content"),
		Token:    params.Get("token"),
		Sign:     params.Get("signature"),
	}
	p.Timestamp, _ = strconv.ParseInt(params.Get("timestamp"), 10, 64)
	bs, _ := json.Marshal(p)
	payload := strings.NewReader(string(bs))
	req, _ := http.NewRequest("POST", "http://bap.bilibili.co/api/v1/message/add", payload)
	req.Header.Add("content-type", "application/json; charset=utf-8")
	v := &resp{}
	if err = d.httpClient.Do(context.TODO(), req, v); err != nil {
		log.Error("s.httpClient.Do error(%v)", err)
	}
	return
}
