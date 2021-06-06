package metadata

import (
	"fmt"
	"strings"
	"testing"
)

func TestMedata(t *testing.T) {
	fmt.Println(strings.ToLower("Content-Type"))
	fmt.Println(strings.ToLower("content-type"))
	fmt.Println(strings.ToLower("content-Type"))
	fmt.Println(strings.ToLower("Content-type"))

}
