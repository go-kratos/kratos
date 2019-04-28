package archive

import (
	"context"
	"database/sql"

	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_confSQL          = "SELECT value FROM archive_config WHERE state=0 AND name=?"
	_confForAuditType = "wait_audit_arctype"
)

// AuditTypesConf get audit conf
func (d *Dao) AuditTypesConf(c context.Context) (atps map[int16]struct{}, err error) {
	row := d.db.QueryRow(c, _confSQL, _confForAuditType)
	var (
		value   string
		typeIDs []int64
	)
	if err = row.Scan(&value); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
		return
	}
	typeIDs, err = xstr.SplitInts(value)
	if err != nil {
		log.Error("archive_config value(%s) xstr.SplitInts error(%v)", value, err)
		return
	}
	atps = map[int16]struct{}{}
	for _, typeid := range typeIDs {
		atps[int16(typeid)] = struct{}{}
	}
	return
}
