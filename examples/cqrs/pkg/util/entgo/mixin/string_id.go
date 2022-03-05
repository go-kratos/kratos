package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"regexp"
)

type StringId struct {
	mixin.Schema
}

func (StringId) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Comment("id").
			MaxLen(25).
			NotEmpty().
			Unique().
			Immutable().
			Match(regexp.MustCompile("^[0-9a-zA-Z_\\-]+$")).
			StructTag(`json:"id,omitempty"`),
	}
}

// Indexes of the StringId.
func (StringId) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id"),
	}
}
