package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "push.go",
        "service.go",
    ],
    importpath = "go-common/app/service/bbq/notice-service/internal/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/bbq/notice-service/api/v1:go_default_library",
        "//app/service/bbq/notice-service/internal/conf:go_default_library",
        "//app/service/bbq/notice-service/internal/dao:go_default_library",
        "//app/service/bbq/notice-service/internal/model:go_default_library",
        "//app/service/bbq/push/api/grpc/v1:go_default_library",
        "//library/log:go_default_library",
        "//library/log/infoc:go_default_library",
        "//library/net/trace:go_default_library",
        "//vendor/github.com/json-iterator/go:go_default_library",
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
