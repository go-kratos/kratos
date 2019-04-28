package http

import (
	"encoding/json"

	item "go-common/app/service/openplatform/ticket-item/api/grpc/v1"
	"go-common/app/service/openplatform/ticket-item/model"
	bm "go-common/library/net/http/blademaster"

	"github.com/pkg/errors"
)

// @params SeatInfoParam
// @router post /openplatform/internal/ticket/item/seatInfo
// @response SeatInfoReply
func seatInfo(c *bm.Context) {
	arg := new(model.SeatInfoParam)
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	var seats []*item.AreaSeatInfo
	json.Unmarshal([]byte(arg.Seats), &seats)
	c.JSON(itemSvc.SeatInfo(c, &item.SeatInfoRequest{
		Area:      arg.Area,
		SeatsNum:  arg.SeatsNum,
		Seats:     seats,
		Width:     arg.Width,
		Height:    arg.Height,
		RowList:   arg.RowList,
		SeatStart: arg.SeatStart,
	}))
}

// @params SeatStockParam
// @router post /openplatform/internal/ticket/item/seatStock
// @response SeatStockReply
func seatStock(c *bm.Context) {
	arg := new(model.SeatStockParam)
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	var seatinfo []*item.SeatPrice
	json.Unmarshal([]byte(arg.SeatInfo), &seatinfo)
	c.JSON(itemSvc.SeatStock(c, &item.SeatStockRequest{
		Screen:   arg.Screen,
		Area:     arg.Area,
		SeatInfo: seatinfo,
	}))
}

// @params RemoveSeatOrdersParam
// @router post /openplatform/internal/ticket/item/RemoveSeatOrders
// @response RemoveSeatOrdersReply
func removeSeatOrders(c *bm.Context) {
	arg := new(model.RemoveSeatOrdersParam)
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	c.JSON(itemSvc.RemoveSeatOrders(c, &item.RemoveSeatOrdersRequest{
		Price: arg.Price,
	}))
}
