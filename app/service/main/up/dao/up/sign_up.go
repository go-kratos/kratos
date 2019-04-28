package up

import (
	"context"
	"go-common/app/service/main/up/dao/global"
	"go-common/app/service/main/up/model"
	"go-common/library/log"
	"go-common/library/xstr"

	"time"
)

const (
	// _getHighAllyUpsSql .
	_getHighAllyUpsSql = "SELECT id, mid, state, begin_date, end_date FROM sign_up WHERE mid IN (?) AND end_date >= ? AND state <> 100"
)

// GetHighAllyUps dao get high ally ups
func (d *Dao) GetHighAllyUps(c context.Context, mids []int64) (res []*model.SignUp, err error) {
	midStr := xstr.JoinInts(mids)
	now := time.Now().Format(DateLayout)
	rows, err := global.GetUpCrmDB().Query(c, _getHighAllyUpsSql, midStr, now)
	if err != nil {
		log.Error("d.GetHighAllyUps error (%v)", err)
		return
	}
	for rows.Next() {
		r := &model.SignUp{}
		err = rows.Scan(&r.ID, &r.Mid, &r.State, &r.BeginDate, &r.EndDate)
		if err != nil {
			log.Error("d.GetHighAllyUps scan error (%v)", err)
			return
		}
		res = append(res, r)
	}
	return
}
