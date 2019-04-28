package manager

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	upgrpc "go-common/app/service/main/up/api/v1"
	"go-common/app/service/main/up/dao"
	"go-common/app/service/main/up/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_insertGroup     = "INSERT INTO up_group (name, tag, short_tag, colors, remark) VALUES (?,?,?,?,?)"
	_checkGroupExist = "SELECT COUNT(0) FROM up_group WHERE (name=? OR tag=? OR short_tag=?) "
	_updateGroupByID = "UPDATE up_group SET name=?, tag=?, short_tag=?, colors=?, remark=? WHERE id=?"
	_removeGroupByID = "UPDATE up_group SET state=0 WHERE id=?"
	_selectGroup     = "SELECT id, name, tag, short_tag, remark, colors, state FROM up_group "
	_upGroupsSQL     = "SELECT id, name, tag, short_tag, remark, colors FROM up_group WHERE state = 1"
)

//AddGroup add group in db
func (d *Dao) AddGroup(c context.Context, groupAddInfo *model.AddGroupArg) (res sql.Result, err error) {
	var color = fmt.Sprintf("%s|%s", groupAddInfo.FontColor, groupAddInfo.BgColor)
	res, err = prepareAndExec(c, d.managerDB, _insertGroup, groupAddInfo.Name, groupAddInfo.Tag, groupAddInfo.ShortTag, color, groupAddInfo.Remark)
	return
}

//CheckGroupExist check if group exist
func (d *Dao) CheckGroupExist(c context.Context, groupAddInfo *model.AddGroupArg, exceptid int64) (exist bool, err error) {
	var sqlStr = _checkGroupExist
	var args = []interface{}{groupAddInfo.Name, groupAddInfo.Tag, groupAddInfo.ShortTag}
	if exceptid != 0 {
		sqlStr += " AND id != ?"
		args = append(args, exceptid)
	}
	rows, err := prepareAndQuery(c, d.managerDB, sqlStr, args...)
	if err != nil {
		return
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		err = rows.Scan(&count)
		break
	}
	if err != nil {
		return
	}
	exist = count > 0
	return
}

//UpdateGroup update group
func (d *Dao) UpdateGroup(c context.Context, groupAddInfo *model.EditGroupArg) (res sql.Result, err error) {
	if groupAddInfo.AddArg == nil {
		return
	}
	var color = fmt.Sprintf("%s|%s", groupAddInfo.AddArg.FontColor, groupAddInfo.AddArg.BgColor)
	res, err = prepareAndExec(c, d.managerDB, _updateGroupByID, groupAddInfo.AddArg.Name, groupAddInfo.AddArg.Tag, groupAddInfo.AddArg.ShortTag, color, groupAddInfo.AddArg.Remark, groupAddInfo.ID)
	return
}

//RemoveGroup remove group
func (d *Dao) RemoveGroup(c context.Context, arg *model.RemoveGroupArg) (res sql.Result, err error) {
	res, err = prepareAndExec(c, d.managerDB, _removeGroupByID, arg.ID)
	return
}

//GetGroup get group
func (d *Dao) GetGroup(c context.Context, arg *model.GetGroupArg) (res []*model.UpGroup, err error) {

	var con = dao.Condition{
		Key:      "state",
		Operator: "=",
		Value:    arg.State,
	}

	var conditionStr, vals, _ = dao.ConcatCondition(con)

	rows, err := prepareAndQuery(c, d.managerDB, _selectGroup+"WHERE "+conditionStr, vals...)
	if err != nil {
		return
	}
	defer rows.Close()

	// id, name, tag, short_tag, remark, colors
	var colorStr string
	for rows.Next() {
		var group = model.UpGroup{}
		err = rows.Scan(&group.ID, &group.Name, &group.Tag, &group.ShortTag, &group.Remark, &colorStr, &group.State)
		if err != nil {
			log.Error("scan row failed, err=%v", err)
			break
		}
		var colors = strings.Split(colorStr, "|")
		if len(colors) >= 2 {
			group.FontColor = colors[0]
			group.BgColor = colors[1]
		}
		res = append(res, &group)
	}
	return
}

// UpGroups get up special group data.
func (d *Dao) UpGroups(c context.Context) (mug map[int64]*upgrpc.UpGroup, err error) {
	rows, err := d.managerDB.Query(c, _upGroupsSQL)
	if err != nil {
		return
	}
	defer rows.Close()
	mug = make(map[int64]*upgrpc.UpGroup)
	for rows.Next() {
		var (
			colorStr string
			colors   []string
			ug       = new(upgrpc.UpGroup)
		)
		if err = rows.Scan(&ug.ID, &ug.Name, &ug.Tag, &ug.ShortTag, &ug.Note, &colorStr); err != nil {
			if err == xsql.ErrNoRows {
				err = nil
			}
			return
		}
		colors = strings.Split(colorStr, "|")
		if len(colors) >= 2 {
			ug.FontColor = colors[0]
			ug.BgColor = colors[1]
		}
		mug[ug.ID] = ug
	}
	return
}
