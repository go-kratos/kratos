package dao

import (
	"context"
	xsql "database/sql"
	"encoding/json"
	"fmt"

	"go-common/app/service/openplatform/abtest/model"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"

	"github.com/pkg/errors"
)

const (
	_selAllByStatus         = "SELECT `id`,`name`,`desc`,`stra`,`seed`,`result`,`status`,`version`,`group`,`author_name`,`modifier_name`,`ctime`,`mtime` FROM `abtest` WHERE `deleted`=0 AND `status` IN (%s) ORDER BY mtime DESC LIMIT %d, %d"
	_selAllByStatusAndGroup = "SELECT `id`,`name`,`desc`,`stra`,`seed`,`result`,`status`,`version`,`group`,`author_name`,`modifier_name`,`ctime`,`mtime` FROM `abtest` WHERE `group`=? AND `deleted`=0 AND `status` IN (%s) ORDER BY mtime DESC LIMIT %d, %d"
	_selActByGroup          = "SELECT `id`,`name`,`desc`,`stra`,`seed`,`result`,`status`,`version`,`group`,`author_name`,`modifier_name`,`ctime`,`mtime` FROM `abtest` WHERE `group`=? AND `status`=1 AND `deleted`=0"
	_selCnt                 = "SELECT count(*) as `count` FROM `abtest` WHERE `deleted`=0 AND `status` IN (%s)"
	_selCntByGroup          = "SELECT count(*) as `count` FROM `abtest` WHERE `group`=? AND `deleted`=0 AND `status` IN (%s)"
	_selByID                = "SELECT `id`,`name`,`desc`,`stra`,`seed`,`result`,`status`,`version`,`group`,`author_name`,`modifier_name`,`ctime`,`mtime` FROM `abtest` WHERE `id`=? AND `deleted`=0"
	_selByIDAndGroup        = "SELECT `id`,`name`,`desc`,`stra`,`seed`,`result`,`status`,`version`,`group`,`author_name`,`modifier_name`,`ctime`,`mtime` FROM `abtest` WHERE `id`=? AND `group`=? AND `deleted`=0"
	_insAB                  = "INSERT INTO `abtest` (`name`,`desc`,`stra`,`seed`,`result`,`status`,`group`,`author_name`,`modifier_name`) VALUES(?,?,?,?,?,0,?,?,?)"
	_upAB                   = "UPDATE `abtest` SET `name`=?,`desc`=?,`stra`=?,`result`=?,`modifier_name`=?,`version`=? WHERE `id`=?"
	_upStatus               = "UPDATE `abtest` SET `status`=?,`modifier_name`=? WHERE `id`=?"
	_delAB                  = "UPDATE `abtest` SET `deleted`=1 WHERE `id`=? AND `status`!=1"
)

//ActByGroup 根据group获取当前激活项目
func (d *Dao) ActByGroup(c context.Context, group int) (res []*model.AB, err error) {
	var (
		rows    *sql.Rows
		straStr string
	)
	if rows, err = d.db.Query(c, _selActByGroup, group); err != nil {
		err = errors.Wrap(err, fmt.Sprintf("[dao.abtest|ActByGroup] d.db.Query err: %v", err))
		return
	}
	defer rows.Close()
	res = make([]*model.AB, 0)
	for rows.Next() {
		ele := &model.AB{}
		if err = rows.Scan(&ele.ID, &ele.Name, &ele.Desc, &straStr, &ele.Seed, &ele.Result, &ele.Status, &ele.Version, &ele.Group, &ele.Author, &ele.Modifier, &ele.CreateTime, &ele.ModifyTime); err != nil {
			err = errors.Wrap(err, fmt.Sprintf("[dao.abtest|ActByGroup] rows.Scan err: %v", err))
			return
		}
		if err = json.Unmarshal([]byte(straStr), &ele.Stra); err != nil {
			err = errors.Wrap(err, fmt.Sprintf("[dao.abtest|ActByGroup] json.Unmarshal err: %v", err))
			return
		}
		res = append(res, ele)
	}
	return
}

//Ab 获取单个测试
func (d *Dao) Ab(c context.Context, id int) (res *model.AB, err error) {
	var straStr string
	row := d.db.QueryRow(c, _selByID, id)
	res = &model.AB{}
	if err = row.Scan(&res.ID, &res.Name, &res.Desc, &straStr, &res.Seed, &res.Result, &res.Status, &res.Version, &res.Group, &res.Author, &res.Modifier, &res.CreateTime, &res.ModifyTime); err != nil {
		err = errors.Wrap(err, fmt.Sprintf("[dao.abtest|Ab] row.Scan err: %v", err))
		return
	}
	err = json.Unmarshal([]byte(straStr), &res.Stra)
	return
}

//AbByIDAndGroup 获取单个测试
func (d *Dao) AbByIDAndGroup(c context.Context, id int, group int) (res *model.AB, err error) {
	var (
		straStr string
		row     *sql.Row
	)
	row = d.db.QueryRow(c, _selByIDAndGroup, id, group)
	res = &model.AB{}
	if err = row.Scan(&res.ID, &res.Name, &res.Desc, &straStr, &res.Seed, &res.Result, &res.Status, &res.Version, &res.Group, &res.Author, &res.Modifier, &res.CreateTime, &res.ModifyTime); err != nil {
		err = errors.Wrap(err, fmt.Sprintf("[dao.abtest|AbByIDAndGroup] row.Scan err: %v", err))
		return
	}
	err = json.Unmarshal([]byte(straStr), &res.Stra)
	return
}

//CountAb 获取测试数量
func (d *Dao) CountAb(c context.Context, mstatus string) (count int, err error) {
	row := d.db.QueryRow(c, fmt.Sprintf(_selCnt, mstatus))
	if err = row.Scan(&count); err != nil {
		err = errors.Wrap(err, fmt.Sprintf("[dao.abtest|CountAb] row.Scan err: %v", err))
		return
	}
	return
}

//CountAbByGroup 获取测试数量
func (d *Dao) CountAbByGroup(c context.Context, mstatus string, group int) (count int, err error) {
	row := d.db.QueryRow(c, fmt.Sprintf(_selCntByGroup, mstatus), group)
	if err = row.Scan(&count); err != nil {
		err = errors.Wrap(err, fmt.Sprintf("[dao.abtest|CountAbByGroup] row.Scan err: %v", err))
		return
	}
	return
}

//ListAb 分页获取所有测试列表
func (d *Dao) ListAb(c context.Context, offset, size int, mstatus string) (res []*model.AB, count int, err error) {
	var (
		rows    *sql.Rows
		straStr string
	)
	if rows, err = d.db.Query(c, fmt.Sprintf(_selAllByStatus, mstatus, offset, size)); err != nil {
		err = errors.Wrap(err, fmt.Sprintf("[dao.abtest|ListAb] d.db.Query err: %v", err))
		return
	}
	defer rows.Close()
	res = make([]*model.AB, 0)
	for rows.Next() {
		ele := &model.AB{}
		if err = rows.Scan(&ele.ID, &ele.Name, &ele.Desc, &straStr, &ele.Seed, &ele.Result, &ele.Status, &ele.Version, &ele.Group, &ele.Author, &ele.Modifier, &ele.CreateTime, &ele.ModifyTime); err != nil {
			err = errors.Wrap(err, fmt.Sprintf("[dao.abtest|ListAb] rows.Scan err: %v", err))
			return
		}
		err = json.Unmarshal([]byte(straStr), &ele.Stra)
		if err != nil {
			err = errors.Wrap(err, fmt.Sprintf("[dao.abtest|ListAb] json.Unmarshal err: %v", err))
			continue
		}
		res = append(res, ele)
	}
	count, err = d.CountAb(c, mstatus)
	return
}

//ListAbByGroup 分页获取分组测试列表
func (d *Dao) ListAbByGroup(c context.Context, offset, size int, mstatus string, group int) (res []*model.AB, count int, err error) {
	var (
		rows    *sql.Rows
		straStr string
	)
	if rows, err = d.db.Query(c, fmt.Sprintf(_selAllByStatusAndGroup, mstatus, offset, size), group); err != nil {
		err = errors.Wrap(err, fmt.Sprintf("[dao.abtest|ListAbByGroup] d.db.Query err: %v", err))
		return
	}
	defer rows.Close()
	res = make([]*model.AB, 0)
	for rows.Next() {
		ele := &model.AB{}
		if err = rows.Scan(&ele.ID, &ele.Name, &ele.Desc, &straStr, &ele.Seed, &ele.Result, &ele.Status, &ele.Version, &ele.Group, &ele.Author, &ele.Modifier, &ele.CreateTime, &ele.ModifyTime); err != nil {
			err = errors.Wrap(err, fmt.Sprintf("[dao.abtest|ListAbByGroup] rows.Scan err: %v", err))
			return
		}
		err = json.Unmarshal([]byte(straStr), &ele.Stra)
		if err != nil {
			err = errors.Wrap(err, fmt.Sprintf("[dao.abtest|ListAbByGroup] json.Unmarshal err: %v", err))
			continue
		}
		res = append(res, ele)
	}
	count, err = d.CountAbByGroup(c, mstatus, group)
	return
}

//AddAb 添加AB实验
func (d *Dao) AddAb(c context.Context, name, desc string, stra string, seed, result, group int, username string) (newID int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, _insAB, name, desc, stra, seed, result, group, username, username); err != nil {
		return
	}
	newID, err = res.LastInsertId()
	return
}

//DelAb 删除AB实验
func (d *Dao) DelAb(c context.Context, id int) (rowsAffected int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, _delAB, id); err != nil {
		return
	}
	return res.RowsAffected()
}

//UpAb 更新AB实验
func (d *Dao) UpAb(c context.Context, id int, name, desc string, stra string, result int, username string, newVersion, status, group int) (rowsAffected int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, _upAB, name, desc, stra, result, username, newVersion, id); err != nil {
		return
	}
	return res.RowsAffected()
}

//UpStatus 更新AB实验状态
func (d *Dao) UpStatus(c context.Context, id, status int, username string, group int) (rowsAffected int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, _upStatus, status, username, id); err != nil {
		return
	}
	return res.RowsAffected()
}

//Incr stat result, save to redis
func (d *Dao) Incr(c context.Context, key string) {
	conn := d.redis.Get(c)
	defer conn.Close()
	conn.Do("incr", key)
	conn.Do("expire", key, 86400*7)
}

//Move backup key
func (d *Dao) Move(c context.Context, key string) {
	conn := d.redis.Get(c)
	defer conn.Close()
	i, _ := redis.Int(conn.Do("get", key))
	if _, err := conn.Do("setex", "O:"+key, 86400*7, i); err == nil {
		conn.Do("del", key)
	}
}

// GetFromRedis .
func (d *Dao) GetFromRedis(c context.Context, key string) (value int, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	value, err = redis.Int(conn.Do("get", key))
	return
}
