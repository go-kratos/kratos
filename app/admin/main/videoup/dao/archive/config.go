package archive

import (
	"context"
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"

	"go-common/app/admin/main/videoup/model/archive"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_confSQL   = "SELECT value FROM archive_config WHERE state=0 AND name=?"
	_upconfSQL = "UPDATE archive_config SET value=?,remark=? WHERE name=?"
	_inconfSQL = "INSERT archive_config(value,remark,name,state) VALUE (?,?,?,0)"
)

// FansConf is fan round check types config.
func (d *Dao) FansConf(c context.Context) (fans int64, err error) {
	row := d.rddb.QueryRow(c, _confSQL, archive.ConfForClick)
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

// RoundTypeConf is typeid round check types config.
func (d *Dao) RoundTypeConf(c context.Context) (roundTypes map[int16]struct{}, err error) {
	roundTypes = map[int16]struct{}{}
	row := d.rddb.QueryRow(c, _confSQL, archive.ConfForRoundType)
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
			log.Error("strconv.ParseInt(%s) error(%v)", tid, err)
			return
		}
		roundTypes[int16(tid)] = struct{}{}
	}
	return
}

// ThresholdConf is second types opposite first types.
func (d *Dao) ThresholdConf(c context.Context) (tpThr map[int16]int, err error) {
	row := d.rddb.QueryRow(c, _confSQL, archive.ConfForThreshold)
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

// AuditTypesConf is audit types.
func (d *Dao) AuditTypesConf(c context.Context) (atps map[int16]struct{}, err error) {
	row := d.rddb.QueryRow(c, _confSQL, archive.ConfForWaitAudit)
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

// WeightVC 获取权重分值
func (d *Dao) WeightVC(c context.Context) (wvc *archive.WeightVC, err error) {
	var value []byte
	row := d.rddb.QueryRow(c, _confSQL, archive.ConfForWeightVC)
	if err = row.Scan(&value); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
		return
	}
	wvc = new(archive.WeightVC)
	if err = json.Unmarshal(value, wvc); err != nil {
		log.Error("json.Unmarshal error(%v)", err)
		wvc = nil
	}
	return
}

// SetWeightVC 设置权重分值
func (d *Dao) SetWeightVC(c context.Context, wvc *archive.WeightVC, desc string) (rows int64, err error) {
	var (
		valueb []byte
		res    sql.Result
	)
	if valueb, err = json.Marshal(wvc); err != nil {
		log.Error("json.Marshal(%+v) error(%v)", wvc, err)
		return
	}

	if res, err = d.db.Exec(c, _upconfSQL, string(valueb), desc, archive.ConfForWeightVC); err != nil {
		log.Error("d.db.Exec(%s, %s, %s, %s) error(%v)", _upconfSQL, string(valueb), desc, archive.ConfForWeightVC, err)
		return
	}
	return res.RowsAffected()
}

// InWeightVC 插入
func (d *Dao) InWeightVC(c context.Context, wvc *archive.WeightVC, desc string) (rows int64, err error) {
	var (
		valueb []byte
		res    sql.Result
	)
	if valueb, err = json.Marshal(wvc); err != nil {
		log.Error("json.Marshal(%+v) error(%v)", wvc, err)
		return
	}

	if res, err = d.db.Exec(c, _inconfSQL, string(valueb), desc, archive.ConfForWeightVC); err != nil {
		log.Error("d.db.Exec(%s, %s, %s, %s) error(%v)", _inconfSQL, string(valueb), desc, archive.ConfForWeightVC, err)
		return
	}
	return res.LastInsertId()
}
