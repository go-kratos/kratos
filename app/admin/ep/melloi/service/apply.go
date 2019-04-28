package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/ep/melloi/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// QueryApply query apply by apply object
func (s *Service) QueryApply(qar *model.QueryApplyRequest) (*model.QueryApplyResponse, error) {
	return s.dao.QueryApply(&qar.Apply, qar.PageNum, qar.PageSize)
}

// QueryUserApplyList query user apply list
func (s *Service) QueryUserApplyList(userName string) ([]*model.Apply, error) {
	return s.dao.QueryUserApplyList(userName)
}

// CheckRunPermission check user run permission
func (s *Service) CheckRunPermission(userName string) (ret bool) {
	var (
		startTime int64
		err       error
		endTime   int64
		applyList []*model.Apply
	)
	currentTime := time.Now().Unix()

	//白名单的人，不需要做check，直接返回true
	if ExistsInSlice(userName, s.c.Melloi.Executor) {
		ret = true
		return
	}

	if applyList, err = s.dao.QueryUserApplyList(userName); err != nil {
		log.Error("s.dao.QueryUserApplyList err :(%v)", err)
		return
	}

	for _, apply := range applyList {
		if startTime, err = strconv.ParseInt(apply.StartTime, 10, 64); err != nil {
			return
		}
		if endTime, err = strconv.ParseInt(apply.EndTime, 10, 64); err != nil {
			return
		}
		if currentTime >= startTime && currentTime <= endTime {
			ret = true
			break
		}
	}
	return ret
}

// CheckRunTime check perf time
func (s *Service) CheckRunTime() (ret bool) {
	currentTime := time.Now()
	// 1:30 ~ 12:00
	if currentTime.Hour() >= 1 && currentTime.Hour() < 12 {
		ret = true
	}
	// 14:00 ~ 17:00
	if currentTime.Hour() >= 14 && currentTime.Hour() < 17 {
		ret = true
	}

	return
}

// UpdateApply update apply info
func (s *Service) UpdateApply(cookie string, apply *model.Apply) (err error) {
	var user *model.User
	if apply.ID == 0 {
		return ecode.MelloiApplyRequestErr
	}

	if user, err = s.QueryUser(apply.From); err != nil {
		log.Error("update apply query user error:（%v)", err)
		return
	}
	if user.ID != 0 {
		user.Accept = apply.Status
		if err = s.dao.UpdateUser(user); err != nil {
			return err
		}
		// 判断是审批操作，发送微信通知
		if apply.Status == 1 {
			applyMsg, _ := s.dao.QueryApplyByID(apply.ID)
			startTime, _ := strconv.Atoi(applyMsg.StartTime)
			st := time.Unix(int64(startTime), 0)
			endTime, _ := strconv.Atoi(applyMsg.EndTime)
			et := time.Unix(int64(endTime), 0)

			content := "[MELLOI]压测申请处理完成 通知 \n 压测服务: " + applyMsg.Path + "\n" + "压测时间：" + st.Format("2006-01-02 15:04:05") + "\n" + "压测结束时间：" + et.Format("2006-01-02 15:04:05") + "\n" + "申请人：" +
				applyMsg.From + "\n" + "审批人：" + applyMsg.To + "\n" + "审批时间：" + time.Now().Format("2006-01-02 15:04:05")
			log.Info("content:(%s)", content)
			// 发送申请通过到群
			go s.AddWechatSend(context.TODO(), cookie, content)
			// 给申请人发送邮件
			go s.SendMail(apply.From+"@bilibili.com", "[MELLOI]压测申请通知", content)
		}
		return s.dao.UpdateApply(apply)
	}
	return ecode.MelloiUpdateUserErr
}

// AddApply add new apply
func (s *Service) AddApply(c context.Context, cookie string, apply *model.Apply) (err error) {
	apply.Status = -1
	apply.Active = 1

	// 添加apply到db
	if err = s.dao.AddApply(apply); err != nil {
		return ecode.MelloiApplyRequestErr
	}

	// 发送微信消息 & 发送邮件通知
	//加密 apply.id|apply.from|apply.to
	applyID := strconv.FormatInt(apply.ID, 10)
	beStr := applyID + "|" + apply.From + "|" + apply.To
	base64Str := base64.StdEncoding.EncodeToString([]byte(beStr))
	// 将时间戳转成日期
	startTime, _ := strconv.Atoi(apply.StartTime)
	st := time.Unix(int64(startTime), 0)
	endTime, _ := strconv.Atoi(apply.EndTime)
	et := time.Unix(int64(endTime), 0)

	// 增加依赖服务列表
	var (
		userService map[string][]string
		serviceList = make(map[string][]string)
		serviceDep  string
		serviceName string
	)

	serviceName = strings.Replace(apply.Path, "bilibili.", "", 1)
	if userService, err = s.QueryDependServiceAdmins(c, serviceName, s.getSessionInCookie(cookie)); err != nil {
		log.Error("query depend service admin error(%v)", err)
		return
	}
	for _, v := range userService {
		for _, service := range v {
			serviceList[service] = nil
		}
	}
	for k := range serviceList {
		serviceDep += "\n" + k
	}

	// 拼接消息体，amd=base64Str
	content := fmt.Sprintf("[MELLOI]压测申请处理 通知 \n 压测服务:%s\n压测开始时间段:%s\n压测结束时间段:%s\n申请人:%s\n申请时间:%s\n依赖服务:%s\n审批地址:http://melloi.bilibili.co#/apply-m?platform=mb&amd=%s",
		apply.Path, st.Format("2006-01-02 15:04:05"), et.Format("2006-01-02 15:04:05"),
		apply.From, time.Now().Format("2006-01-02 15:04:05"), serviceDep, base64Str)
	// 消息接收人
	var touser []string
	touser = append(touser, apply.To)
	// 发送微信
	go s.dao.PushWechatMsgToPerson(context.TODO(), cookie, touser, content)
	// 发送邮件
	subject := "Melloi压测申请"
	go s.SendMail(apply.To+"@bilibili.com", subject, content)
	return
}

// DeleteApply delete apply
func (s *Service) DeleteApply(id int64) error {
	return s.dao.DeleteApply(id)
}
