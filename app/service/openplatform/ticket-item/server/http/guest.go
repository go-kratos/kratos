package http

import (
	"go-common/app/service/openplatform/ticket-item/model"
	bm "go-common/library/net/http/blademaster"
)

/** test http
func guestInfo(c *bm.Context) {
	arg := new(model.GuestParam)
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}

	c.JSON(itemSvc.GuestInfo(c, &model.GuestInfoRequest{ID: arg.ID, Name: arg.Name, GuestImg: arg.GuestImg, Description: arg.Description, GuestID: arg.GuestID}))
}

func guestStatus(c *bm.Context) {
	arg := new(model.GuestStatusParam)
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}

	c.JSON(itemSvc.GuestStatus(c, arg.ID, arg.Status))
}**/

// @params ParamID
// @router get /openplatform/internal/ticket/item/getguests
// @response Guest
/**func getGuests(c *bm.Context) {
	arg := new(model.ParamID)
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	c.JSON(itemSvc.GetGuests(c, &arg.ID))
}**/

// @params GuestSearchParam
// @router get /openplatform/internal/ticket/item/guest/search
// @response GuestSearchList
func guestSearch(c *bm.Context) {
	arg := new(model.GuestSearchParam)
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(itemSvc.GuestSearch(c, arg))
}
