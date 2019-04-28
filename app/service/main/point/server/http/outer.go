package http

import (
	bm "go-common/library/net/http/blademaster"
)

// outerRouter init outer router api path.
func outerRouter(e *bm.Engine) {
	// api v1
	g := e.Group("/x/point", bm.CORS())
	{
		g.GET("/info", authSvc.User, pointInfo)
		g.GET("/history", authSvc.User, pointHistory)
	}

}
