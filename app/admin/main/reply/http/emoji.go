package http

import (
	bm "go-common/library/net/http/blademaster"
)

func listEmoji(c *bm.Context) {
	v := new(struct {
		Pid int64 `form:"pid" default:"0"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	emojis, err := rpSvc.EmojiList(c, v.Pid)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(emojis, nil)
}

func createEmoji(c *bm.Context) {
	v := new(struct {
		Pid    int64  `form:"pid" validate:"min=0"`
		Name   string `form:"name" validate:"required"`
		URL    string `form:"url" validate:"required"`
		Sort   int32  `form:"sort" validate:"min=0"`
		Remark string `form:"remark" validate:"required"`
		State  int32  `form:"state" validate:"min=0"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	if err := rpSvc.EmojiByName(c, v.Name); err != nil {
		c.JSON(nil, err)
		return
	}
	id, err := rpSvc.CreateEmoji(c, v.Pid, v.Name, v.URL, v.Sort, v.State, v.Remark)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(id, nil)
}

func upEmojiSort(c *bm.Context) {
	v := new(struct {
		IDs string `form:"ids" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	err := rpSvc.UpEmojiSort(c, v.IDs)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func upEmojiState(c *bm.Context) {
	v := new(struct {
		ID    int64 `form:"id" validate:"min=0"`
		State int32 `form:"state" validate:"min=0"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	id, err := rpSvc.UpEmojiState(c, v.State, v.ID)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(id, nil)
}

func upEmoji(c *bm.Context) {
	v := new(struct {
		ID     int64  `form:"id" validate:"min=0"`
		Name   string `form:"name" validate:"required"`
		Remark string `form:"remark" validate:"required"`
		URL    string `form:"url" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	id, err := rpSvc.UpEmoji(c, v.Name, v.Remark, v.URL, v.ID)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(id, nil)
}
