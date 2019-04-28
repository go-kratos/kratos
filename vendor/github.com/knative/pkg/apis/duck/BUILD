package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "cached.go",
        "doc.go",
        "enqueue.go",
        "interface.go",
        "patch.go",
        "proxy.go",
        "register.go",
        "typed.go",
        "unstructured.go",
        "verify.go",
    ],
    importpath = "github.com/knative/pkg/apis/duck",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/evanphx/json-patch:go_default_library",
        "//vendor/github.com/google/go-cmp/cmp:go_default_library",
        "//vendor/github.com/knative/pkg/apis:go_default_library",
        "//vendor/github.com/mattbaird/jsonpatch:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/apis/meta/v1/unstructured:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime/schema:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/watch:go_default_library",
        "//vendor/k8s.io/client-go/dynamic:go_default_library",
        "//vendor/k8s.io/client-go/tools/cache:go_default_library",
    ],
)

filegroup(
    name = "package-srcs",
    srcs = glob(["**"]),
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
)

filegroup(
    name = "all-srcs",
    srcs = [
        ":package-srcs",
        "//vendor/github.com/knative/pkg/apis/duck/v1alpha1:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
