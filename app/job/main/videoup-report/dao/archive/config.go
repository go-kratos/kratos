package archive

import (
	"context"
	"database/sql"
	"encoding/json"

	"go-common/app/job/main/videoup-report/model/task"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_confSQL            = "SELECT value FROM archive_config WHERE state=0 AND name=?"
	_confForAuditType   = "wait_audit_arctype"
	_confForWeightValue = "weight_conf_values"
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

//WeightValueConf 权重数值配置
func (d *Dao) WeightValueConf(c context.Context) (wvconf *task.WeightValueConf, err error) {
	var value []byte

	if err = d.db.QueryRow(c, _confSQL, _confForWeightValue).Scan(&value); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan() error(%v)", err)
		}
		return
	}
	wvconf = new(task.WeightValueConf)
	if err = json.Unmarshal(value, wvconf); err != nil {
		log.Error("json.Unmarshal error(%v)", err)
		wvconf = nil
		return
	}
	wvconf.Nsum9 = wvconf.Nlv5 * 3
	wvconf.Nsum15 = (wvconf.Nlv1 * 2) + wvconf.Nsum9
	wvconf.Nsum27 = (wvconf.Nlv2 * 4) + wvconf.Nsum15
	wvconf.Nsum45 = (wvconf.Nlv3 * 6) + wvconf.Nsum27
	wvconf.Tsum2h = wvconf.Tlv1 * 40
	wvconf.Tsum1h = (wvconf.Tlv2 * 20) + wvconf.Tsum2h
	wvconf.MinWeight = -(wvconf.Tsum1h + wvconf.Tlv3*10)
	return
}
