package grpc

import (
	"reflect"
	"testing"

	"google.golang.org/protobuf/types/known/structpb"
)

func TestCodec(t *testing.T) {
	in, err := structpb.NewStruct(map[string]any{"Golang": "Kratos"})
	if err != nil {
		t.Errorf("grpc codec create input data error:%v", err)
	}
	c := codec{}
	if c.Name() != "json" {
		t.Errorf("grpc codec name want %v, got %v", "json", c.Name())
	}
	data, err := c.Marshal(in)
	if err != nil {
		t.Errorf("grpc codec marshal error:%v", err)
	}
	if string(data) != `{"Golang":"Kratos"}` {
		t.Errorf("grpc codec marshal want %v, got %v", `{"Golang":"Kratos"}`, string(data))
	}
	out := &structpb.Struct{}
	err = c.Unmarshal(data, out)
	if err != nil {
		t.Errorf("grpc codec unmarshal error:%v", err)
	}
	if !reflect.DeepEqual(in, out) {
		t.Errorf("grpc codec want %v, got %v", in, out)
	}
}
