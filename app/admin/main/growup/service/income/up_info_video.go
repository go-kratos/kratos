package income

import (
	"go-common/library/database/sql"
	"go-common/library/log"
)

// TxUpdateUpInfoScore update up_info_video score
func (s *Service) TxUpdateUpInfoScore(tx *sql.Tx, mid int64, score int) (err error) {
	_, err = s.dao.TxUpdateUpInfoScore(tx, "up_info_video", score, mid)
	if err != nil {
		tx.Rollback()
		log.Error("s.dao.TxUpdateUpInfoScore error(%v)", err)
		return err
	}

	_, err = s.dao.TxUpdateUpInfoScore(tx, "up_info_column", score, mid)
	if err != nil {
		tx.Rollback()
		log.Error("s.dao.TxUpdateUpInfoScore error(%v)", err)
	}
	return
}
