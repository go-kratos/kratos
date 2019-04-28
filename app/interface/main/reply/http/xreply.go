package http

import (
	xmodel "go-common/app/interface/main/reply/model/xreply"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

func xreply(c *bm.Context) {
	v := new(xmodel.ReplyReq)
	if err := c.Bind(v); err != nil {
		return
	}
	v.Mid = metadata.Int64(c, metadata.Mid)
	v.IP = metadata.String(c, metadata.RemoteIP)
	if !v.Cursor.Legal() {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(rpSvr.Xreply(c, v))
}

func subFolder(c *bm.Context) {
	v := new(xmodel.SubFolderReq)
	if err := c.Bind(v); err != nil {
		return
	}
	v.Mid = metadata.Int64(c, metadata.Mid)
	v.IP = metadata.String(c, metadata.RemoteIP)
	if !v.Cursor.Legal() || v.Cursor.Backward() {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(rpSvr.SubFoldedReply(c, v))
}

func rootFolder(c *bm.Context) {
	v := new(xmodel.RootFolderReq)
	if err := c.Bind(v); err != nil {
		return
	}
	v.Mid = metadata.Int64(c, metadata.Mid)
	v.IP = metadata.String(c, metadata.RemoteIP)
	if !v.Cursor.Legal() || v.Cursor.Backward() {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(rpSvr.RootFoldedReply(c, v))
}
