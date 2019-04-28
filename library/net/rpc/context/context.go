package context

import (
	ctx "context"
	"time"
)

// Context web context interface
type Context interface {
	ctx.Context

	// Now get current time.
	Now() time.Time

	// Seq implement Context method Seq.
	Seq() uint64

	// ServiceMethod implement Context method ServiceMethod.
	ServiceMethod() string

	// User get caller user.
	User() string
}

type rpcCtx struct {
	ctx.Context
	now           time.Time
	seq           uint64
	serviceMethod string
	user          string
}

// NewContext new a rpc context.
func NewContext(c ctx.Context, m, u string, s uint64) Context {
	rc := &rpcCtx{Context: c, now: time.Now(), seq: s, serviceMethod: m, user: u}
	return rc
}

func (c *rpcCtx) Seq() uint64 {
	return c.seq
}

func (c *rpcCtx) ServiceMethod() string {
	return c.serviceMethod
}

func (c *rpcCtx) Now() time.Time {
	return c.now
}

func (c *rpcCtx) User() string {
	return c.user
}
