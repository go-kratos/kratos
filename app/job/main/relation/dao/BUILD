package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "cache_update_test.go",
        "dao_test.go",
        "follower_test.go",
        "following_test.go",
        "mysql_test.go",
        "stat_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/job/main/relation/conf:go_default_library",
        "//app/job/main/relation/model:go_default_library",
        "//app/service/main/relation/model:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "cache_update.go",
        "dao.go",
        "follower.go",
        "following.go",
        "mysql.go",
        "stat.go",
    ],
    importpath = "go-common/app/job/main/relation/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/job/main/relation/conf:go_default_library",
        "//app/job/main/relation/model:go_default_library",
        "//app/job/main/relation/model/i64b:go_default_library",
        "//app/service/main/relation/model:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/time:go_default_library",
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
