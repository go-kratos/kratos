// Copyright 2012 Gary Burd
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package redis

import (
	"context"

	"github.com/bilibili/kratos/pkg/container/pool"
	xtime "github.com/bilibili/kratos/pkg/time"
)

// Error represents an error returned in a command reply.
type Error string

func (err Error) Error() string { return string(err) }

// Config client settings.
type Config struct {
	*pool.Config

	Name         string // redis name, for trace
	Proto        string
	Addr         string
	Auth         string
	DialTimeout  xtime.Duration
	ReadTimeout  xtime.Duration
	WriteTimeout xtime.Duration
	SlowLog      xtime.Duration
}

type Redis struct {
	pool *Pool
	conf *Config
}

func NewRedis(c *Config, options ...DialOption) *Redis {
	return &Redis{
		pool: NewPool(c, options...),
		conf: c,
	}
}

// Do gets a new conn from pool, then execute Do with this conn, finally close this conn.
// ATTENTION: Don't use this method with transaction command like MULTI etc. Because every Do will close conn automatically, use r.Conn to get a raw conn for this situation.
func (r *Redis) Do(ctx context.Context, commandName string, args ...interface{}) (reply interface{}, err error) {
	conn := r.pool.Get(ctx)
	defer conn.Close()
	reply, err = conn.Do(commandName, args...)
	return
}

// Close closes connection pool
func (r *Redis) Close() error {
	return r.pool.Close()
}

// Conn direct gets a connection
func (r *Redis) Conn(ctx context.Context) Conn {
	return r.pool.Get(ctx)
}

func (r *Redis) Pipeline() (p Pipeliner) {
	return &pipeliner{
		pool: r.pool,
	}
}
