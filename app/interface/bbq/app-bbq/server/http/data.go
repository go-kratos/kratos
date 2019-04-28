package http

import (
	"go-common/app/interface/bbq/app-bbq/model"
	bm "go-common/library/net/http/blademaster"
)

func videoPlay(c *bm.Context) {
	uiLog(c, model.ActionPlay, nil)
}
