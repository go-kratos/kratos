package binding

import (
	"fmt"
	"log"
	"net/http"
)

type Arg struct {
	Max   int64 `form:"max" validate:"max=10"`
	Min   int64 `form:"min" validate:"min=2"`
	Range int64 `form:"range" validate:"min=1,max=10"`
	// use split option to split arg 1,2,3 into slice [1 2 3]
	// otherwise slice type with parse  url.Values (eg:a=b&a=c) default.
	Slice []int64 `form:"slice,split" validate:"min=1"`
}

func ExampleBinding() {
	req := initHTTP("max=9&min=3&range=3&slice=1,2,3")
	arg := new(Arg)
	if err := Form.Bind(req, arg); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("arg.Max %d\narg.Min %d\narg.Range %d\narg.Slice %v", arg.Max, arg.Min, arg.Range, arg.Slice)
	// Output:
	// arg.Max 9
	// arg.Min 3
	// arg.Range 3
	// arg.Slice [1 2 3]
}

func initHTTP(params string) (req *http.Request) {
	req, _ = http.NewRequest("GET", "http://api.bilibili.com/test?"+params, nil)
	req.ParseForm()
	return
}
