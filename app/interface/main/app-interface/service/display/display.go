package display

import (
	"context"
	"go-common/library/net/metadata"
	"strconv"
	"time"

	locmdl "go-common/app/service/main/location/model"
	"go-common/library/log"
	xip "go-common/library/net/ip"
)

// DisplayID is display id .
func (s *Service) DisplayID(c context.Context, mid int64, buvid string, now time.Time) (id string) {
	if mid == 0 {
		id = buvid + "-" + strconv.FormatInt(now.Unix(), 10)
	} else {
		id = strconv.FormatInt(mid, 10) + "-" + strconv.FormatInt(now.Unix(), 10)
	}
	return
}

// Zone is zone id and district info .
func (s *Service) Zone(c context.Context, now time.Time) (zone *xip.Zone) {
	var (
		info *locmdl.Info
		err  error
	)
	zone = &xip.Zone{}
	if info, err = s.loc.Info(c, metadata.String(c, metadata.RemoteIP)); err != nil || info == nil {
		log.Error("error %v or info is nil", err)
		return
	}
	zone.ID = info.ZoneID
	zone.Addr = info.Addr
	zone.ISP = info.ISP
	zone.Country = info.Country
	zone.Province = info.Province
	zone.City = info.City
	zone.Latitude = info.Latitude
	zone.Longitude = info.Longitude
	zone.CountryCode = info.CountryCode
	return
}
