package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["conf.go"],
    importpath = "go-common/app/interface/main/growup/conf",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//library/cache/redis:go_default_library",
        "//library/conf:go_default_library",
        "//library/database/hbase.v2:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/rpc:go_default_library",
        "//library/net/rpc/warden:go_default_library",
        "//library/net/trace:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/BurntSushi/toml:go_default_library",
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
