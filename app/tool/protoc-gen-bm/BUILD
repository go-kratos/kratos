package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_binary",
    "go_library",
)

go_binary(
    name = "protoc-gen-bm",
    embed = [":go_default_library"],
    tags = ["automanaged"],
)

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importpath = "go-common/app/tool/protoc-gen-bm",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/tool/protoc-gen-bm/codegenerator:go_default_library",
        "//app/tool/protoc-gen-bm/genbm:go_default_library",
        "//app/tool/protoc-gen-bm/util:go_default_library",
        "//vendor/github.com/golang/glog:go_default_library",
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
        "//app/tool/protoc-gen-bm/codegenerator:all-srcs",
        "//app/tool/protoc-gen-bm/genbm:all-srcs",
        "//app/tool/protoc-gen-bm/generator:all-srcs",
        "//app/tool/protoc-gen-bm/jsonpb:all-srcs",
        "//app/tool/protoc-gen-bm/util:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
