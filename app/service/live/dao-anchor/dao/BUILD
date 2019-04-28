package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "dao.go",
        "mysql.go",
        "redis.go",
        "redisKey.go",
    ],
    importpath = "go-common/app/service/live/dao-anchor/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/common/live/library/lrucache:go_default_library",
        "//app/service/live/dao-anchor/api/grpc/v1:go_default_library",
        "//app/service/live/dao-anchor/conf:go_default_library",
        "//app/service/live/dao-anchor/model:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/queue/databus:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//vendor/github.com/json-iterator/go:go_default_library",
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

go_test(
    name = "go_default_test",
    srcs = [
        "dao_test.go",
        "mysql_test.go",
    ],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = [
        "//app/service/live/dao-anchor/api/grpc/v1:go_default_library",
        "//app/service/live/dao-anchor/conf:go_default_library",
        "//app/service/live/dao-anchor/model:go_default_library",
        "//library/log:go_default_library",
        "//library/queue/databus:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)
