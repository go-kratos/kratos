package http

import bm "go-common/library/net/http/blademaster"

func relation(c *bm.Context) {
	v := new(struct {
		Mid int64 `form:"mid" validate:"min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	data, err := spcSvc.Relation(c, v.Mid)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(struct {
		Follower int64 `json:"follower"`
	}{Follower: data}, nil)
}
