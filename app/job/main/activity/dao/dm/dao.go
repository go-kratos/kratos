package dm

import (
	"go-common/app/job/main/activity/conf"
	bm "go-common/library/net/http/blademaster"
)

// Dao dao struct.
type Dao struct {
	// http
	broadcastURL string
	httpCli      *bm.Client
}

// New return dm dao instance.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		broadcastURL: "http://api.bilibili.co/x/internal/chat/push/room",
		httpCli:      bm.NewClient(c.HTTPClient),
	}
	return
}
