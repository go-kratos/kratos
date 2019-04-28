package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "coin_test.go",
        "exp_test.go",
        "service_test.go",
        "usercoin_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/coin/api:go_default_library",
        "//app/service/main/coin/conf:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "archive.go",
        "check.go",
        "coin.go",
        "exp.go",
        "service.go",
        "usercoin.go",
    ],
    importpath = "go-common/app/service/main/coin/service",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/coin/api:go_default_library",
        "//app/service/main/coin/conf:go_default_library",
        "//app/service/main/coin/dao:go_default_library",
        "//app/service/main/coin/model:go_default_library",
        "//app/service/main/member/api:go_default_library",
        "//app/service/main/member/model:go_default_library",
        "//app/service/main/member/model/block:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/sync/errgroup.v2:go_default_library",
        "//library/sync/pipeline/fanout:go_default_library",
        "//library/time:go_default_library",
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
