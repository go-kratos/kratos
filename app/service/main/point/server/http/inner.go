package http

import bm "go-common/library/net/http/blademaster"

// innerRouter init local router api path.
func innerRouter(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)
	ig := e.Group("/x/internal/point", verSvc.Verify)
	{
		ig.GET("/info", pointInfoInner)
		ig.POST("/consume", pointConsume)
		ig.POST("/add", pointAddByBp)
		ig.GET("/configs", configs)
		ig.GET("/config", config)
		ig.GET("/old/history", oldPointHistory)
		ig.POST("/addbyadmin", pointAdd)
	}

}
