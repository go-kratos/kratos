package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "appeal.go",
        "blocked.go",
        "credit.go",
        "http.go",
        "jury.go",
        "labour.go",
    ],
    importpath = "go-common/app/interface/main/credit/http",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/credit/conf:go_default_library",
        "//app/interface/main/credit/model:go_default_library",
        "//app/interface/main/credit/service:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/middleware/antispam:go_default_library",
        "//library/net/http/blademaster/middleware/auth:go_default_library",
        "//library/net/http/blademaster/middleware/verify:go_default_library",
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
