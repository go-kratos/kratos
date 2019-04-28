package dao

import (
	"context"

	"go-common/app/service/main/vip/model"
	"go-common/library/database/sql"

	"github.com/pkg/errors"
)

const (
	privilegeListSQL          = "SELECT id, privileges_name, title, explains, privileges_type, operator, state, deleted, icon_url, icon_gray_url, order_num, lang_type, ctime, mtime FROM vip_privileges WHERE deleted = 0 AND state = 1 ORDER BY order_num;"
	privilegeResourcesListSQL = "SELECT id,pid,link,image_url,resources_type,ctime,mtime FROM vip_privileges_resources;"
)

// PrivilegeList query .
func (d *Dao) PrivilegeList(c context.Context) (res []*model.Privilege, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, privilegeListSQL); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.Privilege)
		if err = rows.Scan(&r.ID, &r.Name, &r.Title, &r.Explain, &r.Type, &r.Operator, &r.State, &r.Deleted, &r.IconURL, &r.IconGrayURL, &r.Order, &r.LangType, &r.Ctime, &r.Mtime); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

// PrivilegeResourcesList query privilege resources .
func (d *Dao) PrivilegeResourcesList(c context.Context) (res []*model.PrivilegeResources, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, privilegeResourcesListSQL); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.PrivilegeResources)
		if err = rows.Scan(&r.ID, &r.PID, &r.Link, &r.ImageURL, &r.Type, &r.Ctime, &r.Mtime); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}
