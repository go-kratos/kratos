package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "dM.go",
        "data.go",
        "validate.go",
    ],
    importpath = "go-common/app/interface/live/app-room/service/v1/dm",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/live/app-room/api/http/v1:go_default_library",
        "//app/interface/live/app-room/conf:go_default_library",
        "//app/interface/live/app-room/model:go_default_library",
        "//app/service/live/live-dm/api/grpc/v1:go_default_library",
        "//app/service/live/live_riskcontrol/api/grpc/v1:go_default_library",
        "//app/service/live/xcaptcha/api/grpc/v1:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/metadata:go_default_library",
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
