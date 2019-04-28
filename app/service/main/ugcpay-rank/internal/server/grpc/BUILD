package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["elec.go"],
    importpath = "go-common/app/service/main/ugcpay-rank/internal/server/grpc",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/ugcpay-rank/api/v1:go_default_library",
        "//app/service/main/ugcpay-rank/internal/conf:go_default_library",
        "//app/service/main/ugcpay-rank/internal/model:go_default_library",
        "//app/service/main/ugcpay-rank/internal/service:go_default_library",
        "//app/service/main/ugcpay-rank/internal/service/rank:go_default_library",
        "//library/log:go_default_library",
        "//library/net/rpc/warden:go_default_library",
        "@io_bazel_rules_go//proto/wkt:empty_go_proto",
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
