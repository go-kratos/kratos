package http

import (
	"net/http"

	"go-common/app/admin/ep/merlin/conf"
	"go-common/app/admin/ep/merlin/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
)

const (
	_sessIDKey = "_AJSESSIONID"
)

var (
	svc     *service.Service
	authSvc *permit.Permit
)

// Init init
func Init(c *conf.Config, s *service.Service) {
	svc = s
	authSvc = permit.New(c.Auth)

	engine := bm.DefaultServer(c.BM)
	engine.Ping(ping)
	outerRouter(engine)
	if err := engine.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
}

// outerRouter init outer router api path.
func outerRouter(e *bm.Engine) {
	e.GET("/ep/admin/merlin/version", getVersion)
	e.GET("/ep/admin/merlin/conf/version", confVersion)

	base := e.Group("/ep/admin/merlin", authSvc.Permit(""))
	{
		v1 := base.Group("/v1")
		{
			cluster := v1.Group("/cluster")
			{
				cluster.GET("/query", queryCluster)
			}

			machine := v1.Group("/machine")
			{
				machine.POST("/gen", genMachines)
				machine.GET("/del", delMachine)
				machine.GET("/query/detail", queryMachineDetail)
				machine.GET("/query", queryMachines)
				machine.GET("/query/status", queryMachineStatus)
				machine.GET("/transfer", transferMachine)

				machinePackage := machine.Group("/package")
				{
					machinePackage.GET("/query", queryMachinePackages)
				}

				machineLog := machine.Group("/log")
				{
					machineLog.GET("/query", queryMachineLogs)
				}

				machineNode := machine.Group("/node")
				{
					machineNode.POST("update", updateNodes)
					machineNode.GET("query", queryNodes)
				}
			}

			image := v1.Group("image")
			{
				image.GET("/query", queryImage)
				image.POST("/add", addImage)
				image.POST("/update", updateImage)
				image.POST("/del", delImage)
			}

			serviceTree := v1.Group("/tree")
			{
				serviceTree.GET("/query", userTree)
				serviceTree.GET("/container/query", userTreeContainer)
				serviceTree.GET("/auditors/query", treeAuditors)

			}

			audit := v1.Group("/audit")
			{
				auditEndTime := audit.Group("/endTime")
				{
					auditEndTime.GET("/delay", delayMachineEndTime)                               //手动延期 done ok
					auditEndTime.POST("/apply", applyMachineEndTime)                              //申请延期 done
					auditEndTime.GET("/cancel", cancelMachineEndTime)                             //取消延期 done
					auditEndTime.POST("/audit", auditMachineEndTime)                              //审批 通过或驳回 done
					auditEndTime.GET("/query/applyList", queryApplicationRecordsByMachineID)      //done ok
					auditEndTime.GET("/query/user/applyList", queryApplicationRecordsByApplicant) //done ok
					auditEndTime.GET("/query/user/auditList", queryApplicationRecordsByAuditor)   //done ok
				}

			}

			user := v1.Group("/user")
			{
				user.GET("/query", queryUserInfo)
			}

			mobileDevice := v1.Group("/mobiledevice")
			{
				mobileDevice.POST("/query", queryMobileDevice)
				mobileDevice.GET("/refresh", refreshMobileDeviceDetail)
				mobileDevice.GET("/bind", bindMobileDevice)
				mobileDevice.GET("/release", releaseMobileDevice)
				mobileDevice.GET("/isbind", isBindByTheUser)

				mobileDevice.GET("/pullout", lendOutMobileDevice)
				mobileDevice.GET("/return", returnMobileDevice)

				mobileDevice.GET("/start", startMobileDevice)
				mobileDevice.GET("/shutdown", shutDownMobileDevice)
				mobileDevice.GET("/syncall", syncMobileDevice)

				mobileDevice.GET("/category/query", queryCategory)
				mobileDevice.GET("/superuser/query", queryDeviceFarmSuperUser)

				mobileDeviceLog := mobileDevice.Group("/log")
				{
					mobileDeviceLog.GET("/query", queryMobileMachineLogs)
					mobileDeviceLog.GET("/lendout/query", queryMobileMachineLendOut)
				}

				mobileDeviceErrorLog := mobileDevice.Group("/error/log")
				{
					mobileDeviceErrorLog.GET("/query", queryMobileMachineErrorLogs)
					mobileDeviceErrorLog.POST("/report", reportMobileDeviceError)
				}
			}

			biliHub := v1.Group("/bilihub")
			{
				biliHub.GET("/auth", authHub)
				biliHub.GET("/auth/access", accessAuthHub)

				biliHub.GET("/projects/accesspull", accessPullProjects)
				biliHub.GET("/projects", projects)

				biliHub.GET("/repos", repos)
				biliHub.GET("/repotags", tags)
				biliHub.GET("/repos/delete", deleteRepo)
				biliHub.GET("/repotags/delete", deleteRepoTag)

				biliHub.GET("/snapshot", snapshot)
				biliHub.GET("/snapshot/query", querySnapshot)

				biliHub.POST("/machine2image", machine2image)
				biliHub.GET("/machine2image/forcefailed", machine2imageForceFailed)
				biliHub.GET("/machine2image/log/query", queryMachine2ImageLog)

				image := biliHub.Group("/image")
				{
					image.GET("/all", allImage)
					image.POST("/addtag", addTag)
					image.POST("/push", push)
					image.POST("/retag", reTag)
					image.POST("/pull", pull)

					conf := image.Group("/conf")
					{
						conf.POST("/update", updateImageConf)
						conf.GET("/query", queryImageConf)
					}
				}
			}

			dashboard := v1.Group("/dashboard")
			{
				machine := dashboard.Group("/machine")
				{
					machine.GET("/lifecycle", queryMachineLifeCycle)
					machine.GET("/count", queryMachineCount)
					machine.GET("/time", queryMachineTime)
					machine.GET("/usage", queryMachineUsage)
				}

				deviceFarm := dashboard.Group("/devicefarm")
				{
					deviceFarm.GET("/usagecount", queryMobileMachineUsageCount)
					deviceFarm.GET("/modecount", queryMobileMachineModeCount)
					deviceFarm.GET("/usagetime", queryMobileMachineUsageTime)
				}
			}
		}

		v2 := base.Group("/v2")
		{
			machine := v2.Group("/machine")
			{
				machine.POST("/gen", genMachinesV2)
			}
		}
	}

	callback := e.Group("/ep/admin/merlin/callback")
	{
		v1 := callback.Group("/v1")
		{

			v1.POST("/bilihub/snapshot", callbackSnapshot)
			v1.POST("/mobiledevice/error", callbackMobileDeviceError)
		}

	}
}

func ping(c *bm.Context) {
	if err := svc.Ping(c); err != nil {
		log.Error("ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func getVersion(c *bm.Context) {
	v := new(struct {
		Version string `json:"version"`
	})
	v.Version = "v.1.5.9.3"
	c.JSON(v, nil)

}

func confVersion(c *bm.Context) {
	v := new(struct {
		Version string `json:"version"`
	})
	v.Version = svc.ConfVersion(c)
	c.JSON(v, nil)

}
