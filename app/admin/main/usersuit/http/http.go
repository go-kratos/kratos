package http

import (
	"net/http"

	"go-common/app/admin/main/usersuit/conf"
	"go-common/app/admin/main/usersuit/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
)

var (
	svc     *service.Service
	authSvc *permit.Permit
)

// Init fot init open service
func Init(c *conf.Config, s *service.Service) {
	svc = s
	authSvc = permit.New(c.Auth)
	engine := bm.DefaultServer(c.BM)
	innerRouter(engine)
	// init internal server
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

// innerRouter
func innerRouter(e *bm.Engine) {
	// ping monitor
	e.Ping(ping)
	// internal api
	bg := e.Group("/x/admin/usersuit")
	{
		// invite
		bg.POST("/invite/generate", authSvc.Permit("INVITE_INFO_ADMIN"), generate)
		bg.GET("/invite/list", authSvc.Permit("INVITE_INFO_SEARCH"), list)

		// new pendant
		bg.GET("/pendant/info/list", authSvc.Permit("PENDANT_INFO_LIST"), pendantInfoList)                    // 挂件列表(带分页)
		bg.GET("/pendant/info/id", authSvc.Permit("PENDANT_INFO_ID"), pendantInfoID)                          // 挂件详情
		bg.GET("/pendant/group/id", authSvc.Permit("PENDANT_GROUP_ID"), pendantGroupID)                       // 挂件分组详情
		bg.GET("/pendant/group/list", authSvc.Permit("PENDANT_GROUP_LIST"), pendantGroupList)                 // 挂件分组列表(带分页)
		bg.GET("/pendant/group/all", authSvc.Permit("PENDANT_GROUP_ALL"), pendantGroupAll)                    // 挂件分组列表(不带分页) 这个筛选列表可不加权限
		bg.GET("/pendant/info/all/no/page", authSvc.Permit("PENDANT_INFO_ALL_NO_PAGE"), pendantInfoAllNoPage) // 在售挂件(不带分页) 这个筛选列表可不加权限
		bg.POST("/pendant/add/info", authSvc.Permit("PENDANT_ADD_INFO"), addPendantInfo)                      // 添加挂件
		bg.POST("/pendant/up/info", authSvc.Permit("PENDANT_UP_INFO"), upPendantInfo)                         // 更新挂件
		bg.POST("/pendant/up/group/status", authSvc.Permit("PENDANT_UP_GROUP_STATUS"), upPendantGroupStatus)  // 更新挂件分组状态
		bg.POST("/pendant/up/info/status", authSvc.Permit("PENDANT_UP_INFO_STATUS"), upPendantInfoStatus)     // 更新挂件状态
		bg.POST("/pendant/add/group", authSvc.Permit("PENDANT_ADD_GROUP"), addPendantGroup)                   // 添加挂件分组
		bg.POST("/pendant/up/group", authSvc.Permit("PENDANT_UP_GROUP"), upPendantGroup)                      // 更新挂件分组
		bg.GET("/pendant/orders", authSvc.Permit("PENDANT_ORDERS"), pendantOrders)                            // 挂件订单列表
		bg.POST("/pendant/equip", authSvc.Permit("PENDANT_EQUIP"), equipPendant)                              // 装备挂件
		bg.GET("/user/pendant/pkg", authSvc.Permit("USER_PENDANT_PKG"), userPendantPKG)                       // 用户可用挂件背包列表
		bg.GET("/user/pkg/details", authSvc.Permit("USER_PKG_DETAILS"), userPKGDetails)                       // 用户所属挂件详情
		bg.POST("/pendant/up/pkg", authSvc.Permit("PENDANT_UP_PKG"), upPendantPKG)                            // 更新或激活用户挂件
		bg.POST("/pendant/mutli/send", authSvc.Permit("PENDANT_MUTLI_SEND"), mutliSend)                       // 多用户发送挂件
		bg.GET("/pendant/oper/log", authSvc.Permit("PENDANT_OPER_LOG"), pendantOperlog)                       // 发放挂件操作日志
		bg.POST("/face/upload", upload)                                                                       // 上传图片
		//medal
		bg.GET("/medal", authSvc.Permit("MEDAL_INFO_LIST"), medalList)                                // 勋章列表
		bg.GET("/medal/id", authSvc.Permit("MEDAL_INFO_ID"), medalView)                               // 勋章查看
		bg.POST("/medal/add", authSvc.Permit("MEDAL_INFO_ADD"), medalAdd)                             // 勋章添加
		bg.POST("/medal/edit", authSvc.Permit("MEDAL_INFO_EDIT"), medalEdit)                          // 勋章编辑
		bg.GET("/medal/group", authSvc.Permit("MEDAL_GROUP_LIST"), medalGroup)                        // 勋章分组列表
		bg.GET("/medal/group/view", authSvc.Permit("MEDAL_GROUP_VIEW"), medalGroupView)               // 勋章分组查看单个
		bg.GET("/medal/group/parent", authSvc.Permit("MEDAL_GROUP_PARENT"), medalGroupParent)         // 勋章父辈分组列表
		bg.POST("/medal/group/add", authSvc.Permit("MEDAL_GROUP_ADD"), medalGroupAdd)                 // 勋章分组添加
		bg.POST("/medal/group/edit", authSvc.Permit("MEDAL_GROUP_EDIT"), medalGroupEdit)              // 勋章分组编辑
		bg.GET("/medal/member/mid", authSvc.Permit("MEDAL_MEMBER_MID"), medalMemberMID)               // 勋章会员管理
		bg.POST("/medal/member/edit", authSvc.Permit("MEDAL_MEMBER_EDIT"), medalOwnerUpActivated)     // 勋章会员编辑激活
		bg.GET("/medal/member/add/list", authSvc.Permit("MEDAL_MEMBER_ADD_LIST"), medalMemberAddList) // 勋章会员可添加列表
		bg.POST("/medal/member/add", authSvc.Permit("MEDAL_MEMBER_ADD"), medalMemberAdd)              // 勋章会员添加
		bg.POST("/medal/member/del", authSvc.Permit("MEDAL_MEMBER_DEL"), medalMemberDel)              // 勋章会员删除
		bg.POST("/medal/batch/add", authSvc.Permit("MEDAL_BATCH_ADD"), medalBatchAdd)                 // 勋章批量添加
		bg.GET("/medal/oper/log", authSvc.Permit("MEDAL_OPER_LOG"), medalOperlog)                     // 勋章操作日志

	}
}

// ping check server ok.
func ping(c *bm.Context) {
	err := svc.Ping(c)
	if err != nil {
		log.Error("usersuit admin ping error")
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
