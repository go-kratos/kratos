package tag

import (
	"context"

	"go-common/app/interface/main/app-intl/conf"
	tagmdl "go-common/app/interface/main/app-intl/model/tag"
	tag "go-common/app/interface/main/tag/model"
	tagrpc "go-common/app/interface/main/tag/rpc/client"
	"go-common/library/ecode"
	httpx "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

// Dao struct
type Dao struct {
	// conf
	// http client
	mInfo  string
	client *httpx.Client
	tagRPC *tagrpc.Service
}

// New a dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		client: httpx.NewClient(c.HTTPTag),
		mInfo:  c.Host.APICo + _mInfo,
		tagRPC: tagrpc.New2(c.TagRPC),
	}
	return
}

// InfoByIDs is.
func (d *Dao) InfoByIDs(c context.Context, mid int64, tids []int64) (tm map[int64]*tag.Tag, err error) {
	var ts []*tag.Tag
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &tag.ArgIDs{IDs: tids, Mid: mid, RealIP: ip}
	if ts, err = d.tagRPC.InfoByIDs(c, arg); err != nil {
		if err == ecode.NothingFound {
			err = nil
			return
		}
		err = errors.Wrapf(err, "%v", arg)
		return
	}
	tm = make(map[int64]*tag.Tag, len(ts))
	for _, t := range ts {
		tm[t.ID] = t
	}
	return
}

// ArcTags get tags data from api.
func (d *Dao) ArcTags(c context.Context, aid, mid int64) (tags []*tagmdl.Tag, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &tag.ArgAid{Aid: aid, Mid: mid, RealIP: ip}
	res, err := d.tagRPC.ArcTags(c, arg)
	if err != nil {
		return
	}
	if len(res) == 0 {
		return
	}
	tags = make([]*tagmdl.Tag, 0, len(res))
	for _, t := range res {
		tag := &tagmdl.Tag{ID: t.ID, Name: t.Name, Cover: t.Cover, Likes: t.Likes, Hates: t.Hates, Liked: t.Liked, Hated: t.Hated, Attribute: t.Attribute}
		tags = append(tags, tag)
	}
	return
}
