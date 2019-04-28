package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["service.go"],
    importpath = "go-common/app/interface/main/app-intl/service/player",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/app-intl/conf:go_default_library",
        "//app/interface/main/app-intl/dao/account:go_default_library",
        "//app/interface/main/app-intl/dao/archive:go_default_library",
        "//app/interface/main/app-intl/dao/player:go_default_library",
        "//app/interface/main/app-intl/dao/resource:go_default_library",
        "//app/interface/main/app-intl/dao/ugcpay:go_default_library",
        "//app/interface/main/app-intl/model:go_default_library",
        "//app/interface/main/app-intl/model/player:go_default_library",
        "//app/interface/main/app-intl/model/player/archive:go_default_library",
        "//app/service/main/account/model:go_default_library",
        "//app/service/main/ugcpay/api/grpc/v1:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
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
