package http

import (
	"go-common/app/admin/main/passport/model"
	bm "go-common/library/net/http/blademaster"
)

// @params UserBindLogReq
// @router get /x/admin/passport/userBindLog
// @response UserBindLogRes
func userBindLog(c *bm.Context) {
	userActLogReq := new(model.UserBindLogReq)
	c.Bind(userActLogReq)

	c.JSON(srv.UserBindLog(c, userActLogReq))
}

// @params DecryptBindLogParam
// @router get /x/admin/passport/user_bind_log/decrypt
// @response map[]string
func decryptBindLog(c *bm.Context) {
	param := new(model.DecryptBindLogParam)
	c.Bind(param)
	c.JSON(srv.DecryptBindLog(c, param))
}
