package music

import (
	"context"
	"fmt"
	"net/url"

	"go-common/app/interface/main/favorite/conf"
	"go-common/app/interface/main/favorite/model"
	"go-common/library/ecode"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"go-common/library/xstr"
)

const _music = "http://api.bilibili.co/x/internal/v1/audio/songs/batch"

// Dao defeine fav Dao
type Dao struct {
	httpClient *httpx.Client
}

// New return fav dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		httpClient: httpx.NewClient(c.HTTPClient),
	}
	return
}

// MusicMap return the music map data(all state).
func (d *Dao) MusicMap(c context.Context, musicIds []int64) (data map[int64]*model.Music, err error) {
	params := url.Values{}
	params.Set("level", "1")
	params.Set("ids", xstr.JoinInts(musicIds))
	res := new(model.MusicResult)
	ip := metadata.String(c, metadata.RemoteIP)
	if err = d.httpClient.Get(c, _music, ip, params, res); err != nil {
		log.Error("d.HTTPClient.Get(%s?%s) error(%v)", _music, params.Encode())
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("d.HTTPClient.Get(%s?%s) code:%d msg:%s", _music, params.Encode(), res.Code)
		err = fmt.Errorf("Get Music failed!code:=%v", res.Code)
		return
	}
	return res.Data, nil
}
