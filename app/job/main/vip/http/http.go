package http

import (
	"go-common/app/job/main/vip/conf"
	"go-common/app/job/main/vip/model"
	"go-common/app/job/main/vip/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var (
	s *service.Service
)

// Init init http sever instance.
func Init(c *conf.Config, ss *service.Service) {
	// init inner router
	engine := bm.DefaultServer(c.BM)
	innerRouter(engine)
	// init inner server
	if err := engine.Start(); err != nil {
		log.Error("engine start error(%v)", err)
		panic(err)
	}
	s = ss
}

// innerRouter init inner router.
func innerRouter(r *bm.Engine) {
	r.Ping(ping)
	r.GET("/scanUserInfo", checkscanUserInfo)
	r.GET("/handlerOrder", handlerOrder)
	r.GET("/handlerChangeHistory", handlerVipChangeHistory)
	r.GET("/handlerVipSendBcoin", handlerVipSendBcoin)
	r.GET("/sendBcoinJob", sendBcoinJob)
	r.GET("/hadExpiredJob", hadExpireJob)
	r.GET("/willExpiredJob", willExpireJob)
	r.GET("/sendMessageJob", sendMessageJob)
	r.GET("/autoRenewJob", autoRenewJob)
	r.GET("/syncvipdata", syncVipInfoData)
	r.GET("/clearcache", clearUserCache)
	r.GET("/scansalarylog", scanSalaryLog)
	r.GET("/checkuserdata", checkUserData)
	r.GET("/checkBcoinSalary", checkBcoinSalary)
	r.GET("/checkChangeHistory", checkHistory)

	r.GET("/sync/all/user", syncAllUser)
	r.GET("/frozen", frozen)
}

func syncAllUser(c *bm.Context) {
	log.Info("syncAllUser start........................................")
	s.SyncAllUser(c)
	log.Info("syncAllUser end........................................")
}

func checkHistory(c *bm.Context) {
	log.Info("check history info start........................................")
	mids, err := s.CheckChangeHistory(c)
	log.Info("check history info end..............error mids(%+v) error(%+v)", mids, err)
	c.JSON(mids, err)
}

func checkBcoinSalary(c *bm.Context) {
	log.Info("check bcoin info start........................................")
	mids, err := s.CheckBcoinData(c)
	log.Info("check bcoin info end..............error mids(%+v) error(%+v)", mids, err)
	c.JSON(mids, err)
}

func autoRenewJob(c *bm.Context) {
	//s.AutoRenewJob()
}

func sendBcoinJob(c *bm.Context) {
	//s.SendBcoinJob()
}

func hadExpireJob(c *bm.Context) {
	//s.HadExpiredMsgJob()
}

func willExpireJob(c *bm.Context) {
	//s.WillExpiredMsgJob()
}

func sendMessageJob(c *bm.Context) {
	//s.SendMessageJob()
}

// ping check server ok.
func ping(c *bm.Context) {}

func handlerOrder(c *bm.Context) {
	log.Info("handler order start.........................................")
	s.HandlerPayOrder()
	log.Info("handler order end ............................................")
}

func handlerVipChangeHistory(c *bm.Context) {
	log.Info("handler vip change history start ...................... ")
	s.HandlerVipChangeHistory()
	log.Info("handler vip change history end ...................... ")
}

func handlerVipSendBcoin(c *bm.Context) {
	log.Info(" handler vip send bcoin start ..............")
	s.HandlerBcoin()
	log.Info("handler vip send bcoin end ...............")
}

func checkscanUserInfo(c *bm.Context) {
	log.Info("scan user info start ..........................")
	s.ScanUserInfo(c)
	log.Info("scan user info end ...........................")
}

func syncVipInfoData(c *bm.Context) {
	var err error
	arg := new(struct {
		Mid int64 `form:"mid" validate:"required"`
	})
	if err = c.Bind(arg); err != nil {
		log.Error("c.Bind err(%+v)", err)
		return
	}
	if err = s.SyncUserInfoByMid(c, arg.Mid); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func clearUserCache(c *bm.Context) {
	var err error
	arg := new(struct {
		Mid int64 `form:"mid" validate:"required"`
	})
	if err = c.Bind(arg); err != nil {
		log.Error("c.Bind err(%+v)", err)
		return
	}
	s.ClearUserCache(arg.Mid)
	c.JSON(nil, nil)
}

func scanSalaryLog(c *bm.Context) {
	log.Info("scan salary log start ..........................")
	var err error
	if err = s.ScanSalaryLog(c); err != nil {
		log.Error("scan salary log err(%+v)", err)
		c.JSON(nil, err)
		return
	}
	log.Info("scan salary log end ...........................")
	c.JSON(nil, nil)
}

func checkUserData(c *bm.Context) {
	log.Info("check vip_user_info data start ..........................")
	var (
		err   error
		diffs map[int64]string
	)
	if diffs, err = s.CheckUserData(c); err != nil {
		c.JSON(diffs, err)
		return
	}
	log.Info("check vip_user_info data end diffs(%v) size(%d)...........................", diffs, len(diffs))
	c.JSON(diffs, err)
}

func frozen(c *bm.Context) {
	var err error
	arg := new(model.LoginLog)
	if err = c.Bind(arg); err != nil {
		log.Error("c.Bind err(%+v)", err)
		return
	}
	c.JSON(nil, s.Frozen(c, arg))
}
