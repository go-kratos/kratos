package dao

import (
	"bytes"
	"context"
	"encoding/binary"
	"io"
	"strconv"
	"time"

	"go-common/app/service/main/passport/model"
	"go-common/library/log"

	"github.com/tsuna/gohbase/hrpc"
)

const (
	_tLoginLog        = "ugc:AsoLoginLog"
	_fLoginLogInfo    = "f"
	_cLoginLogMid     = "mid"
	_cLoginLogTs      = "ts"
	_cLoginLogLoginIP = "ip"
	_cLoginLogType    = "t"
	_cLoginLogServer  = "s"
)

var (
	_tLoginLogB        = []byte(_tLoginLog)
	_fLoginLogInfoB    = []byte(_fLoginLogInfo)
	_cLoginLogMidB     = []byte(_cLoginLogMid)
	_cLoginLogTsB      = []byte(_cLoginLogTs)
	_cLoginLogLoginIPB = []byte(_cLoginLogLoginIP)
	_cLoginLogTypeB    = []byte(_cLoginLogType)
	_cLoginLogServerB  = []byte(_cLoginLogServer)
)

// LoginLogs get last limit login logs.
func (d *Dao) LoginLogs(c context.Context, mid int64, limit int) (res []*model.LoginLog, err error) {
	ctx, cancel := context.WithTimeout(c, time.Duration(d.c.HBase.LoginLog.ReadTimeout))
	defer cancel()
	st := rowKeyLoginLog(mid, _int64Max)
	ed := rowKeyLoginLog(mid, 0)
	var scaner hrpc.Scanner
	scaner, err = d.loginLogHBase.ScanRange(ctx, _tLoginLogB, st, ed)
	if err != nil {
		log.Error("hbase.ScanRange(%s, %s, %s) error(%v)", _tLoginLogB, st, ed, err)
		return
	}

	res = make([]*model.LoginLog, 0)
	for ; limit > 0; limit-- {
		var u *model.LoginLog
		var r *hrpc.Result
		r, err = scaner.Next()
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			}
			return
		}
		if u, err = scanLoginLog(r.Cells); err != nil {
			scaner.Close()
			return
		}
		if u != nil {
			res = append(res, u)
		}
	}
	if err := scaner.Close(); err != nil {
		log.Error("hbase.Scanner.Close error(%v)", err)
	}
	return
}

// rowKeyLoginLog generate row key of login log.
func rowKeyLoginLog(mid, ts int64) (res []byte) {
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

func reverse(b []byte) {
	l := len(b)
	for i := 0; i < l/2; i++ {
		t := b[i]
		b[i] = b[l-1-i]
		b[l-1-i] = t
	}
}

func scanLoginLog(cells []*hrpc.Cell) (res *model.LoginLog, err error) {
	if len(cells) == 0 {
		return
	}
	res = new(model.LoginLog)
	for _, cell := range cells {
		if bytes.Equal(cell.Family, _fLoginLogInfoB) {
			switch {
			case bytes.Equal(cell.Qualifier, _cLoginLogMidB):
				if res.Mid, err = strconv.ParseInt(string(cell.Value), 10, 64); err != nil {
					log.Error("failed to parse mid from cell, strconv.ParseInt(%s, 10, 64) error(%v)", cell.Value, err)
					return
				}
			case bytes.Equal(cell.Qualifier, _cLoginLogTsB):
				if res.Timestamp, err = strconv.ParseInt(string(cell.Value), 10, 64); err != nil {
					log.Error("failed to parse timestamp from cell, strconv.ParseInt(%s, 10, 64) error(%v)", cell.Value, err)
					return
				}
			case bytes.Equal(cell.Qualifier, _cLoginLogLoginIPB):
				if res.LoginIP, err = strconv.ParseInt(string(cell.Value), 10, 64); err != nil {
					log.Error("failed to parse loginip from cell, strconv.ParseInt(%s, 10, 64) error(%v)", cell.Value, err)
					return
				}
			case bytes.Equal(cell.Qualifier, _cLoginLogTypeB):
				if res.Type, err = strconv.ParseInt(string(cell.Value), 10, 64); err != nil {
					log.Error("failed to parse type from cell, strconv.ParseInt(%s, 10, 64) error(%v)", cell.Value, err)
					return
				}
			case bytes.Equal(cell.Qualifier, _cLoginLogServerB):
				res.Server = string(cell.Value)
			}
		}
	}
	return
}
