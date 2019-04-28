package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "forward_index.go",
        "inverted_index.go",
        "loader.go",
    ],
    importpath = "go-common/app/service/bbq/recsys-recall/service/index",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/job/bbq/recall/proto:go_default_library",
        "//app/service/bbq/recsys-recall/conf:go_default_library",
        "//library/log:go_default_library",
        "//vendor/github.com/golang/snappy:go_default_library",
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
    name = "go_default_xtest",
    srcs = ["inverted_index_test.go"],
    tags = ["automanaged"],
    deps = ["//app/service/bbq/recsys-recall/service/index:go_default_library"],
)
