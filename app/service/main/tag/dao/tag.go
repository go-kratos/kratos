package dao

import (
	"bytes"
	"context"
	xsql "database/sql"
	"fmt"
	"strings"

	"go-common/app/service/main/tag/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

func replace(name string) string {
	var (
		there bool
		rb    []byte
	)
	sb := []byte(strings.Trim(name, " "))
	for _, b := range sb {
		if b < 0x20 || b == 0x7f {
			there = true
			continue
		}
		rb = append(rb, b)
	}
	if there {
		log.Warn("There are invisible characters,tag byte(%v),name(%v)", sb, name)
	}
	return string(rb)
}

var _tagSQL = "SELECT id,type,name,cover,head_cover,content,short_content,attr,state,ctime,mtime FROM tag WHERE id=?"

// Tag return a tag by tid from msyql
func (d *Dao) Tag(c context.Context, tid int64) (t *model.Tag, err error) {
	t = new(model.Tag)
	row := d.db.QueryRow(c, _tagSQL, tid)
	if err = row.Scan(&t.ID, &t.Type, &t.Name, &t.Cover, &t.HeadCover, &t.Content, &t.ShortContent, &t.Attr, &t.State, &t.CTime, &t.MTime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan() error(%v)", err)
		}
		t = nil
		return
	}
	if t.State == model.TagStateDelete {
		log.Warn("TagStateDelete:%+v", t)
	}
	t.Name = replace(t.Name)
	return
}

var _tagsSQL = "SELECT id,type,name,cover,head_cover,content,short_content,attr,state,ctime,mtime FROM tag WHERE id IN (%s)"

// Tags return tags by tids from mysql.
func (d *Dao) Tags(c context.Context, tids []int64) (tags []*model.Tag, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_tagsSQL, xstr.JoinInts(tids)))
	if err != nil {
		log.Error("db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		t := new(model.Tag)
		if err = rows.Scan(&t.ID, &t.Type, &t.Name, &t.Cover, &t.HeadCover, &t.Content, &t.ShortContent, &t.Attr, &t.State, &t.CTime, &t.MTime); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			tags = nil
			return
		}
		if t.State == model.TagStateDelete {
			log.Warn("TagStateDelete:%+v", t)
		}
		t.Name = replace(t.Name)
		tags = append(tags, t)
	}
	return
}

var _tagNameSQL = "SELECT id,type,name,cover,head_cover,content,short_content,attr,state,ctime,mtime FROM tag WHERE name=?"

// TagByName return a tag by name from mysql.
func (d *Dao) TagByName(c context.Context, name string) (t *model.Tag, err error) {
	t = new(model.Tag)
	row := d.db.QueryRow(c, _tagNameSQL, name)
	if err = row.Scan(&t.ID, &t.Type, &t.Name, &t.Cover, &t.HeadCover, &t.Content, &t.ShortContent, &t.Attr, &t.State, &t.CTime, &t.MTime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan() error(%v)", err)
		}
		t = nil
		return
	}
	if t.State == model.TagStateDelete {
		log.Warn("TagStateDelete:%+v", t)
	}
	t.Name = replace(t.Name)
	return
}

var _tagsNamesSQL = "SELECT id,type,name,cover,head_cover,content,short_content,attr,state,ctime,mtime FROM tag WHERE name IN (%s)"

// TagsByNames return tags by names from mysql.
func (d *Dao) TagsByNames(c context.Context, names []string) (tags []*model.Tag, err error) {
	var (
		buf  = bytes.NewBuffer(nil)
		args = make([]interface{}, len(names))
	)
	for i, name := range names {
		buf.WriteByte('?')
		buf.WriteByte(',')
		args[i] = name
	}
	buf.Truncate(buf.Len() - 1)
	rows, err := d.db.Query(c, fmt.Sprintf(_tagsNamesSQL, buf.String()), args...)
	if err != nil {
		log.Error("d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		t := new(model.Tag)
		if err = rows.Scan(&t.ID, &t.Type, &t.Name, &t.Cover, &t.HeadCover, &t.Content, &t.ShortContent, &t.Attr, &t.State, &t.CTime, &t.MTime); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		if t.State == model.TagStateDelete {
			log.Warn("TagStateDelete:%+v", t)
		}
		t.Name = replace(t.Name)
		tags = append(tags, t)
	}
	return
}

var _addTag = "INSERT INTO tag (name,cover,head_cover,content,short_content,type,verify,attr,state) VALUES (?,?,?,?,?,?,0,0,0)"

// TxAddTag add a tag into mysql.
func (d *Dao) TxAddTag(tx *sql.Tx, t *model.Tag) (id int64, err error) {
	res, err := tx.Exec(_addTag, t.Name, t.Cover, t.HeadCover, t.Content, t.ShortContent, t.Type)
	if err != nil {
		log.Error("tx.Exec(%v) error(%v)", t, err)
		return
	}
	return res.LastInsertId()
}

var _addTagCountSQL = "INSERT INTO count (tid,bind,sub) VALUES (?,0,0)"

// CreateTag .
func (d *Dao) CreateTag(c context.Context, t *model.Tag) (err error) {
	var (
		tx  *sql.Tx
		res xsql.Result
	)
	tx, err = d.BeginTran(c)
	if err != nil {
		log.Error("tx.BeginTran error(%v)", err)
		return
	}
	res, err = tx.Exec(_addTag, t.Name, t.Cover, t.HeadCover, t.Content, t.ShortContent, t.Type)
	if err != nil {
		log.Error("tx.Exec(%v) error(%v)", t, err)
		tx.Rollback()
		return
	}
	t.ID, err = res.LastInsertId()
	if err != nil || t.ID == 0 {
		log.Error("tx.Exec(%v)  error(%v)", t, err)
		tx.Rollback()
		return
	}
	_, err = tx.Exec(_addTagCountSQL, t.ID)
	if err != nil {
		log.Error("tx.Exec(%v)  error(%v)", t, err)
		tx.Rollback()
		return
	}
	if err = tx.Commit(); err != nil {
		tx.Rollback()
	}
	return
}

// CreateTags .
func (d *Dao) CreateTags(c context.Context, ts []*model.Tag) (err error) {
	var (
		tx  *sql.Tx
		res xsql.Result
	)
	tx, err = d.BeginTran(c)
	if err != nil {
		log.Error("tx.BeginTran error(%v)", err)
		return
	}
	for _, t := range ts {
		res, err = tx.Exec(_addTag, t.Name, t.Cover, t.HeadCover, t.Content, t.ShortContent, t.Type)
		if err != nil {
			log.Error("tx.Exec(%v) error(%v)", t, err)
			tx.Rollback()
			return
		}
		t.ID, err = res.LastInsertId()
		if err != nil || t.ID == 0 {
			log.Error("tx.Exec(%v) error(%v)", t, err)
			tx.Rollback()
			return
		}
		_, err = tx.Exec(_addTagCountSQL, t.ID)
		if err != nil {
			log.Error("tx.Exec(%v) error(%v)", t, err)
			tx.Rollback()
			return
		}
	}
	if err = tx.Commit(); err != nil {
		tx.Rollback()
	}
	return
}

var _countSelectSQL = "select tid,bind,sub from platform_tag.count where tid = ?;"

// Count .
func (d *Dao) Count(c context.Context, tid int64) (res *model.Count, err error) {
	row := d.db.QueryRow(c, _countSelectSQL, tid)
	res = &model.Count{}
	if err = row.Scan(&res.Tid, &res.Bind, &res.Sub); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		log.Error("rows.Scan() err:%v", err)
	}
	return
}

var _countsSelectSQL = "select tid,bind,sub from platform_tag.count where tid in (%s);"

// Counts .
func (d *Dao) Counts(c context.Context, tids []int64) (res []*model.Count, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_countsSelectSQL, xstr.JoinInts(tids)))
	if err != nil {
		log.Error("d.db.Query(%v),err:%v", tids, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.Count{}
		if err = rows.Scan(&r.Tid, &r.Bind, &r.Sub); err != nil {
			log.Error("rows.Scan() err:%v", err)
			continue
		}
		if r.Tid != 0 {
			res = append(res, r)
		}
	}
	return
}

var _upTagBindCount = "INSERT IGNORE INTO count (tid,bind) VALUES (?,1) ON DUPLICATE KEY UPDATE bind=bind+?"

// TxUpTagBindCount update tag bind count into mysql.
func (d *Dao) TxUpTagBindCount(tx *sql.Tx, tid, count int64) (rows int64, err error) {
	res, err := tx.Exec(_upTagBindCount, tid, count)
	if err != nil {
		log.Error("tx.Exec(%d) error(%v)", tid, err)
		return
	}
	return res.RowsAffected()
}

var _upTagSubCount = "INSERT IGNORE INTO count (tid,sub) VALUES (?,1) ON DUPLICATE KEY UPDATE sub=sub+?"

// TxUpTagSubCount update tag sub count into mysql.
func (d *Dao) TxUpTagSubCount(tx *sql.Tx, tid, count int64) (rows int64, err error) {
	res, err := tx.Exec(_upTagSubCount, tid, count)
	if err != nil {
		log.Error("tx.Exec(%d) error(%v)", tid, err)
		return
	}
	return res.RowsAffected()
}

var _selecHoTagsSQL = "SELECT rid,tid,tname,highlight FROM rank_result WHERE rid=? AND type=? ORDER BY rank"

// Hots .
func (d *Dao) Hots(c context.Context, rid, hotType int64) (hots []*model.HotTag, err error) {
	rows, err := d.db.Query(c, _selecHoTagsSQL, rid, hotType)
	if err != nil {
		log.Error("d.db.Query(%d,%d) error(%v)", rid, hotType, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		hot := &model.HotTag{}
		if err = rows.Scan(&hot.Rid, &hot.Tid, &hot.Tname, &hot.HighLight); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			hots = nil
			return
		}
		hots = append(hots, hot)
	}
	return
}

var _selectRidsSQL = "SELECT DISTINCT prid,rid FROM rank_result"

// Rids .
func (d *Dao) Rids(c context.Context) (cm map[int64][]int64, pridMap map[int64]int64, rids []int64, err error) {
	rows, err := d.db.Query(c, _selectRidsSQL)
	if err != nil {
		log.Error("d.db.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	pridMap = make(map[int64]int64)
	cm = make(map[int64][]int64)
	for rows.Next() {
		var prid, rid int64
		if err = rows.Scan(&prid, &rid); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			cm = nil
			return
		}
		cm[prid] = append(cm[prid], rid)
		pridMap[rid] = prid
		rids = append(rids, rid)
	}
	return
}

var _selectHotSQL = "SELECT rid,tid FROM rank_result WHERE type=0 ORDER BY rank"

// HotMap .
func (d *Dao) HotMap(c context.Context) (res map[int16][]int64, err error) {
	rows, err := d.db.Query(c, _selectHotSQL)
	if err != nil {
		log.Error("Hot d.hot.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	res = make(map[int16][]int64)
	for rows.Next() {
		var rid int16
		var tid int64
		if err = rows.Scan(&rid, &tid); err != nil {
			log.Error("Hot rows.Scan() error(%v)", err)
			res = nil
			return
		}
		res[rid] = append(res[rid], tid)
	}
	return
}

var _selectPridsSQL = "SELECT DISTINCT prid FROM rank_result"

// Prids .
func (d *Dao) Prids(c context.Context) (res []int16, err error) {
	rows, err := d.db.Query(c, _selectPridsSQL)
	if err != nil {
		log.Error("Rids d.rids.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var prid int16
		if err = rows.Scan(&prid); err != nil {
			log.Error("Rids rows.Scan() error(%v)", err)
			res = nil
			return
		}
		res = append(res, prid)
	}
	return
}

var _tagGroupSQL = "SELECT ptid,tid FROM `group` WHERE ptid !=0 LIMIT 5000"

// TagGroup .
func (d *Dao) TagGroup(c context.Context) (resMap map[int64][]int64, err error) {
	rows, err := d.db.Query(c, _tagGroupSQL)
	if err != nil {
		log.Error("d.db.Query(%s) error(%v)", _tagGroupSQL, err)
		return
	}
	defer rows.Close()
	resMap = make(map[int64][]int64)
	for rows.Next() {
		var ptid, tid int64
		if err = rows.Scan(&ptid, &tid); err != nil {
			if err == sql.ErrNoRows {
				err = nil
			} else {
				log.Error("rows.Scan() error(%v)", err)
			}
			return
		}
		resMap[ptid] = append(resMap[ptid], tid)
	}
	return
}

var _tagFilterSQL = "SELECT rid,tid,tname,rank FROM rank_filter WHERE type=1"

//RecommandTagFilter recommend tag filter
func (d *Dao) RecommandTagFilter(c context.Context) (ts []*model.UploadTag, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _tagFilterSQL); err != nil {
		log.Error(" d.db.Query(c,_tagFilter) error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		t := &model.UploadTag{}
		if err = rows.Scan(&t.Rid, &t.Tid, &t.Tname, &t.Rank); err != nil {
			if err == sql.ErrNoRows {
				err = nil
			} else {
				log.Error("rows.Scan(_tagFilter) error(%v)", err)
			}
			return
		}
		ts = append(ts, t)
	}

	return
}

var _tagTopSQL = "SELECT rid,tid,tname,is_business,rank FROM rank_top WHERE type=1"

//RecommandTagTop recommand tag top
func (d *Dao) RecommandTagTop(c context.Context) (ts []*model.UploadTag, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _tagTopSQL); err != nil {
		log.Error(" d.db.Query(c,_tagTop) error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		t := &model.UploadTag{}
		if err = rows.Scan(&t.Rid, &t.Tid, &t.Tname, &t.IsBusiness, &t.Rank); err != nil {
			if err == sql.ErrNoRows {
				err = nil
			} else {
				log.Error("rows.Scan(_tagTop) error(%v)", err)
			}
			return
		}
		ts = append(ts, t)
	}
	return
}

var _upTagStateSQL = "UPDATE tag SET state=? WHERE id=? AND state!=?"

// UpTagState  .
func (d *Dao) UpTagState(c context.Context, tid int64, state int32) (affected int64, err error) {
	res, err := d.db.Exec(c, _upTagStateSQL, state, tid, state)
	if err != nil {
		log.Error("UpTagState d.Exec(%d,%v) error(%v)", tid, state, err)
		return
	}
	return res.RowsAffected()
}
