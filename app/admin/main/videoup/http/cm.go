package http

import (
	"encoding/json"
	"io/ioutil"

	"go-common/app/admin/main/videoup/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	xtime "go-common/library/time"
)

// upArchiveArr update archive attribute.
func upCMArr(c *bm.Context) {
	req := c.Request
	// read
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error("ioutil.ReadAll() error(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	req.Body.Close()
	// params
	var ap struct {
		Aid      int64 `json:"aid"`
		AdminID  int64 `json:"admin_id"`
		RankAttr struct {
			Main      int32 `json:"main"`
			RecentArc int32 `json:"recent_arc"`
			AllArc    int32 `json:"all_arc"`
		} `json:"rank_attr"`
		DynamicAttr struct {
			Main int32 `json:"main"`
		} `json:"dynamic_attr"`
		RecommendAttr struct {
			Main int32 `json:"main"`
		} `json:"recommend_attr"`
	}
	if err = json.Unmarshal(bs, &ap); err != nil {
		log.Error("http upArchiveArr() json.Unmarshal(%s) error(%v)", string(bs), err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if ap.Aid == 0 {
		log.Error("aid==0")
		c.JSON(nil, ecode.RequestErr)
		return
	}
	attrs := make(map[uint]int32, 6)
	attrs[archive.AttrBitNoRank] = ap.RankAttr.Main
	attrs[archive.AttrBitNoDynamic] = ap.DynamicAttr.Main
	attrs[archive.AttrBitNoRecommend] = ap.RecommendAttr.Main
	// forbid
	forbidAttrs := make(map[string]map[uint]int32, 3)
	forbidAttrs[archive.ForbidRank] = map[uint]int32{
		archive.ForbidRankMain:      ap.RankAttr.Main,
		archive.ForbidRankRecentArc: ap.RankAttr.RecentArc,
		archive.ForbidRankAllArc:    ap.RankAttr.AllArc,
	}
	forbidAttrs[archive.ForbidDynamic] = map[uint]int32{
		archive.ForbidDynamicMain: ap.DynamicAttr.Main,
	}
	forbidAttrs[archive.ForbidRecommend] = map[uint]int32{
		archive.ForbidRecommendMain: ap.RecommendAttr.Main,
	}
	// update attrs and forbid
	c.JSON(nil, vdaSvc.UpArchiveAttr(c, ap.Aid, ap.AdminID, attrs, forbidAttrs, ""))
}

// upCMArrDelay up cm archive delaytime
func upCMArcDelay(c *bm.Context) {
	req := c.Request
	// read
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error("ioutil.ReadAll() error(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	req.Body.Close()
	// params
	var ap struct {
		Aid   int64      `json:"aid"`
		Dtime xtime.Time `json:"dtime"`
	}
	if err = json.Unmarshal(bs, &ap); err != nil {
		log.Error("http struct aid json.Unmarshal(%s) error(%v)", string(bs), err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if ap.Aid == 0 {
		log.Error("aid==0")
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = vdaSvc.UpArcDtime(c, ap.Aid, ap.Dtime); err != nil {
		log.Error("vdaSvc.UpArcDtime() error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}
