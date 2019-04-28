package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "captchaAnti.go",
        "dao.go",
        "gee.go",
        "liveCaptcha.go",
        "redis.go",
    ],
    importpath = "go-common/app/service/live/xcaptcha/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/live/captcha/api/liverpc:go_default_library",
        "//app/service/live/captcha/api/liverpc/v0:go_default_library",
        "//app/service/live/captcha/api/liverpc/v1:go_default_library",
        "//app/service/live/xcaptcha/conf:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/log:go_default_library",
        "//library/net/rpc/liverpc:go_default_library",
        "//library/queue/databus:go_default_library",
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
