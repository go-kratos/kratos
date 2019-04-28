package http

import (
	"strconv"

	"go-common/app/interface/main/web/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func ranking(c *bm.Context) {
	var (
		rid                    int64
		rankType, day, arcType int
		err                    error
	)
	params := c.Request.Form
	ridStr := params.Get("rid")
	rankTypeStr := params.Get("type")
	dayStr := params.Get("day")
	arcTypeStr := params.Get("arc_type")
	if rid, err = strconv.ParseInt(ridStr, 10, 64); err != nil || rid < 0 {
		rid = 0
	}
	if rankType, err = strconv.Atoi(rankTypeStr); err != nil {
		rankType = 1
	}
	if day, err = strconv.Atoi(dayStr); err != nil {
		day = 3
	}
	if arcType, err = strconv.Atoi(arcTypeStr); err != nil {
		arcType = 0
	}
	if err = checkType(rid, rankType, day, arcType); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(webSvc.Ranking(c, int16(rid), rankType, model.DayType[day], model.ArcType[arcType]))
}

func checkType(rid int64, rankType, day, arcType int) (err error) {
	if _, ok := model.RankType[rankType]; !ok {
		err = ecode.RequestErr
		return
	}
	if _, ok := model.DayType[day]; !ok {
		err = ecode.RequestErr
		return
	}
	if _, ok := model.ArcType[arcType]; !ok {
		err = ecode.RequestErr
	}
	// bangumi and rookie not have recent contribution
	if (rid == 33 || rankType == 3) && arcType == 1 {
		err = ecode.RequestErr
	}
	return
}

func rankingIndex(c *bm.Context) {
	var (
		day int
		err error
	)
	params := c.Request.Form
	dayStr := params.Get("day")
	if day, err = strconv.Atoi(dayStr); err != nil {
		day = 3
	}
	if err = checkIndexDay(day); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(webSvc.RankingIndex(c, day))
}

func rankingRegion(c *bm.Context) {
	var (
		day, original int
		rid           int64
		err           error
	)
	params := c.Request.Form
	ridStr := params.Get("rid")
	dayStr := params.Get("day")
	originalStr := params.Get("original")
	if rid, err = strconv.ParseInt(ridStr, 10, 64); err != nil || rid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if day, err = strconv.Atoi(dayStr); err != nil {
		day = 3
	}
	if _, ok := model.DayType[day]; !ok {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if original, err = strconv.Atoi(originalStr); err != nil {
		original = 0
	} else if original != 1 && original != 0 {
		original = 0
	}
	c.JSON(webSvc.RankingRegion(c, int16(rid), day, original))
}

func rankingRecommend(c *bm.Context) {
	var (
		rid int64
		err error
	)
	params := c.Request.Form
	ridStr := params.Get("rid")
	if rid, err = strconv.ParseInt(ridStr, 10, 64); err != nil || rid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(webSvc.RankingRecommend(c, int16(rid)))
}

func rankingTag(c *bm.Context) {
	var (
		rid, tagID int64
		err        error
	)
	params := c.Request.Form
	ridStr := params.Get("rid")
	tagIDStr := params.Get("tag_id")
	if rid, err = strconv.ParseInt(ridStr, 10, 64); err != nil || rid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if tagID, err = strconv.ParseInt(tagIDStr, 10, 64); err != nil || tagID <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(webSvc.RankingTag(c, int16(rid), tagID))
}

func regionCustom(c *bm.Context) {
	c.JSON(webSvc.RegionCustom(c))
}

func checkIndexDay(day int) (err error) {
	err = ecode.RequestErr
	for _, dayItem := range model.IndexDayType {
		if day == dayItem {
			err = nil
			return
		}
	}
	return
}
