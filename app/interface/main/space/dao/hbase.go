package dao

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"time"

	"go-common/app/interface/main/space/model"
	"go-common/library/log"

	"github.com/tsuna/gohbase/hrpc"
)

const (
	_hBaseArticleTable      = "read_auth_stats_daily"
	_hBaseUpStatTablePrefix = "up_stats_"
)

func hbaseMd5Key(mid int64) string {
	hasher := md5.New()
	hasher.Write([]byte(strconv.Itoa(int(mid))))
	return hex.EncodeToString(hasher.Sum(nil))
}

// UpArcStat get up archive stat.
func (d *Dao) UpArcStat(c context.Context, mid int64, date string) (stat *model.UpArcStat, err error) {
	var (
		result      *hrpc.Result
		ctx, cancel = context.WithTimeout(c, time.Duration(d.c.HBase.ReadTimeout))
		key         = hbaseMd5Key(mid)
		tableName   = _hBaseUpStatTablePrefix + date // change table at 12:00am
	)
	defer cancel()
	if result, err = d.hbase.GetStr(ctx, tableName, key); err != nil {
		log.Error("UpArcStat d.hbase.GetStr tableName(%s)|mid(%d)|key(%v)|error(%v)", tableName, mid, key, err)
		return
	}
	if result == nil {
		return
	}
	stat = &model.UpArcStat{}
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
		}
	}
	return
}

// UpArtStat get up article stat.
func (d *Dao) UpArtStat(c context.Context, mid int64) (stat *model.UpArtStat, err error) {
	var (
		result      *hrpc.Result
		ctx, cancel = context.WithTimeout(c, time.Duration(d.c.HBase.ReadTimeout))
		key         = hbaseMd5Key(mid)
		tableName   = _hBaseArticleTable
	)
	defer cancel()
	if result, err = d.hbase.GetStr(ctx, tableName, key); err != nil {
		log.Error("UpArtStat d.hbase.GetStr tableName(%s)|mid(%d)|key(%v)|error(%v)", tableName, mid, key, err)
		return
	}
	if result == nil {
		return
	}
	stat = &model.UpArtStat{}
	for _, c := range result.Cells {
		if c == nil {
			continue
		}
		v, _ := strconv.ParseInt(string(c.Value[:]), 10, 64)
		if !bytes.Equal(c.Family, []byte("r")) {
			continue
		}
		switch {
		case bytes.Equal(c.Qualifier, []byte("view1")):
			stat.View = v
		case bytes.Equal(c.Qualifier, []byte("reply1")):
			stat.Reply = v
		case bytes.Equal(c.Qualifier, []byte("coin1")):
			stat.Coin = v
		case bytes.Equal(c.Qualifier, []byte("like1")):
			stat.Like = v
		case bytes.Equal(c.Qualifier, []byte("fav1")):
			stat.Fav = v
		case bytes.Equal(c.Qualifier, []byte("view0")):
			stat.PreView = v
		case bytes.Equal(c.Qualifier, []byte("reply0")):
			stat.PreReply = v
		case bytes.Equal(c.Qualifier, []byte("coin0")):
			stat.PreCoin = v
		case bytes.Equal(c.Qualifier, []byte("like0")):
			stat.PreLike = v
		case bytes.Equal(c.Qualifier, []byte("fav0")):
			stat.PreFav = v
		}
	}
	stat.IncrView = stat.View - stat.PreView
	stat.IncrReply = stat.Reply - stat.PreReply
	stat.IncrCoin = stat.Coin - stat.PreCoin
	stat.IncrLike = stat.Like - stat.PreLike
	stat.IncrFav = stat.Fav - stat.PreFav
	return
}
