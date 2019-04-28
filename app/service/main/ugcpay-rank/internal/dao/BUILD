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
        "grpc.go",
        "localcache.go",
        "mc.cache.go",
        "mc.extra.go",
        "mysql.go",
    ],
    importpath = "go-common/app/service/main/ugcpay-rank/internal/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/account/api:go_default_library",
        "//app/service/main/account/model:go_default_library",
        "//app/service/main/ugcpay-rank/internal/conf:go_default_library",
        "//app/service/main/ugcpay-rank/internal/model:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/log:go_default_library",
        "//library/stat/prom:go_default_library",
        "//library/xstr:go_default_library",
        "//vendor/github.com/bluele/gcache:go_default_library",
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
        "grpc_test.go",
        "localcache_test.go",
        "mc.cache_test.go",
        "mc.extra_test.go",
        "mysql_test.go",
    ],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = [
        "//app/service/main/ugcpay-rank/internal/conf:go_default_library",
        "//app/service/main/ugcpay-rank/internal/model:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)
