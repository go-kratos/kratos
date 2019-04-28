package archive

import (
	arcwar "go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
)

// validView distinguishes whether an view is valid
func validView(vp *arcwar.ViewReply, checkAttr bool) (valid bool) {
	if vp == nil {
		return
	}
	if vp.Arc == nil {
		return
	}
	if vp.Arc.Aid == 0 {
		return
	}
	if len(vp.Pages) == 0 {
		return
	}
	if checkAttr {
		if vp.Arc.AttrVal(archive.AttrBitIsMovie) == archive.AttrYes {
			return // regard it as none
		}
	}
	return true
}
