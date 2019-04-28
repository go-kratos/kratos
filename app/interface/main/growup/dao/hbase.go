package dao

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"strconv"

	"go-common/app/interface/main/growup/model"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/tsuna/gohbase/hrpc"
)

const (
	// HBaseUpStatTablePrefix hbase up_stat_date
	HBaseUpStatTablePrefix = "up_stats_"
)

func hbaseMd5Key(aid int64) string {
	hasher := md5.New()
	hasher.Write([]byte(strconv.Itoa(int(aid))))
	return hex.EncodeToString(hasher.Sum(nil))
}

// BaseUpStat get base up stat.
func (d *Dao) BaseUpStat(c context.Context, mid int64, date string) (stat *model.UpBaseStat, err error) {
	var (
		result      *hrpc.Result
		ctx, cancel = context.WithTimeout(c, d.hbaseTimeOut)
		tableName   = HBaseUpStatTablePrefix + date // change table at 12:00am
		key         = hbaseMd5Key(mid)
	)
	defer cancel()
	if result, err = d.hbase.GetStr(ctx, tableName, key); err != nil {
		log.Error("BaseUpStat d.hbase.GetStr tableName(%s)|mid(%d)|key(%v)|error(%v)", tableName, mid, key, err)
		err = ecode.CreativeDataErr
		return
	}
	if result == nil {
		return
	}
	stat = &model.UpBaseStat{}
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
		case bytes.Equal(c.Qualifier, []byte("sh")):
			stat.Share = v
		}
	}
	return
}
