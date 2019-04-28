package dao

import (
	"context"

	"go-common/app/admin/main/workflow/model"
	"go-common/app/admin/main/workflow/model/search"
	"go-common/library/log"
)

// ChallTagsCount will retrive challenge tags count from group ids and group by tag id
// group id to challenge id to tag count
func (d *Dao) ChallTagsCount(c context.Context, gids []int64) (counts map[int64]map[int64]int64, err error) {
	counts = make(map[int64]map[int64]int64, len(gids))
	if len(gids) <= 0 {
		return
	}

	rows, err := d.ReadORM.Table("workflow_chall").Where("gid IN (?)", gids).Select("gid,tid,count(tid)").
		Group("gid,tid").Rows()
	if err != nil {
		return
	}

	defer rows.Close()
	for rows.Next() {
		var gtc struct {
			Gid   int64
			Tid   int64
			Count int64
		}
		if err = rows.Scan(&gtc.Gid, &gtc.Tid, &gtc.Count); err != nil {
			return
		}
		if _, ok := counts[gtc.Gid]; !ok {
			counts[gtc.Gid] = make(map[int64]int64)
		}
		counts[gtc.Gid][gtc.Tid] = gtc.Count
	}

	return
}

// ChallTagsCountV3 .
func (d *Dao) ChallTagsCountV3(c context.Context, gids []int64) (counts map[int64]map[int64]int64, err error) {
	var result []*search.ChallSearchCommonData
	counts = make(map[int64]map[int64]int64, len(gids))
	cond := &search.ChallSearchCommonCond{
		Fields: []string{"gid", "tid"},
		Gids:   gids,
		States: []int64{int64(model.Pending)},
	}
	if result, err = d.SearchChallengeMultiPage(c, cond); err != nil {
		log.Error("s.dao.SearchChallengeMultiPage(%+v) error(%v)", cond, err)
		return
	}
	for _, r := range result {
		if _, ok := counts[r.Gid]; !ok {
			counts[r.Gid] = make(map[int64]int64)
		}
		counts[r.Gid][r.Tid]++
	}
	return
}
