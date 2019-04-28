package app

import (
	"context"
	"fmt"
	"net/url"

	model "go-common/app/job/main/tv/model/pgc"
	"go-common/library/log"
)

const _type = "mp4"
const _maxBackup = 0
const _otype = "json"

// Playurl calls the api of playurl to get the url to play the video
func (d *Dao) Playurl(ctx context.Context, cid int) (playurl string, err error) {
	var (
		result = model.PlayurlResp{}
		params = url.Values{}
		api    = d.conf.Sync.PlayURL.API
	)
	params.Set("cid", fmt.Sprintf("%d", cid))
	params.Set("type", _type)                               // to get one piece
	params.Set("max_backup", fmt.Sprintf("%d", _maxBackup)) // no backup url needed
	params.Set("otype", _otype)                             // json format response
	params.Set("qn", d.conf.Sync.PlayURL.Qn)                // quality fix to 16
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
