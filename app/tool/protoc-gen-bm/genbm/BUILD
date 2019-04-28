package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "bm_generate.go",
        "genbm.go",
        "http_descriptor.go",
    ],
    importpath = "go-common/app/tool/protoc-gen-bm/genbm",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/tool/protoc-gen-bm/generator:go_default_library",
        "//app/tool/protoc-gen-bm/util:go_default_library",
        "//vendor/github.com/golang/glog:go_default_library",
        "//vendor/google.golang.org/genproto/googleapis/api/annotations:go_default_library",
        "@com_github_golang_protobuf//proto:go_default_library",
        "@com_github_golang_protobuf//protoc-gen-go/descriptor:go_default_library",
        "@com_github_golang_protobuf//protoc-gen-go/plugin:go_default_library",
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
