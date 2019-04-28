package archive

import (
	"context"
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"

	"go-common/app/job/main/videoup/model/archive"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_confSQL = "SELECT value FROM archive_config WHERE state=0 AND name=?"
)

//RoundEndConf round_delay_time 正常状态的配置
func (d *Dao) RoundEndConf(c context.Context) (days int64, err error) {
	row := d.db.QueryRow(c, _confSQL, archive.ConfForRoundEnd)
	var val string
	if err = row.Scan(&val); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
		return
	}
	if days, err = strconv.ParseInt(val, 10, 64); err != nil {
		log.Error("srconv.ParseInt(%s) error(%v)", val, err)
	}
	return
}

//FansConf round_limit_fans正常状态的配置
func (d *Dao) FansConf(c context.Context) (fans int64, err error) {
	row := d.db.QueryRow(c, _confSQL, archive.ConfForClick)
	var val string
	if err = row.Scan(&val); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
		return
	}
	if fans, err = strconv.ParseInt(val, 10, 64); err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", val, err)
	}
	return
}

//RoundTypeConf round_limit_tids 正常状态的配置
func (d *Dao) RoundTypeConf(c context.Context) (roundTypes map[int16]struct{}, err error) {
	roundTypes = map[int16]struct{}{}
	row := d.db.QueryRow(c, _confSQL, archive.ConfForRoundType)
	var (
		val  string
		tids []string
		tid  int64
	)
	if err = row.Scan(&val); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
		return
	}
	tids = strings.Split(val, ",")
	for _, tidStr := range tids {
		if tid, err = strconv.ParseInt(tidStr, 10, 64); err != nil {
			log.Error("strconv.ParseInt(%d) error(%v)", tid, err)
			return
		}
		roundTypes[int16(tid)] = struct{}{}
	}
	return
}

//AuditTypesConf wait_audit_arctype 状态正常的配置
func (d *Dao) AuditTypesConf(c context.Context) (atps map[int16]struct{}, err error) {
	row := d.db.QueryRow(c, _confSQL, archive.ConfForAuditType)
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

//ThresholdConf ThresholdConf is second types opposite first types.
func (d *Dao) ThresholdConf(c context.Context) (tpThr map[int16]int, err error) {
	row := d.db.QueryRow(c, _confSQL, archive.ConfForThreshold)
	var value string
	if err = row.Scan(&value); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan() error(%v)", err)
		}
		return
	}
	if err = json.Unmarshal([]byte(value), &tpThr); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", value, err)
		return
	}
	return
}
