package archive

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"go-common/app/interface/main/creative/dao/tool"
	"go-common/app/interface/main/videoup/conf"
	"go-common/app/interface/main/videoup/model/archive"
	pordermdl "go-common/app/interface/main/videoup/model/porder"
	upapi "go-common/app/service/main/up/api/v1"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

const (
	_viewURL     = "/videoup/view"
	_addURL      = "/videoup/add"
	_editURL     = "/videoup/edit"
	_tagUpURL    = "/videoup/tag/up"
	_applyStaffs = "/videoup/staff/archive/applys"
	// StaffWhiteGroupID const
	StaffWhiteGroupID = int64(24)
)

// View get archive and videos.
func (d *Dao) View(c context.Context, aid int64, ip string) (a *archive.Archive, vs []*archive.Video, err error) {
	params := url.Values{}
	params.Set("aid", strconv.FormatInt(aid, 10))
	var res struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			Archive *archive.Archive `json:"archive"`
			Videos  []*archive.Video `json:"videos"`
		} `json:"data"`
	}
	if err = d.httpR.Get(c, d.viewURI, ip, params, &res); err != nil {
		log.Error("videoup view archive error(%v) | viewUri(%s) aid(%d) ip(%s) params(%v)", err, d.viewURI+"?"+params.Encode(), aid, ip, params)
		err = ecode.CreativeArchiveAPIErr
		return
	}
	if res.Code != 0 {
		err = ecode.Error(ecode.Int(res.Code), res.Message)
		log.Error("videoup view archive res code nq zero, res.Code(%d) | viewUri(%s) aid(%d) ip(%s) params(%v) res(%v)", res.Code, d.viewURI+"?"+params.Encode(), aid, ip, params, res)
		return
	}
	a = res.Data.Archive
	vs = res.Data.Videos
	return
}

// ApplyStaffs fn
func (d *Dao) ApplyStaffs(c context.Context, aid int64, ip string) (staffs []*archive.Staff, err error) {
	params := url.Values{}
	params.Set("aid", strconv.FormatInt(aid, 10))
	var res struct {
		Code    int                  `json:"code"`
		Message string               `json:"message"`
		Data    []*archive.StaffView `json:"data"`
	}
	if err = d.httpR.Get(c, d.applyStaffs, ip, params, &res); err != nil {
		log.Error("videoup view archive error(%v) | viewUri(%s) aid(%d) ip(%s) params(%v)", err, d.viewURI+"?"+params.Encode(), aid, ip, params)
		err = ecode.CreativeArchiveAPIErr
		return
	}
	if res.Code != 0 {
		log.Error("videoup view archive res code nq zero, res.Code(%d) | viewUri(%s) aid(%d) ip(%s) params(%v) res(%v)", res.Code, d.viewURI+"?"+params.Encode(), aid, ip, params, res)
		err = ecode.CreativeArchiveAPIErr
		return
	}
	for _, v := range res.Data {
		staff := &archive.Staff{
			Mid:   v.ApMID,
			Title: v.ApTitle,
		}
		staffs = append(staffs, staff)
	}
	return
}

// Add add archive and videos.
func (d *Dao) Add(c context.Context, ap *archive.ArcParam, ip string) (aid int64, err error) {
	params := url.Values{}
	params.Set("appkey", d.c.App.Key)
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	mh := md5.Sum([]byte(params.Encode() + d.c.App.Secret))
	params.Set("sign", hex.EncodeToString(mh[:]))
	var (
		uri = d.addURI + "?" + params.Encode()
	)
	bs, err := json.Marshal(ap)
	if err != nil {
		log.Error("json.Marshal error(%v) | ap(%v) ap.Mid(%d) ap.videos(%v)", err, ap, ap.Mid, ap.Videos)
		return
	}
	req, err := http.NewRequest("POST", uri, bytes.NewReader(bs))
	if err != nil {
		log.Error("http.NewRequest error(%v) | uri(%s) ap.Mid(%d)", err, uri, ap.Mid)
		return
	}
	req.Header.Set("X-BACKEND-BILI-REAL-IP", ip)
	var res struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			Aid int64 `json:"aid"`
		} `json:"data"`
	}
	if err = d.httpW.Do(c, req, &res); err != nil {
		log.Error("d.Add error(%v) | uri(%s) ap(%+v)", err, uri, ap)
		err = ecode.CreativeArchiveAPIErr
		return
	}
	if res.Code != 0 {
		err = ecode.Error(ecode.Int(res.Code), res.Message)
		log.Error("d.Add nq zero (%v)|(%v)|(%v)|(%v)|uri(%s),ap(%+v)", res.Code, res.Message, res.Data, err, uri, ap)
		return
	}
	log.Info("d.Add (%s)|res.Data.Aid(%d) ip(%s) ", string(bs), res.Data.Aid, ip)
	aid = res.Data.Aid
	return
}

// Edit edit archive and videos.
func (d *Dao) Edit(c context.Context, ap *archive.ArcParam, ip string) (err error) {
	params := url.Values{}
	params.Set("appkey", conf.Conf.App.Key)
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	mh := md5.Sum([]byte(params.Encode() + d.c.App.Secret))
	params.Set("sign", hex.EncodeToString(mh[:]))
	// uri
	var (
		uri = d.editURI + "?" + params.Encode()
	)
	// new request
	bs, err := json.Marshal(ap)
	if err != nil {
		log.Error("json.Marshal ap error (%v) | ap(%v)", err, ap)
		return
	}
	req, err := http.NewRequest("POST", uri, bytes.NewReader(bs))
	if err != nil {
		log.Error("http.NewRequest error(%v) | uri(%s) ap(%v)", err, uri, ap)
		return
	}
	req.Header.Set("X-BACKEND-BILI-REAL-IP", ip)
	var res struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err = d.httpW.Do(c, req, &res); err != nil {
		log.Error("d.Edit error(%v) | uri(%s) ap(%+v)", err, uri, ap)
		err = ecode.CreativeArchiveAPIErr
		return
	}
	if res.Code != 0 {
		err = ecode.Error(ecode.Int(res.Code), res.Message)
		log.Error("d.Add nq zero (%v)|(%v)|(%v)|uri(%s),ap(%+v)", res.Code, res.Message, err, uri, ap)
		return
	}
	log.Info("d.Edit(%s) | ip(%s)", string(bs), ip)
	return
}

// DescFormat fn
func (d *Dao) DescFormat(c context.Context) (descFormats map[int]*archive.DescFormat, err error) {
	var res struct {
		Code    int                   `json:"code"`
		Message string                `json:"message"`
		Data    []*archive.DescFormat `json:"data"`
	}
	descFormats = make(map[int]*archive.DescFormat)
	if err = d.httpR.Get(c, d.descFormatURI, "", nil, &res); err != nil {
		log.Error("videoup descFormat error(%v) | descFormatURI(%s) err(%v)", err, d.descFormatURI, err)
		err = ecode.CreativeArchiveAPIErr
		return
	}
	if res.Code != 0 {
		err = ecode.Error(ecode.Int(res.Code), res.Message)
		log.Error("videoup descFormat res.Code(%d) | descFormatURI(%s) res(%v) err(%v)", res.Code, d.descFormatURI, res, err)
		return
	}
	for _, v := range res.Data {
		descFormats[v.ID] = v
	}
	return
}

// TagUp fn
func (d *Dao) TagUp(c context.Context, aid int64, tag, ip string) (err error) {
	params := url.Values{}
	params.Set("aid", strconv.Itoa(int(aid)))
	params.Set("tag", tag)
	var res struct {
		Code int `json:"code"`
	}
	if err = d.httpW.Post(c, d.tagUpURI, ip, params, &res); err != nil {
		log.Error("Post(%s,%s,%s) err(%v)", d.tagUpURI, ip, params.Encode(), err)
		err = ecode.CreativeArchiveAPIErr
		return
	}
	if res.Code != 0 {
		log.Error("Code(%s,%s,%s) err(%v), code(%d)", d.tagUpURI, ip, params.Encode(), err, res.Code)
		err = ecode.CreativeArchiveAPIErr
		return
	}
	return
}

// PorderCfgList fn
func (d *Dao) PorderCfgList(c context.Context) (cfgs map[int64]*pordermdl.Config, err error) {
	cfgs = make(map[int64]*pordermdl.Config)
	var res struct {
		Code int                 `json:"code"`
		Data []*pordermdl.Config `json:"data"`
	}
	if err = d.httpR.Get(c, d.porderConfigURL, "", nil, &res); err != nil {
		log.Error("archive.porderConfigURL url(%s) error(%v)", d.porderConfigURL, err)
		err = ecode.CreativeArchiveAPIErr
		return
	}
	if res.Code != 0 {
		log.Error("archive.porderConfigURL url(%s) res(%v)", d.porderConfigURL, res)
		err = ecode.CreativeArchiveAPIErr
		return
	}
	for _, cfg := range res.Data {
		cfgs[cfg.ID] = cfg
	}
	return
}

// GameList fn
func (d *Dao) GameList(c context.Context) (gameMap map[int64]*pordermdl.Game, err error) {
	gameMap = make(map[int64]*pordermdl.Game)
	params := url.Values{}
	params.Set("appkey", conf.Conf.Game.App.Key)
	params.Set("appsecret", conf.Conf.Game.App.Secret)
	params.Set("ts", strconv.FormatInt(time.Now().UnixNano()/1000000, 10))
	var res struct {
		Code int               `json:"code"`
		Data []*pordermdl.Game `json:"data"`
	}
	var (
		query, _ = tool.Sign(params)
		url      string
	)
	url = d.gameListURL + "?" + query
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error("http.NewRequest(%s) error(%v), ", url, err)
		err = ecode.CreativeGameOpenAPIErr
		return
	}
	if err = d.httpR.Do(c, req, &res); err != nil {
		log.Error("d.httpR.Do(%s) error(%v);", url, err)
		err = ecode.CreativeGameOpenAPIErr
		return
	}
	log.Info("GameList url(%+v)|gameLen(%+v)", url, len(res.Data))
	if res.Code != 0 {
		log.Error("GameList api url(%s) res(%v);, code(%d)", url, res, res.Code)
		err = ecode.CreativeGameOpenAPIErr
		return
	}
	for _, data := range res.Data {
		gameMap[data.GameBaseID] = data
	}
	return
}

//StaffUps 联合投稿白名单
func (d *Dao) StaffUps(c context.Context) (ups map[int64]int64, err error) {
	return d.UpSpecial(c, StaffWhiteGroupID)
}

// StaffTypeConfig 获取联合投稿分区配置
func (d *Dao) StaffTypeConfig(c context.Context) (isGary bool, typeConf map[int16]*archive.StaffTypeConf, err error) {
	typeConf = make(map[int16]*archive.StaffTypeConf)
	params := url.Values{}
	var res struct {
		Code int `json:"code"`
		Data struct {
			IsGary   bool                     `json:"is_gary"`
			TypeList []*archive.StaffTypeConf `json:"typelist"`
		} `json:"data"`
	}
	if err = d.httpR.Get(c, d.staffConfigURI, "", params, &res); err != nil {
		log.Error("StaffTypeConfig error(%v) | staffConfigURI(%s)", err, d.staffConfigURI)
		err = ecode.CreativeArchiveAPIErr
		return
	}
	log.Info("StaffTypeConfig url(%+v)|res(%+v)", d.staffConfigURI, res)
	if res.Code != 0 || res.Data.TypeList == nil {
		log.Error("StaffTypeConfig api url(%s) res(%+v);, code(%d)", d.staffConfigURI, res, res.Code)
		err = ecode.CreativeArchiveAPIErr
		return
	}
	isGary = res.Data.IsGary
	for _, v := range res.Data.TypeList {
		typeConf[v.TypeID] = v
	}
	return
}

// UpSpecial 获取UP主的特殊用户组
func (d *Dao) UpSpecial(c context.Context, gpid int64) (ups map[int64]int64, err error) {
	var (
		res  *upapi.UpGroupMidsReply
		page int
		g    errgroup.Group
		l    sync.RWMutex
	)
	if res, err = d.UpClient.UpGroupMids(c, &upapi.UpGroupMidsReq{
		GroupID: gpid,
		Pn:      1,
		Ps:      1,
	}); err != nil {
		log.Error("UpSpecial d.UpSpecial gpid(%d)|error(%v)", gpid, err)
		return
	}
	log.Warn("UpSpecial get total: gpid(%d)|total(%d)", gpid, res.Total)
	if res.Total <= 0 {
		return
	}
	ups = make(map[int64]int64, res.Total)
	ps := int(10000)
	pageNum := res.Total / ps
	if res.Total%ps != 0 {
		pageNum++
	}
	for page = 1; page <= pageNum; page++ {
		tmpPage := page
		g.Go(func() (err error) {
			resgg, err := d.UpClient.UpGroupMids(c, &upapi.UpGroupMidsReq{
				GroupID: gpid,
				Pn:      tmpPage,
				Ps:      ps,
			})
			if err != nil {
				log.Error("d.UpGroupMids gg (%d,%d,%d) error(%v) ", gpid, tmpPage, ps, err)
				err = nil
				return
			}
			for _, mid := range resgg.Mids {
				l.Lock()
				ups[mid] = mid
				l.Unlock()
			}
			return
		})
	}
	g.Wait()
	log.Warn("UpSpecial get result: gpid,total,midslen,upslens (%d)|(%d)|(%d)", gpid, res.Total, len(ups))
	return
}
