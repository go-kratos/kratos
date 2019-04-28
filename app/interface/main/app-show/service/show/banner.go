package show

import (
	"context"

	"go-common/app/interface/main/app-show/model"
	"go-common/app/interface/main/app-show/model/banner"
	resource "go-common/app/service/main/resource/model"
	"go-common/library/log"
)

var (
	_banners = map[string]map[int8]map[string]int{
		"-1": map[int8]map[string]int{
			model.PlatIPhone: map[string]int{
				"bottom": 467,
			},
			model.PlatAndroid: map[string]int{
				"bottom": 631,
			},
			model.PlatIPad: map[string]int{
				"bottom": 771,
			},
			model.PlatIPhoneI: map[string]int{
				"bottom": 947,
			},
			model.PlatAndroidG: map[string]int{
				"bottom": 1285,
			},
			model.PlatAndroidI: map[string]int{
				"bottom": 1707,
			},
			model.PlatIPadI: map[string]int{
				"bottom": 1117,
			},
		},
		"0": map[int8]map[string]int{
			model.PlatIPhone: map[string]int{
				"top":    467,
				"center": 468,
				"bottom": 469,
			},
			model.PlatAndroid: map[string]int{
				"top":    631,
				"center": 632,
				"bottom": 633,
			},
			model.PlatIPad: map[string]int{
				"top":    771,
				"center": 772,
				"bottom": 773,
			},
			model.PlatIPhoneI: map[string]int{
				"top":    947,
				"center": 952,
				"bottom": 957,
			},
			model.PlatAndroidG: map[string]int{
				"top":    1285,
				"center": 1290,
				"bottom": 1295,
			},
			model.PlatAndroidI: map[string]int{
				"top":    1707,
				"center": 1712,
				"bottom": 1717,
			},
			model.PlatIPadI: map[string]int{
				"top":    1117,
				"center": 1122,
				"bottom": 1127,
			},
		},
		"65537": map[int8]map[string]int{
			model.PlatIPhone: map[string]int{
				"bottom": 482,
			},
			model.PlatAndroid: map[string]int{
				"bottom": 646,
			},
			model.PlatIPad: map[string]int{
				"bottom": 786,
			},
			model.PlatIPhoneI: map[string]int{
				"bottom": 1013,
			},
			model.PlatAndroidG: map[string]int{
				"bottom": 1351,
			},
			model.PlatAndroidI: map[string]int{
				"bottom": 1773,
			},
			model.PlatIPadI: map[string]int{
				"bottom": 1183,
			},
		},
		"13": map[int8]map[string]int{
			model.PlatIPhone: map[string]int{
				"bottom": 471,
			},
			model.PlatAndroid: map[string]int{
				"bottom": 635,
			},
			model.PlatIPad: map[string]int{
				"bottom": 775,
			},
			model.PlatIPhoneI: map[string]int{
				"bottom": 967,
			},
			model.PlatAndroidG: map[string]int{
				"bottom": 1305,
			},
			model.PlatAndroidI: map[string]int{
				"bottom": 1727,
			},
			model.PlatIPadI: map[string]int{
				"bottom": 1137,
			},
		},
		"1": map[int8]map[string]int{
			model.PlatIPhone: map[string]int{
				"bottom": 470,
			},
			model.PlatAndroid: map[string]int{
				"bottom": 634,
			},
			model.PlatIPad: map[string]int{
				"bottom": 774,
			},
			model.PlatIPhoneI: map[string]int{
				"bottom": 962,
			},
			model.PlatAndroidG: map[string]int{
				"bottom": 1300,
			},
			model.PlatAndroidI: map[string]int{
				"bottom": 1722,
			},
			model.PlatIPadI: map[string]int{
				"bottom": 1132,
			},
		},
		"3": map[int8]map[string]int{
			model.PlatIPhone: map[string]int{
				"bottom": 472,
			},
			model.PlatAndroid: map[string]int{
				"bottom": 636,
			},
			model.PlatIPad: map[string]int{
				"bottom": 776,
			},
			model.PlatIPhoneI: map[string]int{
				"bottom": 971,
			},
			model.PlatAndroidG: map[string]int{
				"bottom": 1309,
			},
			model.PlatAndroidI: map[string]int{
				"bottom": 1731,
			},
			model.PlatIPadI: map[string]int{
				"bottom": 1141,
			},
		},
		"129": map[int8]map[string]int{
			model.PlatIPhone: map[string]int{
				"bottom": 473,
			},
			model.PlatAndroid: map[string]int{
				"bottom": 637,
			},
			model.PlatIPad: map[string]int{
				"bottom": 777,
			},
			model.PlatIPhoneI: map[string]int{
				"bottom": 975,
			},
			model.PlatAndroidG: map[string]int{
				"bottom": 1313,
			},
			model.PlatAndroidI: map[string]int{
				"bottom": 1735,
			},
			model.PlatIPadI: map[string]int{
				"bottom": 1145,
			},
		},
		"4": map[int8]map[string]int{
			model.PlatIPhone: map[string]int{
				"bottom": 474,
			},
			model.PlatAndroid: map[string]int{
				"bottom": 638,
			},
			model.PlatIPad: map[string]int{
				"bottom": 778,
			},
			model.PlatIPhoneI: map[string]int{
				"bottom": 979,
			},
			model.PlatAndroidG: map[string]int{
				"bottom": 1317,
			},
			model.PlatAndroidI: map[string]int{
				"bottom": 1739,
			},
			model.PlatIPadI: map[string]int{
				"bottom": 1149,
			},
		},
		"36": map[int8]map[string]int{
			model.PlatIPhone: map[string]int{
				"bottom": 475,
			},
			model.PlatAndroid: map[string]int{
				"bottom": 639,
			},
			model.PlatIPad: map[string]int{
				"bottom": 779,
			},
			model.PlatIPhoneI: map[string]int{
				"bottom": 983,
			},
			model.PlatAndroidG: map[string]int{
				"bottom": 1321,
			},
			model.PlatAndroidI: map[string]int{
				"bottom": 1706,
			},
			model.PlatIPadI: map[string]int{
				"bottom": 1153,
			},
		},
		"160": map[int8]map[string]int{
			model.PlatIPhone: map[string]int{
				"bottom": 476,
			},
			model.PlatAndroid: map[string]int{
				"bottom": 640,
			},
			model.PlatIPad: map[string]int{
				"bottom": 780,
			},
			model.PlatIPhoneI: map[string]int{
				"bottom": 987,
			},
			model.PlatAndroidG: map[string]int{
				"bottom": 1325,
			},
			model.PlatAndroidI: map[string]int{
				"bottom": 1747,
			},
			model.PlatIPadI: map[string]int{
				"bottom": 1157,
			},
		},
		"119": map[int8]map[string]int{
			model.PlatIPhone: map[string]int{
				"bottom": 477,
			},
			model.PlatAndroid: map[string]int{
				"bottom": 641,
			},
			model.PlatIPad: map[string]int{
				"bottom": 781,
			},
			model.PlatIPhoneI: map[string]int{
				"bottom": 992,
			},
			model.PlatAndroidG: map[string]int{
				"bottom": 1330,
			},
			model.PlatAndroidI: map[string]int{
				"bottom": 1752,
			},
			model.PlatIPadI: map[string]int{
				"bottom": 1162,
			},
		},
		"155": map[int8]map[string]int{
			model.PlatIPhone: map[string]int{
				"bottom": 478,
			},
			model.PlatAndroid: map[string]int{
				"bottom": 642,
			},
			model.PlatIPad: map[string]int{
				"bottom": 782,
			},
			model.PlatIPhoneI: map[string]int{
				"bottom": 997,
			},
			model.PlatAndroidG: map[string]int{
				"bottom": 1335,
			},
			model.PlatAndroidI: map[string]int{
				"bottom": 1757,
			},
			model.PlatIPadI: map[string]int{
				"bottom": 1167,
			},
		},
		"5": map[int8]map[string]int{
			model.PlatIPhone: map[string]int{
				"bottom": 479,
			},
			model.PlatAndroid: map[string]int{
				"bottom": 643,
			},
			model.PlatIPad: map[string]int{
				"bottom": 783,
			},
			model.PlatIPhoneI: map[string]int{
				"bottom": 1001,
			},
			model.PlatAndroidG: map[string]int{
				"bottom": 1339,
			},
			model.PlatAndroidI: map[string]int{
				"bottom": 1761,
			},
			model.PlatIPadI: map[string]int{
				"bottom": 1171,
			},
		},
		"23": map[int8]map[string]int{
			model.PlatIPhone: map[string]int{
				"bottom": 480,
			},
			model.PlatAndroid: map[string]int{
				"bottom": 644,
			},
			model.PlatIPad: map[string]int{
				"bottom": 784,
			},
			model.PlatIPhoneI: map[string]int{
				"bottom": 1005,
			},
			model.PlatAndroidG: map[string]int{
				"bottom": 1343,
			},
			model.PlatAndroidI: map[string]int{
				"bottom": 1765,
			},
			model.PlatIPadI: map[string]int{
				"bottom": 1175,
			},
		},
		"11": map[int8]map[string]int{
			model.PlatIPhone: map[string]int{
				"bottom": 481,
			},
			model.PlatAndroid: map[string]int{
				"bottom": 645,
			},
			model.PlatIPad: map[string]int{
				"bottom": 785,
			},
			model.PlatIPhoneI: map[string]int{
				"bottom": 1009,
			},
			model.PlatAndroidG: map[string]int{
				"bottom": 1347,
			},
			model.PlatAndroidI: map[string]int{
				"bottom": 1769,
			},
			model.PlatIPadI: map[string]int{
				"bottom": 1179,
			},
		},
		"165": map[int8]map[string]int{
			model.PlatIPhone: map[string]int{
				"bottom": 1643,
			},
			model.PlatAndroid: map[string]int{
				"bottom": 1639,
			},
			model.PlatIPad: map[string]int{
				"bottom": 1647,
			},
			model.PlatIPhoneI: map[string]int{
				"bottom": 1643,
			},
			model.PlatAndroidG: map[string]int{
				"bottom": 1639,
			},
			model.PlatIPadI: map[string]int{
				"bottom": 1647,
			},
		},
		"167": map[int8]map[string]int{
			model.PlatIPhone: map[string]int{
				"bottom": 1950,
			},
			model.PlatAndroid: map[string]int{
				"bottom": 1952,
			},
			model.PlatIPad: map[string]int{
				"bottom": 1951,
			},
		},
		"181": map[int8]map[string]int{
			model.PlatIPhone: map[string]int{
				"bottom": 2245,
			},
			model.PlatAndroid: map[string]int{
				"bottom": 2249,
			},
			model.PlatIPad: map[string]int{
				"bottom": 2253,
			},
		},
		"177": map[int8]map[string]int{
			model.PlatIPhone: map[string]int{
				"bottom": 2295,
			},
			model.PlatAndroid: map[string]int{
				"bottom": 2299,
			},
			model.PlatIPad: map[string]int{
				"bottom": 2303,
			},
		},
	}
	_bannersIndex = map[string]map[int8]map[string]int{
		"-1": map[int8]map[string]int{
			model.PlatIPhone: map[string]int{
				"top": 467,
			},
			model.PlatAndroid: map[string]int{
				"top": 631,
			},
			model.PlatIPad: map[string]int{
				"top": 771,
			},
			model.PlatIPhoneI: map[string]int{
				"top": 947,
			},
			model.PlatAndroidG: map[string]int{
				"top": 1285,
			},
			model.PlatAndroidI: map[string]int{
				"top": 1707,
			},
			model.PlatIPadI: map[string]int{
				"top": 1117,
			},
		},
		"0": map[int8]map[string]int{
			model.PlatIPhone: map[string]int{
				"top":    467,
				"center": 468,
				"bottom": 469,
			},
			model.PlatAndroid: map[string]int{
				"top":    631,
				"center": 632,
				"bottom": 633,
			},
			model.PlatIPad: map[string]int{
				"top":    771,
				"center": 772,
				"bottom": 773,
			},
			model.PlatIPhoneI: map[string]int{
				"top":    947,
				"center": 952,
				"bottom": 957,
			},
			model.PlatAndroidG: map[string]int{
				"top":    1285,
				"center": 1290,
				"bottom": 1295,
			},
			model.PlatAndroidI: map[string]int{
				"top":    1707,
				"center": 1712,
				"bottom": 1717,
			},
			model.PlatIPadI: map[string]int{
				"top":    1117,
				"center": 1122,
				"bottom": 1127,
			},
		},
		"65537": map[int8]map[string]int{
			model.PlatIPhone: map[string]int{
				"top": 482,
			},
			model.PlatAndroid: map[string]int{
				"top": 646,
			},
			model.PlatIPad: map[string]int{
				"top": 786,
			},
			model.PlatIPhoneI: map[string]int{
				"top": 1013,
			},
			model.PlatAndroidG: map[string]int{
				"top": 1351,
			},
			model.PlatAndroidI: map[string]int{
				"top": 1773,
			},
			model.PlatIPadI: map[string]int{
				"top": 1183,
			},
		},
		"13": map[int8]map[string]int{
			model.PlatIPhone: map[string]int{
				"top": 471,
			},
			model.PlatAndroid: map[string]int{
				"top": 635,
			},
			model.PlatIPad: map[string]int{
				"top": 775,
			},
			model.PlatIPhoneI: map[string]int{
				"top": 967,
			},
			model.PlatAndroidG: map[string]int{
				"top": 1305,
			},
			model.PlatAndroidI: map[string]int{
				"top": 1727,
			},
			model.PlatIPadI: map[string]int{
				"top": 1137,
			},
		},
		"1": map[int8]map[string]int{
			model.PlatIPhone: map[string]int{
				"top": 470,
			},
			model.PlatAndroid: map[string]int{
				"top": 634,
			},
			model.PlatIPad: map[string]int{
				"top": 774,
			},
			model.PlatIPhoneI: map[string]int{
				"top": 962,
			},
			model.PlatAndroidG: map[string]int{
				"top": 1300,
			},
			model.PlatAndroidI: map[string]int{
				"top": 1722,
			},
			model.PlatIPadI: map[string]int{
				"top": 1132,
			},
		},
		"3": map[int8]map[string]int{
			model.PlatIPhone: map[string]int{
				"top": 472,
			},
			model.PlatAndroid: map[string]int{
				"top": 636,
			},
			model.PlatIPad: map[string]int{
				"top": 776,
			},
			model.PlatIPhoneI: map[string]int{
				"top": 971,
			},
			model.PlatAndroidG: map[string]int{
				"top": 1309,
			},
			model.PlatAndroidI: map[string]int{
				"top": 1731,
			},
			model.PlatIPadI: map[string]int{
				"top": 1141,
			},
		},
		"129": map[int8]map[string]int{
			model.PlatIPhone: map[string]int{
				"top": 473,
			},
			model.PlatAndroid: map[string]int{
				"top": 637,
			},
			model.PlatIPad: map[string]int{
				"top": 777,
			},
			model.PlatIPhoneI: map[string]int{
				"top": 975,
			},
			model.PlatAndroidG: map[string]int{
				"top": 1313,
			},
			model.PlatAndroidI: map[string]int{
				"top": 1735,
			},
			model.PlatIPadI: map[string]int{
				"top": 1145,
			},
		},
		"4": map[int8]map[string]int{
			model.PlatIPhone: map[string]int{
				"top": 474,
			},
			model.PlatAndroid: map[string]int{
				"top": 638,
			},
			model.PlatIPad: map[string]int{
				"top": 778,
			},
			model.PlatIPhoneI: map[string]int{
				"top": 979,
			},
			model.PlatAndroidG: map[string]int{
				"top": 1317,
			},
			model.PlatAndroidI: map[string]int{
				"top": 1739,
			},
			model.PlatIPadI: map[string]int{
				"top": 1149,
			},
		},
		"36": map[int8]map[string]int{
			model.PlatIPhone: map[string]int{
				"top": 475,
			},
			model.PlatAndroid: map[string]int{
				"top": 639,
			},
			model.PlatIPad: map[string]int{
				"top": 779,
			},
			model.PlatIPhoneI: map[string]int{
				"top": 983,
			},
			model.PlatAndroidG: map[string]int{
				"top": 1321,
			},
			model.PlatAndroidI: map[string]int{
				"top": 1706,
			},
			model.PlatIPadI: map[string]int{
				"top": 1153,
			},
		},
		"160": map[int8]map[string]int{
			model.PlatIPhone: map[string]int{
				"top": 476,
			},
			model.PlatAndroid: map[string]int{
				"top": 640,
			},
			model.PlatIPad: map[string]int{
				"top": 780,
			},
			model.PlatIPhoneI: map[string]int{
				"top": 987,
			},
			model.PlatAndroidG: map[string]int{
				"top": 1325,
			},
			model.PlatAndroidI: map[string]int{
				"top": 1747,
			},
			model.PlatIPadI: map[string]int{
				"top": 1157,
			},
		},
		"119": map[int8]map[string]int{
			model.PlatIPhone: map[string]int{
				"top": 477,
			},
			model.PlatAndroid: map[string]int{
				"top": 641,
			},
			model.PlatIPad: map[string]int{
				"top": 781,
			},
			model.PlatIPhoneI: map[string]int{
				"top": 992,
			},
			model.PlatAndroidG: map[string]int{
				"top": 1330,
			},
			model.PlatAndroidI: map[string]int{
				"top": 1752,
			},
			model.PlatIPadI: map[string]int{
				"top": 1162,
			},
		},
		"155": map[int8]map[string]int{
			model.PlatIPhone: map[string]int{
				"top": 478,
			},
			model.PlatAndroid: map[string]int{
				"top": 642,
			},
			model.PlatIPad: map[string]int{
				"top": 782,
			},
			model.PlatIPhoneI: map[string]int{
				"top": 997,
			},
			model.PlatAndroidG: map[string]int{
				"top": 1335,
			},
			model.PlatAndroidI: map[string]int{
				"top": 1757,
			},
			model.PlatIPadI: map[string]int{
				"top": 1167,
			},
		},
		"5": map[int8]map[string]int{
			model.PlatIPhone: map[string]int{
				"top": 479,
			},
			model.PlatAndroid: map[string]int{
				"top": 643,
			},
			model.PlatIPad: map[string]int{
				"top": 783,
			},
			model.PlatIPhoneI: map[string]int{
				"top": 1001,
			},
			model.PlatAndroidG: map[string]int{
				"top": 1339,
			},
			model.PlatAndroidI: map[string]int{
				"top": 1761,
			},
			model.PlatIPadI: map[string]int{
				"top": 1171,
			},
		},
		"23": map[int8]map[string]int{
			model.PlatIPhone: map[string]int{
				"top": 480,
			},
			model.PlatAndroid: map[string]int{
				"top": 644,
			},
			model.PlatIPad: map[string]int{
				"top": 784,
			},
			model.PlatIPhoneI: map[string]int{
				"top": 1005,
			},
			model.PlatAndroidG: map[string]int{
				"top": 1343,
			},
			model.PlatAndroidI: map[string]int{
				"top": 1765,
			},
			model.PlatIPadI: map[string]int{
				"top": 1175,
			},
		},
		"11": map[int8]map[string]int{
			model.PlatIPhone: map[string]int{
				"top": 481,
			},
			model.PlatAndroid: map[string]int{
				"top": 645,
			},
			model.PlatIPad: map[string]int{
				"top": 785,
			},
			model.PlatIPhoneI: map[string]int{
				"top": 1009,
			},
			model.PlatAndroidG: map[string]int{
				"top": 1347,
			},
			model.PlatAndroidI: map[string]int{
				"top": 1769,
			},
			model.PlatIPadI: map[string]int{
				"top": 1179,
			},
		},
		"165": map[int8]map[string]int{
			model.PlatIPhone: map[string]int{
				"top": 1643,
			},
			model.PlatAndroid: map[string]int{
				"top": 1639,
			},
			model.PlatIPad: map[string]int{
				"top": 1647,
			},
			model.PlatIPhoneI: map[string]int{
				"top": 1643,
			},
			model.PlatAndroidG: map[string]int{
				"top": 1639,
			},
			model.PlatIPadI: map[string]int{
				"top": 1647,
			},
		},
		"167": map[int8]map[string]int{
			model.PlatIPhone: map[string]int{
				"top": 1950,
			},
			model.PlatAndroid: map[string]int{
				"top": 1952,
			},
			model.PlatIPad: map[string]int{
				"top": 1951,
			},
		},
		"181": map[int8]map[string]int{
			model.PlatIPhone: map[string]int{
				"top": 2245,
			},
			model.PlatAndroid: map[string]int{
				"top": 2249,
			},
			model.PlatIPad: map[string]int{
				"top": 2253,
			},
		},
		"177": map[int8]map[string]int{
			model.PlatIPhone: map[string]int{
				"top": 2295,
			},
			model.PlatAndroid: map[string]int{
				"top": 2299,
			},
			model.PlatIPad: map[string]int{
				"top": 2303,
			},
		},
	}
	_bannersPlat = map[int8]string{
		model.PlatIPhone:   "467,482,471,470,472,473,474,475,476,477,478,479,480,481,1643,1950,2245,2295",
		model.PlatAndroid:  "631,646,635,634,636,637,638,639,640,641,642,643,644,645,1639,1952,2249,2299",
		model.PlatIPad:     "771,786,775,774,776,777,778,779,780,781,782,783,784,785,1647,1951,2253,2303",
		model.PlatIPhoneI:  "947,1013,967,962,971,975,979,983,987,992,997,1001,1005,1009,1643",
		model.PlatAndroidG: "1285,1351,1305,1300,1309,1313,1317,1321,1325,1330,1335,1339,1343,1347,1639",
		model.PlatAndroidI: "1707,1773,1727,1722,1731,1735,1739,1706,1747,1752,1757,1761,1765,1769,1639",
		model.PlatIPadI:    "1117,1183,1137,1132,1141,1145,1149,1153,1157,1162,1167,1171,1175,1179,1647",
	}
)

// getBanners get banners by plat, build channel, ip.
func (s *Service) getBanners(c context.Context, plat int8, build int, module, channel, ip string, resbs map[int][]*resource.Banner, isIndex bool) (res map[string][]*banner.Banner) {
	var (
		bannerIds = _banners
	)
	if isIndex {
		bannerIds = _bannersIndex
	}
	res = map[string][]*banner.Banner{}
	for pos, bID := range bannerIds[module][plat] {
		if rbs, ok := resbs[bID]; ok {
			var bs []*banner.Banner
			for _, rb := range rbs {
				b := &banner.Banner{}
				b.ResChangeBanner(rb)
				bs = append(bs, b)
			}
			res[pos] = bs
		}
	}
	return
}

// resBannersplat
func (s *Service) resBanners(c context.Context, plat int8, build int, mid int64, resIDStr, channel, ip, buvid, network, mobiApp, device, adExtra string) (res map[int][]*resource.Banner) {
	var (
		plm  = s.bannerCache[plat] // operater banner
		err  error
		isAd = true
	)
	if plat == model.PlatAndroid && build <= 430000 {
		isAd = false
	}
	if res, err = s.res.ResBanner(c, plat, build, mid, resIDStr, channel, ip, buvid, network, mobiApp, device, adExtra, isAd); err != nil || len(res) == 0 {
		log.Error("s.res.ResBanner is null or err(%v)", err)
		res = plm
		return
	}
	return
}
