package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["http.go"],
    importpath = "go-common/app/interface/live/lottery-interface/internal/server/http",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/live/lottery-interface/internal/conf:go_default_library",
        "//app/interface/live/lottery-interface/internal/service:go_default_library",
        "//app/interface/live/lottery-interface/internal/service/v1:go_default_library",
        "//app/service/live/live_riskcontrol/api/grpc/v1:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/middleware/auth:go_default_library",
        "//library/net/http/blademaster/middleware/verify:go_default_library",
        "//library/net/http/blademaster/render:go_default_library",
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
