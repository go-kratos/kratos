package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "clear_cache.go",
        "dahanghai.go",
        "guard.go",
        "log_debug.go",
    ],
    importpath = "go-common/app/service/live/xuser/service/guard/v1",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/live/room/api/liverpc:go_default_library",
        "//app/service/live/room/api/liverpc/v2:go_default_library",
        "//app/service/live/xuser/api/grpc/v1:go_default_library",
        "//app/service/live/xuser/conf:go_default_library",
        "//app/service/live/xuser/dao/account:go_default_library",
        "//app/service/live/xuser/dao/exp:go_default_library",
        "//app/service/live/xuser/dao/guard:go_default_library",
        "//app/service/live/xuser/model:go_default_library",
        "//app/service/live/xuser/model/dhh:go_default_library",
        "//library/cache:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/queue/databus:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
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
