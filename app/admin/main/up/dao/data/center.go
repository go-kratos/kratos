package data

import (
	"bytes"
	"strconv"
	"time"

	"go-common/app/admin/main/up/model/datamodel"
	"go-common/app/admin/main/up/util/hbaseutil"
	"go-common/app/interface/main/creative/model/data"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/tsuna/gohbase/hrpc"
	"golang.org/x/net/context"
)

// ViewerBase visitor data analysis.
func (d *Dao) ViewerBase(c context.Context, mid int64, dt time.Time) (res *datamodel.ViewerBaseInfo, err error) {
	var (
		result       *hrpc.Result
		ctx, cancel  = context.WithTimeout(c, time.Duration(d.c.HBase.ReadTimeout))
		key          = hbaseMd5Key(mid)
		genTableName = generateTableNameFunc(HBaseUpViewerBase, dt, -7)
	)
	defer cancel()

	if result, err = d.getDataWithBackup(ctx, d.hbase, genTableName, 2, key); err != nil {
		log.Error("ViewerBase d.hbase.GetStr tableName(%s)|mid(%d)|key(%v)|error(%v)", HBaseUpViewerBase, mid, key, err)
		err = ecode.CreativeDataErr
		return
	}
	if result == nil || len(result.Cells) == 0 {
		log.Warn("ViewerBase no data tableName(%s)|mid(%d)|key(%v)", HBaseUpViewerBase, mid, key)
		return
	}
	res = &datamodel.ViewerBaseInfo{}
	families := make(map[string]string, 2)
	families["f"] = ""
	families["g"] = ""
	fb := &datamodel.ViewerBaseData{}
	gb := &datamodel.ViewerBaseData{}
	for _, c := range result.Cells {
		if c == nil {
			continue
		}
		if _, ok := families[string(c.Family)]; ok {
			var data = gb
			if string(c.Family) == "f" {
				data = fb
			}
			v, _ := strconv.ParseInt(string(c.Value[:]), 10, 64)
			/*
				f:plat0	web-pc播放
				f:plat1	web-h5播放
				f:plat2	站外播放
				f:plat3	ios播放
				f:plat4	android播放
			*/
			switch {
			case bytes.Equal(c.Qualifier, []byte("male")):
				data.Male = v
			case bytes.Equal(c.Qualifier, []byte("female")):
				data.Female = v
			case bytes.Equal(c.Qualifier, []byte("age1")):
				data.Age1 = v
			case bytes.Equal(c.Qualifier, []byte("age2")):
				data.Age2 = v
			case bytes.Equal(c.Qualifier, []byte("age3")):
				data.Age3 = v
			case bytes.Equal(c.Qualifier, []byte("age4")):
				data.Age4 = v
			case bytes.Equal(c.Qualifier, []byte("plat0")):
				data.PlatPC = v
			case bytes.Equal(c.Qualifier, []byte("plat1")):
				data.PlatH5 = v
			case bytes.Equal(c.Qualifier, []byte("plat2")):
				data.PlatOut = v
			case bytes.Equal(c.Qualifier, []byte("plat3")):
				data.PlatIOS = v
			case bytes.Equal(c.Qualifier, []byte("plat4")):
				data.PlatAndroid = v
			case bytes.Equal(c.Qualifier, []byte("else")):
				data.PlatOtherApp = v
			}
		}
	}
	res.Fans = fb
	res.Guest = gb
	return
}

// ViewerArea visitor area data analysis.
func (d *Dao) ViewerArea(c context.Context, mid int64, dt time.Time) (res *datamodel.ViewerAreaInfo, err error) {
	var (
		result       *hrpc.Result
		ctx, cancel  = context.WithTimeout(c, time.Duration(d.c.HBase.ReadTimeout))
		key          = hbaseMd5Key(mid)
		genTableName = generateTableNameFunc(HBaseUpViewerArea, dt, -7)
	)
	defer cancel()
	if result, err = d.getDataWithBackup(ctx, d.hbase, genTableName, 2, key); err != nil {
		log.Error("ViewerArea d.hbase.GetStr tableName(%s)|mid(%d)|key(%v)|error(%v)", HBaseUpViewerArea, mid, key, err)
		err = ecode.CreativeDataErr
		return
	}
	if result == nil || len(result.Cells) == 0 {
		log.Warn("ViewerArea no data tableName(%s)|mid(%d)|key(%v)", HBaseUpViewerArea, mid, key)
		return
	}
	res = &datamodel.ViewerAreaInfo{}
	families := make(map[string]string, 2)
	families["f"] = ""
	families["g"] = ""
	for _, c := range result.Cells {
		if c == nil {
			continue
		}
		if _, ok := families[string(c.Family)]; ok {
			a := string(c.Qualifier[:])
			var list = &res.Guest
			v, _ := strconv.ParseInt(string(c.Value[:]), 10, 64)
			if string(c.Family) == "f" {
				list = &res.Fans
			}
			*list = append(*list, datamodel.ViewerAreaData{Area: a, Viewers: v})
		}
	}

	return
}

// ViewerTrend visitor trend data analysis.
func (d *Dao) ViewerTrend(c context.Context, mid int64, dt time.Time) (res map[string]*data.Trend, err error) {
	var (
		result       *hrpc.Result
		ctx, cancel  = context.WithTimeout(c, time.Duration(d.c.HBase.ReadTimeout))
		genTableName = generateTableNameFunc(HBaseUpViewerTrend, dt, -7)
		key          = hbaseMd5Key(mid)
	)
	defer cancel()
	if result, err = d.getDataWithBackup(ctx, d.hbase, genTableName, 2, key); err != nil {
		log.Error("ViewerTrend d.hbase.GetStr tableName(%s)|mid(%d)|key(%v)|error(%v)", HBaseUpViewerTrend, mid, key, err)
		err = ecode.CreativeDataErr
		return
	}
	if result == nil || len(result.Cells) == 0 {
		log.Warn("ViewerTrend no data tableName(%s)|mid(%d)|key(%v)", HBaseUpViewerTrend, mid, key)
		return
	}
	families := make(map[string]string, 4)
	families["fs"] = ""
	families["ft"] = ""
	families["gs"] = ""
	families["gt"] = ""
	res = make(map[string]*data.Trend)
	ftd := &data.Trend{}
	gtd := &data.Trend{}
	ty := make(map[int]int64)
	tg := make(map[int]int64)
	nty := make(map[int]int64)
	ntg := make(map[int]int64)
	for _, c := range result.Cells {
		if c == nil {
			continue
		}
		if _, ok := families[string(c.Family)]; ok {
			if string(c.Family) == "fs" {
				tid, _ := strconv.Atoi(string(c.Qualifier[:]))
				v, _ := strconv.ParseInt(string(c.Value[:]), 10, 64)
				ty[tid] = v
				ftd.Ty = ty
			} else if string(c.Family) == "gs" {
				tid, _ := strconv.Atoi(string(c.Qualifier[:]))
				v, _ := strconv.ParseInt(string(c.Value[:]), 10, 64)
				nty[tid] = v
				gtd.Ty = nty
			}
			if string(c.Family) == "ft" {
				o, _ := strconv.Atoi(string(c.Qualifier[:]))
				v, _ := strconv.ParseInt(string(c.Value[:]), 10, 64)
				tg[o] = v
				ftd.Tag = tg
			} else if string(c.Family) == "gt" {
				o, _ := strconv.Atoi(string(c.Qualifier[:]))
				v, _ := strconv.ParseInt(string(c.Value[:]), 10, 64)
				ntg[o] = v
				gtd.Tag = ntg
			}
		}
	}
	res["fan"] = ftd
	res["guest"] = gtd
	return
}

// RelationFansDay up relation 30 days analysis.
func (d *Dao) RelationFansDay(c context.Context, mid int64) (res datamodel.RelationFanHistoryData, err error) {
	var (
		result      *hrpc.Result
		ctx, cancel = context.WithTimeout(c, time.Duration(d.c.HBase.ReadTimeout))
		tableName   = HBaseUpRelationFansDay
		key         = hbaseMd5Key(mid)
		parser      = hbaseutil.Parser{
			ParseIntFunc: hbaseutil.StringToUint,
		}
	)
	defer cancel()
	if result, err = d.hbase.GetStr(ctx, tableName, key); err != nil {
		log.Error("RelationFansDay d.hbase.GetStr tableName(%s)|mid(%d)|key(%v)|error(%v)", tableName, mid, key, err)
		err = ecode.CreativeDataErr
		return
	}
	if result == nil || len(result.Cells) == 0 {
		log.Warn("RelationFansDay no data tableName(%s)|mid(%d)|key(%v)", tableName, mid, key)
		return
	}

	if err = parser.Parse(result.Cells, &res); err != nil {
		log.Error("parse cell fail, err=%v", err)
		return
	}
	/*
		families := make(map[string]string, 2)
		families["a"] = "" //某日新增关注量
		families["u"] = "" //某日取关量
		fd := make(map[string]int)
		nfd := make(map[string]int)
		for _, c := range result.Cells {
			if c == nil {
				continue
			}
			if _, ok := families[string(c.Family)]; ok {
				if string(c.Family) == "a" {
					k := string(c.Qualifier[:])
					v, _ := strconv.Atoi(string(c.Value[:]))
					fd[k] = v
				} else if string(c.Family) == "u" {
					k := string(c.Qualifier[:])
					v, _ := strconv.Atoi(string(c.Value[:]))
					nfd[k] = v
				}
			}
		}
		res["follow"] = fd
		res["unfollow"] = nfd
	*/
	return
}

// RelationFansHistory up relation history.
func (d *Dao) RelationFansHistory(c context.Context, mid int64, month string) (res datamodel.RelationFanHistoryData, err error) {
	var (
		result      *hrpc.Result
		ctx, cancel = context.WithTimeout(c, time.Duration(d.c.HBase.ReadTimeout))
		tableName   = HBaseUpRelationFansHistory
		key         = string(hbaseMd5Key(mid)) + "_" + month
		parser      = hbaseutil.Parser{ParseIntFunc: hbaseutil.StringToUint}
	)
	defer cancel()
	if result, err = d.hbase.GetStr(ctx, tableName, key); err != nil {
		log.Error("RelationFansHistory d.hbase.GetStr tableName(%s)|mid(%d)|key(%v)|error(%v)", tableName, mid, key, err)
		err = ecode.CreativeDataErr
		return
	}
	if result == nil || len(result.Cells) == 0 {
		log.Warn("RelationFansHistory no data tableName(%s)|mid(%d)|key(%v)", tableName, mid, key)
		return
	}

	if err = parser.Parse(result.Cells, &res); err != nil {
		log.Error("parse cell fail, err=%v", err)
		return
	}

	return
}

// RelationFansMonth up relation 400 days analysis.
func (d *Dao) RelationFansMonth(c context.Context, mid int64) (res datamodel.RelationFanHistoryData, err error) {
	var (
		result      *hrpc.Result
		ctx, cancel = context.WithTimeout(c, time.Duration(d.c.HBase.ReadTimeout))
		tableName   = HBaseUpRelationFansMonth
		key         = hbaseMd5Key(mid)
		parser      = hbaseutil.Parser{ParseIntFunc: hbaseutil.StringToUint}
	)
	defer cancel()
	if result, err = d.hbase.GetStr(ctx, tableName, key); err != nil {
		log.Error("RelationFansMonth d.hbase.GetStr tableName(%s)|mid(%d)|key(%v)|error(%v)", tableName, mid, key, err)
		err = ecode.CreativeDataErr
		return
	}
	if result == nil || len(result.Cells) == 0 {
		log.Warn("RelationFansMonth no data tableName(%s)|mid(%d)|key(%v)", tableName, mid, key)
		return
	}
	if err = parser.Parse(result.Cells, &res); err != nil {
		log.Error("parse cell fail, err=%v", err)
		return
	}
	return
}

// ViewerActionHour visitor action hour analysis.
func (d *Dao) ViewerActionHour(c context.Context, mid int64, dt string) (res map[string]*data.ViewerActionHour, err error) {
	var (
		result      *hrpc.Result
		ctx, cancel = context.WithTimeout(c, time.Duration(d.c.HBase.ReadTimeout))
		tableName   = HBaseUpViewerActionHour + dt
		key         = hbaseMd5Key(mid)
	)
	defer cancel()
	if result, err = d.hbase.GetStr(ctx, tableName, key); err != nil {
		log.Error("ViewerActionHour d.hbase.GetStr tableName(%s)|mid(%d)|key(%v)|error(%v)", tableName, mid, key, err)
		err = ecode.CreativeDataErr
		return
	}
	if result == nil || len(result.Cells) == 0 {
		log.Warn("ViewerActionHour no data tableName(%s)|mid(%d)|key(%v)", tableName, mid, key)
		return
	}
	families := make(map[string]string, 4)
	families["fp"] = "" //播放数
	families["fr"] = "" //评论数
	families["fd"] = "" //弹幕数
	families["fe"] = "" //充电数
	families["fs"] = "" //承包数
	res = make(map[string]*data.ViewerActionHour)
	view := make(map[int]int)
	reply := make(map[int]int)
	danmu := make(map[int]int)
	elec := make(map[int]int)
	con := make(map[int]int)
	fah := &data.ViewerActionHour{}
	gah := &data.ViewerActionHour{}
	for _, c := range result.Cells {
		if c == nil {
			continue
		}
		if _, ok := families[string(c.Family)]; ok {
			if string(c.Family) == "fp" {
				k, _ := strconv.Atoi(string(c.Qualifier[:]))
				v, _ := strconv.Atoi(string(c.Value[:]))
				view[k] = v
				fah.View = view
			} else if string(c.Family) == "gp" {
				k, _ := strconv.Atoi(string(c.Qualifier[:]))
				v, _ := strconv.Atoi(string(c.Value[:]))
				view[k] = v
				gah.View = view
			}
			if string(c.Family) == "fr" {
				k, _ := strconv.Atoi(string(c.Qualifier[:]))
				v, _ := strconv.Atoi(string(c.Value[:]))
				reply[k] = v
				fah.Reply = reply
			} else if string(c.Family) == "gr" {
				k, _ := strconv.Atoi(string(c.Qualifier[:]))
				v, _ := strconv.Atoi(string(c.Value[:]))
				reply[k] = v
				gah.Reply = reply
			}
			if string(c.Family) == "fd" {
				k, _ := strconv.Atoi(string(c.Qualifier[:]))
				v, _ := strconv.Atoi(string(c.Value[:]))
				danmu[k] = v
				fah.Dm = danmu
			} else if string(c.Family) == "gd" {
				k, _ := strconv.Atoi(string(c.Qualifier[:]))
				v, _ := strconv.Atoi(string(c.Value[:]))
				danmu[k] = v
				gah.Dm = danmu
			}
			if string(c.Family) == "fe" {
				k, _ := strconv.Atoi(string(c.Qualifier[:]))
				v, _ := strconv.Atoi(string(c.Value[:]))
				elec[k] = v
				fah.Elec = elec
			} else if string(c.Family) == "ge" {
				k, _ := strconv.Atoi(string(c.Qualifier[:]))
				v, _ := strconv.Atoi(string(c.Value[:]))
				elec[k] = v
				gah.Elec = elec
			}
			if string(c.Family) == "fs" {
				k, _ := strconv.Atoi(string(c.Qualifier[:]))
				v, _ := strconv.Atoi(string(c.Value[:]))
				con[k] = v
				fah.Contract = con
			} else if string(c.Family) == "gs" {
				k, _ := strconv.Atoi(string(c.Qualifier[:]))
				v, _ := strconv.Atoi(string(c.Value[:]))
				con[k] = v
				gah.Contract = con
			}
		}
	}
	res["fan"] = fah
	res["not_fan"] = gah
	return
}

// UpIncr for Play/Dm/Reply/Fav/Share/Elec/Coin incr.
func (d *Dao) UpIncr(c context.Context, mid int64, ty int8, now string) (res *data.UpDataIncrMeta, err error) {
	var (
		result      *hrpc.Result
		ctx, cancel = context.WithTimeout(c, d.hbaseTimeOut)
		tableName   string
		IncrKey     string
		key         = hbaseMd5Key(mid)
	)
	defer cancel()
	switch {
	case ty == data.Play:
		tableName = HBaseUpPlayInc + now
	case ty == data.Dm:
		tableName = HBaseUpDmInc + now
	case ty == data.Reply:
		tableName = HBaseUpReplyInc + now
	case ty == data.Share:
		tableName = HBaseUpShareInc + now
	case ty == data.Coin:
		tableName = HBaseUpCoinInc + now
	case ty == data.Fav:
		tableName = HBaseUpFavInc + now
	case ty == data.Elec:
		tableName = HBaseUpElecInc + now
	}
	IncrKey, _ = data.IncrTy(ty)
	if result, err = d.hbase.GetStr(ctx, tableName, key); err != nil {
		log.Error("UpIncr d.hbase.GetStr tableName(%s)|mid(%d)|key(%v)|error(%v)", tableName, mid, key, err)
		err = ecode.CreativeDataErr
		return
	}
	if result == nil || len(result.Cells) == 0 {
		log.Warn("UpIncr no data tableName(%s)|mid(%d)|key(%v)", tableName, mid, key)
		return
	}
	families := make(map[string]string, 4)
	families["u"] = ""  //单日播放增量 u:play ...
	families["av"] = "" //播放top稿件 av:1
	families["v"] = ""  //top稿件播放增量 v:1
	families["rk"] = "" //up主播放量排名 rk: [tid]
	aids := make(map[int]int64)
	incs := make(map[int]int)
	rk := make(map[int]int)
	var incr int
	for _, c := range result.Cells {
		if c == nil {
			continue
		}
		if _, ok := families[string(c.Family)]; ok {
			if string(c.Family) == "av" {
				k, _ := strconv.Atoi(string(c.Qualifier[:]))
				v, _ := strconv.ParseInt(string(c.Value[:]), 10, 64)
				aids[k] = v
			} else if string(c.Family) == "v" {
				k, _ := strconv.Atoi(string(c.Qualifier[:]))
				v, _ := strconv.Atoi(string(c.Value[:]))
				incs[k] = v
			} else if string(c.Family) == "rk" {
				k, _ := strconv.Atoi(string(c.Qualifier[:]))
				v, _ := strconv.Atoi(string(c.Value[:]))
				rk[k] = v
			} else if string(c.Family) == "u" {
				if bytes.Equal(c.Qualifier, []byte(IncrKey)) {
					v, _ := strconv.Atoi(string(c.Value[:]))
					if v < 0 {
						v = 0
					}
					incr = v
				}
			}

		}
	}
	res = &data.UpDataIncrMeta{}
	res.Incr = incr
	res.TopAIDList = aids
	res.TopIncrList = incs
	res.Rank = rk
	return
}

// ThirtyDayArchive for Play/Dm/Reply/Fav/Share/Elec/Coin for archive 30 days.
func (d *Dao) ThirtyDayArchive(c context.Context, mid int64, ty int8) (res []*data.ThirtyDay, err error) {
	var (
		result      *hrpc.Result
		ctx, cancel = context.WithTimeout(c, time.Duration(d.c.HBase.ReadTimeout))
		tableName   string
		key         = hbaseMd5Key(mid)
	)
	defer cancel()
	switch {
	case ty == data.Play:
		tableName = HBasePlayArc
	case ty == data.Dm:
		tableName = HBaseDmArc
	case ty == data.Reply:
		tableName = HBaseReplyArc
	case ty == data.Share:
		tableName = HBaseShareArc
	case ty == data.Coin:
		tableName = HBaseCoinArc
	case ty == data.Fav:
		tableName = HBaseFavArc
	case ty == data.Elec:
		tableName = HBaseElecArc
	case ty == data.Like:
		tableName = HBaseLikeArc
	}
	if result, err = d.hbase.GetStr(ctx, tableName, key); err != nil {
		log.Error("ThirtyDayArchive d.hbase.GetStr tableName(%s)|mid(%d)|key(%v)|error(%v)", tableName, mid, key, err)
		err = ecode.CreativeDataErr
		return
	}
	if result == nil || len(result.Cells) == 0 {
		log.Warn("ThirtyDay no data tableName(%s)|mid(%d)|key(%v)", tableName, mid, key)
		return
	}
	res = make([]*data.ThirtyDay, 0, len(result.Cells))
	for _, c := range result.Cells {
		if c == nil {
			continue
		}
		qual := string(c.Qualifier[:])
		val := string(c.Value[:])
		if string(c.Family) == "u" {
			t, v, err := parseKeyValue(qual, val)
			if err != nil {
				break
			}
			td := &data.ThirtyDay{}
			td.DateKey = t
			td.TotalIncr = v
			res = append(res, td)
		}
	}
	log.Info("ThirtyDayArchive mid(%d) type(%d) return data(%+v)", mid, ty, res)
	return
}

func parseKeyValue(k string, v string) (timestamp, value int64, err error) {
	tm, err := time.Parse("20060102", k)
	if err != nil {
		log.Error("time.Parse error(%v)", err)
		return
	}
	timestamp = tm.Unix()
	value, err = strconv.ParseInt(v, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt error(%v)", err)
	}
	return
}
