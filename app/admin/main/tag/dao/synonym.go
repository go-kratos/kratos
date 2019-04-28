package dao

import (
	"context"
	"fmt"
	"go-common/app/admin/main/tag/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
	"strings"
)

const (
	_synonymsSQL       = "SELECT g.id,g.ptid,g.tid,g.uname,g.ctime,g.mtime,t.name FROM `group` as g LEFT OUTER JOIN tag as t ON t.id = g.tid WHERE g.ptid = 0 %s ORDER BY g.mtime DESC LIMIT ?,?"
	_synonymIDsSQL     = "SELECT g.id,g.ptid,g.tid,g.uname,g.ctime,g.mtime,t.name FROM `group` as g LEFT OUTER JOIN tag as t ON t.id = g.tid WHERE ptid IN (%s)"
	_synonymCountSQL   = "SELECT count(*) FROM `group` as g LEFT OUTER JOIN tag as t ON t.id = g.tid WHERE g.ptid = 0 %s"
	_synonymByNameSQL  = "SELECT g.id,g.ptid,g.tid,g.uname,g.ctime,g.mtime,t.name FROM `group` as g LEFT OUTER JOIN tag as t ON t.id = g.tid WHERE t.name=? AND g.ptid = 0"
	_insertSynonymSQL  = "INSERT INTO platform_tag.group(ptid,tid,uname) VALUES (?,?,?)"
	_delSynonymSonsSQL = "DELETE FROM `group` WHERE tid in (%s) AND ptid = ?"
	_delSynonymSQL     = "DELETE FROM `group` WHERE tid=? AND ptid = 0"
	_insertSynonymsSQL = "INSERT INTO platform_tag.group(ptid,tid,uname) VALUES %s "
	_synonymSQL        = "SELECT g.id,g.ptid,g.tid,g.uname,g.ctime,g.mtime FROM `group` as g WHERE tid =? OR ptid=?"
)

func synonymLike(keyWord string) string {
	if keyWord == "" || len(keyWord) == 0 {
		return ""
	}
	return "AND t.name LIKE \"%" + keyWord + "%\""
}

// SynonymCount count synonym.
func (d *Dao) SynonymCount(c context.Context, keyWord string) (count int64, err error) {
	row := d.db.QueryRow(c, fmt.Sprintf(_synonymCountSQL, synonymLike(keyWord)))
	if err = row.Scan(&count); err != nil {
		log.Error("d.db.QueryRow() error(%v)", err)
	}
	return
}

// Synonyms Synonyms.
func (d *Dao) Synonyms(c context.Context, keyWord string, start, end int32) (stagMap map[int64]*model.SynonymTag, stag []*model.SynonymTag, ids []int64, err error) {
	stagMap = make(map[int64]*model.SynonymTag)
	rows, err := d.db.Query(c, fmt.Sprintf(_synonymsSQL, synonymLike(keyWord)), start, end)
	if err != nil {
		log.Error("query synonyms(%s) error(%v)", keyWord, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		st := &model.SynonymTag{}
		if err = rows.Scan(&st.ID, &st.Ptid, &st.Tid, &st.UName, &st.CTime, &st.MTime, &st.TName); err != nil {
			log.Error("rows.Scan(model.SynonymTag{}) error(%v)", err)
			err = nil
			continue
		}
		ids = append(ids, st.Tid)
		stag = append(stag, st)
		stagMap[st.Tid] = st
	}
	return
}

// SynonymIDs get Synonym By IDs.
func (d *Dao) SynonymIDs(c context.Context, ids []int64) (mapST map[int64]*model.SynonymTag, stag []*model.SynonymTag, sids []int64, err error) {
	mapST = make(map[int64]*model.SynonymTag)
	rows, err := d.db.Query(c, fmt.Sprintf(_synonymIDsSQL, xstr.JoinInts(ids)))
	if err != nil {
		log.Error("d.db.Query(%v) error(%v)", ids, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		st := &model.SynonymTag{}
		if err = rows.Scan(&st.ID, &st.Ptid, &st.Tid, &st.UName, &st.CTime, &st.MTime, &st.TName); err != nil {
			log.Error("rows.Scan(model.SynonymTag{}) error(%v)", err)
			return
		}
		stag = append(stag, st)
		sids = append(sids, st.Tid)
		mapST[st.Tid] = st
	}
	return
}

// SynonymByName get one synonym info by name.
func (d *Dao) SynonymByName(c context.Context, tname string) (res *model.SynonymTag, err error) {
	res = new(model.SynonymTag)
	row := d.db.QueryRow(c, _synonymByNameSQL, tname)
	if err = row.Scan(&res.ID, &res.Ptid, &res.Tid, &res.UName, &res.CTime, &res.MTime, &res.TName); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			res = nil
		} else {
			log.Error("row.Scan OneSynonymInfoByName(%s) error(%v)", tname, err)
		}
	}
	return
}

// InsertSynonym insert synonym.
func (d *Dao) InsertSynonym(c context.Context, uname string, ptid, tid int64) (id int64, err error) {
	res, err := d.db.Exec(c, _insertSynonymSQL, ptid, tid, uname)
	if err != nil {
		log.Error("d.db.Exec(%d,%d,%s) error(%v)", ptid, tid, uname, err)
		return
	}
	return res.LastInsertId()
}

// DelSynonym delete synonym.
func (d *Dao) DelSynonym(c context.Context, tid int64) (affect int64, err error) {
	res, err := d.db.Exec(c, _delSynonymSQL, tid)
	if err != nil {
		log.Error("d.db.Exec(%d) error(%v)", tid, err)
		return
	}
	return res.RowsAffected()
}

// DelSynonymSon 删除二级子类数据.
func (d *Dao) DelSynonymSon(c context.Context, ptid int64, tids []int64) (affect int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_delSynonymSonsSQL, xstr.JoinInts(tids)), ptid)
	if err != nil {
		log.Error("d.db.Exec(%d,%v) error(%v)", ptid, tids, err)
		return
	}
	return res.RowsAffected()
}

// InsertSynonyms insert synonyms.
func (d *Dao) InsertSynonyms(c context.Context, uname string, ptid int64, tids []int64) (id int64, err error) {
	var (
		sql    []string
		sqlStr = " (%d,%d,%q) "
	)
	for _, tid := range tids {
		s := fmt.Sprintf(sqlStr, ptid, tid, uname)
		sql = append(sql, s)
	}
	insertSQL := strings.Join(sql, " , ")
	res, err := d.db.Exec(c, fmt.Sprintf(_insertSynonymsSQL, insertSQL))
	if err != nil {
		log.Error("d.db.Exec(%v) error(%v)", sql, err)
		return
	}
	return res.LastInsertId()
}

// Synonym Synonym Info By ID.
func (d *Dao) Synonym(c context.Context, tid int64) (res []*model.SynonymTag, err error) {
	rows, err := d.db.Query(c, _synonymSQL, tid, tid)
	if err != nil {
		log.Error("d.db.Query(%d) error(%v)", tid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		st := &model.SynonymTag{}
		if err = rows.Scan(&st.ID, &st.Ptid, &st.Tid, &st.UName, &st.CTime, &st.MTime); err != nil {
			log.Error("rows.Scan(model.SynonymTag{}) error(%v)", err)
			return
		}
		res = append(res, st)
	}
	return
}
