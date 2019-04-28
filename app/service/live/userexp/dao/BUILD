load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "dao.go",
        "exp.go",
        "exp_log.go",
        "mc_level.go",
    ],
    importpath = "go-common/app/service/live/userexp/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/live/userexp/conf:go_default_library",
        "//app/service/live/userexp/model:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/log:go_default_library",
        "//library/queue/databus/report:go_default_library",
        "//library/time:go_default_library",
        "//library/xstr:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "exp_test.go",
        "mc_level_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/service/live/userexp/conf:go_default_library",
        "//app/service/live/userexp/model:go_default_library",
        "//library/log:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
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
