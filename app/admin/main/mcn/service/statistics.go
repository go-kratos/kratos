package service

import (
	"context"
	"fmt"
	"time"

	"go-common/app/admin/main/mcn/model"
	dtmdl "go-common/app/interface/main/mcn/model/datamodel"
	accgrpc "go-common/app/service/main/account/api"
	arcgrpc "go-common/app/service/main/archive/api"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
	xtime "go-common/library/time"

	"github.com/pkg/errors"
)

// ArcTopDataStatistics .
func (s *Service) ArcTopDataStatistics(c context.Context, arg *model.McnGetRankReq) (res *model.McnGetRankUpFansReply, err error) {
	return s.dao.ArcTopDataStatistics(c, arg)
}

// McnsTotalDatas .
func (s *Service) McnsTotalDatas(c context.Context, arg *model.TotalMcnDataReq) (res *model.TotalMcnDataInfo, err error) {
	var (
		signUpsTotal, FanTotal, videoupTotal, playTotal           int64
		mids, fansMids, ArcsMids, avids, tids, arcsTids, typeTids []int64
		m                                                         *model.McnDataOverview
		mrf                                                       map[int8][]*model.McnRankFansOverview
		ras                                                       []*model.McnRankArchiveLikesOverview
		mmd                                                       map[string][]*model.McnDataTypeSummary
	)
	td := new(thirdDataMap)
	res = new(model.TotalMcnDataInfo)
	res.TopInfo = new(model.McnDataTopInfo)
	res.TypesInfo = new(model.McnDataTypesInfo)
	date := xtime.Time(time.Date(arg.Date.Time().Year(), arg.Date.Time().Month(), arg.Date.Time().Day()-1, 0, 0, 0, 0, time.Local).Unix())
	if m, err = s.dao.McnDataOverview(c, date); err != nil {
		return
	}
	res.BaseInfo = m
	if mrf, fansMids, err = s.dao.McnRankFansOverview(c, model.DataTypeDay, date, model.TopDataLenth); err != nil {
		return
	}
	if ras, ArcsMids, avids, arcsTids, err = s.dao.McnRankArchiveLikesOverview(c, model.DataTypeDay, date, model.TopDataLenth); err != nil {
		return
	}
	if mmd, typeTids, err = s.dao.McnDataTypeSummary(c, date); err != nil {
		return
	}
	mids = append(mids, fansMids...)
	mids = append(mids, ArcsMids...)
	tids = append(tids, arcsTids...)
	tids = append(tids, typeTids...)
	if td.infosReply, err = s.accGRPC.Infos3(c, &accgrpc.MidsReq{Mids: mids}); err != nil {
		log.Error("s.accGRPC.Infos3(%+v) error(%+v)", mids, err)
		err = nil
	}
	if td.archivesReply, err = s.arcGRPC.Arcs(c, &arcgrpc.ArcsRequest{Aids: avids}); err != nil {
		log.Error("s.arcGRPC.Arcs(%+v) error(%+v)", avids, err)
		err = nil
	}
	td.tpNames = s.videoup.GetTidName(tids)
	for view, data := range mrf {
		for _, v := range data {
			var (
				rateIncr      int64
				name, _, _, _ = td.getTypeName(v.Mid, 0, 0, 0)
			)
			if v.Fans != 0 {
				rateIncr = (10000 * v.FansIncr) / v.Fans
			}
			fr := &model.FansRankIncr{Mid: v.Mid, Name: name, Rank: v.Rank, FansIncr: v.FansIncr, Fans: v.Fans, RateIncr: rateIncr, SignID: v.SignID}
			switch model.DataViewFansTop(view) {
			case model.McnFansIncr:
				res.TopInfo.McnFansIncr = append(res.TopInfo.McnFansIncr, fr)
			case model.McnFansIncrRate:
				res.TopInfo.McnFansRateIncr = append(res.TopInfo.McnFansRateIncr, fr)
			case model.UpFansIncr:
				res.TopInfo.UpFansIncr = append(res.TopInfo.UpFansIncr, fr)
			case model.UpFansIncrRate:
				res.TopInfo.UpFansRateIncr = append(res.TopInfo.UpFansRateIncr, fr)
			}
		}
	}
	for _, v := range ras {
		var mcnName, upName, avTitle, tpName = td.getTypeName(v.McnMid, v.UpMid, v.Avid, int64(v.Tid))
		lr := &model.LikesRankIncr{McnMid: v.McnMid, McnName: mcnName,
			UpMid: v.UpMid, UpName: upName, AVID: v.Avid, AVTitle: avTitle,
			TID: v.Tid, TypeName: tpName, LikesIncr: v.Likes, PlayIncr: v.Plays, SignID: v.SignID}
		res.TopInfo.ArcLikesIncr = append(res.TopInfo.ArcLikesIncr, lr)
	}
	for vt, data := range mmd {
		for _, v := range data {
			var _, _, _, tpName = td.getTypeName(0, 0, 0, int64(v.Tid))
			dt := &model.DataTypes{TID: v.Tid, TypeName: tpName, Amount: v.Amount}
			switch vt {
			case fmt.Sprintf("%d-%d", model.SignUpsAccumulate, model.DataTypeAccumulate):
				signUpsTotal += v.Amount
				res.TypesInfo.SignUps = append(res.TypesInfo.SignUps, dt)
			case fmt.Sprintf("%d-%d", model.FansIncr, model.DataTypeDay):
				FanTotal += v.Amount
				res.TypesInfo.FansIncr = append(res.TypesInfo.FansIncr, dt)
			case fmt.Sprintf("%d-%d", model.VideoUpsIncr, model.DataTypeDay):
				videoupTotal += v.Amount
				res.TypesInfo.VideoupIncr = append(res.TypesInfo.VideoupIncr, dt)
			case fmt.Sprintf("%d-%d", model.PlaysIncr, model.DataTypeDay):
				playTotal += v.Amount
				res.TypesInfo.PlayIncr = append(res.TypesInfo.PlayIncr, dt)
			}
		}
	}
	statTypeRate(res.TypesInfo.SignUps, signUpsTotal)
	statTypeRate(res.TypesInfo.FansIncr, FanTotal)
	statTypeRate(res.TypesInfo.VideoupIncr, videoupTotal)
	statTypeRate(res.TypesInfo.PlayIncr, playTotal)
	return
}

type thirdDataMap struct {
	tpNames       map[int64]string
	infosReply    *accgrpc.InfosReply
	archivesReply *arcgrpc.ArcsReply
}

func (td *thirdDataMap) getTypeName(mcnMid, upMid, avid, tid int64) (mcnName, upName, avTitle, tpName string) {
	if td.infosReply != nil {
		infos := td.infosReply.Infos
		if info, ok := infos[mcnMid]; ok {
			mcnName = info.Name
		}
		if info, ok := infos[upMid]; ok {
			upName = info.Name
		}
	}
	if td.archivesReply != nil {
		archives := td.archivesReply.Arcs
		if archive, ok := archives[avid]; ok {
			avTitle = archive.Title
		}
	}
	tpName = td.tpNames[tid]
	return
}

func statTypeRate(dts []*model.DataTypes, total int64) {
	for _, v := range dts {
		if total != 0 {
			v.Rate = (10000 * v.Amount) / total
		}
		v.Total = total
	}
}

// McnFansAnalyze .
func (s *Service) McnFansAnalyze(c context.Context, arg *model.McnCommonReq) (res *model.McnGetMcnFansReply, err error) {
	var (
		f       *dtmdl.DmConMcnFansD
		sex     *dtmdl.DmConMcnFansSexW
		age     *dtmdl.DmConMcnFansAgeW
		playWay *dtmdl.DmConMcnFansPlayWayW
		areas   []*dtmdl.DmConMcnFansAreaW
		types   []*dtmdl.DmConMcnFansTypeW
		tags    []*dtmdl.DmConMcnFansTagW
		g       errgroup.Group
	)
	res = new(model.McnGetMcnFansReply)
	g.Go(func() (err error) {
		if f, err = s.dao.DataFans(c, arg); err != nil {
			log.Error("s.dao.DataFans(%+v) error(%v)", arg, err)
			return
		}
		res.FansOverview = f
		return
	})
	g.Go(func() (err error) {
		if sex, age, playWay, err = s.dao.DataFansBaseAttr(c, arg); err != nil {
			log.Error("s.dao.DataFansBaseAttr(%+v) error(%v)", arg, err)
			return
		}
		res.FansSex = sex
		res.FansAge = age
		res.FansPlayWay = playWay
		return
	})
	g.Go(func() (err error) {
		if areas, err = s.dao.DataFansArea(c, arg); err != nil {
			log.Error("s.dao.DataFansArea(%+v) error(%v)", arg, err)
			return
		}
		res.FansArea = areas
		return
	})
	g.Go(func() (err error) {
		if types, err = s.dao.DataFansType(c, arg); err != nil {
			log.Error("s.dao.DataFansType(%+v) error(%v)", arg, err)
			return
		}
		res.FansType = types
		return
	})
	g.Go(func() (err error) {
		if tags, err = s.dao.DataFansTag(c, arg); err != nil {
			log.Error("s.dao.DataFansTag(%+v) error(%v)", arg, err)
			return
		}
		res.FansTag = tags
		return
	})
	if err = g.Wait(); err != nil {
		err = errors.WithStack(err)
	}
	return
}
