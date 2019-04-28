package tag

import (
	"context"

	"go-common/app/interface/main/app-view/conf"
	tagmdl "go-common/app/interface/main/app-view/model/tag"
	tag "go-common/app/interface/main/tag/model"
	tagrpc "go-common/app/interface/main/tag/rpc/client"
	"go-common/library/net/metadata"
)

// Dao is tag dao
type Dao struct {
	tagRPC *tagrpc.Service
}

// New initial tag dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		tagRPC: tagrpc.New2(c.TagRPC),
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
		tag := &tagmdl.Tag{TagID: t.ID, Name: t.Name, Cover: t.Cover, Likes: t.Likes, Hates: t.Hates, Liked: t.Liked, Hated: t.Hated, Attribute: t.Attribute}
		tags = append(tags, tag)
	}
	return
}
