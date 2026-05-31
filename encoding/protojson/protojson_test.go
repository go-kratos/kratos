package protojson

import (
	"reflect"
	"strings"
	"testing"

	testData "github.com/go-kratos/kratos/v3/internal/testdata/encoding"
)

func TestName(t *testing.T) {
	c := new(codec)
	if !reflect.DeepEqual(c.Name(), "protojson") {
		t.Errorf("expected %v, got %v", "protojson", c.Name())
	}
}

func TestCodec(t *testing.T) {
	c := new(codec)
	model := testData.TestModel{
		Id:    1,
		Name:  "go-kratos",
		Hobby: []string{"1", "2"},
	}

	data, err := c.Marshal(&model)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	if got, want := strings.ReplaceAll(string(data), " ", ""), `{"id":"1","name":"go-kratos","hobby":["1","2"],"attrs":{}}`; got != want {
		t.Errorf("Marshal() = %s, want %s", got, want)
	}

	var out testData.TestModel
	if err := c.Unmarshal([]byte(`{"id":"1","name":"go-kratos","hobby":["1","2"],"unknown":"discarded"}`), &out); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if !reflect.DeepEqual(out.Id, model.Id) {
		t.Errorf("Id = %d, want %d", out.Id, model.Id)
	}
	if !reflect.DeepEqual(out.Name, model.Name) {
		t.Errorf("Name = %s, want %s", out.Name, model.Name)
	}
	if !reflect.DeepEqual(out.Hobby, model.Hobby) {
		t.Errorf("Hobby = %v, want %v", out.Hobby, model.Hobby)
	}
}

func TestCodecRejectsNonProtoMessage(t *testing.T) {
	c := new(codec)

	if _, err := c.Marshal(struct{}{}); err == nil {
		t.Fatal("Marshal() error = nil, want error")
	}
	if err := c.Unmarshal([]byte("{}"), &struct{}{}); err == nil {
		t.Fatal("Unmarshal() error = nil, want error")
	}
}
