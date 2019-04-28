package http

import (
	"strconv"

	"go-common/app/admin/main/workflow/model"
	"go-common/app/admin/main/workflow/model/param"
	"go-common/app/admin/main/workflow/model/search"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
	"go-common/library/net/http/blademaster/middleware/permit"
)

func challList(ctx *bm.Context) {
	params := ctx.Request.Form

	cidStr := params.Get("cid")
	gidStr := params.Get("gid")
	midStr := params.Get("mid")
	tidStr := params.Get("tid")
	roundsStr := params.Get("rounds")
	statesStr := params.Get("states")
	ctimeFrom := params.Get("ctime_from")
	ctimeTo := params.Get("ctime_to")
	order := params.Get("order")
	sort := params.Get("sort_order")
	pageStr := params.Get("pn")
	pagesizeStr := params.Get("ps")

	cc := &search.ChallSearchCommonCond{}
	numsmap := []*intsParam{
		{value: cidStr, p: &cc.IDs},
		{value: gidStr, p: &cc.Gids},
		{value: midStr, p: &cc.Mids},
		{value: tidStr, p: &cc.Tids},
		{value: statesStr, p: &cc.States},
		{value: roundsStr, p: &cc.Rounds},
	}
	var pn, ps int64
	nummap := []*intParam{
		{value: pageStr, p: &pn},
		{value: pagesizeStr, p: &ps},
	}
	if err := dealNumsmap(numsmap); err != nil {
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	if err := dealNummap(nummap); err != nil {
		ctx.JSON(nil, ecode.RequestErr)
		return
	}

	cc.Order = order
	if cc.Order == "" {
		cc.Order = "id"
	}
	if cc.Order == "last_time" {
		cc.Order = "ctime"
	}

	cc.Sort = adjustOrder("challenge", sort)
	if cc.Sort != "asc" {
		cc.Sort = "desc"
	}
	if ctimeFrom != "" {
		cc.CTimeFrom = ctimeFrom
	}
	if ctimeTo != "" {
		cc.CTimeTo = ctimeTo
	}
	cc.FormatState()
	cc.Fields = []string{"id"}
	cc.PN, _ = strconv.Atoi(pageStr)
	cc.PS, _ = strconv.Atoi(pagesizeStr)
	ctx.JSON(wkfSvc.ChallList(ctx, cc))
}

func challListCommon(ctx *bm.Context) {
	var (
		IPers interface{}
		ok    bool
		pers  []string
		err   error
	)
	v := new(param.ChallengeListCommonParam)
	if err = ctx.Bind(v); err != nil {
		return
	}

	cc := &search.ChallSearchCommonCond{}
	cc.Business = v.Business
	cc.IDs = v.IDs
	cc.Oids = v.Oids
	cc.Mids = v.Mids
	cc.ObjectMids = v.ObjectMids
	cc.Gids = v.Gids
	cc.TypeIDs = v.TypeIDs
	cc.Tids = v.Tids
	cc.Rounds = v.Rounds
	cc.AssigneeAdminIDs = v.AssigneeAdminIDs
	cc.AssigneeAdminNames = v.AssigneeAdminNames
	cc.AdminIDs = v.AdminIDs
	cc.CTimeTo = v.CTimeTo
	cc.CTimeFrom = v.CTimeFrom
	cc.Order = v.Order
	cc.Sort = v.Sort
	cc.States = v.States
	cc.BusinessStates = v.BusinessStates
	cc.PN = v.PN
	cc.PS = v.PS
	cc.CTimeFrom = v.CTimeFrom
	cc.CTimeTo = v.CTimeTo

	cc.FormatState()

	if v.Title != "" {
		cc.KWFields = append(cc.KWFields, "title")
		cc.KW = append(cc.KW, v.Title)
	}
	if v.Content != "" {
		cc.KWFields = append(cc.KWFields, "content")
		cc.KW = append(cc.KW, v.Content)
	}

	if IPers, ok = ctx.Get(permit.CtxPermissions); ok {
		pers = IPers.([]string)
	}

	if ok = isPermitChallList(pers, cc); !ok {
		ctx.JSON(nil, ecode.AccessDenied)
		ctx.Abort()
		return
	}
	cc.Fields = []string{"id", "gid"}
	ctx.JSON(wkfSvc.ChallListCommon(ctx, cc))
}

func challListV3(ctx *bm.Context) {
	v := new(param.ChallengeListV3Param)
	if err := ctx.Bind(v); err != nil {
		return
	}
	cc := &search.ChallSearchCommonCond{}
	cc.Fields = []string{"id", "oid", "mid"}
	cc.Business = v.Business
	cc.IDs = v.IDs
	cc.Oids = v.Oids
	cc.Mids = v.Mids
	cc.Gids = v.Gids
	cc.TypeIDs = v.TypeIDs
	cc.Tids = v.Tids
	cc.Rounds = v.Roles
	cc.AssigneeAdminIDs = v.AssigneeAdminIDs
	cc.AssigneeAdminNames = v.AssigneeAdminNames
	cc.AdminIDs = v.AdminIDs
	cc.CTimeTo = v.CTimeTo
	cc.CTimeFrom = v.CTimeFrom
	cc.Order = v.Order
	cc.Sort = v.Sort
	cc.States = v.States
	cc.BusinessStates = v.BusinessStates
	cc.PN = v.PN
	cc.PS = v.PS
	cc.CTimeFrom = v.CTimeFrom
	cc.CTimeTo = v.CTimeTo
	cc.KW = v.KW
	cc.KWFields = v.KWField
	cc.FormatState()

	ctx.JSON(wkfSvc.ChallListV3(ctx, cc))
}

func challDetail(ctx *bm.Context) {
	v := &struct {
		Cid int64 `form:"cid" validate:"required,gt=0"`
	}{}
	if err := ctx.Bind(v); err != nil {
		return
	}

	ctx.JSON(wkfSvc.ChallDetail(ctx, v.Cid))
}

func upChallBusState(ctx *bm.Context) {
	reqUpFields := &struct {
		Cid             int64 `form:"cid" json:"cid"`
		AssigneeAdminid int64 `json:"adminid"`
		BusState        int8  `form:"business_state" json:"business_state" validate:"min=0,max=14"`
	}{}
	if err := ctx.BindWith(reqUpFields, binding.FormPost); err != nil {
		return
	}
	adminID, adminName := adminInfo(ctx)

	ctx.JSON(nil, wkfSvc.UpChallBusState(ctx, reqUpFields.Cid, adminID, adminName, reqUpFields.BusState))
}

func batchUpChallBusState(ctx *bm.Context) {
	var (
		err         error
		reqUpFields = new(struct {
			Cids            []int64 `form:"cids,split" json:"cids" validate:"required,gt=0"`
			AssigneeAdminid int64   `json:"adminid"`
			BusState        int8    `form:"business_state" json:"business_state" validate:"min=0,max=14"`
		})
	)
	if err = ctx.BindWith(reqUpFields, binding.FormPost); err != nil {
		return
	}
	adminID, adminName := adminInfo(ctx)

	ctx.JSON(nil, wkfSvc.BatchUpChallBusState(ctx, reqUpFields.Cids, adminID, adminName, reqUpFields.BusState))
}

func upChallBusStateV3(ctx *bm.Context) {
	bcbsp := new(param.BatchChallBusStateParam)
	if err := ctx.BindWith(bcbsp, binding.FormPost); err != nil {
		return
	}
	bcbsp.AssigneeAdminID, bcbsp.AssigneeAdminName = adminInfo(ctx)
	ctx.JSON(nil, wkfSvc.SetChallBusState(ctx, bcbsp))
}

func upBusChallsBusState(ctx *bm.Context) {
	var (
		err         error
		reqUpFields struct {
			Business        int8                   `form:"business" json:"business" validate:"required,min=1"`
			Oid             int64                  `form:"oid" json:"oid" validate:"required,min=1"`
			AssigneeAdminid int64                  `form:"adminid" json:"adminid" validate:"required,min=1"`
			State           int8                   `form:"business_state" json:"business_state" validate:"min=0,max=14"`
			PreStates       []int8                 `form:"pre_business_states" json:"pre_business_states" validate:"dive,gt=-1"` //  business_state修改前的状态
			Extra           map[string]interface{} `form:"extra" json:"extra"`
		}
	)

	if err = ctx.BindWith(&reqUpFields, binding.JSON); err != nil {
		log.Error("/business/busState/update bind failed error(%v)", err)
		return
	}

	if len(reqUpFields.PreStates) <= 0 {
		reqUpFields.PreStates = append(reqUpFields.PreStates, int8(0))
	}

	upCids, err := wkfSvc.UpBusChallsBusState(ctx, reqUpFields.Business, reqUpFields.State, reqUpFields.PreStates, reqUpFields.Oid, reqUpFields.AssigneeAdminid, reqUpFields.Extra)

	log.Info("call upBusChallsBusState param(%v) upcids(%v)", reqUpFields, upCids)
	ctx.JSON(map[string]interface{}{"cids": upCids}, err)
}

func setChallResult(ctx *bm.Context) {
	crp := &param.ChallResParam{}
	if err := ctx.BindWith(crp, binding.FormPost); err != nil {
		return
	}
	crp.AdminID, crp.AdminName = adminInfo(ctx)

	ctx.JSON(nil, wkfSvc.SetChallResult(ctx, crp))
}

func batchSetChallResult(ctx *bm.Context) {
	bcrp := &param.BatchChallResParam{}
	if err := ctx.BindWith(bcrp, binding.FormPost); err != nil {
		return
	}
	bcrp.AdminID, bcrp.AdminName = adminInfo(ctx)

	ctx.JSON(nil, wkfSvc.BatchSetChallResult(ctx, bcrp))
}

func setChallStateV3(ctx *bm.Context) {
	bcrp := &param.BatchChallResParam{}
	if err := ctx.BindWith(bcrp, binding.FormPost); err != nil {
		return
	}
	bcrp.AdminID, bcrp.AdminName = adminInfo(ctx)
	ctx.JSON(nil, wkfSvc.BatchSetChallResult(ctx, bcrp))
}

func rstChallResult(ctx *bm.Context) {
	crp := new(param.ChallRstParam)
	if err := ctx.BindWith(crp, binding.FormPost); err != nil {
		return
	}
	crp.AdminID, crp.AdminName = adminInfo(ctx)
	//force to pending
	crp.State = model.Pending
	ctx.JSON(nil, wkfSvc.RstChallResult(ctx, crp))
}

func rstChallResultV3(ctx *bm.Context) {
	crp := new(param.ChallRstParam)
	if err := ctx.BindWith(crp, binding.FormPost); err != nil {
		return
	}
	crp.AdminID, crp.AdminName = adminInfo(ctx)
	// TODO(zhoujiahui): force to pending now
	crp.State = model.Pending
	ctx.JSON(nil, wkfSvc.RstChallResult(ctx, crp))
}

func upChallExtra(ctx *bm.Context) {
	cep := &param.ChallExtraParam{}
	if err := ctx.BindWith(cep, binding.JSON); err != nil {
		return
	}
	cep.AdminID, cep.AdminName = adminInfo(ctx)
	ctx.JSON(nil, wkfSvc.UpChallExtraV2(ctx, cep))
}

func upChallExtraV3(ctx *bm.Context) {
	cep3 := &param.ChallExtraParamV3{}
	if err := ctx.BindWith(cep3, binding.Form); err != nil {
		return
	}
	cep3.AdminID, cep3.AdminName = adminInfo(ctx)
	ctx.JSON(nil, wkfSvc.UpChallExtraV3(ctx, cep3))
}

func batchUpChallExtra(ctx *bm.Context) {
	bcep := new(param.BatchChallExtraParam)
	if err := ctx.BindWith(bcep, binding.JSON); err != nil {
		return
	}
	bcep.AdminID, bcep.AdminName = adminInfo(ctx)

	ctx.JSON(nil, wkfSvc.BatchUpChallExtraV2(ctx, bcep))
}

func listChallBusiness(ctx *bm.Context) {
	v := new(struct {
		Cids []int64 `form:"cids,split" validate:"required,gt=0"`
	})
	if err := ctx.Bind(v); err != nil {
		return
	}
	ctx.JSON(wkfSvc.BusinessList(ctx, v.Cids))
}

func upChall(ctx *bm.Context) {
	cup := new(param.ChallUpParam)
	if err := ctx.BindWith(cup, binding.FormPost); err != nil {
		return
	}
	cup.AdminID, cup.AdminName = adminInfo(ctx)
	ctx.JSON(nil, wkfSvc.UpChall(ctx, cup))
}

func upChallV3(ctx *bm.Context) {
	cup := new(param.ChallUpParam)
	if err := ctx.BindWith(cup, binding.FormPost); err != nil {
		return
	}
	cup.AdminID, cup.AdminName = adminInfo(ctx)
	ctx.JSON(nil, wkfSvc.UpChall(ctx, cup))
}

func platformChallCount(ctx *bm.Context) {
	var (
		assigneeAdminID int64
		ok              bool
		IUid            interface{}
		IPers           interface{}
		permissionMap   map[int8]int64
	)
	if IUid, ok = ctx.Get("uid"); ok {
		assigneeAdminID = IUid.(int64)
	}

	if IPers, ok = ctx.Get(permit.CtxPermissions); ok {
		permissionMap = parsePermission(IPers.([]string))
	}

	ctx.JSON(wkfSvc.PlatformChallCount(ctx, assigneeAdminID, permissionMap))
}

func platformChallListPending(ctx *bm.Context) {
	var (
		err             error
		pclp            *param.ChallListParam
		assigneeAdminID int64
		ok              bool
		IPers           interface{}
		IUid            interface{}
		permissionMap   map[int8]int64
	)
	pclp = new(param.ChallListParam)
	if err = ctx.Bind(pclp); err != nil {
		return
	}
	if len(pclp.Businesses) != len(pclp.AssignNum) {
		ctx.JSON("business and AssignNum length not equal", ecode.RequestErr)
		return
	}
	if pclp.PS == 0 {
		pclp.PS = 10
	}

	if IUid, ok = ctx.Get("uid"); ok {
		assigneeAdminID = IUid.(int64)
	}

	if IPers, ok = ctx.Get(permit.CtxPermissions); ok {
		permissionMap = parsePermission(IPers.([]string))
	}

	ctx.JSON(wkfSvc.PlatformChallListPending(ctx, assigneeAdminID, permissionMap, pclp))
}

func platformHandlingChalllist(ctx *bm.Context) {
	var (
		err             error
		assigneeAdminID int64
		ok              bool
		permissionMap   map[int8]int64
		IUid            interface{}
		IPers           interface{}
	)
	chdlp := new(param.ChallHandlingDoneListParam)
	if err = ctx.Bind(chdlp); err != nil {
		return
	}

	if IUid, ok = ctx.Get("uid"); ok {
		assigneeAdminID = IUid.(int64)
	}

	if IPers, ok = ctx.Get(permit.CtxPermissions); ok {
		permissionMap = parsePermission(IPers.([]string))
	}

	ctx.JSON(wkfSvc.PlatformChallListHandlingDone(ctx, chdlp, permissionMap, assigneeAdminID, model.PlatformStateHandling))
}

func platformDoneChallList(ctx *bm.Context) {
	var (
		err             error
		assigneeAdminID int64
		ok              bool
		permissionMap   map[int8]int64
		IUid            interface{}
		IPers           interface{}
	)
	chdlp := new(param.ChallHandlingDoneListParam)
	if err = ctx.Bind(chdlp); err != nil {
		return
	}
	if IUid, ok = ctx.Get("uid"); ok {
		assigneeAdminID = IUid.(int64)
	}

	if IPers, ok = ctx.Get(permit.CtxPermissions); ok {
		permissionMap = parsePermission(IPers.([]string))
	}
	ctx.JSON(wkfSvc.PlatformChallListHandlingDone(ctx, chdlp, permissionMap, assigneeAdminID, model.PlatformStateDone))
}

func platformCreatedChallList(ctx *bm.Context) {
	var (
		err     error
		cclp    *param.ChallCreatedListParam
		adminID int64
		IUid    interface{}
		ok      bool
	)

	cclp = new(param.ChallCreatedListParam)
	if err = ctx.Bind(cclp); err != nil {
		return
	}
	if cclp.PS == 0 {
		cclp.PS = 10
	}

	if IUid, ok = ctx.Get("uid"); ok {
		adminID = IUid.(int64)
	}

	cond := new(search.ChallSearchCommonCond)
	cond.Fields = []string{"id", "gid"}
	cond.Business = cclp.Businesses
	cond.AdminIDs = []int64{adminID}
	cond.Order = cclp.Order
	cond.Sort = cclp.Sort
	cond.PS = cclp.PS
	cond.PN = cclp.PN

	ctx.JSON(wkfSvc.PlatformChallListCreated(ctx, cond))

}

func platformRelease(ctx *bm.Context) {
	var (
		exist           bool
		IUid            interface{}
		IPers           interface{}
		permissionMap   map[int8]int64
		assigneeAdminID int64
	)
	if IUid, exist = ctx.Get("uid"); !exist {
		ctx.JSON(nil, ecode.UserNotExist)
		return
	}
	assigneeAdminID = IUid.(int64)

	if IPers, exist = ctx.Get(permit.CtxPermissions); !exist {
		ctx.JSON(nil, ecode.MethodNoPermission)
		return
	}
	permissionMap = parsePermission(IPers.([]string))
	ctx.JSON(nil, wkfSvc.PlatformRelease(ctx, permissionMap, assigneeAdminID))
}

func platformCheckIn(ctx *bm.Context) {
	var (
		exist           bool
		IUid            interface{}
		assigneeAdminID int64
	)
	if IUid, exist = ctx.Get("uid"); !exist {
		ctx.JSON(nil, ecode.UserNotExist)
		return
	}
	assigneeAdminID = IUid.(int64)

	ctx.JSON(nil, wkfSvc.PlatformCheckIn(ctx, assigneeAdminID))
}

func isPermitChallList(pers []string, cond *search.ChallSearchCommonCond) (ok bool) {
	if cond.Business == 0 {
		return
	}
	var (
		business int8
		round    int64
	)
	if len(cond.Rounds) != 0 {
		round = cond.Rounds[0]
	}
	business = cond.Business

	switch business {
	case 2: //稿件申诉
		switch round {
		case 0:
			for _, per := range pers {
				if per == ArchiveAppealRound1 {
					cond.Rounds = append(cond.Rounds, 1)
					ok = true
				}
				if per == ArchiveAppealRound2 {
					cond.Rounds = append(cond.Rounds, 2)
					ok = true
				}
				if per == ArchiveAppealRound3 {
					cond.Rounds = append(cond.Rounds, 3)
					ok = true
				}
			}
		case 1:
			for _, per := range pers {
				if per == ArchiveAppealRound1 {
					cond.Rounds = append(cond.Rounds, 1)
					ok = true
				}
			}
		case 2:
			for _, per := range pers {
				if per == ArchiveAppealRound2 {
					cond.Rounds = append(cond.Rounds, 2)
					ok = true
				}
			}
		case 3:
			for _, per := range pers {
				if per == ArchiveAppealRound3 {
					cond.Rounds = append(cond.Rounds, 3)
					ok = true
				}
			}
		}

	case 3: // 短点评投诉
		switch round {
		case 0:
			for _, per := range pers {
				if per == ReviewShortComplainRoun1 {
					cond.Rounds = append(cond.Rounds, 1)
					ok = true
				}
				if per == ReviewShortComplainRoun2 {
					cond.Rounds = append(cond.Rounds, 2)
					ok = true
				}
			}
		case 1:
			for _, per := range pers {
				if per == ReviewShortComplainRoun1 {
					cond.Rounds = append(cond.Rounds, 1)
					ok = true
				}
			}
		case 2:
			for _, per := range pers {
				if per == ReviewShortComplainRoun2 {
					cond.Rounds = append(cond.Rounds, 2)
					ok = true
				}
			}
		}
	case 4: // 长点评投诉
		switch round {
		case 0:
			for _, per := range pers {
				if per == ReviewLongComplainRoun1 {
					cond.Rounds = append(cond.Rounds, 1)
					ok = true
				}
				if per == ReviewLongComplainRoun2 {
					cond.Rounds = append(cond.Rounds, 2)
					ok = true
				}
			}
		case 1:
			for _, per := range pers {
				if per == ReviewLongComplainRoun1 {
					cond.Rounds = append(cond.Rounds, 1)
					ok = true
				}
			}
		case 2:
			for _, per := range pers {
				if per == ReviewLongComplainRoun2 {
					cond.Rounds = append(cond.Rounds, 2)
					ok = true
				}
			}
		}
	case 5: // 小黑屋申诉
		switch round {
		case 0:
			for _, per := range pers {
				if per == CreditAppealRoun1 {
					cond.Rounds = append(cond.Rounds, 1)
					ok = true
				}
				if per == CreditAppealRoun2 {
					cond.Rounds = append(cond.Rounds, 2)
					ok = true
				}
				if per == CreditAppealRoun3 {
					cond.Rounds = append(cond.Rounds, 3)
					ok = true
				}
				if per == CreditAppealRoun4 {
					cond.Rounds = append(cond.Rounds, 4)
					ok = true
				}
				if per == CreditAppealRoun5 {
					cond.Rounds = append(cond.Rounds, 5)
					ok = true
				}
				if per == CreditAppealRoun6 {
					cond.Rounds = append(cond.Rounds, 6)
					ok = true
				}
				if per == CreditAppealRoun7 {
					cond.Rounds = append(cond.Rounds, 7)
					ok = true
				}
				if per == CreditAppealRoun8 {
					cond.Rounds = append(cond.Rounds, 8)
					ok = true
				}
			}
		case 1:
			for _, per := range pers {
				if per == CreditAppealRoun1 {
					cond.Rounds = append(cond.Rounds, 1)
					ok = true
				}
			}
		case 2:
			for _, per := range pers {
				if per == CreditAppealRoun2 {
					cond.Rounds = append(cond.Rounds, 2)
					ok = true
				}
			}
		case 3:
			for _, per := range pers {
				if per == CreditAppealRoun3 {
					cond.Rounds = append(cond.Rounds, 2)
					ok = true
				}
			}
		case 4:
			for _, per := range pers {
				if per == CreditAppealRoun4 {
					cond.Rounds = append(cond.Rounds, 2)
					ok = true
				}
			}
		case 5:
			for _, per := range pers {
				if per == CreditAppealRoun5 {
					cond.Rounds = append(cond.Rounds, 2)
					ok = true
				}
			}
		case 6:
			for _, per := range pers {
				if per == CreditAppealRoun6 {
					cond.Rounds = append(cond.Rounds, 2)
					ok = true
				}
			}
		case 7:
			for _, per := range pers {
				if per == CreditAppealRoun7 {
					cond.Rounds = append(cond.Rounds, 2)
					ok = true
				}
			}
		case 8:
			for _, per := range pers {
				if per == CreditAppealRoun8 {
					cond.Rounds = append(cond.Rounds, 2)
					ok = true
				}
			}
		}
	case 6: // 稿件审核
		switch round {
		case 0:
			for _, per := range pers {
				if per == ArchiveAuditRound1 {
					cond.Rounds = append(cond.Rounds, 1)
					ok = true
				}
			}
		case 1:
			for _, per := range pers {
				if per == ArchiveAuditRound1 {
					cond.Rounds = append(cond.Rounds, 1)
					ok = true
				}
			}
		}
	case 9: //频道举报
		switch round {
		case 0:
			for _, per := range pers {
				if per == ChannelComplainRound1 {
					cond.Rounds = append(cond.Rounds, 1)
					ok = true
				}
				if per == ChannelComplainRound1 {
					cond.Rounds = append(cond.Rounds, 2)
					ok = true
				}
			}
		case 1:
			for _, per := range pers {
				if per == ChannelComplainRound1 {
					cond.Rounds = append(cond.Rounds, 1)
					ok = true
				}
			}
		case 2:
			for _, per := range pers {
				if per == ChannelComplainRound2 {
					cond.Rounds = append(cond.Rounds, 2)
					ok = true
				}
			}
		}
	}

	return
}
