package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "coinexchange_test.go",
        "coinstream_test.go",
        "mc_wallet_test.go",
        "metal_test.go",
        "pub_test.go",
        "redis_test.go",
        "tx_test.go",
        "wallet_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/service/live/wallet/conf:go_default_library",
        "//app/service/live/wallet/model:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "coinexchange.go",
        "coinstream.go",
        "dao.go",
        "mc_wallet.go",
        "metal.go",
        "pub.go",
        "redis.go",
        "tx.go",
        "wallet.go",
    ],
    importpath = "go-common/app/service/live/wallet/dao",
    tags = ["automanaged"],
    deps = [
        "//app/service/live/wallet/conf:go_default_library",
        "//app/service/live/wallet/model:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/queue/databus:go_default_library",
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
