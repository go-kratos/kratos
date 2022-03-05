package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	_mixin "entgo.io/ent/schema/mixin"
)

type CreatorId struct {
	_mixin.Schema
}

func (CreatorId) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("creator_id").
			Comment("创建者用户ID"),
	}
}
