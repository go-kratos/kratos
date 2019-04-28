package dao

import (
	"context"
	"go-common/library/net/metadata"
)

//ReplyReg 评论注册/冻结
func (d *Dao) ReplyReg(c context.Context, req map[string]interface{}) (err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	_, err = replyHTTPCommon(c, d.HTTPClient, d.c.URLs["reply_reg"], "POST", req, ip)
	return
}
