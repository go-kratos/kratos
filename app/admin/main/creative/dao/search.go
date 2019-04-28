package dao

import (
	"context"
	"go-common/app/admin/main/creative/model/academy"
	"go-common/library/database/elastic"
	"go-common/library/log"
)

// ArchivesWithES search archives by es.
func (d *Dao) ArchivesWithES(c context.Context, aca *academy.EsParam) (res *academy.SearchResult, err error) {
	r := d.es.NewRequest("academy_archive").Fields("oid", "tid")
	r.Index("academy_archive").WhereEq("state", academy.StateNormal).WhereEq("business", aca.Business).Pn(aca.Pn).Ps(aca.Ps).Order("id", "desc")
	if aca.Business == academy.BusinessForArchvie && aca.State != academy.DefaultState { //arc_state 稿件原始状态 state 创作学院稿件状态
		r.WhereEq("arc_state", aca.State)
	}
	if aca.Business == academy.BusinessForArticle { //只筛选开放浏览的专栏
		r.WhereEq("arc_state", 0).WhereEq("deleted_time", 0)
	}
	if aca.Keyword != "" {
		r.WhereLike([]string{"title", "tid_name"}, []string{aca.Keyword}, true, "low")
	}
	if aca.Uname != "" {
		r.WhereLike([]string{"uname"}, []string{aca.Uname}, true, "low")
	}
	if aca.OID > 0 {
		r.WhereEq("oid", aca.OID)
	}

	if len(aca.TidsMap) > 0 {
		for _, v := range aca.TidsMap {
			cmb := &elastic.Combo{}
			tids := make([]interface{}, 0, len(v))
			for _, tid := range v {
				tids = append(tids, tid)
			}
			cmb.ComboIn([]map[string][]interface{}{
				{"tid": tids},
			}).MinIn(1).MinAll(1)
			r.WhereCombo(cmb)
		}
	}

	if aca.Business == academy.BusinessForArchvie {
		if aca.Copyright != 3 {
			r.WhereEq("copyright", aca.Copyright) //投稿类型
		} else {
			r.WhereIn("copyright", []int8{0, 1, 2})
		}
	}
	res = &academy.SearchResult{}
	if err = r.Scan(c, res); err != nil {
		log.Error("ArchivesWithES r.Scan params(%s)|error(%v)", r.Params(), err)
	}
	log.Info("ArchivesWithES params(%s)", r.Params())
	return
}
