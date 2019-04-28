package data

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"sort"
	"strconv"
	"time"

	"go-common/app/interface/main/creative/model/data"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/tsuna/gohbase/hrpc"

	"golang.org/x/net/context"
)

func hbaseMd5Key(aid int64) string {
	hasher := md5.New()
	hasher.Write([]byte(strconv.Itoa(int(aid))))
	return hex.EncodeToString(hasher.Sum(nil))
}

// VideoQuitPoints get video quit points.
func (d *Dao) VideoQuitPoints(c context.Context, cid int64) (res []int64, err error) {
	var (
		ctx, cancel     = context.WithTimeout(c, time.Duration(d.c.HBase.ReadTimeout))
		tableName       = HBaseVideoTablePrefix + time.Now().AddDate(0, 0, -1).Add(-12*time.Hour).Format("20060102") // change table at 12:00am
		backupTableName = HBaseVideoTablePrefix + time.Now().AddDate(0, 0, -2).Add(-12*time.Hour).Format("20060102") // change table at 12:00am
		key             = hbaseMd5Key(cid)
	)
	defer cancel()
	result, err := d.hbase.GetStr(ctx, tableName, key)
	if err != nil {
		if result, err = d.hbase.GetStr(ctx, backupTableName, key); err != nil {
			log.Error("VideoQuitPoints d.hbase.GetStr backupTableName(%s)|cid(%d)|key(%v)|error(%v)", backupTableName, cid, key, err)
			err = ecode.CreativeDataErr
			return
		}
	}
	if result == nil {
		return
	}
	// get parts and max part for fill res
	partMap := make(map[int]int64)
	maxPart := 0
	for _, c := range result.Cells {
		if c != nil {
			part, _ := strconv.Atoi(string(c.Qualifier[:]))
			per, _ := strconv.ParseInt(string(c.Value[:]), 10, 64)
			partMap[part] = per
			if part > maxPart {
				maxPart = part
			}
		}
	}
	var restPercent int64 = 10000 // start from 100%
	for i := 1; i <= maxPart; i++ {
		if _, ok := partMap[i]; ok {
			restPercent = restPercent - partMap[i]
		}
		res = append(res, restPercent)
	}
	return
}

// ArchiveStat get the stat of archive.
func (d *Dao) ArchiveStat(c context.Context, aid int64) (stat *data.ArchiveData, err error) {
	var (
		ctx, cancel     = context.WithTimeout(c, time.Duration(d.c.HBase.ReadTimeout))
		tableName       = HBaseArchiveTablePrefix + time.Now().AddDate(0, 0, -1).Add(-12*time.Hour).Format("20060102") // change table at 12:00am
		backupTableName = HBaseArchiveTablePrefix + time.Now().AddDate(0, 0, -2).Add(-12*time.Hour).Format("20060102") // change table at 12:00am
		key             = hbaseMd5Key(aid)
	)
	defer cancel()
	result, err := d.hbase.GetStr(ctx, tableName, key)
	if err != nil {
		if result, err = d.hbase.GetStr(ctx, backupTableName, key); err != nil {
			log.Error("ArchiveStat d.hbase.GetStr backupTableName(%s)|aid(%d)|key(%v)|error(%v)", backupTableName, aid, key, err)
			err = ecode.CreativeDataErr
			return
		}
	}
	if result == nil {
		return
	}
	stat = &data.ArchiveData{}
	stat.ArchiveSource = &data.ArchiveSource{}
	stat.ArchiveGroup = &data.ArchiveGroup{}
	stat.ArchiveStat = &data.ArchiveStat{}
	for _, c := range result.Cells {
		if c == nil {
			continue
		}
		v, _ := strconv.ParseInt(string(c.Value[:]), 10, 64)
		if !bytes.Equal(c.Family, HBaseFamilyPlat) {
			continue
		}
		switch {
		case bytes.Equal(c.Qualifier, HBaseColumnWebPC):
			stat.ArchiveSource.WebPC = v
		case bytes.Equal(c.Qualifier, HBaseColumnWebH5):
			stat.ArchiveSource.WebH5 = v
		case bytes.Equal(c.Qualifier, HBaseColumnOutsite):
			stat.ArchiveSource.Outsite = v
		case bytes.Equal(c.Qualifier, HBaseColumnIOS):
			stat.ArchiveSource.IOS = v
		case bytes.Equal(c.Qualifier, HBaseColumnAndroid):
			stat.ArchiveSource.Android = v
		case bytes.Equal(c.Qualifier, HBaseColumnElse):
			stat.ArchiveSource.Others = v
		case bytes.Equal(c.Qualifier, HBaseColumnFans):
			stat.ArchiveGroup.Fans = v
		case bytes.Equal(c.Qualifier, HBaseColumnGuest):
			stat.ArchiveGroup.Guest = v
		case bytes.Equal(c.Qualifier, HBaseColumnAll):
			stat.ArchiveStat.Play = v
		case bytes.Equal(c.Qualifier, HBaseColumnCoin):
			stat.ArchiveStat.Coin = v
		case bytes.Equal(c.Qualifier, HBaseColumnElec):
			stat.ArchiveStat.Elec = v
		case bytes.Equal(c.Qualifier, HBaseColumnFav):
			stat.ArchiveStat.Fav = v
		case bytes.Equal(c.Qualifier, HBaseColumnShare):
			stat.ArchiveStat.Share = v
			//【稿件分析】增加：播放、弹幕、评论、点赞 v:play v:danmu v:reply v:likes
		case bytes.Equal(c.Qualifier, []byte("play")):
			stat.ArchiveStat.Play = v
		case bytes.Equal(c.Qualifier, []byte("danmu")):
			stat.ArchiveStat.Dm = v
		case bytes.Equal(c.Qualifier, []byte("reply")):
			stat.ArchiveStat.Reply = v
		case bytes.Equal(c.Qualifier, []byte("likes")):
			stat.ArchiveStat.Like = v
		}
	}
	stat.ArchiveSource.Mainsite = stat.ArchiveSource.WebPC + stat.ArchiveSource.WebH5
	stat.ArchiveSource.Mobile = stat.ArchiveSource.Android + stat.ArchiveSource.IOS
	return
}

// ArchiveArea get the count of area.
func (d *Dao) ArchiveArea(c context.Context, aid int64) (res []*data.ArchiveArea, err error) {
	var (
		ctx, cancel     = context.WithTimeout(c, time.Duration(d.c.HBase.ReadTimeout))
		tableName       = HBaseAreaTablePrefix + time.Now().AddDate(0, 0, -1).Add(-12*time.Hour).Format("20060102") // change table at 12:00am
		backupTableName = HBaseAreaTablePrefix + time.Now().AddDate(0, 0, -2).Add(-12*time.Hour).Format("20060102") // change table at 12:00am
		key             = hbaseMd5Key(aid)
	)
	defer cancel()
	result, err := d.hbase.GetStr(ctx, tableName, key)
	if err != nil {
		if result, err = d.hbase.GetStr(ctx, backupTableName, key); err != nil {
			log.Error("ArchiveArea d.hbase.GetStr backupTableName(%s)|aid(%d)|key(%v)|error(%v)", backupTableName, aid, key, err)
			err = ecode.CreativeDataErr
			return
		}
	}
	if result == nil {
		return
	}
	var countArr []int
	countMap := make(map[int64][]*data.ArchiveArea)
	countSet := make(map[int]struct{}) // empty struct{} for saving memory
	for _, c := range result.Cells {
		if c != nil {
			area := &data.ArchiveArea{}
			area.Location = string(c.Qualifier[:])
			area.Count, _ = strconv.ParseInt(string(c.Value[:]), 10, 64)
			countMap[area.Count] = append(countMap[area.Count], area)
			countSet[int(area.Count)] = struct{}{}
		}
	}
	for key := range countSet {
		countArr = append(countArr, key)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(countArr)))
	for _, c := range countArr {
		if _, ok := countMap[int64(c)]; ok {
			res = append(res, countMap[int64(c)]...)
		}
		if len(res) >= 10 {
			res = res[:10] // exact 10 item
			break
		}
	}
	return
}

// BaseUpStat get base up stat.
func (d *Dao) BaseUpStat(c context.Context, mid int64, date string) (stat *data.UpBaseStat, err error) {
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
		case bytes.Equal(c.Qualifier, []byte("sh")):
			stat.Share = v
		case bytes.Equal(c.Qualifier, []byte("elec")): //【视频数据总览】增加：硬币、充电
			stat.Elec = v
		case bytes.Equal(c.Qualifier, []byte("coin")):
			stat.Coin = v
		}
	}
	log.Info("BaseUpStat d.hbase.GetStr tableName(%s)|mid(%d)|key(%v)|stat(%+v)", tableName, mid, key, stat)
	return
}

// UpArchiveStatQuery 获取最高播放/评论/弹幕/...数
func (d *Dao) UpArchiveStatQuery(c context.Context, mid int64, date string) (res *data.ArchiveMaxStat, err error) {
	var (
		result      *hrpc.Result
		ctx, cancel = context.WithTimeout(c, time.Duration(d.c.HBase.ReadTimeout))
		tableName   = HBaseupArchiveStatQuery + date // change table at 12:00am
		rowkey      = hbaseMd5Key(mid)
	)
	defer cancel()
	log.Info("UpArchiveStatQuery aid(%d)|rowkey(%s)", mid, rowkey)
	if result, err = d.hbase.GetStr(ctx, tableName, rowkey); err != nil {
		log.Error("UpArchiveStatQuery d.hbase.GetStr tableName(%s)|mid(%d)|rowkey(%+v)|error(%v)", tableName, mid, rowkey, err)
		err = ecode.CreativeDataErr
		return
	}
	if result == nil || len(result.Cells) == 0 {
		log.Warn("UpArchiveStatQuery no data tableName(%s)|mid(%d)|rowkey(%+v)", tableName, mid, rowkey)
		return
	}
	var cells data.ArchiveMaxStat
	err = parser.Parse(result.Cells, &cells)
	if err != nil {
		log.Error("UpArchiveStatQuery parser.Parse tableName(%s)|mid(%d)|rowkey(%+v)|error(%v)", tableName, mid, rowkey, err)
		err = ecode.CreativeDataErr
		return
	}
	res = &cells
	log.Info("UpArchiveStatQuery mid(%d)|rowkey(%s)|res(%+v)", mid, rowkey, *res)
	return
}
