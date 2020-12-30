package config

import (
	"encoding/json"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/pelletier/go-toml"
)

// ApplyJSON unmarshals a JSON string into a proto message.
// Unknown fields are allowed
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

// ApplyYAML unmarshals a YAML string into a proto message.
// Unknown fields are allowed.
func ApplyYAML(yml string, pb proto.Message) error {
	js, err := yaml.YAMLToJSON([]byte(yml))
	if err != nil {
		return err
	}
	return ApplyJSON(string(js), pb)
}

// ApplyTOML unmarshals a TOML string into a proto message.
func ApplyTOML(tm string, pb proto.Message) error {
	tree, err := toml.Load(tm)
	if err != nil {
		return err
	}
	js, err := json.Marshal(tree.ToMap())
	if err != nil {
		return err
	}
	return ApplyJSON(string(js), pb)
}
