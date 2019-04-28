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
        "elec_mc.go",
        "elec_mysql.go",
        "grpc.go",
        "http.go",
        "mc.go",
        "mysql.go",
        "redis.go",
    ],
    importpath = "go-common/app/job/main/ugcpay/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/job/main/ugcpay/conf:go_default_library",
        "//app/job/main/ugcpay/model:go_default_library",
        "//app/service/main/ugcpay-rank/api/v1:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/stat/prom:go_default_library",
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
        "elec_mc_test.go",
        "elec_mysql_test.go",
        "grpc_test.go",
        "http_test.go",
        "mc_test.go",
        "mysql_test.go",
        "redis_test.go",
    ],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = [
        "//app/job/main/ugcpay/conf:go_default_library",
        "//app/job/main/ugcpay/model:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)
