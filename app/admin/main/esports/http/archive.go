package http

import (
	"encoding/csv"
	"io/ioutil"
	"strings"

	"go-common/app/admin/main/esports/conf"
	"go-common/app/admin/main/esports/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
	"strconv"
)

func arcList(c *bm.Context) {
	var (
		list []*model.ArcResult
		cnt  int
		err  error
	)
	v := new(model.ArcListParam)
	if err = c.Bind(v); err != nil {
		return
	}
	if v.Pn == 0 {
		v.Pn = 1
	}
	if v.Ps == 0 {
		v.Ps = 20
	}
	if list, cnt, err = esSvc.ArcList(c, v); err != nil {
		c.JSON(nil, err)
		return
	}
	data := make(map[string]interface{}, 2)
	page := map[string]int{
		"num":   v.Pn,
		"size":  v.Ps,
		"total": cnt,
	}
	data["page"] = page
	data["list"] = list
	c.JSON(data, nil)
}

func batchAddArc(c *bm.Context) {
	v := new(model.ArcAddParam)
	if err := c.Bind(v); err != nil {
		return
	}
	if len(v.Aids) > conf.Conf.Rule.MaxBatchArcLimit {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, esSvc.BatchAddArc(c, v))
}

func batchEditArc(c *bm.Context) {
	v := new(model.ArcAddParam)
	if err := c.Bind(v); err != nil {
		return
	}
	if len(v.Aids) > conf.Conf.Rule.MaxBatchArcLimit {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, esSvc.BatchEditArc(c, v))
}

func editArc(c *bm.Context) {
	v := new(model.ArcImportParam)
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(nil, esSvc.EditArc(c, v))
}

func arcImportCSV(c *bm.Context) {
	var (
		err  error
		data []byte
	)
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		log.Error("arcImportCSV upload err(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	defer file.Close()
	data, err = ioutil.ReadAll(file)
	if err != nil {
		log.Error("arcImportCSV ioutil.ReadAll err(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	r := csv.NewReader(strings.NewReader(string(data)))
	records, err := r.ReadAll()
	if err != nil {
		log.Error("r.ReadAll() err(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if l := len(records); l > conf.Conf.Rule.MaxCSVRows || l <= 1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var arcs []*model.ArcImportParam
	aidMap := make(map[int64]int64, len(arcs))
	for _, v := range records {
		arc := new(model.ArcImportParam)
		// aid
		if aid, err := strconv.ParseInt(v[0], 10, 64); err != nil || aid <= 0 {
			log.Warn("arcImportCSV strconv.ParseInt(%s) error(%v)", v[0], err)
			continue
		} else {
			if _, ok := aidMap[aid]; ok {
				continue
			}
			arc.Aid = aid
			aidMap[aid] = aid
		}
		// gids
		if gids, err := xstr.SplitInts(v[1]); err != nil {
			log.Warn("arcImportCSV gids xstr.SplitInts(%s) aid(%d) error(%v)", v[1], arc.Aid, err)
		} else {
			for _, id := range gids {
				if id > 0 {
					arc.Gids = append(arc.Gids, id)
				}
			}
		}
		// match ids
		if matchIDs, err := xstr.SplitInts(v[2]); err != nil {
			log.Warn("arcImportCSV match xstr.SplitInts(%s) aid(%d) error(%v)", v[2], arc.Aid, err)
		} else {
			for _, id := range matchIDs {
				if id > 0 {
					arc.MatchIDs = append(arc.MatchIDs, id)
				}
			}
		}
		// team ids
		if teamIDs, err := xstr.SplitInts(v[3]); err != nil {
			log.Warn("arcImportCSV team xstr.SplitInts(%s) aid(%d) error(%v)", v[3], arc.Aid, err)
		} else {
			for _, id := range teamIDs {
				if id > 0 {
					arc.TeamIDs = append(arc.TeamIDs, id)
				}
			}
		}
		// tag ids
		if tagIDs, err := xstr.SplitInts(v[4]); err != nil {
			log.Warn("arcImportCSV tag xstr.SplitInts(%s) aid(%d) error(%v)", v[4], arc.Aid, err)
		} else {
			for _, id := range tagIDs {
				if id > 0 {
					arc.TagIDs = append(arc.TagIDs, id)
				}
			}
		}
		// years
		if years, err := xstr.SplitInts(v[5]); err != nil {
			log.Warn("arcImportCSV year xstr.SplitInts(%s) aid(%d) error(%v)", v[5], arc.Aid, err)
		} else {
			for _, id := range years {
				if id > 0 {
					arc.Years = append(arc.Years, id)
				}
			}
		}
		arcs = append(arcs, arc)
	}
	if len(arcs) == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, esSvc.ArcImportCSV(c, arcs))
}

func batchDelArc(c *bm.Context) {
	v := new(struct {
		Aids []int64 `form:"aids,split" validate:"dive,gt=1,required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(nil, esSvc.BatchDelArc(c, v.Aids))
}
