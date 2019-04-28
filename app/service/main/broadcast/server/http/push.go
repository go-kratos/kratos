package http

import (
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

const maxMsgSize = 1024 * 1024 * 8

func pushKeys(c *bm.Context) {
	v := new(struct {
		Op          int32    `form:"operation" validate:"required"`
		Keys        []string `form:"keys,split" validate:"required"`
		Msg         string   `form:"message" validate:"required"`
		ContentType int32    `form:"content_type"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	if len(v.Msg) > maxMsgSize {
		c.JSON(ecode.FileTooLarge, nil)
	}
	c.JSON(nil, srv.PushKeys(c, v.Op, v.Keys, v.Msg, v.ContentType))
}

func pushMids(c *bm.Context) {
	v := new(struct {
		Op          int32   `form:"operation" validate:"required"`
		Mids        []int64 `form:"mids,split" validate:"required"`
		Msg         string  `form:"message" validate:"required"`
		ContentType int32   `form:"content_type"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	if len(v.Msg) > maxMsgSize {
		c.JSON(nil, ecode.FileTooLarge)
		return
	}
	c.JSON(nil, srv.PushMids(c, v.Op, v.Mids, v.Msg, v.ContentType))
}

func pushRoom(c *bm.Context) {
	v := new(struct {
		Op          int32  `form:"operation" validate:"required"`
		Room        string `form:"room" validate:"required"`
		Msg         string `form:"message" validate:"required"`
		ContentType int32  `form:"content_type"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	if len(v.Msg) > maxMsgSize {
		c.JSON(nil, ecode.FileTooLarge)
		return
	}
	c.JSON(nil, srv.PushRoom(c, v.Op, v.Room, v.Msg, v.ContentType))
}

func pushAll(c *bm.Context) {
	v := new(struct {
		Op          int32  `form:"operation" validate:"required"`
		Speed       int32  `form:"speed" validate:"min=0"`
		Msg         string `form:"message" validate:"required"`
		Platform    string `form:"platform"`
		ContentType int32  `form:"content_type"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	if len(v.Msg) > maxMsgSize {
		c.JSON(nil, ecode.FileTooLarge)
		return
	}
	c.JSON(nil, srv.PushAll(c, v.Op, v.Speed, v.Msg, v.Platform, v.ContentType))
}
