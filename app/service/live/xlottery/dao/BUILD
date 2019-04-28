package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "cache_test.go",
        "storm_test.go",
    ],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = [
        "//app/service/live/xlottery/model:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/container/pool:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/time:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "api.go",
        "award.go",
        "cache.go",
        "capsule.go",
        "capsule_extra.go",
        "coin.go",
        "coin_config.go",
        "dao.go",
        "pool.go",
        "pool_prize.go",
        "prize_extra.go",
        "prize_log.go",
        "storm.go",
    ],
    importpath = "go-common/app/service/live/xlottery/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/live/gift/api/liverpc:go_default_library",
        "//app/service/live/rc/api/liverpc:go_default_library",
        "//app/service/live/rc/api/liverpc/v1:go_default_library",
        "//app/service/live/userext/api/liverpc:go_default_library",
        "//app/service/live/xlottery/api/grpc/v1:go_default_library",
        "//app/service/live/xlottery/conf:go_default_library",
        "//app/service/live/xlottery/model:go_default_library",
        "//app/service/main/account/rpc/client:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/log:go_default_library",
        "//library/net/rpc/liverpc:go_default_library",
        "//library/queue/databus/report:go_default_library",
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
