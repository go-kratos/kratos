package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
)

type AutoIncrementId struct {
	mixin.Schema
}

func (AutoIncrementId) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id").
			Comment("id").
			Positive().
			Immutable().
			Unique().
			StructTag(`json:"id,omitempty"`),
	}
}

// Indexes of the AutoIncrementId.
func (AutoIncrementId) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id"),
	}
}
