package dao

import (
	"context"
	"crypto/tls"
	"errors"
	"time"

	"go-common/app/service/main/push/dao/apns2"
	"go-common/app/service/main/push/dao/fcm"
	"go-common/app/service/main/push/dao/huawei"
	"go-common/app/service/main/push/dao/jpush"
	"go-common/app/service/main/push/dao/mi"
	"go-common/app/service/main/push/dao/oppo"
	"go-common/app/service/main/push/model"
	"go-common/library/conf/env"
	"go-common/library/log"
)

var errNoClinets = errors.New("no clients")

func (d *Dao) loadClients() {
	var cnt int
	for cnt < 3 {
		auths, err := d.auths(context.Background())
		if err != nil {
			log.Error("d.auths() error(%v)", err)
			time.Sleep(time.Second)
			cnt++
			continue
		}
		if len(auths) == 0 {
			return
		}
		for _, a := range auths {
			log.Info("new push client. app(%d) platform(%d)", a.APPID, a.PlatformID)
			i := fmtRoundIndex(a.APPID, a.PlatformID)
			d.clientsIndex[i] = new(uint32)
			switch a.PlatformID {
			case model.PlatformIPhone:
				d.clientsIPhone[a.APPID] = d.newApnsClients(model.PlatformIPhone, a.Value, a.Key, a.BundleID)
				d.clientsLen[i] = len(d.clientsIPhone[a.APPID])
			case model.PlatformIPad:
				d.clientsIPad[a.APPID] = d.newApnsClients(model.PlatformIPad, a.Value, a.Key, a.BundleID)
				d.clientsLen[i] = len(d.clientsIPad[a.APPID])
			case model.PlatformHuawei:
				cs := d.newHuaweiClients(a.APPID, a.BundleID)
				if len(cs) > 0 {
					d.clientsHuawei[a.APPID] = cs
					d.clientsLen[i] = len(d.clientsHuawei)
				}
			case model.PlatformOppo:
				cs := d.newOppoClients(a.APPID, a.BundleID)
				if len(cs) > 0 {
					d.clientsOppo[a.APPID] = cs
					d.clientsLen[i] = len(d.clientsOppo)
				}
			case model.PlatformXiaomi:
				d.clientsMi[a.APPID] = d.newMiClients(a.Key, a.Value)
				d.clientsLen[i] = len(d.clientsMi[a.APPID])
				d.clientMiByMids[a.APPID] = d.newMiClientByMids(a.Key, a.Value)
			case model.PlatformJpush:
				d.clientsJpush[a.APPID] = d.newJpushClients(a.Key, a.Value)
				d.clientsLen[i] = len(d.clientsJpush[a.APPID])
			case model.PlatformFCM:
				d.clientsFCM[a.APPID] = d.newFcmClients(a.Key)
				d.clientsLen[i] = len(d.clientsFCM[a.APPID])
			default:
				log.Warn("unknown platform (%+v)", a)
			}
		}
		return
	}
}

func (d *Dao) newMiClients(pkg, auth string) (cs []*mi.Client) {
	for i := 0; i < d.c.Android.PoolSize; i++ {
		c := mi.NewClient(pkg, auth, time.Duration(d.c.Android.Timeout))
		if env.DeployEnv == env.DeployEnvDev {
			c.SetDevelopmentURL(mi.RegURL)
		} else {
			if d.c.Android.MiUseVip == model.SwitchOn {
				c.SetVipURL(mi.RegURL)
			} else {
				c.SetProductionURL(mi.RegURL)
			}
		}
		cs = append(cs, c)
	}
	return
}

func (d *Dao) newMiClientByMids(pkg, auth string) (c *mi.Client) {
	c = mi.NewClient(pkg, auth, time.Duration(d.c.Android.Timeout))
	if env.DeployEnv == env.DeployEnvDev {
		c.SetDevelopmentURL(mi.AccountURL)
	} else {
		if d.c.Android.MiUseVip == model.SwitchOn {
			c.SetVipURL(mi.RegURL)
		} else {
			c.SetProductionURL(mi.RegURL)
		}
	}
	return
}

func (d *Dao) newApnsClients(platform int, cert, key, bundleID string) (res []*apns2.Client) {
	var (
		err         error
		certificate tls.Certificate
	)
	if certificate, err = tls.X509KeyPair([]byte(cert), []byte(key)); err != nil {
		log.Error("tls.X509KeyPair(%s,%s) error(%v)", cert, key, err)
		PromError("client:加载证书")
		return
	}
	poolSize := d.c.Apns.PoolSize
	if platform == model.PlatformIPad {
		poolSize /= 5 // iPad量少，只有iPhone的不到20%
	}
	for i := 0; i < poolSize; i++ {
		var c *apns2.Client
		if env.DeployEnv == env.DeployEnvDev {
			if d.c.Apns.Proxy == model.SwitchOn {
				c = apns2.NewClientWithProxy(certificate, time.Duration(d.c.Apns.Timeout), d.c.Apns.ProxySocket).Development()
			} else {
				c = apns2.NewClient(certificate, time.Duration(d.c.Apns.Timeout)).Development()
			}
		} else {
			if d.c.Apns.Proxy == model.SwitchOn {
				c = apns2.NewClientWithProxy(certificate, time.Duration(d.c.Apns.Timeout), d.c.Apns.ProxySocket).Production()
			} else {
				c = apns2.NewClient(certificate, time.Duration(d.c.Apns.Timeout)).Production()
			}
		}
		c.BoundID = bundleID
		res = append(res, c)
	}
	return
}

func (d *Dao) newHuaweiClients(appid int64, pkg string) (cs []*huawei.Client) {
	retry := _retry
	for retry > 0 {
		if d.huaweiAuth[appid] != nil {
			break
		}
		retry--
		log.Info("retry huawei auth (%d)", retry)
		time.Sleep(3 * time.Second)
	}
	if d.huaweiAuth[appid] == nil {
		log.Error("no huawei auth app(%d)", appid)
		return
	}
	for i := 0; i < d.c.Android.PoolSize; i++ {
		c := huawei.NewClient(pkg, d.huaweiAuth[appid], time.Duration(d.c.Android.Timeout))
		cs = append(cs, c)
	}
	return
}

func (d *Dao) newOppoClients(appid int64, activity string) (cs []*oppo.Client) {
	retry := _retry
	for retry > 0 {
		if d.oppoAuth[appid] != nil {
			break
		}
		retry--
		log.Info("retry oppo auth (%d)", retry)
		time.Sleep(3 * time.Second)
	}
	if d.oppoAuth[appid] == nil {
		log.Error("no oppo auth app(%d)", appid)
		return
	}
	for i := 0; i < d.c.Android.PoolSize; i++ {
		c := oppo.NewClient(d.oppoAuth[appid], activity, time.Duration(d.c.Android.Timeout))
		cs = append(cs, c)
	}
	return
}

func (d *Dao) newJpushClients(appKey, secret string) (cs []*jpush.Client) {
	for i := 0; i < d.c.Android.PoolSize; i++ {
		cs = append(cs, jpush.NewClient(appKey, secret, time.Duration(d.c.Android.Timeout)))
	}
	return
}

func (d *Dao) newFcmClients(key string) (cs []*fcm.Client) {
	for i := 0; i < d.c.Android.PoolSize; i++ {
		cs = append(cs, fcm.NewClient(key, time.Duration(d.c.Android.Timeout)))
	}
	return
}
