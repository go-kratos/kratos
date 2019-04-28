package auth

import (
	"strconv"
	"time"

	xsql "go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	xauth "go-common/library/net/http/blademaster/middleware/auth"
)

const (
	_queryBannedUser = "select `expire_time`,`forbid_status` from `forbid_user` where `mid` = ?;"
)

// BannedAuth .
type BannedAuth struct {
	*xauth.Auth
	db *xsql.DB
}

// NewBannedAuth .
func NewBannedAuth(authConf *xauth.Config, dbConf *xsql.Config) *BannedAuth {
	return &BannedAuth{
		Auth: xauth.New(authConf),
		db:   xsql.NewMySQL(dbConf),
	}
}

// User .
func (ba *BannedAuth) User(ctx *bm.Context) {
	ba.Auth.User(ctx)

	// 只过滤写接口
	if ctx.Request.Method != "POST" {
		return
	}

	// 获取MID
	var mid int64
	tmp, _ := ctx.Get("mid")
	switch tmp.(type) {
	case int64:
		mid = tmp.(int64)
	}
	if mid == 0 {
		return
	}

	// 过滤封禁用户
	var expired, status int
	row := ba.db.QueryRow(ctx, _queryBannedUser, mid)
	row.Scan(&expired, &status)
	if status == 1 && expired > 0 && expired > int(time.Now().Unix()) {
		// 封禁用户
		log.Infov(ctx, log.KV("log", "user banned "+strconv.Itoa(int(mid))), log.KV("err", ecode.BBQUserBanned))
		ctx.Error = ecode.BBQUserBanned
		ctx.JSON(nil, ecode.BBQUserBanned)
		ctx.Abort()
	}
}
