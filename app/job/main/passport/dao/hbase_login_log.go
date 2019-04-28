package dao

import (
	"bytes"
	"context"
	"encoding/binary"
	"strconv"
	"time"

	"go-common/app/job/main/passport/model"
	"go-common/library/log"
)

const (
	_tLoginLog        = "ugc:AsoLoginLog"
	_fLoginLogInfo    = "f"
	_cLoginLogMid     = "mid"
	_cLoginLogTs      = "ts"
	_cLoginLogLoginIP = "ip"
	_cLoginLogType    = "t"
	_cLoginLogServer  = "s"

	_int64Max = 0x7fffffffffffffff
)

// AddLoginLogHBase add login log.
func (d *Dao) AddLoginLogHBase(c context.Context, loginLog *model.LoginLog) (err error) {
	fvs := make(map[string][]byte)
	fvs[_cLoginLogMid] = []byte(strconv.FormatInt(loginLog.Mid, 10))
	fvs[_cLoginLogTs] = []byte(strconv.FormatInt(loginLog.Timestamp, 10))
	fvs[_cLoginLogLoginIP] = []byte(strconv.FormatInt(loginLog.LoginIP, 10))
	fvs[_cLoginLogType] = []byte(strconv.FormatInt(loginLog.Type, 10))
	fvs[_cLoginLogServer] = []byte(loginLog.Server)
	values := map[string]map[string][]byte{_fLoginLogInfo: fvs}
	key := rowKeyLoginLog(loginLog.Mid, loginLog.Timestamp)
	ctx, cancel := context.WithTimeout(c, time.Duration(d.c.HBase.LoginLog.WriteTimeout))
	defer cancel()
	if _, err = d.loginLogHBase.PutStr(ctx, _tLoginLog, string(key), values); err != nil {
		log.Error("dao.hbase.Put(%v) error(%v)", err)
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
