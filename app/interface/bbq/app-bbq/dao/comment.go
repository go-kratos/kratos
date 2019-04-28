package dao

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"

	"go-common/app/interface/bbq/app-bbq/model"
	user "go-common/app/service/bbq/user/api"
	"go-common/library/log"
	"go-common/library/net/metadata"

	"github.com/json-iterator/go"
)

// ReplySubCursor http请求子评论游标获取评论
func (d *Dao) ReplySubCursor(c context.Context, req map[string]interface{}) (res *model.SubCursorRes, err error) {
	var r []byte
	ip := metadata.String(c, metadata.RemoteIP)
	res = new(model.SubCursorRes)
	r, err = replyHTTPCommon(c, d.httpClient, d.c.URLs["reply_subcursor"], "GET", req, ip)
	if err != nil {
		log.Infov(c,
			log.KV("log", fmt.Sprintf("replyHTTPCommon err [%v]", err)),
		)
		return
	}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	decoder := json.NewDecoder(bytes.NewBuffer(r))
	decoder.UseNumber()
	err = decoder.Decode(&res)
	// err = json.Unmarshal(r, &res)
	if err != nil {
		log.Errorv(c, log.KV("log", fmt.Sprintf("json unmarlshal err data[%s]", string(r))))
		return
	}
	//获取所有mid
	var mid []int64
	if res.Root != nil {
		mid = append(mid, res.Root.Mid)
	}
	if res.Root.Replies != nil {
		mid = append(mid, d.getMidInReplys(res.Root.Replies)...)
	}
	//批量获取userinfo
	var userinfo map[int64]*user.UserBase
	userinfo, err = d.JustGetUserBase(c, mid)
	if err != nil {
		log.Errorv(c, log.KV("log", fmt.Sprintf("d.UserBase err[%v]", err)))
		return
	}
	if res.Root != nil {
		d.replaceMemberInReply(res.Root, userinfo)
	}
	return
}

// ReplyCursor http请求游标获取评论
func (d *Dao) ReplyCursor(c context.Context, req map[string]interface{}) (res *model.CursorRes, err error) {
	var r []byte
	ip := metadata.String(c, metadata.RemoteIP)
	res = new(model.CursorRes)
	r, err = replyHTTPCommon(c, d.httpClient, d.c.URLs["reply_cursor"], "GET", req, ip)
	if err != nil {
		log.Infov(c,
			log.KV("log", fmt.Sprintf("replyHTTPCommon err [%v]", err)),
		)
		return
	}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	decoder := json.NewDecoder(bytes.NewBuffer(r))
	decoder.UseNumber()
	err = decoder.Decode(&res)
	if err != nil {
		log.Errorv(c, log.KV("log", fmt.Sprintf("json unmarlshal err data[%s]", string(r))))
		return
	}
	//获取所有mid
	var mid []int64
	if res.Replies != nil {
		mid = append(mid, d.getMidInReplys(res.Replies)...)
	}
	if res.Hots != nil {
		mid = append(mid, d.getMidInReplys(res.Hots)...)
	}
	if res.Top != nil {
		if res.Top.Admin != nil {
			mid = append(mid, res.Top.Admin.Mid)
		}
		if res.Top.Upper != nil {
			mid = append(mid, res.Top.Upper.Mid)
		}
	}
	//批量获取userinfo
	var userinfo map[int64]*user.UserBase
	userinfo, err = d.JustGetUserBase(c, mid)
	if err != nil {
		log.Errorv(c, log.KV("log", fmt.Sprintf("d.UserBase err[%v]", err)))
		return
	}
	if res.Replies != nil {
		d.replaceMemberInReplys(res.Replies, userinfo)
	}
	if res.Hots != nil {
		d.replaceMemberInReplys(res.Hots, userinfo)
	}
	if res.Top != nil {
		if res.Top.Admin != nil {
			d.replaceMemberInReplys(res.Top.Admin.Replies, userinfo)
		}
		if res.Top.Upper != nil {
			d.replaceMemberInReplys(res.Top.Upper.Replies, userinfo)
		}
	}
	return
}

// getMidInReplys 批量获取评论列表中mid
func (d *Dao) getMidInReplys(r []*model.Reply) (l []int64) {
	for _, v := range r {
		l = append(l, v.Mid)
		for _, c := range v.Content.Members {
			cmid, err := strconv.ParseInt(c.Mid, 10, 64)
			if err != nil {
				log.Errorv(
					context.Background(),
					log.KV("log", fmt.Sprintf("strconv err [%v] data[%s]", err, c.Mid)),
				)
				continue
			}
			l = append(l, cmid)
		}
		if len(v.Replies) > 0 {
			sl := d.getMidInReplys(v.Replies)
			l = append(l, sl...)
		}
	}
	return
}

// replaceMemberInReplys 批量替换评论列表中用户信息
func (d *Dao) replaceMemberInReplys(r []*model.Reply, umap map[int64]*user.UserBase) {
	for _, v := range r {
		if u, ok := umap[v.Mid]; ok {
			v.Member.BInfo = u

		} else {
			v.Member.BInfo = new(user.UserBase)
			v.Member.BInfo.Mid = v.Mid
			v.Member.BInfo.Uname = v.Member.Name
			v.Member.BInfo.Face = v.Member.Avatar
		}
		for _, c := range v.Content.Members {
			cmid, err := strconv.ParseInt(c.Mid, 10, 64)
			if err != nil {
				log.Errorv(
					context.Background(),
					log.KV("log", fmt.Sprintf("strconv err [%v] data[%s]", err, c.Mid)),
				)
				continue
			}
			if u, ok := umap[cmid]; ok {
				c.BInfo = u
			} else {
				c.BInfo = new(user.UserBase)
				c.BInfo.Mid = cmid
				c.BInfo.Uname = c.Name
				c.BInfo.Face = c.Avatar
			}
		}
		if len(v.Replies) > 0 {
			d.replaceMemberInReplys(v.Replies, umap)
		}
	}
}

// replaceMemberInReply 替换评论中用户信息
func (d *Dao) replaceMemberInReply(r *model.Reply, umap map[int64]*user.UserBase) {
	if u, ok := umap[r.Mid]; ok {
		r.Member.BInfo = u
	} else {
		r.Member.BInfo = new(user.UserBase)
		r.Member.BInfo.Mid = r.Mid
		r.Member.BInfo.Uname = r.Member.Name
		r.Member.BInfo.Face = r.Member.Avatar
	}
	for _, c := range r.Content.Members {
		cmid, err := strconv.ParseInt(c.Mid, 10, 64)
		if err != nil {
			log.Errorv(
				context.Background(),
				log.KV("log", fmt.Sprintf("strconv err [%v] data[%s]", err, c.Mid)),
			)
			continue
		}
		if u, ok := umap[cmid]; ok {
			c.BInfo = u
		} else {
			c.BInfo = new(user.UserBase)
			c.BInfo.Mid = cmid
			c.BInfo.Uname = c.Name
			c.BInfo.Face = c.Avatar
		}
	}
	if len(r.Replies) > 0 {
		d.replaceMemberInReplys(r.Replies, umap)
	}
}

// ReplyAdd http请求添加评论
func (d *Dao) ReplyAdd(c context.Context, req map[string]interface{}) (res *model.AddRes, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	res = new(model.AddRes)
	var r []byte
	r, err = replyHTTPCommon(c, d.httpClient, d.c.URLs["reply_add"], "POST", req, ip)
	if err != nil {
		log.Infov(c,
			log.KV("log", fmt.Sprintf("replyHTTPCommon err [%v]", err)),
		)
	}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	json.Unmarshal(r, &res)
	return
}

//ReplyLike http请求评论点赞
func (d *Dao) ReplyLike(c context.Context, req map[string]interface{}) (err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	_, err = replyHTTPCommon(c, d.httpClient, d.c.URLs["reply_like"], "POST", req, ip)
	return
}

// ReplyList 评论列表
func (d *Dao) ReplyList(c context.Context, req map[string]interface{}) (res *model.ReplyList, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	res = new(model.ReplyList)
	var r []byte
	r, err = replyHTTPCommon(c, d.httpClient, d.c.URLs["reply_list"], "GET", req, ip)
	if err != nil {
		log.Infov(c,
			log.KV("log", fmt.Sprintf("replyHTTPCommon err [%v]", err)),
		)
		return
	}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	decoder := json.NewDecoder(bytes.NewBuffer(r))
	decoder.UseNumber()
	err = decoder.Decode(&res)
	if err != nil {
		log.Errorv(c, log.KV("log", fmt.Sprintf("json unmarlshal err data[%s]", string(r))))
		return
	}
	//获取所有mid
	var mid []int64
	if res.Replies != nil {
		mid = append(mid, d.getMidInReplys(res.Replies)...)
	}
	if res.Hots != nil {
		mid = append(mid, d.getMidInReplys(res.Hots)...)
	}
	if res.Top != nil {
		mid = append(mid, res.Top.Mid)
	}
	//批量获取userinfo
	var userinfo map[int64]*user.UserBase
	userinfo, err = d.JustGetUserBase(c, mid)
	if err != nil {
		log.Errorv(c, log.KV("log", fmt.Sprintf("d.UserBase err[%v]", err)))
		return
	}
	if res.Replies != nil {
		d.replaceMemberInReplys(res.Replies, userinfo)
	}
	if res.Hots != nil {
		d.replaceMemberInReplys(res.Hots, userinfo)
	}
	if res.Top != nil {
		d.replaceMemberInReplys(res.Top.Replies, userinfo)
	}
	return
}

// ReplyReport 评论举报
func (d *Dao) ReplyReport(c context.Context, req map[string]interface{}) (err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	_, err = replyHTTPCommon(c, d.httpClient, d.c.URLs["reply_report"], "POST", req, ip)
	return
}

// ReplyCounts 批量评论数
func (d *Dao) ReplyCounts(c context.Context, ids []int64, t int64) (res map[int64]*model.ReplyCount, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	oidStr := strings.Replace(strings.Trim(fmt.Sprint(ids), "[]"), " ", ",", -1)
	req := map[string]interface{}{
		"type": t,
		"oid":  oidStr,
	}
	res = make(map[int64]*model.ReplyCount)
	var r []byte
	r, err = replyHTTPCommon(c, d.httpClient, d.c.URLs["reply_counts"], "GET", req, ip)
	if err != nil {
		log.Infov(c,
			log.KV("log", fmt.Sprintf("replyHTTPCommon err [%v]", err)),
		)
		return
	}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	decoder := json.NewDecoder(bytes.NewBuffer(r))
	decoder.UseNumber()
	err = decoder.Decode(&res)
	if err != nil {
		log.Errorv(c, log.KV("log", fmt.Sprintf("json unmarlshal err data[%s]", string(r))))
	}
	return
}

// ReplyMinfo 批量请求评论
func (d *Dao) ReplyMinfo(c context.Context, svID int64, rpids []int64) (res map[int64]*model.Reply, err error) {
	res = make(map[int64]*model.Reply)
	if len(rpids) == 0 {
		return
	}
	ip := metadata.String(c, metadata.RemoteIP)
	var rpidStr string
	for _, rpid := range rpids {
		if len(rpidStr) != 0 {
			rpidStr += fmt.Sprintf(",%d", rpid)
		} else {
			rpidStr += fmt.Sprintf("%d", rpid)
		}

	}
	log.V(1).Infov(c, log.KV("log", "reply minfo: rpids="+rpidStr))

	req := map[string]interface{}{
		"type": model.DefaultCmType,
		"oid":  svID,
		"rpid": rpidStr,
	}
	var r []byte
	r, err = replyHTTPCommon(c, d.httpClient, d.c.URLs["reply_minfo"], "GET", req, ip)
	if err != nil {
		log.Errorv(c,
			log.KV("log", fmt.Sprintf("replyHTTPCommon err [%v]", err)),
		)
		return
	}

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	decoder := json.NewDecoder(bytes.NewBuffer(r))
	decoder.UseNumber()
	err = decoder.Decode(&res)
	if err != nil {
		log.Errorv(c, log.KV("log", fmt.Sprintf("json unmarlshal err data[%s]", string(r))))
	}
	log.V(1).Infov(c, log.KV("log", "get reply minfo"), log.KV("req_size", len(rpids)), log.KV("rsp_size", len(res)))
	return

}

// ReplyHot 获取批量热评
func (d *Dao) ReplyHot(c context.Context, mid int64, oids []int64) (res map[int64][]*model.Reply, err error) {
	res = make(map[int64][]*model.Reply)
	if len(oids) == 0 {
		return
	}
	ip := metadata.String(c, metadata.RemoteIP)
	var oidStr string
	for _, oid := range oids {
		if len(oidStr) != 0 {
			oidStr += fmt.Sprintf(",%d", oid)
		} else {
			oidStr += fmt.Sprintf("%d", oid)
		}

	}
	log.V(1).Infov(c, log.KV("log", "reply minfo: oids="+oidStr))

	req := map[string]interface{}{
		"type": model.DefaultCmType,
		"oids": oidStr,
		"mid":  mid,
	}
	var r []byte
	r, err = replyHTTPCommon(c, d.httpClient, d.c.URLs["reply_hots"], "GET", req, ip)
	if err != nil {
		log.Errorv(c,
			log.KV("log", fmt.Sprintf("replyHTTPCommon err [%v]", err)),
		)
		return
	}

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	decoder := json.NewDecoder(bytes.NewBuffer(r))
	decoder.UseNumber()
	err = decoder.Decode(&res)
	if err != nil {
		log.Errorv(c, log.KV("log", fmt.Sprintf("json unmarlshal err data[%s]", string(r))))
	}
	log.V(1).Infov(c, log.KV("log", "get reply hots"), log.KV("req_size", len(oids)), log.KV("rsp_size", len(res)))

	// 替换bbq的用户信息
	var mIDs []int64
	for _, replies := range res {
		for _, reply := range replies {
			mIDs = append(mIDs, reply.Mid)
		}
	}
	var userBases map[int64]*user.UserBase
	userBases, err = d.JustGetUserBase(c, mIDs)
	if err != nil {
		log.Errorv(c, log.KV("log", fmt.Sprintf("d.UserBase err[%v]", err)))
		return
	}
	for _, replies := range res {
		if len(replies) > 0 {
			d.replaceMemberInReplys(replies, userBases)
		}
	}

	return

}
