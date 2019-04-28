package app

import (
	"context"
	"encoding/json"
	"fmt"
	"go-common/app/interface/main/creative/conf"
	appMdl "go-common/app/interface/main/creative/model/app"
	resMdl "go-common/app/interface/main/creative/model/resource"
	"go-common/library/log"
)

// Portals for app portal config.
func (s *Service) Portals(c context.Context, mid int64, artAuthor, build, ty int, plat string, resMdlPlat int8) (pts []*appMdl.Portal, err error) {
	var apms []*appMdl.PortalMeta
	if ty == appMdl.PortalIntro {
		if resMdlPlat == resMdl.PlatIPad {
			plat = "ipad"
		}
		apms = s.PortalIntro
	} else if ty == appMdl.PortalNotice {
		apms = s.PortalNotice
	}
	pts = make([]*appMdl.Portal, 0, len(apms))
	for _, v := range apms {
		shouldAppend := false
		if v.Title == appMdl.MyArticle || v.Title == appMdl.OpenArticle {
			if artAuthor == 1 && v.Title == appMdl.MyArticle && v.AllowMaterial(plat, build) {
				shouldAppend = true
			} else if artAuthor == 0 && v.Title == appMdl.OpenArticle && v.AllowMaterial(plat, build) {
				shouldAppend = true
			}
		} else if v.AllowMaterial(plat, build) {
			shouldAppend = true
		}
		if ty == appMdl.PortalIntro && !s.AllowWhiteCategory(mid, v) {
			shouldAppend = false
		}
		if shouldAppend {
			pts = append(pts, &appMdl.Portal{
				Icon:     v.Icon,
				Title:    v.Title,
				Pos:      v.Pos,
				URL:      v.URL,
				New:      v.Mark,
				More:     v.More,
				SubTitle: v.SubTitle,
				MTime:    v.MTime,
			})
		}
	}
	return
}

// AllowWhiteCategory fn
func (s *Service) AllowWhiteCategory(mid int64, p *appMdl.PortalMeta) (ret bool) {
	if len(p.WhiteExps) == 0 {
		return true
	}
	var (
		percentOK, groupOK bool
	)
	for _, v := range p.WhiteExps {
		// for WhitePercentType
		if v.TP == appMdl.WhitePercentType {
			mod := mid % 100
			fac := v.Value * 10
			switch v.Value {
			case appMdl.WhitePercentV10, appMdl.WhitePercentV20, appMdl.WhitePercentV50:
				log.Warn("whiteexp WhitePercentV p(%+v), mid(%d), percentOK(%+v)", p.WhiteExps, mid, percentOK)
				if mod <= fac {
					percentOK = true
				}
			case appMdl.WhitePercentV00:
				percentOK = true
			default:
				percentOK = false
			}
			log.Warn("whiteexp WhitePercentType p(%+v), mid(%d), percentOK(%+v)", p.WhiteExps, mid, percentOK)
		}
		// for WhiteGroupType
		if v.TP == appMdl.WhiteGroupType {
			if midMaps, ok := s.p.AppWhiteMidsByGroups[v.Value]; ok {
				isWhite := false
				for _, m := range midMaps {
					if m == mid {
						isWhite = true
						break
					}
				}
				if isWhite {
					groupOK = true
				}
			}
			log.Warn("whiteexp WhiteGroupType p(%+v), mid(%d), groupOK(%+v)", p.WhiteExps, mid, groupOK)
		}
	}
	return percentOK || groupOK
}

// UploadMaterial fn
func (s *Service) UploadMaterial(c context.Context, aid int64, editors []*appMdl.Editor) (err error) {
	if len(editors) < 0 || aid < 0 {
		return nil
	}
	for _, editor := range editors {
		materialDatas := s.splitMaterials(aid, editor)
		for _, data := range materialDatas {
			s.app.AddMaterialData(c, data)
		}
	}
	return
}

func (s *Service) splitMaterials(aid int64, pe *appMdl.Editor) (pds []*appMdl.EditorData) {
	pds = []*appMdl.EditorData{}
	if pe.CID <= 0 {
		return
	}
	subStr := pe.ParseSubtitles()
	if len(subStr) > 0 {
		pds = append(pds, &appMdl.EditorData{
			AID:  aid,
			CID:  pe.CID,
			Type: appMdl.TypeSubtitle,
			Data: subStr,
		})
	}
	fonStr := pe.ParseFonts()
	if len(fonStr) > 0 {
		pds = append(pds, &appMdl.EditorData{
			AID:  aid,
			CID:  pe.CID,
			Type: appMdl.TypeFont,
			Data: fonStr,
		})
	}
	filStr := pe.ParseFilters()
	if len(filStr) > 0 {
		pds = append(pds, &appMdl.EditorData{
			AID:  aid,
			CID:  pe.CID,
			Type: appMdl.TypeFilter,
			Data: filStr,
		})
	}
	bgmStr := pe.ParseBgms()
	if len(bgmStr) > 0 {
		pds = append(pds, &appMdl.EditorData{
			AID:  aid,
			CID:  pe.CID,
			Type: appMdl.TypeBGM,
			Data: bgmStr,
		})
	}
	stiStr := pe.ParseStickers()
	if len(stiStr) > 0 {
		pds = append(pds, &appMdl.EditorData{
			AID:  aid,
			CID:  pe.CID,
			Type: appMdl.TypeSticker,
			Data: stiStr,
		})
	}
	vstStr := pe.ParseVStickers()
	if len(vstStr) > 0 {
		pds = append(pds, &appMdl.EditorData{
			AID:  aid,
			CID:  pe.CID,
			Type: appMdl.TypeVideoupSticker,
			Data: vstStr,
		})
	}
	traStr := pe.ParseTransitions()
	if len(traStr) > 0 {
		pds = append(pds, &appMdl.EditorData{
			AID:  aid,
			CID:  pe.CID,
			Type: appMdl.TypeTransition,
			Data: traStr,
		})
	}
	cooStr := pe.ParseCooperates()
	if len(cooStr) > 0 {
		pds = append(pds, &appMdl.EditorData{
			AID:  aid,
			CID:  pe.CID,
			Type: appMdl.TypeCooperate,
			Data: cooStr,
		})
	}
	theStr := pe.ParseThemes()
	if len(theStr) > 0 {
		pds = append(pds, &appMdl.EditorData{
			AID:  aid,
			CID:  pe.CID,
			Type: appMdl.TypeTheme,
			Data: theStr,
		})
	}
	switchMap := make(map[string]int8)
	if pe.AudioRecord == 1 {
		switchMap["audio_record"] = 1
	}
	if pe.Camera == 1 {
		switchMap["camera"] = 1
	}
	if pe.Speed == 1 {
		switchMap["speed"] = 1
	}
	if pe.CameraRotate == 1 {
		switchMap["camera_rotate"] = 1
	}
	if len(switchMap) > 0 {
		switchMapData, _ := json.Marshal(switchMap)
		pds = append(pds, &appMdl.EditorData{
			AID:  aid,
			CID:  pe.CID,
			Type: appMdl.TypeSwitch,
			Data: string(switchMapData),
		})
	}
	return
}

// Icons fn
func (s *Service) Icons() (res *conf.AppIcon) {
	return s.c.AppIcon
}

// H5Page fn
func (s *Service) H5Page() (res *conf.H5Page) {
	return s.c.H5Page
}

// BlockIntros fn
func (s *Service) BlockIntros(build int, platStr string) (res map[string]string) {
	res = make(map[string]string)
	res["college"] = s.c.H5Page.CreativeCollege
	if (build >= 5350000 && platStr == "android") || (build >= 8240 && platStr == "ios") {
		res["college"] = fmt.Sprintf("%s?from=rcmd&navhide=1", res["college"])
	} else {
		res["college"] = fmt.Sprintf("%s?from=rcmd", res["college"])
	}
	return
}
