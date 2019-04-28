package archive

import (
	"context"
	a "go-common/app/admin/main/videoup/model/archive"
)

// EditRules fn
func (s *Service) EditRules(c context.Context, white int, state int8, lotteryBind bool) (rules map[string]bool) {
	var (
		exist           bool
		rulesByArcState = make(map[int8]map[string]bool)
	)
	groupAllCan := map[string]bool{
		"tid":          true,
		"title":        true,
		"tag":          true,
		"desc":         true,
		"dynamic":      true,
		"del_video":    true,
		"elec":         true,
		"add_video":    true,
		"dtime":        true,
		"source":       true,
		"no_reprint":   true,
		"cover":        true,
		"copyright":    true,
		"mission_tag":  true,
		"bind_lottery": false,
	}
	groupAllForbid := map[string]bool{
		"tid":          false,
		"title":        false,
		"tag":          false,
		"desc":         false,
		"dynamic":      false,
		"del_video":    false,
		"elec":         false,
		"add_video":    false,
		"dtime":        false,
		"source":       false,
		"no_reprint":   false,
		"cover":        false,
		"copyright":    false,
		"mission_tag":  false,
		"bind_lottery": false,
	}
	groupForbidTidAndCopyright := map[string]bool{
		"tid":          false,
		"title":        true,
		"tag":          true,
		"desc":         true,
		"dynamic":      true,
		"del_video":    true,
		"elec":         true,
		"add_video":    true,
		"dtime":        true,
		"source":       true,
		"no_reprint":   true,
		"cover":        true,
		"copyright":    false,
		"mission_tag":  false,
		"bind_lottery": false,
	}
	groupForbidTidAndCopyrightDtime := map[string]bool{
		"tid":          false,
		"title":        true,
		"tag":          true,
		"desc":         true,
		"dynamic":      true,
		"del_video":    true,
		"elec":         true,
		"add_video":    true,
		"dtime":        false,
		"source":       true,
		"no_reprint":   true,
		"cover":        true,
		"copyright":    false,
		"mission_tag":  false,
		"bind_lottery": false,
	}
	rulesByArcState[a.StateOrange] = groupForbidTidAndCopyrightDtime
	rulesByArcState[a.StateOpen] = groupForbidTidAndCopyrightDtime
	rulesByArcState[a.StateForbidWait] = groupForbidTidAndCopyright
	rulesByArcState[a.StateForbidAdminDelay] = groupForbidTidAndCopyright
	rulesByArcState[a.StateForbidSubmit] = groupForbidTidAndCopyright
	rulesByArcState[a.StateForbidUserDelay] = groupForbidTidAndCopyrightDtime
	rulesByArcState[a.StateForbidXcodeFail] = groupForbidTidAndCopyright
	rulesByArcState[a.StateForbidPolice] = groupAllForbid
	rulesByArcState[a.StateForbidLock] = groupAllForbid
	rulesByArcState[a.StateForbidFackLock] = groupAllForbid
	rulesByArcState[a.StateForbidUpDelete] = groupAllForbid
	rulesByArcState[a.StateForbitUpLoad] = groupForbidTidAndCopyright
	rulesByArcState[a.StateForbidOnlyComment] = groupForbidTidAndCopyright
	rulesByArcState[a.StateForbidDispatch] = groupForbidTidAndCopyright
	rulesByArcState[a.StateForbidFixing] = groupForbidTidAndCopyright
	rulesByArcState[a.StateForbidStorageFail] = groupForbidTidAndCopyright
	rulesByArcState[a.StateForbidWaitXcode] = groupForbidTidAndCopyright
	rulesByArcState[a.StateForbidTmpRecicle] = groupAllCan
	rulesByArcState[a.StateForbidRecycle] = groupAllCan

	rulesByArcState[a.StateForbidFixed] = groupForbidTidAndCopyrightDtime
	rulesByArcState[a.StateForbidLater] = groupForbidTidAndCopyrightDtime
	rulesByArcState[a.StateForbidPatched] = groupForbidTidAndCopyrightDtime

	if rules, exist = rulesByArcState[state]; exist {
		if white == 0 {
			rules["add_video"] = false
		}
		rules["bind_lottery"] = lotteryBind
	}
	return
}
