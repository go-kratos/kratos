package charge

import (
	"context"
	"fmt"
	"time"

	"go-common/app/job/main/growup/conf"
	chargeD "go-common/app/job/main/growup/dao/charge"
	"go-common/app/job/main/growup/dao/email"
	"go-common/app/job/main/growup/service"
	"go-common/app/job/main/growup/service/ctrl"
	"go-common/library/log"
)

// Service struct
type Service struct {
	conf  *conf.Config
	dao   *chargeD.Dao
	email *email.Dao
}

// New fn
func New(c *conf.Config, executor ctrl.Executor) (s *Service) {
	s = &Service{
		conf:  c,
		dao:   chargeD.New(c),
		email: email.New(c),
	}
	log.Info("charge service start")
	executor.Submit(
		s.calDailyCreativeCharge,
		s.checkColumnDailyCharge,
	)
	return s
}

func (s *Service) calDailyCreativeCharge(ctx context.Context) {
	for {
		time.Sleep(service.NextDay(17, 0, 0))
		log.Info("calDailyCreativeCharge begin:%v", time.Now().Format("2006-01-02 15:04:05"))
		now := time.Now()
		date := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, time.Local)
		err := s.RunAndSendMail(context.TODO(), date)
		if err != nil {
			log.Error("s.RunAndSendMail error(%v)", err)
		}
		log.Info("calDailyCreativeCharge end:%v", time.Now().Format("2006-01-02 15:04:05"))
	}
}

func (s *Service) checkColumnDailyCharge(ctx context.Context) {
	for {
		time.Sleep(service.NextDay(13, 0, 0))
		msg := ""
		log.Info("checkColumnDailyCharge begin:%v", time.Now().Format("2006-01-02 15:04:05"))
		now := time.Now()
		date := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, time.Local)
		err := s.CheckTaskColumn(context.TODO(), date.Format(_layout))
		if err != nil {
			log.Error("s.CheckTaskColumn error(%v)", err)
		}
		if err != nil {
			msg = fmt.Sprintf("CheckTaskColumn error(%v)", err)
		} else {
			msg = "Success"
		}
		err = s.email.SendMail(date, msg, "创作激励专栏补贴数据%d年%d月%d日", "shaozhenyu@bilibili.com", "gaopeng@bilibili.com", "limengqing@bilibili.com")
		if err != nil {
			log.Error("s.email.SendMail error(%v)", err)
		}
		log.Info("checkColumnDailyCharge end:%v", time.Now().Format("2006-01-02 15:04:05"))
	}
}

// Ping check dao health.
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}
