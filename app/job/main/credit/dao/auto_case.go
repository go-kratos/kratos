package dao

import (
	"context"
	"database/sql"

	"go-common/app/job/main/credit/model"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_autoCaseConfigSQL = "SELECT reasons,report_score,likes FROM blocked_auto_case WHERE platform = ?"
)

// AutoCaseConf  get auto case config.
func (d *Dao) AutoCaseConf(c context.Context, otype int8) (ac *model.AutoCaseConf, err error) {
	ac = &model.AutoCaseConf{}
	row := d.db.QueryRow(c, _autoCaseConfigSQL, otype)
	if err = row.Scan(&ac.ReasonStr, &ac.ReportScore, &ac.Likes); err != nil {
		if err != sql.ErrNoRows {
			log.Error("d.AutoCaseConf err(%v)", err)
			return
		}
		ac = nil
		err = nil
	}
	if ac != nil && ac.ReasonStr != "" {
		var reasonSlice []int64
		if reasonSlice, err = xstr.SplitInts(ac.ReasonStr); err != nil {
			log.Error("xstr.SplitInts(%s) err(%v)", ac.ReasonStr, err)
			return
		}
		ac.Reasons = make(map[int8]struct{}, len(reasonSlice))
		for _, v := range reasonSlice {
			ac.Reasons[int8(v)] = struct{}{}
		}
	}
	return
}
