package service

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/satori/go.uuid"
)

func generateToken(mid int64, ct time.Time, dc int) string {
	return generateAK(mid, int(ct.Month()), dc)
}

func generateRefresh(mid int64, ct time.Time, dc int) string {
	return generateRK(mid, int(ct.Month()), dc)
}

// oldSession get old session
func (s *Service) oldSession(mid, expires int64, month int) string {
	session := generateSD(mid, month, s.c.DC.Num)
	return fmt.Sprintf("%s,%d,%s", session[:8], expires, session[24:])
}

func decodeSession(sd string) (res []byte, err error) {
	// if is new sd
	if len(sd) == _newSessionHexLen {
		return hex.DecodeString(sd)
	}
	// else if is old sd
	return []byte(sd), nil
}

func generateCSRF(mid int64) (res string) {
	return md5Hex(fmt.Sprintf("%d%d%d", rand.Int63n(100000000), time.Now().Nanosecond(), mid))
}

func generateRK(mid int64, month, dc int) string {
	return generateByAdditional(mid, month, dc, "refresh")
}

func generateAK(mid int64, month, dc int) string {
	return generateByAdditional(mid, month, dc, "token")
}

func generateSD(mid int64, month, dc int) string {
	return generateByAdditional(mid, month, dc, "session")
}

func generateByAdditional(mid int64, month, dc int, additional string) string {
	t := md5Hex(fmt.Sprintf("%s,%d,%s", uuid.NewV4().String(), mid, additional))
	// [0, 29] + 1 + 1
	return t[:30] + formatHex(month) + formatHex(dc)
}

func formatHex(n int) string {
	return fmt.Sprintf("%x", n)
}

func calcMonDelta(ak string, now time.Time) (res int, err error) {
	var parsedMon int
	if parsedMon, err = parseMonth(ak); err != nil {
		return
	}

	curMon := int(now.Month())
	// check if cur mon
	if curMon == parsedMon {
		return 0, nil
	}

	delta := curMon - parsedMon
	if delta < 0 {
		delta += 12
	}
	return delta, nil
}

func monDiff(t time.Time, delta int) time.Time {
	if delta == 0 {
		return t
	}
	year, month, _ := t.Date()
	thisMonthFirstDay := time.Date(year, month, 1, 1, 1, 1, 1, t.Location())
	return thisMonthFirstDay.AddDate(0, delta, 0)
}

func parseMonth(ak string) (int, error) {
	n, err := strconv.ParseInt(ak[30:31], 16, 64)
	if err != nil {
		return 0, err
	}
	return int(n), nil
}
