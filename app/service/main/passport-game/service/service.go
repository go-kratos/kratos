package service

import (
	"context"
	"crypto/rsa"

	"go-common/app/service/main/passport-game/conf"
	"go-common/app/service/main/passport-game/dao"
	"go-common/app/service/main/passport-game/model"
	"go-common/library/log"
	"go-common/library/stat"
	"go-common/library/stat/prom"
)

// Service service.
type Service struct {
	c                  *conf.Config
	d                  *dao.Dao
	missch             chan func()
	appMap             map[string]*model.App
	originPubKey       *rsa.PublicKey
	cloudRSAKey        *RSAKey
	proxy              bool
	currentRegion      string
	oauth              map[string]string
	renewToken         map[string]string
	regionItems        []*model.RegionInfo
	dispatcherErrStats stat.Stat
}

// RSAKey RSA public and private key.
type RSAKey struct {
	pub  *rsa.PublicKey
	priv *rsa.PrivateKey
}

// New new a service instance.
func New(c *conf.Config) (s *Service) {
	d := dao.New(c)
	// load apps
	apps, err := d.Apps(context.TODO())
	if err != nil {
		panic(err)
	}
	if len(apps) == 0 {
		panic("apps loaded from db is empty, abort")
	}
	appMap := make(map[string]*model.App, len(apps))
	for _, app := range apps {
		appMap[app.AppKey] = app
	}
	// load origin public key
	key, err := d.RSAKeyOrigin(context.TODO())
	if err != nil {
		panic(err)
	}
	originPub, err := pubKey([]byte(key.Key))
	if err != nil {
		panic(err)
	}
	// load origin public and private key for cloud
	pub, err := pubKey([]byte(_originPubKey))
	if err != nil {
		panic(err)
	}
	priv, err := privKey([]byte(_originPrivKey))
	if err != nil {
		panic(err)
	}
	// init
	s = &Service{
		c:            c,
		d:            d,
		missch:       make(chan func(), 10240),
		appMap:       appMap,
		originPubKey: originPub,
		cloudRSAKey: &RSAKey{
			pub:  pub,
			priv: priv,
		},
		proxy:              c.Proxy,
		currentRegion:      c.Dispatcher.Name,
		oauth:              c.Dispatcher.Oauth,
		renewToken:         c.Dispatcher.RenewToken,
		regionItems:        c.Dispatcher.RegionInfos,
		dispatcherErrStats: prom.BusinessErrCount,
	}
	go s.cacheproc()
	return
}

// APP query appkey
func (s *Service) APP(app string) (v *model.App, ok bool) {
	v, ok = s.appMap[app]
	return
}

func (s *Service) addCache(f func()) {
	select {
	case s.missch <- f:
	default:
		log.Warn("cache chan full")
	}
}

// cacheproc is a routine for executing closure.
func (s *Service) cacheproc() {
	for {
		f := <-s.missch
		f()
	}
}

// Ping check server ok.
func (s *Service) Ping(c context.Context) (err error) {
	err = s.d.Ping(c)
	return
}

// Close close service, including closing dao.
func (s *Service) Close() (err error) {
	s.d.Close()
	return
}
