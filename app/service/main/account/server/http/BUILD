package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "acc.go",
        "http.go",
        "v1.go",
        "v2.go",
    ],
    importpath = "go-common/app/service/main/account/server/http",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/account/api:go_default_library",
        "//app/service/main/account/conf:go_default_library",
        "//app/service/main/account/model:go_default_library",
        "//app/service/main/account/service:go_default_library",
        "//app/service/main/member/model:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/middleware/verify:go_default_library",
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
