package http

import (
	"testing"

	"google.golang.org/protobuf/types/known/fieldmaskpb"

	"github.com/go-kratos/kratos/v3/internal/testdata/binding"
)

func TestBuildPath(t *testing.T) {
	tests := []struct {
		name         string
		pathTemplate string
		request      *binding.HelloRequest
		opts         []BuildPathOption
		want         string
	}{
		{
			name:         "path",
			pathTemplate: "/helloworld/{name}",
			request:      &binding.HelloRequest{Name: "kratos"},
			want:         "/helloworld/kratos",
		},
		{
			name:         "query",
			pathTemplate: "/helloworld/{name}",
			request:      &binding.HelloRequest{Name: "kratos", Sub: &binding.Sub{Name: "go"}},
			opts:         []BuildPathOption{WithQueryParams()},
			want:         "/helloworld/kratos?sub.naming=go",
		},
		{
			name:         "resource name",
			pathTemplate: "/v1/{name=publishers/*/books/*}",
			request:      &binding.HelloRequest{Name: "publishers/go/books/kratos"},
			opts:         []BuildPathOption{WithQueryParams()},
			want:         "/v1/publishers/go/books/kratos",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildPath(tt.pathTemplate, tt.request, tt.opts...); got != tt.want {
				t.Errorf("expected %s got %s", tt.want, got)
			}
		})
	}
}

func TestBuildPathCompatibility(t *testing.T) {
	tests := []struct {
		pathTemplate string
		request      *binding.HelloRequest
		opts         []BuildPathOption
		want         string
	}{
		{
			pathTemplate: "http://helloworld.Greeter/helloworld/{name}/sub/{sub.naming}",
			request:      &binding.HelloRequest{Name: "test", Sub: &binding.Sub{Name: "2233!!!!"}},
			want:         "http://helloworld.Greeter/helloworld/test/sub/2233!!!!",
		},
		{
			pathTemplate: "http://helloworld.Greeter/helloworld/{name}/sub/{sub.naming}",
			request:      nil,
			want:         "http://helloworld.Greeter/helloworld/{name}/sub/{sub.naming}",
		},
		{
			pathTemplate: "http://helloworld.Greeter/helloworld/{}/sub/{sub.naming}",
			request:      &binding.HelloRequest{Name: "test", Sub: &binding.Sub{Name: "hello"}},
			want:         "http://helloworld.Greeter/helloworld/{}/sub/hello",
		},
		{
			pathTemplate: "http://helloworld.Greeter/helloworld/{}/sub/{sub.name.cc}",
			request:      &binding.HelloRequest{Name: "test", Sub: &binding.Sub{Name: "hello"}},
			want:         "http://helloworld.Greeter/helloworld/{}/sub/",
		},
		{
			pathTemplate: "http://helloworld.Greeter/helloworld/{}/sub/{test_repeated}",
			request:      &binding.HelloRequest{Name: "test", Sub: &binding.Sub{Name: "hello"}, TestRepeated: []string{"123", "456"}},
			want:         "http://helloworld.Greeter/helloworld/{}/sub/123",
		},
		{
			pathTemplate: "http://helloworld.Greeter/helloworld/{name}/sub/{sub.naming}",
			request:      &binding.HelloRequest{Name: "test", Sub: &binding.Sub{Name: "5566!!!"}},
			want:         "http://helloworld.Greeter/helloworld/test/sub/5566!!!",
		},
		{
			pathTemplate: "http://helloworld.Greeter/helloworld/sub",
			request:      &binding.HelloRequest{Name: "test", Sub: &binding.Sub{Name: "2233!!!"}},
			want:         "http://helloworld.Greeter/helloworld/sub",
		},
		{
			pathTemplate: "http://helloworld.Greeter/helloworld/{name}/sub/{sub.name}",
			request:      &binding.HelloRequest{Name: "test"},
			want:         "http://helloworld.Greeter/helloworld/test/sub/",
		},
		{
			pathTemplate: "http://helloworld.Greeter/helloworld/{name}/sub",
			request:      &binding.HelloRequest{Name: "go", Sub: &binding.Sub{Name: "kratos"}},
			opts:         []BuildPathOption{WithQueryParams()},
			want:         "http://helloworld.Greeter/helloworld/go/sub?sub.naming=kratos",
		},
		{
			pathTemplate: "http://helloworld.Greeter/helloworld/{name=publishers/*/books/*}/sub",
			request:      &binding.HelloRequest{Name: "publishers/go/books/kratos", Sub: &binding.Sub{Name: "kratos"}},
			opts:         []BuildPathOption{WithQueryParams()},
			want:         "http://helloworld.Greeter/helloworld/publishers/go/books/kratos/sub?sub.naming=kratos",
		},
		{
			pathTemplate: "http://helloworld.Greeter/helloworld/{name=**}/sub",
			request:      &binding.HelloRequest{Name: "publishers/go/books/kratos", Sub: &binding.Sub{Name: "kratos"}},
			opts:         []BuildPathOption{WithQueryParams()},
			want:         "http://helloworld.Greeter/helloworld/publishers/go/books/kratos/sub?sub.naming=kratos",
		},
		{
			pathTemplate: "http://helloworld.Greeter/helloworld/{sub.naming=publishers/*}",
			request:      &binding.HelloRequest{Sub: &binding.Sub{Name: "publishers/kratos"}},
			want:         "http://helloworld.Greeter/helloworld/publishers/kratos",
		},
		{
			pathTemplate: "http://helloworld.Greeter/helloworld/sub/{sub.naming}",
			request:      &binding.HelloRequest{Sub: &binding.Sub{Name: "kratos"}, UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"name", "sub.naming"}}},
			want:         "http://helloworld.Greeter/helloworld/sub/kratos?updateMask=name,sub.naming",
		},
		{
			pathTemplate: "http://helloworld.Greeter/helloworld/sub/[{sub.naming}]",
			request:      &binding.HelloRequest{Sub: &binding.Sub{Name: "kratos"}},
			want:         "http://helloworld.Greeter/helloworld/sub/[kratos]",
		},
		{
			pathTemplate: "http://helloworld.Greeter/helloworld/[{name}]/sub/[{sub.naming}]",
			request:      &binding.HelloRequest{Name: "test", Sub: &binding.Sub{Name: "kratos"}},
			want:         "http://helloworld.Greeter/helloworld/[test]/sub/[kratos]",
		},
		{
			pathTemplate: "http://helloworld.Greeter/helloworld/[{}]/sub/[{sub.naming}]",
			request:      &binding.HelloRequest{Sub: &binding.Sub{Name: "kratos"}},
			want:         "http://helloworld.Greeter/helloworld/[{}]/sub/[kratos]",
		},
		{
			pathTemplate: "http://helloworld.Greeter/helloworld/[{}]/sub/[{sub.naming}]/{[]}",
			request:      &binding.HelloRequest{Sub: &binding.Sub{Name: "kratos"}},
			want:         "http://helloworld.Greeter/helloworld/[{}]/sub/[kratos]/{[]}",
		},
		{
			pathTemplate: "http://helloworld.Greeter/helloworld/{[sub]}/[{sub.naming}]",
			request:      &binding.HelloRequest{Sub: &binding.Sub{Name: "kratos"}},
			want:         "http://helloworld.Greeter/helloworld/{[sub]}/[kratos]",
		},
		{
			pathTemplate: "http://helloworld.Greeter/helloworld/{[name]}/[{sub.naming}]",
			request:      &binding.HelloRequest{Name: "test", Sub: &binding.Sub{Name: "kratos"}},
			want:         "http://helloworld.Greeter/helloworld/{[name]}/[kratos]",
		},
		{
			pathTemplate: "http://helloworld.Greeter/helloworld/{}/[]/[{sub.naming}]",
			request:      &binding.HelloRequest{Sub: &binding.Sub{Name: "kratos"}},
			want:         "http://helloworld.Greeter/helloworld/{}/[]/[kratos]",
		},
	}

	for _, test := range tests {
		if path := BuildPath(test.pathTemplate, test.request, test.opts...); path != test.want {
			t.Fatalf("want: %s, got: %s", test.want, path)
		}
	}
}

func BenchmarkBuildPath(b *testing.B) {
	benchmarks := []struct {
		name         string
		pathTemplate string
		msg          *binding.HelloRequest
		opts         []BuildPathOption
	}{
		{
			name:         "NoParams",
			pathTemplate: "http://helloworld.Greeter/helloworld/sub",
			msg: &binding.HelloRequest{
				Name: "test",
				Sub:  &binding.Sub{Name: "kratos"},
			},
		},
		{
			name:         "NoParamsWithQuery",
			pathTemplate: "http://helloworld.Greeter/helloworld/sub",
			msg: &binding.HelloRequest{
				Name: "test",
				Sub:  &binding.Sub{Name: "kratos"},
				UpdateMask: &fieldmaskpb.FieldMask{
					Paths: []string{"name", "sub.naming"},
				},
			},
			opts: []BuildPathOption{WithQueryParams()},
		},
		{
			name:         "WithParams",
			pathTemplate: "http://helloworld.Greeter/helloworld/{name}/sub/{sub.naming}",
			msg: &binding.HelloRequest{
				Name: "test",
				Sub:  &binding.Sub{Name: "kratos"},
			},
		},
		{
			name:         "WithParamsAndQuery",
			pathTemplate: "http://helloworld.Greeter/helloworld/{name}/sub/{sub.naming}",
			msg: &binding.HelloRequest{
				Name: "test",
				Sub:  &binding.Sub{Name: "kratos"},
				UpdateMask: &fieldmaskpb.FieldMask{
					Paths: []string{"name", "sub.naming"},
				},
			},
			opts: []BuildPathOption{WithQueryParams()},
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				BuildPath(bm.pathTemplate, bm.msg, bm.opts...)
			}
		})
	}
}
