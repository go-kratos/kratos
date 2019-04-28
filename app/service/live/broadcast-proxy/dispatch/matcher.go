package dispatch

import (
	"encoding/json"
	"errors"
	"github.com/ipipdotnet/ipdb-go"
	"go-common/app/service/live/broadcast-proxy/conf"
	"go-common/app/service/live/broadcast-proxy/expr"
	"go-common/library/log"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"sync"
)

type Matcher struct {
	ipDataV4             *ipdb.City
	ipDataV6             *ipdb.City
	heapPool             sync.Pool
	Config               string
	MaxLimit             int    `json:"ip_max_limit"`
	DefaultDomain        string `json:"default_domain"`
	WildcardDomainSuffix string `json:"wildcard_domain_suffix"`
	CommonDispatch       struct {
		ChinaDispatch struct {
			ChinaTelecom *CommonBucket `json:"china_telecom"`
			ChinaUnicom  *CommonBucket `json:"china_unicom"`
			CMCC         *CommonBucket `json:"cmcc"`
			ChinaOther   *CommonBucket `json:"other"`
		} `json:"china"`
		OverseaDispatch     []*CommonRuleBucket `json:"oversea"`
		UnknownAreaDispatch *CommonBucket       `json:"unknown"`
	} `json:"danmaku_common_dispatch"`
	VIPDispatch []*VIPRuleBucket    `json:"danmaku_vip_dispatch"`
	ServerGroup map[string][]string `json:"danmaku_comet_group"`
	ServerHost  map[string]string   `json:"danmaku_comet_host"`
	IPBlack     []string            `json:"ip_black"`
	TempV6      []string            `json:"temp_v6"`
	forbiddenIP map[string]struct{}
}

type CommonBucket struct {
	Master map[string]int `json:"master"`
	Slave  map[string]int `json:"slave"`
}

type CommonRuleBucket struct {
	CommonBucket
	Rule     string `json:"rule"`
	RuleExpr expr.Expr
}

type VIPRuleBucket struct {
	Rule     string `json:"rule"`
	RuleExpr expr.Expr
	IP       []string `json:"ip"`
	Group    []string `json:"group"`
}

func NewMatcher(matcherConfig []byte, ipDataV4 *ipdb.City, ipDataV6 *ipdb.City, dispatchConfig *conf.DispatchConfig) (*Matcher, error) {
	matcher := new(Matcher)
	matcher.heapPool = sync.Pool{
		New: func() interface{} {
			return NewMinHeap()
		},
	}
	matcher.forbiddenIP = make(map[string]struct{})
	if ipDataV4 == nil || ipDataV6 == nil {
		return nil, errors.New("invalid IP database")
	}
	matcher.ipDataV4 = ipDataV4
	matcher.ipDataV6 = ipDataV6
	matcher.Config = string(matcherConfig)
	if err := json.Unmarshal(matcherConfig, matcher); err != nil {
		return nil, err
	}
	for _, ip := range matcher.IPBlack {
		matcher.forbiddenIP[ip] = struct{}{}
	}
	parser := expr.NewExpressionParser()
	for _, oversea := range matcher.CommonDispatch.OverseaDispatch {
		if oversea.Rule == "" {
			oversea.Rule = "true"
		}
		if err := parser.Parse(oversea.Rule); err != nil {
			log.Error("[Matcher] Parse rule expr:%s, error:%+v", oversea.Rule, err)
			return nil, err
		}
		for _, variable := range parser.GetVariable() {
			if variable != "$lng" && variable != "$lat" {
				return nil, errors.New("oversea dispatch only supports variable $lng and $lat")
			}
		}
		oversea.RuleExpr = parser.GetExpr()
	}
	for _, vip := range matcher.VIPDispatch {
		if err := parser.Parse(vip.Rule); err != nil {
			log.Error("[Matcher] Parse rule expr:%s, error:%+v", vip.Rule, err)
			return nil, err
		}
		for _, variable := range parser.GetVariable() {
			if variable != "$uid" {
				return nil, errors.New("vip dispatch only supports variable $uid")
			}
		}
		if len(parser.GetVariable()) == 0 {
			return nil, errors.New("vip dispatch must contains variable $uid")
		}
		vip.RuleExpr = parser.GetExpr()
	}

	if matcher.MaxLimit == 0 && dispatchConfig != nil {
		matcher.MaxLimit = dispatchConfig.MaxLimit
	}
	if matcher.DefaultDomain == "" && dispatchConfig != nil {
		matcher.DefaultDomain = dispatchConfig.DefaultDomain
	}
	if matcher.WildcardDomainSuffix == "" && dispatchConfig != nil {
		matcher.WildcardDomainSuffix = dispatchConfig.WildcardDomainSuffix
	}
	return matcher, nil
}

func (matcher *Matcher) GetConfig() string {
	return matcher.Config
}

func (matcher *Matcher) Dispatch(ip string, uid int64) ([]string, []string) {
	danmakuIP := matcher.dispatchInternal(ip, uid)
	danmakuHost := make([]string, 0, len(danmakuIP))
	for _, singleDanmakuIP := range danmakuIP {
		if host, ok := matcher.ServerHost[singleDanmakuIP]; ok {
			danmakuHost = append(danmakuHost, host+matcher.WildcardDomainSuffix)
		}
	}
	danmakuIP = append(danmakuIP, matcher.DefaultDomain)
	danmakuHost = append(danmakuHost, matcher.DefaultDomain)
	return danmakuIP, danmakuHost
}

func (matcher *Matcher) dispatchInternal(ip string, uid int64) []string {
	if _, ok := matcher.forbiddenIP[ip]; ok {
		return []string{}
	}
	// VIP Dispatch
	vipDispatchEnv := make(map[expr.Var]interface{})
	vipDispatchEnv[expr.Var("$uid")] = uid
	for _, vip := range matcher.VIPDispatch {
		if v, err := expr.SafetyEvalBool(vip.RuleExpr, vipDispatchEnv); v && err == nil {
			return matcher.pickFromVIPRuleBucket(vip)
		} else {
			if err != nil {
				log.Error("[Matcher] VIP dispatch, uid:%d, eval rule expr:%s error:%+v", uid, vip.Rule, err)
			}
		}
	}
	// Common Dispatch
	var ipDatabase *ipdb.City
	for i := 0; i < len(ip); i++ {
		if ip[i] == '.' {
			ipDatabase = matcher.ipDataV4
			break
		} else if ip[i] == ':' {
			ipDatabase = matcher.ipDataV6
			//break
			//TODO: this is temp solution, replace this block with "break" here when all server supports IPv6
			return matcher.randomPickN(matcher.TempV6, matcher.MaxLimit)
		}
	}
	if ipDatabase == nil {
		return matcher.pickFromCommonBucket(matcher.CommonDispatch.UnknownAreaDispatch)
	}

	detail, err := ipDatabase.FindMap(ip, "EN")
	if err != nil {
		return matcher.pickFromCommonBucket(matcher.CommonDispatch.UnknownAreaDispatch)
	}
	country := strings.TrimSpace(detail["country_name"])
	province := strings.TrimSpace(detail["region_name"])
	isp := strings.TrimSpace(detail["isp_domain"])
	latitude, _ := strconv.ParseFloat(detail["latitude"], 64)
	longitude, _ := strconv.ParseFloat(detail["longitude"], 64)

	if country != "China" && country != "Reserved" && country != "LAN Address" && country != "Loopback" {
		return matcher.pickFromCommonRuleBucket(matcher.CommonDispatch.OverseaDispatch, latitude, longitude)
	} else if country == "China" {
		if province == "Hong Kong" || province == "Macau" || province == "Taiwan" {
			return matcher.pickFromCommonRuleBucket(matcher.CommonDispatch.OverseaDispatch, latitude, longitude)
		} else {
			switch isp {
			case "ChinaTelecom":
				return matcher.pickFromCommonBucket(matcher.CommonDispatch.ChinaDispatch.ChinaTelecom)
			case "ChinaMobile":
				return matcher.pickFromCommonBucket(matcher.CommonDispatch.ChinaDispatch.CMCC)
			case "ChinaUnicom":
				return matcher.pickFromCommonBucket(matcher.CommonDispatch.ChinaDispatch.ChinaUnicom)
			default:
				return matcher.pickFromCommonBucket(matcher.CommonDispatch.ChinaDispatch.ChinaOther)
			}
		}
	} else {
		return matcher.pickFromCommonBucket(matcher.CommonDispatch.UnknownAreaDispatch)
	}
}

func (matcher *Matcher) pickFromCommonRuleBucket(overseaBucket []*CommonRuleBucket, latitude float64, longitude float64) []string {
	overseaDispatchEnv := make(map[expr.Var]interface{})
	overseaDispatchEnv[expr.Var("$lat")] = latitude
	overseaDispatchEnv[expr.Var("$lng")] = longitude
	for _, bucket := range overseaBucket {
		if v, err := expr.SafetyEvalBool(bucket.RuleExpr, overseaDispatchEnv); v && err == nil {
			return matcher.pickFromCommonBucket(&bucket.CommonBucket)
		}
	}
	return []string{}
}

func (matcher *Matcher) pickOneFromWeightedGroup(groupWeightDict map[string]int) (string, string) {
	var luckyKey float64
	var luckyGroup string
	for group, weight := range groupWeightDict {
		if weight > 0 {
			key := math.Pow(rand.Float64(), 1.0/float64(weight))
			if key >= luckyKey {
				luckyKey = key
				luckyGroup = group
			}
		}
	}
	luckyIP := matcher.ServerGroup[luckyGroup]
	if len(luckyIP) == 0 {
		return "", ""
	}
	return matcher.randomPickOne(luckyIP), luckyGroup
}

func (matcher *Matcher) pickNFromWeightedGroup(groupWeightDict map[string]int, n int, groupIgnore string) []string {
	h := matcher.heapPool.Get().(*MinHeap)
	for group, weight := range groupWeightDict {
		if group != groupIgnore && weight > 0 {
			key := math.Pow(rand.Float64(), 1.0/float64(weight))
			if h.HeapLength() < n {
				h.HeapPush(group, key)
			} else {
				_, top, _ := h.HeapTop()
				if key > top {
					h.HeapPush(group, key)
					h.HeapPop()
				}
			}
		}
	}
	r := make([]string, 0, n)
	for h.HeapLength() > 0 {
		v, _, _ := h.HeapPop()
		member := matcher.ServerGroup[v.(string)]
		if len(member) > 0 {
			r = append(r, matcher.randomPickOne(member))
		}
	}
	matcher.heapPool.Put(h)
	return r
}

func (matcher *Matcher) pickFromCommonBucket(b *CommonBucket) []string {
	r := make([]string, 0, matcher.MaxLimit)
	masterIP, masterGroup := matcher.pickOneFromWeightedGroup(b.Master)
	if masterIP != "" {
		r = append(r, masterIP)
	}
	for _, slaveIP := range matcher.pickNFromWeightedGroup(b.Slave, matcher.MaxLimit-len(r), masterGroup) {
		r = append(r, slaveIP)
	}
	return r
}

func (matcher *Matcher) pickFromVIPRuleBucket(b *VIPRuleBucket) []string {
	var length int
	for _, group := range b.Group {
		length += len(matcher.ServerGroup[group])
	}
	length += len(b.IP)
	candidate := make([]string, length)
	i := 0
	for _, group := range b.Group {
		i += copy(candidate[i:], matcher.ServerGroup[group])
	}
	i += copy(candidate[i:], b.IP)
	return matcher.randomPickN(candidate, matcher.MaxLimit)
}

func (matcher *Matcher) randomPickOne(s []string) string {
	return s[rand.Intn(len(s))]
}

func (matcher *Matcher) randomPickN(s []string, n int) []string {
	var r []string
	if n > len(s) {
		n = len(s)
	}
	for _, v := range rand.Perm(len(s))[0:n] {
		r = append(r, s[v])
	}
	return r
}
