package dao

import (
	"context"

	"go-common/app/admin/main/dm/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_countAdvsSQLAll      = "SELECT count(*) FROM dm_advancecomment WHERE dm_inid=? "
	_countAdvsSQLMode     = "SELECT count(*) FROM dm_advancecomment WHERE dm_inid=? AND mode=? "
	_countAdvsSQLType     = "SELECT count(*) FROM dm_advancecomment WHERE dm_inid=? AND type=? "
	_countAdvsSQLTypeMode = "SELECT count(*) FROM dm_advancecomment WHERE dm_inid=? AND type=? AND mode=? "
	_selAdvsSQLAll        = "SELECT buy_id,type,mode,mid,timestamp FROM dm_advancecomment WHERE dm_inid=? ORDER BY buy_id DESC limit ?,?"
	_selAdvsSQLMode       = "SELECT buy_id,type,mode,mid,timestamp FROM dm_advancecomment WHERE dm_inid=? AND mode=? ORDER BY buy_id DESC limit ?,?"
	_selAdvsSQLType       = "SELECT buy_id,type,mode,mid,timestamp FROM dm_advancecomment WHERE dm_inid=? AND type=? ORDER BY buy_id DESC limit ?,?"
	_selAdvsSQLTypeMode   = "SELECT buy_id,type,mode,mid,timestamp FROM dm_advancecomment WHERE dm_inid=? AND type=? AND mode=? ORDER BY buy_id DESC ?,?"
)

// Advances 获取高级弹幕申请
func (d *Dao) Advances(c context.Context, dmInid int64, typ, mode string, pn, ps int64) (res []*model.Advance, total int64, err error) {
	var rows *sql.Rows
	res = make([]*model.Advance, 0)
	if typ == model.AdvTypeAll {
		if mode == model.AdvModeAll {
			countRow := d.biliDM.QueryRow(c, _countAdvsSQLAll, dmInid)
			if err = countRow.Scan(&total); err != nil {
				log.Error("row.ScanCount error(%v)", err)
				return
			}
			rows, err = d.biliDM.Query(c, _selAdvsSQLAll, dmInid, (pn-1)*ps, ps)
			if err != nil {
				log.Error("d.dbDM.Query(%s,%d,%d,%d) error(%v)", _selAdvsSQLAll, dmInid, (pn-1)*ps, ps, err)
				return
			}
		} else {
			countRow := d.biliDM.QueryRow(c, _countAdvsSQLMode, dmInid, mode)
			if err = countRow.Scan(&total); err != nil {
				log.Error("row.ScanCount error(%v)", err)
				return
			}
			rows, err = d.biliDM.Query(c, _selAdvsSQLMode, dmInid, mode, (pn-1)*ps, ps)
			if err != nil {
				log.Error("d.dbDM.Query(%s,%d,%s,%d,%d) error(%v)", _selAdvsSQLMode, dmInid, mode, (pn-1)*ps, ps, err)
				return
			}
		}
	} else {
		if mode == model.AdvModeAll {
			countRow := d.biliDM.QueryRow(c, _countAdvsSQLType, dmInid, typ)
			if err = countRow.Scan(&total); err != nil {
				log.Error("row.ScanCount error(%v)", err)
				return
			}
			rows, err = d.biliDM.Query(c, _selAdvsSQLType, dmInid, typ, (pn-1)*ps, ps)
			if err != nil {
				log.Error("d.dbDM.Query(%s,%d,%s,%d,%d) error(%v)", _selAdvsSQLType, dmInid, typ, (pn-1)*ps, ps, err)
				return
			}
		} else {
			countRow := d.biliDM.QueryRow(c, _countAdvsSQLTypeMode, dmInid, typ)
			if err = countRow.Scan(&total); err != nil {
				log.Error("row.ScanCount error(%v)", err)
				return
			}
			rows, err = d.biliDM.Query(c, _selAdvsSQLTypeMode, dmInid, typ, mode, (pn-1)*ps, ps)
			if err != nil {
				log.Error("d.dbDM.Query(%s,%d,%s,%s,%d,%d) error(%v)", _selAdvsSQLTypeMode, dmInid, typ, mode, (pn-1)*ps, ps, err)
				return
			}
		}
	}
	defer rows.Close()
	for rows.Next() {
		adv := &model.Advance{}
		if err = rows.Scan(&adv.ID, &adv.Type, &adv.Mode, &adv.Mid, &adv.Timestamp); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		res = append(res, adv)
	}
	return
}
