package http

import (
	"encoding/json"
	"go-common/app/admin/main/videoup/model/oversea"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"strconv"
	"strings"
)

// policyGroups 策略组列表
func policyGroups(c *bm.Context) {
	var (
		uid  int64
		err  error
		data *oversea.PolicyGroupData
	)
	v := new(struct {
		UName   string `form:"username"`
		GroupID int64  `form:"group_id"`
		Type    int8   `form:"type"`
		State   int8   `form:"state" default:"-1"`
		Pn      int64  `form:"pn" default:"1"`
		Ps      int64  `form:"ps" default:"20"`
		Order   string `form:"order"`
		Sort    string `form:"sort"`
	})
	if err = c.Bind(v); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if v.UName != "" {
		if uid, err = vdaSvc.GetUID(c, v.UName); err != nil {
			log.Warn("vdaSvc.GetUID(%s) error(%v)", v.UName, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	data, err = vdaSvc.PolicyGroups(c, uid, v.GroupID, v.Type, v.State, v.Ps, v.Pn, v.Order, v.Sort)
	if err != nil {
		log.Error("vdaSvc.PolicyGroups() err(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(data, err)
}

// archiveGroups 稿件的策略组
func archiveGroups(c *bm.Context) {
	var (
		err error
	)
	v := new(struct {
		Aid int64 `form:"aid" validate:"required"`
	})
	if err = c.Bind(v); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	groups, err := vdaSvc.ArchiveGroups(c, v.Aid)
	if err != nil {
		log.Error("vdaSvc.ArchiveGroups() err(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(groups, err)
}

// addPolicyGroup 添加策略组
func addPolicyGroup(c *bm.Context) {
	var (
		v = new(struct {
			Name   string `form:"name" validate:"required"`
			Type   int8   `form:"type" validate:"required"`
			Remark string `form:"remark" default:""`
		})
		group  = &oversea.PolicyGroup{}
		uid, _ = getUIDName(c)
		err    error
	)
	if uid == 0 {
		c.JSON(nil, ecode.Unauthorized)
		return
	}
	if err = c.Bind(v); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if v.Name == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	group.Name = v.Name
	group.Type = v.Type
	group.UID = uid
	group.Remark = v.Remark
	err = vdaSvc.AddPolicyGroup(c, group)
	if err != nil {
		log.Error("vdaSvc.AddPolicyGroup() err(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(group, nil)
}

// editPolicyGroup 编辑策略组
func editPolicyGroup(c *bm.Context) {
	var (
		v = new(struct {
			ID     int64  `form:"id" validate:"required"`
			Name   string `form:"name" validate:"required"`
			Remark string `form:"remark" default:""`
		})
		attrs  = make(map[string]interface{})
		uid, _ = getUIDName(c)
		err    error
	)
	if uid == 0 {
		c.JSON(nil, ecode.Unauthorized)
		return
	}
	if err = c.Bind(v); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	attrs["name"] = v.Name
	attrs["uid"] = uid
	attrs["remark"] = v.Remark
	err = vdaSvc.UpdatePolicyGroup(c, v.ID, attrs)
	if err != nil {
		log.Error("vdaSvc.UpdatePolicyGroup() err(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}

// delPolicyGroups 删除策略组
func delPolicyGroups(c *bm.Context) {
	upPolicyGroupStates(c, oversea.StateDeleted)
}

// restorePolicyGroups 恢复策略组
func restorePolicyGroups(c *bm.Context) {
	upPolicyGroupStates(c, oversea.StateOK)
}

// upPolicyGroupStates 修改策略组状态
func upPolicyGroupStates(c *bm.Context, state int8) {
	var (
		v = new(struct {
			IDStr string `form:"ids" validate:"required"`
		})
		attrs  = make(map[string]interface{})
		intIDs []int64
		strIDs []string
		uid, _ = getUIDName(c)
		err    error
	)
	if uid == 0 {
		c.JSON(nil, ecode.Unauthorized)
		return
	}
	if err = c.Bind(v); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	strIDs = strings.Split(v.IDStr, ",")
	intIDs = make([]int64, len(strIDs))
	for i, id := range strIDs {
		intIDs[i], err = strconv.ParseInt(id, 10, 64)
		if err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	attrs["state"] = state
	attrs["uid"] = uid
	err = vdaSvc.UpdatePolicyGroups(c, intIDs, attrs)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}

// policies 策略组下的策略
func policies(c *bm.Context) {
	var (
		err error
		v   = new(struct {
			Gid int64 `form:"group_id" validate:"required"`
		})
	)
	if err = c.Bind(v); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if v.Gid == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	items, err := vdaSvc.PolicyItems(c, v.Gid)
	if err != nil {
		log.Error("vdaSvc.PolicyItems() err(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(items, nil)
}

// addPolicies 添加策略
func addPolicies(c *bm.Context) {
	var (
		err error
		v   = new(struct {
			Gid  int64  `form:"group_id" validate:"required"`
			JSON string `form:"items" validate:"required"`
		})
		items  []*oversea.PolicyParams
		uid, _ = getUIDName(c)
	)
	if err = c.Bind(v); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = json.Unmarshal([]byte(v.JSON), &items); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if v.Gid == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	err = vdaSvc.AddPolicies(c, uid, v.Gid, items)
	if err != nil {
		log.Error("vdaSvc.AddPolicies() err(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}

// delPolicies 删除策略
func delPolicies(c *bm.Context) {
	var (
		v = new(struct {
			Gid   int64  `form:"group_id"  validate:"required"`
			IDStr string `form:"ids"  validate:"required"`
		})
		uid, _ = getUIDName(c)
		intIDs []int64
		strIDs []string
		err    error
	)
	if err = c.Bind(v); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	strIDs = strings.Split(v.IDStr, ",")
	intIDs = make([]int64, len(strIDs))
	for i, id := range strIDs {
		intIDs[i], err = strconv.ParseInt(id, 10, 64)
		if err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	err = vdaSvc.DelPolices(c, uid, v.Gid, intIDs)
	if err != nil {
		log.Error("vdaSvc.DelPolices() err(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}
