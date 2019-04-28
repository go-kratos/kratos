package dao

import (
	"bytes"
	"context"
	"encoding/binary"
	"io"
	"time"

	"go-common/app/service/main/passport/model"
	"go-common/library/log"

	"github.com/tsuna/gohbase/hrpc"
)

const (
	_tPwdLog        = "ugc:PwdLog"
	_fPwdLog        = "pwdlog"
	_cPwdLogOldPwd  = "old_pwd"
	_cPwdLogOldSalt = "old_salt"
)

var (
	_tPwdLogB        = []byte(_tPwdLog)
	_fPwdLogB        = []byte(_fPwdLog)
	_cPwdLogOldPwdB  = []byte(_cPwdLogOldPwd)
	_cPwdLogOldSaltB = []byte(_cPwdLogOldSalt)
)

// HistoryPwds get history pwd
func (d *Dao) HistoryPwds(c context.Context, mid int64) (res []*model.HistoryPwd, err error) {
	ctx, cancel := context.WithTimeout(c, time.Duration(d.c.HBase.LoginLog.ReadTimeout))
	defer cancel()
	st := RowKeyPwdLog(mid, _int64Max)
	ed := RowKeyPwdLog(mid, 0)
	var scaner hrpc.Scanner
	scaner, err = d.pwdLogHBase.ScanRange(ctx, _tPwdLogB, st, ed)
	if err != nil {
		log.Error("hbase.ScanRange(%s, %s, %s) error(%v)", _tPwdLogB, st, ed, err)
		return
	}

	res = make([]*model.HistoryPwd, 0)
	for {
		var pwd *model.HistoryPwd
		var r *hrpc.Result
		r, err = scaner.Next()
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			}
			return
		}
		if pwd, err = scanPwdLog(r.Cells); err != nil {
			scaner.Close()
			return
		}
		if pwd != nil {
			res = append(res, pwd)
		}
	}
	if err := scaner.Close(); err != nil {
		log.Error("hbase.Scanner.Close error(%v)", err)
	}
	return
}

// RowKeyPwdLog generate row key of pwd log.
func RowKeyPwdLog(mid, ts int64) (res []byte) {
	buf := bytes.Buffer{}
	b := make([]byte, 8)

	// reverse mid bytes
	binary.BigEndian.PutUint64(b, uint64(mid))
	reverse(b)
	buf.Write(b)

	// (int64_max - ts) bytes
	binary.BigEndian.PutUint64(b, uint64(_int64Max-ts))
	buf.Write(b)

	res = buf.Bytes()
	return
}

func scanPwdLog(cells []*hrpc.Cell) (res *model.HistoryPwd, err error) {
	if len(cells) == 0 {
		return
	}
	res = new(model.HistoryPwd)
	for _, cell := range cells {
		if bytes.Equal(cell.Family, _fPwdLogB) {
			switch {
			case bytes.Equal(cell.Qualifier, _cPwdLogOldPwdB):
				res.OldPwd = string(cell.Value)
			case bytes.Equal(cell.Qualifier, _cPwdLogOldSaltB):
				res.OldSalt = string(cell.Value)
			}
		}
	}
	return
}
