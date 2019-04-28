package service

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"go-common/app/admin/main/push/model"
	pushmdl "go-common/app/service/main/push/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_testPushURL  = "http://api.bilibili.co/x/internal/push-service/push"
	_testTokenURL = "http://api.bilibili.co/x/internal/push-service/test/token"
)

// AddTask add task.
func (s *Service) AddTask(c context.Context, task *model.Task) (err error) {
	if err = s.DB.Create(task).Error; err != nil {
		s.dao.SendWechat("推送后台新建任务失败，原因：写入DB失败")
		log.Error("s.AddTask(%+v) error(%v)", task, err)
	}
	return
}

// TestPushMid test push by mid.
func (s *Service) TestPushMid(c context.Context, mid string, task *model.Task) (err error) {
	params := url.Values{}
	params.Set("app_id", fmt.Sprintf("%d", task.AppID))
	params.Set("platform", task.Platform)
	params.Set("alert_title", task.Title)
	params.Set("alert_body", task.Summary)
	params.Set("link_type", strconv.Itoa(task.LinkType))
	params.Set("link_value", task.LinkValue)
	params.Set("expire_time", fmt.Sprintf("%d", task.ExpireTimeUnix))
	params.Set("builds", task.Build)
	params.Set("sound", strconv.Itoa(task.Sound))
	params.Set("vibration", strconv.Itoa(task.Vibration))
	params.Set("mid", mid)
	params.Set("image_url", task.ImageURL)
	var res struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err = s.httpClient.Post(c, _testPushURL, "", params, &res); err != nil {
		log.Error("s.TestPush(%s) httpClient.Get(%s,%v) error(%v)", mid, _testPushURL, params, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("测试推送失败 mid(%s) code(%d)", mid, res.Code)
		err = ecode.Int(res.Code)
	}
	return
}

// TestPushToken test push by token.
func (s *Service) TestPushToken(c context.Context, info *pushmdl.PushInfo, platform int, token string) (err error) {
	params := url.Values{}
	params.Set("app_id", fmt.Sprintf("%d", info.APPID))
	params.Set("platform", strconv.Itoa(platform))
	params.Set("alert_title", info.Title)
	params.Set("alert_body", info.Summary)
	params.Set("link_type", fmt.Sprintf("%d", info.LinkType))
	params.Set("link_value", info.LinkValue)
	params.Set("expire_time", fmt.Sprintf("%d", int64(info.ExpireTime)))
	params.Set("sound", strconv.Itoa(info.Sound))
	params.Set("vibration", strconv.Itoa(info.Vibration))
	params.Set("pass_through", strconv.Itoa(info.PassThrough))
	params.Set("token", token)
	var res struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err = s.httpClient.Post(c, _testTokenURL, "", params, &res); err != nil {
		log.Error("s.TestToken(%s) httpClient.Get(%s) error(%v)", token, params.Encode(), err)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("测试推送失败 token(%s) code(%d)", token, res.Code)
		err = ecode.Int(res.Code)
	}
	return
}

// 上传说明：
// 前端是批量上传，会随机按内容长度切割文件进行分批上传，有可能会切断原始行内容
// 如果想在上传的时候同步判断文件格式，产生错误时需要忽略首行和末行

// CheckUploadMid checks uploaded mid validation.
func (s *Service) CheckUploadMid(c context.Context, data []byte) (err error) {
	var (
		mid     int64
		lineNum int
		lines   = strings.Split(string(data), "\n")
		total   = len(lines)
	)
	for _, v := range lines {
		lineNum++
		v = strings.Trim(v, " \r\t")
		if v == "" {
			continue
		}
		if mid, err = strconv.ParseInt(v, 10, 64); err != nil {
			log.Error("CheckUploadMid data(%s) error(%v)", v, err)
			return ecode.PushUploadInvalidErr
		}
		if mid <= 0 {
			if lineNum == 1 || lineNum == total {
				continue
			}
			log.Error("CheckUploadMid data(%s) error(%v)", v, err)
			return ecode.PushUploadInvalidErr
		}
	}
	return
}

// CheckUploadToken checks uploaded token validation.
func (s *Service) CheckUploadToken(c context.Context, data []byte) (err error) {
	var (
		plat    int
		lineNum int
		lines   = strings.Split(string(data), "\n")
		total   = len(lines)
	)
	for _, v := range lines {
		lineNum++
		if lineNum == 1 || lineNum == total {
			continue
		}
		v = strings.Trim(v, " \r")
		if v == "" {
			continue
		}
		res := strings.Split(v, "\t")
		if len(res) != 2 {
			log.Error("CheckUploadToken data(%s)", v)
			return ecode.PushUploadInvalidErr
		}
		if res[0] == "" || res[1] == "" {
			log.Error("CheckUploadToken data(%s)", v)
			return ecode.PushUploadInvalidErr
		}
		if plat, err = strconv.Atoi(res[0]); err != nil || plat <= 0 {
			log.Error("CheckUploadToken data(%s) error(%v)", v, err)
			return ecode.PushUploadInvalidErr
		}
	}
	return
}
