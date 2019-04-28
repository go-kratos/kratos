package region

import (
	"context"
	"strconv"
	"time"

	"go-common/app/interface/main/app-show/model"
	"go-common/app/interface/main/app-show/model/banner"
	resource "go-common/app/service/main/resource/model"
	"go-common/library/log"
)

var (
	_banners = map[int]map[int8]int{
		13: map[int8]int{
			model.PlatIPhone:   454,
			model.PlatIPad:     788,
			model.PlatAndroid:  617,
			model.PlatIPhoneI:  1022,
			model.PlatAndroidG: 1360,
			model.PlatAndroidI: 1791,
			model.PlatIPadI:    1192,
		},
		1: map[int8]int{
			model.PlatIPhone:   453,
			model.PlatIPad:     787,
			model.PlatAndroid:  616,
			model.PlatIPhoneI:  1017,
			model.PlatAndroidG: 1355,
			model.PlatAndroidI: 1785,
			model.PlatIPadI:    1187,
		},
		3: map[int8]int{
			model.PlatIPhone:   455,
			model.PlatIPad:     789,
			model.PlatAndroid:  618,
			model.PlatIPhoneI:  1028,
			model.PlatAndroidG: 1366,
			model.PlatAndroidI: 1798,
			model.PlatIPadI:    1198,
		},
		129: map[int8]int{
			model.PlatIPhone:   456,
			model.PlatIPad:     790,
			model.PlatAndroid:  619,
			model.PlatIPhoneI:  1033,
			model.PlatAndroidG: 1371,
			model.PlatAndroidI: 1804,
			model.PlatIPadI:    1203,
		},
		4: map[int8]int{
			model.PlatIPhone:   457,
			model.PlatIPad:     791,
			model.PlatAndroid:  620,
			model.PlatIPhoneI:  1038,
			model.PlatAndroidG: 1376,
			model.PlatAndroidI: 1810,
			model.PlatIPadI:    1208,
		},
		36: map[int8]int{
			model.PlatIPhone:   458,
			model.PlatIPad:     792,
			model.PlatAndroid:  621,
			model.PlatIPhoneI:  1043,
			model.PlatAndroidG: 1381,
			model.PlatAndroidI: 1816,
			model.PlatIPadI:    1213,
		},
		160: map[int8]int{
			model.PlatIPhone:   459,
			model.PlatIPad:     793,
			model.PlatAndroid:  622,
			model.PlatIPhoneI:  1048,
			model.PlatAndroidG: 1386,
			model.PlatAndroidI: 1822,
			model.PlatIPadI:    1218,
		},
		119: map[int8]int{
			model.PlatIPhone:   460,
			model.PlatIPad:     794,
			model.PlatAndroid:  623,
			model.PlatIPhoneI:  1053,
			model.PlatAndroidG: 1391,
			model.PlatAndroidI: 1828,
			model.PlatIPadI:    1223,
		},
		155: map[int8]int{
			model.PlatIPhone:   462,
			model.PlatIPad:     795,
			model.PlatAndroid:  624,
			model.PlatIPhoneI:  1058,
			model.PlatAndroidG: 1396,
			model.PlatAndroidI: 1834,
			model.PlatIPadI:    1228,
		},
		5: map[int8]int{
			model.PlatIPhone:   463,
			model.PlatIPad:     796,
			model.PlatAndroid:  625,
			model.PlatIPhoneI:  1063,
			model.PlatAndroidG: 1401,
			model.PlatAndroidI: 1840,
			model.PlatIPadI:    1233,
		},
		23: map[int8]int{
			model.PlatIPhone:   464,
			model.PlatIPad:     797,
			model.PlatAndroid:  626,
			model.PlatIPhoneI:  1068,
			model.PlatAndroidG: 1406,
			model.PlatAndroidI: 1846,
			model.PlatIPadI:    1238,
		},
		11: map[int8]int{
			model.PlatIPhone:   465,
			model.PlatIPad:     798,
			model.PlatAndroid:  627,
			model.PlatIPhoneI:  1073,
			model.PlatAndroidG: 1411,
			model.PlatAndroidI: 1852,
			model.PlatIPadI:    1243,
		},
		655: map[int8]int{
			model.PlatIPhone:   466,
			model.PlatIPad:     799,
			model.PlatAndroid:  628,
			model.PlatIPhoneI:  1079,
			model.PlatAndroidG: 1417,
			model.PlatAndroidI: 1859,
			model.PlatIPadI:    1249,
		},
		165: map[int8]int{
			model.PlatIPhone:   1473,
			model.PlatIPad:     1485,
			model.PlatAndroid:  1479,
			model.PlatIPhoneI:  1491,
			model.PlatAndroidG: 1497,
			model.PlatAndroidI: 1873,
			model.PlatIPadI:    1503,
		},
		167: map[int8]int{
			model.PlatIPhone:  1934,
			model.PlatIPad:    1932,
			model.PlatAndroid: 1933,
		},
		181: map[int8]int{
			model.PlatIPhone:  2225,
			model.PlatIPad:    2239,
			model.PlatAndroid: 2232,
		},
		177: map[int8]int{
			model.PlatIPhone:  2275,
			model.PlatIPad:    2289,
			model.PlatAndroid: 2282,
		},
		188: map[int8]int{
			model.PlatIPhone:   2996,
			model.PlatIPad:     3008,
			model.PlatAndroid:  3002,
			model.PlatIPhoneI:  3014,
			model.PlatAndroidG: 3020,
			model.PlatAndroidI: 3032,
			model.PlatIPadI:    3026,
		},
	}
	_bannersPlat = map[int8]string{
		model.PlatIPhone:   "454,453,455,456,457,458,459,460,462,463,464,465,466,1473,1934,2225,2275",
		model.PlatIPad:     "788,787,789,790,791,792,793,794,795,796,797,798,799,1485,1932,2239,2289",
		model.PlatAndroid:  "617,616,618,619,620,621,622,623,624,625,626,627,628,1479,1933,2232,2282",
		model.PlatIPhoneI:  "1022,1017,1028,1033,1038,1043,1048,1053,1058,1063,1068,1073,1079,1491",
		model.PlatAndroidG: "1360,1355,1366,1371,1376,1381,1386,1391,1396,1401,1406,1411,1417,1497",
		model.PlatAndroidI: "1791,1785,1798,1804,1810,1816,1822,1828,1834,1840,1846,1852,1859,1873",
		model.PlatIPadI:    "1192,1187,1198,1203,1208,1213,1218,1223,1228,1233,1238,1243,1249,1503",
	}
	_bannersPGC = map[int8]map[int]int{
		model.PlatAndroid: map[int]int{
			13:  83,
			167: 85,
			177: 232,
			11:  220,
			23:  49,
		},
		model.PlatIPhone: map[int]int{
			13:  97,
			167: 98,
			177: 233,
			11:  221,
			23:  50,
		},
		model.PlatIPad: map[int]int{
			13:  332,
			167: 333,
			177: 334,
			11:  336,
			23:  335,
		},
	}
)

// getBanners get banners by plat, build channel, ip.
func (s *Service) getBanners(c context.Context, plat int8, build, rid int, mid int64, channel, ip, buvid, network, mobiApp, device, adExtra string) (res map[string][]*banner.Banner) {
	var (
		resID = _banners[rid][plat]
		bs    []*banner.Banner
	)
	res = map[string][]*banner.Banner{}
	if bs = s.bgmBanners(c, plat, rid); len(bs) == 0 {
		bs = s.resBanners(c, plat, build, mid, resID, channel, ip, buvid, network, mobiApp, device, adExtra)
	}
	if len(bs) > 0 {
		res["top"] = bs
	}
	return
}

// resBannersplat
func (s *Service) resBanners(c context.Context, plat int8, build int, mid int64, resID int, channel, ip, buvid, network, mobiApp, device, adExtra string) (res []*banner.Banner) {
	var (
		plm   = s.bannerCache[plat] // operater banner
		err   error
		resbs map[int][]*resource.Banner
		tmp   []*resource.Banner
	)
	resIDStr := strconv.Itoa(resID)
	if resbs, err = s.res.ResBanner(c, plat, build, mid, resIDStr, channel, ip, buvid, network, mobiApp, device, adExtra, true); err != nil || len(resbs) == 0 {
		log.Error("s.res.ResBanner is null or err(%v)", err)
		resbs = plm
	}
	tmp = resbs[resID]
	for _, rb := range tmp {
		b := &banner.Banner{}
		b.ResChangeBanner(rb)
		res = append(res, b)
	}
	return
}

// bgmBanners bangumi banner
func (s *Service) bgmBanners(c context.Context, plat int8, rid int) (bgmBanner []*banner.Banner) {
	var (
		bgmb  = s.bannerBmgCache[plat][rid]
		resID = _banners[rid][plat]
	)
	for i, bb := range bgmb {
		b := &banner.Banner{}
		b.BgmChangeBanner(bb)
		b.RequestId = strconv.FormatInt(time.Now().UnixNano()/1000000, 10)
		b.Index = i + 1
		b.ResourceID = resID
		bgmBanner = append(bgmBanner, b)
	}
	return
}
