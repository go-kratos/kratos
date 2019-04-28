package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "abtest.go",
        "audit.go",
        "broadcast.go",
        "domain.go",
        "http.go",
        "interest.go",
        "local.go",
        "module.go",
        "notice.go",
        "param.go",
        "plugin.go",
        "show.go",
        "sidebar.go",
        "splash.go",
        "static.go",
        "version.go",
        "white.go",
    ],
    importpath = "go-common/app/interface/main/app-resource/http",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/app-resource/conf:go_default_library",
        "//app/interface/main/app-resource/model:go_default_library",
        "//app/interface/main/app-resource/model/module:go_default_library",
        "//app/interface/main/app-resource/model/version:go_default_library",
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
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/middleware/auth:go_default_library",
        "@com_github_golang_protobuf//proto:go_default_library",
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
