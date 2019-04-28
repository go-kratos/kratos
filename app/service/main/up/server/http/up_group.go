package http

import (
	"go-common/app/service/main/up/model"
	"go-common/app/service/main/up/service"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/http/blademaster"

	"github.com/go-sql-driver/mysql"
)

func addGroup(c *blademaster.Context) {
	var res interface{}
	var err error
	var errMsg string
	var r = new(model.AddGroupArg)
	switch {
	default:
		if err = c.Bind(r); err != nil {
			errMsg = err.Error()
			err = ecode.RequestErr
			log.Error("request error, err=%v, req=%v", err, r)
			break
		}

		dbresult, err2 := Svc.AddGroup(c, r)
		if err2 != nil {
			log.Error("add group error, err=%v, req=%v", err2, r)
			if _, ok := err2.(*mysql.MySQLError); ok {
				errMsg = "db error"
			} else {
				errMsg = err2.Error()
			}
			err = ecode.ServerErr
			break
		}
		var lastID, _ = dbresult.LastInsertId()
		log.Info("add group ok, req=%+v, last id=%d", r, lastID)
	}
	if err != nil {
		service.BmHTTPErrorWithMsg(c, err, errMsg)
	} else {
		c.JSON(res, err)
	}
}

func updateGroup(c *blademaster.Context) {
	var res interface{}
	var err error
	var errMsg string
	var r = new(model.EditGroupArg)
	switch {
	default:
		r.AddArg = new(model.AddGroupArg)
		if err = c.Bind(r.AddArg); err != nil {
			errMsg = err.Error()
			err = ecode.RequestErr
			log.Error("request error, err=%v, req=%v", err, r)
			break
		}

		if err = c.Bind(r); err != nil {
			errMsg = err.Error()
			err = ecode.RequestErr
			log.Error("request error, err=%v, req=%v", err, r)
			break
		}

		dbresult, err2 := Svc.UpdateGroup(c, r)
		if err2 != nil {
			log.Error("update group error, err=%v, req=%v", err, r)
			if me, ok := err2.(*mysql.MySQLError); ok {
				switch me.Number {
				case 1062:
					errMsg = "重复数据已存在"
				default:
					errMsg = "db error"
				}
			} else {
				errMsg = err2.Error()
			}
			err = ecode.ServerErr
			break
		}
		var affected, _ = dbresult.RowsAffected()
		log.Info("update group ok, req=%+v, affected count=%d", r, affected)
	}
	if err != nil {
		service.BmHTTPErrorWithMsg(c, err, errMsg)
	} else {
		c.JSON(res, err)
	}
}

func removeGroup(c *blademaster.Context) {
	var res interface{}
	var err error
	var errMsg string
	var r = new(model.RemoveGroupArg)
	switch {
	default:
		if err = c.Bind(r); err != nil {
			errMsg = err.Error()
			err = ecode.RequestErr
			log.Error("request error, err=%v, req=%v", err, r)
			break
		}

		dbresult, err2 := Svc.RemoveGroup(c, r)
		if err2 != nil {
			log.Error("remove group error, err=%v, req=%v", err, r)
			errMsg = err2.Error()
			err = ecode.ServerErr
			break
		}
		var affected, _ = dbresult.RowsAffected()
		log.Info("remove group ok, req=%+v, affected count=%d", r, affected)
	}
	if err != nil {
		service.BmHTTPErrorWithMsg(c, err, errMsg)
	} else {
		c.JSON(res, err)
	}
}

func getGroup(c *blademaster.Context) {
	var res interface{}
	var err error
	var errMsg string
	var r = &model.GetGroupArg{State: 1}
	switch {
	default:
		if err = c.Bind(r); err != nil {
			errMsg = err.Error()
			err = ecode.RequestErr
			log.Error("request error, err=%v, req=%v", err, r)
			break
		}
		groups, e := Svc.GetGroup(c, r)
		if e != nil {
			log.Error("get group error, err=%v", e)
			errMsg = e.Error()
			err = ecode.ServerErr
			break
		}
		res = groups
		log.Info("get group ok, count=%d", len(groups))
	}
	if err != nil {
		service.BmHTTPErrorWithMsg(c, err, errMsg)
	} else {
		c.JSON(res, err)
	}
}
