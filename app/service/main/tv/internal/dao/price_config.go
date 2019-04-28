package dao

import (
	"context"
	"go-common/app/service/main/tv/internal/model"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	_getPriceConfigsByStatusAndPid       = "SELECT `id`, `pid`, `platform`, `product_name`, `product_id`, `suit_type`, `month`, `sub_type`, `price`, `selected`, `remark`, `status`, `superscript`, `operator`, `oper_id`, `stime`, `etime`, `ctime`, `mtime` FROM `tv_price_config` WHERE `status`=? AND `pid`=? ORDER BY `ctime` ASC "
	_countPriceConfigByStatusAndPid      = "SELECT count(*) FROM `tv_price_config` WHERE `status`=? AND `pid`=?"
	_getSaledPriceConfigsByStatusAndPid  = "SELECT `id`, `pid`, `platform`, `product_name`, `product_id`, `suit_type`, `month`, `sub_type`, `price`, `selected`, `remark`, `status`, `superscript`, `operator`, `oper_id`, `stime`, `etime`, `ctime`, `mtime` FROM `tv_price_config` WHERE `status`=? AND `pid`!=0 AND `stime`<=now() AND `etime`>now() ORDER BY `ctime` ASC "
	_countSaledPriceConfigByStatusAndPid = "SELECT count(*) FROM `tv_price_config` WHERE  `status`=? AND `pid`!=0 AND `stime`<=now() AND `etime`>now()"
)

// PriceConfigsByStatus quires rows from tv_price_config.
func (d *Dao) PriceConfigsByStatus(c context.Context, status int8) (res []*model.PriceConfig, total int, err error) {
	res = make([]*model.PriceConfig, 0)
	pid := 0
	totalRow := d.db.QueryRow(c, _countPriceConfigByStatusAndPid, status, pid)
	if err = totalRow.Scan(&total); err != nil {
		log.Error("row.ScanCount error(%v)", err)
		err = errors.WithStack(err)
		return
	}
	rows, err := d.db.Query(c, _getPriceConfigsByStatusAndPid, status, pid)
	if err != nil {
		log.Error("db.Query(%s) error(%v)", _getPriceConfigsByStatusAndPid, err)
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		pc := &model.PriceConfig{}
		if err = rows.Scan(&pc.ID, &pc.Pid, &pc.Platform, &pc.ProductName, &pc.ProductId, &pc.SuitType, &pc.Month, &pc.SubType, &pc.Price, &pc.Selected, &pc.Remark, &pc.Status, &pc.Superscript, &pc.Operator, &pc.OperId, &pc.Stime, &pc.Etime, &pc.Ctime, &pc.Mtime); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			err = errors.WithStack(err)
			return
		}
		res = append(res, pc)
	}
	return
}

// SaledPriceConfigsByStatus quires rows from tv_price_config.
func (d *Dao) SaledPriceConfigsByStatus(c context.Context, status int8) (res []*model.PriceConfig, total int, err error) {
	res = make([]*model.PriceConfig, 0)
	totalRow := d.db.QueryRow(c, _countSaledPriceConfigByStatusAndPid, status)
	if err = totalRow.Scan(&total); err != nil {
		log.Error("row.ScanCount error(%v)", err)
		err = errors.WithStack(err)
		return
	}
	rows, err := d.db.Query(c, _getSaledPriceConfigsByStatusAndPid, status)
	if err != nil {
		log.Error("db.Query(%s) error(%v)", _getPriceConfigsByStatusAndPid, err)
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		pc := &model.PriceConfig{}
		if err = rows.Scan(&pc.ID, &pc.Pid, &pc.Platform, &pc.ProductName, &pc.ProductId, &pc.SuitType, &pc.Month, &pc.SubType, &pc.Price, &pc.Selected, &pc.Remark, &pc.Status, &pc.Superscript, &pc.Operator, &pc.OperId, &pc.Stime, &pc.Etime, &pc.Ctime, &pc.Mtime); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			err = errors.WithStack(err)
			return
		}
		res = append(res, pc)
	}
	return
}
