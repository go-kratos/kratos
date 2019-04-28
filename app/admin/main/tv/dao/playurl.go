package dao

import (
	"context"
	"fmt"
	"net/url"

	"go-common/app/admin/main/tv/model"
	"go-common/library/log"
)

const _type = "mp4"
const _maxBackup = 0
const _otype = "json"
const _qn = "16"

// Playurl Def.
func (d *Dao) Playurl(ctx context.Context, cid int) (playurl string, err error) {
	var (
		result = model.PlayurlResp{}
		params = url.Values{}
		api    = d.c.Cfg.PlayurlAPI
	)
	params.Set("cid", fmt.Sprintf("%d", cid))
	params.Set("type", _type)                               // to get one piece
	params.Set("max_backup", fmt.Sprintf("%d", _maxBackup)) // no backup url needed
	params.Set("otype", _otype)                             // json format response
	params.Set("qn", _qn)                                   // json format response
	if err = d.client.Get(ctx, api, "", params, &result); err != nil {
		log.Error("ClientGet error[%v]", err)
		return
	}
	if result.Code != 0 { // logic error
		err = fmt.Errorf("Resp Code:[%v], Message:[%v]", result.Code, result.Message)
		return
	}
	if len(result.Durl) < 1 { // result empty
		err = fmt.Errorf("Playurl Result is Empty! Resp (%v)", result)
		return
	}
	playurl = result.Durl[0].URL
	return
}

//UPlayurl ugc play url
func (d *Dao) UPlayurl(ctx context.Context, cid int) (playurl string, err error) {
	var (
		result = model.UPlayURLR{}
		params = url.Values{}
		api    = d.c.Cfg.UPlayurlAPI
	)
	params.Set("cid", fmt.Sprintf("%d", cid))
	params.Set("type", "mp4")
	params.Set("max_backup", fmt.Sprintf("%d", 0))
	params.Set("otype", "json")
	params.Set("qn", "16")
	params.Set("platform", "tvproj")
	if err = d.client.Get(ctx, api, "", params, &result); err != nil {
		log.Error("UPlayurl ClientGet error[%v]", err)
		return
	}
	if result.Code != 0 { // logic error
		err = fmt.Errorf("UPlayurl Resp Code:[%v], Message:[%v], Result:[%v]", result.Code, result.Message, result.Result)
		return
	}
	if len(result.Durl) < 1 { // result empty
		err = fmt.Errorf("UPlayurl Result is Empty! Resp (%v)", result)
		return
	}
	playurl = result.Durl[0].URL
	return
}
