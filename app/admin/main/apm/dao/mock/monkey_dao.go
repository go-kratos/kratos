package mock

import (
	"context"
	"go-common/app/admin/main/apm/dao"
	cml "go-common/app/admin/main/apm/model/canal"
	"go-common/app/admin/main/apm/model/need"

	"go-common/app/admin/main/apm/model/ecode"
	"go-common/app/admin/main/apm/model/pprof"
	"go-common/app/admin/main/apm/model/ut"

	"reflect"

	"github.com/bouk/monkey"
)

// MockDaoSetConfigID .
func MockDaoSetConfigID(d *dao.Dao, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "SetConfigID", func(_ *dao.Dao, _ int64, _ string) error {
		return err
	})
}

// MockDaoCanalInfoCounts .
func MockDaoCanalInfoCounts(d *dao.Dao, cnt int, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "CanalInfoCounts", func(_ *dao.Dao, _ *cml.ConfigReq) (int, error) {
		return cnt, err
	})
}

// MockDaoCanalInfoEdit .
func MockDaoCanalInfoEdit(d *dao.Dao, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "CanalInfoEdit", func(_ *dao.Dao, _ *cml.ConfigReq) error {
		return err
	})
}

// MockDaoCanalApplyCounts .
func MockDaoCanalApplyCounts(d *dao.Dao, cnt int, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "CanalApplyCounts", func(_ *dao.Dao, _ *cml.ConfigReq) (int, error) {
		return cnt, err
	})
}

// MockDaoCanalApplyEdit .
func MockDaoCanalApplyEdit(d *dao.Dao, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "CanalApplyEdit", func(_ *dao.Dao, _ *cml.ConfigReq, _ string) error {
		return err
	})
}

// MockDaoCanalApplyCreate .
func MockDaoCanalApplyCreate(d *dao.Dao, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "CanalApplyCreate", func(_ *dao.Dao, _ *cml.ConfigReq, _ string) error {
		return err
	})
}

// MockDaoGetCodes .
func MockDaoGetCodes(d *dao.Dao, data []*codes.Codes, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "GetCodes", func(_ *dao.Dao, _ context.Context, _ string, _ string) ([]*codes.Codes, error) {
		return data, err
	})
}

// MockDaoNeedInfoAdd .
func MockDaoNeedInfoAdd(d *dao.Dao, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "NeedInfoAdd", func(_ *dao.Dao, _ *need.NAddReq, _ string) error {
		return err
	})
}

// MockDaoNeedInfoList .
func MockDaoNeedInfoList(d *dao.Dao, res []*need.NInfo, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "NeedInfoList", func(_ *dao.Dao, _ *need.NListReq) ([]*need.NInfo, error) {
		return res, err
	})
}

// MockDaoNeedInfoCount .
func MockDaoNeedInfoCount(d *dao.Dao, count int64, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "NeedInfoCount", func(_ *dao.Dao, _ *need.NListReq) (int64, error) {
		return count, err
	})
}

// MockDaoGetNeedInfo .
func MockDaoGetNeedInfo(d *dao.Dao, r *need.NInfo, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "GetNeedInfo", func(_ *dao.Dao, _ int64) (*need.NInfo, error) {
		return r, err
	})
}

// MockDaoNeedInfoEdit .
func MockDaoNeedInfoEdit(d *dao.Dao, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "NeedInfoEdit", func(_ *dao.Dao, _ *need.NEditReq) error {
		return err
	})
}

// MockDaoNeedVerify .
func MockDaoNeedVerify(d *dao.Dao, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "NeedVerify", func(_ *dao.Dao, _ *need.NVerifyReq) error {
		return err
	})
}

// MockDaoLikeCountsStats .
func MockDaoLikeCountsStats(d *dao.Dao, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "LikeCountsStats", func(_ *dao.Dao, _ *need.Likereq, _ int, _ int) error {
		return err
	})
}

// MockDaoGetVoteInfo .
func MockDaoGetVoteInfo(d *dao.Dao, u *need.UserLikes, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "GetVoteInfo", func(_ *dao.Dao, _ *need.Likereq, _ string) (*need.UserLikes, error) {
		return u, err
	})
}

// MockDaoUpdateVoteInfo .
func MockDaoUpdateVoteInfo(d *dao.Dao, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "UpdateVoteInfo", func(_ *dao.Dao, _ *need.Likereq, _ string) error {
		return err
	})
}

// MockDaoAddVoteInfo .
func MockDaoAddVoteInfo(d *dao.Dao, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "AddVoteInfo", func(_ *dao.Dao, _ *need.Likereq, _ string) error {
		return err
	})
}

// MockDaoVoteInfoList .
func MockDaoVoteInfoList(d *dao.Dao, res []*need.UserLikes, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "VoteInfoList", func(_ *dao.Dao, _ *need.Likereq) ([]*need.UserLikes, error) {
		return res, err
	})
}

// MockDaoVoteInfoCounts .
func MockDaoVoteInfoCounts(d *dao.Dao, count int64, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "VoteInfoCounts", func(_ *dao.Dao, _ *need.Likereq) (int64, error) {
		return count, err
	})
}

// MockDaoInstances .
func MockDaoInstances(d *dao.Dao, ins *pprof.Ins, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "Instances", func(_ *dao.Dao, _ context.Context, _ string) (*pprof.Ins, error) {
		return ins, err
	})
}

// MockDaoGitLabFace .
func MockDaoGitLabFace(d *dao.Dao, avatarURL string, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "GitLabFace", func(_ *dao.Dao, _ context.Context, _ string) (string, error) {
		return avatarURL, err
	})
}

// MockDaoUploadProxy .
func MockDaoUploadProxy(d *dao.Dao, url string, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "UploadProxy", func(_ *dao.Dao, _ context.Context, _ string, _ int64, _ []byte) (string, error) {
		return url, err
	})
}

// MockDaoParseUTFiles .
func MockDaoParseUTFiles(d *dao.Dao, pkgs []*ut.PkgAnls, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "ParseUTFiles", func(_ *dao.Dao, _ context.Context, _ string) ([]*ut.PkgAnls, error) {
		return pkgs, err
	})
}

// MockDaoSendWechatToUsers .
func MockDaoSendWechatToUsers(d *dao.Dao, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "SendWechatToUsers", func(_ *dao.Dao, _ context.Context, _ []string, _ string) error {
		return err
	})
}

// MockDaoSendWechatToGroup .
func MockDaoSendWechatToGroup(d *dao.Dao, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "SendWechatToGroup", func(_ *dao.Dao, _ context.Context, _ string, _ string) error {
		return err
	})
}

// MockDaoGitLabCommits .
func MockDaoGitLabCommits(d *dao.Dao, commit *ut.GitlabCommit, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "GitLabCommits", func(_ *dao.Dao, _ context.Context, _ string) (*ut.GitlabCommit, error) {
		return commit, err
	})
}

// MockDaoGetCoverage .
func MockDaoGetCoverage(d *dao.Dao, cov float64, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "GetCoverage", func(_ *dao.Dao, _ context.Context, _ string, _ string) (float64, error) {
		return cov, err
	})
}

// MockDaoSetAppCovCache .
func MockDaoSetAppCovCache(d *dao.Dao, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "SetAppCovCache", func(_ *dao.Dao, _ context.Context) error {
		return err
	})
}

// MockDaoGetAppCovCache .
func MockDaoGetAppCovCache(d *dao.Dao, coverage float64, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "GetAppCovCache", func(_ *dao.Dao, _ context.Context, _ string) (float64, error) {
		return coverage, err
	})
}
