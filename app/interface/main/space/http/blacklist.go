package http

import (
	"go-common/library/conf/env"
	bm "go-common/library/net/http/blademaster"
)

var (
	_EmptyBlacklistValue = make(map[int64]struct{})
)

func blacklist(c *bm.Context) {
	black := spcSvc.BlacklistValue
	if env.DeployEnv == env.DeployEnvProd {
		if len(black) == 0 {
			black = _EmptyBlacklistValue
		}
	} else {
		black = _EmptyBlacklistValue
	}
	c.JSON(black, nil)
}
