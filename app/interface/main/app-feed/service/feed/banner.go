package feed

import (
	"context"
	resource "go-common/app/service/main/resource/model"
	"strconv"

	"go-common/app/interface/main/app-card/model/card/banner"
	"go-common/app/interface/main/app-feed/model"
)

var (
	_banners = map[int8]int{
		model.PlatIPhoneB:  467,
		model.PlatIPhone:   467,
		model.PlatAndroid:  631,
		model.PlatIPad:     771,
		model.PlatIPhoneI:  947,
		model.PlatAndroidG: 1285,
		model.PlatAndroidI: 1707,
		model.PlatIPadI:    1117,
	}
)

// banners get banners by plat, build channel, ip.
func (s *Service) banners(c context.Context, plat int8, build int, mid int64, buvid, network, mobiApp, device, openEvent, adExtra, hash string) (bs []*banner.Banner, version string, err error) {
	plat = model.PlatAPPBuleChange(plat)
	var (
		rscID = _banners[plat]
		bm    map[int][]*resource.Banner
	)
	if bm, version, err = s.rsc.Banner(c, plat, build, mid, strconv.Itoa(rscID), "", buvid, network, mobiApp, device, true, openEvent, adExtra, hash); err != nil {
		return
	}
	for _, rb := range bm[rscID] {
		b := &banner.Banner{}
		b.Change(rb)
		bs = append(bs, b)
	}
	return
}
