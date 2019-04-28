package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go-common/app/admin/main/videoup/model/archive"
	"strings"
)

func (s *Service) addVideoOper(c context.Context, oper *archive.VideoOper) (err error) {
	if oldOper, _ := s.arc.VideoOper(c, oper.VID); oldOper != nil && oldOper.LastID == 1 {
		oper.LastID = oldOper.ID
		s.arc.AddVideoOper(c, oper.AID, oper.UID, oper.VID, oper.Attribute, oper.Status, oper.LastID, oper.Content, oper.Remark)
		return
	}
	if lastID, _ := s.arc.AddVideoOper(c, oper.AID, oper.UID, oper.VID, oper.Attribute, oper.Status, oper.LastID, oper.Content, oper.Remark); lastID > 0 {
		s.arc.UpVideoOper(c, lastID, lastID)
		return
	}
	return
}

func (s *Service) addArchiveOper(c context.Context, oper *archive.ArcOper) (err error) {
	if oldOper, _ := s.arc.ArchiveOper(c, oper.Aid); oldOper != nil && oldOper.LastID == 1 {
		oper.LastID = oldOper.ID
	}
	s.arc.AddArcOper(c, oper.Aid, oper.UID, oper.Attribute, oper.TypeID, oper.State, oper.Round, oper.LastID, oper.Content, oper.Remark)
	return
}

func (s *Service) diffVideoOper(vp *archive.VideoParam) (conts []string) {
	if vp.TagID > 0 {
		var operType int8
		if vp.Status >= archive.VideoStatusOpen {
			operType = archive.OperTypeOpenTag
		} else {
			operType = archive.OperTypeRecicleTag
		}
		conts = append(conts, archive.Operformat(operType, "tagid", vp.TagID, archive.OperStyleTwo))
	}
	if vp.Reason != "" {
		conts = append(conts, archive.Operformat(archive.OperTypeAduitReason, "reason", vp.Reason, archive.OperStyleTwo))
	}
	if vp.TaskID > 0 {
		conts = append(conts, archive.Operformat(archive.OperTypeTaskID, "task", vp.TaskID, archive.OperStyleTwo))
	}
	return
}

func (s *Service) diffArchiveOper(ap *archive.ArcParam, a *archive.Archive, addit *archive.Addit, forbid *archive.ForbidAttr) (conts []string, changeTypeID, changeCopyright, changeTitle, changeCover bool) {
	if ap.CanCelMission {
		conts = append(conts, archive.Operformat(archive.OperTypeMission, addit.MissionID, 0, archive.OperStyleOne))
	}
	if ap.Cover != a.Cover {
		if strings.HasPrefix(a.Cover, "http://") && strings.Contains(a.Cover, ap.Cover) {
			changeCover = false
		} else {
			changeCover = true
		}
	}
	if ap.Title != a.Title {
		changeTitle = true
	}
	if ap.Copyright != a.Copyright {
		changeCopyright = true
		conts = append(conts, archive.Operformat(archive.OperTypeCopyright, archive.CopyrightsDesc(a.Copyright), archive.CopyrightsDesc(ap.Copyright), archive.OperStyleOne))
	}
	if cont, _ := s.diffTypeID(ap.TypeID, a.TypeID, ap.State); cont != "" {
		changeTypeID = true
		conts = append(conts, cont)
	}
	if ap.RejectReason != "" && ap.RejectReason != a.RejectReason {
		if a.RejectReason == "" {
			a.RejectReason = "无"
		}
		var operType int8
		if a.Round > 20 {
			operType = archive.OperTypeRejectReason
		} else {
			operType = archive.OperTypeAduitReason
		}
		conts = append(conts, archive.Operformat(operType, a.RejectReason, ap.RejectReason, archive.OperStyleOne))
	}
	if ap.Forward != a.Forward {
		conts = append(conts, archive.Operformat(archive.OperTypeForwardID, a.Forward, ap.Forward, archive.OperStyleOne))
	}
	if ap.Notify {
		conts = append(conts, archive.Operformat(archive.OperNotify, "无", "发送通知", archive.OperStyleOne))
	}
	if forbid == nil || (forbid.OnFlowID != ap.OnFlowID) {
		if forbid != nil {
			conts = append(conts, archive.Operformat(archive.OperTypeFlowID, s.flowsCache[forbid.OnFlowID], s.flowsCache[ap.OnFlowID], archive.OperStyleOne))
		} else {
			conts = append(conts, archive.Operformat(archive.OperTypeFlowID, "无", s.flowsCache[ap.OnFlowID], archive.OperStyleOne))
		}
	}
	if ap.PTime != a.PTime {
		conts = append(conts, archive.Operformat(archive.OperTypePtime, time.Unix(int64(a.PTime), 0).Format("2006-01-02 15:04:05"), time.Unix(int64(ap.PTime), 0).Format("2006-01-02 15:04:05"), archive.OperStyleOne))
	}
	if ap.Access != a.Access {
		conts = append(conts, archive.Operformat(archive.OperTypeAccess, archive.AccessDesc(a.Access), archive.AccessDesc(ap.Access), archive.OperStyleOne))
	}
	if ap.Dynamic != addit.Dynamic {
		conts = append(conts, archive.Operformat(archive.OperTypeDynamic, addit.Dynamic, ap.Dynamic, archive.OperStyleOne))
	}
	return
}

func (s *Service) diffBatchArchiveOper(ap *archive.ArcParam, a *archive.Archive) (conts []string) {
	if ap.Access != a.Access {
		conts = append(conts, archive.Operformat(archive.OperTypeAccess, archive.AccessDesc(a.Access), archive.AccessDesc(ap.Access), archive.OperStyleOne))
	}
	if ap.PTime != a.PTime {
		conts = append(conts, archive.Operformat(archive.OperTypePtime, time.Unix(int64(a.PTime), 0).Format("2006-01-02 15:04:05"), time.Unix(int64(ap.PTime), 0).Format("2006-01-02 15:04:05"), archive.OperStyleOne))
	}
	if ap.FlagCopyright && ap.Copyright != a.Copyright {
		conts = append(conts, archive.Operformat(archive.OperTypeCopyright, archive.CopyrightsDesc(a.Copyright), archive.CopyrightsDesc(ap.Copyright), archive.OperStyleOne))
	}
	if ap.RejectReason != "" && ap.RejectReason != a.RejectReason {
		if a.RejectReason == "" {
			a.RejectReason = "无"
		}
		var operType int8
		if a.Round > 20 {
			operType = archive.OperTypeRejectReason
		} else {
			operType = archive.OperTypeAduitReason
		}
		conts = append(conts, archive.Operformat(operType, a.RejectReason, ap.RejectReason, archive.OperStyleOne))
	}
	return
}

func (s *Service) diffTypeID(newTypeID, oldTypeID int16, state int8) (cont string, changeTypeID bool) {
	if newTypeID != oldTypeID {
		changeTypeID = true
		var oldCont, newCont string
		if ok := s.isTypeID(oldTypeID); ok {
			oldCont = s.typeCache[oldTypeID].Name
		} else {
			oldCont = strconv.Itoa(int(oldTypeID))
		}
		if ok := s.isTypeID(newTypeID); ok {
			newCont = s.typeCache[newTypeID].Name
		} else {
			newCont = strconv.Itoa(int(newTypeID))
		}
		cont = archive.Operformat(archive.OperTypeTypeID, oldCont, newCont, archive.OperStyleOne)
		if state < 0 {
			cont = fmt.Sprintf("%s,过审后生效", cont)
		}
	}
	return
}

//私单修改日志
func (s *Service) diffPorder(c context.Context, aid int64, ap *archive.ArcParam) (conts []string, porder *archive.Porder) {
	porder, _ = s.arc.Porder(c, aid)
	if porder.IndustryID > 0 {
		var yesOrNo = map[int8]string{int8(1): "是", int8(0): "否"}
		if porder.IndustryID != ap.IndustryID {
			conts = append(conts, archive.Operformat(archive.OperPorderIndustryID, s.porderConfigCache[porder.IndustryID].Name, s.porderConfigCache[ap.IndustryID].Name, archive.OperStyleOne))
		}
		if porder.Official != ap.Official {
			conts = append(conts, archive.Operformat(archive.OperPorderOfficial, yesOrNo[porder.Official], yesOrNo[ap.Official], archive.OperStyleOne))
		}
		//game ID
		if porder.BrandID != ap.BrandID {
			conts = append(conts, archive.Operformat(archive.OperPorderBrandID, porder.BrandID, ap.BrandID, archive.OperStyleOne))
		}
		//custom brandName
		if porder.BrandName != ap.BrandName {
			conts = append(conts, archive.Operformat(archive.OperPorderBrandName, porder.BrandName, ap.BrandName, archive.OperStyleOne))
		}
		//porderConfigCache
		if porder.ShowType != ap.ShowType {
			conts = append(conts, archive.Operformat(archive.OperPorderShowType, porder.ShowType, ap.ShowType, archive.OperStyleOne))
		}
		if porder.Advertiser != ap.Advertiser {
			conts = append(conts, archive.Operformat(archive.OperPorderAdvertiser, porder.Advertiser, ap.Advertiser, archive.OperStyleOne))
		}
		if porder.Agent != ap.Agent {
			conts = append(conts, archive.Operformat(archive.OperPorderAgent, porder.Agent, ap.Agent, archive.OperStyleOne))
		}
	}
	return
}
