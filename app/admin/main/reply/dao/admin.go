package dao

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/main/reply/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_adminNameURL = "http://manager.bilibili.co/x/admin/manager/users/unames"

	_selAdminLogByRpIDSQL  = "SELECT id,oid,type,rpid,adminid,result,remark,isnew,isreport,state,ctime,mtime FROM reply_admin_log WHERE rpid=? and isnew=1"
	_selAdminLogsByRpIDSQL = "SELECT id,oid,type,rpid,adminid,result,remark,isnew,isreport,state,ctime,mtime FROM reply_admin_log WHERE rpid=? order by ctime desc"
	_upAdminLogSQL         = "UPDATE reply_admin_log SET isnew=0,mtime=? WHERE rpID IN(%s) AND isnew=1"
	_addAdminLogSQL        = "INSERT INTO reply_admin_log (oid,type,rpid,adminid,result,remark,isnew,isreport,state,ctime,mtime) VALUES %s"
	_addAdminLogFormat     = `(%d,%d,%d,%d,'%s','%s',%d,%d,%d,'%s','%s')`
	_logTimeLayout         = "2006-01-02 15:04:05"
)

// AddAdminLog add admin log to mysql.
func (d *Dao) AddAdminLog(c context.Context, oids, rpIDs []int64, adminID int64, typ, isNew, isReport, state int32, result, remark string, now time.Time) (rows int64, err error) {
	var vals []string
	for i := range oids {
		vals = append(vals, fmt.Sprintf(_addAdminLogFormat, oids[i], typ, rpIDs[i], adminID, result, remark, isNew, isReport, state, now.Format(_logTimeLayout), now.Format(_logTimeLayout)))
	}
	insertSQL := strings.Join(vals, ",")
	res, err := d.db.Exec(c, fmt.Sprintf(_addAdminLogSQL, insertSQL))
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// UpAdminNotNew update admin log to not new.
func (d *Dao) UpAdminNotNew(c context.Context, rpID []int64, now time.Time) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_upAdminLogSQL, xstr.JoinInts(rpID)), now)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// AdminLog return  admin log by rpid
func (d *Dao) AdminLog(c context.Context, rpID int64) (r *model.AdminLog, err error) {
	row := d.db.QueryRow(c, _selAdminLogByRpIDSQL, rpID)
	r = new(model.AdminLog)
	if err = row.Scan(&r.ID, &r.Oid, &r.Type, &r.ReplyID, &r.AdminID, &r.Result, &r.Remark, &r.IsNew, &r.IsReport, &r.State, &r.CTime, &r.MTime); err != nil {
		return
	}
	return
}

// AdminLogsByRpID return operation log list.
func (d *Dao) AdminLogsByRpID(c context.Context, rpID int64) (res []*model.AdminLog, err error) {
	rows, err := d.db.Query(c, _selAdminLogsByRpIDSQL, rpID)
	if err != nil {
		log.Error("query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.AdminLog)
		if err = rows.Scan(&r.ID, &r.Oid, &r.Type, &r.ReplyID, &r.AdminID, &r.Result, &r.Remark, &r.IsNew, &r.IsReport, &r.State, &r.CTime, &r.MTime); err != nil {
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

// AdminName get admin name by id.
func (d *Dao) AdminName(c context.Context, admins map[int64]string) (err error) {
	if len(admins) == 0 {
		return
	}
	var res struct {
		Code int               `json:"code"`
		Data map[string]string `json:"data"`
	}
	var uids []int64
	for key := range admins {
		uids = append(uids, key)
	}
	params := url.Values{}
	params.Set("uids", xstr.JoinInts(uids))
	if err = d.httpClient.Get(c, _adminNameURL, "", params, &res); err != nil {
		log.Error("ReportLog httpClient.Get(%s) error(%v) res(%v)", _adminNameURL, err, res)
		return
	}
	if ec := ecode.Int(res.Code); !ecode.OK.Equal(ec) {
		log.Error("ReportLog not ok(%s) error(%v)  res(%v)", _adminNameURL, ec, res)
		err = ec
		return
	}
	for k, v := range res.Data {
		id, err := strconv.ParseInt(k, 10, 64)
		if err == nil {
			admins[id] = v
		} else {
			log.Error("ReportLog(%s) strconv error(%v) res(%v)", _adminNameURL, err, res)
		}
	}
	return
}
