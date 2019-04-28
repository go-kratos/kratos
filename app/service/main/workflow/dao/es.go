package dao

import (
	"context"

	"go-common/library/log"
)

// UpdateUserTag 更新举报对象的第一个用户tag
func (d *Dao) UpdateUserTag(c context.Context, gid int64, userTid int32) (err error) {
	us := d.es.NewUpdate("workflow_group_common").Insert()
	data := map[string]int64{"id": gid, "first_user_tid": int64(userTid)}
	us.AddData("workflow_group_common", data)
	if err = us.Do(c); err != nil {
		log.Error("failed upsert elastic error(%v)", err)
		return
	}
	log.Info("gid(%d) has first user complain tid(%d) success upsert", gid, userTid)
	return
}
