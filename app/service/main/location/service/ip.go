package service

import (
	"context"
	"errors"
	"net"
	"strconv"

	"go-common/app/service/main/location/model"
	"go-common/library/ecode"
	"go-common/library/log"
	xip "go-common/library/net/ip"
)

// Info get ip info.
func (s *Service) Info(c context.Context, addr string) (res *model.Info, err error) {
	var cityInfo map[string]string
	if cityInfo, err = s.find(addr); err != nil {
		log.Error("%v", err)
		return
	}
	ic, _ := strconv.Atoi(cityInfo["idd_code"])
	la, _ := strconv.ParseFloat(cityInfo["latitude"], 64)
	lo, _ := strconv.ParseFloat(cityInfo["longitude"], 64)
	res = &model.Info{
		Addr:        addr,
		ZoneID:      xip.ZoneID(cityInfo["country_name"], cityInfo["region_name"], cityInfo["city_name"]),
		Country:     cityInfo["country_name"],
		ISP:         cityInfo["isp_domain"],
		Province:    cityInfo["region_name"],
		City:        cityInfo["city_name"],
		Latitude:    la,
		Longitude:   lo,
		CountryCode: ic,
	}
	return
}

// Infos get ip infos.
func (s *Service) Infos(c context.Context, addrs []string) (res map[string]*model.Info, err error) {
	res = make(map[string]*model.Info, len(addrs))
	for _, addr := range addrs {
		var ri *model.Info
		if ri, err = s.Info(c, addr); err != nil {
			log.Error("%v", err)
			res = make(map[string]*model.Info, len(addrs))
			break
		}
		res[addr] = ri
	}
	return
}

// InfoComplete find get whole ip info.
func (s *Service) InfoComplete(c context.Context, addr string) (res *model.InfoComplete, err error) {
	var cityInfo map[string]string
	if cityInfo, err = s.find(addr); err != nil {
		log.Error("%v", err)
		return
	}
	ic, _ := strconv.Atoi(cityInfo["idd_code"])
	la, _ := strconv.ParseFloat(cityInfo["latitude"], 64)
	lo, _ := strconv.ParseFloat(cityInfo["longitude"], 64)
	res = &model.InfoComplete{
		Addr:        addr,
		Country:     cityInfo["country_name"],
		Province:    cityInfo["region_name"],
		City:        cityInfo["city_name"],
		ISP:         cityInfo["isp_domain"],
		Latitude:    la,
		Longitude:   lo,
		CountryCode: ic,
	}
	res.ZoneID = s.zoneIDs(cityInfo["country_name"], cityInfo["region_name"], cityInfo["city_name"])
	return
}

// InfosComplete finds get whole ips infos.
func (s *Service) InfosComplete(c context.Context, ipsStr []string) (res map[string]*model.InfoComplete, err error) {
	res = make(map[string]*model.InfoComplete, len(ipsStr))
	for _, ipStr := range ipsStr {
		var ti *model.InfoComplete
		if ti, err = s.InfoComplete(c, ipStr); err != nil {
			log.Error("%v", err)
			res = make(map[string]*model.InfoComplete, len(ipsStr))
			break
		}
		res[ipStr] = ti
	}
	return
}

// zoneIDs make zoneids
func (s *Service) zoneIDs(country, region, city string) []int64 {
	cZid := xip.ZoneID(country, "", "")
	cpZid := xip.ZoneID(country, region, "")
	cpcZid := xip.ZoneID(country, region, city)
	zoneids := []int64{0, cZid, cpZid, cpcZid}
	return zoneids
}

func (s *Service) find(addr string) (cityInfo map[string]string, err error) {
	ipv := net.ParseIP(addr)
	if ip := ipv.To4(); ip != nil {
		if cityInfo, err = s.v4.FindMap(addr, "CN"); err != nil {
			log.Error("%v", err)
			return
		}
	} else if ip := ipv.To16(); ip != nil {
		if cityInfo, err = s.v6.FindMap(addr, "CN"); err != nil {
			log.Error("%v", err)
			return
		}
	} else {
		err = errors.New("query ip format error")
		return
	}
	if cityInfo == nil {
		err = ecode.NothingFound
		return
	}
	// ex.: from 中国 台湾 花莲市 to 台湾 花莲市 ”“
	if cityInfo["region_name"] == "香港" || cityInfo["region_name"] == "澳门" || cityInfo["region_name"] == "台湾" {
		cityInfo["country_name"] = cityInfo["region_name"]
		cityInfo["region_name"] = cityInfo["city_name"]
		cityInfo["city_name"] = ""
	}
	// ex.: from 中国 中国 ”“ to 中国 ”“ ”“
	if cityInfo["country_name"] == cityInfo["region_name"] {
		cityInfo["region_name"] = ""
		cityInfo["city_name"] = ""
	} else if cityInfo["region_name"] == cityInfo["city_name"] {
		// ex.: from 中国 北京 北京 to 中国 北京 ”“
		cityInfo["city_name"] = ""
	}
	return
}
