package http

import (
	"net/http"
	"strconv"

	"go-common/app/admin/main/tv/conf"
	"go-common/app/admin/main/tv/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/http/blademaster/render"
)

var (
	tvSrv   *service.Service
	vfySvc  *verify.Verify
	authSvc *permit.Permit
)

const (
	_errIDNotFound = "ids not found"
	_errTitleExist = "Title exists already"
)

// Init init http sever instance.
func Init(c *conf.Config, s *service.Service) {
	initService(c, s)
	engine := bm.DefaultServer(c.HTTPServer)
	innerRouter(engine)
	// init internal server
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

// initService init service
func initService(c *conf.Config, s *service.Service) {
	tvSrv = s
	vfySvc = verify.New(nil)
	authSvc = permit.New(c.Auth)
}

func parseInt(value string) int64 {
	intval, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		intval = 0
	}
	return intval
}

func atoi(value string) (intval int) {
	intval, err := strconv.Atoi(value)
	if err != nil {
		intval = 0
	}
	return intval
}

// innerRouter
func innerRouter(e *bm.Engine) {
	// ping monitor
	e.GET("/monitor/ping", ping)
	// internal api
	bg := e.Group("/x/admin/tv")
	{
		// cms content edit
		cont := bg.Group("/cont", vfySvc.Verify)
		{
			cont.POST("/online", online)
			cont.POST("/hidden", hidden)
		}
		// pgc ep inject
		epIn := bg.Group("/ep", vfySvc.Verify)
		{
			epIn.POST("/create", createEP)
			epIn.POST("/remove", removeEP)
		}
		// pgc season inject
		snIn := bg.Group("/season", vfySvc.Verify)
		{
			snIn.POST("/create", createSeason)
			snIn.POST("/remove", removeSeason)
		}
		// intervsRank edit
		interv := bg.Group("/intervs", vfySvc.Verify)
		{
			rank := interv.Group("/rank")
			{
				rank.GET("/list", intervsRank)
				rank.POST("/publish", rankPub)
			}
			module := interv.Group("/module")
			{
				module.GET("/list", intervsMod)
				module.POST("/publish", modPub)
			}
			index := interv.Group("/index")
			{
				index.GET("/list", intervsIndex)
				index.POST("/publish", indexPub)
			}
		}
		// audit result
		audit := bg.Group("/audit_result")
		{
			aud := audit.Group("")
			{
				aud.GET("/ep", epResult)
				aud.GET("/season", seasonResult)
				aud.GET("/archive", arcResult)
				aud.GET("/video", videoResult)
				aud.GET("/ugctypes", auditCategory)
				audit.GET("/abnor_export", abnorExport)
				audit.GET("/abnor_debug", abnorDebug)
			}
			audit.POST("/unshelve", authSvc.Permit("TV_MEDIA_DEL"), unShelve)
		}
		// content repository
		crepo := bg.Group("/contrepo", vfySvc.Verify)
		{
			crepo.GET("/list", contList)
			crepo.GET("/info", contInfo)
			crepo.POST("/save", saveCont)
			crepo.GET("/preview", preview)
			crepo.POST("/online", contOnline)
			crepo.POST("/hidden", contHidden)
			crepo.POST("/upload", upbfs)
		}
		// season repo
		srepo := bg.Group("/searepo", vfySvc.Verify)
		{
			srepo.GET("/list", seasonList)
			srepo.GET("/info", seasonInfo)
			srepo.POST("/save", saveSeason)
			srepo.POST("/online", seasonOnline)
			srepo.POST("/hidden", seasonHidden)
			// ugc
			crugc := srepo.Group("/ugc")
			{
				//archive
				crugc.GET("/archive/lists", arcList)
				crugc.GET("/archive/category", arcCategory)
				crugc.POST("/archive/online", arcOnline)
				crugc.POST("/archive/hidden", arcHidden)
				crugc.GET("/archive/arcTypeRPC", arcTypeRPC)
				crugc.POST("/archive/update", arcUpdate)
				//video
				crugc.GET("/video/lists", VideoList)
				crugc.POST("/video/online", VideoOnline)
				crugc.POST("/video/hidden", VideoHidden)
				crugc.GET("/video/preview", VideoPreview)
				crugc.POST("/video/update", videoUpdate)
			}
		}
		// version mgt
		ver := bg.Group("/version", vfySvc.Verify)
		{
			ver.GET("/list", versionList)
			ver.GET("/info", versionInfo)
			ver.POST("/save", saveVersion)
			ver.POST("/add", addVersion)
			ver.POST("/delete", versionDel)
		}
		// version update mgt
		verup := bg.Group("/verupdate", vfySvc.Verify)
		{
			verup.GET("/list", verUpdateList)
			verup.POST("/add", addVerUpdate)
			verup.POST("/save", saveVerUpdate)
			verup.POST("/enable", verUpdateEnable)
			verup.GET("/full", fullPackageImport)
		}
		// channel mgt
		chl := bg.Group("/channel", vfySvc.Verify)
		{
			chl.GET("/list", chlList)
			chl.GET("/info", chlInfo)
			chl.POST("/edit", chlEdit)
			chl.POST("/add", chlAdd)
			chl.POST("/delete", chlDel)
		}
		upper := bg.Group("/upper", authSvc.Permit("TV_AUDIT_MGT"))
		{
			upper.POST("/add", upAdd)
			upper.POST("/import", upImport)
			upper.POST("/del", upDel)
			upper.GET("", upList)
		}
		upCMS := bg.Group("upcms", vfySvc.Verify)
		{
			upCMS.GET("/list", upcmsList)
			upCMS.POST("/audit", upcmsAudit)
			upCMS.POST("/edit", upcmsEdit)
		}
		//search intervene
		si := bg.Group("/searchInter", vfySvc.Verify)
		{
			si.GET("/list", searInterList)
			si.POST("/add", searInterAdd)
			si.POST("/edit", searInterEdit)
			si.POST("/delete", searInterDel)
			si.POST("/rank", searInterRank)
			si.POST("/publish", searInterPublish)
			si.POST("/publishList", searInterPubList)
		}
		bg.POST("/archive/add", authSvc.Permit("TV_AUDIT_MGT"), arcAdd)
		//modules manager
		mod := bg.Group("/modules", vfySvc.Verify)
		{
			mod.POST("/add", modulesAdd)
			mod.GET("/list", modulesList)
			mod.GET("/editGet", modulesEditGet)
			mod.POST("/editPost", modulesEditPost)
			mod.POST("/publish", modulesPublish)
			mod.GET("/sup_cat", supCat)
		}
		//watermark
		wr := bg.Group("/watermark", authSvc.Permit("TV_AUDIT_MGT"))
		{
			wr.GET("/list", waterMarklist)
			wr.POST("/add", waterMarkAdd)
			wr.POST("/delete", waterMarkDelete)
		}
		mango := bg.Group("/mango", authSvc.Permit("TV_AUDIT_MGT"))
		{
			mango.GET("/list", mangoList)
			mango.POST("/add", mangoAdd)
			mango.POST("/edit", mangoEdit)
			mango.POST("/delete", mangoDel)
			mango.POST("/publish", mangoPub)
		}
		trans := bg.Group("/trans", authSvc.Permit("TV_AUDIT_MGT"))
		{
			trans.GET("/list", transList)
		}
		label := bg.Group("/label", vfySvc.Verify)
		{
			label.POST("/act", actLabels)
			label.POST("/edit", editLabel)
			label.POST("/publish", pubLabel)
			ugcLabel := label.Group("/ugc")
			{
				ugcLabel.POST("/add_time", addTime)
				ugcLabel.POST("/edit_time", editTime)
				ugcLabel.POST("/del_time", delTmLabels)
				ugcLabel.GET("/list", ugcLabels)
			}
			pgcLabel := label.Group("/pgc")
			{
				pgcLabel.GET("/list", pgcLabels)
				pgcLabel.GET("/types", pgcLblTps)
			}
		}
		// app manager
		app := bg.Group("/app", vfySvc.Verify)
		{
			// region manager
			reg := app.Group("/region")
			{
				reg.GET("/list", reglist)
				reg.POST("/sort", regSort)
				reg.POST("/save", saveReg)
				reg.POST("/publish", upState)
			}
		}
		// vip tv-vip
		vip := bg.Group("/vip", authSvc.Permit("TV_VIP"))
		{
			// user info
			vip.GET("/user/info", userInfo)
			//order info
			vip.GET("/order/list", orderList)
			//panel info
			panel := vip.Group("/panel")
			{
				panel.GET("/info", panelInfo)
				panel.POST("/status", panelStatus)
				panel.POST("/save", savePanel)
				panel.GET("/list", panelList)
			}
		}
	}
}

// ping check server ok.
func ping(c *bm.Context) {
}

func renderErrMsg(c *bm.Context, code int, msg string) {
	data := map[string]interface{}{
		"code":    code,
		"message": msg,
	}
	c.Render(http.StatusOK, render.MapJSON(data))
}
