package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "av.go",
        "av_storage.go",
        "rank.go",
        "tool.go",
        "up.go",
        "up_storage.go",
    ],
    importpath = "go-common/app/service/main/ugcpay-rank/internal/service/rank",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/account/model:go_default_library",
        "//app/service/main/ugcpay-rank/internal/dao:go_default_library",
        "//app/service/main/ugcpay-rank/internal/model:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/log:go_default_library",
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
