package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "archive.go",
        "play.go",
        "service.go",
        "upload.go",
        "video.go",
    ],
    importpath = "go-common/app/service/bbq/video/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/bbq/common/db/bbq:go_default_library",
        "//app/service/bbq/topic/api:go_default_library",
        "//app/service/bbq/video/api/grpc/v1:go_default_library",
        "//app/service/bbq/video/api/http/v1:go_default_library",
        "//app/service/bbq/video/conf:go_default_library",
        "//app/service/bbq/video/dao:go_default_library",
        "//app/service/bbq/video/model:go_default_library",
        "//app/service/bbq/video/model/grpc:go_default_library",
        "//library/conf/env:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
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

go_test(
    name = "go_default_test",
    srcs = ["service_test.go"],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = [
        "//app/service/bbq/video/api/grpc/v1:go_default_library",
        "//app/service/bbq/video/conf:go_default_library",
    ],
)
