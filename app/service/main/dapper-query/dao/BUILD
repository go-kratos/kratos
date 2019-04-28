load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "dao.go",
        "util.go",
    ],
    importpath = "go-common/app/service/main/dapper-query/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/dapper-query/conf:go_default_library",
        "//app/service/main/dapper-query/model:go_default_library",
        "//library/log:go_default_library",
        "//vendor/github.com/dgryski/go-farm:go_default_library",
        "//vendor/github.com/influxdata/influxdb/client/v2:go_default_library",
        "//vendor/github.com/tsuna/gohbase:go_default_library",
        "//vendor/github.com/tsuna/gohbase/filter:go_default_library",
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

go_test(
    name = "go_default_test",
    srcs = ["dao_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/dapper-query/conf:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
        "//vendor/golang.org/x/sys/unix:go_default_library",
    ],
)
