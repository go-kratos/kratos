package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	_mixin "entgo.io/ent/schema/mixin"
)

type TenantId struct {
	_mixin.Schema
}

// Fields of the TenantId.
func (TenantId) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("tenant_id").
			Comment("租户ID").
			Default(0),
	}
}

// Indexes of the TenantId.
func (TenantId) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("tenant_id"),
	}
}
