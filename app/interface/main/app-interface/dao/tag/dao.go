package tag

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/interface/main/app-interface/conf"
	tagmdl "go-common/app/interface/main/app-interface/model/tag"
	tag "go-common/app/interface/main/tag/model"
	tagrpc "go-common/app/interface/main/tag/rpc/client"
	"go-common/library/ecode"
	httpx "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_mInfo = "/x/internal/tag/minfo"
)

// Dao is tag dao
type Dao struct {
	tagRPC *tagrpc.Service
	client *httpx.Client
	mInfo  string
}

// New initial tag dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		tagRPC: tagrpc.New2(c.TagRPC),
		client: httpx.NewClient(c.HTTPClient),
		mInfo:  c.Host.APICo + _mInfo,
	}
	return
}

// ArcTags get tags data from api.
func (d *Dao) ArcTags(c context.Context, aid, mid int64, ip string) (ts []*tagmdl.Tag, err error) {
	arg := &tag.ArgAid{Aid: aid, Mid: mid, RealIP: ip}
	tags, err := d.tagRPC.ArcTags(c, arg)
	if err != nil {
		err = errors.Wrapf(err, "%v", arg)
		return
	}
	if len(tags) == 0 {
		return
	}
	ts = make([]*tagmdl.Tag, 0, len(tags))
	for _, t := range tags {
		tag := &tagmdl.Tag{
			TagID:     t.ID,
			Name:      t.Name,
			Cover:     t.Cover,
			Likes:     t.Likes,
			Hates:     t.Hates,
			Liked:     t.Liked,
			Hated:     t.Hated,
			Attribute: t.Attribute,
		}
		ts = append(ts, tag)
	}
	return
}

// TagInfos get tag infos by tagIds
func (d *Dao) TagInfos(c context.Context, tags []int64, mid int64) (tagMyInfo []*tagmdl.Tag, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("tag_id", xstr.JoinInts(tags))
	params.Set("mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int           `json:"code"`
		Data []*tagmdl.Tag `json:"data"`
	}
	if err = d.client.Get(c, d.mInfo, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.mInfo+"?"+params.Encode())
		return
	}
	tagMyInfo = res.Data
	return
}
