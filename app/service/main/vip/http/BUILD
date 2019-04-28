package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "activity.go",
        "bp.go",
        "coupon.go",
        "http.go",
        "notify.go",
        "order.go",
        "order2.go",
        "panel.go",
        "point.go",
        "resouce.go",
        "sync.go",
        "tips.go",
        "user.go",
        "vip.go",
    ],
    importpath = "go-common/app/service/main/vip/http",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/vip/conf:go_default_library",
        "//app/service/main/vip/model:go_default_library",
        "//app/service/main/vip/service:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/middleware/verify:go_default_library",
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
