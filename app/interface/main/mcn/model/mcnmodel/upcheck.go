package mcnmodel

import "go-common/app/interface/main/mcn/model"

//isUpBindable
// up主为未绑定，已驳回，已解除，已申请的up主可被申请，up主拒绝后可重复申请，
// 审核中，已签约，已冻结不可被申请
func isUpBindable(state model.MCNUPState) bool {
	switch state {
	case model.MCNUPStateOnBlock, model.MCNUPStateOnCooling, model.MCNUPStateOnSign, model.MCNUPStateOnReview, model.MCNUPStateOnPreOpen:
		return false
	}
	return true
}
