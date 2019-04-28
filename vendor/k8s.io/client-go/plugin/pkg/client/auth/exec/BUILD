package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["exec.go"],
    importpath = "k8s.io/client-go/plugin/pkg/client/auth/exec",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/golang/glog:go_default_library",
        "//vendor/golang.org/x/crypto/ssh/terminal:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime/schema:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime/serializer:go_default_library",
        "//vendor/k8s.io/client-go/pkg/apis/clientauthentication:go_default_library",
        "//vendor/k8s.io/client-go/pkg/apis/clientauthentication/v1alpha1:go_default_library",
        "//vendor/k8s.io/client-go/pkg/apis/clientauthentication/v1beta1:go_default_library",
        "//vendor/k8s.io/client-go/tools/clientcmd/api:go_default_library",
        "//vendor/k8s.io/client-go/transport:go_default_library",
        "//vendor/k8s.io/client-go/util/connrotation:go_default_library",
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
    srcs = [":package-srcs"],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
