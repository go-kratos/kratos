package dao

import (
	"context"
	"database/sql"
	"fmt"

	"go-common/app/admin/main/usersuit/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	_sharding = 10

	_selMedal     = "SELECT id,name,description,cond,gid,level,level_rank,sort,is_online FROM medal_info"
	_selMedalByID = "SELECT id,name,description,image,image_small,cond,gid,level,level_rank,sort,is_online FROM medal_info WHERE id= ?"
	_insertMedal  = "INSERT INTO medal_info (name,description,cond,gid,level,sort,level_rank,is_online,image,image_small) VALUES(?,?,?,?,?,?,?,?,?,?)"
	_updateMedal  = "UPDATE medal_info SET name=?,description=?,cond=?,gid=?,level=?,sort=?,level_rank=?,is_online=?,image=?,image_small=? WHERE id=?"

	_sellMedalGroup       = "SELECT id,name,pid,rank,is_online FROM medal_group ORDER BY id ASC, rank ASC"
	_sellMedalGroupInfo   = "SELECT m1.id,m1.name,m1.pid,m1.rank,m1.is_online,ifnull(m2.name,'无') FROM medal_group m1 left join medal_group m2 on m1.pid=m2.id  ORDER BY pid ASC, rank ASC"
	_selMedalGroupParent  = "SELECT id,name,pid,rank,is_online FROM medal_group WHERE pid=0"
	_selMedalGroupByID    = "SELECT id,name,pid,rank,is_online FROM medal_group WHERE id= ?"
	_insertMedalGroup     = "INSERT INTO medal_group (name,pid,rank,is_online) VALUES (?,?,?,?)"
	_updateMedalGroupByID = "UPDATE medal_group SET name=?,pid=?,rank=?,is_online=? WHERE id =?"

	_selMedalOwnerByMID         = "SELECT o.id,o.nid,o.is_activated,o.is_del,i.name FROM medal_owner_%s o LEFT JOIN medal_info i ON o.nid=i.id WHERE o.mid=?"
	_selMedalAddList            = "SELECT id,name FROM medal_info WHERE id NOT IN (SELECT nid FROM medal_owner_%s WHERE mid=?)"
	_countOwnerBYNidMidSQL      = "SELECT COUNT(*) FROM medal_owner_%s WHERE mid=? AND nid=?"
	_insertMedalOwner           = "INSERT INTO medal_owner_%s (mid,nid) VALUES (?,?)"
	_updatMedalOwnerActivated   = "UPDATE medal_owner_%s SET is_activated=1 WHERE mid=? AND nid=?"
	_updatMedalOwnerNoActivated = "UPDATE medal_owner_%s SET is_activated=0 WHERE mid=? AND nid!=?"
	_updatMedalOwnerDel         = "UPDATE medal_owner_%s SET is_del=? WHERE mid=? AND nid=?"

	_insertMedalOperationLogSQL = "INSERT INTO medal_operation_log(oper_id,mid,medal_id,source_type,action) VALUES (?,?,?,?,?)"
	_medalOperationLogTotalSQL  = "SELECT COUNT(*) FROM medal_operation_log WHERE mid=?"
	_medalOperationLogSQL       = "SELECT oper_id,action,mid,medal_id,source_type,ctime,mtime FROM medal_operation_log WHERE mid=? ORDER BY mtime DESC LIMIT ?,?"
)

func (d *Dao) hit(id int64) string {
	return fmt.Sprintf("%d", id%_sharding)
}

// Medal medal info .
func (d *Dao) Medal(c context.Context) (ms []*model.Medal, err error) {
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, _selMedal); err != nil {
		log.Error("Medal, d.db.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		m := &model.Medal{}
		if err = rows.Scan(&m.ID, &m.Name, &m.Description, &m.Condition, &m.GID, &m.Level, &m.LevelRank, &m.Sort, &m.IsOnline); err != nil {
			log.Error("Medal, rows.Scan(%+v) error(%v)", m, err)
			return
		}
		ms = append(ms, m)
	}
	err = rows.Err()
	return
}

// MedalByID medal info by id .
func (d *Dao) MedalByID(c context.Context, id int64) (m *model.Medal, err error) {
	var row = d.db.QueryRow(c, _selMedalByID, id)
	m = &model.Medal{}
	if err = row.Scan(&m.ID, &m.Name, &m.Description, &m.Image, &m.ImageSmall, &m.Condition, &m.GID, &m.Level, &m.LevelRank, &m.Sort, &m.IsOnline); err != nil {
		log.Error("MedalByID, rows.Scan(%+v) error(%v)", m, err)
	}
	return
}

// AddMedal add medal .
func (d *Dao) AddMedal(c context.Context, m *model.Medal) (affected int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _insertMedal, m.Name, m.Description, m.Condition, m.GID, m.Level, m.Sort, m.LevelRank, m.IsOnline, m.Image, m.ImageSmall); err != nil {
		log.Error("AddMedal, rows.Exec(%+v) error(%v)", m, err)
		return
	}
	return res.RowsAffected()
}

// UpMedal update nameplate .
func (d *Dao) UpMedal(c context.Context, id int64, m *model.Medal) (affected int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _updateMedal, m.Name, m.Description, m.Condition, m.GID, m.Level, m.Sort, m.LevelRank, m.IsOnline, m.Image, m.ImageSmall, id); err != nil {
		log.Error("UpMedal, d.db.Exec(%+v,%d) error(%v)", m, id, err)
		return
	}
	return res.RowsAffected()
}

// MedalGroup return medal group all .
func (d *Dao) MedalGroup(c context.Context) (res map[int64]*model.MedalGroup, err error) {
	res = make(map[int64]*model.MedalGroup)
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, _sellMedalGroup); err != nil {
		log.Error("MedalGroup, d.db.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		re := &model.MedalGroup{}
		if err = rows.Scan(&re.ID, &re.Name, &re.PID, &re.Rank, &re.IsOnline); err != nil {
			log.Error("MedalGroup, rows.Scan(%+v) error(%v)", re, err)
			return
		}
		res[re.ID] = re
	}
	err = rows.Err()
	return
}

// MedalGroupInfo return medal group all info include parent group name .
func (d *Dao) MedalGroupInfo(c context.Context) (res []*model.MedalGroup, err error) {
	res = make([]*model.MedalGroup, 0)
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, _sellMedalGroupInfo); err != nil {
		log.Error("MedalGroupInfo, d.db.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		re := &model.MedalGroup{}
		if err = rows.Scan(&re.ID, &re.Name, &re.PID, &re.Rank, &re.IsOnline, &re.PName); err != nil {
			log.Error("MedalGroupInfo, rows.Scan(%+v) error(%v)", re, err)
			return
		}
		res = append(res, re)
	}
	err = rows.Err()
	return
}

// MedalGroupParent return medal group parent info .
func (d *Dao) MedalGroupParent(c context.Context) (res []*model.MedalGroup, err error) {
	res = make([]*model.MedalGroup, 0)
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, _selMedalGroupParent); err != nil {
		log.Error("MedalGroupInfo, d.db.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		re := &model.MedalGroup{}
		if err = rows.Scan(&re.ID, &re.Name, &re.PID, &re.Rank, &re.IsOnline); err != nil {
			log.Error("MedalGroupParent, rows.Scan(%+v) error(%v)", re, err)
			return
		}
		res = append(res, re)
	}
	err = rows.Err()
	return
}

// MedalGroupByID medal group by gid .
func (d *Dao) MedalGroupByID(c context.Context, id int64) (mg *model.MedalGroup, err error) {
	var row = d.db.QueryRow(c, _selMedalGroupByID, id)
	mg = &model.MedalGroup{}
	if err = row.Scan(&mg.ID, &mg.Name, &mg.PID, &mg.Rank, &mg.IsOnline); err != nil {
		log.Error("MedalGroupByID, rows.Scan(%+v) error(%v)", mg, err)
		return
	}
	return
}

// MedalGroupAdd add medal group .
func (d *Dao) MedalGroupAdd(c context.Context, mg *model.MedalGroup) (affected int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _insertMedalGroup, mg.Name, mg.PID, mg.Rank, mg.IsOnline); err != nil {
		log.Error("MedalGroupAdd, d.db.Exec(%+v) error(%v)", mg, err)
		return
	}
	return res.RowsAffected()
}

// MedalGroupUp update name group .
func (d *Dao) MedalGroupUp(c context.Context, id int64, mg *model.MedalGroup) (affected int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _updateMedalGroupByID, mg.Name, mg.PID, mg.Rank, mg.IsOnline, id); err != nil {
		log.Error("MedalGroupUp, d.db.Exec(%+v %d) error(%v)", mg, id, err)
		return
	}
	return res.RowsAffected()
}

// MedalOwner medal owner .
func (d *Dao) MedalOwner(c context.Context, mid int64) (res []*model.MedalMemberMID, err error) {
	var (
		rows   *xsql.Rows
		sqlStr string
	)
	res = make([]*model.MedalMemberMID, 0)
	sqlStr = fmt.Sprintf(_selMedalOwnerByMID, d.hit(mid))
	if rows, err = d.db.Query(c, sqlStr, mid); err != nil {
		log.Error("MedalOwner, d.db.Query(%s %d) error(%v)", sqlStr, mid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		re := &model.MedalMemberMID{}
		if err = rows.Scan(&re.ID, &re.NID, &re.IsActivated, &re.IsDel, &re.MedalName); err != nil {
			log.Error("MedalOwner, rows.Scan(%+v) error(%v)", re, err)
			return
		}
		res = append(res, re)
	}
	err = rows.Err()
	return
}

// CountOwnerBYNidMid retun number of medal_owner by mid and nid.
func (d *Dao) CountOwnerBYNidMid(c context.Context, mid, nid int64) (count int64, err error) {
	row := d.db.QueryRow(c, fmt.Sprintf(_countOwnerBYNidMidSQL, d.hit(mid)), mid, nid)
	if err = row.Scan(&count); err != nil {
		if err != sql.ErrNoRows {
			err = errors.Wrap(err, "CountOwnerBYNidMid")
			return
		}
		count = 0
		err = nil
	}
	return
}

// MedalOwnerAdd add medal owner .
func (d *Dao) MedalOwnerAdd(c context.Context, mid, nid int64) (affected int64, err error) {
	var (
		res    sql.Result
		sqlStr string
	)
	sqlStr = fmt.Sprintf(_insertMedalOwner, d.hit(mid))
	if res, err = d.db.Exec(c, sqlStr, mid, nid); err != nil {
		log.Error("MedalOwnerAdd, d.db.Exec(%s %d %d) error(%v)", sqlStr, mid, nid, err)
		return
	}
	return res.RowsAffected()
}

// MedalAddList .
func (d *Dao) MedalAddList(c context.Context, mid int64) (ms []*model.MedalMemberAddList, err error) {
	var (
		rows   *xsql.Rows
		sqlStr string
	)
	sqlStr = fmt.Sprintf(_selMedalAddList, d.hit(mid))
	if rows, err = d.db.Query(c, sqlStr, mid); err != nil {
		log.Error("MedalAddList, d.db.Query(%s %d) error(%v)", sqlStr, mid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		m := &model.MedalMemberAddList{}
		if err = rows.Scan(&m.ID, &m.MedalName); err != nil {
			log.Error("MedalAddList, rows.Scan(%+v) error(%v)", m, err)
			return
		}
		ms = append(ms, m)
	}
	err = rows.Err()
	return
}

// MedalOwnerUpActivated update medal owner is_activated=1.
func (d *Dao) MedalOwnerUpActivated(c context.Context, mid, nid int64) (affected int64, err error) {
	var (
		res    sql.Result
		sqlStr string
	)
	sqlStr = fmt.Sprintf(_updatMedalOwnerActivated, d.hit(mid))
	if res, err = d.db.Exec(c, sqlStr, mid, nid); err != nil {
		log.Error("MedalOwnerUpActivated, d.db.Exec(%s %d %d) error(%v)", sqlStr, mid, nid, err)
		return
	}
	return res.RowsAffected()
}

// MedalOwnerUpNotActivated update medal owner is_activated=0.
func (d *Dao) MedalOwnerUpNotActivated(c context.Context, mid, nid int64) (affected int64, err error) {
	var (
		res    sql.Result
		sqlStr string
	)
	sqlStr = fmt.Sprintf(_updatMedalOwnerNoActivated, d.hit(mid))
	if res, err = d.db.Exec(c, sqlStr, mid, nid); err != nil {
		log.Error("MedalOwnerUpNotActivated, d.db.Exec(%s %d %d) error(%v)", sqlStr, mid, nid, err)
		return
	}
	return res.RowsAffected()
}

// MedalOwnerDel update medal owner is_del=1.
func (d *Dao) MedalOwnerDel(c context.Context, mid, nid int64, isDel int8) (affected int64, err error) {
	var (
		res    sql.Result
		sqlStr string
	)
	sqlStr = fmt.Sprintf(_updatMedalOwnerDel, d.hit(mid))
	if res, err = d.db.Exec(c, sqlStr, isDel, mid, nid); err != nil {
		log.Error("MedalOwnerDel, d.db.Exec(%s %d %d %d) error(%v)", sqlStr, isDel, mid, nid, err)
		return
	}
	return res.RowsAffected()
}

// AddMedalOperLog insert medal operation log.
func (d *Dao) AddMedalOperLog(c context.Context, oid int64, mid int64, medalID int64, action string) (affected int64, err error) {
	var res sql.Result
	// oper_id,mid,medal_id,source_type,action
	if res, err = d.db.Exec(c, _insertMedalOperationLogSQL, oid, mid, medalID, model.MedalSourceTypeAdmin, action); err != nil {
		log.Error("MedalGroupAdd, d.db.Exec() error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// MedalOperLog get medal operation log.
func (d *Dao) MedalOperLog(c context.Context, mid int64, pn, ps int) (opers []*model.MedalOperLog, uids []int64, err error) {
	var (
		rows   *xsql.Rows
		offset = (pn - 1) * ps
	)
	if rows, err = d.db.Query(c, _medalOperationLogSQL, mid, offset, ps); err != nil {
		err = errors.Wrapf(err, "MedalOperLog d.db.Query(%d,%d,%d)", mid, offset, ps)
		return
	}
	defer rows.Close()
	for rows.Next() {
		oper := new(model.MedalOperLog)
		if err = rows.Scan(&oper.OID, &oper.Action, &oper.MID, &oper.MedalID, &oper.SourceType, &oper.CTime, &oper.MTime); err != nil {
			err = errors.Wrap(err, "MedalOperLog row.Scan()")
			return
		}
		opers = append(opers, oper)
		uids = append(uids, oper.MID)
	}
	err = rows.Err()
	return
}

// MedalOperationLogTotal medal operation log  total.
func (d *Dao) MedalOperationLogTotal(c context.Context, mid int64) (count int64, err error) {
	row := d.db.QueryRow(c, _medalOperationLogTotalSQL, mid)
	if err = row.Scan(&count); err != nil {
		err = errors.Wrap(err, "d.dao.MedalOperationLogTotal")
		return
	}
	return
}
