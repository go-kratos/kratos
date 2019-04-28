package service

import (
	"context"
	"strconv"
	"time"

	"go-common/app/admin/ep/melloi/model"
)

//AddWechatSend  add wechat send
func (s *Service) AddWechatSend(c context.Context, cookie, content string) (msgSendRes *model.MsgSendRes, err error) {
	return s.dao.AddWechatSend(c, cookie, content)
}

// AddWechatContent Add Wechat Content
func AddWechatContent(ptestParam model.DoPtestParam, reportSuID int, jobName string, userService map[string][]string) (content string) {
	var (
		url            string
		lay            = "2006-01-02 15:04:05"
		ptestDetailURL string
		serviceList    = make(map[string][]string)
		serviceDep     string
		serviceName    string
	)
	if ptestParam.Type == model.PROTOCOL_HTTP || ptestParam.Type == model.PROTOCOL_SCENE {
		ptestDetailURL = "http://melloi.bilibili.co/#/ptest-detail?reportSuId=" + strconv.Itoa(reportSuID)
	}
	if ptestParam.Type == model.PROTOCOL_GRPC {
		ptestDetailURL = "http://melloi.bilibili.co/#/ptest-detail-grpc?reportSuId=" + strconv.Itoa(reportSuID)
	}
	url = ptestParam.URL
	if ptestParam.Type == model.PROTOCOL_SCENE {
		for _, script := range ptestParam.Scripts {
			url = url + "\n" + script.URL
		}
	}

	// 增加依赖服务列表
	for _, v := range userService {
		for _, service := range v {
			serviceList[service] = nil
		}
	}
	for k := range serviceList {
		serviceDep += "\n" + k
	}
	loadTime := strconv.Itoa(ptestParam.LoadTime) + "s"
	if ptestParam.Upload {
		loadTime = "脚本用户上传，时间1800s以内"
		url = "脚本用户上传，url 未知"
	}
	serviceName = ptestParam.Department + "." + ptestParam.Project + "." + ptestParam.APP
	content = "执行人：" + ptestParam.UserName + "\n压测服务：" + serviceName + "\n" + "压测接口：" + url + "\n开始时间：" + time.Now().Format(lay) + "\n持续时间：" +
		loadTime + "\n压测容器：" + jobName + "\n报告地址：" + ptestDetailURL + "\n压测依赖服务：" + serviceDep
	return
}

// AddWechatDependServiceContent add wechat depend Service Content
func AddWechatDependServiceContent(ptestParam model.DoPtestParam, userService map[string][]string, reportSuId int, user string) (content string) {
	var (
		url            string
		lay            = "2006-01-02 15:04:05"
		ptestDetailURL string
		serviceList    string
	)
	if ptestParam.Type == model.PROTOCOL_HTTP || ptestParam.Type == model.PROTOCOL_SCENE {
		ptestDetailURL = "http://melloi.bilibili.co/#/ptest-detail?reportSuId=" + strconv.Itoa(reportSuId)
	}

	if ptestParam.Type == model.PROTOCOL_GRPC {
		ptestDetailURL = "http://melloi.bilibili.co/#/ptest-detail-grpc?reportSuId=" + strconv.Itoa(reportSuId)
	}

	url = ptestParam.URL
	if ptestParam.Type == model.PROTOCOL_SCENE {
		for _, script := range ptestParam.Scripts {
			url = url + "\n" + script.URL
		}
	}

	for _, service := range userService[user] {
		serviceList += "\n" + service
	}

	serviceName := ptestParam.Department + "." + ptestParam.Project + "." + ptestParam.APP
	content = "[Melloi压测依赖提醒] \n 压测服务:" + serviceName + "\n 压测接口:" + ptestParam.URL + "\n 压测时间：" + time.Now().Format(lay) + "\n 压测时长: " +
		strconv.Itoa(ptestParam.LoadTime) + "\n 报告地址：" + ptestDetailURL + "\n 依赖服务:" + serviceList
	return
}
