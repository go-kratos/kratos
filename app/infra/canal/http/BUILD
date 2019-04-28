package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "canal.go",
        "http.go",
        "infoc.go",
    ],
    importpath = "go-common/app/infra/canal/http",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/infra/canal/conf:go_default_library",
        "//app/infra/canal/infoc:go_default_library",
        "//app/infra/canal/service:go_default_library",
        "//library/conf:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//vendor/github.com/BurntSushi/toml:go_default_library",
        "//vendor/github.com/siddontang/go-mysql/canal:go_default_library",
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
