package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_binary",
    "go_library",
)

go_binary(
    name = "cmd",
    embed = [":go_default_library"],
    tags = ["automanaged"],
)

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    data = [
        "app-resource-example.toml",
        "app-resource-test.toml",
    ],
    importpath = "go-common/app/interface/main/app-resource/cmd",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/app-resource/conf:go_default_library",
        "//app/interface/main/app-resource/grpc:go_default_library",
        "//app/interface/main/app-resource/http:go_default_library",
        "//app/interface/main/app-resource/service/abtest:go_default_library",
        "//app/interface/main/app-resource/service/audit:go_default_library",
        "//app/interface/main/app-resource/service/broadcast:go_default_library",
        "//app/interface/main/app-resource/service/domain:go_default_library",
        "//app/interface/main/app-resource/service/guide:go_default_library",
        "//app/interface/main/app-resource/service/module:go_default_library",
        "//app/interface/main/app-resource/service/notice:go_default_library",
        "//app/interface/main/app-resource/service/param:go_default_library",
        "//app/interface/main/app-resource/service/ping:go_default_library",
        "//app/interface/main/app-resource/service/plugin:go_default_library",
        "//app/interface/main/app-resource/service/show:go_default_library",
        "//app/interface/main/app-resource/service/sidebar:go_default_library",
        "//app/interface/main/app-resource/service/splash:go_default_library",
        "//app/interface/main/app-resource/service/static:go_default_library",
        "//app/interface/main/app-resource/service/version:go_default_library",
        "//app/interface/main/app-resource/service/white:go_default_library",
        "//library/ecode/tip:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster/middleware/auth:go_default_library",
        "//library/net/trace:go_default_library",
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
