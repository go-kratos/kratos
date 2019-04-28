package service

import (
	"context"
	"fmt"

	"go-common/library/log"
)

const (
	ISO8601Date = "2006-01-02"
)

// SendEmail is send the finished Task info  to reciever
func (s *Service) SendEmail(c context.Context, taskID int64) (err error) {
	task, err := s.dao.DetailTask(c, taskID)
	if err != nil {
		log.Error("s.SendEmail() error(%v)", err)
		return
	}
	createAt := task.CTime.Time().Format(ISO8601Date)
	var sourceDesc string
	if task.SourceType == 1 {
		sourceDesc = "创作姬"
	} else {
		sourceDesc = "其他"
	}
	var appStr string
	if task.Platform == 1 {
		appStr = "IOS"
	} else if task.Platform == 2 {
		appStr = "Android"
	}
	date := task.LogDate.Time().Format(ISO8601Date)
	subject := fmt.Sprintf(" %s 创建的日志上报完成通知", createAt)
	body := fmt.Sprintf("你于%s创建的一条日志上报任务（上报来源：%s，%s App端，采集的日志文件日期：%s，指定MID：%d），现已上报完毕。", createAt, sourceDesc, appStr, date, task.MID)
	err = s.dao.SendEmail(subject, task.ContactEmail, body)
	return
}
