package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "business.go",
        "dialog.go",
        "http.go",
        "jointly.go",
        "month.go",
        "order.go",
        "panel.go",
        "platform.go",
        "privilege.go",
        "push.go",
        "resouce.go",
        "tips.go",
        "version.go",
        "vip.go",
        "welfare.go",
    ],
    importpath = "go-common/app/admin/main/vip/http",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/vip/conf:go_default_library",
        "//app/admin/main/vip/model:go_default_library",
        "//app/admin/main/vip/service:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/binding:go_default_library",
        "//library/net/http/blademaster/middleware/permit:go_default_library",
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
