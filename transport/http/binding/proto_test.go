package binding

import (
	"fmt"
	"testing"
)

func TestProtoPath(t *testing.T) {
	url := EncodePath("POST", "http://helloworld.Greeter/helloworld/{name}/sub/{sub.name}", &HelloRequest{Name: "test", Sub: &Sub{Name: "2233!!!"}})
	fmt.Println(url)
	if url != `http://helloworld.Greeter/helloworld/test/sub/2233!!!` {
		t.Fatalf("proto path not expected!actual: %s ", url)
	}

	url = EncodePath("POST", "http://helloworld.Greeter/helloworld/sub", &HelloRequest{Name: "test", Sub: &Sub{Name: "2233!!!"}})
	fmt.Println(url)
	if url != `http://helloworld.Greeter/helloworld/sub` {
		t.Fatalf("proto path not expected!actual: %s ", url)
	}

	url = EncodePath("POST", "http://helloworld.Greeter/helloworld/{name}/sub/{sub.name}", &HelloRequest{Name: "test"})
	fmt.Println(url)
	if url != `http://helloworld.Greeter/helloworld/test/sub/` {
		t.Fatalf("proto path not expected!actual: %s ", url)
	}

	url = EncodePath("POST", "http://helloworld.Greeter/helloworld/{name}/sub/{sub.name33}", &HelloRequest{Name: "test"})
	fmt.Println(url)
	if url != `http://helloworld.Greeter/helloworld/test/sub/{sub.name33}` {
		t.Fatalf("proto path not expected!actual: %s ", url)
	}
}
