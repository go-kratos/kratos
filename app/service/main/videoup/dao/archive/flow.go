package archive

import (
	"context"
	"fmt"
	"time"

	"go-common/app/service/main/videoup/model/archive"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_upStateFlowSQL     = "UPDATE flow_design SET state =? where id=?"
	_inFlowSQL          = "INSERT into flow_design(pool,oid,group_id,uid,remark)  VALUES (?,?,?,?,?)"
	_inFlowLogSQL       = "INSERT into flow_design_log(pool,oid,group_id,uid,action,remark)  VALUES (?,?,?,?,1,?)"
	_flowsSQL           = "SELECT id,rank,type,value,name,ctime FROM flow_group WHERE  state=0 order by rank desc"
	_whiteMidSQL        = "SELECT oid,uid FROM flow_design WHERE pool=1 AND state=0 AND group_id=11"
	_isFlowGroupIDInSQL = "SELECT id FROM flow_design WHERE pool=? AND state=0 AND group_id=? AND oid=? limit 1"
	//稿件 mid 配置
	_midsForbidSQL         = "SELECT flow_design.id,oid,value FROM flow_design left join flow_group on flow_design.group_id=flow_group.id WHERE flow_design.pool=1 AND flow_design.state=0 AND flow_design.parent=0 AND flow_design.group_id>=12 and flow_group.type=4 "
	_isActGroupIDSQL       = "SELECT state FROM flow_group WHERE id=?"
	_isMidIDSQL            = "SELECT id,pool,oid,group_id,parent,uid,remark,ctime,mtime FROM flow_design WHERE pool=1 AND state=0  AND oid=? AND parent=?"
	_findGroupIDByScopeSQL = "SELECT group_id FROM flow_scope WHERE  pool= ? AND industry_id=? AND brand_id=? AND official=? AND state=0  order by id desc limit 1;"
	_appFlowsSQL           = "SELECT oid FROM flow_design WHERE pool=0 AND state=0 AND mtime>=? AND mtime<=? AND group_id=11"
	_flowGroupPool         = "SELECT id, pool FROM flow_group WHERE id IN (%s)"
	_flowsByOIDSQL         = "SELECT fd.id,fd.pool,fd.oid,fd.group_id,fd.parent,fd.state,fg.value FROM flow_design fd LEFT JOIN flow_group fg ON fd.group_id=fg.id WHERE fd.oid=? AND fd.state=0 AND fg.state=0"
	_flowsByGIDSQL         = "SELECT fd.id,fd.pool,fd.oid,fd.group_id,fd.parent,fd.state,fg.value FROM flow_design fd LEFT JOIN flow_group fg ON fd.group_id=fg.id WHERE fd.pool=? AND fd.group_id=? AND fd.state=0 AND fg.state=0 LIMIT ?,?"
	_flowUniqueSQL         = "SELECT id,pool,oid,group_id,parent,state FROM flow_design WHERE oid=? AND pool=? AND group_id=? LIMIT 1"
	_flowCountSQL          = "SELECT count(*) FROM flow_design fd LEFT JOIN flow_group fg ON fd.group_id=fg.id WHERE fd.pool=? AND fd.group_id=? AND fd.state=0 AND fg.state=0 "
	_flowOidsByGidSQL      = "SELECT fd.id,fd.pool,fd.oid,fd.group_id,fd.parent,fd.state,fg.value FROM flow_design fd LEFT JOIN flow_group fg ON fd.group_id=fg.id WHERE fd.pool=? AND fd.group_id=? AND fd.state=0 AND fg.state=0 AND fd.oid IN (%s) "
)

// TxAddFlow tx add flow_design.
func (d *Dao) TxAddFlow(tx *sql.Tx, old, uid, groupID int64, pool int8, remark string) (id int64, err error) {
	res, err := tx.Exec(_inFlowSQL, pool, old, groupID, uid, remark)
	if err != nil {
		log.Error("d._inFlow.Exec() error(%v)", err)
		return
	}
	id, err = res.LastInsertId()
	return
}

// TxUpFlowState tx set flow_design.state=1.
func (d *Dao) TxUpFlowState(tx *sql.Tx, state int8, id int64) (rows int64, err error) {
	res, err := tx.Exec(_upStateFlowSQL, state, id)
	if err != nil {
		log.Error("d.TxUpFlowState.Exec() error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxAddFlowLog tx add flow_design log.
func (d *Dao) TxAddFlowLog(tx *sql.Tx, old, uid, groupID int64, pool int8, remark string) (rows int64, err error) {
	res, err := tx.Exec(_inFlowLogSQL, pool, old, groupID, uid, remark)
	if err != nil {
		log.Error("d._inFlowLog.Exec() error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// Flows get flow_control id and remark.
func (d *Dao) Flows(c context.Context) (fs []*archive.Flow, err error) {
	rows, err := d.db.Query(c, _flowsSQL)
	if err != nil {
		log.Error("d.db.Query(%s) error(%v)", _flowsSQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		f := &archive.Flow{}
		if err = rows.Scan(&f.ID, &f.Rank, &f.Type, &f.Value, &f.Remark, &f.CTime); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		fs = append(fs, f)
	}
	return
}

// WhiteMids get white mids.
func (d *Dao) WhiteMids(c context.Context) (mids map[int64]int64, err error) {
	rows, err := d.db.Query(c, _whiteMidSQL)
	if err != nil {
		log.Error("d.db.Query(%s) error(%v)", _whiteMidSQL, err)
		return
	}
	defer rows.Close()
	mids = make(map[int64]int64)
	for rows.Next() {
		var mid, uid int64
		if err = rows.Scan(&mid, &uid); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		mids[mid] = uid
	}
	return
}

// CheckActGroupID check active GroupID
func (d *Dao) CheckActGroupID(c context.Context, groupID int64) (state int8, err error) {
	row := d.db.QueryRow(c, _isActGroupIDSQL, groupID)
	if err = row.Scan(&state); err != nil {
		if err != sql.ErrNoRows {
			log.Error("row.Scan error(%v)", err)
			return
		}
		err = nil
		state = archive.FlowDelete
	}
	return
}

// CheckFlowGroupID check active GroupID
func (d *Dao) CheckFlowGroupID(c context.Context, pool int8, oid, groupID int64) (flowID int64, err error) {
	row := d.db.QueryRow(c, _isFlowGroupIDInSQL, pool, groupID, oid)
	if err = row.Scan(&flowID); err != nil {
		if err != sql.ErrNoRows {
			log.Error("row.Scan error(%v)", err)
			return
		}
		err = nil
	}
	return
}

//CheckFlowMid  通用 mid 流量导向到flow_design
func (d *Dao) CheckFlowMid(c context.Context, pool int8, oid int64) (flows []*archive.FlowData, err error) {
	rows, err := d.db.Query(c, _isMidIDSQL, oid, pool)
	if err != nil {
		log.Error("d.db.Query (%s)|(%d)|(%d) error(%v)", _isMidIDSQL, oid, pool, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		item := &archive.FlowData{}
		if err = rows.Scan(&item.ID, &item.Pool, &item.OID, &item.GroupID, &item.Parent, &item.UID, &item.Remark, &item.CTime, &item.MTime); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		flows = append(flows, item)
	}
	log.Info("flowDesign mids design (%+v)", flows)
	return
}

// FindGroupIDByScope check active GroupID
func (d *Dao) FindGroupIDByScope(c context.Context, pool int8, IndustryID, brandID int64, official int8) (groupID int64, err error) {
	row := d.db.QueryRow(c, _findGroupIDByScopeSQL, pool, IndustryID, brandID, official)
	if err = row.Scan(&groupID); err != nil {
		if err != sql.ErrNoRows {
			log.Error("row.Scan error(%v)", err)
			return
		}
		err = nil
		groupID = 1
	}
	return
}

//ForbidMids 稿件 mid 禁止配置 ForbidMids get forbid mids.
func (d *Dao) ForbidMids(c context.Context) (mids map[int64][]string, err error) {
	rows, err := d.db.Query(c, _midsForbidSQL)
	if err != nil {
		log.Error("d.db.Query(%s) error(%v)", _midsForbidSQL, err)
		return
	}
	defer rows.Close()
	mids = make(map[int64][]string)
	for rows.Next() {
		var (
			id    int64
			mid   int64
			value string
		)
		if err = rows.Scan(&id, &mid, &value); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		mids[mid] = append(mids[mid], value)
	}
	log.Info("mids (%+v)", mids)
	return
}

// AppFeedAids get aids  by  appFeed flow.
func (d *Dao) AppFeedAids(c context.Context, startTime, endTime time.Time) (aids []int64, err error) {
	rows, err := d.db.Query(c, _appFlowsSQL, startTime, endTime)
	if err != nil {
		log.Error("d.db.Query(%s|%v|%v) error(%v)", _appFlowsSQL, startTime, endTime, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var aid int64
		if err = rows.Scan(&aid); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		if aid == 0 {
			continue
		}
		aids = append(aids, aid)
	}
	return
}

//FlowGroupPools 获取指定流量套餐的pool
func (d *Dao) FlowGroupPools(c context.Context, ids []int64) (res map[int64]int8, err error) {
	var (
		rows *sql.Rows
		id   int64
		pool int8
	)
	res = map[int64]int8{}
	idstr := xstr.JoinInts(ids)
	if rows, err = d.db.Query(context.TODO(), fmt.Sprintf(_flowGroupPool, idstr)); err != nil {
		log.Error("FlowGroupPools d.db.Query error(%v) ids(%s)", err, idstr)
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&id, &pool); err != nil {
			log.Error("FlowGroupPools rows.Scan error(%v) ids(%s)", err, idstr)
			return
		}
		res[id] = pool
	}
	return
}

//FlowsByOID 获取所有命中的流量套餐记录
func (d *Dao) FlowsByOID(c context.Context, oid int64) (res []*archive.FlowData, err error) {
	var (
		rows *sql.Rows
	)
	res = []*archive.FlowData{}
	if rows, err = d.db.Query(context.TODO(), _flowsByOIDSQL, oid); err != nil {
		log.Error("FlowsByOID d.db.Query error(%v) oid(%d)", err, oid)
		return
	}
	defer rows.Close()

	for rows.Next() {
		f := &archive.FlowData{}
		if err = rows.Scan(&f.ID, &f.Pool, &f.OID, &f.GroupID, &f.Parent, &f.State, &f.GroupValue); err != nil {
			log.Error("FlowsByOID rows.Scan error(%v) oid(%d)", err, oid)
			return
		}
		res = append(res, f)
	}
	return
}

//FlowUnique 获取命中 指定流量套餐的记录
func (d *Dao) FlowUnique(c context.Context, oid, groupID int64, pool int8) (f *archive.FlowData, err error) {
	f = &archive.FlowData{}
	if err = d.db.QueryRow(context.TODO(), _flowUniqueSQL, oid, pool, groupID).Scan(&f.ID, &f.Pool, &f.OID, &f.GroupID, &f.Parent, &f.State); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			f = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

//OidsFlowByGID  判断指定oids 是否是gid禁止的
func (d *Dao) OidsFlowByGID(c context.Context, pool, gid int64, oids string) (res []*archive.FlowData, err error) {
	var (
		rows *sql.Rows
	)
	rows, err = d.db.Query(c, fmt.Sprintf(_flowOidsByGidSQL, oids), pool, gid)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		f := &archive.FlowData{}
		if err = rows.Scan(&f.ID, &f.Pool, &f.OID, &f.GroupID, &f.Parent, &f.State, &f.GroupValue); err != nil {
			log.Error("FlowsByOID rows.Scan error(%v) oid(%d)", err, gid)
			return
		}
		res = append(res, f)
	}
	return
}

// CountByGID  count buy state.
func (d *Dao) CountByGID(c context.Context, pool, gid int64) (count int64, err error) {
	row := d.db.QueryRow(c, _flowCountSQL, pool, gid)
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			err = errors.WithStack(err)
		}
	}
	return
}

// FlowPage page.
func (d *Dao) FlowPage(c context.Context, pool, gid, pn, ps int64) (res []*archive.FlowData, err error) {
	var rows *sql.Rows
	rows, err = d.db.Query(c, _flowsByGIDSQL, pool, gid, (pn-1)*ps, ps)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		f := &archive.FlowData{}
		if err = rows.Scan(&f.ID, &f.Pool, &f.OID, &f.GroupID, &f.Parent, &f.State, &f.GroupValue); err != nil {
			log.Error("FlowsByOID rows.Scan error(%v) oid(%d)", err, gid)
			return
		}
		res = append(res, f)
	}
	err = rows.Err()
	return
}
