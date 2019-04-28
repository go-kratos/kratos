package service

import (
	"context"
	"fmt"

	"go-common/app/interface/main/videoup/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
)

func (s *Service) dealElec(c context.Context, openElec int8, aid, mid int64, ip string) (err error) {
	show, err := s.elec.ArcShow(c, mid, aid, ip)
	if err != nil {
		log.Error("s.elec.ArcShow(%d, %d, %d, %d) error(%v)", mid, aid, openElec, show, err)
		return
	}
	if show != (openElec == 1) {
		s.elec.ArcUpdate(c, mid, aid, openElec, ip)
	}
	return
}

func (s *Service) dealOrder(c context.Context, mid, aid, orderID int64, ip string) (err error) {
	if orderID == 0 {
		return
	}
	if err = s.order.BindOrder(c, mid, aid, orderID, ip); err != nil {
		log.Error("s.order.ExecuteOrder mid(%d) aid(%d) orderId(%d)  error(%v)", mid, aid, orderID, err)
		err = ecode.VideoupOrderAPIErr
	}
	return
}

func (s *Service) dealTag(c context.Context, mid, aid int64, srcTag, descTag, ip string, typeID int16) (err error) {
	if srcTag != descTag {
		typeName := ""
		if tp, ok := s.typeCache[typeID]; ok && tp != nil {
			typeName = tp.Name
			if tp, ok = s.typeCache[tp.PID]; ok && tp != nil {
				typeName = fmt.Sprintf("%s,%s", typeName, tp.Name)
			}
		}
		if err = s.tag.UpBind(c, mid, aid, descTag, typeName, ip); err != nil {
			log.Error("s.tag.UpBind(%d, %d, %s, %s,%s) error(%d)", mid, aid, srcTag, descTag, typeName, err)
			return
		}
	}
	return
}

func (s *Service) dealWaterMark(c context.Context, mid int64, wm *archive.Watermark, ip string) (err error) {
	if wm != nil {
		if err = s.creative.SetWatermark(c, mid, wm.State, wm.Ty, wm.Pos, ip); err != nil {
			log.Error("s.creative.SetWatermark(%d,%+v,%+v) error(%d)", mid, wm, err)
			return
		}
	}
	return
}

func (s *Service) freshFavs(c context.Context, mid int64, ap *archive.ArcParam, ip string) (err error) {
	if err = s.arc.FreshFavTypes(c, mid, int(ap.TypeID)); err != nil {
		log.Error("s.arc.FreshFavTypes(%d,%+v,%+v) error(%d)", mid, ap, err)
		return
	}
	return
}

func (s *Service) uploadVideoEditInfo(c context.Context, ap *archive.ArcParam, aid, mid int64, ip string) (err error) {
	ap.EmptyVideoEditInfo()
	editors := make([]*archive.Editor, 0)
	for _, v := range ap.Videos {
		if v.Editor != nil && v.Cid > 0 {
			v.Editor.UpFrom = ap.UpFrom
			v.Editor.CID = v.Cid
			editors = append(editors, v.Editor)
		}
	}
	if len(editors) > 0 {
		if err = s.creative.UploadMaterial(c, editors, aid, mid, ip); err != nil {
			log.Error("s.creative.UploadMaterial (%+v,%d,%d,%s) error(%+v)", editors, aid, mid, ip, err)
			return
		}
	}
	return
}

func (s *Service) lotteryBind(c context.Context, lotteryID, aid, mid int64, ip string) (err error) {
	ck, _ := s.dynamic.UserCheck(c, mid, ip)
	if lotteryID > 0 && (ck == 1) {
		if err = s.dynamic.LotteryBind(c, lotteryID, aid, mid, ip); err != nil {
			log.Error("s.dynamic.LotteryBind (%+v,%d,%d,%s) error(%d)", lotteryID, aid, mid, ip, err)
			return
		}
	}
	return
}

func (s *Service) addFollowing(c context.Context, mid int64, fids []int64, upfrom int8, ip string) (err error) {
	if len(fids) > 0 {
		var src int
		if upfrom == archive.UpFromAPPAndroid {
			src = 173
		} else if upfrom == archive.UpFromAPPiOS || upfrom == archive.UpFromIpad {
			src = 183
		} else {
			src = 173
		}
		for _, fid := range fids {
			s.acc.AddFollowing(context.Background(), mid, fid, src, ip)
		}
	}
	return
}
