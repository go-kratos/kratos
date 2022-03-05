package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

type Remark struct {
	mixin.Schema
}

func (Remark) Fields() []ent.Field {
	return []ent.Field{
		field.String("remark").
			Comment("说明").
			Default("").
			MaxLen(256),
	}
}
