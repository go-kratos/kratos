package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"go-common/app/interface/bbq/app-bbq/api/http/v1"
	"go-common/app/interface/bbq/app-bbq/model"
	"go-common/app/interface/bbq/app-bbq/model/grpc"
	rec "go-common/app/service/bbq/recsys/api/grpc/v1"
	video "go-common/app/service/bbq/video/api/grpc/v1"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"

	xgrpc "google.golang.org/grpc"
)

const (
	_defaultPlatform     = "html5"
	_playBcNum           = 1
	_queryList           = "select `avid`, `cid`, `svid`, `title`, `mid`, `content`, `pubtime`,`duration`,`tid`,`sub_tid`,`ctime`,`cover_url`,`cover_width`,`cover_height`,`state` from video where svid in (%s)"
	_queryRand           = "select `svid` from video where state in (4,5) order by mtime desc limit 400;"
	_queryStatisticsList = "select `svid`, `play`, `subtitles`, `like`, `share`, `report` from video_statistics where svid in (%s)"
	_querySvPlay         = "select `svid`,`path`,`resolution_retio`,`code_rate`,`video_code`,`file_size`,`duration` from %s where svid in (%s) and is_deleted = 0 order by code_rate desc"
)

// GetList 获取列表，按照db排序，按page返回，用于推荐的降级
func (d *Dao) GetList(c context.Context, pageSize int64) (result []int64, err error) {
	rows, err := d.db.Query(c, _queryRand)
	if err != nil {
		log.Error("Query(%s) error(%v)", _queryRand, err)
		return
	}
	defer rows.Close()

	var list []int64
	for rows.Next() {
		sv := new(model.SvInfo)
		if err = rows.Scan(&sv.SVID); err != nil {
			log.Error("row.Scan() error(%v)", err)
			return
		}
		if tag, ok := d.c.Tmap[strconv.FormatInt(sv.TID, 10)]; ok {
			sv.Tag = tag
		}
		list = append(list, sv.SVID)
	}

	size := len(list)
	if size > int(pageSize) {
		size = int(pageSize)
	}
	rand.Seed(time.Now().Unix())
	sort.Slice(list, func(i, j int) bool {
		return rand.Float32() > 0.5
	})
	result = list[:size]

	return
}

// AttentionRecList 关注页为空时的推荐
func (d *Dao) AttentionRecList(ctx context.Context, size int64, mid int64, buvid string) (svIDs []int64, err error) {
	svIDs, err = d.abstractRawRecList(ctx, d.recsysClient.UpsRecService, size, mid, buvid)
	log.V(1).Infow(ctx, "log", "get ups rec service", "method", "UpsRecService")
	return
}

// RawRecList 获取推荐列表
func (d *Dao) RawRecList(ctx context.Context, size int64, mid int64, buvid string) (svIDs []int64, err error) {
	svIDs, err = d.abstractRawRecList(ctx, d.recsysClient.RecService, size, mid, buvid)
	log.V(1).Infow(ctx, "log", "get ups rec service", "method", "RecService")
	return
}

type recsysFunc func(ctx context.Context, in *rec.RecsysRequest, opts ...xgrpc.CallOption) (*rec.RecsysResponse, error)

// abstractRawRecList 由于访问recsys的方法都是同样的请求&回包，同时过程也一样，因此
func (d *Dao) abstractRawRecList(ctx context.Context, f recsysFunc, size int64, mid int64, buvid string) (svIDs []int64, err error) {
	var (
		res *rec.RecsysResponse
	)

	req := &rec.RecsysRequest{
		MID:   mid,
		BUVID: buvid,
		Limit: int32(size),
	}
	if tmp, ok := ctx.(*bm.Context).Get("BBQBase"); ok && tmp != nil {
		switch tmp.(type) {
		case *v1.Base:
			base := tmp.(*v1.Base)
			req.App = base.App
			req.AppVersion = base.Version
		}
	}
	if tmp, ok := ctx.(*bm.Context).Get("QueryID"); ok && tmp != nil {
		switch tmp.(type) {
		case string:
			req.QueryID = tmp.(string)
		}
	}
	log.Info("Rec请求 params: [%v]", req)
	//res, err = d.recsysClient.RecService(ctx, req)
	res, err = f(ctx, req)
	debug := ctx.(*bm.Context).Request.Header.Get("debug")
	if err != nil {
		log.Errorv(ctx,
			log.KV("log", fmt.Sprintf("d.recsysClient.RecService err [%v]", err)),
		)
		// 降级（推荐服务已挂）
		svIDs = d.GetRandSvList(int(size))
		return
	} else if len(res.List) == 0 || debug == "1" {
		log.Warnv(ctx,
			log.KV("log", fmt.Sprintf("d.recsysClient.RecService return empty [%v]", res.List)),
		)

		// 降级（推荐接口返回空）
		svIDs = d.GetRandSvList(int(size))
		return
	} else {
		num := len(res.List)
		if int64(num) != size {
			log.Warnv(ctx,
				log.KV("log", fmt.Sprintf("d.recsysClient.RecService return num[%d] not match size[%d]", num, size)),
			)
		}
		for n, sv := range res.List {
			if int64(n) > size {
				break
			}
			svIDs = append(svIDs, sv.Svid)
		}
	}
	return
}

// GetVideoDetail 从数据库video中获取svid相应的信息
func (d *Dao) GetVideoDetail(ctx context.Context, svIDs []int64) (list []*model.SvInfo, retIDs []int64, err error) {
	list = make([]*model.SvInfo, 0)
	num := len(svIDs)
	if num == 0 {
		return
	}
	var IDsStr string
	for i, id := range svIDs {
		if i < num-1 {
			IDsStr += strconv.FormatInt(id, 10) + ","
		} else {
			IDsStr += strconv.FormatInt(id, 10)
		}
	}
	query := fmt.Sprintf(_queryList, IDsStr)
	rows, err := d.db.Query(ctx, query)
	if err != nil {
		log.Error("Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		sv := new(model.SvInfo)
		if err = rows.Scan(&sv.AVID, &sv.CID, &sv.SVID, &sv.Title, &sv.MID, &sv.Content, &sv.Pubtime, &sv.Duration, &sv.TID, &sv.SubTID, &sv.Ctime, &sv.CoverURL, &sv.CoverWidth, &sv.CoverHeight, &sv.State); err != nil {
			log.Error("row.Scan() error(%v)", err)
			return
		}
		if tag, ok := d.c.Tmap[strconv.FormatInt(sv.TID, 10)]; ok {
			sv.Tag = tag
		}
		retIDs = append(retIDs, sv.SVID)
		list = append(list, sv)
	}
	return
}

// RawVideos 从数据库批量获取视频信息
func (d *Dao) RawVideos(ctx context.Context, svIDs []int64) (res map[int64]*model.SvInfo, err error) {
	res = make(map[int64]*model.SvInfo)
	num := len(svIDs)
	if num == 0 {
		return
	}
	var IDsStr string
	for i, id := range svIDs {
		if i < num-1 {
			IDsStr += strconv.FormatInt(id, 10) + ","
		} else {
			IDsStr += strconv.FormatInt(id, 10)
		}
	}
	query := fmt.Sprintf(_queryList, IDsStr)
	rows, err := d.db.Query(ctx, query)
	if err != nil {
		log.Error("Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		sv := new(model.SvInfo)
		if err = rows.Scan(&sv.AVID, &sv.CID, &sv.SVID, &sv.Title, &sv.MID, &sv.Content, &sv.Pubtime, &sv.Duration, &sv.TID, &sv.SubTID, &sv.Ctime, &sv.CoverURL, &sv.CoverWidth, &sv.CoverHeight, &sv.State); err != nil {
			log.Error("row.Scan() error(%v)", err)
			return
		}
		res[sv.SVID] = sv
	}
	return
}

// RawVideoStatistic get video statistics
func (d *Dao) RawVideoStatistic(c context.Context, svids []int64) (res map[int64]*model.SvStInfo, err error) {
	const maxIDNum = 20
	var (
		idStr string
	)
	res = make(map[int64]*model.SvStInfo)
	if len(svids) > maxIDNum {
		svids = svids[:maxIDNum]
	}
	l := len(svids)
	for k, svid := range svids {
		if k < l-1 {
			idStr += strconv.FormatInt(svid, 10) + ","
		} else {
			idStr += strconv.FormatInt(svid, 10)
		}
		res[svid] = &model.SvStInfo{}
	}
	rows, err := d.db.Query(c, fmt.Sprintf(_queryStatisticsList, idStr))
	if err != nil {
		log.Error("query error(%s)", err.Error())
		return
	}
	defer rows.Close()
	for rows.Next() {
		ssv := new(model.SvStInfo)
		if err = rows.Scan(&ssv.SVID, &ssv.Play, &ssv.Subtitles, &ssv.Like, &ssv.Share, &ssv.Report); err != nil {
			log.Error("RawVideoStatistic rows.Scan() error(%v)", err)
			return
		}
		res[ssv.SVID] = ssv
	}
	cmtCount, _ := d.ReplyCounts(c, svids, model.DefaultCmType)
	for id, cmt := range cmtCount {
		if _, ok := res[id]; ok {
			res[id].Reply = cmt.Count
		}
	}
	return
}

// RawPlayURLs 批量获取cid playurl
func (d *Dao) RawPlayURLs(c context.Context, cids []int64, qn int64, plat string) (res map[int64]*v1.CVideo, err error) {
	res = make(map[int64]*v1.CVideo)
	var cs string
	// transfer cid array to string
	l := len(cids)
	for i := 0; i < l; i++ {
		if i != l-1 {
			cs += strconv.FormatInt(cids[i], 10) + ","
		} else {
			cs += strconv.FormatInt(cids[i], 10)
		}
	}
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("cid", cs)
	params.Set("qn", strconv.FormatInt(qn, 10))
	params.Set("platform", plat)
	var ret struct {
		Code int                              `json:"code"`
		Data map[string]map[string]*v1.CVideo `json:"data"`
	}
	err = d.httpClient.Get(c, d.c.URLs["bvc_batch"], ip, params, &ret)
	if err != nil || ret.Data["cids"] == nil {
		log.Error("http Get err %v", err)
		return
	}
	for id, v := range ret.Data["cids"] {
		var cid int64
		cid, err = strconv.ParseInt(id, 10, 64)
		if err != nil {
			log.Error("strconv.ParseInt err %v", err)
		}
		if _, ok := res[cid]; !ok {
			res[cid] = new(v1.CVideo)
		}
		res[cid] = v
	}
	return
}

// RelPlayURLs 相对地址批量获取playurl
func (d *Dao) RelPlayURLs(c context.Context, addrs []string) (res map[string]*grpc.VideoKeyItem, err error) {
	res = make(map[string]*grpc.VideoKeyItem)
	req := &grpc.RequestMsg{
		Keys:     addrs,
		Backup:   uint32(_playBcNum),
		Platform: _defaultPlatform,
		UIP:      metadata.String(c, metadata.RemoteIP),
	}
	_str, _ := json.Marshal(req)
	log.V(5).Infov(c, log.KV("log", fmt.Sprintf("bvc play req (%s)", string(_str))))
	r, err := d.bvcPlayClient.ProtobufPlayurl(c, req)
	_str, _ = json.Marshal(r)
	if err != nil {
		log.Error("bvc play err[%v] ret[%s]", err, string(_str))
		return
	}
	log.V(5).Infov(c, log.KV("log", fmt.Sprintf("bvc play ret (%s)", string(_str))))
	res = r.Data
	return
}

//RawSVBvcKey 批量获取playurl相对地址
func (d *Dao) RawSVBvcKey(c context.Context, svids []int64) (res map[int64][]*model.SVBvcKey, err error) {
	var (
		tb   map[string][]string
		rows *sql.Rows
	)
	res = make(map[int64][]*model.SVBvcKey)
	tb = make(map[string][]string)
	tName := "video_bvc_%02d"
	for _, v := range svids {
		if v <= 0 {
			continue
		}
		tbName := fmt.Sprintf(tName, v%100)
		tb[tbName] = append(tb[tbName], strconv.FormatInt(v, 10))
	}
	for k, v := range tb {
		query := fmt.Sprintf(_querySvPlay, k, strings.Join(v, ","))
		if rows, err = d.db.Query(c, query); err != nil {
			log.Errorv(c, log.KV("log", "RawSVBvcKey query sql"), log.KV("err", err))
			continue
		}
		for rows.Next() {
			tmp := model.SVBvcKey{}
			if err = rows.Scan(&tmp.SVID, &tmp.Path, &tmp.ResolutionRetio, &tmp.CodeRate, &tmp.VideoCode, &tmp.FileSize, &tmp.Duration); err != nil {
				log.Errorv(c, log.KV("log", "RawSVBvcKey scan"), log.KV("err", err))
				continue
			}
			res[tmp.SVID] = append(res[tmp.SVID], &tmp)
		}
	}
	return
}

// RelRecList 相关推荐列表
func (d *Dao) RelRecList(ctx context.Context, req *rec.RecsysRequest) (svIDs []int64, err error) {
	log.V(1).Infov(ctx, log.KV("log", fmt.Sprintf("RelatedRecService req [%+v]", req)))
	res, err := d.recsysClient.RelatedRecService(ctx, req)
	if err != nil {
		log.Errorv(ctx,
			log.KV("log", fmt.Sprintf("RelatedRecService err [%v]", err)),
		)
		return
	}
	num := len(res.List)
	if int32(num) != req.Limit {
		log.Errorv(ctx,
			log.KV("log", fmt.Sprintf("RelatedRecService ret num[%d] not match req size[%d]", num, req.Limit)),
		)
	}
	for n, sv := range res.List {
		if int32(n) > req.Limit {
			break
		}
		svIDs = append(svIDs, sv.Svid)
	}
	log.V(1).Infov(ctx, log.KV("log", fmt.Sprintf("RelatedRecService svid [%+v]", svIDs)))
	return
}

// SvDel 视频删除
func (d *Dao) SvDel(c context.Context, in *video.VideoDeleteRequest) (interface{}, error) {
	return d.videoClient.VideoDelete(c, in)
}
