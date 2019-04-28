package http

import (
	"strconv"

	"go-common/app/admin/main/manager/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

// addBusiness .
func addBusiness(c *bm.Context) {
	b := &model.Business{}
	if err := c.Bind(b); err != nil {
		return
	}
	c.JSON(nil, mngSvc.AddBusiness(c, b))
}

// updateBusiness .
func updateBusiness(c *bm.Context) {
	b := &model.Business{}
	if err := c.Bind(b); err != nil {
		return
	}
	if b.ID <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, mngSvc.UpdateBusiness(c, b))
}

// addRole .
func addRole(c *bm.Context) {
	br := &model.BusinessRole{}
	if err := c.Bind(br); err != nil {
		return
	}
	c.JSON(nil, mngSvc.AddRole(c, br))
}

// updateRole .
func updateRole(c *bm.Context) {
	br := &model.BusinessRole{}
	if err := c.Bind(br); err != nil {
		return
	}
	if br.ID <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, mngSvc.UpdateRole(c, br))
}

// addUser .
func addUser(c *bm.Context) {
	res := map[string]interface{}{}
	bur := &model.BusinessUserRole{}
	if err := c.Bind(bur); err != nil {
		return
	}
	if cuid, exists := c.Get("uid"); exists {
		bur.CUID = cuid.(int64)
	}
	if err := mngSvc.AddUser(c, bur); err != nil {
		if err == ecode.ManagerUIDNOTExist {
			res["message"] = "这个uid是不存在的, 请你不要乱来哦!"
			c.JSONMap(res, ecode.ManagerUIDNOTExist)
			return
		}
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

// updateUser .
func updateUser(c *bm.Context) {
	bur := &model.BusinessUserRole{}
	if err := c.Bind(bur); err != nil {
		return
	}
	if bur.ID <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, mngSvc.UpdateUser(c, bur))
}

// stateUpdate .
func updateState(c *bm.Context) {
	su := &model.StateUpdate{}
	res := map[string]interface{}{}
	if err := c.Bind(su); err != nil {
		return
	}
	if err := mngSvc.UpdateState(c, su); err != nil {
		if err == ecode.ManagerFlowForbidden {
			res["message"] = "business flow is forbidden"
			c.JSONMap(res, ecode.ManagerFlowForbidden)
			return
		}
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

// businessList .
func businessList(c *bm.Context) {
	bp := &model.BusinessListParams{}
	if err := c.Bind(bp); err != nil {
		return
	}
	if uid, exists := c.Get("uid"); exists {
		bp.UID = uid.(int64)
	}
	c.JSON(mngSvc.BusinessList(c, bp))
}

// flowList .
func flowList(c *bm.Context) {
	bp := &model.BusinessListParams{}
	if err := c.Bind(bp); err != nil {
		return
	}
	if uid, exists := c.Get("uid"); exists {
		bp.UID = uid.(int64)
	}
	c.JSON(mngSvc.FlowList(c, bp))
}

// userList .
func userList(c *bm.Context) {
	isAdmin := false
	u := &model.UserListParams{}
	if err := c.Bind(u); err != nil {
		return
	}
	if uid, exists := c.Get("uid"); exists {
		isAdmin = mngSvc.IsAdmin(c, u.BID, uid.(int64))
	}
	data, total, err := mngSvc.UserList(c, u)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	page := map[string]int64{
		"num":   u.PN,
		"size":  u.PS,
		"total": total,
	}
	c.JSON(map[string]interface{}{
		"page":    page,
		"data":    data,
		"isAdmin": isAdmin,
	}, err)
}

// roleList .
func roleList(c *bm.Context) {
	br := &model.BusinessRole{}
	if err := c.Bind(br); err != nil {
		return
	}
	if uid, exists := c.Get("uid"); exists {
		br.UID = uid.(int64)
	}
	c.JSON(mngSvc.RoleList(c, br))
}

// deleteUser .
func deleteUser(c *bm.Context) {
	bur := &model.BusinessUserRole{}
	if err := c.Bind(bur); err != nil {
		return
	}
	if bur.ID <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, mngSvc.DeleteUser(c, bur))
}

// userRole .
func userRole(c *bm.Context) {
	brl := &model.BusinessUserRoleList{}
	if err := c.Bind(brl); err != nil {
		return
	}
	if brl.BID <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if brl.UID <= 0 {
		if uid, exists := c.Get("uid"); exists {
			uid = uid.(int64)
		}
	}
	c.JSON(mngSvc.UserRole(c, brl))
}

// userRoles .
func userRoles(c *bm.Context) {
	uid, _ := strconv.ParseInt(c.Request.Form.Get("uid"), 10, 64)
	if uid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(mngSvc.UserRoles(c, uid))
}

// userRoles .
func stateUp(c *bm.Context) {
	p := &model.UserStateUp{}
	if err := c.Bind(p); err != nil {
		return
	}
	c.JSON(nil, mngSvc.StateUp(c, p))
}

// allRoles .
func allRoles(c *bm.Context) {
	pid, _ := strconv.ParseInt(c.Request.Form.Get("pid"), 10, 64)
	if pid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	uid, _ := strconv.ParseInt(c.Request.Form.Get("uid"), 10, 64)
	if uid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(mngSvc.AllRoles(c, pid, uid))
}
