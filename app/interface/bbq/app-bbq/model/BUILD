load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "app.go",
        "cache.go",
        "comment.go",
        "danmu.go",
        "invite_code.go",
        "location.go",
        "model.go",
        "redundance.go",
        "report.go",
        "search.go",
        "sv.go",
        "util.go",
        "video.go",
    ],
    importpath = "go-common/app/interface/bbq/app-bbq/model",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/bbq/user/api:go_default_library",
        "//library/time:go_default_library",
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
    srcs = [
        ":package-srcs",
        "//app/interface/bbq/app-bbq/model/grpc:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
