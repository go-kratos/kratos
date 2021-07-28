package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTarget(t *testing.T) {
	target, err := parseTarget("localhost:8000", false)
	assert.Nil(t, err)
	assert.Equal(t, &Target{Scheme: "http", Authority: "localhost:8000"}, target)

	target, err = parseTarget("discovery:///demo", false)
	assert.Nil(t, err)
	assert.Equal(t, &Target{Scheme: "discovery", Authority: "", Endpoint: "demo"}, target)

	target, err = parseTarget("127.0.0.1:8000", false)
	assert.Nil(t, err)
	assert.Equal(t, &Target{Scheme: "http", Authority: "127.0.0.1:8000"}, target)

	target, err = parseTarget("https://127.0.0.1:8000", false)
	assert.Nil(t, err)
	assert.Equal(t, &Target{Scheme: "https", Authority: "127.0.0.1:8000"}, target)

	target, err = parseTarget("127.0.0.1:8000", true)
	assert.Nil(t, err)
	assert.Equal(t, &Target{Scheme: "https", Authority: "127.0.0.1:8000"}, target)
}
