package dao

import (
	"context"
	"crypto/md5"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"go-common/app/interface/main/history/model"
	"go-common/library/log"
)

var (
	tableInfo = "ugc:history"
	family    = "info"
)

// hashRowKey create rowkey(md5(mid)[:2]+mid) for histroy by mid .
func hashRowKey(mid int64) string {
	var bs = make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, uint64(mid))
	rk := md5.Sum(bs)
	return fmt.Sprintf("%x%d", rk[:2], mid)
}

func (d *Dao) column(aid int64, typ int8) string {
	if typ < model.TypeArticle {
		return strconv.FormatInt(aid, 10)
	}
	return fmt.Sprintf("%d_%d", aid, typ)
}

// Add add history list.
func (d *Dao) Add(ctx context.Context, h *model.History) error {
	valueByte, err := json.Marshal(h)
	if err != nil {
		log.Error("json.Marshal(%v) error(%v)", h, err)
		return err
	}
	fValues := make(map[string][]byte)
	column := d.column(h.Aid, h.TP)
	fValues[column] = valueByte
	key := hashRowKey(h.Mid)
	values := map[string]map[string][]byte{family: fValues}
	ctx, cancel := context.WithTimeout(ctx, time.Duration(d.conf.Info.WriteTimeout))
	defer cancel()
	if _, err = d.info.PutStr(ctx, tableInfo, key, values); err != nil {
		log.Error("info.PutStr error(%v)", err)
	}
	return nil
}
