package parallel

import (
	"context"
	"unsafe"

	"go-common/library/cache/redis"
	"go-common/library/log"
)

// RedisTask .
type RedisTask struct {
	ctx  *context.Context
	name string
	pool *redis.Pool
	cmd  string
	args []interface{}
}

// NewRedisTaskWithName new redis parallel task
func NewRedisTaskWithName(ctx *context.Context, name string, pool *redis.Pool, cmd string, args ...interface{}) *RedisTask {
	return &RedisTask{
		ctx:  ctx,
		name: name,
		pool: pool,
		cmd:  cmd,
		args: args,
	}
}

// NewRedisTask new redis parallel task
func NewRedisTask(ctx *context.Context, pool *redis.Pool, cmd string, args ...interface{}) *RedisTask {
	return &RedisTask{
		ctx:  ctx,
		pool: pool,
		cmd:  cmd,
		args: args,
	}
}

// Run .
func (rt *RedisTask) Run() (result *[]byte) {
	conn := rt.pool.Get(*rt.ctx)
	defer conn.Close()

	reply, err := conn.Do(rt.cmd, rt.args...)
	if err != nil {
		log.Error("RedisTask Run error:[%+v]", err)
		return
	}

	switch reply := reply.(type) {
	case []byte:
		result = &reply
	case string:
		b := []byte(reply)
		result = &b
	default:
		result = (*[]byte)(unsafe.Pointer(&reply))
	}

	return
}
