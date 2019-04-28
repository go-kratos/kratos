package region

import (
	"context"
	"strconv"
	"time"

	"go-common/app/interface/main/app-show/model"
	"go-common/app/interface/main/app-show/model/banner"
	"go-common/app/interface/main/app-show/model/region"
	"go-common/app/service/main/archive/api"
	"go-common/library/log"
)

var (
	_auditRids = map[int8]map[int]struct{}{
		model.PlatIPad: map[int]struct{}{
			153: struct{}{},
			168: struct{}{},
			169: struct{}{},
			170: struct{}{},
			33:  struct{}{},
			32:  struct{}{},
			51:  struct{}{},
			152: struct{}{},
			37:  struct{}{},
			178: struct{}{},
			179: struct{}{},
			180: struct{}{},
			147: struct{}{},
			145: struct{}{},
			146: struct{}{},
			83:  struct{}{},
			185: struct{}{},
			187: struct{}{},
			13:  struct{}{},
			167: struct{}{},
			177: struct{}{},
			23:  struct{}{},
			1:   struct{}{},
			160: struct{}{},
			119: struct{}{},
			155: struct{}{},
			165: struct{}{},
			5:   struct{}{},
			181: struct{}{},
		},
		model.PlatIPhone: map[int]struct{}{
			153: struct{}{},
			168: struct{}{},
			169: struct{}{},
			170: struct{}{},
			33:  struct{}{},
			32:  struct{}{},
			51:  struct{}{},
			152: struct{}{},
			37:  struct{}{},
			178: struct{}{},
			179: struct{}{},
			180: struct{}{},
			147: struct{}{},
			145: struct{}{},
			146: struct{}{},
			83:  struct{}{},
			185: struct{}{},
			187: struct{}{},
			1:   struct{}{},
			160: struct{}{},
			119: struct{}{},
			155: struct{}{},
			165: struct{}{},
			5:   struct{}{},
			181: struct{}{},
		},
	}
)

// Audit region data list.
func (s *Service) auditRegion(mobiApp string, plat int8, build, rid int) (isAudit bool) {
	if plats, ok := s.auditCache[mobiApp]; ok {
		if _, ok = plats[build]; ok {
			if rids, ok := _auditRids[plat]; ok {
				if _, ok = rids[rid]; ok {
					return true
				}
			}
		}
	}
	return false
}

func (s *Service) loadAuditCache() {
	as, err := s.adt.Audits(context.TODO())
	if err != nil {
		log.Error("s.adt.Audits error(%v)", err)
		return
	}
	s.auditCache = as
}

// Audit check audit plat and ip, then return audit data.
func (s *Service) Audit(c context.Context, mobiApp string, plat int8, build, rid int, isShow bool) (res *region.Show, ok bool) {
	if plats, ok := s.auditCache[mobiApp]; ok {
		if _, ok = plats[build]; ok {
			res = s.auditData(c, plat, rid, auditShowAids)
			if isShow {
				res.Banner = audirRegionBanners[rid]
			}
			return res, true
		}
	}
	return nil, false
}

// AuditChild check audit plat and ip, then return audit data.
func (s *Service) AuditChild(c context.Context, mobiApp, order string, plat int8, build, rid, tid int) (res *region.Show, ok bool) {
	if plats, ok := s.auditCache[mobiApp]; ok {
		if _, ok = plats[build]; ok {
			res = s.auditData(c, plat, rid, auditChildShowAids)
			res.New = s.auditRegionRPCList(c, rid, 1, 8)
			return res, true
		}
	}
	return nil, false
}

// AuditChildList check audit plat and ip, then return audit data.
func (s *Service) AuditChildList(c context.Context, mobiApp, order string, plat int8, build, rid, tid, pn, ps int) (res []*region.ShowItem, ok bool) {
	if plats, ok := s.auditCache[mobiApp]; ok {
		if _, ok = plats[build]; ok {
			res = s.auditRegionChildList(c, rid, tid, pn, ps)
			return res, true
		}
	}
	return nil, false
}

// auditData some data for audit.
func (s *Service) auditData(c context.Context, p int8, rid int, auditAids map[int][]int64) (res *region.Show) {
	aids := auditAids[rid]
	// archive
	as, err := s.arc.ArchivesPB(c, aids)
	if err != nil {
		log.Error("s.arc.ArchivesPB error(%v)", err)
		as = map[int64]*api.Arc{}
	}
	res = &region.Show{}
	for _, aid := range aids {
		if aid == 0 {
			continue
		}
		item := &region.ShowItem{}
		item.Goto = model.GotoAv
		item.Param = strconv.FormatInt(aid, 10)
		item.URI = model.FillURI(item.Goto, item.Param, nil)
		if a, ok := as[aid]; ok {
			item.Title = a.Title
			item.Cover = a.Pic
			item.Name = a.Author.Name
			item.Play = int(a.Stat.View)
			item.Danmaku = int(a.Stat.Danmaku)
			item.Reply = int(a.Stat.Reply)
			item.Fav = int(a.Stat.Fav)
		}
		res.Recommend = append(res.Recommend, item)
	}
	return
}

func (s *Service) auditRegionChildList(c context.Context, rid, tid, pn, ps int) (res []*region.ShowItem) {
	if tid == 0 {
		arcs, _, err := s.arc.RanksArcs(c, rid, pn, ps)
		if err != nil {
			log.Error("s.rcmmnd.RegionArcList(%d, %d, %d, %d) error(%v)", rid, pn, ps, err)
			return
		}
		res = s.fromArchivesPBOsea(arcs, false)
	} else {
		as, err := s.tag.NewArcs(c, rid, tid, pn, ps, time.Now())
		if err != nil {
			log.Error("s.tag.NewArcs(%d, %d) error(%v)", rid, tid, err)
			return
		}
		res = s.fromAidsOsea(c, as, false, false, 0)
	}
	return
}

func (s *Service) auditRegionRPCList(c context.Context, rid, pn, ps int) (res []*region.ShowItem) {
	arcs, err := s.arc.RankTopArcs(c, rid, pn, ps)
	if err != nil {
		log.Error("s.arc.RankTopArcs(%d) error(%v)", rid, err)
		return
	}
	res, _ = s.fromArchivesPB(arcs)
	return
}

var (
	auditShowAids = map[int][]int64{
		// rid
		1:   []int64{575891, 744286, 663583, 666946, 559050, 744299},
		3:   []int64{881693, 756287, 785484, 402851, 887618, 853895},
		4:   []int64{861290, 861306, 861410, 861538, 861711, 861945},
		5:   []int64{791621, 795406, 797933, 800658, 832103, 833520},
		11:  []int64{1961205, 2028734},
		13:  []int64{2434272, 7408756, 2222558, 845204, 862063, 845034},
		36:  []int64{834839, 838077, 872364, 852955, 877423, 881182},
		119: []int64{638240, 1959692, 78287, 1979757},
		129: []int64{966192, 936016, 1958897, 886841},
	}

	auditChildShowAids = map[int][]int64{
		20:  []int64{936016, 886841, 1773160, 1958897, 1406019, 1935680, 1976153, 1985297, 1984555, 1964367, 29013765, 27379226, 25886650, 27684044, 20203945},
		21:  []int64{689694, 829135, 743922, 876565, 690522, 686220, 286616, 339727, 668054, 288602},
		22:  []int64{1911041, 1976535, 913421},
		24:  []int64{258271, 462832, 430248},
		25:  []int64{190257, 432195},
		26:  []int64{638240, 1959692, 78287, 1979757},
		27:  []int64{775898, 199852, 539880, 2469560, 306718, 2460323, 851414, 2471090, 591021, 286678},
		28:  []int64{221107, 221106, 884789, 364379, 465230, 26437, 29009413, 28965015, 28087847, 27837553, 24691347},
		29:  []int64{1984330, 1966586, 1984971, 28935962, 28818825, 26514923, 23288906, 18043554},
		30:  []int64{308040, 850424, 360940, 482844, 887861, 539600, 869576, 400161, 644935, 333069, 28659609, 24929108, 23068834, 26659364, 25386207},
		31:  []int64{1968681, 1986904, 1986802, 2473751, 2473083, 24910218, 25409335, 25043881, 27384682, 23474776},
		37:  []int64{1968901, 1969254, 1971484},
		47:  []int64{364103, 621797, 557774, 620545, 291630, 853831, 627451, 789570, 582598, 666971},
		54:  []int64{2294239, 2210977, 21297755, 21678914, 22000250, 19929241, 18039794},
		59:  []int64{1969748, 1966643, 1964781, 1969527, 25814802, 25991412, 26577780, 23922472, 28934467},
		71:  []int64{1986816, 1985288, 1986516, 1985717},
		75:  []int64{200595, 721477, 668533, 803294, 708986, 581574, 588820, 718877, 6336, 592586},
		76:  []int64{800617, 817625, 853774, 808176, 810174, 737783, 792994, 811825, 794302, 817814},
		95:  []int64{880857, 26317616, 26697725, 24670946, 13562204, 24136940},
		96:  []int64{2313588, 2314237, 2316089, 28917042, 20177394, 27839524, 25866526, 22021244},
		98:  []int64{875076, 873174, 580862, 289024, 28868117, 26404621, 17229132, 28810408, 27710623},
		122: []int64{1986932, 1985610, 22034719, 19980487, 19841525, 23328696, 29249512},
		124: []int64{842756, 875624, 880558, 862316, 876708, 883418, 403120, 855131, 876867, 833785, 29064835, 27464818, 28055879, 18081681, 22968172},
		126: []int64{1636345, 1985956, 1975358, 1982533},
		127: []int64{1743126, 1625784, 1986533, 1727650},
		128: []int64{2031210, 2034983, 1916941, 2030610, 2015734, 2016150, 1982576, 2039658, 1981156, 1964927},
		130: []int64{1984887, 1985685, 1985886, 25276379, 17119215, 24949925, 25058065, 2929013},
		131: []int64{1980280, 1975409},
		137: []int64{2316922, 2318219},
		138: []int64{2317125, 2317283, 2315385, 2317914, 2317194},
		153: []int64{2429129, 7408756, 2426501, 2425990, 2387429, 2425770, 2219211, 878914, 880182, 2240189},
		154: []int64{1960912, 1984928, 29240625, 26192654, 24211477, 23746281, 23871787},
		156: []int64{28960012, 26624032, 25520347, 23567968, 23706035},
		17:  []int64{28989880, 25158325, 23947116, 27052563, 24237900},
		171: []int64{29027059, 26486853, 19641793, 26432920, 27107785},
		172: []int64{28280704, 28667051, 27462689, 22870782, 17703340},
		65:  []int64{28386938, 28832756, 27894258, 23401066, 24434703},
		173: []int64{28938351, 28945212, 28149415, 17717106, 26227357},
		121: []int64{27148774, 24729449, 24544576, 23651344, 21672258},
		136: []int64{26033272, 26422598, 26804826, 25773023, 22961192},
		19:  []int64{28942929, 26325161, 24502096, 22364954, 19289951},
		39:  []int64{29120624, 26300313, 25504214, 9447066, 20786390},
		176: []int64{28811447, 27569816, 11984355, 10788852, 29346662},
	}

	audirRegionBanners = map[int]map[string][]*banner.Banner{
		1: map[string][]*banner.Banner{
			"top": []*banner.Banner{
				&banner.Banner{
					Title: "四月再见",
					Image: "http://i0.hdslb.com/bfs/archive/8bbc82a30720f8c2cdcca1576e25917f7bbdfb96.jpg",
					Hash:  "db6e4dcc120fcd954a5c2d454b618f09",
					URI:   "http://www.bilibili.com/video/av2471080/",
				},
			},
		},
		3: map[string][]*banner.Banner{
			"top": []*banner.Banner{
				&banner.Banner{
					Title: "【洛天依原创】夜寂",
					Image: "http://i0.hdslb.com/bfs/archive/6fa8a51c9adf6eeda36636ed7fffae5b1888c154.jpg",
					Hash:  "c925b57dbaa1198e8cdedc84c4781313",
					URI:   "http://www.bilibili.com/video/av2126431/",
				},
			},
		},
		4: map[string][]*banner.Banner{
			"top": []*banner.Banner{
				&banner.Banner{
					Title: "释放内心中的熊孩子吧",
					Image: "http://i0.hdslb.com/bfs/archive/b94d053b289184d498236de100af383bd25cfb13.jpg",
					Hash:  "f246e2f10d19e30dc7311c9f1ee8385e",
					URI:   "http://www.bilibili.com/video/av2459834/",
				},
			},
		},
		5: map[string][]*banner.Banner{
			"top": []*banner.Banner{
				&banner.Banner{
					Title: "巅峰料理对决~~~",
					Image: "http://i0.hdslb.com/bfs/archive/a91501598fb180f61f234e31f94731b74235b461.jpg",
					Hash:  "8664c0bd979f62c02cf6711ac0a55219",
					URI:   "http://www.bilibili.com/video/av2607073/?br",
				},
				&banner.Banner{
					Title: "首轮淘汰赛，谁将会离开",
					Image: "http://i0.hdslb.com/bfs/archive/ae6d1ef420a5bfdca969a31ecd7449384cfcd580.jpg",
					Hash:  "346b456dea567a32a11a9427ebe3246f",
					URI:   "http://www.bilibili.com/video/av2609994/?br",
				},
				&banner.Banner{
					Title: "孙红雷罗志祥乡村女装秀",
					Image: "http://i0.hdslb.com/bfs/archive/3ea355584e26376df6cccbb1a2574f03f7e0a41d.jpg",
					Hash:  "45e70e0f72dd1fa4621e13986df03b30",
					URI:   "http://www.bilibili.com/video/av2598211/?br",
				},
				&banner.Banner{
					Title: "帅哥萌妹齐驾到 HK君强势出境！",
					Image: "http://i0.hdslb.com/bfs/archive/ae6d1ef420a5bfdca969a31ecd7449384cfcd580.jpg",
					Hash:  "7e7f8fa57dfffa0ca141e12f43088851",
					URI:   "http://www.bilibili.com/video/av2598658/?br",
				},
				&banner.Banner{
					Title: "林丹谢霆锋上演锅铲大战 容祖儿情绪崩溃大哭",
					Image: "http://i0.hdslb.com/bfs/archive/959930c687c172a28d9f24e6a53aceb7fca4f728.jpg",
					Hash:  "512a11b81e053d52c2ad836c453c18ad",
					URI:   "http://www.bilibili.com/video/av2588446/?br",
				},
				&banner.Banner{
					Title: "【绅士大概一分钟】尽情舞蹈吧少年",
					Image: "http://i0.hdslb.com/bfs/archive/40878881827105f576e0346932cf693c3033e1ca.jpg",
					Hash:  "147763f8003e6f88385cd438e8b6c7e4",
					URI:   "http://www.bilibili.com/video/av2614367/?br",
				},
			},
		},
		11: map[string][]*banner.Banner{
			"top": []*banner.Banner{
				&banner.Banner{
					Title: "亚瑟王",
					Image: "http://i0.hdslb.com/bfs/archive/60b1339a3eeb8d0de287b7c305e0671082946bfc.jpg",
					Hash:  "634684fe0fd4fb3b7501daf9a9b4ab5d",
					URI:   "http://www.bilibili.com/video/av2128802/",
				},
			},
		},
		13: map[string][]*banner.Banner{
			"top": []*banner.Banner{
				&banner.Banner{
					Title: "少女终末旅行",
					Image: "http://i0.hdslb.com/bfs/archive/c0c33be60527c377277048c04ee222c9ec76a82c.jpg",
					Hash:  "2e772627851aa7da8a75f7b5403a5ed3",
					URI:   "http://bangumi.bilibili.com/anime/6463",
				},
			},
		},
		23: map[string][]*banner.Banner{
			"top": []*banner.Banner{
				&banner.Banner{
					Title: "像素大战",
					Image: "http://i0.hdslb.com/bfs/archive/0c8f1e05dfdba3b58fc15159523d0ccceed1e9ac.jpg",
					Hash:  "f62bbb0578beb4bc63550d5632960480",
					URI:   "http://www.bilibili.com/video/av2124091/",
				},
			},
		},
		36: map[string][]*banner.Banner{
			"top": []*banner.Banner{
				&banner.Banner{
					Title: "梦回仙剑",
					Image: "http://i0.hdslb.com/bfs/archive/76145de97ff917a6e603009376f4ca174dd4ed51.jpg",
					Hash:  "6334c8c08ffa8eaa6c13f8c14bd0fae0",
					URI:   "http://www.bilibili.com/video/av2448057/",
				},
			},
		},
		119: map[string][]*banner.Banner{
			"top": []*banner.Banner{
				&banner.Banner{
					Title: "魔都地铁偷走了重要的东西",
					Image: "http://i0.hdslb.com/bfs/archive/f84dd391351a00d69cfb44616c1a64419ad4611c.jpg",
					Hash:  "5de70c2b24a155d7958943b85bf8facc",
					URI:   "http://www.bilibili.com/video/av2106417/",
				},
			},
		},
		129: map[string][]*banner.Banner{
			"top": []*banner.Banner{
				&banner.Banner{
					Title: "元气少女",
					Image: "http://i0.hdslb.com/bfs/archive/014c0793bdaf5930e0edca54755da3c25eafcb2e.jpg",
					Hash:  "46e48cee8d5b344af3a41636e96a60ca",
					URI:   "http://www.bilibili.com/video/av2448328/",
				},
			},
		},
	}
)
