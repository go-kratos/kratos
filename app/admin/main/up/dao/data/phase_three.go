package data

import (
	"context"
	"encoding/binary"
	"strconv"
	"time"

	"go-common/app/admin/main/up/model/datamodel"
	"go-common/app/admin/main/up/util/hbaseutil"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/tsuna/gohbase/hrpc"
)

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

// 粉丝管理 - up_fans_analysis
// mid倒置补0 （补满共10位）+0 +yyyyMMdd          – 累计数据
// mid倒置补0 （补满共10位）+1 +yyyyMMdd          – 7日数据
// mid倒置补0 （补满共10位）+2 +yyyyMMdd        – 30日数据
// mid倒置补0 （补满共10位）+3 +yyyyMMdd        – 90日数据
func fansRowKey(id int64, ty int) string {
	idStr := strconv.FormatInt(id, 10)
	s := reverseString(idStr)
	s = s + strconv.Itoa(ty) + time.Now().AddDate(0, 0, -1).Add(-12*time.Hour).Format("20060102")
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

// UpFansAnalysis for web fans analysis.
func (d *Dao) UpFansAnalysis(c context.Context, mid int64, ty int) (res *datamodel.FanInfo, err error) {
	var (
		result      *hrpc.Result
		ctx, cancel = context.WithTimeout(c, time.Duration(d.c.HBase.ReadTimeout))
		tableName   = HBaseUpFansAnalysis
		//summary      = make(map[string]int32)            //total & classify proportion
		//rankDr       = make(map[string]int32)            //播放时长 mid list
		//rankVideoAct = make(map[string]int32)            //视频互动 mid list
		//rankDyAct    = make(map[string]int32)            //动态互动 mid list
		rankMap = make(map[string]map[string]int32) //top 10 rank list
		source  = make(map[string]int32)            //source  proportion
		rowkey  = fansRowKey(mid, ty)
		parser  = hbaseutil.Parser{}
	)
	res = new(datamodel.FanInfo)
	defer cancel()
	log.Info("UpFansAnalysisForWeb mid(%d)|rowkey(%s)", mid, rowkey)
	if result, err = d.hbase.GetStr(ctx, tableName, rowkey, hrpc.Families(map[string][]string{"f": nil})); err != nil {
		log.Error("d.hbase.GetStr tableName(%s)|mid(%d)|rowkey(%v)|error(%v)", tableName, mid, rowkey, err)
		err = ecode.CreativeDataErr
		return
	}
	if result == nil || len(result.Cells) == 0 {
		log.Warn("UpFansAnalysisForWeb no data tableName(%s)|mid(%d)|rowkey(%v)", tableName, mid, rowkey)
		return
	}
	var summary datamodel.FanSummaryData
	err = parser.Parse(result.Cells, &summary)
	if err != nil {
		log.Error("convert from hbase fail, err=%v", err)
		return
	}

	/*
		for _, c := range result.Cells {
			if c == nil {
				continue
			}
			if string(c.Family) == "f" {

				k := string(c.Qualifier[:])
				var v int32
				if len(c.Value) == 4 {
					v = byteToInt32(c.Value)
					switch k {
					case "all": //总粉丝
						summary.Total = v
						//summary["total"] = v
					case "inc": //新增粉丝
						summary.Inc = v
					case "act": //活跃粉丝
						summary.Active = v
						//	summary["active"] = v
						//case "mdl": //领取勋章粉丝
						//	summary["medal"] = v
						//case "elec": //充电粉丝
						//	summary["elec"] = v
						//case "act_diff": //活跃粉丝（增量）
						//	summary["act_diff"] = v
						//case "mdl_diff": //领取勋章粉丝（增量）
						//	summary["mdl_diff"] = v
						//case "elec_diff": //充电粉丝（增量）
						//	summary["elec_diff"] = v
						//case "v": //播放粉丝占比*10000
						//	summary["v"] = v
						//case "dm": //弹幕粉丝占比*10000
						//	summary["dm"] = v
						//case "r": //评论粉丝占比*10000
						//	summary["r"] = v
						//case "c": //投币粉丝占比*10000
						//	summary["c"] = v
						//	//新增粉丝活跃占比
						//case "inter": //互动活跃度*10000
						//	summary["inter"] = v
						//case "vv": //观看活跃度*10000
						//	summary["vv"] = v
						//case "da": //弹幕粉丝占比*10000
						//	summary["da"] = v
						//case "re": //评论粉丝占比*10000
						//	summary["re"] = v
						//case "co": //投币粉丝占比*10000
						//	summary["co"] = v
						//case "fv": //收藏粉丝占比*10000
						//	summary["fv"] = v
						//case "sh": //分享粉丝占比*10000
						//	summary["sh"] = v
						//case "lk": //点赞粉丝占比*10000
						//	summary["lk"] = v
					}
				} else {
					log.Error("UpFansAnalysisForWeb family[f] get dirty tableName(%s)|mid(%d)|rowkey(%s)|value(%+v)", tableName, mid, rowkey, c.Value)
				}
			}
			if string(c.Family) == "t" {
				var v int32
				if len(c.Value) == 4 {
					v = int32(binary.BigEndian.Uint32(c.Value))
				} else {
					log.Error("UpFansAnalysisForWeb family t get dirty tableName(%s)|mid(%d)|rowkey(%s)|value(%+v)", tableName, mid, fansRowKey(mid, ty), c.Value)
				}
				if strings.Contains(string(c.Qualifier[:]), "dr") {
					rankDr[string(c.Qualifier[:])] = v
				}
				if strings.Contains(string(c.Qualifier[:]), "act") {
					rankVideoAct[string(c.Qualifier[:])] = v
				}
				if strings.Contains(string(c.Qualifier[:]), "dy") {
					rankDyAct[string(c.Qualifier[:])] = v
				}
			}
			if string(c.Family) == "s" {
				var v int32
				if len(c.Value) == 4 {
					v = int32(binary.BigEndian.Uint32(c.Value))
				} else {
					log.Error("UpFansAnalysisForWeb family[t] get dirty data tableName(%s)|mid(%d)|rowkey(%s)|value(%+v)", tableName, mid, rowkey, c.Value)
				}
				if strings.Contains(string(c.Qualifier[:]), "pf") {
					if pk, ok := fanSource[string(c.Qualifier[:])]; ok {
						source[pk] = v
					}
				}
			}
		}
		for _, k := range fanSource { //粉丝来源页面占比如果数据平台未返回,则设置对应的key
			if _, ok := source[k]; !ok {
				source[k] = 0
			}
		}
	*/
	log.Info("UpFansAnalysisForWebRankMap mid(%d)|rowkey(%s)|summary(%+v)|rankMap(%+v)|source(%+v)", mid, rowkey, summary, rankMap, source)
	res.Summary = summary
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
