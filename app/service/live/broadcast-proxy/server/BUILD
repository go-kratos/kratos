package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "dispatch.go",
        "http.go",
        "proxy.go",
        "zk.go",
    ],
    importpath = "go-common/app/service/live/broadcast-proxy/server",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/live/broadcast-proxy/conf:go_default_library",
        "//app/service/live/broadcast-proxy/dispatch:go_default_library",
        "//app/service/live/broadcast-proxy/grocery:go_default_library",
        "//library/log:go_default_library",
        "//vendor/github.com/ipipdotnet/ipdb-go:go_default_library",
        "//vendor/github.com/samuel/go-zookeeper/zk:go_default_library",
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
