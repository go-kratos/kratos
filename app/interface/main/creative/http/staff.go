package http

import (
	bm "go-common/library/net/http/blademaster"
)

// staffConfig 获取联合投稿配置
func staffConfig(c *bm.Context) {
	res := staffSvc.Config()
	c.JSON(res, nil)
}
