package dao

import (
	"context"
	"fmt"

	"time"

	"go-common/app/admin/main/laser/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/xstr"
	"net/url"
)

const (
	_UrlUnames                 = "/x/admin/manager/users/unames"
	_UrlUids                   = "/x/admin/manager/users/uids"
	_queryArchiveStatSQL       = " SELECT stat_date, business, stat_type, typeid, uid, stat_value FROM archive_stat WHERE stat_date = '%s' AND business = %d %s "
	_queryArchiveAuditCargoSQL = " SELECT uid, stat_date, receive_value, audit_value FROM archive_audit_cargo_hour %s "
	_queryArchiveStatStreamSQL = " SELECT stat_time, business, stat_type, typeid, uid, stat_value FROM archive_stat_stream WHERE stat_time = '%s' AND business = %d %s "
)

// StatArchiveStat is stat archive data.
func (d *Dao) StatArchiveStat(c context.Context, business int, typeIDS []int64, uids []int64, statTypes []int64, statDate time.Time) (statNodes []*model.StatNode, err error) {
	var queryStmt string
	if len(statTypes) != 0 {
		queryStmt = queryStmt + fmt.Sprintf(" And stat_type in ( %s ) ", xstr.JoinInts(statTypes))
	}
	if len(typeIDS) != 0 {
		queryStmt = queryStmt + fmt.Sprintf(" AND typeid IN ( %s ) ", xstr.JoinInts(typeIDS))
	}
	if len(uids) != 0 {
		queryStmt = queryStmt + fmt.Sprintf(" AND uid IN ( %s ) ", xstr.JoinInts(uids))
	}
	rows, err := d.laserDB.Query(c, fmt.Sprintf(_queryArchiveStatSQL, statDate.Format("2006-01-02"), business, queryStmt))
	if err != nil {
		log.Error("d.laserDB.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		item := &model.StatNode{}
		if err = rows.Scan(&item.StatDate, &item.Business, &item.StatType, &item.TypeID, &item.UID, &item.StatValue); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		statNodes = append(statNodes, item)
	}
	return

}

// QueryArchiveCargo is query archive audit and receive value.
func (d *Dao) QueryArchiveCargo(c context.Context, statTime time.Time, uids []int64) (items []*model.CargoDetail, err error) {
	whereStmt := fmt.Sprintf(" WHERE stat_date = '%s' ", statTime.Format("2006-01-02 15:04:05"))
	if len(uids) != 0 {
		uidStr := xstr.JoinInts(uids)
		whereStmt = whereStmt + fmt.Sprintf(" AND uid in ( %s ) ", uidStr)
	}
	rows, err := d.laserDB.Query(c, fmt.Sprintf(_queryArchiveAuditCargoSQL, whereStmt))
	if err != nil {
		log.Error("d.laserDB.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		item := &model.CargoDetail{}
		if err = rows.Scan(&item.UID, &item.StatDate, &item.ReceiveValue, &item.AuditValue); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		items = append(items, item)
	}
	return
}

//GetUIDByNames is query uids by uname array separated by comma.
func (d *Dao) GetUIDByNames(c context.Context, unames string) (res map[string]int64, err error) {
	var param = url.Values{}
	param.Set("unames", unames)
	var httpRes struct {
		Code    int              `json:"code"`
		Data    map[string]int64 `json:"data"`
		Message string           `json:"message"`
	}

	err = d.HTTPClient.Get(c, d.c.Host.Manager+_UrlUids, "", param, &httpRes)
	if err != nil {
		log.Error("d.client.Get(%s) error(%v)", d.c.Host.Manager+_UrlUids+"?"+param.Encode(), err)
		return
	}
	if httpRes.Code != ecode.OK.Code() {
		log.Error("url(%s) error(%v), code(%d), message(%s)", d.c.Host.Manager+_UrlUids+"?"+param.Encode(), err, httpRes.Code, httpRes.Message)
	}
	res = httpRes.Data
	return
}

//GetUNamesByUids is query usernames by uids.
func (d *Dao) GetUNamesByUids(c context.Context, uids []int64) (res map[int64]string, err error) {
	var param = url.Values{}
	var uidStr = xstr.JoinInts(uids)
	param.Set("uids", uidStr)

	var httpRes struct {
		Code    int              `json:"code"`
		Data    map[int64]string `json:"data"`
		Message string           `json:"message"`
	}

	err = d.HTTPClient.Get(c, d.c.Host.Manager+_UrlUnames, "", param, &httpRes)
	if err != nil {
		log.Error("d.client.Get(%s) error(%v)", d.c.Host.Manager+_UrlUnames+"?"+param.Encode(), err)
		return
	}
	if httpRes.Code != 0 {
		log.Error("url(%s) error(%v), code(%d), message(%s)", d.c.Host.Manager+_UrlUnames+"?"+param.Encode(), err, httpRes.Code, httpRes.Message)
	}
	res = httpRes.Data
	return
}

// StatArchiveStatStream is stat archive data.
func (d *Dao) StatArchiveStatStream(c context.Context, business int, typeIDS []int64, uids []int64, statTypes []int64, statDate time.Time) (statNodes []*model.StatNode, err error) {
	var queryStmt string
	if len(statTypes) != 0 {
		queryStmt = queryStmt + fmt.Sprintf(" And stat_type in ( %s ) ", xstr.JoinInts(statTypes))
	}
	if len(typeIDS) != 0 {
		queryStmt = queryStmt + fmt.Sprintf(" AND typeid IN ( %s ) ", xstr.JoinInts(typeIDS))
	}
	if len(uids) != 0 {
		queryStmt = queryStmt + fmt.Sprintf(" AND uid IN ( %s ) ", xstr.JoinInts(uids))
	}
	rows, err := d.laserDB.Query(c, fmt.Sprintf(_queryArchiveStatStreamSQL, statDate.Format("2006-01-02"), business, queryStmt))
	if err != nil {
		log.Error("d.laserDB.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		node := &model.StatNode{}
		if err = rows.Scan(&node.StatDate, &node.Business, &node.StatType, &node.TypeID, &node.UID, &node.StatValue); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		statNodes = append(statNodes, node)
	}
	return
}
