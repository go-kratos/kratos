package audio

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/interface/main/app-interface/conf"
	"go-common/app/interface/main/app-interface/model/audio"
	"go-common/library/ecode"
	httpx "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_audios      = "/audio/music-service-c/songs/internal/upsongslist"
	_allAudio    = "/audio/music-service-c/songs/internal/uppersongs-preload"
	_audioDetail = "/audio/music-service-c/songs/internal/uppersongs-batch"
	_favAudio    = "/audio/music-service-c/collections"
	_upperCert   = "/audio/music-service-c/internal/upper-cert"
	_card        = "/x/internal/v1/audio/privilege/mcard"
	_fav         = "/x/internal/v1/audio/personal/coll"
)

type Dao struct {
	client      *httpx.Client
	audios      string
	allAudio    string
	audioDetail string
	favAudio    string
	upperCert   string
	card        string
	fav         string
}

func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		client:      httpx.NewClient(c.HTTPClient),
		audios:      c.Host.APICo + _audios,
		allAudio:    c.Host.APICo + _allAudio,
		audioDetail: c.Host.APICo + _audioDetail,
		favAudio:    c.Host.APICo + _favAudio,
		upperCert:   c.Host.APICo + _upperCert,
		card:        c.Host.APICo + _card,
		fav:         c.Host.APICo + _fav,
	}
	return
}

// Audios
func (d *Dao) Audios(c context.Context, mid int64, pn, ps int) (audios []*audio.Audio, total int, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("uid", strconv.FormatInt(mid, 10))
	params.Set("pageIndex", strconv.Itoa(pn))
	params.Set("pageSize", strconv.Itoa(ps))
	var res struct {
		Code int `json:"code"`
		Data *struct {
			Total int            `json:"total"`
			List  []*audio.Audio `json:"list"`
		} `json:"data"`
	}
	if err = d.client.Get(c, d.audios, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.audios+"?"+params.Encode())
		return
	}
	if res.Data != nil {
		total = res.Data.Total
		audios = res.Data.List
	}
	return
}

// AllAudio get 100 audio by ctime desc
func (d *Dao) AllAudio(c context.Context, vmid int64) (aus []*audio.Audio, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("uid", strconv.FormatInt(vmid, 10))
	var res struct {
		Code int            `json:"code"`
		Data []*audio.Audio `json:"data"`
	}
	if err = d.client.Get(c, d.allAudio, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.allAudio+"?"+params.Encode())
		return
	}
	aus = res.Data
	return
}

func (d *Dao) AudioDetail(c context.Context, ids []int64) (aum map[int64]*audio.Audio, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("ids", xstr.JoinInts(ids))
	var res struct {
		Code int                    `json:"code"`
		Data map[int64]*audio.Audio `json:"data"`
	}
	if err = d.client.Get(c, d.audioDetail, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.audioDetail+"?"+params.Encode())
		return
	}
	aum = res.Data
	return
}

func (d *Dao) FavAudio(c context.Context, accessKey string, mid int64, pn, ps int) (aus []*audio.FavAudio, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("access_key", accessKey)
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("sort", "-1")
	params.Set("page_index", strconv.Itoa(pn))
	params.Set("page_size", strconv.Itoa(ps))
	var res struct {
		Code int `json:"code"`
		Data *struct {
			List []*audio.FavAudio `json:"list"`
		} `json:"data"`
	}
	if err = d.client.Get(c, d.favAudio, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.favAudio+"?"+params.Encode())
		return
	}
	if res.Data != nil {
		aus = res.Data.List
	}
	return
}

func (d *Dao) UpperCert(c context.Context, uid int64) (cert *audio.UpperCert, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("uid", strconv.FormatInt(uid, 10))
	var res struct {
		Code int              `json:"code"`
		Data *audio.UpperCert `json:"data"`
	}
	if err = d.client.Get(c, d.upperCert, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.upperCert+"?"+params.Encode())
		return
	}
	cert = res.Data
	return
}

func (d *Dao) Card(c context.Context, mid ...int64) (cardm map[int64]*audio.Card, err error) {
	params := url.Values{}
	params.Set("mid", xstr.JoinInts(mid))
	var res struct {
		Code int                   `json:"code"`
		Data map[int64]*audio.Card `json:"data"`
	}
	if err = d.client.Get(c, d.card, "", params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.card+"?"+params.Encode())
		return
	}
	cardm = res.Data
	return
}

func (d *Dao) Fav(c context.Context, mid int64) (fav *audio.Fav, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int        `json:"code"`
		Data *audio.Fav `json:"data"`
	}
	if err = d.client.Get(c, d.fav, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.fav+"?"+params.Encode())
		return
	}
	fav = res.Data
	return
}
