package http

import (
	"go-common/app/service/main/videoup/conf"
	"go-common/app/service/main/videoup/service"
	"go-common/library/log"
	"go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/antispam"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	vfySvc *verify.Verify
	anti   *antispam.Antispam
	vdpSvc *service.Service
	config *conf.Config
)

// Init init server.
func Init(c *conf.Config, s *service.Service) {
	vfySvc = verify.New(nil)
	anti = antispam.New(c.AntispamConf)
	vdpSvc = s
	config = c
	eng := blademaster.DefaultServer(c.BM)
	route(eng)

	if err := eng.Start(); err != nil {
		log.Error("eng.Start error(%v)", err)
		panic(err)
	}
}

func route(e *blademaster.Engine) {
	e.Ping(ping)

	vp := e.Group("/videoup", vfySvc.Verify)
	{
		// ugc && pgc
		vp.GET("/simplearchive", simpleArchive)
		vp.GET("/simplevideos", simpleVideos)
		vp.GET("/view", viewArchive)
		vp.GET("/views", viewArchives)
		vp.GET("/up/archives", archivesByMid)
		vp.GET("/cid/archives", archivesByCids)
		vp.GET("/archives/rejected", rejectedArchives)
		vp.GET("/history/view", arcHistory)
		vp.GET("/history/list", arcHistorys)
		vp.GET("/types", types)
		vp.GET("/flows", flows)
		vp.GET("/flow/forbid", flowForbid)
		vp.GET("/query/cid", queryCid)
		vp.GET("/up/special", specialUps)
		vp.GET("/feed/aids", appFeedAids)
		vp.GET("/desc/format", descFormats)
		vp.POST("/archive/report", anti.Handler(), arcReport)
		vp.GET("/archive/reason/tag", arcReasonTag)
		vp.GET("/archive/addit", archiveAddit)
		vp.POST("/del", delArchive)
		// recommend archive
		vp.GET("/recos", Recos)
		vp.POST("/reco/update", RecoUpdate)
		// obtain cid
		vp.POST("/obtain/cid", obtainCid)
		// ugc
		vp.GET("/cid", videoBycid)
		vp.POST("/add", addArchive)
		vp.POST("/edit", editArchive)
		vp.POST("/tag/up", upArchiveTag)

		//setting output
		setting := vp.Group("/setting")
		{
			setting.GET("/dynamic", queryDynamic)
		}

		//ugc only
		ugc := vp.Group("/ugc")
		{
			ugc.POST("/edit/mission", editMissionByUGC)
		}
		// pgc
		pgc := vp.Group("/pgc")
		{
			pgc.POST("/add", addByPGC)
			pgc.POST("/edit", editByPGC)
			pgc.POST("/add/secret", saddByPGC)
			pgc.POST("/add/coopera", caddByPGC)
			pgc.POST("/edit/coopera", ceditByPGC)
		}
		vp.POST("/ns/md5", nsMd5)
		// Get video traffic jam level
		vp.GET("/video/jam", videoJam)
		//ad
		porder := vp.Group("/porder")
		{
			porder.GET("/config/list", porderCfgList)
			porder.GET("/arc/list", porderArcList)
		}
		flow := vp.Group("/flow")
		{
			flow.POST("/entry/mid", addByMid)
			flow.POST("/entry/oid", addByOid)
			//分页查询
			flow.GET("/list", listFlows)
			//指定查询
			flow.GET("/list/judge", listJudgeFlows)
			flow.GET("/info", queryByOid)
		}

		staff := vp.Group("/staff")
		{
			/*
					1. up主申请(添加/修改) staff member  edit其中涉及到 可能已经accept 变成 staffON 的case
					2. 生成多个apply list
					3. 受邀请者 查看 applys
					4.同意  accept  拒绝  refuse
					5.up主自己可以继续编辑  覆盖式 ！！！
				 	发起人维度 稿件页 是批量的
				   	具体申请单页 是单个操作
			*/
			//查看apply
			staff.GET("/apply", viewApply)
			//mid applys 计数
			staff.GET("/mid/applys", checkMid)
			//管理后台
			staff.POST("/apply/batch", batchApplys)
			//apply申请单交互
			staff.POST("/apply/submit", addApply)
			//apply申请单列表
			staff.GET("/apply/list", applys)
			staff.GET("/apply/filter", filterApplys)
			staff.GET("/archive/applys/", archiveApplys)
			//staff正式列表
			staff.GET("", staffs)
		}
	}
}
