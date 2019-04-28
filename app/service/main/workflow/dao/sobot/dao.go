package sobot

import (
	"go-common/app/service/main/workflow/conf"
	bm "go-common/library/net/http/blademaster"
)

// Dao it
type Dao struct {
	c               *conf.Config
	ticketInfoURL   string
	ticketAddURL    string
	ticketModifyURL string
	replyAddURL     string
	// sobot httpclient
	httpSobot *bm.Client
}

// New Dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:               c,
		httpSobot:       bm.NewClient(c.HTTPClient.Sobot),
		ticketInfoURL:   c.Host.ServiceURI + _sobotTicketInfoURL,
		ticketAddURL:    c.Host.ServiceURI + _sobotAddTicketURL,
		ticketModifyURL: c.Host.ServiceURI + _sobotTicketModifyURL,
		replyAddURL:     c.Host.ServiceURI + _sobotAddReplyURL,
	}
	return
}
