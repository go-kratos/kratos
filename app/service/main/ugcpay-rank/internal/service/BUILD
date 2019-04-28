package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "elec.go",
        "rank.go",
        "service.go",
    ],
    importpath = "go-common/app/service/main/ugcpay-rank/internal/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/ugcpay-rank/internal/conf:go_default_library",
        "//app/service/main/ugcpay-rank/internal/dao:go_default_library",
        "//app/service/main/ugcpay-rank/internal/model:go_default_library",
        "//app/service/main/ugcpay-rank/internal/service/rank:go_default_library",
        "//library/log:go_default_library",
        "//library/stat/prom:go_default_library",
        "//library/sync/pipeline/fanout:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
        "//vendor/golang.org/x/sync/singleflight:go_default_library",
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
    srcs = [
        ":package-srcs",
        "//app/service/main/ugcpay-rank/internal/service/rank:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
