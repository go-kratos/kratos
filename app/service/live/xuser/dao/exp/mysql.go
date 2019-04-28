package exp

import (
	"context"
	"database/sql"
	"fmt"
	expm "go-common/app/service/live/xuser/model/exp"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_shard            = 10
	_insExp           = "INSERT IGNORE INTO user_exp_%d (uid,uexp,rexp) VALUES(?,?,?)"
	_selExp           = "SELECT uid,uexp,rexp,ctime,mtime FROM user_exp_%d where uid=?"
	_inSelExp         = "SELECT uid,uexp,rexp,ctime,mtime FROM user_exp_%d where uid IN (%s)"
	_addUexp          = "INSERT INTO user_exp_%d(uid,uexp,rexp) VALUES(?,?,0) ON DUPLICATE KEY UPDATE uexp=uexp+%d"
	_addRexp          = "INSERT INTO user_exp_%d(uid,uexp,rexp) VALUES(?,0,?) ON DUPLICATE KEY UPDATE rexp=rexp+%d"
	_errorDBLogPrefix = "xuser.exp.dao.mysql"
)

// InitExp 初始化用户经验,用于首次查询
func (d *Dao) InitExp(c context.Context, uid int64, uexp int64, rexp int64) (row int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_insExp, uid%_shard), uid, uexp, rexp)
	if err != nil {
		log.Error(_errorDBLogPrefix+"|InitExp d.exp.Exec err: %v", err)
		return
	}
	return res.RowsAffected()
}

// Exp 查询一条记录
func (d *Dao) Exp(c context.Context, uid int64) (exp *expm.Exp, err error) {
	row := d.db.QueryRow(c, fmt.Sprintf(_selExp, uid%_shard), uid)
	exp = &expm.Exp{}
	if err = row.Scan(&exp.UID, &exp.Uexp, &exp.Rexp, &exp.CTime, &exp.MTime); err == sql.ErrNoRows {
		// 查询结果为空时，初始化数据
		_, err = d.InitExp(c, uid, 0, 0)
	}
	if err != nil {
		log.Error(_errorDBLogPrefix+"|Exp row.Scan err: %v", err)
		return
	}
	return
}

// MultiExp 批量查询
func (d *Dao) MultiExp(c context.Context, uids []int64) (exps []*expm.Exp, err error) {
	exps = make([]*expm.Exp, 0)
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
		rows, err1 := d.db.Query(c, fmt.Sprintf(_inSelExp, index, xstr.JoinInts(uids)))
		if err1 != nil {
			err = err1
			log.Error(_errorDBLogPrefix+"|MultiExp d.exp.Query err: %v", err)
			return
		}

		for rows.Next() {
			ele := &expm.Exp{}
			if err = rows.Scan(&ele.UID, &ele.Uexp, &ele.Rexp, &ele.CTime, &ele.MTime); err != nil {
				log.Error(_errorDBLogPrefix+"|MultiExp rows.Scan err: %v", err)
				return
			}
			exps = append(exps, ele)
			delete(um, ele.UID)
		}
	}

	// 初始化不存在的数据，补齐数据
	for uid := range um {
		d.InitExp(c, uid, 0, 0)
		ele := &expm.Exp{UID: uid, Uexp: 0, Rexp: 0}
		exps = append(exps, ele)
	}
	return
}

// AddUexp 添加用户经验
func (d *Dao) AddUexp(c context.Context, uid int64, uexp int64) (affect int64, err error) {
	upSQL := fmt.Sprintf(_addUexp, uid%_shard, uexp)
	res, err := d.db.Exec(c, upSQL, uid, uexp)
	if err != nil {
		log.Error(_errorDBLogPrefix+"|Exec(%s) error(%v)", upSQL, err)
		return
	}
	return res.RowsAffected()
}

// AddRexp 添加主播经验
func (d *Dao) AddRexp(c context.Context, uid int64, rexp int64) (affect int64, err error) {
	upSQL := fmt.Sprintf(_addRexp, uid%_shard, rexp)
	res, err := d.db.Exec(c, upSQL, uid, rexp)
	if err != nil {
		log.Error(_errorDBLogPrefix+"|AddRexp|Exec(%s) error(%v)", upSQL, err)
		return
	}
	return res.RowsAffected()
}
