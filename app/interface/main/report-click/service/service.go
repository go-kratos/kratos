package service

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	hmdl "go-common/app/interface/main/history/model"
	hrpc "go-common/app/interface/main/history/rpc/client"
	"go-common/app/interface/main/report-click/conf"
	"go-common/app/interface/main/report-click/dao"
	"go-common/app/interface/main/report-click/service/crypto/aes"
	"go-common/app/interface/main/report-click/service/crypto/padding"
	accmdl "go-common/app/service/main/account/model"
	accrpc "go-common/app/service/main/account/rpc/client"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/stat/prom"
	"go-common/library/sync/pipeline/fanout"
)

// 0 1 2 3 4 5 6 7 8 9 : <-> d w o i k p s x m q l
// 48 49 50 51 52 53 54 55 56 57 58 <-> 100 119 111 105 107 112 115 120 109 113 108
var ecKeys = map[rune]rune{
	48:  119,
	49:  111,
	50:  105,
	51:  107,
	52:  112,
	53:  115,
	54:  120,
	55:  109,
	56:  113,
	57:  108,
	58:  100,
	119: 48,
	111: 49,
	105: 50,
	107: 51,
	112: 52,
	115: 53,
	120: 54,
	109: 55,
	113: 56,
	108: 57,
	100: 58,
}

// Service service struct info.
type Service struct {
	c        *conf.Config
	d        *dao.Dao
	accRPC   *accrpc.Service3
	hisRPC   *hrpc.Service
	cache    *fanout.Fanout
	promErr  *prom.Prom
	promInfo *prom.Prom
}

// New service.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:        c,
		d:        dao.New(c),
		accRPC:   accrpc.New3(c.AccRPC),
		hisRPC:   hrpc.New(c.HisRPC),
		cache:    fanout.New("cache", fanout.Worker(1), fanout.Buffer(1024)),
		promErr:  prom.BusinessErrCount,
		promInfo: prom.BusinessInfoCount,
	}
	return
}

// FlashSigned flash Signed.
func (s *Service) FlashSigned(params url.Values, secret string, now time.Time) (err error) {
	st := params.Get("stime")
	stime, err := strconv.ParseInt(st, 10, 64)
	if err != nil {
		err = ecode.ClickQueryFormatErr
		return
	}
	if now.Unix()-stime > 60 {
		err = ecode.ClickServerTimeout
		return
	}
	sign := params.Get("sign")
	params.Del("sign")
	mh := md5.Sum([]byte(strings.ToLower(params.Encode()) + secret))
	if hex.EncodeToString(mh[:]) != sign {
		err = ecode.ClickQuerySignErr
	}
	return
}

// Decrypt decrypt bytes by aes key and iv.
func (s *Service) Decrypt(src []byte, aesKey, aesIv string) (res []byte, err error) {
	res, err = aes.CBCDecrypt(src, []byte(aesKey), []byte(aesIv), padding.PKCS5)
	if err != nil {
		log.Error("aes.CBCDecrypt(%s, %s, %s) error(%v)", base64.StdEncoding.EncodeToString(src), s.c.Click.AesKey, s.c.Click.AesIv, err)
		err = ecode.ClickAesDecryptErr
	}
	return
}

// Verify verify bytes from post body.
func (s *Service) Verify(src []byte, aesSalt string, now time.Time) (p url.Values, err error) {
	p, err = url.ParseQuery(string(src))
	if err != nil {
		err = ecode.ClickQueryFormatErr
		return
	}
	// check server time
	st := p.Get("stime")
	stime, err := strconv.ParseInt(st, 10, 64)
	if err != nil {
		err = ecode.ClickQueryFormatErr
		return
	}
	if now.Unix()-stime > 60*3 {
		err = ecode.ClickServerTimeout
		return
	}
	// verify sign
	sign := p.Get("sign")
	sbs, err := hex.DecodeString(sign)
	if err != nil {
		log.Error("hex.DecodeString(%s) error(%v)", sign, err)
		err = ecode.ClickQuerySignErr
		return
	}
	p.Del("sign")
	// sha 256
	h := sha256.New()
	//	h.Write([]byte(strings.ToLower(p.Encode())))
	h.Write([]byte(p.Encode()))
	h.Write([]byte(aesSalt))
	bs := h.Sum(nil)
	// bytes queal
	if !bytes.Equal(sbs, bs) {
		log.Error("hmac.Equal(%s, %x) params(%s) not equal", sign, bs, p.Encode())
		err = ecode.ClickHmacSignErr
	}
	return
}

// Play send play count to kafka.
func (s *Service) Play(c context.Context, plat, aid, cid, part, mid, level, ftime, stime, did, ip, agent, buvid, cookieSid, refer, typeID, subType, sid, epid, playMode, platform, device, mobiAapp, autoPlay, session string) {
	if aid == "" || aid == "0" {
		return
	}
	m, errP := strconv.ParseInt(mid, 10, 64)
	if errP != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", mid, errP)
		mid = "0"
	}
	if m != 0 {
		arg := &accmdl.ArgMid{Mid: m}
		res, err := s.accRPC.Card3(c, arg)
		if err != nil {
			log.Error("s.accRPC.UserInfo() error(%v)", err)
			return
		}
		if res.Silence == 1 {
			log.Warn("user mid(%d) spacesta(%d) too lower", m, res.Silence)
			return
		}
		level = fmt.Sprintf("%d", res.Level)
	}
	s.cache.Do(c, func(ctx context.Context) {
		s.d.Play(ctx, plat, aid, cid, part, mid, level, ftime, stime, did, ip, agent, buvid, cookieSid, refer, typeID, subType, sid, epid, playMode, platform, device, mobiAapp, autoPlay, session)
	})
}

// GenDid gen did.
func (s *Service) GenDid(ip string, now time.Time) string {
	var src string
	ft := now.Unix() - int64(now.Second())
	uip, ok := parseIP(ip)
	if uip == nil {
		return ""
	}
	if ok {
		src = fmt.Sprintf("%d:%d", netAtoN(uip), ft)
		return myEncryptDecrypt(src)
	}
	fs := encode(uint64(ft))
	ipRes := ipv6AtoN(uip)
	if len(ipRes) > 25 { // total 32, 25 for ip, 1 for :, 6 for ftime
		ipRes = ipRes[:25]
	}
	return fmt.Sprintf("%s:%s", ipRes, fs)
}

// CheckDid check did.
func (s *Service) CheckDid(did string) (ip, ft string) {
	params := strings.Split(did, ":")
	if len(params) == 4 {
		log.Warn("report click did:%s", did)
		return ntoIPv6(params[:3]), fmt.Sprintf("%d", decode([]byte(params[3])))
	}
	dst := myEncryptDecrypt(did)
	params = strings.Split(dst, ":")
	if len(params) != 2 {
		return
	}
	ipInt, _ := strconv.ParseInt(params[0], 10, 64)
	ip = netNtoA(uint32(ipInt))
	ft = params[1]
	return
}

func myEncryptDecrypt(src string) (dst string) {
	var tmp []rune
	for _, k := range src {
		if _, ok := ecKeys[k]; !ok {
			return ""
		}
		tmp = append(tmp, ecKeys[k])
	}
	dst = string(tmp)
	return
}

// Report report to history.
func (s *Service) Report(c context.Context, proStr, cidStr, tpStr, subType, realtimeStr, aidStr, midstr, sidStr, epidStr, dtStr, tsStr string) (err error) {
	var (
		tp, stp, dt int
	)
	if tp, err = strconv.Atoi(tpStr); err != nil {
		log.Warn("Report type:%s", tpStr)
	}
	stp, _ = strconv.Atoi(subType)
	mid, _ := strconv.ParseInt(midstr, 10, 64)
	aid, _ := strconv.ParseInt(aidStr, 10, 64)
	sid, _ := strconv.ParseInt(sidStr, 10, 64)
	epid, _ := strconv.ParseInt(epidStr, 10, 64)
	cid, _ := strconv.ParseInt(cidStr, 10, 64)
	if aid == 0 && cid == 0 {
		return ecode.RequestErr
	}
	pro, _ := strconv.ParseInt(proStr, 10, 64)
	realtime, _ := strconv.ParseInt(realtimeStr, 10, 64)
	if dt, err = strconv.Atoi(dtStr); err != nil {
		dt = 2
	}
	ip := metadata.String(c, metadata.RemoteIP)
	ts, _ := strconv.ParseInt(tsStr, 10, 64)
	history := &hmdl.History{Aid: aid, Sid: sid, Epid: epid, TP: int8(tp), STP: int8(stp), Cid: cid, DT: int8(dt), Pro: pro, Unix: ts}
	arg := &hmdl.ArgHistory{Mid: mid, Realtime: realtime, RealIP: ip, History: history}
	return s.hisRPC.Add(c, arg)
}
