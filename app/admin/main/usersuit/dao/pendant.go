package dao

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"go-common/app/admin/main/usersuit/model"
	xsql "go-common/library/database/sql"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_inPdGroupSQL        = "INSERT INTO pendant_group(name,rank,status) VALUES (?,?,?)"
	_inPdGroupRefSQL     = "INSERT INTO pendant_group_ref(gid,pid) VALUES (?,?)"
	_inPdInfoSQL         = "INSERT INTO pendant_info(name,image,image_model,status,rank) VALUES (?,?,?,?,?)"
	_inPdPriceSQL        = "INSERT INTO pendant_price(pid,type,price) VALUES (?,?,?) ON DUPLICATE KEY UPDATE pid=?,type=?,price=?"
	_inPdPKGSQL          = "INSERT INTO user_pendant_pkg(mid,pid,expires,type,status) VALUES (?,?,?,4,1) ON DUPLICATE KEY UPDATE expires = ?, type = 4, status = 1"
	_inPdPKGsSQL         = "INSERT INTO user_pendant_pkg(mid,pid,expires,type,status) VALUES %s"
	_inPdEquipSQL        = "INSERT INTO user_pendant_equip(mid,pid,expires) VALUES (?,?,?) ON DUPLICATE KEY UPDATE pid = ? ,expires = ?"
	_inPdOperationLogSQL = "INSERT INTO pendant_operation_log(oper_id,mid,pid,source_type,action) VALUES %s"
	_upPdGroupRefSQL     = "UPDATE pendant_group_ref SET gid = ? WHERE pid = ?"
	_upPdPKGsSQL         = "UPDATE user_pendant_pkg SET expires = CASE id %s END, type = 4, status = 1 WHERE id IN (%s)"
	_upPdGroupSQL        = "UPDATE pendant_group SET name=?,rank=?,status=? WHERE id=?"
	_upPdGroupStatusSQL  = "UPDATE pendant_group SET status=? WHERE id=?"
	_upPdInfoSQL         = "UPDATE pendant_info SET name=?,image=?,image_model=?,status=?,rank=? WHERE id=?"
	_upPdInfoStatusSQL   = "UPDATE pendant_info SET status=? WHERE id=?"
	_pdInfoAllSQL        = `SELECT i.id,i.name,i.image,i.image_model,i.status,i.rank,g.id,g.name,g.rank FROM pendant_info AS i INNER JOIN pendant_group_ref AS r 
	                           ON i.id = r.pid LEFT JOIN pendant_group AS g ON g.id = r.gid ORDER BY g.rank ,i.rank ASC,i.id DESC LIMIT ?,?`
	_pdGroupInfoTotalSQL    = "SELECT COUNT(*) FROM pendant_info INNER JOIN pendant_group_ref ON pendant_info.id = pendant_group_ref.pid"
	_pdGroupsTotalSQL       = "SELECT COUNT(*) FROM pendant_group"
	_pdGroupRefsTotalSQL    = "SELECT COUNT(*) FROM pendant_group_ref"
	_pdGroupRefsGidTotalSQL = "SELECT COUNT(*) FROM pendant_group_ref WHERE gid = ?"
	_pdGroupsSQL            = "SELECT id,name,rank,status FROM pendant_group ORDER BY rank ASC LIMIT ?,?"
	_pdGroupAllSQL          = "SELECT id,name,rank,status FROM pendant_group ORDER BY rank ASC"
	_pdGroupIDsSQL          = "SELECT id,name,rank,status FROM pendant_group WHERE id IN (%s)"
	_pdGroupIDSQL           = "SELECT id,name,rank,status FROM pendant_group WHERE id = ?"
	_pdInfoIDsSQL           = "SELECT id,name,image,image_model,status,rank FROM pendant_info WHERE id IN (%s) ORDER BY rank ASC"
	_pdPriceIDsSQL          = "SELECT pid,type,price FROM pendant_price WHERE pid IN (%s)"
	_pdGroupRefRanksSQL     = "SELECT pr.gid,pr.pid FROM pendant_group_ref AS pr INNER JOIN pendant_group AS pg WHERE pr.gid = pg.id ORDER BY pg.rank ASC LIMIT ?,?"
	_pdGroupRefsSQL         = "SELECT pid FROM pendant_group_ref WHERE gid = ? LIMIT ?,?"
	_pdInfoIDSQL            = "SELECT id,name,image,image_model,status,rank FROM pendant_info WHERE id = ?"
	_pdInfoAllNoPageSQL     = "SELECT id,name,image,image_model,status,rank FROM pendant_info"
	_maxOrderHistorysSQL    = "SELECT MAX(id) FROM user_pendant_order"
	_countOrderHistorysSQL  = "SELECT COUNT(*) FROM user_pendant_order %s"
	_orderHistorysSQL       = "SELECT mid,order_id,pay_id,appid,status,pid,time_length,cost,buy_time,pay_type FROM user_pendant_order %s"
	_pdPKGsUIDSQL           = "SELECT mid,pid,expires,type,status,is_vip FROM user_pendant_pkg WHERE mid = ? ORDER BY mtime DESC"
	_pdPKGUIDsSQL           = "SELECT id,mid,pid,expires,type,status,is_vip FROM user_pendant_pkg WHERE mid IN (%s) AND pid = ?"
	_pdPKGUIDSQL            = "SELECT id,mid,pid,expires,type,status,is_vip FROM user_pendant_pkg WHERE mid = ? AND pid = ?"
	_pdEquipUIDSQL          = "SELECT pid,expires FROM user_pendant_equip WHERE mid = ? AND expires >= ?"
	_pdOperationLogTotalSQL = "SELECT MAX(id) FROM pendant_operation_log"
	_pdOperationLogSQL      = "SELECT oper_id,action,mid,pid,source_type,ctime,mtime FROM pendant_operation_log ORDER BY mtime DESC LIMIT ?,?"
)

// AddPendantGroup insert pendant group .
func (d *Dao) AddPendantGroup(c context.Context, pg *model.PendantGroup) (gid int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _inPdGroupSQL, pg.Name, pg.Rank, pg.Status); err != nil {
		err = errors.Wrapf(err, "AddPendantGroup d.db.Exec(%s,%d)", pg.Name, pg.Rank)
		return
	}
	return res.LastInsertId()
}

// TxAddPendantGroupRef tx insert pendant group ref.
func (d *Dao) TxAddPendantGroupRef(tx *xsql.Tx, pr *model.PendantGroupRef) (affected int64, err error) {
	var res sql.Result
	if res, err = tx.Exec(_inPdGroupRefSQL, pr.GID, pr.PID); err != nil {
		err = errors.Wrapf(err, "TxAddPendantGroupRef itx.Exec(%d,%d)", pr.GID, pr.PID)
		return
	}
	return res.RowsAffected()
}

// TxAddPendantInfo insert pendant info.
func (d *Dao) TxAddPendantInfo(tx *xsql.Tx, pi *model.PendantInfo) (pid int64, err error) {
	var res sql.Result
	if res, err = tx.Exec(_inPdInfoSQL, pi.Name, pi.Image, pi.ImageModel, pi.Status, pi.Rank); err != nil {
		err = errors.Wrapf(err, "TxAddPendantInfo tx.Exec(%s,%s,%s,%d)", pi.Name, pi.Image, pi.ImageModel, pi.Rank)
		return
	}
	return res.LastInsertId()
}

// TxAddPendantPrices insert pendant prices.
func (d *Dao) TxAddPendantPrices(tx *xsql.Tx, pp *model.PendantPrice) (affected int64, err error) {
	var res sql.Result
	if res, err = tx.Exec(_inPdPriceSQL, pp.PID, pp.TP, pp.Price, pp.PID, pp.TP, pp.Price); err != nil {
		err = errors.Wrapf(err, "TxAddPendantPrices tx.Exec(%d,%d,%d)", pp.PID, pp.TP, pp.Price)
		return
	}
	return res.RowsAffected()
}

// AddPendantPKG  insert pendant pkg.
func (d *Dao) AddPendantPKG(c context.Context, pkg *model.PendantPKG) (affected int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _inPdPKGSQL, pkg.UID, pkg.PID, pkg.Expires, pkg.Expires); err != nil {
		err = errors.Wrapf(err, "AddPendantPKG d.db.Exec(%d,%d,%d)", pkg.UID, pkg.PID, pkg.Expires)
		return
	}
	return res.RowsAffected()
}

// TxAddPendantPKGs multi insert pendant pkg.
func (d *Dao) TxAddPendantPKGs(tx *xsql.Tx, pkgs []*model.PendantPKG) (affected int64, err error) {
	var (
		uids []int64
		pids map[int64]struct{}
	)
	l := len(pkgs)
	valueStrings := make([]string, 0, l)
	valueArgs := make([]interface{}, 0, l*3)
	pids = make(map[int64]struct{})
	for _, pkg := range pkgs {
		valueStrings = append(valueStrings, "(?,?,?,4,1)")
		valueArgs = append(valueArgs, strconv.FormatInt(pkg.UID, 10))
		valueArgs = append(valueArgs, strconv.FormatInt(pkg.PID, 10))
		valueArgs = append(valueArgs, strconv.FormatInt(pkg.Expires, 10))
		uids = append(uids, pkg.UID)
		pids[pkg.PID] = struct{}{}
	}
	stmt := fmt.Sprintf(_inPdPKGsSQL, strings.Join(valueStrings, ","))
	_, err = tx.Exec(stmt, valueArgs...)
	if err != nil {
		err = errors.Wrapf(err, "TxAddPendantPKGs tx.Exec(%s,%+v)", xstr.JoinInts(uids), reflect.ValueOf(pids).MapKeys())
	}
	return
}

// AddPendantEquip  insert pendant equip.
func (d *Dao) AddPendantEquip(c context.Context, pkg *model.PendantPKG) (affected int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _inPdEquipSQL, pkg.UID, pkg.PID, pkg.Expires, pkg.PID, pkg.Expires); err != nil {
		err = errors.Wrapf(err, "AddPendantEquip d.db.Exec(%d,%d,%d)", pkg.UID, pkg.PID, pkg.Expires)
		return
	}
	return res.RowsAffected()
}

// AddPendantOperLog insert pendant operation log.
func (d *Dao) AddPendantOperLog(c context.Context, oid int64, uids []int64, pid int64, action string) (affected int64, err error) {
	var res sql.Result
	l := len(uids)
	valueStrings := make([]string, 0, l)
	valueArgs := make([]interface{}, 0, l*5)
	for _, uid := range uids {
		valueStrings = append(valueStrings, "(?,?,?,?,?)")
		valueArgs = append(valueArgs, strconv.FormatInt(oid, 10))
		valueArgs = append(valueArgs, strconv.FormatInt(uid, 10))
		valueArgs = append(valueArgs, strconv.FormatInt(pid, 10))
		valueArgs = append(valueArgs, strconv.FormatInt(int64(model.PendantSourceTypeAdmin), 10))
		valueArgs = append(valueArgs, action)
	}
	stmt := fmt.Sprintf(_inPdOperationLogSQL, strings.Join(valueStrings, ","))
	res, err = d.db.Exec(c, stmt, valueArgs...)
	if err != nil {
		err = errors.Errorf("AddPendantOperLog tx.Exec(%s,%d,%s) error(%+v)", xstr.JoinInts(uids), pid, action, err)
		return
	}
	return res.RowsAffected()
}

// TxUpPendantGroupRef update pendant group ref.
func (d *Dao) TxUpPendantGroupRef(tx *xsql.Tx, gid, pid int64) (affected int64, err error) {
	var res sql.Result
	if res, err = tx.Exec(_upPdGroupRefSQL, gid, pid); err != nil {
		err = errors.Wrapf(err, "UpPendantGroupRef tx.Exec(%d,%d)", gid, pid)
		return
	}
	return res.RowsAffected()
}

// TxUpPendantPKGs multi update pendant pkg.
func (d *Dao) TxUpPendantPKGs(tx *xsql.Tx, pkgs []*model.PendantPKG) (affected int64, err error) {
	var ids []int64
	l := len(pkgs)
	valueStrings := make([]string, 0, l)
	valueArgs := make([]interface{}, 0, l*2)
	for _, pkg := range pkgs {
		valueStrings = append(valueStrings, "WHEN ? THEN ? ")
		valueArgs = append(valueArgs, strconv.FormatInt(pkg.ID, 10))
		valueArgs = append(valueArgs, strconv.FormatInt(pkg.Expires, 10))
		ids = append(ids, pkg.ID)
	}
	stmt := fmt.Sprintf(_upPdPKGsSQL, strings.Join(valueStrings, " "), xstr.JoinInts(ids))
	_, err = tx.Exec(stmt, valueArgs...)
	if err != nil {
		err = errors.Wrapf(err, "TxUpPendantPKGs tx.Exec(%s)", xstr.JoinInts(ids))
	}
	return
}

// UpPendantGroup update pendant group.
func (d *Dao) UpPendantGroup(c context.Context, pg *model.PendantGroup) (affected int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _upPdGroupSQL, pg.Name, pg.Rank, pg.Status, pg.ID); err != nil {
		err = errors.Wrapf(err, "UpPendantGroup tx.Exec(%s,%d,%d,%d)", pg.Name, pg.Rank, pg.Status, pg.ID)
		return
	}
	return res.RowsAffected()
}

// UpPendantGroupStatus update pendant group status.
func (d *Dao) UpPendantGroupStatus(c context.Context, gid int64, status int8) (affected int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _upPdGroupStatusSQL, status, gid); err != nil {
		err = errors.Wrapf(err, "UpPendantGroupStatus tx.Exec(%d,%d)", status, gid)
		return
	}
	return res.RowsAffected()
}

// TxUpPendantInfo update pendant info.
func (d *Dao) TxUpPendantInfo(tx *xsql.Tx, pi *model.PendantInfo) (affected int64, err error) {
	var res sql.Result
	if res, err = tx.Exec(_upPdInfoSQL, pi.Name, pi.Image, pi.ImageModel, pi.Status, pi.Rank, pi.ID); err != nil {
		err = errors.Wrapf(err, "TxAddPendantPrices tx.Exec(%s,%s,%s,%d,%d,%d)", pi.Name, pi.Image, pi.ImageModel, pi.Status, pi.Rank, pi.ID)
		return
	}
	return res.RowsAffected()
}

// UpPendantInfoStatus update pendant info status.
func (d *Dao) UpPendantInfoStatus(c context.Context, pid int64, status int8) (affected int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _upPdInfoStatusSQL, status, pid); err != nil {
		err = errors.Wrapf(err, "UpPendantGroupStatus tx.Exec(%d,%d)", status, pid)
		return
	}
	return res.RowsAffected()
}

// PendantInfoAll pendant info all.
func (d *Dao) PendantInfoAll(c context.Context, pn, ps int) (pis []*model.PendantInfo, pids []int64, err error) {
	var (
		gid, groupRank sql.NullInt64
		groupName      sql.NullString
		rows           *xsql.Rows
		offset         = (pn - 1) * ps
	)
	if rows, err = d.db.Query(c, _pdInfoAllSQL, offset, ps); err != nil {
		err = errors.Wrapf(err, "PendantInfoAll d.db.Query(%d,%d)", offset, ps)
		return
	}
	defer rows.Close()
	for rows.Next() {
		pi := new(model.PendantInfo)
		if err = rows.Scan(&pi.ID, &pi.Name, &pi.Image, &pi.ImageModel, &pi.Status, &pi.Rank, &gid, &groupName, &groupRank); err != nil {
			err = errors.Wrap(err, "PendantInfoAll row.Scan()")
			return
		}
		pi.GID = gid.Int64
		pi.GroupName = groupName.String
		pi.GroupRank = int16(groupRank.Int64)
		pids = append(pids, pi.ID)
		pis = append(pis, pi)
	}
	err = rows.Err()
	return
}

// PendantGroupInfoTotal pendant group info total.
func (d *Dao) PendantGroupInfoTotal(c context.Context) (count int64, err error) {
	row := d.db.QueryRow(c, _pdGroupInfoTotalSQL)
	if err = row.Scan(&count); err != nil {
		err = errors.Wrap(err, "d.dao.PendantGroupInfoTotal")
	}
	return
}

// PendantGroupsTotal pendant group total.
func (d *Dao) PendantGroupsTotal(c context.Context) (count int64, err error) {
	row := d.db.QueryRow(c, _pdGroupsTotalSQL)
	if err = row.Scan(&count); err != nil {
		err = errors.Wrap(err, "d.dao.PendantGroupsTotal")
	}
	return
}

// PendantGroupRefsTotal pendant group refs total.
func (d *Dao) PendantGroupRefsTotal(c context.Context) (count int64, err error) {
	row := d.db.QueryRow(c, _pdGroupRefsTotalSQL)
	if err = row.Scan(&count); err != nil {
		err = errors.Wrap(err, "d.dao.PendantGroupRefsTotal")
	}
	return
}

// PendantGroupRefsGidTotal pendant group refs total by gid.
func (d *Dao) PendantGroupRefsGidTotal(c context.Context, gid int64) (count int64, err error) {
	row := d.db.QueryRow(c, _pdGroupRefsGidTotalSQL, gid)
	if err = row.Scan(&count); err != nil {
		err = errors.Wrap(err, "d.dao.PendantGroupRefsGidTotal")
	}
	return
}

// PendantGroups pendant group pagesize.
func (d *Dao) PendantGroups(c context.Context, pn, ps int) (pgs []*model.PendantGroup, err error) {
	var (
		rows   *xsql.Rows
		offset = (pn - 1) * ps
	)
	if rows, err = d.db.Query(c, _pdGroupsSQL, offset, ps); err != nil {
		err = errors.Wrapf(err, "PendantGroups d.db.Query(%d,%d)", offset, ps)
		return
	}
	defer rows.Close()
	for rows.Next() {
		pg := new(model.PendantGroup)
		if err = rows.Scan(&pg.ID, &pg.Name, &pg.Rank, &pg.Status); err != nil {
			err = errors.Wrap(err, "PendantGroups row.Scan()")
			return
		}
		pgs = append(pgs, pg)
	}
	err = rows.Err()
	return
}

// PendantGroupAll pendant all group .
func (d *Dao) PendantGroupAll(c context.Context) (pgs []*model.PendantGroup, err error) {
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, _pdGroupAllSQL); err != nil {
		err = errors.Wrap(err, "PendantGroupAll d.db.Query(%d,%d)")
		return
	}
	defer rows.Close()
	for rows.Next() {
		pg := new(model.PendantGroup)
		if err = rows.Scan(&pg.ID, &pg.Name, &pg.Rank, &pg.Status); err != nil {
			err = errors.Wrap(err, "PendantGroupAll row.Scan()")
			return
		}
		pgs = append(pgs, pg)
	}
	err = rows.Err()
	return
}

// PendantGroupIDs pendant group in ids.
func (d *Dao) PendantGroupIDs(c context.Context, ids []int64) (pgs []*model.PendantGroup, err error) {
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, fmt.Sprintf(_pdGroupIDsSQL, xstr.JoinInts(ids))); err != nil {
		err = errors.Wrapf(err, "PendantGroupIDs d.db.Query(%s)", xstr.JoinInts(ids))
		return
	}
	defer rows.Close()
	for rows.Next() {
		pg := new(model.PendantGroup)
		if err = rows.Scan(&pg.ID, &pg.Name, &pg.Rank, &pg.Status); err != nil {
			err = errors.Wrap(err, "PendantGroupIDs row.Scan()")
			return
		}
		pgs = append(pgs, pg)
	}
	err = rows.Err()
	return
}

// PendantGroupID pendant group by id.
func (d *Dao) PendantGroupID(c context.Context, id int64) (pg *model.PendantGroup, err error) {
	row := d.db.QueryRow(c, _pdGroupIDSQL, id)
	if err != nil {
		err = errors.Wrapf(err, "PendantGroupID d.db.Query(%d)", id)
		return
	}
	pg = &model.PendantGroup{}
	if err = row.Scan(&pg.ID, &pg.Name, &pg.Rank, &pg.Status); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			pg = nil
			return
		}
		err = errors.Wrap(err, "PendantGroupID row.Scan")
	}
	return
}

// PendantInfoIDs pendant info in ids.
func (d *Dao) PendantInfoIDs(c context.Context, ids []int64) (pis []*model.PendantInfo, pim map[int64]*model.PendantInfo, err error) {
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, fmt.Sprintf(_pdInfoIDsSQL, xstr.JoinInts(ids))); err != nil {
		err = errors.Wrapf(err, "PendantInfoIDs d.db.Query(%s)", xstr.JoinInts(ids))
		return
	}
	defer rows.Close()
	pim = make(map[int64]*model.PendantInfo, len(ids))
	for rows.Next() {
		pi := new(model.PendantInfo)
		if err = rows.Scan(&pi.ID, &pi.Name, &pi.Image, &pi.ImageModel, &pi.Status, &pi.Rank); err != nil {
			err = errors.Wrap(err, "PendantInfoIDs row.Scan()")
			return
		}
		pis = append(pis, pi)
		pim[pi.ID] = pi
	}
	err = rows.Err()
	return
}

// PendantPriceIDs pendant price in ids.
func (d *Dao) PendantPriceIDs(c context.Context, ids []int64) (ppm map[int64][]*model.PendantPrice, err error) {
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, fmt.Sprintf(_pdPriceIDsSQL, xstr.JoinInts(ids))); err != nil {
		err = errors.Wrapf(err, "PendantPriceIDs d.db.Query(%s)", xstr.JoinInts(ids))
		return
	}
	defer rows.Close()
	ppm = make(map[int64][]*model.PendantPrice, len(ids))
	for rows.Next() {
		pp := new(model.PendantPrice)
		if err = rows.Scan(&pp.PID, &pp.TP, &pp.Price); err != nil {
			err = errors.Wrap(err, "PendantPriceIDs row.Scan()")
			return
		}
		ppm[pp.PID] = append(ppm[pp.PID], pp)
	}
	err = rows.Err()
	return
}

// PendantGroupRefRanks pendant group ref pagesize by rank.
func (d *Dao) PendantGroupRefRanks(c context.Context, pn, ps int) (prs []*model.PendantGroupRef, err error) {
	var (
		rows   *xsql.Rows
		offset = (pn - 1) * ps
	)
	if rows, err = d.db.Query(c, _pdGroupRefRanksSQL, offset, ps); err != nil {
		err = errors.Wrapf(err, "PendantGroupRefRanks d.db.Query(%d,%d)", offset, ps)
		return
	}
	defer rows.Close()
	for rows.Next() {
		pr := new(model.PendantGroupRef)
		if err = rows.Scan(&pr.GID, &pr.PID); err != nil {
			err = errors.Wrap(err, "PendantGroupRefRanks row.Scan()")
			return
		}
		prs = append(prs, pr)
	}
	err = rows.Err()
	return
}

// PendantGroupPIDs pendant group ref pagesize.
func (d *Dao) PendantGroupPIDs(c context.Context, gid int64, pn, ps int) (pids []int64, err error) {
	var (
		rows   *xsql.Rows
		offset = (pn - 1) * ps
	)
	if rows, err = d.db.Query(c, _pdGroupRefsSQL, gid, offset, ps); err != nil {
		err = errors.Wrapf(err, "PendantGroupPIDs d.db.Query(%d,%d,%d)", gid, offset, ps)
		return
	}
	defer rows.Close()
	var pid int64
	for rows.Next() {
		if err = rows.Scan(&pid); err != nil {
			err = errors.Wrap(err, "PendantGroupPIDs row.Scan()")
			return
		}
		pids = append(pids, pid)
	}
	err = rows.Err()
	return
}

// PendantInfoID pendant info.
func (d *Dao) PendantInfoID(c context.Context, id int64) (pi *model.PendantInfo, err error) {
	row := d.db.QueryRow(c, _pdInfoIDSQL, id)
	if err != nil {
		err = errors.Wrapf(err, "PendantInfoID d.db.QueryRow(%d)", id)
		return
	}
	pi = &model.PendantInfo{}
	if err = row.Scan(&pi.ID, &pi.Name, &pi.Image, &pi.ImageModel, &pi.Status, &pi.Rank); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			pi = nil
			return
		}
		err = errors.Wrap(err, "PendantInfoID row.Scan")
	}
	return
}

// PendantInfoAllNoPage pendant info no page.
func (d *Dao) PendantInfoAllNoPage(c context.Context) (pis []*model.PendantInfo, err error) {
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, _pdInfoAllNoPageSQL); err != nil {
		err = errors.Wrap(err, "PendantInfoAllOnSale d.db.Query()")
		return
	}
	defer rows.Close()
	for rows.Next() {
		pi := new(model.PendantInfo)
		if err = rows.Scan(&pi.ID, &pi.Name, &pi.Image, &pi.ImageModel, &pi.Status, &pi.Rank); err != nil {
			err = errors.Wrap(err, "PendantInfoAllOnSale row.Scan()")
			return
		}
		pis = append(pis, pi)
	}
	err = rows.Err()
	return
}

// BuildOrderInfoSQL build a order sql string.
func (d *Dao) BuildOrderInfoSQL(c context.Context, arg *model.ArgPendantOrder, tp string) (sql string, values []interface{}) {
	values = make([]interface{}, 0, 5)
	var (
		cond    []string
		condStr string
	)
	if arg.UID != 0 {
		cond = append(cond, "mid = ?")
		values = append(values, arg.UID)
	}
	if arg.PID != 0 {
		cond = append(cond, "pid = ?")
		values = append(values, arg.PID)
	}
	if arg.Status != 0 {
		cond = append(cond, "status = ?")
		values = append(values, arg.Status)
	}
	if arg.PayID != "" {
		cond = append(cond, "pay_id = ?")
		values = append(values, arg.PayID)
	}
	if arg.Start != 0 {
		cond = append(cond, "mtime >= ?")
		values = append(values, arg.Start)
	}
	if arg.End != 0 {
		cond = append(cond, "mtime <= ?")
		values = append(values, arg.End)
	}
	if tp == "info" {
		condStr = d.joinStrings(cond)
		if condStr != "" {
			sql = fmt.Sprintf(_orderHistorysSQL+" %s %s ", "WHERE", condStr, "ORDER BY mtime DESC LIMIT ?,?")
		} else {
			sql = fmt.Sprintf(_orderHistorysSQL+" %s ", condStr, "ORDER BY mtime DESC LIMIT ?,?")
		}
		values = append(values, (arg.PN-1)*arg.PS, arg.PS)
	} else if tp == "count" {
		condStr = d.joinStrings(cond)
		if condStr != "" {
			sql = fmt.Sprintf(_countOrderHistorysSQL+" %s", "WHERE", condStr)
		} else {
			sql = fmt.Sprintf(_countOrderHistorysSQL, condStr)
		}
	}
	return
}

func (d *Dao) joinStrings(is []string) string {
	if len(is) == 0 {
		return ""
	}
	if len(is) == 1 {
		return is[0]
	}
	var bfPool = sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer([]byte{})
		},
	}
	buf := bfPool.Get().(*bytes.Buffer)
	for _, i := range is {
		buf.WriteString(i)
		buf.WriteString(" AND ")
	}
	if buf.Len() > 0 {
		buf.Truncate(buf.Len() - 4)
	}
	s := buf.String()
	buf.Reset()
	bfPool.Put(buf)
	return s
}

// MaxOrderHistory max order history.
func (d *Dao) MaxOrderHistory(c context.Context) (max int64, err error) {
	row := d.db.QueryRow(c, _maxOrderHistorysSQL)
	if err != nil {
		err = errors.Wrap(err, "MaxOrderHistory d.db.QueryRow()")
		return
	}
	if err = row.Scan(&max); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			max = 0
			return
		}
		err = errors.Wrap(err, "MaxOrderHistory row.Scan")
	}
	return
}

// CountOrderHistory count order history.
func (d *Dao) CountOrderHistory(c context.Context, arg *model.ArgPendantOrder) (total int64, err error) {
	sqlstr, values := d.BuildOrderInfoSQL(c, arg, "count")
	row := d.db.QueryRow(c, sqlstr, values...)
	if err != nil {
		err = errors.Wrap(err, "CountOrderHistory d.db.QueryRow()")
		return
	}
	if err = row.Scan(&total); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			total = 0
			return
		}
		err = errors.Wrap(err, "CountOrderHistory row.Scan")
	}
	return
}

// OrderHistorys get order historys.
func (d *Dao) OrderHistorys(c context.Context, arg *model.ArgPendantOrder) (pos []*model.PendantOrder, pids []int64, err error) {
	var rows *xsql.Rows
	sqlstr, values := d.BuildOrderInfoSQL(c, arg, "info")
	if rows, err = d.db.Query(c, sqlstr, values...); err != nil {
		err = errors.Wrap(err, "OrderHistorys d.db.Query()")
		return
	}
	defer rows.Close()
	for rows.Next() {
		po := new(model.PendantOrder)
		if err = rows.Scan(&po.UID, &po.OrderID, &po.PayID, &po.AppID, &po.Status, &po.PID, &po.TimeLength, &po.Cost, &po.BuyTime, &po.PayType); err != nil {
			err = errors.Wrap(err, "OrderHistorys row.Scan()")
			return
		}
		pos = append(pos, po)
		pids = append(pids, po.PID)
	}
	err = rows.Err()
	return
}

// PendantPKGs get pendant pkgs.
func (d *Dao) PendantPKGs(c context.Context, uid int64) (pkgs []*model.PendantPKG, err error) {
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, _pdPKGsUIDSQL, uid); err != nil {
		err = errors.Wrapf(err, "PendantPKGs d.db.Query(%d)", uid)
		return
	}
	defer rows.Close()
	for rows.Next() {
		pkg := new(model.PendantPKG)
		if err = rows.Scan(&pkg.UID, &pkg.PID, &pkg.Expires, &pkg.TP, &pkg.Status, &pkg.IsVip); err != nil {
			err = errors.Wrap(err, "PendantPKGs row.Scan()")
			return
		}
		pkgs = append(pkgs, pkg)
	}
	err = rows.Err()
	return
}

// PendantPKGUIDs get pendant pkgs by muilti uid.
func (d *Dao) PendantPKGUIDs(c context.Context, uids []int64, pid int64) (pkgs []*model.PendantPKG, err error) {
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, fmt.Sprintf(_pdPKGUIDsSQL, xstr.JoinInts(uids)), pid); err != nil {
		err = errors.Wrapf(err, "PendantPKGUIDs d.db.Query(%s,%d)", xstr.JoinInts(uids), pid)
		return
	}
	defer rows.Close()
	for rows.Next() {
		pkg := new(model.PendantPKG)
		if err = rows.Scan(&pkg.ID, &pkg.UID, &pkg.PID, &pkg.Expires, &pkg.TP, &pkg.Status, &pkg.IsVip); err != nil {
			err = errors.Wrap(err, "PendantPKGUIDs row.Scan()")
			return
		}
		pkgs = append(pkgs, pkg)
	}
	err = rows.Err()
	return
}

// PendantPKG get  pendant in pkg.
func (d *Dao) PendantPKG(c context.Context, uid, pid int64) (pkg *model.PendantPKG, err error) {
	row := d.db.QueryRow(c, _pdPKGUIDSQL, uid, pid)
	if err != nil {
		err = errors.Wrapf(err, "PendantPKG d.db.QueryRow(%d,%d)", uid, pid)
		return
	}
	pkg = &model.PendantPKG{}
	if err = row.Scan(&pkg.ID, &pkg.UID, &pkg.PID, &pkg.Expires, &pkg.TP, &pkg.Status, &pkg.IsVip); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			pkg = nil
			return
		}
		err = errors.Wrap(err, "PendantPKG row.Scan")
	}
	return
}

// PendantEquipUID  pendant equip by uid.
func (d *Dao) PendantEquipUID(c context.Context, uid int64) (pkg *model.PendantPKG, err error) {
	row := d.db.QueryRow(c, _pdEquipUIDSQL, uid, time.Now().Unix())
	if err != nil {
		err = errors.Wrapf(err, "PendantEquipUID d.db.QueryRow(%d)", uid)
		return
	}
	pkg = &model.PendantPKG{}
	if err = row.Scan(&pkg.PID, &pkg.Expires); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			pkg = nil
			return
		}
		err = errors.Wrap(err, "PendantEquipUID row.Scan")
	}
	return
}

// PendantOperLog get pendant operation log.
func (d *Dao) PendantOperLog(c context.Context, pn, ps int) (opers []*model.PendantOperLog, uids []int64, err error) {
	var (
		rows   *xsql.Rows
		offset = (pn - 1) * ps
	)
	if rows, err = d.db.Query(c, _pdOperationLogSQL, offset, ps); err != nil {
		err = errors.Wrapf(err, "PendantOperLog d.db.Query(%d,%d)", offset, ps)
		return
	}
	defer rows.Close()
	for rows.Next() {
		oper := new(model.PendantOperLog)
		if err = rows.Scan(&oper.OID, &oper.Action, &oper.UID, &oper.PID, &oper.SourceType, &oper.CTime, &oper.MTime); err != nil {
			err = errors.Wrap(err, "PendantOperLog row.Scan()")
			return
		}
		opers = append(opers, oper)
		uids = append(uids, oper.UID)
	}
	err = rows.Err()
	return
}

// PendantOperationLogTotal pendant operation log  total.
func (d *Dao) PendantOperationLogTotal(c context.Context) (count int64, err error) {
	row := d.db.QueryRow(c, _pdOperationLogTotalSQL)
	if err = row.Scan(&count); err != nil {
		err = errors.Wrap(err, "d.dao.PendantOperationLogTotal")
	}
	return
}
