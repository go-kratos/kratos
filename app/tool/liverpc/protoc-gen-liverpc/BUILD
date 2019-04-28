package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_binary",
    "go_test",
    "go_library",
)

go_binary(
    name = "protoc-gen-liverpc",
    embed = [":go_default_library"],
    tags = ["automanaged"],
)

go_test(
    name = "go_default_test",
    srcs = [
        "command_line_test.go",
        "generator_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "@com_github_golang_protobuf//proto:go_default_library",
        "@com_github_golang_protobuf//protoc-gen-go/plugin:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "command_line.go",
        "generator.go",
        "go_naming.go",
        "main.go",
    ],
    importpath = "go-common/app/tool/liverpc/protoc-gen-liverpc",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/tool/liverpc/protoc-gen-liverpc/gen:go_default_library",
        "//app/tool/liverpc/protoc-gen-liverpc/gen/stringutils:go_default_library",
        "//app/tool/liverpc/protoc-gen-liverpc/gen/typemap:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
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
        "//app/tool/liverpc/protoc-gen-liverpc/gen:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
