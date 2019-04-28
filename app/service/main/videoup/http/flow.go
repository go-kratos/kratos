package http

import (
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func addByMid(c *bm.Context) {
	v := new(struct {
		Business int8  `form:"business"`
		Mid      int64 `form:"mid" validate:"required"`
		Oid      int64 `form:"oid" validate:"required"`
		State    int8  `form:"state"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	log.Info("flowDesign data(%v)", v)
	c.JSON(nil, vdpSvc.AddByMid(c, v.Business, v.Mid, v.Oid, v.State))
}

func addByOid(c *bm.Context) {
	v := new(struct {
		Business    int8  `form:"business" validate:"required"`
		OID         int64 `form:"oid" validate:"required"`
		UID         int64 `form:"uid" default:"399"`
		NoTimeline  int32 `form:"no_timeline" default:"-1"`
		NoOtt       int32 `form:"no_ott" default:"-1"`
		NoRank      int32 `form:"no_rank" default:"-1"`
		NoRecommend int32 `form:"no_recommend"  default:"-1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	log.Info("flowDesign data(%v)", v)
	c.JSON(nil, vdpSvc.AddByOid(c, v.Business, v.OID, v.UID, map[string]int32{
		"notimeline":  v.NoTimeline,
		"noott":       v.NoOtt,
		"norank":      v.NoRank,
		"norecommend": v.NoRecommend,
	}))
}

func queryByOid(c *bm.Context) {
	v := new(struct {
		Business int8  `form:"business" validate:"required"`
		Oid      int64 `form:"oid" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	log.Info("flowDesign data(%v)", v)
	c.JSON(vdpSvc.HitFlowGroups(c, v.Oid, []int8{v.Business}))
}

func listJudgeFlows(c *bm.Context) {
	//按照禁止项过滤  对外 一个禁止项属于一个groupId
	v := new(struct {
		Business int64   `form:"business" validate:"required"`
		Gid      int64   `form:"gid" validate:"required"`
		Oids     []int64 `form:"oids,split" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	if len(v.Oids) > 200 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	log.Info("flowDesign data(%v)", v)
	c.JSON(vdpSvc.TagrgetFlows(c, v.Business, v.Gid, v.Oids))
}

func listFlows(c *bm.Context) {
	//按照禁止项过滤  对外 一个禁止项属于一个groupId
	v := new(struct {
		Business int64 `form:"business" validate:"required"`
		Gid      int64 `form:"gid" validate:"required"`
		Pn       int64 `form:"page" default:"1"`
		Ps       int64 `form:"pagesize" default:"100" validate:"max=1000"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	log.Info("flowDesign data(%v)", v)
	c.JSON(vdpSvc.FlowPage(c, v.Business, v.Gid, v.Pn, v.Ps))
}
