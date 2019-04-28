package service

import (
	"context"
	"fmt"
	"time"

	"go-common/app/admin/ep/merlin/model"

	"gopkg.in/gomail.v2"
)

const (
	_merlinUIAddr = "http://merlin.bilibili.co"
)

// SendMail send mail.
func (s *Service) SendMail(mailType int, machine *model.Machine) (err error) {
	var (
		mailSendHead    = ""
		mailSendContext = ""
		m               = gomail.NewMessage()
		user            *model.User
	)

	switch mailType {
	//将要过期提醒
	case model.MailTypeMachineWillExpired:
		if user, err = s.dao.FindUserByUserName(machine.Username); err != nil {
			return
		}
		m.SetHeader("To", user.EMail)
		mailSendHead = fmt.Sprintf("machine named [%s] will be expired on %s", machine.Name, machine.EndTime.Format(model.TimeFormat))
		m.SetHeader("Subject", mailSendHead)

		delayLink := fmt.Sprintf("http://merlin.bilibili.co/#/machine-list?machine_name=%s&username=&page_num=1&page_size=10", machine.Name)
		mailSendContext = fmt.Sprintf("<br>可前往Merlin平台申请延期</br><br>链接: <a href=%s>点击</a></br>", delayLink)

		//机器删除提醒
	case model.MailTypeMachineDeleted:
		if user, err = s.dao.FindUserByUserName(machine.Username); err != nil {
			return
		}
		m.SetHeader("To", user.EMail)

		mailSendHead = fmt.Sprintf("machine named [%s] has been deleted on %s", machine.Name, time.Now().Format(model.TimeFormat))
		m.SetHeader("Subject", mailSendHead)

		//机器删除失败提醒
	case model.MailTypeTaskDeleteMachineFailed:
		m.SetHeader("To", s.c.Mail.NoticeOwner[0])

		mailSendHead = fmt.Sprintf("machine named [%s] has been failed to delete by task on %s", machine.Name, time.Now().Format(model.TimeFormat))
		m.SetHeader("Subject", mailSendHead)

	}

	if mailSendHead != "" {
		ml := &model.MailLog{
			ReceiverName: m.GetHeader("To")[0],
			MailType:     mailType,
			SendHead:     mailSendHead,
			SendContext:  mailSendContext,
		}
		s.dao.SendMail(m)
		s.dao.InsertMailLog(ml)
	}
	return
}

// SendMailDeleteMachine Send Mail Delete Machine.
func (s *Service) SendMailDeleteMachine(username string, machine *model.Machine) (err error) {
	var (
		mailSendHead = ""
		m            = gomail.NewMessage()
		user         *model.User
	)

	if user, err = s.dao.FindUserByUserName(machine.Username); err != nil {
		return
	}

	//是不是删自己机器
	if username == machine.Username {
		m.SetHeader("To", user.EMail)
		mailSendHead = fmt.Sprintf("Machine named [%s] has been deleted on %s", machine.Name, time.Now().Format(model.TimeFormat))
		m.SetHeader("Subject", mailSendHead)

	} else {
		var delUser *model.User
		if delUser, err = s.dao.FindUserByUserName(username); err != nil {
			m.SetHeader("To", user.EMail)
		} else {
			m.SetHeader("To", user.EMail, delUser.EMail)
		}

		mailSendHead = fmt.Sprintf("Machine named [%s] has been deleted on %s by %s", machine.Name, time.Now().Format(model.TimeFormat), username)
		m.SetHeader("Subject", mailSendHead)
	}

	for _, header := range m.GetHeader("To") {
		ml := &model.MailLog{
			ReceiverName: header,
			MailType:     model.MailTypeMachineDeleted,
			SendHead:     mailSendHead,
			SendContext:  "",
		}
		s.dao.SendMail(m)
		s.dao.InsertMailLog(ml)
	}

	return
}

// SendMailApplyDelayMachineEndTime send mail apply delay machine end time.
func (s *Service) SendMailApplyDelayMachineEndTime(ctx context.Context, auditor, applicant string, machineID int64, currentEndTime, applyEndTime time.Time) (err error) {
	var (
		userInfo *model.User
		machine  *model.Machine
	)

	if userInfo, err = s.QueryUserInfo(ctx, auditor); err != nil {
		return
	}

	if machine, err = s.dao.QueryMachine(machineID); err != nil {
		return
	}

	timeNow := time.Now().Format(model.TimeFormat)
	currentET := currentEndTime.Format(model.TimeFormat)
	applyET := applyEndTime.Format(model.TimeFormat)
	mailSendHead, mailSendContext := applyDelayMachineEndTimeHeadAndContext(machine.Name, applicant, timeNow, currentET, applyET, _merlinUIAddr)

	err = s.sendMail(userInfo.EMail, mailSendHead, mailSendContext, model.MailTypeApplyDelayMachineEndTime)
	return
}

// SendMailAuditResult send mail audit result.
func (s *Service) SendMailAuditResult(ctx context.Context, auditor, applicant string, machineID int64, auditResult bool) (err error) {
	var (
		userInfo *model.User
		machine  *model.Machine
	)

	if userInfo, err = s.QueryUserInfo(ctx, applicant); err != nil {
		return
	}

	if machine, err = s.dao.QueryMachine(machineID); err != nil {
		return
	}

	timeNow := time.Now().Format(model.TimeFormat)
	mailSendHead, mailSendContext := auditResultHeadAndContext(machine.Name, timeNow, auditor, _merlinUIAddr, auditResult)

	err = s.sendMail(userInfo.EMail, mailSendHead, mailSendContext, model.MailTypeAuditDelayMachineEndTime)
	return
}

func (s *Service) sendMail(receiver, header, body string, mailType int) (err error) {
	m := gomail.NewMessage()
	m.SetHeader("To", receiver)
	m.SetHeader("Subject", header)
	m.SetBody("text/html", body)

	s.dao.SendMail(m)

	ml := &model.MailLog{
		ReceiverName: receiver,
		MailType:     mailType,
		SendHead:     header,
		SendContext:  body,
	}
	err = s.dao.InsertMailLog(ml)
	return
}

func applyDelayMachineEndTimeHeadAndContext(machineName, applicant, timeNow, currentEndTime, applyEndTime, operateLink string) (head, context string) {
	head = fmt.Sprintf("申请延期机器 机器名:[%s] %s", machineName, timeNow)
	context = fmt.Sprintf("<br>申请延期机器: %s</br>"+
		"<br>申请者: %s</br>"+
		"<br>当前过期时间: %s</br>"+
		"<br>申请延期时间: %s</br>"+
		"<br>操作链接: <a href=%s>点击</a></br>", machineName, applicant, currentEndTime, applyEndTime, operateLink)
	return
}

func auditResultHeadAndContext(machineName, timeNow, auditor, operateLink string, auditResult bool) (head, context string) {
	auditResultStr := "驳回"
	if auditResult {
		auditResultStr = "通过"
	}

	head = fmt.Sprintf("申请延期机器审批结果 机器名:[%s] %s", machineName, timeNow)
	context = fmt.Sprintf("<br>申请延期机器: %s</br>"+
		"<br>审批者: %s</br>"+
		"<br>申请结果: %s</br>"+
		"<br>查看链接: <a href=%s>点击</a></br>", machineName, auditor, auditResultStr, operateLink)
	return
}

// SendMailForMultiUsers Send Mail For Multi Users.
func (s *Service) SendMailForMultiUsers(ctx context.Context, receivers []string, mailSendHead string) (err error) {
	for _, receiver := range receivers {
		var userInfo *model.User
		if userInfo, err = s.QueryUserInfo(ctx, receiver); err != nil {
			continue
		}
		err = s.sendMail(userInfo.EMail, mailSendHead, "", model.MailTypeMachineTransfer)
	}
	return
}
