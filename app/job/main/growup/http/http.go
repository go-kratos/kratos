package http

import (
	"net/http"

	"go-common/app/job/main/growup/conf"
	"go-common/app/job/main/growup/service"
	"go-common/app/job/main/growup/service/charge"
	"go-common/app/job/main/growup/service/ctrl"
	"go-common/app/job/main/growup/service/income"
	"go-common/app/job/main/growup/service/tag"

	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var (
	executor  *ctrl.UnboundedExecutor
	svr       *service.Service
	tagSvr    *tag.Service
	incomeSrv *income.Service
	chargeSrv *charge.Service
)

// Init init
func Init(c *conf.Config) {
	// service
	executor = ctrl.NewUnboundedExecutor()
	svr = service.New(c, executor)
	tagSvr = tag.New(c, executor)
	incomeSrv = income.New(c, executor)
	chargeSrv = charge.New(c, executor)

	// bm
	engine := bm.DefaultServer(c.BM)
	initRouter(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

// Close close service
func Close() {
	svr.Close()
	tagSvr.Close()
}

func initRouter(e *bm.Engine) {
	e.Ping(ping)
	mr := e.Group("/x/internal/job/growup")
	tag := mr.Group("tag")
	{
		tag.POST("/income", tagIncome)
		tag.POST("/archive/ratio", tagRatio)
		tag.POST("/ups", tagUps)
		tag.POST("/extra", tagExtraIncome)
	}
	mr.GET("/email/combine", combineMails)
	mr.GET("/email/tagincome", sendTagIncome)

	mr.POST("/blacklist/init/mid", initBlacklistMID)
	mr.POST("/blacklist/update", updateBlacklist)

	mr.POST("/data/update/withdraw", updateWithdraw)
	mr.GET("/data/up/income/statis", getUpIncomeStatis)
	mr.GET("/data/av/income/statis", getAvIncomeStatis)
	mr.GET("/cheat", updateCheat)
	mr.POST("/data/fix/tag", updateTagIncome)
	mr.POST("/data/fix/upincome", fixUpIncome)
	mr.POST("/data/fix/up/av/statis", fixUpAvStatis)
	mr.POST("/data/fix/tag/adjust", fixTagAdjust)
	mr.POST("/data/fix/income", fixIncome)
	mr.POST("/data/fix/account/type", fixAccountType)
	mr.POST("/data/fix/upaccount", fixUpAccount)
	mr.POST("/data/fix/baseincome", fixBaseIncome)
	mr.POST("/data/fix/av/breach", fixAvBreach)
	mr.POST("/data/fix/up/totalincome", fixUpTotalIncome)
	mr.POST("/data/fix/up/business", updateBusinessIncome)
	mr.POST("/data/up/sync/pgc", syncUpPGC)
	mr.POST("/data/up/sync/avbaseincome", syncAvBaseIncome)
	mr.POST("/data/column/tag", updateColumnTag)
	mr.POST("/data/del", delDataLimit)
	creative := mr.Group("/creative")
	{
		creative.POST("/income", creativeIncome)
		creative.POST("/charge", creativeCharge)
		creative.POST("/statis", creativeStatis)
		creative.POST("/bill", creativeBill)
		creative.POST("/budget", creativeBudget)
		creative.POST("/activity", creativeActivity)
	}
	mr.POST("/up/info/update", updateUpInfoVideo)
	mr.POST("/credit/sync", syncCreditScore)
	mr.POST("/bgm/sync", syncBGM)
	mr.POST("/bgm/statis", calBgmStatis)
	mr.POST("/bgm/base", calBgmBaseIncome)
	mr.POST("/sync/account", syncUpAccount)
	auto := mr.Group("/auto")
	{
		auto.POST("/archive/breach", autoBreach)
		auto.POST("/up/punish", autoPunish)
		auto.POST("/up/examination", autoExamination)
	}

	taskStatus := mr.Group("/task")
	{
		taskStatus.POST("/status", updateTaskStatus)
		taskStatus.POST("/column", checkTaskColumn)
	}

	incomeAdjust := mr.Group("/income_bubble")
	{
		incomeAdjust.POST("/meta/sync", syncBubbleMeta)
		incomeAdjust.POST("/task/meta", syncBubbleMetaTask)
		incomeAdjust.POST("/task/snapshot", snapshotBubbleTask)
	}
	// delete
	mr.GET("/avratio", execAvRatio)
	mr.GET("/income", execIncome)
}

// ping check server ok.
func ping(c *bm.Context) {
	if err := svr.Ping(c); err != nil {
		log.Error("service ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
