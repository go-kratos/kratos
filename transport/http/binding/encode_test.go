package binding

import (
	"fmt"
	"testing"

	"github.com/go-kratos/kratos/v2/internal/testdata/binding"
)

func TestProtoPath(t *testing.T) {
	url := EncodeURL("http://helloworld.Greeter/helloworld/{name}/sub/{sub.name}", &binding.HelloRequest{Name: "test", Sub: &binding.Sub{Name: "2233!!!"}}, false)
	fmt.Println(url)
	if url != `http://helloworld.Greeter/helloworld/test/sub/2233!!!` {
		t.Fatalf("proto path not expected!actual: %s ", url)
	}

	url = EncodeURL("http://helloworld.Greeter/helloworld/{name}/sub/{sub.naming}", &binding.HelloRequest{Name: "test", Sub: &binding.Sub{Name: "5566!!!"}}, false)
	fmt.Println(url)
	if url != `http://helloworld.Greeter/helloworld/test/sub/5566!!!` {
		t.Fatalf("proto path not expected!actual: %s ", url)
	}

	url = EncodeURL("http://helloworld.Greeter/helloworld/sub", &binding.HelloRequest{Name: "test", Sub: &binding.Sub{Name: "2233!!!"}}, false)
	fmt.Println(url)
	if url != `http://helloworld.Greeter/helloworld/sub` {
		t.Fatalf("proto path not expected!actual: %s ", url)
	}

	url = EncodeURL("http://helloworld.Greeter/helloworld/{name}/sub/{sub.name}", &binding.HelloRequest{Name: "test"}, false)
	fmt.Println(url)
	if url != `http://helloworld.Greeter/helloworld/test/sub/` {
		t.Fatalf("proto path not expected!actual: %s ", url)
	}

	url = EncodeURL("http://helloworld.Greeter/helloworld/{name}/sub/{sub.name33}", &binding.HelloRequest{Name: "test"}, false)
	fmt.Println(url)
	if url != `http://helloworld.Greeter/helloworld/test/sub/{sub.name33}` {
		t.Fatalf("proto path not expected!actual: %s ", url)
	}

	url = EncodeURL("http://helloworld.Greeter/helloworld/{name}/sub", &binding.HelloRequest{
		Name: "go",
		Sub:  &binding.Sub{Name: "kratos"},
	}, true)
	fmt.Println(url)
	if url != `http://helloworld.Greeter/helloworld/go/sub?sub.naming=kratos` {
		t.Fatalf("proto path not expected!actual: %s ", url)
	}
}
