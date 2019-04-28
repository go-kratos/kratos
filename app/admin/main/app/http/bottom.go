package http

import (
	"time"

	"go-common/app/admin/main/app/model/bottom"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

// bottoms select all
func bottoms(c *bm.Context) {
	c.JSON(bottomSvc.Bottoms(c))
}

// bottomByID select by id
func bottomByID(c *bm.Context) {
	v := &bottom.Param{}
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(bottomSvc.BottomByID(c, v.ID))
}

// bottomInsert insert or update
func bottomInsert(c *bm.Context) {
	var (
		err error
		v   = &bottom.Param{}
	)
	if err = c.Bind(v); err != nil {
		return
	}
	if v.ID > 0 {
		err = bottomSvc.Update(c, v, time.Now())
	} else {
		err = bottomSvc.Insert(c, v, time.Now())
	}
	c.JSON(nil, err)
}

// publish update state
func publishBottom(c *bm.Context) {
	v := &bottom.Param{}
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(nil, bottomSvc.Publish(c, v.IDs, time.Now()))
}

// delBottom delete
func delBottom(c *bm.Context) {
	v := &bottom.Param{}
	if err := c.Bind(v); err != nil {
		return
	}
	if v.ID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, bottomSvc.Delete(c, v.ID))
}
