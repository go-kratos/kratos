package rank

import (
	"context"
	"fmt"

	"go-common/library/log"
)

var (
	_ PrepRank = &ElecPrepUPRank{}
	_ PrepRank = &ElecPrepAVRank{}
)

// PrepRank interface .
type PrepRank interface {
	fmt.Stringer
	// 从中间存储中拉取榜单数据
	Load(ctx context.Context) (rank interface{}, err error)
	// 从中间存储中存储榜单数据
	Save(ctx context.Context, rank interface{}) (err error)
	// 从DB中恢复榜单数据
	Rebuild(ctx context.Context) (rank interface{}, err error)
	// 从订单更新榜单
	UpdateOrder(ctx context.Context, rank interface{}, payMID int64, fee int64) (res interface{}, err error)
	// 从留言更新榜单
	UpdateMessage(ctx context.Context, rank interface{}, payMID int64, message string, hidden bool) (res interface{}, err error)
}

// Rank interface .
type Rank interface {
	// 榜单描述信息，一般用于log输出调试信息
	fmt.Stringer
	Storager
	// 从 prepRank 中恢复榜单数据
	Rebuild(ctx context.Context, prepRank interface{}) (rank interface{}, err error)
}

// Storager .
type Storager interface {
	// 从中间存储中拉取榜单数据
	Load(ctx context.Context) (rank interface{}, err error)
	// 从中间存储中存储榜单数据
	Save(ctx context.Context, rank interface{}) (err error)
}

func tryHard(fn func() (ok bool, e error), fname string, tryTimes int) (err error) {
	var ok bool
	for tryTimes > 0 {
		tryTimes--
		ok, err = fn()
		if err != nil {
			log.Error("tryHard func: %s, ok: %t err: %+v, try times left: %d", fname, ok, err, tryTimes)
			return
		}
		if !ok {
			log.Error("tryHard func: %s, ok: %t err: %+v, try times left: %d", fname, ok, err, tryTimes)
			continue
		}
		// log.Info("tryHard func: %s, run success, try times left: %d", fname, tryTimes)
		return
	}
	return
}
