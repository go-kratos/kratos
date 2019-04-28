package coin

import (
	"context"
	"time"

	"go-common/app/interface/main/app-interface/conf"
	"go-common/app/service/main/archive/api"
	arcrpc "go-common/app/service/main/archive/api/gorpc"
	"go-common/app/service/main/archive/model/archive"
	coinclient "go-common/app/service/main/coin/api"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

// Dao is coin dao
type Dao struct {
	coinClient coinclient.CoinClient
	arcRPC     *arcrpc.Service2
}

// New initial coin dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		arcRPC: arcrpc.New2(c.ArchiveRPC),
	}
	var err error
	if d.coinClient, err = coinclient.NewClient(c.CoinClient); err != nil {
		panic(err)
	}
	return
}

//CoinList coin archive list
func (d *Dao) CoinList(c context.Context, mid int64, pn, ps int) (coinArc []*api.Arc, count int, err error) {
	var (
		coinReply *coinclient.ListReply
		aids      []int64
		arcs      map[int64]*api.Arc
		ip        = metadata.String(c, metadata.RemoteIP)
	)
	coinArc = make([]*api.Arc, 0)
	if coinReply, err = d.coinClient.List(c, &coinclient.ListReq{Mid: mid, Business: "archive", Ts: time.Now().Unix()}); err != nil {
		log.Error("CoinList s.coinClient.List(%d) error(%v)", mid, err)
		err = nil
		return
	}
	existAids := make(map[int64]int64, len(coinReply.List))
	for _, v := range coinReply.List {
		if _, ok := existAids[v.Aid]; ok {
			continue
		}
		aids = append(aids, v.Aid)
		existAids[v.Aid] = v.Aid
	}
	count = len(aids)
	start := (pn - 1) * ps
	end := pn * ps
	switch {
	case start > count:
		aids = aids[:0]
	case end >= count:
		aids = aids[start:]
	default:
		aids = aids[start:end]
	}
	if len(aids) == 0 {
		return
	}
	if arcs, err = d.arcRPC.Archives3(c, &archive.ArgAids2{Aids: aids, RealIP: ip}); err != nil {
		log.Error("CoinList s.arc.Archives3(%v) error(%v)", aids, err)
		err = nil
		return
	}
	for _, aid := range aids {
		if arc, ok := arcs[aid]; ok && arc.IsNormal() {
			if arc.Access >= 10000 {
				arc.Stat.View = 0
			}
			coinArc = append(coinArc, arc)
		}
	}
	return
}
