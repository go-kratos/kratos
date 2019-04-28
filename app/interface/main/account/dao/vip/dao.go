package vip

import (
	"go-common/app/interface/main/account/conf"
	bm "go-common/library/net/http/blademaster"
)

// Dao dao
type Dao struct {
	c               *conf.Config
	client          *bm.Client
	clientSlow      *bm.Client
	infoURL         string
	codeOpenURL     string
	codeVerifyURL   string
	tipsURL         string
	cancelCouponURL string
	codeOpenedURL   string
	cl              *Clientl
}

// New new
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:               c,
		client:          bm.NewClient(c.HTTPClient.Normal),
		clientSlow:      bm.NewClient(c.HTTPClient.Slow),
		infoURL:         c.Host.Vip + _vipInfo,
		codeOpenURL:     c.Host.API + _vipCodeOpen,
		codeVerifyURL:   c.Host.API + _vipCodeVerify,
		tipsURL:         c.Host.API + _viptips,
		cancelCouponURL: c.Host.API + _couponCancel,
		codeOpenedURL:   c.Host.API + _vipCodeOpened,
	}
	// http client for had url md5 sign.
	d.cl = NewClientl(c.Vipproperty.OAuthClient, d.client)
	return
}
