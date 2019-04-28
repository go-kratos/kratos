package tag

import (
	"context"
	"net/url"
	"strconv"

	tagmdl "go-common/app/interface/main/app-interface/model/tag"
	"go-common/library/ecode"
	"go-common/library/net/metadata"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_mInfo = "/x/internal/tag/minfo"
)

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
