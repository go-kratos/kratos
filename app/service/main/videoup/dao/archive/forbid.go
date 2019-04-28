package archive

import (
	"go-common/app/service/main/videoup/model/archive"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_upFlowIDSQL = "INSERT INTO archive_forbid (aid,on_flow_id) VALUES (?,?) ON DUPLICATE KEY UPDATE on_flow_id=?"
	_inForbidSQL = `INSERT INTO archive_forbid (aid,rank_attr,recommend_attr,dynamic_attr,show_attr) VALUES (?,?,?,?,?) ON DUPLICATE KEY UPDATE
					rank_attr=?,recommend_attr=?,dynamic_attr=?,show_attr=?`
)

// TxUpForbid update archive addit.
func (d *Dao) TxUpForbid(tx *xsql.Tx, aid, flowID int64) (rows int64, err error) {
	res, err := tx.Exec(_upFlowIDSQL, aid, flowID, flowID)
	if err != nil {
		log.Error("d.upFlowID.Exec() error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxUpForbidAttr insert archive forbid.
func (d *Dao) TxUpForbidAttr(tx *xsql.Tx, af *archive.ForbidAttr) (rows int64, err error) {
	res, err := tx.Exec(_inForbidSQL, af.Aid, af.RankV, af.RecommendV, af.DynamicV, af.ShowV, af.RankV, af.RecommendV, af.DynamicV, af.ShowV)
	if err != nil {
		log.Error("d.inForbid.Exec error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}
