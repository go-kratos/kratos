package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultDecoder(t *testing.T) {
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
