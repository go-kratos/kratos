package service

import (
	"context"
	"fmt"
	"net/mail"
	"net/smtp"
	"strings"
	"time"

	"go-common/app/interface/openplatform/monitor-end/model"
	"go-common/app/interface/openplatform/monitor-end/model/monitor"
	"go-common/library/log"

	"github.com/scorredoira/email"
)

var text = `消息时间:%s
最近%d秒内，%s端%s服务%s的%s出现异常（%s)，异常数量超过告警阀值 %d
`

func (s *Service) mail(c context.Context, p *monitor.Log, t *model.Target, curr int, code string) {
	if t == nil {
		return
	}
	var groups = t.Groups
	if groups == nil || len(t.Groups) == 0 {
		product := s.productKeys[productKey(t.Product)]
		if product == nil || len(product.Groups) == 0 {
			return
		}
		groups = product.Groups
	}
	for _, g := range groups {
		if g.Name == "" || g.Interval == 0 {
			continue
		}
		if ok, err := s.dao.GetMailLock(c, g.Name, g.Interval, t, code); err != nil || !ok {
			continue
		}
		go s.mailByGroup(c, g.Receivers, p, curr, t.Threshold, t.Duration, code)
	}
	return
}

func (s *Service) mailByGroup(c context.Context, receivers string, p *monitor.Log, curr int, threshold int, duration int, code string) {
	tos := strings.Split(receivers, ",")
	if len(tos) == 0 {
		return
	}
	for i, t := range tos {
		tos[i] = t + "@bilibili.com"
	}
	source := sourceFromLog(p)
	now := time.Now().Format("2006-01-02 15:03:04")
	title := fmt.Sprintf("【端监控告警】%s端%s出现异常", source, p.Product)
	body := fmt.Sprintf(text, now, duration, source, p.Product, p.Event, p.SubEvent, code, threshold)
	if err := send(tos, title, body, 0); err != nil {
		log.Error("s.mailByGroup.send error(%+v), mail to(%s), title(%s), body(%s)", err, receivers, title, body)
	} else {
		log.Info("s.mailByGroup.send successed, mail to(%s), title(%s), body(%s)", receivers, title, body)
	}
}

/*
    * const HOST = 'smtp.exmail.qq.com';
    * const USER = 'show@bilibili.com';
    * const PASS = 'Kfpt2017';
	* const NAME = 'bilibili演出票务';
*/
func send(to []string, title string, body string, mode int) error {
	var (
		m    *email.Message
		host = "smtp.exmail.qq.com:25"
	)
	if mode == 0 {
		m = email.NewMessage(title, body)
	} else {
		m = email.NewHTMLMessage(title, body)
	}
	m.From = mail.Address{
		Name:    "kfc监控告警",
		Address: "show@bilibili.com",
	}
	m.To = to
	return email.Send(host, smtp.PlainAuth("", "show@bilibili.com", "Kfpt2017", "smtp.exmail.qq.com"), m)
}
