package http

import (
	"go-common/app/service/openplatform/abtest/model"
	"go-common/app/service/openplatform/abtest/model/validator"
	bm "go-common/library/net/http/blademaster"
)

//add group
func addGroup(c *bm.Context) {
	params := new(validator.AddGroupParams)
	if err := c.Bind(params); err != nil {
		return
	}
	g := model.Group{
		Name: params.Name,
		Desc: params.Desc,
	}
	c.JSON(abSvr.AddGroup(c, g))
}

//list group
func listGroup(c *bm.Context) {
	c.JSON(abSvr.ListGroup(c))
}

//update group
func updateGroup(c *bm.Context) {
	params := new(validator.UpdateGroupParams)
	if err := c.Bind(params); err != nil {
		return
	}
	g := model.Group{
		ID:   params.ID,
		Name: params.Name,
		Desc: params.Desc,
	}
	c.JSON(abSvr.UpdateGroup(c, g))
}

//delete group
func deleteGroup(c *bm.Context) {
	params := new(validator.DeleteGroupParams)
	if err := c.Bind(params); err != nil {
		return
	}
	c.JSON(abSvr.DeleteGroup(c, params.ID))
}
