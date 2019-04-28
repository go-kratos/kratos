package email

import (
	"fmt"
	"strconv"

	"go-common/app/job/main/videoup-report/model/email"
	"go-common/library/log"
)

//NotifyEmailTemplate 优质UP主/时政UP主/企业UP主/十万粉丝报备邮件
func (d *Dao) NotifyEmailTemplate(params map[string]string) (tpl *email.Template) {
	headers := map[string][]string{
		email.FROM: {d.c.Mail.Username},
	}

	//to
	typeIDStr := params["typeId"]
	if len(d.emailAddr[typeIDStr]) == 0 {
		log.Info("archive(%s) type(%s) don't config email address.", params["aid"], typeIDStr)
		return
	}
	headers[email.TO] = d.emailAddr[typeIDStr]

	//subject
	headers[email.SUBJECT] = []string{fmt.Sprintf("优质/十万粉稿件处理报备[%s]--操作人: %s[%s]", params["upName"], params["username"], params["department"])}
	//body
	body := `
	稿件标题：%s
	up主：%s
	稿件链接：http://www.bilibili.com/video/av%s
	触发条件：%s
	处理操作：%s
	`
	body = fmt.Sprintf(body, params["title"], params["upName"], params["aid"], params["condition"], params["change"])
	fromVideo, err := strconv.ParseBool(params["fromVideo"])
	if err != nil {
		log.Error("NotifyEmailTemplate get email template: strconv.ParseBool error(%v) aid(%s) fromVideo(%s)", err, params["aid"], params["fromVideo"])
		return
	}
	//视频追踪信息还没上线，先不写
	if !fromVideo {
		body += fmt.Sprintf("稿件追踪：http://manager.bilibili.co/#!/archive_utils/arc-track?aid=%s", params["aid"])
	}

	aid, _ := strconv.ParseInt(params["aid"], 10, 64)
	uid, _ := strconv.ParseInt(params["uid"], 10, 64)
	tpl = &email.Template{
		Headers:     headers,
		Body:        body,
		ContentType: "text/plain",
		Type:        email.EmailUP,
		AID:         aid,
		UID:         uid,
		Username:    params["username"],
		Department:  params["department"],
	}
	log.Info("NotifyEmailTemplate: email template(%+v)", tpl)
	return
}

//PrivateEmailTemplate 私单报备邮件模板
func (d *Dao) PrivateEmailTemplate(params map[string]string) (tpl *email.Template) {
	headers := map[string][]string{
		email.FROM: {d.c.Mail.Username},
	}

	//to
	to := d.PrivateAddr[params["typeId"]]
	if len(to) == 0 {
		log.Error("PrivateEmailTemplate lack email address config: typeId(%s), params(%v)", params["typeId"], params)
		return
	}
	headers[email.TO] = to

	//cc
	cc := d.PrivateAddr["CC"]
	if len(cc) > 0 {
		headers[email.CC] = cc
	}

	subject := fmt.Sprintf("私单稿件报备_%s_av%s", params["upName"], params["aid"])
	headers[email.SUBJECT] = []string{subject}

	body := `稿件标题： %s
	稿件状态： %s
	禁止项状态： 排行禁止:%s ；动态禁止:%s ； 推荐禁止:%s
	UP主： %s
	粉丝量：%s
	操作人： %s [%s]
	备注： %s`
	body = fmt.Sprintf(body, params["arcTitle"], params["arcState"], params["noRankAttr"], params["noDynamicAttr"], params["noRecommendAttr"],
		params["upName"], params["upFans"], params["mngName"], params["mngDepartment"], params["note"])

	aid, _ := strconv.ParseInt(params["aid"], 10, 64)
	uid, _ := strconv.ParseInt(params["uid"], 10, 64)
	tpl = &email.Template{
		Headers:     headers,
		Body:        body,
		ContentType: "text/plain",
		Type:        params["emailType"],
		AID:         aid,
		UID:         uid,
		Username:    params["mngName"],
		Department:  params["mngDepartment"],
	}
	log.Info("PrivateEmailTemplate: email template(%+v)", tpl)
	return
}

// MinitorNotifyTeamplate 审核监控报警邮件模板
func (d *Dao) MonitorNotifyTemplate(subject string, body string, toEmails []string) (tpl *email.Template) {
	headers := map[string][]string{
		email.FROM: {d.c.Mail.Username},
	}
	headers[email.TO] = toEmails
	headers[email.SUBJECT] = []string{subject}
	tpl = &email.Template{
		Headers:     headers,
		Body:        body,
		ContentType: "text/plain",
		Type:        email.EmailMonitor,
		AID:         0,
		UID:         0,
		Username:    "",
		Department:  "",
	}
	log.Info("MinitorNotifyTeamplate: email template(%+v)", tpl)
	return
}
