package dao

import (
	"context"
	"database/sql"
	"fmt"

	"go-common/app/job/main/click/model"
	"go-common/library/log"
)

const (
	_cliSQL          = "SELECT aid,web,h5,outside,ios,android,android_tv FROM %s WHERE aid=?"
	_addCliSQL       = "INSERT INTO %s(aid,web,h5,outside,ios,android,android_tv) VALUES(?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE web=?,h5=?,outside=?,ios=?,android=?,android_tv=?"
	_upCliSQL        = "UPDATE %s SET web=web+?,h5=h5+?,outside=outside+?,ios=ios+?,android=android+?,android_tv=android_tv+? WHERE aid=?"
	_upSpecialCliSQL = "UPDATE %s SET %s=? WHERE aid=?"
)

func getTable(aid int64) string {
	return fmt.Sprintf("archive_click_%02d", aid%100)
}

// Click get click
func (d *Dao) Click(c context.Context, aid int64) (cli *model.ClickInfo, err error) {
	row := d.db.QueryRow(c, fmt.Sprintf(_cliSQL, getTable(aid)), aid)
	cli = &model.ClickInfo{}
	if err = row.Scan(&cli.Aid, &cli.Web, &cli.H5, &cli.Outer, &cli.Ios, &cli.Android, &cli.AndroidTV); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			cli = nil
		} else {
			log.Error("row.Scan error(%v)")
		}
	}
	return
}

// AddClick add av clicks
func (d *Dao) AddClick(c context.Context, aid, web, h5, out, ios, android, androidTV int64) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_addCliSQL, getTable(aid)), aid, web, h5, out, ios, android, androidTV, web, h5, out, ios, android, androidTV)
	if err != nil {
		log.Error("d.addCliStmt.Exec(%d) error(%v)", aid, err)
		return
	}
	return res.RowsAffected()
}

// UpClick update av clicks
func (d *Dao) UpClick(c context.Context, cli *model.ClickInfo) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_upCliSQL, getTable(cli.Aid)), cli.Web, cli.H5, cli.Outer, cli.Ios, cli.Android, cli.AndroidTV, cli.Aid)
	if err != nil {
		log.Error("d.upCliStmt.Exec(+%v) error(%v)", cli, err)
		return
	}
	return res.RowsAffected()
}

// UpSpecial update special platform click
func (d *Dao) UpSpecial(c context.Context, aid int64, tp string, num int64) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_upSpecialCliSQL, getTable(aid), tp), num, aid)
	if err != nil {
		log.Error("d.db.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}
