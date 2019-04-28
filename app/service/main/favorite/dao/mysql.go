package dao

import (
	"bytes"
	"context"
	"fmt"

	"go-common/app/service/main/favorite/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_folderSharding   int64 = 100 // folder by mid
	_relationSharding int64 = 500 // objects in folder by mid
	_usersSharding    int64 = 500 // objects faved by oid
	_countSharding    int64 = 50  // objects count by oid
	// folder table
	_cntFolderSQL    = "SELECT COUNT(id) FROM fav_folder_%s WHERE type=? AND mid=? AND state=0"
	_folderSQL       = "SELECT id,type,mid,name,cover,description,count,attr,state,ctime,mtime FROM fav_folder_%s WHERE id=? AND type=? AND mid=?"
	_userFoldersSQL  = "SELECT id,type,mid,name,cover,description,count,attr,state,ctime,mtime FROM fav_folder_%s WHERE type=? AND mid=? AND state=0"
	_folderByNameSQL = "SELECT id,type,mid,name,cover,description,count,attr,state,ctime,mtime FROM fav_folder_%s WHERE name=? AND type=? AND mid=? AND state=0"
	_folderByIdsSQL  = "SELECT id,type,mid,name,cover,description,count,attr,state,ctime,mtime FROM fav_folder_%s WHERE id in (%s)"
	_defFolderSQL    = "SELECT id,type,mid,name,cover,description,count,attr,state,ctime,mtime FROM fav_folder_%s WHERE type=? AND mid=? AND attr&2=0"
	_addFolderSQL    = `INSERT INTO fav_folder_%s (type,mid,name,cover,description,count,attr,state,ctime,mtime) VALUES (?,?,?,?,?,?,?,?,?,?)
						ON DUPLICATE KEY UPDATE name=?,cover=?,description=?,count=?,attr=?,state=?,ctime=?,mtime=?`
	_delFolderSQL    = "UPDATE fav_folder_%s SET state=1 WHERE type=? AND id=?"
	_updateFolderSQL = "UPDATE fav_folder_%s SET name=?,description=?,cover=?,attr=?,state=?,mtime=? WHERE type=? AND id=?"
	_upFolderNameSQL = "UPDATE IGNORE fav_folder_%s SET name=? WHERE id=?"
	_upFolderAttrSQL = "UPDATE fav_folder_%s SET attr=? WHERE id=?"

	// relation table
	_relationSQL       = "SELECT id,type,oid,mid,fid,state,ctime,mtime,sequence FROM fav_relation_%s WHERE type=? AND mid=? AND fid=? AND oid=? AND state=0"
	_relationsSQL      = "SELECT id,type,oid,mid,fid,state,ctime,mtime,sequence FROM fav_relation_%s FORCE INDEX(ix_fid_state_type_mtime) WHERE fid=? AND state=0 AND type=? ORDER BY mtime DESC LIMIT ?,?"
	_allRelationsSQL   = "SELECT id,type,oid,mid,fid,state,ctime,mtime,sequence FROM fav_relation_%s FORCE INDEX(ix_fid_state_sequence) WHERE fid=? AND state=0 ORDER BY sequence DESC LIMIT ?,?"
	_relationFidsSQL   = "SELECT fid FROM fav_relation_%s WHERE type=? AND mid=? AND oid=? AND state=0"
	_FidsByOidsSQL     = "SELECT oid,fid FROM fav_relation_%s WHERE type=? AND mid=? AND oid in (%s) AND state=0"
	_cntRelationSQL    = "SELECT COUNT(id) FROM fav_relation_%s FORCE INDEX(ix_fid_state_type_mtime) WHERE fid=? AND state=0 AND type=?"
	_cntAllRelationSQL = "SELECT COUNT(id) FROM fav_relation_%s FORCE INDEX(ix_fid_state_sequence) WHERE fid=? AND state=0"
	_addRelationSQL    = "INSERT INTO fav_relation_%s (type,oid,mid,fid,state,ctime,mtime) VALUES(?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE state=?,mtime=?"
	_maddRelationsSQL  = "INSERT INTO fav_relation_%s (type,oid,mid,fid,state) VALUES %s ON DUPLICATE KEY UPDATE state=0"
	_delRelationSQL    = "UPDATE fav_relation_%s SET state=1 WHERE type=? AND fid=? AND oid=?"
	_delRelationsSQL   = "UPDATE fav_relation_%s SET state=1 WHERE type=? AND fid=? AND oid in (%s)"
	_recentOidsSQL     = "SELECT oid,type FROM fav_relation_%s FORCE INDEX(ix_fid_state_mtime)  WHERE fid=? AND state=0 ORDER BY mtime DESC LIMIT 3"
	_copyRelationsSQL  = `INSERT IGNORE INTO fav_relation_%s (type,oid,mid,fid) 
						 SELECT %d,oid,%d,%d FROM fav_relation_%s WHERE type=? AND mid=? AND fid=? AND oid in (%s) AND state=0
						 ON DUPLICATE KEY UPDATE state=0`
	_batchOidsSQL = `SELECT oid FROM fav_relation_%s WHERE type=? AND mid=? AND state=0 LIMIT ?`

	// users table
	_cntUsersSQL = "SELECT COUNT(id) FROM fav_users_%s WHERE type=? AND oid=? AND state=0"
	_usersSQL    = "SELECT id,type,oid,mid,state,ctime,mtime FROM fav_users_%s WHERE type=? AND oid=? AND state=0 ORDER BY mtime DESC LIMIT ?,?"

	// folderSort table
	_folderSortSQL    = "SELECT id,type,mid,sort,ctime,mtime FROM fav_folder_sort WHERE type=? AND mid=?"
	_setFolderSortSQL = "INSERT INTO fav_folder_sort (type,mid,sort,ctime,mtime) VALUES (?,?,?,?,?) ON DUPLICATE KEY UPDATE sort=?,mtime=?"

	// count table
	_countSQL      = "SELECT count FROM fav_count_%s WHERE type=? AND oid=?"
	_countsSQL     = "SELECT oid,count FROM fav_count_%s WHERE type=? AND oid in (%s)"
	_folderStatSQL = "SELECT fid,play,fav,share from fav_folder_stat_%s WHERE fid in (%s)"
)

// folderHit hit table by the mod of mid and _folderSharding.
func folderHit(mid int64) string {
	return fmt.Sprintf("%02d", mid%_folderSharding)
}

// relationHit hit table by the mod of mid and _relationSharding.
func relationHit(mid int64) string {
	return fmt.Sprintf("%03d", mid%_relationSharding)
}

// usersHit hit table by the mod of oid and _userSharding.
func usersHit(oid int64) string {
	return fmt.Sprintf("%03d", oid%_usersSharding)
}

// countHit hit table by the mod of oid and _countSharding.
func countHit(oid int64) string {
	return fmt.Sprintf("%02d", oid%_countSharding)
}

// pingMySQL check mysql connection.
func (d *Dao) pingMySQL(c context.Context) error {
	if err := d.dbRead.Ping(c); err != nil {
		return err
	}
	if err := d.dbPush.Ping(c); err != nil {
		return err
	}
	return d.db.Ping(c)
}

// Folder get a Folder by fid from mysql.
func (d *Dao) Folder(c context.Context, tp int8, mid, fid int64) (f *model.Folder, err error) {
	f = new(model.Folder)
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

// FolderByName get a Folder by name from mysql.
func (d *Dao) FolderByName(c context.Context, tp int8, mid int64, name string) (f *model.Folder, err error) {
	f = &model.Folder{}
	row := d.db.QueryRow(c, fmt.Sprintf(_folderByNameSQL, folderHit(mid)), name, tp, mid)
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

// DefaultFolder get default folder from mysql.
func (d *Dao) DefaultFolder(c context.Context, tp int8, mid int64) (f *model.Folder, err error) {
	f = new(model.Folder)
	row := d.db.QueryRow(c, fmt.Sprintf(_defFolderSQL, folderHit(mid)), tp, mid)
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

// AddFolder add a new favorite folder to mysql.
func (d *Dao) AddFolder(c context.Context, f *model.Folder) (fid int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_addFolderSQL, folderHit(f.Mid)), f.Type, f.Mid, f.Name, f.Cover, f.Description, f.Count, f.Attr, f.State, f.CTime, f.MTime,
		f.Name, f.Cover, f.Description, f.Count, f.Attr, f.State, f.CTime, f.MTime)
	if err != nil {
		log.Error("db.Exec error(%v)", err)
		return
	}
	return res.LastInsertId()
}

// UpdateFolder add a new favorite folder to mysql.
func (d *Dao) UpdateFolder(c context.Context, f *model.Folder) (fid int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_updateFolderSQL, folderHit(f.Mid)), f.Name, f.Description, f.Cover, f.Attr, f.State, f.MTime, f.Type, f.ID)
	if err != nil {
		log.Error("db.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// UpFolderName rename user's folder name to mysql.
func (d *Dao) UpFolderName(c context.Context, typ int8, mid, fid int64, name string) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_upFolderNameSQL, folderHit(mid)), name, fid)
	if err != nil {
		log.Error("db.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// UpFolderAttr update user's folder attr to mysql.
func (d *Dao) UpFolderAttr(c context.Context, typ int8, mid, fid int64, attr int32) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_upFolderAttrSQL, folderHit(mid)), attr, fid)
	if err != nil {
		log.Error("db.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// FolderRelations get favorite relations from mysql.
func (d *Dao) FolderRelations(c context.Context, typ int8, mid, fid int64, start, end int) (fr []*model.Favorite, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_relationsSQL, relationHit(mid)), fid, typ, start, end)
	if err != nil {
		log.Error("d.db.Query(%s,%d,%d,%d,%d,%d) error(%v)", fmt.Sprintf(_relationsSQL, relationHit(mid)), typ, mid, fid, start, end, err)
		return
	}
	defer rows.Close()
	fr = make([]*model.Favorite, 0)
	for rows.Next() {
		var r = &model.Favorite{}
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

// FolderAllRelations get favorite relations from mysql.
func (d *Dao) FolderAllRelations(c context.Context, mid, fid int64, start, end int) (fr []*model.Favorite, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_allRelationsSQL, relationHit(mid)), fid, start, end)
	if err != nil {
		log.Error("d.db.Query(%s,%d,%d,%d,%d,%d) error(%v)", fmt.Sprintf(_allRelationsSQL, relationHit(mid)), mid, fid, start, end, err)
		return
	}
	defer rows.Close()
	fr = make([]*model.Favorite, 0)
	for rows.Next() {
		var r = &model.Favorite{}
		if err = rows.Scan(&r.ID, &r.Type, &r.Oid, &r.Mid, &r.Fid, &r.State, &r.CTime, &r.MTime, &r.Sequence); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		if r.Mid != mid {
			log.Error("dirty data relations(%d,%d,%d)", mid, fid)
			continue
		}
		fr = append(fr, r)
	}
	err = rows.Err()
	return
}

// Folders get folders from mysql.
func (d *Dao) Folders(c context.Context, fvmids []*model.ArgFVmid) (fs map[string]*model.Folder, err error) {
	tblMap := make(map[string][]int64, len(fvmids))
	for _, fvmid := range fvmids {
		folderHits := folderHit(fvmid.Vmid)
		tblMap[folderHits] = append(tblMap[folderHits], fvmid.Fid)
	}
	fs = make(map[string]*model.Folder, len(fvmids))
	for folderHit, fids := range tblMap {
		fidsStr := xstr.JoinInts(fids)
		var rows *sql.Rows
		if rows, err = d.db.Query(c, fmt.Sprintf(_folderByIdsSQL, folderHit, fidsStr)); err != nil {
			log.Error("d.db.Query(%s,%s,%s) error(%v)", _folderByIdsSQL, folderHit, fidsStr, err)
			return
		}
		defer rows.Close()
		for rows.Next() {
			f := new(model.Folder)
			if err = rows.Scan(&f.ID, &f.Type, &f.Mid, &f.Name, &f.Cover, &f.Description, &f.Count, &f.Attr, &f.State, &f.CTime, &f.MTime); err != nil {
				log.Error("rows.Scan error(%v)", err)
				return
			}
			fmid := xstr.JoinInts([]int64{f.ID, f.Mid})
			fs[fmid] = f
		}
		err = rows.Err()
		if err != nil {
			return
		}
	}
	return
}

// RelationFidsByOid get favortied folders in relations by oid from mysql.
func (d *Dao) RelationFidsByOid(c context.Context, tp int8, mid, oid int64) (fids []int64, err error) {
	var fid int64
	rows, err := d.db.Query(c, fmt.Sprintf(_relationFidsSQL, relationHit(mid)), tp, mid, oid)
	if err != nil {
		log.Error("d.db.Query(%d,%d,%d) error(%v)", tp, mid, oid, err)
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

// RelationFidsByOids get favortied folders in relations by oid from mysql.
func (d *Dao) RelationFidsByOids(c context.Context, tp int8, mid int64, oids []int64) (fidsMap map[int64][]int64, err error) {
	oidsStr := xstr.JoinInts(oids)
	rows, err := d.dbRead.Query(c, fmt.Sprintf(_FidsByOidsSQL, relationHit(mid), oidsStr), tp, mid)
	if err != nil {
		log.Error("d.db.Query(%d,%d,%d) error(%v)", mid, mid, oidsStr, err)
		return
	}
	defer rows.Close()
	var oid, fid int64
	fidsMap = make(map[int64][]int64, len(oids))
	for rows.Next() {
		if err = rows.Scan(&oid, &fid); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		fidsMap[oid] = append(fidsMap[oid], fid)
	}
	err = rows.Err()
	return
}

// CntRelations get favoried folders count in relation from mysql.
func (d *Dao) CntRelations(c context.Context, mid, fid int64, typ int8) (count int, err error) {
	row := d.db.QueryRow(c, fmt.Sprintf(_cntRelationSQL, relationHit(mid)), fid, typ)
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// CntRelations get favoried folders count in relation from mysql.
func (d *Dao) CntAllRelations(c context.Context, mid, fid int64) (count int, err error) {
	row := d.db.QueryRow(c, fmt.Sprintf(_cntAllRelationSQL, relationHit(mid)), fid)
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// FolderCnt get user's folder count.
func (d *Dao) FolderCnt(c context.Context, tp int8, mid int64) (count int, err error) {
	row := d.db.QueryRow(c, fmt.Sprintf(_cntFolderSQL, folderHit(mid)), tp, mid)
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// AddFav add a favorite relation to mysql.
func (d *Dao) AddFav(c context.Context, fr *model.Favorite) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_addRelationSQL, relationHit(fr.Mid)), fr.Type, fr.Oid, fr.Mid, fr.Fid, fr.State, fr.CTime, fr.MTime, fr.State, fr.MTime)
	if err != nil {
		log.Error("db.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// DelFav delete a favorite relation to mysql.
func (d *Dao) DelFav(c context.Context, tp int8, mid, fid, oid int64) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_delRelationSQL, relationHit(mid)), tp, fid, oid)
	if err != nil {
		log.Error("d.db.Exec(%d,%d,%d,%d) error(%v)", mid, tp, fid, oid, err)
		return
	}
	return res.RowsAffected()
}

// AddRelation add a favorite relation to mysql.
func (d *Dao) AddRelation(c context.Context, fr *model.Favorite) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_addRelationSQL, relationHit(fr.Mid)), fr.Type, fr.Oid, fr.Mid, fr.Fid, fr.State, fr.CTime, fr.MTime, fr.State, fr.MTime)
	if err != nil {
		log.Error("db.Exec error(%v)", err)
		return
	}
	return res.LastInsertId()
}

// Relation get a relation from mysql.
func (d *Dao) Relation(c context.Context, tp int8, mid, fid, oid int64) (m *model.Favorite, err error) {
	m = &model.Favorite{}
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

// DelRelation delete a favorite relation to mysql.
func (d *Dao) DelRelation(c context.Context, tp int8, mid, fid, oid int64) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_delRelationSQL, relationHit(mid)), tp, fid, oid)
	if err != nil {
		log.Error("d.db.Exec(%d,%d,%d,%d) error(%v)", mid, tp, fid, oid, err)
		return
	}
	return res.RowsAffected()
}

// MultiDelRelations delete relations to mysql.
func (d *Dao) MultiDelRelations(c context.Context, tp int8, mid, fid int64, oids []int64) (rows int64, err error) {
	if len(oids) <= 0 {
		return
	}
	oidsStr := xstr.JoinInts(oids)
	res, err := d.db.Exec(c, fmt.Sprintf(_delRelationsSQL, relationHit(mid), oidsStr), tp, fid)
	if err != nil {
		log.Error("d.db.Exec(%d,%d,%d,%s) error(%v)", mid, tp, fid, oidsStr, err)
		return
	}
	return res.RowsAffected()
}

// TxMultiDelRelations delete relations to mysql.
func (d *Dao) TxMultiDelRelations(tx *sql.Tx, tp int8, mid, fid int64, oids []int64) (rows int64, err error) {
	if len(oids) <= 0 {
		return
	}
	oidsStr := xstr.JoinInts(oids)
	res, err := tx.Exec(fmt.Sprintf(_delRelationsSQL, relationHit(mid), oidsStr), tp, fid)
	if err != nil {
		log.Error("d.db.Exec(%d,%d,%d,%s) error(%v)", mid, tp, fid, oidsStr, err)
		return
	}
	return res.RowsAffected()
}

// MultiAddRelations insert relations to mysql.
func (d *Dao) MultiAddRelations(c context.Context, tp int8, mid, fid int64, oids []int64) (rows int64, err error) {
	var buffer bytes.Buffer
	valuesTpl := "(%d,%d,%d,%d,%d),"
	for _, oid := range oids {
		buffer.WriteString(fmt.Sprintf(valuesTpl, tp, oid, mid, fid, 0))
	}
	buffer.Truncate(buffer.Len() - 1)
	res, err := d.db.Exec(c, fmt.Sprintf(_maddRelationsSQL, relationHit(mid), buffer.String()))
	if err != nil {
		log.Error("d.db.Exec(%d,%d,%d,%s) error(%v)", mid, tp, fid, oids, err)
		return
	}
	return res.RowsAffected()
}

// DelFolder delete a folder to mysql.
func (d *Dao) DelFolder(c context.Context, tp int8, mid, fid int64) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_delFolderSQL, folderHit(mid)), tp, fid)
	if err != nil {
		log.Error("d.db.Exec(mid:%d,tp:%d,fid:%d) error(%v)", mid, tp, fid, err)
		return
	}
	return res.RowsAffected()
}

// FolderStats get folders from mysql.
func (d *Dao) FolderStats(c context.Context, fvmids []*model.ArgFVmid) (fs map[int64]*model.Folder, err error) {
	tblMap := make(map[int64][]int64, len(fvmids))
	for _, fvmid := range fvmids {
		tblMap[fvmid.Vmid] = append(tblMap[fvmid.Vmid], fvmid.Fid)
	}
	fs = make(map[int64]*model.Folder, len(fvmids))
	for vmid, fids := range tblMap {
		fidsStr := xstr.JoinInts(fids)
		var rows *sql.Rows
		if rows, err = d.db.Query(c, fmt.Sprintf(_folderStatSQL, folderHit(vmid), fidsStr)); err != nil {
			log.Error("d.db.Query(%s,%s,%s) error(%v)", _folderStatSQL, folderHit(vmid), fidsStr, err)
			return
		}
		defer rows.Close()
		for rows.Next() {
			stat := new(model.Folder)
			if err = rows.Scan(&stat.ID, &stat.PlayCount, &stat.FavedCount, &stat.ShareCount); err != nil {
				log.Error("rows.Scan error(%v)", err)
				return
			}
			fs[stat.ID*100+vmid%100] = stat
		}
		err = rows.Err()
		if err != nil {
			return
		}
	}
	return
}

// UserFolders get user's folders.
func (d *Dao) UserFolders(c context.Context, typ int8, mid int64) (fs map[int64]*model.Folder, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_userFoldersSQL, folderHit(mid)), typ, mid)
	if err != nil {
		log.Error("d.db.Query(mid:%d,typ:%d) error(%v)", mid, typ, err)
		return
	}
	defer rows.Close()
	fs = make(map[int64]*model.Folder)
	for rows.Next() {
		f := new(model.Folder)
		if err = rows.Scan(&f.ID, &f.Type, &f.Mid, &f.Name, &f.Cover, &f.Description, &f.Count, &f.Attr, &f.State, &f.CTime, &f.MTime); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		fs[f.ID] = f
	}
	err = rows.Err()
	return
}

// FolderSort return user's folders sort by type & mid.
func (d *Dao) FolderSort(c context.Context, typ int8, mid int64) (fst *model.FolderSort, err error) {
	var b []byte
	fst = new(model.FolderSort)
	row := d.db.QueryRow(c, _folderSortSQL, typ, mid)
	if err = row.Scan(&fst.ID, &fst.Type, &fst.Mid, &b, &fst.CTime, &fst.MTime); err != nil {
		if err == sql.ErrNoRows {
			fst = nil
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
		return
	}
	fst.SetIndex(b)
	return
}

// UpFolderSort update user's folder sort.
func (d *Dao) UpFolderSort(c context.Context, fst *model.FolderSort) (rows int64, err error) {
	sort := fst.Index()
	res, err := d.db.Exec(c, _setFolderSortSQL, fst.Type, fst.Mid, sort, fst.CTime, fst.MTime, sort, fst.MTime)
	if err != nil {
		log.Error("d.db.Exec() error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// RecentOids return user's three newest fav from a folder.
func (d *Dao) RecentOids(c context.Context, mid, fid int64, typ int8) (oids []int64, err error) {
	rows, err := d.dbRead.Query(c, fmt.Sprintf(_recentOidsSQL, relationHit(mid)), fid)
	if err != nil {
		log.Error("d.db.Query(%d,%d) error(%v)", mid, fid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var oid int64
		var otyp int8
		if err = rows.Scan(&oid, &otyp); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		if otyp == typ {
			oids = append(oids, oid)
		}
	}
	err = rows.Err()
	return
}

// RecentOids return user's three newest fav from a folder.
func (d *Dao) RecentRes(c context.Context, mid, fid int64) (res []*model.Resource, err error) {
	rows, err := d.dbRead.Query(c, fmt.Sprintf(_recentOidsSQL, relationHit(mid)), fid)
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
		res = append(res, &model.Resource{Oid: oid, Typ: typ})
	}
	err = rows.Err()
	return
}

// TxCopyRelations copy resources from oldfid to newfid by oids.
func (d *Dao) TxCopyRelations(tx *sql.Tx, typ int8, oldmid, mid, oldfid, newfid int64, oids []int64) (rows int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_copyRelationsSQL, relationHit(mid), typ, mid, newfid, relationHit(oldmid), xstr.JoinInts(oids)), typ, oldmid, oldfid)
	if err != nil {
		log.Error("db.Exec() error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// CopyRelations copy resources from oldfid to newfid by oids.
func (d *Dao) CopyRelations(c context.Context, typ int8, oldmid, mid, oldfid, newfid int64, oids []int64) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_copyRelationsSQL, relationHit(mid), typ, mid, newfid, relationHit(oldmid), xstr.JoinInts(oids)), typ, oldmid, oldfid)
	if err != nil {
		log.Error("db.Exec() error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// CntUsers get favoried users count from mysql.
func (d *Dao) CntUsers(c context.Context, typ int8, oid int64) (count int, err error) {
	row := d.dbPush.QueryRow(c, fmt.Sprintf(_cntUsersSQL, usersHit(oid)), typ, oid)
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// Users get favorite users from mysql.
func (d *Dao) Users(c context.Context, typ int8, oid int64, start, end int) (us []*model.User, err error) {
	rows, err := d.dbPush.Query(c, fmt.Sprintf(_usersSQL, usersHit(oid)), typ, oid, start, end)
	if err != nil {
		log.Error("d.db.Query(%d,%d,%d,%d) error(%v)", typ, oid, start, end, err)
		return
	}
	defer rows.Close()
	us = make([]*model.User, 0)
	for rows.Next() {
		var u = new(model.User)
		if err = rows.Scan(&u.ID, &u.Type, &u.Oid, &u.Mid, &u.State, &u.CTime, &u.MTime); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		us = append(us, u)
	}
	err = rows.Err()
	return
}

// OidCount get oid's fav count from mysql.
func (d *Dao) OidCount(c context.Context, typ int8, oid int64) (count int64, err error) {
	row := d.dbPush.QueryRow(c, fmt.Sprintf(_countSQL, countHit(oid)), typ, oid)
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// OidsCount get oids's count from mysql.
func (d *Dao) OidsCount(c context.Context, typ int8, oids []int64) (counts map[int64]int64, err error) {
	tblMap := make(map[string][]int64, len(oids))
	for _, oid := range oids {
		countHit := countHit(oid)
		tblMap[countHit] = append(tblMap[countHit], oid)
	}
	counts = make(map[int64]int64, len(oids))
	for countHit, oids := range tblMap {
		oidsStr := xstr.JoinInts(oids)
		var rows *sql.Rows
		if rows, err = d.db.Query(c, fmt.Sprintf(_countsSQL, countHit, oidsStr), typ); err != nil {
			log.Error("d.db.Query(%s,%d) error(%v)", fmt.Sprintf(_countsSQL, countHit, oidsStr), typ, err)
			return
		}
		defer rows.Close()
		for rows.Next() {
			var oid, count int64
			if err = rows.Scan(&oid, &count); err != nil {
				log.Error("rows.Scan error(%v)", err)
				return
			}
			counts[oid] = count
		}
		if err = rows.Err(); err != nil {
			log.Error("rows.Err() error(%v)", err)
			return
		}
	}
	return
}

// BatchOids get batch oids from mysql.
func (d *Dao) BatchOids(c context.Context, typ int8, mid int64, limit int) (oids []int64, err error) {
	rows, err := d.dbRead.Query(c, fmt.Sprintf(_batchOidsSQL, relationHit(mid)), typ, mid, limit)
	if err != nil {
		log.Error("d.db.Query(%d,%d,%d) error(%v)", typ, mid, limit, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var oid int64
		if err = rows.Scan(&oid); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		oids = append(oids, oid)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}
