package tag

import (
	"context"
	"fmt"
	"time"

	"go-common/app/job/main/growup/conf"
	"go-common/app/job/main/growup/dao/email"
	dao "go-common/app/job/main/growup/dao/tag"
	"go-common/app/job/main/growup/service"
	"go-common/app/job/main/growup/service/ctrl"

	"go-common/library/log"
)

// Service struct
type Service struct {
	conf  *conf.Config
	dao   *dao.Dao
	email *email.Dao
}

// New fn
func New(c *conf.Config, executor ctrl.Executor) (s *Service) {
	s = &Service{
		conf:  c,
		dao:   dao.New(c),
		email: email.New(c),
	}
	log.Info("tag service start")
	executor.Submit(
		s.calDailyTagRatio,
		s.calDailyTagIncome,
	)
	return s
}

// Ping check dao health.
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close dao
func (s *Service) Close() {
	s.dao.Close()
}

func (s *Service) calDailyTagRatio(ctx context.Context) {
	for {
		time.Sleep(service.NextDay(14, 0, 0))
		msg := ""
		log.Info("calDailyTagRatio begin:%v", time.Now().Format("2006-01-02 15:04:05"))
		date := time.Now().AddDate(0, 0, -1)
		err := s.TagRatioAll(context.TODO(), date.Format("2006-01-02"))
		if err != nil {
			msg = fmt.Sprintf("calDailyTagRatio error(%v)", err)
		} else {
			msg = "Success"
		}
		err = s.email.SendMail(date, msg, "标签每日计算%d年%d月%d日", "shaozhenyu@bilibili.com")
		if err != nil {
			log.Error("s.email.SendMail error(%v)", err)
		}
		log.Info("calDailyTagRatio end:%v", time.Now().Format("2006-01-02 15:04:05"))
	}
}

func (s *Service) calDailyTagIncome(ctx context.Context) {
	for {
		time.Sleep(service.NextDay(19, 0, 0))
		msg := ""
		log.Info("calDailyTagIncome begin:%v", time.Now().Format("2006-01-02 15:04:05"))
		date := time.Now().AddDate(0, 0, -1)
		err := s.TagIncomeAll(context.TODO(), date.Format("2006-01-02"))
		if err != nil {
			msg = fmt.Sprintf("calDailyTagIncome error(%v)", err)
		} else {
			msg = "Success"
		}
		err = s.TagUps(context.TODO(), date)
		if err != nil {
			msg = fmt.Sprintf("s.TagUps error(%v)", err)
		}
		err = s.email.SendMail(date, msg, "标签每日收入计算%d年%d月%d日", "shaozhenyu@bilibili.com")
		if err != nil {
			log.Error("s.email.SendMail error(%v)", err)
		}
		log.Info("calDailyTagIncome end:%v", time.Now().Format("2006-01-02 15:04:05"))
	}
}
