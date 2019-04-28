package http

import (
	"strings"

	item "go-common/app/service/openplatform/ticket-item/api/grpc/v1"
	"go-common/app/service/openplatform/ticket-item/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"

	"github.com/pkg/errors"
)

// @params ParamID
// @router get /openplatform/internal/ticket/item/info
// @response InfoReply
func info(c *bm.Context) {
	arg := new(model.ParamID)
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	c.JSON(itemSvc.Info(c, &item.InfoRequest{ID: arg.ID}))
}

// @params ParamCards
// @router get /openplatform/internal/ticket/item/cards
// @response CardsReply
func cards(c *bm.Context) {
	arg := new(model.ParamCards)
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	ids := model.UniqueInt64(model.String2Int64(strings.Split(arg.IDs, ",")))
	if len(ids) == 0 {
		err := ecode.RequestErr
		errors.Wrap(err, "参数验证失败")
		log.Error("ItemID invalid %v", arg.IDs)
		return
	}

	c.JSON(itemSvc.Cards(c, &item.CardsRequest{IDs: ids}))
}

// @params ParamBill
// @router get /openplatform/internal/ticket/item/billinfo
// @response BillReply
func billInfo(c *bm.Context) {
	arg := new(model.ParamBill)
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	ids := model.UniqueInt64(model.String2Int64(strings.Split(arg.IDs, ",")))
	if len(ids) == 0 {
		err := ecode.RequestErr
		errors.Wrap(err, "参数验证失败")
		log.Error("ItemID empty %v", arg.IDs)
		return
	}

	sids := model.UniqueInt64(model.String2Int64(strings.Split(arg.Sids, ",")))
	if len(ids) == 0 {
		log.Info("ScreenID empty %v", arg.Sids)
		return
	}

	tids := model.UniqueInt64(model.String2Int64(strings.Split(arg.Tids, ",")))
	if len(ids) == 0 {
		log.Info("TicketID empty %v", arg.Tids)
		return
	}

	c.JSON(itemSvc.BillInfo(c, &item.BillRequest{IDs: ids, ScIDs: sids, TkIDs: tids}))
}

// @params WishRequest
// @router post /openplatform/internal/ticket/item/wishstatus
// @response WishReply
func wish(c *bm.Context) {
	arg := new(item.WishRequest)
	if err := c.BindWith(arg, binding.JSON); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}

	c.JSON(itemSvc.Wish(c, arg))
}

// @params FavRequest
// @router post /openplatform/internal/ticket/item/favstatus
// @response FavReply
func fav(c *bm.Context) {
	arg := new(item.FavRequest)
	if err := c.BindWith(arg, binding.JSON); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}

	c.JSON(itemSvc.Fav(c, arg))
}
