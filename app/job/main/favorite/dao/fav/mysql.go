package fav

import (
	"context"
	"database/sql"
	"fmt"

	favmdl "go-common/app/service/main/favorite/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	xtime "go-common/library/time"
	"go-common/library/xstr"
)

const (
	_countSharding    int64 = 50  // count by oid
	_folderSharding   int64 = 100 // folder by mid
	_relationSharding int64 = 500 // relations in folder by mid
	_usersSharding    int64 = 500 // users by oid

	// folder
	_folderSQL      = "SELECT id,type,mid,name,cover,description,count,attr,state,ctime,mtime FROM fav_folder_%s WHERE id=? AND type=? AND mid=? AND state=0"
	_upFolderCntSQL = "UPDATE fav_folder_%s SET count=?,mtime=? WHERE id=?"
	// relation
	_relationSQL           = "SELECT id,type,oid,mid,fid,state,ctime,mtime,sequence FROM fav_relation_%s  WHERE type=? AND mid=? AND fid=? AND oid=? AND state=0"
	_relationsSQL          = "SELECT id,type,oid,mid,fid,state,ctime,mtime,sequence FROM fav_relation_%s  FORCE INDEX(ix_fid_state_type_mtime)  WHERE fid=? AND mtime>=? AND state=0 AND type=? ORDER BY mtime ASC LIMIT ?"
	_allRelationsSQL       = "SELECT id,type,oid,mid,fid,state,ctime,mtime,sequence FROM fav_relation_%s  FORCE INDEX(ix_fid_state_mtime)  WHERE fid=? AND mtime>=? AND state=0 ORDER BY mtime ASC LIMIT ?"
	_relationOidsSQL       = "SELECT oid FROM fav_relation_%s FORCE INDEX(ix_fid_state_type_mtime) WHERE fid=? AND state=0 AND type=? ORDER BY mtime DESC LIMIT ?,?"
	_relationFidsByOidSQL  = "SELECT fid FROM fav_relation_%s WHERE type=? AND mid=? AND oid=? and state=0"
	_relationFidsSQL       = "SELECT oid,fid FROM fav_relation_%s WHERE type=? AND mid=? AND state=0"
	_cntRelationSQL        = "SELECT COUNT(id) FROM fav_relation_%s FORCE INDEX(ix_fid_state_sequence) WHERE fid=? and state=0"
	_addRelationSQL        = "INSERT INTO fav_relation_%s (type,oid,mid,fid,state,ctime,mtime,sequence) VALUES(?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE state=?,mtime=?,sequence=?"
	_delRelationSQL        = "UPDATE fav_relation_%s SET state=1,mtime=? WHERE type=? AND fid=? AND oid=?"
	_delRelationsByOidsSQL = "UPDATE fav_relation_%s SET state=1,mtime=? WHERE type=? AND fid=? AND oid in (%s)"
	_selRelationsByOidsSql = "SELECT id,type,oid,mid,fid,state,ctime,mtime,sequence FROM fav_relation_%s  WHERE type=? AND fid=? AND oid in (%s)"

	_updateAllRelationSeq = "UPDATE fav_relation_%s SET mtime=mtime,sequence=(case id %s end) WHERE id in (%s)"
	_updateRelationSeq    = "UPDATE fav_relation_%s SET sequence=?,mtime=? WHERE fid=? and oid=? and type=?"
	_selectMaxSequence    = "SELECT id,type,oid,mid,fid,state,ctime,mtime,sequence FROM fav_relation_%s FORCE INDEX(ix_fid_state_sequence) WHERE fid=? AND state=0 order BY sequence desc limit 1"
	_recentOidsSQL        = "SELECT oid,type FROM fav_relation_%s FORCE INDEX(ix_fid_state_mtime) WHERE fid=? AND state=0 ORDER BY mtime DESC LIMIT 3"
	// users
	_addUserSQL = "INSERT INTO fav_users_%s (type,oid,mid,state,ctime,mtime) VALUES(?,?,?,?,?,?) ON DUPLICATE KEY UPDATE state=?,mtime=?"
	_delUserSQL = "UPDATE fav_users_%s SET state=1,mtime=? WHERE type=? AND oid=? AND mid=?"
	// stat
	_statCntSQL   = "SELECT count from fav_count_%s WHERE type=? AND oid=?"
	_upStatCntSQL = "INSERT INTO fav_count_%s (type,oid,count,ctime,mtime) VALUES(?,?,?,?,?) ON DUPLICATE KEY UPDATE count=count+?,mtime=?"
)

func cntHit(mid int64) string {
	return fmt.Sprintf("%02d", mid%_countSharding)
}

func folderHit(mid int64) string {
	return fmt.Sprintf("%02d", mid%_folderSharding)
}

func relationHit(mid int64) string {
	return fmt.Sprintf("%03d", mid%_relationSharding)
}

func usersHit(oid int64) string {
	return fmt.Sprintf("%03d", oid%_usersSharding)
}

// pingMySQL check mysql connection.
func (d *Dao) pingMySQL(c context.Context) error {
	return d.db.Ping(c)
}

func (d *Dao) TxUpdateFavSequence(tx *xsql.Tx, mid int64, fid, oid int64, typ int8, sequence uint64, mtime xtime.Time) (rows int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_updateRelationSeq, relationHit(mid)), sequence, mtime, fid, oid, typ)
	if err != nil {
		log.Error("d.db.Exec(%s,%d,%d,%d,%d,%d) error(%v)", _updateRelationSeq, sequence, mtime, fid, oid, typ, err)
		return
	}
	return res.RowsAffected()
}

func (d *Dao) BatchUpdateSeq(c context.Context, mid int64, favs []*favmdl.Favorite) (int64, error) {
	var caseStr string
	var ids []int64
	for _, fav := range favs {
		ids = append(ids, fav.ID)
		caseStr += fmt.Sprintf("when %d then %d ", fav.ID, fav.Sequence)
	}
	if len(ids) <= 0 {
		return 0, nil
	}
	res, err := d.db.Exec(c, fmt.Sprintf(_updateAllRelationSeq, relationHit(mid), caseStr, xstr.JoinInts(ids)))
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// RecentRes return user's three newest fav from a folder.
func (d *Dao) RecentRes(c context.Context, mid, fid int64) (res []*favmdl.Resource, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_recentOidsSQL, relationHit(mid)), fid)
	if err != nil {
		log.Error("d.db.Query(%d,%d) error(%v)", mid, fid, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var oid int64
		var typ int32
		if err = rows.Scan(&oid, &typ); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		res = append(res, &favmdl.Resource{Oid: oid, Typ: typ})
	}
	err = rows.Err()
	return
}

// Folder get a Folder by fid from mysql.
func (d *Dao) Folder(c context.Context, tp int8, mid, fid int64) (f *favmdl.Folder, err error) {
	f = &favmdl.Folder{}
	row := d.db.QueryRow(c, fmt.Sprintf(_folderSQL, folderHit(mid)), fid, tp, mid)
	if err = row.Scan(&f.ID, &f.Type, &f.Mid, &f.Name, &f.Cover, &f.Description, &f.Count, &f.Attr, &f.State, &f.CTime, &f.MTime); err != nil {
		if err == sql.ErrNoRows {
			f = nil
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// UpFolderCnt update folder count to mysql.
func (d *Dao) UpFolderCnt(c context.Context, mid, fid int64, cnt int, now xtime.Time) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_upFolderCntSQL, folderHit(mid)), cnt, now, fid)
	if err != nil {
		log.Error("db.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// MaxRelation get a max sequence relation from mysql.
func (d *Dao) MaxRelation(c context.Context, mid, fid int64) (m *favmdl.Favorite, err error) {
	m = &favmdl.Favorite{}
	row := d.db.QueryRow(c, fmt.Sprintf(_selectMaxSequence, relationHit(mid)), fid)
	if err = row.Scan(&m.ID, &m.Type, &m.Oid, &m.Mid, &m.Fid, &m.State, &m.CTime, &m.MTime, &m.Sequence); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			m = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// Relation get a relation from mysql.
func (d *Dao) Relation(c context.Context, tp int8, mid, fid, oid int64) (m *favmdl.Favorite, err error) {
	m = &favmdl.Favorite{}
	row := d.db.QueryRow(c, fmt.Sprintf(_relationSQL, relationHit(mid)), tp, mid, fid, oid)
	if err = row.Scan(&m.ID, &m.Type, &m.Oid, &m.Mid, &m.Fid, &m.State, &m.CTime, &m.MTime, &m.Sequence); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			m = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// AllRelations get favorite relations from mysql.
func (d *Dao) AllRelations(c context.Context, mid, fid int64, mtime xtime.Time, limit int) (fr []*favmdl.Favorite, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_allRelationsSQL, relationHit(mid)), fid, mtime, limit)
	if err != nil {
		log.Error("d.db.Query(%s,%d,%d,%d,%d) error(%v)", fmt.Sprintf(_allRelationsSQL, relationHit(mid)), mid, fid, mtime, limit, err)
		return
	}
	defer rows.Close()
	fr = make([]*favmdl.Favorite, 0)
	for rows.Next() {
		var r = &favmdl.Favorite{}
		if err = rows.Scan(&r.ID, &r.Type, &r.Oid, &r.Mid, &r.Fid, &r.State, &r.CTime, &r.MTime, &r.Sequence); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		if r.Mid != mid {
			log.Error("dirty data relations(%d,%d)", mid, fid)
			continue
		}
		fr = append(fr, r)
	}
	err = rows.Err()
	return
}

func (d *Dao) RelationsByOids(c context.Context, typ int8, mid, fid int64, oids []int64) (fr []*favmdl.Favorite, err error) {
	if len(oids) <= 0 {
		return
	}
	rows, err := d.db.Query(c, fmt.Sprintf(_selRelationsByOidsSql, relationHit(mid), xstr.JoinInts(oids)), typ, fid)
	if err != nil {
		log.Error("d.db.Query(%s,%d,%d,%v) error(%v)", fmt.Sprintf(_selRelationsByOidsSql, relationHit(mid), xstr.JoinInts(oids)), typ, fid, err)
		return
	}
	defer rows.Close()
	fr = make([]*favmdl.Favorite, 0)
	for rows.Next() {
		var r = &favmdl.Favorite{}
		if err = rows.Scan(&r.ID, &r.Type, &r.Oid, &r.Mid, &r.Fid, &r.State, &r.CTime, &r.MTime, &r.Sequence); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		if r.Type != typ || r.Mid != mid {
			log.Error("dirty data relations(%d,%d,%d)", typ, mid, fid)
			continue
		}
		fr = append(fr, r)
	}
	err = rows.Err()
	return
}

// Relations get favorite relations from mysql.
func (d *Dao) Relations(c context.Context, typ int8, mid, fid int64, mtime xtime.Time, limit int) (fr []*favmdl.Favorite, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_relationsSQL, relationHit(mid)), fid, mtime, typ, limit)
	if err != nil {
		log.Error("d.db.Query(%s,%d,%d,%d,%d,%d) error(%v)", fmt.Sprintf(_relationsSQL, relationHit(mid)), typ, mid, fid, mtime, limit, err)
		return
	}
	defer rows.Close()
	fr = make([]*favmdl.Favorite, 0)
	for rows.Next() {
		var r = &favmdl.Favorite{}
		if err = rows.Scan(&r.ID, &r.Type, &r.Oid, &r.Mid, &r.Fid, &r.State, &r.CTime, &r.MTime, &r.Sequence); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		if r.Mid != mid {
			log.Error("dirty data relations(%d,%d,%d)", typ, mid, fid)
			continue
		}
		fr = append(fr, r)
	}
	err = rows.Err()
	return
}

// RelationFidsByOid get favortied folders in relations by oid from mysql.
func (d *Dao) RelationFidsByOid(c context.Context, tp int8, mid, oid int64) (fids []int64, err error) {
	var (
		fid int64
		idx = relationHit(mid)
	)
	rows, err := d.db.Query(c, fmt.Sprintf(_relationFidsByOidSQL, idx), tp, mid, oid)
	if err != nil {
		log.Error("d.db.QueryRow(%d,%d,%d) error(%v)", mid, mid, oid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&fid); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		fids = append(fids, fid)
	}
	err = rows.Err()
	return
}

// RelationCnt get favoried folders count in relation from mysql.
func (d *Dao) RelationCnt(c context.Context, mid, fid int64) (cnt int, err error) {
	row := d.db.QueryRow(c, fmt.Sprintf(_cntRelationSQL, relationHit(mid)), fid)
	if err = row.Scan(&cnt); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// AddRelation add a favorite relation to mysql.
func (d *Dao) AddRelation(c context.Context, fr *favmdl.Favorite) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_addRelationSQL, relationHit(fr.Mid)), fr.Type, fr.Oid, fr.Mid, fr.Fid, fr.State, fr.CTime, fr.MTime, fr.Sequence, fr.State, fr.MTime, fr.Sequence)
	if err != nil {
		log.Error("db.Exec error(%v)", err)
		return
	}
	return res.LastInsertId()
}

// DelRelation delete a favorite relation to mysql.
func (d *Dao) DelRelation(c context.Context, tp int8, mid, fid, oid int64, now xtime.Time) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_delRelationSQL, relationHit(mid)), now, tp, fid, oid)
	if err != nil {
		log.Error("d.db.Exec(%d,%d,%d,%d) error(%v)", mid, tp, fid, oid, err)
		return
	}
	return res.RowsAffected()
}

// UpStatCnt update stat count to mysql.
func (d *Dao) UpStatCnt(c context.Context, tp int8, oid int64, incr int, now xtime.Time) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_upStatCntSQL, cntHit(oid)), tp, oid, incr, now, now, incr, now)
	if err != nil {
		log.Error("db.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// StatCnt return stat count from mysql.
func (d *Dao) StatCnt(c context.Context, tp int8, oid int64) (cnt int, err error) {
	row := d.db.QueryRow(c, fmt.Sprintf(_statCntSQL, cntHit(oid)), tp, oid)
	if err = row.Scan(&cnt); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// RelationFids get the map [oid -> "fid,fid"] for user.
func (d *Dao) RelationFids(c context.Context, tp int8, mid int64) (rfids map[int64][]int64, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_relationFidsSQL, relationHit(mid)), tp, mid)
	if err != nil {
		log.Error("db.Query(%d) error(%v)", mid, err)
		return
	}
	defer rows.Close()
	rfids = make(map[int64][]int64, 128)
	for rows.Next() {
		var oid, fid int64
		if err = rows.Scan(&oid, &fid); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		rfids[oid] = append(rfids[oid], fid)
	}
	err = rows.Err()
	return
}

// OidsByFid get oids by fid
func (d *Dao) OidsByFid(c context.Context, typ int8, mid, fid int64, offset, limit int) (oids []int64, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_relationOidsSQL, relationHit(mid)), fid, typ, offset, limit)
	if err != nil {
		log.Error("d.db.Query(%d,%d,%d,%d,%d) error(%v)", typ, mid, fid, offset, limit, err)
		return
	}
	defer rows.Close()
	oid := int64(0)
	oids = make([]int64, 0, limit)
	for rows.Next() {
		if err = rows.Scan(&oid); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		oids = append(oids, oid)
	}
	err = rows.Err()
	return
}

// DelRelationsByOids update del state by fid and oids.
func (d *Dao) DelRelationsByOids(c context.Context, typ int8, mid, fid int64, oids []int64, now xtime.Time) (rows int64, err error) {
	if len(oids) == 0 {
		return
	}
	oidsStr := xstr.JoinInts(oids)
	res, err := d.db.Exec(c, fmt.Sprintf(_delRelationsByOidsSQL, relationHit(mid), oidsStr), now, typ, fid)
	if err != nil {
		log.Error("d.db.Exec(%d,%d,%d,%s) error(%v)", typ, mid, fid, oidsStr, now, err)
		return
	}
	return res.RowsAffected()
}

// AddUser add a favorite user to mysql.
func (d *Dao) AddUser(c context.Context, u *favmdl.User) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_addUserSQL, usersHit(u.Oid)), u.Type, u.Oid, u.Mid, u.State, u.CTime, u.MTime, u.State, u.MTime)
	if err != nil {
		log.Error("db.Exec error(%v)", err)
		return
	}
	return res.LastInsertId()
}

// DelUser delete a favorite user.
func (d *Dao) DelUser(c context.Context, u *favmdl.User) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_delUserSQL, usersHit(u.Oid)), u.MTime, u.Type, u.Oid, u.Mid)
	if err != nil {
		log.Error("d.db.Exec(%d,%d,%d,%d) error(%v)", u.MTime, u.Type, u.Oid, u.Mid, err)
		return
	}
	return res.RowsAffected()
}
