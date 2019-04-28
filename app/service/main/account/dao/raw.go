package dao

import (
	"context"

	v1 "go-common/app/service/main/account/api"
	"go-common/app/service/main/account/model"
	mml "go-common/app/service/main/member/model"
	bml "go-common/app/service/main/member/model/block"
	sml "go-common/app/service/main/usersuit/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup"

	"github.com/pkg/errors"
)

// RawInfo raw info.
func (d *Dao) RawInfo(c context.Context, mid int64) (res *v1.Info, err error) {
	var base *mml.BaseInfo
	if base, err = d.mRPC.Base(c, &mml.ArgMemberMid{Mid: mid}); err != nil {
		// err = errors.Wrap(err, "dao raw info")
		// err = errors.WithStack(err)
		return
	}
	res = &v1.Info{
		Mid:  base.Mid,
		Name: base.Name,
		Sex:  base.SexStr(),
		Face: base.Face,
		Sign: base.Sign,
		Rank: int32(base.Rank),
	}
	return
}

// RawInfos raw infos.
func (d *Dao) RawInfos(c context.Context, mids []int64) (res map[int64]*v1.Info, err error) {
	var bm map[int64]*mml.BaseInfo
	if bm, err = d.mRPC.Bases(c, &mml.ArgMemberMids{Mids: mids}); err != nil {
		// err = errors.Wrap(err, "dao raw info")
		err = errors.WithStack(err)
		return
	}
	res = map[int64]*v1.Info{}
	for _, base := range bm {
		i := &v1.Info{
			Mid:  base.Mid,
			Name: base.Name,
			Sex:  base.SexStr(),
			Face: base.Face,
			Sign: base.Sign,
			Rank: int32(base.Rank),
		}
		res[i.Mid] = i
	}
	return
}

// RawCard get card by mid.
func (d *Dao) RawCard(c context.Context, mid int64) (res *v1.Card, err error) {
	eg, _ := errgroup.WithContext(c)
	var mb *mml.Member
	eg.Go(func() (e error) {
		if mb, e = d.mRPC.Member(c, &mml.ArgMemberMid{Mid: mid}); e != nil {
			log.Error("d.mRPC.Member(%d) err(%v)", mid, e)
			// e = ecode.Degrade
		}
		return
	})
	var medal *sml.MedalInfo
	eg.Go(func() (e error) {
		if medal, e = d.suitRPC.MedalActivated(c, &sml.ArgMid{Mid: mid}); e != nil {
			log.Error("s.suitRPC.MedalActivated(%d) err %v", mid, e)
			e = ecode.Degrade
		}
		return
	})
	var block *bml.RPCResInfo
	eg.Go(func() (e error) {
		if block, e = d.mRPC.BlockInfo(c, &bml.RPCArgInfo{MID: mid}); e != nil {
			log.Error("d.block.BlockInfo(%d) err %v", mid, e)
			e = ecode.Degrade
		}
		return
	})
	var vip *v1.VipInfo
	eg.Go(func() (e error) {
		if vip, e = d.Vip(c, mid); e != nil {
			log.Error("d.Vip(%d) err(%v)", mid, e)
			e = ecode.Degrade
		}
		return
	})
	var pendant *sml.PendantEquip
	eg.Go(func() (e error) {
		if pendant, e = d.suitRPC.Equipment(c, &sml.ArgEquipment{Mid: mid}); e != nil {
			log.Error("d.suitRPC.Equipment(%d) err(%v)", mid, e)
			e = ecode.Degrade
		}
		return
	})
	if err = eg.Wait(); err != nil && err != ecode.Degrade {
		return
	}
	card := &v1.Card{Mid: mid}
	if medal != nil {
		card.Nameplate.Nid = int(medal.ID)
		card.Nameplate.Name = medal.Name
		card.Nameplate.Image = medal.Image
		card.Nameplate.ImageSmall = medal.ImageSmall
		card.Nameplate.Level = medal.LevelDesc
		card.Nameplate.Condition = medal.Condition
	}
	if block != nil {
		card.Silence = blockStatusToSilence(block.BlockStatus)
	}
	if mb != nil {
		card.Name = mb.Name
		card.Sign = mb.Sign
		card.Sex = mb.SexStr()
		card.Rank = int32(mb.Rank)
		card.Face = mb.Face
		if mb.OfficialInfo != nil {
			// card.Official = *mb.OfficialInfo
			card.Official.DeepCopyFromOfficialInfo(mb.OfficialInfo)
		}
		card.Level = mb.Cur
	}
	if vip != nil {
		card.Vip = *vip
	}
	if pendant != nil {
		card.Pendant.Pid = int(pendant.Pid)
		card.Pendant.Expire = int(pendant.Expires)
		if pendant.Pendant != nil {
			card.Pendant.Name = pendant.Pendant.Name
			card.Pendant.Image = fullImage(mid, pendant.Pendant.Image)
		}
	}
	res = card
	return
}

// RawCards get card by mid.
func (d *Dao) RawCards(c context.Context, mids []int64) (res map[int64]*v1.Card, err error) {
	eg, _ := errgroup.WithContext(c)
	var medals map[int64]*sml.MedalInfo
	eg.Go(func() (e error) {
		if medals, e = d.suitRPC.MedalActivatedMulti(c, &sml.ArgMids{Mids: mids}); e != nil {
			log.Error("s.suitRPC.MedalActivatedMulti(%v) err %v", mids, e)
			e = ecode.Degrade
		}
		return
	})
	var blocks map[int64]*bml.RPCResInfo
	eg.Go(func() (e error) {
		var bs []*bml.RPCResInfo
		if bs, e = d.mRPC.BlockBatchInfo(c, &bml.RPCArgBatchInfo{MIDs: mids}); e != nil {
			log.Error("d.block.BlockBatchInfo(%v) err %v", mids, e)
			e = ecode.Degrade
		}
		blocks = make(map[int64]*bml.RPCResInfo, len(bs))
		for _, block := range bs {
			blocks[block.MID] = block
		}
		return
	})
	var mbs map[int64]*mml.Member
	eg.Go(func() (e error) {
		if mbs, e = d.mRPC.Members(c, &mml.ArgMemberMids{Mids: mids}); e != nil {
			log.Error("d.mRPC.Members(%v) err(%v)", mids, e)
			// e = ecode.Degrade
		}
		return
	})
	var vipm map[int64]*v1.VipInfo
	eg.Go(func() (e error) {
		if vipm, e = d.Vips(c, mids); e != nil {
			log.Error("d.CpVips(%v) err(%v)", mids, e)
			e = ecode.Degrade
		}
		return
	})
	var pendantm map[int64]*sml.PendantEquip
	eg.Go(func() (e error) {
		if pendantm, e = d.suitRPC.Equipments(c, &sml.ArgEquipments{Mids: mids}); e != nil {
			log.Error("d.suitRPC.Equipments(%v) err(%v)", mids, e)
			e = ecode.Degrade
		}
		return
	})
	if err = eg.Wait(); err != nil && err != ecode.Degrade {
		return
	}
	res = map[int64]*v1.Card{}
	for _, mid := range mids {
		card := &v1.Card{Mid: mid}
		if mb, ok := mbs[mid]; ok && mb != nil {
			card.Name = mb.Name
			card.Sign = mb.Sign
			card.Sex = mb.SexStr()
			card.Rank = int32(mb.Rank)
			card.Face = mb.Face
			if mb.OfficialInfo != nil {
				// card.Official = *mb.OfficialInfo
				card.Official.DeepCopyFromOfficialInfo(mb.OfficialInfo)
			}
			card.Level = mb.Cur
		} else {
			continue
		}
		if block, ok := blocks[mid]; ok && block != nil {
			card.Silence = blockStatusToSilence(block.BlockStatus)
		}
		if medal, ok := medals[mid]; ok && medal != nil {
			card.Nameplate.Nid = int(medal.ID)
			card.Nameplate.Name = medal.Name
			card.Nameplate.Image = medal.Image
			card.Nameplate.ImageSmall = medal.ImageSmall
			card.Nameplate.Level = medal.LevelDesc
			card.Nameplate.Condition = medal.Condition
		}
		if vip, ok := vipm[mid]; ok && vip != nil {
			card.Vip = *vip
		}
		if pendant, ok := pendantm[mid]; ok && pendant != nil {
			card.Pendant.Pid = int(pendant.Pid)
			card.Pendant.Expire = int(pendant.Expires)
			if pendant.Pendant != nil {
				card.Pendant.Name = pendant.Pendant.Name
				card.Pendant.Image = fullImage(mid, pendant.Pendant.Image)
			}
		}
		res[mid] = card
	}
	return
}

// RawProfile get profile by mid.
func (d *Dao) RawProfile(c context.Context, mid int64) (res *v1.Profile, err error) {
	eg, _ := errgroup.WithContext(c)
	var detail *model.PassportDetail
	eg.Go(func() (e error) {
		if detail, e = d.PassportDetail(c, mid); e != nil {
			log.Error("d.PassPortDetail(%d) err %v", mid, e)
			// e = ecode.Degrade
		}
		return
	})
	var mb *mml.Member
	eg.Go(func() (e error) {
		if mb, e = d.mRPC.Member(c, &mml.ArgMemberMid{Mid: mid}); e != nil {
			log.Error("d.mRPC.Member(%d) err(%v)", mid, e)
			// e = ecode.Degrade
		}
		return
	})
	var moral *mml.Moral
	eg.Go(func() (e error) {
		if moral, e = d.mRPC.Moral(c, &mml.ArgMemberMid{Mid: mid}); e != nil {
			log.Error("d.mRPC.Member(%d) err(%v)", mid, e)
			e = ecode.Degrade
		}
		return
	})
	var realNameStatus *mml.RealnameStatus
	eg.Go(func() (e error) {
		if realNameStatus, e = d.mRPC.RealnameStatus(c, &mml.ArgMemberMid{Mid: mid}); e != nil {
			log.Error("d.mRPC.RealnameStatus(%d) err(%v)", mid, e)
			e = ecode.Degrade
		}
		return
	})
	var medal *sml.MedalInfo
	eg.Go(func() (e error) {
		if medal, e = d.suitRPC.MedalActivated(c, &sml.ArgMid{Mid: mid}); e != nil {
			log.Error("s.suitRPC.MedalActivated(%d) err %v", mid, e)
			e = ecode.Degrade
		}
		return
	})
	var block *bml.RPCResInfo
	eg.Go(func() (e error) {
		if block, e = d.mRPC.BlockInfo(c, &bml.RPCArgInfo{MID: mid}); e != nil {
			log.Error("s.dao.BlockInfo(%d) err %v", mid, e)
			e = ecode.Degrade
		}
		return
	})
	var vip *v1.VipInfo
	eg.Go(func() (e error) {
		if vip, e = d.Vip(c, mid); e != nil {
			log.Error("d.Vip(%d) err(%v)", mid, e)
			e = ecode.Degrade
		}
		return
	})
	var pendant *sml.PendantEquip
	eg.Go(func() (e error) {
		if pendant, e = d.suitRPC.Equipment(c, &sml.ArgEquipment{Mid: mid}); e != nil {
			log.Error("d.suitRPC.Equipment(%d) err(%v)", mid, e)
			e = ecode.Degrade
		}
		return
	})
	if err = eg.Wait(); err != nil && err != ecode.Degrade {
		return
	}
	pfl := &v1.Profile{Mid: mid}
	if mb != nil {
		pfl.Name = mb.Name
		pfl.Sign = mb.Sign
		pfl.Sex = mb.SexStr()
		pfl.Rank = int32(mb.Rank)
		pfl.Face = mb.Face
		if mb.OfficialInfo != nil {
			// pfl.Official = *mb.OfficialInfo
			pfl.Official.DeepCopyFromOfficialInfo(mb.OfficialInfo)
		}
		pfl.Level = mb.Cur
		pfl.Birthday = mb.Birthday
	}
	if block != nil {
		pfl.Silence = blockStatusToSilence(block.BlockStatus)
	}
	if detail != nil {
		pfl.JoinTime = int32(detail.JoinTime)
		pfl.EmailStatus = bindEmailStatus(detail.Email, detail.Spacesta)
		pfl.TelStatus = bindPhoneStatus(detail.Phone)
		pfl.IsTourist = boolToInt32(detail.IsTourist)
	}
	if realNameStatus != nil {
		pfl.Identification = identificationStatus(*realNameStatus)
	}
	pfl.Moral = parseMoral(moral)
	if medal != nil {
		pfl.Nameplate.Nid = int(medal.ID)
		pfl.Nameplate.Name = medal.Name
		pfl.Nameplate.Image = medal.Image
		pfl.Nameplate.ImageSmall = medal.ImageSmall
		pfl.Nameplate.Level = medal.LevelDesc
		pfl.Nameplate.Condition = medal.Condition
	}
	if vip != nil {
		pfl.Vip = *vip
	}
	if pendant != nil {
		pfl.Pendant.Pid = int(pendant.Pid)
		pfl.Pendant.Expire = int(pendant.Expires)
		if pendant.Pendant != nil {
			pfl.Pendant.Name = pendant.Pendant.Name
			pfl.Pendant.Image = fullImage(mid, pendant.Pendant.Image)
		}
	}
	res = pfl
	return
}

func blockStatusToSilence(status bml.BlockStatus) int32 {
	return boolToInt32(status == bml.BlockStatusForever || status == bml.BlockStatusLimit)
}

func bindEmailStatus(email string, spacesta int8) int32 {
	return boolToInt32(spacesta > -10 && len(email) > 0)
}

func bindPhoneStatus(phone string) int32 {
	return boolToInt32(len(phone) > 0)
}

func parseMoral(moral *mml.Moral) int32 {
	m := int32(mml.DefaultMoral)
	if moral != nil {
		m = int32(moral.Moral)
	}
	return m / 100
}

func identificationStatus(realNameStatus mml.RealnameStatus) int32 {
	return boolToInt32(realNameStatus == mml.RealnameStatusTrue)
}

func boolToInt32(b bool) int32 {
	if b {
		return 1
	}
	return 0
}
