package dao

import (
	"context"
	"fmt"

	model "go-common/app/interface/main/credit/model"
	gmc "go-common/library/cache/memcache"
)

const (
	_prefixCaseInfo     = "ca_in_%d"
	_prefixVoteCaseInfo = "vc_in_%d_%d"
	_caseVoteTop        = "ca_vo_top_%d"
)

func caseInfoKey(cid int64) string {
	return fmt.Sprintf(_prefixCaseInfo, cid)
}

func voteCaseInfoKey(mid int64, cid int64) string {
	return fmt.Sprintf(_prefixVoteCaseInfo, mid, cid)
}

func caseVoteTopKey(mid int64) string {
	return fmt.Sprintf(_caseVoteTop, mid)
}

// SetCaseInfoCache set case info.
func (d *Dao) SetCaseInfoCache(c context.Context, cid int64, bc *model.BlockedCase) (err error) {
	var (
		conn = d.mc.Get(c)
		item = &gmc.Item{Key: caseInfoKey(cid), Object: bc, Expiration: d.minCommonExpire, Flags: gmc.FlagJSON}
	)
	defer conn.Close()
	err = conn.Set(item)
	return
}

// CaseInfoCache get case info cache
func (d *Dao) CaseInfoCache(c context.Context, cid int64) (bc *model.BlockedCase, err error) {
	var (
		reply *gmc.Item
		conn  = d.mc.Get(c)
	)
	defer conn.Close()
	if reply, err = conn.Get(caseInfoKey(cid)); err != nil {
		if err == gmc.ErrNotFound {
			err = nil
		}
		return
	}
	bc = &model.BlockedCase{}
	err = conn.Scan(reply, &bc)
	return
}

// SetVoteInfoCache set vote case info.
func (d *Dao) SetVoteInfoCache(c context.Context, mid, cid int64, vi *model.VoteInfo) (err error) {
	var (
		conn = d.mc.Get(c)
		item = &gmc.Item{Key: voteCaseInfoKey(mid, cid), Object: vi, Expiration: d.userExpire, Flags: gmc.FlagJSON}
	)
	defer conn.Close()
	err = conn.Set(item)
	return
}

// VoteInfoCache get vote case info cache
func (d *Dao) VoteInfoCache(c context.Context, mid, cid int64) (vi *model.VoteInfo, err error) {
	var (
		reply *gmc.Item
		conn  = d.mc.Get(c)
	)
	defer conn.Close()
	reply, err = conn.Get(voteCaseInfoKey(mid, cid))
	if err != nil {
		if err == gmc.ErrNotFound {
			err = nil
		}
		return
	}
	vi = &model.VoteInfo{}
	err = conn.Scan(reply, &vi)
	return
}

// CaseVoteTopCache get case votes by top 100.
func (d *Dao) CaseVoteTopCache(c context.Context, mid int64) (bs []*model.BlockedCase, err error) {
	var (
		reply *gmc.Item
		conn  = d.mc.Get(c)
	)
	defer conn.Close()
	reply, err = conn.Get(caseVoteTopKey(mid))
	if err != nil {
		if err == gmc.ErrNotFound {
			err = nil
		}
		return
	}
	bs = make([]*model.BlockedCase, 0)
	err = conn.Scan(reply, &bs)
	return
}

// SetCaseVoteTopCache set case votes by top 100.
func (d *Dao) SetCaseVoteTopCache(c context.Context, mid int64, bs []*model.BlockedCase) (err error) {
	var (
		conn = d.mc.Get(c)
		item = &gmc.Item{Key: caseVoteTopKey(mid), Object: bs, Expiration: d.userExpire, Flags: gmc.FlagJSON}
	)
	defer conn.Close()
	err = conn.Set(item)
	return
}
