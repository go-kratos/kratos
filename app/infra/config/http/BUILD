package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "config.go",
        "config2.go",
        "http.go",
        "local.go",
    ],
    importpath = "go-common/app/infra/config/http",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/infra/config/conf:go_default_library",
        "//app/infra/config/model:go_default_library",
        "//app/infra/config/service/v1:go_default_library",
        "//app/infra/config/service/v2:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/middleware/antispam:go_default_library",
        "//library/net/http/blademaster/middleware/verify:go_default_library",
        "//vendor/github.com/dgryski/go-farm:go_default_library",
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
