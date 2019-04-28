package http

import (
	"net/http"

	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/render"

	"go-common/library/ecode"
)

// check username and dashboard sessionid
func checkCookie(c *bm.Context) (username, sid string, err error) {
	var r = c.Request
	var name *http.Cookie
	if name, err = r.Cookie("username"); err == nil {
		username = name.Value
	}
	var session *http.Cookie
	if session, err = r.Cookie("_AJSESSIONID"); err == nil {
		sid = session.Value
	}
	if username == "" || sid == "" {
		err = ecode.Unauthorized
	}
	return
}

func getAuthorityUserPrivileges(c *bm.Context) {
	username, _, err := checkCookie(c)
	if err != nil {
		c.JSON(nil, err)
		log.Error("growup checkCookie error(%v)", err)
		return
	}

	data, err := svr.GetAuthorityUserPrivileges(username)
	if err != nil {
		c.JSON(nil, err)
		log.Error("growup svr.GetAuthorityUserPrivileges error(%v)", err)
		return
	}
	c.JSON(data, nil)
}

func getAuthorityUserGroup(c *bm.Context) {
	username, _, err := checkCookie(c)
	if err != nil {
		c.JSON(nil, err)
		log.Error("growup checkCookie error(%v)", err)
		return
	}
	data, err := svr.GetAuthorityUserGroup(username)
	if err != nil {
		c.JSON(nil, err)
		log.Error("growup svr.GetAuthorityUserGroup error(%v)", err)
		return
	}
	c.JSON(data, nil)
}

func listAuthorityUsers(c *bm.Context) {
	v := new(struct {
		Username string `form:"username"`
		From     int    `form:"from" validate:"min=0" default:"0"`
		Limit    int    `form:"limit" validate:"min=1" default:"20"`
		Sort     string `form:"sort"`
	})

	if err := c.Bind(v); err != nil {
		return
	}
	users, total, err := svr.ListAuthorityUsers(v.Username, v.From, v.Limit, v.Sort)
	if err != nil {
		c.JSON(nil, err)
		log.Error("growup svr.ListAuthorityUsers error(%v)", err)
		return
	}
	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code":    0,
		"message": "0",
		"data":    users,
		"paging": map[string]int{
			"page_size": v.Limit,
			"total":     total,
		},
	}))
}

func addAuthorityUser(c *bm.Context) {
	v := new(struct {
		Username string `form:"username"`
		Nickname string `form:"nickname"`
	})

	if err := c.Bind(v); err != nil {
		return
	}
	err := svr.AddAuthorityUser(v.Username, v.Nickname)
	if err != nil {
		log.Error("growup svr.AddAuthorityUser error(%v)", err)
	}
	c.JSON(nil, err)
}

func updateAuthorityUserInfo(c *bm.Context) {
	v := new(struct {
		ID       int64  `form:"id"`
		Nickname string `form:"nickname"`
	})

	if err := c.Bind(v); err != nil {
		return
	}
	err := svr.UpdateAuthorityUserInfo(v.ID, v.Nickname)
	if err != nil {
		log.Error("growup svr.UpdateAuthorityUserInfo error(%v)", err)
	}
	c.JSON(nil, err)
}

func updateAuthorityUserAuth(c *bm.Context) {
	v := new(struct {
		ID      int64  `form:"id"`
		GroupID string `form:"group_id"`
		RoleID  string `form:"role_id"`
	})

	if err := c.Bind(v); err != nil {
		return
	}
	err := svr.UpdateAuthorityUserAuth(v.ID, v.GroupID, v.RoleID)
	if err != nil {
		log.Error("growup svr.UpdateAuthorityUserAuth error(%v)", err)
	}
	c.JSON(nil, err)
}

func deleteAuthorityUser(c *bm.Context) {
	v := new(struct {
		ID int64 `form:"id"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	err := svr.DeleteAuthorityUser(v.ID)
	if err != nil {
		log.Error("growup svr.DeleteAuthorityUser error(%v)", err)
	}
	c.JSON(nil, err)
}

func listAuthorityTaskGroups(c *bm.Context) {
	v := new(struct {
		From  int    `form:"from" validate:"min=0" default:"0"`
		Limit int    `form:"limit" validate:"min=1" default:"20"`
		Sort  string `form:"sort"`
	})

	if err := c.Bind(v); err != nil {
		return
	}

	groups, total, err := svr.ListAuthorityTaskGroups(v.From, v.Limit, v.Sort)
	if err != nil {
		c.JSON(nil, err)
		log.Error("growup svr.ListAuthorityTaskGroups error(%v)", err)
		return
	}
	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code":    0,
		"message": "0",
		"data":    groups,
		"paging": map[string]int{
			"page_size": v.Limit,
			"total":     total,
		},
	}))
}

func addAuthorityTaskGroup(c *bm.Context) {
	v := new(struct {
		Name string `form:"name"`
		Desc string `form:"desc"`
	})

	if err := c.Bind(v); err != nil {
		return
	}

	err := svr.AddAuthorityTaskGroup(v.Name, v.Desc)
	if err != nil {
		log.Error("growup svr.AddAuthorityTaskGroup error(%v)", err)
	}
	c.JSON(nil, err)
}

func addAuthorityTaskGroupUser(c *bm.Context) {
	v := new(struct {
		Username string `form:"username"`
		GroupID  string `form:"group_id"`
	})

	if err := c.Bind(v); err != nil {
		return
	}

	err := svr.AddAuthorityTaskGroupUser(v.Username, v.GroupID)
	if err != nil {
		log.Error("growup svr.AddAuthorityTaskGroupUser error(%v)", err)
	}
	c.JSON(nil, err)
}

func updateAuthorityTaskGroupInfo(c *bm.Context) {
	v := new(struct {
		GroupID int64  `form:"group_id"`
		Name    string `form:"name"`
		Desc    string `form:"desc"`
	})

	if err := c.Bind(v); err != nil {
		return
	}

	err := svr.UpdateAuthorityTaskGroupInfo(v.GroupID, v.Name, v.Desc)
	if err != nil {
		log.Error("growup svr.UpdateAuthorityTaskGroupInfo error(%v)", err)
	}
	c.JSON(nil, err)
}

func deleteAuthorityTaskGroup(c *bm.Context) {
	v := new(struct {
		GroupID int64 `form:"group_id"`
	})

	if err := c.Bind(v); err != nil {
		return
	}

	err := svr.DeleteAuthorityTaskGroup(v.GroupID)
	if err != nil {
		log.Error("growup svr.DeleteAuthorityTaskGroup error(%v)", err)
	}
	c.JSON(nil, err)
}

func deleteAuthorityTaskGroupUser(c *bm.Context) {
	v := new(struct {
		ID      int64 `form:"id"`
		GroupID int64 `form:"group_id"`
	})

	if err := c.Bind(v); err != nil {
		return
	}

	err := svr.DeleteAuthorityTaskGroupUser(v.ID, v.GroupID)
	if err != nil {
		log.Error("growup svr.DeleteAuthorityTaskGroupUser error(%v)", err)
	}
	c.JSON(nil, err)
}

func listAuthorityGroupPrivilege(c *bm.Context) {
	v := new(struct {
		GroupID  int64 `form:"group_id"`
		FatherID int64 `form:"father_id" validate:"required"`
	})

	if err := c.Bind(v); err != nil {
		return
	}

	data, err := svr.ListAuthorityGroupPrivilege(v.GroupID, v.FatherID)
	if err != nil {
		log.Error("growup svr.ListAuthorityGroupPrivilege error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func updateAuthorityGroupPrivilege(c *bm.Context) {
	v := new(struct {
		Add     string `form:"add"`
		Minus   string `form:"minus"`
		GroupID int64  `form:"group_id"`
		Type    int    `form:"type"` // 1 数据源
	})
	if err := c.Bind(v); err != nil {
		return
	}

	err := svr.UpdateAuthorityGroupPrivilege(v.GroupID, v.Add, v.Minus, v.Type)
	if err != nil {
		log.Error("growup svr.UpdateAuthorityGroupPrivilege error(%v)", err)
	}
	c.JSON(nil, err)
}

func listAuthorityTaskRoles(c *bm.Context) {
	username, _, err := checkCookie(c)
	if err != nil {
		c.JSON(nil, err)
		log.Error("growup checkCookie error(%v)", err)
		return
	}
	v := new(struct {
		From  int    `form:"from" validate:"min=0" default:"0"`
		Limit int    `form:"limit" validate:"min=1" default:"20"`
		Sort  string `form:"sort"`
	})

	if err = c.Bind(v); err != nil {
		return
	}

	roles, total, err := svr.ListAuthorityTaskRoles(username, v.From, v.Limit, v.Sort)
	if err != nil {
		c.JSON(nil, err)
		log.Error("growup svr.ListAuthorityTaskRoles error(%v)", err)
		return
	}

	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code":    0,
		"message": "0",
		"data":    roles,
		"paging": map[string]int{
			"page_size": v.Limit,
			"total":     total,
		},
	}))
}

func addAuthorityTaskRole(c *bm.Context) {
	v := new(struct {
		GroupID int64  `form:"group_id" validate:"required"`
		Name    string `form:"name"`
		Desc    string `form:"desc"`
	})

	if err := c.Bind(v); err != nil {
		return
	}

	err := svr.AddAuthorityTaskRole(v.GroupID, v.Name, v.Desc)
	if err != nil {
		log.Error("growup svr.AddAuthorityTaskRole error(%v)", err)
	}
	c.JSON(nil, err)
}

func addAuthorityTaskRoleUser(c *bm.Context) {
	v := new(struct {
		Username string `form:"username"`
		RoleID   string `form:"role_id"`
	})

	if err := c.Bind(v); err != nil {
		return
	}

	err := svr.AddAuthorityTaskRoleUser(v.Username, v.RoleID)
	if err != nil {
		log.Error("growup svr.AddAuthorityTaskRoleUser error(%v)", err)
	}
	c.JSON(nil, err)
}

func updateAuthorityTaskRoleInfo(c *bm.Context) {
	v := new(struct {
		RoleID int64  `form:"role_id"`
		Name   string `form:"name"`
		Desc   string `form:"desc"`
	})

	if err := c.Bind(v); err != nil {
		return
	}

	err := svr.UpdateAuthorityTaskRoleInfo(v.RoleID, v.Name, v.Desc)
	if err != nil {
		log.Error("growup svr.UpdateAuthorityTaskRoleInfo error(%v)", err)
	}
	c.JSON(nil, err)
}

func deleteAuthorityTaskRole(c *bm.Context) {
	v := new(struct {
		RoleID int64 `form:"role_id"`
	})

	if err := c.Bind(v); err != nil {
		return
	}

	err := svr.DeleteAuthorityTaskRole(v.RoleID)
	if err != nil {
		log.Error("growup svr.DeleteAuthorityTaskRole error(%v)", err)
	}
	c.JSON(nil, err)
}

func deleteAuthorityTaskRoleUser(c *bm.Context) {
	v := new(struct {
		ID     int64 `form:"id"`
		RoleID int64 `form:"role_id"`
	})

	if err := c.Bind(v); err != nil {
		return
	}

	err := svr.DeleteAuthorityTaskRoleUser(v.ID, v.RoleID)
	if err != nil {
		log.Error("growup svr.DeleteAuthorityTaskRoleUser error(%v)", err)
	}
	c.JSON(nil, err)
}

func listAuthorityRolePrivilege(c *bm.Context) {
	v := new(struct {
		GroupID  int64 `form:"group_id"`
		RoleID   int64 `form:"role_id"`
		FatherID int64 `form:"father_id" validate:"required"`
	})

	if err := c.Bind(v); err != nil {
		return
	}

	data, err := svr.ListAuthorityRolePrivilege(v.GroupID, v.RoleID, v.FatherID)
	if err != nil {
		c.JSON(nil, err)
		log.Error("growup svr.ListAuthorityRolePrivilege error(%v)", err)
		return
	}
	c.JSON(data, nil)
}

func updateAuthorityRolePrivilege(c *bm.Context) {
	v := new(struct {
		Add    string `form:"add"`
		Minus  string `form:"minus"`
		RoleID int64  `form:"role_id"`
	})

	if err := c.Bind(v); err != nil {
		return
	}

	err := svr.UpdateAuthorityRolePrivilege(v.RoleID, v.Add, v.Minus)
	if err != nil {
		log.Error("growup svr.UpdateAuthorityRolePrivilege error(%v)", err)
	}
	c.JSON(nil, err)
}

func listAuthorityGroupAndRole(c *bm.Context) {
	groups, roles, err := svr.ListGroupAndRole()
	if err != nil {
		c.JSON(nil, err)
		log.Error("growup svr.ListGroupAndRole error(%v)", err)
		return
	}
	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code":    0,
		"message": "0",
		"data": map[string]interface{}{
			"groups": groups,
			"roles":  roles,
		},
	}))
}

func listPrivilege(c *bm.Context) {
	data, err := svr.ListPrivilege()
	if err != nil {
		c.JSON(nil, err)
		log.Error("growup svr.AddPrivilege error(%v)", err)
		return
	}
	c.JSON(data, nil)
}

func addPrivilege(c *bm.Context) {
	v := new(struct {
		Level    int64  `form:"level" validate:"required"`
		Name     string `form:"name" validate:"required"`
		FatherID int64  `form:"father_id"`
		IsRouter uint8  `form:"is_router"`
	})

	if err := c.Bind(v); err != nil {
		return
	}
	if v.Level > 1 && v.FatherID == 0 {
		c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
			"code":    ecode.RequestErr,
			"message": "privilege > 1 but father_id = 0",
		}))
		return
	}
	err := svr.AddPrivilege(v.Name, v.Level, v.FatherID, v.IsRouter)
	if err != nil {
		log.Error("growup svr.AddPrivilege error(%v)", err)
	}
	c.JSON(nil, err)
}

func updatePrivilege(c *bm.Context) {
	v := new(struct {
		ID       int64  `form:"id" validate:"required"`
		Level    int64  `form:"level" validate:"required"`
		Name     string `form:"name" validate:"required"`
		FatherID int64  `form:"father_id"`
		IsRouter uint8  `form:"is_router"`
	})
	if err := c.Bind(v); err != nil {
		return
	}

	if v.Level > 1 && v.FatherID == 0 {
		c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
			"code":    ecode.RequestErr,
			"message": "privilege > 1 but father_id = 0",
		}))
		return
	}

	err := svr.UpdatePrivilege(v.ID, v.Name, v.Level, v.FatherID, v.IsRouter)
	if err != nil {
		log.Error("growup svr.UpdatePrivilege error(%v)", err)
	}
	c.JSON(nil, err)
}

func busPrivilege(c *bm.Context) {
	username, _, err := checkCookie(c)
	if err != nil {
		c.JSON(nil, err)
		log.Error("growup checkCookie error(%v)", err)
		return
	}
	v := new(struct {
		Type string `form:"type"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	data, err := svr.BusPrivilege(c, username, v.Type)
	if err != nil {
		c.JSON(nil, err)
		log.Error("growup svr.BusPrivilege error(%v)", err)
		return
	}
	c.JSON(data, nil)
}
