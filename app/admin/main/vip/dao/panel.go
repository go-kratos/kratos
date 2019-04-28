package dao

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go-common/app/admin/main/vip/model"
	xsql "go-common/library/database/sql"
	xtime "go-common/library/time"

	"github.com/pkg/errors"
)

const (
	_inVipPriceConfigSQL = "INSERT INTO vip_price_config_v2(platform,product_name,product_id,suit_type,month,sub_type,original_price,selected,remark,status,operator,oper_id,superscript,start_build,end_build) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"

	_inVipDPriceConfigSQL          = "INSERT INTO vip_price_discount_config_v2(vpc_id,product_id,discount_price,stime,etime,remark,operator,oper_id,first_price) VALUES(?,?,?,?,?,?,?,?,?)"
	_upVipPriceConfigIDSQL         = "UPDATE vip_price_config_v2 SET platform = ?, product_name = ?, product_id = ?, suit_type = ?, month = ?, sub_type = ?, original_price = ?, remark = ?, status = ?, operator = ?, oper_id = ?, selected = ?, superscript = ?, start_build = ?, end_build = ? WHERE id = ?"
	_upVipDPriceConfigSQL          = "UPDATE vip_price_discount_config_v2 SET product_id = ?, discount_price = ?, stime = ?, etime =?, remark = ?, operator = ?, oper_id = ?,first_price = ?  WHERE id = ?"
	_delVipPriceConfigIDSQL        = "DELETE FROM vip_price_config_v2 WHERE id = ?"
	_delVipDPriceConfigIDSQL       = "DELETE FROM vip_price_discount_config_v2 WHERE id = ?"
	_selVipPriceConfigUQCheckSQL   = "SELECT COUNT(*) FROM vip_price_config_v2 WHERE platform = ? AND month = ? AND sub_type = ? AND suit_type = ? %s;"
	_selVipPriceConfigsSQL         = "SELECT id,platform,product_name,product_id,suit_type,month,sub_type,original_price,selected,remark,status,operator,oper_id,ctime,mtime,superscript,start_build,end_build FROM vip_price_config_v2 ORDER BY id DESC"
	_selVipPriceConfigIDSQL        = "SELECT id,platform,product_name,product_id,suit_type,month,sub_type,original_price,selected,remark,status,operator,oper_id,ctime,mtime,superscript,start_build,end_build FROM vip_price_config_v2 WHERE id = ?"
	_selVipDPriceConfigsSQL        = "SELECT id,vpc_id,product_id,discount_price,stime,etime,remark,operator,oper_id,ctime,mtime,first_price FROM vip_price_discount_config_v2 WHERE vpc_id = ? ORDER BY stime ASC"
	_selVipDPriceConfigIDSQL       = "SELECT id,vpc_id,product_id,discount_price,stime,etime,remark,operator,oper_id,ctime,mtime,first_price FROM vip_price_discount_config_v2 WHERE id = ?"
	_selVipDPriceConfigUQTimeSQL   = "SELECT id,vpc_id,product_id,discount_price,stime,etime,remark,ctime,mtime FROM vip_price_discount_config_v2 WHERE vpc_id = ? AND ((etime <> '1970-01-01 08:00:00' AND (etime >= ? OR etime >= ?) AND stime <= ?) OR (etime = '1970-01-01 08:00:00' AND stime <= ?))"
	_selVipPriceDiscountConfigsSQL = "SELECT vpc_id,product_id,discount_price,stime,etime,remark,ctime,mtime FROM vip_price_discount_config_v2 WHERE stime <= ? AND ((etime > ? AND  etime <> '1970-01-01 08:00:00') OR (etime = '1970-01-01 08:00:00'))"
	_selVipMaxPriceDiscountSQL     = "SELECT MAX(discount_price) FROM vip_price_discount_config_v2 WHERE vpc_id = ?"
	_countVipPriceConfigByplatSQL  = "SELECT COUNT(*) FROM vip_price_config_v2 WHERE platform = ?"
)

// AddVipPriceConfig insert vip price config .
func (d *Dao) AddVipPriceConfig(c context.Context, v *model.VipPriceConfig) (err error) {
	if _, err = d.db.Exec(c, _inVipPriceConfigSQL, v.Plat, v.PdName, v.PdID, v.SuitType, v.Month, v.SubType, v.OPrice, v.Selected, v.Remark, v.Status, v.Operator, v.OpID, v.Superscript, v.StartBuild, v.EndBuild); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// AddVipDPriceConfig insert vip discount price config .
func (d *Dao) AddVipDPriceConfig(c context.Context, v *model.VipDPriceConfig) (err error) {
	if _, err = d.db.Exec(c, _inVipDPriceConfigSQL, v.ID, v.PdID, v.DPrice, v.STime, v.ETime, v.Remark, v.Operator, v.OpID, v.FirstPrice); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// UpVipPriceConfig update vip price config .
func (d *Dao) UpVipPriceConfig(c context.Context, v *model.VipPriceConfig) (err error) {
	if _, err = d.db.Exec(c, _upVipPriceConfigIDSQL, v.Plat, v.PdName, v.PdID, v.SuitType, v.Month, v.SubType, v.OPrice, v.Remark, v.Status, v.Operator, v.OpID, v.Selected, v.Superscript, v.StartBuild, v.EndBuild, v.ID); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// UpVipDPriceConfig update vip discount price config .
func (d *Dao) UpVipDPriceConfig(c context.Context, v *model.VipDPriceConfig) (err error) {
	if _, err = d.db.Exec(c, _upVipDPriceConfigSQL, v.PdID, v.DPrice, v.STime, v.ETime, v.Remark, v.Operator, v.OpID, v.FirstPrice, v.DisID); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// DelVipPriceConfig delete vip price config .
func (d *Dao) DelVipPriceConfig(c context.Context, arg *model.ArgVipPriceID) (err error) {
	if _, err = d.db.Exec(c, _delVipPriceConfigIDSQL, arg.ID); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// DelVipDPriceConfig delete vip discount price config .
func (d *Dao) DelVipDPriceConfig(c context.Context, arg *model.ArgVipDPriceID) (err error) {
	if _, err = d.db.Exec(c, _delVipDPriceConfigIDSQL, arg.DisID); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// VipPriceConfigUQCheck count vip price config unquie check.
func (d *Dao) VipPriceConfigUQCheck(c context.Context, arg *model.ArgAddOrUpVipPrice) (count int64, err error) {
	sqlPostfix := ""
	if arg.EndBuild > 0 {
		sqlPostfix += fmt.Sprintf(" AND (start_build<= %d )", arg.EndBuild)
	}
	if arg.StartBuild > 0 {
		sqlPostfix += fmt.Sprintf(" AND (end_build >= %d OR end_build = 0 )", arg.StartBuild)
	}
	if arg.ID > 0 {
		// for update
		sqlPostfix += fmt.Sprintf(" AND id != %d ", arg.ID)
	}
	if err = d.db.QueryRow(c,
		fmt.Sprintf(_selVipPriceConfigUQCheckSQL, sqlPostfix),
		arg.Plat,
		arg.Month,
		arg.SubType,
		arg.SuitType).Scan(&count); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// VipPriceConfigs vip price configs
func (d *Dao) VipPriceConfigs(c context.Context) (vpcs []*model.VipPriceConfig, err error) {
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, _selVipPriceConfigsSQL); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		vpc := new(model.VipPriceConfig)
		if err = rows.Scan(&vpc.ID, &vpc.Plat, &vpc.PdName, &vpc.PdID, &vpc.SuitType, &vpc.Month, &vpc.SubType,
			&vpc.OPrice, &vpc.Selected, &vpc.Remark, &vpc.Status, &vpc.Operator, &vpc.OpID, &vpc.CTime, &vpc.MTime,
			&vpc.Superscript, &vpc.StartBuild, &vpc.EndBuild); err != nil {
			if err == xsql.ErrNoRows {
				err = nil
				vpc = nil
			} else {
				err = errors.WithStack(err)
			}
			return
		}
		vpcs = append(vpcs, vpc)
	}
	return
}

// VipPriceConfigID vip price config by id
func (d *Dao) VipPriceConfigID(c context.Context, arg *model.ArgVipPriceID) (vpc *model.VipPriceConfig, err error) {
	row := d.db.QueryRow(c, _selVipPriceConfigIDSQL, arg.ID)
	vpc = new(model.VipPriceConfig)
	if err = row.Scan(&vpc.ID, &vpc.Plat, &vpc.PdName, &vpc.PdID, &vpc.SuitType, &vpc.Month, &vpc.SubType, &vpc.OPrice, &vpc.Selected,
		&vpc.Remark, &vpc.Status, &vpc.Operator, &vpc.OpID, &vpc.CTime, &vpc.MTime, &vpc.Superscript, &vpc.StartBuild, &vpc.EndBuild); err != nil {
		if err != sql.ErrNoRows {
			err = errors.WithStack(err)
			return
		}
		err = nil
		vpc = nil
	}
	return
}

// VipDPriceConfigs vip discount price configs
func (d *Dao) VipDPriceConfigs(c context.Context, arg *model.ArgVipPriceID) (res []*model.VipDPriceConfig, err error) {
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, _selVipDPriceConfigsSQL, arg.ID); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		v := new(model.VipDPriceConfig)
		if err = rows.Scan(&v.DisID, &v.ID, &v.PdID, &v.DPrice, &v.STime, &v.ETime, &v.Remark, &v.Operator, &v.OpID, &v.CTime, &v.MTime, &v.FirstPrice); err != nil {
			if err == xsql.ErrNoRows {
				err = nil
				res = nil
			} else {
				err = errors.WithStack(err)
			}
			return
		}
		res = append(res, v)
	}
	return
}

// VipDPriceConfigID vip discount price config by id
func (d *Dao) VipDPriceConfigID(c context.Context, arg *model.ArgVipDPriceID) (res *model.VipDPriceConfig, err error) {
	row := d.db.QueryRow(c, _selVipDPriceConfigIDSQL, arg.DisID)
	res = new(model.VipDPriceConfig)
	if err = row.Scan(&res.DisID, &res.ID, &res.PdID, &res.DPrice, &res.STime, &res.ETime, &res.Remark, &res.Operator, &res.OpID, &res.CTime, &res.MTime, &res.FirstPrice); err != nil {
		if err != xsql.ErrNoRows {
			err = errors.WithStack(err)
			return
		}
		err = nil
		res = nil
	}
	return
}

// VipDPriceConfigUQTime count vip discount price config unquie check time.
func (d *Dao) VipDPriceConfigUQTime(c context.Context, arg *model.ArgAddOrUpVipDPrice) (mvd map[int64]*model.VipDPriceConfig, err error) {
	var (
		rows  *xsql.Rows
		etime = arg.ETime
	)
	if etime == 0 {
		etime = xtime.Time(time.Now().AddDate(1000, 0, 0).Unix())
	}
	mvd = make(map[int64]*model.VipDPriceConfig)
	if rows, err = d.db.Query(c, _selVipDPriceConfigUQTimeSQL, arg.ID, arg.STime, etime, etime, arg.STime); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		vdc := new(model.VipDPriceConfig)
		if err = rows.Scan(&vdc.DisID, &vdc.ID, &vdc.PdID, &vdc.DPrice, &vdc.STime, &vdc.ETime, &vdc.Remark, &vdc.CTime, &vdc.MTime); err != nil {
			if err != xsql.ErrNoRows {
				err = errors.WithStack(err)
				return
			}
			mvd = nil
			err = nil
			return
		}
		mvd[vdc.DisID] = vdc
	}
	err = rows.Err()
	return
}

// VipPriceDiscountConfigs get vip price discount configs.
func (d *Dao) VipPriceDiscountConfigs(c context.Context) (mvd map[int64]*model.VipDPriceConfig, err error) {
	var (
		rows *xsql.Rows
		now  = time.Now()
	)
	mvd = make(map[int64]*model.VipDPriceConfig)
	if rows, err = d.db.Query(c, _selVipPriceDiscountConfigsSQL, now, now); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		vdc := new(model.VipDPriceConfig)
		if err = rows.Scan(&vdc.ID, &vdc.PdID, &vdc.DPrice, &vdc.STime, &vdc.ETime, &vdc.Remark, &vdc.CTime, &vdc.MTime); err != nil {
			if err != xsql.ErrNoRows {
				err = errors.WithStack(err)
				return
			}
			mvd = nil
			err = nil
			return
		}
		mvd[vdc.ID] = vdc
	}
	err = rows.Err()
	return
}

// VipMaxPriceDiscount max price discount.
func (d *Dao) VipMaxPriceDiscount(c context.Context, arg *model.ArgAddOrUpVipPrice) (res float64, err error) {
	var max sql.NullFloat64
	row := d.db.QueryRow(c, _selVipMaxPriceDiscountSQL, arg.ID)
	if err = row.Scan(&max); err != nil {
		return
	}
	res = max.Float64
	return
}

// CountVipPriceConfigByPlat count vip price config by platform id.
func (d *Dao) CountVipPriceConfigByPlat(c context.Context, plat int64) (count int64, err error) {
	row := d.db.QueryRow(c, _countVipPriceConfigByplatSQL, plat)
	if err = row.Scan(&count); err != nil {
		err = errors.WithStack(err)
	}
	return
}
