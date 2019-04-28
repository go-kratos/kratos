package dao

import (
	"context"

	"go-common/app/service/main/vip/model"
	"go-common/library/database/sql"

	"github.com/pkg/errors"
)

const (
	_effectiveAssociateVipsSQL = "SELECT title,remark,link,associate_platform FROM  vip_associate_vip WHERE state = 0 AND deleted = 0 ORDER BY order_num DESC;"
)

// EffectiveAssociateVips effective associate vips.
func (d *Dao) EffectiveAssociateVips(c context.Context) (res []*model.AssociateVipResp, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _effectiveAssociateVipsSQL); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.AssociateVipResp)
		if err = rows.Scan(&r.Title, &r.Remark, &r.Link, &r.AssociatePlatform); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}
