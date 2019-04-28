package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "history.go",
        "http.go",
        "param.go",
        "toview.go",
    ],
    importpath = "go-common/app/interface/main/history/http",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/history/conf:go_default_library",
        "//app/interface/main/history/model:go_default_library",
        "//app/interface/main/history/service:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/log/anticheat:go_default_library",
        "//library/log/infoc:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/middleware/auth:go_default_library",
        "//library/net/http/blademaster/middleware/verify:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/xstr:go_default_library",
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
