package http

import (
	"go-common/app/service/main/msm/model"
	bm "go-common/library/net/http/blademaster"
)

func scope(c *bm.Context) {
	var (
		err      error
		scopeMap map[int64]*model.Scope
		param    = new(struct {
			AppTreeID int64 `form:"app_tree_id" validate:"gt=0"`
		})
	)
	if err = c.Bind(param); err != nil {
		return
	}
	if scopeMap, err = svr.ServiceScopes(c, param.AppTreeID); err != nil {
		c.JSON(nil, err)
		return
	}
	data := make(map[string]interface{}, 2)
	data["service_tree_id"] = param.AppTreeID
	data["scopes"] = scopeMap
	c.JSON(data, nil)
}
