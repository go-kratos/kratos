package http

import (
	"encoding/json"
	"io/ioutil"

	"go-common/app/admin/main/videoup/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

// submit second_round submit.
func submit(c *bm.Context) {
	var (
		req = c.Request
		bs  []byte
		err error
	)
	if bs, err = ioutil.ReadAll(req.Body); err != nil {
		log.Error("ioutil.ReadAll() error(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	req.Body.Close()
	// params
	var ap = &archive.ArcParam{}
	if err = json.Unmarshal(bs, ap); err != nil {
		log.Error("http submit() json.Unmarshal(%s) error(%v)", string(bs), err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// TODO check data.
	if err = vdaSvc.Submit(c, ap); err != nil {
		log.Error("vdaSvc.Submit() error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

// batchSubmit batch submit by async.
func batchArchive(c *bm.Context) {
	var (
		req = c.Request
		bs  []byte
		err error
		aps []*archive.ArcParam
	)
	if bs, err = ioutil.ReadAll(req.Body); err != nil {
		log.Error("ioutil.ReadAll() error(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	req.Body.Close()
	// params
	if err = json.Unmarshal(bs, &aps); err != nil {
		log.Error("http submit() json.Unmarshal(%s) error(%v)", string(bs), err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if ok := vdaSvc.CheckArchive(aps); !ok {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// TODO check data.
	if err = vdaSvc.BatchArchive(c, aps, archive.ActionArchiveSubmit); err != nil {
		log.Error("vdaSvc.BatchSubmit() error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

// batchArchiveSecondRound batch submit by async.
func batchArchiveSecondRound(c *bm.Context) {
	var (
		req = c.Request
		bs  []byte
		err error
		aps []*archive.ArcParam
	)
	if bs, err = ioutil.ReadAll(req.Body); err != nil {
		log.Error("ioutil.ReadAll() error(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	req.Body.Close()
	// params
	if err = json.Unmarshal(bs, &aps); err != nil {
		log.Error("http submit() json.Unmarshal(%s) error(%v)", string(bs), err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if ok := vdaSvc.CheckArchive(aps); !ok {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// TODO check data.
	if err = vdaSvc.BatchArchive(c, aps, archive.ActionArchiveSecondRound); err != nil {
		log.Error("vdaSvc.BatchSubmit() error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

// batchAttrs batch attr by async.
func batchAttrs(c *bm.Context) {
	var (
		req = c.Request
		bs  []byte
		err error
		aps []*archive.ArcParam
	)
	if bs, err = ioutil.ReadAll(req.Body); err != nil {
		log.Error("ioutil.ReadAll() error(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	req.Body.Close()
	// params
	if err = json.Unmarshal(bs, &aps); err != nil {
		log.Error("http submit() json.Unmarshal(%s) error(%v)", string(bs), err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if ok := vdaSvc.CheckArchive(aps); !ok {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// TODO check data.
	if err = vdaSvc.BatchArchive(c, aps, archive.ActionArchiveAttr); err != nil {
		log.Error("vdaSvc.BatchSubmit() error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

// batchTypeIDs batch typeid by async.
func batchTypeIDs(c *bm.Context) {
	var (
		req = c.Request
		bs  []byte
		err error
		aps []*archive.ArcParam
	)
	if bs, err = ioutil.ReadAll(req.Body); err != nil {
		log.Error("ioutil.ReadAll() error(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	req.Body.Close()
	// params
	if err = json.Unmarshal(bs, &aps); err != nil {
		log.Error("http submit() json.Unmarshal(%s) error(%v)", string(bs), err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if ok := vdaSvc.CheckArchive(aps); !ok {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// TODO check data.
	if err = vdaSvc.BatchArchive(c, aps, archive.ActionArchiveTypeID); err != nil {
		log.Error("vdaSvc.BatchSubmit() error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

// // batchZlimit batch modify zlimit.
// func batchZlimit(c wctx.Context) {
// 	var (
// 		req = c.Request()
// 		res = c.Result()
// 		bs  []byte
// 		err error
// 	)
// 	if bs, err = ioutil.ReadAll(req.Body); err != nil {
// 		log.Error("ioutil.ReadAll() error(%v)", err)
// 		res["code"] = ecode.RequestErr
// 		return
// 	}
// 	req.Body.Close()
// 	// params
// 	var ap = &archive.ArcParam{}
// 	if err = json.Unmarshal(bs, ap); err != nil {
// 		log.Error("http submit() json.Unmarshal(%s) error(%v)", string(bs), err)
// 		res["code"] = ecode.RequestErr
// 		return
// 	}
// 	if len(ap.Aids) == 0 {
// 		res["code"] = ecode.RequestErr
// 		return
// 	}
// 	if err = vdaSvc.BatchZlimit(c, ap); err != nil {
// 		log.Error("vdaSvc.submit() error(%v)", err)
// 		res["code"] = err
// 		return
// 	}
// }

// upAuther batch modify zlimit.
func upAuther(c *bm.Context) {
	var (
		req = c.Request
		bs  []byte
		err error
	)
	if bs, err = ioutil.ReadAll(req.Body); err != nil {
		log.Error("ioutil.ReadAll() error(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	req.Body.Close()
	// params
	var ap = &archive.ArcParam{}
	if err = json.Unmarshal(bs, ap); err != nil {
		log.Error("http submit() json.Unmarshal(%s) error(%v)", string(bs), err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = vdaSvc.UpAuther(c, ap); err != nil {
		log.Error("vdaSvc.UpAuther() error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

// upAuther batch modify zlimit.
func upAccess(c *bm.Context) {
	var (
		req = c.Request
		bs  []byte
		err error
	)
	if bs, err = ioutil.ReadAll(req.Body); err != nil {
		log.Error("ioutil.ReadAll() error(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	req.Body.Close()
	// params
	var ap = &archive.ArcParam{}
	if err = json.Unmarshal(bs, ap); err != nil {
		log.Error("http submit() json.Unmarshal(%s) error(%v)", string(bs), err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = vdaSvc.UpAccess(c, ap); err != nil {
		log.Error("vdaSvc.UpAccess() error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func upArcTag(c *bm.Context) {
	pm := new(archive.TagParam)
	if err := c.Bind(pm); err != nil {
		c.JSON(nil, err)
		return
	}

	uid, _ := getUIDName(c)
	c.JSON(nil, vdaSvc.UpArcTag(c, uid, pm))
}

//通用批量tag接口 添加tag或者删除tag接口
// a.支持频道回查 								form_list = channel_review
// b.支持adminBind/upBind(默认走adminBind)    is_up_bind=true
// c.支持同步隐藏tag 							sync_hidden_tag=true
//todo 因为稿件服务不cache tags 也不需要发force_sync.未来计划砍掉审核库archive.tag
func batchTag(c *bm.Context) {
	var err error
	pm := new(archive.BatchTagParam)
	if err = c.Bind(pm); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	//仅同步隐藏tag时  批量tag操作 pm.tags为空 && SyncHiddenTag = true
	if (pm.Action != "" && pm.Action != "add" && pm.Action != "delete") || (pm.Action != "" && pm.Tags == "" && !pm.SyncHiddenTag) {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if pm.Action == "" && pm.Tags == "" && pm.FromList != archive.FromListChannelReview {
		c.JSON(nil, nil)
		return
	}

	uid, _ := getUIDName(c)
	data := map[string]interface{}{}
	data["message"], err = vdaSvc.BatchUpTag(c, uid, pm)
	c.JSONMap(data, err)
}

func channelInfo(c *bm.Context) {
	aid := c.Request.Form.Get("aid")
	aids, err := xstr.SplitInts(aid)
	if err != nil {
		log.Error("channelInfo xstr.SplitInts(%s) error(%v)", aid, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	//批量上线50个
	if len(aids) > 50 || len(aids) <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	data := map[string]interface{}{}
	data["data"], err = vdaSvc.GetChannelInfo(c, aids)
	if err != nil {
		data["message"] = "频道查询API异常"
	}
	c.JSONMap(data, err)
}

func aiTrack(c *bm.Context) {
	v := new(struct {
		Aid []int64 `form:"aid,split" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(vdaSvc.AITrack(c, v.Aid))
}
