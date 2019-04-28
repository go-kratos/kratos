package data

import (
	"context"
	"encoding/binary"
	"strconv"
	"time"

	"go-common/app/admin/main/up/util/hbaseutil"
	"go-common/app/interface/main/creative/model/data"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/tsuna/gohbase/hrpc"
)

var (
	//播放端占比
	playProportion = map[string]string{
		"pc":  "pc",      //pc端播放占比*10000
		"h5":  "h5",      //h5端播放占比*10000
		"out": "out",     //站外播放占比*10000
		"adr": "android", //android端播放占比*10000
		"ios": "ios",     //ios端播放占比*10000
	}
	//播放来源页面占比
	pageSource = map[string]string{
		"pv0":  "other",         //其他
		"pv1":  "tenma",         //天马推荐
		"pv2":  "related_video", //相关视频
		"pv3":  "search",        //搜索
		"pv4":  "h5",            //H5页面
		"pv5":  "space",         //空间
		"pv6":  "dynamic",       //动态
		"pv7":  "history",       //播放历史
		"pv8":  "tag",           //标签
		"pv9":  "cache",         //离线缓存
		"pv10": "rank",          //排行榜
		"pv11": "type",          //分区
	}
	//粉丝来源页面占比
	fanSource = map[string]string{
		"pf0": "other",      //其他
		"pf1": "space",      //主站空间
		"pf2": "main",       //主站播放页
		"pf3": "main_other", //主站其他
		"pf4": "live",       //直播
		"pf5": "audio",      //	音乐
		"pf6": "article",    //	文章
	}
	parser = hbaseutil.Parser{}
)

func sourceOtherMerge(v string) bool { //合并这些页面来源为其他
	if v == "other" || v == "h5" || v == "history" || v == "cache" ||
		v == "rank" || v == "type" || v == "tag" {
		return true
	}
	return false
}

// reverse for string.
func reverseString(s string) string {
	rs := []rune(s)
	l := len(rs)
	for f, t := 0, l-1; f < t; f, t = f+1, t-1 {
		rs[f], rs[t] = rs[t], rs[f]
	}
	ns := string(rs)
	if l < 10 {
		for i := 0; i < 10-l; i++ {
			ns = ns + "0"
		}
	}
	return ns
}

// 播放来源 - up_play_analysis
// mid倒置补（10位）+ yyyyMMdd
func playSourceKey(id int64) string {
	idStr := strconv.FormatInt(id, 10)
	s := reverseString(idStr)
	s = s + time.Now().AddDate(0, 0, -1).Add(-12*time.Hour).Format("20060102")
	return s
}

// 平均观看时长、播放用户数、留存率 - up_archive_play_analysis
// avid倒置补（10位）+ yyyyMMdd
func arcPlayKey(id int64) string {
	idStr := strconv.FormatInt(id, 10)
	s := reverseString(idStr)
	s = s + time.Now().AddDate(0, 0, -1).Add(-12*time.Hour).Format("20060102")
	return s
}

// 稿件索引表 - up_archive_query
// [mid倒置补0（10位）] + [投稿年月（4位）] + [原创/转载（1位）]
func arcQueryKey(id int64, dt string, cp int) string {
	idStr := strconv.FormatInt(id, 10)
	cpStr := strconv.Itoa(cp)
	s := reverseString(idStr)
	s = s + dt + cpStr
	return s
}

// UpPlaySourceAnalysis for play analysis.
func (d *Dao) UpPlaySourceAnalysis(c context.Context, mid int64) (res *data.PlaySource, err error) {
	var (
		result      *hrpc.Result
		ctx, cancel = context.WithTimeout(c, time.Duration(d.c.HBase.ReadTimeout))
		tableName   = HBaseUpPlaySourceAnalysis
		rowkey      = playSourceKey(mid)
	)
	defer cancel()
	log.Info("UpPlaySourceAnalysis mid(%d)|rowkey(%s)", mid, rowkey)
	if result, err = d.hbase.GetStr(ctx, tableName, rowkey); err != nil {
		log.Error("d.hbase.GetStr tableName(%s)|mid(%d)|rowkey(%s)|error(%v)", tableName, mid, rowkey, err)
		err = ecode.CreativeDataErr
		return
	}
	if result == nil || len(result.Cells) == 0 {
		log.Warn("UpPlaySourceAnalysis no data tableName(%s)|mid(%d)|rowkey(%s)", tableName, mid, rowkey)
		return
	}
	pp := make(map[string]int32)
	ps := make(map[string]int32)
	for _, c := range result.Cells {
		if c == nil {
			continue
		}
		if string(c.Family) == "s" {
			var v int32
			if len(c.Value) == 4 {
				v = int32(binary.BigEndian.Uint32(c.Value))
			} else {
				log.Error("UpPlaySourceAnalysis family[s] get dirty value  tableName(%s)|mid(%d)|rowkey(%s)", tableName, mid, rowkey)
			}
			if pyk, ok := playProportion[string(c.Qualifier[:])]; ok {
				pp[pyk] = v
			}
			if pk, ok := pageSource[string(c.Qualifier[:])]; ok {
				ps[pk] = v
			}
		}
	}
	for _, k := range playProportion { //播放平台设置数据平台未返回的key
		if _, ok := pp[k]; !ok {
			pp[k] = 0
		}
	}
	for _, k := range pageSource { //播放页面设置数据平台未返回的key
		if _, ok := ps[k]; !ok {
			ps[k] = 0
		}
	}
	var other int32
	for k, v := range ps {
		if sourceOtherMerge(k) { //如果该页面来源被计算入其他则删除该页面来源对应的key
			other = other + v
			delete(ps, k)
		}
	}
	ps["other"] = other
	res = &data.PlaySource{
		PlayProportion: pp,
		PageSource:     ps,
	}
	log.Info("UpPlaySourceAnalysis PlayProportion(%+v)|PageSource(%+v)|rowkey(%s)", pp, ps, rowkey)
	return
}

// UpArcPlayAnalysis for  arc play analysis.
func (d *Dao) UpArcPlayAnalysis(c context.Context, aid int64) (res *data.ArchivePlay, err error) {
	var (
		result      *hrpc.Result
		ctx, cancel = context.WithTimeout(c, time.Duration(d.c.HBase.ReadTimeout))
		tableName   = HBaseUpArcPlayAnalysis
		rowkey      = arcPlayKey(aid)
	)
	defer cancel()
	log.Info("UpArcPlayAnalysis aid(%d)|rowkey(%s)", aid, rowkey)
	if result, err = d.hbase.GetStr(ctx, tableName, rowkey); err != nil {
		log.Error("d.hbase.GetStr tableName(%s)|aid(%d)|rowkey(%+v)|error(%v)", tableName, aid, rowkey, err)
		err = ecode.CreativeDataErr
		return
	}
	if result == nil || len(result.Cells) == 0 {
		log.Warn("UpArcPlayAnalysis no data tableName(%s)|aid(%d)|rowkey(%+v)", tableName, aid, rowkey)
		return
	}
	ap := &data.ArchivePlay{}
	for _, c := range result.Cells {
		if c == nil {
			continue
		}
		if string(c.Family) == "v" {
			ap.AID = aid
			switch string(c.Qualifier[:]) {
			case "ctime":
				if len(c.Value) == 4 {
					ap.CTime = int32(binary.BigEndian.Uint32(c.Value))
				} else {
					log.Error("UpArcPlayAnalysis family[ctime] get dirty value  tableName(%s)|aid(%d)|rowkey(%s)", tableName, aid, rowkey)
				}
			case "uv":
				if len(c.Value) == 4 {
					ap.View = int32(binary.BigEndian.Uint32(c.Value))
				} else {
					log.Error("UpArcPlayAnalysis family[uv] get dirty value  tableName(%s)|aid(%d)|rowkey(%s)", tableName, aid, rowkey)
				}
			case "dur":
				if len(c.Value) == 8 {
					ap.Duration = int64(binary.BigEndian.Uint64(c.Value))
				} else {
					log.Error("UpArcPlayAnalysis family[dur] get dirty value  tableName(%s)|aid(%d)|rowkey(%s)", tableName, aid, rowkey)
				}
			case "avg_dur":
				if len(c.Value) == 8 {
					ap.AvgDuration = int64(binary.BigEndian.Uint64(c.Value))
				} else {
					log.Error("UpArcPlayAnalysis family[avg_dur] get dirty value  tableName(%s)|aid(%d)|rowkey(%s)", tableName, aid, rowkey)
				}
			case "rate":
				if len(c.Value) == 4 {
					ap.Rate = int32(binary.BigEndian.Uint32(c.Value))
				} else {
					log.Error("UpArcPlayAnalysis family[rate] get dirty value  tableName(%s)|aid(%d)|rowkey(%s)", tableName, aid, rowkey)
				}
			}
		}
	}
	res = ap
	log.Info("UpArcPlayAnalysis aid(%d)|rowkey(%s)|res(%+v)", aid, rowkey, res)
	return
}

// UpArcQuery for play aids by mid.
func (d *Dao) UpArcQuery(c context.Context, mid int64, dt string, cp int) (res []int64, err error) {
	var (
		result      *hrpc.Result
		ctx, cancel = context.WithTimeout(c, time.Duration(d.c.HBase.ReadTimeout))
		tableName   = HBaseUpArcQuery
		rowkey      = arcQueryKey(mid, dt, cp)
	)
	defer cancel()
	log.Info("UpArcQuery mid(%d)|rowkey(%s)", mid, rowkey)
	if result, err = d.hbase.GetStr(ctx, tableName, rowkey); err != nil {
		log.Error("d.hbase.GetStr tableName(%s)|mid(%d)|rowkey(%s)|error(%v)", tableName, mid, rowkey, err)
		err = ecode.CreativeDataErr
		return
	}
	if result == nil || len(result.Cells) == 0 {
		log.Warn("UpArcQuery no data tableName(%s)|mid(%d)|rowkey(%s)", tableName, mid, rowkey)
		return
	}
	res = make([]int64, 0)
	for _, c := range result.Cells {
		if c == nil {
			continue
		}
		if string(c.Family) == "av" {
			var v int64
			if len(c.Value) == 4 {
				v = int64(binary.BigEndian.Uint32(c.Value))
			} else {
				log.Error("UpArcQuery family[av] get dirty value  tableName(%s)|rowkey(%s)", tableName, rowkey)
			}
			res = append(res, v)
		}
	}
	log.Info("UpArcQuery mid(%d)|rowkey(%s)|res(%+v)", mid, rowkey, res)
	return
}
