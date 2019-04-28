package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "codec_factory.go",
        "negotiated_codec.go",
        "protobuf_extension.go",
    ],
    importpath = "k8s.io/apimachinery/pkg/runtime/serializer",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/k8s.io/apimachinery/pkg/runtime:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime/schema:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime/serializer/json:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime/serializer/protobuf:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime/serializer/recognizer:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime/serializer/versioning:go_default_library",
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
        "//vendor/k8s.io/apimachinery/pkg/runtime/serializer/json:all-srcs",
        "//vendor/k8s.io/apimachinery/pkg/runtime/serializer/protobuf:all-srcs",
        "//vendor/k8s.io/apimachinery/pkg/runtime/serializer/recognizer:all-srcs",
        "//vendor/k8s.io/apimachinery/pkg/runtime/serializer/streaming:all-srcs",
        "//vendor/k8s.io/apimachinery/pkg/runtime/serializer/versioning:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
