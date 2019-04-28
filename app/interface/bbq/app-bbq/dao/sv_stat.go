package dao

import (
	"context"
	xsql "database/sql"
	"go-common/app/interface/bbq/app-bbq/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_insertVideoStatistics      = "insert into video_statistics (`svid`,`play`,`subtitles`,`like`,`share`,`report`) values (?,?,?,?,?,?)"
	_incrVideoStatisticsLike    = "update video_statistics set `like` = `like` + 1 where `svid` = ?"
	_decrVideoStatisticsLike    = "update video_statistics set `like` = `like` - 1 where `svid` = ? and `like` > 0"
	_insertVideoStatisticsShare = "INSERT INTO video_statistics (`svid`,`share`) VALUES (?,1) ON DUPLICATE KEY UPDATE `share` = `share` + 1;"
	_incrVideoStatisticsShare   = "UPDATE `video_statistics` SET `share` = `share` + 1 WHERE `svid` = ?;"
	_selectVideoStatisticsLike  = "select `svid`,`play`,`subtitles`,`like`,`share`,`report` from video_statistics where `svid` = ?"
	_videoStatisticsShare       = "select `share`+1 from video_statistics where `svid`= ?"
)

//TxIncrVideoStatisticsLike .
func (d *Dao) TxIncrVideoStatisticsLike(tx *sql.Tx, svid int64) (num int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(_incrVideoStatisticsLike, svid); err != nil {
		log.Error("update video_statistics like err(%v)", err)
		return
	}
	return res.RowsAffected()
}

//TxDecrVideoStatisticsLike .
func (d *Dao) TxDecrVideoStatisticsLike(tx *sql.Tx, svid int64) (num int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(_decrVideoStatisticsLike, svid); err != nil {
		log.Error("update video_statistics like err(%v)", err)
		return
	}
	return res.RowsAffected()
}

//IncrVideoStatisticsShare .
func (d *Dao) IncrVideoStatisticsShare(ctx context.Context, svid int64) (num int32, err error) {
	row := d.db.QueryRow(ctx, _videoStatisticsShare, svid)
	err = row.Scan(&num)
	if err != nil || num == 0 {
		if _, err = d.db.Exec(ctx, _insertVideoStatisticsShare, svid); err != nil {
			log.Error("insert video_statistics share error: [%+v] svid: [%d]", err, svid)
		}
		num = 1
	} else if _, err = d.db.Exec(ctx, _incrVideoStatisticsShare, svid); err != nil {
		log.Error("update video_statistics share error: [%+v] svid: [%d]", err, svid)
	}
	return
}

//TxAddVideoStatistics .
func (d *Dao) TxAddVideoStatistics(tx *sql.Tx, videoStatistics *model.VideoStatistics) (num int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(_insertVideoStatistics, videoStatistics.SVID, videoStatistics.Play, videoStatistics.Subtitles, videoStatistics.Like, videoStatistics.Share, videoStatistics.Report); err != nil {
		log.Error("insert video_statistics  err(%v)", err)
		return
	}
	return res.LastInsertId()
}

//RawVideoStatistics .
func (d *Dao) RawVideoStatistics(c context.Context, svid int64) (res *model.VideoStatistics, err error) {
	res = new(model.VideoStatistics)
	row := d.db.QueryRow(c, _selectVideoStatisticsLike, svid)
	err = row.Scan(&res.SVID, &res.Play, &res.Subtitles, &res.Like, &res.Share, &res.Report)
	if err == sql.ErrNoRows {
		err = nil
	}
	return
}
