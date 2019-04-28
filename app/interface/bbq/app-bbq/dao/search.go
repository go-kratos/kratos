package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"go-common/app/interface/bbq/app-bbq/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"net/url"
	"strconv"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

const (
	httpProto      = "http://"
	searchBaseURI  = "/bbq/search"
	sugBaseURI     = "/main/suggest"
	getSVIDbyRelID = "select id,svid from video where id in (%s)"
)

// SearchBBQ 搜索视频
func (d *Dao) SearchBBQ(c context.Context, sreq *model.SearchBaseReq) (ret *model.RawSearchRes, err error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	ret = new(model.RawSearchRes)
	path := httpProto + d.c.Search.Host + searchBaseURI
	params := url.Values{}
	d.preSetSearchParam(c, &params)
	params.Set("keyword", sreq.KeyWord)
	params.Set("page", strconv.FormatInt(sreq.Page, 10))
	params.Set("pagesize", strconv.FormatInt(sreq.PageSize, 10))
	params.Set("highlight", strconv.FormatInt(sreq.Highlight, 10))
	params.Set("search_type", sreq.Type)
	log.Infov(c, log.KV("log", fmt.Sprintf("search url(%s)", path+"?"+params.Encode())))
	req, err := d.httpClient.NewRequest("GET", path, params.Get("ip"), params)
	if err != nil {
		log.Error("search url(%s) error(%v)", path+"?"+params.Encode(), err)
		return
	}
	if err = d.httpClient.Do(c, req, &ret); err != nil {
		log.Errorv(c, log.KV("log", fmt.Sprintf("search url(%s) error(%v)", path+"?"+params.Encode(), err)))
		return
	}
	if ret.Code != 0 {
		err = ecode.Int(ret.Code)
		log.Errorv(c, log.KV("log", fmt.Sprintf("search url(%s) error(%v)", path+"?"+params.Encode(), err)))
		return
	}
	_str, _ := json.Marshal(ret)
	log.Infov(c,
		log.KV("log", fmt.Sprintf("Search req[%s] ret[%s]", path+"?"+params.Encode(), _str)))
	return
}

// SugBBQ 搜索视频
func (d *Dao) SugBBQ(c context.Context, sreq *model.SugBaseReq) (ret json.RawMessage, err error) {
	path := httpProto + d.c.Search.Host + sugBaseURI
	params := url.Values{}
	d.preSetSearchParam(c, &params)
	params.Set("term", sreq.Term)
	params.Set("suggest_type", sreq.SuggestType)
	params.Set("highlight", strconv.FormatInt(sreq.Highlight, 10))
	params.Set("main_ver", sreq.MainVer)
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	log.Infov(c, log.KV("log", fmt.Sprintf("sug url(%s)", path+"?"+params.Encode())))
	req, err := d.httpClient.NewRequest("GET", path, params.Get("ip"), params)
	if err != nil {
		log.Error("sug url(%s) error(%v)", path+"?"+params.Encode(), err)
		return
	}
	var res struct {
		Code   int    `json:"code"`
		Stoken string `json:"stoken"`
		Res    struct {
			Tag json.RawMessage `json:"tag"`
		} `json:"result"`
	}
	if err = d.httpClient.Do(c, req, &res); err != nil {
		log.Errorv(c, log.KV("log", fmt.Sprintf("sug url(%s) error(%v)", path+"?"+params.Encode(), err)))
		return
	}
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Errorv(c, log.KV("log", fmt.Sprintf("sug url(%s) error(%v)", path+"?"+params.Encode(), err)))
		return
	}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	_str, _ := json.Marshal(ret)
	log.Infov(c, log.KV("log", fmt.Sprintf("Sug req[%s] ret[%s]", path+"?"+params.Encode(), _str)))
	ret = res.Res.Tag
	return
}

func (d *Dao) preSetSearchParam(c context.Context, params *url.Values) {
	device := c.Value("device")
	if device != nil {
		dev := device.(*bm.Device)
		if dev.RawPlatform != "" {
			params.Set("platform", dev.RawPlatform)
		}
		if dev.RawMobiApp != "" {
			params.Set("mobi_app", dev.RawMobiApp)
		}
		if dev.Device != "" {
			params.Set("device", dev.Device)
		}
		if dev.Build > 0 {

			params.Set("device", strconv.FormatInt(dev.Build, 10))
		}
	}
	ip := metadata.String(c, metadata.RemoteIP)
	if ip != "" {
		params.Set("clientip", ip)
	}
}

// ConvID2SVID 转换搜索相对id到svid
func (d *Dao) ConvID2SVID(c context.Context, ids []int64) (res map[int64]int64, err error) {
	var idList []string
	res = make(map[int64]int64)
	if len(ids) == 0 {
		return
	}
	for _, id := range ids {
		idList = append(idList, strconv.FormatInt(id, 10))
	}
	if len(idList) == 0 {
		log.Warn("empty query list relID [%v]", ids)
		return
	}
	idStr := strings.Join(idList, ",")
	sql := fmt.Sprintf(getSVIDbyRelID, idStr)
	rows, err := d.db.Query(c, sql)
	if err != nil {
		log.Warn("Query Err [%v]", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var svid int64
		var id int64
		err = rows.Scan(&id, &svid)
		if err != nil {
			log.Warn("Scan Err [%v]", err)
			return
		}
		res[id] = svid
	}
	return
}

// ParseRel2ID 转换搜索相对id到自增id
func (d *Dao) ParseRel2ID(relID []int32) (idList []int64) {
	for _, id := range relID {
		idList = append(idList, int64(id/100))
	}
	return
}
