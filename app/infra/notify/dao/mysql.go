package dao

import (
	"context"
	"encoding/json"

	"go-common/app/infra/notify/model"
)

const (
	_loadNotify = "SELECT n.id,topic.cluster,topic.topic,a.group,n.callback,n.concurrent,n.filter,n.mtime FROM notify AS n LEFT JOIN auth2 as a ON a.id=n.gid LEFT JOIN topic ON a.topic_id=topic.id WHERE n.state=1 AND n.zone=?"
	_loadPub    = `SELECT a.group,topic.cluster,topic.topic,a.operation,app2.app_secret FROM auth2 AS a LEFT JOIN app2 ON app2.id=a.app_id LEFT JOIN topic ON topic.id=a.topic_id`
	_selFilter  = "SELECT `filters` FROM filters WHERE nid=?"
	_addFailBk  = "INSERT INTO fail_backup(topic,`group`,`cluster`,msg,`index`) VALUE (?,?,?,?,?)"
	_delFailBk  = "DELETE FROM fail_backup where id=?"
	_loadFailBk = "SELECT id,topic,`group`,`cluster`,`msg`,`offset`,`index` FROM fail_backup"
)

// LoadNotify load all notify config.
func (d *Dao) LoadNotify(c context.Context, zone string) (ns []*model.Watcher, err error) {
	rows, err := d.db.Query(c, _loadNotify, zone)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		n := new(model.Watcher)
		if err = rows.Scan(&n.ID, &n.Cluster, &n.Topic, &n.Group, &n.Callback, &n.Concurrent, &n.Filter, &n.Mtime); err != nil {
			return
		}
		ns = append(ns, n)
	}
	return
}

// LoadPub load all pub config.
func (d *Dao) LoadPub(c context.Context) (ps []*model.Pub, err error) {
	rows, err := d.db.Query(c, _loadPub)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		n := new(model.Pub)
		if err = rows.Scan(&n.Group, &n.Cluster, &n.Topic, &n.Operation, &n.AppSecret); err != nil {
			return
		}
		ps = append(ps, n)
	}
	return
}

// Filters get filter condition.
func (d *Dao) Filters(c context.Context, id int64) (fs []*model.Filter, err error) {
	rows := d.db.QueryRow(c, _selFilter, id)
	if err != nil {
		return
	}
	var filters string
	if err = rows.Scan(&filters); err != nil {
		return
	}
	err = json.Unmarshal([]byte(filters), &fs)
	return
}

// AddFailBk add fail msg to fail backup.
func (d *Dao) AddFailBk(c context.Context, topic, group, cluster, msg string, index int64) (id int64, err error) {
	res, err := d.db.Exec(c, _addFailBk, topic, group, cluster, msg, index)
	if err != nil {
		return
	}
	return res.LastInsertId()
}

// DelFailBk del msg from fail backup.
func (d *Dao) DelFailBk(c context.Context, id int64) (affected int64, err error) {
	res, err := d.db.Exec(c, _delFailBk, id)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// LoadFailBk load all fail backup msg.
func (d *Dao) LoadFailBk(c context.Context) (fbs []*model.FailBackup, err error) {
	rows, err := d.db.Query(c, _loadFailBk)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		fb := new(model.FailBackup)
		if err = rows.Scan(&fb.ID, &fb.Topic, &fb.Group, &fb.Cluster, &fb.Msg, &fb.Offset, &fb.Index); err != nil {
			return
		}
		fbs = append(fbs, fb)
	}
	return
}
