package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"kratos-cqrs/pkg/util/sonyflake"
)

type SnowflackId struct {
	mixin.Schema
}

func (SnowflackId) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id").
			Comment("id").
			DefaultFunc(sonyflake.GenerateSonyflake()).
			Positive().
			Immutable().
			StructTag(`json:"id,omitempty"`),
	}
}

// Indexes of the SnowflackId.
func (SnowflackId) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id"),
	}
}
