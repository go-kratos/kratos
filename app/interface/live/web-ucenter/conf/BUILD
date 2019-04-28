package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["conf.go"],
    importpath = "go-common/app/interface/live/web-ucenter/conf",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//library/cache/memcache:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/conf:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode/tip:go_default_library",
        "//library/log:go_default_library",
        "//library/log/infoc:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/middleware/auth:go_default_library",
        "//library/net/http/blademaster/middleware/verify:go_default_library",
        "//library/net/rpc:go_default_library",
        "//library/net/rpc/liverpc:go_default_library",
        "//library/net/rpc/warden:go_default_library",
        "//library/net/trace:go_default_library",
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
