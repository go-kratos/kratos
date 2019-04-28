package dao

import (
	"bytes"
	"context"
	"io"
	"strconv"
	"time"

	"go-common/app/service/main/passport/model"
	"go-common/library/log"

	"github.com/tsuna/gohbase/filter"
	"github.com/tsuna/gohbase/hrpc"
)

const (
	_tFace       = "account:user_face"
	_fFaceApply  = "c"
	_cMid        = "mid"
	_cApplyTime  = "at"
	_cModifyTime = "mt"
	_cNewFace    = "nf"
	_cOldFace    = "of"
	_cOperator   = "op"
	_cStatus     = "s"

	_maxIDLen = 10
)

var (
	_tFaceB       = []byte(_tFace)
	_cMidB        = []byte(_cMid)
	_fFaceApplyB  = []byte(_fFaceApply)
	_cApplyTimeB  = []byte(_cApplyTime)
	_cModifyTimeB = []byte(_cModifyTime)
	_cNewFaceB    = []byte(_cNewFace)
	_cOldFaceB    = []byte(_cOldFace)
	_cOperatorB   = []byte(_cOperator)
	_cStatusB     = []byte(_cStatus)
)

// FaceApplies get face applies from hbase.
func (d *Dao) FaceApplies(c context.Context, mid, from, to int64, status, operator string) (res []*model.FaceApply, err error) {
	midStr := strconv.FormatInt(mid, 10)
	if !checkIDLen(midStr) {
		log.Error("midInt64: %d, midStr: %s, len(midStr): %d, exceed max length %d", mid, midStr, len(midStr), _maxIDLen)
		return
	}

	ctx, cancel := context.WithTimeout(c, time.Duration(d.c.HBase.FaceApply.ReadTimeout))
	defer cancel()
	st := rowKeyFaceApplyMts(midStr, to, _uint32Max)
	if len(st) == 0 {
		return
	}
	ed := rowKeyFaceApplyMts(midStr, from, 0)
	if len(st) == 0 {
		return
	}

	fts := filter.NewList(filter.MustPassAll)
	if status != "" {
		statusFt := filter.NewSingleColumnValueFilter(_fFaceApplyB, _cStatusB, filter.Equal, filter.NewBinaryComparator(filter.NewByteArrayComparable([]byte(status))), true, true)
		fts.AddFilters(statusFt)
	}
	if operator != "" {
		operatorFt := filter.NewSingleColumnValueFilter(_fFaceApplyB, _cOperatorB, filter.Equal, filter.NewBinaryComparator(filter.NewByteArrayComparable([]byte(operator))), true, true)
		fts.AddFilters(operatorFt)
	}
	var options []func(hrpc.Call) error
	if len(fts.Filters) > 0 {
		options = append(options, hrpc.Filters(fts))
	}
	var scaner hrpc.Scanner
	scaner, err = d.hbase.ScanRange(ctx, _tFaceB, st, ed, options...)
	if err != nil {
		log.Error("hbase.ScanRange(%s, %s, %s) error(%v)", _tFaceB, st, ed, err)
		return
	}

	var rs []*hrpc.Result
	for {
		var r *hrpc.Result
		r, err = scaner.Next()
		if err != nil {
			if err == io.EOF {
				// set err nil
				err = nil
				break
			}
			return
		}
		rs = append(rs, r)
	}

	if len(rs) == 0 {
		return
	}
	res = make([]*model.FaceApply, 0, len(rs))
	for _, r := range rs {
		var u *model.FaceApply
		if u, err = scanFaceRecord(r.Cells); err != nil {
			return
		}
		if u != nil {
			res = append(res, u)
		}
	}
	return
}

// rowKeyFaceApplyMts get row key for face apply using this schema:
// mid string reverse with right fill 0 +
// (int64_max - mts) string cut last 10 digit +
// (unsigned_int32_max - id) string with left fill 0.
func rowKeyFaceApplyMts(midStr string, mts, id int64) []byte {
	rMid := reverseID(midStr, 10)
	if len(rMid) == 0 {
		return nil
	}
	b := bytes.Buffer{}
	b.WriteString(rMid)

	rMTS := diffTs(mts)
	b.WriteString(rMTS)

	rID := diffID(id)
	b.WriteString(rID)

	return b.Bytes()
}

func scanFaceRecord(cells []*hrpc.Cell) (res *model.FaceApply, err error) {
	if len(cells) == 0 {
		return
	}
	res = new(model.FaceApply)
	for _, cell := range cells {
		if bytes.Equal(cell.Family, _fFaceApplyB) {
			switch {
			case bytes.Equal(cell.Qualifier, _cMidB):
				if res.Mid, err = strconv.ParseInt(string(cell.Value), 10, 64); err != nil {
					log.Error("failed to parse mid from cell, strconv.ParseInt(%s, 10, 64) error(%v)", cell.Value, err)
					return
				}
			case bytes.Equal(cell.Qualifier, _cOldFaceB):
				res.OldFace = string(cell.Value)
			case bytes.Equal(cell.Qualifier, _cNewFaceB):
				res.NewFace = string(cell.Value)
			case bytes.Equal(cell.Qualifier, _cApplyTimeB):
				if res.ApplyTime, err = strconv.ParseInt(string(cell.Value), 10, 64); err != nil {
					log.Error("failed to parse apply_time from cell, strconv.ParseInt(%s, 10, 64) error(%v)", cell.Value, err)
					return
				}
			case bytes.Equal(cell.Qualifier, _cStatusB):
				res.Status = string(cell.Value)
			case bytes.Equal(cell.Qualifier, _cOperatorB):
				res.Operator = string(cell.Value)
			case bytes.Equal(cell.Qualifier, _cModifyTimeB):
				res.ModifyTime = string(cell.Value)
			}
		}
	}
	return
}
