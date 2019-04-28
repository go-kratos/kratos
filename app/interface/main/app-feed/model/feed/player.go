package feed

import (
	"go-common/app/interface/main/app-feed/model"
	"go-common/app/service/main/archive/model/archive"
	"strconv"
)

func (i *Item) FromPlayerAv(a *archive.ArchiveWithPlayer) {
	if i.Title == "" {
		i.Title = a.Title
	}
	if i.Cover == "" {
		i.Cover = model.CoverURLHTTPS(a.Pic)
	} else {
		i.Cover = model.CoverURLHTTPS(i.Cover)
	}
	i.Param = strconv.FormatInt(a.Aid, 10)
	i.Goto = model.GotoAv
	i.URI = model.FillURI(i.Goto, i.Param, 0, 0, model.AvPlayHandler(a.Archive3, a.PlayerInfo))
	i.Cid = a.FirstCid
	i.Rid = a.TypeID
	i.TName = a.TypeName
	i.Desc = strconv.Itoa(int(a.Stat.Danmaku)) + "弹幕"
	i.fillArcStat(a.Archive3)
	i.Duration = a.Duration
	i.Mid = a.Author.Mid
	i.Name = a.Author.Name
	i.Face = a.Author.Face
	i.CTime = a.PubDate
	i.Cid = a.FirstCid
	i.Autoplay = a.Rights.Autoplay
}
