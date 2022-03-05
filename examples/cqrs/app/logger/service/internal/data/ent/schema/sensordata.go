package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// SensorData holds the schema definition for the SensorData entity.
type SensorData struct {
	ent.Schema
}

func (SensorData) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "sensor_data",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
	}
}

// Fields of the SensorData.
func (SensorData) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("time").
			Comment("时间戳").
			Nillable().Optional(),
		field.Int("sensor_id").
			Comment("传感器ID"),
		field.Float("temperature").
			Comment("温度"),
		field.Float("cpu").
			Comment("CPU使用率"),
	}
}

// Edges of the SensorData.
func (SensorData) Edges() []ent.Edge {
	return nil
}
