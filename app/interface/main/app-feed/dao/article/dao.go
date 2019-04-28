package article

import (
	"context"

	"go-common/app/interface/main/app-feed/conf"
	article "go-common/app/interface/openplatform/article/model"
	artrpc "go-common/app/interface/openplatform/article/rpc/client"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

// Dao is archive dao.
type Dao struct {
	// rpc
	artRPC *artrpc.Service
}

// New new a archive dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// rpc
		artRPC: artrpc.New(c.ArticleRPC),
	}
	return
}

func (d *Dao) Articles(c context.Context, aids []int64) (ms map[int64]*article.Meta, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &article.ArgAids{Aids: aids, RealIP: ip}
	if ms, err = d.artRPC.ArticleMetas(c, arg); err != nil {
		err = errors.Wrapf(err, "%v", aids)
	}
	return
}
