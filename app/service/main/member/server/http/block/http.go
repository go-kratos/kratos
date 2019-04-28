package block

import (
	model "go-common/app/service/main/member/model/block"
	service "go-common/app/service/main/member/service/block"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	v "go-common/library/net/http/blademaster/middleware/verify"

	"github.com/pkg/errors"
)

var (
	svc *service.Service
)

// Setup is
func Setup(blockSvc *service.Service, engine *bm.Engine, v *v.Verify) {
	svc = blockSvc

	blkGroup := engine.Group("/x/internal/block", v.Verify)
	blkGroup.POST("/block", block)
	blkGroup.POST("/remove", remove)
	blkGroup.GET("/info", info)

	blkGroupBatch := engine.Group("/x/internal/block/batch", v.Verify)
	blkGroupBatch.POST("/block", batchBlock)
	blkGroupBatch.POST("/remove", batchRemove)
	blkGroupBatch.POST("/info", batchInfo)
	blkGroupBatch.GET("/detail", batchDetail)
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
