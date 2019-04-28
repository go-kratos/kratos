package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "code.go",
        "dao.go",
        "http.go",
        "memcache.go",
        "mysql.go",
        "remote.go",
    ],
    importpath = "go-common/app/admin/main/coupon/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/coupon/conf:go_default_library",
        "//app/admin/main/coupon/model:go_default_library",
        "//app/service/openplatform/pgc-season/api/grpc/season/v1:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/database/elastic:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
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
        "code_test.go",
        "dao_test.go",
        "http_test.go",
        "memcache_test.go",
        "mysql_test.go",
        "remote_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/coupon/conf:go_default_library",
        "//app/admin/main/coupon/model:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
        "//vendor/gopkg.in/h2non/gock.v1:go_default_library",
    ],
)
