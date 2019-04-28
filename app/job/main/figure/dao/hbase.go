package dao

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"time"

	"go-common/app/job/main/figure/conf"
	"go-common/app/job/main/figure/model"
	"go-common/library/log"

	"github.com/pkg/errors"
)

func (d *Dao) rowKey(mid int64) (res string) {
	res = fmt.Sprintf("%d", mid)
	return
}

func (d *Dao) rowVerKey(mid int64, now time.Time) (res string) {
	res = fmt.Sprintf("%d_%d", mid, d.Version(now))
	return
}

// PutSpyScore add spy score info.
func (d *Dao) PutSpyScore(c context.Context, mid int64, score int8) (err error) {
	var (
		key         = d.rowKey(mid)
		scoreB      = make([]byte, 8)
		ctx, cancel = context.WithTimeout(c, time.Duration(conf.Conf.HBase.WriteTimeout))
	)
	defer cancel()
	log.Info("Put spy score key [%s] score [%d]", key, score)
	binary.BigEndian.PutUint64(scoreB, uint64(score))
	values := map[string]map[string][]byte{model.USFamilyUser: map[string][]byte{model.USColumnSpyScore: scoreB}}
	if _, err = d.hbase.PutStr(ctx, model.UserInfoTable, key, values); err != nil {
		log.Error("hbase.Put error(%v)", err)
	}
	return
}

// PutReplyAct add spy score info.
func (d *Dao) PutReplyAct(c context.Context, mid int64, column string, incr int64) (err error) {
	var (
		key         = d.rowVerKey(mid, time.Now())
		ctx, cancel = context.WithTimeout(c, time.Duration(conf.Conf.HBase.WriteTimeout))
	)
	defer cancel()
	log.Info("Put reply act key [%s] c [%s] incr [%d]", key, column, incr)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, incr)
	values := map[string]map[string][]byte{model.ACFamilyUser: map[string][]byte{column: bytesBuffer.Bytes()}}
	if _, err = d.hbase.Increment(ctx, model.ActionCounterTable, key, values); err != nil {
		err = errors.Wrapf(err, "msg(%d,%s,%d), hbase.Increment(key: %s values: %v)", mid, column, incr, key, values)
	}
	return
}

// PutCoinUnusual coin unusual.
func (d *Dao) PutCoinUnusual(c context.Context, mid int64, column string) (err error) {
	var (
		key         = d.rowVerKey(mid, time.Now())
		incrBytes   = make([]byte, 8)
		ctx, cancel = context.WithTimeout(c, time.Duration(conf.Conf.HBase.WriteTimeout))
	)
	defer cancel()
	log.Info("Put coin unusual key [%s] c [%s]", key, column)
	binary.BigEndian.PutUint64(incrBytes, uint64(1))
	values := map[string]map[string][]byte{model.ACFamilyUser: map[string][]byte{column: incrBytes}}
	if _, err = d.hbase.Increment(ctx, model.ActionCounterTable, key, values); err != nil {
		err = errors.Wrapf(err, "msg(%d,%s), hbase.Increment(key: %s values: %v)", mid, column, key, values)
	}
	return
}

// PutCoinCount coin count.
func (d *Dao) PutCoinCount(c context.Context, mid int64) (err error) {
	var (
		key         = d.rowVerKey(mid, time.Now())
		incrBytes   = make([]byte, 8)
		ctx, cancel = context.WithTimeout(c, time.Duration(conf.Conf.HBase.WriteTimeout))
	)
	defer cancel()
	log.Info("Put coin count key [%s]", key)
	binary.BigEndian.PutUint64(incrBytes, uint64(1))
	values := map[string]map[string][]byte{model.ACFamilyUser: map[string][]byte{model.ACColumnCoins: incrBytes}}
	if _, err = d.hbase.Increment(ctx, model.ActionCounterTable, key, values); err != nil {
		err = errors.Wrapf(err, "msg(%d), hbase.Increment(key: %s values: %v)", mid, key, values)
	}
	return
}

// PayOrderInfo user pay order info.
func (d *Dao) PayOrderInfo(c context.Context, column string, mid, money int64) (err error) {
	var (
		key         = d.rowVerKey(mid, time.Now())
		incrBytes   = make([]byte, 8)
		ctx, cancel = context.WithTimeout(c, time.Duration(conf.Conf.HBase.WriteTimeout))
	)
	defer cancel()
	log.Info("Pay Order key [%s]", key)
	binary.BigEndian.PutUint64(incrBytes, uint64(money))
	values := map[string]map[string][]byte{model.ACFamilyUser: map[string][]byte{column: incrBytes}}
	if _, err = d.hbase.Increment(ctx, model.ActionCounterTable, key, values); err != nil {
		err = errors.Wrapf(err, "msg(%d,%s,%d), hbase.Increment(key: %s values: %v)", mid, column, key, values)
	}
	return
}
