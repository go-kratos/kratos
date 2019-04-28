package service

import (
	"context"
	"math/rand"
	"time"

	"go-common/app/job/main/ugcpay/conf"
	"go-common/app/job/main/ugcpay/dao"
	"go-common/app/job/main/ugcpay/service/pay"
	"go-common/library/queue/databus"

	"github.com/robfig/cron"
)

var (
	ctx = context.Background()
)

// Service struct
type Service struct {
	c            *conf.Config
	dao          *dao.Dao
	binlogMQ     *databus.Databus
	elecBinlogMQ *databus.Databus
	cron         *cron.Cron
	pay          *pay.Pay
	taskLog      *taskLog
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:            c,
		dao:          dao.New(c),
		binlogMQ:     databus.New(c.BinlogMQ),
		elecBinlogMQ: databus.New(c.ElecBinlogMQ),
		cron:         cron.New(),
		pay: &pay.Pay{
			ID:                     conf.Conf.Biz.Pay.ID,
			Token:                  conf.Conf.Biz.Pay.Token,
			RechargeShellNotifyURL: conf.Conf.Biz.Pay.RechargeCallbackURL,
		},
	}
	s.taskLog = &taskLog{
		d: s.dao,
	}

	// 创建日账单任务
	taskBillDaily := &taskBillDaily{
		dao:        s.dao,
		pay:        s.pay,
		rnd:        rand.New(rand.NewSource(time.Now().Unix())),
		dayOffset:  conf.Conf.Biz.Task.DailyBillOffset,
		namePrefix: conf.Conf.Biz.Task.DailyBillPrefix,
		tl:         s.taskLog,
	}

	if err := s.cron.AddFunc(conf.Conf.Biz.Cron.TaskDailyBill, s.wrapDisProc(taskBillDaily)); err != nil {
		panic(err)
	}

	// 创建up虚拟账户入账任务
	taskAccountUser := &taskAccountUser{
		dao:        s.dao,
		taskPre:    taskBillDaily, // 前置任务
		dayOffset:  conf.Conf.Biz.Task.DailyBillOffset,
		namePrefix: conf.Conf.Biz.Task.AccountUserPrefix,
		tl:         s.taskLog,
	}

	if err := s.cron.AddFunc(conf.Conf.Biz.Cron.TaskAccountUser, s.wrapDisProc(taskAccountUser)); err != nil {
		panic(err)
	}

	// 创建资金池入账任务
	taskAccountBiz := &taskAccountBiz{
		dao:        s.dao,
		taskPre:    taskBillDaily, // 前置任务
		dayOffset:  conf.Conf.Biz.Task.DailyBillOffset,
		namePrefix: conf.Conf.Biz.Task.AccountBizPrefix,
		tl:         s.taskLog,
	}

	if err := s.cron.AddFunc(conf.Conf.Biz.Cron.TaskAccountBiz, s.wrapDisProc(taskAccountBiz)); err != nil {
		panic(err)
	}

	// 创建月账单任务
	taskBillMonthly := &taskBillMonthly{
		dao:         s.dao,
		rnd:         rand.New(rand.NewSource(time.Now().Unix())),
		monthOffset: conf.Conf.Biz.Task.MonthBillOffset,
		namePrefix:  conf.Conf.Biz.Task.MonthBillPrefix,
		tl:          s.taskLog,
	}

	if err := s.cron.AddFunc(conf.Conf.Biz.Cron.TaskMonthlyBill, s.wrapDisProc(taskBillMonthly)); err != nil {
		panic(err)
	}

	// 创建转贝壳任务
	taskRechargeShell := &taskRechargeShell{
		dao:         s.dao,
		pay:         s.pay,
		rnd:         rand.New(rand.NewSource(time.Now().Unix())),
		monthOffset: conf.Conf.Biz.Task.RechargeShellOffset,
		namePrefix:  conf.Conf.Biz.Task.RechargeShellPrefix,
		tl:          s.taskLog,
	}

	if err := s.cron.AddFunc(conf.Conf.Biz.Cron.TaskRechargeShell, s.wrapDisProc(taskRechargeShell)); err != nil {
		panic(err)
	}
	s.cron.Start()

	go s.binlogproc()
	go s.elecbinlogproc()

	// go s.repairOrderUser() 修复订单用
	return s
}

// Ping Service
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close Service
func (s *Service) Close() {
	s.dao.Close()
}
