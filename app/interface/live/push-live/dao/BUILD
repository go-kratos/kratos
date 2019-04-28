package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "config_test.go",
        "d_test.go",
        "filter_test.go",
        "hbase_test.go",
        "push_test.go",
        "switch_test.go",
        "task_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/interface/live/push-live/conf:go_default_library",
        "//app/interface/live/push-live/model:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/log:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "blacklist.go",
        "config.go",
        "dao.go",
        "filter.go",
        "hbase.go",
        "push.go",
        "switch.go",
        "task.go",
    ],
    importpath = "go-common/app/interface/live/push-live/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/live/push-live/conf:go_default_library",
        "//app/interface/live/push-live/model:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/database/hbase.v2:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/stat/prom:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//library/xstr:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
        "//vendor/github.com/tsuna/gohbase/hrpc:go_default_library",
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
