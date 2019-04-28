package service

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"go-common/app/interface/main/push/model"
	pushmdl "go-common/app/service/main/push/model"
	"go-common/library/log"
)

// PubReport pub report.
func (s *Service) PubReport(c context.Context, r *pushmdl.Report) (err error) {
	err = s.dao.PubReport(c, r)
	return
}

// ReportOld report old version app
func (s *Service) ReportOld(ctx context.Context, token, buvid, version string, mid int64, pid, timezone int) (err error) {
	platform := translatePlatform(pid)
	build := ver2build(version, platform)
	if build == 0 {
		return
	}
	// 接新上报后的版本不再使用老的上报数据了
	switch platform {
	case pushmdl.PlatformXiaomi:
		// version 5.16
		if build >= 516000 {
			return
		}
	case pushmdl.PlatformIPhone:
		// version 5.16
		if build >= 6140 {
			return
		}
	case pushmdl.PlatformIPad:
		// version 1.50
		if build >= 12040 {
			return
		}
	default:
		// 未识别的平台
		return
	}
	if platform == pushmdl.PlatformIPad && build < 10000 {
		platform = pushmdl.PlatformIPhone
	}
	r := &pushmdl.Report{
		APPID:        pushmdl.APPIDBBPhone,
		PlatformID:   platform,
		Mid:          mid,
		Buvid:        buvid,
		Build:        build,
		DeviceToken:  token,
		TimeZone:     timezone,
		NotifySwitch: 1,
	}
	err = s.dao.PubReport(ctx, r)
	log.Info("pub old report(%+v)", r)
	return
}

func translatePlatform(platformID int) int {
	switch platformID {
	case model.OldPlatformIPhone, model.OldPlatformIPad:
		return pushmdl.PlatformIPhone
	case model.OldPlatformAndroid, model.OldPlatformAndroidNow:
		return pushmdl.PlatformXiaomi
	case model.OldPlatformIPadHD:
		return pushmdl.PlatformIPad
	}
	return pushmdl.PlatformUnknown
}

var buildRegex, _ = regexp.Compile(`\((\d+)\)`)

func ver2build(versionAndBuild string, platform int) (res int) {
	version := versionAndBuild
	// example: 5.12.1(6050) remove '(' suffix
	i := strings.Index(version, "(")
	if i != -1 {
		version = version[0:i]
	}
	// example: 5.14.0-preview remove '-' suffix
	i = strings.Index(version, "-")
	if i != -1 {
		version = version[0:i]
	}
	switch platform {
	case pushmdl.PlatformIPhone:
		res = model.VersionsIPhone[version]
	case pushmdl.PlatformIPad:
		res = model.VersionsIPad[version]
	default:
		p := strings.Split(version, ".")
		if len(p) < 3 {
			return
		}
		res, _ = strconv.Atoi(p[0] + p[1] + fmt.Sprintf("%03s", p[2]))
	}
	if res == 0 {
		// match as 2_5.10(5960)
		matches := buildRegex.FindSubmatch([]byte(versionAndBuild))
		if len(matches) > 1 {
			res, _ = strconv.Atoi(string(matches[1]))
		}
	}
	return
}
