package service

import (
	"context"
	"sort"
	"strconv"

	"go-common/app/service/main/location/model"
	"go-common/library/log"
	xip "go-common/library/net/ip"
)

// TmpInfo get ip info.
func (s *Service) TmpInfo(addr string) (ti *model.TmpInfo, err error) {
	var cityInfo map[string]string
	if cityInfo, err = s.find(addr); err != nil {
		log.Error("%v", err)
		return
	}
	ti = &model.TmpInfo{
		Addr:     addr,
		ZoneID:   xip.ZoneID(cityInfo["country_name"], cityInfo["region_name"], cityInfo["city_name"]),
		Country:  cityInfo["country_name"],
		Province: cityInfo["region_name"],
		City:     cityInfo["city_name"],
	}
	return
}

// TmpInfos get ip infos.
func (s *Service) TmpInfos(c context.Context, addrs ...string) (zone []*xip.Zone, err error) {
	addrs = RemoveDuplicatesAndEmpty(addrs)
	zone = make([]*xip.Zone, len(addrs))
	for ide, addr := range addrs {
		var tmp *xip.Zone
		if tmp, err = s.TmpInfo2(c, addr); err != nil {
			log.Error("%v", err)
			zone = make([]*xip.Zone, len(addrs))
			break
		}
		zone[ide] = tmp
	}
	return
}

// TmpInfo2 get ip zone from ip
func (s *Service) TmpInfo2(c context.Context, addr string) (zone *xip.Zone, err error) {
	var cityInfo map[string]string
	if cityInfo, err = s.find(addr); err != nil {
		log.Error("%v", err)
		return
	}
	ic, _ := strconv.Atoi(cityInfo["idd_code"])
	la, _ := strconv.ParseFloat(cityInfo["latitude"], 64)
	lo, _ := strconv.ParseFloat(cityInfo["longitude"], 64)
	zone = &xip.Zone{
		ID:          xip.ZoneID(cityInfo["country_name"], cityInfo["region_name"], cityInfo["city_name"]),
		Addr:        addr,
		ISP:         cityInfo["isp_domain"],
		Country:     cityInfo["country_name"],
		Province:    cityInfo["region_name"],
		City:        cityInfo["city_name"],
		Latitude:    la,
		Longitude:   lo,
		CountryCode: ic,
	}
	return
}

// RemoveDuplicatesAndEmpty string去重 去空
func RemoveDuplicatesAndEmpty(a []string) (ret []string) {
	sort.Strings(a)
	for i := 0; i < len(a); i++ {
		if (i > 0 && a[i-1] == a[i]) || len(a[i]) == 0 {
			continue
		}
		ret = append(ret, a[i])
	}
	return
}
