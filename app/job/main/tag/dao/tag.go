package dao

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"

	"go-common/app/job/main/tag/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_tagsByNameSQL  = "SELECT id,name,type,verify,attr,state FROM tag WHERE name in (%s)"
	_tagSQL         = "SELECT id,name,type,verify,attr,state FROM tag WHERE id in (%s)"
	_addTagSQL      = "INSERT INTO tag (name,type,verify,attr,state) VALUES (?,?,?,?,?)"
	_addTagCountSQL = "INSERT INTO `count` (tid,bind,sub) VALUES (?,0,0)"
)

// TagByNames tag by names.
func (d *Dao) TagByNames(c context.Context, names []string) (res []*model.Tag, missed []string, err error) {
	res = make([]*model.Tag, 0, len(names))
	tagMap := make(map[string]struct{}, len(names))
	var (
		buf  = bytes.NewBuffer(nil)
		args = make([]interface{}, len(names))
	)
	for i, name := range names {
		tagMap[name] = struct{}{}
		buf.WriteByte('?')
		buf.WriteByte(',')
		args[i] = name
	}
	buf.Truncate(buf.Len() - 1)
	rows, err := d.platform.Query(c, fmt.Sprintf(_tagsByNameSQL, buf.String()), args...)
	if err != nil {
		log.Error("d.TagByNames(%v) error(%v)", names, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		t := &model.Tag{}
		if err = rows.Scan(&t.ID, &t.Name, &t.Type, &t.Verify, &t.Attr, &t.State); err != nil {
			log.Error("d.TagByNames(%v) scan() error(%v)", names, err)
			return
		}
		if _, ok := tagMap[t.Name]; ok {
			delete(tagMap, t.Name)
		}
		res = append(res, t)
	}
	err = rows.Err()
	missed = make([]string, 0, len(tagMap))
	for name := range tagMap {
		missed = append(missed, name)
	}
	return
}

// Tags tag by ids.
func (d *Dao) Tags(c context.Context, ids []int64) (res []*model.Tag, err error) {
	res = make([]*model.Tag, 0, len(ids))
	rows, err := d.platform.Query(c, fmt.Sprintf(_tagSQL, xstr.JoinInts(ids)))
	if err != nil {
		log.Error("d.Tags(%v) error(%v)", ids, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		t := &model.Tag{}
		if err = rows.Scan(&t.ID, &t.Name, &t.Type, &t.Verify, &t.Attr, &t.State); err != nil {
			log.Error("d.Tags(%v) scan() error(%v)", ids, err)
			return
		}
		t.Name = replace(t.Name)
		res = append(res, t)
	}
	err = rows.Err()
	return
}

// InsertTags insert tags.
func (d *Dao) InsertTags(c context.Context, tags []*model.Tag) (err error) {
	var (
		tx  *xsql.Tx
		res sql.Result
	)
	if tx, err = d.BeginTran(c); err != nil {
		log.Error("d.InsertTags() BeginTran(%v) error(%v)", tags, err)
		return
	}
	for _, t := range tags {
		if res, err = tx.Exec(_addTagSQL, t.Name, t.Type, t.Verify, t.Attr, t.State); err != nil {
			log.Error("d.InsertTags(%v) add tag(%v) error(%v)", tags, t, err)
			tx.Rollback()
			return
		}
		if t.ID, err = res.LastInsertId(); err != nil || t.ID <= 0 {
			log.Error("d.InsertTags(%v) add tag(%v) get id <=0 || error(%v)", tags, t, err)
			tx.Rollback()
			return
		}
		if _, err = tx.Exec(_addTagCountSQL, t.ID); err != nil {
			log.Error("d.InsertTags(%v) add tag_count(%v) error(%v)", tags, t, err)
			tx.Rollback()
			return
		}
	}
	if err = tx.Commit(); err != nil {
		log.Error("d.InsertTags(%v) Commit() error(%v)", tags, err)
		tx.Rollback()
	}
	return
}
