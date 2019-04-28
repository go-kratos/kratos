package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = ["api_test.go"],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = [
        "//app/service/live/xuser/conf:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = ["api.go"],
    importpath = "go-common/app/service/live/xuser/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/live/banned_service/api/liverpc:go_default_library",
        "//app/service/live/fans_medal/api/liverpc:go_default_library",
        "//app/service/live/room/api/liverpc:go_default_library",
        "//app/service/live/xuser/conf:go_default_library",
        "//library/net/rpc/liverpc:go_default_library",
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
        "//app/service/live/xuser/dao/account:all-srcs",
        "//app/service/live/xuser/dao/exp:all-srcs",
        "//app/service/live/xuser/dao/guard:all-srcs",
        "//app/service/live/xuser/dao/roomAdmin:all-srcs",
        "//app/service/live/xuser/dao/vip:all-srcs",
        "//app/service/live/xuser/dao/xanchor:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
