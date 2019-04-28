package mail

import (
	"bytes"
	"fmt"
	"text/template"

	"go-common/app/tool/saga/conf"
	"go-common/app/tool/saga/model"
	"go-common/library/log"

	"github.com/pkg/errors"
	gomail "gopkg.in/gomail.v2"
)

// SendMail2 ...
func SendMail2(addr *model.MailAddress, subject string, data string) (err error) {
	var (
		msg = gomail.NewMessage()
	)
	msg.SetAddressHeader("From", conf.Conf.Property.Mail.Address, conf.Conf.Property.Mail.Name)
	msg.SetHeader("To", msg.FormatAddress(addr.Address, addr.Name))
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", data)
	d := gomail.NewDialer(
		conf.Conf.Property.Mail.Host,
		conf.Conf.Property.Mail.Port,
		conf.Conf.Property.Mail.Address,
		conf.Conf.Property.Mail.Pwd,
	)
	if err = d.DialAndSend(msg); err != nil {
		err = errors.WithMessage(err, fmt.Sprintf("Send mail (%+v,%s,%s) failed", addr, subject, data))
		return
	}
	log.Info("Send mail  (%+v,%s,%s) success", addr, subject, data)
	return
}

// SendMail  send mail
func SendMail(m *model.Mail, data *model.MailData) (err error) {
	var (
		toUsers []string
		msg     = gomail.NewMessage()
		buf     = &bytes.Buffer{}
	)
	msg.SetAddressHeader("From", conf.Conf.Property.Mail.Address, conf.Conf.Property.Mail.Name) // 发件人
	for _, ads := range m.ToAddress {
		toUsers = append(toUsers, msg.FormatAddress(ads.Address, ads.Name))
	}
	t := template.New("MR Mail")
	if t, err = t.Parse(mailTPL); err != nil {
		log.Error("tpl.Parse(%s) error(%+v)", mailTPL, errors.WithStack(err))
		return
	}
	err = t.Execute(buf, data)
	if err != nil {
		log.Error("t.Execute error(%+v)", errors.WithStack(err))
		return
	}
	msg.SetHeader("To", toUsers...)
	msg.SetHeader("Subject", m.Subject)    // 主题
	msg.SetBody("text/html", buf.String()) // 正文
	d := gomail.NewDialer(
		conf.Conf.Property.Mail.Host,
		conf.Conf.Property.Mail.Port,
		conf.Conf.Property.Mail.Address,
		conf.Conf.Property.Mail.Pwd,
	)
	if err = d.DialAndSend(msg); err != nil {
		log.Error("Send mail Fail(%v) diff(%s)", msg, err)
		return
	}
	return
}

// SendMail3 SendMail all parameter
func SendMail3(from string, sender string, senderPwd string, m *model.Mail, data *model.MailData) (err error) {
	var (
		toUsers []string
		msg     = gomail.NewMessage()
		buf     = &bytes.Buffer{}
	)
	msg.SetAddressHeader("From", from, sender) // 发件人
	for _, ads := range m.ToAddress {
		toUsers = append(toUsers, msg.FormatAddress(ads.Address, ads.Name))
	}
	t := template.New("MR Mail")
	if t, err = t.Parse(mailTPL3); err != nil {
		log.Error("tpl.Parse(%s) error(%+v)", mailTPL3, errors.WithStack(err))
		return
	}
	err = t.Execute(buf, data)
	if err != nil {
		log.Error("t.Execute error(%+v)", errors.WithStack(err))
		return
	}
	msg.SetHeader("To", toUsers...)
	msg.SetHeader("Subject", m.Subject)    // 主题
	msg.SetBody("text/html", buf.String()) // 正文
	d := gomail.NewDialer(
		"smtp.exmail.qq.com",
		465,
		from,
		senderPwd,
	)
	if err = d.DialAndSend(msg); err != nil {
		log.Error("Send mail Fail(%v) diff(%s)", msg, err)
		return
	}
	return
}
