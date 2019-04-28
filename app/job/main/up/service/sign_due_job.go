package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"time"

	//"go-common/app/job/main/up/conf"
	"go-common/app/admin/main/up/util"
	"go-common/app/job/main/up/dao/upcrm"
	"go-common/app/job/main/up/model/signmodel"
	"go-common/app/job/main/up/model/upcrmmodel"
	account "go-common/app/service/main/account/model"
	"go-common/library/log"

	"go-common/app/job/main/up/conf"
)

var (
	//ErrNoAdminName no admin name
	ErrNoAdminName = errors.New("no admin name")

	tmplSignDueTitle   *template.Template
	tmplSignDueContent *template.Template

	tmplPayDueTitle   *template.Template
	tmplPayDueContent *template.Template

	tmplTaskDueTitle   *template.Template
	tmplTaskDueContent *template.Template
)

// use for template function call
var funcHelper = template.FuncMap{
	"Now": time.Now,
}

func (s *Service) initEmailTemplate() (err error) {
	if conf.Conf.MailTemplateConf.SignTmplTitle == "" ||
		conf.Conf.MailTemplateConf.SignTmplContent == "" ||
		conf.Conf.MailTemplateConf.PayTmplTitle == "" ||
		conf.Conf.MailTemplateConf.PayTmplContent == "" ||
		conf.Conf.MailTemplateConf.TaskTmplTitle == "" ||
		conf.Conf.MailTemplateConf.TaskTmplContent == "" {
		err = fmt.Errorf(`mail template conf is invalid, check mail-template.toml file, make sure all the following has value:
		TaskTmplContent
		TaskTmplTitle
		PayTmplContent
		PayTmplTitle
		SignTmplContent
		SignTmplTitle`)
		return
	}
	tmplSignDueTitle, err = template.New("signTitle").Funcs(funcHelper).Parse(conf.Conf.MailTemplateConf.SignTmplTitle)
	if err != nil {
		log.Error("parse template fail, err=%v", err)
		return
	}

	tmplSignDueContent, err = template.New("signContent").Funcs(funcHelper).Parse(conf.Conf.MailTemplateConf.SignTmplContent)
	if err != nil {
		log.Error("parse template fail, err=%v", err)
		return
	}

	tmplPayDueTitle, err = template.New("payTitle").Funcs(funcHelper).Parse(conf.Conf.MailTemplateConf.PayTmplTitle)
	if err != nil {
		log.Error("parse template fail, err=%v", err)
		return
	}
	tmplPayDueContent, err = template.New("payContent").Funcs(funcHelper).Parse(conf.Conf.MailTemplateConf.PayTmplContent)
	if err != nil {
		log.Error("parse template fail, err=%v", err)
		return
	}
	tmplTaskDueTitle, err = template.New("taskTitle").Funcs(funcHelper).Parse(conf.Conf.MailTemplateConf.TaskTmplTitle)
	if err != nil {
		log.Error("parse template fail, err=%v", err)
		return
	}
	tmplTaskDueContent, err = template.New("taskContent").Funcs(funcHelper).Parse(conf.Conf.MailTemplateConf.TaskTmplContent)
	if err != nil {
		log.Error("parse template fail, err=%v", err)
		return
	}
	return
}

//CheckDateDueJob check task due
/*
	快到期的job提醒
*/
func (s *Service) CheckDateDueJob(date time.Time) {
	log.Info("start run CheckDateDueJob, date=%s", date)
	s.checkSignUpDue(date)
	log.Info("finish run CheckDateDueJob, date=%s", date)
}

type dueData struct {
	Signs []*upcrm.SignWithName
	Pays  []*upcrm.PayWithAdmin
	Tasks []*upcrm.TaskWithAdmin
}

func (d *dueData) addSign(sign *upcrm.SignWithName) {
	d.Signs = append(d.Signs, sign)
}

func (d *dueData) addPay(pay *upcrm.PayWithAdmin) {
	d.Pays = append(d.Pays, pay)
}

func (d *dueData) addTask(task *upcrm.TaskWithAdmin) {
	d.Tasks = append(d.Tasks, task)
}

func getOrCreate(dataMap map[string]*dueData, key string) *dueData {
	var data, ok = dataMap[key]
	if !ok {
		data = &dueData{}
		dataMap[key] = data
	}
	return data
}

func getName(infoMap map[int64]*account.Info, mid int64) string {
	if info, ok := infoMap[mid]; ok {
		return info.Name
	}
	return ""
}
func (s *Service) checkSignUpDue(date time.Time) {
	// 30天内到期 sign
	list, err := s.crmdb.GetDueSignUp(date, 30)
	s.crmdb.StartTask(upcrmmodel.TaskTypeSignCheckDue, date)

	defer func() {
		if err == nil {
			s.crmdb.FinishTask(upcrmmodel.TaskTypeSignCheckDue, date, upcrmmodel.TaskStateFinish)
		} else {
			s.crmdb.FinishTask(upcrmmodel.TaskTypeSignCheckDue, date, upcrmmodel.TaskStateError)
		}
	}()
	if err != nil {
		log.Error("fail to get due sign, date=%+v, err=%+v", date, err)
		return
	}

	var adminDueDataMap = make(map[string]*dueData)
	var ids []int64
	for _, v := range list {
		ids = append(ids, v.Mid)
		var data = getOrCreate(adminDueDataMap, v.AdminName)
		data.addSign(v)
	}

	// 7天内到期的pay
	listPayDue, err := s.crmdb.GetDuePay(date, 7)
	if err != nil {
		log.Error("fail to get due pay, date=%+v, err=%+v", date, err)
		return
	}
	for _, v := range listPayDue {
		ids = append(ids, v.Mid)

		var data = getOrCreate(adminDueDataMap, v.AdminName)
		data.addPay(v)
	}

	// 到期的任务，
	listTaskDue, err := s.crmdb.GetDueTask(date)
	if err != nil {
		log.Error("fail to get due task, date=%+v, err=%+v", date, err)
		return
	}

	for _, v := range listTaskDue {
		ids = append(ids, v.Mid)
		var data = getOrCreate(adminDueDataMap, v.AdminName)
		data.addTask(v)
	}
	ids = util.Unique(ids)
	infoMap, e := s.acc.GetCachedInfos(context.Background(), ids, "")

	if e == nil {
		for _, v := range list {
			v.Name = getName(infoMap, v.Mid)
		}

		for _, v := range listPayDue {
			v.Name = getName(infoMap, v.Mid)
		}

		for _, v := range listTaskDue {
			v.Name = getName(infoMap, v.Mid)
		}
	}

	for admin, v := range adminDueDataMap {
		var adminAll = append(conf.Conf.MailConf.DueMailReceivers, admin)
		// 发送sign到期邮件
		var due = v
		s.worker.Add(func() {
			var succIds []uint32
			if len(due.Signs) > 0 {
				var e = s.sendMailWithTemplate(due.Signs, tmplSignDueTitle, tmplSignDueContent, adminAll...)
				if e == nil {
					for _, data := range due.Signs {
						succIds = append(succIds, data.ID)
					}
				} else {
					log.Warn("fail to send email, err=%v", e)
				}

				// 更新邮件发送标记
				s.crmdb.UpdateEmailState(signmodel.TableNameSignUp, succIds, signmodel.EmailStateSendSucc)
			}

			if len(due.Pays) > 0 {
				// 发送pay到期邮件
				succIds = nil
				e = s.sendMailWithTemplate(due.Pays, tmplPayDueTitle, tmplPayDueContent, adminAll...)
				if e == nil {
					for _, data := range due.Pays {
						succIds = append(succIds, data.ID)
					}
				} else {
					log.Warn("fail to send email, err=%v", e)
				}

				// 更新邮件发送标记
				s.crmdb.UpdateEmailState(signmodel.TableNameSignPay, succIds, signmodel.EmailStateSendSucc)
			}

			if len(due.Tasks) > 0 {
				// 发送task到期邮件
				e = s.sendMailWithTemplate(due.Tasks, tmplTaskDueTitle, tmplTaskDueContent, adminAll...)
				if e != nil {
					log.Warn("fail to send email, err=%v", e)
				}
				// 这个没有邮件发送的标记
				//s.crmdb.UpdateEmailState(signmodel.TableNameSignPay, succIds, signmodel.EmailStateSendSucc)
			}
		})
	}
}

// data, data to generate email content
// contentTmpl, template to generate email content
// adminname, slice for all admin name
//
func (s *Service) sendMailWithTemplate(data interface{}, subjectTmpl, contentTmpl *template.Template, adminName ...string) (err error) {
	if contentTmpl == nil {
		err = fmt.Errorf("template for email is nil, data=%+v", data)
		log.Error("%s", err)
		return
	}

	var contentBuf = bytes.NewBuffer(nil)
	err = contentTmpl.Execute(contentBuf, data)
	if err != nil {
		log.Error("template fail to execute, err=%v", err)
		return
	}

	var subjectBuf = bytes.NewBuffer(nil)
	err = subjectTmpl.Execute(subjectBuf, data)
	if err != nil {
		log.Error("template fail to execute, err=%v", err)
		return
	}

	var addrs []string
	for _, v := range adminName {
		if v == "" {
			log.Warn("admin name is empty")
			continue
		}
		var addr = fmt.Sprintf("%s@bilibili.com", v)
		addrs = append(addrs, addr)
	}
	//log.Info("email sub=%s, content=%s", subjectBuf.String(),  contentBuf.String())

	if len(addrs) == 0 {
		log.Error("admin name is empty, cannot send email, data=%+v", data)
		err = ErrNoAdminName
		return
	}
	log.Info("email send , sub=%s, admin=%s, data=%+v", subjectBuf.String(), adminName, data)
	err = s.maildao.SendMail(contentBuf.String(), subjectBuf.String(), addrs)
	return
}
