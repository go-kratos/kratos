package bnj

import (
	"time"

	"go-common/app/job/main/activity/conf"
	"go-common/library/cache/memcache"
	"go-common/library/net/http/blademaster"
)

// Dao .
type Dao struct {
	c                *conf.Config
	client           *blademaster.Client
	mc               *memcache.Pool
	broadcastURL     string
	messageURL       string
	timeFinishExpire int32
	lessTimeExpire   int32
}

// New .
func New(c *conf.Config) *Dao {
	d := &Dao{
		c:      c,
		client: blademaster.NewClient(c.HTTPClient),
		mc:     memcache.NewPool(c.Memcache.Like),
	}
	d.broadcastURL = d.c.Host.APICo + _broadURL
	d.messageURL = d.c.Host.MsgCo + _messageURL
	d.timeFinishExpire = int32(time.Duration(c.Memcache.TimeFinishExpire) / time.Second)
	d.lessTimeExpire = int32(time.Duration(c.Memcache.LessTimeExpire) / time.Second)
	return d
}
