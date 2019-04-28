package data

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"go-common/app/service/main/up/model/data"
	"go-common/library/ecode"
	"go-common/library/log"
	"strconv"
)

func hbaseMd5Key(aid int64) []byte {
	hasher := md5.New()
	hasher.Write([]byte(strconv.Itoa(int(aid))))
	return []byte(hex.EncodeToString(hasher.Sum(nil)))
}

// BaseUpStat get base up stat.
func (d *Dao) BaseUpStat(c context.Context, mid int64, date string) (stat *data.UpBaseStat, err error) {
	var (
		ctx, cancel = context.WithTimeout(c, d.hbaseTimeOut)
		tableName   = HBaseUpStatTablePrefix + date // change table at 12:00am
	)
	defer cancel()
	result, err := d.hbase.Get(ctx, []byte(tableName), hbaseMd5Key(mid))
	if err != nil {
		log.Error("BaseUpStat d.hbase.Get BackupTable(%s, %d) error(%v)", tableName, mid, err)
		err = ecode.ServerErr
		return
	}
	if result == nil {
		log.Error("BaseUpStat d.hbase.Get BackupTable(%s, %d) result nil", tableName, mid)
		return
	}
	stat = &data.UpBaseStat{}
	for _, c := range result.Cells {
		if c == nil {
			continue
		}
		v, _ := strconv.ParseInt(string(c.Value[:]), 10, 64)
		if !bytes.Equal(c.Family, []byte("u")) {
			continue
		}
		switch {
		case bytes.Equal(c.Qualifier, []byte("play")):
			stat.View = v
		case bytes.Equal(c.Qualifier, []byte("dm")):
			stat.Dm = v
		case bytes.Equal(c.Qualifier, []byte("reply")):
			stat.Reply = v
		case bytes.Equal(c.Qualifier, []byte("fans")):
			stat.Fans = v
		case bytes.Equal(c.Qualifier, []byte("fav")):
			stat.Fav = v
		case bytes.Equal(c.Qualifier, []byte("like")):
			stat.Like = v
		}
	}
	return
}
