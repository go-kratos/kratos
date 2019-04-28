package dao

import (
	"context"
	"database/sql"
	"fmt"

	"go-common/app/service/live/userexp/model"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_shard    = 10
	_insExp   = "INSERT IGNORE INTO user_exp_%d (uid,uexp,rexp) VALUES(?,?,?)"
	_selExp   = "SELECT uid,uexp,rexp FROM user_exp_%d where uid=?"
	_inSelExp = "SELECT uid,uexp,rexp FROM user_exp_%d where uid IN (%s)"
	_addUexp  = "INSERT INTO user_exp_%d(uid,uexp,rexp) VALUES(?,?,0) ON DUPLICATE KEY UPDATE uexp=uexp+%d"
	_addRexp  = "INSERT INTO user_exp_%d(uid,uexp,rexp) VALUES(?,0,?) ON DUPLICATE KEY UPDATE rexp=rexp+%d"
)

// InitExp 初始化用户经验,用于首次查询
func (d *Dao) InitExp(c context.Context, uid int64, uexp int64, rexp int64) (row int64, err error) {
	res, err := d.expDb.Exec(c, fmt.Sprintf(_insExp, uid%_shard), uid, uexp, rexp)
	if err != nil {
		log.Error("[dao.exp|InitExp] d.exp.Exec err: %v", err)
		return
	}
	return res.RowsAffected()
}

// Exp 查询一条记录
func (d *Dao) Exp(c context.Context, uid int64) (exp *model.Exp, err error) {
	row := d.expDb.QueryRow(c, fmt.Sprintf(_selExp, uid%_shard), uid)
	exp = &model.Exp{}
	if err = row.Scan(&exp.Uid, &exp.Uexp, &exp.Rexp); err == sql.ErrNoRows {
		// 查询结果为空时，初始化数据
		_, err = d.InitExp(c, uid, 0, 0)
	}
	if err != nil {
		log.Error("[dao.exp|Exp] row.Scan err: %v", err)
		return
	}
	return
}

// MultiExp 批量查询
func (d *Dao) MultiExp(c context.Context, uids []int64) (exps []*model.Exp, err error) {
	var (
		suffix   int64
		uidGroup [_shard][]int64
		um       = make(map[int64]struct{}, len(uids))
	)
	for _, uid := range uids {
		suffix = uid % _shard
		uidGroup[suffix] = append(uidGroup[suffix], uid)
		um[uid] = struct{}{}
	}
	for index, uids := range uidGroup {
		if 0 == len(uids) {
			continue
		}
		rows, err1 := d.expDb.Query(c, fmt.Sprintf(_inSelExp, index, xstr.JoinInts(uids)))
		if err1 != nil {
			err = err1
			log.Error("[dao.exp|MultiExp] d.exp.Query err: %v", err)
			return
		}

		for rows.Next() {
			ele := &model.Exp{}
			if err = rows.Scan(&ele.Uid, &ele.Uexp, &ele.Rexp); err != nil {
				log.Error("[dao.exp|MultiExp] rows.Scan err: %v", err)
				return
			}
			exps = append(exps, ele)

			delete(um, ele.Uid)
		}
	}

	// 初始化不存在的数据，补齐数据
	for uid := range um {
		d.InitExp(c, uid, 0, 0)
		ele := &model.Exp{Uid: uid, Uexp: 0, Rexp: 0}
		exps = append(exps, ele)
	}
	return
}

// AddUexp 添加用户经验
func (d *Dao) AddUexp(c context.Context, uid int64, uexp int64) (affect int64, err error) {
	upSQL := fmt.Sprintf(_addUexp, uid%_shard, uexp)
	res, err := d.expDb.Exec(c, upSQL, uid, uexp)
	if err != nil {
		log.Error("db.Exec(%s) error(%v)", upSQL, err)
		return
	}
	return res.RowsAffected()
}

// AddRexp 添加主播经验
func (d *Dao) AddRexp(c context.Context, uid int64, rexp int64) (affect int64, err error) {
	upSQL := fmt.Sprintf(_addRexp, uid%_shard, rexp)
	res, err := d.expDb.Exec(c, upSQL, uid, rexp)
	if err != nil {
		log.Error("db.Exec(%s) error(%v)", upSQL, err)
		return
	}
	return res.RowsAffected()
}
