package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "foo.go",
        "foo2.go",
    ],
    importpath = "go-common/app/interface/live/live-demo/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/live/live-demo/api/http:go_default_library",
        "//app/interface/live/live-demo/conf:go_default_library",
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
        "//app/interface/live/live-demo/service/v2:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
