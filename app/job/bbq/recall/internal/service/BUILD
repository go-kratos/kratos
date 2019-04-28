package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "bloomfilter.go",
        "forward_index.go",
        "inverted_index.go",
        "service.go",
    ],
    importpath = "go-common/app/job/bbq/recall/internal/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/job/bbq/recall/internal/conf:go_default_library",
        "//app/job/bbq/recall/internal/dao:go_default_library",
        "//app/job/bbq/recall/proto:go_default_library",
        "//library/log:go_default_library",
        "//vendor/github.com/golang/snappy:go_default_library",
        "//vendor/github.com/robfig/cron:go_default_library",
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
