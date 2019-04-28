package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "logging.go",
        "main.go",
        "version.go",
        "wrappers.go",
    ],
    importpath = "go-common/app/tool/liverpc/protoc-gen-liverpc/gen",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/tool/liverpc/protoc-gen-liverpc/gen/stringutils:go_default_library",
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
    srcs = [
        ":package-srcs",
        "//app/tool/liverpc/protoc-gen-liverpc/gen/stringutils:all-srcs",
        "//app/tool/liverpc/protoc-gen-liverpc/gen/typemap:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
