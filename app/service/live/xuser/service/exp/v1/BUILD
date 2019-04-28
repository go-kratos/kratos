package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "audit.go",
        "check_param.go",
        "clear_cache.go",
        "level_rule.go",
        "log_debug.go",
        "update_db.go",
        "user_exp.go",
    ],
    importpath = "go-common/app/service/live/xuser/service/exp/v1",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/live/dao-anchor/api/grpc/v1:go_default_library",
        "//app/service/live/xuser/api/grpc/v1:go_default_library",
        "//app/service/live/xuser/conf:go_default_library",
        "//app/service/live/xuser/dao/exp:go_default_library",
        "//app/service/live/xuser/dao/xanchor:go_default_library",
        "//app/service/live/xuser/model/exp:go_default_library",
        "//library/cache:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/queue/databus/report:go_default_library",
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
