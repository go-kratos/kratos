package income

import (
	"context"
	"time"

	model "go-common/app/job/main/growup/model/income"

	"go-common/library/log"
	xtime "go-common/library/time"
)

const (
	_columnChargeSQL = "SELECT id,aid,title,mid,tag_id,upload_time,inc_charge,view_c,date FROM column_daily_charge WHERE id > ? AND date = ? AND inc_charge > 0 ORDER BY id LIMIT ?"
)

// ColumnDailyCharge get column daily charge by date
func (d *Dao) ColumnDailyCharge(c context.Context, date time.Time, id int64, limit int) (columns []*model.ColumnCharge, err error) {
	columns = make([]*model.ColumnCharge, 0)
	rows, err := d.db.Query(c, _columnChargeSQL, id, date, limit)
	if err != nil {
		log.Error("ColumnDailyCharge d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		column := &model.ColumnCharge{}
		var uploadTime int64
		err = rows.Scan(&column.ID, &column.ArticleID, &column.Title, &column.MID, &column.TagID, &uploadTime, &column.IncCharge, &column.IncViewCount, &column.Date)
		if err != nil {
			log.Error("ColumnDailyCharge rows.Scan error(%v)", err)
			return
		}
		column.UploadTime = xtime.Time(uploadTime)
		columns = append(columns, column)
	}
	return
}
