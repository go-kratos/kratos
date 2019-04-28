package http

import (
	"go-common/app/admin/main/creative/conf"
	"go-common/app/admin/main/creative/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
)

var (
	svc     *service.Service
	authSrc *permit.Permit
)

// Init http server
func Init(c *conf.Config) {
	svc = service.New(c)
	authSrc = permit.New(c.Auth)
	engine := bm.DefaultServer(c.BM)
	innerRouter(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

func innerRouter(e *bm.Engine) {
	e.GET("/monitor/ping", moPing)
	b := e.Group("/x/admin/creative")
	{
		innerMusicRouter(e, b)
		innerMaterialRouter(e, b)
		app := b.Group("/app")
		{
			app.GET("/portal", viewPortal)
			app.POST("/portal/add", addPortal)
			app.POST("/portal/update", upPortal)
			app.POST("/portal/down", downPortal)
			app.GET("/portal/list", portalList)
		}
		notice := b.Group("/notice")
		{
			notice.GET("/view", viewNotice)
			notice.GET("/list", listNotice)
			notice.POST("/add", addNotice)
			notice.POST("/update", upNotice)
			notice.POST("/delete", delNotice)
		}
		arc := b.Group("/oper/col_arc", authSrc.Verify())
		{
			arc.GET("/list", listCollectArcOper)
			arc.POST("/add", addCollectArcOper)
			// blew is copy from notice for support
			arc.GET("/view", viewNotice)
			arc.POST("/update", upNotice)
			arc.POST("/delete", delNotice)
		}
		whitelist := b.Group("/whitelist")
		{
			whitelist.GET("/view", viewWhiteList)
			whitelist.GET("/list", listWhiteList)
			whitelist.GET("/export.so", exportWhiteList)
			whitelist.POST("/add", addWhiteList)
			whitelist.POST("/add/batch", addBatchWhiteList)
			whitelist.POST("/update", upWhiteList)
			whitelist.POST("/delete", delWhiteList)
		}

		academy := b.Group("/academy", authSrc.Verify(), authSrc.Permit("ACADEMY_MANAGER")) //创作学院
		// academy := b.Group("/academy") //创作学院
		{
			academy.POST("/tag/update/fix", fixTag)         //清理脏数据
			academy.POST("/archive/update/fix", fixArchive) //清理脏数据
			//tag
			academy.POST("/tag/add", addTag)
			academy.POST("/tag/update", upTag)
			academy.POST("/tag/bind", bindTag)
			academy.GET("/tag/view", viewTag)
			academy.GET("/tag/list", listTag)
			academy.POST("/tag/order", orderTag)
			//archive
			academy.POST("/archive/add", addArc)
			academy.POST("/archive/update", upArcTag)
			academy.POST("/archive/remove", removeArcTag)
			academy.POST("/archive/batch/add", batchAddArc)
			academy.POST("/archive/batch/update", batchUpArc)
			academy.POST("/archive/batch/remove", batchRemoveArc)
			academy.GET("/archive/view", viewArc)
			academy.GET("/archive/list", listArc)
			//occupation & skill
			academy.POST("/occupation/add", addOccupation)
			academy.POST("/occupation/update", upOccupation)
			academy.POST("/occupation/bind", bindOccupation)
			academy.POST("/occupation/order", orderOccupation)
			academy.GET("/occupation/list", listOccupation)
			academy.POST("/skill/add", addSkill)
			academy.POST("/skill/update", upSkill)
			academy.POST("/skill/bind", bindSkill)
			academy.POST("/software/add", addSoftware)
			academy.POST("/software/update", upSoftware)
			academy.POST("/software/bind", bindSoftware)
			//arc & skill
			academy.GET("/skill/archive/view", viewSkillArc)
			academy.GET("/skill/archive/list", skillArcList)
			academy.POST("/skill/archive/add", addSkillArc)
			academy.POST("/skill/archive/update", upSkillArc)
			academy.POST("/skill/archive/bind", bindSkillArc)
			//search keywords
			academy.GET("/search/keywords", searchKeywords)
			academy.POST("/search/keywords/sub", subSearchKeywords)
		}

		task := b.Group("/task", authSrc.Verify(), authSrc.Permit("CREATIVE_TASK_MANAGER")) //任务系统
		// task := b.Group("/task") //任务系统
		{
			//task list
			task.GET("/pre", taskPre)
			task.GET("/list", taskList)
			task.POST("/online", batchOnline)
			//group
			task.GET("/group/view", viewGroup)
			task.POST("/group/add", addGroup)
			task.POST("/group/edit", editGroup)
			task.POST("/group/order", orderGroup)
			task.POST("/group/upstate", upStateGroup)
			//sub
			task.GET("/sub/view", viewSubtask)
			task.POST("/sub/add", addSubtask)
			task.POST("/sub/edit", editSubtask)
			task.POST("/sub/order", orderSubtask)
			task.POST("/sub/upstate", upStateSubtask)
			task.POST("/sub/transfer", transferSubtask)
			//reward
			task.GET("/reward/list", rewardTree)
			task.GET("/reward/view", viewReward)
			task.POST("/reward/add", addReward)
			task.POST("/reward/edit", editReward)
			task.POST("/reward/upstate", upStateReward)
			//gift
			task.GET("/gift/list", listGiftReward)
			task.GET("/gift/view", viewGiftReward)
			task.POST("/gift/add", addGiftReward)
			task.POST("/gift/edit", editGiftReward)
			task.POST("/gift/upstate", upStateGiftReward)
		}
	}
}

//素材库 db creative
func innerMaterialRouter(e *bm.Engine, group *bm.RouterGroup) {
	if group == nil {
		return
	}
	//字幕库 字体库 滤镜库
	material := group.Group("/material", authSrc.Verify())
	{
		material.GET("", authSrc.Permit("MATERIAL_READ"), infoMaterial)
		material.GET("/search", authSrc.Permit("MATERIAL_READ"), searchMaterialDb)
		//支持新增和修改
		material.POST("/add", authSrc.Permit("MATERIAL_WRITE"), syncMaterial)
		//支持批量修改
		material.POST("/state", authSrc.Permit("MATERIAL_WRITE"), stateMaterial)
		//仅支持 image/ zip 上传
		material.POST("/upload", upload)
		//素材库分类
		material.GET("/category", authSrc.Permit("MATERIAL_READ"), category)
		material.POST("/category/add", authSrc.Permit("MATERIAL_WRITE"), addMCategory)
		material.POST("/category/edit", authSrc.Permit("MATERIAL_WRITE"), editMCategory)
		material.POST("/category/index", authSrc.Permit("MATERIAL_WRITE"), indexMCategory)
		material.POST("/category/delete", authSrc.Permit("MATERIAL_WRITE"), delMCategory)
		material.GET("/category/search", authSrc.Permit("MATERIAL_READ"), searchMCategory)
	}
}

//音频库 db archive
func innerMusicRouter(e *bm.Engine, group *bm.RouterGroup) {
	if group == nil {
		return
	}
	groupMusic := group.Group("/music")
	{
		//音乐管理及同步
		groupMusic.POST("/add", authSrc.Permit("MUSIC_UPDATE"), syncMusic)
		groupMusic.POST("/up/frontname", authSrc.Permit("MUSIC_UPDATE"), editMusicFrontName)
		groupMusic.POST("/up/tags", authSrc.Permit("MUSIC_UPDATE"), editMusicTags)
		groupMusic.POST("/edit", authSrc.Permit("MUSIC_UPDATE"), editMusic)
		groupMusic.POST("/batch/tags", authSrc.Permit("MUSIC_UPDATE"), batchEditMusicTags)
		groupMusic.POST("/up/timeline", authSrc.Permit("MUSIC_UPDATE"), editMusicTimeline)
		groupMusic.GET("/search", authSrc.Permit("MUSIC_READ"), searchMusic)
		//音乐分类
		groupMusic.GET("/category", authSrc.Permit("MUSIC_CATEGORY_READ"), categoryInfo)
		groupMusic.POST("/category/add", authSrc.Permit("MUSIC_CATEGORY_UPDATE"), addCategory)
		groupMusic.POST("/category/edit", authSrc.Permit("MUSIC_CATEGORY_UPDATE"), editCategory)
		groupMusic.POST("/category/index", authSrc.Permit("MUSIC_CATEGORY_UPDATE"), indexCategory)
		groupMusic.POST("/category/delete", authSrc.Permit("MUSIC_CATEGORY_UPDATE"), delCategory)
		groupMusic.GET("/category/search", authSrc.Permit("MUSIC_CATEGORY_READ"), searchCategory)
		//素材分类
		groupMusic.GET("/material", authSrc.Permit("MUSIC_MATERIAL_READ"), materialInfo)
		groupMusic.POST("/material/add", authSrc.Permit("MUSIC_MATERIAL_UPDATE"), addMaterial)
		groupMusic.POST("/material/edit", authSrc.Permit("MUSIC_MATERIAL_UPDATE"), editMaterial)
		groupMusic.POST("/material/delete", authSrc.Permit("MUSIC_MATERIAL_UPDATE"), delMaterial)
		groupMusic.POST("/material/batch/delete", authSrc.Permit("MUSIC_MATERIAL_UPDATE"), batchDeleteMaterial)
		groupMusic.GET("/material/search", authSrc.Permit("MUSIC_MATERIAL_READ"), searchMaterial)
		//音乐及素材 管理端
		groupMusic.GET("/material/relation", authSrc.Permit("MUSIC_WITH_MATERIAL_READ"), musicMaterialRelationInfo)
		groupMusic.POST("/material/relation/add", authSrc.Permit("MUSIC_WITH_MATERIAL_UPDATE"), addMaterialRelation)
		groupMusic.POST("/material/relation/batch/add", authSrc.Permit("MUSIC_WITH_MATERIAL_UPDATE"), batchAddMaterialRelation)
		groupMusic.POST("/material/relation/edit", authSrc.Permit("MUSIC_WITH_MATERIAL_UPDATE"), editMaterialRelation)
		//音乐及分类 app端
		groupMusic.GET("/category/relation", authSrc.Permit("MUSIC_WITH_CATEGORY_READ"), musicCategoryRelationInfo)
		groupMusic.POST("/category/relation/add", authSrc.Permit("MUSIC_WITH_CATEGORY_UPDATE"), addCategoryRelation)
		groupMusic.POST("/category/relation/batch/add", authSrc.Permit("MUSIC_WITH_CATEGORY_UPDATE"), batchAddCategoryRelation)
		groupMusic.POST("/category/relation/index", authSrc.Permit("MUSIC_WITH_CATEGORY_UPDATE"), indexCategoryRelation)
		groupMusic.POST("/category/relation/edit", authSrc.Permit("MUSIC_WITH_CATEGORY_UPDATE"), editCategoryRelation)
		groupMusic.POST("/category/relation/batch/delete", authSrc.Permit("MUSIC_WITH_CATEGORY_UPDATE"), batchDeleteCategoryRelation)
		groupMusic.POST("/category/relation/delete", authSrc.Permit("MUSIC_WITH_CATEGORY_UPDATE"), delCategoryRelation)
		groupMusic.GET("/category/relation/search", authSrc.Permit("MUSIC_WITH_CATEGORY_READ"), searchCategoryRelation)
	}

}
