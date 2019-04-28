package http

import (
	"go-common/app/admin/ep/melloi/model"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)

// get administrator for order by current username
func queryOrderAdmin(c *bm.Context) {
	userName, _ := c.Request.Cookie("username")
	c.JSON(srv.QueryOrderAdmin(userName.Value))
}

// add administrator for order
func addOrderAdmin(c *bm.Context) {
	admin := model.OrderAdmin{}
	if err := c.BindWith(&admin, binding.Form); err != nil {
		return
	}
	c.JSON(nil, srv.AddOrderAdmin(&admin))
}
