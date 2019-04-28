package server

import (
	v1 "go-common/app/service/main/account/api"
	"go-common/app/service/main/account/model"
	"go-common/library/net/rpc/context"
)

// Info3 receive ArgMid contains mid and real ip, then init user info.
func (r *RPC) Info3(c context.Context, arg *model.ArgMid, res *v1.Info) (err error) {
	var info *v1.Info
	if info, err = r.s.Info(c, arg.Mid); err == nil && info != nil {
		*res = *info
	}
	return
}

// Infos3 receive ArgMids contains mids and real ip, then multi init user info.
func (r *RPC) Infos3(c context.Context, a *model.ArgMids, res *map[int64]*v1.Info) (err error) {
	*res, err = r.s.Infos(c, a.Mids)
	return
}

// InfosByName3 receive ArgMids contains mids and real ip, then multi init user info.
func (r *RPC) InfosByName3(c context.Context, a *model.ArgNames, res *map[int64]*v1.Info) (err error) {
	*res, err = r.s.InfosByName(c, a.Names)
	return
}

// Card3 receive ArgMid contains mid and real ip, then init user card.
func (r *RPC) Card3(c context.Context, arg *model.ArgMid, res *v1.Card) (err error) {
	var card *v1.Card
	if card, err = r.s.Card(c, arg.Mid); err == nil && res != nil {
		*res = *card
	}
	return
}

// Cards3 receive ArgMids contains mids and real ip, then multi init user card.
func (r *RPC) Cards3(c context.Context, a *model.ArgMids, res *map[int64]*v1.Card) (err error) {
	*res, err = r.s.Cards(c, a.Mids)
	return
}

// Profile3 get user audit info.
func (r *RPC) Profile3(c context.Context, arg *model.ArgMid, res *v1.Profile) (err error) {
	var p *v1.Profile
	if p, err = r.s.Profile(c, arg.Mid); err == nil && p != nil {
		*res = *p
	}
	return
}

// ProfileWithStat3 get user audit info.
func (r *RPC) ProfileWithStat3(c context.Context, arg *model.ArgMid, res *model.ProfileStat) (err error) {
	var p *model.ProfileStat
	if p, err = r.s.ProfileWithStat(c, arg.Mid); err == nil && p != nil {
		*res = *p
	}
	return
}

// AddExp3 add exp for user.
func (r *RPC) AddExp3(c context.Context, a *model.ArgExp, res *struct{}) (err error) {
	err = r.s.AddExp(c, a.Mid, a.Exp, a.Operater, a.Operate, a.Reason)
	return
}

// AddMoral3 receive ArgMoral contains mid, moral, oper, reason and remark, then add moral for user.
func (r *RPC) AddMoral3(c context.Context, a *model.ArgMoral, res *struct{}) (err error) {
	err = r.s.AddMoral(c, a.Mid, a.Moral, a.Oper, a.Reason, a.Remark)
	return
}

// Relation3 get friend relation.
func (r *RPC) Relation3(c context.Context, a *model.ArgRelation, res *model.Relation) (err error) {
	var rl *model.Relation
	if rl, err = r.s.Relation(c, a.Mid, a.Owner); err == nil && rl != nil {
		*res = *rl
	}
	return
}

// Attentions3 get attentions list ,including following and whisper.
func (r *RPC) Attentions3(c context.Context, a *model.ArgMid, res *[]int64) (err error) {
	*res, err = r.s.Attentions(c, a.Mid)
	return
}

// Blacks3 get user black list.
func (r *RPC) Blacks3(c context.Context, a *model.ArgMid, res *map[int64]struct{}) (err error) {
	*res, err = r.s.Blacks(c, a.Mid)
	return
}

// Relations3 get friend relations.
func (r *RPC) Relations3(c context.Context, a *model.ArgRelations, res *map[int64]*model.Relation) (err error) {
	*res, err = r.s.Relations(c, a.Mid, a.Owners)
	return
}

// RichRelations3 get friend relations.
func (r *RPC) RichRelations3(c context.Context, a *model.ArgRichRelation, res *map[int64]int) (err error) {
	*res, err = r.s.RichRelations2(c, a.Owner, a.Mids)
	return
}
