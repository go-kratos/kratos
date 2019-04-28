package http

import (
	v12 "go-common/app/service/live/xcaptcha/api/grpc/v1"
	"go-common/library/net/http/blademaster"
)

// captchaVerify
func captchaVerify(ctx *blademaster.Context) {
	req := new(v12.XVerifyReq)
	if err := ctx.Bind(req); err != nil {
		return
	}
	resp, err := xCaptchaService.Verify(ctx, req)
	ctx.JSON(resp, err)
}
