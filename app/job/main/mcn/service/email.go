package service

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"time"

	"go-common/app/job/main/mcn/conf"
	"go-common/app/job/main/mcn/model"
	accgrpc "go-common/app/service/main/account/api"
	"go-common/library/log"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

// var .
var (
	// ErrNoAdminName no admin name
	ErrNoAdminName = errors.New("no admin name")

	tmplSignDueTitle   *template.Template
	tmplSignDueContent *template.Template

	tmplPayDueTitle   *template.Template
	tmplPayDueContent *template.Template
)

// use for template function call
var funcHelper = template.FuncMap{
	"Now": time.Now,
}

func (s *Service) initEmailTemplate() (err error) {
	if conf.Conf.MailTemplateConf.SignTmplTitle == "" ||
		conf.Conf.MailTemplateConf.SignTmplContent == "" ||
		conf.Conf.MailTemplateConf.PayTmplTitle == "" ||
		conf.Conf.MailTemplateConf.PayTmplContent == "" {
		err = fmt.Errorf(`mail template conf is invalid, check mail-template.toml file, make sure all the following has value:
		TaskTmplContent
		TaskTmplTitle
		PayTmplContent
		PayTmplTitle`)
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
	return
}

// CheckDateDueCron .
func (s *Service) CheckDateDueCron() {
	log.Info("start run CheckDateDueJob, date=%s", time.Now().Format(model.TimeFormatSec))
	s.checkSignUpDue()
	log.Info("finish run CheckDateDueJob, date=%s", time.Now().Format(model.TimeFormatSec))
}

type stateFunc func(context.Context, []int64) (int64, error)

type emailData struct {
	IDs            []int64
	AdminName      []string
	Data           interface{}
	UpStateFunc    stateFunc
	Title, Content *template.Template
}

func (e *emailData) addEmailDatas(es *[]*emailData) {
	*es = append(*es, e)
}

// buildEmail .
func buildEmail(ids []int64, data interface{}, title, content *template.Template, upStateFunc stateFunc, adminName ...string) *emailData {
	return &emailData{
		IDs:         ids,
		Data:        data,
		Title:       title,
		Content:     content,
		UpStateFunc: upStateFunc,
		AdminName:   adminName,
	}
}

type dueData struct {
	Signs      []*model.MCNSignInfo
	Pays       []*model.SignPayInfo
	Sids, Pids []int64
}

func (d *dueData) addSign(sign *model.MCNSignInfo) {
	d.Signs = append(d.Signs, sign)
	d.Sids = append(d.Sids, sign.SignID)
}

func (d *dueData) addPay(pay *model.SignPayInfo) {
	d.Pays = append(d.Pays, pay)
	d.Pids = append(d.Pids, pay.SignPayID)
}

func (d *dueData) addName(infoMap map[int64]*accgrpc.Info) {
	for _, v := range d.Signs {
		v.McnName = getName(infoMap, v.McnMid)
	}
	for _, v := range d.Pays {
		v.McnName = getName(infoMap, v.McnMid)
	}
}

// func getOrCreate(dataMap map[string]*dueData, key string) *dueData {
// 	var data, ok = dataMap[key]
// 	if !ok {
// 		data = &dueData{}
// 		dataMap[key] = data
// 	}
// 	return data
// }

func getName(infoMap map[int64]*accgrpc.Info, mid int64) string {
	if info, ok := infoMap[mid]; ok {
		return info.Name
	}
	return ""
}

func (s *Service) checkSignUpDue() {
	var (
		mids       []int64
		emailDatas []*emailData
		data       = &dueData{}
		c          = context.Background()
		infoMap    map[int64]*accgrpc.Info
	)
	// 30天内到期 sign
	listDue, err := s.dao.McnSignDues(c)
	if err != nil {
		log.Error("s.dao.McnSignDues error(%+v)", err)
		return
	}
	for _, v := range listDue {
		mids = append(mids, v.McnMid)
		data.addSign(v)
	}
	// 7天内到期的pay
	listPayDue, err := s.dao.McnSignPayDues(c)
	if err != nil {
		log.Error("s.dao.McnSignPayDues  error(%+v)", err)
		return
	}
	for _, v := range listPayDue {
		mids = append(mids, v.McnMid)
		data.addPay(v)
	}
	mids = uniqNoEmpty(mids)
	infosReply, err := s.accGRPC.Infos3(c, &accgrpc.MidsReq{Mids: mids})
	if err != nil {
		log.Error("s.accGRPC.Infos3(%s) error(%+v)", xstr.JoinInts(mids), err)
		err = nil
	} else {
		infoMap = infosReply.Infos
	}
	emailDatas = make([]*emailData, 0)
	data.addName(infoMap)
	buildEmail(data.Sids, data.Signs, tmplSignDueTitle, tmplSignDueContent, s.dao.UpMcnSignEmailState, conf.Conf.MailConf.DueAuthorityGroups...).addEmailDatas(&emailDatas)
	buildEmail(data.Pids, data.Pays, tmplPayDueTitle, tmplPayDueContent, s.dao.UpMcnSignPayEmailState, conf.Conf.MailConf.DueAuthorityGroups...).addEmailDatas(&emailDatas)
	for _, e := range emailDatas {
		s.doSendEmailFunc(c, e)
	}
}

func (s *Service) doSendEmailFunc(c context.Context, e *emailData) {
	s.worker.Do(c, func(c context.Context) {
		if len(e.IDs) == 0 {
			log.Warn("not need to update")
			return
		}
		var err error
		if err = s.sendMailWithTemplate(e.Data, e.Title, e.Content, e.AdminName...); err != nil {
			log.Error("s.sendMailWithTemplate(%+v,%+v,%+v,%+v) error(%+v)", e.Data, e.Title, e.Content, e.AdminName, err)
			return
		}
		if _, err = e.UpStateFunc(c, e.IDs); err != nil {
			log.Error("upfunc(%+v,%s) error(%+v)", e.UpStateFunc, xstr.JoinInts(e.IDs), err)
			return
		}
		log.Info("func(%s) update succ", xstr.JoinInts(e.IDs))
	})
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
		addrs = append(addrs, v)
	}

	if len(addrs) == 0 {
		log.Error("admin name is empty, cannot send email, data=%+v", data)
		err = ErrNoAdminName
		return
	}
	if err = s.dao.SendMail(contentBuf.String(), subjectBuf.String(), addrs); err != nil {
		log.Error("s.dao.SendMail(%s,%s,%+v) error(%+v)", contentBuf.String(), subjectBuf.String(), addrs, err)
		return
	}
	log.Info("email send succ, sub=%s, admin=%s", subjectBuf.String(), adminName)
	return
}

func chain(ids ...[]int64) []int64 {
	res := make([]int64, 0, len(ids))
	for _, l := range ids {
		res = append(res, l...)
	}
	return res
}

// func uniq(ids ...[]int64) []int64 {
// 	hm := make(map[int64]struct{})
// 	for _, i := range chain(ids...) {
// 		hm[i] = struct{}{}
// 	}
// 	res := make([]int64, 0, len(ids))
// 	for i := range hm {
// 		res = append(res, i)
// 	}
// 	return res
// }

func uniqNoEmpty(ids ...[]int64) []int64 {
	hm := make(map[int64]struct{})
	for _, i := range chain(ids...) {
		hm[i] = struct{}{}
	}
	res := make([]int64, 0, len(ids))
	for i := range hm {
		if i > 0 {
			res = append(res, i)
		}
	}
	return res
}
