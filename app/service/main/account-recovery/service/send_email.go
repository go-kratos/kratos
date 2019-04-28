package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go-common/app/service/main/account-recovery/model"
	"go-common/library/log"
	"go-common/library/net/metadata"

	"github.com/microcosm-cc/bluemonday"
)

type userLog struct {
	Business  int    `json:"business"`
	Type      int64  `json:"type"`
	IP        string `json:"ip"`
	CTime     string `json:"ctime"`
	Str0      string `json:"str_0"`
	ExtraData string `json:"extra_data"`
	Mid       int64  `json:"mid"`
}
type mailExtra struct {
	Subject string `json:"subject"`
	Content string `json:"content"`
}

// SendMailLog send mailLog
func (s *Service) SendMailLog(c context.Context, mid int64, mailType int, linkMail string, params ...string) (err error) {
	subject, content, err := buildMailData(true, mailType, params...)
	if err != nil {
		log.Error("SendMailLog err(%v)", err)
		return
	}
	mailData := &mailExtra{
		Subject: subject,
		Content: content,
	}
	mailDataBytes, err := json.Marshal(mailData)
	if err != nil {
		log.Error("mailData (%v) json marshal err(%v)", mailDataBytes, err)
		return
	}
	uLog := userLog{
		Mid:       mid,
		Business:  51,
		Type:      1,
		IP:        metadata.String(c, metadata.RemoteIP),
		CTime:     time.Now().Format("2006-01-02 15:04:05"),
		Str0:      linkMail,
		ExtraData: string(mailDataBytes),
	}
	for {
		if err = s.userActLogPub.Send(c, linkMail, uLog); err != nil {
			log.Error("databus send(%v) error(%v)", uLog, err)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		log.Info("SendMailLog success  uLog: %+v", uLog)
		break
	}
	return
}

// SendMailM send more type mail
func (s *Service) SendMailM(c context.Context, mailType int, linkMail string, params ...string) (err error) {
	subject, body, err := buildMailData(false, mailType, params...)
	if err != nil {
		log.Error("SendMailLog err(%v)", err)
		return
	}

	content := fmt.Sprintf("<html><body>"+
		"</body>"+
		" %s "+
		"</html>",
		body,
	)

	if s.d.SendMail(content, subject, linkMail); err != nil {
		err = errors.New("send mial fail")
		log.Error("SendMail fail,mailType=%d", mailType)
	}
	return
}

func buildMailData(hide bool, mailType int, params ...string) (subject, content string, err error) {
	switch mailType {
	case model.VerifyMail:
		subject = "【哔哩哔哩】账号找回-验证邮件（点开查看，无需回复）"
		// 验证码：params[0]
		content = generateHTML(mailType, params...)
		if hide {
			content = generateHTML(mailType, params[0][:2]+strings.Repeat("*", 4))
		}
	case model.CommitMail:
		subject = "【哔哩哔哩】账号找回-申述已受理（点开查看，无需回复）"
		// 打码后的UID： params[0] , 提交时间：params[1] , 案件ID：params[2]
		content = generateHTML(mailType, params...)
	case model.RejectMail:
		subject = "【哔哩哔哩】账号找回-申诉结果（点开查看，无需回复）"
		// 打码后的UID： params[0] , 提交时间：params[1] , 案件ID：params[2]
		content = generateHTML(mailType, params...)
	case model.AgreeMail:
		subject = "【哔哩哔哩】账号找回-申诉结果（点开查看，无需回复）"
		// 打码后的UID： params[0] , 提交时间：params[1] , 案件ID：params[2], 用户名：params[3] , 初始密码：params[4]
		content = generateHTML(mailType, params...)
		if hide {
			content = generateHTML(mailType, params[0], params[1], params[2], params[3], model.HIDEALL)
		}
	default:
		err = errors.New("没有该类型的邮件")
	}

	if hide {
		content = xssFilter(content)
	}

	return subject, content, err
}

// generateHTML generate mail body
func generateHTML(mailType int, params ...string) (content string) {
	var body bytes.Buffer
	switch mailType {
	case model.VerifyMail:
		body.WriteString("<div>尊敬的用户，你好：<br />")
		body.WriteString("你正在哔哩哔哩进行相应账号找回申诉，本次请求的邮件验证码为：%s,<br />")
		body.WriteString("本验证码30分钟内有效，请及时输入。<br />")
		body.WriteString("如非本人操作，请忽略该邮件。<br />")
		body.WriteString("祝在哔哩哔哩收获愉快！<br /></div>")
		content = fmt.Sprintf(body.String(), params[0])

	case model.CommitMail:
		body.WriteString("<div>尊敬的用户，你好：<br />")
		body.WriteString("账号【UID：%s】的找回申诉已受理（提交时间%s），案件ID为【%s】，请耐心等待，结果将在最长7个工作日内发送到本邮箱，请注意查收。")
		body.WriteString("（这是一封自动发送的邮件，请不要直接回复）<br /></div>")
		content = fmt.Sprintf(body.String(), params[0], params[1], params[2])

	case model.RejectMail:
		body.WriteString("<div>尊敬的用户，你好：<br />")
		body.WriteString("账号【UID：%s】的找回申诉已处理完毕（提交时间%s），案件ID为【%s】，由于你提供的信息与此账号信息不一致，申诉未能通过审核。请补充提供更准确完整的资料再行尝试。")
		body.WriteString("<a href='https://account.bilibili.com/appeal/home.html#/find'>点此重新申诉</a><br /><br />")
		body.WriteString("祝在哔哩哔哩收获愉快！<br /></div>")
		content = fmt.Sprintf(body.String(), params[0], params[1], params[2])

	case model.AgreeMail:
		body.WriteString("<div>尊敬的用户，你好：<br />")
		body.WriteString("账号【UID：%s】的找回申诉已处理完毕（提交时间%s），案件ID为【%s】，审核通过，密码已重置，请使用如下用户名与初始密码登录：<br />")
		body.WriteString("用户名：%s <br />")
		body.WriteString("初始密码：%s <br />")
		body.WriteString("为了你的账号安全，请及时修改密码与相应密保工具。 <br />")
		body.WriteString("（这是一封自动发送的邮件，请不要直接回复）<br /></div>")
		content = fmt.Sprintf(body.String(), params[0], params[1], params[2], params[3], params[4])
	}
	return content
}

// SendMailMany send more type mail
func (s *Service) SendMailMany(c context.Context, mailType int, batchRes []*model.BatchAppeal, userMap map[string]*model.User) (err error) {
	var (
		content = ""
		subject = ""
	)
	switch mailType {
	case model.RejectMail:
		subject = "【哔哩哔哩】账号找回-申诉结果（点开查看，无需回复）"
	case model.AgreeMail:
		subject = "【哔哩哔哩】账号找回-申诉结果（点开查看，无需回复）"
	default:
		err = errors.New("没有该类型的邮件")
	}
	if err != nil {
		log.Error("SendMailMany err(%v)", err)
		return
	}
	//循环发送多封邮件
	for _, appealInfo := range batchRes {
		params := make([]string, 0)
		rid := appealInfo.Rid
		mid := appealInfo.Mid
		linkMail := appealInfo.LinkMail
		ctime := appealInfo.Ctime.Time().Format("2006-01-02 15:04:05")

		params = append(params, hideUID(mid), ctime, rid)
		if model.AgreeMail == mailType {
			userID := userMap[mid].UserID
			pwd := userMap[mid].Pwd
			params = append(params, userID, pwd)

		}
		content = fmt.Sprintf("<html><body>"+
			"</body>"+
			" %s "+
			"</html>",
			generateHTML(mailType, params...),
		)
		if s.d.SendMail(content, subject, linkMail); err != nil {
			err = errors.New("send mial fail")
			log.Error("SendMailMany SendMail fail,mailType=%d", mailType)
		}
		midInt, _ := strconv.ParseInt(mid, 10, 64)
		//发送邮件日志
		s.SendMailLog(c, midInt, mailType, linkMail, params...)

	}
	return
}

// xss filter
func xssFilter(content string) string {
	p := bluemonday.StrictPolicy()
	return p.Sanitize(content)
}
