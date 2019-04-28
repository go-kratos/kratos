package http

import (
	"go-common/app/service/openplatform/ticket-item/model"
	bm "go-common/library/net/http/blademaster"
)

// @params VersionSearchParam
// @router get /openplatform/internal/ticket/item/version/search
// @response VersionSearchList
func versionSearch(c *bm.Context) {
	req := &model.VersionSearchParam{}
	if err := c.Bind(req); err != nil {
		return
	}
	c.JSON(itemSvc.VersionSearch(c, req))
}
