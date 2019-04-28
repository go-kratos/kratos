package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "dao.cache_test.go",
        "dao_test.go",
        "mc.cache_test.go",
        "up_group_test.go",
        "up_mc_test.go",
        "up_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/up/api/v1:go_default_library",
        "//app/service/main/up/conf:go_default_library",
        "//app/service/main/up/dao:go_default_library",
        "//app/service/main/up/dao/global:go_default_library",
        "//app/service/main/up/model:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "dao.cache.go",
        "dao.go",
        "mc.cache.go",
        "up.go",
        "up_group.go",
        "up_mc.go",
    ],
    importpath = "go-common/app/service/main/up/dao/manager",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/up/api/v1:go_default_library",
        "//app/service/main/up/conf:go_default_library",
        "//app/service/main/up/dao:go_default_library",
        "//app/service/main/up/dao/global:go_default_library",
        "//app/service/main/up/model:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/stat/prom:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//library/sync/pipeline/fanout:go_default_library",
        "//library/xstr:go_default_library",
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
