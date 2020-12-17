package example

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	c := new(ExampleConfig)
	err := config.ApplyYAML(`
address: localhost
timeout: 99
`, c)
	fmt.Printf("pb %v\n", c)
	if err != nil {
		panic(err)
	}

	x, err := New(c.Address, ApplyOptions(c)...)
	if err != nil {
		panic(err)
	}
	fmt.Printf("after %v", x)
	assert.NotNil(t, x)
}