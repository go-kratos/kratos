package data

import (
	"context"
	"encoding/binary"

	"strconv"
	"strings"
	"time"

	"go-common/app/interface/main/creative/model/data"

	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/tsuna/gohbase/hrpc"
)

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

// up主粉丝勋章 - up_fans_medal
// mid倒置补0（补满共10位）+1 +yyyyMMdd
func upFansMedalRowKey(id int64, ty int) string {
	idStr := strconv.FormatInt(id, 10)
	s := reverseString(idStr)
	s = s + strconv.Itoa(ty) + time.Now().AddDate(0, 0, -1).Add(-12*time.Hour).Format("20060102")
	return s
}

func byteToInt32(b []byte) int32 {
	return int32(b[3]) | int32(b[2])<<8 | int32(b[1])<<16 | int32(b[0])<<24
}

// UpFansAnalysisForApp for  app fans analysis.
func (d *Dao) UpFansAnalysisForApp(c context.Context, mid int64, ty int) (res *data.AppFan, err error) {
	var (
		result       *hrpc.Result
		ctx, cancel  = context.WithTimeout(c, time.Duration(d.c.HBase.ReadTimeout))
		tableName    = HBaseUpFansAnalysis
		rowkey       = fansRowKey(mid, ty)
		summary      = make(map[string]int64)            //total & classify proportion
		rankDr       = make(map[string]int32)            //播放时长 mid list
		rankVideoAct = make(map[string]int32)            //视频互动 mid list
		rankMap      = make(map[string]map[string]int32) //top 10 rank list
	)
	defer cancel()
	log.Info("UpFansAnalysisForApp mid(%d)|rowkey(%s)", mid, rowkey)
	if result, err = d.hbase.GetStr(ctx, tableName, rowkey); err != nil {
		log.Error("d.hbase.GetStr tableName(%s)|mid(%d)|rowkey(%v)|error(%v)", tableName, mid, rowkey, err)
		err = ecode.CreativeDataErr
		return
	}
	if result == nil || len(result.Cells) == 0 {
		log.Warn("UpFansAnalysisForApp no data tableName(%s)|mid(%d)|", tableName, mid)
		return
	}
	for _, c := range result.Cells {
		if c == nil {
			continue
		}
		if string(c.Family) == "f" {
			if len(c.Value) != 4 {
				log.Error("UpFansAnalysisForApp family[f] get dirty value tableName(%s)|mid(%d)|rowkey(%s)|value(%+v)", tableName, mid, rowkey, c.Value)
				continue
			}
			v := int64(byteToInt32(c.Value))
			switch string(c.Qualifier[:]) {
			case "all": //总粉丝
				summary["total"] = v
			case "inc": //新增粉丝
				summary["inc"] = v
			case "act": //活跃粉丝
				summary["active"] = v
			case "v":
				summary["play"] = v
			case "dm":
				summary["dm"] = v
			case "r":
				summary["reply"] = v
			case "c":
				summary["coin"] = v
			case "inter": //互动活跃度*10000
				summary["inter"] = v
			case "vv": //观看活跃度*10000
				summary["vv"] = v
			case "da": //弹幕粉丝占比*10000
				summary["da"] = v
			case "re": //评论粉丝占比*10000
				summary["re"] = v
			case "co": //投币粉丝占比*10000
				summary["co"] = v
			case "fv": //收藏粉丝占比*10000
				summary["fv"] = v
			case "sh": //分享粉丝占比*10000
				summary["sh"] = v
			case "lk": //点赞粉丝占比*10000
				summary["lk"] = v
			}
			log.Info("UpFansAnalysisForApp family[f] value tableName(%s)|mid(%d)|rowkey(%s)|key(%s)|value(%+v|Uint32[%+v]|int32[%+v])", tableName, mid, rowkey, string(c.Qualifier[:]), c.Value, binary.BigEndian.Uint32(c.Value), v)
		}
		if string(c.Family) == "t" {
			var v int32
			if len(c.Value) == 4 {
				v = int32(binary.BigEndian.Uint32(c.Value))
			} else {
				log.Error("UpFansAnalysisForApp family t get dirty tableName(%s)|mid(%d)|rowkey(%s)|value(%+v)", tableName, mid, fansRowKey(mid, ty), c.Value)
			}
			if strings.Contains(string(c.Qualifier[:]), "dr") {
				rankDr[string(c.Qualifier[:])] = v
			}
			if strings.Contains(string(c.Qualifier[:]), "act") {
				rankVideoAct[string(c.Qualifier[:])] = v
			}
		}
	}
	rankMap[data.PlayDuration] = rankDr
	rankMap[data.VideoAct] = rankVideoAct
	log.Info("UpFansAnalysisForApp mid(%d)|rowkey(%s)|summary(%+v)|rankMap(%+v)", mid, rowkey, summary, rankMap)
	res = &data.AppFan{
		Summary: summary,
		RankMap: rankMap,
	}
	return
}

// UpFansAnalysisForWeb for web fans analysis.
func (d *Dao) UpFansAnalysisForWeb(c context.Context, mid int64, ty int) (res *data.WebFan, err error) {
	var (
		result       *hrpc.Result
		ctx, cancel  = context.WithTimeout(c, time.Duration(d.c.HBase.ReadTimeout))
		tableName    = HBaseUpFansAnalysis
		summary      = make(map[string]int32)            //total & classify proportion
		rankDr       = make(map[string]int32)            //播放时长 mid list
		rankVideoAct = make(map[string]int32)            //视频互动 mid list
		rankDyAct    = make(map[string]int32)            //动态互动 mid list
		rankMap      = make(map[string]map[string]int32) //top 10 rank list
		source       = make(map[string]int32)            //source  proportion
		rowkey       = fansRowKey(mid, ty)
	)
	defer cancel()
	log.Info("UpFansAnalysisForWeb mid(%d)|rowkey(%s)", mid, rowkey)
	if result, err = d.hbase.GetStr(ctx, tableName, rowkey); err != nil {
		log.Error("d.hbase.GetStr tableName(%s)|mid(%d)|rowkey(%v)|error(%v)", tableName, mid, rowkey, err)
		err = ecode.CreativeDataErr
		return
	}
	if result == nil || len(result.Cells) == 0 {
		log.Warn("UpFansAnalysisForWeb no data tableName(%s)|mid(%d)|rowkey(%v)", tableName, mid, rowkey)
		return
	}
	var cells data.FansAnalysis
	err = parser.Parse(result.Cells, &cells)
	for k, v := range cells.F {
		if k == "all" { //总粉丝
			summary["total"] = v
		} else if k == "act" { //活跃粉丝
			summary["active"] = v
		} else if k == "mdl" { //领取勋章粉丝
			summary["medal"] = v
		} else {
			summary[k] = v
		}
	}
	for k, v := range cells.T {
		if strings.Contains(k, "dr") {
			rankDr[k] = v
		} else if strings.Contains(k, "act") {
			rankVideoAct[k] = v
		} else if strings.Contains(k, "dy") {
			rankDyAct[k] = v
		}
	}
	for k, v := range cells.S {
		if strings.Contains(k, "pf") {
			if pk, ok := fanSource[k]; ok {
				source[pk] = v
			}
		}
	}
	for _, k := range fanSource { //粉丝来源页面占比如果数据平台未返回,则设置对应的key
		if _, ok := source[k]; !ok {
			source[k] = 0
		}
	}
	rankMap[data.PlayDuration] = rankDr
	rankMap[data.VideoAct] = rankVideoAct
	rankMap[data.DynamicAct] = rankDyAct
	log.Info("UpFansAnalysisForWebRankMap mid(%d)|rowkey(%s)|summary(%+v)|rankMap(%+v)|source(%+v)", mid, rowkey, summary, rankMap, source)
	res = &data.WebFan{
		Summary: summary,
		RankMap: rankMap,
		Source:  source,
	}
	return
}

// UpFansMedal get 领取勋章数+佩戴勋章数.
func (d *Dao) UpFansMedal(c context.Context, mid int64) (res *data.UpFansMedal, err error) {
	var (
		result      *hrpc.Result
		ctx, cancel = context.WithTimeout(c, time.Duration(d.c.HBase.ReadTimeout))
		tableName   = HBaseUpFansMedalQuery
		rowkey      = upFansMedalRowKey(mid, 1)
	)
	defer cancel()
	log.Info("UpFansMedal aid(%d)|rowkey(%s)", mid, rowkey)
	if result, err = d.hbase.GetStr(ctx, tableName, rowkey); err != nil {
		log.Error("d.hbase.GetStr tableName(%s)|mid(%d)|rowkey(%+v)|error(%v)", tableName, mid, rowkey, err)
		err = ecode.CreativeDataErr
		return
	}
	if result == nil || len(result.Cells) == 0 {
		log.Warn("UpFansMedal no data tableName(%s)|mid(%d)|rowkey(%+v)", tableName, mid, rowkey)
		return
	}
	var cells data.UpFansMedal
	err = parser.Parse(result.Cells, &cells)
	if err != nil {
		log.Error("parser.Parse tableName(%s)|mid(%d)|rowkey(%+v)|error(%v)", tableName, mid, rowkey, err)
		err = ecode.CreativeDataErr
		return
	}
	res = &cells
	log.Info("UpFansMedal mid(%d)|rowkey(%s)|res(%+v)", mid, rowkey, *res)
	return
}
