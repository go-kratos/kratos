package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

type SwitchStatus struct {
	mixin.Schema
}

func (SwitchStatus) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("status").
			Comment("状态").
			Values(
				"on",
				"off",
			).
			Default("on"),
	}
}
