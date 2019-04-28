package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "action.go",
        "arctag.go",
        "change.go",
        "channel.go",
        "http.go",
        "platrform.go",
        "sub.go",
        "tag.go",
        "tagarc.go",
    ],
    importpath = "go-common/app/interface/main/tag/http",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/tag/conf:go_default_library",
        "//app/interface/main/tag/model:go_default_library",
        "//app/interface/main/tag/service:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//app/service/main/tag/model:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
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
