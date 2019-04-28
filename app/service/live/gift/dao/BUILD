package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "dao_test.go",
        "discount_gift_test.go",
        "discount_plan_test.go",
        "gift_online_test.go",
        "gift_plan_test.go",
        "redis_test.go",
        "ugift_day_test.go",
        "ugift_week_test.go",
        "user_gift_test.go",
    ],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = [
        "//app/service/live/gift/api/grpc/v1:go_default_library",
        "//app/service/live/gift/conf:go_default_library",
        "//app/service/live/gift/model:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "dao.go",
        "discount_gift.go",
        "discount_plan.go",
        "gift_online.go",
        "gift_plan.go",
        "redis.go",
        "ugift_day.go",
        "ugift_week.go",
        "user_gift.go",
    ],
    importpath = "go-common/app/service/live/gift/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/live/activity/api/liverpc:go_default_library",
        "//app/service/live/fans_medal/api/liverpc:go_default_library",
        "//app/service/live/gift/api/grpc/v1:go_default_library",
        "//app/service/live/gift/conf:go_default_library",
        "//app/service/live/gift/model:go_default_library",
        "//app/service/live/live_user/api/liverpc:go_default_library",
        "//app/service/live/room/api/liverpc:go_default_library",
        "//app/service/live/user/api/liverpc:go_default_library",
        "//app/service/live/xuser/api/grpc/v1:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/log:go_default_library",
        "//library/net/rpc/liverpc:go_default_library",
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
