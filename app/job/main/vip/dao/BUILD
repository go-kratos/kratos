package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "associate_test.go",
        "business_test.go",
        "coupon_old_test.go",
        "coupon_test.go",
        "dao_test.go",
        "http_test.go",
        "memcache_test.go",
        "mysql_test.go",
        "oldvip_test.go",
        "order_test.go",
        "push_test.go",
        "redis_test.go",
        "resource_test.go",
        "sync_test.go",
        "vip_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/job/main/vip/conf:go_default_library",
        "//app/job/main/vip/model:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
        "//vendor/gopkg.in/h2non/gock.v1:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "associate.go",
        "business.go",
        "coupon.go",
        "coupon_old.go",
        "dao.go",
        "http.go",
        "memcache.go",
        "mysql.go",
        "oldvip.go",
        "order.go",
        "push.go",
        "redis.go",
        "resource.go",
        "sync.go",
        "vip.go",
    ],
    importpath = "go-common/app/job/main/vip/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/job/main/vip/conf:go_default_library",
        "//app/job/main/vip/model:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/stat/prom:go_default_library",
        "//library/time:go_default_library",
        "//library/xstr:go_default_library",
        "//vendor/github.com/google/uuid:go_default_library",
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
