package http

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/admin/main/aegis/model/task"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_getUIDs        = "/x/admin/manager/users/uids"
	_getUname       = "/x/admin/manager/users/unames"
	_getUdepartment = "/x/admin/manager/users/udepts"
	_getRole        = "/x/admin/manager/internal/user/role"
	_getRoles       = "/x/admin/manager/internal/user/roles"
)

// GetRole 获取用户对应业务下的角色
func (d *Dao) GetRole(c context.Context, bid, uid int64) (roles []*task.Role, err error) {
	roles = []*task.Role{}
	params := url.Values{}
	params.Set("bid", strconv.FormatInt(bid, 10))
	params.Set("uid", strconv.FormatInt(uid, 10))

	res := new(struct {
		Code int          `json:"code"`
		Data []*task.Role `json:"data"`
	})

	if err = d.clientR.Get(c, d.c.Host.Manager+_getRole, "", params, res); err != nil {
		log.Error("GetRole error(%v) url(%s) params(%s)", err, d.c.Host.Manager+_getRole, params.Encode())
		return
	}
	if res.Code != 0 {
		log.Error("GetRole request failed res(%+v) url(%s) params(%s)", res, d.c.Host.Manager+_getRole, params.Encode())
		return
	}
	if len(res.Data) == 0 {
		return
	}
	roles = res.Data
	return
}

// GetUserRoles 获取用户在所有业务下的角色
func (d *Dao) GetUserRoles(c context.Context, uid int64) (roles []*task.Role, err error) {
	roles = []*task.Role{}
	params := url.Values{}
	params.Set("uid", strconv.FormatInt(uid, 10))
	res := new(struct {
		Code int          `json:"code"`
		Data []*task.Role `json:"data"`
	})

	if err = d.clientR.Get(c, d.c.Host.Manager+_getRoles, "", params, res); err != nil {
		log.Error("GetUserRoles error(%v) url(%s) params(%s)", err, d.c.Host.Manager+_getRole, params.Encode())
		return
	}
	if res.Code != 0 {
		log.Error("GetUserRoles request failed res(%+v) url(%s) params(%s)", res, d.c.Host.Manager+_getRole, params.Encode())
		return
	}
	if len(res.Data) == 0 {
		return
	}

	roles = res.Data
	return
}

// GetUnames .
func (d *Dao) GetUnames(c context.Context, uids []int64) (unames map[int64]string, err error) {
	unames = map[int64]string{}
	if len(uids) == 0 {
		return
	}
	params := url.Values{}
	params.Set("uids", xstr.JoinInts(uids))

	res := new(struct {
		Code int              `json:"code"`
		Data map[int64]string `json:"data"`
	})

	if err = d.clientR.Get(c, d.c.Host.Manager+_getUname, "", params, res); err != nil {
		log.Error("GetUnames error(%v) url(%s) params(%s)", err, d.c.Host.Manager+_getUname, params.Encode())
		return
	}
	if res.Code != 0 {
		log.Error("GetUnames request failed res(%+v) url(%s) params(%s)", res, d.c.Host.Manager+_getUname, params.Encode())
		return
	}
	if len(res.Data) == 0 {
		return
	}

	unames = res.Data
	return
}

// GetUIDs .
func (d *Dao) GetUIDs(c context.Context, unames string) (uids map[string]int64, err error) {
	uids = map[string]int64{}
	params := url.Values{}
	params.Set("unames", unames)

	res := new(struct {
		Code int              `json:"code"`
		Data map[string]int64 `json:"data"`
	})

	if err = d.clientR.Get(c, d.c.Host.Manager+_getUIDs, "", params, res); err != nil {
		log.Error("GetUIDs error(%v) url(%s) params(%s)", err, d.c.Host.Manager+_getUname, params.Encode())
		return
	}
	if res.Code != 0 {
		log.Error("GetUIDs request failed res(%+v) url(%s) params(%s)", res, d.c.Host.Manager+_getUname, params.Encode())
		return
	}
	if len(res.Data) == 0 {
		return
	}

	uids = res.Data
	return
}

// GetUdepartment .
func (d *Dao) GetUdepartment(c context.Context, uids []int64) (udepartment map[int64]string, err error) {
	udepartment = map[int64]string{}
	if len(uids) == 0 {
		return
	}
	params := url.Values{}
	params.Set("uids", xstr.JoinInts(uids))

	res := new(struct {
		Code int              `json:"code"`
		Data map[int64]string `json:"data"`
	})

	if err = d.clientR.Get(c, d.c.Host.Manager+_getUdepartment, "", params, res); err != nil {
		log.Error("GetUdepartment error(%v) url(%s) params(%s)", err, d.c.Host.Manager+_getUdepartment, params.Encode())
		return
	}
	if res.Code != 0 {
		log.Error("GetUdepartment request failed res(%+v) url(%s) params(%s)", res, d.c.Host.Manager+_getUname, params.Encode())
		return
	}
	if len(res.Data) == 0 {
		return
	}

	udepartment = res.Data
	return
}
