package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "service.go",
        "video.go",
    ],
    importpath = "go-common/app/service/bbq/search/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/bbq/search/api/grpc/v1:go_default_library",
        "//app/service/bbq/search/conf:go_default_library",
        "//app/service/bbq/search/dao:go_default_library",
        "//app/service/bbq/search/model:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//vendor/github.com/json-iterator/go:go_default_library",
        "//vendor/gopkg.in/olivere/elastic.v5:go_default_library",
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
