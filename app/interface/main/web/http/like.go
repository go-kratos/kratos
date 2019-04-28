package http

import (
	"strconv"

	"go-common/app/interface/main/web/model"
	bm "go-common/library/net/http/blademaster"
)

func like(c *bm.Context) {
	var (
		mid int64
		err error
	)
	v := new(struct {
		Aid  int64 `form:"aid" validate:"min=1,required"`
		Like int8  `form:"like" validate:"min=1,max=4,required"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	midStr, _ := c.Get("mid")
	mid = midStr.(int64)
	upperID, err := webSvc.Like(c, v.Aid, mid, v.Like)
	c.JSON(nil, err)
	if err == nil && webSvc.CheatInfoc != nil {
		webSvc.CheatInfoc.InfoAntiCheat2(c, strconv.FormatInt(upperID, 10), strconv.FormatInt(v.Aid, 10), strconv.FormatInt(mid, 10), strconv.FormatInt(v.Aid, 10), "av", model.LikeType[v.Like], "")
	}
}

func likeTriple(c *bm.Context) {
	var (
		mid       int64
		actTriple = "triplelike"
		err       error
	)
	v := new(struct {
		Aid int64 `form:"aid" validate:"min=1,required"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	midStr, _ := c.Get("mid")
	mid = midStr.(int64)
	triple, err := webSvc.LikeTriple(c, v.Aid, mid)
	c.JSON(triple, err)
	if err != nil {
		return
	}
	if err == nil && webSvc.CheatInfoc != nil && triple.Anticheat {
		webSvc.CheatInfoc.InfoAntiCheat2(c, strconv.FormatInt(triple.UpID, 10), strconv.FormatInt(v.Aid, 10), strconv.FormatInt(mid, 10), strconv.FormatInt(v.Aid, 10), "av", actTriple, "")
	}
}

func hasLike(c *bm.Context) {
	var (
		mid int64
	)
	v := new(struct {
		Aid int64 `form:"aid" validate:"min=1,required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	midStr, _ := c.Get("mid")
	mid = midStr.(int64)
	c.JSON(webSvc.HasLike(c, v.Aid, mid))
}
