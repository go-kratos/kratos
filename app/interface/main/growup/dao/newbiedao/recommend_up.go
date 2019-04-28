package newbiedao

import (
	"context"
	"go-common/app/interface/main/growup/model"

	"go-common/library/log"
)

const _RecommendUpSql = "SELECT mid, tid FROM recommend_up_white"

// GetRecommendUpList get recommend up list
func (d *Dao) GetRecommendUpList(c context.Context) error {
	recUps := make(map[int64]map[int64]*model.RecommendUp)
	rows, err := d.db.Query(c, _RecommendUpSql)
	if err != nil {
		log.Error("d.db.Query recommend up error(%v)", err)
		return err
	}
	defer rows.Close()
	for rows.Next() {
		recUp := new(model.RecommendUp)
		err = rows.Scan(&recUp.Mid, &recUp.Tid)
		if err != nil {
			log.Error("rows scan error(%v)", err)
			return err
		}

		if _, ok := recUps[recUp.Tid]; !ok {
			recUps[recUp.Tid] = make(map[int64]*model.RecommendUp)
		}
		recUps[recUp.Tid][recUp.Mid] = recUp
	}

	RecommendUpList = recUps
	return nil
}
