package tag

import (
	"context"
	"time"

	"go-common/app/interface/main/app-channel/conf"
	tag "go-common/app/interface/main/tag/model"
	tagrpc "go-common/app/interface/main/tag/rpc/client"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

// Dao tag
type Dao struct {
	c      *conf.Config
	client *bm.Client
	tagRPC *tagrpc.Service
}

//New tag
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:      c,
		client: bm.NewClient(conf.Conf.HTTPClient),
		tagRPC: tagrpc.New2(c.TagRPC),
	}
	return
}

// InfoByID by tag id
func (d *Dao) InfoByID(c context.Context, mid, tid int64) (t *tag.Tag, err error) {
	arg := &tag.ArgID{ID: tid, Mid: mid}
	t, err = d.tagRPC.InfoByID(c, arg)
	return
}

// InfoByName by tag name
func (d *Dao) InfoByName(c context.Context, mid int64, tname string) (t *tag.Tag, err error) {
	arg := &tag.ArgName{Name: tname, Mid: mid}
	t, err = d.tagRPC.InfoByName(c, arg)
	return
}

// ChannelDetail channel info by id or nane
func (d *Dao) ChannelDetail(c context.Context, mid, tid int64, tname string, isOversea int32) (t *tag.ChannelDetail, err error) {
	arg := &tag.ReqChannelDetail{Mid: mid, Tid: tid, TName: tname, From: isOversea}
	t, err = d.tagRPC.ChannelDetail(c, arg)
	return
}

// InfoByIDs by tids
func (d *Dao) InfoByIDs(c context.Context, mid int64, tids []int64) (tm map[int64]*tag.Tag, err error) {
	var ts []*tag.Tag
	arg := &tag.ArgIDs{IDs: tids, Mid: mid}
	if ts, err = d.tagRPC.InfoByIDs(c, arg); err != nil {
		if err == ecode.NothingFound {
			err = nil
			return
		}
		return
	}
	tm = make(map[int64]*tag.Tag, len(ts))
	for _, t := range ts {
		tm[t.ID] = t
	}
	return
}

// Resources channel resources aids
func (d *Dao) Resources(c context.Context, plat int8, id, mid int64, name, buvid string, build, requestCnt, loginEvent, displayID int) (res *tag.ChannelResource, err error) {
	arg := &tag.ArgChannelResource{
		Tid:        id,
		Mid:        mid,
		Plat:       int32(plat),
		Build:      int32(build),
		LoginEvent: int32(loginEvent),
		RequestCNT: int32(requestCnt),
		DisplayID:  int32(displayID),
		Type:       3,
		Name:       name,
		Buvid:      buvid,
		From:       0,
	}
	if res, err = d.tagRPC.ChannelResources(c, arg); err != nil {
		return
	}
	return
}

// SubscribeUpdate subscribe update
func (d *Dao) SubscribeUpdate(c context.Context, mid int64, tids string) (err error) {
	arg := &tag.ArgUpdateCustomSort{Tids: tids, Mid: mid, Type: 1}
	err = d.tagRPC.UpdateCustomSortTags(c, arg)
	return
}

// SubscribeAdd subscribe add
func (d *Dao) SubscribeAdd(c context.Context, mid, tagID int64, now time.Time) (err error) {
	arg := &tag.ArgAddSub{Tids: []int64{tagID}, Mid: mid, Now: now}
	err = d.tagRPC.AddSub(c, arg)
	return
}

// SubscribeCancel subscribe add
func (d *Dao) SubscribeCancel(c context.Context, mid, tagID int64, now time.Time) (err error) {
	arg := &tag.ArgCancelSub{Tid: tagID, Mid: mid, Now: now}
	err = d.tagRPC.CancelSub(c, arg)
	return
}

// Recommend func
func (d *Dao) Recommend(c context.Context, mid int64, isOversea int32) (rec []*tag.Channel, err error) {
	arg := &tag.ArgRecommandChannel{Mid: mid, From: isOversea}
	rec, err = d.tagRPC.RecommandChannel(c, arg)
	return
}

// ListByCategory 分类下的频道
func (d *Dao) ListByCategory(c context.Context, id, mid int64, isOversea int32) (list []*tag.Channel, err error) {
	arg := &tag.ArgChanneList{ID: id, Mid: mid, From: isOversea}
	list, err = d.tagRPC.ChanneList(c, arg)
	return
}

// Subscribe 已订阅频道
func (d *Dao) Subscribe(c context.Context, mid int64) (customSort *tag.CustomSortChannel, err error) {
	arg := &tag.ArgCustomSort{Type: 1, Mid: mid, Pn: 1, Ps: 400, Order: -1}
	customSort, err = d.tagRPC.CustomSortChannel(c, arg)
	return
}

//Discover 频道tab页的3个发现频道
func (d *Dao) Discover(c context.Context, mid int64, isOversea int32) (discover []*tag.Channel, err error) {
	arg := &tag.ArgDiscoverChanneList{Mid: mid, From: isOversea}
	discover, err = d.tagRPC.DiscoverChannel(c, arg)
	return
}

//Category channel category
func (d *Dao) Category(c context.Context, isOversea int32) (category []*tag.ChannelCategory, err error) {
	category, err = d.tagRPC.ChannelCategories(c, &tag.ArgChannelCategories{From: isOversea})
	return
}

//Square 频道广场页推荐频道+稿件
func (d *Dao) Square(c context.Context, mid int64, tagNum, oidNum, build int, loginEvent int32, plat int8, buvid string, isOversea int32) (square *tag.ChannelSquare, err error) {
	arg := &tag.ReqChannelSquare{
		Mid:        mid,
		TagNumber:  int32(tagNum),
		OidNumber:  int32(oidNum),
		Type:       3,
		Buvid:      buvid,
		Build:      int32(build),
		Plat:       int32(plat),
		LoginEvent: loginEvent,
		DisplayID:  1,
		From:       isOversea,
	}
	square, err = d.tagRPC.ChannelSquare(c, arg)
	return
}
