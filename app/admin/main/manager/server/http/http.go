package http

import (
	"net/http"

	"go-common/app/admin/main/manager/conf"
	"go-common/app/admin/main/manager/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/rpc/warden"
)

var (
	mngSvc    *service.Service
	vfySvc    *verify.Verify
	permitSvc *permit.Permit
)

// Init init http sever instance.
func Init(c *conf.Config, s *service.Service) {
	mngSvc = s
	vfySvc = verify.New(nil)
	wardenConf := &warden.ClientConfig{
		NonBlock: true,
	}
	permitSvc = permit.New2(wardenConf)
	// init inner router
	engine := bm.DefaultServer(c.HTTPServer)
	innerRouter(engine)
	// init inner server
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

func innerRouter(e *bm.Engine) {
	e.Ping(ping)
	g := e.Group("/x/admin/manager")
	{
		g.GET("/auth", permitSvc.Verify2(), authUser)
		g.GET("/heartbeat", permitSvc.Verify2(), heartbeat)
		g.GET("/logout", logout)
		g.GET("/permission", vfySvc.Verify, permissions)
		log := g.Group("/log", permitSvc.Permit2(""))
		{
			log.GET("/audit", searchLogAudit)
			log.GET("/action", searchLogAction)
		}
		internal := g.Group("/internal", vfySvc.Verify)
		{
			internal.GET("/user/role", userRole)
			internal.GET("/business/role", roleList)
			internal.GET("/tag/list", tagList)
			internal.GET("/control/list", tagControl)
			internal.GET("/user/roles", userRoles)
			internal.GET("/user/state/up", stateUp)
			internal.GET("/user/list", userList)
			internal.GET("/roles", allRoles)
			internal.POST("/role/add", addRole)
		}
		gusers := g.Group("/users", vfySvc.Verify)
		{
			gusers.GET("/department/all", departments)
			gusers.GET("/role/all", roles)
			gusers.GET("/department/relation", usersByDepartment)
			gusers.GET("/role/relation", usersByRole)
			gusers.GET("", users)
			gusers.GET("/total", usersTotal)
			gusers.GET("/unames", usersNames)
			gusers.GET("/uids", userIds)
			gusers.GET("/udepts", usersDepts)
		}
		grank := g.Group("/rank", permitSvc.Verify2())
		{
			grank.GET("/group/info", rankGroup)
			grank.GET("/group/list", rankGroups)
			grank.POST("/group/add", addRankGroup)
			grank.POST("/group/update", updateRankGroup)
			grank.POST("/group/del", delRankGroup)
			grank.POST("/user/add", addRankUser)
			grank.GET("/user/list", rankUsers)
			grank.POST("/user/save", saveRankUser)
			grank.POST("/user/del", delRankUser)
		}
		tag := g.Group("/tag", permitSvc.Permit2(""))
		{
			tag.POST("/type/add", addType)
			tag.POST("/type/update", updateType)
			tag.POST("/type/delete", deleteType)
			tag.POST("/add", addTag)
			tag.POST("/update", updateTag)
			tag.POST("/batch/state/update", batchUpdateState)
			tag.POST("/control/add", addControl)
			tag.POST("/control/update", updateControl)
			tag.POST("/attr/update", attrUpdate)
			tag.GET("/attr/list", attrList)
			tag.GET("/list", tagList)
			tag.GET("/type/list", typeList)
			tag.GET("/control", tagControl)
		}
		reason := g.Group("/reason", permitSvc.Permit2(""))
		{
			reason.GET("/catesecext/list", cateSecExtList)
			reason.POST("/catesecext/add", addCateSecExt)
			reason.POST("/catesecext/update", updateCateSecExt)
			reason.POST("/catesecext/ban", banCateSecExt)
			reason.GET("/association/list", associationList)
			reason.POST("/association/add", addAssociation)
			reason.POST("/association/update", updateAssociation)
			reason.POST("/association/ban", banAssociation)
			reason.GET("/list", reasonList)
			reason.POST("/add", addReason)
			reason.POST("/update", updateReason)
			reason.POST("/batch/state/update", batchUpdateReasonState)
			reason.GET("/dropdown/list", dropDownList)
			reason.GET("/business/attr", businessAttr)
		}
		business := g.Group("/business", permitSvc.Permit2(""))
		{
			business.POST("/add", addBusiness)
			business.POST("/update", updateBusiness)
			business.POST("/role/add", addRole)
			business.POST("/role/update", updateRole)
			business.POST("/user/add", addUser)
			business.POST("/user/update", updateUser)
			business.POST("/state/update", updateState)
			business.POST("/user/delete", deleteUser)
			business.GET("/list", businessList)
			business.GET("/flow/list", flowList)
			business.GET("/user/list", userList)
			business.GET("/role/list", roleList)
			business.GET("/user/role", userRole)
		}
	}
}

// ping check server ok.
func ping(c *bm.Context) {
	var err error
	if err = mngSvc.Ping(c); err != nil {
		log.Error("service ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
