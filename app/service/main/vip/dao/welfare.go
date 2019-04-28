package dao

import (
	"context"
	xsql "database/sql"
	"fmt"
	"strings"

	"go-common/app/service/main/vip/model"
	"go-common/library/database/sql"

	"github.com/pkg/errors"
)

const (
	_getWelfareListSQL              = "SELECT id, welfare_name, homepage_uri, backdrop_uri, tid, rank FROM vip_welfare WHERE state = 0 AND tid = ? AND recommend in (0,?) AND stime < ? AND etime > ? ORDER BY rank LIMIT ?,?"
	_countWelfareSQL                = "SELECT COUNT(1) FROM vip_welfare WHERE state = 0 AND tid = ? AND recommend in (0,?) AND stime < ? AND etime > ?"
	_getRecommendWelfareSQL         = "SELECT id, welfare_name, homepage_uri, backdrop_uri, tid, rank FROM vip_welfare WHERE state = 0 AND recommend = ? AND stime < ? AND etime > ? ORDER BY rank LIMIT ?,?"
	_countRecommendWelfareSQL       = "SELECT COUNT(1) FROM vip_welfare WHERE state = 0 AND recommend = ? AND stime < ? AND etime > ?"
	_getWelfareTypeListSQL          = "SELECT id, name FROM vip_welfare_type WHERE state = 0"
	_getWelfareInfoSQL              = "SELECT id, welfare_name, welfare_desc, receive_rate, homepage_uri, backdrop_uri, usage_form, vip_type, stime, etime FROM vip_welfare WHERE id = ?"
	_getWelfareBatchSQL             = "SELECT id, received_count, count, vtime FROM vip_welfare_code_batch WHERE state = 0 AND wid = ?"
	_getReceivedCodeSQL             = "SELECT id, mtime FROM vip_welfare_code WHERE state = 0 AND wid = ? AND mid = ?"
	_getWelfareCodeListSQL          = "SELECT id, bid, code FROM vip_welfare_code WHERE mid = 0 AND state = 0 AND wid = ? AND bid IN (%s) ORDER BY id limit 10"
	_receiveWelfareSQL              = "UPDATE vip_welfare_code SET mid = ? where id = ? and mid = 0"
	_updateStockSQL                 = "UPDATE vip_welfare_code_batch SET received_count = received_count+1 WHERE received_count < count AND id = ? "
	_getMyWelfareSQL                = "SELECT a.id, a.welfare_name, a.welfare_desc, a.usage_form, a.receive_uri, a.stime, a.etime, b.code FROM vip_welfare a LEFT JOIN vip_welfare_code b ON a.id = b.wid WHERE b.mid = ? ORDER BY b.mtime desc"
	_addReceiveRedirectWelfareSQL   = "INSERT INTO vip_welfare_code (wid, mid) VALUES (?, ?)"
	_countReceiveRedirectWelfareSQL = "SELECT COUNT(1) FROM vip_welfare_code WHERE wid = ? AND mid = ?"
	_insertReceiveRecordSQL         = "INSERT INTO vip_welfare_record (mid, wid, month, count) VALUES (?, ?, ?, 1)"
)

// GetWelfareList get welfare list
func (d *Dao) GetWelfareList(c context.Context, req *model.ArgWelfareList) (res []*model.WelfareListResp, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _getWelfareListSQL, req.Tid, req.Recommend, req.NowTime, req.NowTime, (req.Pn-1)*req.Ps, req.Ps); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		r := &model.WelfareListResp{}
		if err = rows.Scan(&r.ID, &r.Name, &r.HomepageUri, &r.BackdropUri, &r.Tid, &r.Rank); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res = append(res, r)
	}
	return
}

// CountWelfare count welfare list
func (d *Dao) CountWelfare(c context.Context, req *model.ArgWelfareList) (count int64, err error) {
	row := d.db.QueryRow(c, _countWelfareSQL, req.Tid, req.Recommend, req.NowTime, req.NowTime)

	if err = row.Scan(&count); err != nil {
		if sql.ErrNoRows == err {
			err = nil
			count = 0
			return
		}
		err = errors.WithStack(err)
	}

	return
}

// GetRecommendWelfare get recommend welfare
func (d *Dao) GetRecommendWelfare(c context.Context, req *model.ArgWelfareList) (res []*model.WelfareListResp, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _getRecommendWelfareSQL, req.Recommend, req.NowTime, req.NowTime, (req.Pn-1)*req.Ps, req.Ps); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		r := &model.WelfareListResp{}
		if err = rows.Scan(&r.ID, &r.Name, &r.HomepageUri, &r.BackdropUri, &r.Tid, &r.Rank); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res = append(res, r)
	}
	return
}

// CountRecommendWelfare count recommend welfare
func (d *Dao) CountRecommendWelfare(c context.Context, req *model.ArgWelfareList) (count int64, err error) {
	row := d.db.QueryRow(c, _countRecommendWelfareSQL, req.Recommend, req.NowTime, req.NowTime)

	if err = row.Scan(&count); err != nil {
		if sql.ErrNoRows == err {
			err = nil
			count = 0
			return
		}
		err = errors.WithStack(err)
	}

	return
}

// GetWelfareTypeList get welfare type list
func (d *Dao) GetWelfareTypeList(c context.Context) (res []*model.WelfareTypeListResp, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _getWelfareTypeListSQL); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		r := &model.WelfareTypeListResp{}
		if err = rows.Scan(&r.ID, &r.Name); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res = append(res, r)
	}
	return
}

// GetWelfareInfo get welfare info
func (d *Dao) GetWelfareInfo(c context.Context, id int64) (res *model.WelfareInfoResp, err error) {
	res = new(model.WelfareInfoResp)

	if err = d.db.QueryRow(c, _getWelfareInfoSQL, id).
		Scan(&res.ID, &res.Name, &res.Desc, &res.ReceiveRate, &res.HomepageUri, &res.BackdropUri, &res.UsageForm, &res.VipType, &res.Stime, &res.Etime); err != nil {
		if err == sql.ErrNoRows {
			res = nil
			err = nil
			return
		}
		err = errors.Wrapf(err, "dao GetWelfareInfo(%d)", id)
	}
	return
}

// GetWelfareBatch get welfare batch infos
func (d *Dao) GetWelfareBatch(c context.Context, wid int64) (res []*model.WelfareBatchResp, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _getWelfareBatchSQL, wid); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		r := &model.WelfareBatchResp{}
		if err = rows.Scan(&r.Id, &r.ReceivedCount, &r.Count, &r.Vtime); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res = append(res, r)
	}
	return
}

// GetReceivedCode get received code
func (d *Dao) GetReceivedCode(c context.Context, wid, mid int64) (res []*model.ReceivedCodeResp, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _getReceivedCodeSQL, wid, mid); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		r := &model.ReceivedCodeResp{}
		if err = rows.Scan(&r.ID, &r.Mtime); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res = append(res, r)
	}
	return
}

//UpdateWelfareCodeUser user receive welfare
func (d *Dao) UpdateWelfareCodeUser(c context.Context, tx *sql.Tx, id int, mid int64) (affectedRows int64, err error) {
	var (
		res xsql.Result
	)
	if res, err = tx.Exec(_receiveWelfareSQL, mid, id); err != nil {
		err = errors.WithStack(err)
		return
	}
	affectedRows, err = res.RowsAffected()
	return
}

//UpdateWelfareBatch reduce count
func (d *Dao) UpdateWelfareBatch(c context.Context, tx *sql.Tx, bid int) (err error) {
	if _, err = tx.Exec(_updateStockSQL, bid); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// GetWelfareCodeUnReceived get unReceived welfare code
func (d *Dao) GetWelfareCodeUnReceived(c context.Context, wid int64, bids []string) (res []*model.UnReceivedCodeResp, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, fmt.Sprintf(_getWelfareCodeListSQL, strings.Join(bids, ",")), wid); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		r := new(model.UnReceivedCodeResp)
		if err = rows.Scan(&r.Id, &r.Bid, &r.Code); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res = append(res, r)
	}

	return
}

// GetMyWelfare get my welfare infos
func (d *Dao) GetMyWelfare(c context.Context, mid int64) (res []*model.MyWelfareResp, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _getMyWelfareSQL, mid); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		r := &model.MyWelfareResp{}
		if err = rows.Scan(&r.Wid, &r.Name, &r.Desc, &r.UsageForm, &r.ReceiveUri, &r.Stime, &r.Etime, &r.Code); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res = append(res, r)
	}

	return
}

// AddReceiveRedirectWelfare add redirect url
func (d *Dao) AddReceiveRedirectWelfare(c context.Context, wid, mid int64) (err error) {
	if _, err = d.db.Exec(c, _addReceiveRedirectWelfareSQL, wid, mid); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// CountReceiveRedirectWelfare count it
func (d *Dao) CountReceiveRedirectWelfare(c context.Context, wid, mid int64) (count int64, err error) {
	row := d.db.QueryRow(c, _countReceiveRedirectWelfareSQL, wid, mid)

	if err = row.Scan(&count); err != nil {
		if sql.ErrNoRows == err {
			err = nil
			count = 0
			return
		}
		err = errors.WithStack(err)
	}

	return
}

// InsertReceiveRecord to prevent repeated receive
func (d *Dao) InsertReceiveRecord(c context.Context, tx *sql.Tx, mid, wid, monthYear int64) (err error) {
	if _, err = tx.Exec(_insertReceiveRecordSQL, mid, wid, monthYear); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}
