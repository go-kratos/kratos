package dao

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"sync"
	"time"

	"go-common/app/admin/main/workflow/model"
	"go-common/app/admin/main/workflow/model/param"
	credit "go-common/app/interface/main/credit/model"
	"go-common/app/job/main/member/model/block"
	acc "go-common/app/service/main/account/api"
	arc "go-common/app/service/main/archive/api"
	member "go-common/app/service/main/member/model/block"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_upergroupURI     = "http://api.bilibili.co/x/internal/uper/special/get"
	_tagListURI       = "http://manager.bilibili.co/x/admin/manager/internal/tag/list"
	_blockURI         = "http://api.bilibili.co/x/internal/block/batch/block"
	_creditBlockedURI = "http://api.bilibili.co/x/internal/credit/blocked/info/add"
	_blockInfoURI     = "http://api.bilibili.co/x/internal/block/info"
	_blockNumURI      = "http://api.bilibili.co/x/internal/credit/blocked/user/num"
	_blockCaseAddURI  = "http://api.bilibili.co/x/internal/credit/blocked/case/add"
)

// BatchUperSpecial .
// http://info.bilibili.co/pages/viewpage.action?pageId=8479274
func (d *Dao) BatchUperSpecial(c context.Context, mids []int64) (UperTagMap map[int64][]*model.SpecialTag, err error) {
	uri := _upergroupURI
	uperSpecialResp := new(model.UperSpecial)
	UperTagMap = make(map[int64][]*model.SpecialTag)
	uv := url.Values{}
	uv.Set("group_id", "0")
	uv.Set("mids", xstr.JoinInts(mids))
	uv.Set("pn", "1")
	uv.Set("ps", "1000")

	if err = d.httpRead.Get(c, uri, "", uv, uperSpecialResp); err != nil {
		err = errors.Wrap(err, fmt.Sprintf("search uper special tag failed mids(%v)", uv.Get("mids")))
		return
	}

	if uperSpecialResp.Code != ecode.OK.Code() {
		log.Error("call %s result error code(%d), message(%s)", uri, uperSpecialResp.Code, uperSpecialResp.Message)
		err = ecode.Int(uperSpecialResp.Code)
		return
	}

	for _, special := range uperSpecialResp.Data.Items {
		UperTagMap[special.MID] = append(UperTagMap[special.MID], special)
	}
	return
}

// ArchiveRPC .
func (d *Dao) ArchiveRPC(c context.Context, oids []int64) (archives map[int64]*model.Archive, err error) {
	if len(oids) == 0 {
		return
	}
	archives = make(map[int64]*model.Archive, len(oids))
	var res *arc.ArcsReply
	arg := &arc.ArcsRequest{
		Aids: oids,
	}
	if res, err = d.arcRPC.Arcs(c, arg); err != nil {
		log.Error("d.arcRPC.Archives3(%+v) error(%v)", arg, err)
		return
	}
	for oid, arc := range res.Arcs {
		tmplArc := &model.Archive{
			Author: arc.Author.Name,
			State:  arc.State,
			Mid:    arc.Author.Mid,
			TypeID: arc.TypeID,
			Type:   arc.TypeName,
			Title:  arc.Title,
		}
		archives[oid] = tmplArc
	}
	return
}

// TagList .
// http://info.bilibili.co/pages/viewpage.action?pageId=9831467
// tags map[bid][tagid]
func (d *Dao) TagList(c context.Context) (tags map[int8]map[int64]*model.TagMeta, err error) {
	uri := _tagListURI
	uv := url.Values{}
	uv.Set("ps", "1000")
	uv.Set("pn", "1")
	result := new(model.TagListResult)
	tags = make(map[int8]map[int64]*model.TagMeta)
	if err = d.httpRead.Get(c, uri, "", uv, result); err != nil {
		return
	}
	if result.Code != ecode.OK.Code() {
		log.Error("tag list failed: %s?%s, error code(%d), message(%s)", uri, uv.Encode(), result.Code, result.Message)
		err = ecode.Int(result.Code)
		return
	}

	for _, tm := range result.Data.Tags {
		if _, ok := tags[tm.Bid]; !ok {
			tags[tm.Bid] = make(map[int64]*model.TagMeta)
		}
		tags[tm.Bid][tm.TagID] = tm
	}
	return
}

// CommonExtraInfo return common external info
func (d *Dao) CommonExtraInfo(c context.Context, bid int8, uri string, ids, oids, eids []int64) (data map[string]interface{}, err error) {
	log.Info("start call common extra info bid(%d) gids(%v) oids(%v) eids(%v)", bid, ids, oids, eids)
	data = make(map[string]interface{})
	uv := url.Values{}
	uv.Set("bid", strconv.FormatInt(int64(bid), 10))
	uv.Set("ids", xstr.JoinInts(ids))
	uv.Set("oids", xstr.JoinInts(oids))
	uv.Set("eids", xstr.JoinInts(eids))
	res := &model.CommonExtraDataResponse{}
	if err = d.httpRead.Get(c, uri, "", uv, &res); err != nil {
		return data, err
	}
	if res.Code != ecode.OK.Code() {
		log.Error("get extra info failed: url(%s), code(%d), message(%s), bid(%d)", uri+"?"+uv.Encode(), res.Code, res.Message, bid)
		err = ecode.Int(res.Code)
		return
	}
	log.Info("success extra info (%+v) req url(%s)", res, uri+"?"+uv.Encode())
	data = res.Data
	return
}

// AccountInfoRPC .
func (d *Dao) AccountInfoRPC(c context.Context, mids []int64) (authors map[int64]*model.Account) {
	g := &errgroup.Group{}
	mutex := sync.RWMutex{}
	authors = make(map[int64]*model.Account)
	if len(mids) == 0 {
		return
	}
	// distinct mid
	mMap := make(map[int64]bool)
	dMids := make([]int64, 0)
	for _, m := range mids {
		if _, ok := mMap[m]; !ok {
			mMap[m] = false
			dMids = append(dMids, m)
			continue
		}
	}
	for _, mid := range dMids {
		gMid := mid
		g.Go(func() (err error) {
			var res *acc.ProfileStatReply
			start := time.Now()
			arg := &acc.MidReq{Mid: gMid}
			if res, err = d.accRPC.ProfileWithStat3(c, arg); err != nil {
				log.Error("d.accRPC.ProfileWithStat3(%v) error(%v)", arg, err)
				err = nil
				return
			}
			mutex.Lock()
			acc := &model.Account{
				Mid:      res.Profile.Mid,
				Name:     res.Profile.Name,
				Rank:     res.Profile.Rank,
				Follower: res.Follower,
				Official: &model.Official{Role: res.Profile.Official.Role},
			}
			authors[res.Profile.Mid] = acc
			log.Info("mid(%d) data(%+v) wrap success", res.Profile.Mid, acc)
			mutex.Unlock()

			log.Info("account rpc request gmid(%d) time %s", gMid, time.Since(start).String())
			return
		})
	}
	g.Wait()
	return
}

// AddMoral 扣节操
func (d *Dao) AddMoral(c context.Context, mids []int64, gssp *param.GroupStateSetParam) (err error) {
	var errFlag bool
	for _, mid := range mids {
		arg := &acc.MoralReq{
			Mid:    mid,
			Moral:  float64(gssp.DecreaseMoral),
			Oper:   gssp.AdminName,
			Reason: gssp.Reason,
			Remark: "workflow",
		}
		if _, err = d.accRPC.AddMoral3(c, arg); err != nil {
			log.Error("failed decrease moral arg(%+v) error(%v)", arg, err)
			errFlag = true
		}
	}
	if !errFlag {
		log.Info("add moral success mids(%v) param(%+v)", mids, gssp)
	}
	return
}

// AddBlock 发起账号封禁
// http://info.bilibili.co/pages/viewpage.action?pageId=7559616
func (d *Dao) AddBlock(c context.Context, mids []int64, gssp *param.GroupStateSetParam) (err error) {
	uri := _blockURI
	uv := url.Values{}
	uv.Set("mids", xstr.JoinInts(mids))
	source := strconv.Itoa(int(member.BlockSourceBlackHouse))
	uv.Set("source", source) //来源: 后台相关
	if gssp.Business == model.CommentComplain {
		area := strconv.Itoa(int(member.BlockAreaReply))
		uv.Set("area", area) //违规业务
	}
	if gssp.BlockDay > 0 { //限时封禁
		action := strconv.Itoa(int(block.BlockActionLimit))
		uv.Set("action", action)
		uv.Set("duration", strconv.FormatInt(gssp.BlockDay*86400, 10))
	} else { //永久封禁
		action := strconv.Itoa(int(block.BlockActionForever))
		uv.Set("action", action)
	}
	uv.Set("start_time", strconv.FormatInt(time.Now().Unix(), 10))
	uv.Set("op_id", strconv.FormatInt(gssp.AdminID, 10))
	uv.Set("operator", gssp.AdminName)
	reason := model.BlockReason[gssp.BlockReason]
	uv.Set("reason", reason)
	uv.Set("notify", "1")

	resp := &model.CommonResponse{}
	if err = d.httpWrite.Post(c, uri, "", uv, resp); err != nil {
		log.Error("add block url(%s) param(%s) error(%v)", uri, uv.Encode(), err)
		return
	}
	if resp.Code != ecode.OK.Code() {
		err = ecode.Int(resp.Code)
		log.Error("call add block url(%s) param(%s) error res code(%d)", uri, uv.Encode(), resp.Code)
		return
	}
	log.Info("add block success mids(%v) param(%+v)", mids, gssp)
	return
}

// AddCreditBlockInfo 上报封禁信息到小黑屋
// http://info.bilibili.co/pages/viewpage.action?pageId=5417571
func (d *Dao) AddCreditBlockInfo(c context.Context, bus map[int64]*model.Business, gssp *param.GroupStateSetParam) (err error) {
	// 请求小黑屋
	uri := _creditBlockedURI
	for _, b := range bus {
		uv := url.Values{}
		uv.Set("mid", strconv.FormatInt(b.Mid, 10))
		uv.Set("oper_id", strconv.FormatInt(gssp.AdminID, 10))
		uv.Set("origin_content", b.Title)
		if gssp.Business == model.CommentComplain {
			uv.Set("origin_type", strconv.Itoa(int(member.BlockAreaReply))) //违规业务
		}
		if gssp.BlockDay > 0 { //限时封禁
			uv.Set("blocked_days", strconv.FormatInt(gssp.BlockDay, 10))
			uv.Set("blocked_forever", strconv.Itoa(int(credit.NotInBlockedForever)))
			uv.Set("punish_type", strconv.Itoa(int(credit.PunishTypeBlock)))
		} else { //永久封禁
			uv.Set("blocked_forever", strconv.Itoa(int(credit.InBlockedForever)))
			uv.Set("punish_type", strconv.Itoa(int(credit.PunishTypeForever)))
		}

		uv.Set("punish_time", strconv.FormatInt(time.Now().Unix(), 10))
		uv.Set("reason_type", strconv.Itoa(int(gssp.BlockReason)))
		uv.Set("operator_name", gssp.AdminName)

		resp := &model.CommonResponse{}
		if err = d.httpWrite.Post(c, uri, "", uv, resp); err != nil {
			log.Error("add credit block info url(%s) param(%s) error(%v)", uri, uv.Encode(), err)
			continue
		}
		if resp.Code != ecode.OK.Code() {
			err = ecode.Int(resp.Code)
			log.Error("call add credit block info url(%s) param(%s) error res code(%d)", uri, uv.Encode(), resp.Code)
			continue
		}
		log.Info("add credit block info success mid(%v) param(%+v)", b.Mid, gssp)
	}
	return
}

// AddCreditCase 请求风纪委众裁
func (d *Dao) AddCreditCase(c context.Context, uv url.Values) (err error) {
	uri := _blockCaseAddURI
	var res model.CommonResponse
	if err = d.httpWrite.Post(c, uri, "", uv, &res); err != nil {
		log.Error("d.httpWrite.Post(%s) body(%s) error(%v)", uri, uv.Encode(), err)
		err = ecode.WkfSetPublicRefereeFailed
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("call url(%s) body(%s) error code(%d)", uri, uv.Encode(), res.Code)
		err = ecode.WkfSetPublicRefereeFailed
		return
	}
	log.Info("call block add case success url(%s) body(%s)", uri, uv.Encode())
	return
}

// BlockNum 查询封禁次数
func (d *Dao) BlockNum(c context.Context, mid int64) (sum int64, err error) {
	uri := _blockNumURI
	uv := url.Values{}
	uv.Set("mid", strconv.FormatInt(mid, 10))
	var numResp model.BlockNumResp
	if err = d.httpRead.Get(c, uri, "", uv, &numResp); err != nil {
		log.Error("d.httpRead.Get() error(%v) url(%s?%s)", err, uri, uv.Encode())
		err = ecode.WkfGetBlockInfoFailed
		return
	}
	if numResp.Code != ecode.OK.Code() {
		log.Error("call url(%s?%s) response code (%d) error", uri, uv.Encode(), numResp.Code)
		err = ecode.WkfGetBlockInfoFailed
		return
	}
	sum = numResp.Data.BlockedSum
	log.Info("url(%s) mid(%d) block num (%d)", uri, mid, sum)
	return
}

// BlockInfo 查询封禁信息
func (d *Dao) BlockInfo(c context.Context, mid int64) (resp model.BlockInfoResp, err error) {
	uri := _blockInfoURI
	uv := url.Values{}
	uv.Set("mid", strconv.FormatInt(mid, 10))
	if err = d.httpRead.Get(c, uri, "", uv, &resp); err != nil {
		log.Error("d.httpRead.Get() error(%v) url(%s?%s)", err, uri, uv.Encode())
		err = ecode.WkfGetBlockInfoFailed
		return
	}
	if resp.Code != ecode.OK.Code() {
		log.Error("call url(%s?%s) response code (%d) error", uri, uv.Encode(), resp.Code)
		err = ecode.WkfGetBlockInfoFailed
		return
	}
	log.Info("url(%s) mid(%d) block info resp(%+v)", uri, mid, resp)
	return
}

// SourceInfo 返回业务来源
func (d *Dao) SourceInfo(c context.Context, uri string) (data map[string]interface{}, err error) {
	log.Info("start call SourceInfo uri(%s)", uri)
	res := &model.SourceQueryResponse{}
	if err = d.httpRead.Get(c, uri, "", nil, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("get SourceInfo failed: url(%s), code(%d), message(%s)", uri, res.Code, res.Message)
		err = ecode.Int(res.Code)
		return
	}
	log.Info("success SourceInfo (%+v)", res.Data)
	return res.Data, nil
}
