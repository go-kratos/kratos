package favorite

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	"go-common/app/interface/main/app-interface/conf"
	"go-common/app/interface/main/app-interface/model/favorite"
	"go-common/library/ecode"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

const (
	_folder      = "/x/internal/v2/fav/folder"
	_folderVideo = "/x/internal/v2/fav/video"
)

// Dao is favorite dao
type Dao struct {
	client     *httpx.Client
	favor      string
	favorVideo string
}

// New initial favorite dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		client:     httpx.NewClient(c.HTTPClient),
		favor:      c.Host.APICo + _folder,
		favorVideo: c.Host.APICo + _folderVideo,
	}
	return
}

// Folders get favorite floders from api.
func (d *Dao) Folders(c context.Context, mid, vmid int64, mobiApp string, build int, mediaList bool) (fs []*favorite.Folder, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("vmid", strconv.FormatInt(vmid, 10))
	params.Set("mobi_app", mobiApp)
	// params.Set("build", strconv.Itoa(build))
	if mediaList {
		params.Set("medialist", "1")
	}
	var res struct {
		Code int                `json:"code"`
		Data []*favorite.Folder `json:"data"`
	}
	if err = d.client.Get(c, d.favor, ip, params, &res); err != nil {
		return
	}
	b, _ := json.Marshal(&res)
	log.Info("Folders url(%s) response(%s)", d.favor+"?"+params.Encode(), b)
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.favor+"?"+params.Encode())
		return
	}
	fs = res.Data
	return
}

// FolderVideo get favorite floders from UGC api.
func (d *Dao) FolderVideo(c context.Context, accessKey, actionKey, device, mobiApp, platform, keyword, order string, build, tid, pn, ps int, mid, fid, vmid int64) (fav *favorite.Video, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("access_key", accessKey)
	params.Set("actionKey", actionKey)
	params.Set("build", strconv.Itoa(build))
	params.Set("device", device)
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("fid", strconv.FormatInt(fid, 10))
	params.Set("tid", strconv.Itoa(tid))
	params.Set("keyword", keyword)
	params.Set("order", order)
	params.Set("pn", strconv.Itoa(pn))
	params.Set("ps", strconv.Itoa(ps))
	params.Set("mobi_app", mobiApp)
	params.Set("platform", platform)
	params.Set("vmid", strconv.FormatInt(vmid, 10))
	var res struct {
		Code int             `json:"code"`
		Data *favorite.Video `json:"data"`
	}
	if err = d.client.Get(c, d.favorVideo, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.favorVideo+"?"+params.Encode())
		return
	}
	fav = res.Data
	return
}
