package dao

import (
	"context"
	xsql "database/sql"
	notice "go-common/app/service/bbq/notice-service/api/v1"
	"go-common/app/service/bbq/user/api"
	"go-common/app/service/bbq/user/internal/model"
	acc "go-common/app/service/main/account/api"
	"go-common/library/database/sql"
	"go-common/library/time"
	"reflect"

	"github.com/bouk/monkey"
)

// MockUserBase .
func (d *Dao) MockUserBase(res map[int64]*api.UserBase, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "UserBase", func(_ *Dao, _ context.Context, _ []int64) (map[int64]*api.UserBase, error) {
		return res, err
	})
}

// MockPing .
func (d *Dao) MockPing(err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "Ping", func(_ *Dao, _ context.Context) error {
		return err
	})
}

// MockBeginTran .
func (d *Dao) MockBeginTran(p1 *xsql.Tx, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "BeginTran", func(_ *Dao, _ context.Context) (*xsql.Tx, error) {
		return p1, err
	})
}

// MockCreateNotice .
func (d *Dao) MockCreateNotice(err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "CreateNotice", func(_ *Dao, _ context.Context, _ *notice.NoticeBase) error {
		return err
	})
}

// MockFilter .
func (d *Dao) MockFilter(level int32, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "Filter", func(_ *Dao, _ context.Context, _ string, _ string) (int32, error) {
		return level, err
	})
}

// MockForbidUser .
func (d *Dao) MockForbidUser(err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "ForbidUser", func(_ *Dao, _ context.Context, _ uint64, _ uint64) error {
		return err
	})
}

// MockReleaseUser .
func (d *Dao) MockReleaseUser(err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "ReleaseUser", func(_ *Dao, _ context.Context, _ uint64) error {
		return err
	})
}

// MockGetLocation .
func (d *Dao) MockGetLocation(p1 *api.LocationItem, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "GetLocation", func(_ *Dao, _ context.Context, _ int32) (*api.LocationItem, error) {
		return p1, err
	})
}

// MockGetUserBProfile .
func (d *Dao) MockGetUserBProfile(res *acc.ProfileReply, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "GetUserBProfile", func(_ *Dao, _ context.Context, _ *api.PhoneCheckReq) (*acc.ProfileReply, error) {
		return res, err
	})
}

// MockRawUserBase .
func (d *Dao) MockRawUserBase(res map[int64]*api.UserBase, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "RawUserBase", func(_ *Dao, _ context.Context, _ []int64) (map[int64]*api.UserBase, error) {
		return res, err
	})
}

// MockCacheUserBase .
func (d *Dao) MockCacheUserBase(res map[int64]*api.UserBase, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "CacheUserBase", func(_ *Dao, _ context.Context, _ []int64) (map[int64]*api.UserBase, error) {
		return res, err
	})
}

// MockAddCacheUserBase .
func (d *Dao) MockAddCacheUserBase(err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "AddCacheUserBase", func(_ *Dao, _ context.Context, _ map[int64]*api.UserBase) error {
		return err
	})
}

// MockUpdateUserField .
func (d *Dao) MockUpdateUserField(num int64, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "UpdateUserField", func(_ *Dao, _ context.Context, _ *sql.Tx, _ int64, _ string, _ interface{}) (int64, error) {
		return num, err
	})
}

// MockAddUserBase .
func (d *Dao) MockAddUserBase(num int64, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "AddUserBase", func(_ *Dao, _ context.Context, _ *api.UserBase) (int64, error) {
		return num, err
	})
}

// MockUpdateUserBaseUname .
func (d *Dao) MockUpdateUserBaseUname(num int64, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "UpdateUserBaseUname", func(_ *Dao, _ context.Context, _ int64, _ string) (int64, error) {
		return num, err
	})
}

// MockUpdateUserBase .
func (d *Dao) MockUpdateUserBase(num int64, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "UpdateUserBase", func(_ *Dao, _ context.Context, _ int64, _ *api.UserBase) (int64, error) {
		return num, err
	})
}

// MockCheckUname .
func (d *Dao) MockCheckUname(err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "CheckUname", func(_ *Dao, _ context.Context, _ int64, _ string) error {
		return err
	})
}

// MockTxAddUserBlack .
func (d *Dao) MockTxAddUserBlack(num int64, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "TxAddUserBlack", func(_ *Dao, _ context.Context, _ *sql.Tx, _ int64, _ int64) (int64, error) {
		return num, err
	})
}

// MockTxCancelUserBlack .
func (d *Dao) MockTxCancelUserBlack(num int64, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "TxCancelUserBlack", func(_ *Dao, _ context.Context, _ *sql.Tx, _ int64, _ int64) (int64, error) {
		return num, err
	})
}

// MockFetchBlackList .
func (d *Dao) MockFetchBlackList(upMid []int64, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "FetchBlackList", func(_ *Dao, _ context.Context, _ int64) ([]int64, error) {
		return upMid, err
	})
}

// MockFetchPartBlackList .
func (d *Dao) MockFetchPartBlackList(MID2IDMap map[int64]time.Time, blackMIDs []int64, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "FetchPartBlackList", func(_ *Dao, _ context.Context, _ int64, _ model.CursorValue, _ int) (map[int64]time.Time, []int64, error) {
		return MID2IDMap, blackMIDs, err
	})
}

// MockIsBlack .
func (d *Dao) MockIsBlack(MIDMap map[int64]bool) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "IsBlack", func(_ *Dao, _ context.Context, _ int64, _ []int64) map[int64]bool {
		return MIDMap
	})
}

// MockRawUserCard .
func (d *Dao) MockRawUserCard(userCard *model.UserCard, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "RawUserCard", func(_ *Dao, _ context.Context, _ int64) (*model.UserCard, error) {
		return userCard, err
	})
}

// MockRawUserCards .
func (d *Dao) MockRawUserCards(userCards map[int64]*model.UserCard, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "RawUserCards", func(_ *Dao, _ context.Context, _ []int64) (map[int64]*model.UserCard, error) {
		return userCards, err
	})
}

// MockRawUserAccCards .
func (d *Dao) MockRawUserAccCards(res *acc.CardsReply, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "RawUserAccCards", func(_ *Dao, _ context.Context, _ []int64) (*acc.CardsReply, error) {
		return res, err
	})
}

// MockTxAddUserFan .
func (d *Dao) MockTxAddUserFan(num int64, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "TxAddUserFan", func(_ *Dao, _ *sql.Tx, _ int64, _ int64) (int64, error) {
		return num, err
	})
}

// MockTxCancelUserFan .
func (d *Dao) MockTxCancelUserFan(num int64, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "TxCancelUserFan", func(_ *Dao, _ *sql.Tx, _ int64, _ int64) (int64, error) {
		return num, err
	})
}

// MockIsFan .
func (d *Dao) MockIsFan(MIDMap map[int64]bool) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "IsFan", func(_ *Dao, _ context.Context, _ int64, _ []int64) map[int64]bool {
		return MIDMap
	})
}

// MockFetchPartFanList .
func (d *Dao) MockFetchPartFanList(MID2IDMap map[int64]time.Time, followedMIDs []int64, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "FetchPartFanList", func(_ *Dao, _ context.Context, _ int64, _ model.CursorValue, _ int) (map[int64]time.Time, []int64, error) {
		return MID2IDMap, followedMIDs, err
	})
}

// MockTxAddUserFollow .
func (d *Dao) MockTxAddUserFollow(num int64, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "TxAddUserFollow", func(_ *Dao, _ context.Context, _ *sql.Tx, _ int64, _ int64) (int64, error) {
		return num, err
	})
}

// MockTxCancelUserFollow .
func (d *Dao) MockTxCancelUserFollow(num int64, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "TxCancelUserFollow", func(_ *Dao, _ context.Context, _ *sql.Tx, _ int64, _ int64) (int64, error) {
		return num, err
	})
}

// MockFetchFollowList .
func (d *Dao) MockFetchFollowList(upMid []int64, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "FetchFollowList", func(_ *Dao, _ context.Context, _ int64) ([]int64, error) {
		return upMid, err
	})
}

// MockFetchPartFollowList .
func (d *Dao) MockFetchPartFollowList(MID2IDMap map[int64]time.Time, followedMIDs []int64, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "FetchPartFollowList", func(_ *Dao, _ context.Context, _ int64, _ model.CursorValue, _ int) (map[int64]time.Time, []int64, error) {
		return MID2IDMap, followedMIDs, err
	})
}

// MockIsFollow .
func (d *Dao) MockIsFollow(MIDMap map[int64]bool) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "IsFollow", func(_ *Dao, _ context.Context, _ int64, _ []int64) map[int64]bool {
		return MIDMap
	})
}

// MockTxAddUserLike .
func (d *Dao) MockTxAddUserLike(num int64, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "TxAddUserLike", func(_ *Dao, _ *sql.Tx, _ int64, _ int64) (int64, error) {
		return num, err
	})
}

// MockTxCancelUserLike .
func (d *Dao) MockTxCancelUserLike(num int64, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "TxCancelUserLike", func(_ *Dao, _ *sql.Tx, _ int64, _ int64) (int64, error) {
		return num, err
	})
}

// MockCheckUserLike .
func (d *Dao) MockCheckUserLike(res []int64, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "CheckUserLike", func(_ *Dao, _ context.Context, _ int64, _ []int64) ([]int64, error) {
		return res, err
	})
}

// MockGetUserLikeList .
func (d *Dao) MockGetUserLikeList(likeSvs []*api.LikeSv, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "GetUserLikeList", func(_ *Dao, _ context.Context, _ int64, _ bool, _ model.CursorValue, _ int) ([]*api.LikeSv, error) {
		return likeSvs, err
	})
}

// MockRawBatchUserStatistics .
func (d *Dao) MockRawBatchUserStatistics(res map[int64]*api.UserStat, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "RawBatchUserStatistics", func(_ *Dao, _ context.Context, _ []int64) (map[int64]*api.UserStat, error) {
		return res, err
	})
}

// MockTxIncrUserStatisticsFollow .
func (d *Dao) MockTxIncrUserStatisticsFollow(num int64, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "TxIncrUserStatisticsFollow", func(_ *Dao, _ *sql.Tx, _ int64) (int64, error) {
		return num, err
	})
}

// MockTxIncrUserStatisticsFan .
func (d *Dao) MockTxIncrUserStatisticsFan(num int64, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "TxIncrUserStatisticsFan", func(_ *Dao, _ *sql.Tx, _ int64) (int64, error) {
		return num, err
	})
}

// MockTxDecrUserStatisticsFollow .
func (d *Dao) MockTxDecrUserStatisticsFollow(num int64, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "TxDecrUserStatisticsFollow", func(_ *Dao, _ *sql.Tx, _ int64) (int64, error) {
		return num, err
	})
}

// MockTxDecrUserStatisticsFan .
func (d *Dao) MockTxDecrUserStatisticsFan(num int64, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "TxDecrUserStatisticsFan", func(_ *Dao, _ *sql.Tx, _ int64) (int64, error) {
		return num, err
	})
}

// MockTxIncrUserStatisticsField .
func (d *Dao) MockTxIncrUserStatisticsField(rowsAffected int64, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "TxIncrUserStatisticsField", func(_ *Dao, _ context.Context, _ *sql.Tx, _ int64, _ string) (int64, error) {
		return rowsAffected, err
	})
}

// MockTxDescUserStatisticsField .
func (d *Dao) MockTxDescUserStatisticsField(rowsAffected int64, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "TxDescUserStatisticsField", func(_ *Dao, _ context.Context, _ *sql.Tx, _ int64, _ string) (int64, error) {
		return rowsAffected, err
	})
}

// MockUpdateUserVideoView .
func (d *Dao) MockUpdateUserVideoView(err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "UpdateUserVideoView", func(_ *Dao, _ context.Context, _ int64, _ int64) error {
		return err
	})
}
