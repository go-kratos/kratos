package http

import (
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func listEmojiPacks(c *bm.Context) {
	packs, err := rpSvc.EmojiPackageList(c)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(packs, nil)
}

func createEmojiPackage(c *bm.Context) {
	v := new(struct {
		Name   string `form:"name" validate:"required"`
		URL    string `form:"url" validate:"required"`
		Sort   int32  `form:"sort" validate:"min=0"`
		Remark string `form:"remark" default:""`
		State  int32  `form:"state" validate:"min=0"`
	})
	if err := c.Bind(v); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	id, err := rpSvc.CreateEmojiPackage(c, v.Name, v.URL, v.Sort, v.Remark, v.State)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(id, nil)
}

func upEmojiPackageSort(c *bm.Context) {
	v := new(struct {
		IDs string `form:"ids" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err := rpSvc.UpEmojiPackageSort(c, v.IDs)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func editEmojiPack(c *bm.Context) {
	v := new(struct {
		ID     int64  `form:"id" validate:"min=0"`
		Name   string `form:"name" default:""`
		URL    string `form:"url" default:""`
		Remark string `form:"remark" default:""`
		State  int32  `form:"state" validate:"min=0"`
	})
	if err := c.Bind(v); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	id, err := rpSvc.UpEmojiPackage(c, v.Name, v.URL, v.Remark, v.State, v.ID)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(id, nil)
}
