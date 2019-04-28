package http

import (
	"net/http"

	"go-common/app/service/openplatform/ticket-item/conf"
	"go-common/app/service/openplatform/ticket-item/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var (
	itemSvc *service.ItemService
)

// Init init item service.
func Init(c *conf.Config, s *service.ItemService) {
	itemSvc = s
	//engineIn := bm.DefaultServer(c.BM.Inner)
	engineIn := bm.DefaultServer(nil)
	innerRouter(engineIn)
	// init inner server
	if err := engineIn.Start(); err != nil {
		log.Error("engineIn.Start error(%v)", err)
		panic(err)
	}
	/** unused
	engineLocal := bm.DefaultServer(c.BM.Local)
	localRouter(engineLocal)
	// init local router
	if err := engineLocal.Start(); err != nil {
		log.Error("engineLocal.Start error(%v)", err)
		panic(err)
	}**/
}

// innerRouter init inner router.
func innerRouter(e *bm.Engine) {
	e.Ping(ping)
	// item
	group := e.Group("/openplatform/internal/ticket/item")
	{
		group.GET("/info", info)
		group.GET("/billinfo", billInfo)
		group.GET("/cards", cards)
		// group.GET("/guestinfo", guestInfo)
		// group.GET("/gueststatus", guestStatus)
		//group.GET("/getguests", getGuests)
		//group.GET("/getbulletins", getBulletins)
		group.GET("/guest/search", guestSearch)
		group.GET("/venue/search", venueSearch)     // 场馆搜索
		group.GET("/version/search", versionSearch) // 项目版本搜索
		group.POST("/wishstatus", wish)             // 想去项目
		group.POST("/favstatus", fav)               // 收藏/取消收藏项目
		group.POST("/venueinfo", venueInfo)
		group.POST("/placeinfo", placeInfo)
		group.POST("/areainfo", areaInfo)
		group.POST("/seatinfo", seatInfo)
		group.POST("/seatStock", seatStock)
		group.POST("/removeSeatOrders", removeSeatOrders)
	}
}

// localRouter init local router.
/**func localRouter(e *bm.Engine) {
}**/

// ping check server ok.
func ping(c *bm.Context) {
	if err := itemSvc.Ping(c); err != nil {
		log.Error("ticket http service ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
