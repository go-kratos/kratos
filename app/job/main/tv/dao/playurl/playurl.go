package playurl

import (
	"context"
	"fmt"
	"net/url"

	model "go-common/app/job/main/tv/model/pgc"
	"go-common/library/log"
)

const (
	_type      = "mp4"
	_maxBackup = 0
	_otype     = "json"
)

// Playurl calls the api of playurl to get the url to play the video
func (d *Dao) Playurl(ctx context.Context, cid int) (playurl string, hitDead bool, err error) {
	var (
		result    = model.PlayurlResp{}
		params    = url.Values{}
		api       = d.conf.Sync.PlayURL.API
		originURL string
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
		for _, v := range d.conf.Sync.PlayURL.Deadcodes {
			if result.Code == v { // hit dead code
				hitDead = true
				return
			}
		}
		err = fmt.Errorf("Resp Code:[%v], Message:[%v]", result.Code, result.Message)
		return
	}
	if len(result.Durl) < 1 { // result empty
		err = fmt.Errorf("Playurl Result is Empty! Resp (%v)", result)
		return
	}
	originURL = result.Durl[0].URL
	if playurl, err = d.hostChange(originURL); err != nil { // replace the host of the playurl
		log.Error("HostChange Origin: %s, Error: %v", originURL, err)
	}
	return
}

// replace the url's host by TV's dedicated host
func (d *Dao) hostChange(playurl string) (replaced string, err error) {
	var host = d.conf.Sync.PlayURL.PlayPath
	u, err := url.Parse(playurl)
	if err != nil {
		log.Error("hostChange ParseURL error (%v)", err)
		return
	}
	log.Info("[hostChange] for URL: %s, Original Host: %s, Now we change it to: %s", playurl, u.Host, host)
	u.Host = host   // replace the host
	u.RawQuery = "" // remove useless query
	replaced = u.String()
	return
}
