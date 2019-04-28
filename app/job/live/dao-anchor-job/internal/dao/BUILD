package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "av.go",
        "bfs.go",
        "dao.go",
        "daoAnchor.go",
        "video.go",
    ],
    importpath = "go-common/app/job/live/dao-anchor-job/internal/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/job/live/dao-anchor-job/internal/conf:go_default_library",
        "//app/service/live/av/api/liverpc:go_default_library",
        "//app/service/live/av/api/liverpc/v1:go_default_library",
        "//app/service/live/dao-anchor/api/grpc/v0:go_default_library",
        "//app/service/live/dao-anchor/api/grpc/v1:go_default_library",
        "//app/service/video/stream-mng/api/v1:go_default_library",
        "//library/database/bfs:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/rpc/liverpc:go_default_library",
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
