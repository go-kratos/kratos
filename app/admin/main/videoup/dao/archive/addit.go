package archive

import (
	"context"

	"fmt"
	"go-common/app/admin/main/videoup/model/archive"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_inAddRdrSQL         = "INSERT INTO archive_addit (aid,redirect_url) VALUES (?,?) ON DUPLICATE KEY UPDATE redirect_url=?"
	_upAdditSQL          = "UPDATE archive_addit SET mission_id=?,source=?,description=?,dynamic=? WHERE aid=?"
	_AdditSQL            = "SELECT aid,mission_id,from_ip,up_from,recheck_reason,redirect_url,source,order_id,desc_format_id,dynamic,inner_attr FROM archive_addit WHERE aid=?"
	_additBatch          = "SELECT aid,mission_id,from_ip,up_from,recheck_reason,redirect_url,source,order_id,desc_format_id,dynamic,inner_attr FROM archive_addit WHERE aid IN (%s)"
	_inAdditInnerAttrSQL = "INSERT INTO archive_addit (aid, inner_attr) VALUES (?,?) ON DUPLICATE KEY UPDATE inner_attr=?"
)

// TxUpAdditRedirect update archive redirect url.
func (d *Dao) TxUpAdditRedirect(tx *xsql.Tx, aid int64, redirectURL string) (rows int64, err error) {
	res, err := tx.Exec(_inAddRdrSQL, aid, redirectURL, redirectURL)
	if err != nil {
		log.Error("d._inAdditRedirect.Exec() error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxUpAddit update archive_addit mission_id && source by aid.
func (d *Dao) TxUpAddit(tx *xsql.Tx, aid, missionID int64, source, description, dynamic string) (rows int64, err error) {
	res, err := tx.Exec(_upAdditSQL, missionID, source, description, dynamic, aid)
	if err != nil {
		log.Error("d.TxUpAddit.tx.Exec error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// Addit get a archive addit by avid.
func (d *Dao) Addit(c context.Context, aid int64) (ad *archive.Addit, err error) {
	row := d.rddb.QueryRow(c, _AdditSQL, aid)
	ad = &archive.Addit{}
	if err = row.Scan(&ad.Aid, &ad.MissionID, &ad.FromIP, &ad.UpFrom, &ad.RecheckReason, &ad.RedirectURL, &ad.Source, &ad.OrderID, &ad.DescFormatID, &ad.Dynamic, &ad.InnerAttr); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

//TxUpInnerAttr update archive inner_attr
func (d *Dao) TxUpInnerAttr(tx *xsql.Tx, aid int64, attr int64) (id int64, err error) {
	res, err := tx.Exec(_inAdditInnerAttrSQL, aid, attr, attr)
	if err != nil {
		log.Error("TxUpInnerAttr tx.Exec error(%v) aid(%d) attr(%d)", err, aid, attr)
		return 0, err
	}

	return res.LastInsertId()
}

// ArcStateMap get archive id and state map
func (d *Dao) AdditBatch(c context.Context, aids []int64) (sMap map[int64]*archive.Addit, err error) {
	sMap = make(map[int64]*archive.Addit)
	if len(aids) == 0 {
		return
	}
	rows, err := d.rddb.Query(c, fmt.Sprintf(_additBatch, xstr.JoinInts(aids)))
	if err != nil {
		log.Error("db.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		ad := &archive.Addit{}
		if err = rows.Scan(&ad.Aid, &ad.MissionID, &ad.FromIP, &ad.UpFrom, &ad.RecheckReason, &ad.RedirectURL, &ad.Source, &ad.OrderID, &ad.DescFormatID, &ad.Dynamic, &ad.InnerAttr); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		sMap[ad.Aid] = ad
	}
	return
}
