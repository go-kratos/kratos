package dao

import (
	"context"
	cml "go-common/app/admin/main/apm/model/canal"
	"go-common/app/admin/main/apm/model/ecode"
	"go-common/app/admin/main/apm/model/need"
	"go-common/app/admin/main/apm/model/pprof"
	"go-common/app/admin/main/apm/model/ut"
	"reflect"

	"github.com/bouk/monkey"
)

//MockSetConfigID is
func (d *Dao) MockSetConfigID(err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "SetConfigID", func(_ *Dao, _ int64, _ string) error {
		return err
	})
}

//MockCanalInfoCounts is
func (d *Dao) MockCanalInfoCounts(cnt int, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "CanalInfoCounts", func(_ *Dao, _ *cml.ConfigReq) (int, error) {
		return cnt, err
	})
}

//MockCanalInfoEdit is
func (d *Dao) MockCanalInfoEdit(err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "CanalInfoEdit", func(_ *Dao, _ *cml.ConfigReq) error {
		return err
	})
}

//MockCanalApplyCounts is
func (d *Dao) MockCanalApplyCounts(cnt int, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "CanalApplyCounts", func(_ *Dao, _ *cml.ConfigReq) (int, error) {
		return cnt, err
	})
}

//MockCanalApplyEdit is
func (d *Dao) MockCanalApplyEdit(err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "CanalApplyEdit", func(_ *Dao, _ *cml.ConfigReq, _ string) error {
		return err
	})
}

//MockCanalApplyCreate is
func (d *Dao) MockCanalApplyCreate(err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "CanalApplyCreate", func(_ *Dao, _ *cml.ConfigReq, _ string) error {
		return err
	})
}

//MockGetCodes is
func (d *Dao) MockGetCodes(data []*codes.Codes, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "GetCodes", func(_ *Dao, _ context.Context, _ string, _ string) ([]*codes.Codes, error) {
		return data, err
	})
}

//MockNeedInfoAdd is
func (d *Dao) MockNeedInfoAdd(err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "NeedInfoAdd", func(_ *Dao, _ *need.NAddReq, _ string) error {
		return err
	})
}

//MockNeedInfoList is
func (d *Dao) MockNeedInfoList(res []*need.NInfo, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "NeedInfoList", func(_ *Dao, _ *need.NListReq) ([]*need.NInfo, error) {
		return res, err
	})
}

//MockNeedInfoCount is
func (d *Dao) MockNeedInfoCount(count int64, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "NeedInfoCount", func(_ *Dao, _ *need.NListReq) (int64, error) {
		return count, err
	})
}

//MockGetNeedInfo is
func (d *Dao) MockGetNeedInfo(r []*need.NInfo, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "GetNeedInfo", func(_ *Dao, _ int64) ([]*need.NInfo, error) {
		return r, err
	})
}

//MockNeedInfoEdit is
func (d *Dao) MockNeedInfoEdit(err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "NeedInfoEdit", func(_ *Dao, _ *need.NEditReq) error {
		return err
	})
}

//MockNeedVerify is
func (d *Dao) MockNeedVerify(err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "NeedVerify", func(_ *Dao, _ *need.NVerifyReq) error {
		return err
	})
}

//MockLikeCountsStats is
func (d *Dao) MockLikeCountsStats(err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "LikeCountsStats", func(_ *Dao, _ *need.Likereq, _ int, _ int) error {
		return err
	})
}

//MockGetVoteInfo is
func (d *Dao) MockGetVoteInfo(u []*need.UserLikes, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "GetVoteInfo", func(_ *Dao, _ *need.Likereq, _ string) ([]*need.UserLikes, error) {
		return u, err
	})
}

//MockUpdateVoteInfo is
func (d *Dao) MockUpdateVoteInfo(err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "UpdateVoteInfo", func(_ *Dao, _ *need.Likereq, _ string) error {
		return err
	})
}

//MockAddVoteInfo is
func (d *Dao) MockAddVoteInfo(err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "AddVoteInfo", func(_ *Dao, _ *need.Likereq, _ string) error {
		return err
	})
}

//MockVoteInfoList is
func (d *Dao) MockVoteInfoList(res []*need.UserLikes, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "VoteInfoList", func(_ *Dao, _ *need.Likereq) ([]*need.UserLikes, error) {
		return res, err
	})
}

//MockVoteInfoCounts is
func (d *Dao) MockVoteInfoCounts(count int64, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "VoteInfoCounts", func(_ *Dao, _ *need.Likereq) (int64, error) {
		return count, err
	})
}

//MockInstances is
func (d *Dao) MockInstances(ins *pprof.Ins, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "Instances", func(_ *Dao, _ context.Context, _ string) (*pprof.Ins, error) {
		return ins, err
	})
}

//MockGitLabFace is
func (d *Dao) MockGitLabFace(avatarURL string, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "GitLabFace", func(_ *Dao, _ context.Context, _ string) (string, error) {
		return avatarURL, err
	})
}

//MockUploadProxy is
func (d *Dao) MockUploadProxy(url string, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "UploadProxy", func(_ *Dao, _ context.Context, _ string, _ int64, _ []byte) (string, error) {
		return url, err
	})
}

//MockParseUTFiles is
func (d *Dao) MockParseUTFiles(pkgs []*ut.PkgAnls, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "ParseUTFiles", func(_ *Dao, _ context.Context, _ string) ([]*ut.PkgAnls, error) {
		return pkgs, err
	})
}

//MockSendWechatToUsers is
func (d *Dao) MockSendWechatToUsers(err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "SendWechatToUsers", func(_ *Dao, _ context.Context, _ []string, _ string) error {
		return err
	})
}

//MockSendWechatToGroup is
func (d *Dao) MockSendWechatToGroup(err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "SendWechatToGroup", func(_ *Dao, _ context.Context, _ string, _ string) error {
		return err
	})
}

//MockGitLabCommits is
func (d *Dao) MockGitLabCommits(commit *ut.GitlabCommit, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "GitLabCommits", func(_ *Dao, _ context.Context, _ string) (*ut.GitlabCommit, error) {
		return commit, err
	})
}

//MockGetCoverage is
func (d *Dao) MockGetCoverage(cov float64, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "GetCoverage", func(_ *Dao, _ context.Context, _ string, _ string) (float64, error) {
		return cov, err
	})
}
