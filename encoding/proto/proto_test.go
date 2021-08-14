package proto

import (
	"github.com/go-kratos/kratos/v2/encoding/proto/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestName(t *testing.T) {
	c := new(codec)
	assert.Equal(t, c.Name(), "proto")
}

func TestCodec(t *testing.T) {
	c := new(codec)

	model := test.TestModel{
		Id:    1,
		Name:  "kratos",
		Hobby: []string{"study", "eat", "play"},
	}

	m, err := c.Marshal(&model)
	assert.Nil(t, err)

	var res test.TestModel

	err = c.Unmarshal(m, &res)
	assert.Nil(t, err)

	assert.Equal(t, res.Id,model.Id)
	assert.Equal(t, res.Name,model.Name)
	assert.Equal(t, res.Hobby,model.Hobby)
}
