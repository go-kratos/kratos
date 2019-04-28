package pgc

import (
	"fmt"
	"strconv"

	"go-common/app/job/main/tv/conf"
	"go-common/app/job/main/tv/dao/lic"
	model "go-common/app/job/main/tv/model/pgc"
)

var categories = map[int8]string{
	1: "番剧",
	2: "电影",
	3: "纪录片",
	4: "国漫",
	5: "电视剧",
}

var zones = map[int64]string{
	1: "中国",
	2: "日本",
}

const _zoneNotFound = "其他"

// newLic create the skeleton of the license struct
func newLic(Season *model.TVEpSeason, conf *conf.Sync) *model.License {
	// one license stryct oer season
	var (
		ps   []*model.PS
		sign = conf.Sign
		area string
		ok   bool
	)
	if areaInt, _ := strconv.ParseInt(Season.Area, 10, 64); areaInt != 0 { //compatible with old version ( area was int )
		if area, ok = zones[areaInt]; !ok {
			area = _zoneNotFound
		}
	} else { // new logic, directly transform
		area = Season.Area
	}
	var programS = &model.PS{
		ProgramSetID:     conf.AuditPrefix + fmt.Sprintf("%d", Season.ID),
		ProgramSetName:   Season.Title,
		ProgramSetClass:  Season.Style,
		ProgramSetType:   categories[Season.Category],
		PublishDate:      Season.PlayTime.Time().Format("2006-01-02"),
		Copyright:        Season.Copyright,
		ProgramCount:     int(Season.TotalNum),
		CREndData:        "1970-01-01",
		DefinitionType:   "SD",
		CpCode:           conf.LConf.CPCode,
		PayStatus:        Season.Status,
		PrimitiveName:    Season.OriginName,
		Alias:            Season.Alias,
		Zone:             area,
		LeadingRole:      Season.Role,
		ProgramSetDesc:   Season.Desc,
		Staff:            Season.Staff,
		ProgramSetPoster: Season.Cover,
		ProgramList:      &model.ProgramList{},
		Producer:         Season.Producer,
		SubGenre:         Season.Version,
	}
	ps = append(ps, programS)
	return lic.BuildLic(sign, ps, 0)
}
