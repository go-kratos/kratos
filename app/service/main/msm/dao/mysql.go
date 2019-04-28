package dao

import (
	"container/list"
	"context"

	"go-common/app/service/main/msm/model"
	"go-common/library/log"
	xtime "go-common/library/time"
)

const (
	_codesSQL      = "SELECT code, message, mtime FROM codes"
	_diffCodesSQL  = "SELECT code, message, mtime FROM codes WHERE mtime > ? ORDER BY mtime LIMIT 100"
	_codesLangsSQL = "select a.code,a.message,a.mtime,IFNULL(b.locale,''),IFNULL(b.msg,''),IFNULL(b.mtime,'') as bmtime from codes as a left join code_msg as b on a.id=b.code_id"
)

// Codes get all codes.
func (d *Dao) Codes(c context.Context) (codes map[int]string, lcode *model.Code, err error) {
	var (
		code  int
		msg   string
		tmp   int64
		mtime xtime.Time
	)
	rows, err := d.db.Query(c, _codesSQL)
	if err != nil {
		log.Error("d.db.Query(%v) error(%v)", _codesSQL, err)
		return
	}
	defer rows.Close()
	lcode = &model.Code{}
	codes = make(map[int]string)
	for rows.Next() {
		if err = rows.Scan(&code, &msg, &mtime); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		codes[code] = msg
		if int64(mtime) > tmp {
			lcode.Code = code
			lcode.Ver = int64(mtime)
			lcode.Msg = msg
			tmp = int64(mtime)
		}
	}
	return
}

// Diff get change codes.
func (d *Dao) Diff(c context.Context, ver int64) (vers *list.List, err error) {
	var (
		code  int
		msg   string
		mtime xtime.Time
	)
	rows, err := d.db.Query(c, _diffCodesSQL, xtime.Time(ver))
	if err != nil {
		log.Error("d.db.Query(%v) error(%v)", _diffCodesSQL, err)
		return
	}
	defer rows.Close()
	vers = list.New()
	for rows.Next() {
		if err = rows.Scan(&code, &msg, &mtime); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		vers.PushBack(&model.Code{Ver: int64(mtime), Code: code, Msg: msg})
	}
	return
}

// CodesLang get all codes.
func (d *Dao) CodesLang(c context.Context) (codes map[int]map[string]string, lcode *model.CodeLangs, err error) {
	var (
		code    int
		tmp     int64
		mtime   xtime.Time
		message string
		bl      string
		bmsg    string
		bmtime  xtime.Time
	)
	rows, err := d.db.Query(c, _codesLangsSQL)
	if err != nil {
		log.Error("d.db.Query(%v) error(%v)", _codesLangsSQL, err)
		return
	}
	defer rows.Close()
	lcode = &model.CodeLangs{}
	codes = make(map[int]map[string]string)
	for rows.Next() {
		t := make(map[string]string)
		bl = ""
		if err = rows.Scan(&code, &message, &mtime, &bl, &bmsg, &bmtime); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		t["default"] = message
		if len(bl) > 0 {
			t[bl] = bmsg
		}
		codes[code] = t
		if int64(mtime) > tmp {
			lcode.Code = code
			lcode.Ver = int64(mtime)
			lcode.Msg = t
			tmp = int64(mtime)
		}
	}
	err = rows.Err()
	return
}
