package dao

import (
	"context"
	xsql "database/sql"
	"encoding/json"
	"fmt"
	"go-common/app/service/bbq/user/api"
	"go-common/app/service/bbq/user/internal/conf"
	"go-common/app/service/bbq/user/internal/model"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	"strconv"
	"strings"
)

const (
	_incrUserBase = "insert into user_base (`mid`,`uname`,`face`,`birthday`,`exp`,`level`,`user_type`,`complete_degree`,`sex`) values(?,?,?,?,?,?,?,?,?)"
	//_incrUserUname   = "insert into user_base (`mid`,`uname`) values(?,?)"
	// TODO: 用户签名数据暂时隐藏
	// _userBase        = "select `mid`,`uname`,`face`,`birthday`,`exp`,`level`,`ctime`,`mtime`,`signature`,`region`,`sex`, `user_type`,`complete_degree`from user_base where `mid` in (%s)"
	_userBase           = "select `mid`,`uname`,`face`,`birthday`,`exp`,`level`,`ctime`,`mtime`,'',`region`,`sex`, `user_type`,`complete_degree`from user_base where `mid` in (%s)"
	_updateUserUname    = "update user_base set `uname` = ? where `mid` = ?"
	_selectUname        = "select `mid` from user_base where `uname` = ? and `mid` != ?"
	_selectBZhanUpUname = "select `mid` from user_statistics_hive where `uname` = ? and `mid` != ? and `fan_total` > 10000"
	_updateUser         = "update user_base set uname=?, face=?, birthday=?, sex=?, region=?, signature=?, complete_degree=? where mid=?"
	_updateUserField    = "update user_base set `%s` = ? where mid = ?"
)

// keyUserBase 用户基础信息缓存key
func keyUserBase(mid int64) string {
	return fmt.Sprintf(model.CacheKeyUserBase, mid)
}

// RawUserBase 从数据库获取用户基础信息
func (d *Dao) RawUserBase(c context.Context, mids []int64) (res map[int64]*api.UserBase, err error) {
	if len(mids) == 0 {
		return
	}
	var midStr string
	for _, mid := range mids {
		if len(midStr) != 0 {
			midStr += ","
		}
		midStr += strconv.FormatInt(mid, 10)
	}
	querySQL := fmt.Sprintf(_userBase, midStr)
	rows, err := d.db.Query(c, querySQL)
	if err != nil {
		log.Errorv(c, log.KV("event", "mysql_query"), log.KV("error", err), log.KV("sql", querySQL))
		return
	}
	defer rows.Close()
	for rows.Next() {
		userBase := new(api.UserBase)
		if err = rows.Scan(&userBase.Mid,
			&userBase.Uname,
			&userBase.Face,
			&userBase.Birthday,
			&userBase.Exp,
			&userBase.Level,
			&userBase.Ctime,
			&userBase.Mtime,
			&userBase.Signature,
			&userBase.Region,
			&userBase.Sex,
			&userBase.UserType,
			&userBase.CompleteDegree); err != nil {
			log.Errorv(c, log.KV("event", "mysql_scan"), log.KV("error", err), log.KV("sql", querySQL))
			return
		}
		if res == nil {
			res = make(map[int64]*api.UserBase)
		}
		res[userBase.Mid] = userBase
	}
	log.Infov(c, log.KV("event", "mysql_query"), log.KV("row_num", len(res)), log.KV("sql", querySQL))
	return
}

// CacheUserBase  multi get user base from cache.
func (d *Dao) CacheUserBase(c context.Context, mids []int64) (res map[int64]*api.UserBase, err error) {
	if res == nil {
		res = make(map[int64]*api.UserBase)
	}
	keys := make([]string, 0, len(mids))
	keyMidMap := make(map[int64]bool, len(mids))
	for _, mid := range mids {
		key := keyUserBase(mid)
		if _, exist := keyMidMap[mid]; !exist {
			// duplicate mid
			keyMidMap[mid] = true
			keys = append(keys, key)
		}
	}
	conn := d.redis.Get(c)
	defer conn.Close()
	for _, key := range keys {
		conn.Send("GET", key)
	}
	conn.Flush()
	var data []byte
	for i := 0; i < len(keys); i++ {
		if data, err = redis.Bytes(conn.Receive()); err != nil {
			if err == redis.ErrNil {
				err = nil
			} else {
				log.Errorv(c, log.KV("event", "redis_get"), log.KV("key", keys[i]))
			}
			continue
		}
		baseItem := new(api.UserBase)
		json.Unmarshal(data, baseItem)
		res[baseItem.Mid] = baseItem
	}
	log.Infov(c, log.KV("event", "redis_get"), log.KV("row_num", len(res)))
	return
}

// AddCacheUserBase 添加用户缓存
func (d *Dao) AddCacheUserBase(c context.Context, userBases map[int64]*api.UserBase) (err error) {
	keyValueMap := make(map[string][]byte, len(userBases))
	for mid, userBase := range userBases {
		key := keyUserBase(mid)
		if _, exist := keyValueMap[key]; !exist {
			data, _ := json.Marshal(userBase)
			keyValueMap[key] = data
		}
	}
	conn := d.redis.Get(c)
	defer conn.Close()
	for key, value := range keyValueMap {
		conn.Send("SET", key, value, "EX", model.CacheExpireUserBase)
	}
	conn.Flush()
	for i := 0; i < len(keyValueMap); i++ {
		conn.Receive()
	}
	log.Infov(c, log.KV("event", "redis_set"), log.KV("row_num", len(userBases)))
	return
}

//DelCacheUserBase 删除用户缓存
func (d *Dao) DelCacheUserBase(c context.Context, mid int64) {
	var key = keyUserBase(mid)
	conn := d.redis.Get(c)
	defer conn.Close()
	conn.Do("DEL", key)
}

////TxAddUserBase .
//func (d *Dao) TxAddUserBase(c context.Context, tx *sql.Tx, userBase *api.UserBase) (num int64, err error) {
//	var res xsql.Result
//	if res, err = tx.Exec(_incrUserBase, userBase.Mid, userBase.Uname, userBase.Face, userBase.Birthday, userBase.Exp, userBase.Level, userBase.UserType, userBase.CompleteDegree); err != nil {
//		log.Error("incr user base err(%v)", err)
//		return
//	}
//	d.DelCacheUserBase(c, userBase.Mid)
//	return res.LastInsertId()
//}
//

// UpdateUserField .
func (d *Dao) UpdateUserField(c context.Context, tx *sql.Tx, mid int64, field string, f interface{}) (num int64, err error) {
	var res xsql.Result
	querySQL := fmt.Sprintf(_updateUserField, field)
	if res, err = tx.Exec(querySQL, f, mid); err != nil {
		log.Errorv(c, log.KV("log", fmt.Sprintf("update user mid(%d) field(%s) value(%v) err(%v)", mid, field, f, err)))
		return
	}
	log.V(1).Infow(c, "log", "update user field", "mid", mid, "field", field, "value", f)
	d.DelCacheUserBase(c, mid)
	return res.RowsAffected()
}

//AddUserBase .
func (d *Dao) AddUserBase(c context.Context, userBase *api.UserBase) (num int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, _incrUserBase, userBase.Mid, userBase.Uname, userBase.Face, userBase.Birthday, userBase.Exp, userBase.Level, userBase.UserType, userBase.CompleteDegree, userBase.Sex); err != nil {
		log.Error("incr user base err(%v)", err)
		return
	}
	d.DelCacheUserBase(c, userBase.Mid)
	return res.LastInsertId()
}

//UpdateUserBaseUname .
func (d *Dao) UpdateUserBaseUname(c context.Context, mid int64, uname string) (num int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, _updateUserUname, uname, mid); err != nil {
		log.Error("update user base uname err(%v)", err)
		return
	}
	d.DelCacheUserBase(c, mid)
	return res.RowsAffected()
}

// UpdateUserBase 更新用户基础信息
func (d *Dao) UpdateUserBase(c context.Context, mid int64, userBase *api.UserBase) (num int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, _updateUser, userBase.Uname, userBase.Face, userBase.Birthday, userBase.Sex,
		userBase.Region, userBase.Signature, userBase.CompleteDegree, mid); err != nil {
		log.Errorv(c, log.KV("event", "mysql_update"), log.KV("mid", mid), log.KV("error", err))
		return
	}
	d.DelCacheUserBase(c, mid)
	return res.RowsAffected()
}

// CheckUname 检测昵称
func (d *Dao) CheckUname(c context.Context, mid int64, uname string) (err error) {
	// 前缀不能为Qing_
	if strings.HasPrefix(uname, "Qing_") {
		err = ecode.UserUnamePrefixErr
		return
	}
	//特殊字符
	if !model.CheckUnameSpecial(uname) {
		err = ecode.UserUnameSpecial
		return
	}
	//字符长度
	if !model.CheckUnameLength(uname) {
		err = ecode.UserUnameLength
		return
	}
	//bbq是否存在
	tmp := int64(0)
	row := d.db.QueryRow(c, _selectUname, uname, mid)
	if err = row.Scan(&tmp); err != nil && err != sql.ErrNoRows {
		err = ecode.EditUserBaseErr
		return
	}
	if tmp != 0 {
		err = ecode.UserUnameExisted
		log.Infow(c, "log", "uname已存在", "uname", uname, "mid", mid)
		return
	}
	//自己b站的昵称
	var userCard *model.UserCard
	if userCard, err = d.RawUserCard(c, mid); err != nil {
		err = ecode.EditUserBaseErr
		return
	}
	if userCard.Name == uname {
		return nil
	}
	//是否是万粉的昵称
	row2 := d.db.QueryRow(c, _selectBZhanUpUname, uname, mid)
	if tmpErr := row2.Scan(&tmp); tmpErr != nil && tmpErr != sql.ErrNoRows {
		err = ecode.EditUserBaseErr
		log.V(1).Infow(c, "log", "获取B站万粉资料为空", "uname", uname)
		return
	}
	if tmp != 0 {
		err = ecode.UserUnameExisted
		log.Infow(c, "log", "uname命中B站万粉up主", "uname", uname, "mid", mid)
		return
	}

	// 昵称是否包含敏感词
	level, filterErr := d.Filter(c, uname, "BBQ_account")
	if filterErr != nil {
		log.Errorv(c, log.KV("log", "filter fail"))
	} else if level >= 20 {
		err = ecode.UserUnameFilterErr
		log.Warnv(c, log.KV("log", fmt.Sprintf("uname filter fail: uname=%s, level=%d", uname, level)))
		return
	}

	// 运营不允许使用的uname列表中
	if conf.UnameConf.UnameForbidden(uname) {
		log.Infow(c, "log", "hit fobidden uname", "uname", uname)
		err = ecode.UserUnameExisted
		return
	}

	return nil
}
