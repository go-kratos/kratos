package bfs

import (
	"go-common/app/interface/main/mcn/conf"
)

// Dao dao
type Dao struct {
	c      *conf.Config
	bucket string
	key    string
	secret string
	bfs    string
}

// New init mysql db
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:      c,
		bucket: c.BFS.Bucket,
		key:    c.BFS.Key,
		secret: c.BFS.Secret,
		bfs:    c.Host.Bfs,
	}
	return
}
