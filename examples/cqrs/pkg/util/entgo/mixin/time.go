package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"time"
)

type Time struct {
	mixin.Schema
}

func (Time) Fields() []ent.Field {
	return []ent.Field{
		field.Time("create_time").
			Comment("创建时间").
			Immutable().
			Default(time.Now),
		field.Time("update_time").
			Comment("更新时间").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}
