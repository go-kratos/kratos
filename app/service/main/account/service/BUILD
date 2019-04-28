package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "acc.go",
        "cache_delay.go",
        "exp.go",
        "relation.go",
        "service.go",
        "vip.go",
    ],
    importpath = "go-common/app/service/main/account/service",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/account/api:go_default_library",
        "//app/service/main/account/conf:go_default_library",
        "//app/service/main/account/dao:go_default_library",
        "//app/service/main/account/model:go_default_library",
        "//app/service/main/account/model/queue:go_default_library",
        "//app/service/main/coin/api/gorpc:go_default_library",
        "//app/service/main/coin/model:go_default_library",
        "//app/service/main/member/model:go_default_library",
        "//app/service/main/relation/model:go_default_library",
        "//app/service/main/relation/rpc/client:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/stat/prom:go_default_library",
        "//library/sync/errgroup:go_default_library",
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
        "acc_test.go",
        "exp_test.go",
        "relation_test.go",
        "service_test.go",
        "vip_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/account/conf:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)
