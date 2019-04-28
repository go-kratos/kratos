package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "client.go",
        "elastic.go",
        "result.go",
        "sug.go",
        "wild.go",
    ],
    importpath = "go-common/app/interface/main/tv/model/search",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/tv/model:go_default_library",
        "//app/interface/main/tv/model/thirdp:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/siddontang/go-mysql/mysql:go_default_library",
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
