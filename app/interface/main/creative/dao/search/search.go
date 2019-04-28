package search

import (
	"context"
	"strconv"
	"strings"

	"go-common/app/interface/main/creative/model/search"
	"go-common/library/database/elastic"
	"go-common/library/ecode"
	"go-common/library/log"
)

var (
	orderMap = map[string]string{
		"senddate": "pubdate",  //发布时间
		"click":    "click",    //点击数
		"scores":   "review",   //评论
		"stow":     "favorite", //收藏
		"dm_count": "dm_count", //弹幕
	}

	applyStateMap = map[string]string{
		"pending":   "pending",
		"processed": "processed",
		"neglected": "neglected",
	}
)

// ArchivesES search archives by es.
func (d *Dao) ArchivesES(c context.Context, mid int64, tid int16, keyword, order, class, ip string, pn, ps, coop int) (sres *search.Result, err error) {
	r := d.es.NewRequest("creative_archive_staff").Fields(
		"id",
		"pid",
		"typeid",
		"title",
		"state",
		"cover",
		"description",
		"duration",
		"pubdate",
	)
	r.Index("creative_archive").Pn(pn).Ps(ps).OrderScoreFirst(false)
	if mid > 0 && coop == 0 {
		cmbup := &elastic.Combo{}
		cmbup.ComboEQ([]map[string]interface{}{
			{"mid": mid},
		})
		r.WhereCombo(cmbup.MinEQ(1))
	} else if mid > 0 && coop == 1 {
		cmbup := &elastic.Combo{}
		cmbup.ComboEQ([]map[string]interface{}{
			{"mid": mid},
			{"staff_mid": mid},
		})
		r.WhereCombo(cmbup.MinEQ(1))
	}
	if keyword != "" { //筛选稿件标题或者描述
		r.WhereLike([]string{"title", "description"}, []string{keyword}, true, "low")
	}
	if tid > 0 {
		r.WhereEq("pid", tid)
	}
	if class != "" {
		if len(strings.Split(class, ",")) == 1 { //如果筛选全部则不传参数
			r.WhereEq("state", class) //state: is_pubing,pubed,not_pubed（全部） pubed (已通过) not_pubed(未通过) is_pubing（进行中）
		}
	}
	if order != "" {
		if o, ok := orderMap[order]; ok {
			r.Order(o, "desc")
		}
	} else {
		r.Order("pubdate", "desc") //默认按发布时间倒序排
	}
	log.Info("ArchivesES params(%s)", r.Params())
	var res = &search.ArcResult{}
	if err = r.Scan(c, res); err != nil {
		log.Error("ArchivesES r.Scan error(%v)", err)
		err = ecode.CreativeSearchErr
		return
	}
	sres = &search.Result{}
	sres.Page.Pn = res.Page.Num
	sres.Page.Ps = res.Page.Size
	sres.Page.Count = res.Page.Total
	if res.Result.PList != nil {
		sres.Class = &search.ClassCount{ //获取按稿件状态计数
			Pubed:    res.Result.PList.Pubed.Count,
			NotPubed: res.Result.PList.NotPubed.Count,
			Pubing:   res.Result.PList.IsPubing.Count,
		}
	}
	tcs := make(map[int16]*search.TypeCount)
	for _, v := range res.Result.TList { //获取按一级分区稿件计数
		if v != nil {
			key, err := strconv.ParseInt(v.Key, 10, 16)
			if err != nil {
				log.Error("strconv.ParseInt(%s)|error(%v)", v.Key, err)
				return nil, err
			}
			tid = int16(key)
			tc := &search.TypeCount{
				Tid:   tid,
				Count: int64(v.Count),
			}
			tcs[tid] = tc
		}
	}
	sres.Type = tcs
	for _, v := range res.Result.Vlist {
		if v != nil {
			sres.Aids = append(sres.Aids, v.ID)
		}
	}
	return
}

// ArchivesStaffES search staff applies by es.
func (d *Dao) ArchivesStaffES(c context.Context, mid int64, tid int16, keyword, state string, pn, ps int) (sres *search.StaffApplyResult, err error) {
	r := d.es.NewRequest("creative_archive_apply").Fields(
		"id",
		"pid",
		"typeid",
		"title",
		"state",
		"cover",
		"description",
		"duration",
		"pubdate",
	)
	r.Index("creative_archive").Pn(pn).Ps(ps).OrderScoreFirst(false)
	if mid > 0 {
		r.WhereEq("apply_staff.apply_staff_mid", mid)
	}
	if state != "" {
		if o, ok := applyStateMap[state]; ok {
			r.WhereEq("apply_staff.deal_state", o)
		}
	} else {
		r.WhereEq("apply_staff.deal_state", "pending")
	}
	if keyword != "" { //筛选稿件标题或者描述
		r.WhereLike([]string{"title", "description"}, []string{keyword}, true, "low")
	}
	if tid > 0 {
		r.WhereEq("pid", tid)
	}
	log.Info("ArchivesStaffES params(%s)", r.Params())
	var res = &search.ApplyResult{}
	if err = r.Scan(c, res); err != nil {
		log.Error("ArchivesStaffES r.Scan error(%v)", err)
		err = ecode.CreativeSearchErr
		return
	}
	sres = &search.StaffApplyResult{}
	sres.Page.Pn = res.Page.Num
	sres.Page.Ps = res.Page.Size
	sres.Page.Count = res.Page.Total
	//tlist
	if res.Result.ApplyPList != nil {
		sres.StateCount = &search.ApplyStateCount{ //获取按稿件状态计数
			Pending:   res.Result.ApplyPList.Pending.Count,
			Processed: res.Result.ApplyPList.Processed.Count,
			Neglected: res.Result.ApplyPList.Neglected.Count,
		}
	}
	// vlist
	tcs := make(map[int16]*search.TypeCount)
	for _, v := range res.Result.TList { //获取按一级分区稿件计数
		if v != nil {
			key, err := strconv.ParseInt(v.Key, 10, 16)
			if err != nil {
				log.Error("strconv.ParseInt(%s)|error(%v)", v.Key, err)
				return nil, err
			}
			tid = int16(key)
			tc := &search.TypeCount{
				Tid:   tid,
				Count: int64(v.Count),
			}
			tcs[tid] = tc
		}
	}
	sres.Type = tcs
	for _, v := range res.Result.Vlist {
		if v != nil {
			sres.Aids = append(sres.Aids, v.ID)
		}
	}
	return
}
