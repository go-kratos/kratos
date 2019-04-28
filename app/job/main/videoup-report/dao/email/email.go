package email

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go-common/app/job/main/videoup-report/model/email"
	"go-common/library/log"
	"go-common/library/queue/databus/report"
	gomail "gopkg.in/gomail.v2"
)

//SendMail send the email
func (d *Dao) SendMail(tpl *email.Template) {
	var err error
	headers := tpl.Headers
	if len(headers[email.FROM]) == 0 || len(headers[email.TO]) == 0 || len(headers[email.SUBJECT]) == 0 {
		log.Error("email lack From/To/Subject: emailTemplate(%+v)", *tpl)
		return
	}
	if len(tpl.ContentType) == 0 {
		tpl.ContentType = "text/plain"
	}

	log.Info("start send mail: emailTemplate(%+v)", *tpl)
	msg := gomail.NewMessage()
	msg.SetHeaders(headers)
	msg.SetBody(tpl.ContentType, tpl.Body)
	result := email.EmailResOK

	if err = d.email.DialAndSend(msg); err != nil {
		result = email.EmailResFail
		log.Error("s.email.DialAndSend error(%v) emailTemplate(%+v)", err, tpl)
	}
	d.sendEmailLog(tpl, headers[email.TO], headers[email.CC], result)

	//retry
	if err != nil {
		address := headers[email.TO]
		if len(headers[email.CC]) > 0 {
			address = append(address, headers[email.CC]...)
			msg.SetHeader(email.CC)
		}

		for _, addr := range address {
			msg.SetHeader(email.TO, addr)
			result = email.EmailResOK
			if err = d.email.DialAndSend(msg); err != nil {
				result = email.EmailResFail
				log.Error("s.email.DialAndSend error(%v) to(%s) emailTemplate(%+v)", err, addr, tpl)
			}
			d.sendEmailLog(tpl, []string{addr}, []string{}, result)
			time.Sleep(time.Second * 5)
		}
	}
}

func (d *Dao) sendEmailLog(tpl *email.Template, to []string, cc []string, result string) {
	if tpl == nil || len(tpl.Headers) <= 0 || len(tpl.Headers[email.SUBJECT]) <= 0 {
		log.Error("sendEmailLog tpl nil | no headers, tpl(%+v)", tpl)
		return
	}

	address := fmt.Sprintf("to: %s", strings.Join(to, ","))
	if len(cc) > 0 {
		address = fmt.Sprintf("%s\ncc: %s", address, strings.Join(cc, ","))
	}

	item := &report.ManagerInfo{
		Uname:    tpl.Username,
		UID:      tpl.UID,
		Business: email.LogBusEmail,
		Type:     email.LogTypeEmailJob,
		Oid:      tpl.AID,
		Action:   tpl.Type,
		Ctime:    time.Now(),
		Content: map[string]interface{}{
			"subject":    tpl.Headers[email.SUBJECT][0],
			"body":       tpl.Body,
			"address":    address,
			"department": tpl.Department,
			"result":     result,
		},
	}
	report.Manager(item)
	log.Info("sendEmailLog template(%+v) result(%s) log.content(%+v)", tpl, result, item.Content)
}

//PushToRedis start to push email to redis according to speed
func (d *Dao) PushToRedis(c context.Context, tpl *email.Template) (isFast bool, key string, err error) {
	if tpl == nil {
		return
	}
	//探查发邮件速度快慢
	isFast = d.detector.Detect(tpl.UID)

	//超限名单只能被回落或下一次超限名单替代
	if d.detector.IsFastUnique(tpl.UID) {
		key = email.MailFastKey
		d.fastChan <- 1
	} else {
		key = email.MailKey
	}

	if err = d.PushRedis(c, tpl, key); err != nil {
		log.Error("PushToRedis d.PushRedis error(%v) key(%s), tpl(%+v) ", err, key, tpl)
	}
	return
}

//Start get email from redis and send
func (d *Dao) Start(key string) (err error) {
	var (
		bs  []byte
		tpl = &email.Template{}
	)

	bs, err = d.PopRedis(context.TODO(), key)
	if err != nil || bs == nil {
		time.Sleep(5 * time.Second)
		return
	}

	err = json.Unmarshal(bs, tpl)
	if err != nil {
		log.Error("email Start json.unmarshal error(%v) template(%s)", err, string(bs))
		return
	}

	//控制邮件发送频率
	st := <-d.controlChan
	d.SendMail(tpl)
	time.Sleep(time.Second * 5)
	d.controlChan <- st
	return
}
