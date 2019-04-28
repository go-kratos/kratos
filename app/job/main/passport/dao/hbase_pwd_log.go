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
	_tPwdLog        = "ugc:PwdLog"
	_fPwdLog        = "pwdlog"
	_cPwdLogMid     = "mid"
	_cPwdLogOldPwd  = "old_pwd"
	_cPwdLogOldSalt = "old_salt"
	_cPwdLogNewPwd  = "new_pwd"
	_cPwdLogNewSalt = "new_salt"
	_cPwdLogIP      = "ip"
	_cPwdLogTs      = "ts"
)

// AddPwdLogHBase add pwd log.
func (d *Dao) AddPwdLogHBase(c context.Context, pwdLog *model.PwdLog) (err error) {
	fvs := make(map[string][]byte)
	fvs[_cPwdLogMid] = []byte(strconv.FormatInt(pwdLog.Mid, 10))
	fvs[_cPwdLogOldPwd] = []byte(pwdLog.OldPwd)
	fvs[_cPwdLogOldSalt] = []byte(pwdLog.OldSalt)
	fvs[_cPwdLogNewPwd] = []byte(pwdLog.NewPwd)
	fvs[_cPwdLogNewSalt] = []byte(pwdLog.NewSalt)
	fvs[_cPwdLogTs] = []byte(strconv.FormatInt(pwdLog.Timestamp, 10))
	fvs[_cPwdLogIP] = []byte(strconv.FormatInt(pwdLog.IP, 10))

	values := map[string]map[string][]byte{_fPwdLog: fvs}
	key := rowKeyPwdLog(pwdLog.Mid, pwdLog.Timestamp)
	ctx, cancel := context.WithTimeout(c, time.Duration(d.c.HBase.PwdLog.WriteTimeout))
	defer cancel()
	if _, err = d.pwdLogHBase.PutStr(ctx, _tPwdLog, string(key), values); err != nil {
		log.Error("failed to put pwd log to hbase, dao.hbase.Put(%+v) error(%v)", pwdLog, err)
	}
	log.Info("Add pwdLog to HBase, (%+v)", pwdLog)
	return
}

// rowKeyPwdLog generate row key of pwd log.
func rowKeyPwdLog(mid, ts int64) (res []byte) {
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
