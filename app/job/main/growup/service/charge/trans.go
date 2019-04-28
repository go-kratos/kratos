package charge

import (
	model "go-common/app/job/main/growup/model/charge"
)

func transAv2Archive(avs []*model.AvCharge) (archs []*model.Archive) {
	archs = make([]*model.Archive, 0, len(avs))
	for _, av := range avs {
		archs = append(archs, &model.Archive{
			ID:        av.AvID,
			IncCharge: av.IncCharge,
			TagID:     av.TagID,
			Date:      av.Date,
		})
	}
	return
}

func transAvMap2Archive(avs map[int64]*model.AvCharge) (archs []*model.Archive) {
	archs = make([]*model.Archive, 0, len(avs))
	for _, av := range avs {
		archs = append(archs, &model.Archive{
			ID:        av.AvID,
			IncCharge: av.IncCharge,
			TagID:     av.TagID,
			Date:      av.Date,
		})
	}
	return
}

func transCm2Archive(cms []*model.Column) (archs []*model.Archive) {
	archs = make([]*model.Archive, 0, len(cms))
	for _, cm := range cms {
		archs = append(archs, &model.Archive{
			ID:        cm.AID,
			IncCharge: cm.IncCharge,
			TagID:     cm.TagID,
			Date:      cm.Date,
		})
	}
	return
}

func transCmMap2Archive(cms map[int64]*model.Column) (archs []*model.Archive) {
	archs = make([]*model.Archive, 0, len(cms))
	for _, cm := range cms {
		archs = append(archs, &model.Archive{
			ID:        cm.AID,
			IncCharge: cm.IncCharge,
			TagID:     cm.TagID,
			Date:      cm.Date,
		})
	}
	return
}

func transBgm2Archive(bgms []*model.BgmCharge) (archs []*model.Archive) {
	archs = make([]*model.Archive, 0, len(bgms))
	for _, bgm := range bgms {
		archs = append(archs, &model.Archive{
			ID:        bgm.SID,
			IncCharge: bgm.IncCharge,
			TagID:     0,
			Date:      bgm.Date,
		})
	}
	return
}

func transBgmMap2Archive(bgms map[string]*model.BgmCharge) (archs []*model.Archive) {
	archs = make([]*model.Archive, 0, len(bgms))
	for _, bgm := range bgms {
		archs = append(archs, &model.Archive{
			ID:        bgm.SID,
			IncCharge: bgm.IncCharge,
			TagID:     0,
			Date:      bgm.Date,
		})
	}
	return
}
