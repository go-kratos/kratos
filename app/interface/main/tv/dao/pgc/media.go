package pgc

import (
	"context"

	"go-common/app/interface/main/tv/model"
	"go-common/library/log"
)

// Media gets the media detail data from PGC API
func (d *Dao) Media(ctx context.Context, tvParam *model.MediaParam) (detail *model.SeasonDetail, err error) {
	var result model.MediaResp
	if err = d.client.Get(ctx, d.conf.Host.APIMedia, "", tvParam.GenerateUrl(), &result); err != nil {
		log.Error("ClientGet Sid %d, error[%v]", tvParam.SeasonID, err)
		return
	}
	if err = result.CodeErr(); err != nil {
		log.Error("PGC API MediaResp: [CODE:(%d),MESSAGE:(%s)]", result.Code, result.Message)
		return
	}
	detail = result.Result
	return
}

// MediaV2 gets the media detail data from PGC API V2
func (d *Dao) MediaV2(ctx context.Context, tvParam *model.MediaParam) (detail *model.SnDetailV2, err error) {
	var result model.MediaRespV2
	if err = d.client.Get(ctx, d.conf.Host.APIMediaV2, "", tvParam.GenerateUrl(), &result); err != nil {
		log.Error("ClientGet Sid %d, error[%v]", tvParam.SeasonID, err)
		return
	}
	if err = result.CodeErr(); err != nil {
		log.Error("PGC API MediaResp: [CODE:(%d),MESSAGE:(%s)]", result.Code, result.Message)
		return
	}
	detail = result.Result
	return
}
