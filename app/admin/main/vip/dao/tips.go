package dao

import (
	"context"
	xsql "database/sql"
	"fmt"

	"go-common/app/admin/main/vip/model"
	"go-common/library/database/sql"

	"github.com/pkg/errors"
)

const (
	_tipsSQL      = "SELECT `id`,`platform`,`version`,`tip`,`link`,`start_time`,`end_time`,`level`,`judge_type`,`operator`,`deleted`, `position`,`ctime`,`mtime`FROM `vip_tips` WHERE  `deleted` = 0 "
	_tipsByIDSQL  = "SELECT `id`,`platform`,`version`,`tip`,`link`,`start_time`,`end_time`,`level`,`judge_type`,`operator`,`deleted`, `position`,`ctime`,`mtime`FROM `vip_tips` WHERE  `id` = ?;"
	_updateTipSQL = "UPDATE `vip_tips` SET `platform` = ?,`version` = ?,`tip` = ?,`link` = ?,`start_time` = ?,`end_time` = ?,`level` = ?, `judge_type` = ?, `operator` = ?, `position` = ? WHERE `id` = ?;"
	_addTipSQL    = "INSERT INTO `vip_tips`(`platform`,`version`,`tip`,`link`,`start_time`,`end_time`,`level`,`judge_type`,`operator`,`deleted`,`ctime`, `position`)VALUES(?,?,?,?,?,?,?,?,?,?,?,?);"
	_deleteTipSQL = "UPDATE `vip_tips` SET `deleted` = ?,`operator` = ? WHERE `id` = ?;"
	_expireTipSQL = "UPDATE `vip_tips` SET `end_time` = ?,`operator` = ? WHERE `id` = ?;"
)

// TipList tips list.
func (d *Dao) TipList(c context.Context, platform int8, state int8, now int64, position int8) (rs []*model.Tips, err error) {
	var (
		rows *sql.Rows
		sql  = _tipsSQL
	)
	switch state {
	case model.WaitShowTips:
		sql += fmt.Sprintf(" AND `start_time` > %d ", now)
	case model.EffectiveTips:
		sql += fmt.Sprintf(" AND `start_time` < %d AND `end_time`> %d", now, now)
	case model.ExpireTips:
		sql += fmt.Sprintf(" AND `end_time` < %d ", now)
	}
	if platform != 0 {
		sql += fmt.Sprintf(" AND `platform` = %d ", platform)
	}
	if position != 0 {
		sql += fmt.Sprintf(" AND `position` = %d ", position)
	}
	if rows, err = d.db.Query(c, sql); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.Tips)
		if err = rows.Scan(&r.ID, &r.Platform, &r.Version, &r.Tip, &r.Link, &r.StartTime, &r.EndTime, &r.Level, &r.JudgeType, &r.Operator, &r.Deleted, &r.Position, &r.Ctime, &r.Mtime); err != nil {
			rs = nil
			err = errors.WithStack(err)
			return
		}
		rs = append(rs, r)
	}
	if err = rows.Err(); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// TipByID by id .
func (d *Dao) TipByID(c context.Context, id int64) (r *model.Tips, err error) {
	res := d.db.QueryRow(c, _tipsByIDSQL, id)
	r = new(model.Tips)
	if err = res.Scan(&r.ID, &r.Platform, &r.Version, &r.Tip, &r.Link, &r.StartTime, &r.EndTime, &r.Level, &r.JudgeType, &r.Operator, &r.Deleted, &r.Position, &r.Ctime, &r.Mtime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			r = nil
			return
		}
		err = errors.WithStack(err)
	}
	return
}

// TipUpdate tip update.
func (d *Dao) TipUpdate(c context.Context, t *model.Tips) (eff int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, _updateTipSQL, t.Platform, t.Version, t.Tip, t.Link, t.StartTime, t.EndTime, t.Level, t.JudgeType, t.Operator, t.Position, t.ID); err != nil {
		err = errors.WithStack(err)
		return
	}
	if eff, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// AddTip add tip.
func (d *Dao) AddTip(c context.Context, t *model.Tips) (lid int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, _addTipSQL, t.Platform, t.Version, t.Tip, t.Link, t.StartTime, t.EndTime, t.Level, t.JudgeType, t.Operator, t.Deleted, t.Ctime, t.Position); err != nil {
		err = errors.WithStack(err)
		return
	}
	if lid, err = res.LastInsertId(); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// DeleteTip tip delete.
func (d *Dao) DeleteTip(c context.Context, id int64, deleted int8, operator string) (eff int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, _deleteTipSQL, deleted, operator, id); err != nil {
		err = errors.WithStack(err)
		return
	}
	if eff, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// ExpireTip tip expire.
func (d *Dao) ExpireTip(c context.Context, id int64, operator string, t int64) (eff int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, _expireTipSQL, t, operator, id); err != nil {
		err = errors.WithStack(err)
		return
	}
	if eff, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
	}
	return
}
