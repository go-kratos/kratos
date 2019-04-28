package grpc

import (
	v1 "go-common/app/service/main/archive/api"
)

func (s *server) archive3ToArc(a *v1.Arc) (arc *v1.Arc) {
	arc = &v1.Arc{
		Aid:         a.Aid,
		Videos:      a.Videos,
		TypeID:      a.TypeID,
		TypeName:    a.TypeName,
		Copyright:   a.Copyright,
		Pic:         a.Pic,
		Title:       a.Title,
		PubDate:     a.PubDate,
		Ctime:       a.Ctime,
		Desc:        a.Desc,
		State:       a.State,
		Access:      a.Access,
		Attribute:   a.Attribute,
		Duration:    a.Duration,
		MissionID:   a.MissionID,
		OrderID:     a.OrderID,
		RedirectURL: a.RedirectURL,
		Forward:     a.Forward,
		Rights: v1.Rights{
			Bp:            a.Rights.Bp,
			Elec:          a.Rights.Elec,
			Download:      a.Rights.Download,
			Movie:         a.Rights.Movie,
			Pay:           a.Rights.Pay,
			HD5:           a.Rights.HD5,
			NoReprint:     a.Rights.NoReprint,
			Autoplay:      a.Rights.Autoplay,
			UGCPay:        a.Rights.UGCPay,
			IsCooperation: a.Rights.IsCooperation,
		},
		Author: v1.Author{
			Mid:  a.Author.Mid,
			Name: a.Author.Name,
			Face: a.Author.Face,
		},
		Stat: v1.Stat{
			Aid:     a.Stat.Aid,
			View:    a.Stat.View,
			Danmaku: a.Stat.Danmaku,
			Reply:   a.Stat.Reply,
			Fav:     a.Stat.Fav,
			Coin:    a.Stat.Coin,
			Share:   a.Stat.Share,
			NowRank: a.Stat.NowRank,
			HisRank: a.Stat.HisRank,
			Like:    a.Stat.Like,
			DisLike: a.Stat.DisLike,
		},
		ReportResult: a.ReportResult,
		Dynamic:      a.Dynamic,
		FirstCid:     a.FirstCid,
		Dimension: v1.Dimension{
			Width:  a.Dimension.Width,
			Height: a.Dimension.Height,
			Rotate: a.Dimension.Rotate,
		},
	}
	for _, si := range a.StaffInfo {
		if si != nil {
			arc.StaffInfo = append(arc.StaffInfo, &v1.StaffInfo{Mid: si.Mid, Title: si.Title})
		}
	}
	return
}

// ChangeToGrpc is
func (s *server) page3ToPage(v *v1.Page) (vg *v1.Page) {
	vg = &v1.Page{
		Cid:      v.Cid,
		Page:     v.Page,
		From:     v.From,
		Part:     v.Part,
		Duration: v.Duration,
		Vid:      v.Vid,
		Desc:     v.Desc,
		WebLink:  v.WebLink,
		Dimension: v1.Dimension{
			Width:  v.Dimension.Width,
			Height: v.Dimension.Height,
			Rotate: v.Dimension.Rotate,
		},
	}
	return
}
