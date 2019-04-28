package service

import (
	"go-common/app/admin/main/videoup/model/archive"
	"go-common/library/database/sql"
	"go-common/library/log"
)

// TxUpPorder .
func (s *Service) TxUpPorder(tx *sql.Tx, ap *archive.ArcParam) (err error) {
	//区分自首还是审核回查添加
	if _, err = s.arc.TxUpPorder(tx, ap.Aid, ap); err != nil {
		log.Error("s.arc.TxUpPorder(%d,%+v) error(%v)", ap.Aid, ap, err)
		return
	}
	log.Info("TxUpPorder aid(%d) update archive_porder", ap.Aid)
	return
}
