package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "dao.cache.go",
        "dao.go",
        "forbid_user.go",
        "location.go",
        "monkey.go",
        "user.go",
        "user_base.go",
        "user_black.go",
        "user_card.go",
        "user_fan.go",
        "user_follow.go",
        "user_like.go",
        "user_stat.go",
    ],
    importpath = "go-common/app/service/bbq/user/internal/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/bbq/notice-service/api/v1:go_default_library",
        "//app/service/bbq/user/api:go_default_library",
        "//app/service/bbq/user/internal/conf:go_default_library",
        "//app/service/bbq/user/internal/model:go_default_library",
        "//app/service/main/account/api:go_default_library",
        "//app/service/main/filter/api/grpc/v1:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/net/rpc/warden:go_default_library",
        "//library/stat/prom:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//library/sync/pipeline/fanout:go_default_library",
        "//library/time:go_default_library",
        "//library/xstr:go_default_library",
        "//vendor/github.com/bouk/monkey:go_default_library",
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
        "dao.cache_test.go",
        "dao_test.go",
        "user_base_test.go",
        "user_black_test.go",
        "user_card_test.go",
        "user_fan_test.go",
        "user_follow_test.go",
        "user_like_test.go",
        "user_stat_test.go",
    ],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = [
        "//app/service/bbq/user/api:go_default_library",
        "//app/service/bbq/user/internal/conf:go_default_library",
        "//app/service/bbq/user/internal/model:go_default_library",
        "//library/conf/paladin:go_default_library",
        "//library/log:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)
