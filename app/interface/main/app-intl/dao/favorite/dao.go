package favorite

import (
	"context"
	"go-common/app/interface/main/app-intl/conf"
	favrpc "go-common/app/service/main/favorite/api/gorpc"
	fav "go-common/app/service/main/favorite/model"
	httpx "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

const (
	_isFavDef = "/x/internal/v2/fav/video/default"
	_isFav    = "/x/internal/v2/fav/video/favoured"
	_addFav   = "/x/internal/v2/fav/video/add"
)

// Dao is favorite dao
type Dao struct {
	client   *httpx.Client
	isFavDef string
	isFav    string
	addFav   string
	favRPC   *favrpc.Service
}

// New initial favorite dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		client:   httpx.NewClient(c.HTTPClient),
		isFavDef: c.Host.APICo + _isFavDef,
		isFav:    c.Host.APICo + _isFav,
		addFav:   c.Host.APICo + _addFav,
		favRPC:   favrpc.New2(c.FavoriteRPC),
	}
	return
}

// AddVideo add favorite
func (d *Dao) AddVideo(c context.Context, mid int64, fids []int64, aid int64, ak string) error {
	ip := metadata.String(c, metadata.RemoteIP)
	return d.favRPC.AddVideo(c, &fav.ArgAddVideo{Mid: mid, Fids: fids, Aid: aid, AccessKey: ak, RealIP: ip})
}
