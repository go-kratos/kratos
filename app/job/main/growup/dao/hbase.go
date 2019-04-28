package dao

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"time"

	"go-common/app/job/main/growup/model"
	"go-common/library/log"
)

var (
	//HBaseArchiveTablePrefix 分类分端播放
	HBaseArchiveTablePrefix = "video_play_category_"
	//HBaseFamilyPlat  family
	HBaseFamilyPlat = []byte("v")
)

func hbaseMd5Key(aid int64) string {
	hasher := md5.New()
	hasher.Write([]byte(strconv.Itoa(int(aid))))
	return hex.EncodeToString(hasher.Sum(nil))
}

// ArchiveStat get the stat of archive.
func (d *Dao) ArchiveStat(c context.Context, aid int64, date time.Time) (stat *model.ArchiveStat, err error) {
	var (
		ctx, cancel = context.WithTimeout(c, time.Duration(d.c.HBase.ReadTimeout))
		tableName   = HBaseArchiveTablePrefix + date.Format("20060102")
		key         = hbaseMd5Key(aid)
	)
	defer cancel()
	result, err := d.hbase.GetStr(ctx, tableName, key)
	if err != nil {
		log.Error("ArchiveStat d.hbase.GetStr tableName(%s)|aid(%d)|key(%v)|error(%v)", tableName, aid, key, err)
		return
	}
	if result == nil {
		return
	}
	stat = &model.ArchiveStat{}
	for _, c := range result.Cells {
		if c == nil {
			continue
		}
		v, _ := strconv.ParseInt(string(c.Value[:]), 10, 64)
		if !bytes.Equal(c.Family, HBaseFamilyPlat) {
			continue
		}
		switch {
		case bytes.Equal(c.Qualifier, []byte("play")):
			stat.Play = v
		case bytes.Equal(c.Qualifier, []byte("dm")):
			stat.Dm = v
		case bytes.Equal(c.Qualifier, []byte("reply")):
			stat.Reply = v
		case bytes.Equal(c.Qualifier, []byte("like")):
			stat.Like = v
		case bytes.Equal(c.Qualifier, []byte("sh")):
			stat.Share = v
		}
	}
	return
}
