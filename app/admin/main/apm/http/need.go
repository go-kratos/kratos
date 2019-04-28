package http

import (
	"context"
	"go-common/app/admin/main/apm/conf"
	"go-common/app/admin/main/apm/model/need"

	bm "go-common/library/net/http/blademaster"
)

// @params NListReq
// @router get /x/admin/apm/need/list
// @response NListResp
func needList(c *bm.Context) {
	var (
		infos []*need.NInfo
		total int64
		err   error
	)
	username := name(c)
	arg := new(need.NListReq)
	if err = c.Bind(arg); err != nil {
		return
	}
	if infos, total, err = apmSvc.NeedInfoList(c, arg, username); err != nil {
		c.JSON(nil, err)
		return
	}
	res := &need.NListResp{
		Data:  infos,
		Total: total,
	}
	c.JSON(res, nil)
}

// @params NAddReq
// @router post /x/admin/apm/need/add
// @response EmpResp
func needAdd(c *bm.Context) {
	var err error
	username := name(c)

	req := new(need.NAddReq)
	if err = c.Bind(req); err != nil {
		return
	}
	if err = apmSvc.NeedInfoAdd(c, req, username); err != nil {
		c.JSON(nil, err)
		return
	}
	go apmSvc.SendWeMessage(context.Background(), req.Title, need.VerifyType[need.NeedApply], "", username, conf.Conf.Superman)
	go apmSvc.SendWeMessage(context.Background(), req.Title, need.VerifyType[need.NeedVerify], "", username, []string{username})

	c.JSON(nil, err)

}

// @params NEditReq
// @router post /x/admin/apm/need/edit
// @response EmpResp
func needEdit(c *bm.Context) {
	var err error
	req := new(need.NEditReq)
	if err = c.Bind(req); err != nil {
		return
	}
	username := name(c)

	if err = apmSvc.NeedInfoEdit(c, req, username); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, err)

}

// @params NVerifyReq
// @router post /x/admin/apm/need/verify
// @response EmpResp
func needVerify(c *bm.Context) {
	var (
		err error
		ni  *need.NInfo
	)
	req := new(need.NVerifyReq)
	if err = c.Bind(req); err != nil {
		return
	}
	username := name(c)
	if ni, err = apmSvc.NeedInfoVerify(c, req); err != nil {
		c.JSON(nil, err)
		return
	}
	go apmSvc.SendWeMessage(context.Background(), ni.Title, need.VerifyType[need.NeedReview], need.VerifyType[req.Status], username, []string{username})

	c.JSON(nil, err)

}

// @params Likereq
// @router post /x/admin/apm/need/thumbsup
// @response EmpResp
func needThumbsUp(c *bm.Context) {
	var err error
	req := new(need.Likereq)
	if err = c.Bind(req); err != nil {
		return
	}
	username := name(c)
	if err = apmSvc.NeedInfoVote(c, req, username); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, err)

}

// @params Likereq
// @router get /x/admin/need/vote/list
// @response VoteListResp
func needVoteList(c *bm.Context) {
	var (
		votes []*need.UserLikes
		total int64
		err   error
	)
	req := new(need.Likereq)
	if err = c.Bind(req); err != nil {
		return
	}
	if votes, total, err = apmSvc.NeedVoteList(c, req); err != nil {
		c.JSON(nil, err)
		return
	}
	res := &need.VoteListResp{
		Data:  votes,
		Total: total,
	}
	c.JSON(res, nil)
}
