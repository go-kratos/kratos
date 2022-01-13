package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoParameters(t *testing.T) {
	path := "/test/noparams"
	m := buildPathVars(path)
	assert.Emptyf(t, m, "Map should be empty")
}

func TestSingleParam(t *testing.T) {
	path := "/test/{message.id}"
	m := buildPathVars(path)
	assert.Len(t, m, 1)
	assert.Empty(t, m["message.id"])
}

func TestTwoParametersReplacement(t *testing.T) {
	path := "/test/{message.id}/{message.name=messages/*}"
	m := buildPathVars(path)
	assert.Len(t, m, 2)
	assert.Empty(t, m["message.id"])
	assert.NotEmpty(t, m["message.name"])
	assert.Equal(t, *m["message.name"], "messages/*")
}

func TestNoReplacePath(t *testing.T) {
	path := "/test/{message.id=test}"
	assert.Equal(t, "/test/{message.id:test}", replacePath("message.id", "test", path))

	path = "/test/{message.id=test/*}"
	assert.Equal(t, "/test/{message.id:test/.*}", replacePath("message.id", "test/*", path))
}

func TestReplacePath(t *testing.T) {
	path := "/test/{message.id}/{message.name=messages/*}"
	newPath := replacePath("message.name", "messages/*", path)
	assert.Equal(t, "/test/{message.id}/{message.name:messages/.*}", newPath)
}

func TestIteration(t *testing.T) {
	path := "/test/{message.id}/{message.name=messages/*}"
	vars := buildPathVars(path)
	for v, s := range vars {
		if s != nil {
			path = replacePath(v, *s, path)
		}
	}
	assert.Equal(t, "/test/{message.id}/{message.name:messages/.*}", path)
}

func TestIterationMiddle(t *testing.T) {
	path := "/test/{message.name=messages/*}/books"
	vars := buildPathVars(path)
	for v, s := range vars {
		if s != nil {
			path = replacePath(v, *s, path)
		}
	}
	assert.Equal(t, "/test/{message.name:messages/.*}/books", path)
}
