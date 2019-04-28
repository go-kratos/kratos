package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["service.go"],
    importpath = "go-common/app/service/bbq/sys-msg/internal/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/bbq/sys-msg/api/v1:go_default_library",
        "//app/service/bbq/sys-msg/internal/conf:go_default_library",
        "//app/service/bbq/sys-msg/internal/dao:go_default_library",
        "//library/log:go_default_library",
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
