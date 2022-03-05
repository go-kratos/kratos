package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Sensor holds the schema definition for the Sensor entity.
type Sensor struct {
	ent.Schema
}

func (Sensor) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "sensors",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
	}
}

// Fields of the Sensor.
func (Sensor) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			SchemaType(map[string]string{
				dialect.MySQL:    "INT",
				dialect.Postgres: "SERIAL",
			}),
		field.String("type").
			Comment("传感器类型").
			Default("").
			MaxLen(50),
		field.String("location").
			Comment("所在位置").
			Default("").
			MaxLen(50),
	}
}

// Edges of the Sensor.
func (Sensor) Edges() []ent.Edge {
	return nil
}

// Indexes of the Sensor.
func (Sensor) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id").StorageKey("sensor_pkey"),
	}
}
