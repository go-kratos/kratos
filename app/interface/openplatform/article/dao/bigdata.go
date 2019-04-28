package dao

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"net/url"
	"strconv"
	"time"

	"go-common/app/interface/main/creative/model/data"
	"go-common/app/interface/openplatform/article/model"
	"go-common/library/ecode"
	"go-common/library/log"

	"golang.org/x/net/context"
)

var (
	//HBaseArticleTable 文章作者概况
	HBaseArticleTable = "read_auth_stats_daily"
)

func hbaseMd5Key(aid int64) []byte {
	hasher := md5.New()
	hasher.Write([]byte(strconv.Itoa(int(aid))))
	return []byte(hex.EncodeToString(hasher.Sum(nil)))
}

// UpStat get the stat of article.
func (d *Dao) UpStat(c context.Context, mid int64) (stat model.UpStat, err error) {
	var (
		tableName = HBaseArticleTable
	)
	result, err := d.hbase.Get(c, []byte(tableName), hbaseMd5Key(mid))
	if err != nil {
		log.Error("bigdata: d.hbase.Get BackupTable(%s, %d) error(%+v)", tableName, mid, err)
		PromError("bigdata:hbase")
		err = ecode.CreativeDataErr
		return
	}
	if result == nil {
		return
	}
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
		case bytes.Equal(c.Qualifier, []byte("share1")):
			stat.Share = v
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
		case bytes.Equal(c.Qualifier, []byte("share0")):
			stat.PreShare = v
		}
	}
	stat.IncrView = stat.View - stat.PreView
	stat.IncrReply = stat.Reply - stat.PreReply
	stat.IncrCoin = stat.Coin - stat.PreCoin
	stat.IncrLike = stat.Like - stat.PreLike
	stat.IncrFav = stat.Fav - stat.PreFav
	stat.IncrShare = stat.Share - stat.PreShare
	d.AddCacheUpStatDaily(c, mid, &stat)
	return
}

// ThirtyDayArticle for Read/Reply/Like/Fav/Coin for article 30 days.
func (d *Dao) ThirtyDayArticle(c context.Context, mid int64) (res []*model.ThirtyDayArticle, err error) {
	var (
		tableName = "read_auth_stats" //文章30天数据
	)
	result, err := d.hbase.Get(c, []byte(tableName), hbaseMd5Key(mid))
	if err != nil {
		log.Error("bigdata: d.hbase.Get tableName(%s) mid(%d) error(%+v)", tableName, mid, err)
		PromError("bigdata:30天数据")
		err = ecode.CreativeDataErr
		return
	}
	if result == nil || len(result.Cells) == 0 {
		log.Warn("bigdata: ThirtyDay article no data (%s, %d)", tableName, mid)
		PromError("bigdata:30天数据")
		return
	}
	res = make([]*model.ThirtyDayArticle, 0, 5)
	vtds := make([]*data.ThirtyDay, 0, 30)
	ptds := make([]*data.ThirtyDay, 0, 30)
	ltds := make([]*data.ThirtyDay, 0, 30)
	ftds := make([]*data.ThirtyDay, 0, 30)
	ctds := make([]*data.ThirtyDay, 0, 30)
	view := &model.ThirtyDayArticle{Category: "view"}
	reply := &model.ThirtyDayArticle{Category: "reply"}
	like := &model.ThirtyDayArticle{Category: "like"}
	fav := &model.ThirtyDayArticle{Category: "fav"}
	coin := &model.ThirtyDayArticle{Category: "coin"}
	for _, c := range result.Cells {
		if c == nil {
			continue
		}
		family := string(c.Family)
		qual := string(c.Qualifier[:])
		val := string(c.Value[:])
		switch family {
		case "v": //"阅读量"
			t, v, err := parseKeyValue(qual, val)
			if err != nil {
				break
			}
			td := &data.ThirtyDay{}
			td.DateKey = t
			td.TotalIncr = v
			vtds = append(vtds, td)
			view.ThirtyDay = vtds
		case "p": //"评论量"
			t, v, err := parseKeyValue(qual, val)
			if err != nil {
				break
			}
			td := &data.ThirtyDay{}
			td.DateKey = t
			td.TotalIncr = v
			ptds = append(ptds, td)
			reply.Category = "reply"
			reply.ThirtyDay = ptds
		case "l": //"点赞量"
			t, v, err := parseKeyValue(qual, val)
			if err != nil {
				break
			}
			td := &data.ThirtyDay{}
			td.DateKey = t
			td.TotalIncr = v
			ltds = append(ltds, td)
			like.Category = "like"
			like.ThirtyDay = ltds
		case "f": //"收藏量"
			t, v, err := parseKeyValue(qual, val)
			if err != nil {
				break
			}
			td := &data.ThirtyDay{}
			td.DateKey = t
			td.TotalIncr = v
			ftds = append(ftds, td)
			fav.Category = "fav"
			fav.ThirtyDay = ftds
		case "c": //"投币量"
			t, v, err := parseKeyValue(qual, val)
			if err != nil {
				break
			}
			td := &data.ThirtyDay{}
			td.DateKey = t
			td.TotalIncr = v
			ctds = append(ctds, td)
			coin.Category = "coin"
			coin.ThirtyDay = ctds
		}
	}
	res = append(res, view)
	res = append(res, reply)
	res = append(res, like)
	res = append(res, fav)
	res = append(res, coin)
	return
}

func parseKeyValue(k string, v string) (timestamp, value int64, err error) {
	tm, err := time.Parse("20060102", k)
	if err != nil {
		log.Error("time.Parse error(%+v)", err)
		return
	}
	timestamp = tm.Unix()
	value, err = strconv.ParseInt(v, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt error(%+v)", err)
	}
	return
}

// SkyHorse sky horse
func (d *Dao) SkyHorse(c context.Context, mid int64, build int, buvid string, plat int8, ps int) (res *model.SkyHorseResp, err error) {
	if buvid == "" {
		err = ecode.NothingFound
		return
	}
	params := url.Values{}
	params.Set("cmd", "article")
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("buvid", buvid)
	params.Set("build", strconv.Itoa(build))
	params.Set("plat", strconv.FormatInt(int64(plat), 10))
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	params.Set("request_cnt", strconv.Itoa(ps))
	params.Set("from", "8")
	res = &model.SkyHorseResp{}
	err = d.httpClient.Get(c, d.c.Article.SkyHorseURL, "", params, &res)
	if err != nil {
		PromError("bigdata:天马接口")
		log.Error("bigdata: d.client.Get(%s) error(%+v)", d.c.Article.SkyHorseURL+"?"+params.Encode(), err)
		return
	}
	// -3: 数量不足
	if res.Code != 0 && res.Code != -3 {
		PromError("bigdata:天马接口")
		log.Error("bigdata: url(%s) res: %+v", d.c.Article.SkyHorseURL+"?"+params.Encode(), res)
		err = ecode.Int(res.Code)
		return
	}
	if len(res.Data) == 0 {
		PromError("bigdata:天马返回空")
		log.Warn("bigdata: url(%s) res: %+v", d.c.Article.SkyHorseURL+"?"+params.Encode(), res)
	}
	return
}
