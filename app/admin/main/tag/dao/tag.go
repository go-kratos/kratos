package dao

import (
	"bytes"
	"context"
	"fmt"

	"go-common/app/admin/main/tag/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_tagSQL            = "SELECT id,type,name,cover,content,verify,attr,state,ctime,mtime,head_cover,short_content FROM tag WHERE id=?"
	_tagsSQL           = "SELECT id,type,name,cover,content,verify,attr,state,ctime,mtime,head_cover,short_content FROM tag WHERE id in (%s)"
	_tagNamesSQL       = "SELECT id,type,name,cover,content,verify,attr,state,ctime,mtime,head_cover,short_content FROM tag WHERE name in (%s)"
	_tagNameSQL        = "SELECT id,type,name,cover,content,verify,attr,state,ctime,mtime,head_cover,short_content FROM tag WHERE name=?"
	_insertTagSQL      = "INSERT INTO tag(name,cover,content,type,verify,attr,state,head_cover,short_content) VALUES (?,?,?,?,?,?,?,?,?)"
	_insertTagCountSQL = "INSERT INTO platform_tag.count(tid) VALUES (?)"
	_updateTagSQL      = "UPDATE tag t SET t.cover=?, t.content=?, t.type=?, t.verify=?, t.attr=?, t.state=?, t.head_cover=?, t.short_content=? WHERE id=?"
	_upVerifySQL       = "UPDATE tag SET verify=? WHERE id=?"
	_upTagStateSQL     = "UPDATE tag SET state=? WHERE id=?"
	_tagCountSQL       = "SELECT tid,bind,sub FROM platform_tag.count WHERE tid=?"
	_tagCountsSQL      = "SELECT tid,bind,sub FROM platform_tag.count WHERE tid in (%s)"
)

// Tag query tag by id.
func (d *Dao) Tag(c context.Context, tid int64) (tag *model.Tag, err error) {
	tag = new(model.Tag)
	row := d.db.QueryRow(c, _tagSQL, tid)
	if err = row.Scan(&tag.ID, &tag.Type, &tag.Name, &tag.Cover, &tag.Content, &tag.Verify, &tag.Attr, &tag.State, &tag.CTime, &tag.MTime, &tag.HeadCover, &tag.ShortContent); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Error("row.Scan() error(%v)", err)
	}
	return
}

// Tags batch query tag by tids.
func (d *Dao) Tags(c context.Context, tids []int64) (tags []*model.Tag, tagMap map[int64]*model.Tag, err error) {
	tagMap = make(map[int64]*model.Tag)
	rows, err := d.db.Query(c, fmt.Sprintf(_tagsSQL, xstr.JoinInts(tids)))
	if err != nil {
		log.Error("db query tags(%v) error(%v)", tids, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		tag := &model.Tag{}
		if err = rows.Scan(&tag.ID, &tag.Type, &tag.Name, &tag.Cover, &tag.Content, &tag.Verify, &tag.Attr, &tag.State, &tag.CTime, &tag.MTime, &tag.HeadCover, &tag.ShortContent); err != nil {
			log.Error("rows.Scan(model.tag{}) error(%v)", err)
			return
		}
		tags = append(tags, tag)
		tagMap[tag.ID] = tag
	}
	return
}

// TagByName query tag by tname.
func (d *Dao) TagByName(c context.Context, tname string) (tag *model.Tag, err error) {
	tag = new(model.Tag)
	row := d.db.QueryRow(c, _tagNameSQL, tname)
	if err = row.Scan(&tag.ID, &tag.Type, &tag.Name, &tag.Cover, &tag.Content, &tag.Verify, &tag.Attr, &tag.State, &tag.CTime, &tag.MTime, &tag.HeadCover, &tag.ShortContent); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Error("row.Scan() error(%v)", err)
	}
	return
}

// TagByNames query tag by tnames.
func (d *Dao) TagByNames(c context.Context, tnames []string) (tags []*model.Tag, tagMap map[int64]*model.Tag, tagNameMap map[string]*model.Tag, err error) {
	tagMap = make(map[int64]*model.Tag)
	tagNameMap = make(map[string]*model.Tag)
	var (
		buf  = bytes.NewBuffer(nil)
		args = make([]interface{}, len(tnames))
	)
	for i, name := range tnames {
		buf.WriteByte('?')
		buf.WriteByte(',')
		args[i] = name
	}
	buf.Truncate(buf.Len() - 1)
	rows, err := d.db.Query(c, fmt.Sprintf(_tagNamesSQL, buf.String()), args...)
	if err != nil {
		log.Error("db query tags(%s) error(%v)", tnames, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		tag := &model.Tag{}
		if err = rows.Scan(&tag.ID, &tag.Type, &tag.Name, &tag.Cover, &tag.Content, &tag.Verify, &tag.Attr, &tag.State, &tag.CTime, &tag.MTime, &tag.HeadCover, &tag.ShortContent); err != nil {
			log.Error("rows.Scan(model.tag{}) error(%v)", err)
			return
		}
		tags = append(tags, tag)
		tagMap[tag.ID] = tag
		tagNameMap[tag.Name] = tag
	}
	return
}

// UpVerifyState update tag verify state by tid.
func (d *Dao) UpVerifyState(c context.Context, tid int64, verifyState int32) (affect int64, err error) {
	res, err := d.db.Exec(c, _upVerifySQL, verifyState, tid)
	if err != nil {
		log.Error("db update verify(%d,%d) error(%v)", tid, verifyState, err)
		return
	}
	return res.RowsAffected()
}

// TxInsertTag transaction insert tag.
func (d *Dao) TxInsertTag(tx *sql.Tx, tag *model.Tag) (id int64, err error) {
	res, err := tx.Exec(_insertTagSQL, tag.Name, tag.Cover, tag.Content, tag.Type, tag.Verify, tag.Attr, tag.State, tag.HeadCover, tag.ShortContent)
	if err != nil {
		log.Error("db insert tag(%+v) error(%v)", tag, err)
		return
	}
	return res.LastInsertId()
}

// TxInsertTagCount transaction insert tag_count.
func (d *Dao) TxInsertTagCount(tx *sql.Tx, tid int64) (id int64, err error) {
	res, err := tx.Exec(_insertTagCountSQL, tid)
	if err != nil {
		log.Error("db insert tag_count(%d) error(%v)", tid, err)
		return
	}
	return res.LastInsertId()
}

// UpdateTag update tag.
func (d *Dao) UpdateTag(c context.Context, tag *model.Tag) (affect int64, err error) {
	res, err := d.db.Exec(c, _updateTagSQL, tag.Cover, tag.Content, tag.Type, tag.Verify, tag.Attr, tag.State, tag.HeadCover, tag.ShortContent, tag.ID)
	if err != nil {
		log.Error("db update tag(%+v) error(%v)", tag, err)
		return
	}
	return res.RowsAffected()
}

// TxUpdateTag update tag.
func (d *Dao) TxUpdateTag(tx *sql.Tx, tag *model.Tag) (affect int64, err error) {
	res, err := tx.Exec(_updateTagSQL, tag.Cover, tag.Content, tag.Type, tag.Verify, tag.Attr, tag.State, tag.HeadCover, tag.ShortContent, tag.ID)
	if err != nil {
		log.Error("db update tag(%+v) error(%v)", tag, err)
		return
	}
	return res.RowsAffected()
}

// UpTagState update tag state.
func (d *Dao) UpTagState(c context.Context, tid int64, state int32) (affect int64, err error) {
	res, err := d.db.Exec(c, _upTagStateSQL, state, tid)
	if err != nil {
		log.Error("db update tag state(%d,%d) error(%v)", tid, state, err)
		return
	}
	return res.RowsAffected()
}

// TagCount query tag count by tid.
func (d *Dao) TagCount(c context.Context, tid int64) (res *model.TagCount, err error) {
	res = new(model.TagCount)
	row := d.db.QueryRow(c, _tagCountSQL, tid)
	if err = row.Scan(&res.Tid, &res.Bind, &res.Sub); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			res = nil
			return
		}
		log.Error("row.Scan() error(%v)", err)
	}
	return
}

// TagCounts batch query tag count.
func (d *Dao) TagCounts(c context.Context, tids []int64) (res map[int64]*model.TagCount, err error) {
	res = make(map[int64]*model.TagCount)
	rows, err := d.db.Query(c, fmt.Sprintf(_tagCountsSQL, xstr.JoinInts(tids)))
	if err != nil {
		log.Error("db tag counts(%v) error(%v)", tids, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		c := &model.TagCount{}
		if err = rows.Scan(&c.Tid, &c.Bind, &c.Sub); err != nil {
			log.Error("rows.Scan(TagCount{}) error(%v)", err)
			return
		}
		res[c.Tid] = c
	}
	return
}
