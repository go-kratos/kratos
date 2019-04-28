package model

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildFaceURL(t *testing.T) {
	face := BuildFaceURL("/bfs/facepri/4a259f715b63157f24f76521231480438058436e.jpg")
	u, err := url.Parse(face)
	assert.Nil(t, err)
	assert.NotNil(t, u.Query().Get("token"))
}
