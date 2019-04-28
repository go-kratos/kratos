package service

import (
	"context"
	"fmt"
	"time"

	"go-common/app/job/main/growup/conf"
	"go-common/app/job/main/growup/dao"
	"go-common/app/job/main/growup/dao/charge"
	"go-common/app/job/main/growup/dao/dataplatform"
	"go-common/app/job/main/growup/dao/email"
	"go-common/app/job/main/growup/dao/income"
	"go-common/app/job/main/growup/dao/tag"
	"go-common/app/job/main/growup/service/ctrl"

	"go-common/library/log"
	"go-common/library/queue/databus"
)

// Service struct
type Service struct {
	conf   *conf.Config
	dao    *dao.Dao
	email  *email.Dao
	dp     *dataplatform.Dao
	tag    *tag.Dao
	income *income.Dao
	charge *charge.Dao
	// databus sub
	arcSub *databus.Databus
}

// New fn
func New(c *conf.Config, executor ctrl.Executor) (s *Service) {
	s = &Service{
		conf:   c,
		dao:    dao.New(c),
		email:  email.New(c),
		dp:     dataplatform.New(c),
		tag:    tag.New(c),
		income: income.New(c),
		charge: charge.New(c),
		arcSub: databus.New(c.ArchiveSub),
	}

	// init task status service
	taskSvr = &taskService{
		dao: s.dao,
		dp:  s.dp,
	}
	log.Info("service start")
	executor.Submit(
		s.checkExpired,
		s.updateDateBlacklist,
		s.sendMail,
		s.updateCheat,
		s.mailTagIncome,
		s.calCreativeActivity,
		s.updateUpInfoVideo,
		s.creativeBudget,
		s.checkAvBreach,
		s.checkUpPunish,
		s.archiveConsume,
		s.syncIncomeBubbleMetaTask,
		s.snapshotIncomeBubbleTask,
	)
	return s
}

// check the account state of video up
func (s *Service) checkExpired(ctx context.Context) {
	for {
		// up_info_video
		s.expiredCheck(0, 5, 0, "up_info_video") // check all expired from quit to default
		s.expiredCheck(0, 7, 3, "up_info_video") // check all expired from forbidden to signed
		s.expiredCheck(1, 4, 0, "up_info_video") // check ugc expired from reject to default
		s.expiredCheck(2, 4, 1, "up_info_video") // check pgc expired from reject to pre-audit

		// up_info_column
		s.expiredCheck(0, 5, 0, "up_info_column") // check all expired from quit to default
		s.expiredCheck(0, 7, 3, "up_info_column") // check all expired from forbidden to signed
		s.expiredCheck(1, 4, 0, "up_info_column") // check ugc expired from reject to default
		s.expiredCheck(2, 4, 1, "up_info_column") // check pgc expired from reject to pre-audit
		time.Sleep(NextDay(0, 0, 0))
	}
}

// updateDateBlacklist update  blacklist at 17:00 every day
func (s *Service) updateDateBlacklist(ctx context.Context) {
	for {
		time.Sleep(NextDay(17, 0, 0))
		msg := ""
		log.Info("Exec growup-job updateDateBlacklist begin:%v", time.Now().Format("2006-01-02 15:04:05"))
		err := s.UpdateBlacklist(context.TODO())
		if err != nil {
			msg = fmt.Sprintf("UpdateBlacklist error(%v)", err)
		} else {
			msg = "Success"
		}
		err = s.email.SendMail(time.Now(), msg, "创作激励同步黑名单%d年%d月%d日", "shaozhenyu@bilibili.com", "gaopeng@bilibili.com", "limengqing@bilibili.com")
		if err != nil {
			log.Error("s.email.SendMail error(%v)", err)
		}
		log.Info("Exec growup-job updateDateBlacklist end:%v", time.Now().Format("2006-01-02 15:04:05"))
	}
}

func (s *Service) updateUpInfoVideo(ctx context.Context) {
	var mailReceivers []string
	for _, v := range s.conf.Mail.Send {
		if v.Type == 4 {
			mailReceivers = v.Addr
			break
		}
	}
	for {
		time.Sleep(NextDay(11, 0, 0))
		log.Info("Exec growup-job updateUpInfoVideo begin:%v", time.Now().Format("2006-01-02 15:04:05"))
		msg := ""
		err := s.UpdateUpInfo(context.TODO())
		if err != nil {
			msg = fmt.Sprintf("s.UpdateUpInfo error(%v)", err)
		} else {
			msg, err = s.autoExamination(context.TODO())
			if err != nil {
				msg = fmt.Sprintf("s.autoExamination error(%v)", err)
			}
		}
		err = s.email.SendMail(time.Now(), msg, "创作激励自动过审%d年%d月%d日", mailReceivers...)
		if err != nil {
			log.Error("s.email.SendMail error(%v)", err)
		}
		log.Info("Exec growup-job updateUpInfoVideo end:%v", time.Now().Format("2006-01-02 15:04:05"))
	}
}

func (s *Service) sendMail(ctx context.Context) {
	for {
		time.Sleep(NextDay(12, 0, 0))
		log.Info("Exec growup-job sendMail begin:%v", time.Now().Format("2006-01-02 15:04:05"))
		s.CombineMails()
		log.Info("Exec growup-job sendMail end:%v", time.Now().Format("2006-01-02 15:04:05"))
	}
}

func (s *Service) updateCheat(ctx context.Context) {
	for {
		time.Sleep(NextDay(13, 0, 0))
		log.Info("Exec growup-job updateCheat begin: %v", time.Now().Format("2006-01-02 15:04:05"))
		t := time.Now().Add(-24 * time.Hour)
		s.CheatStatistics(context.TODO(), t)
		log.Info("Exec growup-job updateCheat end: %v", time.Now().Format("2006-01-02 15:04:05"))
	}
}

func (s *Service) mailTagIncome(ctx context.Context) {
	for {
		time.Sleep(NextDay(19, 30, 0))
		log.Info("Exec growup-job mailTagIncome begin: %v", time.Now().Format("2006-01-02 15:04:05"))
		t := time.Now().Add(-24 * time.Hour)
		date := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
		s.sendTagIncome(context.TODO(), date)
		log.Info("End growup-job mailTagIncome: %v", time.Now().Format("2006-01-02 15:04:05"))
	}
}

func (s *Service) calCreativeActivity(ctx context.Context) {
	for {
		time.Sleep(NextDay(15, 30, 0))
		msg := ""
		log.Info("Exec growup-job calCreativeActivity begin: %v", time.Now().Format("2006-01-02 15:04:05"))
		t := time.Now().Add(-24 * time.Hour)
		date := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
		err := s.CreativeActivity(context.TODO(), date)
		if err != nil {
			msg = fmt.Sprintf("calCreativeActivity error(%v)", err)
		} else {
			msg = "Success"
		}
		err = s.email.SendMail(date, msg, "创作激励活动每日计算%d年%d月%d日", "shaozhenyu@bilibili.com", "gaopeng@bilibili.com", "limengqing@bilibili.com")
		if err != nil {
			log.Error("s.email.SendMail error(%v)", err)
		}
		log.Info("End growup-job calCreativeActivity: %v", time.Now().Format("2006-01-02 15:04:05"))
	}
}

func (s *Service) creativeBudget(ctx context.Context) {
	for {
		time.Sleep(NextDay(20, 0, 0))
		msg := ""
		log.Info("Exec growup-job creativeBudget begin: %v", time.Now().Format("2006-01-02 15:04:05"))
		t := time.Now().Add(-24 * time.Hour)
		date := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
		err := s.CreativeBudget(context.TODO(), date)
		if err != nil {
			msg = fmt.Sprintf("creativeBudget error(%v)", err)
		} else {
			msg = "Success"
		}
		err = s.email.SendMail(date, msg, "创作激励预算每日计算%d年%d月%d日", "shaozhenyu@bilibili.com", "gaopeng@bilibili.com", "limengqing@bilibili.com")
		if err != nil {
			log.Error("s.email.SendMail error(%v)", err)
		}
		log.Info("End growup-job creativeBudget: %v", time.Now().Format("2006-01-02 15:04:05"))
	}
}

func (s *Service) checkAvBreach(ctx context.Context) {
	var mailReceivers []string
	for _, v := range s.conf.Mail.Send {
		if v.Type == 4 {
			mailReceivers = v.Addr
			break
		}
	}
	for {
		time.Sleep(NextDay(0, 0, 0))
		log.Info("Exec growup-job check and auto breach av begin: %v", time.Now().Format("2006-01-02 15:04:05"))
		t := time.Now().Add(-48 * time.Hour)
		date := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
		msg, err := s.autoAvBreach(context.TODO(), date.Format(_layout))
		if err != nil {
			msg = fmt.Sprintf("autoAvBreach error(%v)", err)
		}
		err = s.email.SendMail(date, msg, "创作激励自制转转载每日扣除%d年%d月%d日", mailReceivers...)
		if err != nil {
			log.Error("s.email.SendMail error(%v)", err)
		}
		log.Info("End growup-job autoAvBreach: %v", time.Now().Format("2006-01-02 15:04:05"))
	}
}

func (s *Service) checkUpPunish(ctx context.Context) {
	var mailReceivers []string
	for _, v := range s.conf.Mail.Send {
		if v.Type == 4 {
			mailReceivers = v.Addr
			break
		}
	}
	for {
		time.Sleep(NextDay(10, 0, 0))
		if time.Now().Weekday() != time.Monday {
			continue
		}
		log.Info("Exec growup-job checkUpPunish begin: %v", time.Now().Format("2006-01-02 15:04:05"))
		msg, err := s.autoUpPunish(context.TODO())
		if err != nil {
			msg = fmt.Sprintf("s.autoUpPunish error(%v)", err)
		}
		err = s.email.SendMail(time.Now(), msg, "创作激励自制转转载处罚%d年%d月%d日", mailReceivers...)
		if err != nil {
			log.Error("s.email.SendMail error(%v)", err)
		}
		log.Info("End growup-job checkUpPunish: %v", time.Now().Format("2006-01-02 15:04:05"))
	}
}

func (s *Service) syncIncomeBubbleMetaTask(ctx context.Context) {
	for {
		msg := "ok"
		time.Sleep(NextDay(12, 0, 0))
		log.Info("Exec growup-job syncIncomeBubbleMetaTask begin: %v", time.Now().Format("2006-01-02 15:04:05"))
		t := time.Now().Add(-24 * time.Hour)
		date := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
		err := s.SyncIncomeBubbleMetaTask(context.TODO(), date)
		if err != nil {
			msg = fmt.Sprintf("s.syncIncomeBubbleMetaTask error(%v)", err)
		}
		err = s.email.SendMail(time.Now(), msg, "创作激励同步动态转发抽奖%d年%d月%d日", "shaozhenyu@bilibili.com", "gaopeng@bilibili.com", "limengqing@bilibili.com")
		if err != nil {
			log.Error("s.email.SendMail error(%v)", err)
		}
		log.Info("End growup-job syncIncomeBubbleMetaTask: %v", time.Now().Format("2006-01-02 15:04:05"))
	}
}

func (s *Service) snapshotIncomeBubbleTask(ctx context.Context) {
	for {
		msg := "ok"
		time.Sleep(NextDay(19, 0, 0))
		log.Info("Exec growup-job snapshotIncomeBubbleTask begin: %v", time.Now().Format("2006-01-02 15:04:05"))
		t := time.Now().Add(-24 * time.Hour)
		date := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
		err := s.SnapshotBubbleIncomeTask(context.TODO(), date)
		if err != nil {
			msg = fmt.Sprintf("s.snapshotIncomeBubbleTask error(%v)", err)
		}
		err = s.email.SendMail(time.Now(), msg, "创作激励动态转发抽奖收入同步%d年%d月%d日", "shaozhenyu@bilibili.com", "gaopeng@bilibili.com", "limengqing@bilibili.com")
		if err != nil {
			log.Error("s.email.SendMail error(%v)", err)
		}
		log.Info("End growup-job snapshotIncomeBubbleTask: %v", time.Now().Format("2006-01-02 15:04:05"))
	}
}

// NextDay  next day x hours
func NextDay(hour, min, second int) time.Duration {
	n := time.Now()
	d := time.Date(n.Year(), n.Month(), n.Day(), hour, min, second, 0, n.Location())
	for !d.After(n) {
		d = d.AddDate(0, 0, 1)
	}
	return time.Until(d)
}

// Close close the service
func (s *Service) Close() {
	s.dao.Close()
	s.tag.Close()
	s.income.Close()
	s.charge.Close()
}

// Ping check server ok.
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}
