package search

import (
	"context"
	"fmt"
	"go-common/app/interface/main/creative/dao/tool"
	resMdl "go-common/app/interface/main/creative/model/resource"
	"go-common/app/interface/main/creative/model/search"
	"go-common/library/database/elastic"
	"go-common/library/ecode"
	"go-common/library/log"
	"time"
)

var (
	// ReplyOrderMap map
	ReplyOrderMap = map[string]string{
		"ctime": "ctime",
		"like":  "like",
		"count": "count",
	}
	daysMap = map[string]int{
		filterCtimeOneDayAgo:   -1,
		filterCtimeOneWeekAgo:  -7,
		filterCtimeOneMonthAgo: -30,
		filterCtimeOneYearAgo:  -365,
	}
	filterCtimeOneDayAgo   = "0"
	filterCtimeOneWeekAgo  = "1"
	filterCtimeOneMonthAgo = "2"
	filterCtimeOneYearAgo  = "3"
)

// ReplyES fn
// order: ctime/like/count
// filter: 1/2/3
func (d *Dao) ReplyES(c context.Context, p *search.ReplyParam) (sres *search.Replies, err error) {
	r := d.es.NewRequest("creative_reply").Fields(
		"count",
		"ctime",
		"floor",
		"hate",
		"id",
		"like",
		"message",
		"mid",
		"mtime",
		"o_mid",
		"oid",
		"parent",
		"rcount",
		"root",
		"state",
		"type",
	)
	idx := fmt.Sprintf("%s_%02d", "creative_reply", p.OMID%100)
	r.Index(idx).Pn(p.Pn).Ps(p.Ps).OrderScoreFirst(true)
	if p.FilterCtime != "" {
		if dayNum, ok := daysMap[p.FilterCtime]; ok {
			begin := time.Now().AddDate(0, 0, dayNum).Format("2006-01-02 15:04:05")
			r.WhereRange("ctime", begin, "", elastic.RangeScopeLcRc)
		}
	}
	if p.IsReport > 0 {
		r.WhereEq("isreport", p.IsReport)
	}
	// 如果指定了oid就不需要传递o_mid
	if p.OID > 0 {
		r.WhereEq("oid", p.OID)
	} else {
		r.WhereEq("o_mid", p.OMID)
	}
	if p.Kw != "" {
		r.WhereLike([]string{"message"}, []string{p.Kw}, true, elastic.LikeLevelLow)
	}
	if p.Type > 0 {
		r.WhereEq("type", p.Type)
	}
	if p.ResMdlPlat == resMdl.PlatIPad {
		r.WhereIn("type", []int8{search.Article})
		r.WhereNot(elastic.NotTypeIn, "type")
	}
	r.WhereIn("state", []int{0, 1, 2, 5, 6})
	if p.Order != "" {
		if o, ok := ReplyOrderMap[p.Order]; ok {
			r.Order(o, "desc")
		} else {
			r.Order("ctime", "desc") //默认按发布时间倒序排
		}
	} else {
		r.Order("ctime", "desc") //默认按发布时间倒序排
	}
	log.Info("ReplyES params(%+v)|p(%+v)", r.Params(), p)
	var res = &search.ReliesES{}
	if err = r.Scan(c, res); err != nil {
		log.Error("ReplyES r.Scan p(%+v)|error(%v)", p, err)
		err = ecode.CreativeSearchErr
		return
	}
	sres = &search.Replies{
		Order:      res.Order,
		Keyword:    p.Kw,
		Repliers:   []int64{},
		DeriveOids: []int64{},
		DeriveIds:  []int64{},
		Oids:       []int64{},
		TyOids:     make(map[int][]int64),
		Result:     make([]*search.Reply, 0),
	}
	for _, v := range res.Result {
		sres.Result = append(sres.Result, &search.Reply{
			Message: v.Message,
			ID:      v.ID,
			Floor:   v.Floor,
			Count:   v.Count,
			Root:    v.Root,
			Oid:     v.Oid,
			CTime:   v.CTime,
			MTime:   v.MTime,
			State:   v.State,
			Parent:  v.Parent,
			Mid:     v.Mid,
			Like:    v.Like,
			Type:    v.Type,
		})
	}
	if res.Page != nil {
		sres.Total = res.Page.Total
		sres.PageCount = res.Page.Num
	}
	replyMids := make(map[int64]int64, len(res.Result))
	for _, v := range res.Result {
		_, ok := replyMids[v.Mid]
		if !ok {
			sres.Repliers = append(sres.Repliers, v.Mid)
			replyMids[v.Mid] = v.Mid
		}
		sres.Oids = append(sres.Oids, v.Oid)
		sres.TyOids[v.Type] = sres.Oids
		if v.Root > 0 {
			sres.DeriveOids = append(sres.DeriveOids, v.Oid)
			sres.DeriveIds = append(sres.DeriveIds, v.Parent)
		}
	}
	sres.Oids = tool.DeDuplicationSlice(sres.Oids)
	for k, v := range sres.TyOids {
		sres.TyOids[k] = tool.DeDuplicationSlice(v)
	}
	return
}
