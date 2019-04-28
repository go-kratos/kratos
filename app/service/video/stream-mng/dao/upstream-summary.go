package dao

import (
	"context"
	"github.com/pkg/errors"
	"go-common/app/service/video/stream-mng/model"
	"go-common/library/database/sql"
)

// 存储上行调度信息

const (
	_insertUpStreamInfo         = "INSERT INTO `upstream_info` (room_id,cdn,platform,ip,country,city,isp) values(?,?,?,?,?,?,?)"
	_getSummaryUpStreamRtmp     = "SELECT cdn, count(id) as value FROM upstream_info where mtime >= FROM_UNIXTIME(?) AND mtime <= FROM_UNIXTIME(?) group by cdn"
	_getSummaryUpStreamISP      = "SELECT isp, count(id) as value FROM upstream_info where mtime >= FROM_UNIXTIME(?) AND mtime <= FROM_UNIXTIME(?) group by isp"
	_getSummaryUpStreamCountry  = "SELECT country, count(id) as value FROM upstream_info where mtime >= FROM_UNIXTIME(?) AND mtime <= FROM_UNIXTIME(?) group by country"
	_getSummaryUpStreamPlatform = "SELECT platform, count(id) as value FROM upstream_info where mtime >= FROM_UNIXTIME(?) AND mtime <= FROM_UNIXTIME(?) group by platform"
	_getSummaryUpStreamCity     = "SELECT city, count(id) as value FROM upstream_info where mtime >= FROM_UNIXTIME(?) AND mtime <= FROM_UNIXTIME(?) group by city"
)

// CreateUpStreamDispatch 创建一条上行调度信息
func (d *Dao) CreateUpStreamDispatch(c context.Context, info *model.UpStreamInfo) error {
	_, err := d.stmtUpStreamDispatch.Exec(c, info.RoomID, info.CDN, info.PlatForm, info.IP, info.Country, info.City, info.ISP)

	return err
}

// GetSummaryUpStreamRtmp 得到统计信息
func (d *Dao) GetSummaryUpStreamRtmp(c context.Context, start int64, end int64) (infos []*model.SummaryUpStreamRtmp, err error) {
	res := []*model.SummaryUpStreamRtmp{}
	var rows *sql.Rows

	if rows, err = d.tidb.Query(c, _getSummaryUpStreamRtmp, start, end); err != nil {
		err = errors.WithStack(err)
		return
	}

	defer rows.Close()
	for rows.Next() {
		info := new(model.SummaryUpStreamRtmp)
		if err = rows.Scan(&info.CDN, &info.Count); err != nil {
			err = errors.WithStack(err)
			infos = nil
			return
		}
		res = append(res, info)
	}
	err = rows.Err()
	return res, err
}

// GetSummaryUpStreamISP 得到ISP统计信息
func (d *Dao) GetSummaryUpStreamISP(c context.Context, start int64, end int64) (infos []*model.SummaryUpStreamRtmp, err error) {
	res := []*model.SummaryUpStreamRtmp{}
	var rows *sql.Rows

	if rows, err = d.tidb.Query(c, _getSummaryUpStreamISP, start, end); err != nil {
		err = errors.WithStack(err)
		return
	}

	defer rows.Close()
	for rows.Next() {
		info := new(model.SummaryUpStreamRtmp)
		if err = rows.Scan(&info.ISP, &info.Count); err != nil {
			err = errors.WithStack(err)
			infos = nil
			return
		}
		res = append(res, info)
	}
	err = rows.Err()
	return res, err
}

// GetSummaryUpStreamCountry 得到Country统计信息
func (d *Dao) GetSummaryUpStreamCountry(c context.Context, start int64, end int64) (infos []*model.SummaryUpStreamRtmp, err error) {
	res := []*model.SummaryUpStreamRtmp{}
	var rows *sql.Rows

	if rows, err = d.tidb.Query(c, _getSummaryUpStreamCountry, start, end); err != nil {
		err = errors.WithStack(err)
		return
	}

	defer rows.Close()
	for rows.Next() {
		info := new(model.SummaryUpStreamRtmp)
		if err = rows.Scan(&info.Country, &info.Count); err != nil {
			err = errors.WithStack(err)
			infos = nil
			return
		}
		res = append(res, info)
	}
	err = rows.Err()
	return res, err
}

// GetSummaryUpStreamPlatform 得到Platform统计信息
func (d *Dao) GetSummaryUpStreamPlatform(c context.Context, start int64, end int64) (infos []*model.SummaryUpStreamRtmp, err error) {
	res := []*model.SummaryUpStreamRtmp{}
	var rows *sql.Rows

	if rows, err = d.tidb.Query(c, _getSummaryUpStreamPlatform, start, end); err != nil {
		err = errors.WithStack(err)
		return
	}

	defer rows.Close()
	for rows.Next() {
		info := new(model.SummaryUpStreamRtmp)
		if err = rows.Scan(&info.PlatForm, &info.Count); err != nil {
			err = errors.WithStack(err)
			infos = nil
			return
		}
		res = append(res, info)
	}
	err = rows.Err()
	return res, err
}

// GetSummaryUpStreamCity 得到City统计信息
func (d *Dao) GetSummaryUpStreamCity(c context.Context, start int64, end int64) (infos []*model.SummaryUpStreamRtmp, err error) {
	res := []*model.SummaryUpStreamRtmp{}
	var rows *sql.Rows

	if rows, err = d.tidb.Query(c, _getSummaryUpStreamCity, start, end); err != nil {
		err = errors.WithStack(err)
		return
	}

	defer rows.Close()
	for rows.Next() {
		info := new(model.SummaryUpStreamRtmp)
		if err = rows.Scan(&info.City, &info.Count); err != nil {
			err = errors.WithStack(err)
			infos = nil
			return
		}
		res = append(res, info)
	}
	err = rows.Err()
	return res, err
}
