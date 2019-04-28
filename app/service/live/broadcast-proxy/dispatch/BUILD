package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = ["matcher_test.go"],
    embed = [":go_default_library"],
    tags = ["automanaged"],
)

go_library(
    name = "go_default_library",
    srcs = [
        "matcher.go",
        "min_heap.go",
        "sinaip.go",
    ],
    importpath = "go-common/app/service/live/broadcast-proxy/dispatch",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/live/broadcast-proxy/conf:go_default_library",
        "//app/service/live/broadcast-proxy/expr:go_default_library",
        "//library/log:go_default_library",
        "//vendor/github.com/ipipdotnet/ipdb-go:go_default_library",
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
