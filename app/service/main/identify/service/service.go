package service

import (
	"context"
	"net"
	"strings"

	"go-common/app/service/main/identify/conf"
	"go-common/app/service/main/identify/dao"
	"go-common/library/cache"
	"go-common/library/log"
	"go-common/library/queue/databus"

	"github.com/pkg/errors"
)

var (
	emptyStringArray = []string{""}
)

// Service is a identify service.
type Service struct {
	c *conf.Config
	d *dao.Dao
	// cache proc
	cache *cache.Cache
	// login log
	loginlogchan chan func()
	// loginLogDataBus
	loginLogDataBus *databus.Databus
	// loginCacheExpires
	loginCacheExpires int32
	// intranetCIDR
	intranetCIDR []*net.IPNet
}

// New new a identify service.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:     c,
		d:     dao.New(c),
		cache: cache.New(1, 1024),
		// login log chan
		loginlogchan:    make(chan func(), 10240),
		loginLogDataBus: databus.New(c.DataBus.UserLog),
		// loginCacheExpires
		loginCacheExpires: 3600,
	}
	cidrs := make([]*net.IPNet, 0, len(c.Identify.IntranetCIDR))
	for _, raw := range c.Identify.IntranetCIDR {
		_, inet, err := net.ParseCIDR(raw)
		if err != nil {
			panic(errors.Wrapf(err, "Invalid CIDR: %s", raw))
		}
		cidrs = append(cidrs, inet)
	}
	s.intranetCIDR = cidrs
	if c.Identify.LoginCacheExpires > 0 {
		s.loginCacheExpires = c.Identify.LoginCacheExpires
	}
	for i := 0; i < s.c.Identify.LoginLogConsumerSize; i++ {
		go func() {
			for f := range s.loginlogchan {
				f()
			}
		}()
	}
	return
}

// Close dao.
func (s *Service) Close() {
	s.d.Close()
}

// Ping check server ok.
func (s *Service) Ping(c context.Context) (err error) {
	return s.d.Ping(c)
}

func (s *Service) addLoginLog(f func()) {
	select {
	case s.loginlogchan <- f:
	default:
		log.Warn("loginlogchan chan full")
	}
}

// readCookiesVal parses all "Cookie" values from the input line and
// returns the successfully parsed values.
//
// if filter isn't empty, only cookies of that name are returned
func readCookiesVal(line, filter string) []string {
	if line == "" {
		return emptyStringArray
	}

	var cookies []string
	parts := strings.Split(strings.TrimSpace(line), ";")
	if len(parts) == 1 && parts[0] == "" {
		return emptyStringArray
	}
	// Per-line attributes
	for i := 0; i < len(parts); i++ {
		parts[i] = strings.TrimSpace(parts[i])
		if len(parts[i]) == 0 {
			continue
		}
		name, val := parts[i], ""
		if j := strings.Index(name, "="); j >= 0 {
			name, val = name[:j], name[j+1:]
		}
		// if !isCookieNameValid(name) {
		// 	continue
		// }
		if filter != "" && filter != name {
			continue
		}
		val, ok := parseCookieValue(val, true)
		if !ok {
			continue
		}
		cookies = append(cookies, val)
	}
	if len(cookies) == 0 {
		return emptyStringArray
	}
	return cookies
}

func parseCookieValue(raw string, allowDoubleQuote bool) (string, bool) {
	// Strip the quotes, if present.
	if allowDoubleQuote && len(raw) > 1 && raw[0] == '"' && raw[len(raw)-1] == '"' {
		raw = raw[1 : len(raw)-1]
	}
	// for i := 0; i < len(raw); i++ {
	// 	if !validCookieValueByte(raw[i]) {
	// 		return "", false
	// 	}
	// }
	return raw, true
}

func (s *Service) isIntranetIP(ip string) (exists bool) {
	for _, c := range s.intranetCIDR {
		if c.Contains(net.ParseIP(ip)) {
			exists = true
			return
		}
	}
	return
}
