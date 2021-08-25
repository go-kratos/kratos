package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_defaultDecoder(t *testing.T) {
	src := &KeyValue{
		Key:    "service",
		Value:  []byte("config"),
		Format: "",
	}
	target := make(map[string]interface{}, 0)
	err := defaultDecoder(src, target)
	assert.Nil(t, err)
	assert.Equal(t, map[string]interface{}{
		"service": []byte("config"),
	}, target)

	src = &KeyValue{
		Key:    "service.name.alias",
		Value:  []byte("2233"),
		Format: "",
	}
	target = make(map[string]interface{}, 0)
	err = defaultDecoder(src, target)
	assert.Nil(t, err)
	assert.Equal(t, map[string]interface{}{
		"service": map[string]interface{}{
			"name": map[string]interface{}{
				"alias": []byte("2233"),
			},
		},
	}, target)
}

func Test_defaultResolver(t *testing.T) {
	type args struct {
		input map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := defaultResolver(tt.args.input); (err != nil) != tt.wantErr {
				t.Errorf("defaultResolver() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
