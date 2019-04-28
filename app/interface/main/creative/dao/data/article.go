package data

import (
	"context"
	"strconv"
	"time"

	"go-common/app/interface/main/creative/model/data"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/tsuna/gohbase/hrpc"
)

var (
	trendTBLNameMap = map[byte]string{
		data.ArtView:    ArtViewTBL,
		data.ArtReply:   ArtReplyTBL,
		data.ArtShare:   ArtShareTBL,
		data.ArtCoin:    ArtCoinTBL,
		data.ArtFavTBL:  ArtFavTBL,
		data.ArtLikeTBL: ArtLikeTBL,
	}

	rankTBLNameMap = map[byte]string{
		data.ArtView:    ArtViewIncTBL,
		data.ArtReply:   ArtReplyIncTBL,
		data.ArtShare:   ArtShareIncTBL,
		data.ArtCoin:    ArtCoinIncTBL,
		data.ArtFavTBL:  ArtFavIncTBL,
		data.ArtLikeTBL: ArtLikeIncTBL,
	}
)

// ArtThirtyDay for article trend 30 days.
func (d *Dao) ArtThirtyDay(c context.Context, mid int64, ty byte) (res []*data.ArtTrend, err error) {
	var (
		result      *hrpc.Result
		ctx, cancel = context.WithTimeout(c, time.Duration(d.c.HBase.ReadTimeout))
		key         = hbaseMd5Key(mid)
	)
	defer cancel()

	tableName, ok := trendTBLNameMap[ty]
	if !ok {
		log.Error("ArtThirtyDay not exist type(%d)|mid(%d)", ty, mid)
		return
	}

	log.Info("ArtThirtyDay mid(%d)|tableName(%s)|family(u)|rowkey(%s)", mid, tableName, key)
	if result, err = d.hbase.GetStr(ctx, tableName, key); err != nil {
		log.Error("ArtThirtyDay d.hbase.GetStr tableName(%s)|mid(%d)|key(%v)|error(%v)", tableName, mid, key, err)
		err = ecode.CreativeDataErr
		return
	}
	if result == nil || len(result.Cells) == 0 {
		log.Warn("ArtThirtyDay no data tableName(%s)|mid(%d)|key(%v)", tableName, mid, key)
		return
	}

	res = make([]*data.ArtTrend, 0, len(result.Cells))
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
			td := &data.ArtTrend{}
			td.DateKey = t
			td.TotalIncr = v
			res = append(res, td)
		}
	}
	log.Info("ArtThirtyDay mid(%d)|tableName(%s)|family(u)|rowkey(%s)|res(%+v)", mid, tableName, key, res)
	return
}

// ArtRank for article rank
func (d *Dao) ArtRank(c context.Context, mid int64, ty byte, date string) (res *data.ArtRankMap, err error) {
	var (
		result      *hrpc.Result
		ctx, cancel = context.WithTimeout(c, time.Duration(d.c.HBase.ReadTimeout))
		key         = hbaseMd5Key(mid)
	)
	defer cancel()

	tableName, ok := rankTBLNameMap[ty]
	if !ok {
		log.Error("ArtRank not exist type(%d)|mid(%d)", ty, mid)
		return
	}
	tableName += date

	log.Info("ArtRank mid(%d)|tableName(%s)|family(rd,v)|rowkey(%s)", mid, tableName, key)
	if result, err = d.hbase.GetStr(ctx, tableName, key); err != nil {
		log.Error("ArtRank d.hbase.GetStr tableName(%s)|mid(%d)|key(%v)|error(%v)", tableName, mid, key, err)
		err = ecode.CreativeDataErr
		return
	}
	if result == nil || len(result.Cells) == 0 {
		log.Warn("ArtRank no data tableName(%s)|mid(%d)|key(%v)", tableName, mid, key)
		return
	}

	aids := make(map[int]int64)
	incrs := make(map[int]int)
	for _, c := range result.Cells {
		if c == nil {
			continue
		}
		if string(c.Family) == "rd" {
			k, _ := strconv.Atoi(string(c.Qualifier[:]))
			v, _ := strconv.ParseInt(string(c.Value[:]), 10, 64)
			aids[k] = v
		} else if string(c.Family) == "v" {
			k, _ := strconv.Atoi(string(c.Qualifier[:]))
			v, _ := strconv.Atoi(string(c.Value[:]))
			incrs[k] = v
		}
	}
	res = &data.ArtRankMap{}
	res.AIDs = aids
	res.Incrs = incrs
	log.Info("ArtRank mid(%d)|tableName(%s)|family(rd,v)|rowkey(%s)|res(%+v)", mid, tableName, key, res)
	return
}

// 专栏阅读来源分析 rowkey mid倒置补（10位）+ yyyyMMdd
func readSourceKey(id int64) string {
	idStr := strconv.FormatInt(id, 10)
	s := reverseString(idStr)
	s = s + time.Now().AddDate(0, 0, -1).Add(-12*time.Hour).Format("20060102")
	return s
}

// ReadAnalysis for article read source.
func (d *Dao) ReadAnalysis(c context.Context, mid int64) (res *data.ArtRead, err error) {
	var (
		result      *hrpc.Result
		ctx, cancel = context.WithTimeout(c, time.Duration(d.c.HBase.ReadTimeout))
		tableName   = ArtReadAnalysisTBL
		key         = readSourceKey(mid)
	)
	defer cancel()

	log.Info("ReadAnalysis mid(%d)|tableName(%s)|family(f)|rowkey(%s)", mid, tableName, key)
	if result, err = d.hbase.GetStr(ctx, tableName, key); err != nil {
		log.Error("ReadAnalysis d.hbase.GetStr tableName(%s)|mid(%d)|key(%v)|error(%v)", tableName, mid, key, err)
		err = ecode.CreativeDataErr
		return
	}
	if result == nil || len(result.Cells) == 0 {
		log.Warn("ReadAnalysis no data tableName(%s)|mid(%d)|key(%v)", tableName, mid, key)
		return
	}

	var cells data.ArtRead
	err = parser.Parse(result.Cells, &cells)
	if err != nil {
		log.Error("ReadAnalysis parser.Parse tableName(%s)|mid(%d)|rowkey(%+v)|error(%v)", tableName, mid, key, err)
		err = ecode.CreativeDataErr
		return
	}
	res = &cells
	log.Info("ReadAnalysis mid(%d)|tableName(%s)|family(f)|rowkey(%s)|res(%+v)", mid, tableName, key, res)
	return
}
