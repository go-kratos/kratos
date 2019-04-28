package service

import (
	"go-common/app/job/bbq/video/model"

	gomail "gopkg.in/gomail.v2"
)

//SendMail ...
func (s *Service) SendMail(mailType int) (err error) {
	var (
		m       = gomail.NewMessage()
		message = map[string][]string{}
		cType   string
		cBody   string
	)

	switch mailType {
	//运营i后台脚本导入完成推送邮件
	case model.JobFinishNotice:
		message["To"] = s.c.Mail.To
		message["Subject"] = []string{"同步运营筛选视频任务已完成"}
		cBody = "运营筛选视频已经导入完成！"
		cType = "text/plain"
	case 2:
		message["To"] = []string{"write your address"}
	}

	m.SetHeaders(message)
	m.SetBody(cType, cBody)

	err = s.dao.SendMail(m)
	return
}
