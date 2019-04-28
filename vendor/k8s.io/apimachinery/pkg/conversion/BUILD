package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "converter.go",
        "deep_equal.go",
        "doc.go",
        "helper.go",
    ],
    importpath = "k8s.io/apimachinery/pkg/conversion",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = ["//vendor/k8s.io/apimachinery/third_party/forked/golang/reflect:go_default_library"],
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
        "//vendor/k8s.io/apimachinery/pkg/conversion/queryparams:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
