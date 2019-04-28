package archive

import (
	"context"

	arcmdl "go-common/app/service/main/archive/model/archive"
	"go-common/library/log"
)

const (
	_reportResultSQL = "SELECT aid,result,state,adminid,after_deal FROM archive_report_result WHERE is_show=1"
)

// ReportResults get all report
func (d *Dao) ReportResults(c context.Context) (res map[int64]*arcmdl.ReportResult, err error) {
	d.infoProm.Incr("ReportResults")
	rows, err := d.reportResultStmt.Query(c)
	if err != nil {
		log.Error("d.reportResultStmt.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	res = make(map[int64]*arcmdl.ReportResult)
	for rows.Next() {
		rr := &arcmdl.ReportResult{}
		if err = rows.Scan(&rr.Aid, &rr.Result, &rr.State, &rr.AdminID, &rr.AfterDeal); err != nil {
			log.Error("row.Scan error(%v)", err)
			d.errProm.Incr("report_result_db")
			return
		}
		res[rr.Aid] = rr
	}
	return
}
