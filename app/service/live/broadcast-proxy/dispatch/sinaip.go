package dispatch

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"
	"syscall"
)

type SinaIP struct {
	country  map[uint16]*Country
	province map[uint16]*Province
	city     map[uint16]*City
	isp      map[uint16]*ISP
	district map[uint16]*District
	segment  []*IPSegment
}

type IPDetail struct {
	Country   string
	Province  string
	City      string
	ISP       string
	District  string
	Latitude  float64
	Longitude float64
}

type IPSegment struct {
	start     uint32
	end       uint32
	country   *Country
	province  *Province
	city      *City
	isp       *ISP
	district  *District
	latitude  float64
	longitude float64
}

type Country struct {
	name string
}

type Province struct {
	name string
}

type City struct {
	name string
}

type ISP struct {
	name string
}

type District struct {
	name string
}

func NewSinaIP(file string) (*SinaIP, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	data, err := syscall.Mmap(int(f.Fd()), 0, int(info.Size()), syscall.PROT_READ, syscall.MAP_PRIVATE)
	if err != nil {
		return nil, err
	}
	defer syscall.Munmap(data)

	var (
		segmentOffset  uint32
		countryOffset  uint32
		provinceOffset uint32
		cityOffset     uint32
		ispOffset      uint32
		districtOffset uint32
	)
	reader := bytes.NewReader(data[0:24])
	if err := binary.Read(reader, binary.BigEndian, &segmentOffset); err != nil {
		return nil, err
	}
	if err := binary.Read(reader, binary.BigEndian, &countryOffset); err != nil {
		return nil, err
	}
	if err := binary.Read(reader, binary.BigEndian, &provinceOffset); err != nil {
		return nil, err
	}
	if err := binary.Read(reader, binary.BigEndian, &cityOffset); err != nil {
		return nil, err
	}
	if err := binary.Read(reader, binary.BigEndian, &ispOffset); err != nil {
		return nil, err
	}
	if err := binary.Read(reader, binary.BigEndian, &districtOffset); err != nil {
		return nil, err
	}

	biliIP := new(SinaIP)
	if err := biliIP.loadCity(data[cityOffset:]); err != nil {
		return nil, err
	}
	if err := biliIP.loadCountry(data[countryOffset:]); err != nil {
		return nil, err
	}
	if err := biliIP.loadProvince(data[provinceOffset:]); err != nil {
		return nil, err
	}
	if err := biliIP.loadISP(data[ispOffset:]); err != nil {
		return nil, err
	}
	if err := biliIP.loadDistrict(data[districtOffset:]); err != nil {
		return nil, err
	}
	if err := biliIP.loadIPSegment(data[segmentOffset:]); err != nil {
		return nil, err
	}
	return biliIP, nil
}

func (b *SinaIP) DoQuery(ip uint32) *IPDetail {
	left := 0
	right := len(b.segment) - 1
	var r *IPDetail
	for left <= right {
		middle := left + (right-left)/2
		s := b.segment[middle]
		if ip >= s.start && ip <= s.end {
			r = new(IPDetail)
			if s.country != nil {
				r.Country = s.country.name
			}
			if s.province != nil {
				r.Province = s.province.name
			}
			if s.city != nil {
				r.City = s.city.name
			}
			if s.isp != nil {
				r.ISP = s.isp.name
			}
			if s.district != nil {
				r.District = s.district.name
			}
			r.Latitude = s.latitude
			r.Longitude = s.longitude
			break
		} else if ip < s.start {
			right = middle - 1
		} else if ip > s.end {
			left = middle + 1
		}
	}
	return r
}

func (b *SinaIP) loadCountry(data []byte) error {
	reader := bytes.NewReader(data)

	var count uint16
	if err := binary.Read(reader, binary.BigEndian, &count); err != nil {
		return err
	}
	b.country = make(map[uint16]*Country, count)
	for i := uint16(0); i < count; i++ {
		var code uint16
		if err := binary.Read(reader, binary.BigEndian, &code); err != nil {
			return err
		}
		var length uint8
		if err := binary.Read(reader, binary.BigEndian, &length); err != nil {
			return err
		}
		country := make([]byte, length)
		if _, err := io.ReadFull(reader, country); err != nil {
			return err
		}
		b.country[code] = &Country{
			name: string(country),
		}
	}
	return nil
}

func (b *SinaIP) loadProvince(data []byte) error {
	reader := bytes.NewReader(data)

	var count uint16
	if err := binary.Read(reader, binary.BigEndian, &count); err != nil {
		return err
	}
	b.province = make(map[uint16]*Province, count)
	for i := uint16(0); i < count; i++ {
		var code uint16
		if err := binary.Read(reader, binary.BigEndian, &code); err != nil {
			return err
		}
		var length uint8
		if err := binary.Read(reader, binary.BigEndian, &length); err != nil {
			return err
		}
		province := make([]byte, length)
		if _, err := io.ReadFull(reader, province); err != nil {
			return err
		}
		b.province[code] = &Province{
			name: string(province),
		}
	}
	return nil
}

func (b *SinaIP) loadCity(data []byte) error {
	reader := bytes.NewReader(data)

	var count uint16
	if err := binary.Read(reader, binary.BigEndian, &count); err != nil {
		return err
	}
	b.city = make(map[uint16]*City, count)
	for i := uint16(0); i < count; i++ {
		var code uint16
		if err := binary.Read(reader, binary.BigEndian, &code); err != nil {
			return err
		}
		var length uint8
		if err := binary.Read(reader, binary.BigEndian, &length); err != nil {
			return err
		}
		city := make([]byte, length)
		if _, err := io.ReadFull(reader, city); err != nil {
			return err
		}
		b.city[code] = &City{
			name: string(city),
		}
	}
	return nil
}

func (b *SinaIP) loadISP(data []byte) error {
	reader := bytes.NewReader(data)

	var count uint16
	if err := binary.Read(reader, binary.BigEndian, &count); err != nil {
		return err
	}
	b.isp = make(map[uint16]*ISP, count)
	for i := uint16(0); i < count; i++ {
		var code uint16
		if err := binary.Read(reader, binary.BigEndian, &code); err != nil {
			return err
		}
		var length uint8
		if err := binary.Read(reader, binary.BigEndian, &length); err != nil {
			return err
		}
		isp := make([]byte, length)
		if _, err := io.ReadFull(reader, isp); err != nil {
			return err
		}
		b.isp[code] = &ISP{
			name: string(isp),
		}
	}
	return nil
}

func (b *SinaIP) loadDistrict(data []byte) error {
	reader := bytes.NewReader(data)

	var count uint16
	if err := binary.Read(reader, binary.BigEndian, &count); err != nil {
		return err
	}
	b.district = make(map[uint16]*District, count)
	for i := uint16(0); i < count; i++ {
		var code uint16
		if err := binary.Read(reader, binary.BigEndian, &code); err != nil {
			return err
		}
		var length uint8
		if err := binary.Read(reader, binary.BigEndian, &length); err != nil {
			return err
		}
		district := make([]byte, length)
		if _, err := io.ReadFull(reader, district); err != nil {
			return err
		}
		b.district[code] = &District{
			name: string(district),
		}
	}
	return nil
}

func (b *SinaIP) loadIPSegment(data []byte) error {
	reader := bytes.NewReader(data)

	var count uint32
	if err := binary.Read(reader, binary.BigEndian, &count); err != nil {
		return err
	}
	b.segment = make([]*IPSegment, 0, count)
	for i := uint32(0); i < count; i++ {
		segment := new(IPSegment)
		if err := binary.Read(reader, binary.BigEndian, &segment.start); err != nil {
			return err
		}
		if err := binary.Read(reader, binary.BigEndian, &segment.end); err != nil {
			return err
		}
		var (
			countryCode  uint16
			provinceCode uint16
			cityCode     uint16
			ispCode      uint16
			districtCode uint16
			latitude     int32
			longitude    int32
		)
		if err := binary.Read(reader, binary.BigEndian, &countryCode); err != nil {
			return err
		} else {
			segment.country = b.country[countryCode]
		}

		if err := binary.Read(reader, binary.BigEndian, &provinceCode); err != nil {
			return err
		} else {
			segment.province = b.province[provinceCode]
		}

		if err := binary.Read(reader, binary.BigEndian, &cityCode); err != nil {
			return err
		} else {
			segment.city = b.city[cityCode]
		}

		if err := binary.Read(reader, binary.BigEndian, &ispCode); err != nil {
			return err
		} else {
			segment.isp = b.isp[ispCode]
		}

		if err := binary.Read(reader, binary.BigEndian, &districtCode); err != nil {
			return err
		} else {
			segment.district = b.district[districtCode]
		}

		if err := binary.Read(reader, binary.BigEndian, &latitude); err != nil {
			return err
		} else {
			segment.latitude = float64(latitude) / float64(10000)
		}
		if err := binary.Read(reader, binary.BigEndian, &longitude); err != nil {
			return err
		} else {
			segment.longitude = float64(longitude) / float64(10000)
		}
		b.segment = append(b.segment, segment)
	}
	return nil
}
