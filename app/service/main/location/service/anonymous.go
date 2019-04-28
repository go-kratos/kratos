package service

import (
	"net"

	"go-common/app/service/main/location/model"
	"go-common/library/log"

	maxminddb "github.com/oschwald/maxminddb-golang"
)

// NewAnonym load bin file.
func (s *Service) NewAnonym(path string) (reader *maxminddb.Reader, err error) {
	if reader, err = maxminddb.Open(path); err != nil {
		log.Error("NewAnonym bin file error(%v)", err)
	}
	return
}

// Anonym .
func (s *Service) Anonym(ipStr string) (res *model.AnonymousIP, err error) {
	nip := net.ParseIP(ipStr)
	res = &model.AnonymousIP{}
	if nip == nil || nip.To4() == nil {
		log.Error("invalid ip(%s) format", ipStr)
		return
	}
	if err = s.anonym.Lookup(nip, res); err != nil {
		log.Error("Anonym bin file ip(%v) error(%v)", ipStr, err)
	}
	return
}
