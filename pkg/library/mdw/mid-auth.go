package mdw

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/pkg/log"
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"github.com/dgrijalva/jwt-go"
)

type Dao interface {
	CheckAuthToken(ctx context.Context, userId int64, token string) (reply bool, err error)
	CheckPermission(ctx context.Context, userId int64, router string, action string) (reply bool, err error)
}

/**
 * 认证jwt、权限
 */
func AuthJwt(dao Dao, skippers ...SkipperFunc) bm.HandlerFunc {
	return func(ctx *bm.Context) {
		auth, exist := ctx.Get("jwt")
		if !exist {
			ctx.AbortWithStatus(401)
			return
		}
		u := auth.(*jwt.Token) //获取 token 信息
		jwtuserid := int64(u.Claims.(jwt.MapClaims)["userid"].(float64))
		ok, err := dao.CheckAuthToken(ctx, jwtuserid, u.Raw)
		if err != nil {
			log.Error("auth token failed: %v", err)
		}
		if !ok {
			ctx.AbortWithStatus(401)
			return
		}
		if SkipHandler(ctx, skippers...) {
			ctx.Set("auth_user_id", jwtuserid)
			ctx.Next()
			return
		}
		router := fmt.Sprintf("%s", ctx.Request.URL.Path)
		action := fmt.Sprintf("%s", ctx.Request.Method)
		if ok, _ := dao.CheckPermission(ctx, jwtuserid, router, action); !ok {
			ctx.AbortWithStatus(401)
			return
		}
		ctx.Set("auth_user_id", jwtuserid)
		ctx.Next()
	}
}

// SkipperFunc 定义中间件跳过函数
type SkipperFunc func(*bm.Context) bool

// AllowPathPrefixSkipper 检查请求路径是否包含指定的前缀，如果包含则跳过
func AllowPathPrefix(prefixes ...string) SkipperFunc {
	return func(c *bm.Context) bool {
		path := c.Request.URL.Path
		pathLen := len(path)

		for _, p := range prefixes {
			if pl := len(p); pathLen >= pl && path[:pl] == p {
				return true
			}
		}
		return false
	}
}

// AllowPathPrefixNoSkipper 检查请求路径是否包含指定的前缀，如果包含则不跳过
func AllowPathPrefixNo(prefixes ...string) SkipperFunc {
	return func(c *bm.Context) bool {
		path := c.Request.URL.Path
		pathLen := len(path)
		for _, p := range prefixes {
			if pl := len(p); pathLen >= pl && path[:pl] == p {
				return false
			}
		}
		return true
	}
}

// SkipHandler 统一处理跳过函数
func SkipHandler(c *bm.Context, skippers ...SkipperFunc) bool {
	for _, skipper := range skippers {
		if skipper(c) {
			return true
		}
	}
	return false
}
