package charge

import (
	"context"
	"fmt"
	"time"

	model "go-common/app/job/main/growup/model/charge"
	task "go-common/app/job/main/growup/service"

	"go-common/library/log"
	"golang.org/x/sync/errgroup"
)

const (
	_layout    = "2006-01-02"
	_layoutSec = "2006-01-02 15:04:05"
)

var (
	_dbInsert = 1
	_dbUpdate = 2

	_batchSize = 2000
	_limitSize = 2000

	startWeeklyDate  time.Time
	startMonthlyDate time.Time
)

// RunAndSendMail run and send mail
func (s *Service) RunAndSendMail(c context.Context, date time.Time) (err error) {
	var msg, msgVideo, msgColumn, msgBgm string
	mailReceivers := []string{"shaozhenyu@bilibili.com", "gaopeng@bilibili.com", "limengqing@bilibili.com"}
	var (
		eg          errgroup.Group
		avBgmCharge = make(chan []*model.AvCharge, 1000)
	)
	// check task
	defer func() {
		task.GetTaskService().SetTaskStatus(c, task.TaskCreativeCharge, date.Format(_layout), err)
		if err != nil {
			msg = err.Error()
		}
		msgErr := s.email.SendMail(date, msg, "创作激励每日补贴%d年%d月%d日", mailReceivers...)
		if msgErr != nil {
			log.Error("s.email.SendMail error(%v)", msgErr)
		}
	}()

	err = task.GetTaskService().TaskReady(c, date.Format("2006-01-02"), task.TaskAvCharge, task.TaskCmCharge, task.TaskBgmSync)
	if err != nil {
		return
	}

	eg.Go(func() (err error) {
		startTime := time.Now().Unix()
		if err = s.runVideo(c, date, avBgmCharge); err != nil {
			log.Error("s.runVideo error(%v)", err)
		} else {
			msgVideo = fmt.Sprintf("%s 视频补贴计算完成，耗时%ds\n", date.Format("2006-01-02"), time.Now().Unix()-startTime)
		}
		return
	})
	eg.Go(func() (err error) {
		startTime := time.Now().Unix()
		if err = s.runColumn(c, date); err != nil {
			log.Error("s.runColumn error(%v)", err)
		} else {
			msgColumn = fmt.Sprintf("%s 专栏补贴计算完成，耗时%ds\n", date.Format("2006-01-02"), time.Now().Unix()-startTime)
		}
		return
	})
	eg.Go(func() (err error) {
		startTime := time.Now().Unix()
		if err = s.runBgm(c, date, avBgmCharge); err != nil {
			log.Error("s.runBgm error(%v)", err)
		} else {
			msgBgm = fmt.Sprintf("%s 素材补贴计算完成，耗时%ds\n", date.Format("2006-01-02"), time.Now().Unix()-startTime)
		}
		return
	})
	if err = eg.Wait(); err != nil {
		log.Error("run readGroup.Wait error(%v)", err)
	}
	msg = fmt.Sprintf("%s,%s,%s", msgVideo, msgColumn, msgBgm)
	return
}

func getStartWeeklyDate(date time.Time) time.Time {
	for date.Weekday() != time.Monday {
		date = date.AddDate(0, 0, -1)
	}
	return date
}

func getStartMonthlyDate(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.Local)
}
