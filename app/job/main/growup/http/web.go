package http

import (
	"time"

	"go-common/app/job/main/growup/conf"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func execAvRatio(c *bm.Context) {
	log.Info("begin update av charge ratio")
	v := new(struct {
		Date string `form:"date" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	t, err := time.Parse("2006-01-02", v.Date)
	if err != nil {
		log.Error("execAvRatio date error!date:%s", v.Date)
		return
	}
	log.Info("ratio,hour:%d, num:%d, sleep:%d", conf.Conf.Ratio.Hour, conf.Conf.Ratio.Num, conf.Conf.Ratio.Sleep)
	for {
		rows, _ := svr.DeleteAvRatio(c, conf.Conf.Ratio.Num)
		time.Sleep(time.Duration(conf.Conf.Ratio.Sleep) * time.Millisecond)
		if rows == 0 {
			break
		}
	}
	err = svr.ExecRatioForHTTP(c, t.Year(), int(t.Month()), t.Day())
	if err != nil {
		log.Error("Exec avRatio from http error!(%v)", err)
	} else {
		log.Info("Exec avRatio from http succeed!")
	}
	c.JSON(nil, err)
}

func execIncome(c *bm.Context) {
	log.Info("begin update up income")
	v := new(struct {
		Date string `form:"date" validate:"required" `
	})

	if err := c.Bind(v); err != nil {
		return
	}
	t, err := time.Parse("2006-01-02", v.Date)
	if err != nil {
		log.Error("execIncome date error!date:%s", v.Date)
		return
	}
	log.Info("income,hour:%d, num:%d, sleep:%d", conf.Conf.Income.Hour, conf.Conf.Income.Num, conf.Conf.Income.Sleep)
	err = svr.ExecIncomeForHTTP(c, t.Year(), int(t.Month()), t.Day())
	if err != nil {
		log.Error("ExecIncomeForHTTP error!(%v)", err)
	} else {
		log.Info("Exec Income from http succeed!")
	}
	c.JSON(nil, err)
}

func getUpIncomeStatis(c *bm.Context) {
	log.Info("begin calculate up income statis")
	v := new(struct {
		Date        string `form:"date" validate:"required"`
		HasWithdraw int    `form:"has_withdraw"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	err := svr.GetUpIncomeStatis(c, v.Date, v.HasWithdraw)
	if err != nil {
		log.Error("(job growup svr.GetUpIncomeStatis error(%v)", err)
	}
	c.JSON(nil, err)
}

func getAvIncomeStatis(c *bm.Context) {
	v := new(struct {
		Date string `form:"date" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	err := svr.GetAvIncomeStatis(c, v.Date)
	if err != nil {
		log.Error("(job growup svr.GetAvIncomeStatis error(%v)", err)
	}
	c.JSON(nil, err)
}

func updateTagIncome(c *bm.Context) {
	log.Info("begin fix tag income")
	v := new(struct {
		Date string `form:"date" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	err := svr.UpdateTagIncome(c, v.Date)
	if err != nil {
		log.Error("(job growup svr.UpdateTagIncome error(%v)", err)
	}
	c.JSON(nil, err)
}

func fixUpIncome(c *bm.Context) {
	log.Info("begin fix up income from tag")
	v := new(struct {
		Date        string `form:"date" validate:"required"`
		TagID       int64  `form:"tag_id" validate:"required"`
		AddCount    int    `form:"add_count" validate:"required"`
		TotalIncome int    `form:"total_income" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	err := svr.FixUpIncome(c, v.Date, v.TagID, v.AddCount, v.TotalIncome)
	if err != nil {
		log.Error("(job growup svr.FixUpIncome error(%v)", err)
	}
	c.JSON(nil, err)
}

func fixIncome(c *bm.Context) {
	log.Info("begin fix income")
	err := svr.FixIncome(c)
	if err != nil {
		log.Error("(job growup svr.FixIncome error(%v)", err)
	}
	c.JSON(nil, err)
}

func fixUpAvStatis(c *bm.Context) {
	log.Info("begin fix up av statis")
	v := new(struct {
		Count int `form:"count" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	err := svr.FixUpAvStatis(c, v.Count)
	if err != nil {
		log.Error("(job growup svr.FixUpAvStatis error(%v)", err)
	}
	c.JSON(nil, err)
}

func updateWithdraw(c *bm.Context) {
	v := new(struct {
		OldDate string `form:"old_date" validate:"required"`
		NewDate string `form:"new_date" validate:"required"`
		Count   int64  `form:"count" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	err := svr.UpdateWithdraw(c, v.OldDate, v.NewDate, v.Count)
	if err != nil {
		log.Error("(svr.UpdateWithdraw error(%v)", err)
	}
	c.JSON(nil, err)
}

func creativeIncome(c *bm.Context) {
	v := new(struct {
		Date string `form:"date" validate:"required"`
	})

	if err := c.Bind(v); err != nil {
		return
	}
	date, err := time.ParseInLocation("2006-01-02", v.Date, time.Local)
	if err != nil {
		log.Error("time.Parse date(%s) error", v.Date)
		c.JSON(nil, err)
		return
	}

	log.Info("start run %s income", v.Date)
	err = incomeSrv.RunAndSendMail(c, date)
	if err != nil {
		log.Error("incomeSrv.Run error(%v)", err)
	} else {
		log.Info("run %s income success", v.Date)
	}
	c.JSON(nil, err)
}

func creativeCharge(c *bm.Context) {
	v := new(struct {
		Date string `form:"date" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	date, err := time.ParseInLocation("2006-01-02", v.Date, time.Local)
	if err != nil {
		log.Error("time.Parse date(%s) error", v.Date)
		c.JSON(nil, err)
		return
	}

	log.Info("start run %s charge", v.Date)
	err = chargeSrv.RunAndSendMail(c, date)
	if err != nil {
		log.Error("chargeSrv.RunAndSendMail error(%v)", err)
	} else {
		log.Info("run %s charge success", v.Date)
	}
	c.JSON(nil, err)
}

func creativeStatis(c *bm.Context) {
	v := new(struct {
		Date string `form:"date" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	date, err := time.ParseInLocation("2006-01-02", v.Date, time.Local)
	if err != nil {
		log.Error("time.Parse date(%s) error", v.Date)
		c.JSON(nil, err)
		return
	}

	log.Info("start run %s statis", v.Date)
	err = incomeSrv.RunStatis(c, date)
	if err != nil {
		log.Error("incomeSrv.RunStatis error(%v)", err)
	} else {
		log.Info("run %s statis success", v.Date)
	}
	c.JSON(nil, err)
}

func creativeBill(c *bm.Context) {
	v := new(struct {
		StartDate string `form:"start_date" validate:"required"`
		EndDate   string `form:"end_date" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	start, err := time.ParseInLocation("2006-01-02", v.StartDate, time.Local)
	if err != nil {
		log.Error("time.Parse date(%s) error", v.StartDate)
		c.JSON(nil, err)
		return
	}
	end, err := time.ParseInLocation("2006-01-02", v.EndDate, time.Local)
	if err != nil {
		log.Error("time.Parse date(%s) error", v.EndDate)
		c.JSON(nil, err)
		return
	}

	if err = svr.CreativeUpBill(c, start, end); err != nil {
		log.Error("svr.CreativeUpBill error(%v)", err)
	}
	c.JSON(nil, err)
}

func creativeBudget(c *bm.Context) {
	v := new(struct {
		Date string `form:"date" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	date, err := time.ParseInLocation("2006-01-02", v.Date, time.Local)
	if err != nil {
		log.Error("time.Parse date(%s) error", v.Date)
		c.JSON(nil, err)
		return
	}

	err = svr.CreativeBudget(c, date)
	if err != nil {
		log.Error("svr.CreativeBudget error(%v)", err)
	}
	c.JSON(nil, err)
}

func creativeActivity(c *bm.Context) {
	v := new(struct {
		Date string `form:"date" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	date, err := time.ParseInLocation("2006-01-02", v.Date, time.Local)
	if err != nil {
		log.Error("time.Parse date(%s) error", v.Date)
		c.JSON(nil, err)
		return
	}
	err = svr.CreativeActivity(c, date)
	if err != nil {
		log.Error("svr.CreativeActivity error(%v)", err)
	}
	c.JSON(nil, err)
}

func updateUpInfoVideo(c *bm.Context) {
	err := svr.UpdateUpInfo(c)
	if err != nil {
		log.Error("svr.UpdateUpInfo error(%v)", err)
	} else {
		log.Info("run %s UpdateUpInfo success", time.Now().Format("2006-01-02"))
	}
	c.JSON(nil, err)
}

func fixTagAdjust(c *bm.Context) {
	v := new(struct {
		ID int64 `form:"id" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	err := svr.UpdateTagAdjust(c, v.ID)
	if err != nil {
		log.Error("svr.UpdateTagAdjust error(%v)", err)
	}
	c.JSON(nil, err)
}

func fixAccountType(c *bm.Context) {
	v := new(struct {
		MID         int64 `form:"mid" validate:"required"`
		AccountType int   `form:"account_type" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	err := svr.UpdateAccountType(c, v.MID, v.AccountType)
	if err != nil {
		log.Error("svr.UpdateAccountType error!(%v)", err)
	}
	c.JSON(nil, err)
}

func syncBGM(c *bm.Context) {
	err := incomeSrv.SyncBgmInfo(c)
	if err != nil {
		log.Error("svr.SyncBgmInfo error!(%v)", err)
	}
	c.JSON(nil, err)
}

func fixUpAccount(c *bm.Context) {
	v := new(struct {
		MID              int64 `form:"mid"`
		TotalIncome      int64 `form:"total_income"`
		UnwithdrawIncome int64 `form:"unwithdraw_income"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	err := svr.UpdateUpAccountMoney(c, v.MID, v.TotalIncome, v.UnwithdrawIncome)
	if err != nil {
		log.Error("svr.UpdateUpAccountMoney error!(%v)", err)
	}
	c.JSON(nil, err)
}

func fixBaseIncome(c *bm.Context) {
	v := new(struct {
		MID  int64  `form:"mid" validate:"required"`
		Base int64  `form:"base" validate:"required"`
		Date string `form:"date" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	err := svr.FixBaseIncome(c, v.Base, v.MID, v.Date)
	if err != nil {
		log.Error("svr.FixBaseIncome error!(%v)", err)
	}
	c.JSON(nil, err)
}

func fixAvBreach(c *bm.Context) {
	v := new(struct {
		MID   int64  `form:"mid" validate:"required"`
		Date  string `form:"date" validate:"required"`
		Count int    `form:"count" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	err := svr.FixAvBreach(c, v.MID, v.Date, v.Count)
	if err != nil {
		log.Error("svr.FixAvBreach error!(%v)", err)
	}
	c.JSON(nil, err)
}

func fixUpTotalIncome(c *bm.Context) {
	v := new(struct {
		Table string `form:"table" validate:"required"`
		Date  string `form:"date" validate:"required"`
		Count int    `form:"count" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	err := svr.FixUpTotalIncome(c, v.Table, v.Date, v.Count)
	if err != nil {
		log.Error("svr.FixUpTotalIncome error!(%v)", err)
	}
	c.JSON(nil, err)
}

func syncUpPGC(c *bm.Context) {
	err := svr.SyncUpPGC(c)
	if err != nil {
		log.Error("svr.SyncUpPGC error!(%v)", err)
	}
	c.JSON(nil, err)
}

func syncAvBaseIncome(c *bm.Context) {
	v := new(struct {
		Table string `form:"table" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	err := svr.SyncAvBaseIncome(c, v.Table)
	if err != nil {
		log.Error("svr.SyncAvBaseIncome error!(%v)", err)
	}
	c.JSON(nil, err)
}

func updateColumnTag(c *bm.Context) {
	v := new(struct {
		Table string `form:"table" validate:"required"`
		New   int    `form:"new" validate:"required"`
		Old   string `form:"old" validate:"required"`
		Count int64  `form:"count" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	err := svr.UpdateColumnTag(c, v.Table, v.Old, v.New, v.Count)
	if err != nil {
		log.Error("svr.UpdateColumnTag error!(%v)", err)
	}
	c.JSON(nil, err)
}

func syncCreditScore(c *bm.Context) {
	err := svr.SyncCreditScore(c)
	if err != nil {
		log.Error("svr.SyncCreditScore error!(%v)", err)
	}
	c.JSON(nil, err)
}

func calBgmStatis(c *bm.Context) {
	err := svr.FixBgmIncomeStatis(c)
	if err != nil {
		log.Error("svr.FixBgmIncomeStatis error!(%v)", err)
	}
	c.JSON(nil, err)
}

func calBgmBaseIncome(c *bm.Context) {
	v := new(struct {
		MID  int64  `form:"mid" validate:"required"`
		Date string `form:"date" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	err := svr.FixBgmBaseIncome(c, v.MID, v.Date)
	if err != nil {
		log.Error("svr.FixBgmIncomeStatis error!(%v)", err)
	}
	c.JSON(nil, err)
}

func updateBusinessIncome(c *bm.Context) {
	v := new(struct {
		Date string `form:"date" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	err := incomeSrv.UpdateBusinessIncome(c, v.Date)
	if err != nil {
		log.Error("svr.FixBusinessIncome error!(%v)", err)
	}
	c.JSON(nil, err)
}

func delDataLimit(c *bm.Context) {
	v := new(struct {
		Table string `form:"table" validate:"required"`
		Count int64  `form:"count" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	err := svr.DelDataLimit(c, v.Table, v.Count)
	if err != nil {
		log.Error("svr.DelDataLimit error(%v)", err)
	}
	c.JSON(nil, err)
}

func autoBreach(c *bm.Context) {
	v := new(struct {
		Date string `form:"date" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	_, err := svr.AutoBreach(c, v.Date)
	if err != nil {
		log.Error("svr.AutoBreach error(%v)", err)
	}
	c.JSON(nil, err)
}

func autoPunish(c *bm.Context) {
	_, err := svr.AutoPunish(c)
	if err != nil {
		log.Error("svr.AutoPunish error(%v)", err)
	}
	c.JSON(nil, err)
}

func autoExamination(c *bm.Context) {
	_, err := svr.AutoExamination(c)
	if err != nil {
		log.Error("svr.AutoExamination error(%v)", err)
	}
	c.JSON(nil, err)
}

func syncUpAccount(c *bm.Context) {
	err := svr.SyncUpAccount(c)
	if err != nil {
		log.Error("svr.SyncUpAccount error(%v)", err)
	}
	c.JSON(nil, err)
}
