package dao

import (
	"context"
	"time"

	"go-common/app/interface/main/push-archive/model"
	"go-common/library/log"
)

const (
	_inStatisticsSQL      = "INSERT INTO `push_statistics` (`aid`, `group`, `type`, `mids`, `mids_counter`, `ctime`, `mtime`) VALUES(?,?,?,?,?,?,?);"
	_statisticsIDRangeSQL = "SELECT coalesce(min(id), 0), coalesce(max(id) , 0) FROM `push_statistics` WHERE `ctime` < ?"
	_delStatisticsByIDSQL = "DELETE FROM `push_statistics` WHERE `id` >=? AND `id`<=?;"
)

//SetStatistics 插入一条记录
func (d *Dao) SetStatistics(ctx context.Context, st *model.PushStatistic) (rows int64, err error) {
	res, err := d.setStatisticsStmt.Exec(ctx, st.Aid, st.Group, st.Type, st.Mids, st.MidsCounter, st.CTime, time.Now())
	if err != nil {
		log.Error("SetStatistics() d.setStatisticsStmt.Exec error(%v), pushstatistic(%v)", err, st)
		PromError("db:保存统计数据")
		return
	}

	rows, err = res.RowsAffected()
	return
}

//GetStatisticsIDRange get id range
func (d *Dao) GetStatisticsIDRange(ctx context.Context, deadline time.Time) (min int64, max int64, err error) {
	if err = d.db.QueryRow(ctx, _statisticsIDRangeSQL, deadline).Scan(&min, &max); err != nil {
		log.Error("GetStatisticsIDRange() error(%v), deadline(%v)", err, deadline)
		PromError("db:查询统计数据")
	}
	return
}

//DelStatisticsByID delete by id range
func (d *Dao) DelStatisticsByID(ctx context.Context, min, max int64) (rows int64, err error) {
	res, err := d.db.Exec(ctx, _delStatisticsByIDSQL, min, max)
	if err != nil {
		log.Error("DelStatistics() error(%v), min(%d) max(%d)", err, min, max)
		PromError("db:删除统计数据")
		return
	}
	rows, err = res.RowsAffected()
	return
}
