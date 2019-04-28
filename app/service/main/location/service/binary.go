package service

import (
	"fmt"
	"io/ioutil"
	"strconv"

	"go-common/library/net/ip"

	"github.com/pkg/errors"
)

const (
	_HeaderInfoSize = 24
)

// HeadIndex index pos.
type HeadIndex struct {
	Index    int
	Country  int
	Province int
	City     int
	ISP      int
	District int
}

// ZoneInfo all kinds of zone info.
type ZoneInfo struct {
	code int
	name string
}

func formByte(data []byte, pos int) (res uint32) {
	res = (uint32(data[pos]&0xff) << 8) + uint32(data[pos+1]&0xff)
	return
}

func formByte2(data []byte, pos int) (res int32) {
	res = (int32(data[pos]&0xff) << 24) + (int32(data[pos+1]&0xff) << 16) + (int32(data[pos+2]&0xff) << 8) + int32(data[pos+3]&0xff)
	return
}

// newBin new binary ip.
func (s *Service) newBin(path string) (list *ip.List, err error) {
	var (
		contents  []byte
		countrys  map[int]*ZoneInfo
		provinces map[int]*ZoneInfo
		citys     map[int]*ZoneInfo
		isps      map[int]*ZoneInfo
		districts map[int]*ZoneInfo
	)
	if contents, err = ioutil.ReadFile(path); err != nil {
		err = errors.WithStack(err)
		return
	}
	if len(contents) < _HeaderInfoSize {
		err = errors.New("文件异常")
		return
	}
	index := s.loadHeader(contents)
	if list, err = s.loadIP(contents, index); err != nil {
		err = errors.WithStack(err)
		return
	}
	if countrys, err = s.loadCounty(contents, index); err != nil {
		err = errors.WithStack(err)
		return
	}
	if provinces, err = s.loadProvince(contents, index); err != nil {
		err = errors.WithStack(err)
		return
	}
	if citys, err = s.loadCity(contents, index); err != nil {
		err = errors.WithStack(err)
		return
	}
	if isps, err = s.loadISP(contents, index); err != nil {
		err = errors.WithStack(err)
		return
	}
	if districts, err = s.loadDistrict(contents, index); err != nil {
		err = errors.WithStack(err)
		return
	}
	var (
		zone *ZoneInfo
		ok   bool
	)
	for _, info := range list.IPs {
		if zone, ok = countrys[info.CountryCode]; ok {
			info.Country = zone.name
		}
		if zone, ok = provinces[info.ProvinceCode]; ok {
			info.Province = zone.name
		}
		if zone, ok = citys[info.CityCode]; ok {
			info.City = zone.name
		}
		if info.Province == "香港" || info.Province == "澳门" || info.Province == "台湾" {
			info.Country = info.Province
			info.Province = info.City
			info.City = ""
		}
		if info.Country == info.Province {
			info.Province = ""
			info.City = ""
		}
		if info.Province == info.City {
			info.City = ""
		}
		if zone, ok = isps[info.ISPCode]; ok {
			info.ISP = zone.name
		}
		if zone, ok = districts[info.DistrictCode]; ok {
			info.District = zone.name
		}
	}
	return
}

func (s *Service) loadHeader(data []byte) (index *HeadIndex) {
	var pos = 0
	index = &HeadIndex{}
	index.Index = int(formByte2(data, pos))
	pos += 4
	index.Country = int(formByte2(data, pos))
	pos += 4
	index.Province = int(formByte2(data, pos))
	pos += 4
	index.City = int(formByte2(data, pos))
	pos += 4
	index.ISP = int(formByte2(data, pos))
	pos += 4
	index.District = int(formByte2(data, pos))
	return
}

func (s *Service) loadIP(data []byte, index *HeadIndex) (list *ip.List, err error) {
	var (
		pos   = index.Index
		count = int(formByte2(data, pos))
	)
	pos += 4
	list = new(ip.List)
	for i := 0; i < count; i++ {
		if pos > len(data) {
			break
		}
		item := &ip.IP{}
		item.Begin = uint32(formByte2(data, pos))
		pos += 4
		item.End = uint32(formByte2(data, pos))
		pos += 4
		item.CountryCode = int(formByte(data, pos))
		pos += 2
		item.ProvinceCode = int(formByte(data, pos))
		pos += 2
		item.CityCode = int(formByte(data, pos))
		pos += 2
		item.ISPCode = int(formByte(data, pos))
		pos += 2
		item.DistrictCode = int(formByte(data, pos))
		pos += 2
		latitude, _ := strconv.ParseFloat(fmt.Sprintf("%.4f", float64(formByte2(data, pos))/float64(10000)), 64)
		item.Latitude = latitude
		pos += 4
		longitude, _ := strconv.ParseFloat(fmt.Sprintf("%.4f", float64(formByte2(data, pos))/float64(10000)), 64)
		item.Longitude = longitude
		pos += 4
		list.IPs = append(list.IPs, item)
	}
	if len(list.IPs) != count {
		err = errors.New("loadIP error")
	}
	return
}

func (s *Service) loadCounty(data []byte, index *HeadIndex) (countrys map[int]*ZoneInfo, err error) {
	var (
		pos   = index.Country
		count = int(formByte(data, pos))
	)
	pos += 2
	countrys = make(map[int]*ZoneInfo)
	for i := 0; i < count; i++ {
		if pos > len(data) {
			break
		}
		item := new(ZoneInfo)
		item.code = int(formByte(data, pos))
		pos += 2
		lenght := int(int32(data[pos]))
		pos++
		item.name = string(data[pos : pos+lenght])
		pos += lenght
		countrys[item.code] = item
	}
	if len(countrys) != count {
		err = errors.New("loadCounty error")
		return
	}
	return
}

func (s *Service) loadProvince(data []byte, index *HeadIndex) (provinces map[int]*ZoneInfo, err error) {
	var (
		pos   = index.Province
		count = int(formByte(data, pos))
	)
	pos += 2
	provinces = make(map[int]*ZoneInfo)
	for i := 0; i < count; i++ {
		if pos > len(data) {
			break
		}
		item := new(ZoneInfo)
		item.code = int(formByte(data, pos))
		pos += 2
		lenght := int(int32(data[pos]))
		pos++
		item.name = string(data[pos : pos+lenght])
		pos += lenght
		provinces[item.code] = item
	}
	if len(provinces) != count {
		err = errors.New("loadProvince error")
		return
	}
	return
}

func (s *Service) loadCity(data []byte, index *HeadIndex) (citys map[int]*ZoneInfo, err error) {
	var (
		pos   = index.City
		count = int(formByte(data, pos))
	)
	pos += 2
	citys = make(map[int]*ZoneInfo)
	for i := 0; i < count; i++ {
		if pos > len(data) {
			break
		}
		item := new(ZoneInfo)
		item.code = int(formByte(data, pos))
		pos += 2
		lenght := int(int32(data[pos]))
		pos++
		item.name = string(data[pos : pos+lenght])
		pos += lenght
		citys[item.code] = item
	}
	if len(citys) != count {
		err = errors.New("loadCity error")
		return
	}
	return
}

func (s *Service) loadISP(data []byte, index *HeadIndex) (isps map[int]*ZoneInfo, err error) {
	var (
		pos   = index.ISP
		count = int(formByte(data, pos))
	)
	pos += 2
	isps = make(map[int]*ZoneInfo)
	for i := 0; i < count; i++ {
		if pos > len(data) {
			break
		}
		item := new(ZoneInfo)
		item.code = int(formByte(data, pos))
		pos += 2
		lenght := int(int32(data[pos]))
		pos++
		item.name = string(data[pos : pos+lenght])
		pos += lenght
		isps[item.code] = item
	}
	if len(isps) != count {
		err = errors.New("loadISP error")
		return
	}
	return
}

func (s *Service) loadDistrict(data []byte, index *HeadIndex) (districts map[int]*ZoneInfo, err error) {
	var (
		pos   = index.District
		count = int(formByte(data, pos))
	)
	pos += 2
	districts = make(map[int]*ZoneInfo)
	for i := 0; i < count; i++ {
		if pos > len(data) {
			break
		}
		item := new(ZoneInfo)
		item.code = int(formByte(data, pos))
		pos += 2
		lenght := int(int32(data[pos]))
		pos++
		item.name = string(data[pos : pos+lenght])
		pos += lenght
		districts[item.code] = item
	}
	if len(districts) != count {
		err = errors.New("loadDistrict error")
		return
	}
	return
}
