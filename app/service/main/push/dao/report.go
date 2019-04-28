package dao

import (
	"context"
	"time"

	"go-common/app/service/main/push/dao/mi"
	"go-common/app/service/main/push/model"
	"go-common/library/log"
)

// DelMiInvalid .
func (d *Dao) DelMiInvalid(c context.Context) (err error) {
	log.Info("delete xiaomi invalid report start, apps(%d)", len(d.clientsMi))
	for appid, clients := range d.clientsMi {
		log.Info("clients info app(%d) len(%d)", appid, len(clients))
		if len(clients) == 0 {
			log.Warn("no clients app(%d)", appid)
			continue
		}
		log.Info("del mi invalid start, app(%d)", appid)
		var res *mi.Response
		if res, err = clients[0].InvalidTokens(); err != nil {
			log.Error("client.InvalidTokens() error(%v)", err)
			PromError("report:获取小米无效token")
			continue
		}
		if res == nil || len(res.Data.List) == 0 {
			log.Warn("no tokens app(%d)", appid)
			continue
		}
		if err = d.delInvalidMiReports(c, appid, res.Data.List); err != nil {
			PromError("report:主动删除xiaomi无效上报")
			continue
		}
		log.Info("already del mi invalid stop, app(%d) count(%d)", appid, len(res.Data.List))
	}
	PromInfo("report:主动删除xiaomi无效上报")
	return
}

// DelMiUninstalled deletes mi uninstalled tokens.
func (d *Dao) DelMiUninstalled(c context.Context) (err error) {
	log.Info("delete xiaomi uninstalled tokens start, apps(%d)", len(d.clientsMi))
	for appid, clients := range d.clientsMi {
		log.Info("clients info app(%d) len(%d)", appid, len(clients))
		if len(clients) == 0 {
			log.Warn("no clients app(%d)", appid)
			continue
		}
		log.Info("del mi uninstalled tokens start, app(%d)", appid)
		var res *mi.UninstalledResponse
		if res, err = clients[0].UninstalledTokens(); err != nil {
			log.Error("client.UninstalledTokens() error(%v)", err)
			PromError("report:获取小米卸载token")
			continue
		}
		if res.Code == mi.ResultCodeNoMsgInEmq {
			log.Info("no tokens app(%d)", appid)
			continue
		}
		if res.Code != mi.ResultCodeOk {
			log.Error("get uninstalled tokens error resp(%+v)", res)
			continue
		}
		if len(res.Data) == 0 {
			log.Warn("no tokens app(%d)", appid)
			continue
		}
		if err = d.delInvalidMiReports(c, appid, res.Data); err != nil {
			PromError("report:主动删除xiaomi卸载token")
			continue
		}
		log.Info("already del mi uninstalled stop, app(%d) count(%d)", appid, len(res.Data))
	}
	PromInfo("report:主动删除xiaomi卸载token")
	return
}

func (d *Dao) delInvalidMiReports(c context.Context, appid int64, tokens []string) (err error) {
	var rs []*model.Report
	if rs, err = d.Reports(c, tokens); err != nil {
		log.Error("d.Reports(%v) error(%v)", tokens, err)
		return
	} else if len(rs) == 0 {
		log.Warn("reports can not be found by tokens(%d)", len(tokens))
		log.Warn("reports can not be found by tokens(%v)", tokens)
		return
	}
	for _, r := range rs {
		log.Info("deleted invalid mi report, app(%d) mid(%d) token(%s)", appid, r.Mid, r.DeviceToken)
		var (
			i     int
			e     error
			retry = _retry
		)
		for i < retry {
			if _, e = d.DelReport(c, r.DeviceToken); e == nil {
				break
			}
			time.Sleep(time.Second)
			i++
			log.Warn("retry delete report, mid(%d) token(%s)", r.Mid, r.DeviceToken)
		}
		if e != nil || r.Mid <= 0 {
			continue
		}
		i = 0
		for i < retry {
			if e = d.DelReportCache(c, r.Mid, r.APPID, r.DeviceToken); e == nil {
				break
			}
			log.Warn("retry delete report cache, mid(%d) token(%s)", r.Mid, r.DeviceToken)
			time.Sleep(time.Second)
			i++
		}
	}
	log.Info("del invalid mi report, app(%d) tokens(%d)", appid, len(rs))
	return
}
