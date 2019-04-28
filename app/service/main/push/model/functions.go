package model

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"math"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go-common/library/log"

	"github.com/dgryski/go-farm"
)

// SplitInts splts string to int-slice by ,
func SplitInts(s string) (res []int) {
	if s == "" {
		return
	}
	ints := strings.Split(s, ",")
	for _, v := range ints {
		i, _ := strconv.Atoi(v)
		res = append(res, i)
	}
	return
}

// JoinInts merges int slice to string.
func JoinInts(ints []int) string {
	if len(ints) == 0 {
		return ""
	}
	if len(ints) == 1 {
		return strconv.Itoa(ints[0])
	}
	buf := bytes.Buffer{}
	for _, v := range ints {
		buf.WriteString(strconv.Itoa(v))
		buf.WriteString(",")
	}
	if buf.Len() > 0 {
		buf.Truncate(buf.Len() - 1)
	}
	return buf.String()
}

// ExistsInt judge if item in the ints.
func ExistsInt(ints []int, item int) (exists bool) {
	for _, i := range ints {
		if i == item {
			return true
		}
	}
	return false
}

// HashToken gets token's hash value.
func HashToken(token string) int64 {
	return int64(farm.Hash64([]byte(token)) % math.MaxInt64)
}

// RealTime culculates real time by timezone.
func RealTime(reportZone int) time.Time {
	now := time.Now()
	_, offset := now.Zone()
	return now.Add(time.Duration(reportZone-offset/3600) * time.Hour)
}

// Scheme gets uri scheme.
func Scheme(typ int8, val string, platform, build int) (uri string) {
	switch typ {
	case LinkTypeBangumi: // 番剧
		if platform == PlatformAndroid {
			uri = SchemeBangumiSeasonAndroid + val
		} else {
			uri = SchemeBangumiSeasonIOS + val
		}
	case LinkTypeVideo: // 视频
		if platform == PlatformAndroid {
			uri = SchemeVideoAndroid + val
		} else {
			uri = SchemeVideoIOS + val
		}
	case LinkTypeLive:
		var (
			param string
			parts = strings.Split(val, ",") // 值可能为 1 或者 1,0
		)
		if len(parts) == 2 {
			param = "?broadcast_type=" + parts[1]
		}
		uri = SchemeLive + parts[0] + param
		if platform == PlatformAndroid && build < 5290000 {
			uri = SchemeLiveAndroid + parts[0]
		}
	case LinkTypeSplist: // 专题
		uri = SchemeSplist + val
	case LinkTypeAuthor: // 个人空间
		if platform == PlatformAndroid {
			uri = SchemeAuthorAndroid + val
		} else {
			uri = SchemeAuthorIOS + val
		}
	case LinkTypeSearch: // 搜索
		if platform == PlatformAndroid {
			uri = SchemeSearchAndroid + val
		} else {
			uri = SchemeSearchIOS + val
		}
	case LinkTypeBrowser: // H5
		if platform == PlatformAndroid {
			uri = SchemeBrowserAndroid + url.QueryEscape(val)
		} else {
			// 容错逻辑，标准写法是 SchemeBrowserIOS + val，且 val 需要业务方进行 urlencode
			// 但是老客户端有bug，客户端会强制encode，客户端从 5.28 开始修了这个bug
			// 版本覆盖完全后，可改成标准写法
			uri = val
		}
	case LinkTypeVipBuy:
		uri = SchemeVipBuy + val
	case LinkTypeCustom:
		uri = val
	default:
		uri = ""
	}
	return
}

// ParseBuild parses string to build struct.
func ParseBuild(s string) (builds map[int]*Build) {
	builds = make(map[int]*Build)
	if s == "" {
		return
	}
	temp := make(map[string]*Build)
	if err := json.Unmarshal([]byte(s), &temp); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", s, err)
		return
	}
	for plat, build := range temp {
		p, _ := strconv.Atoi(plat)
		builds[p] = build
	}
	return
}

// TempTaskID gen temporary task ID.
func TempTaskID() string {
	n := time.Now().UnixNano()
	m := md5.Sum([]byte(strconv.FormatInt(n, 10)))
	return TempTaskPrefix + fmt.Sprintf("%x", m)[:8] // 要把taskid当作jobkey参数，jobkey要求长度最多9位, 1位prefix+8位时间hash值前段
}

// JobName gen job name.
func JobName(timestamp int64, content, linkValue, group string) int64 {
	s := []byte(fmt.Sprintf("%d%s%s%s%s", timestamp, time.Now().Format("20060102"), content, linkValue, group))
	return int64(farm.Hash64(s) % math.MaxInt64)
}

// Hash gen hash value by solt.
func Hash(salt string) string {
	s := salt + strconv.FormatInt(time.Now().UnixNano(), 10)
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

// 免打扰时间默认值
const (
	_defaultSilentBeginHour   = 22
	_defaultSilentEndHour     = 8
	_defaultSilentBeginMinute = 0
	_defaultSilentEndMinute   = 0
)

// ParseSilentTime .
func ParseSilentTime(s string) (st BusinessSilentTime) {
	st = BusinessSilentTime{
		BeginHour:   _defaultSilentBeginHour,
		EndHour:     _defaultSilentEndHour,
		BeginMinute: _defaultSilentBeginMinute,
		EndMinute:   _defaultSilentEndMinute,
	}
	s = strings.Trim(s, " ")
	if s == "" {
		return
	}
	r := strings.Split(s, "-")
	if len(r) != 2 {
		return
	}
	begin := strings.Split(r[0], ":")
	if len(begin) == 2 {
		st.BeginHour, _ = strconv.Atoi(begin[0])
		st.BeginMinute, _ = strconv.Atoi(begin[1])
	}
	end := strings.Split(r[1], ":")
	if len(end) == 2 {
		st.EndHour, _ = strconv.Atoi(end[0])
		st.EndMinute, _ = strconv.Atoi(end[1])
	}
	return st
}

// IsAndroid .
func IsAndroid(platformID int) bool {
	m := map[int]bool{
		PlatformIPhone: true,
		PlatformIPad:   true,
	}
	return !m[platformID]
}

// ValidateBuild checks token&platform valid.
func ValidateBuild(platform, build int, builds map[int]*Build) bool {
	if len(builds) == 0 {
		return true
	}
	if IsAndroid(platform) {
		platform = PlatformAndroid
	}
	if builds[platform] == nil {
		return true
	}
	c := builds[platform].Condition
	b := builds[platform].Build
	switch c {
	case "gt":
		return build > b
	case "gte":
		return build >= b
	case "lt":
		return build < b
	case "lte":
		return build <= b
	case "eq":
		return build == b
	case "ne":
		return build != b
	}
	return false
}
