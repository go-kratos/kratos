package block

import (
	model "go-common/app/admin/main/member/model/block"
	service "go-common/app/admin/main/member/service/block"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"

	"github.com/pkg/errors"
)

var (
	svc *service.Service
)

// Setup is.
func Setup(blockSvc *service.Service, e *bm.Engine, authSvc *permit.Permit) {
	svc = blockSvc
	cb := e.Group("/x/admin/block", authSvc.Permit("BLOCK_SEARCH"))
	{
		cb.POST("/search", blockSearch)
		cb.GET("/history", history)
	}
	cb = e.Group("/x/admin/block", authSvc.Permit("BLOCK_BLOCK"))
	{
		cb.POST("", batchBlock)
	}
	cb = e.Group("/x/admin/block", authSvc.Permit("BLOCK_REMOVE"))
	{
		cb.POST("/remove", batchRemove)
	}
}

func bind(c *bm.Context, v model.ParamValidator) (err error) {
	if err = c.Bind(v); err != nil {
		err = errors.WithStack(err)
		return
	}
	if !v.Validate() {
		err = ecode.RequestErr
		c.JSON(nil, ecode.RequestErr)
		return
	}
	return
}
