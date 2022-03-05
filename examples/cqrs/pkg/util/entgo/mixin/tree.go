package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
)

type Tree struct {
	mixin.Schema
}

// Fields of the Tree.
func (Tree) Fields() []ent.Field {
	return []ent.Field{
		/*
		* 树结构编码,用于快速查找, 每一层由4位字符组成,用-分割
		* 如第一层:0001 第二层:0001-0001 第三层:0001-0001-0001
		 */
		field.Uint64("parent_id").
			Comment("父级类别").
			Unique().
			Immutable(),
		field.String("tree_path").
			Comment("树路径"),
		field.String("tree_index").
			Comment("排序序号"),
		field.Uint32("tree_level").
			Comment("树层级"),
	}
}

// Indexes of the Tree.
func (Tree) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("tree_path"),
	}
}
