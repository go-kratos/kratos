package config

import (
	"github.com/ghodss/yaml"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"strings"
)


type Spec interface{}

func ApplyYAML(yml string, pb proto.Message) error {
	js, err := yaml.YAMLToJSON([]byte(yml))
	if err != nil {
		return err
	}
	return ApplyJSON(string(js), pb)
}


func ApplyJSON(js string, pb proto.Message) error {
	reader := strings.NewReader(js)
	m := jsonpb.Unmarshaler{}
	if err := m.Unmarshal(reader, pb); err != nil {
		m.AllowUnknownFields = true
		reader.Reset(js)
		return m.Unmarshal(reader, pb)
	}
	return nil
}
