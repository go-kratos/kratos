package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "account.go",
        "change_history.go",
        "contract.go",
        "dao.cache.go",
        "dao.go",
        "lock.go",
        "main_vip.go",
        "mc.cache.go",
        "mc.cache.key.go",
        "pay_order.go",
        "price_config.go",
        "user_info.go",
        "yst.go",
    ],
    importpath = "go-common/app/service/main/tv/internal/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/account/api:go_default_library",
        "//app/service/main/tv/internal/conf:go_default_library",
        "//app/service/main/tv/internal/model:go_default_library",
        "//app/service/main/tv/internal/pkg:go_default_library",
        "//app/service/main/vipinfo/api:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/stat/prom:go_default_library",
        "//library/sync/pipeline/fanout:go_default_library",
        "//library/time:go_default_library",
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
        "account_test.go",
        "change_history_test.go",
        "contract_test.go",
        "dao.cache_test.go",
        "dao_test.go",
        "lock_test.go",
        "main_vip_test.go",
        "mc.cache.key_test.go",
        "mc.cache_test.go",
        "pay_order_test.go",
        "price_config_test.go",
        "user_info_test.go",
        "yst_test.go",
    ],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = [
        "//app/service/main/tv/internal/conf:go_default_library",
        "//app/service/main/tv/internal/model:go_default_library",
        "//library/ecode:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)
