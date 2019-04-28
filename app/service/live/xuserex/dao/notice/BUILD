package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "dao.cache.go",
        "mc.cache.go",
        "roomNotice.go",
    ],
    importpath = "go-common/app/service/live/xuserex/dao/notice",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/live/resource/sdk:go_default_library",
        "//app/service/live/xuser/api/grpc/v1:go_default_library",
        "//app/service/live/xuserex/api/grpc/v1:go_default_library",
        "//app/service/live/xuserex/conf:go_default_library",
        "//app/service/live/xuserex/model/roomNotice:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/database/hbase.v2:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/stat/prom:go_default_library",
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

go_test(
    name = "go_default_test",
    srcs = [
        "dao.cache_test.go",
        "mc.cache_test.go",
        "roomNotice_test.go",
    ],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = [
        "//app/service/live/xuserex/conf:go_default_library",
        "//app/service/live/xuserex/model/roomNotice:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)
