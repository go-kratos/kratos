package template

import (
	"context"

	"go-common/app/interface/main/creative/model/template"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	// insert
	_addTplSQL = "INSERT INTO archive_templates (mid,name,title,typeid,tag,copyright,content,state,ctime,mtime) VALUES (?,?,?,?,?,?,?,?,?,?)"
	// update
	_upTplSQL  = "UPDATE archive_templates SET name=?,title=?,typeid=?,tag=?,copyright=?,content=?,mtime=? WHERE id=? AND mid=?"
	_delTplSQL = "UPDATE archive_templates SET state=?,mtime=? WHERE id=? AND mid=?"
	// select
	_getTplSQL      = "SELECT id,name,title,typeid,tag,copyright,content,state,ctime,mtime FROM archive_templates WHERE id=? AND mid=?"
	_getMutilTplSQL = "SELECT id,name,title,typeid,tag,copyright,content FROM archive_templates WHERE mid=? AND state=0 ORDER BY ctime DESC"
	_getCntSQL      = "SELECT count(id) FROM archive_templates WHERE mid = ? AND state = 0"
)

// templates get all Template from db.
func (d *Dao) templates(c context.Context, mid int64) (tps []*template.Template, err error) {
	tps = make([]*template.Template, 0)
	rows, err := d.getMutilTplStmt.Query(c, mid)
	if err != nil {
		log.Error("mysqlDB.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		t := &template.Template{}
		if err = rows.Scan(&t.ID, &t.Name, &t.Title, &t.TypeID, &t.Tag, &t.Copyright, &t.Content); err != nil {
			log.Error("row.Scan error(%v)", err)
			return
		}
		tps = append(tps, t)
	}
	return
}

// Template get single template.
func (d *Dao) Template(c context.Context, id, mid int64) (t *template.Template, err error) {
	row := d.getTplStmt.QueryRow(c, id, mid)
	t = &template.Template{}
	if err = row.Scan(&t.ID, &t.Name, &t.Title, &t.TypeID, &t.Tag, &t.Copyright, &t.Content, &t.State, &t.CTime, &t.MTime); err != nil {
		if err == sql.ErrNoRows {
			t = nil
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// Templates get all Template
func (d *Dao) Templates(c context.Context, mid int64) (tps []*template.Template, err error) {
	// try cache
	if tps, _ = d.tplCache(c, mid); tps != nil {
		return
	}
	// from db
	if tps, err = d.templates(c, mid); tps != nil {
		d.addCache(func() {
			d.addTplCache(context.TODO(), mid, tps)
		})
	}
	return
}

// AddTemplate add Template
func (d *Dao) AddTemplate(c context.Context, mid int64, tp *template.Template) (id int64, err error) {
	res, err := d.addTplStmt.Exec(c, mid, tp.Name, tp.Title, tp.TypeID, tp.Tag, tp.Copyright, tp.Content, tp.State, tp.CTime, tp.MTime)
	if err != nil {
		log.Error("d.AddTemplate.Exec(%+v) error(%v)", tp, err)
		return
	}
	if id, err = res.LastInsertId(); err != nil {
		log.Error("res.LastInsertId(%d) error(%v)", id, err)
		return
	}
	d.addCache(func() {
		d.delTplCache(context.TODO(), mid)
	})
	return
}

// UpTemplate update Template
func (d *Dao) UpTemplate(c context.Context, mid int64, tp *template.Template) (rows int64, err error) {
	res, err := d.upTplStmt.Exec(c, tp.Name, tp.Title, tp.TypeID, tp.Tag, tp.Copyright, tp.Content, tp.MTime, tp.ID, mid)
	if err != nil {
		log.Error("d.upTplStmt.Exec(%d, %d) error(%v)", mid, tp.ID, err)
		return
	}
	if rows, err = res.RowsAffected(); err != nil {
		log.Error("res.RowsAffected rows(%d) error(%v)", rows, err)
		return
	}
	d.addCache(func() {
		d.delTplCache(context.TODO(), mid)
	})
	return
}

// DelTemplate delete Template
func (d *Dao) DelTemplate(c context.Context, mid int64, tp *template.Template) (rows int64, err error) {
	res, err := d.delTplStmt.Exec(c, tp.State, tp.MTime, tp.ID, mid)
	if err != nil {
		log.Error("d.delTplStmt.Exec(%d ) error(%v)", mid, err)
		return
	}
	if rows, err = res.RowsAffected(); err != nil {
		log.Error("res.RowsAffected rows(%d) error(%v)", rows, err)
		return
	}
	d.addCache(func() {
		d.delTplCache(context.TODO(), mid)
	})
	return
}

// Count count all state Template
func (d *Dao) Count(c context.Context, mid int64) (count int64, err error) {
	row := d.getCntStmt.QueryRow(c, mid)
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			count = 0
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}
