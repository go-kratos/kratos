package client

import (
	"context"

	"go-common/app/interface/main/dm2/model"
)

const (
	_subtitleGet              = "RPC.SubtitleGet"
	_subtitleSujectSubmit     = "RPC.SubtitleSujectSubmit"
	_subtitleSubjectSubmitGet = "RPC.SubtitleSubjectSubmitGet"
)

// SubtitleGet get mask list
func (s *Service) SubtitleGet(c context.Context, arg *model.ArgSubtitleGet) (res *model.VideoSubtitles, err error) {
	err = s.client.Call(c, _subtitleGet, arg, &res)
	return
}

// SubtitleSujectSubmit .
func (s *Service) SubtitleSujectSubmit(c context.Context, arg *model.ArgSubtitleAllowSubmit) (err error) {
	err = s.client.Call(c, _subtitleSujectSubmit, arg, _noArg)
	return
}

// SubtitleSubjectSubmitGet .
func (s *Service) SubtitleSubjectSubmitGet(c context.Context, arg *model.ArgArchiveID) (res *model.SubtitleSubjectReply, err error) {
	err = s.client.Call(c, _subtitleSubjectSubmitGet, arg, &res)
	return
}
