package academy

import (
	"context"
	"time"

	"go-common/app/interface/main/creative/model/academy"
	"go-common/library/database/elastic"
	"go-common/library/ecode"
	"go-common/library/log"
)

// ArchivesWithES search archives by es.
func (d *Dao) ArchivesWithES(c context.Context, aca *academy.EsParam) (res *academy.SearchResult, err error) {
	r := d.es.NewRequest("academy_archive").Index("academy_archive").Fields("oid", "tid", "business")

	if aca.Business > 0 {
		r.WhereEq("business", aca.Business)
	}
	r.WhereEq("state", 0).Pn(aca.Pn).Ps(aca.Ps) //state 创作学院稿件状态
	r.WhereEq("arc_state", 0)                   //arc_state 原始稿件状态(视频、专栏)
	r.WhereEq("deleted_time", 0)                //过滤删除的专栏

	if aca.Seed > 0 {
		r.OrderRandomSeed(time.Unix(aca.Seed, 0).Format("2006-01-02 15:04:05")) //随机推荐
	}
	if aca.Keyword != "" {
		r.WhereLike([]string{"title", "tid_name"}, []string{aca.Keyword}, true, "low").Highlight(true)
	}
	if aca.Order != "" {
		r.Order(aca.Order, "desc").OrderScoreFirst(false) //order: click (最多点击数), fav(最多收藏数), pubtime(最新发布时间), hot(最热值)
	}

	if aca.Duration > 0 { //h5 时长筛选 1(1-10分钟 ) 2(10-30分钟) 3(30-60分钟) 4(60分钟+)
		switch aca.Duration {
		case 1:
			r.WhereRange("duration", 1*60, 10*60, elastic.RangeScopeLcRc)
		case 2:
			r.WhereRange("duration", 10*60, 30*60, elastic.RangeScopeLcRc)
		case 3:
			r.WhereRange("duration", 30*60, 60*60, elastic.RangeScopeLcRc)
		case 4:
			r.WhereRange("duration", 60*60, nil, elastic.RangeScopeLcRo)
		}

		if aca.Order == "" {
			r.Order("duration", "desc")
		}
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
	res = &academy.SearchResult{}
	log.Info("ArchivesWithES r.Scan params(%s)", r.Params())
	if err = r.Scan(c, res); err != nil {
		log.Error("ArchivesWithES r.Scan|error(%v)", err)
		err = ecode.CreativeSearchErr
		return
	}
	return
}

//Keywords get all keywords.
func (d *Dao) Keywords(c context.Context) (res []*academy.SearchKeywords, err error) {
	_getKWSQL := "SELECT id, rank, parent_id, state, name, comment FROM academy_search_keywords WHERE state=0 ORDER BY rank ASC"
	rows, err := d.db.Query(c, _getKWSQL)
	if err != nil {
		log.Error("Keywords d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	res = make([]*academy.SearchKeywords, 0)
	for rows.Next() {
		o := &academy.SearchKeywords{}
		if err = rows.Scan(&o.ID, &o.Rank, &o.ParentID, &o.State, &o.Name, &o.Comment); err != nil {
			log.Error("Keywords rows.Scan error(%v)", err)
			return
		}
		if o.Name == "" {
			continue
		}
		res = append(res, o)
	}
	return
}
