package http

import (
	"go-common/app/service/main/passport/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func historyPwdCheck(c *bm.Context) {
	param := new(model.HistoryPwdCheckParam)
	c.Bind(param)
	if param.Mid <= 0 || param.Pwd == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(passportSvc.HistoryPwdCheck(c, param))
}
