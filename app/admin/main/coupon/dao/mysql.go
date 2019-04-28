package dao

import (
	"bytes"
	"context"
	xsql "database/sql"
	"fmt"
	"strconv"
	"strings"

	"go-common/app/admin/main/coupon/model"
	"go-common/library/database/sql"

	"github.com/pkg/errors"
)

const (
	_addbatch                       = "INSERT INTO coupon_batch_info(app_id,name,batch_token,max_count,current_count,start_time,expire_time,ver,ctime,limit_count,operator)VALUES(?,?,?,?,?,?,?,?,?,?,?);"
	_addAllowancebatch              = "INSERT INTO coupon_batch_info(app_id,name,batch_token,max_count,current_count,start_time,expire_time,expire_day,ver,ctime,limit_count,operator,full_amount,amount,state,coupon_type,platform_limit,product_limit_month,product_limit_renewal)VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);"
	_batchList                      = "SELECT id,app_id,name,batch_token,max_count,current_count,start_time,expire_time,expire_day,ver,ctime,mtime,limit_count,operator,full_amount,amount,state,coupon_type,platform_limit,product_limit_month,product_limit_renewal FROM coupon_batch_info WHERE 1=1 "
	_appAppInfoSQL                  = "SELECT id,name,app_key,notify_url,ctime,mtime FROM coupon_app_info;"
	_updateAllowanceBatchSQL        = "UPDATE coupon_batch_info SET app_id = ?,name = ?,max_count = ?,limit_count = ?,operator = ?,platform_limit = ?,product_limit_month = ?,product_limit_renewal = ? WHERE id = ?;"
	_updateBatchStatusSQL           = "UPDATE coupon_batch_info SET state = ?, operator = ?  WHERE id = ?;"
	_batchInfoSQL                   = "SELECT id,app_id,name,batch_token,max_count,current_count,start_time,expire_time,expire_day,ver,ctime,mtime,limit_count,operator,full_amount,amount,state,coupon_type,platform_limit,product_limit_month,product_limit_renewal FROM coupon_batch_info WHERE batch_token = ? "
	_batchInfoByIDSQL               = "SELECT id,app_id,name,batch_token,max_count,current_count,start_time,expire_time,expire_day,ver,ctime,mtime,limit_count,operator,full_amount,amount,state,coupon_type,platform_limit,product_limit_month,product_limit_renewal FROM coupon_batch_info WHERE id = ? "
	_updateAllowanceStateSQL        = "UPDATE coupon_allowance_info_%02d SET state = ?,ver = ver+1 WHERE coupon_token = ? AND ver = ?;"
	_couponAllowanceByTokenSQL      = "SELECT id,coupon_token,mid,state,start_time,expire_time,origin,ver,batch_token,order_no,amount,full_amount,ctime,mtime FROM coupon_allowance_info_%02d WHERE coupon_token = ?;"
	_addCouponAllowanceChangeLogSQL = "INSERT INTO coupon_allowance_change_log_%02d (coupon_token,order_no,mid,state,ctime, change_type) VALUES(?,?,?,?,?,?);"
	_couponAllowancePageSQL         = "SELECT id,coupon_token,mid,state,start_time,expire_time,origin,ver,batch_token,order_no,amount,full_amount,ctime,mtime,remark  FROM coupon_allowance_info_%02d WHERE  mid = ? %s ORDER BY id DESC"
	_batchAddAllowanceCouponSQL     = "INSERT INTO coupon_allowance_info_%02d(coupon_token,mid,state,start_time,expire_time,origin,batch_token,amount,full_amount,ctime,app_id) VALUES "
	_updateBatchSQL                 = "UPDATE coupon_batch_info SET current_count = current_count + ? WHERE batch_token = ?;"
	_updateCodeBatchSQL             = "UPDATE coupon_batch_info SET app_id = ?,name = ?,limit_count = ?,operator = ?,platform_limit = ?,product_limit_month = ?,product_limit_renewal = ? WHERE id = ?;"

	//view
	_addViewBatchSQL    = "INSERT INTO coupon_batch_info(app_id,name,batch_token,max_count,current_count,start_time,expire_time,ver,limit_count,coupon_type,operator)VALUES(?,?,?,?,?,?,?,?,?,?,?)"
	_updateViewBatchSQL = "UPDATE coupon_batch_info SET app_id=?,name=?,max_count=?,limit_count=?,operator = ?,ver=? WHERE id=?"
	_updateViewSQL      = "UPDATE coupon_info_%02d SET state=? WHERE coupon_token=?"
	_searchViewSQL      = "SELECT coupon_token,mid,batch_token,state,order_no,oid,start_time,expire_time,ctime,mtime FROM coupon_info_%02d WHERE 1=1"
	_searchViewCountSQL = "SELECT COUNT(1) FROM coupon_info_%02d WHERE 1=1"
	_viewInfoSQL        = "SELECT coupon_token,mid,batch_token,state,order_no,oid,start_time,expire_time,ctime FROM coupon_info_%02d WHERE coupon_token=?"
	_addViewChangeLog   = "INSERT INTO coupon_change_log_%02d(coupon_token,mid,state) VALUES(?,?,?)"
)

// BeginTran begin transaction.
func (d *Dao) BeginTran(c context.Context) (*sql.Tx, error) {
	return d.db.Begin(c)
}

func hitAllowanceInfo(mid int64) int64 {
	return mid % 10
}

func hitAllowanceChangeLog(mid int64) int64 {
	return mid % 10
}

func hitViewInfo(mid int64) int64 {
	return mid % 100
}

// BatchList query batch list.
func (d *Dao) BatchList(c context.Context, appid int64, t int8) (res []*model.CouponBatchInfo, err error) {
	var (
		rows *sql.Rows
		sql  = _batchList
	)
	if appid != 0 {
		sql += fmt.Sprintf(" AND `app_id` = %d", appid)
	}
	if t != 0 {
		sql += fmt.Sprintf(" AND `coupon_type` = %d", t)
	}
	if rows, err = d.db.Query(c, sql); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.CouponBatchInfo{}
		if err = rows.Scan(&r.ID, &r.AppID, &r.Name, &r.BatchToken, &r.MaxCount, &r.CurrentCount, &r.StartTime, &r.ExpireTime, &r.ExpireDay, &r.Ver,
			&r.Ctime, &r.Mtime, &r.LimitCount, &r.Operator, &r.FullAmount, &r.Amount, &r.State, &r.CouponType, &r.PlatformLimit, &r.ProdLimMonth, &r.ProdLimRenewal); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

// BatchViewList query batch list.
func (d *Dao) BatchViewList(c context.Context, appid int64, batchToken string, t int8) (res []*model.CouponBatchInfo, err error) {
	var (
		rows *sql.Rows
		sql  = _batchList
	)
	if appid != 0 {
		sql += fmt.Sprintf(" AND app_id = %d", appid)
	}
	if t != 0 {
		sql += fmt.Sprintf(" AND coupon_type = %d", t)
	}
	if len(batchToken) > 0 {
		sql += fmt.Sprintf(" AND batch_token='%v'", batchToken)
	}
	if rows, err = d.db.Query(c, sql); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.CouponBatchInfo{}
		if err = rows.Scan(&r.ID, &r.AppID, &r.Name, &r.BatchToken, &r.MaxCount, &r.CurrentCount, &r.StartTime, &r.ExpireTime, &r.ExpireDay, &r.Ver,
			&r.Ctime, &r.Mtime, &r.LimitCount, &r.Operator, &r.FullAmount, &r.Amount, &r.State, &r.CouponType, &r.PlatformLimit, &r.ProdLimMonth, &r.ProdLimRenewal); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

// AddBatchInfo add batch info.
func (d *Dao) AddBatchInfo(c context.Context, b *model.CouponBatchInfo) (a int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, _addbatch, b.AppID, b.Name, b.BatchToken, b.MaxCount, b.CurrentCount, b.StartTime, b.ExpireTime,
		b.Ver, b.Ctime, b.LimitCount, b.Operator); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// AllAppInfo all app info.
func (d *Dao) AllAppInfo(c context.Context) (res []*model.AppInfo, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _appAppInfoSQL); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.AppInfo{}
		if err = rows.Scan(&r.ID, &r.Name, &r.Appkey, &r.NotifyURL, &r.Ctime, &r.Mtime); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

// AddAllowanceBatchInfo add allowance batch info.
func (d *Dao) AddAllowanceBatchInfo(c context.Context, b *model.CouponBatchInfo) (a int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, _addAllowancebatch, b.AppID, b.Name, b.BatchToken, b.MaxCount, b.CurrentCount, b.StartTime, b.ExpireTime, b.ExpireDay,
		b.Ver, b.Ctime, b.LimitCount, b.Operator, b.FullAmount, b.Amount, b.State, b.CouponType, b.PlatformLimit, b.ProdLimMonth, b.ProdLimRenewal); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// UpdateAllowanceBatchInfo update allowance batch info.
func (d *Dao) UpdateAllowanceBatchInfo(c context.Context, b *model.CouponBatchInfo) (a int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, _updateAllowanceBatchSQL, b.AppID, b.Name, b.MaxCount, b.LimitCount, b.Operator, b.PlatformLimit, b.ProdLimMonth, b.ProdLimRenewal, b.ID); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// UpdateCodeBatchInfo update code batch info.
func (d *Dao) UpdateCodeBatchInfo(c context.Context, b *model.CouponBatchInfo) (a int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, _updateCodeBatchSQL, b.AppID, b.Name, b.LimitCount, b.Operator, b.PlatformLimit, b.ProdLimMonth, b.ProdLimRenewal, b.ID); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// UpdateBatchStatus update batch status.
func (d *Dao) UpdateBatchStatus(c context.Context, status int8, operator string, id int64) (a int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, _updateBatchStatusSQL, status, operator, id); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
	}
	return
}

//BatchInfo batch info.
func (d *Dao) BatchInfo(c context.Context, token string) (r *model.CouponBatchInfo, err error) {
	var row *sql.Row
	r = new(model.CouponBatchInfo)
	row = d.db.Master().QueryRow(c, _batchInfoSQL, token)
	if err = row.Scan(&r.ID, &r.AppID, &r.Name, &r.BatchToken, &r.MaxCount, &r.CurrentCount, &r.StartTime, &r.ExpireTime, &r.ExpireDay, &r.Ver,
		&r.Ctime, &r.Mtime, &r.LimitCount, &r.Operator, &r.FullAmount, &r.Amount, &r.State, &r.CouponType, &r.PlatformLimit, &r.ProdLimMonth, &r.ProdLimRenewal); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			r = nil
			return
		}
		err = errors.WithStack(err)
		return
	}
	return
}

//BatchInfoByID batch info by id.
func (d *Dao) BatchInfoByID(c context.Context, id int64) (r *model.CouponBatchInfo, err error) {
	var row *sql.Row
	r = new(model.CouponBatchInfo)
	row = d.db.QueryRow(c, _batchInfoByIDSQL, id)
	if err = row.Scan(&r.ID, &r.AppID, &r.Name, &r.BatchToken, &r.MaxCount, &r.CurrentCount, &r.StartTime, &r.ExpireTime, &r.ExpireDay, &r.Ver,
		&r.Ctime, &r.Mtime, &r.LimitCount, &r.Operator, &r.FullAmount, &r.Amount, &r.State, &r.CouponType, &r.PlatformLimit, &r.ProdLimMonth, &r.ProdLimRenewal); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			r = nil
			return
		}
		err = errors.WithStack(err)
		return
	}
	return
}

// UpdateAllowanceStatus update allowance status.
func (d *Dao) UpdateAllowanceStatus(c context.Context, tx *sql.Tx, state int8, mid int64, token string, ver int64) (a int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(fmt.Sprintf(_updateAllowanceStateSQL, hitAllowanceInfo(mid)), state, token, ver); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// AllowanceByToken query coupon by token.
func (d *Dao) AllowanceByToken(c context.Context, mid int64, token string) (r *model.CouponAllowanceInfo, err error) {
	var row *sql.Row
	r = &model.CouponAllowanceInfo{}
	row = d.db.QueryRow(c, fmt.Sprintf(_couponAllowanceByTokenSQL, hitAllowanceInfo(mid)), token)
	if err = row.Scan(&r.ID, &r.CouponToken, &r.Mid, &r.State, &r.StartTime, &r.ExpireTime, &r.Origin, &r.Ver, &r.BatchToken,
		&r.OrderNO, &r.Amount, &r.FullAmount, &r.CTime, &r.MTime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			r = nil
			return
		}
		err = errors.WithStack(err)
		return
	}
	return
}

//InsertCouponAllowanceHistory insert coupon history .
func (d *Dao) InsertCouponAllowanceHistory(c context.Context, tx *sql.Tx, l *model.CouponAllowanceChangeLog) (a int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(fmt.Sprintf(_addCouponAllowanceChangeLogSQL, hitAllowanceChangeLog(l.Mid)), l.CouponToken, l.OrderNO, l.Mid, l.State, l.Ctime, l.ChangeType); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
	}
	return
}

//AllowanceList allowance list.
func (d *Dao) AllowanceList(c context.Context, arg *model.ArgAllowanceSearch) (res []*model.CouponAllowanceInfo, err error) {
	var (
		rows     *sql.Rows
		whereSQL = " "
	)
	if arg.AppID != 0 {
		whereSQL += fmt.Sprintf(" AND `app_id` = %d ", arg.AppID)
	}
	if arg.CouponToken != "" {
		whereSQL += fmt.Sprintf(" AND `coupon_token` = '%s' ", arg.CouponToken)
	}
	if arg.OrderNO != "" {
		whereSQL += fmt.Sprintf(" AND `order_no` = '%s' ", arg.OrderNO)
	}
	if arg.BatchToken != "" {
		whereSQL += fmt.Sprintf(" AND `batch_token` =  '%s' ", arg.BatchToken)
	}
	if rows, err = d.db.Query(c, fmt.Sprintf(_couponAllowancePageSQL, hitAllowanceInfo(arg.Mid), whereSQL), arg.Mid); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.CouponAllowanceInfo)
		if err = rows.Scan(&r.ID, &r.CouponToken, &r.Mid, &r.State, &r.StartTime, &r.ExpireTime, &r.Origin, &r.Ver, &r.BatchToken, &r.OrderNO, &r.Amount, &r.FullAmount,
			&r.CTime, &r.MTime, &r.Remark); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

//AddViewBatch add view batch.
func (d *Dao) AddViewBatch(c context.Context, arg *model.ArgCouponViewBatch) (err error) {
	if _, err = d.db.Exec(c, _addViewBatchSQL, arg.AppID, arg.Name, arg.BatchToken, arg.MaxCount, arg.CurrentCount, arg.StartTime, arg.ExpireTime, arg.Ver, arg.LimitCount, arg.CouponType, arg.Operator); err != nil {
		err = errors.WithStack(err)
	}
	return
}

//UpdateViewBatch update viewBatch
func (d *Dao) UpdateViewBatch(c context.Context, arg *model.ArgCouponViewBatch) (err error) {
	if _, err = d.db.Exec(c, _updateViewBatchSQL, arg.AppID, arg.Name, arg.MaxCount, arg.LimitCount, arg.Operator, arg.Ver, arg.ID); err != nil {
		err = errors.WithStack(err)
	}
	return
}

//TxUpdateViewInfo update view info.
func (d *Dao) TxUpdateViewInfo(tx *sql.Tx, status int8, couponToken string, mid int64) (err error) {
	if _, err = tx.Exec(fmt.Sprintf(_updateViewSQL, hitViewInfo(mid)), status, couponToken); err != nil {
		err = errors.WithStack(err)
	}
	return
}

//TxCouponViewLog tx add view log.
func (d *Dao) TxCouponViewLog(tx *sql.Tx, arg *model.CouponChangeLog) (err error) {
	if _, err = tx.Exec(fmt.Sprintf(_addViewChangeLog, hitViewInfo(arg.Mid)), arg.CouponToken, arg.Mid, arg.State); err != nil {
		err = errors.WithStack(err)
	}
	return
}

//CouponViewInfo .
func (d *Dao) CouponViewInfo(c context.Context, couponToken string, mid int64) (r *model.CouponInfo, err error) {
	row := d.db.QueryRow(c, fmt.Sprintf(_viewInfoSQL, hitViewInfo(mid)), couponToken)
	r = new(model.CouponInfo)
	if err = row.Scan(&r.CouponToken, &r.Mid, &r.BatchToken, &r.State, &r.OrderNo, &r.OID, &r.StartTime, &r.ExpireTime, &r.Ctime); err != nil {
		err = errors.WithStack(err)
	}
	return
}

//SearchViewCouponCount search view count.
func (d *Dao) SearchViewCouponCount(c context.Context, arg *model.ArgSearchCouponView) (count int64, err error) {
	whereSQL := fmt.Sprintf(_searchViewCountSQL, hitViewInfo(arg.Mid))
	whereSQL += fmt.Sprintf(" AND mid=%v", arg.Mid)
	if len(arg.CouponToken) > 0 {
		whereSQL += fmt.Sprintf(" AND coupon_token='%v'", arg.CouponToken)
	}
	if len(arg.BatchTokens) > 0 {
		whereSQL += fmt.Sprintf(" AND batch_token IN('%v')", strings.Join(arg.BatchTokens, "','"))
	}
	row := d.db.QueryRow(c, whereSQL)
	if err = row.Scan(&count); err != nil {
		err = errors.WithStack(err)
	}
	return
}

//SearchViewCouponInfo search view coupon info .
func (d *Dao) SearchViewCouponInfo(c context.Context, arg *model.ArgSearchCouponView) (res []*model.CouponInfo, err error) {
	var rows *sql.Rows
	whereSQL := fmt.Sprintf(_searchViewSQL, hitViewInfo(arg.Mid))
	whereSQL += fmt.Sprintf(" AND mid=%v", arg.Mid)
	if len(arg.CouponToken) > 0 {
		whereSQL += fmt.Sprintf(" AND coupon_token='%v'", arg.CouponToken)
	}
	if len(arg.BatchTokens) > 0 {
		whereSQL += fmt.Sprintf(" AND batch_token IN('%v')", strings.Join(arg.BatchTokens, "','"))
	}
	if arg.PN <= 0 {
		arg.PN = 1
	}
	if arg.PS <= 0 || arg.PS >= 100 {
		arg.PS = 20
	}
	whereSQL += fmt.Sprintf(" ORDER BY ID DESC LIMIT %v,%v", (arg.PN-1)*arg.PS, arg.PS)
	if rows, err = d.db.Query(c, whereSQL); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.CouponInfo)
		if err = rows.Scan(&r.CouponToken, &r.Mid, &r.BatchToken, &r.State, &r.OrderNo, &r.OID, &r.StartTime, &r.ExpireTime, &r.Ctime, &r.Mtime); err != nil {
			err = errors.WithStack(err)
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

// BatchAddAllowanceCoupon batch add allowance coupon.
func (d *Dao) BatchAddAllowanceCoupon(c context.Context, tx *sql.Tx, cps []*model.CouponAllowanceInfo) (a int64, err error) {
	var (
		buf bytes.Buffer
		res xsql.Result
		sql string
	)
	buf.WriteString(fmt.Sprintf(_batchAddAllowanceCouponSQL, hitAllowanceInfo(cps[0].Mid)))
	for _, v := range cps {
		buf.WriteString("('")
		buf.WriteString(v.CouponToken)
		buf.WriteString("',")
		buf.WriteString(strconv.FormatInt(v.Mid, 10))
		buf.WriteString(",")
		buf.WriteString(fmt.Sprintf("%d", v.State))
		buf.WriteString(",")
		buf.WriteString(strconv.FormatInt(v.StartTime, 10))
		buf.WriteString(",")
		buf.WriteString(strconv.FormatInt(v.ExpireTime, 10))
		buf.WriteString(",")
		buf.WriteString(fmt.Sprintf("%d", v.Origin))
		buf.WriteString(",'")
		buf.WriteString(v.BatchToken)
		buf.WriteString("',")
		buf.WriteString(fmt.Sprintf("%f", v.Amount))
		buf.WriteString(",")
		buf.WriteString(fmt.Sprintf("%f", v.FullAmount))
		buf.WriteString(",'")
		buf.WriteString(v.CTime.Time().Format("2006-01-02 15:04:05"))
		buf.WriteString("',")
		buf.WriteString(strconv.FormatInt(v.AppID, 10))
		buf.WriteString("),")
	}
	sql = buf.String()
	if res, err = tx.Exec(sql[0 : len(sql)-1]); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// UpdateBatchInfo update batch info.
func (d *Dao) UpdateBatchInfo(c context.Context, tx *sql.Tx, token string, count int) (a int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(_updateBatchSQL, count, token); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}
