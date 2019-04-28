package util

import (
	"context"

	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

const (
	_ajsSessioID = "_AJSESSIONID"
)

type sessionKeyT string

var sessionKey sessionKeyT = "sessionID"

// SessionIDMiddleware extrace session from cookie set to context
func SessionIDMiddleware(c *bm.Context) {
	cookie, err := c.Request.Cookie(_ajsSessioID)
	if err != nil {
		c.JSON(nil, ecode.AccessDenied)
		c.Abort()
		return
	}
	c.Context = SessionIDWithContext(c.Context, cookie.Value)
}

// SessionIDFromContext get session id from context
func SessionIDFromContext(ctx context.Context) string {
	if val := ctx.Value(sessionKey); val != nil {
		if sessionID, ok := val.(string); ok {
			return sessionID
		}
	}
	return ""
}

// SessionIDWithContext set session id to context
func SessionIDWithContext(ctx context.Context, sessionID string) context.Context {
	return context.WithValue(ctx, sessionKey, sessionID)
}
