package service

import (
	"context"
	"go-common/app/service/bbq/user/api"
	"go-common/library/ecode"
	"go-common/library/log"
	"testing"
	"time"
)

func TestService_AddUserFollow(t *testing.T) {
	ctx := context.Background()
	// 初始化状态
	s.ModifyRelation(ctx, &api.ModifyRelationReq{Mid: 88895104, UpMid: 88895134, Action: 2})
	s.ModifyRelation(ctx, &api.ModifyRelationReq{Mid: 88895104, UpMid: 88895134, Action: 4})

	originStat1 := getUserStat(88895104)
	originStat2 := getUserStat(88895134)
	s.ModifyRelation(ctx, &api.ModifyRelationReq{Mid: 88895104, UpMid: 88895134, Action: 1})
	// 关注是否生效
	res := s.dao.IsFollow(ctx, 88895104, []int64{88895134})
	log.Info("res: %v", res)
	if len(res) == 0 {
		t.Errorf("user follow fail")
	}
	res = s.dao.IsFan(ctx, 88895134, []int64{88895104})
	log.Info("res: %v", res)
	if len(res) == 0 {
		t.Errorf("user follow fail")
	}

	// 关注对统计数据的影响
	curStat1 := getUserStat(88895104)
	if originStat1.Follow+1 != curStat1.Follow {
		t.Errorf("follow stat fail: origin=%d, cur=%d", originStat1.Follow, curStat1.Follow)
	}
	curStat2 := getUserStat(88895134)
	if originStat2.Fan+1 != curStat2.Fan {
		t.Errorf("fan stat fail: origin=%d, cur=%d", originStat2.Fan, curStat2.Fan)
	}

	// 拉黑状态下不能关注
	s.ModifyRelation(ctx, &api.ModifyRelationReq{Mid: 88895104, UpMid: 88895134, Action: 2})
	s.ModifyRelation(ctx, &api.ModifyRelationReq{Mid: 88895104, UpMid: 88895134, Action: 4})
	s.ModifyRelation(ctx, &api.ModifyRelationReq{Mid: 88895104, UpMid: 88895134, Action: 3})
	_, err := s.ModifyRelation(ctx, &api.ModifyRelationReq{Mid: 88895104, UpMid: 88895134, Action: 1})
	if err != ecode.UserAlreadyBlackFollowErr {
		t.Errorf("follow error when up_mid in black: err=%v", err)
	}

}

func TestService_CancelUserFollow(t *testing.T) {
	ctx := context.Background()
	// 初始化状态
	s.ModifyRelation(ctx, &api.ModifyRelationReq{Mid: 88895104, UpMid: 88895134, Action: 2})
	s.ModifyRelation(ctx, &api.ModifyRelationReq{Mid: 88895104, UpMid: 88895134, Action: 4})

	s.ModifyRelation(ctx, &api.ModifyRelationReq{Mid: 88895104, UpMid: 88895134, Action: 1})

	originStat1 := getUserStat(88895104)
	originStat2 := getUserStat(88895134)
	// 取关后的效果
	s.ModifyRelation(ctx, &api.ModifyRelationReq{Mid: 88895104, UpMid: 88895134, Action: 2})
	res := s.dao.IsFollow(ctx, 88895104, []int64{88895134})
	log.Info("res: %v", res)
	if len(res) != 0 {
		t.Errorf("cancel user follow fail")
	}
	res = s.dao.IsFan(ctx, 88895134, []int64{88895104})
	log.Info("res: %v", res)
	if len(res) != 0 {
		t.Errorf("cancel user follow fail")
	}

	// 对统计数据的影响
	curStat1 := getUserStat(88895104)
	if originStat1.Follow-1 != curStat1.Follow {
		t.Errorf("follow stat fail: origin=%d, cur=%d", originStat1.Follow, curStat1.Follow)
	}
	curStat2 := getUserStat(88895134)
	if originStat2.Fan-1 != curStat2.Fan {
		t.Errorf("fan stat fail: origin=%d, cur=%d", originStat2.Fan, curStat2.Fan)
	}

}

func getUserStat(mid int64) (stat *api.UserStat) {
	ctx := context.Background()
	res, _ := s.dao.RawBatchUserStatistics(ctx, []int64{mid})
	stat = res[mid]
	return
}

// 测试多次重复的关注取关效果
func TestDuplicateFollow(t *testing.T) {
	ctx := context.Background()
	// 清除
	s.ModifyRelation(ctx, &api.ModifyRelationReq{Mid: 88895104, UpMid: 88895134, Action: 2})
	s.ModifyRelation(ctx, &api.ModifyRelationReq{Mid: 88895104, UpMid: 88895134, Action: 4})

	originStat1 := getUserStat(88895104)
	originStat2 := getUserStat(88895134)
	s.ModifyRelation(ctx, &api.ModifyRelationReq{Mid: 88895104, UpMid: 88895134, Action: 1})
	s.ModifyRelation(ctx, &api.ModifyRelationReq{Mid: 88895104, UpMid: 88895134, Action: 1})
	s.ModifyRelation(ctx, &api.ModifyRelationReq{Mid: 88895104, UpMid: 88895134, Action: 2})
	s.ModifyRelation(ctx, &api.ModifyRelationReq{Mid: 88895104, UpMid: 88895134, Action: 2})
	s.ModifyRelation(ctx, &api.ModifyRelationReq{Mid: 88895104, UpMid: 88895134, Action: 1})
	s.ModifyRelation(ctx, &api.ModifyRelationReq{Mid: 88895104, UpMid: 88895134, Action: 1})
	s.ModifyRelation(ctx, &api.ModifyRelationReq{Mid: 88895104, UpMid: 88895134, Action: 1})
	s.ModifyRelation(ctx, &api.ModifyRelationReq{Mid: 88895104, UpMid: 88895134, Action: 2})
	curStat1 := getUserStat(88895104)
	curStat2 := getUserStat(88895134)
	if curStat1.Follow != originStat1.Follow {
		t.Errorf("follow stat fail: cur=%d, origin=%d", curStat1.Follow, originStat1.Follow)
	}
	if curStat2.Fan != originStat2.Fan {
		t.Errorf("fan stat fail: cur=%d, origin=%d", curStat1.Follow, originStat1.Follow)
	}
}

func TestUserStatistics(t *testing.T) {

	ctx := context.Background()
	// 初始化
	s.ModifyRelation(ctx, &api.ModifyRelationReq{Mid: 88895104, UpMid: 88895134, Action: 2})
	s.ModifyRelation(ctx, &api.ModifyRelationReq{Mid: 88895104, UpMid: 88895134, Action: 4})

	s.ModifyRelation(ctx, &api.ModifyRelationReq{Mid: 88895104, UpMid: 88895134, Action: 2})
	res, _ := s.dao.RawBatchUserStatistics(ctx, []int64{88895104, 88895134})
	stat := res[88895104]
	follow := stat.Follow
	stat = res[88895134]
	fan := stat.Fan

	s.ModifyRelation(ctx, &api.ModifyRelationReq{Mid: 88895104, UpMid: 88895134, Action: 1})
	time.Sleep(time.Second)
	res, _ = s.dao.RawBatchUserStatistics(ctx, []int64{88895104, 88895134})
	stat = res[88895104]
	if follow+1 != stat.Follow {
		t.Errorf("follow statistics fail")
	}
	stat = res[88895134]
	if fan+1 != stat.Fan {
		t.Errorf("fan statistics fail")
	}

	s.ModifyRelation(ctx, &api.ModifyRelationReq{Mid: 88895104, UpMid: 88895134, Action: 2})
	time.Sleep(time.Second)
	res, _ = s.dao.RawBatchUserStatistics(ctx, []int64{88895104, 88895134})
	log.Info("stat: %v", res)
	stat = res[88895104]
	if follow != stat.Follow {
		t.Error("follow statistics fail", follow, stat.Follow)
	}
	stat = res[88895134]
	if fan != stat.Fan {
		t.Error("fan statistics fail", fan, stat.Fan)
	}

}

func TestAddUserBlack(t *testing.T) {
	ctx := context.Background()
	s.ModifyRelation(ctx, &api.ModifyRelationReq{Mid: 88895104, UpMid: 88895134, Action: 3})
	res := s.dao.IsBlack(ctx, 88895104, []int64{88895134})
	log.Info("res: %v", res)
	if len(res) == 0 {
		t.Errorf("user follow fail")
	}
}
func TestCancelUserBlack(t *testing.T) {
	ctx := context.Background()
	s.ModifyRelation(ctx, &api.ModifyRelationReq{Mid: 88895104, UpMid: 88895134, Action: 4})
	res := s.dao.IsBlack(ctx, 88895104, []int64{88895134})
	log.Info("res: %v", res)
	if len(res) != 0 {
		t.Errorf("cancel user follow fail")
	}
}

func TestUserBlack(t *testing.T) {
	ctx := context.Background()
	// 初始化清除关系
	s.ModifyRelation(ctx, &api.ModifyRelationReq{Mid: 88895104, UpMid: 88895134, Action: 2})
	s.ModifyRelation(ctx, &api.ModifyRelationReq{Mid: 88895104, UpMid: 88895134, Action: 4})

	// 关注时拉黑，会清楚拉黑状态
	TestService_AddUserFollow(t)
	_, err := s.addUserBlack(ctx, 88895104, 88895134)
	if err != nil {
		t.Errorf("follow state black case fail")
	}
	res1 := s.dao.IsBlack(ctx, 88895104, []int64{88895134})
	if len(res1) == 0 {
		t.Errorf("Black fail")
	}
	res1 = s.dao.IsFollow(ctx, 88895104, []int64{88895134})
	if len(res1) != 0 {
		t.Errorf("Black fail")
	}
	TestService_CancelUserFollow(t)

	// 拉黑不能关注
	s.ModifyRelation(ctx, &api.ModifyRelationReq{Mid: 88895104, UpMid: 88895134, Action: 3})
	_, err = s.addUserFollow(ctx, 88895104, 88895134)
	if err == nil {
		t.Errorf("black stat follow case fail")
	}
	s.ModifyRelation(ctx, &api.ModifyRelationReq{Mid: 88895104, UpMid: 88895134, Action: 4})

	// 拉黑数量
	res, _ := s.dao.RawBatchUserStatistics(ctx, []int64{88895104})
	stat := res[88895104]
	black := stat.Black

	s.ModifyRelation(ctx, &api.ModifyRelationReq{Mid: 88895104, UpMid: 88895134, Action: 3})
	res, _ = s.dao.RawBatchUserStatistics(ctx, []int64{88895104})
	stat = res[88895104]
	if black+1 != stat.Black {
		t.Error("black stat fail")
	}

	s.ModifyRelation(ctx, &api.ModifyRelationReq{Mid: 88895104, UpMid: 88895134, Action: 4})
	res, _ = s.dao.RawBatchUserStatistics(ctx, []int64{88895104})
	stat = res[88895104]
	if black != stat.Black {
		t.Error("black stat fail")
	}

	s.ModifyRelation(ctx, &api.ModifyRelationReq{Mid: 88895104, UpMid: 88895134, Action: 3})
	s.ModifyRelation(ctx, &api.ModifyRelationReq{Mid: 88895104, UpMid: 88895134, Action: 4})
	s.ModifyRelation(ctx, &api.ModifyRelationReq{Mid: 88895104, UpMid: 88895134, Action: 3})
	s.ModifyRelation(ctx, &api.ModifyRelationReq{Mid: 88895104, UpMid: 88895134, Action: 4})
	s.ModifyRelation(ctx, &api.ModifyRelationReq{Mid: 88895104, UpMid: 88895134, Action: 3})
	s.ModifyRelation(ctx, &api.ModifyRelationReq{Mid: 88895104, UpMid: 88895134, Action: 4})
	stat = getUserStat(88895104)
	if black != stat.Black {
		t.Error("black stat fail")
	}

}

func TestService_ListFollowUserInfo(t *testing.T) {

	ctx := context.Background()
	s.ListFollowUserInfo(ctx, &api.ListRelationUserInfoReq{Mid: 88895134, UpMid: 88895104})
	// TODO:
}

func TestService_ListFanUserInfo(t *testing.T) {
	ctx := context.Background()
	s.ListFanUserInfo(ctx, &api.ListRelationUserInfoReq{Mid: 88895104, UpMid: 88895134})
	// TODO:
}

func TestService_ListBlackUserInfo(t *testing.T) {
	ctx := context.Background()
	s.ListBlackUserInfo(ctx, &api.ListRelationUserInfoReq{Mid: 88895134, UpMid: 88895104})
}

func TestService_ListFollow(t *testing.T) {
	ctx := context.Background()
	s.ModifyRelation(ctx, &api.ModifyRelationReq{Mid: 88895204, UpMid: 88895134, Action: 1})
	stat := getUserStat(88895204)
	reply, _ := s.ListFollow(ctx, &api.ListRelationReq{Mid: 88895204})
	if stat.Follow != int64(len(reply.List)) {
		t.Errorf("get follow list fail")
	}

}

func TestService_ListBlack(t *testing.T) {
	ctx := context.Background()
	s.ModifyRelation(ctx, &api.ModifyRelationReq{Mid: 88895204, UpMid: 88895134, Action: 3})
	stat := getUserStat(88895204)
	reply, _ := s.ListBlack(ctx, &api.ListRelationReq{Mid: 88895204})
	if stat.Black != int64(len(reply.List)) {
		t.Errorf("get follow list fail")
	}
}
