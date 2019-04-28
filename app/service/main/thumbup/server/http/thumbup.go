package http

import (
	"go-common/app/service/main/thumbup/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
	"go-common/library/net/metadata"
)

func like(c *bm.Context) {
	var (
		err error
	)
	v := new(struct {
		Business  string `form:"business" validate:"required"`
		OriginID  int64  `form:"origin_id" validate:"min=0"`
		MessageID int64  `form:"message_id" validate:"min=1,required"`
		Type      int8   `form:"type" validate:"required"`
		Mid       int64  `form:"mid" validate:"min=1,required"`
		UpMid     int64  `form:"up_mid" validate:"omitempty,min=1"`
	})
	if err = c.Bind(v); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(likeSrv.Like(c, v.Business, v.Mid, v.OriginID, v.MessageID, v.Type, v.UpMid))
}

func hasLike(c *bm.Context) {
	v := new(struct {
		Business   string  `form:"business" validate:"required"`
		OriginID   int64   `form:"origin_id" validate:"min=0"`
		MessageIDs []int64 `form:"message_ids,split" validate:"required"`
		Mid        int64   `form:"mid" validate:"min=1,required"`
	})
	if err := c.Bind(v); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	res, _, err := likeSrv.HasLike(c, v.Business, v.Mid, v.MessageIDs)
	c.JSON(res, err)
}

func stats(c *bm.Context) {
	v := new(struct {
		Business   string  `form:"business" validate:"required"`
		OriginID   int64   `form:"origin_id" validate:"min=0"`
		MessageIDs []int64 `form:"message_ids,split" validate:"required"`
		Mid        int64   `form:"mid"`
	})
	if err := c.Bind(v); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if v.Mid > 0 {
		c.JSON(likeSrv.StatsWithLike(c, v.Business, v.Mid, v.OriginID, v.MessageIDs))
		return
	}
	c.JSON(likeSrv.Stats(c, v.Business, v.OriginID, v.MessageIDs))
}

func userLikes(c *bm.Context) {
	var (
		err  error
		data []*model.ItemLikeRecord
	)
	v := new(struct {
		Business string `form:"business" validate:"required"`
		Type     string `form:"type" validate:"required"`
		Mid      int64  `form:"mid" validate:"min=1,required"`
		Pn       int    `form:"pn" default:"1" validate:"omitempty,min=1"`
		Ps       int    `form:"ps" default:"20" validate:"omitempty,min=1"`
	})
	if err = c.Bind(v); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if v.Type == "like" {
		data, err = likeSrv.UserLikes(c, v.Business, v.Mid, v.Pn, v.Ps)
	} else {
		data, err = likeSrv.UserDislikes(c, v.Business, v.Mid, v.Pn, v.Ps)
	}
	if data == nil {
		data = make([]*model.ItemLikeRecord, 0)
	}
	c.JSON(data, err)
}

func itemLikes(c *bm.Context) {
	var (
		err  error
		data []*model.UserLikeRecord
	)
	v := new(struct {
		Business  string `form:"business" validate:"required"`
		OriginID  int64  `form:"origin_id" validate:"min=0"`
		MessageID int64  `form:"message_id" validate:"min=1,required"`
		Mid       int64  `form:"mid" validate:"omitempty,min=1"`
		Type      string `form:"type" validate:"required"`
		Pn        int    `form:"pn" default:"1" validate:"omitempty,min=1"`
		Ps        int    `form:"ps" default:"20" validate:"omitempty,min=1"`
	})
	if err = c.Bind(v); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if v.Type == "like" {
		data, err = likeSrv.ItemLikes(c, v.Business, v.OriginID, v.MessageID, v.Pn, v.Ps, v.Mid)
	} else {
		data, err = likeSrv.ItemDislikes(c, v.Business, v.OriginID, v.MessageID, v.Pn, v.Ps, v.Mid)
	}
	if data == nil {
		data = make([]*model.UserLikeRecord, 0)
	}
	c.JSON(data, err)
}

func multiStats(c *bm.Context) {
	v := new(model.MultiBusiness)
	if err := c.BindWith(v, binding.JSON); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(likeSrv.MultiStatsWithLike(c, v))
}

func updateCount(c *bm.Context) {
	v := new(struct {
		Business      string `form:"business" validate:"required"`
		OriginID      int64  `form:"origin_id" validate:"min=0"`
		MessageID     int64  `form:"message_id" validate:"min=1,required"`
		LikeChange    int64  `form:"like_change"`
		DislikeChange int64  `form:"dislike_change"`
		Operator      string `form:"operator" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	ip := metadata.String(c, metadata.RemoteIP)
	c.JSON(nil, likeSrv.UpdateCount(c, v.Business, v.OriginID, v.MessageID, v.LikeChange, v.DislikeChange, ip, v.Operator))
}

func rawStats(c *bm.Context) {
	v := new(struct {
		Business  string `form:"business" validate:"required"`
		OriginID  int64  `form:"origin_id" validate:"min=0"`
		MessageID int64  `form:"message_id" validate:"min=1,required"`
	})
	if err := c.Bind(v); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(likeSrv.RawStats(c, v.Business, v.OriginID, v.MessageID))
}

func updateUpMids(c *bm.Context) {
	v := new(struct {
		Business string             `json:"business" validate:"required"`
		Data     []*model.UpMidsReq `json:"data" validate:"required,gte=1,lte=100"`
	})
	if err := c.BindWith(v, binding.JSON); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, likeSrv.UpdateUpMids(c, v.Business, v.Data))
}

func itemHasLike(c *bm.Context) {
	v := new(struct {
		Business  string  `json:"business" validate:"required"`
		OriginID  int64   `json:"origin_id" validate:"min=0"`
		MessageID int64   `json:"message_id" validate:"min=1"`
		Mids      []int64 `json:"mids,split" validate:"required"`
	})
	if err := c.BindWith(v, binding.JSON); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(likeSrv.ItemHasLike(c, v.Business, v.OriginID, v.MessageID, v.Mids))
}
